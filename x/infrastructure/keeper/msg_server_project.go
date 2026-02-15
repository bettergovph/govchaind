package keeper

import (
	"context"
	"fmt"

	"govchain/x/infrastructure/types"

	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateProject(ctx context.Context, msg *types.MsgCreateProject) (*types.MsgCreateProjectResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	if msg.ProjectId == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidProjectID, "project ID cannot be empty")
	}

	// Check if project already exists
	_, err := k.Project.Get(ctx, msg.ProjectId)
	if err == nil {
		return nil, errorsmod.Wrap(types.ErrProjectAlreadyExists, msg.ProjectId)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	txBytes := sdkCtx.TxBytes()
	txHash := fmt.Sprintf("%X", tmhash.Sum(txBytes))

	project := types.InfraProject{
		Id:                       msg.ProjectId,
		Title:                    msg.Title,
		Description:              msg.Description,
		Status:                   msg.Status,
		Type:                     msg.ProjectType,
		Sector:                   msg.Sector,
		Purpose:                  msg.Purpose,
		Budget:                   msg.Budget,
		Period:                   msg.Period,
		Locations:                msg.Locations,
		Parties:                  msg.Parties,
		PublicAuthority:          msg.PublicAuthority,
		AdditionalClassifications: msg.AdditionalClassifications,
		Language:                 "en",
		PublisherAddress:         msg.Creator,
		BlockHeight:             fmt.Sprintf("%d", sdkCtx.BlockHeight()),
		TxHash:                  txHash,
		Owner:                   msg.Creator,
		AuthorizedPublishers:    []string{msg.Creator},
		CreatedAtBlock:          uint64(sdkCtx.BlockHeight()),
		UpdatedAtBlock:          uint64(sdkCtx.BlockHeight()),
	}

	if err := k.Project.Set(ctx, msg.ProjectId, project); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to store project")
	}

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"project_created",
		sdk.NewAttribute("project_id", msg.ProjectId),
		sdk.NewAttribute("title", msg.Title),
		sdk.NewAttribute("status", msg.Status),
		sdk.NewAttribute("publisher", msg.Creator),
	))

	return &types.MsgCreateProjectResponse{ProjectId: msg.ProjectId}, nil
}

func (k msgServer) UpdateProject(ctx context.Context, msg *types.MsgUpdateProject) (*types.MsgUpdateProjectResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	project, err := k.Project.Get(ctx, msg.ProjectId)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrProjectNotFound, msg.ProjectId)
	}

	if project.Owner != msg.Creator && !contains(project.AuthorizedPublishers, msg.Creator) {
		return nil, errorsmod.Wrap(types.ErrUnauthorizedPublisher, msg.Creator)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	txHash := fmt.Sprintf("%X", tmhash.Sum(sdkCtx.TxBytes()))

	// Update fields if provided
	if msg.Title != "" {
		project.Title = msg.Title
	}
	if msg.Description != "" {
		project.Description = msg.Description
	}
	if msg.Status != "" {
		project.Status = msg.Status
	}
	if msg.ProjectType != "" {
		project.Type = msg.ProjectType
	}
	if len(msg.Sector) > 0 {
		project.Sector = msg.Sector
	}
	if msg.Purpose != "" {
		project.Purpose = msg.Purpose
	}
	if msg.Budget != nil {
		project.Budget = msg.Budget
	}
	if msg.Period != nil {
		project.Period = msg.Period
	}
	if len(msg.Locations) > 0 {
		project.Locations = msg.Locations
	}
	if len(msg.Parties) > 0 {
		project.Parties = msg.Parties
	}
	if msg.Completion != nil {
		project.Completion = msg.Completion
	}

	project.PublisherAddress = msg.Creator
	project.BlockHeight = fmt.Sprintf("%d", sdkCtx.BlockHeight())
	project.TxHash = txHash
	project.UpdatedAtBlock = uint64(sdkCtx.BlockHeight())

	if err := k.Project.Set(ctx, msg.ProjectId, project); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update project")
	}

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"project_updated",
		sdk.NewAttribute("project_id", msg.ProjectId),
		sdk.NewAttribute("status", project.Status),
		sdk.NewAttribute("publisher", msg.Creator),
	))

	return &types.MsgUpdateProjectResponse{}, nil
}

func (k msgServer) UpdateProgress(ctx context.Context, msg *types.MsgUpdateProgress) (*types.MsgUpdateProgressResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	project, err := k.Project.Get(ctx, msg.ProjectId)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrProjectNotFound, msg.ProjectId)
	}

	if project.Owner != msg.Creator && !contains(project.AuthorizedPublishers, msg.Creator) {
		return nil, errorsmod.Wrap(types.ErrUnauthorizedPublisher, msg.Creator)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Merge metrics: update existing or append new
	for _, newMetric := range msg.Metrics {
		found := false
		for i, existing := range project.Metrics {
			if existing.Id == newMetric.Id {
				// Append new observations to existing metric
				project.Metrics[i].Observations = append(project.Metrics[i].Observations, newMetric.Observations...)
				found = true
				break
			}
		}
		if !found {
			project.Metrics = append(project.Metrics, newMetric)
		}
	}

	project.UpdatedAtBlock = uint64(sdkCtx.BlockHeight())
	project.PublisherAddress = msg.Creator
	project.BlockHeight = fmt.Sprintf("%d", sdkCtx.BlockHeight())

	if err := k.Project.Set(ctx, msg.ProjectId, project); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update project progress")
	}

	for _, m := range msg.Metrics {
		measure := ""
		if len(m.Observations) > 0 {
			measure = m.Observations[0].Measure
		}
		sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
			"progress_updated",
			sdk.NewAttribute("project_id", msg.ProjectId),
			sdk.NewAttribute("metric_id", m.Id),
			sdk.NewAttribute("measure", measure),
			sdk.NewAttribute("publisher", msg.Creator),
		))
	}

	return &types.MsgUpdateProgressResponse{}, nil
}

func (k msgServer) AttachDocument(ctx context.Context, msg *types.MsgAttachDocument) (*types.MsgAttachDocumentResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	project, err := k.Project.Get(ctx, msg.ProjectId)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrProjectNotFound, msg.ProjectId)
	}

	if project.Owner != msg.Creator && !contains(project.AuthorizedPublishers, msg.Creator) {
		return nil, errorsmod.Wrap(types.ErrUnauthorizedPublisher, msg.Creator)
	}

	if msg.Document == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "document cannot be nil")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	project.Documents = append(project.Documents, msg.Document)
	project.UpdatedAtBlock = uint64(sdkCtx.BlockHeight())
	project.PublisherAddress = msg.Creator
	project.BlockHeight = fmt.Sprintf("%d", sdkCtx.BlockHeight())

	if err := k.Project.Set(ctx, msg.ProjectId, project); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to attach document")
	}

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"document_attached",
		sdk.NewAttribute("project_id", msg.ProjectId),
		sdk.NewAttribute("document_id", msg.Document.Id),
		sdk.NewAttribute("document_type", msg.Document.DocumentType),
		sdk.NewAttribute("publisher", msg.Creator),
	))

	return &types.MsgAttachDocumentResponse{DocumentId: msg.Document.Id}, nil
}

func (k msgServer) LinkContractingProcess(ctx context.Context, msg *types.MsgLinkContractingProcess) (*types.MsgLinkContractingProcessResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	project, err := k.Project.Get(ctx, msg.ProjectId)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrProjectNotFound, msg.ProjectId)
	}

	if project.Owner != msg.Creator && !contains(project.AuthorizedPublishers, msg.Creator) {
		return nil, errorsmod.Wrap(types.ErrUnauthorizedPublisher, msg.Creator)
	}

	if msg.ContractingProcess == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "contracting process cannot be nil")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	project.ContractingProcesses = append(project.ContractingProcesses, msg.ContractingProcess)
	project.UpdatedAtBlock = uint64(sdkCtx.BlockHeight())
	project.PublisherAddress = msg.Creator
	project.BlockHeight = fmt.Sprintf("%d", sdkCtx.BlockHeight())

	if err := k.Project.Set(ctx, msg.ProjectId, project); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to link contracting process")
	}

	return &types.MsgLinkContractingProcessResponse{}, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

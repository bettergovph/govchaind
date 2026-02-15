package keeper

import (
	"context"
	"fmt"

	"govchain/x/procurement/types"

	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) SubmitRelease(ctx context.Context, msg *types.MsgSubmitRelease) (*types.MsgSubmitReleaseResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	if msg.Release == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "release cannot be nil")
	}

	if msg.Release.Ocid == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidOCID, "OCID cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	txBytes := sdkCtx.TxBytes()
	txHash := fmt.Sprintf("%X", tmhash.Sum(txBytes))
	blockHeight := fmt.Sprintf("%d", sdkCtx.BlockHeight())

	// Set on-chain metadata
	release := *msg.Release
	release.PublisherAddress = msg.Creator
	release.BlockHeight = blockHeight
	release.TxHash = txHash
	release.ChainTimestamp = nil // will be set by block time

	// Generate a composite key: ocid/release_id
	releaseKey := release.Ocid + "/" + release.Id
	if err := k.Release.Set(ctx, releaseKey, release); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to store release")
	}

	// Update or create the contracting process
	process, err := k.Process.Get(ctx, release.Ocid)
	if err != nil {
		// Create new process
		process = types.ContractingProcess{
			Ocid:                 release.Ocid,
			CompiledRelease:      &release,
			ReleaseIds:           []string{release.Id},
			CurrentStage:         "planning",
			Owner:                msg.Creator,
			AuthorizedPublishers: []string{msg.Creator},
			CreatedAtBlock:       uint64(sdkCtx.BlockHeight()),
			UpdatedAtBlock:       uint64(sdkCtx.BlockHeight()),
		}
	} else {
		// Update existing process
		process.CompiledRelease = &release
		process.ReleaseIds = append(process.ReleaseIds, release.Id)
		process.UpdatedAtBlock = uint64(sdkCtx.BlockHeight())

		// Update stage based on release tags
		if len(release.Tag) > 0 {
			process.CurrentStage = deriveStageFromTag(release.Tag[0])
		}
	}

	if err := k.Process.Set(ctx, release.Ocid, process); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update contracting process")
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"release_submitted",
		sdk.NewAttribute("ocid", release.Ocid),
		sdk.NewAttribute("release_id", release.Id),
		sdk.NewAttribute("publisher", msg.Creator),
	))

	return &types.MsgSubmitReleaseResponse{
		Ocid:        release.Ocid,
		ReleaseId:   release.Id,
		BlockHeight: uint64(sdkCtx.BlockHeight()),
		TxHash:      txHash,
	}, nil
}

func (k msgServer) CreateProcess(ctx context.Context, msg *types.MsgCreateProcess) (*types.MsgCreateProcessResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	if msg.Ocid == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidOCID, "OCID cannot be empty")
	}

	// Check if process already exists
	_, err := k.Process.Get(ctx, msg.Ocid)
	if err == nil {
		return nil, errorsmod.Wrap(types.ErrProcessAlreadyExists, msg.Ocid)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Build the initial release
	release := types.Release{
		Ocid:             msg.Ocid,
		Id:               msg.Ocid + "-planning-001",
		Tag:              []string{"planning"},
		InitiationType:   "tender",
		Buyer:            msg.Buyer,
		Planning:         msg.Planning,
		Parties:          msg.Parties,
		PublisherAddress:  msg.Creator,
		BlockHeight:      fmt.Sprintf("%d", sdkCtx.BlockHeight()),
	}

	// Store release
	releaseKey := release.Ocid + "/" + release.Id
	if err := k.Release.Set(ctx, releaseKey, release); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to store release")
	}

	// Create process
	process := types.ContractingProcess{
		Ocid:                 msg.Ocid,
		CompiledRelease:      &release,
		ReleaseIds:           []string{release.Id},
		CurrentStage:         "planning",
		CurrentStatus:        "active",
		Owner:                msg.Creator,
		AuthorizedPublishers: []string{msg.Creator},
		CreatedAtBlock:       uint64(sdkCtx.BlockHeight()),
		UpdatedAtBlock:       uint64(sdkCtx.BlockHeight()),
	}

	if err := k.Process.Set(ctx, msg.Ocid, process); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to store contracting process")
	}

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"process_created",
		sdk.NewAttribute("ocid", msg.Ocid),
		sdk.NewAttribute("publisher", msg.Creator),
		sdk.NewAttribute("stage", "planning"),
	))

	return &types.MsgCreateProcessResponse{Ocid: msg.Ocid}, nil
}

func (k msgServer) PublishTender(ctx context.Context, msg *types.MsgPublishTender) (*types.MsgPublishTenderResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	process, err := k.Process.Get(ctx, msg.Ocid)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrProcessNotFound, msg.Ocid)
	}

	if process.Owner != msg.Creator && !contains(process.AuthorizedPublishers, msg.Creator) {
		return nil, errorsmod.Wrap(types.ErrUnauthorizedPublisher, msg.Creator)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	txHash := fmt.Sprintf("%X", tmhash.Sum(sdkCtx.TxBytes()))

	tenderId := msg.Tender.Id
	if tenderId == "" {
		tenderId = msg.Ocid + "-tender-001"
	}

	release := types.Release{
		Ocid:             msg.Ocid,
		Id:               msg.Ocid + "-tender-" + tenderId,
		Tag:              []string{"tender"},
		InitiationType:   "tender",
		Tender:           msg.Tender,
		Parties:          msg.AdditionalParties,
		PublisherAddress:  msg.Creator,
		BlockHeight:      fmt.Sprintf("%d", sdkCtx.BlockHeight()),
		TxHash:           txHash,
	}

	releaseKey := release.Ocid + "/" + release.Id
	if err := k.Release.Set(ctx, releaseKey, release); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to store release")
	}

	// Update process
	process.CompiledRelease.Tender = msg.Tender
	if msg.AdditionalParties != nil {
		process.CompiledRelease.Parties = append(process.CompiledRelease.Parties, msg.AdditionalParties...)
	}
	process.ReleaseIds = append(process.ReleaseIds, release.Id)
	process.CurrentStage = "tender"
	process.UpdatedAtBlock = uint64(sdkCtx.BlockHeight())

	if err := k.Process.Set(ctx, msg.Ocid, process); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update contracting process")
	}

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"tender_published",
		sdk.NewAttribute("ocid", msg.Ocid),
		sdk.NewAttribute("tender_id", tenderId),
		sdk.NewAttribute("publisher", msg.Creator),
	))

	return &types.MsgPublishTenderResponse{Ocid: msg.Ocid, TenderId: tenderId}, nil
}

func (k msgServer) SubmitAward(ctx context.Context, msg *types.MsgSubmitAward) (*types.MsgSubmitAwardResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	process, err := k.Process.Get(ctx, msg.Ocid)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrProcessNotFound, msg.Ocid)
	}

	if process.Owner != msg.Creator && !contains(process.AuthorizedPublishers, msg.Creator) {
		return nil, errorsmod.Wrap(types.ErrUnauthorizedPublisher, msg.Creator)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	txHash := fmt.Sprintf("%X", tmhash.Sum(sdkCtx.TxBytes()))

	awardId := msg.Award.Id
	release := types.Release{
		Ocid:             msg.Ocid,
		Id:               msg.Ocid + "-award-" + awardId,
		Tag:              []string{"award"},
		InitiationType:   "tender",
		Awards:           []*types.Award{msg.Award},
		Parties:          msg.AdditionalParties,
		PublisherAddress:  msg.Creator,
		BlockHeight:      fmt.Sprintf("%d", sdkCtx.BlockHeight()),
		TxHash:           txHash,
	}

	releaseKey := release.Ocid + "/" + release.Id
	if err := k.Release.Set(ctx, releaseKey, release); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to store release")
	}

	// Update process
	process.CompiledRelease.Awards = append(process.CompiledRelease.Awards, msg.Award)
	if msg.AdditionalParties != nil {
		process.CompiledRelease.Parties = append(process.CompiledRelease.Parties, msg.AdditionalParties...)
	}
	process.ReleaseIds = append(process.ReleaseIds, release.Id)
	process.CurrentStage = "award"
	process.UpdatedAtBlock = uint64(sdkCtx.BlockHeight())

	if err := k.Process.Set(ctx, msg.Ocid, process); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update contracting process")
	}

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"award_made",
		sdk.NewAttribute("ocid", msg.Ocid),
		sdk.NewAttribute("award_id", awardId),
		sdk.NewAttribute("publisher", msg.Creator),
	))

	return &types.MsgSubmitAwardResponse{Ocid: msg.Ocid, AwardId: awardId}, nil
}

func (k msgServer) SignContract(ctx context.Context, msg *types.MsgSignContract) (*types.MsgSignContractResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	process, err := k.Process.Get(ctx, msg.Ocid)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrProcessNotFound, msg.Ocid)
	}

	if process.Owner != msg.Creator && !contains(process.AuthorizedPublishers, msg.Creator) {
		return nil, errorsmod.Wrap(types.ErrUnauthorizedPublisher, msg.Creator)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	txHash := fmt.Sprintf("%X", tmhash.Sum(sdkCtx.TxBytes()))

	contractId := msg.Contract.Id
	release := types.Release{
		Ocid:             msg.Ocid,
		Id:               msg.Ocid + "-contract-" + contractId,
		Tag:              []string{"contract"},
		InitiationType:   "tender",
		Contracts:        []*types.Contract{msg.Contract},
		PublisherAddress:  msg.Creator,
		BlockHeight:      fmt.Sprintf("%d", sdkCtx.BlockHeight()),
		TxHash:           txHash,
	}

	releaseKey := release.Ocid + "/" + release.Id
	if err := k.Release.Set(ctx, releaseKey, release); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to store release")
	}

	process.CompiledRelease.Contracts = append(process.CompiledRelease.Contracts, msg.Contract)
	process.ReleaseIds = append(process.ReleaseIds, release.Id)
	process.CurrentStage = "contract"
	process.UpdatedAtBlock = uint64(sdkCtx.BlockHeight())

	if err := k.Process.Set(ctx, msg.Ocid, process); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update contracting process")
	}

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"contract_signed",
		sdk.NewAttribute("ocid", msg.Ocid),
		sdk.NewAttribute("contract_id", contractId),
		sdk.NewAttribute("publisher", msg.Creator),
	))

	return &types.MsgSignContractResponse{Ocid: msg.Ocid, ContractId: contractId}, nil
}

func (k msgServer) UpdateImplementation(ctx context.Context, msg *types.MsgUpdateImplementation) (*types.MsgUpdateImplementationResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	process, err := k.Process.Get(ctx, msg.Ocid)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrProcessNotFound, msg.Ocid)
	}

	if process.Owner != msg.Creator && !contains(process.AuthorizedPublishers, msg.Creator) {
		return nil, errorsmod.Wrap(types.ErrUnauthorizedPublisher, msg.Creator)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	txHash := fmt.Sprintf("%X", tmhash.Sum(sdkCtx.TxBytes()))

	release := types.Release{
		Ocid:             msg.Ocid,
		Id:               fmt.Sprintf("%s-impl-%d", msg.Ocid, sdkCtx.BlockHeight()),
		Tag:              []string{"implementation"},
		InitiationType:   "tender",
		PublisherAddress:  msg.Creator,
		BlockHeight:      fmt.Sprintf("%d", sdkCtx.BlockHeight()),
		TxHash:           txHash,
	}

	releaseKey := release.Ocid + "/" + release.Id
	if err := k.Release.Set(ctx, releaseKey, release); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to store release")
	}

	// Update implementation on the matching contract
	for i, c := range process.CompiledRelease.Contracts {
		if c.Id == msg.ContractId {
			process.CompiledRelease.Contracts[i].Implementation = msg.Implementation
			break
		}
	}

	process.ReleaseIds = append(process.ReleaseIds, release.Id)
	process.CurrentStage = "implementation"
	process.UpdatedAtBlock = uint64(sdkCtx.BlockHeight())

	if err := k.Process.Set(ctx, msg.Ocid, process); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update contracting process")
	}

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"implementation_updated",
		sdk.NewAttribute("ocid", msg.Ocid),
		sdk.NewAttribute("contract_id", msg.ContractId),
		sdk.NewAttribute("publisher", msg.Creator),
	))

	return &types.MsgUpdateImplementationResponse{}, nil
}

// deriveStageFromTag maps a release tag string to a lifecycle stage
func deriveStageFromTag(tag string) string {
	switch tag {
	case "planning", "planningUpdate":
		return "planning"
	case "tender", "tenderAmendment", "tenderUpdate", "tenderCancellation":
		return "tender"
	case "award", "awardUpdate", "awardCancellation":
		return "award"
	case "contract", "contractUpdate", "contractAmendment":
		return "contract"
	case "implementation", "implementationUpdate":
		return "implementation"
	default:
		return "planning"
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

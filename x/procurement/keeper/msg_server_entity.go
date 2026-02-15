package keeper

import (
	"context"
	"fmt"

	"govchain/x/procurement/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RegisterEntity(ctx context.Context, msg *types.MsgRegisterEntity) (*types.MsgRegisterEntityResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	if msg.Organization == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "organization cannot be nil")
	}

	entityId := msg.Organization.Id
	if entityId == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "organization id cannot be empty")
	}

	// Check if entity already exists
	_, err := k.Entity.Get(ctx, entityId)
	if err == nil {
		return nil, errorsmod.Wrap(types.ErrEntityAlreadyExists, entityId)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Set on-chain address if not already set
	if msg.Organization.OnChainAddress == "" {
		msg.Organization.OnChainAddress = msg.Creator
	}

	entity := types.EntityRegistry{
		EntityId:          entityId,
		Organization:      msg.Organization,
		Registrar:         msg.Creator,
		RegisteredAtBlock: uint64(sdkCtx.BlockHeight()),
		Active:            true,
	}

	if err := k.Entity.Set(ctx, entityId, entity); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to store entity")
	}

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"entity_registered",
		sdk.NewAttribute("entity_id", entityId),
		sdk.NewAttribute("name", msg.Organization.Name),
		sdk.NewAttribute("registrar", msg.Creator),
	))

	return &types.MsgRegisterEntityResponse{EntityId: entityId}, nil
}

func (k msgServer) UpdateEntity(ctx context.Context, msg *types.MsgUpdateEntity) (*types.MsgUpdateEntityResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	entity, err := k.Entity.Get(ctx, msg.EntityId)
	if err != nil {
		if err == collections.ErrNotFound {
			return nil, errorsmod.Wrap(types.ErrEntityNotFound, msg.EntityId)
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get entity")
	}

	// Only the original registrar can update
	if entity.Registrar != msg.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the registrar can update the entity")
	}

	entity.Organization = msg.Organization

	if err := k.Entity.Set(ctx, msg.EntityId, entity); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update entity")
	}

	return &types.MsgUpdateEntityResponse{}, nil
}

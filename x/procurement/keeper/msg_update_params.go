package keeper

import (
	"context"

	"govchain/x/procurement/types"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	authorityAddr, err := k.addressCodec.BytesToString(k.GetAuthority())
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid authority address")
	}
	if msg.Authority != authorityAddr {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "invalid authority; expected %s, got %s", authorityAddr, msg.Authority)
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	if err := k.Params.Set(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

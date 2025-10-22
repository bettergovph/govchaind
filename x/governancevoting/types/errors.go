package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrDelegationNotFound    = sdkerrors.Register(ModuleName, 1, "delegation not found")
	ErrDelegationInvalidAddr = sdkerrors.Register(ModuleName, 2, "invalid delegation address")
)

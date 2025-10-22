package types

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrDuplicateContractID = sdkerrors.Register(ModuleName, 1, "duplicate contract identifier in genesis state")
)

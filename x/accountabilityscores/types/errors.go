package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrScorecardNotFound    = sdkerrors.Register(ModuleName, 1, "scorecard not found")
	ErrScorecardUpdaterAddr = sdkerrors.Register(ModuleName, 2, "invalid scorecard updated_by address")
	ErrInvalidApproval      = sdkerrors.Register(ModuleName, 3, "invalid accountability approval")
)

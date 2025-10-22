package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrBudgetNotFound           = sdkerrors.Register(ModuleName, 1, "budget not found")
	ErrInvalidBudgetTransition  = sdkerrors.Register(ModuleName, 2, "invalid budget status transition")
	ErrBudgetCreatorInvalidAddr = sdkerrors.Register(ModuleName, 3, "invalid budget creator address")
	ErrBudgetLiabilityMissing   = sdkerrors.Register(ModuleName, 4, "liability contract required")
	ErrBudgetApprovalSignerAddr = sdkerrors.Register(ModuleName, 5, "invalid budget approval signer")
	ErrBudgetApprovalDuplicate  = sdkerrors.Register(ModuleName, 6, "duplicate budget approval role or signer")
	ErrBudgetApprovalThreshold  = sdkerrors.Register(ModuleName, 7, "insufficient budget approvals for liability requirements")
	ErrBudgetApprovalScore      = sdkerrors.Register(ModuleName, 8, "budget approval accountability score below threshold")
)

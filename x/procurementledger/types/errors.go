package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrProcurementNotFound          = sdkerrors.Register(ModuleName, 1, "procurement not found")
	ErrInvalidProcurementTransition = sdkerrors.Register(ModuleName, 2, "invalid procurement status transition")
	ErrProcurementOfficerAddr       = sdkerrors.Register(ModuleName, 3, "invalid procurement officer address")
	ErrBudgetMissingForProcurement  = sdkerrors.Register(ModuleName, 4, "budget not found for procurement")
	ErrProcurementLiabilityMissing  = sdkerrors.Register(ModuleName, 5, "liability contract required for procurement")
	ErrProcurementApprovalSigner    = sdkerrors.Register(ModuleName, 6, "invalid procurement approval signer")
	ErrProcurementApprovalDuplicate = sdkerrors.Register(ModuleName, 7, "duplicate procurement approval role or signer")
	ErrProcurementApprovalThreshold = sdkerrors.Register(ModuleName, 8, "insufficient procurement approvals for liability requirements")
	ErrProcurementApprovalScore     = sdkerrors.Register(ModuleName, 9, "procurement approval accountability score below threshold")
)

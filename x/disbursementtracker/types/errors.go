package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrDisbursementNotFound          = sdkerrors.Register(ModuleName, 1, "disbursement not found")
	ErrInvalidDisbursementTransition = sdkerrors.Register(ModuleName, 2, "invalid disbursement status transition")
	ErrPerformedByAddress            = sdkerrors.Register(ModuleName, 3, "invalid performed_by address")
	ErrProcurementMissing            = sdkerrors.Register(ModuleName, 4, "procurement not found for disbursement")
	ErrDisbursementAmountExceeded    = sdkerrors.Register(ModuleName, 5, "disbursement amount exceeds procurement total")
	ErrDisbursementLiabilityMissing  = sdkerrors.Register(ModuleName, 6, "liability contract required for disbursement")
	ErrDisbursementApprovalSigner    = sdkerrors.Register(ModuleName, 7, "invalid disbursement approval signer")
	ErrDisbursementApprovalDuplicate = sdkerrors.Register(ModuleName, 8, "duplicate disbursement approval role or signer")
	ErrDisbursementApprovalThreshold = sdkerrors.Register(ModuleName, 9, "insufficient disbursement approvals for liability requirements")
	ErrDisbursementApprovalScore     = sdkerrors.Register(ModuleName, 10, "disbursement approval accountability score below threshold")
)

package types

import (
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"

	accountabilitytypes "govchain/x/accountabilityscores/types"
)

// DisbursementStatus tracks the reconciliation lifecycle of a disbursement.
type DisbursementStatus string

const (
	DisbursementStatusScheduled DisbursementStatus = "DISBURSEMENT_STATUS_SCHEDULED"
	DisbursementStatusReleased  DisbursementStatus = "DISBURSEMENT_STATUS_RELEASED"
	DisbursementStatusVerified  DisbursementStatus = "DISBURSEMENT_STATUS_VERIFIED"
	DisbursementStatusDisputed  DisbursementStatus = "DISBURSEMENT_STATUS_DISPUTED"
	DisbursementStatusResolved  DisbursementStatus = "DISBURSEMENT_STATUS_RESOLVED"
)

var disbursementTransitions = map[DisbursementStatus][]DisbursementStatus{
	DisbursementStatusScheduled: {DisbursementStatusReleased, DisbursementStatusDisputed},
	DisbursementStatusReleased:  {DisbursementStatusVerified, DisbursementStatusDisputed},
	DisbursementStatusVerified:  {DisbursementStatusResolved},
	DisbursementStatusDisputed:  {DisbursementStatusResolved},
	DisbursementStatusResolved:  {},
}

// Disbursement records a release of funds tied to a procurement.
type Disbursement struct {
	Id                   uint64                         `json:"id" yaml:"id"`
	ProcurementId        uint64                         `json:"procurement_id" yaml:"procurement_id"`
	Amount               sdkmath.Int                    `json:"amount" yaml:"amount"`
	Currency             string                         `json:"currency" yaml:"currency"`
	Status               DisbursementStatus             `json:"status" yaml:"status"`
	EvidenceURI          string                         `json:"evidence_uri" yaml:"evidence_uri"`
	Notes                string                         `json:"notes" yaml:"notes"`
	PerformedBy          string                         `json:"performed_by" yaml:"performed_by"`
	LiabilityContractURI string                         `json:"liability_contract_uri" yaml:"liability_contract_uri"`
	Approvals            []accountabilitytypes.Approval `json:"approvals" yaml:"approvals"`
}

// ValidateBasic ensures the disbursement entry is consistent.
func (d Disbursement) ValidateBasic() error {
	if d.Id == 0 {
		return fmt.Errorf("disbursement id must be set")
	}
	if d.ProcurementId == 0 {
		return fmt.Errorf("procurement id is required")
	}
	if !d.Amount.IsPositive() {
		return fmt.Errorf("amount must be positive")
	}
	if strings.TrimSpace(d.Currency) == "" {
		return fmt.Errorf("currency is required")
	}
	if strings.TrimSpace(d.PerformedBy) == "" {
		return fmt.Errorf("performed_by is required")
	}
	if _, ok := disbursementTransitions[d.Status]; !ok {
		return fmt.Errorf("unknown disbursement status: %s", d.Status)
	}
	if strings.TrimSpace(d.LiabilityContractURI) == "" && d.Status != DisbursementStatusScheduled {
		return fmt.Errorf("liability contract uri required once disbursement leaves scheduled")
	}
	for _, approval := range d.Approvals {
		if err := approval.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

// CanTransitionTo validates status changes.
func (s DisbursementStatus) CanTransitionTo(next DisbursementStatus) bool {
	allowed, ok := disbursementTransitions[s]
	if !ok {
		return false
	}
	for _, candidate := range allowed {
		if candidate == next {
			return true
		}
	}
	return false
}

// DefaultDisbursementStatus returns the default starting status.
func DefaultDisbursementStatus() DisbursementStatus {
	return DisbursementStatusScheduled
}

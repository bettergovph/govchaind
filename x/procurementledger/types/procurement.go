package types

import (
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"

	accountabilitytypes "govchain/x/accountabilityscores/types"
)

// ProcurementStatus captures the lifecycle of a procurement record.
type ProcurementStatus string

const (
	ProcurementStatusPlanning  ProcurementStatus = "PROCUREMENT_STATUS_PLANNING"
	ProcurementStatusTendering ProcurementStatus = "PROCUREMENT_STATUS_TENDERING"
	ProcurementStatusAwarded   ProcurementStatus = "PROCUREMENT_STATUS_AWARDED"
	ProcurementStatusExecuting ProcurementStatus = "PROCUREMENT_STATUS_EXECUTING"
	ProcurementStatusCompleted ProcurementStatus = "PROCUREMENT_STATUS_COMPLETED"
	ProcurementStatusCancelled ProcurementStatus = "PROCUREMENT_STATUS_CANCELLED"
)

var procurementTransitions = map[ProcurementStatus][]ProcurementStatus{
	ProcurementStatusPlanning:  {ProcurementStatusTendering, ProcurementStatusCancelled},
	ProcurementStatusTendering: {ProcurementStatusAwarded, ProcurementStatusCancelled},
	ProcurementStatusAwarded:   {ProcurementStatusExecuting, ProcurementStatusCancelled},
	ProcurementStatusExecuting: {ProcurementStatusCompleted, ProcurementStatusCancelled},
	ProcurementStatusCompleted: {},
	ProcurementStatusCancelled: {},
}

// Procurement models a tender or contract linked to a budget allocation.
type Procurement struct {
	Id          uint64            `json:"id" yaml:"id"`
	BudgetId    uint64            `json:"budget_id" yaml:"budget_id"`
	Reference   string            `json:"reference" yaml:"reference"`
	Agency      string            `json:"agency" yaml:"agency"`
	Title       string            `json:"title" yaml:"title"`
	Category    string            `json:"category" yaml:"category"`
	Amount      sdkmath.Int       `json:"amount" yaml:"amount"`
	Currency    string            `json:"currency" yaml:"currency"`
	Supplier    string            `json:"supplier" yaml:"supplier"`
	Status      ProcurementStatus `json:"status" yaml:"status"`
	MetadataURI string            `json:"metadata_uri" yaml:"metadata_uri"`
	Officer     string            `json:"officer" yaml:"officer"`
	// LiabilityContractURI references the contract accepted by approvers
	// before the procurement can advance past tendering.
	LiabilityContractURI string                         `json:"liability_contract_uri" yaml:"liability_contract_uri"`
	Approvals            []accountabilitytypes.Approval `json:"approvals" yaml:"approvals"`
}

// ValidateBasic ensures the procurement entry is well formed.
func (p Procurement) ValidateBasic() error {
	if p.Id == 0 {
		return fmt.Errorf("procurement id must be set")
	}
	if p.BudgetId == 0 {
		return fmt.Errorf("budget id is required")
	}
	if strings.TrimSpace(p.Reference) == "" {
		return fmt.Errorf("reference is required")
	}
	if strings.TrimSpace(p.Agency) == "" {
		return fmt.Errorf("agency is required")
	}
	if strings.TrimSpace(p.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if !p.Amount.IsPositive() {
		return fmt.Errorf("amount must be positive")
	}
	if strings.TrimSpace(p.Currency) == "" {
		return fmt.Errorf("currency is required")
	}
	if strings.TrimSpace(p.Officer) == "" {
		return fmt.Errorf("procurement officer is required")
	}
	if _, ok := procurementTransitions[p.Status]; !ok {
		return fmt.Errorf("unknown procurement status: %s", p.Status)
	}
	if strings.TrimSpace(p.LiabilityContractURI) == "" && p.Status != ProcurementStatusPlanning {
		return fmt.Errorf("liability contract uri required once procurement leaves planning")
	}
	for _, approval := range p.Approvals {
		if err := approval.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

// CanTransitionTo checks if a status transition is valid.
func (s ProcurementStatus) CanTransitionTo(next ProcurementStatus) bool {
	allowed, ok := procurementTransitions[s]
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

// DefaultProcurementStatus returns the default status for new procurements.
func DefaultProcurementStatus() ProcurementStatus {
	return ProcurementStatusPlanning
}

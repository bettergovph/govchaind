package types

import (
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"

	accountabilitytypes "govchain/x/accountabilityscores/types"
)

// BudgetStatus captures the lifecycle stage of a budget record.
type BudgetStatus string

const (
	BudgetStatusDraft      BudgetStatus = "BUDGET_STATUS_DRAFT"
	BudgetStatusLegislated BudgetStatus = "BUDGET_STATUS_LEGISLATED"
	BudgetStatusExecuting  BudgetStatus = "BUDGET_STATUS_EXECUTING"
	BudgetStatusCompleted  BudgetStatus = "BUDGET_STATUS_COMPLETED"
	BudgetStatusArchived   BudgetStatus = "BUDGET_STATUS_ARCHIVED"
)

var budgetStatusTransitions = map[BudgetStatus][]BudgetStatus{
	BudgetStatusDraft:      {BudgetStatusLegislated, BudgetStatusArchived},
	BudgetStatusLegislated: {BudgetStatusExecuting, BudgetStatusArchived},
	BudgetStatusExecuting:  {BudgetStatusCompleted, BudgetStatusArchived},
	BudgetStatusCompleted:  {BudgetStatusArchived},
	BudgetStatusArchived:   {},
}

// Budget models a planned allocation for a specific fiscal period.
type Budget struct {
	Id          uint64       `json:"id" yaml:"id"`
	Agency      string       `json:"agency" yaml:"agency"`
	FiscalYear  string       `json:"fiscal_year" yaml:"fiscal_year"`
	Title       string       `json:"title" yaml:"title"`
	Amount      sdkmath.Int  `json:"amount" yaml:"amount"`
	Currency    string       `json:"currency" yaml:"currency"`
	MetadataURI string       `json:"metadata_uri" yaml:"metadata_uri"`
	Status      BudgetStatus `json:"status" yaml:"status"`
	CreatedBy   string       `json:"created_by" yaml:"created_by"`
	// LiabilityContractURI points to the contract that signers agree to
	// when approving this record. It can be populated before or alongside
	// approvals and is required once the budget leaves the draft stage.
	LiabilityContractURI string `json:"liability_contract_uri" yaml:"liability_contract_uri"`
	// Approvals captures the multi-party endorsements tied to the
	// liability contract.
	Approvals []accountabilitytypes.Approval `json:"approvals" yaml:"approvals"`
}

// ValidateBasic performs stateless validation on a budget entry.
func (b Budget) ValidateBasic() error {
	if b.Id == 0 {
		return fmt.Errorf("budget identifier must be set")
	}
	if strings.TrimSpace(b.Agency) == "" {
		return fmt.Errorf("agency is required")
	}
	if strings.TrimSpace(b.FiscalYear) == "" {
		return fmt.Errorf("fiscal year is required")
	}
	if strings.TrimSpace(b.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if !b.Amount.IsPositive() {
		return fmt.Errorf("amount must be positive")
	}
	if strings.TrimSpace(b.Currency) == "" {
		return fmt.Errorf("currency is required")
	}
	if _, ok := budgetStatusTransitions[b.Status]; !ok {
		return fmt.Errorf("unknown budget status: %s", b.Status)
	}
	if strings.TrimSpace(b.CreatedBy) == "" {
		return fmt.Errorf("creator address is required")
	}
	if strings.TrimSpace(b.LiabilityContractURI) == "" && b.Status != BudgetStatusDraft {
		return fmt.Errorf("liability contract uri required once budget leaves draft")
	}
	for _, approval := range b.Approvals {
		if err := approval.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

// CanTransitionTo returns true if a status change is permitted.
func (s BudgetStatus) CanTransitionTo(next BudgetStatus) bool {
	allowed, ok := budgetStatusTransitions[s]
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

// DefaultBudgetStatus returns the default lifecycle state for new budgets.
func DefaultBudgetStatus() BudgetStatus {
	return BudgetStatusDraft
}

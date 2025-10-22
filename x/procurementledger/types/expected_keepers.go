package types

import (
	"context"

	accountabilitytypes "govchain/x/accountabilityscores/types"
	budgettypes "govchain/x/budgetregistry/types"
)

// BudgetKeeper defines the subset of the budget registry keeper used by the procurement module.
type BudgetKeeper interface {
	GetBudget(context.Context, uint64) (budgettypes.Budget, error)
}

// AccountabilityKeeper exposes scorecard lookups for approval validation.
type AccountabilityKeeper interface {
	GetScorecard(context.Context, uint64) (accountabilitytypes.Scorecard, error)
}

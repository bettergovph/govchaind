package keeper

import (
	"context"
	"fmt"

	"govchain/x/budgetregistry/types"
)

// InitGenesis initialises state from genesis data.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := genState.Validate(); err != nil {
		return err
	}

	for _, budget := range genState.Budgets {
		if _, err := k.Budgets.Get(ctx, budget.Id); err == nil {
			return fmt.Errorf("duplicate budget id %d during genesis", budget.Id)
		}
		if err := k.Budgets.Set(ctx, budget.Id, budget); err != nil {
			return err
		}
	}

	if err := k.BudgetSeq.Set(ctx, genState.BudgetCount); err != nil {
		return err
	}
	return nil
}

// ExportGenesis exports the current state into a genesis structure.
func (k Keeper) ExportGenesis(ctx context.Context) (types.GenesisState, error) {
	count, err := k.BudgetSeq.Peek(ctx)
	if err != nil {
		return types.GenesisState{}, err
	}

	budgets := make([]types.Budget, 0)
	if err := k.WalkBudgets(ctx, func(budget types.Budget) (bool, error) {
		budgets = append(budgets, budget)
		return false, nil
	}); err != nil {
		return types.GenesisState{}, err
	}

	return types.GenesisState{
		Budgets:     budgets,
		BudgetCount: count,
	}, nil
}

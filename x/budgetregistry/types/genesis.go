package types

import "fmt"

// GenesisState defines the initial module state for the budget registry module.
type GenesisState struct {
	Budgets     []Budget `json:"budgets" yaml:"budgets"`
	BudgetCount uint64   `json:"budget_count" yaml:"budget_count"`
}

// DefaultGenesis returns a default genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Budgets:     []Budget{},
		BudgetCount: 0,
	}
}

// Validate performs basic genesis state validation checks.
func (gs GenesisState) Validate() error {
	seen := make(map[uint64]struct{}, len(gs.Budgets))
	for _, budget := range gs.Budgets {
		if _, exists := seen[budget.Id]; exists {
			return fmt.Errorf("duplicate budget id %d in genesis", budget.Id)
		}
		seen[budget.Id] = struct{}{}
		if err := budget.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid budget %d: %w", budget.Id, err)
		}
	}
	if uint64(len(gs.Budgets)) > gs.BudgetCount {
		return fmt.Errorf("budget count %d smaller than budgets supplied %d", gs.BudgetCount, len(gs.Budgets))
	}
	return nil
}

package types

import "fmt"

// GenesisState defines the initial governance voting state.
type GenesisState struct {
	Delegations     []Delegation `json:"delegations" yaml:"delegations"`
	DelegationCount uint64       `json:"delegation_count" yaml:"delegation_count"`
}

// DefaultGenesis returns an empty genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Delegations:     []Delegation{},
		DelegationCount: 0,
	}
}

// Validate ensures the genesis data is valid.
func (gs GenesisState) Validate() error {
	seen := make(map[uint64]struct{}, len(gs.Delegations))
	index := make(map[string]struct{}, len(gs.Delegations))
	for _, delegation := range gs.Delegations {
		if _, exists := seen[delegation.Id]; exists {
			return fmt.Errorf("duplicate delegation id %d in genesis", delegation.Id)
		}
		seen[delegation.Id] = struct{}{}
		if err := delegation.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid delegation %d: %w", delegation.Id, err)
		}
		key := delegation.IndexKey()
		if _, exists := index[key]; exists {
			return fmt.Errorf("duplicate delegation index %s in genesis", key)
		}
		index[key] = struct{}{}
	}
	if uint64(len(gs.Delegations)) > gs.DelegationCount {
		return fmt.Errorf("delegation count %d smaller than delegations supplied %d", gs.DelegationCount, len(gs.Delegations))
	}
	return nil
}

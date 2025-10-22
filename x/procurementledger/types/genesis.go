package types

import "fmt"

// GenesisState defines initial state for procurement ledger.
type GenesisState struct {
	Procurements     []Procurement `json:"procurements" yaml:"procurements"`
	ProcurementCount uint64        `json:"procurement_count" yaml:"procurement_count"`
}

// DefaultGenesis returns an empty genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Procurements:     []Procurement{},
		ProcurementCount: 0,
	}
}

// Validate performs basic genesis validation.
func (gs GenesisState) Validate() error {
	seen := make(map[uint64]struct{}, len(gs.Procurements))
	for _, procurement := range gs.Procurements {
		if _, exists := seen[procurement.Id]; exists {
			return fmt.Errorf("duplicate procurement id %d in genesis", procurement.Id)
		}
		seen[procurement.Id] = struct{}{}
		if err := procurement.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid procurement %d: %w", procurement.Id, err)
		}
	}
	if uint64(len(gs.Procurements)) > gs.ProcurementCount {
		return fmt.Errorf("procurement count %d smaller than procurements supplied %d", gs.ProcurementCount, len(gs.Procurements))
	}
	return nil
}

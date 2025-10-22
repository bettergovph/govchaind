package types

import "fmt"

// GenesisState holds the initial disbursement tracker configuration.
type GenesisState struct {
	Disbursements     []Disbursement `json:"disbursements" yaml:"disbursements"`
	DisbursementCount uint64         `json:"disbursement_count" yaml:"disbursement_count"`
}

// DefaultGenesis returns an empty genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Disbursements:     []Disbursement{},
		DisbursementCount: 0,
	}
}

// Validate performs validation on genesis data.
func (gs GenesisState) Validate() error {
	seen := make(map[uint64]struct{}, len(gs.Disbursements))
	for _, disbursement := range gs.Disbursements {
		if _, exists := seen[disbursement.Id]; exists {
			return fmt.Errorf("duplicate disbursement id %d in genesis", disbursement.Id)
		}
		seen[disbursement.Id] = struct{}{}
		if err := disbursement.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid disbursement %d: %w", disbursement.Id, err)
		}
	}
	if uint64(len(gs.Disbursements)) > gs.DisbursementCount {
		return fmt.Errorf("disbursement count %d smaller than disbursements supplied %d", gs.DisbursementCount, len(gs.Disbursements))
	}
	return nil
}

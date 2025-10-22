package types

// DefaultGenesis returns the default genesis state for the wasm module.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		Contracts:     []Contract{},
		ContractCount: 0,
	}
}

// GenesisState defines the initial state of the wasm module.
type GenesisState struct {
	Params        Params     `json:"params" yaml:"params"`
	Contracts     []Contract `json:"contracts" yaml:"contracts"`
	ContractCount uint64     `json:"contract_count" yaml:"contract_count"`
}

// Validate performs basic genesis state validation returning an error upon failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	seen := make(map[uint64]struct{}, len(gs.Contracts))
	for _, contract := range gs.Contracts {
		if err := contract.ValidateBasic(); err != nil {
			return err
		}
		if _, exists := seen[contract.Id]; exists {
			return ErrDuplicateContractID
		}
		seen[contract.Id] = struct{}{}
	}
	return nil
}

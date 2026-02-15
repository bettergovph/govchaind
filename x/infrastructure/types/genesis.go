package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:      DefaultParams(),
		ProjectList: []InfraProject{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	projectIdMap := make(map[string]bool)
	for _, elem := range gs.ProjectList {
		if _, ok := projectIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for infrastructure project")
		}
		projectIdMap[elem.Id] = true
	}

	return gs.Params.Validate()
}

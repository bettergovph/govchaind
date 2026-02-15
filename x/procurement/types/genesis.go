package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:      DefaultParams(),
		ProcessList: []ContractingProcess{},
		EntityList:  []EntityRegistry{},
		ReleaseList: []Release{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	processIdMap := make(map[string]bool)
	for _, elem := range gs.ProcessList {
		if _, ok := processIdMap[elem.Ocid]; ok {
			return fmt.Errorf("duplicated ocid for contracting process")
		}
		processIdMap[elem.Ocid] = true
	}

	entityIdMap := make(map[string]bool)
	for _, elem := range gs.EntityList {
		if _, ok := entityIdMap[elem.EntityId]; ok {
			return fmt.Errorf("duplicated id for entity")
		}
		entityIdMap[elem.EntityId] = true
	}

	return gs.Params.Validate()
}

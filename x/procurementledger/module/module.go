package procurementledger

import "govchain/x/procurementledger/types"

// AppModule will host message routing, queries, and genesis management for procurement data.
type AppModule struct{}

// NewAppModule constructs the placeholder AppModule.
func NewAppModule() AppModule {
	return AppModule{}
}

// Name returns the module name for registration purposes.
func (AppModule) Name() string {
	return types.ModuleName
}

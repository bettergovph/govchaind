package budgetregistry

import "govchain/x/budgetregistry/types"

// AppModule defines the root module type. It will implement the Cosmos SDK interfaces as the
// keeper, message server, and query server mature.
type AppModule struct{}

// NewAppModule returns a placeholder AppModule value.
func NewAppModule() AppModule {
	return AppModule{}
}

// Name satisfies the AppModule interface contract and allows the module to register with the
// application once wiring is introduced.
func (AppModule) Name() string {
	return types.ModuleName
}

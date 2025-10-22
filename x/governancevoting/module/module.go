package governancevoting

import "govchain/x/governancevoting/types"

// AppModule encapsulates the governance voting integration points with the Cosmos SDK.
type AppModule struct{}

// NewAppModule constructs the placeholder governance voting module.
func NewAppModule() AppModule {
	return AppModule{}
}

// Name returns the module name.
func (AppModule) Name() string {
	return types.ModuleName
}

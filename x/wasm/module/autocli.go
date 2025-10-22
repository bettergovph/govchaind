package wasm

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	// Smart contract transactions will be added in a future iteration.
	return &autocliv1.ModuleOptions{}
}

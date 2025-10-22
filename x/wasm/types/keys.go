package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "wasm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	GovModuleName = "gov"
)

// ParamsKey is the prefix to retrieve module parameters.
var ParamsKey = collections.NewPrefix("p_wasm")

var (
	ContractKey      = collections.NewPrefix("contract/value/")
	ContractCountKey = collections.NewPrefix("contract/count/")
)

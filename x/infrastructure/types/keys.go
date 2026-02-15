package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "infrastructure"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	GovModuleName = "gov"
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix("p_infrastructure")

var (
	// InfraProject storage keys
	ProjectKey      = collections.NewPrefix("project/value/")
	ProjectCountKey = collections.NewPrefix("project/count/")
)

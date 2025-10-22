package types

import "cosmossdk.io/collections"

const (
	ModuleName  = "procurementledger"
	StoreKey    = ModuleName
	RouterKey   = ModuleName
	MemStoreKey = "mem_procurementledger"
)

var (
	ProcurementKey      = collections.NewPrefix("procurement/value/")
	ProcurementCountKey = collections.NewPrefix("procurement/count/")
)

// KeyPrefix is retained for backwards compatibility with scaffolding helpers.
func KeyPrefix(key string) []byte {
	return []byte(key)
}

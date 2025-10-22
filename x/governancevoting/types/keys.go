package types

import "cosmossdk.io/collections"

const (
	ModuleName  = "governancevoting"
	StoreKey    = ModuleName
	RouterKey   = ModuleName
	MemStoreKey = "mem_governancevoting"
)

var (
	DelegationKey      = collections.NewPrefix("delegation/value/")
	DelegationCountKey = collections.NewPrefix("delegation/count/")
	DelegationIndexKey = collections.NewPrefix("delegation/index/")
)

// KeyPrefix maintained for compatibility with scaffolding helpers.
func KeyPrefix(key string) []byte {
	return []byte(key)
}

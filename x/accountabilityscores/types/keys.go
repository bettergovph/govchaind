package types

import "cosmossdk.io/collections"

const (
	ModuleName  = "accountabilityscores"
	StoreKey    = ModuleName
	RouterKey   = ModuleName
	MemStoreKey = "mem_accountabilityscores"
)

var (
	ScorecardKey      = collections.NewPrefix("scorecard/value/")
	ScorecardCountKey = collections.NewPrefix("scorecard/count/")
	ScorecardIndexKey = collections.NewPrefix("scorecard/index/")
)

// KeyPrefix converts strings into byte prefixes for compatibility with scaffolding.
func KeyPrefix(key string) []byte {
	return []byte(key)
}

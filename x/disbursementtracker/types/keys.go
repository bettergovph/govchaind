package types

import "cosmossdk.io/collections"

const (
	ModuleName  = "disbursementtracker"
	StoreKey    = ModuleName
	RouterKey   = ModuleName
	MemStoreKey = "mem_disbursementtracker"
)

var (
	DisbursementKey      = collections.NewPrefix("disbursement/value/")
	DisbursementCountKey = collections.NewPrefix("disbursement/count/")
	DisbursementTotalKey = collections.NewPrefix("disbursement/total/")
)

// KeyPrefix kept for compatibility with scaffolding helpers.
func KeyPrefix(key string) []byte {
	return []byte(key)
}

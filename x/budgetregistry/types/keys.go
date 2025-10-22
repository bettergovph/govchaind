package types

import "cosmossdk.io/collections"

const (
	ModuleName  = "budgetregistry"
	StoreKey    = ModuleName
	RouterKey   = ModuleName
	MemStoreKey = "mem_budgetregistry"
)

var (
	BudgetKey      = collections.NewPrefix("budget/value/")
	BudgetCountKey = collections.NewPrefix("budget/count/")
)

// KeyPrefix turns a key string into a byte prefix for persistent storage maps.
func KeyPrefix(key string) []byte {
	return []byte(key)
}

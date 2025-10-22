package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"govchain/x/wasm/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	authority    []byte

	Schema      collections.Schema
	Params      collections.Item[types.Params]
	ContractSeq collections.Sequence
	Contracts   collections.Map[uint64, types.Contract]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		Params:       collections.NewItem(sb, types.ParamsKey, "params", collections.NewJSONValueCodec[types.Params]()),
		Contracts:    collections.NewMap(sb, types.ContractKey, "contract", collections.Uint64Key, collections.NewJSONValueCodec[types.Contract]()),
		ContractSeq:  collections.NewSequence(sb, types.ContractCountKey, "contractSequence"),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority address bytes.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

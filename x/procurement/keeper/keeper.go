package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"govchain/x/procurement/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema       collections.Schema
	Params       collections.Item[types.Params]
	ProcessSeq   collections.Sequence
	Process      collections.Map[string, types.ContractingProcess]
	ReleaseSeq   collections.Sequence
	Release      collections.Map[string, types.Release]
	EntitySeq    collections.Sequence
	Entity       collections.Map[string, types.EntityRegistry]
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

		Params:     collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Process:    collections.NewMap(sb, types.ProcessKey, "process", collections.StringKey, codec.CollValue[types.ContractingProcess](cdc)),
		ProcessSeq: collections.NewSequence(sb, types.ProcessCountKey, "processSequence"),
		Release:    collections.NewMap(sb, types.ReleaseKey, "release", collections.StringKey, codec.CollValue[types.Release](cdc)),
		ReleaseSeq: collections.NewSequence(sb, types.ReleaseCountKey, "releaseSequence"),
		Entity:     collections.NewMap(sb, types.EntityKey, "entity", collections.StringKey, codec.CollValue[types.EntityRegistry](cdc)),
		EntitySeq:  collections.NewSequence(sb, types.EntityCountKey, "entitySequence"),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

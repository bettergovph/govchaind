package keeper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	collcodec "cosmossdk.io/collections/codec"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"govchain/x/governancevoting/types"
)

// Keeper manages proposal metadata, delegation records, and quorum enforcement.
type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec

	Schema          collections.Schema
	DelegationSeq   collections.Sequence
	Delegations     collections.Map[uint64, types.Delegation]
	DelegationIndex collections.Map[string, uint64]
}

// NewKeeper returns a keeper instance with dependencies configured.
func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:    storeService,
		cdc:             cdc,
		addressCodec:    addressCodec,
		Delegations:     collections.NewMap(sb, types.DelegationKey, "delegation", collections.Uint64Key, collections.NewJSONValueCodec[types.Delegation]()),
		DelegationSeq:   collections.NewSequence(sb, types.DelegationCountKey, "delegation_seq"),
		DelegationIndex: collections.NewMap(sb, types.DelegationIndexKey, "delegation_index", collections.StringKey, collcodec.KeyToValueCodec(collcodec.NewUint64Key[uint64]()).WithName("delegation_index_value")),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(fmt.Errorf("failed to build governance voting schema: %w", err))
	}
	k.Schema = schema
	return k
}

// RegisterDelegation records a new delegated accountability relationship.
func (k Keeper) RegisterDelegation(ctx context.Context, delegation types.Delegation) (types.Delegation, error) {
	if _, err := k.addressCodec.StringToBytes(delegation.Delegator); err != nil {
		return types.Delegation{}, fmt.Errorf("%w: %v", types.ErrDelegationInvalidAddr, err)
	}
	if _, err := k.addressCodec.StringToBytes(delegation.Delegatee); err != nil {
		return types.Delegation{}, fmt.Errorf("%w: %v", types.ErrDelegationInvalidAddr, err)
	}
	if delegation.Active == false {
		delegation.Active = true
	}
	key := delegation.IndexKey()
	if existingID, err := k.DelegationIndex.Get(ctx, key); err == nil {
		existing, err := k.Delegations.Get(ctx, existingID)
		if err != nil {
			return types.Delegation{}, err
		}
		delegation.Id = existing.Id
	} else if !errors.Is(err, collections.ErrNotFound) {
		return types.Delegation{}, err
	}

	if delegation.Id == 0 {
		id, err := k.DelegationSeq.Next(ctx)
		if err != nil {
			return types.Delegation{}, err
		}
		delegation.Id = id + 1
	}

	if err := delegation.ValidateBasic(); err != nil {
		return types.Delegation{}, err
	}

	if err := k.Delegations.Set(ctx, delegation.Id, delegation); err != nil {
		return types.Delegation{}, err
	}
	if err := k.DelegationIndex.Set(ctx, key, delegation.Id); err != nil {
		return types.Delegation{}, err
	}
	return delegation, nil
}

// DeactivateDelegation marks a delegation as inactive.
func (k Keeper) DeactivateDelegation(ctx context.Context, id uint64) (types.Delegation, error) {
	delegation, err := k.Delegations.Get(ctx, id)
	if err != nil {
		return types.Delegation{}, fmt.Errorf("%w: %v", types.ErrDelegationNotFound, err)
	}
	if !delegation.Active {
		return delegation, nil
	}
	delegation.Active = false
	if err := k.Delegations.Set(ctx, delegation.Id, delegation); err != nil {
		return types.Delegation{}, err
	}
	return delegation, nil
}

// GetDelegation returns a stored delegation by identifier.
func (k Keeper) GetDelegation(ctx context.Context, id uint64) (types.Delegation, error) {
	delegation, err := k.Delegations.Get(ctx, id)
	if err != nil {
		return types.Delegation{}, fmt.Errorf("%w: %v", types.ErrDelegationNotFound, err)
	}
	return delegation, nil
}

// GetDelegationByScope resolves the delegation for a delegator and scope.
func (k Keeper) GetDelegationByScope(ctx context.Context, delegator, scope string) (types.Delegation, error) {
	key := strings.ToLower(strings.TrimSpace(delegator)) + "|" + strings.ToLower(strings.TrimSpace(scope))
	id, err := k.DelegationIndex.Get(ctx, key)
	if err != nil {
		return types.Delegation{}, fmt.Errorf("%w: %v", types.ErrDelegationNotFound, err)
	}
	return k.GetDelegation(ctx, id)
}

// WalkDelegations iterates all stored delegations.
func (k Keeper) WalkDelegations(ctx context.Context, cb func(types.Delegation) (bool, error)) error {
	return k.Delegations.Walk(ctx, nil, func(_ uint64, delegation types.Delegation) (bool, error) {
		return cb(delegation)
	})
}

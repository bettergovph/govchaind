package keeper

import (
	"context"
	"fmt"

	"govchain/x/governancevoting/types"
)

// InitGenesis loads state from genesis data.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := genState.Validate(); err != nil {
		return err
	}

	for _, delegation := range genState.Delegations {
		if _, err := k.Delegations.Get(ctx, delegation.Id); err == nil {
			return fmt.Errorf("duplicate delegation id %d during genesis", delegation.Id)
		}
		if err := k.Delegations.Set(ctx, delegation.Id, delegation); err != nil {
			return err
		}
		if err := k.DelegationIndex.Set(ctx, delegation.IndexKey(), delegation.Id); err != nil {
			return err
		}
	}

	if err := k.DelegationSeq.Set(ctx, genState.DelegationCount); err != nil {
		return err
	}
	return nil
}

// ExportGenesis exports state to genesis representation.
func (k Keeper) ExportGenesis(ctx context.Context) (types.GenesisState, error) {
	count, err := k.DelegationSeq.Peek(ctx)
	if err != nil {
		return types.GenesisState{}, err
	}

	delegations := make([]types.Delegation, 0)
	if err := k.WalkDelegations(ctx, func(delegation types.Delegation) (bool, error) {
		delegations = append(delegations, delegation)
		return false, nil
	}); err != nil {
		return types.GenesisState{}, err
	}

	return types.GenesisState{
		Delegations:     delegations,
		DelegationCount: count,
	}, nil
}

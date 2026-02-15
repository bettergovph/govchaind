package keeper

import (
	"context"

	"govchain/x/procurement/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.ProcessList {
		if err := k.Process.Set(ctx, elem.Ocid, elem); err != nil {
			return err
		}
	}
	if err := k.ProcessSeq.Set(ctx, genState.ProcessCount); err != nil {
		return err
	}

	for _, elem := range genState.EntityList {
		if err := k.Entity.Set(ctx, elem.EntityId, elem); err != nil {
			return err
		}
	}
	if err := k.EntitySeq.Set(ctx, genState.EntityCount); err != nil {
		return err
	}

	for _, elem := range genState.ReleaseList {
		releaseKey := elem.Ocid + "/" + elem.Id
		if err := k.Release.Set(ctx, releaseKey, elem); err != nil {
			return err
		}
	}
	if err := k.ReleaseSeq.Set(ctx, genState.ReleaseCount); err != nil {
		return err
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	err = k.Process.Walk(ctx, nil, func(key string, elem types.ContractingProcess) (bool, error) {
		genesis.ProcessList = append(genesis.ProcessList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	genesis.ProcessCount, err = k.ProcessSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}

	err = k.Entity.Walk(ctx, nil, func(key string, elem types.EntityRegistry) (bool, error) {
		genesis.EntityList = append(genesis.EntityList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	genesis.EntityCount, err = k.EntitySeq.Peek(ctx)
	if err != nil {
		return nil, err
	}

	err = k.Release.Walk(ctx, nil, func(key string, elem types.Release) (bool, error) {
		genesis.ReleaseList = append(genesis.ReleaseList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	genesis.ReleaseCount, err = k.ReleaseSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}

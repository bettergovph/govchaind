package keeper

import (
	"context"

	"govchain/x/infrastructure/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.ProjectList {
		if err := k.Project.Set(ctx, elem.Id, elem); err != nil {
			return err
		}
	}

	if err := k.ProjectSeq.Set(ctx, genState.ProjectCount); err != nil {
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

	err = k.Project.Walk(ctx, nil, func(key string, elem types.InfraProject) (bool, error) {
		genesis.ProjectList = append(genesis.ProjectList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	genesis.ProjectCount, err = k.ProjectSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}

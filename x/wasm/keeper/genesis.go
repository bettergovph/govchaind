package keeper

import (
	"context"

	"govchain/x/wasm/types"
)

// InitGenesis initializes the module state from genesis data.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	for _, contract := range genState.Contracts {
		if err := k.Contracts.Set(ctx, contract.Id, contract); err != nil {
			return err
		}
	}

	if err := k.ContractSeq.Set(ctx, genState.ContractCount); err != nil {
		return err
	}

	return nil
}

// ExportGenesis exports the current module state for genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	count, err := k.ContractSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}

	contracts := make([]types.Contract, 0)
	if err := k.WalkContracts(ctx, func(contract types.Contract) (bool, error) {
		contracts = append(contracts, contract)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return &types.GenesisState{
		Params:        params,
		Contracts:     contracts,
		ContractCount: count,
	}, nil
}

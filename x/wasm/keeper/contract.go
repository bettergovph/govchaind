package keeper

import (
	"context"

	"govchain/x/wasm/types"
)

// RegisterContract stores a new smart contract metadata entry and returns its identifier.
func (k Keeper) RegisterContract(ctx context.Context, contract types.Contract) (types.Contract, error) {
	id, err := k.ContractSeq.Next(ctx)
	if err != nil {
		return types.Contract{}, err
	}

	contract.Id = id + 1
	if err := contract.ValidateBasic(); err != nil {
		return types.Contract{}, err
	}

	if err := k.Contracts.Set(ctx, contract.Id, contract); err != nil {
		return types.Contract{}, err
	}

	return contract, nil
}

// GetContract returns the stored metadata for a contract identifier.
func (k Keeper) GetContract(ctx context.Context, id uint64) (types.Contract, error) {
	return k.Contracts.Get(ctx, id)
}

// WalkContracts iterates over all stored contracts and executes the provided callback.
func (k Keeper) WalkContracts(ctx context.Context, cb func(types.Contract) (bool, error)) error {
	return k.Contracts.Walk(ctx, nil, func(_ uint64, contract types.Contract) (bool, error) {
		return cb(contract)
	})
}

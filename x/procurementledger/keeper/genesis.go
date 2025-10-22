package keeper

import (
	"context"
	"fmt"

	"govchain/x/procurementledger/types"
)

// InitGenesis initialises the procurement ledger from genesis data.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := genState.Validate(); err != nil {
		return err
	}

	for _, procurement := range genState.Procurements {
		if _, err := k.Procurements.Get(ctx, procurement.Id); err == nil {
			return fmt.Errorf("duplicate procurement id %d during genesis", procurement.Id)
		}
		if err := k.Procurements.Set(ctx, procurement.Id, procurement); err != nil {
			return err
		}
	}

	if err := k.ProcurementSeq.Set(ctx, genState.ProcurementCount); err != nil {
		return err
	}
	return nil
}

// ExportGenesis exports the procurement ledger state to genesis representation.
func (k Keeper) ExportGenesis(ctx context.Context) (types.GenesisState, error) {
	count, err := k.ProcurementSeq.Peek(ctx)
	if err != nil {
		return types.GenesisState{}, err
	}

	procurements := make([]types.Procurement, 0)
	if err := k.WalkProcurements(ctx, func(procurement types.Procurement) (bool, error) {
		procurements = append(procurements, procurement)
		return false, nil
	}); err != nil {
		return types.GenesisState{}, err
	}

	return types.GenesisState{
		Procurements:     procurements,
		ProcurementCount: count,
	}, nil
}

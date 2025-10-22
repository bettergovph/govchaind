package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"

	"govchain/x/disbursementtracker/types"
)

// InitGenesis initializes state from genesis data.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := genState.Validate(); err != nil {
		return err
	}

	totals := make(map[uint64]sdkmath.Int)
	for _, disbursement := range genState.Disbursements {
		if _, err := k.Disbursements.Get(ctx, disbursement.Id); err == nil {
			return fmt.Errorf("duplicate disbursement id %d during genesis", disbursement.Id)
		}
		if err := k.Disbursements.Set(ctx, disbursement.Id, disbursement); err != nil {
			return err
		}
		total := totals[disbursement.ProcurementId]
		total = total.Add(disbursement.Amount)
		totals[disbursement.ProcurementId] = total
	}

	for procurementID, total := range totals {
		if err := k.DisbursementTotals.Set(ctx, procurementID, total); err != nil {
			return err
		}
	}

	if err := k.DisbursementSeq.Set(ctx, genState.DisbursementCount); err != nil {
		return err
	}
	return nil
}

// ExportGenesis exports state into a genesis representation.
func (k Keeper) ExportGenesis(ctx context.Context) (types.GenesisState, error) {
	count, err := k.DisbursementSeq.Peek(ctx)
	if err != nil {
		return types.GenesisState{}, err
	}

	disbursements := make([]types.Disbursement, 0)
	if err := k.WalkDisbursements(ctx, func(disbursement types.Disbursement) (bool, error) {
		disbursements = append(disbursements, disbursement)
		return false, nil
	}); err != nil {
		return types.GenesisState{}, err
	}

	return types.GenesisState{
		Disbursements:     disbursements,
		DisbursementCount: count,
	}, nil
}

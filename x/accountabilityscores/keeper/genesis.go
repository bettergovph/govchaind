package keeper

import (
	"context"
	"fmt"

	"govchain/x/accountabilityscores/types"
)

// InitGenesis initialises the keeper state from genesis data.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := genState.Validate(); err != nil {
		return err
	}

	for _, score := range genState.Scorecards {
		if _, err := k.Scorecards.Get(ctx, score.Id); err == nil {
			return fmt.Errorf("duplicate scorecard id %d during genesis", score.Id)
		}
		if err := k.Scorecards.Set(ctx, score.Id, score); err != nil {
			return err
		}
		if err := k.ScorecardByID.Set(ctx, score.IndexKey(), score.Id); err != nil {
			return err
		}
	}

	if err := k.ScoreSeq.Set(ctx, genState.ScorecardCount); err != nil {
		return err
	}
	return nil
}

// ExportGenesis exports the keeper state into genesis representation.
func (k Keeper) ExportGenesis(ctx context.Context) (types.GenesisState, error) {
	count, err := k.ScoreSeq.Peek(ctx)
	if err != nil {
		return types.GenesisState{}, err
	}

	scorecards := make([]types.Scorecard, 0)
	if err := k.WalkScorecards(ctx, func(score types.Scorecard) (bool, error) {
		scorecards = append(scorecards, score)
		return false, nil
	}); err != nil {
		return types.GenesisState{}, err
	}

	return types.GenesisState{
		Scorecards:     scorecards,
		ScorecardCount: count,
	}, nil
}

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

	"govchain/x/accountabilityscores/types"
)

// Keeper persists accountability scoring data across stakeholders.
type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec

	Schema        collections.Schema
	ScoreSeq      collections.Sequence
	Scorecards    collections.Map[uint64, types.Scorecard]
	ScorecardByID collections.Map[string, uint64]
}

// NewKeeper creates the keeper instance.
func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:  storeService,
		cdc:           cdc,
		addressCodec:  addressCodec,
		Scorecards:    collections.NewMap(sb, types.ScorecardKey, "scorecard", collections.Uint64Key, collections.NewJSONValueCodec[types.Scorecard]()),
		ScoreSeq:      collections.NewSequence(sb, types.ScorecardCountKey, "scorecard_seq"),
		ScorecardByID: collections.NewMap(sb, types.ScorecardIndexKey, "scorecard_index", collections.StringKey, collcodec.KeyToValueCodec(collcodec.NewUint64Key[uint64]()).WithName("scorecard_index_value")),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(fmt.Errorf("failed to build accountability score schema: %w", err))
	}
	k.Schema = schema
	return k
}

// UpsertScorecard creates or updates a scorecard for the subject/metric tuple.
func (k Keeper) UpsertScorecard(ctx context.Context, score types.Scorecard) (types.Scorecard, error) {
	if _, err := k.addressCodec.StringToBytes(score.UpdatedBy); err != nil {
		return types.Scorecard{}, fmt.Errorf("%w: %v", types.ErrScorecardUpdaterAddr, err)
	}

	key := score.IndexKey()
	if score.Id == 0 {
		if existingID, err := k.ScorecardByID.Get(ctx, key); err == nil {
			existing, err := k.Scorecards.Get(ctx, existingID)
			if err != nil {
				return types.Scorecard{}, err
			}
			score.Id = existing.Id
		} else if !errors.Is(err, collections.ErrNotFound) {
			return types.Scorecard{}, err
		}
	}

	if score.Id == 0 {
		id, err := k.ScoreSeq.Next(ctx)
		if err != nil {
			return types.Scorecard{}, err
		}
		score.Id = id + 1
	}

	if err := score.ValidateBasic(); err != nil {
		return types.Scorecard{}, err
	}

	if err := k.Scorecards.Set(ctx, score.Id, score); err != nil {
		return types.Scorecard{}, err
	}
	if err := k.ScorecardByID.Set(ctx, key, score.Id); err != nil {
		return types.Scorecard{}, err
	}
	return score, nil
}

// GetScorecard returns the scorecard by identifier.
func (k Keeper) GetScorecard(ctx context.Context, id uint64) (types.Scorecard, error) {
	score, err := k.Scorecards.Get(ctx, id)
	if err != nil {
		return types.Scorecard{}, fmt.Errorf("%w: %v", types.ErrScorecardNotFound, err)
	}
	return score, nil
}

// GetScorecardBySubjectMetric fetches the scorecard tied to the provided subject/metric combination.
func (k Keeper) GetScorecardBySubjectMetric(ctx context.Context, subject, metric string) (types.Scorecard, error) {
	key := strings.ToLower(strings.TrimSpace(subject)) + "|" + strings.ToLower(strings.TrimSpace(metric))
	id, err := k.ScorecardByID.Get(ctx, key)
	if err != nil {
		return types.Scorecard{}, fmt.Errorf("%w: %v", types.ErrScorecardNotFound, err)
	}
	return k.GetScorecard(ctx, id)
}

// WalkScorecards iterates all stored scorecards.
func (k Keeper) WalkScorecards(ctx context.Context, cb func(types.Scorecard) (bool, error)) error {
	return k.Scorecards.Walk(ctx, nil, func(_ uint64, score types.Scorecard) (bool, error) {
		return cb(score)
	})
}

package keeper_test

import (
	"context"
	"testing"

	address "cosmossdk.io/core/address"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"govchain/x/accountabilityscores/keeper"
	"govchain/x/accountabilityscores/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig()
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_scores")).Ctx

	k := keeper.NewKeeper(storeService, encCfg.Codec, addressCodec)

	return &fixture{
		ctx:          ctx,
		keeper:       k,
		addressCodec: addressCodec,
	}
}

func (f *fixture) newAddress(t *testing.T, raw string) string {
	t.Helper()
	addr, err := f.addressCodec.BytesToString([]byte(raw))
	if err != nil {
		t.Fatalf("failed to convert address: %v", err)
	}
	return addr
}

func TestUpsertScorecard(t *testing.T) {
	f := initFixture(t)
	updater := f.newAddress(t, "scoreUpdaterAddr___________")

	score, err := f.keeper.UpsertScorecard(f.ctx, types.Scorecard{
		Subject:   "DOF",
		Metric:    "Transparency",
		Score:     85,
		Weight:    50,
		UpdatedBy: updater,
	})
	if err != nil {
		t.Fatalf("upsert scorecard: %v", err)
	}
	if score.Id == 0 {
		t.Fatalf("expected id assigned")
	}

	// Update same subject/metric should reuse identifier
	updated, err := f.keeper.UpsertScorecard(f.ctx, types.Scorecard{
		Subject:     "DOF",
		Metric:      "Transparency",
		Score:       90,
		Weight:      60,
		UpdatedBy:   updater,
		EvidenceURI: "ipfs://score",
	})
	if err != nil {
		t.Fatalf("update scorecard: %v", err)
	}
	if updated.Id != score.Id {
		t.Fatalf("expected same id, got %d and %d", updated.Id, score.Id)
	}

	fetched, err := f.keeper.GetScorecardBySubjectMetric(f.ctx, "dof", "transparency")
	if err != nil {
		t.Fatalf("get by subject metric: %v", err)
	}
	if fetched.Score != 90 {
		t.Fatalf("expected score 90, got %d", fetched.Score)
	}
}

package keeper_test

import (
	"context"
	"fmt"
	"testing"

	address "cosmossdk.io/core/address"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	accountabilitykeeper "govchain/x/accountabilityscores/keeper"
	accountabilitytypes "govchain/x/accountabilityscores/types"
	"govchain/x/budgetregistry/keeper"
	"govchain/x/budgetregistry/types"
)

type fixture struct {
	ctx            context.Context
	keeper         keeper.Keeper
	addressCodec   address.Codec
	accountability *accountabilitykeeper.Keeper
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig()
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	keys := map[string]*storetypes.KVStoreKey{
		accountabilitytypes.StoreKey: storetypes.NewKVStoreKey(accountabilitytypes.StoreKey),
		types.StoreKey:               storetypes.NewKVStoreKey(types.StoreKey),
	}
	ctx := testutil.DefaultContextWithKeys(keys, nil, nil)
	storeService := runtime.NewKVStoreService(keys[types.StoreKey])
	accountabilityStore := runtime.NewKVStoreService(keys[accountabilitytypes.StoreKey])

	accountabilityKeeper := accountabilitykeeper.NewKeeper(accountabilityStore, encCfg.Codec, addressCodec)
	k := keeper.NewKeeper(storeService, encCfg.Codec, addressCodec, accountabilityKeeper)

	return &fixture{
		ctx:            ctx,
		keeper:         k,
		addressCodec:   addressCodec,
		accountability: &accountabilityKeeper,
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

func TestRegisterBudget(t *testing.T) {
	f := initFixture(t)

	creator := f.newAddress(t, "budgetCreatorAddr__________")
	amount := sdk.NewInt64Coin("php", 1_000_000).Amount

	registered, err := f.keeper.RegisterBudget(f.ctx, types.Budget{
		Agency:      "DOF",
		FiscalYear:  "2025",
		Title:       "General Appropriations",
		Amount:      amount,
		Currency:    "php",
		MetadataURI: "ipfs://cid",
		CreatedBy:   creator,
	})
	if err != nil {
		t.Fatalf("register budget: %v", err)
	}
	if registered.Id == 0 {
		t.Fatalf("expected budget id to be set")
	}
	if registered.Status != types.BudgetStatusDraft {
		t.Fatalf("unexpected status %s", registered.Status)
	}

	stored, err := f.keeper.GetBudget(f.ctx, registered.Id)
	if err != nil {
		t.Fatalf("get budget: %v", err)
	}
	if stored.Title != "General Appropriations" {
		t.Fatalf("unexpected title %s", stored.Title)
	}
}

func TestUpdateBudgetStatus(t *testing.T) {
	f := initFixture(t)
	creator := f.newAddress(t, "budgetCreatorAddr__________")
	amount := sdk.NewInt64Coin("php", 5_000_000).Amount

	registered, err := f.keeper.RegisterBudget(f.ctx, types.Budget{
		Agency:      "DBM",
		FiscalYear:  "2024",
		Title:       "Budget Execution",
		Amount:      amount,
		Currency:    "php",
		MetadataURI: "ipfs://cid",
		CreatedBy:   creator,
	})
	if err != nil {
		t.Fatalf("register budget: %v", err)
	}

	maintainer := f.newAddress(t, "scoreMaintainerAddr_______")
	auditor := f.newAddress(t, "budgetAuditorAddr__________")
	controller := f.newAddress(t, "budgetControllerAddr_______")

	auditorScore, err := f.accountability.UpsertScorecard(f.ctx, accountabilitytypes.Scorecard{
		Subject:     auditor,
		Metric:      "budget_liability_test",
		Score:       90,
		Weight:      10,
		EvidenceURI: "ipfs://auditor",
		UpdatedBy:   maintainer,
	})
	if err != nil {
		t.Fatalf("upsert auditor score: %v", err)
	}

	controllerScore, err := f.accountability.UpsertScorecard(f.ctx, accountabilitytypes.Scorecard{
		Subject:     controller,
		Metric:      "budget_liability_test",
		Score:       88,
		Weight:      9,
		EvidenceURI: "ipfs://controller",
		UpdatedBy:   maintainer,
	})
	if err != nil {
		t.Fatalf("upsert controller score: %v", err)
	}

	if _, err := f.keeper.RecordBudgetApproval(f.ctx, registered.Id, "ipfs://contract/budget", accountabilitytypes.Approval{
		Signer:       auditor,
		Role:         "Audit",
		ScorecardId:  auditorScore.Id,
		SignatureURI: fmt.Sprintf("ipfs://signatures/%s", auditor),
	}); err != nil {
		t.Fatalf("record auditor approval: %v", err)
	}
	if _, err := f.keeper.RecordBudgetApproval(f.ctx, registered.Id, "", accountabilitytypes.Approval{
		Signer:       controller,
		Role:         "Controller",
		ScorecardId:  controllerScore.Id,
		SignatureURI: fmt.Sprintf("ipfs://signatures/%s", controller),
	}); err != nil {
		t.Fatalf("record controller approval: %v", err)
	}

	updated, err := f.keeper.UpdateBudgetStatus(f.ctx, registered.Id, types.BudgetStatusLegislated)
	if err != nil {
		t.Fatalf("update status: %v", err)
	}
	if updated.Status != types.BudgetStatusLegislated {
		t.Fatalf("unexpected status %s", updated.Status)
	}

	if _, err := f.keeper.UpdateBudgetStatus(f.ctx, registered.Id, types.BudgetStatusDraft); err == nil {
		t.Fatalf("expected invalid transition error")
	}
}

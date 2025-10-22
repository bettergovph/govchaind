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
	budgettypes "govchain/x/budgetregistry/types"
	"govchain/x/procurementledger/keeper"
	"govchain/x/procurementledger/types"
)

type budgetKeeperStub struct {
	budgets map[uint64]budgettypes.Budget
}

func (b budgetKeeperStub) GetBudget(ctx context.Context, id uint64) (budgettypes.Budget, error) {
	if budget, ok := b.budgets[id]; ok {
		return budget, nil
	}
	return budgettypes.Budget{}, budgettypes.ErrBudgetNotFound
}

type fixture struct {
	ctx            context.Context
	keeper         keeper.Keeper
	addressCodec   address.Codec
	budgets        *budgetKeeperStub
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
	budgets := &budgetKeeperStub{budgets: map[uint64]budgettypes.Budget{}}
	k := keeper.NewKeeper(storeService, encCfg.Codec, addressCodec, budgets, accountabilityKeeper)

	return &fixture{
		ctx:            ctx,
		keeper:         k,
		addressCodec:   addressCodec,
		budgets:        budgets,
		accountability: &accountabilityKeeper,
	}
}

func (f *fixture) addBudget(t *testing.T, id uint64, agency string) {
	t.Helper()
	creator, err := f.addressCodec.BytesToString([]byte("budgetCreatorAddr__________"))
	if err != nil {
		t.Fatalf("address conversion: %v", err)
	}
	f.budgets.budgets[id] = budgettypes.Budget{
		Id:         id,
		Agency:     agency,
		FiscalYear: "2024",
		Title:      "Appropriation",
		Amount:     sdk.NewInt64Coin("php", 10_000_000).Amount,
		Currency:   "php",
		Status:     budgettypes.BudgetStatusLegislated,
		CreatedBy:  creator,
	}
}

func (f *fixture) newOfficer(t *testing.T, raw string) string {
	t.Helper()
	addr, err := f.addressCodec.BytesToString([]byte(raw))
	if err != nil {
		t.Fatalf("address conversion: %v", err)
	}
	return addr
}

func TestRegisterProcurement(t *testing.T) {
	f := initFixture(t)
	officer := f.newOfficer(t, "officerAddr________________")

	// No budget registered should error
	if _, err := f.keeper.RegisterProcurement(f.ctx, types.Procurement{
		BudgetId:  1,
		Reference: "REF-2024-001",
		Agency:    "DOTr",
		Title:     "Railway Upgrade",
		Category:  "Infrastructure",
		Amount:    sdk.NewInt64Coin("php", 2_500_000).Amount,
		Currency:  "php",
		Officer:   officer,
	}); err == nil {
		t.Fatalf("expected error without budget")
	}

	f.addBudget(t, 1, "DOTr")

	registered, err := f.keeper.RegisterProcurement(f.ctx, types.Procurement{
		BudgetId:    1,
		Reference:   "REF-2024-001",
		Agency:      "DOTr",
		Title:       "Railway Upgrade",
		Category:    "Infrastructure",
		Amount:      sdk.NewInt64Coin("php", 2_500_000).Amount,
		Currency:    "php",
		Officer:     officer,
		Supplier:    "BuildCo",
		MetadataURI: "ipfs://contract",
	})
	if err != nil {
		t.Fatalf("register procurement: %v", err)
	}
	if registered.Id == 0 {
		t.Fatalf("expected id assigned")
	}
	if registered.Status != types.ProcurementStatusPlanning {
		t.Fatalf("unexpected status %s", registered.Status)
	}
}

func TestUpdateProcurementStatus(t *testing.T) {
	f := initFixture(t)
	officer := f.newOfficer(t, "officerAddr________________")
	f.addBudget(t, 1, "DBM")

	procurement, err := f.keeper.RegisterProcurement(f.ctx, types.Procurement{
		BudgetId:  1,
		Reference: "REF-2024-002",
		Agency:    "DBM",
		Title:     "Software Licensing",
		Category:  "ICT",
		Amount:    sdk.NewInt64Coin("php", 1_000_000).Amount,
		Currency:  "php",
		Officer:   officer,
		Supplier:  "TechCorp",
	})
	if err != nil {
		t.Fatalf("register procurement: %v", err)
	}

	maintainer := f.newOfficer(t, "scoreMaintainerAddr_______")
	oversight := f.newOfficer(t, "oversightAddr______________")
	citizen := f.newOfficer(t, "citizenAddr_______________")

	oversightScore, err := f.accountability.UpsertScorecard(f.ctx, accountabilitytypes.Scorecard{
		Subject:     oversight,
		Metric:      "procurement_liability_test",
		Score:       86,
		Weight:      8,
		EvidenceURI: "ipfs://oversight",
		UpdatedBy:   maintainer,
	})
	if err != nil {
		t.Fatalf("upsert oversight score: %v", err)
	}

	citizenScore, err := f.accountability.UpsertScorecard(f.ctx, accountabilitytypes.Scorecard{
		Subject:     citizen,
		Metric:      "procurement_liability_test",
		Score:       90,
		Weight:      7,
		EvidenceURI: "ipfs://citizen",
		UpdatedBy:   maintainer,
	})
	if err != nil {
		t.Fatalf("upsert citizen score: %v", err)
	}

	if _, err := f.keeper.RecordProcurementApproval(f.ctx, procurement.Id, "ipfs://contract/procurement", accountabilitytypes.Approval{
		Signer:       oversight,
		Role:         "Oversight",
		ScorecardId:  oversightScore.Id,
		SignatureURI: fmt.Sprintf("ipfs://signatures/%s", oversight),
	}); err != nil {
		t.Fatalf("record oversight approval: %v", err)
	}
	if _, err := f.keeper.RecordProcurementApproval(f.ctx, procurement.Id, "", accountabilitytypes.Approval{
		Signer:       citizen,
		Role:         "Citizen",
		ScorecardId:  citizenScore.Id,
		SignatureURI: fmt.Sprintf("ipfs://signatures/%s", citizen),
	}); err != nil {
		t.Fatalf("record citizen approval: %v", err)
	}

	updated, err := f.keeper.UpdateProcurementStatus(f.ctx, procurement.Id, types.ProcurementStatusTendering)
	if err != nil {
		t.Fatalf("update status: %v", err)
	}
	if updated.Status != types.ProcurementStatusTendering {
		t.Fatalf("unexpected status %s", updated.Status)
	}

	if _, err := f.keeper.UpdateProcurementStatus(f.ctx, procurement.Id, types.ProcurementStatusAwarded); err != nil {
		t.Fatalf("update to awarded: %v", err)
	}

	if _, err := f.keeper.UpdateProcurementStatus(f.ctx, procurement.Id, types.ProcurementStatusPlanning); err == nil {
		t.Fatalf("expected invalid transition error")
	}
}

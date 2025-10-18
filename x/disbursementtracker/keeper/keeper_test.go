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

	"govchain/x/disbursementtracker/keeper"
	"govchain/x/disbursementtracker/types"
	procurementtypes "govchain/x/procurementledger/types"
)

type procurementKeeperStub struct {
	procurements map[uint64]procurementtypes.Procurement
}

func (p procurementKeeperStub) GetProcurement(ctx context.Context, id uint64) (procurementtypes.Procurement, error) {
	if procurement, ok := p.procurements[id]; ok {
		return procurement, nil
	}
	return procurementtypes.Procurement{}, procurementtypes.ErrProcurementNotFound
}

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
	procurements *procurementKeeperStub
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig()
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_disbursement")).Ctx

	procurements := &procurementKeeperStub{procurements: map[uint64]procurementtypes.Procurement{}}
	k := keeper.NewKeeper(storeService, encCfg.Codec, addressCodec, procurements)

	return &fixture{
		ctx:          ctx,
		keeper:       k,
		addressCodec: addressCodec,
		procurements: procurements,
	}
}

func (f *fixture) addProcurement(t *testing.T, id uint64, amount int64) {
	t.Helper()
	officer, err := f.addressCodec.BytesToString([]byte("officerAddr________________"))
	if err != nil {
		t.Fatalf("address conversion: %v", err)
	}
	f.procurements.procurements[id] = procurementtypes.Procurement{
		Id:        id,
		BudgetId:  1,
		Reference: "REF-2024-003",
		Agency:    "DICT",
		Title:     "Broadband Expansion",
		Category:  "Infrastructure",
		Amount:    sdk.NewInt64Coin("php", amount).Amount,
		Currency:  "php",
		Status:    procurementtypes.ProcurementStatusAwarded,
		Officer:   officer,
	}
}

func (f *fixture) newPerformer(t *testing.T, raw string) string {
	t.Helper()
	addr, err := f.addressCodec.BytesToString([]byte(raw))
	if err != nil {
		t.Fatalf("address conversion: %v", err)
	}
	return addr
}

func TestRegisterDisbursement(t *testing.T) {
	f := initFixture(t)
	performer := f.newPerformer(t, "performerAddr______________")

	if _, err := f.keeper.RegisterDisbursement(f.ctx, types.Disbursement{
		ProcurementId: 1,
		Amount:        sdk.NewInt64Coin("php", 1_000_000).Amount,
		Currency:      "php",
		PerformedBy:   performer,
	}); err == nil {
		t.Fatalf("expected error when procurement missing")
	}

	f.addProcurement(t, 1, 2_000_000)

	disbursement, err := f.keeper.RegisterDisbursement(f.ctx, types.Disbursement{
		ProcurementId: 1,
		Amount:        sdk.NewInt64Coin("php", 1_000_000).Amount,
		Currency:      "php",
		PerformedBy:   performer,
		EvidenceURI:   "ipfs://evidence",
	})
	if err != nil {
		t.Fatalf("register disbursement: %v", err)
	}
	if disbursement.Status != types.DisbursementStatusScheduled {
		t.Fatalf("unexpected status %s", disbursement.Status)
	}

	total, err := f.keeper.GetDisbursedTotal(f.ctx, 1)
	if err != nil {
		t.Fatalf("get total: %v", err)
	}
	if !total.Equal(sdk.NewInt64Coin("php", 1_000_000).Amount) {
		t.Fatalf("unexpected total %s", total)
	}

	if _, err := f.keeper.RegisterDisbursement(f.ctx, types.Disbursement{
		ProcurementId: 1,
		Amount:        sdk.NewInt64Coin("php", 1_100_000).Amount,
		Currency:      "php",
		PerformedBy:   performer,
	}); err == nil {
		t.Fatalf("expected error exceeding procurement amount")
	}
}

func TestUpdateDisbursementStatus(t *testing.T) {
	f := initFixture(t)
	performer := f.newPerformer(t, "performerAddr______________")
	f.addProcurement(t, 1, 1_000_000)

	disbursement, err := f.keeper.RegisterDisbursement(f.ctx, types.Disbursement{
		ProcurementId: 1,
		Amount:        sdk.NewInt64Coin("php", 500_000).Amount,
		Currency:      "php",
		PerformedBy:   performer,
	})
	if err != nil {
		t.Fatalf("register disbursement: %v", err)
	}

	updated, err := f.keeper.UpdateDisbursementStatus(f.ctx, disbursement.Id, types.DisbursementStatusReleased)
	if err != nil {
		t.Fatalf("update status: %v", err)
	}
	if updated.Status != types.DisbursementStatusReleased {
		t.Fatalf("unexpected status %s", updated.Status)
	}

	if _, err := f.keeper.UpdateDisbursementStatus(f.ctx, disbursement.Id, types.DisbursementStatusScheduled); err == nil {
		t.Fatalf("expected invalid transition error")
	}
}

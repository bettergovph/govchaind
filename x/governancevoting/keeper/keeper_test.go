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

	"govchain/x/governancevoting/keeper"
	"govchain/x/governancevoting/types"
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
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_governance")).Ctx

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

func TestRegisterDelegation(t *testing.T) {
	f := initFixture(t)
	delegator := f.newAddress(t, "delegatorAddr______________")
	delegatee := f.newAddress(t, "delegateeAddr______________")

	delegation, err := f.keeper.RegisterDelegation(f.ctx, types.Delegation{
		Delegator: delegator,
		Delegatee: delegatee,
		Scope:     "budget-oversight",
	})
	if err != nil {
		t.Fatalf("register delegation: %v", err)
	}
	if delegation.Id == 0 {
		t.Fatalf("expected id assigned")
	}
	if !delegation.Active {
		t.Fatalf("expected delegation active")
	}

	fetched, err := f.keeper.GetDelegationByScope(f.ctx, delegator, "budget-oversight")
	if err != nil {
		t.Fatalf("get delegation: %v", err)
	}
	if fetched.Id != delegation.Id {
		t.Fatalf("expected same delegation id")
	}

	updated, err := f.keeper.RegisterDelegation(f.ctx, types.Delegation{
		Delegator:   delegator,
		Delegatee:   delegatee,
		Scope:       "budget-oversight",
		ExpiresAt:   100,
		MetadataURI: "ipfs://delegation",
	})
	if err != nil {
		t.Fatalf("update delegation: %v", err)
	}
	if updated.Id != delegation.Id {
		t.Fatalf("expected update to reuse identifier")
	}
	if updated.ExpiresAt != 100 {
		t.Fatalf("expected expires_at to update")
	}
}

func TestDeactivateDelegation(t *testing.T) {
	f := initFixture(t)
	delegator := f.newAddress(t, "delegatorAddr______________")
	delegatee := f.newAddress(t, "delegateeAddr______________")

	delegation, err := f.keeper.RegisterDelegation(f.ctx, types.Delegation{
		Delegator: delegator,
		Delegatee: delegatee,
		Scope:     "procurement",
	})
	if err != nil {
		t.Fatalf("register delegation: %v", err)
	}

	updated, err := f.keeper.DeactivateDelegation(f.ctx, delegation.Id)
	if err != nil {
		t.Fatalf("deactivate delegation: %v", err)
	}
	if updated.Active {
		t.Fatalf("expected delegation inactive")
	}
}

package accountabilityscores

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"govchain/x/accountabilityscores/keeper"
	"govchain/x/accountabilityscores/types"
)

var _ module.AppModule = AppModule{}
var _ module.AppModuleBasic = AppModule{}
var _ appmodule.AppModule = AppModule{}

// AppModule wires the accountability scoring logic into the Cosmos application.
type AppModule struct {
	cdc    codec.Codec
	keeper keeper.Keeper
}

// NewAppModule returns the app module implementation.
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{cdc: cdc, keeper: keeper}
}

// IsAppModule marks the struct as implementing the app module interface.
func (AppModule) IsAppModule() {}

// Name returns the module's name for routing and registration.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers legacy amino types (noop).
func (AppModule) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

// RegisterInterfaces registers module interface types.
func (AppModule) RegisterInterfaces(_ codectypes.InterfaceRegistry) {}

// RegisterGRPCGatewayRoutes registers gRPC gateway routes (none for this module).
func (AppModule) RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux) {}

// RegisterServices registers gRPC services (none for this module).
func (AppModule) RegisterServices(_ grpc.ServiceRegistrar) error { return nil }

// DefaultGenesis provides default genesis state.
func (am AppModule) DefaultGenesis(codec.JSONCodec) json.RawMessage {
	bz, err := json.Marshal(types.DefaultGenesis())
	if err != nil {
		panic(fmt.Errorf("failed to marshal accountability default genesis: %w", err))
	}
	return bz
}

// ValidateGenesis validates the module genesis state.
func (am AppModule) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var state types.GenesisState
	if err := json.Unmarshal(bz, &state); err != nil {
		return fmt.Errorf("failed to unmarshal accountability genesis: %w", err)
	}
	return state.Validate()
}

// InitGenesis initializes module state from genesis.
func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, bz json.RawMessage) {
	var state types.GenesisState
	if err := json.Unmarshal(bz, &state); err != nil {
		panic(fmt.Errorf("failed to unmarshal accountability genesis: %w", err))
	}
	if err := am.keeper.InitGenesis(ctx, state); err != nil {
		panic(fmt.Errorf("failed to initialize accountability genesis: %w", err))
	}
}

// ExportGenesis exports the module state to genesis.
func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	state, err := am.keeper.ExportGenesis(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to export accountability genesis: %w", err))
	}
	bz, err := json.Marshal(state)
	if err != nil {
		panic(fmt.Errorf("failed to marshal accountability genesis: %w", err))
	}
	return bz
}

// ConsensusVersion returns the consensus version.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock is a no-op for this module.
func (AppModule) BeginBlock(context.Context) error { return nil }

// EndBlock is a no-op for this module.
func (AppModule) EndBlock(context.Context) error { return nil }

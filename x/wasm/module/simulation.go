package wasm

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

// GenerateGenesisState is a no-op placeholder for simulation support.
func (AppModule) GenerateGenesisState(_ *module.SimulationState) {}

// ProposalContents returns the module governance proposals used for simulations.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams returns randomized module parameters for simulations.
func (AppModule) RandomizedParams(_ *rand.Rand) []simtypes.LegacyParamChange {
	return nil
}

// RegisterStoreDecoder registers a decoder for simulation store queries.
func (AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// RegisterLegacyAminoCodec registers any concrete types on the Amino codec.
// The wasm module currently only stores metadata and does not expose legacy interfaces.
func RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

// RegisterInterfaces registers interface implementations for the module.
func RegisterInterfaces(_ codectypes.InterfaceRegistry) {}

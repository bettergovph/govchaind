package types

import (
	"fmt"
)

// Params defines the configurable parameters for the wasm module.
type Params struct {
	EnableSmartContracts bool `json:"enable_smart_contracts" yaml:"enable_smart_contracts"`
}

// DefaultParams returns the default module parameters.
func DefaultParams() Params {
	return Params{EnableSmartContracts: true}
}

// Validate performs basic validation of parameter values.
func (p Params) Validate() error {
	if !p.EnableSmartContracts {
		return fmt.Errorf("smart contracts must remain enabled for wasm module preparation")
	}
	return nil
}

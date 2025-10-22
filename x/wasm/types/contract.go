package types

import (
	"fmt"
)

// Contract stores high-level metadata for a registered smart contract.
type Contract struct {
	Id             uint64 `json:"id" yaml:"id"`
	Address        string `json:"address" yaml:"address"`
	Creator        string `json:"creator" yaml:"creator"`
	Admin          string `json:"admin" yaml:"admin"`
	Label          string `json:"label" yaml:"label"`
	CodeIdentifier string `json:"code_identifier" yaml:"code_identifier"`
}

// ValidateBasic performs stateless validation on the contract metadata.
func (c Contract) ValidateBasic() error {
	if c.Id == 0 {
		return fmt.Errorf("contract identifier must be non-zero")
	}
	if c.Address == "" {
		return fmt.Errorf("contract address is required")
	}
	if c.Creator == "" {
		return fmt.Errorf("contract creator is required")
	}
	return nil
}

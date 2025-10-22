package types

import (
	"fmt"
	"strings"
)

// Delegation captures delegated accountability relationships used for governance voting.
type Delegation struct {
	Id          uint64 `json:"id" yaml:"id"`
	Delegator   string `json:"delegator" yaml:"delegator"`
	Delegatee   string `json:"delegatee" yaml:"delegatee"`
	Scope       string `json:"scope" yaml:"scope"`
	ExpiresAt   uint64 `json:"expires_at" yaml:"expires_at"`
	Active      bool   `json:"active" yaml:"active"`
	MetadataURI string `json:"metadata_uri" yaml:"metadata_uri"`
}

// ValidateBasic performs stateless validation.
func (d Delegation) ValidateBasic() error {
	if d.Id == 0 {
		return fmt.Errorf("delegation id must be set")
	}
	if strings.TrimSpace(d.Delegator) == "" {
		return fmt.Errorf("delegator address is required")
	}
	if strings.TrimSpace(d.Delegatee) == "" {
		return fmt.Errorf("delegatee address is required")
	}
	if strings.TrimSpace(d.Scope) == "" {
		return fmt.Errorf("scope is required")
	}
	return nil
}

// IndexKey returns a unique composite key for the delegation.
func (d Delegation) IndexKey() string {
	return strings.ToLower(strings.TrimSpace(d.Delegator)) + "|" + strings.ToLower(strings.TrimSpace(d.Scope))
}

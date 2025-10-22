package types

import "strings"

// Approval captures a signed endorsement for a budget, procurement, or
// disbursement record. Each approval must be traceable to an accountability
// scorecard and link to a signed liability contract artifact.
type Approval struct {
	Signer       string `json:"signer" yaml:"signer"`
	Role         string `json:"role" yaml:"role"`
	ScorecardId  uint64 `json:"scorecard_id" yaml:"scorecard_id"`
	SignatureURI string `json:"signature_uri" yaml:"signature_uri"`
}

// ValidateBasic performs stateless checks on the approval payload.
func (a Approval) ValidateBasic() error {
	if strings.TrimSpace(a.Signer) == "" {
		return ErrInvalidApproval.Wrap("signer address is required")
	}
	if strings.TrimSpace(a.Role) == "" {
		return ErrInvalidApproval.Wrap("role is required")
	}
	if a.ScorecardId == 0 {
		return ErrInvalidApproval.Wrap("scorecard id must be set")
	}
	if strings.TrimSpace(a.SignatureURI) == "" {
		return ErrInvalidApproval.Wrap("signature_uri is required")
	}
	return nil
}

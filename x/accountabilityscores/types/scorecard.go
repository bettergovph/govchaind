package types

import (
	"fmt"
	"strings"
)

// Scorecard captures a transparency or accountability metric for a subject.
type Scorecard struct {
	Id          uint64 `json:"id" yaml:"id"`
	Subject     string `json:"subject" yaml:"subject"`
	Metric      string `json:"metric" yaml:"metric"`
	Score       uint32 `json:"score" yaml:"score"`
	Weight      uint32 `json:"weight" yaml:"weight"`
	EvidenceURI string `json:"evidence_uri" yaml:"evidence_uri"`
	UpdatedBy   string `json:"updated_by" yaml:"updated_by"`
}

// ValidateBasic performs stateless validation on the scorecard.
func (s Scorecard) ValidateBasic() error {
	if s.Id == 0 {
		return fmt.Errorf("scorecard id must be set")
	}
	if strings.TrimSpace(s.Subject) == "" {
		return fmt.Errorf("subject is required")
	}
	if strings.TrimSpace(s.Metric) == "" {
		return fmt.Errorf("metric is required")
	}
	if s.Score > 100 {
		return fmt.Errorf("score must be between 0 and 100")
	}
	if s.Weight == 0 {
		return fmt.Errorf("weight must be positive")
	}
	if strings.TrimSpace(s.UpdatedBy) == "" {
		return fmt.Errorf("updated_by is required")
	}
	return nil
}

// IndexKey returns the canonical index key for the subject/metric combination.
func (s Scorecard) IndexKey() string {
	return strings.ToLower(strings.TrimSpace(s.Subject)) + "|" + strings.ToLower(strings.TrimSpace(s.Metric))
}

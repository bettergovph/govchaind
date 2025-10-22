package types

import "fmt"

// GenesisState defines the initial accountability score state.
type GenesisState struct {
	Scorecards     []Scorecard `json:"scorecards" yaml:"scorecards"`
	ScorecardCount uint64      `json:"scorecard_count" yaml:"scorecard_count"`
}

// DefaultGenesis returns an empty state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Scorecards:     []Scorecard{},
		ScorecardCount: 0,
	}
}

// Validate ensures the genesis data is consistent.
func (gs GenesisState) Validate() error {
	seen := make(map[uint64]struct{}, len(gs.Scorecards))
	index := make(map[string]struct{}, len(gs.Scorecards))
	for _, score := range gs.Scorecards {
		if _, exists := seen[score.Id]; exists {
			return fmt.Errorf("duplicate scorecard id %d in genesis", score.Id)
		}
		seen[score.Id] = struct{}{}
		if err := score.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid scorecard %d: %w", score.Id, err)
		}
		key := score.IndexKey()
		if _, exists := index[key]; exists {
			return fmt.Errorf("duplicate scorecard index %s in genesis", key)
		}
		index[key] = struct{}{}
	}
	if uint64(len(gs.Scorecards)) > gs.ScorecardCount {
		return fmt.Errorf("scorecard count %d smaller than scorecards supplied %d", gs.ScorecardCount, len(gs.Scorecards))
	}
	return nil
}

package types

import (
	"context"

	accountabilitytypes "govchain/x/accountabilityscores/types"
)

// AccountabilityKeeper defines the subset of the accountability scores keeper
// required by the budget registry.
type AccountabilityKeeper interface {
	GetScorecard(context.Context, uint64) (accountabilitytypes.Scorecard, error)
}

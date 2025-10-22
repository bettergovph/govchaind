package types

import (
	"context"

	accountabilitytypes "govchain/x/accountabilityscores/types"
	procurementtypes "govchain/x/procurementledger/types"
)

// ProcurementKeeper defines the procurement interface consumed by the disbursement module.
type ProcurementKeeper interface {
	GetProcurement(context.Context, uint64) (procurementtypes.Procurement, error)
}

// AccountabilityKeeper exposes scorecard lookup functionality.
type AccountabilityKeeper interface {
	GetScorecard(context.Context, uint64) (accountabilitytypes.Scorecard, error)
}

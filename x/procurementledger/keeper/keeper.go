package keeper

import (
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	accountabilitytypes "govchain/x/accountabilityscores/types"
	"govchain/x/procurementledger/types"
)

const (
	minProcurementApprovals     = 2
	minProcurementApprovalScore = 80
)

// Keeper maintains procurement and contract lifecycle state.
type Keeper struct {
	storeService   corestore.KVStoreService
	cdc            codec.Codec
	addressCodec   address.Codec
	budgets        types.BudgetKeeper
	accountability types.AccountabilityKeeper

	Schema         collections.Schema
	ProcurementSeq collections.Sequence
	Procurements   collections.Map[uint64, types.Procurement]
}

// NewKeeper constructs a keeper instance with its dependencies satisfied.
func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	budgetKeeper types.BudgetKeeper,
	accountability types.AccountabilityKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:   storeService,
		cdc:            cdc,
		addressCodec:   addressCodec,
		budgets:        budgetKeeper,
		accountability: accountability,
		Procurements:   collections.NewMap(sb, types.ProcurementKey, "procurement", collections.Uint64Key, collections.NewJSONValueCodec[types.Procurement]()),
		ProcurementSeq: collections.NewSequence(sb, types.ProcurementCountKey, "procurement_seq"),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(fmt.Errorf("failed to build procurement ledger schema: %w", err))
	}
	k.Schema = schema
	return k
}

// RegisterProcurement records a new procurement entry after validating the referenced budget.
func (k Keeper) RegisterProcurement(ctx context.Context, procurement types.Procurement) (types.Procurement, error) {
	if _, err := k.addressCodec.StringToBytes(procurement.Officer); err != nil {
		return types.Procurement{}, fmt.Errorf("%w: %v", types.ErrProcurementOfficerAddr, err)
	}
	if procurement.Status == "" {
		procurement.Status = types.DefaultProcurementStatus()
	}
	if _, err := k.budgets.GetBudget(ctx, procurement.BudgetId); err != nil {
		return types.Procurement{}, fmt.Errorf("%w: %v", types.ErrBudgetMissingForProcurement, err)
	}

	id, err := k.ProcurementSeq.Next(ctx)
	if err != nil {
		return types.Procurement{}, err
	}
	procurement.Id = id + 1

	if err := procurement.ValidateBasic(); err != nil {
		return types.Procurement{}, err
	}

	if err := k.Procurements.Set(ctx, procurement.Id, procurement); err != nil {
		return types.Procurement{}, err
	}
	return procurement, nil
}

// UpdateProcurementStatus adjusts the lifecycle stage of a procurement.
func (k Keeper) UpdateProcurementStatus(ctx context.Context, id uint64, next types.ProcurementStatus) (types.Procurement, error) {
	procurement, err := k.Procurements.Get(ctx, id)
	if err != nil {
		return types.Procurement{}, fmt.Errorf("%w: %v", types.ErrProcurementNotFound, err)
	}

	if procurement.Status == next {
		return procurement, nil
	}
	if !procurement.Status.CanTransitionTo(next) {
		return types.Procurement{}, fmt.Errorf("%w: %s -> %s", types.ErrInvalidProcurementTransition, procurement.Status, next)
	}

	if requiresProcurementApproval(next) {
		if err := k.ensureProcurementApprovals(ctx, procurement); err != nil {
			return types.Procurement{}, err
		}
	}

	procurement.Status = next
	if err := k.Procurements.Set(ctx, procurement.Id, procurement); err != nil {
		return types.Procurement{}, err
	}
	return procurement, nil
}

// GetProcurement returns the procurement metadata for the supplied identifier.
func (k Keeper) GetProcurement(ctx context.Context, id uint64) (types.Procurement, error) {
	procurement, err := k.Procurements.Get(ctx, id)
	if err != nil {
		return types.Procurement{}, fmt.Errorf("%w: %v", types.ErrProcurementNotFound, err)
	}
	return procurement, nil
}

// WalkProcurements iterates all procurements.
func (k Keeper) WalkProcurements(ctx context.Context, cb func(types.Procurement) (bool, error)) error {
	return k.Procurements.Walk(ctx, nil, func(_ uint64, procurement types.Procurement) (bool, error) {
		return cb(procurement)
	})
}

// RecordProcurementApproval adds or updates a procurement approval tied to the
// shared liability contract.
func (k Keeper) RecordProcurementApproval(
	ctx context.Context,
	id uint64,
	contractURI string,
	approval accountabilitytypes.Approval,
) (types.Procurement, error) {
	procurement, err := k.Procurements.Get(ctx, id)
	if err != nil {
		return types.Procurement{}, fmt.Errorf("%w: %v", types.ErrProcurementNotFound, err)
	}

	if strings.TrimSpace(contractURI) != "" {
		procurement.LiabilityContractURI = contractURI
	}
	if strings.TrimSpace(procurement.LiabilityContractURI) == "" {
		return types.Procurement{}, types.ErrProcurementLiabilityMissing
	}

	seenRoles := map[string]struct{}{}
	seenSigners := map[string]struct{}{}
	roleKey := normalizeKey(approval.Role)
	signerKey := normalizeKey(approval.Signer)
	updated := make([]accountabilitytypes.Approval, 0, len(procurement.Approvals)+1)
	for _, existing := range procurement.Approvals {
		existingRoleKey := normalizeKey(existing.Role)
		existingSignerKey := normalizeKey(existing.Signer)
		if existingRoleKey == roleKey || existingSignerKey == signerKey {
			continue
		}
		seenRoles[existingRoleKey] = struct{}{}
		seenSigners[existingSignerKey] = struct{}{}
		updated = append(updated, existing)
	}

	if err := k.validateProcurementApproval(ctx, approval, seenRoles, seenSigners); err != nil {
		return types.Procurement{}, err
	}

	updated = append(updated, approval)
	procurement.Approvals = updated

	if err := k.Procurements.Set(ctx, procurement.Id, procurement); err != nil {
		return types.Procurement{}, err
	}
	return procurement, nil
}

func (k Keeper) ensureProcurementApprovals(ctx context.Context, procurement types.Procurement) error {
	if strings.TrimSpace(procurement.LiabilityContractURI) == "" {
		return types.ErrProcurementLiabilityMissing
	}
	if len(procurement.Approvals) < minProcurementApprovals {
		return types.ErrProcurementApprovalThreshold.Wrapf("requires at least %d approvals, found %d", minProcurementApprovals, len(procurement.Approvals))
	}

	seenRoles := map[string]struct{}{}
	seenSigners := map[string]struct{}{}
	for _, approval := range procurement.Approvals {
		if err := k.validateProcurementApproval(ctx, approval, seenRoles, seenSigners); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) validateProcurementApproval(
	ctx context.Context,
	approval accountabilitytypes.Approval,
	seenRoles map[string]struct{},
	seenSigners map[string]struct{},
) error {
	if err := approval.ValidateBasic(); err != nil {
		return err
	}
	if _, err := k.addressCodec.StringToBytes(approval.Signer); err != nil {
		return types.ErrProcurementApprovalSigner.Wrap(err.Error())
	}

	roleKey := normalizeKey(approval.Role)
	if _, exists := seenRoles[roleKey]; exists {
		return types.ErrProcurementApprovalDuplicate.Wrap("role already approved")
	}
	seenRoles[roleKey] = struct{}{}

	signerKey := normalizeKey(approval.Signer)
	if _, exists := seenSigners[signerKey]; exists {
		return types.ErrProcurementApprovalDuplicate.Wrap("signer already approved")
	}
	seenSigners[signerKey] = struct{}{}

	scorecard, err := k.accountability.GetScorecard(ctx, approval.ScorecardId)
	if err != nil {
		return types.ErrProcurementApprovalScore.Wrap(err.Error())
	}
	if normalizeKey(scorecard.Subject) != signerKey {
		return types.ErrProcurementApprovalScore.Wrap("scorecard subject does not match signer")
	}
	if scorecard.Score < minProcurementApprovalScore {
		return types.ErrProcurementApprovalScore.Wrapf("score %d below minimum %d", scorecard.Score, minProcurementApprovalScore)
	}
	return nil
}

func requiresProcurementApproval(status types.ProcurementStatus) bool {
	switch status {
	case types.ProcurementStatusAwarded, types.ProcurementStatusExecuting, types.ProcurementStatusCompleted:
		return true
	default:
		return false
	}
}

func normalizeKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

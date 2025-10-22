package keeper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	accountabilitytypes "govchain/x/accountabilityscores/types"
	"govchain/x/budgetregistry/types"
)

const (
	minBudgetApprovals     = 2
	minBudgetApprovalScore = 80
)

// Keeper wires together module dependencies to manage budget registry state.
type Keeper struct {
	storeService   corestore.KVStoreService
	cdc            codec.Codec
	addressCodec   address.Codec
	accountability types.AccountabilityKeeper

	Schema    collections.Schema
	BudgetSeq collections.Sequence
	Budgets   collections.Map[uint64, types.Budget]
}

// NewKeeper creates a new budget registry keeper with the required dependencies.
func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	accountability types.AccountabilityKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:   storeService,
		cdc:            cdc,
		addressCodec:   addressCodec,
		accountability: accountability,
		Budgets:        collections.NewMap(sb, types.BudgetKey, "budget", collections.Uint64Key, collections.NewJSONValueCodec[types.Budget]()),
		BudgetSeq:      collections.NewSequence(sb, types.BudgetCountKey, "budget_seq"),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(fmt.Errorf("failed to build budget registry schema: %w", err))
	}
	k.Schema = schema

	return k
}

// RegisterBudget stores a new budget entry and returns the persisted value.
func (k Keeper) RegisterBudget(ctx context.Context, budget types.Budget) (types.Budget, error) {
	if _, err := k.addressCodec.StringToBytes(budget.CreatedBy); err != nil {
		return types.Budget{}, fmt.Errorf("%w: %v", types.ErrBudgetCreatorInvalidAddr, err)
	}

	if budget.Status == "" {
		budget.Status = types.DefaultBudgetStatus()
	}

	id, err := k.BudgetSeq.Next(ctx)
	if err != nil {
		return types.Budget{}, err
	}
	budget.Id = id + 1

	if err := budget.ValidateBasic(); err != nil {
		return types.Budget{}, err
	}

	if err := k.Budgets.Set(ctx, budget.Id, budget); err != nil {
		return types.Budget{}, err
	}
	return budget, nil
}

// UpdateBudgetStatus changes the lifecycle status of a budget entry.
func (k Keeper) UpdateBudgetStatus(ctx context.Context, id uint64, next types.BudgetStatus) (types.Budget, error) {
	budget, err := k.Budgets.Get(ctx, id)
	if err != nil {
		return types.Budget{}, fmt.Errorf("%w: %v", types.ErrBudgetNotFound, err)
	}

	if budget.Status == next {
		return budget, nil
	}

	if !budget.Status.CanTransitionTo(next) {
		return types.Budget{}, fmt.Errorf("%w: %s -> %s", types.ErrInvalidBudgetTransition, budget.Status, next)
	}

	if next != types.BudgetStatusDraft {
		if err := k.ensureBudgetApprovals(ctx, budget); err != nil {
			return types.Budget{}, err
		}
	}

	budget.Status = next
	if err := k.Budgets.Set(ctx, budget.Id, budget); err != nil {
		return types.Budget{}, err
	}
	return budget, nil
}

// GetBudget returns the stored budget for an identifier.
func (k Keeper) GetBudget(ctx context.Context, id uint64) (types.Budget, error) {
	budget, err := k.Budgets.Get(ctx, id)
	if err != nil {
		return types.Budget{}, fmt.Errorf("%w: %v", types.ErrBudgetNotFound, err)
	}
	return budget, nil
}

// HasBudget reports whether a budget id exists.
func (k Keeper) HasBudget(ctx context.Context, id uint64) (bool, error) {
	_, err := k.Budgets.Get(ctx, id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// WalkBudgets iterates all stored budgets.
func (k Keeper) WalkBudgets(ctx context.Context, cb func(types.Budget) (bool, error)) error {
	return k.Budgets.Walk(ctx, nil, func(_ uint64, budget types.Budget) (bool, error) {
		return cb(budget)
	})
}

// RecordBudgetApproval stores or updates an approval for the budget and ensures
// the signer satisfies accountability and liability requirements.
func (k Keeper) RecordBudgetApproval(
	ctx context.Context,
	id uint64,
	contractURI string,
	approval accountabilitytypes.Approval,
) (types.Budget, error) {
	budget, err := k.Budgets.Get(ctx, id)
	if err != nil {
		return types.Budget{}, fmt.Errorf("%w: %v", types.ErrBudgetNotFound, err)
	}

	if strings.TrimSpace(contractURI) != "" {
		budget.LiabilityContractURI = contractURI
	}
	if strings.TrimSpace(budget.LiabilityContractURI) == "" {
		return types.Budget{}, types.ErrBudgetLiabilityMissing
	}

	seenRoles := map[string]struct{}{}
	seenSigners := map[string]struct{}{}
	// Seed the maps with the approvals that will be preserved so we can
	// maintain uniqueness while allowing updates.
	roleKey := normalizeKey(approval.Role)
	signerKey := normalizeKey(approval.Signer)
	updated := make([]accountabilitytypes.Approval, 0, len(budget.Approvals)+1)
	for _, existing := range budget.Approvals {
		existingRoleKey := normalizeKey(existing.Role)
		existingSignerKey := normalizeKey(existing.Signer)
		if existingRoleKey == roleKey || existingSignerKey == signerKey {
			continue
		}
		seenRoles[existingRoleKey] = struct{}{}
		seenSigners[existingSignerKey] = struct{}{}
		updated = append(updated, existing)
	}

	if err := k.validateBudgetApproval(ctx, approval, seenRoles, seenSigners); err != nil {
		return types.Budget{}, err
	}

	updated = append(updated, approval)
	budget.Approvals = updated

	if err := k.Budgets.Set(ctx, budget.Id, budget); err != nil {
		return types.Budget{}, err
	}
	return budget, nil
}

func (k Keeper) ensureBudgetApprovals(ctx context.Context, budget types.Budget) error {
	if strings.TrimSpace(budget.LiabilityContractURI) == "" {
		return types.ErrBudgetLiabilityMissing
	}
	if len(budget.Approvals) < minBudgetApprovals {
		return types.ErrBudgetApprovalThreshold.Wrapf("requires at least %d approvals, found %d", minBudgetApprovals, len(budget.Approvals))
	}

	seenRoles := map[string]struct{}{}
	seenSigners := map[string]struct{}{}
	for _, approval := range budget.Approvals {
		if err := k.validateBudgetApproval(ctx, approval, seenRoles, seenSigners); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) validateBudgetApproval(
	ctx context.Context,
	approval accountabilitytypes.Approval,
	seenRoles map[string]struct{},
	seenSigners map[string]struct{},
) error {
	if err := approval.ValidateBasic(); err != nil {
		return err
	}
	if _, err := k.addressCodec.StringToBytes(approval.Signer); err != nil {
		return types.ErrBudgetApprovalSignerAddr.Wrap(err.Error())
	}

	roleKey := normalizeKey(approval.Role)
	if _, exists := seenRoles[roleKey]; exists {
		return types.ErrBudgetApprovalDuplicate.Wrap("role already approved")
	}
	seenRoles[roleKey] = struct{}{}

	signerKey := normalizeKey(approval.Signer)
	if _, exists := seenSigners[signerKey]; exists {
		return types.ErrBudgetApprovalDuplicate.Wrap("signer already approved")
	}
	seenSigners[signerKey] = struct{}{}

	scorecard, err := k.accountability.GetScorecard(ctx, approval.ScorecardId)
	if err != nil {
		return types.ErrBudgetApprovalScore.Wrap(err.Error())
	}
	if normalizeKey(scorecard.Subject) != signerKey {
		return types.ErrBudgetApprovalScore.Wrap("scorecard subject does not match signer")
	}
	if scorecard.Score < minBudgetApprovalScore {
		return types.ErrBudgetApprovalScore.Wrapf("score %d below minimum %d", scorecard.Score, minBudgetApprovalScore)
	}
	return nil
}

func normalizeKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

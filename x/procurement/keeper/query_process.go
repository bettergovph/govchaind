package keeper

import (
	"context"
	"errors"

	"govchain/x/procurement/types"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) GetProcess(ctx context.Context, req *types.QueryGetProcessRequest) (*types.QueryGetProcessResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	process, err := q.k.Process.Get(ctx, req.Ocid)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetProcessResponse{Process: process}, nil
}

func (q queryServer) ListProcesses(ctx context.Context, req *types.QueryListProcessesRequest) (*types.QueryListProcessesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	processes, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Process,
		req.Pagination,
		func(_ string, value types.ContractingProcess) (types.ContractingProcess, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListProcessesResponse{Processes: processes, Pagination: pageRes}, nil
}

func (q queryServer) GetRelease(ctx context.Context, req *types.QueryGetReleaseRequest) (*types.QueryGetReleaseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	releaseKey := req.Ocid + "/" + req.ReleaseId
	release, err := q.k.Release.Get(ctx, releaseKey)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetReleaseResponse{Release: release}, nil
}

func (q queryServer) ListReleases(ctx context.Context, req *types.QueryListReleasesRequest) (*types.QueryListReleasesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Get the process to find release IDs
	process, err := q.k.Process.Get(ctx, req.Ocid)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	var releases []types.Release
	for _, rid := range process.ReleaseIds {
		releaseKey := req.Ocid + "/" + rid
		release, err := q.k.Release.Get(ctx, releaseKey)
		if err != nil {
			continue
		}
		releases = append(releases, release)
	}

	return &types.QueryListReleasesResponse{Releases: releases}, nil
}

func (q queryServer) GetEntity(ctx context.Context, req *types.QueryGetEntityRequest) (*types.QueryGetEntityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	entity, err := q.k.Entity.Get(ctx, req.EntityId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetEntityResponse{Entity: entity}, nil
}

func (q queryServer) ListEntities(ctx context.Context, req *types.QueryListEntitiesRequest) (*types.QueryListEntitiesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	entities, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Entity,
		req.Pagination,
		func(_ string, value types.EntityRegistry) (types.EntityRegistry, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListEntitiesResponse{Entities: entities, Pagination: pageRes}, nil
}

func (q queryServer) SearchTenders(ctx context.Context, req *types.QuerySearchTendersRequest) (*types.QuerySearchTendersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Walk all processes and filter by tender criteria
	var matched []types.ContractingProcess
	err := q.k.Process.Walk(ctx, nil, func(key string, process types.ContractingProcess) (bool, error) {
		if process.CompiledRelease == nil || process.CompiledRelease.Tender == nil {
			return false, nil
		}
		tender := process.CompiledRelease.Tender

		if req.Status != "" && tender.Status != req.Status {
			return false, nil
		}
		if req.ProcurementMethod != "" && tender.ProcurementMethod != req.ProcurementMethod {
			return false, nil
		}
		if req.MainProcurementCategory != "" && tender.MainProcurementCategory != req.MainProcurementCategory {
			return false, nil
		}
		if req.ProcuringEntityId != "" && (tender.ProcuringEntity == nil || tender.ProcuringEntity.Id != req.ProcuringEntityId) {
			return false, nil
		}

		matched = append(matched, process)
		return false, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySearchTendersResponse{Processes: matched}, nil
}

func (q queryServer) GetOCDSPackage(ctx context.Context, req *types.QueryGetOCDSPackageRequest) (*types.QueryGetOCDSPackageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	process, err := q.k.Process.Get(ctx, req.Ocid)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	// Collect all releases for this process
	var releases []*types.Release
	for _, rid := range process.ReleaseIds {
		releaseKey := req.Ocid + "/" + rid
		release, err := q.k.Release.Get(ctx, releaseKey)
		if err != nil {
			continue
		}
		r := release
		releases = append(releases, &r)
	}

	pkg := &types.ReleasePackage{
		Version:  "1.1",
		Releases: releases,
		Publisher: &types.Publisher{
			Name:   "OpenGovChain",
			Scheme: "cosmos",
		},
		License: "https://opendatacommons.org/licenses/pddl/1.0/",
	}

	return &types.QueryGetOCDSPackageResponse{Package: pkg}, nil
}

func (q queryServer) GetRecordPackage(ctx context.Context, req *types.QueryGetRecordPackageRequest) (*types.QueryGetRecordPackageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	process, err := q.k.Process.Get(ctx, req.Ocid)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	record := &types.Record{
		Ocid:            process.Ocid,
		CompiledRelease: process.CompiledRelease,
	}

	return &types.QueryGetRecordPackageResponse{
		Records: []*types.Record{record},
		Publisher: &types.Publisher{
			Name:   "OpenGovChain",
			Scheme: "cosmos",
		},
		Version: "1.1",
		License: "https://opendatacommons.org/licenses/pddl/1.0/",
	}, nil
}

package keeper

import (
	"context"
	"errors"
	"strings"

	"govchain/x/infrastructure/types"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) GetProject(ctx context.Context, req *types.QueryGetProjectRequest) (*types.QueryGetProjectResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	project, err := q.k.Project.Get(ctx, req.ProjectId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetProjectResponse{Project: project}, nil
}

func (q queryServer) ListProjects(ctx context.Context, req *types.QueryListProjectsRequest) (*types.QueryListProjectsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	projects, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Project,
		req.Pagination,
		func(_ string, value types.InfraProject) (types.InfraProject, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListProjectsResponse{Projects: projects, Pagination: pageRes}, nil
}

func (q queryServer) ProjectsBySector(ctx context.Context, req *types.QueryProjectsBySectorRequest) (*types.QueryProjectsBySectorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var matched []types.InfraProject
	err := q.k.Project.Walk(ctx, nil, func(key string, project types.InfraProject) (bool, error) {
		for _, s := range project.Sector {
			if strings.EqualFold(s, req.Sector) {
				matched = append(matched, project)
				break
			}
		}
		return false, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProjectsBySectorResponse{Projects: matched}, nil
}

func (q queryServer) ProjectsByStatus(ctx context.Context, req *types.QueryProjectsByStatusRequest) (*types.QueryProjectsByStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var matched []types.InfraProject
	err := q.k.Project.Walk(ctx, nil, func(key string, project types.InfraProject) (bool, error) {
		if strings.EqualFold(project.Status, req.Status) {
			matched = append(matched, project)
		}
		return false, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProjectsByStatusResponse{Projects: matched}, nil
}

func (q queryServer) ProjectsByRegion(ctx context.Context, req *types.QueryProjectsByRegionRequest) (*types.QueryProjectsByRegionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var matched []types.InfraProject
	err := q.k.Project.Walk(ctx, nil, func(key string, project types.InfraProject) (bool, error) {
		for _, loc := range project.Locations {
			if strings.EqualFold(loc.Region, req.Region) {
				matched = append(matched, project)
				break
			}
		}
		return false, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProjectsByRegionResponse{Projects: matched}, nil
}

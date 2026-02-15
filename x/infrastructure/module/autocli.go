package infrastructure

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"govchain/x/infrastructure/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod:      "GetProject",
					Use:            "get-project [project-id]",
					Short:          "Gets an infrastructure project by ID",
					Alias:          []string{"show-project"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "project_id"}},
				},
				{
					RpcMethod: "ListProjects",
					Use:       "list-projects",
					Short:     "List all infrastructure projects",
				},
				{
					RpcMethod:      "ProjectsBySector",
					Use:            "projects-by-sector [sector]",
					Short:          "Query projects by sector",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "sector"}},
				},
				{
					RpcMethod:      "ProjectsByStatus",
					Use:            "projects-by-status [status]",
					Short:          "Query projects by status",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "status"}},
				},
				{
					RpcMethod:      "ProjectsByRegion",
					Use:            "projects-by-region [region]",
					Short:          "Query projects by region",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "region"}},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateProject",
					Use:            "create-project [project-id] [title] [description]",
					Short:          "Create a new infrastructure project",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "project_id"}, {ProtoField: "title"}, {ProtoField: "description"}},
				},
				{
					RpcMethod:      "UpdateProject",
					Use:            "update-project [project-id]",
					Short:          "Update an infrastructure project",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "project_id"}},
				},
				{
					RpcMethod:      "UpdateProgress",
					Use:            "update-progress [project-id]",
					Short:          "Update project progress metrics",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "project_id"}},
				},
				{
					RpcMethod:      "AttachDocument",
					Use:            "attach-document [project-id]",
					Short:          "Attach a document to a project",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "project_id"}},
				},
				{
					RpcMethod:      "LinkContractingProcess",
					Use:            "link-process [project-id]",
					Short:          "Link an OCDS contracting process to a project",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "project_id"}},
				},
			},
		},
	}
}

package procurement

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"govchain/x/procurement/types"
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
					RpcMethod:      "GetProcess",
					Use:            "get-process [ocid]",
					Short:          "Gets a contracting process by OCID",
					Alias:          []string{"show-process"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "ocid"}},
				},
				{
					RpcMethod: "ListProcesses",
					Use:       "list-processes",
					Short:     "List all contracting processes",
				},
				{
					RpcMethod:      "GetRelease",
					Use:            "get-release [ocid] [release-id]",
					Short:          "Gets a specific release",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "ocid"}, {ProtoField: "release_id"}},
				},
				{
					RpcMethod:      "ListReleases",
					Use:            "list-releases [ocid]",
					Short:          "List all releases for a contracting process",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "ocid"}},
				},
				{
					RpcMethod:      "GetEntity",
					Use:            "get-entity [entity-id]",
					Short:          "Gets a registered entity",
					Alias:          []string{"show-entity"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "entity_id"}},
				},
				{
					RpcMethod: "ListEntities",
					Use:       "list-entities",
					Short:     "List all registered entities",
				},
				{
					RpcMethod: "SearchTenders",
					Use:       "search-tenders",
					Short:     "Search active tenders with filters",
				},
				{
					RpcMethod:      "GetOCDSPackage",
					Use:            "ocds-package [ocid]",
					Short:          "Get OCDS release package for a contracting process",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "ocid"}},
				},
				{
					RpcMethod:      "GetRecordPackage",
					Use:            "record-package [ocid]",
					Short:          "Get OCDS record package for a contracting process",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "ocid"}},
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
					RpcMethod: "SubmitRelease",
					Use:       "submit-release",
					Short:     "Submit an OCDS release to the chain",
				},
				{
					RpcMethod: "RegisterEntity",
					Use:       "register-entity",
					Short:     "Register an organization on-chain",
				},
				{
					RpcMethod:      "UpdateEntity",
					Use:            "update-entity [entity-id]",
					Short:          "Update a registered organization",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "entity_id"}},
				},
				{
					RpcMethod:      "CreateProcess",
					Use:            "create-process [ocid]",
					Short:          "Create a new contracting process",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "ocid"}},
				},
				{
					RpcMethod:      "PublishTender",
					Use:            "publish-tender [ocid]",
					Short:          "Publish a tender for a contracting process",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "ocid"}},
				},
				{
					RpcMethod:      "SubmitAward",
					Use:            "submit-award [ocid]",
					Short:          "Record an award decision",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "ocid"}},
				},
				{
					RpcMethod:      "SignContract",
					Use:            "sign-contract [ocid]",
					Short:          "Record a contract signing",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "ocid"}},
				},
				{
					RpcMethod:      "UpdateImplementation",
					Use:            "update-implementation [ocid] [contract-id]",
					Short:          "Record implementation progress",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "ocid"}, {ProtoField: "contract_id"}},
				},
			},
		},
	}
}

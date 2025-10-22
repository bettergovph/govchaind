# DATU Module Operation Examples

This folder hosts a comprehensive integration example that exercises every keeper operation introduced for the DATU-aligned modules. The scenario wires the budget registry, procurement ledger, disbursement tracker, accountability scores, governance voting, and WASM smart-contract registry together against a shared in-memory multistore.

The `TestDATUModuleOperationsExample` test demonstrates how a single registered smart contract can act as the connective tissue for all domain-specific records:

* The WASM keeper registers the supervising contract once and exposes it to the rest of the workflow.
* Budget, procurement, and disbursement entries reference the contract through `wasm://` metadata links to reflect the on-chain controller that authenticates them.
* Accountability scorecards and governance delegations reuse the same contract address to anchor DATU oversight signals back to the audited program.

Running the test exercises each keeper method (register, update, query, walk, and totals) and asserts that the resulting state matches the DATU invariants.

```bash
# Execute the full example
go test ./documentation/datu_examples -run TestDATUModuleOperationsExample -v
```

Use this test as a blueprint for higher-level integration tests or CLI/REST flows that need to stitch the DATU smart contract with the Cosmos SDK modules.

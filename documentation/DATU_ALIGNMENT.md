# Aligning OpenGovChain with the DATU Reference Architecture

This note maps OpenGovChain's current Cosmos SDK implementation to the DATU architecture that was originally articulated for a Stellar-based stack, then identifies concrete engineering tasks required for parity or improvements.

## 1. Consensus & Validator Topology

- **Current state**: OpenGovChain runs on Tendermint BFT with a tokenless, volunteer validator set governed through Cosmos staking and governance modules.【F:documentation/TECHNICAL_IMPLEMENTATION.md†L31-L39】【F:documentation/NETWORK_CONFIG.md†L5-L33】
- **DATU goal**: Tiered quorum slices with government, civil society, and citizen auditors.
- **Action items**:
  - Implement validator metadata and weighting to emulate tiered quorum requirements using custom staking hooks or a dedicated module that tracks institutional credentials, delegation relationships, and minimum representation checks before block proposal (e.g., pre-commit vote extensions).【F:documentation/NETWORK_CONFIG.md†L41-L54】
  - Extend governance policies to enforce multi-tier participation by linking validator groups to proposal voting and slashing conditions aligned with delegated accountability (citizen auditors borrowing weight from institutions without custody transfer).
  - Evaluate Tendermint's proposer-based consensus for deterministic quorum vs. Stellar's FBA; document mitigation strategies (e.g., config-driven validator rotation, accountability weight caps).

## 2. Smart Contract & Module Layer

- **Current state**: A single `datasets` module records IPFS-backed public data with create/update/delete flows and provenance via transaction hashes.【F:x/datasets/keeper/msg_server_entry.go†L16-L137】【F:proto/govchain/datasets/v1/entry.proto†L6-L25】
- **DATU goal**: A suite of Soroban contracts (BudgetRegistry, ProcurementLedger, DisbursementTracker, AccountabilityScores, GovernanceVoting).
- **Action items**:
  - Design Cosmos SDK modules mirroring the DATU functional scope, reusing collection patterns from `datasets` for registry-style state machines while mapping DATU-specific invariants (e.g., Transparency Units) into either a non-transferable token module or typed `sdk.Coin` aliases.
  - Introduce event schemas and document hashing strategies for each module so downstream analytics and audit tools can verify state transitions, leveraging Cosmos events the way Soroban logs were envisioned.
  - Scaffold CosmWasm support if contract flexibility similar to Soroban is required, or continue with native Go modules for deterministic performance.

## 3. Data & Document Anchoring

- **Current state**: IPFS integration is already central, with Helia-based upload flows and CID verification, plus metadata indexing via REST/gRPC queries.【F:documentation/TECHNICAL_IMPLEMENTATION.md†L113-L278】
- **DATU goal**: Anchored time-series for budget line items with Merkle proofs.
- **Action items**:
  - Version dataset entries and add Merkle root snapshots per reporting period; store roots on-chain while larger diffs reside in IPFS/Filecoin.
  - Provide GraphQL/REST endpoints that correlate on-chain entries with off-chain audit documents, reusing the existing query scaffolding.

## 4. Governance & Accountability Mechanics

- **Current state**: Tokenless governance is outlined conceptually, relying on volunteer consensus and community proposals.【F:documentation/NETWORK_CONFIG.md†L47-L54】
- **DATU goal**: Delegated accountability, on-chain governance voting with transparent audit trails.
- **Action items**:
  - Implement delegation transactions akin to `delegateAccountability()` and `reclaimAccountability()` as message types in a new accountability module, persisting delegation metadata for use by Tendermint vote extensions or governance tallies.
  - Extend governance parameters to require cross-tier quorums and publish accountability scores, combining staking power with qualitative metrics.

## 5. Application & Integration Layer

- **Current state**: Documentation covers REST interfaces, CosmJS flows, and deployment pipelines for validators and operators.【F:documentation/TECHNICAL_IMPLEMENTATION.md†L252-L383】
- **DATU goal**: Government budget portal, public transparency explorer, civic engagement APIs.
- **Action items**:
  - Expand API surface with domain-specific endpoints (budgets, procurements, disbursements) once corresponding modules exist.
  - Align UX requirements (role-based access, citizen dashboards) with existing Next.js stack mentioned in docs, ensuring DID-based auth can be integrated later.

## 6. Observability, Security, and Compliance

- **Current state**: Guidance exists for validator requirements, checksum verification, and operational playbooks (deployment, maintenance, troubleshooting).【F:documentation/TECHNICAL_IMPLEMENTATION.md†L150-L400】
- **DATU goal**: HSM-backed keys, formal verification, compliance audits.
- **Action items**:
  - Integrate key management best practices (signing service or KMS) into deployment scripts and document compliance checklists.
  - Add formal verification tooling (property-based tests, model checking) for the new accountability-critical modules.
  - Publish observability dashboards (Prometheus/Grafana) tailored to DATU metrics (budget flow latency, delegation health, cross-tier quorum stats).

---

This plan keeps Cosmos SDK and Tendermint at the core while layering DATU-specific modules and governance logic to approximate the original Stellar FBA blueprint.

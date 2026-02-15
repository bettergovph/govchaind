# OCDS → Protobuf Mapping Reference

## OpenGovChain Procurement Module
**Standard:** Open Contracting Data Standard (OCDS) v1.1.5  
**Target:** Cosmos SDK / Ignite protobuf definitions  
**Module:** `x/procurement`

---

## File Structure

```
proto/ocds/v1/
├── common.proto     # Shared building blocks (Value, Period, Item, Document, etc.)
├── release.proto    # Core OCDS Release schema (Planning, Tender, Award, Contract, Implementation)
├── codelists.proto  # OCDS codelists as enums (closed) and documented strings (open)
├── tx.proto         # Cosmos SDK transaction messages + query definitions
└── state.proto      # On-chain state objects, params, genesis, events
```

---

## Field Mapping: OCDS JSON → Protobuf

### Release (Top Level)

| OCDS JSON Field | Protobuf Field | Type | Proto File |
|---|---|---|---|
| `ocid` | `Release.ocid` | string | release.proto |
| `id` | `Release.id` | string | release.proto |
| `date` | `Release.date` | google.protobuf.Timestamp | release.proto |
| `tag` | `Release.tag` | repeated string | release.proto |
| `initiationType` | `Release.initiation_type` | string | release.proto |
| `language` | `Release.language` | string | release.proto |
| `parties` | `Release.parties` | repeated Organization | release.proto |
| `buyer` | `Release.buyer` | OrganizationReference | release.proto |
| `planning` | `Release.planning` | Planning | release.proto |
| `tender` | `Release.tender` | Tender | release.proto |
| `awards` | `Release.awards` | repeated Award | release.proto |
| `contracts` | `Release.contracts` | repeated Contract | release.proto |
| `relatedProcesses` | `Release.related_processes` | repeated RelatedProcess | release.proto |
| *(chain extension)* | `Release.publisher_address` | string | release.proto |
| *(chain extension)* | `Release.block_height` | string | release.proto |
| *(chain extension)* | `Release.tx_hash` | string | release.proto |

### Organization / Parties

| OCDS JSON Field | Protobuf Field | Type | Proto File |
|---|---|---|---|
| `parties[].id` | `Organization.id` | string | common.proto |
| `parties[].name` | `Organization.name` | string | common.proto |
| `parties[].identifier` | `Organization.identifier` | Identifier | common.proto |
| `parties[].additionalIdentifiers` | `Organization.additional_identifiers` | repeated Identifier | common.proto |
| `parties[].address` | `Organization.address` | Address | common.proto |
| `parties[].contactPoint` | `Organization.contact_point` | ContactPoint | common.proto |
| `parties[].roles` | `Organization.roles` | repeated string | common.proto |
| *(chain extension)* | `Organization.on_chain_address` | string | common.proto |
| *(chain extension)* | `Organization.is_verified` | bool | common.proto |

### Planning

| OCDS JSON Field | Protobuf Field | Type | Proto File |
|---|---|---|---|
| `planning.rationale` | `Planning.rationale` | string | release.proto |
| `planning.budget` | `Planning.budget` | Budget | release.proto |
| `planning.documents` | `Planning.documents` | repeated Document | release.proto |
| `planning.milestones` | `Planning.milestones` | repeated Milestone | release.proto |

### Budget

| OCDS JSON Field | Protobuf Field | Type | Proto File |
|---|---|---|---|
| `budget.id` | `Budget.id` | string | common.proto |
| `budget.description` | `Budget.description` | string | common.proto |
| `budget.amount` | `Budget.amount` | Value | common.proto |
| `budget.project` | `Budget.project` | string | common.proto |
| `budget.projectID` | `Budget.project_id` | string | common.proto |
| `budget.uri` | `Budget.uri` | string | common.proto |

### Tender

| OCDS JSON Field | Protobuf Field | Type | Proto File |
|---|---|---|---|
| `tender.id` | `Tender.id` | string | release.proto |
| `tender.title` | `Tender.title` | string | release.proto |
| `tender.description` | `Tender.description` | string | release.proto |
| `tender.status` | `Tender.status` | string | release.proto |
| `tender.procuringEntity` | `Tender.procuring_entity` | OrganizationReference | release.proto |
| `tender.items` | `Tender.items` | repeated Item | release.proto |
| `tender.value` | `Tender.value` | Value | release.proto |
| `tender.minValue` | `Tender.min_value` | Value | release.proto |
| `tender.procurementMethod` | `Tender.procurement_method` | string | release.proto |
| `tender.procurementMethodDetails` | `Tender.procurement_method_details` | string | release.proto |
| `tender.procurementMethodRationale` | `Tender.procurement_method_rationale` | string | release.proto |
| `tender.mainProcurementCategory` | `Tender.main_procurement_category` | string | release.proto |
| `tender.additionalProcurementCategories` | `Tender.additional_procurement_categories` | repeated string | release.proto |
| `tender.awardCriteria` | `Tender.award_criteria` | string | release.proto |
| `tender.awardCriteriaDetails` | `Tender.award_criteria_details` | string | release.proto |
| `tender.submissionMethod` | `Tender.submission_method` | repeated string | release.proto |
| `tender.submissionMethodDetails` | `Tender.submission_method_details` | string | release.proto |
| `tender.tenderPeriod` | `Tender.tender_period` | Period | release.proto |
| `tender.enquiryPeriod` | `Tender.enquiry_period` | Period | release.proto |
| `tender.hasEnquiries` | `Tender.has_enquiries` | bool | release.proto |
| `tender.eligibilityCriteria` | `Tender.eligibility_criteria` | string | release.proto |
| `tender.awardPeriod` | `Tender.award_period` | Period | release.proto |
| `tender.contractPeriod` | `Tender.contract_period` | Period | release.proto |
| `tender.numberOfTenderers` | `Tender.number_of_tenderers` | int32 | release.proto |
| `tender.tenderers` | `Tender.tenderers` | repeated OrganizationReference | release.proto |
| `tender.documents` | `Tender.documents` | repeated Document | release.proto |
| `tender.milestones` | `Tender.milestones` | repeated Milestone | release.proto |
| `tender.amendments` | `Tender.amendments` | repeated Amendment | release.proto |

### Award

| OCDS JSON Field | Protobuf Field | Type | Proto File |
|---|---|---|---|
| `awards[].id` | `Award.id` | string | release.proto |
| `awards[].title` | `Award.title` | string | release.proto |
| `awards[].description` | `Award.description` | string | release.proto |
| `awards[].status` | `Award.status` | string | release.proto |
| `awards[].date` | `Award.date` | google.protobuf.Timestamp | release.proto |
| `awards[].value` | `Award.value` | Value | release.proto |
| `awards[].suppliers` | `Award.suppliers` | repeated OrganizationReference | release.proto |
| `awards[].items` | `Award.items` | repeated Item | release.proto |
| `awards[].contractPeriod` | `Award.contract_period` | Period | release.proto |
| `awards[].documents` | `Award.documents` | repeated Document | release.proto |
| `awards[].amendments` | `Award.amendments` | repeated Amendment | release.proto |

### Contract

| OCDS JSON Field | Protobuf Field | Type | Proto File |
|---|---|---|---|
| `contracts[].id` | `Contract.id` | string | release.proto |
| `contracts[].awardID` | `Contract.award_id` | string | release.proto |
| `contracts[].title` | `Contract.title` | string | release.proto |
| `contracts[].description` | `Contract.description` | string | release.proto |
| `contracts[].status` | `Contract.status` | string | release.proto |
| `contracts[].period` | `Contract.period` | Period | release.proto |
| `contracts[].value` | `Contract.value` | Value | release.proto |
| `contracts[].items` | `Contract.items` | repeated Item | release.proto |
| `contracts[].dateSigned` | `Contract.date_signed` | google.protobuf.Timestamp | release.proto |
| `contracts[].documents` | `Contract.documents` | repeated Document | release.proto |
| `contracts[].implementation` | `Contract.implementation` | Implementation | release.proto |
| `contracts[].relatedProcesses` | `Contract.related_processes` | repeated RelatedProcess | release.proto |
| `contracts[].milestones` | `Contract.milestones` | repeated Milestone | release.proto |
| `contracts[].amendments` | `Contract.amendments` | repeated Amendment | release.proto |

### Implementation

| OCDS JSON Field | Protobuf Field | Type | Proto File |
|---|---|---|---|
| `implementation.transactions` | `Implementation.transactions` | repeated Transaction | release.proto |
| `implementation.milestones` | `Implementation.milestones` | repeated Milestone | release.proto |
| `implementation.documents` | `Implementation.documents` | repeated Document | release.proto |

### Common Building Blocks

| OCDS JSON Field | Protobuf Message | Proto File |
|---|---|---|
| Value (`amount`, `currency`) | `Value` | common.proto |
| Period (`startDate`, `endDate`, `maxExtentDate`, `durationInDays`) | `Period` | common.proto |
| Identifier (`scheme`, `id`, `legalName`, `uri`) | `Identifier` | common.proto |
| Address (`streetAddress`, `locality`, `region`, `postalCode`, `countryName`) | `Address` | common.proto |
| ContactPoint (`name`, `email`, `telephone`, `faxNumber`, `url`) | `ContactPoint` | common.proto |
| Classification (`scheme`, `id`, `description`, `uri`) | `Classification` | common.proto |
| Unit (`scheme`, `id`, `name`, `value`, `uri`) | `Unit` | common.proto |
| Item (`id`, `description`, `classification`, `additionalClassifications`, `quantity`, `unit`) | `Item` | common.proto |
| Document (`id`, `documentType`, `title`, `description`, `url`, ...) | `Document` | common.proto |
| Milestone (`id`, `title`, `type`, `description`, `code`, `dueDate`, ...) | `Milestone` | common.proto |
| Amendment (`id`, `date`, `rationale`, `description`, ...) | `Amendment` | common.proto |
| RelatedProcess (`id`, `relationship`, `title`, `scheme`, `identifier`, `uri`) | `RelatedProcess` | common.proto |
| Transaction (`id`, `source`, `date`, `value`, `payer`, `payee`, `uri`) | `Transaction` | common.proto |
| OrganizationReference (`id`, `name`) | `OrganizationReference` | common.proto |

---

## Codelist Handling

| Codelist | Type | Enforcement | Proto File |
|---|---|---|---|
| releaseTag | Closed | Enum `ReleaseTag` | codelists.proto |
| initiationType | Closed | Enum `InitiationType` | codelists.proto |
| tenderStatus | Closed | Enum `TenderStatus` | codelists.proto |
| method (procurementMethod) | Closed | Enum `ProcurementMethod` | codelists.proto |
| procurementCategory | Closed | Enum `ProcurementCategory` | codelists.proto |
| awardStatus | Closed | Enum `AwardStatus` | codelists.proto |
| contractStatus | Closed | Enum `ContractStatus` | codelists.proto |
| milestoneStatus | Closed | Enum `MilestoneStatus` | codelists.proto |
| currency | Closed | Enum `Currency` (subset) | codelists.proto |
| partyRole | Open | string (documented) | codelists.proto |
| documentType | Open | string (documented) | codelists.proto |
| awardCriteria | Open | string (documented) | codelists.proto |
| submissionMethod | Open | string (documented) | codelists.proto |
| milestoneType | Open | string (documented) | codelists.proto |
| relatedProcess | Open | string (documented) | codelists.proto |
| itemClassificationScheme | Open | string (documented) | codelists.proto |

---

## On-Chain Extensions (beyond standard OCDS)

These fields exist in the protobuf schema but NOT in standard OCDS JSON.
They provide blockchain-specific value:

| Extension | Location | Purpose |
|---|---|---|
| `publisher_address` | Release | Cosmos SDK account that submitted the release |
| `block_height` | Release | Block number for temporal anchoring |
| `tx_hash` | Release | Transaction hash for verification |
| `chain_timestamp` | Release | Consensus timestamp |
| `content_hash` | Document | SHA-256 hash of document contents |
| `content_hash_method` | Document | Hash algorithm identifier |
| `ipfs_cid` | Document | IPFS content identifier for decentralized storage |
| `on_chain_address` | Organization | Cosmos account address |
| `is_verified` | Organization | On-chain identity verification status |
| `tx_hash` | Transaction | Cosmos tx hash recording the spend |

---

## API Compatibility

The REST/gRPC API layer should transform protobuf state into standard OCDS JSON:

1. **`GET /procurement/releases/{ocid}`** → `QueryGetOCDSPackageRequest` → Standard OCDS Release Package JSON
2. **`GET /procurement/records/{ocid}`** → `QueryGetRecordPackageRequest` → Standard OCDS Record Package JSON
3. **`GET /procurement/tenders?status=active`** → `QuerySearchTendersRequest` → Filtered release packages

This ensures compatibility with existing OCDS tools (Kingfisher, Data Review Tool, Red Flags, etc.)
while the underlying data benefits from blockchain immutability and transparency.

---

## Philippine-Specific Considerations (RA 9184)

The `Params` in `state.proto` include governance-configurable thresholds aligned with
Republic Act 9184 (Government Procurement Reform Act):

- **Public bidding threshold** — Amount above which public bidding is required
- **Alternative method threshold** — Thresholds for alternative procurement methods
- **Minimum tenderers** — Minimum bidders required for open procurement
- **Default currency** — PHP (Philippine Peso)
- **Organization identifier scheme** — PH-SEC, PH-DTI, PH-PhilGEPS

These can be updated through Cosmos SDK governance proposals without chain upgrades.

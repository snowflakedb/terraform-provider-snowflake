# Hybrid Tables Integration Tests Branch - Follow-Up TODO

Tracking file for follow-up changes on `feature/hybrid-tables-integration-tests`.
Based on PR #4454 review comments and tp-research.md compliance rules.

## Status Legend
- [ ] Not started
- [x] Completed
- [~] Deferred (future PR)

---

## Phase 1: Rebase & Compilation Fixes
- [x] Rebase integration branch onto SDK branch (8e75eba5e)
- [x] Resolve rebase conflicts (ext.go, gen.go, gen_test.go, ext_test.go)
- [x] Fix validation compilation errors (remove ModifyColumnAction references)
- [x] Verify SDK and integration test packages compile

## Phase 2: tp-research.md Compliance - DDL Struct Rules (Section 1.1, 1.2)

All DDL-tagged structs MUST be in gen.go (generated from def.go). ext.go is ONLY for convert() methods, helpers, and pure Go types.

- [x] **2a.** Fix `db:"null"` -> `db:"null?"` in `hybridTableDetailsRow` (hybrid_tables_gen.go:267).
  Fixed in-place in gen.go, matching convention used by views_gen.go and materialized_views_gen.go.

## Phase 3: PR Review Comment Follow-Ups

Items from reviewer comments (sfc-gh-asawicki, sfc-gh-jmichalak) not yet applied:

- [~] **3a.** Replace `Field("Name", "string")` with dedicated type methods in def.go
  (Section 1.5: `.Text("Name")`, `.Number("Rows")`, `.Bool("IsUnique")`, `.Time("CreatedOn")`)
  Reviewer comment: "nit: we have dedicated method for types so Field could be replaced"
  **Deferred: SDK branch change, not integration tests branch.**

- [~] **3b.** Move desc structs/converters to proper locations
  Reviewer comment: "The desc structs and converters should be in a separate file"
  **Deferred: SDK branch change, not integration tests branch.**

- [~] **3c.** Indexes should eventually be their own interface (Section 1.13)
  Reviewer comment: "Index operations -> dedicated Indexes interface"
  **Deferred: Larger refactor, future PR.**

- [~] **3d.** ALTER TABLE sharing with Tables interface (Section 1.12)
  Reviewer comment: "ALTER goes through normal ALTER TABLE, maybe piggyback on Tables?"
  **Deferred: Future PR.**

- [x] **3e.** SHOW INDEXES supports LIKE, STARTS_WITH, LIMIT, full IN clause.
  Verified: `ShowIndexesHybridTableOptions` has Like, In (*TableIn), StartsWith, Limit.

- [x] **3f.** DATA_RETENTION_TIME_IN_DAYS is ALTER-only (Section 1.14).
  Verified: Only in `HybridTableSetProperties` and `HybridTableUnsetProperties`, NOT in `CreateHybridTableOptions`.

## Phase 4: Integration Test Fixes

- [x] **4a.** Integration test compiles with correct SDK types.
  Applied 7 bulk type replacements (sed) + manual fixes:
  - `HybridTableColumnsConstraintsAndIndexes` -> `HybridTableColumnsConstraintsAndIndexesRequest` (21 occurrences)
  - `[]sdk.HybridTableColumn{` -> `[]sdk.HybridTableColumnRequest{` (21 occurrences)
  - `sdk.HybridTableColumnInlineConstraint` -> `sdk.ColumnInlineConstraint` (20 occurrences)
  - `sdk.HybridTableOutOfLineConstraint` -> `sdk.HybridTableOutOfLineConstraintRequest` (3 occurrences)
  - Index request function renames (Create/Drop/Show)
  - `ShowHybridTableIndexIn` -> `TableIn`
  - `RandomAccountObjectIdentifier()` -> `RandomSchemaObjectIdentifier()` for index IDs
  - `NewDropIndexHybridTableRequest(tableId, indexName.Name())` -> `NewDropIndexHybridTableRequest(indexName)` (2->1 args)
  - `Table: &tableId` -> `Table: tableId` (pointer->value for TableIn.Table)

- [x] **4b.** Pointer field nil checks added for `HybridTableIndex`.
  - `found.Columns` -> `*found.Columns` with `require.NotNil` guard
  - `found.IsUnique` -> `*found.IsUnique` with `require.NotNil` guard

- [x] **4c.** Helper file (`hybrid_table_client.go`) updated to use Request/DTO types.
  - `HybridTableColumnsConstraintsAndIndexes` -> `HybridTableColumnsConstraintsAndIndexesRequest`
  - `[]HybridTableColumn` -> `[]HybridTableColumnRequest`
  - `HybridTableColumnInlineConstraint` -> `ColumnInlineConstraint`

- [ ] **4d.** Runtime verification of index identifier types (SchemaObjectIdentifier).
  Needs actual Snowflake connection to verify. Cannot be done locally.

## Phase 5: Code Quality & Prechecks

- [x] **5a.** `go build ./...` - PASS (full project build, exit code 0)
- [x] **5b.** `go test ./pkg/sdk/ -run HybridTable` - PASS (all hybrid table unit tests pass)
  Note: `TestClient_*` tests fail (pre-existing, need Snowflake connection config).
- [x] **5c.** `make fmt` - PASS (gofumpt applied to 3 hybrid table files)
- [x] **5d.** `go vet ./pkg/sdk/...` - PASS (no errors)
  `go vet -tags=non_account_level_tests ./pkg/sdk/testint/...` - no hybrid table errors.
  Note: `make lint` crashes due to golangci-lint go1.25 vs go1.26 version mismatch (pre-existing env issue).

---

## Notes

- tp-research.md location: ~/Code/tp-research.md
- SDK branch: feature/hybrid-tables-sdk (commit 8e75eba5e)
- Integration branch: feature/hybrid-tables-integration-tests
- PR #4454: Hybrid Tables SDK
- All `In` -> `TableIn` conversions done for ShowHybridTable and ShowIndexes
- `WithAlterColumnAction` takes `[]HybridTableAlterColumnActionRequest` (slice, not single value)

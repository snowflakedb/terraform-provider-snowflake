# Type-Drift CustomizeDiff Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move the Read-time substitution of `nullable` and `collate` out of `buildHybridColumnStateFromDescribe` and into a CustomizeDiff (for `nullable` on PK columns) plus a `DiffSuppressFunc` (for `collate`), so the Read path writes exactly what Snowflake returned, matching jmichalak's review feedback that Terraform reconciliation belongs in custom diff, not in Read.

**Architecture:** Read writes raw DESCRIBE values into state (no config substitution). A new `suppressNullableForPrimaryKeyColumns` CustomizeDiff suppresses the spurious `nullable: false → true` diff that PK columns produce because Snowflake silently enforces NOT NULL on PK columns. A new `DiffSuppressCaseInsensitive` `DiffSuppressFunc` wired onto the `collate` field tolerates case differences between the user's spelling and Snowflake's canonical form.

**Tech Stack:** Go 1.x, terraform-plugin-sdk v2, Snowflake Go SDK (`pkg/sdk`), terraform-plugin-testing acceptance framework.

---

## File Structure

- Modify: `pkg/resources/hybrid_table.go:100-105` (collate schema field — add `DiffSuppressFunc`)
- Modify: `pkg/resources/hybrid_table.go:250-257` (CustomizeDiff chain — add `suppressNullableForPrimaryKeyColumns()`)
- Modify: `pkg/resources/hybrid_table.go:839+` (`buildHybridColumnStateFromDescribe` — drop config substitution)
- Modify: `pkg/resources/hybrid_table.go:1023+` (add `suppressNullableForPrimaryKeyColumns` next to existing `forceNewIfColumn*` helpers)
- Reuse: `ignoreCaseSuppressFunc` at `pkg/resources/diff_suppressions.go:302` (already exists; do NOT add a new helper)
- Test: `pkg/testacc/resource_hybrid_table_acceptance_test.go` (add `TestAcc_HybridTable_PKNullableNoSpurious` and `TestAcc_HybridTable_CollateCaseInsensitive` — code only; do NOT run)

---

## Task 1: SKIPPED — reuse existing helper

The original plan added a new exported `DiffSuppressCaseInsensitive`. An equivalent private helper `ignoreCaseSuppressFunc` already exists at `pkg/resources/diff_suppressions.go:302`:

```go
func ignoreCaseSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}
```

Sibling helpers in the same file (`ignoreTrimSpaceSuppressFunc`, `ignoreCaseAndTrimSpaceSuppressFunc`, `ignoreAlwaysSuppressFunc`) follow the same private camelCase convention. Adding a new exported helper would be redundant. Task 2 wires the existing `ignoreCaseSuppressFunc` directly.

- [x] **Skipped** — no new helper is added; no commit for Task 1.

---

## Task 2: Wire `DiffSuppressFunc` onto the `collate` field

**Files:**
- Modify: `pkg/resources/hybrid_table.go:100-105`

- [x] **Step 1: Inspect the current schema declaration**

Run: `grep -n -A4 '"collate":' pkg/resources/hybrid_table.go | head -20`
Expected: shows the `collate` block at L100-105 with no `DiffSuppressFunc`.

- [x] **Step 2: Apply edit**

Replace at `pkg/resources/hybrid_table.go:100-105`:

```go
				"collate": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: "Column collation specification, e.g. en-ci.",
				},
```

with:

```go
				"collate": {
					Type:             schema.TypeString,
					Optional:         true,
					Default:          "",
					DiffSuppressFunc: ignoreCaseSuppressFunc,
					Description:      "Column collation specification, e.g. en-ci. Case-insensitive (en-ci and EN-CI are treated as equal).",
				},
```

- [x] **Step 3: Verify the change compiles**

Run: `go build ./pkg/resources/...`
Expected: no errors.

- [x] **Step 4: Verify provider validation still passes**

Run: `go test ./pkg/provider/ -run TestProvider -v -count=1`
Expected: PASS.

- [x] **Step 5: Commit**

```bash
git add pkg/resources/hybrid_table.go
git commit -m "feat(hybrid_table): tolerate case-only collate drift

Wires ignoreCaseSuppressFunc onto column.collate so that a
DESCRIBE-returned EN-CI does not produce a spurious diff against
a user-supplied en-ci.

Reuses the existing private helper from diff_suppressions.go
rather than introducing a new one."
```

---

## Task 3: Add the `suppressNullableForPrimaryKeyColumns` CustomizeDiff

**Files:**
- Modify: `pkg/resources/hybrid_table.go` (add helper next to existing `forceNewIfColumnFieldChanged` at L955-977)
- Modify: `pkg/resources/hybrid_table.go:250-257` (chain it into `customdiff.All`)

- [x] **Step 1: Write the failing acceptance test**

Append to `pkg/testacc/resource_hybrid_table_acceptance_test.go`:

```go
// TestAcc_HybridTable_PKNullableNoSpurious verifies that a primary-key column
// declared as nullable=true (the schema default) does not produce a spurious
// diff after Read, even though Snowflake silently enforces NOT NULL on PK
// columns and DESCRIBE reports null="N". The reconciliation must happen in
// CustomizeDiff, not via Read-time substitution.
func TestAcc_HybridTable_PKNullableNoSpurious(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}
	pk := []sdk.TableColumnSignature{
		{Name: "ID"},
	}
	model := model.HybridTableFromId("test", id, columns, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model),
			},
			// A second apply with no config change must produce a no-op plan.
			{
				Config: accconfig.FromModels(t, model),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
```

- [x] **Step 2: Run the failing test** — SKIPPED (acceptance tests not run per task instructions)

Run: `TF_ACC=1 go test ./pkg/testacc/ -run TestAcc_HybridTable_PKNullableNoSpurious -v -timeout 30m`
Expected: PASS today (because Read currently substitutes nullable=true from config). The test is the contract we want preserved after the refactor.

If it fails, stop — that signals the existing substitution logic has already broken and the rest of this plan needs re-scoping.

- [x] **Step 3: Write the helper (still preserving Read substitution; will be flipped in Task 4)**

Insert into `pkg/resources/hybrid_table.go` immediately after `forceNewIfColumnFieldChanged` (currently ends at L977):

```go
// suppressNullableForPrimaryKeyColumns clears the planned diff on
// `column.<idx>.nullable` when the column is part of `primary_key.0.keys`
// and the only change is from old=false (DESCRIBE-derived) to new=true
// (the schema default). Snowflake silently enforces NOT NULL on PK columns,
// so DESCRIBE always returns null="N" for them — the diff is spurious and
// has no corresponding ALTER (hybrid tables do not support SET/DROP NOT NULL).
//
// This sits in CustomizeDiff rather than DiffSuppressFunc because resolving
// the diff requires looking at sibling fields (primary_key.0.keys), which a
// per-field DiffSuppressFunc cannot access.
func suppressNullableForPrimaryKeyColumns() schema.CustomizeDiffFunc {
	return func(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
		if !diff.HasChange("column") {
			return nil
		}
		pkRaw, ok := diff.Get("primary_key").([]any)
		if !ok || len(pkRaw) == 0 {
			return nil
		}
		pkMap, ok := pkRaw[0].(map[string]any)
		if !ok {
			return nil
		}
		pkKeysRaw, ok := pkMap["keys"].([]any)
		if !ok {
			return nil
		}
		pkKeys := make(map[string]struct{}, len(pkKeysRaw))
		for _, k := range pkKeysRaw {
			if s, ok := k.(string); ok {
				pkKeys[strings.ToUpper(s)] = struct{}{}
			}
		}

		oldRaw, newRaw := diff.GetChange("column")
		oldCols := parseHybridColumns(oldRaw)
		newCols := parseHybridColumns(newRaw)
		for newIdx, n := range newCols {
			if _, isPK := pkKeys[strings.ToUpper(n.name)]; !isPK {
				continue
			}
			for _, o := range oldCols {
				if o.name != n.name {
					continue
				}
				if !o.nullable && n.nullable {
					if err := diff.Clear(fmt.Sprintf("column.%d.nullable", newIdx)); err != nil {
						return err
					}
				}
				break
			}
		}
		return nil
	}
}
```

- [x] **Step 4: Wire it into the CustomizeDiff chain**

Replace at `pkg/resources/hybrid_table.go:250-257`:

```go
		CustomizeDiff: TrackingCustomDiffWrapper(resources.HybridTable, customdiff.All(
			hybridTableParametersCustomDiff,
			ComputedIfAnyAttributeChanged(hybridTableSchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(hybridTableSchema, DescribeOutputAttributeName, "name", "comment", "column"),
			ComputedIfAnyAttributeChanged(hybridTableSchema, FullyQualifiedNameAttributeName, "name"),
			forceNewIfColumnCollateChanged(),
			forceNewIfColumnNullableChanged(),
		)),
```

with:

```go
		CustomizeDiff: TrackingCustomDiffWrapper(resources.HybridTable, customdiff.All(
			hybridTableParametersCustomDiff,
			ComputedIfAnyAttributeChanged(hybridTableSchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(hybridTableSchema, DescribeOutputAttributeName, "name", "comment", "column"),
			ComputedIfAnyAttributeChanged(hybridTableSchema, FullyQualifiedNameAttributeName, "name"),
			suppressNullableForPrimaryKeyColumns(),
			forceNewIfColumnCollateChanged(),
			forceNewIfColumnNullableChanged(),
		)),
```

NOTE: `suppressNullableForPrimaryKeyColumns` runs **before** `forceNewIfColumnNullableChanged` so the spurious nullable diff is cleared before the latter would otherwise trigger ForceNew on it.

- [x] **Step 5: Verify build and provider validation**

Run: `go build ./pkg/resources/... && go test ./pkg/provider/ -run TestProvider -v -count=1`
Expected: build passes, TestProvider passes.

- [x] **Step 6: Verify the acceptance test still passes** — SKIPPED (acceptance tests not run per task instructions)

Run: `TF_ACC=1 go test ./pkg/testacc/ -run TestAcc_HybridTable_PKNullableNoSpurious -v -timeout 30m`
Expected: PASS (the refactor in Task 4 hasn't started yet, so Read still substitutes — but the new CustomizeDiff is now in place and idempotent under the existing setup).

- [x] **Step 7: Commit**

```bash
git add pkg/resources/hybrid_table.go pkg/testacc/resource_hybrid_table_acceptance_test.go
git commit -m "feat(hybrid_table): suppress spurious nullable diff on PK columns

Adds suppressNullableForPrimaryKeyColumns CustomizeDiff so the
old=false → new=true diff on a PK column (where Snowflake silently
enforces NOT NULL) is cleared at plan time instead of being papered
over by Read-time substitution.

Step 1 of moving Snowflake-vs-config reconciliation out of Read."
```

---

## Task 4: Drop the Read-time substitution

**Files:**
- Modify: `pkg/resources/hybrid_table.go:776-847` (`buildHybridColumnStateFromDescribe`)

- [ ] **Step 1: Add an acceptance test for the collate case-only-drift case**

Append to `pkg/testacc/resource_hybrid_table_acceptance_test.go`:

```go
// TestAcc_HybridTable_CollateCaseInsensitive verifies that a config-supplied
// collate of "en-ci" produces no spurious diff even if DESCRIBE returns it
// as "EN-CI" or some other case variant. Reconciliation must come from the
// DiffSuppressFunc on the field, not from Read-time substitution.
func TestAcc_HybridTable_CollateCaseInsensitive(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar, Collate: sdk.String("en-ci")},
	}
	pk := []sdk.TableColumnSignature{
		{Name: "ID"},
	}
	model := model.HybridTableFromId("test", id, columns, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model),
			},
			{
				Config: accconfig.FromModels(t, model),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
```

NOTE: The `model.HybridTableFromId` builder may not yet expose `Collate` per-column. If it doesn't, add a `WithColumnCollate` builder method or — if that pulls in too much scaffolding — write the config inline using `config.HybridTableConfig` raw HCL. Verify with: `grep -n "Collate" pkg/acceptance/bettertestspoc/config/model/hybrid_table_model.go`.

- [ ] **Step 2: Run the failing test**

Run: `TF_ACC=1 go test ./pkg/testacc/ -run TestAcc_HybridTable_CollateCaseInsensitive -v -timeout 30m`
Expected: PASS today (Read substitutes config collate; DiffSuppress is also in place from Task 2). This test pins the contract.

- [ ] **Step 3: Edit `buildHybridColumnStateFromDescribe` to drop the substitution**

Replace at `pkg/resources/hybrid_table.go:776-847`:

```go
func buildHybridColumnStateFromDescribe(details []sdk.HybridTableDetails, d *schema.ResourceData) ([]any, error) {
	flattened := make([]any, 0)

	// Build lookups from config for fields where DESCRIBE returns a different value
	// than the user wrote and the framework cannot reconcile it via DiffSuppressFunc:
	// - collate: DESCRIBE returns "X COLLATE 'Y'" combined; the SDK splits it but the
	//   server-side spelling can differ from what the user wrote (e.g. case).
	// - nullable: DESCRIBE returns NOT NULL for PK columns even when config says nullable=true.
	//   Preserve the config value since Snowflake silently enforces NOT NULL on PK columns.
	//
	// Data type drift (e.g. INTEGER vs NUMBER(38,0)) is handled by DiffSuppressDataTypes
	// on the column.type field — we do not need to substitute the config value here.
	type configColumnInfo struct {
		collate  string
		nullable bool
	}
	configByName := make(map[string]configColumnInfo)
	if configCols, ok := d.GetOk("column"); ok {
		for _, rawCol := range configCols.([]any) {
			colMap := rawCol.(map[string]any)
			if colName, ok := colMap["name"].(string); ok {
				info := configColumnInfo{nullable: true}
				if collate, ok := colMap["collate"].(string); ok {
					info.collate = collate
				}
				if nullable, ok := colMap["nullable"].(bool); ok {
					info.nullable = nullable
				}
				configByName[strings.ToUpper(colName)] = info
			}
		}
	}

	for _, td := range details {
		if td.Kind != "COLUMN" {
			continue
		}

		colInfo, found := configByName[strings.ToUpper(td.Name)]
		if !found {
			// Externally added column not in config: schema default is nullable=true.
			// Without this, zero-value would set nullable=false, which is ForceNew and
			// would produce a spurious delete+create plan instead of the desired Update.
			colInfo = configColumnInfo{nullable: true}
		}
		// Prefer the user's config value for collate (it preserves the exact spelling
		// the user wrote, e.g. "en-ci"). Fall back to the SDK-split Collation from
		// DESCRIBE for imported tables where no config exists.
		collate := colInfo.collate
		if collate == "" && td.Collation != nil {
			collate = *td.Collation
		}
		flat := map[string]any{
			"name":     td.Name,
			"type":     td.Type,
			"nullable": colInfo.nullable,
			"comment":  td.Comment,
			"collate":  collate,
		}

		def, err := toHybridColumnDefaultConfig(td)
		if err != nil {
			return nil, err
		}
		if def != nil {
			flat["default"] = def
		}

		flattened = append(flattened, flat)
	}
	return flattened, nil
}
```

with:

```go
func buildHybridColumnStateFromDescribe(details []sdk.HybridTableDetails, d *schema.ResourceData) ([]any, error) {
	flattened := make([]any, 0)

	// Reconciliation between the raw DESCRIBE values and user config is now done
	// in CustomizeDiff and DiffSuppressFunc, not here:
	// - column.<idx>.nullable: PK columns always come back as NOT NULL; the
	//   spurious old=false → new=true diff is cleared by
	//   suppressNullableForPrimaryKeyColumns.
	// - column.<idx>.collate: case-only drift is absorbed by
	//   DiffSuppressCaseInsensitive on the field.
	// - column.<idx>.type: format normalization (e.g. INTEGER vs NUMBER(38,0))
	//   is absorbed by DiffSuppressDataTypes on the field.
	for _, td := range details {
		if td.Kind != "COLUMN" {
			continue
		}

		collate := ""
		if td.Collation != nil {
			collate = *td.Collation
		}
		flat := map[string]any{
			"name":     td.Name,
			"type":     td.Type,
			"nullable": !strings.EqualFold(td.IsNullable, "N") && td.IsNullable != "" && !descIsNotNull(td),
			"comment":  td.Comment,
			"collate":  collate,
		}

		def, err := toHybridColumnDefaultConfig(td)
		if err != nil {
			return nil, err
		}
		if def != nil {
			flat["default"] = def
		}

		flattened = append(flattened, flat)
	}
	return flattened, nil
}
```

NOTE: `td.IsNullable` is `bool` per `pkg/sdk/hybrid_tables_gen.go:265` — verify the actual type. If it is `bool`, simplify to `flat["nullable"] = td.IsNullable`. The above is defensive in case the field has been switched to `string` in a parallel change. Run: `grep -n "IsNullable" pkg/sdk/hybrid_tables_gen.go`.

If `IsNullable` is `bool`, replace the `nullable` line with:

```go
			"nullable": td.IsNullable,
```

and delete the `descIsNotNull` reference (it should not exist).

- [ ] **Step 4: Verify build**

Run: `go build ./pkg/resources/...`
Expected: no errors.

- [ ] **Step 5: Run the two acceptance tests added in Tasks 3 and 4**

Run:
```bash
TF_ACC=1 go test ./pkg/testacc/ \
  -run 'TestAcc_HybridTable_PKNullableNoSpurious|TestAcc_HybridTable_CollateCaseInsensitive' \
  -v -timeout 60m
```
Expected: both PASS. The CustomizeDiff (Task 3) and the DiffSuppressFunc (Task 2) now do the reconciliation that Read used to do.

- [ ] **Step 6: Run the existing hybrid_table acceptance tests**

Run: `TF_ACC=1 go test ./pkg/testacc/ -run 'TestAcc_HybridTable_' -v -timeout 90m`
Expected: all PASS, including `TestAcc_HybridTable_Basic` and `TestAcc_HybridTable_CompleteUseCase`. Any new failure indicates the refactor regressed an existing scenario.

- [ ] **Step 7: Commit**

```bash
git add pkg/resources/hybrid_table.go pkg/testacc/resource_hybrid_table_acceptance_test.go
git commit -m "refactor(hybrid_table): drop Read-time config substitution

Read now writes raw DESCRIBE values for nullable and collate. Diff
reconciliation lives entirely in CustomizeDiff
(suppressNullableForPrimaryKeyColumns) and on the collate field's
DiffSuppressFunc, matching the Terraform philosophy that Read
reflects what the upstream system reports.

Addresses PR #4689 review feedback from sfc-gh-jmichalak."
```

---

## Self-Review Checklist

1. **Spec coverage:**
   - "Move nullable/collate substitution out of Read" → Task 4
   - "Use CustomizeDiff for nullable on PK columns" → Task 3
   - "Use DiffSuppressFunc for collate" → Tasks 1-2
   - "Don't change behavior for the user" → acceptance tests in Tasks 3-4 pin the contract

2. **Placeholder scan:** None of "TBD", "implement later", "fill in details" appear. Each step has runnable code or commands. NOTE blocks call out conditional follow-ups (verify field type before assuming) — those are concrete instructions, not placeholders.

3. **Type consistency:** `suppressNullableForPrimaryKeyColumns()` returns `schema.CustomizeDiffFunc`, matching the slice signature of `customdiff.All` used at L250. `DiffSuppressCaseInsensitive` matches `schema.SchemaDiffSuppressFunc` signature `(k, old, new string, d *schema.ResourceData) bool`.

---

## Execution Handoff

Plan complete and saved to `docs/superpowers/plans/2026-06-01-pr4689-3-typedrift-customizediff.md`. Two execution options:

1. **Subagent-Driven (recommended)** — fresh subagent per task, review between tasks, fast iteration.
2. **Inline Execution** — execute tasks in this session using `superpowers:executing-plans`, batched with checkpoints.

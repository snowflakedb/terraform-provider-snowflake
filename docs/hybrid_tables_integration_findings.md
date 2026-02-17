# Hybrid Tables Integration Test Findings

**Date**: 2026-02-17
**Tests Run**: 26 passing, 4 blocked by SDK design issue
**File**: `pkg/sdk/testint/hybrid_tables_integration_test.go` (897 lines)

## Critical Findings

### ✅ FINDING #1: NULL Values in Rows/Bytes Columns (FIXED)
**Issue**: SHOW HYBRID TABLES returns NULL for rows/bytes on newly created empty tables
**Expected**: Integer values
**Actual**: NULL values
**Fix Applied**:
- Changed `hybridTableRow.Rows` from `int` to `sql.NullInt64`
- Changed `hybridTableRow.Bytes` from `int` to `sql.NullInt64`
- Updated `convert()` function to handle nullable integers
- **File**: `pkg/sdk/hybrid_tables_gen.go` (lines 127-137)
- **File**: `pkg/sdk/hybrid_tables_ext.go` (convert function)

**Impact**: All SHOW operations now work correctly

---

### ✅ FINDING #2: Column Names Returned in Uppercase (DOCUMENTED)
**Issue**: DESCRIBE TABLE returns column names in uppercase regardless of how defined
**Expected**: Column names as defined (e.g., "id")
**Actual**: Uppercase column names (e.g., "ID")
**Resolution**: This is standard Snowflake behavior - documented in tests
**Test Adjustments**: Updated test assertions to expect uppercase

**Impact**: No code changes needed, test assertions adjusted

---

### ✅ FINDING #3: ALTER COLUMN COMMENT Syntax Issue (FIXED)
**Issue**: Generator produced incorrect SQL syntax for ALTER COLUMN COMMENT
**Expected**: `ALTER COLUMN name COMMENT 'value'` (no equals sign)
**Actual**: `ALTER COLUMN name COMMENT = 'value'` (with equals sign)
**Error**: `001003 (42000): SQL compilation error: syntax error`

**Fix Applied**:
- Changed DDL tag from `ddl:"parameter,single_quotes"` to `ddl:"parameter,no_equals,single_quotes"`
- **File**: `pkg/sdk/hybrid_tables_gen.go` (HybridTableAlterColumnAction struct)
- Added explanatory comment per rule 13

**Impact**: ALTER COLUMN SET COMMENT now works correctly

---

### ❌ FINDING #4: Index Names Cannot Be Schema-Qualified (SDK DESIGN BLOCKER)
**Issue**: Index names require unqualified names but SDK uses SchemaObjectIdentifier
**Expected**: `CREATE INDEX idx_name ON table...` (unqualified index name)
**Actual**: SDK generates `CREATE INDEX db.schema.idx_name...` (qualified)
**Error**: `001003 (42000): SQL compilation error: Explicit qualification of index name is not allowed`

**Affected Operations**:
- CREATE INDEX
- DROP INDEX
- SHOW INDEXES (works but can't test index creation)

**Affected Files**:
- `pkg/sdk/hybrid_tables_ext.go`:
  - `CreateHybridTableIndexRequest.name` (SchemaObjectIdentifier)
  - `DropHybridTableIndexRequest.name` (SchemaObjectIdentifier)

**Recommended Fix**:
1. Change index name type from `SchemaObjectIdentifier` to `string`
2. Document that indexes are scoped to their table, not to schema
3. Update DROP INDEX to use simple string names

**Status**: 🚫 **BLOCKING** - Cannot test INDEX operations until SDK redesigned
**Tests Affected**: 4 tests blocked (CREATE INDEX basic, CREATE INDEX with INCLUDE, DROP INDEX, SHOW INDEXES)

**SNOW Ticket**: TODO - Create ticket for SDK maintainers

---

## Working Features (26 Tests Passing)

### CREATE Operations (9 tests) ✅
- Basic with single primary key
- With table COMMENT
- Composite primary key (multi-column)
- Column with DEFAULT expression
- Column with NOT NULL constraint
- Column with IDENTITY (auto-increment)
- Multiple data types (VARCHAR, DATE, TIMESTAMP_NTZ, BOOLEAN, VARIANT, DECIMAL)
- Out-of-line UNIQUE constraint
- Out-of-line FOREIGN KEY constraint

### ALTER Operations (6 tests) ✅
- ALTER COLUMN SET COMMENT
- ALTER COLUMN UNSET COMMENT
- SET DATA_RETENTION_TIME_IN_DAYS
- SET COMMENT (table level)
- UNSET COMMENT (table level)
- UNSET DATA_RETENTION_TIME_IN_DAYS

### SHOW Filter Operations (5 tests) ✅
- SHOW with LIKE pattern
- SHOW with IN DATABASE
- SHOW with IN SCHEMA
- SHOW with STARTS WITH
- SHOW with LIMIT

### SHOW/DESCRIBE Validation (4 tests) ✅
- SHOW validates all 9 columns (created_on, name, database_name, schema_name, owner, rows, bytes, comment, owner_role_type)
- DESCRIBE validates all 13 columns (name, type, kind, null?, default, primary key, unique key, check, expression, comment, policy name, privacy domain, schema_evolution_record)
- ShowByID retrieval
- DESCRIBE with varied column configurations

### DROP Operations (2 tests) ✅
- Drop basic
- Drop with IF EXISTS (non-existent table)

---

## Test Execution Details

**Total Tests**: 30 (26 passing, 4 blocked)
**Execution Time**: ~167 seconds (2.8 minutes)
**Pass Rate**: 100% of non-blocked tests

**Command Used**:
```bash
SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE=true go test -v -tags=non_account_level_tests ./pkg/sdk/testint -run TestInt_HybridTables
```

---

## Recommendations

### For SDK Maintainers
1. **Critical**: Fix index naming to use unqualified names
2. Review NULL handling patterns across all object types
3. Consider adding DDL validation tests before SQL generation
4. Document uppercase column name behavior in DESCRIBE operations

### For Resource Implementation
1. Use SDK client for all operations (don't use raw SQL)
2. Handle NULL values for rows/bytes in SHOW output
3. Expect uppercase column names from DESCRIBE
4. Wait for index SDK fix before implementing index management

### For Integration Tests
1. All CREATE, ALTER, SHOW, DROP operations fully validated
2. Foreign key relationships tested with parent/child tables
3. All SHOW filters (LIKE, IN, STARTS WITH, LIMIT) working
4. INDEX operations require SDK redesign before testing

---

## Next Steps

1. **Create SNOW Ticket**: Index naming SDK design issue
2. **Update SDK Generator**: Fix ALTER COLUMN COMMENT syntax generation
3. **Document Workarounds**: For teams needing index operations
4. **Resource Implementation**: Can proceed with all non-index features
5. **Unit Tests**: Ready to implement SDK unit tests (Task #20)

---

*Generated from integration test results on 2026-02-17*

# Implementation Plan

| Field | Value |
|-------|-------|
| Task | Solve https://github.com/snowflakedb/terraform-provider-snowflake/issues/3946 |
| Date | 2026-02-24 |
| Agent | task-d0cbd9f2 |
| Repository | snowflakedb/terraform-provider-snowflake |
| PRs | 1 |

## Overview

This is a straightforward bug fix with minimal scope. The issue involves swapping two parameters in the UpdateFailoverGroup function where sdk.NewAccountIdentifier is called with arguments in the wrong order (accountName, organizationName instead of organizationName, accountName). The fix involves: (1) Correcting parameter order in 2 locations (lines 608 and 617 in pkg/resources/failover_group.go), (2) Adding a test case that updates allowed_accounts to verify the fix, (3) Removing the documentation warning about this known issue, (4) Adding migration guide entry about the bug fix. Total estimated diff: ~80-100 lines, well under the 400-line minimum for a single PR. This is an atomic change that should not be split - the bug fix, test, and documentation updates all belong together as a single logical unit.

## PR Stack

### PR 1: Fix allowed_accounts update in failover_group

**Description**: ## Summary
- Fix parameter order bug in UpdateFailoverGroup where sdk.NewAccountIdentifier was called with swapped arguments (accountName, organizationName instead of organizationName, accountName)
- Add acceptance test TestAcc_FailoverGroup_UpdateAllowedAccounts that specifically updates allowed_accounts field to verify the fix works
- Remove documentation warning about this known issue from templates/resources/failover_group.md.tmpl
- Add migration guide entry documenting the bug fix

## Context
Issue #3946 reported that updating allowed_accounts fails. Root cause: parameters to sdk.NewAccountIdentifier() were swapped in the Update function (lines 608 and 617), while Create function (line 210) had correct order. SDK signature is NewAccountIdentifier(organizationName, accountName).

## Test plan
- Run existing failover group acceptance tests to ensure no regressions
- Run new TestAcc_FailoverGroup_UpdateAllowedAccounts test to verify allowed_accounts can be updated
- Verify generated documentation no longer contains the warning

🤖 Generated with [Claude Code](https://claude.com/claude-code)

**Scope**:
1. Fix parameter order in pkg/resources/failover_group.go in the UpdateFailoverGroup function:
   - Line 608: Change `accountIdentifier := sdk.NewAccountIdentifier(accountName, organizationName)` to `accountIdentifier := sdk.NewAccountIdentifier(organizationName, accountName)`
   - Line 617: Change `accountIdentifier := sdk.NewAccountIdentifier(accountName, organizationName)` to `accountIdentifier := sdk.NewAccountIdentifier(organizationName, accountName)`
   - These are in the d.HasChange("allowed_accounts") block where old and new allowed accounts are parsed
   - This matches the correct order used in CreateFailoverGroup at line 210

2. Add test in pkg/testacc/resource_failover_group_acceptance_test.go:
   - Create new test function `TestAcc_FailoverGroup_UpdateAllowedAccounts` following pattern of existing tests (e.g., TestAcc_FailoverGroupBasic)
   - Add test setup: `_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)` and get test account
   - Test should have 3 steps:
     a) Step 1: Create failover group with initial allowed_accounts using existing failoverGroupBasic helper or similar
     b) Step 2: Update config to use different allowed_accounts value (create helper function if needed following pattern of failoverGroupWithChanges)
     c) Step 3: Import test to verify state
   - Use resource.TestCheckResourceAttr to assert allowed_accounts values in each step
   - Follow the structure and patterns from TestAcc_FailoverGroup_issue2544

3. Update templates/resources/failover_group.md.tmpl:
   - Remove line 14 containing: `!> **Updating allowed_accounts** Currently, updating the allowed_accounts field may fail due to an incorrect query being sent (see [#3946]...`
   - Keep line 13 (blank line) to maintain spacing between the preview warning and the page title

4. Update MIGRATION_GUIDE.md:
   - Locate the latest version section (search for ## v pattern at top of file)
   - Add new subsection: `### *(bugfix)* Fixed allowed_accounts update in snowflake_failover_group`
   - Add description: "Issue [#3946](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3946) has been resolved. Previously, updating the allowed_accounts field would fail due to swapped parameters in the SDK call. This has been fixed and allowed_accounts can now be updated correctly without requiring workarounds."
   - Add reference line: `References: [#3946](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3946)`
   - Follow the formatting pattern of other bugfix entries in the guide

**Rationale**: This is a single atomic bug fix where all components must ship together. The 2-line code fix corrects the parameter order bug. The test proves the fix works and prevents regression. The documentation updates inform users the issue is resolved. Splitting would create meaningless PRs - a code fix without tests is unverifiable, tests without the fix would fail, and documentation updates without the fix would be misleading. At ~80-100 total diff lines, this is well under the 400-line threshold and represents the minimal complete unit of work.

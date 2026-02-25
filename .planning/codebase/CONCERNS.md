# Codebase Concerns

**Analysis Date:** 2026-02-25

## Tech Debt

**Identifier Migration and Validation (SNOW-1495079, SNOW-1634872):**
- Issue: Identifier parsing and validation are not standardized across the codebase. Multiple identifier formats are supported (pipe-separated `|`, dot-separated `.`) but validation is skipped for AccountObjectIdentifier due to inability to differentiate identifier parts.
- Files:
  - `pkg/resources/deprecated_identifier_helpers.go` - Contains FormatFullyQualifiedObjectID and ParseFullyQualifiedObjectID marked for replacement
  - `pkg/internal/provider/validators/validators.go` - IsValidIdentifier skips validation for AccountObjectIdentifier
  - `pkg/resources/identifier_state_upgraders.go` - Multiple state upgrader implementations for different resources
- Impact: Potential validation gaps for account-scoped identifiers; identifier parsing logic is duplicated across the codebase; future identifier rework required
- Fix approach: Complete identifier rework as planned in SNOW-1634872; implement unified identifier parsing/validation across all object types; consolidate deprecated identifier helpers into centralized SDK-based approach

**UNSET Clause Not Implemented (SNOW-1515781):**
- Issue: Multiple security integrations and OAuth resources cannot UNSET optional parameters; forced to use SET with empty values which is invalid for many field types, requiring conditional ForceNew logic as workaround
- Files:
  - `pkg/resources/saml2_integration.go` - 9+ TODOs noting UNSET not implemented for optional fields
  - `pkg/resources/api_authentication_integration_common.go` - Similar UNSET limitations
  - `pkg/resources/external_oauth_integration.go` - OAuth parameter handling blocked by UNSET absence
  - `pkg/resources/oauth_integration_for_custom_clients.go` - Custom OAuth client parameters
- Impact: Inability to cleanly remove optional parameters from existing resources; workaround creates fragile state transitions and forces resource recreation unnecessarily; user experience degraded when unsetting previously-set values
- Fix approach: Implement UNSET clause support in SDK client for affected integration types; update resource logic to use UNSET instead of SET with empty values; remove ForceNew workarounds

**Parameter Handling Split Across Resources (SNOW-2298249):**
- Issue: Account parameters are handled in two separate resources with partial overlap: `account_parameter` (single parameter) and `account_parameters` (all parameters). List of supported parameters diverges; parameters available in `current_account` data source are not unified with these resources.
- Files:
  - `pkg/resources/account_parameter.go` - Single parameter resource with 20+ supported parameters listed
  - `pkg/resources/account_parameters.go` - Bulk parameters resource with different parameter set
  - Multiple data source implementations with redundant parameter definitions
- Impact: Inconsistent parameter support across resources; user confusion about which resource handles which parameter; maintenance burden with parameter list duplication; SDK client methods may not support all operations uniformly
- Fix approach: Unify parameter handling across all three implementations (single, bulk, data source); consolidate supported parameters list; update SDK client SetAccountParameter method for consistency

**Data Type Handling Complexity (SNOW-2054240, SNOW-2054235, SNOW-1596962):**
- Issue: Collection types (arrays, objects) with nested data types have complex diff suppression and state management logic. Special handling required for StateFunc behavior; full map-based collections cannot be saved outside pre-created maps; VECTOR data type partially unsupported in row access policy parsing
- Files:
  - `pkg/resources/data_type_handling_commons.go` - General data type handling with multiple TODOs on collection testing
  - `pkg/testfunctional/test_resource_data_type_diff_handling_list.go` - Force new on nested data types pattern
  - `pkg/resources/row_access_policy.go` - VECTOR data type parsing limitation
- Impact: Risk of state drift with complex nested data types; insufficient test coverage for collections; VECTOR type compatibility gap; difficult to maintain and extend data type support
- Fix approach: Expand functional test coverage for all collection type combinations; implement comprehensive VECTOR type support in row access policy parsing; refactor data type state handling to reduce special cases

**Deprecated Identifier Format Migration:**
- Issue: Legacy pipe-separated identifier format (`|`) still supported alongside new dot-separated format (`.`). State upgrader exists but identifier format migration is ongoing and incomplete.
- Files:
  - `pkg/resources/deprecated_identifier_helpers.go` - FormatFullyQualifiedObjectID supports both formats
  - `pkg/resources/identifier_state_upgraders.go` - migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName
  - Multiple resource state upgraders handling format conversion
- Impact: Legacy state parsing code must be maintained indefinitely; potential for format confusion in user configs; increased complexity in identifier comparison logic
- Fix approach: Set deprecation timeline for pipe-separated format; enforce fully-qualified name format in all new resources; remove legacy format support after deprecation period

## Known Bugs

**Row Access Policy Policy Unassignment Not Implemented (SNOW-1818849):**
- Symptoms: Dropping a row access policy may leave dangling policy assignments on tables/views if policies were manually assigned externally
- Files: `pkg/resources/row_access_policy.go:101`
- Trigger: Delete row access policy resource when external policy assignments exist
- Workaround: Manually unassign policies from tables/views before destroying resource, or use SQL ALTER TABLE/VIEW to detach policy
- Fix approach: Add pre-delete step to enumerate and drop all policy assignments before dropping the policy itself

**Grant Ownership Outbound Privileges Behavior Unclear (SNOW-1182623):**
- Symptoms: Comments in code indicate uncertainty about whether REVOKE or COPY should be used for outbound privileges in delete operation
- Files: `pkg/resources/grant_ownership.go:305`
- Trigger: Deleting grant_ownership resource
- Workaround: Currently set to COPY; may need manual cleanup if REVOKE intended
- Fix approach: Clarify expected behavior with Snowflake docs/support; implement configurable outbound_privileges handling in delete operation

## Security Considerations

**Secret Handling in OAuth Integrations:**
- Risk: OAuth client secrets and refresh tokens are handled as plain strings in Terraform state; no explicit masking or sensitive field marking visible in resource definitions
- Files:
  - `pkg/resources/api_authentication_integration_common.go` - oauth_client_secret, oauth_token_endpoint
  - `pkg/resources/secret_with_oauth_client_credentials.go` - oauth_scopes and api_authentication fields
  - `pkg/resources/secret_with_oauth_authorization_code_grant.go` - oauth_refresh_token and expiry times
- Current mitigation: Terraform state files themselves are protected; secrets are only written via API, not in configuration (except when explicitly provided by user)
- Recommendations:
  - Mark all OAuth secret fields as `Sensitive: true` in schema to prevent logging in debug output
  - Add warnings in documentation about state file protection importance
  - Consider implementing secret rotation capabilities

**Credential Validation Gaps:**
- Risk: Limited validation of credentials passed through API integrations; malformed credentials may only be detected at Snowflake API call time
- Files:
  - `pkg/resources/api_authentication_integration_common.go` - Basic string validation only
  - `pkg/resources/git_repository.go` - Secret reference validation
- Current mitigation: Snowflake API will reject invalid credentials
- Recommendations: Add client-side validation for common credential formats before API calls; add more detailed error messages when credential validation fails

**Grant Privilege Mapping Complexity:**
- Risk: Manual mapping of object types to grant functions with potential gaps; missing object type could silently fail
- Files: `pkg/resources/grant_privileges_to_database_role.go:951`, `pkg/resources/grant_privileges_to_account_role_identifier.go:214`
- Current mitigation: TBD marker for unmapped types
- Recommendations: Generate object type to parsing function mappings; add comprehensive test coverage for each object type; consider code generation from SDK definitions

## Performance Bottlenecks

**Large File Sizes with Complex Logic:**
- Problem: Several resources are extremely large with complex conditional logic making performance analysis difficult
- Files:
  - `pkg/testacc/resource_warehouse_acceptance_test.go` (3223 lines)
  - `pkg/testacc/resource_grant_privileges_to_account_role_acceptance_test.go` (3222 lines)
  - `pkg/testacc/resource_task_acceptance_test.go` (2681 lines)
- Cause: Acceptance tests with many scenario variations; resources with numerous optional parameters and special handling
- Improvement path:
  - Refactor large resources into separate focused components
  - Use test parameterization to reduce test file duplication
  - Consider modular resource design with shared behaviors

**Grant Privilege Iteration Complexity:**
- Problem: Grant operations iterate through large privilege lists with nested object type mappings; no visible caching
- Files: `pkg/resources/grant_privileges_to_account_role.go` (1294 lines)
- Cause: Manual privilege mapping for each object type; no pre-computed privilege matrix
- Improvement path: Cache privilege mappings at provider initialization; pre-compute object type compatibility; consider lazy-loading privilege definitions

**Parameter List Lookups:**
- Problem: Parameter lookups search through linear lists without indexing; called frequently during parameter operations
- Files: `pkg/resources/account_parameter.go`, `pkg/resources/account_parameters.go`, resource-specific parameter commons
- Cause: Parameters stored as slices with linear search
- Improvement path: Use parameter name as map key for O(1) lookup; cache parameter metadata at provider initialization

## Fragile Areas

**SAML2 Integration Update Logic (SNOW-1515781):**
- Files: `pkg/resources/saml2_integration.go` (831 lines)
- Why fragile: Multiple forced ForceNew conditions on optional fields due to UNSET unavailability; conditional logic spread across update function creates branching complexity; difficult to maintain when SDK adds UNSET support
- Safe modification: Any change must preserve existing ForceNew conditions; test with complete SAML2 config before and after updates; coordinate with UNSET implementation to remove workarounds
- Test coverage: Accept test likely incomplete due to multiple conditional paths and field combinations

**Grant Ownership Object Type Mapping:**
- Files: `pkg/resources/grant_ownership.go:405` - hardcoded object type mapping function
- Why fragile: Manual string-to-function mapping with no validation; new Snowflake object types require manual code update; mapping can diverge from SDK
- Safe modification: Any new object type addition must update mapping and add acceptance test; validate SDK changes do not add unmapped types
- Test coverage: Gaps for less common object types

**Data Type Diff Suppression:**
- Files: `pkg/resources/data_type_handling_commons.go` with collection handling logic in test files
- Why fragile: DataType equality determined by comparing parsed types with unknown value preservation; multiple code paths for different collection scenarios
- Safe modification: Changes to data type parsing require updating all affected resources and their tests; collection types need functional test coverage before deployment
- Test coverage: Many collection type combinations noted as TODO

**Identifier Validation Bypass for AccountObjectIdentifier:**
- Files: `pkg/internal/provider/validators/validators.go:48-51`
- Why fragile: Deliberate validation skip means identifier validation is incomplete; if SDK validation is removed, invalid identifiers may pass through
- Safe modification: Keep validation logic synchronized with SDK identifier validation; test with identifiers containing special characters and dots; revisit when identifier rework complete
- Test coverage: Insufficient tests for identifiers with special characters

## Scaling Limits

**Acceptance Test Suite Size:**
- Current capacity: 3200+ line acceptance test files for individual resources
- Limit: Tests become difficult to understand and maintain; test execution time grows; developer experience degrades
- Scaling path:
  - Split large test files by scenario (basic, complete, edges)
  - Use test parameterization to reduce duplication
  - Implement test data factory patterns for common setup

**Grant Privilege Scope:**
- Current capacity: Account and database role grants with ~20 privilege types across ~40 object types
- Limit: Manual privilege mapping will not scale to account roles with thousands of privileges; nested privilege hierarchies not fully modeled
- Scaling path: Move privilege definitions to SDK; generate privilege handling code; implement privilege matrix caching

**Parameter Handling Across Resources:**
- Current capacity: 40+ account parameters split across two resources; resource-specific parameters multiply this across 100+ resources
- Limit: Parameter list maintenance becomes unmanageable; parameter SDK updates require multiple resource updates
- Scaling path: Implement parameter handling as plugin/registration system; consolidate parameter definitions in SDK; generate resource parameter schema from SDK

## Dependencies at Risk

**Identifier Helpers Library (internal only):**
- Risk: Two competing identifier implementations exist; state upgraders manage transitions but format incompatibility risks remain
- Impact: Any breaking changes to identifier format require state migration; existing state may become unrecoverable if migration tools removed
- Migration plan: Consolidate to single identifier implementation after identifier rework (SNOW-1634872) complete

**Data Type Parsing Complexity:**
- Risk: Custom data type parser in `pkg/sdk/datatypes` handles parsing but equality and diff logic is in resource layer; parser bugs affect multiple resources
- Impact: Data type changes may not be properly detected; state drift with collections; VECTOR type gaps
- Migration plan: Consolidate data type handling logic into SDK; implement comprehensive data type equality testing

## Missing Critical Features

**Row Access Policy Unassignment:**
- Problem: Cannot automatically remove policy assignments when policy is deleted
- Blocks: Clean policy lifecycle management; requires manual external cleanup
- Priority: Medium - affects operational safety

**UNSET Clause Support:**
- Problem: Partial parameter removal requires resource recreation
- Blocks: Clean parameter updates for all integration types; forces unnecessary downtime
- Priority: High - affects multiple resource types (SAML2, OAuth, APIs)

**Parameter Unification:**
- Problem: Parameters handled inconsistently across account_parameter, account_parameters, and data sources
- Blocks: Consistent parameter management; single source of truth for supported parameters
- Priority: Medium - affects maintainability and user experience

**VECTOR Data Type Support:**
- Problem: Row access policy cannot parse functions with VECTOR arguments
- Blocks: Using modern vector search features in policies
- Priority: Low - relatively new feature; affects advanced use cases

## Test Coverage Gaps

**AccountObjectIdentifier Validation:**
- What's not tested: Identifiers with special characters, case sensitivity, dot handling in account names
- Files: `pkg/internal/provider/validators/validators.go` - IsValidIdentifier validation bypass for AccountObjectIdentifier
- Risk: Invalid identifiers accepted; validation divergence if SDK changes validation logic
- Priority: High

**SAML2 Integration Optional Field Combinations:**
- What's not tested: All combinations of optional field updates that trigger ForceNew; UNSET scenarios when implemented
- Files: `pkg/resources/saml2_integration.go` - Update function with 9+ conditional ForceNew branches
- Risk: Undiscovered state transitions that cause unexpected resource recreation
- Priority: High

**Collection Data Type Diff Handling:**
- What's not tested: Full matrix of collection type combinations; external modifications to nested collections; StateFunc behavior with unknowns
- Files: `pkg/resources/data_type_handling_commons.go`, `pkg/testfunctional/test_resource_data_type_diff_handling_list.go`
- Risk: State drift with complex nested types; undetected external changes
- Priority: Medium - affects feature stability

**Grant Privilege Object Type Coverage:**
- What's not tested: All 40+ object types with all privilege combinations; edge cases with cross-type privileges
- Files: `pkg/resources/grant_privileges_to_account_role.go`, `pkg/resources/grant_privileges_to_database_role.go`
- Risk: Privilege mappings fail silently for untested combinations; manual mapping errors not caught
- Priority: Medium - affects security configuration reliability

**Parameter Validation for Each Resource Type:**
- What's not tested: Valid/invalid parameter combinations per resource type; parameter interactions; parameter change propagation
- Files: Multiple `*_parameters.go` resource files with parameter schema definitions
- Risk: Invalid parameter combinations accepted; parameter dependencies not enforced
- Priority: Low - generally well-structured but gaps exist

---

*Concerns audit: 2026-02-25*

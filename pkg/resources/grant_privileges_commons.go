package resources

import (
	"context"
	"fmt"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// inheritedGrantsRequireExperiment returns a CustomizeDiff that fails the plan when an `inherited`
// block is used in any of the given blocks while the INHERITED_GRANTS experiment is not enabled.
func inheritedGrantsRequireExperiment(blockNames ...string) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
		rawConfig := d.GetRawConfig()
		if rawConfig.IsNull() || !rawConfig.IsKnown() {
			return nil
		}
		if !rawConfigHasInheritedGrant(rawConfig, blockNames...) {
			return nil
		}
		providerCtx := meta.(*provider.Context)
		if !experimentalfeatures.IsExperimentEnabled(experimentalfeatures.InheritedGrants, providerCtx.EnabledExperiments) {
			return fmt.Errorf("using an `inherited` block requires the %q experiment to be enabled. Add it to the `experimental_features_enabled` list in the provider configuration", experimentalfeatures.InheritedGrants)
		}
		return nil
	}
}

// validateInheritedGrantsConfig returns a ValidateRawResourceConfigFunc that rejects `with_grant_option`
// and `always_apply` when used together with an `inherited` block in any of the given blocks.
func validateInheritedGrantsConfig(blockNames ...string) schema.ValidateRawResourceConfigFunc {
	return func(ctx context.Context, req schema.ValidateResourceConfigFuncRequest, resp *schema.ValidateResourceConfigFuncResponse) {
		if !rawConfigHasInheritedGrant(req.RawConfig, blockNames...) {
			return
		}
		// Inherited grants do not support the WITH GRANT OPTION clause.
		if withGrantOption := req.RawConfig.GetAttr("with_grant_option"); !withGrantOption.IsNull() && withGrantOption.IsKnown() && withGrantOption.True() {
			resp.Diagnostics = append(resp.Diagnostics, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid with_grant_option",
				Detail:   "`with_grant_option` cannot be used together with an `inherited` block, because inherited grants do not support the WITH GRANT OPTION clause.",
			})
		}
		// always_apply re-grants the configured privileges on every apply. Inherited grants already
		// cover all current and future objects in the container, so re-granting serves no purpose and
		// is disallowed to keep the resource behavior clear.
		if alwaysApply := req.RawConfig.GetAttr("always_apply"); !alwaysApply.IsNull() && alwaysApply.IsKnown() && alwaysApply.True() {
			resp.Diagnostics = append(resp.Diagnostics, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid always_apply",
				Detail:   "`always_apply` cannot be used together with an `inherited` block. Inherited grants already cover all current and future objects in the container, so re-granting on every apply is unnecessary.",
			})
		}
	}
}

// rawConfigHasInheritedGrant reports whether any of the given blocks configures an `inherited` grant in
// the given raw configuration.
func rawConfigHasInheritedGrant(rawConfig cty.Value, blockNames ...string) bool {
	for _, blockName := range blockNames {
		block := rawConfig.GetAttr(blockName)
		if block.IsNull() || !block.IsKnown() || block.LengthInt() == 0 {
			continue
		}
		for _, elem := range block.AsValueSlice() {
			if elem.IsNull() || !elem.IsKnown() {
				continue
			}
			inherited := elem.GetAttr("inherited")
			if inherited.IsNull() || !inherited.IsKnown() {
				continue
			}
			if inherited.Type() == cty.String {
				if inherited.AsString() != "" {
					return true
				}
				continue
			}
			if inherited.LengthInt() > 0 {
				return true
			}
		}
	}
	return false
}

// computeInheritedPrivileges filters the grants returned by SHOW GRANTS TO [DATABASE] ROLE down to the
// inherited grants (is_inherited = true) that are granted to the given role and match the object type and
// container described by the identifier data, using the inherited-specific columns (inherited_from,
// inherited_from_database, inherited_from_schema).
func computeInheritedPrivileges(data fmt.Stringer, roleName string, granteeType sdk.ObjectType, expectedPrivileges []string, grants []sdk.Grant, strictPrivilegeManagement bool) (actualPrivileges []string) {
	grantedOn, container, database, schema := inheritedGrantScopeFromData(data)

	for _, grant := range grants {
		if grant.IsInherited != nil && !*grant.IsInherited {
			continue
		}
		if (grant.GrantTo != granteeType && grant.GrantedTo != granteeType) || grant.GranteeName.Name() != roleName {
			continue
		}
		if grant.GrantedOn != grantedOn || !inheritedGrantMatchesContainer(grant, container, database, schema) {
			continue
		}
		if !strictPrivilegeManagement && !slices.Contains(expectedPrivileges, grant.Privilege) {
			continue
		}
		actualPrivileges = append(actualPrivileges, grant.Privilege)
	}

	return actualPrivileges
}

// inheritedGrantScopeFromData returns the object type and container an inherited grant's identifier data
// is scoped to, used to match the rows returned by SHOW GRANTS TO [DATABASE] ROLE.
func inheritedGrantScopeFromData(data fmt.Stringer) (grantedOn sdk.ObjectType, container InheritedContainerKind, database *sdk.AccountObjectIdentifier, schema *sdk.DatabaseObjectIdentifier) {
	switch d := data.(type) {
	case *OnAccountObjectInheritedGrantData:
		return d.ObjectNamePlural.Singular(), InAccountInheritedContainerKind, nil, nil
	case *OnSchemaInheritedGrantData:
		return sdk.ObjectTypeSchema, d.Kind, d.DatabaseName, nil
	case *OnSchemaObjectInheritedGrantData:
		return d.ObjectNamePlural.Singular(), d.Kind, d.DatabaseName, d.SchemaName
	default:
		return "", "", nil, nil
	}
}

// inheritedGrantMatchesContainer reports whether the inherited grant row comes from the same
// container (account, database, or schema) the resource is scoped to, based on the inherited_from,
// inherited_from_database, and inherited_from_schema columns.
func inheritedGrantMatchesContainer(grant sdk.Grant, container InheritedContainerKind, database *sdk.AccountObjectIdentifier, schema *sdk.DatabaseObjectIdentifier) bool {
	if grant.InheritedFrom == nil {
		return false
	}
	switch container {
	case InAccountInheritedContainerKind:
		return *grant.InheritedFrom == sdk.GrantInheritedFromAccount
	case InDatabaseInheritedContainerKind:
		return *grant.InheritedFrom == sdk.GrantInheritedFromDatabase &&
			database != nil && grant.InheritedFromDatabase != nil &&
			*grant.InheritedFromDatabase == database.Name()
	case InSchemaInheritedContainerKind:
		return *grant.InheritedFrom == sdk.GrantInheritedFromSchema &&
			schema != nil && grant.InheritedFromDatabase != nil && grant.InheritedFromSchema != nil &&
			*grant.InheritedFromDatabase == schema.DatabaseName() &&
			*grant.InheritedFromSchema == schema.Name()
	default:
		return false
	}
}

// getInheritedGrantContainer reads the in_account / in_database / in_schema attributes of an `inherited`
// block and returns the resolved container kind together with the parsed database/schema identifiers.
func getInheritedGrantContainer(inherited map[string]any) (InheritedContainerKind, *sdk.AccountObjectIdentifier, *sdk.DatabaseObjectIdentifier, error) {
	if inAccount, ok := inherited["in_account"].(bool); ok && inAccount {
		return InAccountInheritedContainerKind, nil, nil, nil
	}
	if inDatabase, ok := inherited["in_database"].(string); ok && len(inDatabase) > 0 {
		databaseId, err := sdk.ParseAccountObjectIdentifier(inDatabase)
		if err != nil {
			return "", nil, nil, err
		}
		return InDatabaseInheritedContainerKind, new(databaseId), nil, nil
	}
	if inSchema, ok := inherited["in_schema"].(string); ok && len(inSchema) > 0 {
		schemaId, err := sdk.ParseDatabaseObjectIdentifier(inSchema)
		if err != nil {
			return "", nil, nil, err
		}
		return InSchemaInheritedContainerKind, nil, new(schemaId), nil
	}
	return "", nil, nil, fmt.Errorf("inherited block must specify exactly one of in_account, in_database, or in_schema")
}

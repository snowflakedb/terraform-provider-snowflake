package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// oauthAllowedRolesListRequiresSecondaryRolesNone validates that allowed_roles_list can
// only be set when oauth_use_secondary_roles is NONE.
func oauthAllowedRolesListRequiresSecondaryRolesNone(_ context.Context, req schema.ValidateResourceConfigFuncRequest, resp *schema.ValidateResourceConfigFuncResponse) {
	rawConfig := req.RawConfig
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return
	}

	allowedRoles := rawConfig.GetAttr("allowed_roles_list")
	if allowedRoles.IsNull() || !allowedRoles.IsKnown() || allowedRoles.LengthInt() == 0 {
		return
	}

	secondaryRoles := rawConfig.GetAttr("oauth_use_secondary_roles")
	if secondaryRoles.IsNull() || !secondaryRoles.IsKnown() {
		return
	}
	if secondaryRoles.AsString() != string(sdk.OauthSecurityIntegrationUseSecondaryRolesOptionNone) {
		resp.Diagnostics = append(resp.Diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid allowed_roles_list",
			Detail:   fmt.Sprintf("allowed_roles_list can only be set when oauth_use_secondary_roles is set to %s", sdk.OauthSecurityIntegrationUseSecondaryRolesOptionNone),
		})
	}
}

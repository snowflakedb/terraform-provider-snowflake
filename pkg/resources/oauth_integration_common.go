package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// oauthAllowedRolesListRequiresSecondaryRolesNone validates that allowed_roles_list can
// only be set when oauth_use_secondary_roles is NONE.
func oauthAllowedRolesListRequiresSecondaryRolesNone(_ context.Context, d *schema.ResourceDiff, _ any) error {
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() {
		return nil
	}
	configMap := rawConfig.AsValueMap()

	allowedRoles, ok := configMap["allowed_roles_list"]
	if !ok || allowedRoles.IsNull() || !allowedRoles.IsKnown() || allowedRoles.LengthInt() == 0 {
		return nil
	}

	secondaryRoles, ok := configMap["oauth_use_secondary_roles"]
	if !ok || secondaryRoles.IsNull() || !secondaryRoles.IsKnown() {
		return nil
	}
	if secondaryRoles.AsString() != string(sdk.OauthSecurityIntegrationUseSecondaryRolesOptionNone) {
		return fmt.Errorf("allowed_roles_list can only be set when oauth_use_secondary_roles is set to %s", sdk.OauthSecurityIntegrationUseSecondaryRolesOptionNone)
	}
	return nil
}

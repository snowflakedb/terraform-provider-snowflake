package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func v094UserStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	oldDefaultSecondaryRoles := rawState["default_secondary_roles"]
	if oldDefaultSecondaryRoles == nil {
		rawState["default_secondary_roles_option"] = string(sdk.SecondaryRolesOptionDefault)
	} else {
		if len(oldDefaultSecondaryRoles.([]any)) > 0 {
			rawState["default_secondary_roles_option"] = string(sdk.SecondaryRolesOptionAll)
		} else {
			rawState["default_secondary_roles_option"] = string(sdk.SecondaryRolesOptionNone)
		}
	}
	delete(rawState, "default_secondary_roles")

	return rawState, nil
}

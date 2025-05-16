package resources

import (
	"context"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func v092ScimIntegrationStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	rawState["run_as_role"] = rawState["provisioner_role"]
	delete(rawState, "provisioner_role")

	if v, ok := rawState["enabled"]; ok {
		rawState["enabled"] = strconv.FormatBool(v.(bool))
	}

	return rawState, nil
}

func v093ScimIntegrationStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	if v, ok := rawState["scim_client"]; ok && strings.EqualFold(strings.TrimSpace(v.(string)), string(sdk.ScimSecurityIntegrationScimClientAzure)) {
		rawState["sync_password"] = BooleanDefault
	}

	return rawState, nil
}

package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func v2_9_0_AuthenticationPolicyStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	idRaw, ok := rawState["id"].(string)
	if !ok {
		return nil, fmt.Errorf("cannot read authentication policy id from the state")
	}
	id, err := sdk.ParseSchemaObjectIdentifier(idRaw)
	if err != nil {
		return nil, err
	}

	client := meta.(*provider.Context).Client
	authenticationPolicyInSnowflake, err := client.AuthenticationPolicies.ShowByIDSafely(ctx, id)
	if err != nil {
		return nil, err
	}
	showOutput := schemas.AuthenticationPolicyToSchema(authenticationPolicyInSnowflake)
	rawState[ShowOutputAttributeName] = []any{showOutput}
	delete(rawState, "mfa_authentication_methods")

	describeOutput, err := client.AuthenticationPolicies.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	describeOutputState := schemas.AuthenticationPolicyDescriptionToSchema(describeOutput)
	rawState[DescribeOutputAttributeName] = []any{describeOutputState}

	return rawState, nil
}

package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func v2_15_0_PasswordPolicyStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	rawState, err := migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName(ctx, rawState, meta)
	if err != nil {
		return nil, err
	}

	// Remove deprecated fields from state
	delete(rawState, "or_replace")
	delete(rawState, "if_not_exists")

	idRaw, ok := rawState["id"].(string)
	if !ok {
		return nil, fmt.Errorf("cannot read password policy id from the state")
	}
	id, err := sdk.ParseSchemaObjectIdentifier(idRaw)
	if err != nil {
		return nil, err
	}

	client := meta.(*provider.Context).Client

	passwordPolicy, err := client.PasswordPolicies.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			// Object was deleted externally; Read will handle cleanup
			return rawState, nil
		}
		return nil, err
	}
	rawState[ShowOutputAttributeName] = []any{schemas.PasswordPolicyToSchema(passwordPolicy)}

	details, err := client.PasswordPolicies.DescribeDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	rawState[DescribeOutputAttributeName] = []any{schemas.PasswordPolicyDetailsToSchema(details)}

	return rawState, nil
}

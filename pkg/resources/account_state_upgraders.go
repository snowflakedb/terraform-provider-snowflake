package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func v0_99_0_AccountStateUpgrader(ctx context.Context, state map[string]any, meta any) (map[string]any, error) {
	if state == nil {
		return state, nil
	}

	client := meta.(*provider.Context).Client
	if v, ok := state["must_change_password"]; ok && v != nil {
		state["must_change_password"] = booleanStringFromBool(v.(bool))
	}
	if v, ok := state["is_org_admin"]; ok && v != nil {
		state["is_org_admin"] = booleanStringFromBool(v.(bool))
	}
	account, err := client.Accounts.ShowByID(ctx, sdk.NewAccountObjectIdentifier(state["name"].(string)))
	if err != nil {
		return nil, err
	}

	state["id"] = helpers.EncodeResourceIdentifier(sdk.NewAccountIdentifier(account.OrganizationName, account.AccountName))

	return state, nil
}

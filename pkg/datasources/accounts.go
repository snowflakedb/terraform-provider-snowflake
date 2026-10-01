package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountsSchema = map[string]*schema.Schema{
	"with_history": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Includes dropped accounts that have not yet been deleted.",
	},
	"with_parameters": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs SHOW PARAMETERS IN ACCOUNT once, for the connected account, and saves the output on that account's row only. Other listed accounts stay empty because Snowflake cannot show parameters for another organization account; the call is not repeated or skipped per row. By default this value is set to true.",
	},
	"like": likeSchema,
	"accounts": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all accounts details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW ACCOUNTS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowAccountSchema,
					},
				},
				resources.ParametersAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW PARAMETERS IN ACCOUNT for the connected account.",
					Elem: &schema.Resource{
						Schema: schemas.ShowAccountParametersSchema,
					},
				},
			},
		},
	},
}

func Accounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.Accounts, ReadAccounts),
		Schema:      accountsSchema,
		Description: "Data source used to get details of filtered accounts. Filtering is aligned with the current possibilities for [SHOW ACCOUNTS](https://docs.snowflake.com/en/sql-reference/sql/show-accounts) query. The results of SHOW and SHOW PARAMETERS IN ACCOUNT (connected account only) are encapsulated in one output collection `accounts`.",
	}
}

func ReadAccounts(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client

	req := sdk.NewShowAccountRequest()
	if likePattern, ok := d.GetOk("like"); ok {
		req = req.WithLike(sdk.Like{Pattern: sdk.String(likePattern.(string))})
	}
	if history, ok := d.GetOk("with_history"); ok && history.(bool) {
		req = req.WithHistory(true)
	}

	accounts, err := client.Accounts.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("accounts")

	var currentAccountParameters []map[string]any
	var currentAccountId sdk.AccountIdentifier
	if d.Get("with_parameters").(bool) {
		sessionDetails, err := client.ContextFunctions.CurrentSessionDetails(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		currentAccountId = sdk.NewAccountIdentifier(sessionDetails.OrganizationName, sessionDetails.AccountName)

		parameters, err := client.Accounts.ShowParameters(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		currentAccountParameters = []map[string]any{schemas.AccountParametersToSchema(parameters, providerCtx)}
	}

	flattenedAccounts := make([]map[string]any, len(accounts))
	for i, account := range accounts {
		var accountParameters []map[string]any
		if d.Get("with_parameters").(bool) && account.AccountID().FullyQualifiedName() == currentAccountId.FullyQualifiedName() {
			accountParameters = currentAccountParameters
		}
		flattenedAccounts[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.AccountToSchema(&account)},
			resources.ParametersAttributeName: accountParameters,
		}
	}

	if err := d.Set("accounts", flattenedAccounts); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

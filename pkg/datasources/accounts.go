package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountsSchema = map[string]*schema.Schema{
	"with_history": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Includes dropped accounts that have not yet been deleted.",
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
				// TODO(SNOW-1348092 - next prs): Add parameters
			},
		},
	},
}

func Accounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.Accounts, ReadAccounts),
		Schema:      accountsSchema,
		Description: "Data source used to get details of filtered accounts. Filtering is aligned with the current possibilities for [SHOW ACCOUNTS](https://docs.snowflake.com/en/sql-reference/sql/show-accounts) query. The results of SHOW are encapsulated in one output collection `accounts`.",
	}
}

func ReadAccounts(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	req := new(sdk.ShowAccountOptions)
	handleLike(d, &req.Like)
	if history, ok := d.GetOk("with_history"); ok && history.(bool) {
		req.History = sdk.Bool(true)
	}

	accounts, err := client.Accounts.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("accounts")

	flattenedAccounts := make([]map[string]any, len(accounts))
	for i, account := range accounts {
		account := account
		flattenedAccounts[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.AccountToSchema(&account)},
		}
	}

	if err := d.Set("accounts", flattenedAccounts); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

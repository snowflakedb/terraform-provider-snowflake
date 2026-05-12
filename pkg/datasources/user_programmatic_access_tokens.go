package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var userProgrammaticAccessTokensSchema = map[string]*schema.Schema{
	"for_user": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Returns programmatic access tokens for the specified user.",
		ValidateDiagFunc: resources.IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"user_programmatic_access_tokens": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all user programmatic access tokens details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW USER PROGRAMMATIC ACCESS TOKENS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowProgrammaticAccessTokenSchema,
					},
				},
			},
		},
	},
}

func UserProgrammaticAccessTokens() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.UserProgrammaticAccessTokens, ReadUserProgrammaticAccessTokens),
		Schema:      userProgrammaticAccessTokensSchema,
		Description: "Data source used to get details of filtered user programmatic access tokens. Filtering is aligned with the current possibilities for [SHOW USER PROGRAMMATIC ACCESS TOKENS](https://docs.snowflake.com/en/sql-reference/sql/show-user-programmatic-access-tokens) query. The results of SHOW are encapsulated in one output collection `user_programmatic_access_tokens`.",
	}
}

func ReadUserProgrammaticAccessTokens(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	userName := d.Get("for_user").(string)
	userId, err := sdk.ParseAccountObjectIdentifier(userName)
	if err != nil {
		return diag.FromErr(err)
	}
	req := sdk.NewShowUserProgrammaticAccessTokenRequest().WithUserName(userId)

	userProgrammaticAccessTokens, err := client.UserProgrammaticAccessTokens.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("user_programmatic_access_tokens_read")

	flattenedUserProgrammaticAccessTokens := make([]map[string]any, len(userProgrammaticAccessTokens))
	for i, userProgrammaticAccessToken := range userProgrammaticAccessTokens {
		flattenedUserProgrammaticAccessTokens[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.ProgrammaticAccessTokenToSchema(&userProgrammaticAccessToken)},
		}
	}
	if err := d.Set("user_programmatic_access_tokens", flattenedUserProgrammaticAccessTokens); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

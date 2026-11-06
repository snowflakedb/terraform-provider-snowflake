package datasources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var sharesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC SHARE for each share returned by SHOW SHARES. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"shares": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all shares details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW SHARES. SHOW SHARES lists at most three accounts of an outbound share, so for a share with more accounts the full list of `to` is read from `SNOWFLAKE.ACCOUNT_USAGE.SHARES` (it needs `SELECT` on that view and a warehouse in the connection, and the view can be up to 3 hours behind). If that fails, a warning is reported and `to` holds the accounts listed by SHOW SHARES. `to_total` is the number of accounts reported by SHOW SHARES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowShareSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE SHARE: a row for every object granted to the share, including the system objects that Snowflake adds to the shared database. A share that cannot be described (for example because of missing privileges) is reported with a warning and has an empty output.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeShareDetailsSchema,
					},
				},
			},
		},
	},
}

func Shares() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.SharesDatasource), TrackingReadWrapper(datasources.Shares, ReadShares)),
		Schema:      sharesSchema,
		Description: docs.WithPreviewFeatureMark() + "Data source used to get details of filtered shares. Filtering is aligned with the current possibilities for [SHOW SHARES](https://docs.snowflake.com/en/sql-reference/sql/show-shares) query." +
			" The results of SHOW and DESCRIBE are encapsulated in one output collection `shares`.",
	}
}

func ReadShares(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowShareRequest()

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)

	shares, err := client.Shares.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("shares_read")

	var diags diag.Diagnostics
	flattenedShares := make([]map[string]any, len(shares))
	for i := range shares {
		share := shares[i]
		if err := share.CompleteTo(ctx, client.Shares); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Share %s has %d recipients but SHOW SHARES lists only %d", share.ID().FullyQualifiedName(), share.ToTotal, len(share.To)),
				Detail: "The full list is read from SNOWFLAKE.ACCOUNT_USAGE.SHARES, which failed (the role needs SELECT on that view, for example through IMPORTED PRIVILEGES on the SNOWFLAKE database, and the connection needs a warehouse). " +
					"Recipients that are not listed are missing from the output. Error: " + err.Error(),
			})
		}
		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := describeShare(ctx, client, &share)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Failed to describe share %s", share.Name),
					Detail:   fmt.Sprintf("The describe_output of the share is empty. Error: %s", err),
				})
			} else {
				describeOutput = schemas.ShareDetailsListToSchema(details)
			}
		}
		flattenedShares[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.ShareToSchema(&share)},
			resources.DescribeOutputAttributeName: describeOutput,
		}
	}
	if err := d.Set("shares", flattenedShares); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

func describeShare(ctx context.Context, client *sdk.Client, share *sdk.Share) ([]sdk.ShareDetails, error) {
	if share.Kind == sdk.ShareKindInbound {
		return client.Shares.DescribeConsumer(ctx, share.ExternalID())
	}
	return client.Shares.DescribeProvider(ctx, share.ID())
}

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

var codeBundlesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC CODE BUNDLE for each code bundle returned by SHOW CODE BUNDLES. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"in": extendedInSchema,
	"code_bundles": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all code bundle details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW CODE BUNDLES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowCodeBundleSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE CODE BUNDLE.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeCodeBundleDetailsSchema,
					},
				},
			},
		},
	},
}

func CodeBundles() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.CodeBundles, ReadCodeBundles),
		Schema:      codeBundlesSchema,
		Description: "Data source used to get details of filtered code bundles. Filtering is aligned with the current possibilities for [SHOW CODE BUNDLES](https://docs.snowflake.com/en/sql-reference/sql/show-code-bundles) query. " +
			"The results of SHOW and DESCRIBE are encapsulated in one output collection `code_bundles`.",
	}
}

func ReadCodeBundles(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowCodeBundleRequest{}

	if err := handleExtendedIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}

	codeBundles, err := client.CodeBundles.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("code_bundles_read")

	flattened := make([]map[string]any, len(codeBundles))
	for i := range codeBundles {
		codeBundle := codeBundles[i]
		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.CodeBundles.DescribeDetails(ctx, codeBundle.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOutput = []map[string]any{schemas.CodeBundleDetailsToSchema(details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.CodeBundleToSchema(&codeBundle)},
			resources.DescribeOutputAttributeName: describeOutput,
		}
	}
	if err := d.Set("code_bundles", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

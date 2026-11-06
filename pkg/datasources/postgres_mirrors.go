package datasources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var postgresMirrorsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESCRIBE_MIRROR for each mirror returned by LIST_MIRRORS. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"postgres_instance": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Lists only the mirrors that replicate from this Postgres instance. The name is case-sensitive. A quoted identifier, such as the instance's `fully_qualified_name`, is also accepted.",
	},
	"postgres_mirrors": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all Postgres mirror details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of LIST_MIRRORS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowPostgresMirrorSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE_MIRROR.",
					Elem: &schema.Resource{
						Schema: schemas.DescribePostgresMirrorDetailsSchema,
					},
				},
			},
		},
	},
}

func PostgresMirrors() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.PostgresMirrorsDatasource), TrackingReadWrapper(datasources.PostgresMirrors, ReadPostgresMirrors)),
		Schema:      postgresMirrorsSchema,
		Description: "Data source used to get details of filtered Postgres mirrors. Filtering is aligned with the current possibilities for the [LIST_MIRRORS](https://docs.snowflake.com/en/user-guide/snowflake-postgres/postgres-data-mirroring-reference) procedure. " +
			"The results of LIST_MIRRORS and DESCRIBE_MIRROR are encapsulated in one output collection `postgres_mirrors`.",
	}
}

func ReadPostgresMirrors(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	arguments := sdk.NewListPostgresMirrorArgumentsRequest()
	if v, ok := d.GetOk("postgres_instance"); ok {
		instanceId, err := sdk.ParseAccountObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		arguments.WithPostgresInstance(instanceId.Name())
	}

	mirrors, err := client.PostgresMirrors.List(ctx, sdk.NewListPostgresMirrorRequest(*arguments))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("postgres_mirrors_read")

	var diags diag.Diagnostics
	flattened := make([]map[string]any, len(mirrors))
	for i := range mirrors {
		mirror := mirrors[i]
		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.PostgresMirrors.Describe(ctx, sdk.NewDescribePostgresMirrorRequest(*sdk.NewDescribePostgresMirrorArgumentsRequest(mirror.Name)))
			if err != nil {
				return diag.FromErr(err)
			}
			if !details.IsMembershipKnown() {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Could not read the tables and schemas of the Postgres mirror.",
					Detail: fmt.Sprintf("DESCRIBE_MIRROR returned no tables or schemas for Postgres mirror %s, either because the source Postgres instance could not be reached or because the SNOWFLAKE application in this account does not report them yet. "+
						"describe_output shows its postgres_tables and postgres_schemas as empty, although the actual tables and schemas are unknown.", mirror.Name),
				})
			}
			describeOutput = []map[string]any{schemas.PostgresMirrorDetailsToSchema(details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.PostgresMirrorToSchema(&mirror)},
			resources.DescribeOutputAttributeName: describeOutput,
		}
	}
	if err := d.Set("postgres_mirrors", flattened); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

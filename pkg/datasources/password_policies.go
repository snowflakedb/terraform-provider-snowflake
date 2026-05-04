package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var passwordPoliciesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC PASSWORD POLICY for each password policy returned by SHOW PASSWORD POLICIES. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like":  likeSchema,
	"in":    inSchema,
	"limit": limitFromSchema,
	"password_policies": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all password policy details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW PASSWORD POLICIES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowPasswordPolicySchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE PASSWORD POLICY.",
					Elem: &schema.Resource{
						Schema: schemas.DescribePasswordPolicyDetailsSchema,
					},
				},
			},
		},
	},
}

func PasswordPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.PasswordPoliciesDatasource), TrackingReadWrapper(datasources.PasswordPolicies, ReadPasswordPolicies)),
		Schema:      passwordPoliciesSchema,
		Description: "Data source used to get details of filtered password policies. Filtering is aligned with the current possibilities for [SHOW PASSWORD POLICIES](https://docs.snowflake.com/en/sql-reference/sql/show-password-policies) query." +
			" The results of SHOW and DESCRIBE are encapsulated in one output collection `password_policies`.",
	}
}

func ReadPasswordPolicies(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowPasswordPolicyRequest{}

	handleLike(d, &req.Like)
	handleLimitFrom(d, &req.Limit)
	if err := handleIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}

	passwordPolicies, err := client.PasswordPolicies.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("password_policies_read")

	flattened := make([]map[string]any, len(passwordPolicies))
	for i := range passwordPolicies {
		pp := passwordPolicies[i]
		var describeOut []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.PasswordPolicies.DescribeDetails(ctx, pp.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOut = []map[string]any{schemas.PasswordPolicyDetailsToSchema(details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.PasswordPolicyToSchema(&pp)},
			resources.DescribeOutputAttributeName: describeOut,
		}
	}
	if err := d.Set("password_policies", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

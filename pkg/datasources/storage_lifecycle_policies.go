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

var storageLifecyclePoliciesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC STORAGE LIFECYCLE POLICY for each storage lifecycle policy returned by SHOW STORAGE LIFECYCLE POLICIES. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like": likeSchema,
	"in":   inSchema,
	"storage_lifecycle_policies": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all storage lifecycle policy details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW STORAGE LIFECYCLE POLICIES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowStorageLifecyclePolicySchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE STORAGE LIFECYCLE POLICY.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeStorageLifecyclePolicyDetailsSchema,
					},
				},
			},
		},
	},
}

func StorageLifecyclePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.StorageLifecyclePoliciesDatasource), TrackingReadWrapper(datasources.StorageLifecyclePolicies, ReadStorageLifecyclePolicies)),
		Schema:      storageLifecyclePoliciesSchema,
		Description: "Data source used to get details of filtered storage lifecycle policies. Filtering is aligned with the current possibilities for [SHOW STORAGE LIFECYCLE POLICIES](https://docs.snowflake.com/en/sql-reference/sql/show-storage-lifecycle-policies) query." +
			" The results of SHOW and DESCRIBE are encapsulated in one output collection `storage_lifecycle_policies`.",
	}
}

func ReadStorageLifecyclePolicies(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowStorageLifecyclePolicyRequest{}

	handleLike(d, &req.Like)
	if err := handleIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}

	storageLifecyclePolicies, err := client.StorageLifecyclePolicies.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("storage_lifecycle_policies_read")

	flattened := make([]map[string]any, len(storageLifecyclePolicies))
	for i := range storageLifecyclePolicies {
		policy := storageLifecyclePolicies[i]
		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.StorageLifecyclePolicies.Describe(ctx, policy.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOutput = []map[string]any{schemas.StorageLifecyclePolicyDetailsToSchema(*details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.StorageLifecyclePolicyToSchema(&policy)},
			resources.DescribeOutputAttributeName: describeOutput,
		}
	}
	if err := d.Set("storage_lifecycle_policies", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

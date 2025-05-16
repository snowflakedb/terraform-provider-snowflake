package datasources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountRolesSchema = map[string]*schema.Schema{
	"like": likeSchema,
	"in_class": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: resources.IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		Description:      "Filters the SHOW GRANTS output by class name.",
	},
	"account_roles": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all account role details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW ROLES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowRoleSchema,
					},
				},
			},
		},
	},
}

func AccountRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.AccountRoles, ReadAccountRoles),
		Schema:      accountRolesSchema,
		Description: "Data source used to get details of filtered account roles. Filtering is aligned with the current possibilities for [SHOW ROLES](https://docs.snowflake.com/en/sql-reference/sql/show-roles) query (`like` and `in_class` are all supported). The results of SHOW are encapsulated in one output collection.",
	}
}

func ReadAccountRoles(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	req := sdk.NewShowRoleRequest()

	handleLike(d, &req.Like)

	if className, ok := d.GetOk("in_class"); ok {
		req.WithInClass(sdk.RolesInClass{
			Class: sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(className.(string)),
		})
	}

	roles, err := client.Roles.Show(ctx, req)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to show account roles",
				Detail:   fmt.Sprintf("Error: %s", err),
			},
		}
	}

	d.SetId("account_roles_read")

	flattenedAccountRoles := make([]map[string]any, len(roles))
	for i, role := range roles {
		role := role
		flattenedAccountRoles[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.RoleToSchema(&role)},
		}
	}

	err = d.Set("account_roles", flattenedAccountRoles)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

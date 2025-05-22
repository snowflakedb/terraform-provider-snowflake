package datasources

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var currentRoleSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The name of the [primary role](https://docs.snowflake.com/en/user-guide/security-access-control-overview.html#label-access-control-role-enforcement) in use for the current session.",
	},
}

func CurrentRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.CurrentRoleDatasource), TrackingReadWrapper(datasources.CurrentRole, ReadCurrentRole)),
		Schema:      currentRoleSchema,
	}
}

func ReadCurrentRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	role, err := client.ContextFunctions.CurrentRole(ctx)
	if err != nil {
		log.Printf("[DEBUG] current_role failed to decode")
		d.SetId("")
		return nil
	}

	d.SetId(helpers.EncodeSnowflakeID(role))
	err = d.Set("name", role.Name())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

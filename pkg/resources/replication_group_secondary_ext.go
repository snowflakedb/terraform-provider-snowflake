package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	replicationGroupSecondaryCreateExt      = TrackingCreateWrapper(resources.ReplicationGroupSecondary, CreateReplicationGroupSecondary)
	replicationGroupSecondaryReadExt        = TrackingReadWrapper(resources.ReplicationGroupSecondary, ReadReplicationGroupSecondaryFunc(true))
	replicationGroupSecondaryDescriptionExt = "Resource used to manage a replication group replica in the target account. For more information, check [replication group documentation](https://docs.snowflake.com/en/sql-reference/sql/create-replication-group)."
	replicationGroupSecondaryImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.ReplicationGroupSecondary, ImportReplicationGroupSecondary),
	}
	replicationGroupSecondaryShowOutputSchemaExt = map[string]*schema.Schema{
		ShowOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW REPLICATION GROUPS` for the given replication group secondary.",
			Elem: &schema.Resource{
				Schema: schemas.ShowReplicationGroupSchema,
			},
		},
	}
)

func replicationGroupSecondaryNewCreateRequestExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateSecondaryReplicationGroupRequest, error) {
	raw, ok := d.GetOk("primary")
	if !ok {
		return nil, fmt.Errorf("primary is required")
	}
	primary, err := sdk.ParseExternalObjectIdentifier(raw.(string))
	if err != nil {
		return nil, err
	}
	return sdk.NewCreateSecondaryReplicationGroupRequest(id, primary), nil
}

func replicationGroupSecondaryCreateInSdkExt(ctx context.Context, meta any, _ sdk.AccountObjectIdentifier, req *sdk.CreateSecondaryReplicationGroupRequest) diag.Diagnostics {
	if err := meta.(*provider.Context).Client.ReplicationGroups.CreateSecondary(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func replicationGroupSecondaryShowOutputSetExt(d *schema.ResourceData, group *sdk.ReplicationGroup) error {
	return d.Set(ShowOutputAttributeName, []map[string]any{schemas.ReplicationGroupToSchema(group)})
}

func replicationGroupSecondarySetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, group *sdk.ReplicationGroup, _ bool) error {
	return errors.Join(
		d.Set("name", id.Name()),
		d.Set("primary", group.Primary.FullyQualifiedName()),
	)
}

package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	replicationGroupSecondaryCreateExt      = PreviewFeatureCreateContextWrapper(string(previewfeatures.ReplicationGroupSecondaryResource), TrackingCreateWrapper(resources.ReplicationGroupSecondary, CreateReplicationGroupSecondary))
	replicationGroupSecondaryReadExt        = PreviewFeatureReadContextWrapper(string(previewfeatures.ReplicationGroupSecondaryResource), TrackingReadWrapper(resources.ReplicationGroupSecondary, ReadReplicationGroupSecondaryFunc(true)))
	replicationGroupSecondaryDescriptionExt = "Resource used to manage a replication group replica in the target account. For more information, check [replication group documentation](https://docs.snowflake.com/en/sql-reference/sql/create-replication-group)."
	replicationGroupSecondaryImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.ReplicationGroupSecondary, ImportReplicationGroupSecondary),
	}
	// Listed on the replication group secondary SSOT object and appended to the generated customdiff.All.
	replicationGroupSecondaryRecreateWhenPrimary = RecreateWhenResourceBoolFieldChangedExternally("is_primary", false)
)

func replicationGroupSecondaryCreateInSdkExt(ctx context.Context, meta any, _ sdk.AccountObjectIdentifier, req *sdk.CreateSecondaryReplicationGroupRequest) diag.Diagnostics {
	if err := meta.(*provider.Context).Client.ReplicationGroups.CreateSecondary(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func replicationGroupSecondarySetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, _ *sdk.ReplicationGroup, _ bool) error {
	return d.Set("name", id.Name())
}

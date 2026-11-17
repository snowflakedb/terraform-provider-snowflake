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
	failoverGroupSecondaryCreateExt      = PreviewFeatureCreateContextWrapper(string(previewfeatures.FailoverGroupSecondaryResource), TrackingCreateWrapper(resources.FailoverGroupSecondary, CreateFailoverGroupSecondary))
	failoverGroupSecondaryReadExt        = PreviewFeatureReadContextWrapper(string(previewfeatures.FailoverGroupSecondaryResource), TrackingReadWrapper(resources.FailoverGroupSecondary, ReadFailoverGroupSecondaryFunc(true)))
	failoverGroupSecondaryDescriptionExt = "Resource used to manage a failover group replica in the target account. For more information, check [failover group documentation](https://docs.snowflake.com/en/sql-reference/sql/create-failover-group)."
	failoverGroupSecondaryImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.FailoverGroupSecondary, ImportFailoverGroupSecondary),
	}
)

func failoverGroupSecondaryCreateInSdkExt(ctx context.Context, meta any, _ sdk.AccountObjectIdentifier, req *sdk.CreateSecondaryFailoverGroupRequest) diag.Diagnostics {
	if err := meta.(*provider.Context).Client.FailoverGroups.CreateSecondary(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func failoverGroupSecondarySetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, _ *sdk.FailoverGroup, _ bool) error {
	return d.Set("name", id.Name())
}

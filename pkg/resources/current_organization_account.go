package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var currentOrganizationAccountSchema = map[string]*schema.Schema{
	"resource_monitor": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      externalChangesNotDetectedFieldDescription("Parameter that specifies the name of the resource monitor used to control all virtual warehouses created in the account."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"password_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies [password policy](https://docs.snowflake.com/en/user-guide/password-authentication#label-using-password-policies) for the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"session_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies [session policy](https://docs.snowflake.com/en/user-guide/session-policies-using) for the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

func CurrentOrganizationAccount() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource used to manage an organization account within the organization you are connected to. See [ALTER ORGANIZATION ACCOUNT](https://docs.snowflake.com/en/sql-reference/sql/alter-organization-account) documentation for more information on resource capabilities.",
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.CurrentOrganizationAccountResource), TrackingCreateWrapper(resources.CurrentOrganizationAccount, CreateCurrentOrganizationAccount)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.CurrentOrganizationAccountResource), TrackingReadWrapper(resources.CurrentOrganizationAccount, ReadCurrentOrganizationAccount)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.CurrentOrganizationAccountResource), TrackingUpdateWrapper(resources.CurrentOrganizationAccount, UpdateCurrentOrganizationAccount)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.CurrentOrganizationAccountResource), TrackingDeleteWrapper(resources.CurrentOrganizationAccount, DeleteCurrentOrganizationAccount)),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.CurrentOrganizationAccount, accountParametersCustomDiff),

		Schema: collections.MergeMaps(currentOrganizationAccountSchema, accountParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.CurrentOrganizationAccount, schema.ImportStatePassthroughContext),
		},

		Timeouts: defaultTimeouts,
	}
}

func CreateCurrentOrganizationAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func ReadCurrentOrganizationAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func UpdateCurrentOrganizationAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func DeleteCurrentOrganizationAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

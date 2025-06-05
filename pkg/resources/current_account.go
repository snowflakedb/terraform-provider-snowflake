package resources

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var currentAccountSchema = map[string]*schema.Schema{
	"resource_monitor": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      relatedResourceDescription("Specifies an existing network policy. This network policy controls network traffic that is attempting to exchange an authorization code for an access or refresh token or to use a refresh token to obtain a new access token.", resources.NetworkPolicy),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"authentication_policy": {
		Type: schema.TypeString,
	},
	"password_policy": {
		Type: schema.TypeString,
	},
	"session_policy": {
		Type: schema.TypeString,
	},
	"feature_policy": {
		Type: schema.TypeString,
	},
	"packages_policy": {
		Type: schema.TypeString,
	},
	"organization_user_group": {},
	"tags": {
		Type: schema.TypeString,
	},
}

func CurrentAccount() *schema.Resource {
	return &schema.Resource{
		Description:   "TODO",
		CreateContext: TrackingCreateWrapper(resources.Account, CreateCurrentAccount),
		ReadContext:   TrackingReadWrapper(resources.Account, ReadCurrentAccount),
		UpdateContext: TrackingUpdateWrapper(resources.Account, UpdateCurrentAccount),
		DeleteContext: TrackingDeleteWrapper(resources.Account, DeleteCurrentAccount),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Account, taskParametersCustomDiff),

		Schema: collections.MergeMaps(currentAccountSchema, accountParametersSchema),

		Timeouts: defaultTimeouts,
	}
}

func CreateCurrentAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return ReadCurrentAccount(false)(ctx, d, meta)
}

func ReadCurrentAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// TODO: Check if withExternalChangesMarking is needed here or if it can be removed
	// client := meta.(*provider.Context).Client

	return nil
}

func UpdateCurrentAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// client := meta.(*provider.Context).Client
	return ReadAccount(false)(ctx, d, meta)
}

func DeleteCurrentAccount(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	d.SetId("")
	return nil
}

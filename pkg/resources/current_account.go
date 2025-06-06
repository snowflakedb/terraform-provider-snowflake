package resources

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var currentAccountSchema = map[string]*schema.Schema{
	"resource_monitor": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"authentication_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"password_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"session_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"feature_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"packages_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"organization_user_group": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		Optional:    true,
		Description: "The list of organization user groups imported into the account.",
	},
	// TODO: Tags are done by tags_association resource
}

func CurrentAccount() *schema.Resource {
	return &schema.Resource{
		Description:   "TODO",
		CreateContext: TrackingCreateWrapper(resources.Account, CreateCurrentAccount),
		ReadContext:   TrackingReadWrapper(resources.Account, ReadCurrentAccount),
		UpdateContext: TrackingUpdateWrapper(resources.Account, UpdateCurrentAccount),
		DeleteContext: TrackingDeleteWrapper(resources.Account, DeleteCurrentAccount),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Account, accountParametersCustomDiff),

		Schema: collections.MergeMaps(currentAccountSchema, accountParametersSchema),

		Timeouts: defaultTimeouts,
	}
}

func CreateCurrentAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	d.SetId("current_account")
	return UpdateCurrentAccount(ctx, d, meta)
}

func ReadCurrentAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// TODO: Check if withExternalChangesMarking is needed here or if it can be removed
	client := meta.(*provider.Context).Client

	// TODO: Get policy links

	parameters, err := client.Accounts.ShowParameters(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := handleAccountParameterRead(d, parameters); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateCurrentAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	if d.HasChange("resource_monitor") {
		opts := new(sdk.AlterAccountOptions)
		if err := accountObjectIdentifierAttributeUpdate(d, "resource_monitor", &opts.Set.ResourceMonitor, &opts.Unset.ResourceMonitor); err != nil {
			return diag.FromErr(err)
		}
		if opts.Set.ResourceMonitor != nil || opts.Unset.ResourceMonitor != nil {
			if err := client.Accounts.Alter(ctx, opts); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("authentication_policy") {
		opts := new(sdk.AlterAccountOptions)
		if err := schemaObjectIdentifierAttributeUpdate(d, "authentication_policy", &opts.Set.AuthenticationPolicy, &opts.Unset.AuthenticationPolicy); err != nil {
			return diag.FromErr(err)
		}
		if opts.Set.AuthenticationPolicy != nil || opts.Unset.AuthenticationPolicy != nil {
			if err := client.Accounts.Alter(ctx, opts); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("password_policy") {
		opts := new(sdk.AlterAccountOptions)
		if err := schemaObjectIdentifierAttributeUpdate(d, "password_policy", &opts.Set.PasswordPolicy, &opts.Unset.PasswordPolicy); err != nil {
			return diag.FromErr(err)
		}
		if opts.Set.PasswordPolicy != nil || opts.Unset.PasswordPolicy != nil {
			if err := client.Accounts.Alter(ctx, opts); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("session_policy") {
		opts := new(sdk.AlterAccountOptions)
		if err := schemaObjectIdentifierAttributeUpdate(d, "session_policy", &opts.Set.SessionPolicy, &opts.Unset.SessionPolicy); err != nil {
			return diag.FromErr(err)
		}
		if opts.Set.SessionPolicy != nil || opts.Unset.SessionPolicy != nil {
			if err := client.Accounts.Alter(ctx, opts); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("packages_policy") {
		opts := new(sdk.AlterAccountOptions)
		if err := schemaObjectIdentifierAttributeUpdate(d, "packages_policy", &opts.Set.PackagesPolicy, &opts.Unset.PackagesPolicy); err != nil {
			return diag.FromErr(err)
		}
		if opts.Set.PackagesPolicy != nil || opts.Unset.PackagesPolicy != nil {
			if err := client.Accounts.Alter(ctx, opts); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	//if errs := errors.Join(
	//	alterAccountIfAttributeChanged(client, ctx, d, schemaObjectIdentifierAttributeUpdate, "feature_policy", func(opts *sdk.AlterAccountOptions) *sdk.SchemaObjectIdentifier { return opts.Set.FeaturePolicy }, func(opts *sdk.AlterAccountOptions) *bool { return opts.Unset.FeaturePolicy }),
	//); errs != nil {
	//	return diag.FromErr(errs)
	//}

	//before, after := d.GetChange("organization_user_group")
	//// TODO: Update if not empty

	setParameters := new(sdk.AccountSet)
	unsetParameters := new(sdk.AccountUnset)
	if diags := handleAccountParametersUpdate(d, setParameters, unsetParameters); diags != nil {
		return diags
	}

	return ReadCurrentAccount(ctx, d, meta)
}

func DeleteCurrentAccount(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	d.SetId("")
	return nil
}

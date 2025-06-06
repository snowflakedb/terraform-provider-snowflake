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

	alterIfIdentifierAttributeChanged := func(set *sdk.AccountSet, unset *sdk.AccountUnset, setId sdk.ObjectIdentifier, unsetBool *bool) diag.Diagnostics {
		if setId != nil {
			if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: set}); err != nil {
				return diag.FromErr(err)
			}
		}
		if unsetBool != nil && *unsetBool {
			if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Unset: unset}); err != nil {
				return diag.FromErr(err)
			}
		}
		return nil
	}

	if d.HasChange("resource_monitor") {
		set, unset := new(sdk.AccountSet), new(sdk.AccountUnset)
		if err := accountObjectIdentifierAttributeUpdate(d, "resource_monitor", &set.ResourceMonitor, &unset.ResourceMonitor); err != nil {
			return diag.FromErr(err)
		}
		if diags := alterIfIdentifierAttributeChanged(set, unset, set.ResourceMonitor, unset.ResourceMonitor); diags != nil {
			return diags
		}
	}

	if d.HasChange("authentication_policy") {
		set, unset := new(sdk.AccountSet), new(sdk.AccountUnset)
		if err := schemaObjectIdentifierAttributeUpdate(d, "authentication_policy", &set.AuthenticationPolicy, &unset.AuthenticationPolicy); err != nil {
			return diag.FromErr(err)
		}
		if diags := alterIfIdentifierAttributeChanged(set, unset, set.AuthenticationPolicy, unset.AuthenticationPolicy); diags != nil {
			return diags
		}
	}

	if d.HasChange("password_policy") {
		set, unset := new(sdk.AccountSet), new(sdk.AccountUnset)
		if err := schemaObjectIdentifierAttributeUpdate(d, "password_policy", &set.PasswordPolicy, &unset.PasswordPolicy); err != nil {
			return diag.FromErr(err)
		}
		if diags := alterIfIdentifierAttributeChanged(set, unset, set.PasswordPolicy, unset.PasswordPolicy); diags != nil {
			return diags
		}
	}

	if d.HasChange("session_policy") {
		set, unset := new(sdk.AccountSet), new(sdk.AccountUnset)
		if err := schemaObjectIdentifierAttributeUpdate(d, "session_policy", &set.SessionPolicy, &unset.SessionPolicy); err != nil {
			return diag.FromErr(err)
		}
		if diags := alterIfIdentifierAttributeChanged(set, unset, set.SessionPolicy, unset.SessionPolicy); diags != nil {
			return diags
		}
	}

	if d.HasChange("packages_policy") {
		set, unset := new(sdk.AccountSet), new(sdk.AccountUnset)
		if err := schemaObjectIdentifierAttributeUpdate(d, "packages_policy", &set.PackagesPolicy, &unset.PackagesPolicy); err != nil {
			return diag.FromErr(err)
		}
		if diags := alterIfIdentifierAttributeChanged(set, unset, set.PackagesPolicy, unset.PackagesPolicy); diags != nil {
			return diags
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

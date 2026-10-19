package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func AuthenticationPolicy() *schema.Resource {
	return &schema.Resource{
		Schema:      authenticationPolicySchema,
		Description: authenticationPolicyDescriptionExt,

		CreateContext: authenticationPolicyCreateExt,
		ReadContext:   authenticationPolicyReadExt,
		UpdateContext: authenticationPolicyUpdateExt,
		DeleteContext: authenticationPolicyDeleteExt,
		Importer:      authenticationPolicyImporterExt,

		CustomizeDiff:  authenticationPolicyCustomizeDiffExt,
		Timeouts:       defaultTimeouts,
		SchemaVersion:  authenticationPolicySchemaVersionExt,
		StateUpgraders: authenticationPolicyStateUpgradersExt,
	}
}

func ImportAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	id, err := authenticationPolicyParseIdExt(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := authenticationPolicyShowByIdInSdkExt(ctx, meta, id)
	if err != nil {
		return nil, err
	}

	if err := authenticationPolicySetFieldsNotSetByReadExt(d, id, details); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := authenticationPolicyParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := authenticationPolicyNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := authenticationPolicyApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if diags := authenticationPolicyCreateInSdkExt(ctx, meta, id, req); diags != nil {
		return diags
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadAuthenticationPolicyFunc(false)(ctx, d, meta)
}

func ReadAuthenticationPolicyFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		id, err := authenticationPolicyParseIdExt(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		authenticationPolicy, diags := authenticationPolicyShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags != nil {
			return diags
		}
		authenticationPolicyDescriptions, err := authenticationPolicyDescribeInSdkExt(ctx, meta, id)
		if err != nil {
			return diag.FromErr(err)
		}
		// DESCRIBE details, and property-read failures stay warnings, so this returns diagnostics.
		diags = authenticationPolicyHandleExternalChangesInDescribeExt(d, authenticationPolicyDescriptions, withExternalChangesMarking)
		if diags.HasError() {
			return diags
		}

		if err = setStateToValuesFromConfig(d, authenticationPolicySchema, authenticationPolicySetStateToValueFromConfigFieldsExt); err != nil {
			// Keep warning diagnostics accumulated above; diag.FromErr alone would drop them.
			return append(diags, diag.FromErr(err)...)
		}
		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.AuthenticationPolicyToSchema(authenticationPolicy)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.AuthenticationPolicyDescriptionsToSchema(authenticationPolicyDescriptions)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			authenticationPolicySetConfigFieldsExt(d, authenticationPolicy),
		)
		if errs != nil {
			return append(diags, diag.FromErr(errs)...)
		}
		return diags
	}
}

func UpdateAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := authenticationPolicyParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := authenticationPolicyHierarchyRenameExt(ctx, d, meta, &id); diags != nil {
		return diags
	}

	setUnset, err := authenticationPolicyApplySetUnsetExt(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := authenticationPolicyAlterSetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	if err := authenticationPolicyAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	return ReadAuthenticationPolicyFunc(false)(ctx, d, meta)
}

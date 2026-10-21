package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Alert() *schema.Resource {
	return &schema.Resource{
		Schema:      alertSchema,
		Description: alertDescriptionExt,

		CreateContext: alertCreateExt,
		ReadContext:   alertReadExt,
		UpdateContext: alertUpdateExt,
		DeleteContext: alertDeleteExt,
		Importer:      alertImporterExt,

		CustomizeDiff: alertCustomizeDiffExt,
		Timeouts:      defaultTimeouts,
	}
}

func CreateAlert(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := alertParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := alertNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := alertApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if diags := alertCreateInSdkExt(ctx, meta, id, req); diags != nil {
		return diags
	}

	// Alert still uses the legacy pipe-delimited identifier.
	d.SetId(helpers.EncodeSnowflakeID(id))
	diags := alertAfterCreateExt(ctx, d, meta, id)
	return append(diags, ReadAlertFunc(false)(ctx, d, meta)...)
}

func ReadAlertFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		id, err := alertParseIdExt(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		alert, diags := alertShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags != nil {
			return diags
		}

		errs := errors.Join(
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			alertSetConfigFieldsExt(d, alert),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateAlert(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := alertParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := alertBeforeAlterExt(ctx, d, meta, id); diags != nil {
		return diags
	}

	setUnset, err := alertApplySetUnsetExt(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := alertAlterSetInSdkExt(ctx, d, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	diags := alertAfterAlterExt(ctx, d, meta, id)
	return append(diags, ReadAlertFunc(false)(ctx, d, meta)...)
}

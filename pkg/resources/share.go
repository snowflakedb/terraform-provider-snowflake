package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Share() *schema.Resource {
	return &schema.Resource{
		Schema:      shareSchema,
		Description: shareDescriptionExt,

		CreateContext: shareCreateExt,
		ReadContext:   shareReadExt,
		UpdateContext: shareUpdateExt,
		DeleteContext: shareDeleteExt,
		Importer:      shareImporterExt,

		CustomizeDiff: shareCustomizeDiffExt,
		Timeouts:      defaultTimeouts,
	}
}

// CreateShare implements schema.CreateFunc.
func CreateShare(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := shareParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := shareNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := shareApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if diags := shareCreateInSdkExt(ctx, meta, id, req); diags != nil {
		return diags
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	if diags := shareAfterCreateExt(ctx, d, meta, id); diags != nil {
		return diags
	}
	return ReadShareFunc(false)(ctx, d, meta)
}

// ReadShareFunc implements schema.ReadFunc. The share does not mark external changes yet, so the flag is not used.
func ReadShareFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		id, err := shareParseIdExt(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		// Differs from the generic Read: the hook can return a warning (the recipients lookup failed) together with the
		// share, so only an error stops the read and the warnings are returned at the end, as before.
		share, diags := shareShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags.HasError() {
			return diags
		}

		errs := errors.Join(
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			shareSetConfigFieldsExt(d, share),
		)
		if errs != nil {
			return append(diags, diag.FromErr(errs)...)
		}
		return diags
	}
}

// UpdateShare implements schema.UpdateFunc.
func UpdateShare(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := shareParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	setUnset, err := shareApplySetUnsetExt(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := shareAlterSetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	return ReadShareFunc(false)(ctx, d, meta)
}

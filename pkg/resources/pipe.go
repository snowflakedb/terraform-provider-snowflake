package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Pipe() *schema.Resource {
	return &schema.Resource{
		Schema:      pipeSchema,
		Description: pipeDescriptionExt,

		CreateContext: pipeCreateExt,
		ReadContext:   pipeReadExt,
		UpdateContext: pipeUpdateExt,
		DeleteContext: pipeDeleteExt,
		Importer:      pipeImporterExt,

		CustomizeDiff: pipeCustomizeDiffExt,
		Timeouts:      defaultTimeouts,
	}
}

// CreatePipe implements schema.CreateFunc.
func CreatePipe(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := pipeParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := pipeNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := pipeApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if diags := pipeCreateInSdkExt(ctx, meta, id, req); diags != nil {
		return diags
	}

	// Pipe still uses the legacy pipe-delimited identifier.
	d.SetId(helpers.EncodeSnowflakeID(id))
	return ReadPipeFunc(false)(ctx, d, meta)
}

// ReadPipeFunc implements schema.ReadFunc. Pipe has no external changes marking.
func ReadPipeFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		id, err := pipeParseIdExt(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		pipe, diags := pipeShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags != nil {
			return diags
		}

		errs := errors.Join(
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			pipeSetConfigFieldsExt(d, pipe),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

// UpdatePipe implements schema.UpdateFunc.
func UpdatePipe(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := pipeParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	setUnset, err := pipeApplySetUnsetExt(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := pipeAlterSetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	if err := pipeAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	return ReadPipeFunc(false)(ctx, d, meta)
}

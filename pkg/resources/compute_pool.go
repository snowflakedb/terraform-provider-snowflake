package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ComputePool() *schema.Resource {
	return &schema.Resource{
		Schema:      computePoolSchema,
		Description: computePoolDescriptionExt,

		CreateContext: computePoolCreateExt,
		ReadContext:   computePoolReadExt,
		UpdateContext: computePoolUpdateExt,
		DeleteContext: computePoolDeleteExt,
		Importer:      computePoolImporterExt,

		CustomizeDiff: computePoolCustomizeDiffExt,
		Timeouts:      defaultTimeouts,
	}
}

func ImportComputePool(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	id, err := computePoolParseIdExt(d.Id())
	if err != nil {
		return nil, err
	}

	computePool, err := computePoolShowByIdInSdkExt(ctx, meta, id)
	if err != nil {
		return nil, err
	}

	if err := computePoolSetFieldsNotSetByReadExt(d, id, computePool); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateComputePool(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := computePoolParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := computePoolNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := computePoolApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if err := computePoolCreateInSdkExt(ctx, meta, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadComputePoolFunc(false)(ctx, d, meta)
}

func ReadComputePoolFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		id, err := computePoolParseIdExt(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		computePool, diags := computePoolShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags != nil {
			return diags
		}
		computePoolDetails, err := computePoolDescribeInSdkExt(ctx, meta, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = computePoolHandleExternalChangesExt(d, computePool, withExternalChangesMarking); err != nil {
			return diag.FromErr(err)
		}

		if err = setStateToValuesFromConfig(d, computePoolSchema, computePoolSetStateToValueFromConfigFieldsExt); err != nil {
			return diag.FromErr(err)
		}
		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.ComputePoolToSchema(computePool)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ComputePoolDetailsToSchema(computePoolDetails)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			computePoolSetConfigFieldsExt(d, computePool),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateComputePool(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := computePoolParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	setUnset, err := computePoolApplySetUnsetExt(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := computePoolAlterSetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	if err := computePoolAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	return ReadComputePoolFunc(false)(ctx, d, meta)
}

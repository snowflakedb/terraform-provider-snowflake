package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func WarehouseAdaptive() *schema.Resource {
	return &schema.Resource{
		Schema:      warehouseAdaptiveSchema,
		Description: warehouseAdaptiveDescriptionExt,

		CreateContext: warehouseAdaptiveCreateExt,
		ReadContext:   warehouseAdaptiveReadExt,
		UpdateContext: warehouseAdaptiveUpdateExt,
		DeleteContext: warehouseAdaptiveDeleteExt,
		Importer:      warehouseAdaptiveImporterExt,

		CustomizeDiff: warehouseAdaptiveCustomizeDiffExt,
		Timeouts:      defaultTimeouts,
	}
}

func ImportWarehouseAdaptive(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	id, err := warehouseAdaptiveParseIdExt(d.Id())
	if err != nil {
		return nil, err
	}

	w, err := warehouseAdaptiveShowByIdInSdkExt(ctx, meta, id)
	if err != nil {
		return nil, err
	}

	// Interactive rejection and the import-only field sets stay in SetFieldsNotSetByReadExt.
	if err := warehouseAdaptiveSetFieldsNotSetByReadExt(d, id, w); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateWarehouseAdaptive(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := warehouseAdaptiveParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := warehouseAdaptiveNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := warehouseAdaptiveApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if diags := warehouseAdaptiveApplyParametersCreateExt(d, req); diags != nil {
		return diags
	}

	// CreateAdaptive and the wrapped create error stay in CreateInSdkExt.
	if diags := warehouseAdaptiveCreateInSdkExt(ctx, meta, id, req); diags != nil {
		return diags
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadWarehouseAdaptiveFunc(false)(ctx, d, meta)
}

func ReadWarehouseAdaptiveFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		id, err := warehouseAdaptiveParseIdExt(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		w, diags := warehouseAdaptiveShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags != nil {
			return diags
		}
		// Describe is omitted: this resource has no describe_output.
		// Details come from SHOW PARAMETERS plus ToWarehouseParametersDetails, not ShowParametersDetails.
		warehouseParameters, err := warehouseAdaptiveShowParametersInSdkExt(ctx, meta, id)
		if err != nil {
			return diag.FromErr(err)
		}
		warehouseParametersDetails, err := sdk.ToWarehouseParametersDetails(warehouseParameters)
		if err != nil {
			return diag.FromErr(err)
		}
		if diags := warehouseAdaptiveSetParametersFieldsExt(d, warehouseParametersDetails); diags != nil {
			return diags
		}
		if err = warehouseAdaptiveHandleExternalChangesExt(d, w, withExternalChangesMarking); err != nil {
			return diag.FromErr(err)
		}

		// setStateToValuesFromConfig is omitted: adaptive Read does not copy config back onto state.
		errs := errors.Join(
			// SHOW returns sdk.Warehouse; show_output uses the adaptive projection.
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.WarehouseAdaptiveToSchema(w.AsAdaptive())}),
			d.Set(ParametersAttributeName, []map[string]any{schemas.WarehouseAdaptiveParametersToSchema(warehouseParameters, meta.(*provider.Context))}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			warehouseAdaptiveSetConfigFieldsExt(d, w),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateWarehouseAdaptive(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := warehouseAdaptiveParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := warehouseAdaptiveHierarchyRenameExt(ctx, d, meta, &id); diags != nil {
		return diags
	}

	setUnset, err := warehouseAdaptiveApplySetUnsetExt(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if diags := warehouseAdaptiveApplyParametersChangesExt(d, setUnset); diags != nil {
		return diags
	}
	// d matches warehouse AlterSetInSdkExt. SET uses AlterWithSuspend inside the hook.
	if err := warehouseAdaptiveAlterSetInSdkExt(ctx, d, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	if err := warehouseAdaptiveAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	return ReadWarehouseAdaptiveFunc(false)(ctx, d, meta)
}

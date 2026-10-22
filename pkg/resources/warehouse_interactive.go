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

func WarehouseInteractive() *schema.Resource {
	return &schema.Resource{
		Schema:      warehouseInteractiveSchema,
		Description: warehouseInteractiveDescriptionExt,

		CreateContext: warehouseInteractiveCreateExt,
		ReadContext:   warehouseInteractiveReadExt,
		UpdateContext: warehouseInteractiveUpdateExt,
		DeleteContext: warehouseInteractiveDeleteExt,
		Importer:      warehouseInteractiveImporterExt,

		CustomizeDiff: warehouseInteractiveCustomizeDiffExt,
		Timeouts:      defaultTimeouts,
	}
}

func ImportWarehouseInteractive(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	id, err := warehouseInteractiveParseIdExt(d.Id())
	if err != nil {
		return nil, err
	}

	w, err := warehouseInteractiveShowByIdInSdkExt(ctx, meta, id)
	if err != nil {
		return nil, err
	}

	if err := warehouseInteractiveSetFieldsNotSetByReadExt(d, id, w); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateWarehouseInteractive(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := warehouseInteractiveParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := warehouseInteractiveNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := warehouseInteractiveApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if diags := warehouseInteractiveApplyParametersCreateExt(d, req); diags != nil {
		return diags
	}

	if diags := warehouseInteractiveCreateInSdkExt(ctx, meta, id, req); diags != nil {
		return diags
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadWarehouseInteractiveFunc(false)(ctx, d, meta)
}

func ReadWarehouseInteractiveFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		id, err := warehouseInteractiveParseIdExt(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		w, diags := warehouseInteractiveShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags != nil {
			return diags
		}
		// SHOW PARAMETERS is converted locally. Read does not call ShowParametersDetails.
		warehouseParameters, err := warehouseInteractiveShowParametersInSdkExt(ctx, meta, id)
		if err != nil {
			return diag.FromErr(err)
		}
		warehouseParametersDetails, err := sdk.ToWarehouseParametersDetails(warehouseParameters)
		if err != nil {
			return diag.FromErr(err)
		}
		if diags := warehouseInteractiveSetParametersFieldsExt(d, warehouseParametersDetails); diags != nil {
			return diags
		}
		if err = warehouseInteractiveHandleExternalChangesExt(d, w, withExternalChangesMarking); err != nil {
			return diag.FromErr(err)
		}

		// No DESCRIBE, and Read does not copy config-backed fields back into state.
		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.WarehouseInteractiveToSchema(w.AsInteractive())}),
			d.Set(ParametersAttributeName, []map[string]any{schemas.WarehouseInteractiveParametersToSchema(warehouseParameters, meta.(*provider.Context))}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			warehouseInteractiveSetConfigFieldsExt(d, w),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateWarehouseInteractive(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := warehouseInteractiveParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := warehouseInteractiveHierarchyRenameExt(ctx, d, meta, &id); diags != nil {
		return diags
	}

	// Tables ADD/DROP runs before object-attribute SET/UNSET.
	if diags := warehouseInteractiveBeforeAlterExt(ctx, d, meta, id); diags != nil {
		return diags
	}

	setUnset, err := warehouseInteractiveApplySetUnsetExt(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if diags := warehouseInteractiveApplyParametersChangesExt(d, setUnset); diags != nil {
		return diags
	}
	// d matches warehouseAlterSetInSdkExt. Size changes suspend through AlterWithSuspend.
	if err := warehouseInteractiveAlterSetInSdkExt(ctx, d, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	if err := warehouseInteractiveAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	return ReadWarehouseInteractiveFunc(false)(ctx, d, meta)
}

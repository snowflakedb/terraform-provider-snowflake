package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Warehouse() *schema.Resource {
	return &schema.Resource{
		Schema:      warehouseSchema,
		Description: warehouseDescriptionExt,

		CreateContext: warehouseCreateExt,
		ReadContext:   warehouseReadExt,
		UpdateContext: warehouseUpdateExt,
		DeleteContext: warehouseDeleteExt,
		Importer:      warehouseImporterExt,

		CustomizeDiff:  warehouseCustomizeDiffExt,
		Timeouts:       defaultTimeouts,
		SchemaVersion:  warehouseSchemaVersionExt,
		StateUpgraders: warehouseStateUpgradersExt,
	}
}

func ImportWarehouse(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	id, err := warehouseParseIdExt(d.Id())
	if err != nil {
		return nil, err
	}

	warehouse, err := warehouseShowByIdInSdkExt(ctx, meta, id)
	if err != nil {
		return nil, err
	}

	if err := warehouseSetFieldsNotSetByReadExt(d, id, warehouse); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateWarehouse(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := warehouseParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := warehouseNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := warehouseApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if diags := warehouseApplyParametersCreateExt(d, req); diags != nil {
		return diags
	}

	if diags := warehouseCreateInSdkExt(ctx, meta, id, req); diags != nil {
		return diags
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadWarehouseFunc(false)(ctx, d, meta)
}

func ReadWarehouseFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		id, err := warehouseParseIdExt(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		warehouse, diags := warehouseShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags != nil {
			return diags
		}
		rawWarehouseParameters, err := warehouseShowParametersInSdkExt(ctx, meta, id)
		if err != nil {
			return diag.FromErr(err)
		}
		warehouseParameters, err := warehouseParametersDetailsFromRawExt(rawWarehouseParameters)
		if err != nil {
			return diag.FromErr(err)
		}
		if diags := warehouseSetParametersFieldsExt(d, warehouseParameters); diags != nil {
			return diags
		}
		if err = warehouseHandleExternalChangesExt(d, warehouse, withExternalChangesMarking); err != nil {
			return diag.FromErr(err)
		}

		if err = setStateToValuesFromConfig(d, warehouseSchema, warehouseSetStateToValueFromConfigFieldsExt); err != nil {
			return diag.FromErr(err)
		}
		errs := errors.Join(
			d.Set(ParametersAttributeName, []map[string]any{schemas.WarehouseParametersToSchema(rawWarehouseParameters, meta.(*provider.Context))}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			warehouseSetConfigFieldsExt(d, warehouse),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateWarehouse(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := warehouseParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := warehouseHierarchyRenameExt(ctx, d, meta, &id); diags != nil {
		return diags
	}

	setUnset, err := warehouseApplySetUnsetExt(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if diags := warehouseApplyParametersChangesExt(d, setUnset); diags != nil {
		return diags
	}
	if err := warehouseAlterSetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	if err := warehouseAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	return ReadWarehouseFunc(false)(ctx, d, meta)
}

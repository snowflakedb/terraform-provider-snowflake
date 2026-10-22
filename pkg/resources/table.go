package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Table() *schema.Resource {
	return &schema.Resource{
		Schema: tableSchema,

		CreateContext: tableCreateExt,
		ReadContext:   tableReadExt,
		UpdateContext: tableUpdateExt,
		DeleteContext: tableDeleteExt,
		Importer:      tableImporterExt,

		CustomizeDiff: tableCustomizeDiffExt,
		Timeouts:      defaultTimeouts,
	}
}

// CreateTable implements schema.CreateFunc.
func CreateTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := tableParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := tableNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := tableApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if diags := tableCreateInSdkExt(ctx, meta, id, req); diags != nil {
		return diags
	}

	// The legacy pipe-delimited ID must stay, because tableParseIdExt decodes it and existing state uses it.
	d.SetId(helpers.EncodeSnowflakeID(id))
	return ReadTableFunc(false)(ctx, d, meta)
}

func ReadTableFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		id, err := tableParseIdExt(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		table, diags := tableShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags != nil {
			return diags
		}
		if done, diags := tableBeforeDescribeExt(ctx, d, meta, id, table); done {
			return diags
		}
		tableDescription, err := tableDescribeInSdkExt(ctx, meta, id)
		if err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			tableSetConfigFieldsExt(d, table),
			tableSetDescribeFieldsExt(d, tableDescription),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

// UpdateTable implements schema.UpdateFunc.
func UpdateTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := tableParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := tableHierarchyRenameExt(ctx, d, meta, &id); diags != nil {
		return diags
	}

	setUnset, err := tableApplySetUnsetExt(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := tableAlterSetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	if err := tableAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}

	if diags := tableAfterAlterExt(ctx, d, meta, id); diags != nil {
		return diags
	}
	return ReadTableFunc(false)(ctx, d, meta)
}

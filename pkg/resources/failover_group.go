package resources

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FailoverGroup() *schema.Resource {
	return &schema.Resource{
		Schema:      failoverGroupSchema,
		Description: failoverGroupDescriptionExt,

		CreateContext: failoverGroupCreateExt,
		ReadContext:   failoverGroupReadExt,
		UpdateContext: failoverGroupUpdateExt,
		DeleteContext: failoverGroupDeleteExt,
		Importer:      failoverGroupImporterExt,

		Timeouts: defaultTimeouts,
	}
}

// CreateFailoverGroup implements schema.CreateFunc.
func CreateFailoverGroup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := failoverGroupParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	done, diags := failoverGroupBeforeCreateExt(ctx, d, meta, id)
	if done {
		return diags
	}

	req, err := failoverGroupNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := failoverGroupApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if diags := failoverGroupCreateInSdkExt(ctx, meta, id, req); diags != nil {
		return diags
	}

	d.SetId(id.Name())

	if diags := failoverGroupAfterCreateExt(ctx, d, meta, id); diags != nil {
		return diags
	}

	return ReadFailoverGroup(ctx, d, meta)
}

// ReadFailoverGroup implements schema.ReadFunc.
func ReadFailoverGroup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := failoverGroupParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	failoverGroup, diags := failoverGroupShowByIdSafelyInSdkExt(ctx, d, meta, id)
	if diags != nil {
		return diags
	}

	errs := errors.Join(
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		failoverGroupSetConfigFieldsExt(ctx, d, meta, failoverGroup),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	return nil
}

// UpdateFailoverGroup implements schema.UpdateFunc.
func UpdateFailoverGroup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := failoverGroupParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := failoverGroupAlterStepsExt(ctx, d, meta, id); diags != nil {
		return diags
	}

	return ReadFailoverGroup(ctx, d, meta)
}

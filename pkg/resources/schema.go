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

func Schema() *schema.Resource {
	return &schema.Resource{
		Schema:      schemaSchema,
		Description: schemaDescriptionExt,

		CreateContext: schemaCreateExt,
		ReadContext:   schemaReadExt,
		UpdateContext: schemaUpdateExt,
		DeleteContext: schemaDeleteExt,
		Importer:      schemaImporterExt,

		CustomizeDiff:  schemaCustomizeDiffExt,
		Timeouts:       defaultTimeouts,
		SchemaVersion:  schemaSchemaVersionExt,
		StateUpgraders: schemaStateUpgradersExt,
	}
}

func ImportSchema(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	id, err := schemaParseIdExt(d.Id())
	if err != nil {
		return nil, err
	}

	s, err := schemaShowByIdInSdkExt(ctx, meta, id)
	if err != nil {
		return nil, err
	}

	if err := schemaSetFieldsNotSetByReadExt(d, id, s); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateSchema(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := schemaParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	done, diags := schemaBeforeCreateExt(ctx, d, meta, id)
	if done {
		return diags
	}

	req, err := schemaNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := schemaApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if diags := schemaApplyParametersCreateExt(d, req); diags != nil {
		return diags
	}

	if err := schemaCreateInSdkExt(ctx, meta, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadSchemaFunc(false)(ctx, d, meta)
}

func ReadSchemaFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		id, err := schemaParseIdExt(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		schema, diags := schemaShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags != nil {
			return diags
		}
		// TODO [next PR]: Unlike compute pool, Describe uses the SHOW object's id and is best-effort: log + skip Set, do not fail Read. Will be addressed with the generation.
		describeResult, describeErr := schemaDescribeInSdkExt(ctx, meta, schema.ID())
		rawSchemaParameters, err := schemaShowParametersInSdkExt(ctx, meta, id)
		if err != nil {
			return diag.FromErr(err)
		}
		schemaParameters, err := schemaShowParametersDetailsInSdkExt(ctx, meta, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if diags := schemaSetParametersFieldsExt(d, schemaParameters); diags != nil {
			return diags
		}
		if err = schemaHandleExternalChangesExt(d, schema, withExternalChangesMarking); err != nil {
			return diag.FromErr(err)
		}

		if err = setStateToValuesFromConfig(d, schemaSchema, schemaSetStateToValueFromConfigFieldsExt); err != nil {
			return diag.FromErr(err)
		}
		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.SchemaToSchema(schema)}),
			schemaSetDescribeOutputExt(d, id, describeResult, describeErr),
			d.Set(ParametersAttributeName, []map[string]any{schemas.SchemaParametersToSchema(rawSchemaParameters, meta.(*provider.Context))}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			schemaSetConfigFieldsExt(d, schema),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateSchema(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := schemaParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := schemaHierarchyRenameExt(ctx, d, meta, &id); diags != nil {
		return diags
	}

	if diags := schemaBeforeAlterExt(ctx, d, meta, id); diags != nil {
		return diags
	}

	setUnset, err := schemaApplySetUnsetExt(d)
	if err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}
	if diags := schemaApplyParametersChangesExt(d, setUnset); diags != nil {
		d.Partial(true)
		return diags
	}
	if err := schemaAlterSetInSdkExt(ctx, meta, id, setUnset); err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}
	if err := schemaAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}
	return ReadSchemaFunc(false)(ctx, d, meta)
}

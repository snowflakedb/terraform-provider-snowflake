package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	tfschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CatalogIntegrationAwsGlue() *tfschema.Resource {
	return &tfschema.Resource{
		Schema:      catalogIntegrationAwsGlueSchema,
		Description: catalogIntegrationAwsGlueDescriptionExt,

		CreateContext: catalogIntegrationAwsGlueCreateExt,
		ReadContext:   catalogIntegrationAwsGlueReadExt,
		UpdateContext: catalogIntegrationAwsGlueUpdateExt,
		DeleteContext: catalogIntegrationAwsGlueDeleteExt,
		Importer:      catalogIntegrationAwsGlueImporterExt,

		CustomizeDiff: catalogIntegrationAwsGlueCustomizeDiffExt,
		Timeouts:      defaultTimeouts,
	}
}

func ImportCatalogIntegrationAwsGlue(ctx context.Context, d *tfschema.ResourceData, meta any) ([]*tfschema.ResourceData, error) {
	id, err := catalogIntegrationAwsGlueParseIdExt(d.Id())
	if err != nil {
		return nil, err
	}

	catalogIntegrationAwsGlue, err := catalogIntegrationAwsGlueShowByIdInSdkExt(ctx, meta, id)
	if err != nil {
		return nil, err
	}

	if err := catalogIntegrationAwsGlueSetFieldsNotSetByReadExt(d, id, catalogIntegrationAwsGlue); err != nil {
		return nil, err
	}
	return []*tfschema.ResourceData{d}, nil
}

func CreateCatalogIntegrationAwsGlue(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {
	id, err := catalogIntegrationAwsGlueParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := catalogIntegrationAwsGlueNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := catalogIntegrationAwsGlueApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if diags := catalogIntegrationAwsGlueCreateInSdkExt(ctx, meta, id, req); diags != nil {
		return diags
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadCatalogIntegrationAwsGlueFunc(false)(ctx, d, meta)
}

func ReadCatalogIntegrationAwsGlueFunc(withExternalChangesMarking bool) tfschema.ReadContextFunc {
	return func(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {
		id, err := catalogIntegrationAwsGlueParseIdExt(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		catalogIntegrationAwsGlue, diags := catalogIntegrationAwsGlueShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags != nil {
			return diags
		}
		catalogIntegrationAwsGlueDetails, err := catalogIntegrationAwsGlueDescribeInSdkExt(ctx, meta, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = catalogIntegrationAwsGlueHandleExternalChangesInDescribeExt(d, catalogIntegrationAwsGlueDetails, withExternalChangesMarking); err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.CatalogIntegrationAwsGlueToSchema(catalogIntegrationAwsGlue)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.CatalogIntegrationAwsGlueDetailsToSchema(catalogIntegrationAwsGlueDetails)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateCatalogIntegrationAwsGlue(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {
	id, err := catalogIntegrationAwsGlueParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	setUnset, err := catalogIntegrationAwsGlueApplySetUnsetExt(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := catalogIntegrationAwsGlueAlterSetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	if err := catalogIntegrationAwsGlueAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	return ReadCatalogIntegrationAwsGlueFunc(false)(ctx, d, meta)
}

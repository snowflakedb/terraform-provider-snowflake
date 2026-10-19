package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	tfschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CatalogIntegrationIcebergRest() *tfschema.Resource {
	return &tfschema.Resource{
		Schema:      catalogIntegrationIcebergRestSchema,
		Description: catalogIntegrationIcebergRestDescriptionExt,

		CreateContext: catalogIntegrationIcebergRestCreateExt,
		ReadContext:   catalogIntegrationIcebergRestReadExt,
		UpdateContext: catalogIntegrationIcebergRestUpdateExt,
		DeleteContext: catalogIntegrationIcebergRestDeleteExt,
		Importer:      catalogIntegrationIcebergRestImporterExt,

		CustomizeDiff: catalogIntegrationIcebergRestCustomizeDiffExt,
		Timeouts:      defaultTimeouts,
	}
}

func ImportCatalogIntegrationIcebergRest(ctx context.Context, d *tfschema.ResourceData, meta any) ([]*tfschema.ResourceData, error) {
	id, err := catalogIntegrationIcebergRestParseIdExt(d.Id())
	if err != nil {
		return nil, err
	}

	catalogIntegrationIcebergRest, err := catalogIntegrationIcebergRestShowByIdInSdkExt(ctx, meta, id)
	if err != nil {
		return nil, err
	}

	if err := catalogIntegrationIcebergRestSetFieldsNotSetByReadExt(d, id, catalogIntegrationIcebergRest); err != nil {
		return nil, err
	}
	return []*tfschema.ResourceData{d}, nil
}

func CreateCatalogIntegrationIcebergRest(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {
	id, err := catalogIntegrationIcebergRestParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := catalogIntegrationIcebergRestNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := catalogIntegrationIcebergRestApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if diags := catalogIntegrationIcebergRestCreateInSdkExt(ctx, meta, id, req); diags != nil {
		return diags
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadCatalogIntegrationIcebergRestFunc(false)(ctx, d, meta)
}

func ReadCatalogIntegrationIcebergRestFunc(withExternalChangesMarking bool) tfschema.ReadContextFunc {
	return func(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {
		id, err := catalogIntegrationIcebergRestParseIdExt(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		catalogIntegrationIcebergRest, diags := catalogIntegrationIcebergRestShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags != nil {
			return diags
		}
		catalogIntegrationIcebergRestDetails, err := catalogIntegrationIcebergRestDescribeInSdkExt(ctx, meta, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = catalogIntegrationIcebergRestHandleExternalChangesInDescribeExt(d, catalogIntegrationIcebergRestDetails, withExternalChangesMarking); err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.CatalogIntegrationIcebergRestToSchema(catalogIntegrationIcebergRest)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.CatalogIntegrationIcebergRestDetailsToSchema(catalogIntegrationIcebergRestDetails)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateCatalogIntegrationIcebergRest(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {
	id, err := catalogIntegrationIcebergRestParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	setUnset, err := catalogIntegrationIcebergRestApplySetUnsetExt(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := catalogIntegrationIcebergRestAlterSetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	if err := catalogIntegrationIcebergRestAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	return ReadCatalogIntegrationIcebergRestFunc(false)(ctx, d, meta)
}

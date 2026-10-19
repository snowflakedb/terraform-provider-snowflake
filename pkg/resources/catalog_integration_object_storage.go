package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	tfschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CatalogIntegrationObjectStorage() *tfschema.Resource {
	return &tfschema.Resource{
		Schema:      catalogIntegrationObjectStorageSchema,
		Description: catalogIntegrationObjectStorageDescriptionExt,

		CreateContext: catalogIntegrationObjectStorageCreateExt,
		ReadContext:   catalogIntegrationObjectStorageReadExt,
		UpdateContext: catalogIntegrationObjectStorageUpdateExt,
		DeleteContext: catalogIntegrationObjectStorageDeleteExt,
		Importer:      catalogIntegrationObjectStorageImporterExt,

		CustomizeDiff: catalogIntegrationObjectStorageCustomizeDiffExt,
		Timeouts:      defaultTimeouts,
	}
}

func ImportCatalogIntegrationObjectStorage(ctx context.Context, d *tfschema.ResourceData, meta any) ([]*tfschema.ResourceData, error) {
	id, err := catalogIntegrationObjectStorageParseIdExt(d.Id())
	if err != nil {
		return nil, err
	}

	catalogIntegrationObjectStorage, err := catalogIntegrationObjectStorageShowByIdInSdkExt(ctx, meta, id)
	if err != nil {
		return nil, err
	}

	if err := catalogIntegrationObjectStorageSetFieldsNotSetByReadExt(d, id, catalogIntegrationObjectStorage); err != nil {
		return nil, err
	}
	return []*tfschema.ResourceData{d}, nil
}

func CreateCatalogIntegrationObjectStorage(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {
	id, err := catalogIntegrationObjectStorageParseIdFromConfigExt(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := catalogIntegrationObjectStorageNewCreateRequestExt(d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := catalogIntegrationObjectStorageApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}

	if diags := catalogIntegrationObjectStorageCreateInSdkExt(ctx, meta, id, req); diags != nil {
		return diags
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadCatalogIntegrationObjectStorageFunc(false)(ctx, d, meta)
}

func ReadCatalogIntegrationObjectStorageFunc(withExternalChangesMarking bool) tfschema.ReadContextFunc {
	return func(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {
		id, err := catalogIntegrationObjectStorageParseIdExt(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		catalogIntegrationObjectStorage, diags := catalogIntegrationObjectStorageShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags != nil {
			return diags
		}
		catalogIntegrationObjectStorageDetails, err := catalogIntegrationObjectStorageDescribeInSdkExt(ctx, meta, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = catalogIntegrationObjectStorageHandleExternalChangesInDescribeExt(d, catalogIntegrationObjectStorageDetails, withExternalChangesMarking); err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.CatalogIntegrationObjectStorageToSchema(catalogIntegrationObjectStorage)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.CatalogIntegrationObjectStorageDetailsToSchema(catalogIntegrationObjectStorageDetails)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateCatalogIntegrationObjectStorage(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {
	id, err := catalogIntegrationObjectStorageParseIdExt(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	setUnset, err := catalogIntegrationObjectStorageApplySetUnsetExt(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := catalogIntegrationObjectStorageAlterSetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	if err := catalogIntegrationObjectStorageAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
	return ReadCatalogIntegrationObjectStorageFunc(false)(ctx, d, meta)
}

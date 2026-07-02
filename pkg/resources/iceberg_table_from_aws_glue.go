package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var icebergTableFromAwsGlueSchema = collections.MergeMaps(
	icebergTableCommonSchema(),
	map[string]*schema.Schema{
		"catalog_table_name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Specifies the name of the table as it appears in the AWS Glue catalog.",
		},
		"catalog_namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Specifies the namespace (or database) in the AWS Glue catalog that the table belongs to. If not specified, the catalog integration's default namespace is used.",
		},
		"auto_refresh": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validateBooleanString,
			Default:          BooleanDefault,
			Description:      booleanStringFieldDescription("Specifies whether Snowflake should periodically refresh the Iceberg table metadata from the AWS Glue catalog."),
			DiffSuppressFunc: SuppressIfAny(
				IgnoreChangeToCurrentSnowflakePlainValueInOutputWithMapping("show_output", "auto_refresh_status", func(x any) any {
					return len(x.([]any)) != 0
				}),
			),
		},
	},
	icebergTableParametersSchema(),
)

func IcebergTableFromAwsGlue() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.IcebergTableFromAwsGlueResource), TrackingCreateWrapper(resources.IcebergTableFromAwsGlue, CreateIcebergTableFromAwsGlue)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.IcebergTableFromAwsGlueResource), TrackingReadWrapper(resources.IcebergTableFromAwsGlue, ReadIcebergTableFromAwsGlueFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.IcebergTableFromAwsGlueResource), TrackingUpdateWrapper(resources.IcebergTableFromAwsGlue, UpdateIcebergTableFromAwsGlue)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.IcebergTableFromAwsGlueResource), TrackingDeleteWrapper(resources.IcebergTableFromAwsGlue, icebergTableDeleteFunc())),

		Description: "Resource used to manage an Iceberg table whose metadata is managed by an AWS Glue catalog. For more information, check [the official documentation](https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-aws-glue).",

		Schema: icebergTableFromAwsGlueSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.IcebergTableFromAwsGlue, importIcebergTableFromAwsGlue),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(icebergTableFromAwsGlueSchema, ShowOutputAttributeName, "comment", "auto_refresh", "catalog_table_name", "catalog_namespace"),
			ComputedIfAnyAttributeChanged(icebergTableFromAwsGlueSchema, ParametersAttributeName, "external_volume", "catalog", "replace_invalid_characters"),
			icebergTableParametersCustomDiff,
		),
	}
}

func importIcebergTableFromAwsGlue(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	table, err := client.IcebergTables.ShowByIDSafely(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
		return nil, err
	}

	if err := d.Set("auto_refresh", booleanStringFromBool(table.AutoRefreshStatus != nil)); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateIcebergTableFromAwsGlue(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	catalogTableName := d.Get("catalog_table_name").(string)

	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)
	req := sdk.NewCreateFromAwsGlueIcebergTableRequest(id, catalogTableName)

	if err := stringAttributeCreate(d, "comment", &req.Comment); err != nil {
		return diag.FromErr(err)
	}
	if err := stringAttributeCreate(d, "catalog_namespace", &req.CatalogNamespace); err != nil {
		return diag.FromErr(err)
	}
	if err := booleanStringAttributeCreate(d, "auto_refresh", &req.AutoRefresh); err != nil {
		return diag.FromErr(err)
	}
	if diags := handleIcebergTableParametersCreate(d, &req.ExternalVolume, &req.Catalog, &req.ReplaceInvalidCharacters); diags.HasError() {
		return diags
	}

	if err := client.IcebergTables.CreateFromAwsGlue(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Iceberg table from AWS Glue catalog (%s), err = %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadIcebergTableFromAwsGlueFunc(false)(ctx, d, meta)
}

func ReadIcebergTableFromAwsGlueFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		return readIcebergTable(ctx, d, meta, func(d *schema.ResourceData, table *sdk.IcebergTable) error {
			var catalogTableName string
			if table.CatalogTableName != nil {
				catalogTableName = *table.CatalogTableName
			}
			if withExternalChangesMarking {
				var autoRefreshSet string
				if table.AutoRefreshStatus != nil {
					autoRefreshSet = table.AutoRefreshStatus.ExecutionState
				}
				if err := handleExternalChangesToObject(
					d,
					"show_output.0.auto_refresh_status",
					outputMapping{"execution_state", "auto_refresh", autoRefreshSet, booleanStringFromBool(table.AutoRefreshStatus != nil), nil},
				); err != nil {
					return err
				}
				var catalogNamespace string
				if table.CatalogNamespace != nil {
					catalogNamespace = *table.CatalogNamespace
				}
				if err := handleExternalChangesToObjectInShow(
					d,
					outputMapping{"catalog_namespace", "catalog_namespace", catalogNamespace, catalogNamespace, nil},
				); err != nil {
					return err
				}
			}
			return errors.Join(
				d.Set("catalog_table_name", catalogTableName),
			)
		})
	}
}

func UpdateIcebergTableFromAwsGlue(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewIcebergTableSetPropertiesRequest()
	unset := sdk.NewIcebergTableUnsetPropertiesRequest()
	if diags := handleIcebergTableCommonUpdate(d, set, unset); diags.HasError() {
		return diags
	}
	if err := booleanStringAttributeUnsetFallbackUpdate(d, "auto_refresh", &set.AutoRefresh, false); err != nil {
		return diag.FromErr(err)
	}

	if err := applyIcebergTableAlter(ctx, client, id, set, unset); err != nil {
		return diag.FromErr(err)
	}

	return ReadIcebergTableFromAwsGlueFunc(false)(ctx, d, meta)
}

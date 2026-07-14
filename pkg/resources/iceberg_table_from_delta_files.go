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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var icebergTableFromDeltaFilesSchema = collections.MergeMaps(
	icebergTableCommonSchema(),
	map[string]*schema.Schema{
		"base_location": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     validation.StringIsNotEmpty,
			Description:      "Specifies the relative path of the Delta table's directory in the external volume. Cannot be changed after creation.",
			DiffSuppressFunc: ignoreDirectoryPathTrailingSlashSuppressFunc,
		},
		"auto_refresh": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validateBooleanString,
			Default:          BooleanDefault,
			Description:      booleanStringFieldDescription("Specifies whether Snowflake should automatically refresh the Iceberg table metadata when new files are added to the Delta table's directory."),
			DiffSuppressFunc: SuppressIfAny(
				IgnoreChangeToCurrentSnowflakePlainValueInOutputWithMapping("show_output", "auto_refresh_status", func(x any) any {
					return len(x.([]any)) != 0
				}),
			),
		},
		ParametersAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW PARAMETERS IN ICEBERG TABLE` for the given Iceberg table.",
			Elem:        &schema.Resource{Schema: schemas.ShowIcebergTableExternallyManagedParametersSchema},
		},
	},
	icebergTableExternalManagedParametersSchema(),
)

func IcebergTableFromDeltaFiles() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.IcebergTableFromDeltaFilesResource), TrackingCreateWrapper(resources.IcebergTableFromDeltaFiles, CreateIcebergTableFromDeltaFiles)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.IcebergTableFromDeltaFilesResource), TrackingReadWrapper(resources.IcebergTableFromDeltaFiles, ReadIcebergTableFromDeltaFilesFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.IcebergTableFromDeltaFilesResource), TrackingUpdateWrapper(resources.IcebergTableFromDeltaFiles, UpdateIcebergTableFromDeltaFiles)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.IcebergTableFromDeltaFilesResource), TrackingDeleteWrapper(resources.IcebergTableFromDeltaFiles, icebergTableDeleteFunc())),

		Description: "Resource used to manage an Iceberg table whose metadata is created from Delta table files in an external volume. For more information, check [the official documentation](https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-delta).",

		Schema: icebergTableFromDeltaFilesSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.IcebergTableFromDeltaFiles, importIcebergTable),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(icebergTableFromDeltaFilesSchema, ShowOutputAttributeName, "comment", "auto_refresh"),
			ComputedIfAnyAttributeChanged(icebergTableFromDeltaFilesSchema, ParametersAttributeName, "external_volume", "catalog", "replace_invalid_characters"),
			icebergTableExternalManagedParametersCustomDiff,
		),
	}
}

func CreateIcebergTableFromDeltaFiles(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	baseLocation := d.Get("base_location").(string)

	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)
	req := sdk.NewCreateFromDeltaLakeIcebergTableRequest(id, baseLocation)

	if err := stringAttributeCreate(d, "comment", &req.Comment); err != nil {
		return diag.FromErr(err)
	}
	if err := booleanStringAttributeCreate(d, "auto_refresh", &req.AutoRefresh); err != nil {
		return diag.FromErr(err)
	}
	if diags := handleIcebergTableParametersCreate(d, &req.ExternalVolume, &req.Catalog, &req.ReplaceInvalidCharacters); diags.HasError() {
		return diags
	}

	if err := client.IcebergTables.CreateFromDeltaLake(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Iceberg table from Delta files (%s), err = %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadIcebergTableFromDeltaFilesFunc(false)(ctx, d, meta)
}

func ReadIcebergTableFromDeltaFilesFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		return readIcebergTable(ctx, d, meta, func(d *schema.ResourceData, table *sdk.IcebergTable, _ []sdk.IcebergTableDetails) error {
			var baseLocation string
			if table.BaseLocation != nil {
				baseLocation = *table.BaseLocation
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
			}
			return errors.Join(
				d.Set("base_location", baseLocation),
			)
		})
	}
}

func UpdateIcebergTableFromDeltaFiles(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	return ReadIcebergTableFromDeltaFilesFunc(false)(ctx, d, meta)
}

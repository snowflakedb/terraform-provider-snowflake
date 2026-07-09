package resources

import (
	"context"
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

var icebergTableFromFilesSchema = collections.MergeMaps(
	icebergTableCommonSchema(),
	map[string]*schema.Schema{
		"metadata_file_path": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  externalChangesNotDetectedFieldDescription("Specifies the relative path of the Iceberg metadata file in the external volume. Cannot be changed after creation."),
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

func IcebergTableFromFiles() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.IcebergTableFromFilesResource), TrackingCreateWrapper(resources.IcebergTableFromFiles, CreateIcebergTableFromFiles)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.IcebergTableFromFilesResource), TrackingReadWrapper(resources.IcebergTableFromFiles, ReadIcebergTableFromFiles)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.IcebergTableFromFilesResource), TrackingUpdateWrapper(resources.IcebergTableFromFiles, UpdateIcebergTableFromFiles)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.IcebergTableFromFilesResource), TrackingDeleteWrapper(resources.IcebergTableFromFiles, icebergTableDeleteFunc())),

		Description: "Resource used to manage an Iceberg table whose metadata is created from an existing Apache Iceberg metadata file in an external volume. For more information, check [the official documentation](https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-iceberg-files).",

		Schema: icebergTableFromFilesSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.IcebergTableFromFiles, importIcebergTable),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(icebergTableFromFilesSchema, ShowOutputAttributeName, "comment"),
			ComputedIfAnyAttributeChanged(icebergTableFromFilesSchema, ParametersAttributeName, "external_volume", "catalog", "replace_invalid_characters"),
			icebergTableExternalManagedParametersCustomDiff,
		),
	}
}

func CreateIcebergTableFromFiles(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	metadataFilePath := d.Get("metadata_file_path").(string)

	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)
	req := sdk.NewCreateFromIcebergFilesIcebergTableRequest(id, metadataFilePath)

	if err := stringAttributeCreate(d, "comment", &req.Comment); err != nil {
		return diag.FromErr(err)
	}
	if diags := handleIcebergTableParametersCreate(d, &req.ExternalVolume, &req.Catalog, &req.ReplaceInvalidCharacters); diags.HasError() {
		return diags
	}

	if err := client.IcebergTables.CreateFromIcebergFiles(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Iceberg table from files (%s), err = %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadIcebergTableFromFiles(ctx, d, meta)
}

func ReadIcebergTableFromFiles(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// TODO(SNOW-3663247): Read metadata_file_path from https://docs.snowflake.com/en/sql-reference/functions/system_get_iceberg_table_information.
	// Currently, it returns the storage location concatenated with the path, so it's impossible to get the original path.
	// `metadata_file_path` is not exposed by SHOW or DESCRIBE.
	return readIcebergTable(ctx, d, meta, nil)
}

func UpdateIcebergTableFromFiles(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	if err := applyIcebergTableAlter(ctx, client, id, set, unset); err != nil {
		return diag.FromErr(err)
	}

	return ReadIcebergTableFromFiles(ctx, d, meta)
}

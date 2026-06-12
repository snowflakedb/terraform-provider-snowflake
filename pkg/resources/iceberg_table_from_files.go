package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

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

var icebergTableFromFilesSchema = map[string]*schema.Schema{
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the Iceberg table."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the Iceberg table."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the Iceberg table; must be unique for the schema in which the Iceberg table is created."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"metadata_file_path": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringIsNotEmpty,
		Description:  "Specifies the relative path of the Iceberg metadata file in the external volume. Cannot be changed after creation.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the Iceberg table.",
	},
	FullyQualifiedNameAttributeName: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW ICEBERG TABLES` for the given Iceberg table. Note that this value will be only recomputed whenever values of fields affecting the output change.",
		Elem:        &schema.Resource{Schema: schemas.ShowIcebergTableSchema},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE ICEBERG TABLE` for the given Iceberg table.",
		Elem:        &schema.Resource{Schema: schemas.DescribeIcebergTableSchema},
	},
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN ICEBERG TABLE` for the given Iceberg table.",
		Elem:        &schema.Resource{Schema: schemas.ShowIcebergTableParametersSchema},
	},
}

func IcebergTableFromFiles() *schema.Resource {
	allSchema := collections.MergeMaps(icebergTableFromFilesSchema, icebergTableFromFilesParametersSchema)
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.IcebergTables.DropSafely
		},
	)
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.IcebergTableFromFilesResource), TrackingCreateWrapper(resources.IcebergTableFromFiles, CreateIcebergTableFromFiles)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.IcebergTableFromFilesResource), TrackingReadWrapper(resources.IcebergTableFromFiles, ReadIcebergTableFromFiles)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.IcebergTableFromFilesResource), TrackingUpdateWrapper(resources.IcebergTableFromFiles, UpdateIcebergTableFromFiles)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.IcebergTableFromFilesResource), TrackingDeleteWrapper(resources.IcebergTableFromFiles, deleteFunc)),

		Description: "Resource used to manage an Iceberg table whose metadata is created from an existing Apache Iceberg metadata file in an external volume. For more information, check [the official documentation](https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-iceberg-files).",

		Schema: allSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.IcebergTableFromFiles, ImportIcebergTableFromFiles),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(icebergTableFromFilesSchema, ShowOutputAttributeName, "comment"),
			ComputedIfAnyAttributeChanged(icebergTableFromFilesParametersSchema, ParametersAttributeName, "external_volume", "catalog", "replace_invalid_characters"),
			icebergTableFromFilesParametersCustomDiff,
		),
	}
}

func ImportIcebergTableFromFiles(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	_, err = client.IcebergTables.ShowByIDSafely(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
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
	for _, diags := range []diag.Diagnostics{
		handleParameterCreateWithMapping(d, sdk.IcebergTableParameterExternalVolume, &req.ExternalVolume, sdk.ParseAccountObjectIdentifier),
		handleParameterCreateWithMapping(d, sdk.IcebergTableParameterCatalog, &req.Catalog, sdk.ParseAccountObjectIdentifier),
		handleParameterCreate(d, sdk.IcebergTableParameterReplaceInvalidCharacters, &req.ReplaceInvalidCharacters),
	} {
		if len(diags) > 0 {
			return diags
		}
	}

	if err := client.IcebergTables.CreateFromIcebergFiles(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Iceberg table from files (%s), err = %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadIcebergTableFromFiles(ctx, d, meta)
}

func ReadIcebergTableFromFiles(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	table, err := client.IcebergTables.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query Iceberg table. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Iceberg table id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	details, err := client.IcebergTables.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe Iceberg table (%s), err = %w", id.FullyQualifiedName(), err))
	}

	parameters, err := client.IcebergTables.ShowParameters(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not show parameters for Iceberg table (%s), err = %w", id.FullyQualifiedName(), err))
	}

	var comment string
	if table.Comment != nil {
		comment = *table.Comment
	}

	providerCtx := meta.(*provider.Context)
	errs := errors.Join(
		d.Set("database", table.DatabaseName),
		d.Set("schema", table.SchemaName),
		d.Set("name", table.Name),
		d.Set("comment", comment),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.IcebergTableToSchema(table)}),
		d.Set(DescribeOutputAttributeName, schemas.IcebergTableDetailsToSchema(details)),
		d.Set(ParametersAttributeName, []map[string]any{schemas.IcebergTableParametersToSchema(parameters, providerCtx)}),
		handleIcebergTableFromFilesParameterRead(d, parameters),
	)
	// `metadata_file_path` is not exposed by SHOW or DESCRIBE.
	return diag.FromErr(errs)
}

func UpdateIcebergTableFromFiles(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewIcebergTableSetPropertiesRequest()
	unset := sdk.NewIcebergTableUnsetPropertiesRequest()

	if errs := errors.Join(
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	); errs != nil {
		return diag.FromErr(errs)
	}
	if diags := handleParameterUpdate(d, sdk.IcebergTableParameterReplaceInvalidCharacters, &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters); len(diags) > 0 {
		return diags
	}

	alterReq := sdk.NewAlterIcebergTableRequest(id)
	if !reflect.DeepEqual(*set, *sdk.NewIcebergTableSetPropertiesRequest()) {
		if err := client.IcebergTables.Alter(ctx, alterReq.WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if !reflect.DeepEqual(*unset, *sdk.NewIcebergTableUnsetPropertiesRequest()) {
		if err := client.IcebergTables.Alter(ctx, alterReq.WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadIcebergTableFromFiles(ctx, d, meta)
}

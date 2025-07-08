package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var externalTableSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier for the external table; must be unique for the database and schema in which the externalTable is created.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the external table.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the external table.",
	},
	"table_format": {
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  `Identifies the external table table type. For now, only "delta" for Delta Lake table format is supported.`,
		ValidateFunc: validation.StringInSlice([]string{"delta"}, true),
	},
	"column": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		ForceNew:    true,
		Description: "Definitions of a column to create in the external table. Minimum one required.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Column name",
					ForceNew:    true,
				},
				"type": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Column type, e.g. VARIANT",
					ForceNew:         true,
					ValidateDiagFunc: IsDataTypeValid,
				},
				"as": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "String that specifies the expression for the column. When queried, the column returns results derived from this expression.",
					ForceNew:    true,
				},
			},
		},
	},
	"location": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies a location for the external table, using its FQDN. You can hardcode it (`\"@MYDB.MYSCHEMA.MYSTAGE\"`), or populate dynamically (`\"@${snowflake_stage.mystage.fully_qualified_name}\"`)",
	},
	"file_format": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the file format for the external table.",
	},
	"pattern": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the file names and/or paths on the external stage to match.",
	},
	"aws_sns_topic": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the aws sns topic for the external table.",
	},
	"partition_by": {
		Type:        schema.TypeList,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		ForceNew:    true,
		Description: "Specifies any partition columns to evaluate for the external table.",
	},
	"refresh_on_create": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies weather to refresh when an external table is created.",
		Default:     true,
		ForceNew:    true,
	},
	"auto_refresh": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether to automatically refresh the external table metadata once, immediately after the external table is created.",
		Default:     true,
		ForceNew:    true,
	},
	"copy_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies to retain the access permissions from the original table when an external table is recreated using the CREATE OR REPLACE TABLE variant",
		Default:     false,
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies a comment for the external table.",
	},
	"owner": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the role that owns the external table.",
	},
	"tag":                           tagReferenceSchema,
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func ExternalTable() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErr[sdk.SchemaObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.ExternalTables.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ExternalTableResource), TrackingCreateWrapper(resources.ExternalTable, CreateExternalTable)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ExternalTableResource), TrackingReadWrapper(resources.ExternalTable, ReadExternalTable)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ExternalTableResource), TrackingUpdateWrapper(resources.ExternalTable, UpdateExternalTable)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ExternalTableResource), TrackingDeleteWrapper(resources.ExternalTable, deleteFunc)),

		Schema: externalTableSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateExternalTable implements schema.CreateFunc.
func CreateExternalTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schema, name)
	location := d.Get("location").(string)
	fileFormat := d.Get("file_format").(string)

	tableColumns := d.Get("column").([]any)
	columnRequests := make([]*sdk.ExternalTableColumnRequest, len(tableColumns))
	for i, col := range tableColumns {
		columnDef := map[string]string{}
		for key, val := range col.(map[string]any) {
			columnDef[key] = val.(string)
		}
		columnRequests[i] = sdk.NewExternalTableColumnRequest(
			columnDef["name"],
			sdk.DataType(columnDef["type"]),
			columnDef["as"],
		)
	}
	autoRefresh := d.Get("auto_refresh").(bool)
	refreshOnCreate := d.Get("refresh_on_create").(bool)
	copyGrants := d.Get("copy_grants").(bool)

	var partitionBy []string
	if v, ok := d.GetOk("partition_by"); ok {
		partitionBy = expandStringList(v.([]any))
	}

	pattern, hasPattern := d.GetOk("pattern")
	awsSnsTopic, hasAwsSnsTopic := d.GetOk("aws_sns_topic")
	comment, hasComment := d.GetOk("comment")

	var tagAssociationRequests []*sdk.TagAssociationRequest
	if _, ok := d.GetOk("tag"); ok {
		tagAssociations := getPropertyTags(d, "tag")
		tagAssociationRequests = make([]*sdk.TagAssociationRequest, len(tagAssociations))
		for i, t := range tagAssociations {
			tagAssociationRequests[i] = sdk.NewTagAssociationRequest(t.Name, t.Value)
		}
	}

	switch {
	case d.Get("table_format").(string) == "delta":
		req := sdk.NewCreateDeltaLakeExternalTableRequest(id, location).
			WithColumns(columnRequests).
			WithPartitionBy(partitionBy).
			WithRefreshOnCreate(refreshOnCreate).
			WithAutoRefresh(autoRefresh).
			WithRawFileFormat(fileFormat).
			WithCopyGrants(copyGrants).
			WithTag(tagAssociationRequests)
		if hasComment {
			req = req.WithComment(comment.(string))
		}
		err := client.ExternalTables.CreateDeltaLake(ctx, req)
		if err != nil {
			return diag.FromErr(err)
		}
	default:
		req := sdk.NewCreateExternalTableRequest(id, location).
			WithColumns(columnRequests).
			WithPartitionBy(partitionBy).
			WithRefreshOnCreate(refreshOnCreate).
			WithAutoRefresh(autoRefresh).
			WithRawFileFormat(fileFormat).
			WithCopyGrants(copyGrants).
			WithTag(tagAssociationRequests)
		if hasPattern {
			req = req.WithPattern(pattern.(string))
		}
		if hasAwsSnsTopic {
			req = req.WithAwsSnsTopic(awsSnsTopic.(string))
		}
		if hasComment {
			req = req.WithComment(comment.(string))
		}
		err := client.ExternalTables.Create(ctx, req)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadExternalTable(ctx, d, meta)
}

// ReadExternalTable implements schema.ReadFunc.
func ReadExternalTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.SchemaObjectIdentifier)

	externalTable, err := client.ExternalTables.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query external table. Marking the resource as removed.",
					Detail:   fmt.Sprintf("External table id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", externalTable.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owner", externalTable.Owner); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// UpdateExternalTable implements schema.UpdateFunc.
func UpdateExternalTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.SchemaObjectIdentifier)

	if d.HasChange("tag") {
		unsetTags, setTags := GetTagsDiff(d, "tag")

		if len(unsetTags) > 0 {
			err := client.ExternalTables.Alter(ctx, sdk.NewAlterExternalTableRequest(id).WithUnsetTag(unsetTags))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error setting tags on %v, err = %w", d.Id(), err))
			}
		}

		if len(setTags) > 0 {
			tagAssociationRequests := make([]*sdk.TagAssociationRequest, len(setTags))
			for i, t := range setTags {
				tagAssociationRequests[i] = sdk.NewTagAssociationRequest(t.Name, t.Value)
			}
			err := client.ExternalTables.Alter(ctx, sdk.NewAlterExternalTableRequest(id).WithSetTag(tagAssociationRequests))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error setting tags on %v, err = %w", d.Id(), err))
			}
		}
	}

	return ReadExternalTable(ctx, d, meta)
}

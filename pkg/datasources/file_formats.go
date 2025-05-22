package datasources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var fileFormatsSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the file formats from.",
	},
	"file_formats": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The file formats in the schema",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"database": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"schema": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"comment": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"format_type": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func FileFormats() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.FileFormatsDatasource), TrackingReadWrapper(datasources.FileFormats, ReadFileFormats)),
		Schema:      fileFormatsSchema,
	}
}

func ReadFileFormats(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	result, err := client.FileFormats.Show(ctx, &sdk.ShowFileFormatsOptions{
		In: &sdk.In{
			Schema: sdk.NewDatabaseObjectIdentifier(databaseName, schemaName),
		},
	})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	fileFormats := []map[string]interface{}{}

	for _, fileFormat := range result {
		fileFormatMap := map[string]interface{}{}

		fileFormatMap["name"] = fileFormat.Name.Name()
		fileFormatMap["database"] = fileFormat.Name.DatabaseName()
		fileFormatMap["schema"] = fileFormat.Name.SchemaName()
		fileFormatMap["comment"] = fileFormat.Comment
		fileFormatMap["format_type"] = fileFormat.Type

		fileFormats = append(fileFormats, fileFormatMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return diag.FromErr(d.Set("file_formats", fileFormats))
}

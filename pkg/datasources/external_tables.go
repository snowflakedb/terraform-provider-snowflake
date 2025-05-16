package datasources

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var externalTablesSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the external tables from.",
	},
	"external_tables": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The external tables in the schema",
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
			},
		},
	},
}

func ExternalTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.ExternalTablesDatasource), TrackingReadWrapper(datasources.ExternalTables, ReadExternalTables)),
		Schema:      externalTablesSchema,
	}
}

func ReadExternalTables(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	schemaId := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)
	showIn := sdk.NewShowExternalTableInRequest().WithSchema(schemaId)
	externalTables, err := client.ExternalTables.Show(ctx, sdk.NewShowExternalTableRequest().WithIn(*showIn))
	if err != nil {
		log.Printf("[DEBUG] failed when searching external tables in schema (%s), err = %s", schemaId.FullyQualifiedName(), err.Error())
		d.SetId("")
		return nil
	}

	externalTablesObjects := make([]map[string]any, len(externalTables))
	for i, externalTable := range externalTables {
		externalTablesObjects[i] = map[string]any{
			"name":     externalTable.Name,
			"database": externalTable.DatabaseName,
			"schema":   externalTable.SchemaName,
			"comment":  externalTable.Comment,
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(schemaId))

	return diag.FromErr(d.Set("external_tables", externalTablesObjects))
}

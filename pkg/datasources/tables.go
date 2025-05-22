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

var tablesSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the tables from.",
	},
	"tables": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The tables in the schema",
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

func Tables() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.TablesDatasource), TrackingReadWrapper(datasources.Tables, ReadTables)),
		Schema:      tablesSchema,
	}
}

func ReadTables(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	schemaId := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)
	extractedTables, err := client.Tables.Show(ctx, sdk.NewShowTableRequest().WithIn(
		&sdk.In{Schema: schemaId},
	))
	if err != nil {
		log.Printf("[DEBUG] failed when searching tables in schema (%s), err = %s", schemaId.FullyQualifiedName(), err.Error())
		d.SetId("")
		return nil
	}

	tables := make([]map[string]any, 0)

	for _, extractedTable := range extractedTables {
		if extractedTable.IsExternal {
			continue
		}

		table := map[string]any{
			"name":     extractedTable.Name,
			"database": extractedTable.DatabaseName,
			"schema":   extractedTable.SchemaName,
			"comment":  extractedTable.Comment,
		}

		tables = append(tables, table)
	}

	d.SetId(helpers.EncodeSnowflakeID(databaseName, schemaName))
	return diag.FromErr(d.Set("tables", tables))
}

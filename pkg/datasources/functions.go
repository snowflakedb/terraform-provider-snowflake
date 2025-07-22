package datasources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var functionsSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the functions from.",
	},
	"functions": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The functions in the schema",
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
				"argument_types": {
					Type:     schema.TypeList,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Optional: true,
					Computed: true,
				},
				"return_type": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func Functions() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.FunctionsDatasource), TrackingReadWrapper(datasources.Functions, ReadContextFunctions)),
		Schema:      functionsSchema,
	}
}

func ReadContextFunctions(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	request := sdk.NewShowFunctionRequest()
	request.WithIn(sdk.ExtendedIn{In: sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)}})
	functions, err := client.Functions.Show(ctx, request)
	if err != nil {
		id := d.Id()

		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Unable to parse functions in schema (%s)", id),
				Detail:   "See our document on design decisions for functions: <LINK (coming soon)>",
			},
		}
	}

	entities := []map[string]interface{}{}
	for _, item := range functions {
		m := map[string]interface{}{}
		m["name"] = item.Name
		m["database"] = databaseName
		m["schema"] = schemaName
		m["comment"] = item.Description
		m["argument_types"] = collections.Map(item.ArgumentsOld, func(a sdk.DataType) string {
			return string(a)
		})
		m["return_type"] = string(item.ReturnTypeOld)

		entities = append(entities, m)
	}
	d.SetId(helpers.EncodeSnowflakeID(databaseName, schemaName))
	if err := d.Set("functions", entities); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

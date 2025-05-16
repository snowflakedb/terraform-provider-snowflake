package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var externalFunctionsSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:         schema.TypeString,
		Optional:     true,
		RequiredWith: []string{"database"},
		Description:  "The schema from which to return the external functions from.",
	},
	"external_functions": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The external functions in the schema",
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
				"language": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func ExternalFunctions() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.ExternalFunctionsDatasource), TrackingReadWrapper(datasources.ExternalFunctions, ReadContextExternalFunctions)),
		Schema:      externalFunctionsSchema,
	}
}

func ReadContextExternalFunctions(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	req := sdk.NewShowExternalFunctionRequest()
	externalFunctions, err := client.ExternalFunctions.Show(ctx, req)
	if err != nil {
		d.SetId("")
		return nil
	}

	externalFunctionsList := []map[string]interface{}{}
	for _, externalFunction := range externalFunctions {
		externalFunctionMap := map[string]interface{}{}
		externalFunctionMap["name"] = externalFunction.Name

		// do we filter by database?
		currentDatabase := strings.Trim(externalFunction.CatalogName, `"`)
		if databaseName != "" {
			if currentDatabase != databaseName {
				continue
			}
			externalFunctionMap["database"] = currentDatabase
		} else {
			externalFunctionMap["database"] = currentDatabase
		}

		// do we filter by schema?
		currentSchema := strings.Trim(externalFunction.SchemaName, `"`)
		if schemaName != "" {
			if currentSchema != schemaName {
				continue
			}
			externalFunctionMap["schema"] = currentSchema
		} else {
			externalFunctionMap["schema"] = currentSchema
		}

		externalFunctionMap["comment"] = externalFunction.Description
		externalFunctionMap["language"] = externalFunction.Language
		externalFunctionsList = append(externalFunctionsList, externalFunctionMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	if err := d.Set("external_functions", externalFunctionsList); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

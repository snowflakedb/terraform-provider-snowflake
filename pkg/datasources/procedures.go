package datasources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var proceduresSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the procedures from.",
	},
	"procedures": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The procedures in the schema",
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

func Procedures() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.ProceduresDatasource), TrackingReadWrapper(datasources.Procedures, ReadContextProcedures)),
		Schema:      proceduresSchema,
	}
}

func ReadContextProcedures(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	req := sdk.NewShowProcedureRequest()
	if databaseName != "" {
		req.WithIn(sdk.ExtendedIn{In: sdk.In{Database: sdk.NewAccountObjectIdentifier(databaseName)}})
	}
	if schemaName != "" {
		req.WithIn(sdk.ExtendedIn{In: sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)}})
	}
	procedures, err := client.Procedures.Show(ctx, req)
	if err != nil {
		id := fmt.Sprintf(`%v|%v`, databaseName, schemaName)

		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Unable to parse procedures in schema (%s)", id),
				Detail:   "See our document on design decisions for procedures: <LINK (coming soon)>",
			},
		}
	}
	proceduresList := []map[string]interface{}{}

	for _, procedure := range procedures {
		procedureMap := map[string]interface{}{}
		procedureMap["name"] = procedure.Name
		procedureMap["database"] = procedure.CatalogName
		procedureMap["schema"] = procedure.SchemaName
		procedureMap["comment"] = procedure.Description
		procedureMap["argument_types"] = collections.Map(procedure.ArgumentsOld, func(a sdk.DataType) string {
			return string(a)
		})
		procedureMap["return_type"] = string(procedure.ReturnTypeOld)
		proceduresList = append(proceduresList, procedureMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	if err := d.Set("procedures", proceduresList); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

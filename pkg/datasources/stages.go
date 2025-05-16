package datasources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var stagesSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the stages from.",
	},
	"stages": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The stages in the schema",
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
				"storage_integration": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func Stages() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.StagesDatasource), TrackingReadWrapper(datasources.Stages, ReadStages)),
		Schema:      stagesSchema,
	}
}

func ReadStages(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	stages, err := client.Stages.Show(ctx, sdk.NewShowStageRequest().WithIn(
		sdk.In{
			Schema: sdk.NewDatabaseObjectIdentifier(databaseName, schemaName),
		},
	))
	if err != nil {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to query stages",
				Detail:   fmt.Sprintf("DatabaseName: %s, SchemaName: %s, Err: %s", databaseName, schemaName, err),
			},
		}
	}

	stagesList := make([]map[string]any, len(stages))
	for i, stage := range stages {
		stagesList[i] = map[string]any{
			"name":                stage.Name,
			"database":            stage.DatabaseName,
			"schema":              stage.SchemaName,
			"comment":             stage.Comment,
			"storage_integration": stage.StorageIntegration,
		}
	}

	if err := d.Set("stages", stagesList); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set stages",
				Detail:   fmt.Sprintf("Err: %s", err),
			},
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(databaseName, schemaName))

	return nil
}

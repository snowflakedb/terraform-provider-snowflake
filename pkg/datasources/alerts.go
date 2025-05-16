package datasources

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var alertsSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The database from which to return the alerts from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The schema from which to return the alerts from.",
	},
	"pattern": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the command output by object name.",
	},
	"alerts": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Lists alerts for the current/specified database or schema, or across the entire account.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the alert.",
				},
				"database_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Database in which the alert is stored.",
				},
				"schema_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Schema in which the alert is stored.",
				},
				"comment": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Comment for the alert.",
				},
				"owner": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Role that owns the alert (i.e. has the OWNERSHIP privilege on the alert)",
				},
				"condition": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The text of the SQL statement that serves as the condition when the alert should be triggered.",
				},
				"action": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The text of the SQL statement that should be executed when the alert is triggered.",
				},
			},
		},
	},
}

// Alerts Snowflake Roles resource.
func Alerts() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.AlertsDatasource), TrackingReadWrapper(datasources.Alerts, ReadAlerts)),
		Schema:      alertsSchema,
	}
}

// ReadAlerts Reads the database metadata information.
func ReadAlerts(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	d.SetId("alerts_read")

	opts := sdk.ShowAlertOptions{}

	if v, ok := d.GetOk("pattern"); ok {
		alertPattern := v.(string)
		opts.Like = &sdk.Like{Pattern: &alertPattern}
	}

	if v, ok := d.GetOk("database"); ok {
		databaseName := v.(string)

		if v, ok := d.GetOk("schema"); ok {
			schemaName := v.(string)
			opts.In = &sdk.In{
				Schema: sdk.NewDatabaseObjectIdentifier(databaseName, schemaName),
			}
		} else {
			opts.In = &sdk.In{
				Database: sdk.NewAccountObjectIdentifier(databaseName),
			}
		}
	}

	listAlerts, err := client.Alerts.Show(ctx, &opts)
	if err != nil {
		log.Printf("[DEBUG] failed to list alerts in schema (%s)", d.Id())
		d.SetId("")
		return diag.FromErr(err)
	}

	alerts := make([]map[string]any, 0, len(listAlerts))
	for _, alert := range listAlerts {
		alertMap := map[string]any{}
		alertMap["name"] = alert.Name
		alertMap["comment"] = alert.Comment
		alertMap["owner"] = alert.Owner
		alerts = append(alerts, alertMap)
	}

	if err := d.Set("alerts", alerts); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var alertSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the alert; must be unique for the database and schema in which the alert is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the alert.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the alert.",
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the alert.",
	},
	"warehouse": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The warehouse the alert will use.",
		ForceNew:    true,
	},
	"alert_schedule": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "The schedule for periodically running an alert.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cron": {
					Type:          schema.TypeList,
					Optional:      true,
					MaxItems:      1,
					ConflictsWith: []string{"alert_schedule.interval"},
					Description:   "Specifies the cron expression for the alert. The cron expression must be in the following format: \"minute hour day-of-month month day-of-week\". The following values are supported: minute: 0-59 hour: 0-23 day-of-month: 1-31 month: 1-12 day-of-week: 0-6 (0 is Sunday)",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"expression": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the cron expression for the alert. The cron expression must be in the following format: \"minute hour day-of-month month day-of-week\". The following values are supported: minute: 0-59 hour: 0-23 day-of-month: 1-31 month: 1-12 day-of-week: 0-6 (0 is Sunday)",
							},
							"time_zone": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the time zone for alert refresh.",
							},
						},
					},
				},
				"interval": {
					Type:          schema.TypeInt,
					Optional:      true,
					ConflictsWith: []string{"alert_schedule.cron"},
					Description:   "Specifies the interval in minutes for the alert schedule. The interval must be greater than 0 and less than 1440 (24 hours).",
				},
			},
		},
	},
	"condition": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "The SQL statement that represents the condition for the alert. (SELECT, SHOW, CALL)",
		ForceNew:         false,
		DiffSuppressFunc: DiffSuppressStatement,
	},
	"action": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "The SQL statement that should be executed if the condition returns one or more rows.",
		ForceNew:         false,
		DiffSuppressFunc: DiffSuppressStatement,
	},
	"enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies if an alert should be 'started' (enabled) after creation or should remain 'suspended' (default).",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func Alert() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErr[sdk.SchemaObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] { return client.Alerts.DropSafely },
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.AlertResource), TrackingCreateWrapper(resources.Alert, CreateAlert)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.AlertResource), TrackingReadWrapper(resources.Alert, ReadAlert)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.AlertResource), TrackingUpdateWrapper(resources.Alert, UpdateAlert)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.AlertResource), TrackingDeleteWrapper(resources.Alert, deleteFunc)),

		Schema: alertSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// ReadAlert implements schema.ReadContextFunc.
func ReadAlert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	alert, err := client.Alerts.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query alert. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Alert id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	if err := d.Set("enabled", alert.State == sdk.AlertStateStarted); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", alert.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database", alert.DatabaseName); err != nil {
		return diag.FromErr(err)
	}

	alertSchedule := alert.Schedule
	if alertSchedule != "" {
		if strings.Contains(alertSchedule, "MINUTE") {
			interval, err := strconv.Atoi(strings.TrimSuffix(alertSchedule, " MINUTE"))
			if err != nil {
				return diag.FromErr(err)
			}
			err = d.Set("alert_schedule", []interface{}{
				map[string]interface{}{
					"interval": interval,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			repScheduleParts := strings.Split(alertSchedule, " ")
			timeZone := repScheduleParts[len(repScheduleParts)-1]
			expression := strings.TrimSuffix(strings.TrimPrefix(alertSchedule, "USING CRON "), " "+timeZone)
			err = d.Set("alert_schedule", []interface{}{
				map[string]interface{}{
					"cron": []interface{}{
						map[string]interface{}{
							"expression": expression,
							"time_zone":  timeZone,
						},
					},
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if err := d.Set("schema", alert.SchemaName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("warehouse", alert.Warehouse); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", alert.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("condition", alert.Condition); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("action", alert.Action); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// CreateAlert implements schema.CreateContextFunc.
func CreateAlert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)

	objectIdentifier := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	alertSchedule := getAlertSchedule(d.Get("alert_schedule"))

	warehouseName := d.Get("warehouse").(string)
	warehouse := sdk.NewAccountObjectIdentifier(warehouseName)

	opts := &sdk.CreateAlertOptions{}

	if v, ok := d.GetOk("comment"); ok {
		opts.Comment = sdk.String(v.(string))
	}

	condition := d.Get("condition").(string)

	action := d.Get("action").(string)

	err := client.Alerts.Create(ctx, objectIdentifier, warehouse, alertSchedule, condition, action, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(objectIdentifier))

	enabled := d.Get("enabled").(bool)
	var diags diag.Diagnostics
	if enabled {
		opts := sdk.AlterAlertOptions{Action: &sdk.AlertActionResume}
		err := client.Alerts.Alter(ctx, objectIdentifier, &opts)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Enabling alert: %s", err.Error()),
			})
		}
	}
	return append(diags, ReadAlert(ctx, d, meta)...)
}

func getAlertSchedule(v interface{}) string {
	var alertSchedule string
	schedule := v.([]interface{})[0].(map[string]interface{})
	if v, ok := schedule["cron"]; ok {
		c := v.([]interface{})
		if len(c) > 0 {
			cron := c[0].(map[string]interface{})
			cronExpression := cron["expression"].(string)
			timeZone := cron["time_zone"].(string)
			alertSchedule = fmt.Sprintf("USING CRON %s %s", cronExpression, timeZone)
		}
	}
	if v, ok := schedule["interval"]; ok {
		interval := v.(int)
		if interval > 0 {
			alertSchedule = fmt.Sprintf("%s MINUTE", strconv.Itoa(interval))
		}
	}
	return alertSchedule
}

// UpdateAlert implements schema.UpdateContextFunc.
func UpdateAlert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	enabled := d.Get("enabled").(bool)
	if d.HasChanges("enabled", "warehouse", "alert_schedule", "condition", "action", "comment") {
		if err := waitSuspendAlert(ctx, client, objectIdentifier); err != nil {
			log.Printf("[WARN] failed to suspend alert %s", objectIdentifier.Name())
		}
	}

	opts := &sdk.AlterAlertOptions{
		Set:   &sdk.AlertSet{},
		Unset: &sdk.AlertUnset{},
	}
	runSetStatement := false

	if d.HasChange("warehouse") {
		runSetStatement = true
		_, v := d.GetChange("warehouse")
		warehouseName := v.(string)
		warehouse := sdk.NewAccountObjectIdentifier(warehouseName)
		opts.Set.Warehouse = &warehouse
	}

	if d.HasChange("alert_schedule") {
		runSetStatement = true
		_, v := d.GetChange("alert_schedule")
		alertSchedule := getAlertSchedule(v)
		opts.Set.Schedule = &alertSchedule
	}

	if d.HasChange("comment") {
		_, v := d.GetChange("comment")
		runSetStatement = true
		newComment := v.(string)
		opts.Set.Comment = &newComment
	}

	if runSetStatement {
		setOptions := &sdk.AlterAlertOptions{Set: opts.Set}
		err := client.Alerts.Alter(ctx, objectIdentifier, setOptions)
		if err != nil {
			return diag.Errorf("error updating alert %v: %v", objectIdentifier.Name(), err)
		}
	}

	if d.HasChange("condition") {
		condition := d.Get("condition").(string)
		alterOptions := &sdk.AlterAlertOptions{}
		alterOptions.ModifyCondition = &[]string{condition}
		err := client.Alerts.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.Errorf("error updating schedule on condition %v: %v", objectIdentifier.Name(), err)
		}
	}

	if d.HasChange("action") {
		action := d.Get("action").(string)
		alterOptions := &sdk.AlterAlertOptions{}
		alterOptions.ModifyAction = &action
		err := client.Alerts.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.Errorf("error updating schedule on action %v: %v", objectIdentifier.Name(), err)
		}
	}
	var diags diag.Diagnostics
	if enabled {
		if err := waitResumeAlert(ctx, client, objectIdentifier); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("failed to resume alert %s", objectIdentifier.Name()),
			})
		}
	} else {
		if err := waitSuspendAlert(ctx, client, objectIdentifier); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("failed to suspend alert %s", objectIdentifier.Name()),
			})
		}
	}
	return append(diags, ReadAlert(ctx, d, meta)...)
}

func waitResumeAlert(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	resumeAlert := func() (error, bool) {
		opts := sdk.AlterAlertOptions{Action: &sdk.AlertActionResume}
		err := client.Alerts.Alter(ctx, id, &opts)
		if err != nil {
			return fmt.Errorf("error resuming alert %v err = %w", id.Name(), err), false
		}
		alert, err := client.Alerts.ShowByID(ctx, id)
		if err != nil {
			return err, false
		}
		return nil, alert.State == sdk.AlertStateStarted
	}
	return util.Retry(5, 10*time.Second, resumeAlert)
}

func waitSuspendAlert(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	suspendAlert := func() (error, bool) {
		opts := sdk.AlterAlertOptions{Action: &sdk.AlertActionSuspend}
		err := client.Alerts.Alter(ctx, id, &opts)
		if err != nil {
			return fmt.Errorf("error suspending alert %v err = %w", id.Name(), err), false
		}
		alert, err := client.Alerts.ShowByID(ctx, id)
		if err != nil {
			return err, false
		}
		return nil, alert.State == sdk.AlertStateSuspended
	}
	return util.Retry(5, 10*time.Second, suspendAlert)
}

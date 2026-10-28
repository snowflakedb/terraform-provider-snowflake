package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	alertCreateExt        = PreviewFeatureCreateContextWrapper(string(previewfeatures.AlertResource), TrackingCreateWrapper(resources.Alert, CreateAlert))
	alertReadExt          = PreviewFeatureReadContextWrapper(string(previewfeatures.AlertResource), TrackingReadWrapper(resources.Alert, ReadAlertFunc(true)))
	alertUpdateExt        = PreviewFeatureUpdateContextWrapper(string(previewfeatures.AlertResource), TrackingUpdateWrapper(resources.Alert, UpdateAlert))
	alertDescriptionExt   = ""
	alertCustomizeDiffExt schema.CustomizeDiffFunc
	alertImporterExt      = &schema.ResourceImporter{
		StateContext: schema.ImportStatePassthroughContext,
	}
)

// =============================================================================
// Request / mapping / Read-Update helpers
// =============================================================================

type alertAlterRequestsExt struct {
	set             *sdk.AlertSetRequest
	modifyCondition *string
	modifyAction    *string
}

func alertNewCreateRequestExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier) (*sdk.CreateAlertRequest, error) {
	warehouse := sdk.NewAccountObjectIdentifier(d.Get("warehouse").(string))
	alertSchedule := getAlertSchedule(d.Get("alert_schedule"))
	condition := d.Get("condition").(string)
	action := d.Get("action").(string)
	return sdk.NewCreateAlertRequest(id, warehouse, alertSchedule, sdk.NewAlertConditionFromString(condition), action), nil
}

func alertSetConfigFieldsExt(d *schema.ResourceData, alert *sdk.Alert) error {
	alertSchedule, err := alertScheduleToState(alert.Schedule)
	if err != nil {
		return err
	}
	return errors.Join(
		d.Set("enabled", alert.State == sdk.AlertStateStarted),
		d.Set("name", alert.Name),
		d.Set("database", alert.DatabaseName),
		setAlertSchedule(d, alertSchedule),
		d.Set("schema", alert.SchemaName),
		d.Set("warehouse", alert.Warehouse),
		d.Set("comment", alert.Comment),
		d.Set("condition", alert.Condition),
		d.Set("action", alert.Action),
	)
}

// setAlertSchedule leaves alert_schedule untouched when SHOW returns no schedule.
func setAlertSchedule(d *schema.ResourceData, alertSchedule []any) error {
	if alertSchedule == nil {
		return nil
	}
	return d.Set("alert_schedule", alertSchedule)
}

func alertScheduleToState(alertSchedule string) ([]any, error) {
	if alertSchedule == "" {
		return nil, nil
	}
	if strings.Contains(alertSchedule, "MINUTE") {
		interval, err := strconv.Atoi(strings.TrimSuffix(alertSchedule, " MINUTE"))
		if err != nil {
			return nil, err
		}
		return []any{
			map[string]any{
				"interval": interval,
			},
		}, nil
	}
	repScheduleParts := strings.Split(alertSchedule, " ")
	timeZone := repScheduleParts[len(repScheduleParts)-1]
	expression := strings.TrimSuffix(strings.TrimPrefix(alertSchedule, "USING CRON "), " "+timeZone)
	return []any{
		map[string]any{
			"cron": []any{
				map[string]any{
					"expression": expression,
					"time_zone":  timeZone,
				},
			},
		},
	}, nil
}

func getAlertSchedule(v any) string {
	var alertSchedule string
	schedule := v.([]any)[0].(map[string]any)
	if v, ok := schedule["cron"]; ok {
		c := v.([]any)
		if len(c) > 0 {
			cron := c[0].(map[string]any)
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

func alertApplySetUnsetExt(d *schema.ResourceData) (*alertAlterRequestsExt, error) {
	reqs := &alertAlterRequestsExt{set: sdk.NewAlertSetRequest()}

	if d.HasChange("warehouse") {
		_, v := d.GetChange("warehouse")
		reqs.set.WithWarehouse(sdk.NewAccountObjectIdentifier(v.(string)))
	}

	if d.HasChange("alert_schedule") {
		_, v := d.GetChange("alert_schedule")
		reqs.set.WithSchedule(getAlertSchedule(v))
	}

	if d.HasChange("comment") {
		_, v := d.GetChange("comment")
		reqs.set.WithComment(v.(string))
	}

	if d.HasChange("condition") {
		reqs.modifyCondition = sdk.String(d.Get("condition").(string))
	}

	if d.HasChange("action") {
		reqs.modifyAction = sdk.String(d.Get("action").(string))
	}
	return reqs, nil
}

// =============================================================================
// SDK calls
// =============================================================================

// alertAfterCreateExt resumes the alert after creation when it is enabled. A failure is only a warning.
func alertAfterCreateExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if !d.Get("enabled").(bool) {
		return nil
	}
	client := meta.(*provider.Context).Client
	if err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).WithAction(sdk.AlertActionResume)); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Enabling alert: %s", err.Error()),
			},
		}
	}
	return nil
}

// alertBeforeAlterExt suspends the alert before any change. A failure is only logged.
func alertBeforeAlterExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if d.HasChanges("enabled", "warehouse", "alert_schedule", "condition", "action", "comment") {
		client := meta.(*provider.Context).Client
		if err := waitSuspendAlert(ctx, client, id); err != nil {
			log.Printf("[WARN] failed to suspend alert %s", id.Name())
		}
	}
	return nil
}

func alertAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier, reqs *alertAlterRequestsExt) error {
	client := meta.(*provider.Context).Client
	if !reflect.DeepEqual(reqs.set, sdk.NewAlertSetRequest()) {
		if err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).WithSet(*reqs.set)); err != nil {
			return fmt.Errorf("error updating alert %v: %w", id.Name(), err)
		}
	}
	if reqs.modifyCondition != nil {
		if err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).WithModifyCondition([]string{*reqs.modifyCondition})); err != nil {
			return fmt.Errorf("error updating schedule on condition %v: %w", id.Name(), err)
		}
	}
	if reqs.modifyAction != nil {
		if err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).WithModifyAction(*reqs.modifyAction)); err != nil {
			return fmt.Errorf("error updating schedule on action %v: %w", id.Name(), err)
		}
	}
	return nil
}

// alertAfterAlterExt resumes or suspends the alert to match enabled. A failure is only a warning.
func alertAfterAlterExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if d.Get("enabled").(bool) {
		if err := waitResumeAlert(ctx, client, id); err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("failed to resume alert %s", id.Name()),
				},
			}
		}
		return nil
	}
	if err := waitSuspendAlert(ctx, client, id); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("failed to suspend alert %s", id.Name()),
			},
		}
	}
	return nil
}

func waitResumeAlert(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	resumeAlert := func() (error, bool) {
		err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).WithAction(sdk.AlertActionResume))
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
		err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).WithAction(sdk.AlertActionSuspend))
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

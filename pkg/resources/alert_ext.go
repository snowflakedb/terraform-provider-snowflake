package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	alertCreateExt      = PreviewFeatureCreateContextWrapper(string(previewfeatures.AlertResource), TrackingCreateWrapper(resources.Alert, CreateAlert))
	alertReadExt        = PreviewFeatureReadContextWrapper(string(previewfeatures.AlertResource), TrackingReadWrapper(resources.Alert, ReadAlertFunc(true)))
	alertUpdateExt      = PreviewFeatureUpdateContextWrapper(string(previewfeatures.AlertResource), TrackingUpdateWrapper(resources.Alert, UpdateAlert))
	alertDescriptionExt = docs.WithPreviewFeatureMark() + "Resource used to manage alert objects. For more information, check [alert documentation](https://docs.snowflake.com/en/sql-reference/sql/create-alert)."
	alertImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.Alert, ImportAlert),
	}
)

type alertAlterRequestsExt struct {
	set             *sdk.AlertSetRequest
	unset           *sdk.AlertUnsetRequest
	modifyCondition *string
	modifyAction    *string
}

func alertNewCreateRequestExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier) (*sdk.CreateAlertRequest, error) {
	req := sdk.NewCreateAlertRequest(id, sdk.NewAlertConditionFromString(d.Get("condition").(string)), d.Get("action").(string))
	if schedule := alertScheduleFromConfig(d.Get("alert_schedule")); schedule != nil {
		req.WithSchedule(*schedule)
	}
	return req, nil
}

func alertSetConfigFieldsExt(_ context.Context, d *schema.ResourceData, _ any, alert *sdk.Alert) error {
	alertSchedule, err := alertScheduleToState(alert.Schedule)
	if err != nil {
		return err
	}
	warehouse := ""
	if alert.Warehouse != nil {
		warehouse = alert.Warehouse.Name()
	}
	comment := ""
	if alert.Comment != nil {
		comment = *alert.Comment
	}
	return errors.Join(
		d.Set("enabled", alert.State == sdk.AlertStateStarted),
		d.Set("alert_schedule", alertSchedule),
		d.Set("warehouse", warehouse),
		d.Set("comment", comment),
		d.Set("condition", alert.Condition),
		d.Set("action", alert.Action),
	)
}

// alertSetFieldsNotSetByReadExt sets the identifier fields on import. Read sets the other fields from SHOW.
func alertSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier, _ *sdk.Alert, _ bool) error {
	return errors.Join(
		d.Set("name", id.Name()),
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
	)
}

// alertScheduleFromConfig turns the alert_schedule block into the SCHEDULE value, or nil when the block is not set.
func alertScheduleFromConfig(v any) *string {
	blocks, ok := v.([]any)
	if !ok || len(blocks) == 0 || blocks[0] == nil {
		return nil
	}
	schedule := blocks[0].(map[string]any)
	if c, ok := schedule["cron"].([]any); ok && len(c) > 0 {
		cron := c[0].(map[string]any)
		return sdk.String(sdk.FormatAlertCronSchedule(cron["expression"].(string), cron["time_zone"].(string)))
	}
	if interval, ok := schedule["interval"].(int); ok && interval > 0 {
		return sdk.String(sdk.FormatAlertIntervalSchedule(interval))
	}
	return nil
}

// alertScheduleToState turns the SCHEDULE value of SHOW into the alert_schedule block, or nil when it is not set.
func alertScheduleToState(alertSchedule *string) ([]any, error) {
	if alertSchedule == nil || *alertSchedule == "" {
		return nil, nil
	}
	parsed, err := sdk.ParseAlertSchedule(*alertSchedule)
	if err != nil {
		return nil, err
	}
	if parsed.CronExpression != "" {
		return []any{
			map[string]any{
				"cron": []any{
					map[string]any{
						"expression": parsed.CronExpression,
						"time_zone":  parsed.TimeZone,
					},
				},
			},
		}, nil
	}
	return []any{map[string]any{"interval": parsed.Minutes}}, nil
}

// alertForceNewIfScheduleAddedOrRemoved recreates the alert when alert_schedule is added or removed.
// Snowflake cannot turn a scheduled alert into one without a schedule, or the other way around.
var alertForceNewIfScheduleAddedOrRemoved = customdiff.ForceNewIfChange("alert_schedule", func(_ context.Context, oldValue, newValue, _ any) bool {
	return (len(oldValue.([]any)) == 0) != (len(newValue.([]any)) == 0)
})

func alertApplySetUnsetExt(d *schema.ResourceData) (*alertAlterRequestsExt, error) {
	reqs := &alertAlterRequestsExt{set: sdk.NewAlertSetRequest(), unset: sdk.NewAlertUnsetRequest()}
	errs := errors.Join(
		accountObjectIdentifierAttributeUpdate(d, "warehouse", &reqs.set.Warehouse, &reqs.unset.Warehouse),
		stringAttributeUpdateBuilder(d, "comment", reqs.set.WithComment, reqs.unset.WithComment),
	)
	if d.HasChange("alert_schedule") {
		if schedule := alertScheduleFromConfig(d.Get("alert_schedule")); schedule != nil {
			reqs.set.WithSchedule(*schedule)
		}
	}
	if d.HasChange("condition") {
		reqs.modifyCondition = sdk.String(d.Get("condition").(string))
	}
	if d.HasChange("action") {
		reqs.modifyAction = sdk.String(d.Get("action").(string))
	}
	if errs != nil {
		return nil, errs
	}
	return reqs, nil
}

// alertAfterCreateExt resumes the alert after creation when it is enabled. A failure is only a warning.
func alertAfterCreateExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if !d.Get("enabled").(bool) {
		return nil
	}
	client := meta.(*provider.Context).Client
	if err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).WithResume(true)); err != nil {
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

// alertAlterSetInSdkExt runs SET, then MODIFY CONDITION and MODIFY ACTION, which are separate ALTER statements.
func alertAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier, reqs *alertAlterRequestsExt) error {
	client := meta.(*provider.Context).Client
	if !reflect.DeepEqual(*reqs.set, *sdk.NewAlertSetRequest()) {
		if err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).WithSet(*reqs.set)); err != nil {
			return fmt.Errorf("error updating alert %v: %w", id.Name(), err)
		}
	}
	if reqs.modifyCondition != nil {
		if err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).WithModifyCondition([]string{*reqs.modifyCondition})); err != nil {
			return fmt.Errorf("error updating condition of alert %v: %w", id.Name(), err)
		}
	}
	if reqs.modifyAction != nil {
		if err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).WithModifyAction(*reqs.modifyAction)); err != nil {
			return fmt.Errorf("error updating action of alert %v: %w", id.Name(), err)
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
		err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).WithResume(true))
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
		err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).WithSuspend(true))
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

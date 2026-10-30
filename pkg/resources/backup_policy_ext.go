package resources

import (
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	backupPolicyCreateExt      = PreviewFeatureCreateContextWrapper(string(previewfeatures.BackupPolicyResource), TrackingCreateWrapper(resources.BackupPolicy, CreateBackupPolicy))
	backupPolicyReadExt        = PreviewFeatureReadContextWrapper(string(previewfeatures.BackupPolicyResource), TrackingReadWrapper(resources.BackupPolicy, ReadBackupPolicyFunc(true)))
	backupPolicyUpdateExt      = PreviewFeatureUpdateContextWrapper(string(previewfeatures.BackupPolicyResource), TrackingUpdateWrapper(resources.BackupPolicy, UpdateBackupPolicy))
	backupPolicyDescriptionExt = "Resource used to manage backup policies. For more information, check [backup policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-backup-policy)."
	backupPolicyImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.BackupPolicy, ImportBackupPolicy),
	}
)

type backupPolicyAlterRequestsExt struct {
	set   *sdk.BackupPolicySetRequest
	unset *sdk.BackupPolicyUnsetRequest
}

func backupPolicyNewCreateRequestExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier) (*sdk.CreateBackupPolicyRequest, error) {
	req := sdk.NewCreateBackupPolicyRequest(id)
	schedule, err := backupPolicyScheduleFromConfig(d.Get("schedule"))
	if err != nil {
		return nil, err
	}
	if schedule != nil {
		req.WithSchedule(*schedule)
	}
	return req, nil
}

func backupPolicyApplySetUnsetExt(d *schema.ResourceData) (*backupPolicyAlterRequestsExt, error) {
	set, unset := sdk.NewBackupPolicySetRequest(), sdk.NewBackupPolicyUnsetRequest()
	errs := errors.Join(
		intAttributeUpdateBuilder(d, "expire_after_days", set.WithExpireAfterDays, unset.WithExpireAfterDays),
		stringAttributeUpdateBuilder(d, "comment", set.WithComment, unset.WithComment),
	)
	if d.HasChange("schedule") {
		schedule, err := backupPolicyScheduleFromConfig(d.Get("schedule"))
		switch {
		case err != nil:
			errs = errors.Join(errs, err)
		case schedule != nil:
			set.WithSchedule(*schedule)
		default:
			unset.WithSchedule(true)
		}
	}
	if errs != nil {
		return nil, errs
	}
	return &backupPolicyAlterRequestsExt{set: set, unset: unset}, nil
}

func backupPolicySetConfigFieldsExt(d *schema.ResourceData, backupPolicy *sdk.BackupPolicy) error {
	schedule, err := backupPolicyScheduleToState(backupPolicy.Schedule)
	if err != nil {
		return err
	}
	return errors.Join(
		d.Set("comment", backupPolicy.Comment),
		d.Set("schedule", schedule),
	)
}

// backupPolicyScheduleFromConfig turns the schedule block into the SCHEDULE value, or nil when the block is not set.
func backupPolicyScheduleFromConfig(v any) (*string, error) {
	blocks, ok := v.([]any)
	if !ok || len(blocks) == 0 {
		return nil, nil
	}
	block, ok := blocks[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("when setting a schedule one of hours, minutes, or using_cron field should be set")
	}
	if hours, ok := block["hours"].(int); ok && hours > 0 {
		return sdk.String(fmt.Sprintf("%d HOUR", hours)), nil
	}
	if minutes, ok := block["minutes"].(int); ok && minutes > 0 {
		return sdk.String(fmt.Sprintf("%d MINUTE", minutes)), nil
	}
	if cron, ok := block["using_cron"].(string); ok && cron != "" {
		return sdk.String(fmt.Sprintf("USING CRON %s", cron)), nil
	}
	return nil, fmt.Errorf("when setting a schedule one of hours, minutes, or using_cron field should be set")
}

// backupPolicyScheduleToState turns the SCHEDULE value of SHOW into the schedule block, or nil when it is not set.
func backupPolicyScheduleToState(schedule *string) ([]any, error) {
	if schedule == nil || *schedule == "" {
		return nil, nil
	}
	taskSchedule, err := sdk.ParseTaskSchedule(*schedule)
	if err != nil {
		return nil, err
	}
	switch {
	case taskSchedule.Cron != "":
		return []any{map[string]any{"using_cron": taskSchedule.Cron}}, nil
	case taskSchedule.Hours > 0:
		return []any{map[string]any{"hours": taskSchedule.Hours}}, nil
	case taskSchedule.Minutes > 0:
		return []any{map[string]any{"minutes": taskSchedule.Minutes}}, nil
	default:
		return nil, fmt.Errorf("unsupported backup policy schedule: %s", *schedule)
	}
}

func backupPolicySetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier, backupPolicy *sdk.BackupPolicy, importPermissiveMode bool) error {
	schedule, err := backupPolicyScheduleToState(backupPolicy.Schedule)
	if err != nil {
		return err
	}
	if importPermissiveMode {
		schedule = nil
	}
	comment := ""
	if backupPolicy.Comment != nil {
		comment = *backupPolicy.Comment
	}
	expireAfterDays := 0
	if backupPolicy.ExpireAfterDays != nil {
		expireAfterDays = *backupPolicy.ExpireAfterDays
	}
	return errors.Join(
		d.Set("name", id.Name()),
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
		d.Set("schedule", schedule),
		importStringAttribute(d, "comment", importPermissiveMode, comment),
		importIntAttribute(d, "expire_after_days", importPermissiveMode, expireAfterDays),
		importBooleanStringAttribute(d, "with_retention_lock", importPermissiveMode, backupPolicy.HasRetentionLock),
	)
}

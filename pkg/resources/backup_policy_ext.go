package resources

import (
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	backupPolicyCreateExt      = TrackingCreateWrapper(resources.BackupPolicy, CreateBackupPolicy)
	backupPolicyReadExt        = TrackingReadWrapper(resources.BackupPolicy, ReadBackupPolicyFunc(true))
	backupPolicyUpdateExt      = TrackingUpdateWrapper(resources.BackupPolicy, UpdateBackupPolicy)
	backupPolicyDescriptionExt = "Resource used to manage backup policies. For more information, check [backup policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-backup-policy)."
	backupPolicyImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.BackupPolicy, ImportBackupPolicy),
	}
)

type backupPolicyAlterRequestsExt struct {
	set   *sdk.BackupPolicySetRequest
	unset *sdk.BackupPolicyUnsetRequest
}

func backupPolicySetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier, backupPolicy *sdk.BackupPolicy, importPermissiveMode bool) error {
	schedule := ""
	if backupPolicy.Schedule != nil {
		schedule = *backupPolicy.Schedule
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
		importStringAttribute(d, "schedule", importPermissiveMode, schedule),
		importStringAttribute(d, "comment", importPermissiveMode, comment),
		importIntAttribute(d, "expire_after_days", importPermissiveMode, expireAfterDays),
		importBooleanStringAttribute(d, "with_retention_lock", importPermissiveMode, backupPolicy.HasRetentionLock),
	)
}

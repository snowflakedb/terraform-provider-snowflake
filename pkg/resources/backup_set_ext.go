package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	backupSetCreateExt      = TrackingCreateWrapper(resources.BackupSet, CreateBackupSet)
	backupSetReadExt        = TrackingReadWrapper(resources.BackupSet, ReadBackupSetFunc(true))
	backupSetUpdateExt      = TrackingUpdateWrapper(resources.BackupSet, UpdateBackupSet)
	backupSetDescriptionExt = "Resource used to manage backup sets. For more information, check [backup set documentation](https://docs.snowflake.com/en/user-guide/backups)."
	backupSetImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.BackupSet, importBackupSet),
	}
)

func init() {
	backupSetSchema["backup_policy"].DiffSuppressFunc = SuppressIfAny(
		backupSetSchema["backup_policy"].DiffSuppressFunc,
		suppressBackupSetBackupPolicy,
	)
}

type backupSetAlterRequestsExt struct {
	set   *sdk.BackupSetSetRequest
	unset *sdk.BackupSetUnsetRequest
}

// SHOW splits the target and the backup policy across columns, so those attributes stay
// copied from config. This marks show and describe output computed when they change.
func backupSetOutputsComputedWhenTargetOrPolicyChanges(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	return customdiff.All(
		ComputedIfAnyAttributeChanged(backupSetSchema, ShowOutputAttributeName, "backup_policy", "for_table", "for_dynamic_table", "for_schema", "for_database"),
		ComputedIfAnyAttributeChanged(backupSetSchema, DescribeOutputAttributeName, "backup_policy", "for_table", "for_dynamic_table", "for_schema", "for_database"),
	)(ctx, d, meta)
}

func backupSetApplyCreateOptionalsExt(d *schema.ResourceData, request *sdk.CreateBackupSetRequest) error {
	return errors.Join(
		schemaObjectIdentifierAttributeCreate(d, "for_table", &request.ForTable),
		schemaObjectIdentifierAttributeCreate(d, "for_dynamic_table", &request.ForDynamicTable),
		databaseObjectIdentifierAttributeCreate(d, "for_schema", &request.ForSchema),
		accountObjectIdentifierAttributeCreateBuilder(d, "for_database", request.WithForDatabase),
		schemaObjectIdentifierAttributeCreate(d, "backup_policy", &request.WithBackupPolicy),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
}

// Name is not ForceNew. RENAME TO updates the id used by the later SET and UNSET.
func backupSetHierarchyRenameExt(ctx context.Context, d *schema.ResourceData, meta any, id *sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if !d.HasChange("name") || d.GetRawState().IsNull() {
		return nil
	}
	newID := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), d.Get("name").(string))
	if err := meta.(*provider.Context).Client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(*id).WithRenameTo(newID)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(newID))
	*id = newID
	return nil
}

// backup_policy is APPLY BACKUP POLICY, not SET. The public page has no UNSET BACKUP POLICY.
// FORCE is used only when a policy is already attached, and it fails when that policy has a retention lock.
func backupSetBeforeAlterExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if !d.HasChange("backup_policy") {
		return nil
	}
	value := d.Get("backup_policy").(string)
	if value == "" {
		return diag.Errorf("An attached backup policy cannot be removed. Change backup_policy to another policy, or recreate the backup set when the current policy has no retention lock.")
	}
	policyID, err := sdk.ParseSchemaObjectIdentifier(value)
	if err != nil {
		return diag.FromErr(err)
	}
	request := sdk.NewBackupSetApplyBackupPolicyRequest(policyID)
	if old, _ := d.GetChange("backup_policy"); old.(string) != "" {
		request = request.WithForce(true)
	}
	if err := meta.(*provider.Context).Client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).WithApplyBackupPolicy(*request)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func importBackupSet(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	id, err := backupSetParseId(d.Id())
	if err != nil {
		return nil, err
	}
	client := meta.(*provider.Context).Client
	backupSet, err := client.BackupSets.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}
	importPermissiveMode := meta.(*provider.Context).Experiments.IsEnabled(experimentalfeatures.ImportPermissiveMode)
	if err := backupSetSetImportedFields(ctx, client, d, id, backupSet, importPermissiveMode); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func backupSetSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier, backupSet *sdk.BackupSet, importPermissiveMode bool) error {
	return backupSetSetImportedFields(context.Background(), nil, d, id, backupSet, importPermissiveMode)
}

func backupSetSetImportedFields(ctx context.Context, client *sdk.Client, d *schema.ResourceData, id sdk.SchemaObjectIdentifier, backupSet *sdk.BackupSet, importPermissiveMode bool) error {
	errs := []error{
		d.Set("name", id.Name()),
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
		// Read sets comment from SHOW in both modes.
		d.Set("comment", backupSet.Comment),
		importOptionalStringAttribute(d, "backup_policy", importPermissiveMode, backupSetBackupPolicy(backupSet)),
	}
	switch backupSet.ObjectKind {
	case "TABLE":
		if backupSet.ObjectDatabaseName == nil || *backupSet.ObjectDatabaseName == "" || backupSet.ObjectSchemaName == nil || *backupSet.ObjectSchemaName == "" || backupSet.ObjectName == "" {
			return fmt.Errorf("backup set %s is missing the database, schema, or name of its table target", id.FullyQualifiedName())
		}
		if client == nil {
			return fmt.Errorf("backup set %s import needs the table target kind", id.FullyQualifiedName())
		}
		targetID := sdk.NewSchemaObjectIdentifier(*backupSet.ObjectDatabaseName, *backupSet.ObjectSchemaName, backupSet.ObjectName)
		field, err := backupSetTableTargetField(ctx, client, targetID)
		if err != nil {
			return err
		}
		errs = append(errs, d.Set(field, targetID.FullyQualifiedName()))
	case "SCHEMA":
		if backupSet.ObjectDatabaseName == nil || *backupSet.ObjectDatabaseName == "" || backupSet.ObjectName == "" {
			return fmt.Errorf("backup set %s is missing the database or name of its schema target", id.FullyQualifiedName())
		}
		errs = append(errs, d.Set("for_schema", sdk.NewDatabaseObjectIdentifier(*backupSet.ObjectDatabaseName, backupSet.ObjectName).FullyQualifiedName()))
	case "DATABASE":
		if backupSet.ObjectName == "" {
			return fmt.Errorf("backup set %s is missing the name of its database target", id.FullyQualifiedName())
		}
		errs = append(errs, d.Set("for_database", sdk.NewAccountObjectIdentifier(backupSet.ObjectName).FullyQualifiedName()))
	default:
		return fmt.Errorf("backup set %s has unsupported object kind %q", id.FullyQualifiedName(), backupSet.ObjectKind)
	}
	return errors.Join(errs...)
}

// SHOW TABLES reports kind TABLE for a dynamic table, so is_dynamic decides the field.
func backupSetTableTargetField(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) (string, error) {
	table, err := client.Tables.ShowByID(ctx, id)
	if err != nil {
		return "", err
	}
	if table.IsDynamic {
		return "for_dynamic_table", nil
	}
	return "for_table", nil
}

// SHOW splits the policy across three columns. A missing part means no policy is attached.
func backupSetBackupPolicy(backupSet *sdk.BackupSet) *string {
	if backupSet.BackupPolicyDatabaseName == nil || *backupSet.BackupPolicyDatabaseName == "" ||
		backupSet.BackupPolicySchemaName == nil || *backupSet.BackupPolicySchemaName == "" ||
		backupSet.BackupPolicyName == nil || *backupSet.BackupPolicyName == "" {
		return nil
	}
	value := sdk.NewSchemaObjectIdentifier(*backupSet.BackupPolicyDatabaseName, *backupSet.BackupPolicySchemaName, *backupSet.BackupPolicyName).FullyQualifiedName()
	return &value
}

// suppressBackupSetBackupPolicy ignores a config value equal to the policy SHOW reports.
// The attribute is not one show_output column, so the generated suppressor does not apply.
func suppressBackupSetBackupPolicy(_, _, newValue string, d *schema.ResourceData) bool {
	if d.Id() == "" || newValue == "" {
		return false
	}
	parsed, err := sdk.ParseSchemaObjectIdentifier(newValue)
	if err != nil {
		return false
	}
	for _, output := range []string{ShowOutputAttributeName, DescribeOutputAttributeName} {
		policy, ok := backupSetPolicyIdentifierFromOutput(d, output)
		if ok && policy.FullyQualifiedName() == parsed.FullyQualifiedName() {
			return true
		}
	}
	return false
}

func backupSetPolicyIdentifierFromOutput(d *schema.ResourceData, attr string) (sdk.SchemaObjectIdentifier, bool) {
	output, ok := d.GetOk(attr)
	if !ok {
		return sdk.SchemaObjectIdentifier{}, false
	}
	list, ok := output.([]any)
	if !ok || len(list) != 1 {
		return sdk.SchemaObjectIdentifier{}, false
	}
	row, ok := list[0].(map[string]any)
	if !ok {
		return sdk.SchemaObjectIdentifier{}, false
	}
	database, _ := row["backup_policy_database_name"].(string)
	schemaName, _ := row["backup_policy_schema_name"].(string)
	name, _ := row["backup_policy_name"].(string)
	if database == "" || schemaName == "" || name == "" {
		return sdk.SchemaObjectIdentifier{}, false
	}
	return sdk.NewSchemaObjectIdentifier(database, schemaName, name), true
}

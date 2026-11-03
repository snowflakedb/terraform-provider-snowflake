//go:build non_account_level_tests

package testacc

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_BackupSet_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	table, cleanupTable := testClient().Table.CreateWithPredefinedColumns(t)
	t.Cleanup(cleanupTable)
	tableID := table.ID()

	basic := model.BackupSet("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithForTable(tableID.FullyQualifiedName())
	ref := basic.ResourceReference()
	providerModelImportPermissiveMode := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportPermissiveMode)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.BackupSet),
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, backupSetAssertions(t, ref, id, tableID, "", nil)...),
			},
			// import - permissive mode
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   accconfig.FromModels(t, providerModelImportPermissiveMode, basic),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
		},
	})
}

func TestAcc_BackupSet_BasicUseCase_Import(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	table, cleanupTable := testClient().Table.CreateWithPredefinedColumns(t)
	t.Cleanup(cleanupTable)
	tableID := table.ID()

	// The basic config omits comment, so enforcing import plans an unset.
	backupSet, cleanupBackupSet := testClient().BackupSet.CreateWithOptions(t, id, sdk.NewCreateBackupSetRequest(id).
		WithForTable(tableID).
		WithComment(comment))
	t.Cleanup(cleanupBackupSet)
	require.NotNil(t, backupSet.Comment)

	basic := model.BackupSet("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithForTable(tableID.FullyQualifiedName())
	ref := basic.ResourceReference()
	resourceId := helpers.EncodeResourceIdentifier(id)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.BackupSet),
		Steps: []resource.TestStep{
			{
				Config:             accconfig.FromModels(t, basic),
				ResourceName:       ref,
				ImportState:        true,
				ImportStateId:      resourceId,
				ImportStatePersist: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedBackupSetResource(t, resourceId).
						HasName(id.Name()).
						HasDatabase(id.DatabaseName()).
						HasSchema(id.SchemaName()).
						HasForTable(tableID.FullyQualifiedName()).
						HasNoForDynamicTable().
						HasNoForSchema().
						HasNoForDatabase().
						HasComment(*backupSet.Comment).
						HasNoBackupPolicy().
						HasFullyQualifiedName(id.FullyQualifiedName()),
				),
			},
			{
				Config: accconfig.FromModels(t, basic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(*backupSet.Comment), nil),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t,
					resourceassert.BackupSetResource(t, ref).
						HasName(id.Name()).
						HasDatabase(id.DatabaseName()).
						HasSchema(id.SchemaName()).
						HasForTable(tableID.FullyQualifiedName()).
						HasNoForDynamicTable().
						HasNoForSchema().
						HasNoForDatabase().
						HasCommentEmpty().
						HasNoBackupPolicy().
						HasFullyQualifiedName(id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_BackupSet_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	table, cleanupTable := testClient().Table.CreateWithPredefinedColumns(t)
	t.Cleanup(cleanupTable)
	tableID := table.ID()

	policyID := testClient().Ids.RandomSchemaObjectIdentifier()
	t.Cleanup(testClient().BackupPolicy.Create(t, policyID))

	// for_table is ForceNew. comment is SET/UNSET. backup_policy is APPLY and cannot be unset.
	complete := model.BackupSet("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithForTable(tableID.FullyQualifiedName()).
		WithComment(comment).
		WithBackupPolicy(policyID.FullyQualifiedName())
	commentUnset := model.BackupSet("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithForTable(tableID.FullyQualifiedName()).
		WithBackupPolicy(policyID.FullyQualifiedName())
	ref := complete.ResourceReference()
	providerModelImportPermissiveMode := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportPermissiveMode)

	assertComplete := backupSetAssertions(t, ref, id, tableID, comment, &policyID)
	assertCommentUnset := backupSetAssertions(t, ref, id, tableID, "", &policyID)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.BackupSet),
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, complete),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateVerify:        true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedBackupSetResource(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasDatabase(id.DatabaseName()).
						HasSchema(id.SchemaName()).
						HasForTable(tableID.FullyQualifiedName()).
						HasNoForDynamicTable().
						HasNoForSchema().
						HasNoForDatabase().
						HasComment(comment).
						HasBackupPolicy(policyID.FullyQualifiedName()).
						HasFullyQualifiedName(id.FullyQualifiedName()),
					resourceshowoutputassert.ImportedBackupSetShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasComment(comment).
						HasBackupPolicyName(policyID.Name()).
						HasBackupPolicyState(sdk.BackupPolicyStateActive),
				),
			},
			// import - permissive mode
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   accconfig.FromModels(t, providerModelImportPermissiveMode, complete),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
						planchecks.ExpectNoChangeOnField(ref, "backup_policy"),
						planchecks.ExpectNoChangeOnField(ref, "for_table"),
					},
				},
				Config: accconfig.FromModels(t, commentUnset),
				Check:  assertThat(t, assertCommentUnset...),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, commentUnset),
				Destroy:                  true,
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					_, err := testClient().BackupSet.Show(t, id)
					require.ErrorIs(t, err, sdk.ErrObjectNotFound)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, complete),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateVerify:        true,
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					testClient().BackupSet.Alter(t, sdk.NewAlterBackupSetRequest(id).
						WithSet(*sdk.NewBackupSetSetRequest().WithComment(externalComment)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
						planchecks.ExpectNoChangeOnField(ref, "backup_policy"),
						planchecks.ExpectNoChangeOnField(ref, "for_table"),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_BackupSet_DynamicTable(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	_, dynamicID := backupSetTableAndDynamicTable(t)

	dynamic := model.BackupSet("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithForDynamicTable(dynamicID.FullyQualifiedName())
	ref := dynamic.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.BackupSet),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: accconfig.FromModels(t, dynamic),
				Check: assertThat(t,
					resourceassert.BackupSetResource(t, ref).
						HasForDynamicTable(dynamicID.FullyQualifiedName()).
						HasNoForTable().
						HasNoForSchema().
						HasNoForDatabase().
						HasFullyQualifiedName(id.FullyQualifiedName()),
					resourceshowoutputassert.BackupSetShowOutput(t, ref).
						HasObjectKind("TABLE").
						HasObjectName(dynamicID.Name()).
						HasObjectDatabaseName(dynamicID.DatabaseName()).
						HasObjectSchemaName(dynamicID.SchemaName()),
					objectassert.BackupSet(t, id).
						HasObjectKind("TABLE").
						HasObjectName(dynamicID.Name()).
						HasNoBackupPolicyState(),
				),
			},
			// Import uses is_dynamic to choose for_dynamic_table. SHOW TABLES reports kind TABLE for both.
			{
				Config:            accconfig.FromModels(t, dynamic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedBackupSetResource(t, helpers.EncodeResourceIdentifier(id)).
						HasForDynamicTable(dynamicID.FullyQualifiedName()).
						HasNoForTable().
						HasNoForSchema().
						HasNoForDatabase(),
				),
			},
		},
	})
}

func TestAcc_BackupSet_ExternalTargetTypeChange(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	tableID, dynamicID := backupSetTableAndDynamicTable(t)

	tableModel := model.BackupSet("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithForTable(tableID.FullyQualifiedName())
	ref := tableModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.BackupSet),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, tableModel),
				Check: assertThat(t,
					resourceassert.BackupSetResource(t, ref).HasForTable(tableID.FullyQualifiedName()),
					objectassert.BackupSet(t, id).HasObjectName(tableID.Name()).HasObjectKind("TABLE"),
				),
			},
			// The target cannot be altered. Replacing the backup set outside Terraform changes its type.
			{
				PreConfig: func() {
					testClient().BackupSet.DropFunc(t, id)()
					_, cleanup := testClient().BackupSet.CreateWithOptions(t, id, sdk.NewCreateBackupSetRequest(id).WithForDynamicTable(dynamicID))
					t.Cleanup(cleanup)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						// Read writes the remote target (the dynamic table) onto for_dynamic_table.
						// for_table is left untouched at the originally configured table, but
						// for_dynamic_table going from set to absent is enough to force recreation.
						planchecks.ExpectChange(ref, "for_dynamic_table", tfjson.ActionDelete, sdk.String(dynamicID.FullyQualifiedName()), nil),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: accconfig.FromModels(t, tableModel),
				Check: assertThat(t,
					resourceassert.BackupSetResource(t, ref).
						HasForTable(tableID.FullyQualifiedName()).
						HasNoForDynamicTable(),
					objectassert.BackupSet(t, id).HasObjectName(tableID.Name()).HasObjectKind("TABLE"),
				),
			},
		},
	})
}

func TestAcc_BackupSet_TargetTypeChange(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	tableID, dynamicID := backupSetTableAndDynamicTable(t)

	tableModel := model.BackupSet("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithForTable(tableID.FullyQualifiedName())
	dynamicModel := model.BackupSet("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithForDynamicTable(dynamicID.FullyQualifiedName())
	ref := tableModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.BackupSet),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, tableModel),
				Check: assertThat(t,
					resourceassert.BackupSetResource(t, ref).HasForTable(tableID.FullyQualifiedName()).HasNoForDynamicTable(),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(ref, "for_table", tfjson.ActionDelete, sdk.String(tableID.FullyQualifiedName()), nil),
						planchecks.ExpectChange(ref, "for_dynamic_table", tfjson.ActionCreate, nil, sdk.String(dynamicID.FullyQualifiedName())),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: accconfig.FromModels(t, dynamicModel),
				Check: assertThat(t,
					resourceassert.BackupSetResource(t, ref).
						HasForDynamicTable(dynamicID.FullyQualifiedName()).
						HasNoForTable(),
					resourceshowoutputassert.BackupSetShowOutput(t, ref).
						HasObjectKind("TABLE").
						HasObjectName(dynamicID.Name()),
					objectassert.BackupSet(t, id).HasObjectName(dynamicID.Name()).HasObjectKind("TABLE"),
				),
			},
		},
	})
}

func backupSetTableAndDynamicTable(t *testing.T) (sdk.SchemaObjectIdentifier, sdk.SchemaObjectIdentifier) {
	t.Helper()
	table, cleanupTable := testClient().Table.CreateWithPredefinedColumns(t)
	t.Cleanup(cleanupTable)
	dynamicTable, cleanupDynamicTable := testClient().DynamicTable.CreateDynamicTable(t, table.ID())
	t.Cleanup(cleanupDynamicTable)
	return table.ID(), dynamicTable.ID()
}

func backupSetAssertions(t *testing.T, ref string, id sdk.SchemaObjectIdentifier, tableID sdk.SchemaObjectIdentifier, comment string, policyID *sdk.SchemaObjectIdentifier) []assert.TestCheckFuncProvider {
	t.Helper()
	resourceCheck := resourceassert.BackupSetResource(t, ref).
		HasNameString(id.Name()).
		HasDatabaseString(id.DatabaseName()).
		HasSchemaString(id.SchemaName()).
		HasForTableString(tableID.FullyQualifiedName()).
		HasFullyQualifiedNameString(id.FullyQualifiedName())
	showCheck := resourceshowoutputassert.BackupSetShowOutput(t, ref).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasObjectKind("TABLE").
		HasObjectName(tableID.Name()).
		HasObjectDatabaseName(tableID.DatabaseName()).
		HasObjectSchemaName(tableID.SchemaName()).
		HasOwnerRole(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE")
	describeCheck := resourceshowoutputassert.BackupSetDescribeOutput(t, ref).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasObjectKind("TABLE").
		HasObjectName(tableID.Name()).
		HasObjectDatabaseName(tableID.DatabaseName()).
		HasObjectSchemaName(tableID.SchemaName()).
		HasOwnerRole(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE")
	showObject := objectassert.BackupSet(t, id).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasObjectKind("TABLE").
		HasObjectName(tableID.Name()).
		HasObjectDatabaseName(tableID.DatabaseName()).
		HasObjectSchemaName(tableID.SchemaName()).
		HasOwnerRole(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE")
	details := objectassert.BackupSetDetails(t, id).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasObjectKind("TABLE").
		HasObjectName(tableID.Name()).
		HasObjectDatabaseName(tableID.DatabaseName()).
		HasObjectSchemaName(tableID.SchemaName()).
		HasOwnerRole(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE")
	if comment == "" {
		resourceCheck = resourceCheck.HasCommentEmpty()
		showCheck = showCheck.HasComment("")
		describeCheck = describeCheck.HasComment("")
		// SHOW and DESCRIBE return an empty string, not NULL, when comment is unset.
		showObject = showObject.HasComment("")
		details = details.HasComment("")
	} else {
		resourceCheck = resourceCheck.HasCommentString(comment)
		showCheck = showCheck.HasComment(comment)
		describeCheck = describeCheck.HasComment(comment)
		showObject = showObject.HasComment(comment)
		details = details.HasComment(comment)
	}
	if policyID == nil {
		resourceCheck = resourceCheck.HasNoBackupPolicy()
		showCheck = showCheck.
			HasBackupPolicyName("").
			HasBackupPolicyDatabaseName("").
			HasBackupPolicySchemaName("")
		showCheck.StringValueSet("backup_policy_state", "")
		describeCheck = describeCheck.
			HasBackupPolicyName("").
			HasBackupPolicyDatabaseName("").
			HasBackupPolicySchemaName("")
		describeCheck.StringValueSet("backup_policy_state", "")
		// SHOW and DESCRIBE return empty strings, not NULL, when no policy is attached.
		showObject = showObject.
			HasBackupPolicyName("").
			HasBackupPolicyDatabaseName("").
			HasBackupPolicySchemaName("").
			HasNoBackupPolicyState()
		details = details.
			HasBackupPolicyName("").
			HasBackupPolicyDatabaseName("").
			HasBackupPolicySchemaName("").
			HasNoBackupPolicyState()
	} else {
		resourceCheck = resourceCheck.HasBackupPolicyString(policyID.FullyQualifiedName())
		showCheck = showCheck.
			HasBackupPolicyName(policyID.Name()).
			HasBackupPolicyDatabaseName(policyID.DatabaseName()).
			HasBackupPolicySchemaName(policyID.SchemaName()).
			HasBackupPolicyState(sdk.BackupPolicyStateActive)
		describeCheck = describeCheck.
			HasBackupPolicyName(policyID.Name()).
			HasBackupPolicyDatabaseName(policyID.DatabaseName()).
			HasBackupPolicySchemaName(policyID.SchemaName()).
			HasBackupPolicyState(sdk.BackupPolicyStateActive)
		showObject = showObject.
			HasBackupPolicyName(policyID.Name()).
			HasBackupPolicyDatabaseName(policyID.DatabaseName()).
			HasBackupPolicySchemaName(policyID.SchemaName()).
			HasBackupPolicyState(sdk.BackupPolicyStateActive)
		details = details.
			HasBackupPolicyName(policyID.Name()).
			HasBackupPolicyDatabaseName(policyID.DatabaseName()).
			HasBackupPolicySchemaName(policyID.SchemaName()).
			HasBackupPolicyState(sdk.BackupPolicyStateActive)
	}
	return []assert.TestCheckFuncProvider{
		resourceCheck,
		showCheck,
		describeCheck,
		showObject,
		details,
	}
}

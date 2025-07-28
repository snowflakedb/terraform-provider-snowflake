//go:build account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO(SNOW-2197902): Extract common part from current account tests (they are reused here)

func TestAcc_CurrentOrganizationAccount_Parameters(t *testing.T) {
	testClient().EnsureValidNonProdOrganizationAccountIsUsed(t)

	warehouseId := testClient().Ids.WarehouseId()

	eventTable, eventTableCleanup := testClient().EventTable.Create(t)
	t.Cleanup(eventTableCleanup)

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	createNetworkPolicyRequest := sdk.NewCreateNetworkPolicyRequest(testClient().Ids.RandomAccountObjectIdentifier()).WithAllowedIpList([]sdk.IPRequest{*sdk.NewIPRequest("0.0.0.0/0")})
	networkPolicy, networkPolicyCleanup := testClient().NetworkPolicy.CreateNetworkPolicyWithRequest(t, createNetworkPolicyRequest)
	t.Cleanup(networkPolicyCleanup)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	provider := providermodel.SnowflakeProvider().WithWarehouse(testClient().Ids.WarehouseId().FullyQualifiedName())
	currentOrganizationAccountName := testClient().OrganizationAccount.ShowCurrent(t).AccountName
	unsetParametersModel := model.CurrentOrganizationAccount("test", currentOrganizationAccountName)
	setParametersModel := model.CurrentOrganizationAccount("test", currentOrganizationAccountName).WithAllParametersSetToPredefinedValues(warehouseId, eventTable.ID(), externalVolumeId, networkPolicy.ID(), stage.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// resource with unset parameters
			{
				Config: config.FromModels(t, provider, unsetParametersModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, unsetParametersModel.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasAllDefaultParameters(),
				),
			},
			// import with unset parameters
			{
				Config:       config.FromModels(t, provider, unsetParametersModel),
				ResourceName: unsetParametersModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedCurrentOrganizationAccountResource(t, currentOrganizationAccountName).
						HasNameString(currentOrganizationAccountName).
						HasAllDefaultParameters(),
				),
			},
			// set all parameters
			{
				Config: config.FromModels(t, provider, setParametersModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, setParametersModel.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasAllParametersEqualToPredefinedValues(warehouseId, eventTable.ID(), externalVolumeId, networkPolicy.ID(), stage.ID()),
				),
			},
			// import with all parameters set
			{
				Config:       config.FromModels(t, provider, setParametersModel),
				ResourceName: setParametersModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedCurrentOrganizationAccountResource(t, currentOrganizationAccountName).
						HasNameString(currentOrganizationAccountName).
						HasAllParametersEqualToPredefinedValues(warehouseId, eventTable.ID(), externalVolumeId, networkPolicy.ID(), stage.ID()),
				),
			},
			// unset parameters
			{
				Config: config.FromModels(t, provider, unsetParametersModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, unsetParametersModel.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasAllDefaultParameters(),
				),
			},
			// Test for external changes
			{
				PreConfig: func() {
					testClient().Account.Alter(t, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{Parameters: &sdk.AccountParameters{AbortDetachedQuery: sdk.Bool(true)}}})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(setParametersModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, provider, unsetParametersModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, setParametersModel.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasAllDefaultParameters(),
				),
			},
		},
	})
}

func TestAcc_CurrentOrganizationAccount_NonParameterValues(t *testing.T) {
	testClient().EnsureValidNonProdOrganizationAccountIsUsed(t)

	resourceMonitor, resourceMonitorCleanup := testClient().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(resourceMonitorCleanup)

	newResourceMonitor, newResourceMonitorCleanup := testClient().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(newResourceMonitorCleanup)

	passwordPolicy, passwordPolicyCleanup := testClient().PasswordPolicy.CreatePasswordPolicy(t)
	t.Cleanup(passwordPolicyCleanup)

	newPasswordPolicy, newPasswordPolicyCleanup := testClient().PasswordPolicy.CreatePasswordPolicy(t)
	t.Cleanup(newPasswordPolicyCleanup)

	sessionPolicy, sessionPolicyCleanup := testClient().SessionPolicy.CreateSessionPolicy(t)
	t.Cleanup(sessionPolicyCleanup)

	newSessionPolicy, newSessionPolicyCleanup := testClient().SessionPolicy.CreateSessionPolicy(t)
	t.Cleanup(newSessionPolicyCleanup)

	comment := random.Comment()
	newComment := random.Comment()

	provider := providermodel.SnowflakeProvider().WithWarehouse(testClient().Ids.WarehouseId().FullyQualifiedName())

	currentOrganizationAccountName := testClient().OrganizationAccount.ShowCurrent(t).AccountName

	unsetModel := model.CurrentOrganizationAccount("test", currentOrganizationAccountName)

	setModel := model.CurrentOrganizationAccount("test", currentOrganizationAccountName).
		WithComment(comment).
		WithResourceMonitor(resourceMonitor.ID().Name()).
		WithPasswordPolicy(passwordPolicy.ID().FullyQualifiedName()).
		WithSessionPolicy(sessionPolicy.ID().FullyQualifiedName())

	setModelWithDifferentValues := model.CurrentOrganizationAccount("test", currentOrganizationAccountName).
		WithComment(newComment).
		WithResourceMonitor(newResourceMonitor.ID().Name()).
		WithPasswordPolicy(newPasswordPolicy.ID().FullyQualifiedName()).
		WithSessionPolicy(newSessionPolicy.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// create with unset values
			{
				Config: config.FromModels(t, provider, unsetModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, unsetModel.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasCommentEmpty().
						HasNoResourceMonitor().
						HasPasswordPolicyEmpty().
						HasSessionPolicyEmpty(),
				),
			},
			// import
			{
				Config:       config.FromModels(t, provider, unsetModel),
				ResourceName: unsetModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedCurrentOrganizationAccountResource(t, currentOrganizationAccountName).
						HasNameString(currentOrganizationAccountName).
						HasCommentEmpty().
						HasNoResourceMonitor().
						HasPasswordPolicyEmpty().
						HasSessionPolicyEmpty(),
				),
			},
			// set optional values externally
			{
				PreConfig: func() {
					testClient().OrganizationAccount.Alter(t, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithSessionPolicy(sessionPolicy.ID())))
					testClient().OrganizationAccount.Alter(t, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithPasswordPolicy(passwordPolicy.ID())))
					testClient().OrganizationAccount.Alter(t, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithComment(comment)))
				},
				Config: config.FromModels(t, provider, unsetModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, unsetModel.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasCommentEmpty().
						HasNoResourceMonitor().
						HasPasswordPolicyEmpty().
						HasSessionPolicyEmpty(),
				),
			},
			// set optional values
			{
				Config: config.FromModels(t, provider, setModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, setModel.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasCommentString(comment).
						HasResourceMonitorString(resourceMonitor.ID().Name()).
						HasPasswordPolicyString(passwordPolicy.ID().FullyQualifiedName()).
						HasSessionPolicyString(sessionPolicy.ID().FullyQualifiedName()),
				),
			},
			// import
			{
				Config:       config.FromModels(t, provider, setModel),
				ResourceName: setModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedCurrentOrganizationAccountResource(t, currentOrganizationAccountName).
						HasNameString(currentOrganizationAccountName).
						HasCommentString(comment).
						HasNoResourceMonitor().
						HasPasswordPolicyString(passwordPolicy.ID().FullyQualifiedName()).
						HasSessionPolicyString(sessionPolicy.ID().FullyQualifiedName()),
				),
			},
			// set new optional values
			{
				Config: config.FromModels(t, provider, setModelWithDifferentValues),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, setModelWithDifferentValues.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasCommentString(newComment).
						HasResourceMonitorString(newResourceMonitor.ID().Name()).
						HasPasswordPolicyString(newPasswordPolicy.ID().FullyQualifiedName()).
						HasSessionPolicyString(newSessionPolicy.ID().FullyQualifiedName()),
				),
			},
			// unset externally
			{
				PreConfig: func() {
					testClient().OrganizationAccount.Alter(t, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithSessionPolicy(true)))
					testClient().OrganizationAccount.Alter(t, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithPasswordPolicy(true)))
					testClient().OrganizationAccount.Alter(t, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithComment(true)))
				},
				Config: config.FromModels(t, provider, setModelWithDifferentValues),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, setModelWithDifferentValues.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasCommentString(newComment).
						HasResourceMonitorString(newResourceMonitor.ID().Name()).
						HasPasswordPolicyString(newPasswordPolicy.ID().FullyQualifiedName()).
						HasSessionPolicyString(newSessionPolicy.ID().FullyQualifiedName()),
				),
			},
			// unset optional values
			{
				Config: config.FromModels(t, provider, unsetModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, unsetModel.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasCommentEmpty().
						HasResourceMonitorEmpty().
						HasPasswordPolicyEmpty().
						HasSessionPolicyEmpty(),
				),
			},
		},
	})
}

func TestAcc_CurrentOrganizationAccount_Complete(t *testing.T) {
	testClient().EnsureValidNonProdOrganizationAccountIsUsed(t)

	warehouseId := testClient().Ids.WarehouseId()

	eventTable, eventTableCleanup := testClient().EventTable.Create(t)
	t.Cleanup(eventTableCleanup)

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	createNetworkPolicyRequest := sdk.NewCreateNetworkPolicyRequest(testClient().Ids.RandomAccountObjectIdentifier()).WithAllowedIpList([]sdk.IPRequest{*sdk.NewIPRequest("0.0.0.0/0")})
	networkPolicy, networkPolicyCleanup := testClient().NetworkPolicy.CreateNetworkPolicyWithRequest(t, createNetworkPolicyRequest)
	t.Cleanup(networkPolicyCleanup)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	resourceMonitor, resourceMonitorCleanup := testClient().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(resourceMonitorCleanup)

	passwordPolicy, passwordPolicyCleanup := testClient().PasswordPolicy.CreatePasswordPolicy(t)
	t.Cleanup(passwordPolicyCleanup)

	sessionPolicy, sessionPolicyCleanup := testClient().SessionPolicy.CreateSessionPolicy(t)
	t.Cleanup(sessionPolicyCleanup)

	provider := providermodel.SnowflakeProvider().WithWarehouse(testClient().Ids.WarehouseId().FullyQualifiedName())

	currentOrganizationAccountName := testClient().OrganizationAccount.ShowCurrent(t).AccountName
	completeConfigModel := model.CurrentOrganizationAccount("test", currentOrganizationAccountName).
		WithAllParametersSetToPredefinedValues(warehouseId, eventTable.ID(), externalVolumeId, networkPolicy.ID(), stage.ID()).
		WithResourceMonitor(resourceMonitor.ID().Name()).
		WithPasswordPolicy(passwordPolicy.ID().FullyQualifiedName()).
		WithSessionPolicy(sessionPolicy.ID().FullyQualifiedName())

	config.FromModels(t, completeConfigModel)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, provider, completeConfigModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, completeConfigModel.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasAllDefaultParameters().
						HasCommentEmpty().
						HasResourceMonitorString(resourceMonitor.ID().Name()).
						HasPasswordPolicyEmpty().
						HasSessionPolicyEmpty(),
				),
				// If the config is set for all fields, we expect to see them as changes in the plan (because they're different from the default values).
				// They will be applied in the next step.
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(completeConfigModel.ResourceReference(), plancheck.ResourceActionUpdate),
						// Expect changes for all fields that are set in the config: parameters (chose to show only 3 parameters of a different type as an example), resource monitor, password policy, session policy, comment
						planchecks.ExpectChange(completeConfigModel.ResourceReference(), "abort_detached_query", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange(completeConfigModel.ResourceReference(), "binary_input_format", tfjson.ActionUpdate, sdk.String("HEX"), sdk.String("BASE64")),
						planchecks.ExpectChange(completeConfigModel.ResourceReference(), "client_encryption_key_size", tfjson.ActionUpdate, sdk.String("128"), sdk.String("256")),
						planchecks.ExpectChange(completeConfigModel.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String(""), sdk.String("")),
						planchecks.ExpectChange(completeConfigModel.ResourceReference(), "resource_monitor", tfjson.ActionUpdate, sdk.String(resourceMonitor.ID().Name()), sdk.String(resourceMonitor.ID().Name())),
						planchecks.ExpectChange(completeConfigModel.ResourceReference(), "password_policy", tfjson.ActionUpdate, sdk.String(""), sdk.String(passwordPolicy.ID().FullyQualifiedName())),
						planchecks.ExpectChange(completeConfigModel.ResourceReference(), "session_policy", tfjson.ActionUpdate, sdk.String(""), sdk.String(sessionPolicy.ID().FullyQualifiedName())),
					},
				},
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(completeConfigModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, provider, completeConfigModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, completeConfigModel.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasAllParametersEqualToPredefinedValues(warehouseId, eventTable.ID(), externalVolumeId, networkPolicy.ID(), stage.ID()).
						HasResourceMonitorString(resourceMonitor.ID().Name()).
						HasPasswordPolicyString(passwordPolicy.ID().FullyQualifiedName()).
						HasSessionPolicyString(sessionPolicy.ID().FullyQualifiedName()),
				),
			},
			{
				Config:       config.FromModels(t, provider, completeConfigModel),
				ResourceName: completeConfigModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedCurrentOrganizationAccountResource(t, currentOrganizationAccountName).
						HasNameString(currentOrganizationAccountName).
						HasAllParametersEqualToPredefinedValues(warehouseId, eventTable.ID(), externalVolumeId, networkPolicy.ID(), stage.ID()).
						HasNoResourceMonitor().
						HasPasswordPolicyString(passwordPolicy.ID().FullyQualifiedName()).
						HasSessionPolicyString(sessionPolicy.ID().FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_CurrentOrganizationAccount_NameValidationOnCreate(t *testing.T) {
	testClient().EnsureValidNonProdOrganizationAccountIsUsed(t)

	provider := providermodel.SnowflakeProvider().WithWarehouse(testClient().Ids.WarehouseId().FullyQualifiedName())
	organizationAccountName := "INVALID_ORGANIZATION_ACCOUNT_NAME"
	modelWithInvalidOrganizationAccountName := model.CurrentOrganizationAccount("test", organizationAccountName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, provider, modelWithInvalidOrganizationAccountName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("passed name: %s, doesn't match current organization account name: %s, renames can be performed only after resource initialization", organizationAccountName, testClient().OrganizationAccount.ShowCurrent(t).AccountName)),
			},
		},
	})
}

// TestAcc_CurrentOrganizationAccount_NonEmptyComment_OnCreate shows that as long as resource configuration matches the infrastructure state, no changes are expected in the plan after create.
func TestAcc_CurrentOrganizationAccount_NonEmptyComment_OnCreate(t *testing.T) {
	testClient().EnsureValidNonProdOrganizationAccountIsUsed(t)

	comment := random.Comment()

	// We start with an organization account that already has a comment set, variation with initial empty comment is tested in TestAcc_CurrentOrganizationAccount_NonParameterValues
	testClient().OrganizationAccount.Alter(t, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithComment(comment)))
	t.Cleanup(func() {
		testClient().OrganizationAccount.Alter(t, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithComment(true)))
	})

	currentOrganizationAccountName := testClient().OrganizationAccount.ShowCurrent(t).AccountName
	completePropertiesModel := model.CurrentOrganizationAccount("test", currentOrganizationAccountName).WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create with import-like behavior
			{
				Config: config.FromModels(t, completePropertiesModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, completePropertiesModel.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasCommentString(comment),
					resourceshowoutputassert.OrganizationAccountShowOutput(t, completePropertiesModel.ResourceReference()).
						HasComment(comment),
				),
			},
			// Expect no changes as the resource configuration matches infrastructure state (comment was already set)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(completePropertiesModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, completePropertiesModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, completePropertiesModel.ResourceReference()).
						HasNameString(currentOrganizationAccountName).
						HasCommentString(comment),
					resourceshowoutputassert.OrganizationAccountShowOutput(t, completePropertiesModel.ResourceReference()).
						HasComment(comment),
				),
			},
		},
	})
}

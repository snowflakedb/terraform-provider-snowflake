//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
	unsetParametersModel := model.CurrentOrganizationAccount("test")
	setParametersModel := model.CurrentOrganizationAccount("test").WithAllParametersSetToPredefinedValues(warehouseId, eventTable.ID(), externalVolumeId, networkPolicy.ID(), stage.ID())

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
					resourceassert.CurrentOrganizationAccountResource(t, unsetParametersModel.ResourceReference()).HasAllDefaultParameters(),
				),
			},
			// import with unset parameters
			{
				Config:       config.FromModels(t, provider, unsetParametersModel),
				ResourceName: unsetParametersModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedCurrentOrganizationAccountResource(t, "current_organization_account").HasAllDefaultParameters(),
				),
			},
			// set all parameters
			{
				Config: config.FromModels(t, provider, setParametersModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, setParametersModel.ResourceReference()).
						HasAllParametersEqualToPredefinedValues(warehouseId, eventTable.ID(), externalVolumeId, networkPolicy.ID(), stage.ID()),
				),
			},
			// import with all parameters set
			{
				Config:       config.FromModels(t, provider, setParametersModel),
				ResourceName: setParametersModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedCurrentOrganizationAccountResource(t, "current_organization_account").
						HasAllParametersEqualToPredefinedValues(warehouseId, eventTable.ID(), externalVolumeId, networkPolicy.ID(), stage.ID()),
				),
			},
			// unset parameters
			{
				Config: config.FromModels(t, provider, unsetParametersModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, unsetParametersModel.ResourceReference()).HasAllDefaultParameters(),
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
					resourceassert.CurrentOrganizationAccountResource(t, setParametersModel.ResourceReference()).HasAllDefaultParameters(),
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

	provider := providermodel.SnowflakeProvider().WithWarehouse(testClient().Ids.WarehouseId().FullyQualifiedName())

	unsetModel := model.CurrentOrganizationAccount("test")

	setModel := model.CurrentOrganizationAccount("test").
		WithResourceMonitor(resourceMonitor.ID().Name()).
		WithPasswordPolicy(passwordPolicy.ID().FullyQualifiedName()).
		WithSessionPolicy(sessionPolicy.ID().FullyQualifiedName())

	setModelWithDifferentValues := model.CurrentOrganizationAccount("test").
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
						HasNoResourceMonitor().
						HasNoPasswordPolicy().
						HasNoSessionPolicy(),
				),
			},
			// import
			{
				Config:       config.FromModels(t, provider, unsetModel),
				ResourceName: unsetModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedCurrentOrganizationAccountResource(t, "current_organization_account").
						HasNoResourceMonitor().
						HasNoPasswordPolicy().
						HasNoSessionPolicy(),
				),
			},
			// set policies and resource monitor
			{
				Config: config.FromModels(t, provider, setModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, setModel.ResourceReference()).
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
					resourceassert.ImportedCurrentOrganizationAccountResource(t, "current_organization_account").
						HasNoResourceMonitor().
						HasPasswordPolicyString(passwordPolicy.ID().FullyQualifiedName()).
						HasSessionPolicyString(sessionPolicy.ID().FullyQualifiedName()),
				),
			},
			// set new policies
			{
				Config: config.FromModels(t, provider, setModelWithDifferentValues),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, setModelWithDifferentValues.ResourceReference()).
						HasResourceMonitorString(newResourceMonitor.ID().Name()).
						HasPasswordPolicyString(newPasswordPolicy.ID().FullyQualifiedName()).
						HasSessionPolicyString(newSessionPolicy.ID().FullyQualifiedName()),
				),
			},
			// unset policies and resource monitor
			{
				Config: config.FromModels(t, provider, unsetModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, unsetModel.ResourceReference()).
						HasResourceMonitorEmpty().
						HasPasswordPolicyEmpty().
						HasSessionPolicyEmpty(),
				),
			},
			// change externally
			{
				PreConfig: func() {
					testClient().OrganizationAccount.Alter(t, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithResourceMonitor(resourceMonitor.ID())))
					testClient().OrganizationAccount.Alter(t, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithPasswordPolicy(passwordPolicy.ID())))
					testClient().OrganizationAccount.Alter(t, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithSessionPolicy(sessionPolicy.ID())))
				},
				Config: config.FromModels(t, provider, unsetModel),
				Check: assertThat(t,
					resourceassert.CurrentOrganizationAccountResource(t, unsetModel.ResourceReference()).
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
	completeConfigModel := model.CurrentOrganizationAccount("test").
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
					resourceassert.ImportedCurrentOrganizationAccountResource(t, "current_organization_account").
						HasAllParametersEqualToPredefinedValues(warehouseId, eventTable.ID(), externalVolumeId, networkPolicy.ID(), stage.ID()).
						HasNoResourceMonitor().
						HasPasswordPolicyString(passwordPolicy.ID().FullyQualifiedName()).
						HasSessionPolicyString(sessionPolicy.ID().FullyQualifiedName()),
				),
			},
		},
	})
}

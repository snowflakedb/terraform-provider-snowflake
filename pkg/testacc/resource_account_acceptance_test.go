//go:build !account_level_tests

package testacc

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"regexp"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TESTS
// - Create with be
// - Set/Unset be
// - Upgrade to check if any changes are required

func TestAcc_Account_Minimal(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	organizationName := testClient().Context.CurrentAccountId(t).OrganizationName()
	id := sdk.NewAccountObjectIdentifier(random.AdminName())
	accountId := sdk.NewAccountIdentifier(organizationName, id.Name())
	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)
	region := testClient().Context.CurrentRegion(t)

	providerModel := providermodel.SnowflakeProvider().WithRole(snowflakeroles.Orgadmin.Name())

	configModel := model.Account("test", id.Name(), name, string(sdk.EditionStandard), email, 3).
		WithAdminRsaPublicKey(key)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providerModel, configModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminNameString(name).
						HasAdminRsaPublicKeyString(key).
						HasNoAdminUserType().
						HasEmailString(email).
						HasNoFirstName().
						HasNoLastName().
						HasMustChangePasswordString(r.BooleanDefault).
						HasNoRegionGroup().
						HasNoRegion().
						HasNoComment().
						HasIsOrgAdminString(r.BooleanDefault).
						HasGracePeriodInDaysString("3"),
					resourceshowoutputassert.AccountShowOutput(t, configModel.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasSnowflakeRegion(region).
						HasRegionGroup("").
						HasEdition(sdk.EditionStandard).
						HasAccountUrlNotEmpty().
						HasCreatedOnNotEmpty().
						HasComment("SNOWFLAKE").
						HasAccountLocatorNotEmpty().
						HasAccountLocatorUrlNotEmpty().
						HasManagedAccounts(0).
						HasConsumptionBillingEntityNameNotEmpty().
						HasMarketplaceConsumerBillingEntityName("").
						HasMarketplaceProviderBillingEntityNameNotEmpty().
						HasOldAccountURL("").
						HasIsOrgAdmin(false).
						HasAccountOldUrlSavedOnEmpty().
						HasAccountOldUrlLastUsedEmpty().
						HasOrganizationOldUrl("").
						HasOrganizationOldUrlSavedOnEmpty().
						HasOrganizationOldUrlLastUsedEmpty().
						HasIsEventsAccount(false).
						HasIsOrganizationAccount(false).
						HasDroppedOnEmpty().
						HasScheduledDeletionTimeEmpty().
						HasRestoredOnEmpty().
						HasMovedToOrganization("").
						HasMovedOn("").
						HasOrganizationUrlExpirationOnEmpty(),
				),
			},
			{
				ResourceName: configModel.ResourceReference(),
				Config:       config.FromModels(t, providerModel, configModel),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedAccountResource(t, helpers.EncodeResourceIdentifier(accountId)).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasNoAdminName().
						HasNoAdminRsaPublicKey().
						HasNoAdminUserType().
						HasNoEmail().
						HasNoFirstName().
						HasNoLastName().
						HasNoMustChangePassword().
						HasEditionString(string(sdk.EditionStandard)).
						HasNoRegionGroup().
						HasRegionString(region).
						HasCommentString("SNOWFLAKE").
						HasIsOrgAdminString(r.BooleanFalse).
						HasNoGracePeriodInDays(),
				),
			},
		},
	})
}

func TestAcc_Account_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	organizationName := testClient().Context.CurrentAccountId(t).OrganizationName()
	id := sdk.NewAccountObjectIdentifier(random.AdminName())
	accountId := sdk.NewAccountIdentifier(organizationName, id.Name())
	firstName := random.AlphaN(30)
	lastName := random.AlphaN(30)
	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)
	region := testClient().Context.CurrentRegion(t)
	comment := random.Comment()
	// TODO(SNOW-2131939): The default consumption billing entity consists of organization name followed by _DefaultBE
	defaultConsumptionBillingEntity := fmt.Sprintf("%s_DefaultBE", organizationName)

	providerModel := providermodel.SnowflakeProvider().WithRole(snowflakeroles.Orgadmin.Name())

	configModel := model.Account("test", id.Name(), name, string(sdk.EditionStandard), email, 3).
		WithAdminUserTypeEnum(sdk.UserTypePerson).
		WithAdminRsaPublicKey(key).
		WithFirstName(firstName).
		WithLastName(lastName).
		WithMustChangePassword(r.BooleanTrue).
		WithRegionGroup("PUBLIC").
		WithRegion(region).
		WithComment(comment).
		WithConsumptionBillingEntity(defaultConsumptionBillingEntity).
		WithIsOrgAdmin(r.BooleanFalse)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providerModel, configModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(sdk.NewAccountIdentifier(organizationName, id.Name()).FullyQualifiedName()).
						HasAdminNameString(name).
						HasAdminRsaPublicKeyString(key).
						HasAdminUserType(sdk.UserTypePerson).
						HasEmailString(email).
						HasFirstNameString(firstName).
						HasLastNameString(lastName).
						HasMustChangePasswordString(r.BooleanTrue).
						HasRegionGroupString("PUBLIC").
						HasRegionString(region).
						HasCommentString(comment).
						HasConsumptionBillingEntityString(defaultConsumptionBillingEntity).
						HasIsOrgAdminString(r.BooleanFalse).
						HasGracePeriodInDaysString("3"),
					resourceshowoutputassert.AccountShowOutput(t, configModel.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasSnowflakeRegion(region).
						HasRegionGroup("").
						HasEdition(sdk.EditionStandard).
						HasAccountUrlNotEmpty().
						HasCreatedOnNotEmpty().
						HasComment(comment).
						HasAccountLocatorNotEmpty().
						HasAccountLocatorUrlNotEmpty().
						HasManagedAccounts(0).
						HasConsumptionBillingEntityName(defaultConsumptionBillingEntity).
						HasMarketplaceConsumerBillingEntityName("").
						HasMarketplaceProviderBillingEntityNameNotEmpty().
						HasOldAccountURL("").
						HasIsOrgAdmin(false).
						HasAccountOldUrlSavedOnEmpty().
						HasAccountOldUrlLastUsedEmpty().
						HasOrganizationOldUrl("").
						HasOrganizationOldUrlSavedOnEmpty().
						HasOrganizationOldUrlLastUsedEmpty().
						HasIsEventsAccount(false).
						HasIsOrganizationAccount(false).
						HasDroppedOnEmpty().
						HasScheduledDeletionTimeEmpty().
						HasRestoredOnEmpty().
						HasMovedToOrganization("").
						HasMovedOn("").
						HasOrganizationUrlExpirationOnEmpty(),
				),
			},
			{
				ResourceName: configModel.ResourceReference(),
				Config:       config.FromModels(t, providerModel, configModel),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedAccountResource(t, helpers.EncodeResourceIdentifier(accountId)).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(sdk.NewAccountIdentifier(organizationName, id.Name()).FullyQualifiedName()).
						HasNoAdminName().
						HasNoAdminRsaPublicKey().
						HasNoEmail().
						HasNoFirstName().
						HasNoLastName().
						HasNoAdminUserType().
						HasNoMustChangePassword().
						HasEditionString(string(sdk.EditionStandard)).
						HasNoRegionGroup().
						HasRegionString(region).
						HasCommentString(comment).
						HasConsumptionBillingEntityString(defaultConsumptionBillingEntity).
						HasIsOrgAdminString(r.BooleanFalse).
						HasNoGracePeriodInDays(),
				),
			},
		},
	})
}

func TestAcc_Account_Rename(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	organizationName := testClient().Context.CurrentAccountId(t).OrganizationName()
	id := sdk.NewAccountObjectIdentifier(random.AdminName())
	accountId := sdk.NewAccountIdentifier(organizationName, id.Name())

	newId := testClient().Ids.RandomSensitiveAccountObjectIdentifier()
	newAccountId := sdk.NewAccountIdentifier(organizationName, newId.Name())

	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)

	providerModel := providermodel.SnowflakeProvider().WithRole(snowflakeroles.Orgadmin.Name())

	configModel := model.Account("test", id.Name(), name, string(sdk.EditionStandard), email, 3).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key)

	newConfigModel := model.Account("test", newId.Name(), name, string(sdk.EditionStandard), email, 3).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providerModel, configModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminUserType(sdk.UserTypeService),
					resourceshowoutputassert.AccountShowOutput(t, configModel.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(newConfigModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, newConfigModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, newConfigModel.ResourceReference()).
						HasNameString(newId.Name()).
						HasFullyQualifiedNameString(newAccountId.FullyQualifiedName()).
						HasAdminUserType(sdk.UserTypeService),
					resourceshowoutputassert.AccountShowOutput(t, newConfigModel.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(newId.Name()),
				),
			},
		},
	})
}

func TestAcc_Account_IsOrgAdmin(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	organizationName := testClient().Context.CurrentAccountId(t).OrganizationName()
	id := sdk.NewAccountObjectIdentifier(random.AdminName())
	accountId := sdk.NewAccountIdentifier(organizationName, id.Name())

	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)

	providerModel := providermodel.SnowflakeProvider().WithRole(snowflakeroles.Orgadmin.Name())

	configModelWithOrgAdminTrue := model.Account("test", id.Name(), name, string(sdk.EditionStandard), email, 3).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key).
		WithIsOrgAdmin(r.BooleanTrue)

	configModelWithOrgAdminFalse := model.Account("test", id.Name(), name, string(sdk.EditionStandard), email, 3).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key).
		WithIsOrgAdmin(r.BooleanFalse)

	configModelWithoutOrgAdmin := model.Account("test", id.Name(), name, string(sdk.EditionStandard), email, 3).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			// Create with ORGADMIN enabled
			{
				Config: config.FromModels(t, providerModel, configModelWithOrgAdminTrue),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModelWithOrgAdminTrue.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminUserType(sdk.UserTypeService).
						HasIsOrgAdminString(r.BooleanTrue),
					resourceshowoutputassert.AccountShowOutput(t, configModelWithOrgAdminTrue.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasIsOrgAdmin(true),
				),
			},
			// Disable ORGADMIN
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(configModelWithOrgAdminFalse.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, configModelWithOrgAdminFalse),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModelWithOrgAdminFalse.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminUserType(sdk.UserTypeService).
						HasIsOrgAdminString(r.BooleanFalse),
					resourceshowoutputassert.AccountShowOutput(t, configModelWithOrgAdminFalse.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasIsOrgAdmin(false),
				),
			},
			// Remove is_org_admin from the config and go back to default (disabled)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(configModelWithoutOrgAdmin.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, configModelWithoutOrgAdmin),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModelWithoutOrgAdmin.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminUserType(sdk.UserTypeService).
						HasIsOrgAdminString(r.BooleanDefault),
					resourceshowoutputassert.AccountShowOutput(t, configModelWithoutOrgAdmin.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasIsOrgAdmin(false),
				),
			},
			// External change (enable ORGADMIN)
			{
				PreConfig: func() {
					testClient().Account.Alter(t, &sdk.AlterAccountOptions{
						Name: sdk.Pointer(accountId.AsAccountObjectIdentifier()),
						Set:  &sdk.AccountSet{OrgAdmin: sdk.Bool(true)},
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(configModelWithoutOrgAdmin.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, configModelWithoutOrgAdmin),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModelWithoutOrgAdmin.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminUserType(sdk.UserTypeService).
						HasIsOrgAdminString(r.BooleanDefault),
					resourceshowoutputassert.AccountShowOutput(t, configModelWithoutOrgAdmin.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasIsOrgAdmin(false),
				),
			},
		},
	})
}

func TestAcc_Account_UpdatingConsumptionBillingEntity(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	organizationName := testClient().Context.CurrentAccountId(t).OrganizationName()
	id := sdk.NewAccountObjectIdentifier(random.AdminName())
	accountId := sdk.NewAccountIdentifier(organizationName, id.Name())
	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)

	// TODO(SNOW-2131939): The default consumption billing entity consists of organization name followed by _DefaultBE
	defaultConsumptionBillingEntity := fmt.Sprintf("%s_DefaultBE", organizationName)

	providerModel := providermodel.SnowflakeProvider().WithRole(snowflakeroles.Orgadmin.Name())

	configModel := model.Account("test", id.Name(), name, string(sdk.EditionStandard), email, 3).
		WithAdminRsaPublicKey(key)

	configModelWithConsumptionBillingEntity := model.Account("test", id.Name(), name, string(sdk.EditionStandard), email, 3).
		WithAdminRsaPublicKey(key).
		WithConsumptionBillingEntity(defaultConsumptionBillingEntity)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providerModel, configModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasNoConsumptionBillingEntity(),
					resourceshowoutputassert.AccountShowOutput(t, configModel.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasConsumptionBillingEntityName(defaultConsumptionBillingEntity),
				),
			},
			// Set consumption billing entity
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(configModelWithConsumptionBillingEntity.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, configModelWithConsumptionBillingEntity),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModelWithConsumptionBillingEntity.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasConsumptionBillingEntityString(defaultConsumptionBillingEntity),
					resourceshowoutputassert.AccountShowOutput(t, configModelWithConsumptionBillingEntity.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasConsumptionBillingEntityName(defaultConsumptionBillingEntity),
				),
			},
			// Unset consumption billing entity
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(configModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, configModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasConsumptionBillingEntityEmpty(),
					resourceshowoutputassert.AccountShowOutput(t, configModel.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasConsumptionBillingEntityName(defaultConsumptionBillingEntity),
				),
			},
		},
		// TODO(SNOW-2131939): Test with external changes
	})
}

func TestAcc_Account_IgnoreUpdateAfterCreationOnCertainFields(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	organizationName := testClient().Context.CurrentAccountId(t).OrganizationName()
	id := sdk.NewAccountObjectIdentifier(random.AdminName())
	accountId := sdk.NewAccountIdentifier(organizationName, id.Name())

	firstName := random.AlphaN(30)
	lastName := random.AlphaN(30)
	email := random.Email()
	name := random.AdminName()
	pass := random.Password()

	newFirstName := random.AlphaN(30)
	newLastName := random.AlphaN(30)
	newEmail := random.Email()
	newName := random.AdminName()
	newPass := random.Password()

	providerModel := providermodel.SnowflakeProvider().WithRole(snowflakeroles.Orgadmin.Name())

	configModel := model.Account("test", id.Name(), name, string(sdk.EditionStandard), email, 3).
		WithAdminUserTypeEnum(sdk.UserTypePerson).
		WithFirstName(firstName).
		WithLastName(lastName).
		WithMustChangePassword(r.BooleanTrue).
		WithAdminPassword(pass)

	newConfigModel := model.Account("test", id.Name(), newName, string(sdk.EditionStandard), newEmail, 3).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminPassword(newPass).
		WithFirstName(newFirstName).
		WithLastName(newLastName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providerModel, configModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminNameString(name).
						HasAdminPasswordString(pass).
						HasAdminUserType(sdk.UserTypePerson).
						HasEmailString(email).
						HasFirstNameString(firstName).
						HasLastNameString(lastName).
						HasMustChangePasswordString(r.BooleanTrue),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(newConfigModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, providerModel, newConfigModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, newConfigModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminNameString(name).
						HasAdminPasswordString(pass).
						HasAdminUserType(sdk.UserTypePerson).
						HasEmailString(email).
						HasFirstNameString(firstName).
						HasLastNameString(lastName).
						HasMustChangePasswordString(r.BooleanTrue),
				),
			},
		},
	})
}

func TestAcc_Account_TryToCreateWithoutOrgadmin(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	id := sdk.NewAccountObjectIdentifier(random.AdminName())
	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)

	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	t.Setenv(snowflakeenvs.Role, snowflakeroles.Accountadmin.Name())

	providerModel := providermodel.SnowflakeProvider().WithRole(snowflakeroles.Orgadmin.Name())

	configModel := model.Account("test", id.Name(), name, string(sdk.EditionStandard), email, 3).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, providerModel, configModel),
				ExpectError: regexp.MustCompile("Error: current user doesn't have the orgadmin role in session"),
			},
		},
	})
}

func TestAcc_Account_InvalidValues(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	id := sdk.NewAccountObjectIdentifier(random.AdminName())
	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)

	providerModel := providermodel.SnowflakeProvider().WithRole(snowflakeroles.Orgadmin.Name())

	configModelInvalidUserType := model.Account("test", id.Name(), name, string(sdk.EditionStandard), email, 3).
		WithAdminUserType("invalid_user_type").
		WithAdminRsaPublicKey(key)

	configModelInvalidAccountEdition := model.Account("test", id.Name(), name, "invalid_account_edition", email, 3).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key)

	configModelInvalidGracePeriodInDays := model.Account("test", id.Name(), name, string(sdk.EditionStandard), email, 2).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, providerModel, configModelInvalidUserType),
				ExpectError: regexp.MustCompile("invalid user type: invalid_user_type"),
			},
			{
				Config:      config.FromModels(t, providerModel, configModelInvalidAccountEdition),
				ExpectError: regexp.MustCompile("unknown account edition: invalid_account_edition"),
			},
			{
				Config:      config.FromModels(t, providerModel, configModelInvalidGracePeriodInDays),
				ExpectError: regexp.MustCompile(`Error: expected grace_period_in_days to be at least \(3\), got 2`),
			},
		},
	})
}

func TestAcc_Account_UpgradeFrom_v0_99_0(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	id := sdk.NewAccountObjectIdentifier(random.AdminName())
	email := random.Email()
	adminName := random.AdminName()
	adminPassword := random.Password()
	firstName := random.AlphaN(30)
	lastName := random.AlphaN(30)
	region := testClient().Context.CurrentRegion(t)
	comment := random.Comment()

	providerModel := providermodel.SnowflakeProvider().WithRole(snowflakeroles.Orgadmin.Name())

	configModel := model.Account("test", id.Name(), adminName, string(sdk.EditionStandard), email, 3).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminPassword(adminPassword).
		WithFirstName(firstName).
		WithLastName(lastName).
		WithMustChangePasswordValue(tfconfig.BoolVariable(true)).
		WithRegion(region).
		WithIsOrgAdmin(r.BooleanFalse).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.99.0"),
				Config:            accountConfig_v0_99_0(id.Name(), adminName, adminPassword, email, sdk.EditionStandard, firstName, lastName, true, region, 3, comment),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, providerModel, configModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasAdminNameString(adminName).
						HasAdminPasswordString(adminPassword).
						HasEmailString(email).
						HasFirstNameString(firstName).
						HasLastNameString(lastName).
						HasMustChangePasswordString(r.BooleanTrue).
						HasRegionGroupString("").
						HasRegionString(region).
						HasCommentString(comment).
						HasIsOrgAdminString(r.BooleanFalse).
						HasGracePeriodInDaysString("3"),
				),
			},
		},
	})
}

func accountConfig_v0_99_0(
	name string,
	adminName string,
	adminPassword string,
	email string,
	edition sdk.AccountEdition,
	firstName string,
	lastName string,
	mustChangePassword bool,
	region string,
	gracePeriodInDays int,
	comment string,
) string {
	return fmt.Sprintf(`
provider "snowflake" {
	role = "ORGADMIN"
}

resource "snowflake_account" "test" {
	name = "%[1]s"
	admin_name = "%[2]s"
	admin_password = "%[3]s"
	email = "%[4]s"
	edition = "%[5]s"
	first_name = "%[6]s"
	last_name = "%[7]s"
	must_change_password = %[8]t
	region = "%[9]s"
	grace_period_in_days = %[10]d
	comment = "%[11]s"
}
`,
		name,
		adminName,
		adminPassword,
		email,
		edition,
		firstName,
		lastName,
		mustChangePassword,
		region,
		gracePeriodInDays,
		comment,
	)
}

// TODO(SNOW-1875369): add a state upgrader test for an imported account with optional parameters

//go:build account_level_tests

// These tests are temporarily moved to account level tests due to flakiness caused by changes in the higher-level parameters.

package testacc

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Users_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	id := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "1")
	id2 := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "2")
	id3 := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("u", id.Name())
	user2Model := model.User("u2", id2.Name())
	user3Model := model.User("u3", id3.Name())

	usersModelLikeFirstOne := datasourcemodel.Users("test").
		WithWithDescribe(false).
		WithWithParameters(false).
		WithLike(id.Name()).
		WithDependsOn(userModel.ResourceReference(), user2Model.ResourceReference(), user3Model.ResourceReference())

	usersModelStartsWithPrefix := datasourcemodel.Users("test").
		WithWithDescribe(false).
		WithWithParameters(false).
		WithStartsWith(prefix).
		WithDependsOn(userModel.ResourceReference(), user2Model.ResourceReference(), user3Model.ResourceReference())

	usersModelLimitRowsAndFrom := datasourcemodel.Users("test").
		WithWithDescribe(false).
		WithWithParameters(false).
		WithLimitRowsAndFrom(1, prefix).
		WithDependsOn(userModel.ResourceReference(), user2Model.ResourceReference(), user3Model.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModel, user2Model, user3Model, usersModelLikeFirstOne),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(usersModelLikeFirstOne.DatasourceReference(), "users.#", "1"),
				),
			},
			{
				Config: config.FromModels(t, userModel, user2Model, user3Model, usersModelStartsWithPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(usersModelStartsWithPrefix.DatasourceReference(), "users.#", "2"),
				),
			},
			{
				Config: config.FromModels(t, userModel, user2Model, user3Model, usersModelLimitRowsAndFrom),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(usersModelLimitRowsAndFrom.DatasourceReference(), "users.#", "1"),
				),
			},
		},
	})
}

func TestAcc_Users_CompleteUseCase(t *testing.T) {
	personUserId := testClient().Ids.RandomAccountObjectIdentifier()
	serviceUserId := testClient().Ids.RandomAccountObjectIdentifier()
	legacyServiceUserId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	pass := random.Password()
	key1, key1Fp := random.GenerateRSAPublicKey(t)
	key2, key2Fp := random.GenerateRSAPublicKey(t)

	personUserModel := model.User("person", personUserId.Name()).
		WithPassword(pass).
		WithLoginName(personUserId.Name() + "_login").
		WithDisplayName("Display Name").
		WithFirstName("Jan").
		WithMiddleName("Jakub").
		WithLastName("Testowski").
		WithEmail("fake@email.com").
		WithMustChangePassword("true").
		WithDisabled("false").
		WithDaysToExpiry(8).
		WithMinsToUnlock(9).
		WithDefaultWarehouse("some_warehouse").
		WithDefaultNamespace("some.namespace").
		WithDefaultRole("some_role").
		WithDefaultSecondaryRolesOptionEnum(sdk.SecondaryRolesOptionAll).
		WithMinsToBypassMfa(10).
		WithRsaPublicKey(key1).
		WithRsaPublicKey2(key2).
		WithComment(comment).
		WithDisableMfa("true")

	serviceUserModel := model.ServiceUser("service", serviceUserId.Name()).
		WithLoginName(serviceUserId.Name() + "_login").
		WithDisplayName("Service Display Name").
		WithEmail("service@email.com").
		WithDisabled("false").
		WithDaysToExpiry(8).
		WithMinsToUnlock(9).
		WithDefaultWarehouse("some_warehouse").
		WithDefaultNamespace("some.namespace").
		WithDefaultRole("some_role").
		WithDefaultSecondaryRolesOptionEnum(sdk.SecondaryRolesOptionAll).
		WithRsaPublicKey(key1).
		WithRsaPublicKey2(key2).
		WithComment(comment)

	legacyServiceUserModel := model.LegacyServiceUser("legacy", legacyServiceUserId.Name()).
		WithPassword(pass).
		WithLoginName(legacyServiceUserId.Name() + "_login").
		WithDisplayName("Legacy Display Name").
		WithEmail("legacy@email.com").
		WithMustChangePassword("true").
		WithDisabled("false").
		WithDaysToExpiry(8).
		WithMinsToUnlock(9).
		WithDefaultWarehouse("some_warehouse").
		WithDefaultNamespace("some.namespace").
		WithDefaultRole("some_role").
		WithDefaultSecondaryRolesOptionEnum(sdk.SecondaryRolesOptionAll).
		WithRsaPublicKey(key1).
		WithRsaPublicKey2(key2).
		WithComment(comment)

	usersModelPersonWithoutOptionals := datasourcemodel.Users("test").
		WithLike(personUserId.Name()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(personUserModel.ResourceReference())

	usersModelPerson := datasourcemodel.Users("test").
		WithLike(personUserId.Name()).
		WithDependsOn(personUserModel.ResourceReference())

	usersModelService := datasourcemodel.Users("test").
		WithLike(serviceUserId.Name()).
		WithDependsOn(serviceUserModel.ResourceReference())

	usersModelLegacyService := datasourcemodel.Users("test").
		WithLike(legacyServiceUserId.Name()).
		WithDependsOn(legacyServiceUserModel.ResourceReference())

	userCommonShowAssert := func(t *testing.T, datasourceRef string) *resourceshowoutputassert.UserShowOutputAssert {
		t.Helper()
		return resourceshowoutputassert.UsersDatasourceShowOutput(t, datasourceRef).
			HasCreatedOnNotEmpty().
			HasDisabled(false).
			HasSnowflakeLock(false).
			HasDaysToExpiryNotEmpty().
			HasMinsToUnlockNotEmpty().
			HasDefaultWarehouse("some_warehouse").
			HasDefaultNamespace("some.namespace").
			HasDefaultRole("some_role").
			HasDefaultSecondaryRoles(`["ALL"]`).
			HasHasRsaPublicKey(true).
			HasComment(comment).
			HasHasMfa(false)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// Person user WITHOUT additional outputs
			{
				Config: config.FromModels(t, personUserModel, usersModelPersonWithoutOptionals),
				Check: assertThat(t,
					userCommonShowAssert(t, usersModelPersonWithoutOptionals.DatasourceReference()).
						HasName(personUserId.Name()).
						HasType(string(sdk.UserTypePerson)).
						HasLoginName(strings.ToUpper(fmt.Sprintf("%s_LOGIN", personUserId.Name()))).
						HasDisplayName("Display Name").
						HasFirstName("Jan").
						HasLastName("Testowski").
						HasEmail("fake@email.com").
						HasMustChangePassword(true).
						HasMinsToBypassMfaNotEmpty(),

					assert.Check(resource.TestCheckResourceAttr(usersModelPersonWithoutOptionals.DatasourceReference(), "users.0.describe_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPersonWithoutOptionals.DatasourceReference(), "users.0.parameters.#", "0")),
				),
			},
			// Person user WITH additional outputs
			{
				Config: config.FromModels(t, personUserModel, usersModelPerson),
				Check: assertThat(t,
					userCommonShowAssert(t, usersModelPerson.DatasourceReference()).
						HasName(personUserId.Name()).
						HasType(string(sdk.UserTypePerson)).
						HasLoginName(strings.ToUpper(fmt.Sprintf("%s_LOGIN", personUserId.Name()))).
						HasDisplayName("Display Name").
						HasFirstName("Jan").
						HasLastName("Testowski").
						HasEmail("fake@email.com").
						HasMustChangePassword(true).
						HasMinsToBypassMfaNotEmpty(),

					resourceparametersassert.UsersDatasourceParameters(t, usersModelPerson.DatasourceReference()).
						HasAllDefaults(),

					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.name", personUserId.Name())),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.type", string(sdk.UserTypePerson))),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.display_name", "Display Name")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.login_name", strings.ToUpper(fmt.Sprintf("%s_LOGIN", personUserId.Name())))),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.first_name", "Jan")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.middle_name", "Jakub")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.last_name", "Testowski")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.email", "fake@email.com")),
					assert.Check(resource.TestCheckNoResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.password")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.must_change_password", "true")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.disabled", "false")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.snowflake_lock", "false")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.snowflake_support", "false")),
					assert.Check(resource.TestCheckResourceAttrSet(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.days_to_expiry")),
					assert.Check(resource.TestCheckResourceAttrSet(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.mins_to_unlock")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.default_warehouse", "some_warehouse")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.default_namespace", "some.namespace")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.default_role", "some_role")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.default_secondary_roles", `["ALL"]`)),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.ext_authn_duo", "false")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.ext_authn_uid", "")),
					assert.Check(resource.TestCheckResourceAttrSet(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.mins_to_bypass_mfa")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.mins_to_bypass_network_policy", "0")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.rsa_public_key", key1)),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.rsa_public_key_fp", "SHA256:"+key1Fp)),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.rsa_public_key2", key2)),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.rsa_public_key2_fp", "SHA256:"+key2Fp)),
					assert.Check(resource.TestCheckResourceAttrSet(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.password_last_set_time")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.custom_landing_page_url", "")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.custom_landing_page_url_flush_next_ui_load", "false")),
					assert.Check(resource.TestCheckResourceAttr(usersModelPerson.DatasourceReference(), "users.0.describe_output.0.has_mfa", "false")),
				),
			},
			//  Service user WITH additional outputs
			{
				Config: config.FromModels(t, serviceUserModel, usersModelService),
				Check: assertThat(t,
					userCommonShowAssert(t, usersModelService.DatasourceReference()).
						HasName(serviceUserId.Name()).
						HasType(string(sdk.UserTypeService)).
						HasLoginName(strings.ToUpper(fmt.Sprintf("%s_LOGIN", serviceUserId.Name()))).
						HasDisplayName("Service Display Name").
						HasFirstName("").
						HasLastName("").
						HasEmail("service@email.com").
						HasMustChangePassword(false).
						HasMinsToBypassMfa(""),

					resourceparametersassert.UsersDatasourceParameters(t, usersModelService.DatasourceReference()).
						HasAllDefaults(),

					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.name", serviceUserId.Name())),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.type", string(sdk.UserTypeService))),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.display_name", "Service Display Name")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.login_name", strings.ToUpper(fmt.Sprintf("%s_LOGIN", serviceUserId.Name())))),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.first_name", "")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.middle_name", "")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.last_name", "")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.email", "service@email.com")),
					assert.Check(resource.TestCheckNoResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.password")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.must_change_password", "false")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.disabled", "false")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.snowflake_lock", "false")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.snowflake_support", "false")),
					assert.Check(resource.TestCheckResourceAttrSet(usersModelService.DatasourceReference(), "users.0.describe_output.0.days_to_expiry")),
					assert.Check(resource.TestCheckResourceAttrSet(usersModelService.DatasourceReference(), "users.0.describe_output.0.mins_to_unlock")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.default_warehouse", "some_warehouse")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.default_namespace", "some.namespace")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.default_role", "some_role")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.default_secondary_roles", `["ALL"]`)),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.ext_authn_duo", "false")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.ext_authn_uid", "")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.mins_to_bypass_mfa", "0")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.mins_to_bypass_network_policy", "0")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.rsa_public_key", key1)),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.rsa_public_key_fp", "SHA256:"+key1Fp)),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.rsa_public_key2", key2)),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.rsa_public_key2_fp", "SHA256:"+key2Fp)),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.password_last_set_time", "")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.custom_landing_page_url", "")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.custom_landing_page_url_flush_next_ui_load", "false")),
					assert.Check(resource.TestCheckResourceAttr(usersModelService.DatasourceReference(), "users.0.describe_output.0.has_mfa", "false")),
				),
			},
			//  Legacy service user WITH additional outputs
			{
				Config: config.FromModels(t, legacyServiceUserModel, usersModelLegacyService),
				Check: assertThat(t,
					userCommonShowAssert(t, usersModelLegacyService.DatasourceReference()).
						HasName(legacyServiceUserId.Name()).
						HasType(string(sdk.UserTypeLegacyService)).
						HasLoginName(strings.ToUpper(fmt.Sprintf("%s_LOGIN", legacyServiceUserId.Name()))).
						HasDisplayName("Legacy Display Name").
						HasFirstName("").
						HasLastName("").
						HasEmail("legacy@email.com").
						HasMustChangePassword(true).
						HasMinsToBypassMfa(""),

					resourceparametersassert.UsersDatasourceParameters(t, usersModelLegacyService.DatasourceReference()).
						HasAllDefaults(),

					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.name", legacyServiceUserId.Name())),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.type", string(sdk.UserTypeLegacyService))),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.display_name", "Legacy Display Name")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.login_name", strings.ToUpper(fmt.Sprintf("%s_LOGIN", legacyServiceUserId.Name())))),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.first_name", "")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.middle_name", "")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.last_name", "")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.email", "legacy@email.com")),
					assert.Check(resource.TestCheckNoResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.password")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.must_change_password", "true")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.disabled", "false")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.snowflake_lock", "false")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.snowflake_support", "false")),
					assert.Check(resource.TestCheckResourceAttrSet(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.days_to_expiry")),
					assert.Check(resource.TestCheckResourceAttrSet(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.mins_to_unlock")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.default_warehouse", "some_warehouse")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.default_namespace", "some.namespace")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.default_role", "some_role")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.default_secondary_roles", `["ALL"]`)),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.ext_authn_duo", "false")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.ext_authn_uid", "")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.mins_to_bypass_mfa", "0")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.mins_to_bypass_network_policy", "0")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.rsa_public_key", key1)),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.rsa_public_key_fp", "SHA256:"+key1Fp)),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.rsa_public_key2", key2)),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.rsa_public_key2_fp", "SHA256:"+key2Fp)),
					assert.Check(resource.TestCheckResourceAttrSet(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.password_last_set_time")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.custom_landing_page_url", "")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.custom_landing_page_url_flush_next_ui_load", "false")),
					assert.Check(resource.TestCheckResourceAttr(usersModelLegacyService.DatasourceReference(), "users.0.describe_output.0.has_mfa", "false")),
				),
			},
		},
	})
}

func TestAcc_Users_UserNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Users/without_user"),
				ExpectError:     regexp.MustCompile("there should be at least one user"),
			},
		},
	})
}

//go:build account_level_tests

package testacc

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_User_BasicFlows(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	id2 := testClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()
	newComment := random.Comment()

	key1, _ := random.GenerateRSAPublicKey(t)
	key2, _ := random.GenerateRSAPublicKey(t)

	loginName := random.SensitiveAlphanumeric()
	newLoginName := random.SensitiveAlphanumeric()
	pass := random.Password()
	newPass := random.Password()

	userModelNoAttributes := model.User("w", id.Name())
	userModelNoAttributesRenamed := model.User("w", id2.Name()).
		WithComment(newComment)

	userModelAllAttributes := model.User("w", id.Name()).
		WithPassword(pass).
		WithLoginName(loginName).
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

	userModelAllAttributesChanged := func(loginName string) *model.UserModel {
		return model.User("w", id.Name()).
			WithPassword(newPass).
			WithLoginName(loginName).
			WithDisplayName("New Display Name").
			WithFirstName("Janek").
			WithMiddleName("Kuba").
			WithLastName("Terraformowski").
			WithEmail("fake@email.net").
			WithMustChangePassword("false").
			WithDisabled("true").
			WithDaysToExpiry(12).
			WithMinsToUnlock(13).
			WithDefaultWarehouse("other_warehouse").
			WithDefaultNamespace("one_part_namespace").
			WithDefaultRole("other_role").
			WithDefaultSecondaryRolesOptionEnum(sdk.SecondaryRolesOptionAll).
			WithMinsToBypassMfa(14).
			WithRsaPublicKey(key2).
			WithRsaPublicKey2(key1).
			WithComment(newComment).
			WithDisableMfa("false")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// CREATE WITHOUT ATTRIBUTES
			{
				Config: config.FromModels(t, userModelNoAttributes),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasNoPassword().
						HasNoLoginName().
						HasNoDisplayName().
						HasNoFirstName().
						HasNoMiddleName().
						HasNoLastName().
						HasNoEmail().
						HasMustChangePasswordString(r.BooleanDefault).
						HasDisabledString(r.BooleanDefault).
						HasNoDaysToExpiry().
						HasMinsToUnlockString(r.IntDefaultString).
						HasNoDefaultWarehouse().
						HasNoDefaultNamespace().
						HasNoDefaultRole().
						HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionDefault).
						HasMinsToBypassMfaString(r.IntDefaultString).
						HasNoRsaPublicKey().
						HasNoRsaPublicKey2().
						HasNoComment().
						HasDisableMfaString(r.BooleanDefault).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.UserShowOutput(t, userModelNoAttributes.ResourceReference()).
						HasLoginName(strings.ToUpper(id.Name())).
						HasDisplayName(id.Name()),
				),
			},
			// RENAME AND CHANGE ONE PROP
			{
				Config: config.FromModels(t, userModelNoAttributesRenamed),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelNoAttributes.ResourceReference()).
						HasNameString(id2.Name()).
						HasCommentString(newComment),
					// default names stay the same
					resourceshowoutputassert.UserShowOutput(t, userModelNoAttributes.ResourceReference()).
						HasLoginName(strings.ToUpper(id.Name())).
						HasDisplayName(id.Name()),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelNoAttributesRenamed.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "disable_mfa", "days_to_expiry", "mins_to_unlock", "mins_to_bypass_mfa", "login_name", "display_name", "disabled", "must_change_password", "default_secondary_roles_option"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedUserResource(t, id2.Name()).
						HasLoginNameString(strings.ToUpper(id.Name())).
						HasDisplayNameString(id.Name()).
						HasDisabled(false).
						HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionAll).
						HasMustChangePassword(false),
				),
			},
			// DESTROY
			{
				Config:  config.FromModels(t, userModelNoAttributes),
				Destroy: true,
			},
			// CREATE WITH ALL ATTRIBUTES
			{
				Config: config.FromModels(t, userModelAllAttributes),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelAllAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasPasswordString(pass).
						HasLoginNameString(loginName).
						HasDisplayNameString("Display Name").
						HasFirstNameString("Jan").
						HasMiddleNameString("Jakub").
						HasLastNameString("Testowski").
						HasEmailString("fake@email.com").
						HasMustChangePassword(true).
						HasDisabled(false).
						HasDaysToExpiryString("8").
						HasMinsToUnlockString("9").
						HasDefaultWarehouseString("some_warehouse").
						HasDefaultNamespaceString("some.namespace").
						HasDefaultRoleString("some_role").
						HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionAll).
						HasMinsToBypassMfaString("10").
						HasRsaPublicKeyString(key1).
						HasRsaPublicKey2String(key2).
						HasCommentString(comment).
						HasDisableMfaString(r.BooleanTrue).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: config.FromModels(t, userModelAllAttributesChanged(newLoginName)),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelAllAttributesChanged(newLoginName).ResourceReference()).
						HasNameString(id.Name()).
						HasPasswordString(newPass).
						HasLoginNameString(newLoginName).
						HasDisplayNameString("New Display Name").
						HasFirstNameString("Janek").
						HasMiddleNameString("Kuba").
						HasLastNameString("Terraformowski").
						HasEmailString("fake@email.net").
						HasMustChangePassword(false).
						HasDisabled(true).
						HasDaysToExpiryString("12").
						HasMinsToUnlockString("13").
						HasDefaultWarehouseString("other_warehouse").
						HasDefaultNamespaceString("one_part_namespace").
						HasDefaultRoleString("other_role").
						HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionAll).
						HasMinsToBypassMfaString("14").
						HasRsaPublicKeyString(key2).
						HasRsaPublicKey2String(key1).
						HasCommentString(newComment).
						HasDisableMfaString(r.BooleanFalse).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelAllAttributesChanged(newLoginName).ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "disable_mfa", "days_to_expiry", "mins_to_unlock", "mins_to_bypass_mfa", "default_namespace", "login_name", "show_output.0.days_to_expiry"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedUserResource(t, id.Name()).
						HasDefaultNamespaceString("ONE_PART_NAMESPACE").
						HasLoginNameString(strings.ToUpper(newLoginName)),
				),
			},
			// CHANGE PROP TO THE CURRENT SNOWFLAKE VALUE
			{
				PreConfig: func() {
					testClient().User.SetLoginName(t, id, loginName)
				},
				Config: config.FromModels(t, userModelAllAttributesChanged(loginName)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// UNSET ALL
			{
				Config: config.FromModels(t, userModelNoAttributes),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasPasswordString("").
						HasLoginNameString("").
						HasDisplayNameString("").
						HasFirstNameString("").
						HasMiddleNameString("").
						HasLastNameString("").
						HasEmailString("").
						HasMustChangePasswordString(r.BooleanDefault).
						HasDisabledString(r.BooleanDefault).
						HasDaysToExpiryString("0").
						HasMinsToUnlockString(r.IntDefaultString).
						HasDefaultWarehouseString("").
						HasDefaultNamespaceString("").
						HasDefaultRoleString("").
						HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionDefault).
						HasMinsToBypassMfaString(r.IntDefaultString).
						HasRsaPublicKeyString("").
						HasRsaPublicKey2String("").
						HasCommentString("").
						HasDisableMfaString(r.BooleanDefault).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.UserShowOutput(t, userModelNoAttributes.ResourceReference()).
						HasLoginName(strings.ToUpper(id.Name())).
						HasDisplayName(""),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2481 has been fixed
// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2353 has been fixed
// done on user, to not interfere with other parallel tests on the same account
func TestAcc_User_RemovedOutsideOfTerraform(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("u", userId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				PreConfig: testClient().User.DropUserFunc(t, userId),
				Config:    config.FromModels(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					func(state *terraform.State) error {
						if len(state.RootModule().Resources) != 1 {
							return errors.New("user should be created again and present in the state")
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAcc_User_issue2058 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2058 issue.
// The problem was with a dot in user identifier.
// Before the fix it results in panic: interface conversion: sdk.ObjectIdentifier is sdk.DatabaseObjectIdentifier, not sdk.AccountObjectIdentifier error.
func TestAcc_User_issue2058(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifierContaining(".")

	userModel1 := model.User("w", userId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModel1),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel1.ResourceReference()).HasNameString(userId.Name()),
				),
			},
		},
	})
}

func TestAcc_User_AllParameters(t *testing.T) {
	networkPolicy, networkPolicyCleanup := testClient().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
	t.Cleanup(networkPolicyCleanup)

	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("u", userId.Name())
	userModelWithAllParametersSet := model.User("u", userId.Name()).
		WithAbortDetachedQuery(true).
		WithAutocommit(false).
		WithBinaryInputFormatEnum(sdk.BinaryInputFormatUTF8).
		WithBinaryOutputFormatEnum(sdk.BinaryOutputFormatBase64).
		WithClientMemoryLimit(1024).
		WithClientMetadataRequestUseConnectionCtx(true).
		WithClientPrefetchThreads(2).
		WithClientResultChunkSize(48).
		WithClientResultColumnCaseInsensitive(true).
		WithClientSessionKeepAlive(true).
		WithClientSessionKeepAliveHeartbeatFrequency(2400).
		WithClientTimestampTypeMappingEnum(sdk.ClientTimestampTypeMappingNtz).
		WithDateInputFormat("YYYY-MM-DD").
		WithDateOutputFormat("YY-MM-DD").
		WithEnableUnloadPhysicalTypeOptimization(false).
		WithErrorOnNondeterministicMerge(false).
		WithErrorOnNondeterministicUpdate(true).
		WithGeographyOutputFormatEnum(sdk.GeographyOutputFormatWKB).
		WithGeometryOutputFormatEnum(sdk.GeometryOutputFormatWKB).
		WithJdbcTreatDecimalAsInt(false).
		WithJdbcTreatTimestampNtzAsUtc(true).
		WithJdbcUseSessionTimezone(false).
		WithJsonIndent(4).
		WithLockTimeout(21222).
		WithLogLevelEnum(sdk.LogLevelError).
		WithMultiStatementCount(0).
		WithNoorderSequenceAsDefault(false).
		WithOdbcTreatDecimalAsInt(true).
		WithQueryTag("some_tag").
		WithQuotedIdentifiersIgnoreCase(true).
		WithRowsPerResultset(2).
		WithS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
		WithSearchPath("$public, $current").
		WithSimulatedDataSharingConsumer("some_consumer").
		WithStatementQueuedTimeoutInSeconds(10).
		WithStatementTimeoutInSeconds(10).
		WithStrictJsonOutput(true).
		WithTimestampDayIsAlways24h(true).
		WithTimestampInputFormat("YYYY-MM-DD").
		WithTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampTypeMappingEnum(sdk.TimestampTypeMappingLtz).
		WithTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimezone("Europe/Warsaw").
		WithTimeInputFormat("HH24:MI").
		WithTimeOutputFormat("HH24:MI").
		WithTraceLevelEnum(sdk.TraceLevelPropagate).
		WithTransactionAbortOnError(true).
		WithTransactionDefaultIsolationLevelEnum(sdk.TransactionDefaultIsolationLevelReadCommitted).
		WithTwoDigitCenturyStart(1980).
		WithUnsupportedDdlActionEnum(sdk.UnsupportedDDLActionFail).
		WithUseCachedResult(false).
		WithWeekOfYearPolicy(1).
		WithWeekStart(1).
		WithEnableUnredactedQuerySyntaxError(true).
		WithNetworkPolicyId(networkPolicy.ID()).
		WithPreventUnloadToInternalStages(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// create with default values for all the parameters
			{
				Config: config.FromModels(t, userModel),
				Check: assertThat(t,
					objectparametersassert.UserParameters(t, userId).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.UserResourceParameters(t, userModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			// import when no parameter set
			{
				ResourceName: userModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceparametersassert.ImportedUserResourceParameters(t, userId.Name()).
						HasAllDefaults(),
				),
			},
			// set all parameters
			{
				Config: config.FromModels(t, userModelWithAllParametersSet),
				Check: assertThat(t,
					objectparametersassert.UserParameters(t, userId).
						HasAbortDetachedQuery(true).
						HasAutocommit(false).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcTreatDecimalAsInt(false).
						HasJdbcTreatTimestampNtzAsUtc(true).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$public, $current").
						HasSimulatedDataSharingConsumer("some_consumer").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelPropagate).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1).
						HasEnableUnredactedQuerySyntaxError(true).
						HasNetworkPolicy(networkPolicy.ID().Name()).
						HasPreventUnloadToInternalStages(true),
					resourceparametersassert.UserResourceParameters(t, "snowflake_user.u").
						HasAbortDetachedQuery(true).
						HasAutocommit(false).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcTreatDecimalAsInt(false).
						HasJdbcTreatTimestampNtzAsUtc(true).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$public, $current").
						HasSimulatedDataSharingConsumer("some_consumer").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelPropagate).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1).
						HasEnableUnredactedQuerySyntaxError(true).
						HasNetworkPolicy(networkPolicy.ID().Name()).
						HasPreventUnloadToInternalStages(true),
				),
			},
			// import when all parameters set
			{
				ResourceName: userModelWithAllParametersSet.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceparametersassert.ImportedUserResourceParameters(t, userId.Name()).
						HasAbortDetachedQuery(true).
						HasAutocommit(false).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcTreatDecimalAsInt(false).
						HasJdbcTreatTimestampNtzAsUtc(true).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$public, $current").
						HasSimulatedDataSharingConsumer("some_consumer").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelPropagate).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1).
						HasEnableUnredactedQuerySyntaxError(true).
						HasNetworkPolicy(networkPolicy.ID().Name()).
						HasPreventUnloadToInternalStages(true),
				),
			},
			// unset all the parameters
			{
				Config: config.FromModels(t, userModel),
				Check: assertThat(t,
					objectparametersassert.UserParameters(t, userId).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.UserResourceParameters(t, userModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
		},
	})
}

func TestAcc_User_issue2836(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	defaultRole := "SOME ROLE WITH SPACE case sensitive"
	defaultRoleQuoted := fmt.Sprintf(`"%s"`, defaultRole)

	userModel := model.User("u", userId.Name()).
		WithDefaultRole(defaultRoleQuoted)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModel),
				Check: assertThat(t,
					objectassert.User(t, userId).
						HasDefaultRole(defaultRole),
				),
			},
		},
	})
}

func TestAcc_User_issue2970(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	pass := random.Password()
	key, _ := random.GenerateRSAPublicKey(t)
	resourceName := "u"

	newPass := random.Password()
	newKey, _ := random.GenerateRSAPublicKey(t)
	incorrectlyFormattedNewKey := fmt.Sprintf("-invalid----BEGIN PUBLIC KEY-----\n%s-----END PUBLIC KEY-----\n", newKey)

	userModel := model.User(resourceName, userId.Name()).
		WithPassword(pass).
		WithRsaPublicKey(key)

	newUserModelIncorrectNewKey := model.User(resourceName, userId.Name()).
		WithPassword(newPass).
		WithRsaPublicKey(incorrectlyFormattedNewKey)

	newUserModel := model.User(resourceName, userId.Name()).
		WithPassword(newPass).
		WithRsaPublicKey(newKey)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModel),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasPasswordString(pass).
						HasRsaPublicKeyString(key),
				),
			},
			{
				Config: config.FromModels(t, newUserModelIncorrectNewKey),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(newUserModelIncorrectNewKey.ResourceReference(), "password", tfjson.ActionUpdate, sdk.String(pass), sdk.String(newPass)),
						planchecks.ExpectChange(newUserModelIncorrectNewKey.ResourceReference(), "rsa_public_key", tfjson.ActionUpdate, sdk.String(key), sdk.String(incorrectlyFormattedNewKey)),
					},
				},
				ExpectError: regexp.MustCompile("New public key rejected by current policy"),
			},
			{
				Config: config.FromModels(t, newUserModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(newUserModel.ResourceReference(), "password", tfjson.ActionUpdate, sdk.String(pass), sdk.String(newPass)),
						planchecks.ExpectChange(newUserModel.ResourceReference(), "rsa_public_key", tfjson.ActionUpdate, sdk.String(key), sdk.String(newKey)),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, newUserModel.ResourceReference()).
						HasPasswordString(newPass).
						HasRsaPublicKeyString(newKey),
				),
			},
		},
	})
}

func TestAcc_User_issue1572(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.UserWithDefaultMeta(userId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModel),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasDisabledString(r.BooleanDefault),
					objectassert.User(t, userId).HasDisabled(false),
				),
			},
			{
				PreConfig: func() {
					testClient().User.Disable(t, userId)
					objectassert.User(t, userId).HasDisabled(true)
				},
				Config: config.FromModels(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectDrift(userModel.ResourceReference(), "disabled", sdk.String(r.BooleanDefault), sdk.String(r.BooleanTrue)),
						planchecks.ExpectChange(userModel.ResourceReference(), "disabled", tfjson.ActionUpdate, sdk.String(r.BooleanTrue), sdk.String(r.BooleanDefault)),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasDisabledString(r.BooleanDefault),
					objectassert.User(t, userId).HasDisabled(false),
				),
			},
		},
	})
}

func TestAcc_User_issue1535_withNullPassword(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	pass := random.Password()

	userModel := model.UserWithDefaultMeta(userId.Name()).
		WithPassword(pass)

	userWithNullPasswordModel := model.UserWithDefaultMeta(userId.Name()).
		WithNullPassword()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModel),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasPasswordString(pass),
				),
			},
			{
				Config: config.FromModels(t, userWithNullPasswordModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userWithNullPasswordModel.ResourceReference(), "password", tfjson.ActionUpdate, sdk.String(pass), nil),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userWithNullPasswordModel.ResourceReference()).
						HasEmptyPassword(),
				),
			},
		},
	})
}

func TestAcc_User_issue1535_withRemovedPassword(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	pass := random.Password()

	userModel := model.UserWithDefaultMeta(userId.Name()).
		WithPassword(pass)

	userWithoutPasswordModel := model.UserWithDefaultMeta(userId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModel),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasPasswordString(pass),
				),
			},
			{
				Config: config.FromModels(t, userWithoutPasswordModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userWithoutPasswordModel.ResourceReference(), "password", tfjson.ActionUpdate, sdk.String(pass), nil),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userWithoutPasswordModel.ResourceReference()).
						HasEmptyPassword(),
				),
			},
		},
	})
}

func TestAcc_User_issue1155_handleChangesToDaysToExpiry(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModelWithoutDaysToExpiry := model.UserWithDefaultMeta(userId.Name())
	userModelDaysToExpiry10 := model.UserWithDefaultMeta(userId.Name()).WithDaysToExpiry(10)
	userModelDaysToExpiry5 := model.UserWithDefaultMeta(userId.Name()).WithDaysToExpiry(5)
	userModelDaysToExpiry0 := model.UserWithDefaultMeta(userId.Name()).WithDaysToExpiry(0)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// 1. create without days_to_expiry
			{
				Config: config.FromModels(t, userModelWithoutDaysToExpiry),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelWithoutDaysToExpiry.ResourceReference()).HasNoDaysToExpiry(),
					objectassert.User(t, userId).HasDaysToExpiryEmpty(),
				),
			},
			// 2. change to 10 (no plan after)
			{
				Config: config.FromModels(t, userModelDaysToExpiry10),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelDaysToExpiry10.ResourceReference()).HasDaysToExpiryString("10"),
					objectassert.User(t, userId).HasDaysToExpiryNotEmpty(),
				),
			},
			// 3. change externally to 2 (no changes)
			{
				PreConfig: func() {
					testClient().User.SetDaysToExpiry(t, userId, 2)
				},
				Config: config.FromModels(t, userModelDaysToExpiry10),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// 4. change externally to 0 (no changes)
			{
				PreConfig: func() {
					testClient().User.SetDaysToExpiry(t, userId, 0)
				},
				Config: config.FromModels(t, userModelDaysToExpiry10),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// 5. change in config to 5 (change)
			{
				Config: config.FromModels(t, userModelDaysToExpiry5),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userModelDaysToExpiry5.ResourceReference(), "days_to_expiry", tfjson.ActionUpdate, sdk.String("10"), sdk.String("5")),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelDaysToExpiry10.ResourceReference()).HasDaysToExpiryString("5"),
					objectassert.User(t, userId).HasDaysToExpiryNotEmpty(),
				),
			},
			// 6. change in config to 0 (change)
			{
				Config: config.FromModels(t, userModelDaysToExpiry0),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userModelDaysToExpiry0.ResourceReference(), "days_to_expiry", tfjson.ActionUpdate, sdk.String("5"), sdk.String("0")),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelDaysToExpiry10.ResourceReference()).HasDaysToExpiryString("0"),
					objectassert.User(t, userId).HasDaysToExpiryEmpty(),
				),
			},
			// 7. remove from config (no change)
			{
				Config: config.FromModels(t, userModelWithoutDaysToExpiry),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelWithoutDaysToExpiry.ResourceReference()).HasDaysToExpiryString("0"),
					objectassert.User(t, userId).HasDaysToExpiryEmpty(),
				),
			},
		},
	})
}

func TestAcc_User_handleExternalTypeChange(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.UserWithDefaultMeta(userId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModel),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).HasNameString(userId.Name()).HasUserTypeString(""),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).HasType(""),
				),
			},
			{
				PreConfig: func() {
					testClient().User.SetType(t, userId, sdk.UserTypeService)
					objectassert.User(t, userId).HasType(string(sdk.UserTypeService))
				},
				Config: config.FromModels(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).HasNameString(userId.Name()).HasUserTypeString(""),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).HasType(""),
				),
			},
			// no change should happen if the change is to PERSON explicitly
			{
				PreConfig: func() {
					testClient().User.SetType(t, userId, sdk.UserTypePerson)
					objectassert.User(t, userId).HasType("PERSON")
				},
				Config: config.FromModels(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).HasNameString(userId.Name()).HasUserTypeString("PERSON"),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).HasType("PERSON"),
				),
			},
			// no change should happen if we fall back to default
			{
				PreConfig: func() {
					testClient().User.UnsetType(t, userId)
					objectassert.User(t, userId).HasType("")
				},
				Config: config.FromModels(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).HasNameString(userId.Name()).HasUserTypeString(""),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).HasType(""),
				),
			},
		},
	})
}

func TestAcc_User_handleChangesToDefaultSecondaryRoles(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModelEmpty := model.UserWithDefaultMeta(userId.Name())
	userModelWithOptionAll := model.UserWithDefaultMeta(userId.Name()).WithDefaultSecondaryRolesOptionEnum(sdk.SecondaryRolesOptionAll)
	userModelWithOptionNone := model.UserWithDefaultMeta(userId.Name()).WithDefaultSecondaryRolesOptionEnum(sdk.SecondaryRolesOptionNone)
	userModelLowercaseValue := model.UserWithDefaultMeta(userId.Name()).WithDefaultSecondaryRolesOption("all")
	userModelIncorrectValue := model.UserWithDefaultMeta(userId.Name()).WithDefaultSecondaryRolesOption("OTHER")
	userModelEmptyValue := model.UserWithDefaultMeta(userId.Name()).WithDefaultSecondaryRolesOption("")
	userModelNullValue := model.UserWithDefaultMeta(userId.Name()).WithDefaultSecondaryRolesOptionValue(config.ReplacementPlaceholderVariable(config.SnowflakeProviderConfigNull))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// 1. create without default secondary roles option set (DEFAULT will be used)
			{
				Config: config.FromModels(t, userModelEmpty),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelEmpty.ResourceReference()).HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionDefault),
					objectassert.User(t, userId).HasDefaultSecondaryRoles(`["ALL"]`),
				),
			},
			// 2. add default secondary roles ALL (expecting no change)
			{
				Config: config.FromModels(t, userModelWithOptionAll),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// 3. change to lowercase (change because we have DEFAULT in state because previous step was suppressed so none of the suppressors NormalizeAndCompare nor IgnoreChangeToCurrentSnowflakeValueInShowWithMapping suppresses it; it can be made better later)
			{
				Config: config.FromModels(t, userModelLowercaseValue),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userModelWithOptionNone.ResourceReference(), "default_secondary_roles_option", tfjson.ActionUpdate, sdk.String("DEFAULT"), sdk.String("all")),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelEmpty.ResourceReference()).HasDefaultSecondaryRolesOptionString("all"),
					objectassert.User(t, userId).HasDefaultSecondaryRoles(`["ALL"]`),
				),
			},
			// 4. add default secondary roles NONE
			{
				Config: config.FromModels(t, userModelWithOptionNone),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelWithOptionNone.ResourceReference()).HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionNone),
					objectassert.User(t, userId).HasDefaultSecondaryRoles(`[]`),
				),
			},
			// 5. unset externally
			{
				PreConfig: func() {
					testClient().User.UnsetDefaultSecondaryRoles(t, userId)
				},
				Config: config.FromModels(t, userModelWithOptionNone),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userModelWithOptionNone.ResourceReference(), "default_secondary_roles_option", tfjson.ActionUpdate, sdk.String("ALL"), sdk.String("NONE")),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelWithOptionAll.ResourceReference()).HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionNone),
					objectassert.User(t, userId).HasDefaultSecondaryRoles(`[]`),
				),
			},
			// 6. unset in config (change)
			{
				Config: config.FromModels(t, userModelEmpty),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userModelEmpty.ResourceReference(), "default_secondary_roles_option", tfjson.ActionUpdate, sdk.String("NONE"), sdk.String("DEFAULT")),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelEmpty.ResourceReference()).HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionDefault),
					objectassert.User(t, userId).HasDefaultSecondaryRoles(`["ALL"]`),
				),
			},
			// 7. incorrect value used
			{
				Config:      config.FromModels(t, userModelIncorrectValue),
				ExpectError: regexp.MustCompile("invalid secondary roles option: OTHER"),
			},
			// 8. set to empty in config (invalid)
			{
				Config:      config.FromModels(t, userModelEmptyValue),
				ExpectError: regexp.MustCompile("invalid secondary roles option: "),
			},
			// 9. set in config to NONE (change)
			{
				Config: config.FromModels(t, userModelWithOptionNone),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userModelWithOptionNone.ResourceReference(), "default_secondary_roles_option", tfjson.ActionUpdate, sdk.String("DEFAULT"), sdk.String("NONE")),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelWithOptionNone.ResourceReference()).HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionNone),
					objectassert.User(t, userId).HasDefaultSecondaryRoles("[]"),
				),
			},
			// 10. unset in config (change)
			{
				Config: config.FromModels(t, userModelEmpty),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userModelEmpty.ResourceReference(), "default_secondary_roles_option", tfjson.ActionUpdate, sdk.String("NONE"), sdk.String("DEFAULT")),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelEmpty.ResourceReference()).HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionDefault),
					objectassert.User(t, userId).HasDefaultSecondaryRoles(`["ALL"]`),
				),
			},
			// 11. add default secondary roles NONE
			{
				Config: config.FromModels(t, userModelWithOptionNone),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelEmpty.ResourceReference()).HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionNone),
					objectassert.User(t, userId).HasDefaultSecondaryRoles(`[]`),
				),
			},
			// 12. set to null value in config (change)
			{
				Config: config.FromModels(t, userModelNullValue),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userModelNullValue.ResourceReference(), "default_secondary_roles_option", tfjson.ActionUpdate, sdk.String("NONE"), sdk.String("DEFAULT")),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelEmpty.ResourceReference()).HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionDefault),
					objectassert.User(t, userId).HasDefaultSecondaryRoles(`["ALL"]`),
				),
			},
		},
	})
}

func TestAcc_User_migrateFromVersion094_defaultSecondaryRolesSet(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	userModelWithOptionAll := model.UserWithDefaultMeta(id.Name()).WithDefaultSecondaryRolesOptionEnum(sdk.SecondaryRolesOptionAll)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config: fmt.Sprintf(`
resource "snowflake_user" "test" {
	name = "%s"
	default_secondary_roles = ["ALL"]
}`, id.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(userModelWithOptionAll.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(userModelWithOptionAll.ResourceReference(), "default_secondary_roles.#", "1"),
					resource.TestCheckResourceAttr(userModelWithOptionAll.ResourceReference(), "default_secondary_roles.0", "ALL"),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, userModelWithOptionAll),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// we do not have a plancheck yet to validate no changes on the given field; this is a current alternative
						planchecks.ExpectChange(userModelWithOptionAll.ResourceReference(), "default_secondary_roles", tfjson.ActionUpdate, nil, nil),
						planchecks.ExpectChange(userModelWithOptionAll.ResourceReference(), "default_secondary_roles_option", tfjson.ActionUpdate, sdk.String(string(sdk.SecondaryRolesOptionAll)), sdk.String(string(sdk.SecondaryRolesOptionAll))),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(userModelWithOptionAll.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(userModelWithOptionAll.ResourceReference(), "default_secondary_roles_option", string(sdk.SecondaryRolesOptionAll)),
					resource.TestCheckNoResourceAttr(userModelWithOptionAll.ResourceReference(), "default_secondary_roles"),
				),
			},
		},
	})
}

func TestAcc_User_ParameterValidationsAndDiffSuppressions(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("w", id.Name()).
		WithBinaryInputFormat(strings.ToLower(string(sdk.BinaryInputFormatHex))).
		WithBinaryOutputFormat(strings.ToLower(string(sdk.BinaryOutputFormatHex))).
		WithGeographyOutputFormat(strings.ToLower(string(sdk.GeographyOutputFormatGeoJSON))).
		WithGeometryOutputFormat(strings.ToLower(string(sdk.GeometryOutputFormatGeoJSON))).
		WithLogLevel(strings.ToLower(string(sdk.LogLevelInfo))).
		WithTimestampTypeMapping(strings.ToLower(string(sdk.TimestampTypeMappingNtz))).
		WithTraceLevel(strings.ToLower(string(sdk.TraceLevelAlways)))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModel),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasBinaryInputFormatString(string(sdk.BinaryInputFormatHex)).
						HasBinaryOutputFormatString(string(sdk.BinaryOutputFormatHex)).
						HasGeographyOutputFormatString(string(sdk.GeographyOutputFormatGeoJSON)).
						HasGeometryOutputFormatString(string(sdk.GeometryOutputFormatGeoJSON)).
						HasLogLevelString(string(sdk.LogLevelInfo)).
						HasTimestampTypeMappingString(string(sdk.TimestampTypeMappingNtz)).
						HasTraceLevelString(string(sdk.TraceLevelAlways)),
				),
			},
		},
	})
}

func TestAcc_User_LoginNameAndDisplayName(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()

	loginName := random.SensitiveAlphanumeric()
	userModelWithoutBoth := model.User("w", id.Name())
	userModelWithNewId := model.User("w", newId.Name())
	userModelWithBoth := model.User("w", newId.Name()).WithLoginName(loginName).WithDisplayName("display_name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// Create without both set
			{
				Config: config.FromModels(t, userModelWithoutBoth),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelWithoutBoth.ResourceReference()).
						HasNoDisplayName().
						HasNoLoginName(),
					objectassert.User(t, id).
						HasDisplayName(strings.ToUpper(id.Name())).
						HasLoginName(strings.ToUpper(id.Name())),
				),
			},
			// Rename
			{
				Config: config.FromModels(t, userModelWithNewId),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelWithNewId.ResourceReference()).
						HasNoDisplayName().
						HasNoLoginName(),
					objectassert.User(t, newId).
						HasDisplayName(strings.ToUpper(id.Name())).
						HasLoginName(strings.ToUpper(id.Name())),
				),
			},
			// Set both params
			{
				Config: config.FromModels(t, userModelWithBoth),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelWithBoth.ResourceReference()).
						HasDisplayNameString("display_name").
						HasLoginNameString(loginName),
					objectassert.User(t, newId).
						HasDisplayName("display_name").
						HasLoginName(strings.ToUpper(loginName)),
				),
			},
			// Unset externally
			{
				PreConfig: func() {
					testClient().User.Alter(t, newId, &sdk.AlterUserOptions{
						Unset: &sdk.UserUnset{
							ObjectProperties: &sdk.UserObjectPropertiesUnset{
								LoginName:   sdk.Bool(true),
								DisplayName: sdk.Bool(true),
							},
						},
					})
				},
				Config: config.FromModels(t, userModelWithBoth),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelWithBoth.ResourceReference()).
						HasDisplayNameString("display_name").
						HasLoginNameString(loginName),
					objectassert.User(t, newId).
						HasDisplayName("display_name").
						HasLoginName(strings.ToUpper(loginName)),
				),
			},
			// Unset both params
			{
				Config: config.FromModels(t, userModelWithNewId),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelWithNewId.ResourceReference()).
						HasDisplayNameString("").
						HasLoginNameString(""),
					objectassert.User(t, newId).
						HasDisplayName("").
						HasLoginName(strings.ToUpper(newId.Name())),
				),
			},
			// Set externally
			{
				PreConfig: func() {
					testClient().User.Alter(t, newId, &sdk.AlterUserOptions{
						Set: &sdk.UserSet{
							ObjectProperties: &sdk.UserAlterObjectProperties{
								UserObjectProperties: sdk.UserObjectProperties{
									LoginName:   sdk.String("external_" + loginName),
									DisplayName: sdk.String("external_display_name"),
								},
							},
						},
					})
				},
				Config: config.FromModels(t, userModelWithNewId),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelWithNewId.ResourceReference()).
						HasDisplayNameString("").
						HasLoginNameString(""),
					objectassert.User(t, newId).
						HasDisplayName("").
						HasLoginName(strings.ToUpper(newId.Name())),
				),
			},
		},
	})
}

// https://docs.snowflake.com/en/release-notes/bcr-bundles/2024_08/bcr-1798
// https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3125
func TestAcc_User_handleChangesToShowUsers_bcr202408_gh3125(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModelNoAttributes := model.User("w", userId.Name())
	userModelWithNoneDefaultSecondaryRoles := model.User("w", userId.Name()).WithDefaultSecondaryRolesOptionEnum(sdk.SecondaryRolesOptionNone)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModelNoAttributes),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelNoAttributes.ResourceReference()).
						HasAllDefaults(userId, sdk.SecondaryRolesOptionDefault),
				),
			},
			{
				Config: config.FromModels(t, userModelWithNoneDefaultSecondaryRoles),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModelWithNoneDefaultSecondaryRoles.ResourceReference()).
						HasAllDefaults(userId, sdk.SecondaryRolesOptionNone),
				),
			},
		},
	})
}

// https://docs.snowflake.com/en/release-notes/bcr-bundles/2024_08/bcr-1798
// https://docs.snowflake.com/release-notes/bcr-bundles/2024_08/bcr-1692
// https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3125
func TestAcc_User_handleChangesToShowUsers_bcr202408_generallyEnabled(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("w", userId.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.98.0"),
				Config:            config.FromModels(t, userModel),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasAllDefaults(userId, sdk.SecondaryRolesOptionDefault),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasAllDefaults(userId, sdk.SecondaryRolesOptionDefault),
				),
			},
		},
	})
}

// https://docs.snowflake.com/en/release-notes/bcr-bundles/2024_08/bcr-1798
func TestAcc_User_handleChangesToShowUsers_bcr202408_defaults(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("w", userId.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.97.0"),
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				Config:            config.FromModels(t, userModel),
				ExpectError:       regexp.MustCompile("\"default_namespace\": converting NULL to string is unsupported"),
			},
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.98.0"),
				Config:            config.FromModels(t, userModel),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasNoDefaultNamespace(),
				),
			},
		},
	})
}

func TestAcc_User_importPassword(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	pass := random.Password()
	firstName := random.AlphaN(6)

	_, userCleanup := testClient().User.CreateUserWithOptions(t, userId, &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
		Password:  sdk.String(pass),
		FirstName: sdk.String(firstName),
	}})
	t.Cleanup(userCleanup)

	userModel := model.User("w", userId.Name()).WithPassword(pass).WithFirstName(firstName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// IMPORT
			{
				Config:        config.FromModels(t, userModel),
				ResourceName:  userModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: userId.Name(),
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedUserResource(t, userId.Name()).
						HasNoPassword().
						HasFirstNameString(firstName),
				),
				ImportStatePersist: true,
			},
			{
				Config: config.FromModels(t, userModel),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasNotEmptyPassword().
						HasFirstNameString(firstName),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: config.FromModels(t, userModel),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasNotEmptyPassword().
						HasFirstNameString(firstName),
				),
			},
		},
	})
}

// https://github.com/snowflakedb/terraform-provider-snowflake/issues/3522
// The issue arises when the following steps happen:
// 1. Object is created.
// 2. Object's attribute is updated to computed value of other object added in this run.
// diff.HasChange(changedKey) used in ComputedIfAnyAttributeChanged does not trigger, so the whole show_output is not marked as new attribute.
// The original set up would require additional external provider, so to simplify, we use computed value from other object.
// TestAcc_User_gh3522_proof test shows the previous behavior, TestAcc_User_gh3522_fix confirms the fix.
func TestAcc_User_gh3522_proof(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	userId2 := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("1.0.5"),
				Config:            gh3522ConfigFirstStep(userId),
				Check: assertThat(t, resourceassert.UserResource(t, "snowflake_legacy_service_user.one").
					HasNoComment(),
				),
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("1.0.5"),
				Config:            gh3522ConfigSecondStep(userId, userId2, comment),
				// Resulting in:
				// | Error: Provider produced inconsistent final plan
				// |
				// | When expanding the plan for snowflake_legacy_service_user.one to include new
				// | values learned so far during apply, provider
				// | "registry.terraform.io/hashicorp/snowflake" produced an invalid new value for
				// | .show_output: was known, but now unknown.
				// |
				// | This is a bug in the provider, which should be reported in the provider's own
				// | issue tracker.
				ExpectError: regexp.MustCompile("Provider produced inconsistent final plan"),
			},
		},
	})
}

// Check TestAcc_User_gh3522_proof for details.
func TestAcc_User_gh3522_fix(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	userId2 := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: gh3522ConfigFirstStep(userId),
				Check: assertThat(t, resourceassert.UserResource(t, "snowflake_legacy_service_user.one").
					HasNoComment(),
				),
			},
			{
				Config: gh3522ConfigSecondStep(userId, userId2, comment),
				Check: assertThat(t, resourceassert.UserResource(t, "snowflake_legacy_service_user.one").
					HasCommentString(comment),
				),
			},
		},
	})
}

func gh3522ConfigFirstStep(userId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_legacy_service_user" "one" {
  name           = "%[1]s"
}
`, userId.Name())
}

func gh3522ConfigSecondStep(userId sdk.AccountObjectIdentifier, userId2 sdk.AccountObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_legacy_service_user" "one" {
  name    = "%[1]s"
  comment = snowflake_legacy_service_user.two.show_output.0.comment
}

resource "snowflake_legacy_service_user" "two" {
  name    = "%[2]s"
  comment = "%[3]s"
}
`, userId.Name(), userId2.Name(), comment)
}

// Result of https://github.com/snowflakedb/terraform-provider-snowflake/issues/3655
func TestAcc_User_gh3655(t *testing.T) {
	networkPolicyId := testClient().Ids.RandomAccountObjectIdentifier()
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	networkPolicyModel := model.NetworkPolicy("test", networkPolicyId.Name())
	userModel := model.User("test", userId.Name()).
		WithNetworkPolicyValue(config.UnquotedWrapperVariable("snowflake_network_policy.test.fully_qualified_name"))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.0.0"),
				Config:            config.FromModels(t, networkPolicyModel, userModel),
				Check: assertThat(t, resourceassert.UserResource(t, userModel.ResourceReference()).
					HasNetworkPolicyString(networkPolicyId.Name()),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, networkPolicyModel, userModel),
				Check: assertThat(t, resourceassert.UserResource(t, userModel.ResourceReference()).
					HasNetworkPolicyString(networkPolicyId.Name()),
				),
			},
		},
	})
}

func TestAcc_User_BCR_2025_05(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("test", userId.Name())

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					secondaryTestClient().BcrBundles.DisableBcrBundle(t, "2025_05")
				},
				Config: config.FromModels(t, providerModel, userModel),
				Check: assertThat(t, resourceassert.UserResource(t, userModel.ResourceReference()).
					HasNameString(userId.Name()).
					HasUserTypeEmpty(),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).
						HasTypeEmpty(),
				),
			},
			// Prove that after enabling the 2025_05 bundle, the type is still not set, and the plan is empty.
			{
				PreConfig: func() {
					secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_05")
				},
				Config: config.FromModels(t, providerModel, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t, resourceassert.UserResource(t, userModel.ResourceReference()).
					HasNameString(userId.Name()).
					HasUserTypeEmpty(),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).
						HasTypeEmpty(),
				),
			},
			// When the user type is set to PERSON externally, the plan is empty.
			{
				PreConfig: func() {
					secondaryTestClient().User.SetType(t, userId, sdk.UserTypePerson)
				},
				Config: config.FromModels(t, providerModel, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t, resourceassert.UserResource(t, userModel.ResourceReference()).
					HasNameString(userId.Name()).
					HasUserTypeString(string(sdk.UserTypePerson)),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).
						HasType(string(sdk.UserTypePerson)),
				),
			},
		},
	})
}

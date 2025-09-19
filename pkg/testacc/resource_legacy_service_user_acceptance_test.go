//go:build account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_LegacyServiceUser_BasicFlows(t *testing.T) {
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

	userModelNoAttributes := model.LegacyServiceUser("w", id.Name())
	userModelNoAttributesRenamed := model.LegacyServiceUser("w", id2.Name()).
		WithComment(newComment)

	userModelAllAttributes := model.LegacyServiceUser("w", id.Name()).
		WithPassword(pass).
		WithLoginName(loginName).
		WithDisplayName("Display Name").
		WithEmail("fake@email.com").
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

	userModelAllAttributesChanged := func(loginName string) *model.LegacyServiceUserModel {
		return model.LegacyServiceUser("w", id.Name()).
			WithPassword(newPass).
			WithLoginName(loginName).
			WithDisplayName("New Display Name").
			WithEmail("fake@email.net").
			WithMustChangePassword("false").
			WithDisabled("true").
			WithDaysToExpiry(12).
			WithMinsToUnlock(13).
			WithDefaultWarehouse("other_warehouse").
			WithDefaultNamespace("one_part_namespace").
			WithDefaultRole("other_role").
			WithDefaultSecondaryRolesOptionEnum(sdk.SecondaryRolesOptionAll).
			WithRsaPublicKey(key2).
			WithRsaPublicKey2(key1).
			WithComment(newComment)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.LegacyServiceUser),
		Steps: []resource.TestStep{
			// CREATE WITHOUT ATTRIBUTES
			{
				Config: config.FromModels(t, userModelNoAttributes),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasNoPassword().
						HasNoLoginName().
						HasNoDisplayName().
						HasNoEmail().
						HasMustChangePasswordString(r.BooleanDefault).
						HasDisabledString(r.BooleanDefault).
						HasNoDaysToExpiry().
						HasMinsToUnlockString(r.IntDefaultString).
						HasNoDefaultWarehouse().
						HasNoDefaultNamespace().
						HasNoDefaultRole().
						HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionDefault).
						HasNoRsaPublicKey().
						HasNoRsaPublicKey2().
						HasNoComment().
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
					resourceassert.LegacyServiceUserResource(t, userModelNoAttributes.ResourceReference()).
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
				ImportStateVerifyIgnore: []string{"password", "days_to_expiry", "mins_to_unlock", "login_name", "display_name", "disabled", "must_change_password", "default_secondary_roles_option"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedLegacyServiceUserResource(t, id2.Name()).
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
					resourceassert.LegacyServiceUserResource(t, userModelAllAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasPasswordString(pass).
						HasLoginNameString(loginName).
						HasDisplayNameString("Display Name").
						HasEmailString("fake@email.com").
						HasMustChangePassword(true).
						HasDisabled(false).
						HasDaysToExpiryString("8").
						HasMinsToUnlockString("9").
						HasDefaultWarehouseString("some_warehouse").
						HasDefaultNamespaceString("some.namespace").
						HasDefaultRoleString("some_role").
						HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionAll).
						HasRsaPublicKeyString(key1).
						HasRsaPublicKey2String(key2).
						HasCommentString(comment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: config.FromModels(t, userModelAllAttributesChanged(newLoginName)),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelAllAttributesChanged(newLoginName).ResourceReference()).
						HasNameString(id.Name()).
						HasPasswordString(newPass).
						HasLoginNameString(newLoginName).
						HasDisplayNameString("New Display Name").
						HasEmailString("fake@email.net").
						HasMustChangePassword(false).
						HasDisabled(true).
						HasDaysToExpiryString("12").
						HasMinsToUnlockString("13").
						HasDefaultWarehouseString("other_warehouse").
						HasDefaultNamespaceString("one_part_namespace").
						HasDefaultRoleString("other_role").
						HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionAll).
						HasRsaPublicKeyString(key2).
						HasRsaPublicKey2String(key1).
						HasCommentString(newComment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelAllAttributesChanged(newLoginName).ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "days_to_expiry", "mins_to_unlock", "default_namespace", "login_name", "show_output.0.days_to_expiry"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedLegacyServiceUserResource(t, id.Name()).
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
					resourceassert.LegacyServiceUserResource(t, userModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasPasswordString("").
						HasLoginNameString("").
						HasDisplayNameString("").
						HasEmailString("").
						HasMustChangePasswordString(r.BooleanDefault).
						HasDisabledString(r.BooleanDefault).
						HasDaysToExpiryString("0").
						HasMinsToUnlockString(r.IntDefaultString).
						HasDefaultWarehouseString("").
						HasDefaultNamespaceString("").
						HasDefaultRoleString("").
						HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionDefault).
						HasRsaPublicKeyString("").
						HasRsaPublicKey2String("").
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.UserShowOutput(t, userModelNoAttributes.ResourceReference()).
						HasLoginName(strings.ToUpper(id.Name())).
						HasDisplayName(""),
				),
			},
		},
	})
}

func TestAcc_LegacyServiceUser_AllParameters(t *testing.T) {
	networkPolicy, networkPolicyCleanup := testClient().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
	t.Cleanup(networkPolicyCleanup)

	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.LegacyServiceUser("u", userId.Name())
	userModelWithAllParametersSet := model.LegacyServiceUser("u", userId.Name()).
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
		CheckDestroy: CheckDestroy(t, resources.LegacyServiceUser),
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
					resourceparametersassert.UserResourceParameters(t, userModelWithAllParametersSet.ResourceReference()).
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

func TestAcc_LegacyServiceUser_handleExternalTypeChange(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.LegacyServiceUserWithDefaultMeta(userId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.LegacyServiceUser),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModel),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).HasNameString(userId.Name()).HasUserTypeString("LEGACY_SERVICE"),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).HasType("LEGACY_SERVICE"),
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
					resourceassert.UserResource(t, userModel.ResourceReference()).HasNameString(userId.Name()).HasUserTypeString("LEGACY_SERVICE"),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).HasType("LEGACY_SERVICE"),
				),
			},
			{
				PreConfig: func() {
					testClient().User.UnsetType(t, userId)
					objectassert.User(t, userId).HasType("")
				},
				Config: config.FromModels(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).HasNameString(userId.Name()).HasUserTypeString("LEGACY_SERVICE"),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).HasType("LEGACY_SERVICE"),
				),
			},
		},
	})
}

func TestAcc_LegacyServiceUser_setIncompatibleAttributes(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.LegacyServiceUser),
		Steps: []resource.TestStep{
			{
				Config:      legacyServiceUserConfigWithIncompatibleAttribute(userId, "first_name", random.AlphaN(6)),
				ExpectError: regexp.MustCompile("An argument named \"first_name\" is not expected here"),
			},
			{
				Config:      legacyServiceUserConfigWithIncompatibleAttribute(userId, "middle_name", random.AlphaN(6)),
				ExpectError: regexp.MustCompile("An argument named \"middle_name\" is not expected here"),
			},
			{
				Config:      legacyServiceUserConfigWithIncompatibleAttribute(userId, "last_name", random.AlphaN(6)),
				ExpectError: regexp.MustCompile("An argument named \"last_name\" is not expected here"),
			},
			{
				Config:      legacyServiceUserConfigWithIncompatibleAttribute(userId, "mins_to_bypass_mfa", random.AlphaN(6)),
				ExpectError: regexp.MustCompile("An argument named \"mins_to_bypass_mfa\" is not expected here"),
			},
			{
				Config:      legacyServiceUserConfigWithIncompatibleAttribute(userId, "disable_mfa", random.AlphaN(6)),
				ExpectError: regexp.MustCompile("An argument named \"disable_mfa\" is not expected here"),
			},
		},
	})
}

func legacyServiceUserConfigWithIncompatibleAttribute(userId sdk.AccountObjectIdentifier, key string, value string) string {
	return fmt.Sprintf(`
        resource "snowflake_legacy_service_user" "test" {
        	name = %s
			%s = "%s"
        }
	`, userId.FullyQualifiedName(), key, value)
}

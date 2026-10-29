//go:build non_account_level_tests

package testint

import (
	"slices"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// quotaWebhookUrl is never called: admin summaries are sent only when users cross their limits.
const quotaWebhookUrl = "https://hooks.slack.com/services/T00000000/B00000000/AAAAAAAAAAAAAAAAAAAAAAAA"

func TestInt_Quotas(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	sharedResourceNamed := func(t *testing.T, scope sdk.QuotaScope, domain sdk.QuotaSharedResourceDomain) []sdk.QuotaScopeSharedResource {
		t.Helper()
		return slices.DeleteFunc(slices.Clone(scope.SharedResources), func(r sdk.QuotaScopeSharedResource) bool {
			normalized, err := r.NormalizedDomain()
			return err != nil || normalized != domain
		})
	}
	userNames := func(users []sdk.QuotaScopeUser) []string {
		names := make([]string, len(users))
		for i, user := range users {
			names[i] = user.UserName
		}
		return names
	}
	describe := func(t *testing.T, id sdk.SchemaObjectIdentifier) *sdk.QuotaDetails {
		t.Helper()
		details, err := client.Quotas.DescribeDetails(ctx, id)
		require.NoError(t, err)
		return details
	}
	// enabling enforcement on a quota that only applies to a new user can't block anyone else
	quotaScopedToNewUser := func(t *testing.T) (sdk.SchemaObjectIdentifier, *sdk.User) {
		t.Helper()
		quotaId, quotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(quotaCleanup)
		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)
		testClientHelper().Quota.IncludeUsers(t, quotaId, user.ID().Name())
		return quotaId, user
	}
	grantProcedureToSnowflakeApplication := func(t *testing.T, procedureId sdk.SchemaObjectIdentifierWithArguments) {
		t.Helper()
		t.Cleanup(testClientHelper().Grant.GrantUsageOnDatabaseToSnowflakeApplication(t, procedureId.DatabaseId()))
		t.Cleanup(testClientHelper().Grant.GrantUsageOnSchemaToSnowflakeApplication(t, procedureId.SchemaId()))
		t.Cleanup(testClientHelper().Grant.GrantUsageOnProcedureToSnowflakeApplication(t, procedureId))
	}

	t.Run("create, show and drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Quotas.Create(ctx, sdk.NewCreateQuotaRequest(id).WithComment("quota comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Quota.DropFunc(t, id))

		assertThatObject(t, objectassert.Quota(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment("quota comment").
			HasOwnerNotEmpty().
			HasCreatedOnNotEmpty())

		err = client.Quotas.Drop(ctx, sdk.NewDropQuotaRequest(id))
		require.NoError(t, err)

		_, err = client.Quotas.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)

		// dropping a missing quota fails; DropSafely drops it with IF EXISTS
		err = client.Quotas.Drop(ctx, sdk.NewDropQuotaRequest(id))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
		require.NoError(t, client.Quotas.DropSafely(ctx, id))
	})

	t.Run("alter: rename and comment", func(t *testing.T) {
		quotaId, quotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(quotaCleanup)
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		t.Cleanup(testClientHelper().Quota.DropFunc(t, newId))

		err := client.Quotas.Alter(ctx, sdk.NewAlterQuotaRequest(quotaId).WithSet(*sdk.NewQuotaSetRequest().WithComment("changed")))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Quota(t, quotaId).HasComment("changed"))

		err = client.Quotas.Alter(ctx, sdk.NewAlterQuotaRequest(quotaId).WithUnset(*sdk.NewQuotaUnsetRequest().WithComment(true)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Quota(t, quotaId).HasNoComment())

		err = client.Quotas.Alter(ctx, sdk.NewAlterQuotaRequest(quotaId).WithRenameTo(newId))
		require.NoError(t, err)
		_, err = client.Quotas.ShowByID(ctx, quotaId)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
		assertThatObject(t, objectassert.Quota(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()))
	})

	t.Run("show with like and in", func(t *testing.T) {
		quotaId, quotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(quotaCleanup)
		otherQuotaId, otherQuotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(otherQuotaCleanup)

		quotas, err := client.Quotas.Show(ctx, sdk.NewShowQuotaRequest().
			WithLike(sdk.Like{Pattern: new(quotaId.Name())}).
			WithIn(sdk.In{Schema: quotaId.SchemaId()}))
		require.NoError(t, err)
		require.Len(t, quotas, 1)
		assertThatObject(t, objectassert.QuotaFromObject(t, &quotas[0]).
			HasName(quotaId.Name()).
			HasDatabaseName(quotaId.DatabaseName()).
			HasSchemaName(quotaId.SchemaName()))

		quotas, err = client.Quotas.Show(ctx, sdk.NewShowQuotaRequest().WithIn(sdk.In{Schema: quotaId.SchemaId()}))
		require.NoError(t, err)
		ids := make([]sdk.SchemaObjectIdentifier, len(quotas))
		for i, quota := range quotas {
			ids[i] = quota.ID()
		}
		assert.Contains(t, ids, quotaId)
		assert.Contains(t, ids, otherQuotaId)
	})

	t.Run("describe details of a new quota", func(t *testing.T) {
		quotaId, quotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(quotaCleanup)

		details := describe(t, quotaId)
		assert.Equal(t, quotaId, details.Id)
		assert.Empty(t, details.Scope.SharedResources)
		assert.Empty(t, details.Scope.IncludedUsers)
		assert.Empty(t, details.Scope.ExcludedUsers)
		assert.Nil(t, details.Config.PerUserLimit)
		assert.Empty(t, details.NotificationThresholds)
		assert.Empty(t, details.NotificationIntegrations)
		assert.Empty(t, details.CustomActions)
		assert.Nil(t, details.CycleStartAction)
	})

	t.Run("AI domains as shared resources", func(t *testing.T) {
		quotaId, quotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(quotaCleanup)

		for _, domain := range []sdk.QuotaSharedResourceDomain{
			sdk.QuotaSharedResourceDomainCortexCode,
			sdk.QuotaSharedResourceDomainSnowflakeIntelligence,
		} {
			testClientHelper().Quota.AddSharedResource(t, quotaId, domain, "")
		}
		testClientHelper().Quota.AddSharedResource(t, quotaId, sdk.QuotaSharedResourceDomainAiFunction, "AI_CLASSIFY")

		scope := describe(t, quotaId).Scope
		require.Len(t, sharedResourceNamed(t, scope, sdk.QuotaSharedResourceDomainCortexCode), 1)
		assert.True(t, sharedResourceNamed(t, scope, sdk.QuotaSharedResourceDomainCortexCode)[0].IsWholeDomain())
		require.Len(t, sharedResourceNamed(t, scope, sdk.QuotaSharedResourceDomainSnowflakeIntelligence), 1)
		require.Len(t, sharedResourceNamed(t, scope, sdk.QuotaSharedResourceDomainAiFunction), 1)
		aiFunction := sharedResourceNamed(t, scope, sdk.QuotaSharedResourceDomainAiFunction)[0]
		assert.False(t, aiFunction.IsWholeDomain())
		assert.Equal(t, "AI_CLASSIFY", aiFunction.Name)

		args, err := sdk.NewQuotaRemoveSharedResourceArgsRequestFromInstance(sdk.QuotaSharedResourceDomainAiFunction, "AI_CLASSIFY")
		require.NoError(t, err)
		_, err = client.Quotas.RemoveSharedResource(ctx, sdk.NewRemoveSharedResourceQuotaRequest(quotaId, *args))
		require.NoError(t, err)

		scope = describe(t, quotaId).Scope
		assert.Empty(t, sharedResourceNamed(t, scope, sdk.QuotaSharedResourceDomainAiFunction))
		assert.Len(t, scope.SharedResources, 2)
	})

	t.Run("Cortex Agent as a shared resource", func(t *testing.T) {
		quotaId, quotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(quotaCleanup)
		agentId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		t.Cleanup(testClientHelper().CortexAgent.CreateWithId(t, agentId))

		testClientHelper().Quota.AddSharedResource(t, quotaId, sdk.QuotaSharedResourceDomainCortexAgent, agentId.FullyQualifiedName())

		agents := sharedResourceNamed(t, describe(t, quotaId).Scope, sdk.QuotaSharedResourceDomainCortexAgent)
		require.Len(t, agents, 1)
		assert.False(t, agents[0].IsWholeDomain())
		// GET_QUOTA_SCOPE returns the agent's name only; DescribeDetails resolves it
		assert.Equal(t, agentId.Name(), agents[0].Name)
		require.NotNil(t, agents[0].CortexAgentId)
		assert.Equal(t, agentId, *agents[0].CortexAgentId)

		args, err := sdk.NewQuotaRemoveSharedResourceArgsRequestFromInstance(sdk.QuotaSharedResourceDomainCortexAgent, agentId.FullyQualifiedName())
		require.NoError(t, err)
		_, err = client.Quotas.RemoveSharedResource(ctx, sdk.NewRemoveSharedResourceQuotaRequest(quotaId, *args))
		require.NoError(t, err)
		assert.Empty(t, describe(t, quotaId).Scope.SharedResources)
	})

	t.Run("warehouse as a shared resource", func(t *testing.T) {
		quotaId, quotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(quotaCleanup)
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		testClientHelper().Quota.AddSharedResource(t, quotaId, sdk.QuotaSharedResourceDomainWarehouse, warehouse.ID().FullyQualifiedName())

		warehouses := sharedResourceNamed(t, describe(t, quotaId).Scope, sdk.QuotaSharedResourceDomainWarehouse)
		require.Len(t, warehouses, 1)
		assert.Equal(t, warehouse.ID().Name(), warehouses[0].Name)

		args, err := sdk.NewQuotaRemoveSharedResourceArgsRequestFromInstance(sdk.QuotaSharedResourceDomainWarehouse, warehouse.ID().FullyQualifiedName())
		require.NoError(t, err)
		_, err = client.Quotas.RemoveSharedResource(ctx, sdk.NewRemoveSharedResourceQuotaRequest(quotaId, *args))
		require.NoError(t, err)

		assert.Empty(t, describe(t, quotaId).Scope.SharedResources)
	})

	// the instance method calls and SYSTEM$REFERENCE escape double quotes inside names
	t.Run("names with double quotes", func(t *testing.T) {
		quotaId, quotaCleanup := testClientHelper().Quota.CreateWithId(t, testClientHelper().Ids.RandomSchemaObjectIdentifierContaining(`"`))
		t.Cleanup(quotaCleanup)
		warehouseId := testClientHelper().Ids.RandomAccountObjectIdentifierContaining(`"`)
		_, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithRequest(t, sdk.NewCreateWarehouseRequest(warehouseId).WithInitiallySuspended(true))
		t.Cleanup(warehouseCleanup)

		args := sdk.NewQuotaAddSharedResourceArgsRequest(sdk.QuotaSharedResourceDomainWarehouse).WithWarehouse(warehouseId)
		_, err := client.Quotas.AddSharedResource(ctx, sdk.NewAddSharedResourceQuotaRequest(quotaId, *args))
		require.NoError(t, err)

		warehouses := sharedResourceNamed(t, describe(t, quotaId).Scope, sdk.QuotaSharedResourceDomainWarehouse)
		require.Len(t, warehouses, 1)
		assert.Equal(t, warehouseId.Name(), warehouses[0].Name)
	})

	t.Run("user scope", func(t *testing.T) {
		quotaId, quotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(quotaCleanup)
		includedUser, includedUserCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(includedUserCleanup)
		excludedUser, excludedUserCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(excludedUserCleanup)
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		includeUsers := func(users ...string) {
			t.Helper()
			_, err := client.Quotas.IncludeUsers(ctx, sdk.NewIncludeUsersQuotaRequest(quotaId, *sdk.NewQuotaIncludeUsersArgsRequest(users)))
			require.NoError(t, err)
		}
		excludeUsers := func(users ...string) {
			t.Helper()
			_, err := client.Quotas.ExcludeUsers(ctx, sdk.NewExcludeUsersQuotaRequest(quotaId, *sdk.NewQuotaExcludeUsersArgsRequest().
				WithByUser(*sdk.NewQuotaExcludeUsersByUserRequest(users))))
			require.NoError(t, err)
		}

		// names are resolved case-insensitively
		includeUsers(strings.ToLower(includedUser.ID().Name()))
		excludeUsers(excludedUser.ID().Name())
		_, err := client.Quotas.SetUserTags(ctx, sdk.NewSetUserTagsQuotaRequest(quotaId, *sdk.NewQuotaSetUserTagsArgsRequest(
			[]sdk.QuotaUserTagRequest{*sdk.NewQuotaUserTagRequest(tag.ID(), "finance"), *sdk.NewQuotaUserTagRequest(tag.ID(), "Analytics")},
			sdk.QuotaUserTagsOperatorIntersection,
		)))
		require.NoError(t, err)

		scope := describe(t, quotaId).Scope
		assert.Equal(t, []string{includedUser.ID().Name()}, userNames(scope.IncludedUsers))
		assert.Equal(t, []string{excludedUser.ID().Name()}, userNames(scope.ExcludedUsers))
		require.NotNil(t, scope.UserTags)
		assert.Equal(t, string(sdk.QuotaUserTagsOperatorIntersection), scope.UserTags.Operator)
		assert.False(t, scope.UserTags.AreExclusions())
		require.Len(t, scope.UserTags.Tags, 1)
		assert.Equal(t, tag.ID(), scope.UserTags.Tags[0].TagId())
		assert.ElementsMatch(t, []string{"finance", "Analytics"}, scope.UserTags.Tags[0].TagValues)

		// excluding tags replaces the user tags and clears the included and excluded users
		_, err = client.Quotas.ExcludeUsers(ctx, sdk.NewExcludeUsersQuotaRequest(quotaId, *sdk.NewQuotaExcludeUsersArgsRequest().
			WithByTag(*sdk.NewQuotaExcludeUsersByTagRequest([]sdk.QuotaExcludedUserTagRequest{*sdk.NewQuotaExcludedUserTagRequest(tag.ID(), "contractor")}))))
		require.NoError(t, err)

		scope = describe(t, quotaId).Scope
		assert.Empty(t, scope.IncludedUsers)
		assert.Empty(t, scope.ExcludedUsers)
		require.NotNil(t, scope.UserTags)
		assert.True(t, scope.UserTags.AreExclusions())
		require.Len(t, scope.UserTags.Tags, 1)
		assert.Equal(t, []string{"contractor"}, scope.UserTags.Tags[0].TagValues)

		// empty arrays clear the user lists, and excluding no tags resets the quota to all users
		includeUsers(includedUser.ID().Name())
		includeUsers()
		excludeUsers()
		_, err = client.Quotas.ExcludeUsers(ctx, sdk.NewExcludeUsersQuotaRequest(quotaId, *sdk.NewQuotaExcludeUsersArgsRequest().
			WithByTag(*sdk.NewQuotaExcludeUsersByTagRequest(nil))))
		require.NoError(t, err)

		scope = describe(t, quotaId).Scope
		assert.Empty(t, scope.IncludedUsers)
		assert.Empty(t, scope.ExcludedUsers)
		require.NotNil(t, scope.UserTags)
		assert.True(t, scope.UserTags.AreExclusions())
		assert.Empty(t, scope.UserTags.Tags)
	})

	t.Run("per-user limits and block enforcement", func(t *testing.T) {
		quotaId, _ := quotaScopedToNewUser(t)
		testClientHelper().Quota.AddSharedResource(t, quotaId, sdk.QuotaSharedResourceDomainAiFunction, "")

		config := describe(t, quotaId).Config
		assert.Nil(t, config.PerUserLimit)
		assert.False(t, config.BlockEnforcementEnabled)
		assert.True(t, config.PerUserBlockNotificationsEnabled)
		assert.Empty(t, config.AdminEmails)

		testClientHelper().Quota.SetPerUserLimit(t, quotaId, 1000, sdk.QuotaCycleMonthly)
		testClientHelper().Quota.SetPerUserLimit(t, quotaId, 500, sdk.QuotaCycleWeekly)
		testClientHelper().Quota.SetPerUserLimit(t, quotaId, 120, sdk.QuotaCycleDaily)
		_, err := client.Quotas.SetBlockEnforcementEnabled(ctx, sdk.NewSetBlockEnforcementEnabledQuotaRequest(quotaId, *sdk.NewQuotaSetBlockEnforcementEnabledArgsRequest(true).WithNotifyUsers(false)))
		require.NoError(t, err)

		config = describe(t, quotaId).Config
		assert.Equal(t, new(1000), config.PerUserLimit)
		assert.Equal(t, new(500), config.PerUserLimitWeekly)
		assert.Equal(t, new(120), config.PerUserLimitDaily)
		assert.True(t, config.BlockEnforcementEnabled)
		assert.False(t, config.PerUserBlockNotificationsEnabled)

		// limits must be positive
		_, err = client.Quotas.SetPerUserLimit(ctx, sdk.NewSetPerUserLimitQuotaRequest(quotaId, *sdk.NewQuotaSetPerUserLimitArgsRequest(0)))
		require.ErrorContains(t, err, "Invalid per-user limit")

		for _, cycle := range []sdk.QuotaCycle{sdk.QuotaCycleMonthly, sdk.QuotaCycleWeekly} {
			_, err = client.Quotas.UnsetPerUserLimit(ctx, sdk.NewUnsetPerUserLimitQuotaRequest(quotaId, *sdk.NewQuotaUnsetPerUserLimitArgsRequest(cycle)))
			require.NoError(t, err)
		}
		_, err = client.Quotas.SetBlockEnforcementEnabled(ctx, sdk.NewSetBlockEnforcementEnabledQuotaRequest(quotaId, *sdk.NewQuotaSetBlockEnforcementEnabledArgsRequest(false).WithNotifyUsers(true)))
		require.NoError(t, err)

		config = describe(t, quotaId).Config
		assert.Nil(t, config.PerUserLimit)
		assert.Nil(t, config.PerUserLimitWeekly)
		assert.Equal(t, new(120), config.PerUserLimitDaily)
		assert.False(t, config.BlockEnforcementEnabled)
		assert.True(t, config.PerUserBlockNotificationsEnabled)
	})

	t.Run("notification thresholds", func(t *testing.T) {
		quotaId, quotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(quotaCleanup)

		addThreshold := func(threshold int, strategy sdk.QuotaSpendStrategy, notifyUser bool, cycle sdk.QuotaCycle) {
			t.Helper()
			_, err := client.Quotas.AddNotificationThreshold(ctx, sdk.NewAddNotificationThresholdQuotaRequest(quotaId,
				*sdk.NewQuotaAddNotificationThresholdArgsRequest(threshold, strategy, notifyUser, cycle)))
			require.NoError(t, err)
		}
		// thresholds are identified by (threshold, spend strategy, cycle)
		addThreshold(80, sdk.QuotaSpendStrategyProjected, true, sdk.QuotaCycleMonthly)
		addThreshold(80, sdk.QuotaSpendStrategyProjected, false, sdk.QuotaCycleDaily)
		addThreshold(1000, sdk.QuotaSpendStrategyActual, false, sdk.QuotaCycleWeekly)

		thresholds, err := client.Quotas.GetNotificationThresholds(ctx, sdk.NewGetNotificationThresholdsQuotaRequest(quotaId))
		require.NoError(t, err)
		require.Len(t, thresholds, 3)
		assert.True(t, slices.ContainsFunc(thresholds, func(threshold sdk.QuotaNotificationThreshold) bool {
			return threshold.Threshold == 80 && threshold.SpendStrategy == sdk.QuotaSpendStrategyProjected && !threshold.NotifyUser && threshold.Cycle == sdk.QuotaCycleDaily
		}))

		_, err = client.Quotas.RemoveNotificationThreshold(ctx, sdk.NewRemoveNotificationThresholdQuotaRequest(quotaId,
			*sdk.NewQuotaRemoveNotificationThresholdArgsRequest(80, sdk.QuotaSpendStrategyProjected).WithCycle(sdk.QuotaCycleDaily)))
		require.NoError(t, err)

		thresholds, err = client.Quotas.GetNotificationThresholds(ctx, sdk.NewGetNotificationThresholdsQuotaRequest(quotaId))
		require.NoError(t, err)
		require.Len(t, thresholds, 2)
		assert.False(t, slices.ContainsFunc(thresholds, func(threshold sdk.QuotaNotificationThreshold) bool {
			return threshold.Cycle == sdk.QuotaCycleDaily
		}))

		// without the cycle, the monthly threshold is removed
		_, err = client.Quotas.RemoveNotificationThreshold(ctx, sdk.NewRemoveNotificationThresholdQuotaRequest(quotaId,
			*sdk.NewQuotaRemoveNotificationThresholdArgsRequest(80, sdk.QuotaSpendStrategyProjected)))
		require.NoError(t, err)

		thresholds, err = client.Quotas.GetNotificationThresholds(ctx, sdk.NewGetNotificationThresholdsQuotaRequest(quotaId))
		require.NoError(t, err)
		require.Len(t, thresholds, 1)
		assert.Equal(t, 1000, thresholds[0].Threshold)
		assert.Equal(t, sdk.QuotaCycleWeekly, thresholds[0].Cycle)
	})

	// Quotas accept webhook and queue integrations (not email ones), and verified admin emails only; helpers.VerifiedEmail is
	// verified in the accounts the tests run in.
	t.Run("notification integrations and admin emails", func(t *testing.T) {
		quotaId, quotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(quotaCleanup)
		integration, integrationCleanup := testClientHelper().NotificationIntegration.CreateWebhook(t, quotaWebhookUrl)
		t.Cleanup(integrationCleanup)
		t.Cleanup(testClientHelper().Grant.GrantUsageOnIntegrationToSnowflakeApplication(t, integration.ID()))

		_, err := client.Quotas.AddNotificationIntegration(ctx, sdk.NewAddNotificationIntegrationQuotaRequest(quotaId, *sdk.NewQuotaAddNotificationIntegrationArgsRequest(integration.ID())))
		require.NoError(t, err)
		_, err = client.Quotas.SetAdminEmails(ctx, sdk.NewSetAdminEmailsQuotaRequest(quotaId, *sdk.NewQuotaSetAdminEmailsArgsRequestFromEmails(helpers.VerifiedEmail)))
		require.NoError(t, err)

		integrations, err := client.Quotas.GetNotificationIntegrations(ctx, sdk.NewGetNotificationIntegrationsQuotaRequest(quotaId))
		require.NoError(t, err)
		require.Len(t, integrations, 1)
		assert.Equal(t, integration.ID().Name(), integrations[0].IntegrationName)

		config, err := client.Quotas.GetConfig(ctx, sdk.NewGetConfigQuotaRequest(quotaId))
		require.NoError(t, err)
		assert.Equal(t, []string{helpers.VerifiedEmail}, config.AdminEmails)

		_, err = client.Quotas.RemoveNotificationIntegration(ctx, sdk.NewRemoveNotificationIntegrationQuotaRequest(quotaId, *sdk.NewQuotaRemoveNotificationIntegrationArgsRequest(integration.ID())))
		require.NoError(t, err)
		// an empty list clears the admin emails
		_, err = client.Quotas.SetAdminEmails(ctx, sdk.NewSetAdminEmailsQuotaRequest(quotaId, *sdk.NewQuotaSetAdminEmailsArgsRequestFromEmails()))
		require.NoError(t, err)

		integrations, err = client.Quotas.GetNotificationIntegrations(ctx, sdk.NewGetNotificationIntegrationsQuotaRequest(quotaId))
		require.NoError(t, err)
		assert.Empty(t, integrations)
		config, err = client.Quotas.GetConfig(ctx, sdk.NewGetConfigQuotaRequest(quotaId))
		require.NoError(t, err)
		assert.Empty(t, config.AdminEmails)
	})

	t.Run("custom actions", func(t *testing.T) {
		quotaId, quotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(quotaCleanup)
		// the first argument receives the IDs of the users who crossed the threshold; extra arguments fill the remaining ones
		procedure, procedureCleanup := testClientHelper().Procedure.Create(t, sdk.DataTypeVARCHAR)
		t.Cleanup(procedureCleanup)
		grantProcedureToSnowflakeApplication(t, procedure.ID())
		procedureWithArgument, procedureWithArgumentCleanup := testClientHelper().Procedure.Create(t, sdk.DataTypeVARCHAR, sdk.DataTypeVARCHAR)
		t.Cleanup(procedureWithArgumentCleanup)
		t.Cleanup(testClientHelper().Grant.GrantUsageOnProcedureToSnowflakeApplication(t, procedureWithArgument.ID()))

		addAction := func(procedureId sdk.SchemaObjectIdentifierWithArguments, arguments []string, strategy sdk.QuotaSpendStrategy, threshold int) {
			t.Helper()
			_, err := client.Quotas.AddCustomAction(ctx, sdk.NewAddCustomActionQuotaRequest(quotaId,
				*sdk.NewQuotaAddCustomActionArgsRequest(procedureId, arguments, strategy, threshold)))
			require.NoError(t, err)
		}
		addAction(procedure.ID(), nil, sdk.QuotaSpendStrategyActual, 100)
		addAction(procedure.ID(), nil, sdk.QuotaSpendStrategyProjected, 100)
		addAction(procedureWithArgument.ID(), []string{"it's"}, sdk.QuotaSpendStrategyProjected, 90)

		// the number of extra arguments must match the procedure's remaining parameters
		_, err := client.Quotas.AddCustomAction(ctx, sdk.NewAddCustomActionQuotaRequest(quotaId,
			*sdk.NewQuotaAddCustomActionArgsRequest(procedure.ID(), []string{"extra"}, sdk.QuotaSpendStrategyActual, 50)))
		require.ErrorContains(t, err, "ARGUMENTS_SIZE_MISMATCH")

		actions, err := client.Quotas.GetCustomActions(ctx, sdk.NewGetCustomActionsQuotaRequest(quotaId))
		require.NoError(t, err)
		require.Len(t, actions, 3)
		for _, action := range actions {
			procedureId, err := action.ProcedureId()
			require.NoError(t, err)
			assert.Contains(t, []string{procedure.ID().FullyQualifiedName(), procedureWithArgument.ID().FullyQualifiedName()}, procedureId.FullyQualifiedName())
			assert.NotEmpty(t, action.ActionId)
		}

		// removing by threshold and procedure removes the actions of both spend strategies;
		// the procedure is passed as the PROCEDURE_FQN returned by GET_CUSTOM_ACTIONS
		idx := slices.IndexFunc(actions, func(action sdk.QuotaCustomAction) bool { return action.Threshold == 100 })
		require.GreaterOrEqual(t, idx, 0)
		_, err = client.Quotas.RemoveCustomActions(ctx, sdk.NewRemoveCustomActionsQuotaRequest(quotaId,
			*sdk.NewQuotaRemoveCustomActionsArgsRequest().WithThreshold(100).WithProcedureFqn(actions[idx].ProcedureFqn)))
		require.NoError(t, err)

		actions, err = client.Quotas.GetCustomActions(ctx, sdk.NewGetCustomActionsQuotaRequest(quotaId))
		require.NoError(t, err)
		require.Len(t, actions, 1)
		assert.Equal(t, 90, actions[0].Threshold)
		assert.Equal(t, []string{"it's"}, actions[0].ProcedureArgs)

		_, err = client.Quotas.RemoveCustomActions(ctx, sdk.NewRemoveCustomActionsQuotaRequest(quotaId, *sdk.NewQuotaRemoveCustomActionsArgsRequest()))
		require.NoError(t, err)

		actions, err = client.Quotas.GetCustomActions(ctx, sdk.NewGetCustomActionsQuotaRequest(quotaId))
		require.NoError(t, err)
		assert.Empty(t, actions)
	})

	t.Run("cycle-start action", func(t *testing.T) {
		quotaId, quotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(quotaCleanup)
		procedure, procedureCleanup := testClientHelper().Procedure.Create(t)
		t.Cleanup(procedureCleanup)
		grantProcedureToSnowflakeApplication(t, procedure.ID())

		_, err := client.Quotas.SetCycleStartAction(ctx, sdk.NewSetCycleStartActionQuotaRequest(quotaId, *sdk.NewQuotaSetCycleStartActionArgsRequest(procedure.ID(), nil)))
		require.NoError(t, err)

		action := describe(t, quotaId).CycleStartAction
		require.NotNil(t, action)
		procedureId, err := action.ProcedureId()
		require.NoError(t, err)
		assert.Equal(t, procedure.ID().FullyQualifiedName(), procedureId.FullyQualifiedName())
		assert.Empty(t, action.ProcedureArgs)

		_, err = client.Quotas.RemoveCycleStartAction(ctx, sdk.NewRemoveCycleStartActionQuotaRequest(quotaId))
		require.NoError(t, err)

		assert.Nil(t, describe(t, quotaId).CycleStartAction)
	})
}

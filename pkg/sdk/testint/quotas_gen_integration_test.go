//go:build non_account_level_tests

package testint

import (
	"slices"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// quotaWebhookUrl is never called: admin summaries are sent only when users cross their limits.
const quotaWebhookUrl = "https://hooks.slack.com/services/T00000000/B00000000/AAAAAAAAAAAAAAAAAAAAAAAA"

func TestInt_Quotas(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	createQuota := func(t *testing.T) sdk.SchemaObjectIdentifier {
		t.Helper()
		id, cleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(cleanup)
		return id
	}
	describe := func(t *testing.T, id sdk.SchemaObjectIdentifier) *sdk.QuotaDetails {
		t.Helper()
		details, err := client.Quotas.DescribeDetails(ctx, id)
		require.NoError(t, err)
		return details
	}
	sharedResourcesIn := func(t *testing.T, id sdk.SchemaObjectIdentifier, domain sdk.QuotaSharedResourceDomain) []sdk.QuotaScopeSharedResource {
		t.Helper()
		return slices.DeleteFunc(describe(t, id).Scope.SharedResources, func(r sdk.QuotaScopeSharedResource) bool {
			normalized, err := r.NormalizedDomain()
			return err != nil || normalized != domain
		})
	}
	addSharedResource := func(t *testing.T, id sdk.SchemaObjectIdentifier, domain sdk.QuotaSharedResourceDomain, instance string) error {
		t.Helper()
		args, err := sdk.NewQuotaAddSharedResourceArgsRequestFromInstance(domain, instance)
		require.NoError(t, err)
		_, err = client.Quotas.AddSharedResource(ctx, sdk.NewAddSharedResourceQuotaRequest(id, *args))
		return err
	}
	removeSharedResource := func(t *testing.T, id sdk.SchemaObjectIdentifier, domain sdk.QuotaSharedResourceDomain, instance string) {
		t.Helper()
		args, err := sdk.NewQuotaRemoveSharedResourceArgsRequestFromInstance(domain, instance)
		require.NoError(t, err)
		_, err = client.Quotas.RemoveSharedResource(ctx, sdk.NewRemoveSharedResourceQuotaRequest(id, *args))
		require.NoError(t, err)
	}
	userNames := func(users []sdk.QuotaScopeUser) []string {
		return collections.Map(users, func(user sdk.QuotaScopeUser) string { return user.UserName })
	}
	tagValues := func(tags []sdk.QuotaScopeTag) []string {
		var values []string
		for _, tag := range tags {
			values = append(values, tag.TagValues...)
		}
		return values
	}
	// enabling enforcement on a quota that only applies to a new user can't block anyone else
	quotaScopedToNewUser := func(t *testing.T) sdk.SchemaObjectIdentifier {
		t.Helper()
		id := createQuota(t)
		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)
		testClientHelper().Quota.IncludeUsers(t, id, user.ID().Name())
		return id
	}
	grantProcedureToSnowflakeApplication := func(t *testing.T, procedureId sdk.SchemaObjectIdentifierWithArguments) {
		t.Helper()
		t.Cleanup(testClientHelper().Grant.GrantUsageOnDatabaseToSnowflakeApplication(t, procedureId.DatabaseId()))
		t.Cleanup(testClientHelper().Grant.GrantUsageOnSchemaToSnowflakeApplication(t, procedureId.SchemaId()))
		t.Cleanup(testClientHelper().Grant.GrantUsageOnProcedureToSnowflakeApplication(t, procedureId))
	}

	t.Run("create: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Quotas.Create(ctx, sdk.NewCreateQuotaRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Quota.DropFunc(t, id))

		assertThatObject(t, objectassert.Quota(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwnerNotEmpty().
			HasNoComment(),
		)
	})

	t.Run("create: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		err := client.Quotas.Create(ctx, sdk.NewCreateQuotaRequest(id).WithIfNotExists(true).WithComment(comment))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Quota.DropFunc(t, id))

		assertThatObject(t, objectassert.Quota(t, id).
			HasName(id.Name()).
			HasComment(comment),
		)

		// IF NOT EXISTS keeps the existing quota
		err = client.Quotas.Create(ctx, sdk.NewCreateQuotaRequest(id).WithIfNotExists(true).WithComment(random.Comment()))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Quota(t, id).
			HasComment(comment),
		)
	})

	t.Run("create: or replace", func(t *testing.T) {
		id := createQuota(t)
		comment := random.Comment()

		err := client.Quotas.Create(ctx, sdk.NewCreateQuotaRequest(id).WithOrReplace(true).WithComment(comment))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Quota(t, id).
			HasComment(comment),
		)
	})

	t.Run("show: like, in, and limit", func(t *testing.T) {
		id := createQuota(t)
		otherId := createQuota(t)

		shown, err := client.Quotas.Show(ctx, sdk.NewShowQuotaRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}).
			WithIn(sdk.In{Schema: id.SchemaId()}).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		require.Len(t, shown, 1)
		assertThatObject(t, objectassert.QuotaFromObject(t, &shown[0]).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()),
		)

		inDatabase, err := client.Quotas.Show(ctx, sdk.NewShowQuotaRequest().WithIn(sdk.In{Database: id.DatabaseId()}))
		require.NoError(t, err)
		ids := collections.Map(inDatabase, func(quota sdk.Quota) sdk.SchemaObjectIdentifier { return quota.ID() })
		assert.Contains(t, ids, id)
		assert.Contains(t, ids, otherId)

		inAccount, err := client.Quotas.Show(ctx, sdk.NewShowQuotaRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}).
			WithIn(sdk.In{Account: sdk.Bool(true)}))
		require.NoError(t, err)
		require.Len(t, inAccount, 1)
	})

	t.Run("describe: new quota", func(t *testing.T) {
		id := createQuota(t)

		details := describe(t, id)
		assert.Equal(t, id, details.Id)
		assert.Empty(t, details.Scope.SharedResources)
		assert.Empty(t, details.Scope.IncludedUsers)
		assert.Empty(t, details.Scope.ExcludedUsers)
		assert.Empty(t, details.Scope.SelectionTags())
		assert.Empty(t, details.Scope.ExclusionTags())
		assert.Nil(t, details.Config.PerUserLimit)
		assert.Nil(t, details.Config.PerUserLimitWeekly)
		assert.Nil(t, details.Config.PerUserLimitDaily)
		assert.False(t, details.Config.BlockEnforcementEnabled)
		assert.True(t, details.Config.PerUserBlockNotificationsEnabled)
		assert.Empty(t, details.Config.AdminEmails)
		assert.Empty(t, details.NotificationThresholds)
		assert.Empty(t, details.NotificationIntegrations)
		assert.Empty(t, details.CustomActions)
		assert.Nil(t, details.CycleStartAction)
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		_, err := client.Quotas.DescribeDetails(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter: set and unset comment", func(t *testing.T) {
		id := createQuota(t)
		comment := random.Comment()

		err := client.Quotas.Alter(ctx, sdk.NewAlterQuotaRequest(id).WithSet(*sdk.NewQuotaSetRequest().WithComment(comment)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Quota(t, id).
			HasComment(comment),
		)

		err = client.Quotas.Alter(ctx, sdk.NewAlterQuotaRequest(id).WithUnset(*sdk.NewQuotaUnsetRequest().WithComment(true)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Quota(t, id).
			HasNoComment(),
		)
	})

	t.Run("alter: rename", func(t *testing.T) {
		id := createQuota(t)
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		t.Cleanup(testClientHelper().Quota.DropFunc(t, newId))

		err := client.Quotas.Alter(ctx, sdk.NewAlterQuotaRequest(id).WithRenameTo(newId))
		require.NoError(t, err)

		_, err = client.Quotas.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
		assertThatObject(t, objectassert.Quota(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()),
		)
	})

	t.Run("alter: rename to another schema", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)
		id := createQuota(t)
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		t.Cleanup(testClientHelper().Quota.DropFunc(t, newId))

		err := client.Quotas.Alter(ctx, sdk.NewAlterQuotaRequest(id).WithRenameTo(newId))
		require.NoError(t, err)

		_, err = client.Quotas.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
		assertThatObject(t, objectassert.Quota(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()),
		)
	})

	t.Run("drop: basic", func(t *testing.T) {
		id := createQuota(t)

		err := client.Quotas.Drop(ctx, sdk.NewDropQuotaRequest(id))
		require.NoError(t, err)
		_, err = client.Quotas.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop: if exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Quotas.Drop(ctx, sdk.NewDropQuotaRequest(id))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

		err = client.Quotas.Drop(ctx, sdk.NewDropQuotaRequest(id).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("add and remove shared resources: AI domain", func(t *testing.T) {
		id := createQuota(t)

		require.NoError(t, addSharedResource(t, id, sdk.QuotaSharedResourceDomainCortexCode, ""))
		cortexCode := sharedResourcesIn(t, id, sdk.QuotaSharedResourceDomainCortexCode)
		require.Len(t, cortexCode, 1)
		assert.True(t, cortexCode[0].IsWholeDomain())

		require.NoError(t, addSharedResource(t, id, sdk.QuotaSharedResourceDomainAiFunction, "AI_CLASSIFY"))
		aiFunctions := sharedResourcesIn(t, id, sdk.QuotaSharedResourceDomainAiFunction)
		require.Len(t, aiFunctions, 1)
		assert.False(t, aiFunctions[0].IsWholeDomain())
		assert.Equal(t, "AI_CLASSIFY", aiFunctions[0].Name)

		removeSharedResource(t, id, sdk.QuotaSharedResourceDomainAiFunction, "AI_CLASSIFY")
		assert.Empty(t, sharedResourcesIn(t, id, sdk.QuotaSharedResourceDomainAiFunction))
		assert.Len(t, sharedResourcesIn(t, id, sdk.QuotaSharedResourceDomainCortexCode), 1)

		removeSharedResource(t, id, sdk.QuotaSharedResourceDomainCortexCode, "")
		assert.Empty(t, describe(t, id).Scope.SharedResources)
	})

	t.Run("add and remove shared resources: Cortex Agent", func(t *testing.T) {
		id := createQuota(t)
		agentId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		t.Cleanup(testClientHelper().CortexAgent.CreateWithId(t, agentId))

		require.NoError(t, addSharedResource(t, id, sdk.QuotaSharedResourceDomainCortexAgent, agentId.FullyQualifiedName()))
		agents := sharedResourcesIn(t, id, sdk.QuotaSharedResourceDomainCortexAgent)
		require.Len(t, agents, 1)
		assert.False(t, agents[0].IsWholeDomain())
		// GET_QUOTA_SCOPE returns the agent's name only; DescribeDetails resolves it
		assert.Equal(t, agentId.Name(), agents[0].Name)
		require.NotNil(t, agents[0].CortexAgentId)
		assert.Equal(t, agentId, *agents[0].CortexAgentId)

		removeSharedResource(t, id, sdk.QuotaSharedResourceDomainCortexAgent, agentId.FullyQualifiedName())
		assert.Empty(t, describe(t, id).Scope.SharedResources)
	})

	t.Run("add and remove shared resources: warehouse", func(t *testing.T) {
		id := createQuota(t)
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		require.NoError(t, addSharedResource(t, id, sdk.QuotaSharedResourceDomainWarehouse, warehouse.ID().FullyQualifiedName()))
		warehouses := sharedResourcesIn(t, id, sdk.QuotaSharedResourceDomainWarehouse)
		require.Len(t, warehouses, 1)
		assert.Equal(t, warehouse.ID().Name(), warehouses[0].Name)

		removeSharedResource(t, id, sdk.QuotaSharedResourceDomainWarehouse, warehouse.ID().FullyQualifiedName())
		assert.Empty(t, describe(t, id).Scope.SharedResources)
	})

	// ADD_SHARED_RESOURCE takes one resource per call, but a quota holds several
	t.Run("add and remove shared resources: several at once", func(t *testing.T) {
		id := createQuota(t)

		require.NoError(t, addSharedResource(t, id, sdk.QuotaSharedResourceDomainAiFunction, "AI_CLASSIFY"))
		require.NoError(t, addSharedResource(t, id, sdk.QuotaSharedResourceDomainAiFunction, "AI_COMPLETE"))
		require.NoError(t, addSharedResource(t, id, sdk.QuotaSharedResourceDomainCortexCode, ""))
		require.NoError(t, addSharedResource(t, id, sdk.QuotaSharedResourceDomainSnowflakeIntelligence, ""))

		scope := describe(t, id).Scope
		require.Len(t, scope.SharedResources, 4)
		assert.ElementsMatch(t, []string{"AI_CLASSIFY", "AI_COMPLETE"}, collections.Map(sharedResourcesIn(t, id, sdk.QuotaSharedResourceDomainAiFunction), func(r sdk.QuotaScopeSharedResource) string { return r.Name }))
		assert.Len(t, sharedResourcesIn(t, id, sdk.QuotaSharedResourceDomainCortexCode), 1)
		assert.Len(t, sharedResourcesIn(t, id, sdk.QuotaSharedResourceDomainSnowflakeIntelligence), 1)

		removeSharedResource(t, id, sdk.QuotaSharedResourceDomainAiFunction, "AI_COMPLETE")
		assert.Len(t, describe(t, id).Scope.SharedResources, 3)
	})

	t.Run("add shared resources: Snowflake rejects mixed credit types", func(t *testing.T) {
		id := createQuota(t)
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		require.NoError(t, addSharedResource(t, id, sdk.QuotaSharedResourceDomainAiFunction, ""))
		err := addSharedResource(t, id, sdk.QuotaSharedResourceDomainWarehouse, warehouse.ID().FullyQualifiedName())
		require.ErrorContains(t, err, "different credit types")
	})

	t.Run("add shared resources: Snowflake rejects an instance of a whole domain", func(t *testing.T) {
		id := createQuota(t)

		require.NoError(t, addSharedResource(t, id, sdk.QuotaSharedResourceDomainAiFunction, ""))
		err := addSharedResource(t, id, sdk.QuotaSharedResourceDomainAiFunction, "AI_CLASSIFY")
		require.ErrorContains(t, err, "Cannot add individual resources")
	})

	// the instance method calls and SYSTEM$REFERENCE escape double quotes inside names
	t.Run("add and remove shared resources: names with double quotes", func(t *testing.T) {
		id, cleanup := testClientHelper().Quota.CreateWithId(t, testClientHelper().Ids.RandomSchemaObjectIdentifierContaining(`"`))
		t.Cleanup(cleanup)
		warehouseId := testClientHelper().Ids.RandomAccountObjectIdentifierContaining(`"`)
		_, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithRequest(t, sdk.NewCreateWarehouseRequest(warehouseId).WithInitiallySuspended(true))
		t.Cleanup(warehouseCleanup)

		args := sdk.NewQuotaAddSharedResourceArgsRequest(sdk.QuotaSharedResourceDomainWarehouse).WithWarehouse(warehouseId)
		_, err := client.Quotas.AddSharedResource(ctx, sdk.NewAddSharedResourceQuotaRequest(id, *args))
		require.NoError(t, err)
		warehouses := sharedResourcesIn(t, id, sdk.QuotaSharedResourceDomainWarehouse)
		require.Len(t, warehouses, 1)
		assert.Equal(t, warehouseId.Name(), warehouses[0].Name)

		removeArgs := sdk.NewQuotaRemoveSharedResourceArgsRequest(sdk.QuotaSharedResourceDomainWarehouse).WithWarehouse(warehouseId)
		_, err = client.Quotas.RemoveSharedResource(ctx, sdk.NewRemoveSharedResourceQuotaRequest(id, *removeArgs))
		require.NoError(t, err)
		assert.Empty(t, describe(t, id).Scope.SharedResources)
	})

	t.Run("include users", func(t *testing.T) {
		id := createQuota(t)
		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)
		otherUser, otherUserCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(otherUserCleanup)
		includeUsers := func(users ...string) {
			t.Helper()
			_, err := client.Quotas.IncludeUsers(ctx, sdk.NewIncludeUsersQuotaRequest(id, *sdk.NewQuotaIncludeUsersArgsRequest(users)))
			require.NoError(t, err)
		}

		// names are resolved case-insensitively
		includeUsers(strings.ToLower(user.ID().Name()))
		assert.Equal(t, []string{user.ID().Name()}, userNames(describe(t, id).Scope.IncludedUsers))

		// the list accumulates
		includeUsers(otherUser.ID().Name())
		assert.ElementsMatch(t, []string{user.ID().Name(), otherUser.ID().Name()}, userNames(describe(t, id).Scope.IncludedUsers))

		// an empty array clears it
		includeUsers()
		assert.Empty(t, describe(t, id).Scope.IncludedUsers)
	})

	t.Run("exclude users", func(t *testing.T) {
		id := createQuota(t)
		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)
		otherUser, otherUserCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(otherUserCleanup)
		excludeUsers := func(users ...string) {
			t.Helper()
			_, err := client.Quotas.ExcludeUsers(ctx, sdk.NewExcludeUsersQuotaRequest(id, *sdk.NewQuotaExcludeUsersArgsRequest().
				WithByUser(*sdk.NewQuotaExcludeUsersByUserRequest(users))))
			require.NoError(t, err)
		}

		excludeUsers(user.ID().Name())
		assert.Equal(t, []string{user.ID().Name()}, userNames(describe(t, id).Scope.ExcludedUsers))

		excludeUsers(otherUser.ID().Name())
		assert.ElementsMatch(t, []string{user.ID().Name(), otherUser.ID().Name()}, userNames(describe(t, id).Scope.ExcludedUsers))

		excludeUsers()
		assert.Empty(t, describe(t, id).Scope.ExcludedUsers)
	})

	t.Run("set user tags", func(t *testing.T) {
		id := createQuota(t)
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		setUserTags := func(operator sdk.QuotaUserTagsOperator, values ...string) {
			t.Helper()
			tags := collections.Map(values, func(value string) sdk.QuotaUserTagRequest { return *sdk.NewQuotaUserTagRequestFromTag(tag.ID(), value) })
			_, err := client.Quotas.SetUserTags(ctx, sdk.NewSetUserTagsQuotaRequest(id, *sdk.NewQuotaSetUserTagsArgsRequest(tags, operator)))
			require.NoError(t, err)
		}

		setUserTags(sdk.QuotaUserTagsOperatorIntersection, "finance", "Analytics")
		scope := describe(t, id).Scope
		assert.Equal(t, string(sdk.QuotaUserTagsOperatorIntersection), scope.UserTags.Operator)
		require.Len(t, scope.SelectionTags(), 1)
		assert.Equal(t, tag.ID(), scope.SelectionTags()[0].TagId())
		assert.ElementsMatch(t, []string{"finance", "Analytics"}, tagValues(scope.SelectionTags()))

		// the tags are replaced
		setUserTags(sdk.QuotaUserTagsOperatorUnion, "finance")
		scope = describe(t, id).Scope
		assert.Equal(t, string(sdk.QuotaUserTagsOperatorUnion), scope.UserTags.Operator)
		assert.Equal(t, []string{"finance"}, tagValues(scope.SelectionTags()))

		setUserTags(sdk.QuotaUserTagsOperatorUnion)
		assert.Empty(t, describe(t, id).Scope.SelectionTags())
	})

	t.Run("exclude users by tag", func(t *testing.T) {
		id := createQuota(t)
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)
		testClientHelper().Quota.IncludeUsers(t, id, user.ID().Name())
		_, err := client.Quotas.SetUserTags(ctx, sdk.NewSetUserTagsQuotaRequest(id, *sdk.NewQuotaSetUserTagsArgsRequest(
			[]sdk.QuotaUserTagRequest{*sdk.NewQuotaUserTagRequestFromTag(tag.ID(), "finance")}, sdk.QuotaUserTagsOperatorUnion)))
		require.NoError(t, err)
		excludeTags := func(values ...string) {
			t.Helper()
			tags := collections.Map(values, func(value string) sdk.QuotaExcludedUserTagRequest {
				return *sdk.NewQuotaExcludedUserTagRequestFromTag(tag.ID(), value)
			})
			_, err := client.Quotas.ExcludeUsers(ctx, sdk.NewExcludeUsersQuotaRequest(id, *sdk.NewQuotaExcludeUsersArgsRequest().
				WithByTag(*sdk.NewQuotaExcludeUsersByTagRequest(tags))))
			require.NoError(t, err)
		}

		// excluded tags are kept apart from the tags that select users, and the user lists stay
		excludeTags("contractor")
		scope := describe(t, id).Scope
		assert.Equal(t, []string{"contractor"}, tagValues(scope.ExclusionTags()))
		assert.Equal(t, []string{"finance"}, tagValues(scope.SelectionTags()))
		assert.Equal(t, []string{user.ID().Name()}, userNames(scope.IncludedUsers))

		// the excluded tags are replaced
		excludeTags("intern")
		assert.Equal(t, []string{"intern"}, tagValues(describe(t, id).Scope.ExclusionTags()))

		excludeTags()
		scope = describe(t, id).Scope
		assert.Empty(t, scope.ExclusionTags())
		assert.Equal(t, []string{"finance"}, tagValues(scope.SelectionTags()))
	})

	t.Run("set and unset per-user limits", func(t *testing.T) {
		id := createQuota(t)
		setLimit := func(limit int, cycle sdk.QuotaCycle) {
			t.Helper()
			testClientHelper().Quota.SetPerUserLimit(t, id, limit, cycle)
		}
		unsetLimit := func(cycle sdk.QuotaCycle) {
			t.Helper()
			_, err := client.Quotas.UnsetPerUserLimit(ctx, sdk.NewUnsetPerUserLimitQuotaRequest(id, *sdk.NewQuotaUnsetPerUserLimitArgsRequest(cycle)))
			require.NoError(t, err)
		}

		setLimit(1000, sdk.QuotaCycleMonthly)
		assert.Equal(t, new(1000), describe(t, id).Config.PerUserLimit)

		setLimit(500, sdk.QuotaCycleWeekly)
		assert.Equal(t, new(500), describe(t, id).Config.PerUserLimitWeekly)

		setLimit(120, sdk.QuotaCycleDaily)
		assert.Equal(t, new(120), describe(t, id).Config.PerUserLimitDaily)

		_, err := client.Quotas.SetPerUserLimit(ctx, sdk.NewSetPerUserLimitQuotaRequest(id, *sdk.NewQuotaSetPerUserLimitArgsRequest(0)))
		require.ErrorContains(t, err, "Invalid per-user limit")

		unsetLimit(sdk.QuotaCycleMonthly)
		config := describe(t, id).Config
		assert.Nil(t, config.PerUserLimit)
		assert.Equal(t, new(500), config.PerUserLimitWeekly)

		unsetLimit(sdk.QuotaCycleWeekly)
		config = describe(t, id).Config
		assert.Nil(t, config.PerUserLimitWeekly)
		assert.Equal(t, new(120), config.PerUserLimitDaily)
	})

	t.Run("set block enforcement", func(t *testing.T) {
		id := quotaScopedToNewUser(t)
		require.NoError(t, addSharedResource(t, id, sdk.QuotaSharedResourceDomainAiFunction, ""))
		setBlockEnforcement := func(enabled bool, notifyUsers bool) {
			t.Helper()
			_, err := client.Quotas.SetBlockEnforcementEnabled(ctx, sdk.NewSetBlockEnforcementEnabledQuotaRequest(id,
				*sdk.NewQuotaSetBlockEnforcementEnabledArgsRequest(enabled).WithNotifyUsers(notifyUsers)))
			require.NoError(t, err)
		}

		setBlockEnforcement(true, false)
		config := describe(t, id).Config
		assert.True(t, config.BlockEnforcementEnabled)
		assert.False(t, config.PerUserBlockNotificationsEnabled)

		setBlockEnforcement(false, true)
		config = describe(t, id).Config
		assert.False(t, config.BlockEnforcementEnabled)
		assert.True(t, config.PerUserBlockNotificationsEnabled)
	})

	t.Run("add and remove notification thresholds", func(t *testing.T) {
		id := createQuota(t)
		addThreshold := func(threshold int, strategy sdk.QuotaSpendStrategy, notifyUser bool, cycle sdk.QuotaCycle) {
			t.Helper()
			_, err := client.Quotas.AddNotificationThreshold(ctx, sdk.NewAddNotificationThresholdQuotaRequest(id,
				*sdk.NewQuotaAddNotificationThresholdArgsRequest(threshold, strategy, notifyUser, cycle)))
			require.NoError(t, err)
		}
		thresholds := func() []sdk.QuotaNotificationThreshold {
			t.Helper()
			return describe(t, id).NotificationThresholds
		}

		addThreshold(80, sdk.QuotaSpendStrategyProjected, true, sdk.QuotaCycleMonthly)
		require.Len(t, thresholds(), 1)
		assert.Equal(t, sdk.QuotaNotificationThreshold{Threshold: 80, SpendStrategy: sdk.QuotaSpendStrategyProjected, NotifyUser: true, Cycle: sdk.QuotaCycleMonthly}, withoutTimestamp(thresholds()[0]))

		addThreshold(80, sdk.QuotaSpendStrategyProjected, false, sdk.QuotaCycleDaily)
		addThreshold(1000, sdk.QuotaSpendStrategyActual, false, sdk.QuotaCycleWeekly)
		assert.Len(t, thresholds(), 3)

		// Snowflake keeps one threshold per threshold, spend strategy, and cycle: adding it again replaces it
		addThreshold(80, sdk.QuotaSpendStrategyProjected, false, sdk.QuotaCycleMonthly)
		current := thresholds()
		require.Len(t, current, 3)
		assert.True(t, slices.ContainsFunc(current, func(threshold sdk.QuotaNotificationThreshold) bool {
			return threshold.Threshold == 80 && threshold.Cycle == sdk.QuotaCycleMonthly && !threshold.NotifyUser
		}))

		_, err := client.Quotas.RemoveNotificationThreshold(ctx, sdk.NewRemoveNotificationThresholdQuotaRequest(id,
			*sdk.NewQuotaRemoveNotificationThresholdArgsRequest(80, sdk.QuotaSpendStrategyProjected).WithCycle(sdk.QuotaCycleDaily)))
		require.NoError(t, err)
		current = thresholds()
		require.Len(t, current, 2)
		assert.False(t, slices.ContainsFunc(current, func(threshold sdk.QuotaNotificationThreshold) bool { return threshold.Cycle == sdk.QuotaCycleDaily }))

		// without the cycle, the monthly threshold is removed
		_, err = client.Quotas.RemoveNotificationThreshold(ctx, sdk.NewRemoveNotificationThresholdQuotaRequest(id,
			*sdk.NewQuotaRemoveNotificationThresholdArgsRequest(80, sdk.QuotaSpendStrategyProjected)))
		require.NoError(t, err)
		current = thresholds()
		require.Len(t, current, 1)
		assert.Equal(t, 1000, current[0].Threshold)
		assert.Equal(t, sdk.QuotaCycleWeekly, current[0].Cycle)
	})

	// quotas accept webhook and queue integrations, not email ones
	t.Run("add and remove notification integrations", func(t *testing.T) {
		id := createQuota(t)
		integration, integrationCleanup := testClientHelper().NotificationIntegration.CreateWebhook(t, quotaWebhookUrl)
		t.Cleanup(integrationCleanup)
		t.Cleanup(testClientHelper().Grant.GrantUsageOnIntegrationToSnowflakeApplication(t, integration.ID()))

		_, err := client.Quotas.AddNotificationIntegration(ctx, sdk.NewAddNotificationIntegrationQuotaRequest(id, *sdk.NewQuotaAddNotificationIntegrationArgsRequest(integration.ID())))
		require.NoError(t, err)
		integrations := describe(t, id).NotificationIntegrations
		require.Len(t, integrations, 1)
		assert.Equal(t, integration.ID().Name(), integrations[0].IntegrationName)

		_, err = client.Quotas.RemoveNotificationIntegration(ctx, sdk.NewRemoveNotificationIntegrationQuotaRequest(id, *sdk.NewQuotaRemoveNotificationIntegrationArgsRequest(integration.ID())))
		require.NoError(t, err)
		assert.Empty(t, describe(t, id).NotificationIntegrations)
	})

	// admin emails must be verified emails of users in the account; helpers.VerifiedEmail is verified in the test accounts
	t.Run("set admin emails", func(t *testing.T) {
		id := createQuota(t)

		_, err := client.Quotas.SetAdminEmails(ctx, sdk.NewSetAdminEmailsQuotaRequest(id, *sdk.NewQuotaSetAdminEmailsArgsRequestFromEmails(helpers.VerifiedEmail)))
		require.NoError(t, err)
		assert.Equal(t, []string{helpers.VerifiedEmail}, describe(t, id).Config.AdminEmails)

		// an empty list clears them
		_, err = client.Quotas.SetAdminEmails(ctx, sdk.NewSetAdminEmailsQuotaRequest(id, *sdk.NewQuotaSetAdminEmailsArgsRequestFromEmails()))
		require.NoError(t, err)
		assert.Empty(t, describe(t, id).Config.AdminEmails)
	})

	t.Run("add and remove custom actions", func(t *testing.T) {
		id := createQuota(t)
		// the first argument receives the IDs of the users who crossed the threshold; extra arguments fill the remaining ones
		procedure, procedureCleanup := testClientHelper().Procedure.Create(t, sdk.DataTypeVARCHAR)
		t.Cleanup(procedureCleanup)
		grantProcedureToSnowflakeApplication(t, procedure.ID())
		procedureWithArgument, procedureWithArgumentCleanup := testClientHelper().Procedure.Create(t, sdk.DataTypeVARCHAR, sdk.DataTypeVARCHAR)
		t.Cleanup(procedureWithArgumentCleanup)
		t.Cleanup(testClientHelper().Grant.GrantUsageOnProcedureToSnowflakeApplication(t, procedureWithArgument.ID()))
		addAction := func(procedureId sdk.SchemaObjectIdentifierWithArguments, arguments []string, strategy sdk.QuotaSpendStrategy, threshold int) error {
			t.Helper()
			_, err := client.Quotas.AddCustomAction(ctx, sdk.NewAddCustomActionQuotaRequest(id,
				*sdk.NewQuotaAddCustomActionArgsRequest(procedureId, arguments, strategy, threshold)))
			return err
		}
		actions := func() []sdk.QuotaCustomAction {
			t.Helper()
			return describe(t, id).CustomActions
		}

		require.NoError(t, addAction(procedure.ID(), nil, sdk.QuotaSpendStrategyActual, 100))
		current := actions()
		require.Len(t, current, 1)
		procedureId, err := current[0].ProcedureId()
		require.NoError(t, err)
		assert.Equal(t, procedure.ID().FullyQualifiedName(), procedureId.FullyQualifiedName())
		assert.Equal(t, sdk.QuotaSpendStrategyActual, current[0].SpendStrategy)
		assert.Equal(t, 100, current[0].Threshold)
		assert.NotEmpty(t, current[0].ActionId)

		require.NoError(t, addAction(procedure.ID(), nil, sdk.QuotaSpendStrategyProjected, 100))
		assert.Len(t, actions(), 2)

		require.NoError(t, addAction(procedureWithArgument.ID(), []string{"it's"}, sdk.QuotaSpendStrategyProjected, 90))
		current = actions()
		require.Len(t, current, 3)
		idx := slices.IndexFunc(current, func(action sdk.QuotaCustomAction) bool { return action.Threshold == 90 })
		require.GreaterOrEqual(t, idx, 0)
		assert.Equal(t, []string{"it's"}, current[idx].ProcedureArgs)

		// the number of extra arguments must match the procedure's remaining parameters
		err = addAction(procedure.ID(), []string{"extra"}, sdk.QuotaSpendStrategyActual, 50)
		require.ErrorContains(t, err, "ARGUMENTS_SIZE_MISMATCH")
		assert.Len(t, actions(), 3)

		// removing by threshold and procedure removes the actions of both spend strategies;
		// the procedure is passed as the PROCEDURE_FQN returned by GET_CUSTOM_ACTIONS
		idx = slices.IndexFunc(current, func(action sdk.QuotaCustomAction) bool { return action.Threshold == 100 })
		require.GreaterOrEqual(t, idx, 0)
		_, err = client.Quotas.RemoveCustomActions(ctx, sdk.NewRemoveCustomActionsQuotaRequest(id,
			*sdk.NewQuotaRemoveCustomActionsArgsRequest().WithThreshold(100).WithProcedureFqn(current[idx].ProcedureFqn)))
		require.NoError(t, err)
		current = actions()
		require.Len(t, current, 1)
		assert.Equal(t, 90, current[0].Threshold)

		_, err = client.Quotas.RemoveCustomActions(ctx, sdk.NewRemoveCustomActionsQuotaRequest(id, *sdk.NewQuotaRemoveCustomActionsArgsRequest()))
		require.NoError(t, err)
		assert.Empty(t, actions())
	})

	t.Run("set and remove the cycle-start action", func(t *testing.T) {
		id := createQuota(t)
		procedure, procedureCleanup := testClientHelper().Procedure.Create(t)
		t.Cleanup(procedureCleanup)
		grantProcedureToSnowflakeApplication(t, procedure.ID())

		_, err := client.Quotas.SetCycleStartAction(ctx, sdk.NewSetCycleStartActionQuotaRequest(id, *sdk.NewQuotaSetCycleStartActionArgsRequest(procedure.ID(), nil)))
		require.NoError(t, err)
		action := describe(t, id).CycleStartAction
		require.NotNil(t, action)
		procedureId, err := action.ProcedureId()
		require.NoError(t, err)
		assert.Equal(t, procedure.ID().FullyQualifiedName(), procedureId.FullyQualifiedName())
		assert.Empty(t, action.ProcedureArgs)

		_, err = client.Quotas.RemoveCycleStartAction(ctx, sdk.NewRemoveCycleStartActionQuotaRequest(id))
		require.NoError(t, err)
		assert.Nil(t, describe(t, id).CycleStartAction)
	})
}

// withoutTimestamp drops the fields that differ between runs, so that a threshold can be compared as a whole.
func withoutTimestamp(threshold sdk.QuotaNotificationThreshold) sdk.QuotaNotificationThreshold {
	return sdk.QuotaNotificationThreshold{
		Threshold:     threshold.Threshold,
		SpendStrategy: threshold.SpendStrategy,
		NotifyUser:    threshold.NotifyUser,
		Cycle:         threshold.Cycle,
	}
}

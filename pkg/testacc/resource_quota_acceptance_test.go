//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func quotaTestCase(t *testing.T, steps ...resource.TestStep) resource.TestCase {
	t.Helper()
	return resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Quota),
		Steps:        steps,
	}
}

func TestAcc_Quota_BasicUseCase(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	basic := model.QuotaFromId("test", id)
	configured := model.QuotaFromId("test", id).
		WithComment(comment).
		WithSharedResources(
			model.QuotaSharedResourceModel{Domain: sdk.QuotaSharedResourceDomainAiFunction},
			model.QuotaSharedResourceModel{Domain: sdk.QuotaSharedResourceDomainCortexCode, Instance: "CORTEX_CODE_CLI"},
		).
		WithIncludedUsers(user.ID().Name()).
		WithPerUserLimitMonthly(1000).
		WithPerUserLimitDaily(50)
	renamed := model.QuotaFromId("test", newId)

	ref := basic.ResourceReference()

	basicAssertions := func(id sdk.SchemaObjectIdentifier) []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceassert.QuotaResource(t, ref).
				HasDatabaseString(id.DatabaseName()).
				HasSchemaString(id.SchemaName()).
				HasNameString(id.Name()).
				HasFullyQualifiedNameString(id.FullyQualifiedName()).
				HasCommentString("").
				HasSharedResourcesCount(0).
				HasUserTagsEmpty().
				HasExcludedUserTagEmpty().
				HasIncludedUsersEmpty().
				HasExcludedUsersEmpty().
				HasPerUserLimitMonthly(0).
				HasPerUserLimitWeekly(0).
				HasPerUserLimitDaily(0),
			resourceshowoutputassert.QuotaShowOutput(t, ref).
				HasCreatedOnNotEmpty().
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasComment("").
				HasOwnerNotEmpty(),
		}
	}

	configuredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.QuotaResource(t, ref).
			HasNameString(id.Name()).
			HasCommentString(comment).
			HasSharedResourcesCount(2).
			HasSharedResource(sdk.QuotaSharedResourceDomainAiFunction, "").
			HasSharedResource(sdk.QuotaSharedResourceDomainCortexCode, "CORTEX_CODE_CLI").
			HasIncludedUsers(user.ID().Name()).
			HasPerUserLimitMonthly(1000).
			HasPerUserLimitWeekly(0).
			HasPerUserLimitDaily(50),
		resourceshowoutputassert.QuotaShowOutput(t, ref).
			HasComment(comment),
		assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.config.0.per_user_limit", "1000")),
		assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.scope.0.shared_resources.#", "2")),
	}

	resource.Test(t, quotaTestCase(t,
		// create with required fields only
		resource.TestStep{
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
				},
			},
			Config: accconfig.FromModels(t, basic),
			Check:  assertThat(t, basicAssertions(id)...),
		},
		// import with required fields only
		resource.TestStep{
			Config:            accconfig.FromModels(t, basic),
			ResourceName:      ref,
			ImportState:       true,
			ImportStateVerify: true,
		},
		// configure in place
		resource.TestStep{
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
				},
			},
			Config: accconfig.FromModels(t, configured),
			Check:  assertThat(t, configuredAssertions...),
		},
		// external changes are detected and reverted
		resource.TestStep{
			PreConfig: func() {
				testClient().Quota.SetPerUserLimit(t, id, 5000, sdk.QuotaCycleMonthly)
				testClient().Quota.AddSharedResource(t, id, sdk.QuotaSharedResourceDomainSnowflakeIntelligence, "")
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					planchecks.ExpectDrift(ref, "per_user_limit_monthly", sdk.String("1000"), sdk.String("5000")),
					planchecks.ExpectChange(ref, "per_user_limit_monthly", tfjson.ActionUpdate, sdk.String("5000"), sdk.String("1000")),
				},
			},
			Config: accconfig.FromModels(t, configured),
			Check:  assertThat(t, configuredAssertions...),
		},
		// unset everything
		resource.TestStep{
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
				},
			},
			Config: accconfig.FromModels(t, basic),
			Check:  assertThat(t, basicAssertions(id)...),
		},
		// rename in place
		resource.TestStep{
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
				},
			},
			Config: accconfig.FromModels(t, renamed),
			Check:  assertThat(t, basicAssertions(newId)...),
		},
	))
}

func TestAcc_Quota_CompleteUseCase(t *testing.T) {
	includedUser, includedUserCleanup := testClient().User.CreateUser(t)
	t.Cleanup(includedUserCleanup)
	excludedUser, excludedUserCleanup := testClient().User.CreateUser(t)
	t.Cleanup(excludedUserCleanup)
	tag, tagCleanup := testClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)
	agentId := testClient().Ids.RandomSchemaObjectIdentifier()
	t.Cleanup(testClient().CortexAgent.CreateWithId(t, agentId))

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	complete := model.QuotaFromId("test", id).
		WithComment(comment).
		WithSharedResources(
			model.QuotaSharedResourceModel{Domain: sdk.QuotaSharedResourceDomainAiFunction, Instance: "AI_CLASSIFY"},
			model.QuotaSharedResourceModel{Domain: sdk.QuotaSharedResourceDomainCortexAgent, Instance: agentId.FullyQualifiedName()},
			model.QuotaSharedResourceModel{Domain: sdk.QuotaSharedResourceDomainSnowflakeIntelligence},
		).
		WithUserTags(sdk.QuotaUserTagsOperatorIntersection,
			model.QuotaTagValueModel{Tag: tag.ID(), Value: "finance"},
			model.QuotaTagValueModel{Tag: tag.ID(), Value: "Analytics"},
		).
		WithExcludedUserTags(model.QuotaTagValueModel{Tag: tag.ID(), Value: "contractor"}).
		WithIncludedUsers(includedUser.ID().Name()).
		WithExcludedUsers(excludedUser.ID().Name()).
		WithPerUserLimitMonthly(1000).
		WithPerUserLimitWeekly(300).
		WithPerUserLimitDaily(60)

	ref := complete.ResourceReference()

	resource.Test(t, quotaTestCase(t,
		resource.TestStep{
			Config: accconfig.FromModels(t, complete),
			Check: assertThat(t,
				resourceassert.QuotaResource(t, ref).
					HasNameString(id.Name()).
					HasCommentString(comment).
					HasSharedResourcesCount(3).
					HasSharedResource(sdk.QuotaSharedResourceDomainAiFunction, "AI_CLASSIFY").
					HasSharedResource(sdk.QuotaSharedResourceDomainCortexAgent, agentId.FullyQualifiedName()).
					HasSharedResource(sdk.QuotaSharedResourceDomainSnowflakeIntelligence, "").
					HasUserTags(sdk.QuotaUserTagsOperatorIntersection, 2).
					HasUserTag(tag.ID(), "finance").
					HasUserTag(tag.ID(), "Analytics").
					HasExcludedUserTag(tag.ID(), "contractor").
					HasIncludedUsers(includedUser.ID().Name()).
					HasExcludedUsers(excludedUser.ID().Name()).
					HasPerUserLimitMonthly(1000).
					HasPerUserLimitWeekly(300).
					HasPerUserLimitDaily(60),
				resourceshowoutputassert.QuotaShowOutput(t, ref).
					HasName(id.Name()).
					HasComment(comment),
				assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.scope.0.user_tags_operator", "INTERSECTION")),
				assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.scope.0.excluded_user_tags.#", "1")),
			),
		},
		// every attribute is read back
		resource.TestStep{
			Config:            accconfig.FromModels(t, complete),
			ResourceName:      ref,
			ImportState:       true,
			ImportStateVerify: true,
		},
	))
}

// Values that Snowflake returns spelled differently (upper-case enums, quoted identifiers) don't produce diffs; every step
// checks that the plan is empty after the apply.
func TestAcc_Quota_EquivalentSpellings(t *testing.T) {
	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	tag, tagCleanup := testClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	aiDomains := model.QuotaFromId("test", id).
		WithSharedResources(
			model.QuotaSharedResourceModel{Domain: "ai function", Instance: "ai_classify"},
			model.QuotaSharedResourceModel{Domain: "cortex code"},
		).
		WithUserTags("union", model.QuotaTagValueModel{Tag: tag.ID(), Value: "finance"}).
		WithIncludedUsers(user.ID().FullyQualifiedName())
	warehouseDomain := model.QuotaFromId("test", id).
		WithSharedResources(model.QuotaSharedResourceModel{Domain: sdk.QuotaSharedResourceDomainWarehouse, Instance: warehouse.ID().FullyQualifiedName()})

	ref := aiDomains.ResourceReference()

	resource.Test(t, quotaTestCase(t,
		resource.TestStep{
			Config: accconfig.FromModels(t, aiDomains),
			Check: assertThat(t,
				resourceassert.QuotaResource(t, ref).
					HasSharedResourcesCount(2).
					HasUserTags(sdk.QuotaUserTagsOperatorUnion, 1),
			),
		},
		// switching between credit types removes the AI domains before adding the warehouse
		resource.TestStep{
			Config: accconfig.FromModels(t, warehouseDomain),
			Check: assertThat(t,
				resourceassert.QuotaResource(t, ref).
					HasSharedResourcesCount(1),
			),
		},
	))
}

func TestAcc_Quota_UserTags(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)
	tag, tagCleanup := testClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	excludedTags := model.QuotaFromId("test", id).
		WithExcludedUserTags(model.QuotaTagValueModel{Tag: tag.ID(), Value: "contractor"}).
		WithIncludedUsers(user.ID().Name())
	bothTags := model.QuotaFromId("test", id).
		WithUserTags(sdk.QuotaUserTagsOperatorUnion, model.QuotaTagValueModel{Tag: tag.ID(), Value: "finance"}).
		WithExcludedUserTags(model.QuotaTagValueModel{Tag: tag.ID(), Value: "contractor"}).
		WithIncludedUsers(user.ID().Name())
	noTags := model.QuotaFromId("test", id).
		WithIncludedUsers(user.ID().Name())

	ref := excludedTags.ResourceReference()

	resource.Test(t, quotaTestCase(t,
		resource.TestStep{
			Config: accconfig.FromModels(t, excludedTags),
			Check: assertThat(t,
				resourceassert.QuotaResource(t, ref).
					HasUserTagsEmpty().
					HasExcludedUserTag(tag.ID(), "contractor").
					HasIncludedUsers(user.ID().Name()),
			),
		},
		resource.TestStep{
			Config:            accconfig.FromModels(t, excludedTags),
			ResourceName:      ref,
			ImportState:       true,
			ImportStateVerify: true,
		},
		// tags that select users and excluded tags are set independently, and the user lists stay
		resource.TestStep{
			Config: accconfig.FromModels(t, bothTags),
			Check: assertThat(t,
				resourceassert.QuotaResource(t, ref).
					HasUserTags(sdk.QuotaUserTagsOperatorUnion, 1).
					HasUserTag(tag.ID(), "finance").
					HasExcludedUserTag(tag.ID(), "contractor").
					HasIncludedUsers(user.ID().Name()),
			),
		},
		resource.TestStep{
			Config: accconfig.FromModels(t, noTags),
			Check: assertThat(t,
				resourceassert.QuotaResource(t, ref).
					HasUserTagsEmpty().
					HasExcludedUserTagEmpty().
					HasIncludedUsers(user.ID().Name()),
			),
		},
	))
}

func TestAcc_Quota_Domains(t *testing.T) {
	agentId := testClient().Ids.RandomSchemaObjectIdentifier()
	t.Cleanup(testClient().CortexAgent.CreateWithId(t, agentId))
	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	ref := model.QuotaFromId("test", id).ResourceReference()

	// Each step replaces the monitored domain in place; the switch to the WAREHOUSE domain checks that removals happen first.
	steps := make([]resource.TestStep, 0)
	for _, sharedResource := range []model.QuotaSharedResourceModel{
		{Domain: sdk.QuotaSharedResourceDomainAiFunction},
		{Domain: sdk.QuotaSharedResourceDomainAiFunction, Instance: "AI_CLASSIFY"},
		{Domain: sdk.QuotaSharedResourceDomainCortexCode},
		{Domain: sdk.QuotaSharedResourceDomainCortexCode, Instance: "CORTEX_CODE_SNOWSIGHT"},
		{Domain: sdk.QuotaSharedResourceDomainCortexAgent},
		{Domain: sdk.QuotaSharedResourceDomainCortexAgent, Instance: agentId.FullyQualifiedName()},
		{Domain: sdk.QuotaSharedResourceDomainSnowflakeIntelligence},
		{Domain: sdk.QuotaSharedResourceDomainAiGateway},
		{Domain: sdk.QuotaSharedResourceDomainWarehouse},
		{Domain: sdk.QuotaSharedResourceDomainWarehouse, Instance: warehouse.ID().Name()},
	} {
		config := model.QuotaFromId("test", id).WithSharedResources(sharedResource)
		steps = append(steps, resource.TestStep{
			Config: accconfig.FromModels(t, config),
			Check: assertThat(t,
				resourceassert.QuotaResource(t, ref).
					HasSharedResourcesCount(1).
					HasSharedResource(sharedResource.Domain, sharedResource.Instance),
			),
		})
	}
	steps = append(steps, resource.TestStep{
		Config:            accconfig.FromModels(t, model.QuotaFromId("test", id).WithSharedResources(model.QuotaSharedResourceModel{Domain: sdk.QuotaSharedResourceDomainWarehouse, Instance: warehouse.ID().Name()})),
		ResourceName:      ref,
		ImportState:       true,
		ImportStateVerify: true,
	})

	resource.Test(t, quotaTestCase(t, steps...))
}

// The provider validates only the enum values; everything else is left to Snowflake.
func TestAcc_Quota_Validations(t *testing.T) {
	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, quotaTestCase(t,
		resource.TestStep{
			Config:      accconfig.FromModels(t, model.QuotaFromId("test", id).WithSharedResources(model.QuotaSharedResourceModel{Domain: "NOT A DOMAIN"})),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile("invalid quota shared resource domain"),
		},
		resource.TestStep{
			Config: accconfig.FromModels(t, model.QuotaFromId("test", id).WithSharedResources(
				model.QuotaSharedResourceModel{Domain: sdk.QuotaSharedResourceDomainAiFunction},
				model.QuotaSharedResourceModel{Domain: sdk.QuotaSharedResourceDomainWarehouse, Instance: warehouse.ID().Name()},
			)),
			ExpectError: regexp.MustCompile("different credit types"),
		},
	))
}

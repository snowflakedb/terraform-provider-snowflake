//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SessionPolicy_BasicUseCase(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(secondSchema.ID())
	comment := random.Comment()
	externalComment := random.Comment()

	sessionIdleTimeoutMins := random.IntRange(5, 1440)
	externalSessionIdleTimeoutMins := random.IntRange(5, 1440)

	sessionUiIdleTimeoutMins := random.IntRange(5, 1440)
	externalSessionUiIdleTimeoutMins := random.IntRange(5, 1440)

	role1, role1Cleanup := testClient().Role.CreateRole(t)
	t.Cleanup(role1Cleanup)
	role2, role2Cleanup := testClient().Role.CreateRole(t)
	t.Cleanup(role2Cleanup)
	role3, role3Cleanup := testClient().Role.CreateRole(t)
	t.Cleanup(role3Cleanup)

	basic := model.SessionPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name())

	altered := model.SessionPolicy("t", newId.DatabaseName(), newId.SchemaName(), newId.Name()).
		WithSessionIdleTimeoutMins(sessionIdleTimeoutMins).
		WithSessionUiIdleTimeoutMins(sessionUiIdleTimeoutMins).
		WithAllowedSecondaryRolesNone().
		WithBlockedSecondaryRolesAll().
		WithComment(comment)

	altered2 := model.SessionPolicy("t", newId.DatabaseName(), newId.SchemaName(), newId.Name()).
		WithSessionIdleTimeoutMins(sessionIdleTimeoutMins).
		WithSessionUiIdleTimeoutMins(sessionUiIdleTimeoutMins).
		WithAllowedSecondaryRolesAll().
		WithBlockedSecondaryRolesNone().
		WithComment(comment)

	allAttributes := model.SessionPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithSessionIdleTimeoutMins(sessionIdleTimeoutMins).
		WithSessionUiIdleTimeoutMins(sessionUiIdleTimeoutMins).
		WithAllowedSecondaryRoles(role1.Name, role2.Name).
		WithBlockedSecondaryRoles(role3.Name).
		WithComment(comment)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.SessionPolicyResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasNoSessionIdleTimeoutMins().
			HasNoSessionUiIdleTimeoutMins().
			HasAllowedSecondaryRolesEmpty().
			HasBlockedSecondaryRolesEmpty().
			HasCommentEmpty(),
		resourceshowoutputassert.SessionPolicyShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("SESSION_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		resourceshowoutputassert.SessionPolicyDescribeOutput(t, ref).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("").
			HasSessionIdleTimeoutMins(240).
			HasSessionUiIdleTimeoutMins(1080).
			HasAllowedSecondaryRoles("ALL").
			HasNoBlockedSecondaryRoles(),
	}

	basicAssertionsWithZeros := append([]assert.TestCheckFuncProvider{
		resourceassert.SessionPolicyResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasSessionIdleTimeoutMins(0).
			HasSessionUiIdleTimeoutMins(0).
			HasAllowedSecondaryRolesEmpty().
			HasBlockedSecondaryRolesEmpty().
			HasCommentEmpty(),
	}, basicAssertions[1:]...)

	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.SessionPolicyResource(t, ref).
			HasName(newId.Name()).
			HasSchema(newId.SchemaName()).
			HasDatabase(newId.DatabaseName()).
			HasSessionIdleTimeoutMins(sessionIdleTimeoutMins).
			HasSessionUiIdleTimeoutMins(sessionUiIdleTimeoutMins).
			HasNoAllowedSecondaryRoles().
			HasAllBlockedSecondaryRoles().
			HasComment(comment),
		resourceshowoutputassert.SessionPolicyShowOutput(t, ref).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasKind("SESSION_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		resourceshowoutputassert.SessionPolicyDescribeOutput(t, ref).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(comment).
			HasSessionIdleTimeoutMins(sessionIdleTimeoutMins).
			HasSessionUiIdleTimeoutMins(sessionUiIdleTimeoutMins).
			HasNoAllowedSecondaryRoles().
			HasBlockedSecondaryRoles("ALL"),
	}

	alteredAssertions2 := []assert.TestCheckFuncProvider{
		resourceassert.SessionPolicyResource(t, ref).
			HasName(newId.Name()).
			HasSchema(newId.SchemaName()).
			HasDatabase(newId.DatabaseName()).
			HasSessionIdleTimeoutMins(sessionIdleTimeoutMins).
			HasSessionUiIdleTimeoutMins(sessionUiIdleTimeoutMins).
			HasAllAllowedSecondaryRoles().
			HasNoBlockedSecondaryRoles().
			HasComment(comment),
		resourceshowoutputassert.SessionPolicyShowOutput(t, ref).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasKind("SESSION_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		resourceshowoutputassert.SessionPolicyDescribeOutput(t, ref).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(comment).
			HasSessionIdleTimeoutMins(sessionIdleTimeoutMins).
			HasSessionUiIdleTimeoutMins(sessionUiIdleTimeoutMins).
			HasAllowedSecondaryRoles("ALL").
			HasNoBlockedSecondaryRoles(),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.SessionPolicyResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasSessionIdleTimeoutMins(sessionIdleTimeoutMins).
			HasSessionUiIdleTimeoutMins(sessionUiIdleTimeoutMins).
			HasAllowedSecondaryRoles(role1.Name, role2.Name).
			HasBlockedSecondaryRoles(role3.Name).
			HasComment(comment),
		resourceshowoutputassert.SessionPolicyShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("SESSION_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		resourceshowoutputassert.SessionPolicyDescribeOutput(t, ref).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(comment).
			HasSessionIdleTimeoutMins(sessionIdleTimeoutMins).
			HasSessionUiIdleTimeoutMins(sessionUiIdleTimeoutMins).
			HasAllowedSecondaryRoles(role1.Name, role2.Name).
			HasBlockedSecondaryRoles(role3.Name),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SessionPolicy),
		Steps: []resource.TestStep{
			// Create
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change alterable props (including cross-schema rename)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, altered),
				Check:  assertThat(t, alteredAssertions...),
			},
			// Change all secondary roles to none and vice versa
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, altered2),
				Check:  assertThat(t, alteredAssertions2...),
			},
			// Unset
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertionsWithZeros...),
			},
			// Destroy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroy),
					},
				},
				Config:  config.FromModels(t, basic),
				Destroy: true,
			},
			// Create with all attributes
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// Import
			{
				Config:            config.FromModels(t, allAttributes),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"session_idle_timeout_mins",
					"session_ui_idle_timeout_mins",
					"allowed_secondary_roles",
					"blocked_secondary_roles",
				},
			},
			// Change props externally
			{
				PreConfig: func() {
					alterRequest := sdk.NewAlterSessionPolicyRequest(id).WithSet(*sdk.NewSessionPolicySetRequest().
						WithSessionIdleTimeoutMins(externalSessionIdleTimeoutMins).
						WithSessionUiIdleTimeoutMins(externalSessionUiIdleTimeoutMins).
						WithAllowedSecondaryRoles(*sdk.NewSessionPolicySecondaryRolesRequest().WithNone(true)).
						WithBlockedSecondaryRoles(*sdk.NewSessionPolicySecondaryRolesRequest().WithAll(true)).
						WithComment(externalComment),
					)
					testClient().SessionPolicy.Alter(t, alterRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: func() []plancheck.PlanCheck {
						sessionIdle := sdk.String(strconv.Itoa(sessionIdleTimeoutMins))
						sessionUiIdle := sdk.String(strconv.Itoa(sessionUiIdleTimeoutMins))
						externalSessionIdle := sdk.String(strconv.Itoa(externalSessionIdleTimeoutMins))
						externalSessionUiIdle := sdk.String(strconv.Itoa(externalSessionUiIdleTimeoutMins))
						blockedSecondaryRoles := sdk.String(fmt.Sprintf("[%s]", role3.Name))
						return []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
							planchecks.ExpectDrift(ref, "session_idle_timeout_mins", sessionIdle, externalSessionIdle),
							planchecks.ExpectDrift(ref, "session_ui_idle_timeout_mins", sessionUiIdle, externalSessionUiIdle),
							// Don't check allowed_secondary_roles.0.roles, as the order of elements is unpredictable
							planchecks.ExpectDrift(ref, "allowed_secondary_roles.0.none", sdk.String("false"), sdk.String("true")),
							planchecks.ExpectNoChangeOnField(ref, "allowed_secondary_roles.0.all"),
							planchecks.ExpectNoChangeOnField(ref, "blocked_secondary_roles.0.none"),
							planchecks.ExpectDrift(ref, "blocked_secondary_roles.0.roles", blockedSecondaryRoles, sdk.String("[]")),
							planchecks.ExpectDrift(ref, "blocked_secondary_roles.0.all", sdk.String("false"), sdk.String("true")),
							planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
							planchecks.ExpectChange(ref, "session_idle_timeout_mins", tfjson.ActionUpdate, externalSessionIdle, sessionIdle),
							planchecks.ExpectChange(ref, "session_ui_idle_timeout_mins", tfjson.ActionUpdate, externalSessionUiIdle, sessionUiIdle),
							planchecks.ExpectChange(ref, "allowed_secondary_roles.0.none", tfjson.ActionUpdate, sdk.String("true"), nil),
							planchecks.ExpectChange(ref, "blocked_secondary_roles.0.roles", tfjson.ActionUpdate, sdk.String("[]"), blockedSecondaryRoles),
							planchecks.ExpectChange(ref, "blocked_secondary_roles.0.all", tfjson.ActionUpdate, sdk.String("true"), nil),
							planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
						}
					}(),
				},
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
		},
	})
}

func TestAcc_SessionPolicy_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	invalidSessionIdleTimeoutMins := model.SessionPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithSessionIdleTimeoutMins(0)

	invalidSessionUiIdleTimeoutMins := model.SessionPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithSessionUiIdleTimeoutMins(0)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SessionPolicy),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidSessionIdleTimeoutMins),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected session_idle_timeout_mins to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, invalidSessionUiIdleTimeoutMins),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected session_ui_idle_timeout_mins to be at least \(1\), got 0`),
			},
		},
	})
}

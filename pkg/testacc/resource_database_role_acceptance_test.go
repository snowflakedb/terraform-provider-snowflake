//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DatabaseRole_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomDatabaseObjectIdentifier()
	newId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(id.DatabaseId())
	comment := random.Comment()
	currentRole := testClient().Context.CurrentRole(t)

	basic := model.DatabaseRole("test", id.DatabaseName(), id.Name())

	complete := model.DatabaseRole("test", newId.DatabaseName(), newId.Name()).
		WithComment(comment)

	assertBasic := assertThat(t,
		objectassert.DatabaseRole(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasIsDefault(false).
			HasIsCurrent(false).
			HasIsInherited(false).
			HasGrantedToRoles(0).
			HasGrantedToDatabaseRoles(0).
			HasGrantedDatabaseRoles(0).
			HasOwner(currentRole.Name()).
			HasComment("").
			HasOwnerRoleType("ROLE"),

		resourceassert.DatabaseRoleResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasCommentString("").
			HasFullyQualifiedNameString(id.FullyQualifiedName()),

		resourceshowoutputassert.DatabaseRoleShowOutput(t, basic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasIsDefault(false).
			HasIsCurrent(false).
			HasIsInherited(false).
			HasGrantedToRoles(0).
			HasGrantedToDatabaseRoles(0).
			HasGrantedDatabaseRoles(0).
			HasOwnerNotEmpty().
			HasComment("").
			HasOwnerRoleTypeNotEmpty(),
	)

	assertComplete := assertThat(t,
		objectassert.DatabaseRole(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasIsDefault(false).
			HasIsCurrent(false).
			HasIsInherited(false).
			HasGrantedToRoles(0).
			HasGrantedToDatabaseRoles(0).
			HasGrantedDatabaseRoles(0).
			HasOwner(currentRole.Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE"),

		resourceassert.DatabaseRoleResource(t, complete.ResourceReference()).
			HasNameString(newId.Name()).
			HasDatabaseString(newId.DatabaseName()).
			HasCommentString(comment).
			HasFullyQualifiedNameString(newId.FullyQualifiedName()),

		resourceshowoutputassert.DatabaseRoleShowOutput(t, complete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasIsDefault(false).
			HasIsCurrent(false).
			HasIsInherited(false).
			HasGrantedToRoles(0).
			HasGrantedToDatabaseRoles(0).
			HasGrantedDatabaseRoles(0).
			HasOwnerNotEmpty().
			HasComment(comment).
			HasOwnerRoleTypeNotEmpty(),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DatabaseRole),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: config.FromModels(t, basic),
				Check:  assertBasic,
			},

			// Import - without optionals
			{
				Config:       config.FromModels(t, basic),
				ResourceName: basic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedDatabaseRoleResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasCommentString(""),
					resourceshowoutputassert.ImportedWarehouseShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasComment(""),
				),
			},

			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertComplete,
			},

			// Import - with optionals
			{
				Config:       config.FromModels(t, complete),
				ResourceName: complete.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedDatabaseRoleResource(t, helpers.EncodeResourceIdentifier(newId)).
						HasNameString(newId.Name()).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedWarehouseShowOutput(t, helpers.EncodeResourceIdentifier(newId)).
						HasName(newId.Name()).
						HasComment(comment),
				),
			},

			// Update - unset optionals (back to basic)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertBasic,
			},

			// Update - detect external changes
			{
				PreConfig: func() {
					testClient().DatabaseRole.Alter(t, sdk.NewAlterDatabaseRoleRequest(id).WithSet(*sdk.NewDatabaseRoleSetRequest(random.Comment())))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertBasic,
			},

			// Create - with optionals
			{
				Taint: []string{complete.ResourceReference()},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertComplete,
			},
		},
	})
}

func TestAcc_DatabaseRole_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := testClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()
	databaseRoleModelWithComment := model.DatabaseRole("test", id.DatabaseName(), id.Name()).WithComment(comment)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            config.FromModels(t, databaseRoleModelWithComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database_role.test", "id", fmt.Sprintf(`%s|%s`, id.DatabaseName(), id.Name())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, databaseRoleModelWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database_role.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database_role.test", "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_DatabaseRole_IdentifierQuotingDiffSuppression(t *testing.T) {
	id := testClient().Ids.RandomDatabaseObjectIdentifier()
	quotedDatabaseRoleId := fmt.Sprintf(`"%s"`, id.Name())
	comment := random.Comment()
	databaseRoleModelWithComment := model.DatabaseRole("test", id.DatabaseName(), quotedDatabaseRoleId).WithComment(comment)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             config.FromModels(t, databaseRoleModelWithComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database_role.test", "database", id.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_database_role.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database_role.test", "id", fmt.Sprintf(`%s|%s`, id.DatabaseName(), id.Name())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, databaseRoleModelWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database_role.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database_role.test", "database", id.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_database_role.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database_role.test", "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}

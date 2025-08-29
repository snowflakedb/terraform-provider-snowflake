//go:build !account_level_tests

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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DatabaseRole_Basic(t *testing.T) {
	id := testClient().Ids.RandomDatabaseObjectIdentifier()
	newId := testClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()
	databaseRoleModel := model.DatabaseRole("test", id.DatabaseName(), id.Name())
	databaseRoleModelWithComment := model.DatabaseRole("test", id.DatabaseName(), id.Name()).WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DatabaseRole),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseRoleModel),
				Check: assertThat(t,
					resourceassert.DatabaseRoleResource(t, "snowflake_database_role.test").
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.DatabaseRoleShowOutput(t, "snowflake_database_role.test").
						HasName(id.Name()).
						HasComment(""),
					objectassert.DatabaseRole(t, id).
						HasName(id.Name()).
						HasComment(""),
				),
			},
			{
				ResourceName: "snowflake_database_role.test",
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
			// set comment
			{
				Config: config.FromModels(t, databaseRoleModelWithComment),
				Check: assertThat(t,
					resourceassert.DatabaseRoleResource(t, "snowflake_database_role.test").
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasCommentString(comment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.DatabaseRoleShowOutput(t, "snowflake_database_role.test").
						HasName(id.Name()).
						HasComment(comment),
					objectassert.DatabaseRole(t, id).
						HasName(id.Name()).
						HasComment(comment),
				),
			},
			{
				ResourceName: "snowflake_database_role.test",
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedDatabaseRoleResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedWarehouseShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasComment(comment),
				),
			},
			// unset comment
			{
				Config: config.FromModels(t, databaseRoleModel),
				Check: assertThat(t,
					resourceassert.DatabaseRoleResource(t, "snowflake_database_role.test").
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.DatabaseRoleShowOutput(t, "snowflake_database_role.test").
						HasName(id.Name()).
						HasComment(""),
					objectassert.DatabaseRole(t, id).
						HasName(id.Name()).
						HasComment(""),
				),
			},
			{
				ResourceName: "snowflake_database_role.test",
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
			// rename
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database_role.test", plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, databaseRoleModel.WithName(newId.Name())),
				Check: assertThat(t,
					resourceassert.DatabaseRoleResource(t, "snowflake_database_role.test").
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasCommentString("").
						HasFullyQualifiedNameString(newId.FullyQualifiedName()),
					resourceshowoutputassert.DatabaseRoleShowOutput(t, "snowflake_database_role.test").
						HasName(newId.Name()).
						HasComment(""),
					objectassert.DatabaseRole(t, newId).
						HasName(newId.Name()).
						HasComment(""),
				),
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

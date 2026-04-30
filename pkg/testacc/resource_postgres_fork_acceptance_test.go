//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_PostgresFork_BasicUseCase(t *testing.T) {
	sourceId := testClient().Ids.RandomAccountObjectIdentifier()
	forkId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	// Create source postgres instance for forking
	_, sourceCleanup := testClient().PostgresInstance.CreateWithRequest(t,
		sdk.NewCreatePostgresInstanceRequest(sourceId, "STANDARD_1", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres))
	t.Cleanup(sourceCleanup)

	modelBasic := model.PostgresFork("test", forkId.Name(), sourceId.Name())

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.PostgresForkResource(t, modelBasic.ResourceReference()).
			HasNameString(forkId.Name()).
			HasForkFromString(sourceId.Name()).
			HasNoComment().
			HasNoAtTimestamp().
			HasNoAtOffset().
			HasNoBeforeTimestamp().
			HasNoBeforeOffset().
			HasFullyQualifiedNameString(forkId.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, modelBasic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(forkId.Name()).
			HasType("FORK").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"),
	}

	modelWithComment := model.PostgresFork("test", forkId.Name(), sourceId.Name()).
		WithComment(comment)

	assertWithComment := []assert.TestCheckFuncProvider{
		resourceassert.PostgresForkResource(t, modelWithComment.ResourceReference()).
			HasNameString(forkId.Name()).
			HasForkFromString(sourceId.Name()).
			HasCommentString(comment).
			HasFullyQualifiedNameString(forkId.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, modelWithComment.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(forkId.Name()).
			HasType("FORK").
			HasComment(comment).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresFork),
		Steps: []resource.TestStep{
			// Create fork - basic
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, assertBasic...),
			},
			// Import
			{
				Config:            accconfig.FromModels(t, modelBasic),
				ResourceName:      modelBasic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"at_timestamp",
					"at_offset",
					"before_timestamp",
					"before_offset",
					"fork_from",
				},
			},
			// Update - set comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithComment.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelWithComment),
				Check:  assertThat(t, assertWithComment...),
			},
			// Update - unset comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.PostgresForkResource(t, modelBasic.ResourceReference()).
						HasNameString(forkId.Name()).
						HasCommentString(""),
				),
			},
			// Destroy
			{
				Destroy: true,
				Config:  accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					invokeactionassert.PostgresForkDoesNotExist(t, forkId),
				),
			},
		},
	})
}

func TestAcc_PostgresFork_WithAtTimestamp(t *testing.T) {
	sourceId := testClient().Ids.RandomAccountObjectIdentifier()
	forkId := testClient().Ids.RandomAccountObjectIdentifier()

	// Create source postgres instance for forking
	_, sourceCleanup := testClient().PostgresInstance.CreateWithRequest(t,
		sdk.NewCreatePostgresInstanceRequest(sourceId, "STANDARD_1", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres))
	t.Cleanup(sourceCleanup)

	modelFork := model.PostgresFork("test", forkId.Name(), sourceId.Name()).
		WithAtTimestamp("2025-01-15 12:00:00")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresFork),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelFork),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelFork.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.PostgresForkResource(t, modelFork.ResourceReference()).
						HasNameString(forkId.Name()).
						HasForkFromString(sourceId.Name()).
						HasAtTimestampString("2025-01-15 12:00:00"),
					resourceshowoutputassert.PostgresInstanceShowOutput(t, modelFork.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(forkId.Name()).
						HasType("FORK"),
				),
			},
		},
	})
}

func TestAcc_PostgresFork_UpdateAfterFork(t *testing.T) {
	sourceId := testClient().Ids.RandomAccountObjectIdentifier()
	forkId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	// Create source postgres instance for forking
	_, sourceCleanup := testClient().PostgresInstance.CreateWithRequest(t,
		sdk.NewCreatePostgresInstanceRequest(sourceId, "STANDARD_1", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres))
	t.Cleanup(sourceCleanup)

	modelBasic := model.PostgresFork("test", forkId.Name(), sourceId.Name())
	modelUpdated := model.PostgresFork("test", forkId.Name(), sourceId.Name()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresFork),
		Steps: []resource.TestStep{
			// Create
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.PostgresForkResource(t, modelBasic.ResourceReference()).
						HasNameString(forkId.Name()),
					resourceshowoutputassert.PostgresInstanceShowOutput(t, modelBasic.ResourceReference()).
						HasType("FORK"),
				),
			},
			// Update - set comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelUpdated.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelUpdated),
				Check: assertThat(t,
					resourceassert.PostgresForkResource(t, modelUpdated.ResourceReference()).
						HasNameString(forkId.Name()).
						HasCommentString(comment),
					resourceshowoutputassert.PostgresInstanceShowOutput(t, modelUpdated.ResourceReference()).
						HasType("FORK").
						HasComment(comment),
				),
			},
		},
	})
}

func TestAcc_PostgresFork_Validations(t *testing.T) {
	sourceId := testClient().Ids.RandomAccountObjectIdentifier()
	forkId := testClient().Ids.RandomAccountObjectIdentifier()

	// Model with both at_timestamp and before_timestamp set — should fail validation
	modelConflict := model.PostgresFork("test", forkId.Name(), sourceId.Name()).
		WithAtTimestamp("2025-01-15 12:00:00").
		WithBeforeTimestamp("2025-01-15 11:00:00")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresFork),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelConflict),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"at_timestamp": conflicts with before_timestamp`),
			},
		},
	})
}

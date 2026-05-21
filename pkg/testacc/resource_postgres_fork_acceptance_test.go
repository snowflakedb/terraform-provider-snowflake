//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// createSourceForFork creates a postgres instance suitable for forking and registers cleanup.
// It waits for the instance to be fully ready to accept fork operations.
func createSourceForFork(t *testing.T) sdk.AccountObjectIdentifier {
	t.Helper()
	sourceId := testClient().Ids.RandomAccountObjectIdentifier()
	_, sourceCleanup := testClient().PostgresInstance.CreateWithRequest(t,
		sdk.NewCreatePostgresInstanceRequest(sourceId, "STANDARD_M", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres))
	t.Cleanup(sourceCleanup)
	testClient().PostgresInstance.WaitForForkReady(t, sourceId, 5*time.Minute)
	return sourceId
}

func TestAcc_PostgresFork_Basic(t *testing.T) {
	sourceId := createSourceForFork(t)
	forkId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

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
		postgresShowOutputBaseAssert(t, modelBasic.ResourceReference(), forkId.Name()),
	}

	modelWithComment := model.PostgresFork("test", forkId.Name(), sourceId.Name()).
		WithComment(comment)

	assertWithComment := []assert.TestCheckFuncProvider{
		resourceassert.PostgresForkResource(t, modelWithComment.ResourceReference()).
			HasNameString(forkId.Name()).
			HasForkFromString(sourceId.Name()).
			HasCommentString(comment).
			HasFullyQualifiedNameString(forkId.FullyQualifiedName()),
		postgresShowOutputBaseAssert(t, modelWithComment.ResourceReference(), forkId.Name()).
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresFork),
		Steps: []resource.TestStep{
			// Step 1: Create fork with only required params
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertBasic...),
			},
			// Step 2: Import
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
			// Step 3: Update — set comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithComment.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelWithComment),
				Check:  assertThat(t, assertWithComment...),
			},
			// Step 4: Update — unset comment
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
			// Step 5: Destroy verification
			{
				Destroy: true,
				Config:  accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					invokeactionassert.PostgresInstanceDoesNotExist(t, forkId),
				),
			},
		},
	})
}

func TestAcc_PostgresFork_Complete(t *testing.T) {
	sourceId := createSourceForFork(t)
	forkId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	// Use a recent timestamp (1 hour ago) — Snowflake requires fork timestamps within the past 10 days
	recentTimestamp := time.Now().Add(-1 * time.Hour).UTC().Format("2006-01-02 15:04:05")

	modelComplete := model.PostgresFork("test", forkId.Name(), sourceId.Name()).
		WithComment(comment).
		WithAtTimestamp(recentTimestamp).
		WithHighAvailability(false).
		WithPostgresSettings(`{"work_mem": "64MB"}`)

	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.PostgresForkResource(t, modelComplete.ResourceReference()).
			HasNameString(forkId.Name()).
			HasForkFromString(sourceId.Name()).
			HasCommentString(comment).
			HasAtTimestampString(recentTimestamp).
			HasHighAvailabilityString("false").
			HasFullyQualifiedNameString(forkId.FullyQualifiedName()),
		postgresShowOutputBaseAssert(t, modelComplete.ResourceReference(), forkId.Name()).
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresFork),
		Steps: []resource.TestStep{
			// Step 1: Create fork with all options
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check:  assertThat(t, assertComplete...),
			},
			// Step 2: Import
			{
				Config:            accconfig.FromModels(t, modelComplete),
				ResourceName:      modelComplete.ResourceReference(),
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

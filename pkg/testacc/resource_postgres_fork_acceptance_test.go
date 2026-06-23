//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"
	"time"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
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

// createSourceForFork creates a postgres instance suitable for forking and registers cleanup.
// It waits for the instance to reach READY state; the resource itself retries the fork
// operation until the backend accepts it.
func createSourceForFork(t *testing.T) sdk.AccountObjectIdentifier {
	t.Helper()
	sourceId := testClient().Ids.RandomAccountObjectIdentifier()
	_, sourceCleanup := testClient().PostgresInstance.CreateWithRequest(t,
		sdk.NewCreatePostgresInstanceRequest(sourceId, "STANDARD_M", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres))
	t.Cleanup(sourceCleanup)
	testClient().PostgresInstance.WaitForReady(t, sourceId, 5*time.Minute)
	return sourceId
}

func TestAcc_PostgresFork_BasicUseCase(t *testing.T) {
	sourceId := createSourceForFork(t)
	forkId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	modelBasic := model.PostgresFork("test", forkId.Name(), sourceId.Name())

	modelWithComment := model.PostgresFork("test", forkId.Name(), sourceId.Name()).
		WithComment(comment)

	ref := modelBasic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.PostgresForkResource(t, ref).
			HasNameString(forkId.Name()).
			HasForkFromString(sourceId.Name()).
			HasNoComment().
			HasFullyQualifiedNameString(forkId.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(forkId.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"),
		resourceshowoutputassert.PostgresInstanceDescribeOutput(t, ref).
			HasName(forkId.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasStorageSizeGb(10).
			HasHighAvailability(false).
			HasNoComment(),
	}

	assertWithComment := []assert.TestCheckFuncProvider{
		resourceassert.PostgresForkResource(t, ref).
			HasNameString(forkId.Name()).
			HasForkFromString(sourceId.Name()).
			HasCommentString(comment).
			HasFullyQualifiedNameString(forkId.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(forkId.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(comment),
		resourceshowoutputassert.PostgresInstanceDescribeOutput(t, ref).
			HasName(forkId.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasStorageSizeGb(10).
			HasHighAvailability(false).
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
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, assertBasic...),
			},
			// Step 2: Import
			{
				Config:            accconfig.FromModels(t, modelBasic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"at",
					"before",
					"fork_from",
				},
			},
			// Step 3: Update — set comment
			{
				Config: accconfig.FromModels(t, modelWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, assertWithComment...),
			},
			// Step 4: External drift — alter comment externally, Terraform reverts to configured value
			{
				PreConfig: func() {
					testClient().PostgresInstance.Alter(t, sdk.NewAlterPostgresInstanceRequest(forkId).WithSet(
						*sdk.NewPostgresInstanceSetRequest().
							WithComment(externalComment)))
				},
				Config: accconfig.FromModels(t, modelWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
					},
				},
				Check: assertThat(t, assertWithComment...),
			},
			// Step 5: Unset — back to basic model (no comment)
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, assertBasic...),
			},
		},
	})
}

func TestAcc_PostgresFork_CompleteUseCase(t *testing.T) {
	sourceId := createSourceForFork(t)
	forkId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	// Use a recent timestamp (1 hour ago) — Snowflake requires fork timestamps within the past 10 days
	recentTime := time.Now().Add(-1 * time.Hour).UTC()

	modelComplete := model.PostgresFork("test", forkId.Name(), sourceId.Name()).
		WithComment(comment).
		WithAtTimestamp(recentTime).
		WithComputeFamily("STANDARD_M").
		WithStorageSizeGb(20).
		WithHighAvailability("false").
		WithPostgresSettings(`{"work_mem": "64MB"}`)

	ref := modelComplete.ResourceReference()

	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.PostgresForkResource(t, ref).
			HasNameString(forkId.Name()).
			HasForkFromString(sourceId.Name()).
			HasCommentString(comment).
			HasComputeFamilyString("STANDARD_M").
			HasStorageSizeGbString("20").
			HasHighAvailabilityString("false").
			HasPostgresSettingsString(`{"work_mem": "64MB"}`).
			HasFullyQualifiedNameString(forkId.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(forkId.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasComment(comment),
		resourceshowoutputassert.PostgresInstanceDescribeOutput(t, ref).
			HasName(forkId.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasStorageSizeGb(20).
			HasHighAvailability(false).
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
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					assertThat(t, assertComplete...),
					resource.TestCheckResourceAttr(ref, "at.0.timestamp", recentTime.Format("2006-01-02 15:04:05")),
				),
			},
			// Step 2: Import
			{
				Config:            accconfig.FromModels(t, modelComplete),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"at",
					"before",
					"fork_from",
				},
			},
		},
	})
}

func TestAcc_PostgresFork_Validations(t *testing.T) {
	sourceId := testClient().Ids.RandomAccountObjectIdentifier()
	forkId := testClient().Ids.RandomAccountObjectIdentifier()

	// Model with both at and before set — should fail validation
	modelConflict := model.PostgresFork("test", forkId.Name(), sourceId.Name()).
		WithAtTimestamp(time.Now().UTC().Add(-time.Hour)).
		WithBeforeTimestamp(time.Now().UTC().Add(-2 * time.Hour))

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
				ExpectError: regexp.MustCompile(`"at": conflicts with before`),
			},
		},
	})
}

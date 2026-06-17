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

func TestAcc_PostgresInstance_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	modelBasic := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 10).
		WithHighAvailability(false)

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.PostgresInstanceResource(t, modelBasic.ResourceReference()).
			HasNameString(id.Name()).
			HasComputeFamilyString("STANDARD_M").
			HasStorageSizeGbString("10").
			HasAuthenticationAuthorityString("POSTGRES").
			HasNoComment().
			HasNoNetworkPolicy().
			HasNoStorageIntegration().
			HasNoPostgresSettings().
			HasNoMaintenanceWindowStart().
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, modelBasic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasAuthenticationAuthority("POSTGRES"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresInstance),
		Steps: []resource.TestStep{
			// Create with only required params
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
					"show_output.0.updated_on",
					"describe_output.0.updated_on",
				},
			},
			// Destroy
			{
				Destroy: true,
				Config:  accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					invokeactionassert.PostgresInstanceDoesNotExist(t, id),
				),
			},
		},
	})
}

func TestAcc_PostgresInstance_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	commentUpdated := random.Comment()

	modelComplete := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 10).
		WithComment(comment).
		WithHighAvailability(false).
		WithPostgresSettings(`{"work_mem": "64MB"}`)

	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.PostgresInstanceResource(t, modelComplete.ResourceReference()).
			HasNameString(id.Name()).
			HasComputeFamilyString("STANDARD_M").
			HasStorageSizeGbString("10").
			HasAuthenticationAuthorityString("POSTGRES").
			HasCommentString(comment).
			HasHighAvailabilityString("false").
			HasPostgresSettingsString(`{"work_mem": "64MB"}`).
			HasNoNetworkPolicy().
			HasNoStorageIntegration().
			HasNoMaintenanceWindowStart().
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, modelComplete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasAuthenticationAuthority("POSTGRES").
			HasComment(comment),
	}

	modelWithMaintenance := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 10).
		WithComment(commentUpdated).
		WithHighAvailability(false).
		WithPostgresSettings(`{"work_mem": "64MB"}`).
		WithMaintenanceWindowStart("10")

	assertWithMaintenance := []assert.TestCheckFuncProvider{
		resourceassert.PostgresInstanceResource(t, modelWithMaintenance.ResourceReference()).
			HasNameString(id.Name()).
			HasCommentString(commentUpdated).
			HasMaintenanceWindowStartString("10").
			HasPostgresSettingsString(`{"work_mem": "64MB"}`).
			HasHighAvailabilityString("false").
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, modelWithMaintenance.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(commentUpdated),
	}

	modelBasic := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 10).
		WithHighAvailability(false)

	assertAfterUnset := []assert.TestCheckFuncProvider{
		resourceassert.PostgresInstanceResource(t, modelBasic.ResourceReference()).
			HasNameString(id.Name()).
			HasComputeFamilyString("STANDARD_M").
			HasStorageSizeGbString("10").
			HasAuthenticationAuthorityString("POSTGRES").
			HasCommentString("").
			HasNoMaintenanceWindowStart().
			HasNoPostgresSettings().
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, modelBasic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresInstance),
		Steps: []resource.TestStep{
			// Step 1: Create with optional params
			{
				Config: accconfig.FromModels(t, modelComplete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, assertComplete...),
			},
			// Step 2: Import
			{
				Config:            accconfig.FromModels(t, modelComplete),
				ResourceName:      modelComplete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"show_output.0.updated_on",
					"describe_output.0.updated_on",
				},
			},
			// Step 3: Update - set maintenance_window_start (alter-only field) and change comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithMaintenance.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelWithMaintenance),
				Check:  assertThat(t, assertWithMaintenance...),
			},
			// Step 4: Update - unset comment (back to basic model without optionals)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertAfterUnset...),
			},
			// Step 5: External change detection
			{
				PreConfig: func() {
					testClient().PostgresInstance.Alter(t, sdk.NewAlterPostgresInstanceRequest(id).WithSet(
						*sdk.NewPostgresInstanceSetRequest().
							WithComment(comment)))
				},
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, assertAfterUnset...),
			},
		},
	})
}

func TestAcc_PostgresInstance_Rename(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()

	modelBasic := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 10).
		WithHighAvailability(false)
	modelWithChangedName := model.PostgresInstance("test", newId.Name(), "POSTGRES", "STANDARD_M", 10).
		WithHighAvailability(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresInstance),
		Steps: []resource.TestStep{
			// create object
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.PostgresInstanceResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasComputeFamilyString("STANDARD_M").
						HasStorageSizeGbString("10").
						HasAuthenticationAuthorityString("POSTGRES"),
					resourceshowoutputassert.PostgresInstanceShowOutput(t, modelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComputeFamily("STANDARD_M").
						HasAuthenticationAuthority("POSTGRES"),
				),
			},
			// rename object - Snowflake backend does not currently support
			// ALTER POSTGRES INSTANCE ... RENAME TO (returns SQL compilation error)
			{
				Config:      accconfig.FromModels(t, modelWithChangedName),
				ExpectError: regexp.MustCompile(`error renaming Postgres instance`),
			},
		},
	})
}

func TestAcc_PostgresInstance_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	modelInvalidStorageSizeGb := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 0)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresInstance),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelInvalidStorageSizeGb),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected storage_size_gb to be at least \(1\), got 0`),
			},
		},
	})
}

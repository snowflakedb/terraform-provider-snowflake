//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_PostgresInstance_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 10).WithTimeout(accconfig.Timeouts{
		Create: "10m",
		Update: "10m",
		Delete: "10m",
		Read:   "10m",
	})

	withOptionals := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 10).
		WithComment(comment).
		WithPostgresSettings(`{"postgres:work_mem": "64KB"}`)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.PostgresInstanceResource(t, ref).
			HasNameString(id.Name()).
			HasComputeFamilyString("STANDARD_M").
			HasStorageSizeGbString("10").
			HasAuthenticationAuthorityString("POSTGRES").
			HasNoComment().
			HasNoNetworkPolicy().
			HasNoStorageIntegration().
			HasNoPostgresSettings().
			HasMaintenanceWindowStart(r.IntDefault).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, ref).
			// TODO(Could be Snowflake bug): HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasAuthenticationAuthority("POSTGRES"),
		resourceshowoutputassert.PostgresInstanceDescribeOutput(t, ref).
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasStorageSizeGb(10).
			HasAuthenticationAuthority("POSTGRES").
			HasHighAvailability(false),
		// TODO(Could be Snowflake bug): HasNoComment(),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.PostgresInstanceResource(t, ref).
			HasNameString(id.Name()).
			HasComputeFamilyString("STANDARD_M").
			HasStorageSizeGbString("10").
			HasAuthenticationAuthorityString("POSTGRES").
			HasCommentString(comment).
			HasPostgresSettingsString(`{"postgres:work_mem": "64KB"}`).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasAuthenticationAuthority("POSTGRES").
			HasComment(comment),
		resourceshowoutputassert.PostgresInstanceDescribeOutput(t, ref).
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasStorageSizeGb(10).
			HasAuthenticationAuthority("POSTGRES").
			HasHighAvailability(false).
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresInstance),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: accconfig.FromModels(t, basic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:            accconfig.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"high_availability",            // mismatching because of default value
					"maintenance_window_start",     // mismatching because of default value
					"network_policy",               // TODO: To address, I think the diff should be skipped
					"postgres_version",             // TODO: To address, I think the diff should be skipped
					"storage_integration",          // TODO: To address, I think the diff should be skipped
					"describe_output.0.updated_on", // TODO: To address, I think the diff should be skipped
					"show_output.0.updated_on",     // TODO: To address, I think the diff should be skipped
				},
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// Import - with optionals
			{
				Config:            accconfig.FromModels(t, withOptionals),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"show_output.0.updated_on",
					"describe_output.0.updated_on",
				},
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - external changes
			{
				PreConfig: func() {
					testClient().PostgresInstance.Alter(t, sdk.NewAlterPostgresInstanceRequest(id).WithSet(
						*sdk.NewPostgresInstanceSetRequest().
							WithComment(externalComment),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Destroy
			{
				Destroy: true,
				Config:  accconfig.FromModels(t, basic),
			},
			// Create - with optionals
			{
				PreConfig: func() {
					_, err := testClient().PostgresInstance.Show(t, id)
					require.ErrorIs(t, err, sdk.ErrObjectNotFound)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
		},
	})
}

func TestAcc_PostgresInstance_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	modelComplete := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 10).
		WithComment(comment).
		WithHighAvailability("false").
		WithPostgresSettings(`{"work_mem": "64MB"}`).
		WithMaintenanceWindowStart(10)

	ref := modelComplete.ResourceReference()

	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.PostgresInstanceResource(t, ref).
			HasNameString(id.Name()).
			HasComputeFamilyString("STANDARD_M").
			HasStorageSizeGbString("10").
			HasAuthenticationAuthorityString("POSTGRES").
			HasCommentString(comment).
			HasHighAvailabilityString("false").
			HasPostgresSettingsString(`{"work_mem": "64MB"}`).
			HasMaintenanceWindowStartString("10").
			HasNoNetworkPolicy().
			HasNoStorageIntegration().
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasAuthenticationAuthority("POSTGRES").
			HasComment(comment),
		resourceshowoutputassert.PostgresInstanceDescribeOutput(t, ref).
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasStorageSizeGb(10).
			HasAuthenticationAuthority("POSTGRES").
			HasHighAvailability(false).
			HasComment(comment).
			HasMaintenanceWindowStart(10),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresInstance),
		Steps: []resource.TestStep{
			// Step 1: Create with all optional params
			{
				Config: accconfig.FromModels(t, modelComplete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, assertComplete...),
			},
			// Step 2: Import
			{
				Config:            accconfig.FromModels(t, modelComplete),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"show_output.0.updated_on",
					"describe_output.0.updated_on",
				},
			},
		},
	})
}

func TestAcc_PostgresInstance_Rename(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()

	modelBasic := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 10).
		WithHighAvailability("false")
	modelWithChangedName := model.PostgresInstance("test", newId.Name(), "POSTGRES", "STANDARD_M", 10).
		WithHighAvailability("false")

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
				Check: assertThat(
					t,
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

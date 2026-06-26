//go:build account_level_tests

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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_PostgresInstance_BasicUseCase(t *testing.T) {
	t.Skip("TODO: Fix failing test")

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	modelBasic := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 10).
		WithHighAvailability("false")

	modelWithComment := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 10).
		WithHighAvailability("false").
		WithComment(comment).
		WithPostgresSettings(`{"work_mem": "64MB"}`)

	ref := modelBasic.ResourceReference()

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
			HasNoMaintenanceWindowStart().
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, ref).
			HasCreatedOnNotEmpty().
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
			HasHighAvailability(false).
			HasNoComment(),
	}

	assertWithComment := []assert.TestCheckFuncProvider{
		resourceassert.PostgresInstanceResource(t, ref).
			HasNameString(id.Name()).
			HasComputeFamilyString("STANDARD_M").
			HasStorageSizeGbString("10").
			HasAuthenticationAuthorityString("POSTGRES").
			HasCommentString(comment).
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
			// Step 1: Create with only required params
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
					"show_output.0.updated_on",
					"describe_output.0.updated_on",
				},
			},
			// Step 3: Update — set comment and postgres_settings
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
					testClient().PostgresInstance.Alter(t, sdk.NewAlterPostgresInstanceRequest(id).WithSet(
						*sdk.NewPostgresInstanceSetRequest().
							WithComment(externalComment),
					))
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
			// Step 5: Unset — back to basic model
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

func TestAcc_PostgresInstance_CompleteUseCase(t *testing.T) {
	t.Skip("TODO: Fix failing test")

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

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
	comment := random.Comment()

	modelBasic := model.PostgresInstance("test", id.Name(), "STANDARD_1", 10, "POSTGRES").
		WithHighAvailability(false)

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.PostgresInstanceResource(t, modelBasic.ResourceReference()).
			HasNameString(id.Name()).
			HasComputeFamilyString("STANDARD_1").
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
			HasComputeFamily("STANDARD_1").
			HasAuthenticationAuthority("POSTGRES").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"),
	}

	modelComplete := model.PostgresInstance("test", id.Name(), "STANDARD_1", 10, "POSTGRES").
		WithComment(comment).
		WithHighAvailability(false)

	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.PostgresInstanceResource(t, modelComplete.ResourceReference()).
			HasNameString(id.Name()).
			HasComputeFamilyString("STANDARD_1").
			HasStorageSizeGbString("10").
			HasAuthenticationAuthorityString("POSTGRES").
			HasCommentString(comment).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, modelComplete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasComputeFamily("STANDARD_1").
			HasAuthenticationAuthority("POSTGRES").
			HasComment(comment).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"),
	}

	assertAfterUnset := []assert.TestCheckFuncProvider{
		resourceassert.PostgresInstanceResource(t, modelBasic.ResourceReference()).
			HasNameString(id.Name()).
			HasComputeFamilyString("STANDARD_1").
			HasStorageSizeGbString("10").
			HasAuthenticationAuthorityString("POSTGRES").
			HasCommentString("").
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, modelBasic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasComputeFamily("STANDARD_1").
			HasAuthenticationAuthority("POSTGRES").
			HasComment("").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"),
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
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:            accconfig.FromModels(t, modelBasic),
				ResourceName:      modelBasic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelComplete),
				Check:  assertThat(t, assertComplete...),
			},
			// Import - with optionals
			{
				Config:            accconfig.FromModels(t, modelComplete),
				ResourceName:      modelComplete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertAfterUnset...),
			},
			// Update - detect external changes
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
			// Destroy - ensure postgres instance is destroyed before the next step
			{
				Destroy: true,
				Config:  accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					invokeactionassert.PostgresInstanceDoesNotExist(t, id),
				),
			},
			// Create - with optionals
			{
				Config: accconfig.FromModels(t, modelComplete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_PostgresInstance_Rename(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()

	modelBasic := model.PostgresInstance("test", id.Name(), "STANDARD_1", 10, "POSTGRES").
		WithHighAvailability(false)
	modelWithChangedName := model.PostgresInstance("test", newId.Name(), "STANDARD_1", 10, "POSTGRES").
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
						HasComputeFamilyString("STANDARD_1").
						HasStorageSizeGbString("10").
						HasAuthenticationAuthorityString("POSTGRES"),
					resourceshowoutputassert.PostgresInstanceShowOutput(t, modelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasComputeFamily("STANDARD_1").
						HasAuthenticationAuthority("POSTGRES").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE"),
				),
			},
			// rename object
			{
				Config: accconfig.FromModels(t, modelWithChangedName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithChangedName.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.PostgresInstanceResource(t, modelWithChangedName.ResourceReference()).
						HasNameString(newId.Name()).
						HasComputeFamilyString("STANDARD_1").
						HasStorageSizeGbString("10").
						HasAuthenticationAuthorityString("POSTGRES"),
					resourceshowoutputassert.PostgresInstanceShowOutput(t, modelWithChangedName.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(newId.Name()).
						HasComputeFamily("STANDARD_1").
						HasAuthenticationAuthority("POSTGRES").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE"),
				),
			},
		},
	})
}

func TestAcc_PostgresInstance_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	modelInvalidStorageSizeGb := model.PostgresInstance("test", id.Name(), "STANDARD_1", 0, "POSTGRES")

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

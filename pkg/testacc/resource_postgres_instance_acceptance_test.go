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
	t.Skip("TODO(SNOW-3765941): Skipped until Alter retry logic is not fixed")

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	networkPolicy, networkPolicyCleanup := testClient().NetworkPolicy.CreateNetworkPolicyForPostgres(t, testClient().NetworkRule)
	t.Cleanup(networkPolicyCleanup)

	basic := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 18, 10).WithTimeout(accconfig.Timeouts{
		Create: "10m",
		Update: "10m",
		Delete: "10m",
		Read:   "10m",
	})

	withOptionals := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 18, 10).
		WithComment(comment).
		WithHighAvailability("false").
		WithNetworkPolicy(networkPolicy.Name).
		// TODO(SNOW-3580377): storage_integration requires POSTGRES_EXTERNAL_STORAGE type; no pre-created integration available.
		WithMaintenanceWindowStart(10).
		WithPostgresSettings(`{"postgres:work_mem":"64KB"}`)

	ref := basic.ResourceReference()

	// basicAssertions builds the check slice for the "basic" (no optionals) config.
	// After initial Create, optional fields are absent from state (HasNoX).
	// After Update that unsets optionals, TypeString fields carry "" due to SDK v2 planned-value semantics — use HasXEmpty in that case.
	basicAssertions := func(optionalsWereSet bool) []assert.TestCheckFuncProvider {
		postgresInstanceResourceAssert := resourceassert.PostgresInstanceResource(t, ref).
			HasNameString(id.Name()).
			HasComputeFamilyString("STANDARD_M").
			HasStorageSizeGbString("10").
			HasAuthenticationAuthorityString("POSTGRES").
			HasHighAvailability(r.BooleanDefault).
			HasNoStorageIntegration().
			HasPostgresVersion(18).
			HasMaintenanceWindowStart(r.IntDefault).
			HasFullyQualifiedNameString(id.FullyQualifiedName())
		if optionalsWereSet {
			postgresInstanceResourceAssert = postgresInstanceResourceAssert.HasCommentEmpty().
				HasNetworkPolicyEmpty().
				HasPostgresSettingsEmpty()
		} else {
			postgresInstanceResourceAssert = postgresInstanceResourceAssert.HasNoComment().
				HasNoNetworkPolicy().
				HasNoPostgresSettings()
		}
		return []assert.TestCheckFuncProvider{
			postgresInstanceResourceAssert,
			resourceshowoutputassert.PostgresInstanceShowOutput(t, ref).
				HasCreatedOnNotEmpty().
				HasName(id.Name()).
				HasOwner(snowflakeroles.Accountadmin.Name()).
				HasOwnerRoleType("ROLE").
				HasComputeFamily("STANDARD_M").
				HasAuthenticationAuthority("POSTGRES").
				HasStorageSize(10).
				HasIsHighlyAvailable(false).
				HasState(sdk.PostgresInstanceStateReady),
			resourceshowoutputassert.PostgresInstanceDescribeOutput(t, ref).
				HasCreatedOnNotEmpty().
				HasName(id.Name()).
				HasOwner(snowflakeroles.Accountadmin.Name()).
				HasOwnerRoleType("ROLE").
				HasComputeFamily("STANDARD_M").
				HasStorageSizeGb(10).
				HasAuthenticationAuthority("POSTGRES").
				HasHighAvailability(false).
				HasPostgresVersion(18).
				HasState("READY"),
		}
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.PostgresInstanceResource(t, ref).
			HasNameString(id.Name()).
			HasComputeFamilyString("STANDARD_M").
			HasStorageSizeGbString("10").
			HasAuthenticationAuthorityString("POSTGRES").
			HasCommentString(comment).
			HasHighAvailability("false").
			HasNetworkPolicy(networkPolicy.Name).
			HasNoStorageIntegration().
			HasPostgresVersion(18).
			HasMaintenanceWindowStart(10).
			HasPostgresSettingsString(`{"postgres:work_mem":"64KB"}`).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.PostgresInstanceShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasAuthenticationAuthority("POSTGRES").
			HasComment(comment).
			HasStorageSize(10).
			HasIsHighlyAvailable(false).
			HasState(sdk.PostgresInstanceStateReady),
		resourceshowoutputassert.PostgresInstanceDescribeOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComputeFamily("STANDARD_M").
			HasStorageSizeGb(10).
			HasAuthenticationAuthority("POSTGRES").
			HasHighAvailability(false).
			HasComment(comment).
			HasPostgresVersion(18).
			// Not asserted: Snowflake DESCRIBE has a propagation lag for network_policy,
			// maintenance_window_start, and postgres_settings immediately after ALTER SET —
			// DESCRIBE returns stale/empty values until propagation completes.
			HasState("READY"),
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
				Check: assertThat(t, basicAssertions(false)...),
			},
			// Import - without optionals
			{
				Config:            accconfig.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// BooleanDefault sentinel cannot be recovered from Snowflake's true/false value
					"high_availability",
					// Snowflake returns "{}" for postgres_settings when unset; import stores it
					// but the pre-import state has it as empty — same reason step 8 ignores this.
					"postgres_settings",
					"describe_output.0.updated_on", // changes between reads
					"show_output.0.updated_on",     // changes between reads
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
			// Update - unset optionals (run directly after set optionals to use its clean state;
			// import with optionals is tested on the freshly-created instance at the end)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions(true)...),
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
				Check:  assertThat(t, basicAssertions(true)...),
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
			// Import - with optionals (done on a freshly created instance to avoid stale SHOW state
			// left over from the earlier SET POSTGRES_SETTINGS / COMMENT operations)
			{
				Config:            accconfig.FromModels(t, withOptionals),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// DESCRIBE has propagation lag for these fields after CREATE/ALTER SET — the value
					// in state right after Apply may differ from what a later Import Read sees.
					"postgres_settings",
					"describe_output.0.maintenance_window_start",
					"describe_output.0.updated_on", // changes between reads
					"show_output.0.updated_on",     // changes between reads
				},
			},
		},
	})
}

func TestAcc_PostgresInstance_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	modelInvalidStorageSizeGb := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 18, 0)
	modelInvalidMaintenanceWindowStart := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 18, 10).WithMaintenanceWindowStart(24)
	modelInvalidPostgresVersion := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 0, 10)
	modelInvalidHighAvailability := model.PostgresInstance("test", id.Name(), "POSTGRES", "STANDARD_M", 18, 10).WithHighAvailability("invalid")
	modelInvalidComputeFamily := model.PostgresInstance("test", id.Name(), "POSTGRES", "INVALID", 18, 10)
	modelInvalidAuthenticationAuthority := model.PostgresInstance("test", id.Name(), "INVALID", "STANDARD_M", 18, 10)

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
			{
				Config:      accconfig.FromModels(t, modelInvalidMaintenanceWindowStart),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected maintenance_window_start to be in the range \(0 - 23\), got 24`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidPostgresVersion),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected postgres_version to be at least \(1\), got 0`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidHighAvailability),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected.*high_availability.*to be one of \["true" "false"\], got invalid`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidComputeFamily),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid postgres instance compute family: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidAuthenticationAuthority),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid postgres instance authentication authority: INVALID`),
			},
		},
	})
}

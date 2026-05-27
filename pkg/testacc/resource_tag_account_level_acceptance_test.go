//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tag_CompleteUseCase_OnConflict_Bcr2291(t *testing.T) {
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() != testenvs.SnowflakeNonProdEnvironment {
		t.Skip("Tag propagation tests are only supported in non-prod environments")
	}

	id := secondaryTestClient().Ids.RandomSchemaObjectIdentifier()

	// Ensure the bundle is always re-enabled after the test, even on failure.
	t.Cleanup(func() {
		secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2026_03")
	})

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)

	withPropagateOnly := model.TagBase("test", id).
		WithPropagateEnum(sdk.TagPropagationOnDependency)

	withPropagateAndOnConflict := model.TagBase("test", id).
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithOnConflictCustomValue("conflict_value")

	resourceRef := withPropagateAndOnConflict.ResourceReference()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Step 1: disable the 2026_03 bundle and create the tag with propagate + on_conflict.
			// Creating with `on_conflict` works regardless of the bundle; only SHOW TAGS is affected by BCR-2291.
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				PreConfig: func() {
					secondaryTestClient().BcrBundles.DisableBcrBundle(t, "2026_03")
				},
				Config: config.FromModels(t, providerModel, withPropagateAndOnConflict),
				Check: assertThat(t,
					resourceassert.TagResource(t, resourceRef).
						HasPropagateEnum(sdk.TagPropagationOnDependency).
						HasOnConflictCustomValue("conflict_value"),
				),
			},
			// Step 2: with the bundle disabled, externally change ON_CONFLICT. The column is
			// absent from SHOW TAGS, so Read cannot see the drift and the plan is a no-op.
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				PreConfig: func() {
					secondaryTestClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).
						WithSet(*sdk.NewTagSetRequest().WithPropagate(
							*sdk.NewTagPropagateRequest(sdk.TagPropagationOnDependency).
								WithOnConflict(sdk.TagOnConflict{CustomValue: sdk.String("other_value_no_bcr")}),
						)),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: config.FromModels(t, providerModel, withPropagateAndOnConflict),
				Check: assertThat(t,
					resourceassert.TagResource(t, resourceRef).
						HasOnConflictCustomValue("conflict_value"),
				),
			},
			// Step 3: enable the 2026_03 bundle. Read now observes the external value set in
			// step 2 and the plan reconciles the resource back to the configured value.
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				PreConfig: func() {
					secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2026_03")
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceRef, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withPropagateAndOnConflict),
				Check: assertThat(t,
					resourceassert.TagResource(t, resourceRef).
						HasOnConflictCustomValue("conflict_value"),
				),
			},
			// Step 4: with the bundle still enabled, externally set a different ON_CONFLICT while
			// the config drops the attribute. The drift is detected and the value is cleared.
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				PreConfig: func() {
					secondaryTestClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).
						WithSet(*sdk.NewTagSetRequest().WithPropagate(
							*sdk.NewTagPropagateRequest(sdk.TagPropagationOnDependency).
								WithOnConflict(sdk.TagOnConflict{CustomValue: sdk.String("external_value")}),
						)),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withPropagateOnly.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withPropagateOnly),
				Check: assertThat(t,
					resourceassert.TagResource(t, withPropagateOnly.ResourceReference()).
						HasOnConflictEmpty(),
				),
			},
		},
	})
}

func TestAcc_Tag_CompleteUseCase_OnConflictAllowedValuesSequence_Bcr2291(t *testing.T) {
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() != testenvs.SnowflakeNonProdEnvironment {
		t.Skip("Tag propagation tests are only supported in non-prod environments")
	}

	id := secondaryTestClient().Ids.RandomSchemaObjectIdentifier()

	// Ensure the bundle is always re-enabled after the test, even on failure.
	t.Cleanup(func() {
		secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2026_03")
	})

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)

	withAllowedValues := model.TagBase("test", id).
		WithOrderedAllowedValues("confidential", "internal", "public").
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithOnConflictAllowedValuesSequence()

	withPropagateOnly := model.TagBase("test", id).
		WithPropagateEnum(sdk.TagPropagationOnDependency)

	resourceRef := withAllowedValues.ResourceReference()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Step 1: with the bundle enabled, create a tag with allowed_values + allowed_values_sequence on_conflict.
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				PreConfig: func() {
					secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2026_03")
				},
				Config: config.FromModels(t, providerModel, withAllowedValues),
				Check: assertThat(t,
					resourceassert.TagResource(t, resourceRef).
						HasPropagateEnum(sdk.TagPropagationOnDependency).
						HasOnConflictAllowedValuesSequence(),
				),
			},
			// Step 2: drop on_conflict; the value is cleared.
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withPropagateOnly.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withPropagateOnly),
				Check: assertThat(t,
					resourceassert.TagResource(t, withPropagateOnly.ResourceReference()).
						HasOnConflictEmpty(),
				),
			},
		},
	})
}

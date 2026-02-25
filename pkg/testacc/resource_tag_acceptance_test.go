//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_Tag_BasicUseCase(t *testing.T) {
	maskingPolicy, maskingPolicyCleanup := testClient().MaskingPolicy.CreateMaskingPolicy(t)
	t.Cleanup(maskingPolicyCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	basic := model.TagBase("test", id)

	complete := model.TagBase("test", newId).
		WithComment(comment).
		WithAllowedValues("value1", "value2").
		WithMaskingPolicies(maskingPolicy.ID())

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.Tag(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment("").
			HasAllowedValues(),

		resourceassert.TagResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCommentString("").
			HasAllowedValuesEmpty().
			HasMaskingPoliciesEmpty(),

		resourceshowoutputassert.TagShowOutput(t, basic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment("").
			HasOwnerRoleType("ROLE").
			HasNoAllowedValues(),
	}

	assertComplete := []assert.TestCheckFuncProvider{
		objectassert.Tag(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment(comment).
			HasAllowedValuesUnordered("value1", "value2"),

		resourceassert.TagResource(t, complete.ResourceReference()).
			HasNameString(newId.Name()).
			HasFullyQualifiedNameString(newId.FullyQualifiedName()).
			HasDatabaseString(newId.DatabaseName()).
			HasSchemaString(newId.SchemaName()).
			HasCommentString(comment).
			HasAllowedValues("value1", "value2").
			HasMaskingPolicies(maskingPolicy.ID().FullyQualifiedName()),

		resourceshowoutputassert.TagShowOutput(t, complete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasAllowedValues("value1", "value2"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tagsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Tag),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      basic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			// Import - with optionals
			{
				Config:            config.FromModels(t, complete),
				ResourceName:      complete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - detect external changes
			{
				PreConfig: func() {
					testClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).WithSet(sdk.NewTagSetRequest().WithMaskingPolicies([]sdk.SchemaObjectIdentifier{maskingPolicy.ID()})))
					testClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).WithSet(sdk.NewTagSetRequest().WithComment(comment)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Destroy - ensure tag is destroyed before the next step
			{
				Destroy: true,
				Config:  config.FromModels(t, basic),
				Check: assertThat(t,
					invokeactionassert.TagDoesNotExist(t, id),
				),
			},
			// Create - with optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_Tag_CompleteUseCase_AllowedValuesOrdering(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	basic := model.TagBase("test", id).WithAllowedValues("foo", "", "bar")
	basicWithDifferentValues := model.TagBase("test", id).WithAllowedValues("", "bar", "foo")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tagsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Tag),
		Steps: []resource.TestStep{
			// Create - with allowed_values
			{
				Config: config.FromModels(t, basic),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasAllowedValuesUnordered("", "bar", "foo"),

					resourceassert.TagResource(t, basic.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasAllowedValues("", "bar", "foo"),

					resourceshowoutputassert.TagShowOutput(t, basic.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasAllowedValues("", "bar", "foo"),
				),
			},
			// Import - with allowed_values
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      basic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - with different ordering
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicWithDifferentValues.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, basicWithDifferentValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasAllowedValuesUnordered("", "bar", "foo"),

					resourceassert.TagResource(t, basicWithDifferentValues.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasAllowedValues("", "bar", "foo"),

					resourceshowoutputassert.TagShowOutput(t, basicWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasAllowedValues("", "bar", "foo"),
				),
			},
			// Import - with different ordering
			{
				Config:            config.FromModels(t, basicWithDifferentValues),
				ResourceName:      basicWithDifferentValues.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_Tag_migrateFromVersion_0_98_0(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	tagModel := model.TagBase("test", id).
		WithAllowedValuesValue(tfconfig.ListVariable(tfconfig.StringVariable("foo"), tfconfig.StringVariable("bar")))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.98.0"),
				Config:            tagV098(id),
				Check: assertThat(t, resourceassert.TagResource(t, tagModel.ResourceReference()).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttr(tagModel.ResourceReference(), "allowed_values.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(tagModel.ResourceReference(), "allowed_values.0", "bar")),
					assert.Check(resource.TestCheckResourceAttr(tagModel.ResourceReference(), "allowed_values.1", "foo")),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: tagsProviderFactory,
				Config:                   config.FromModels(t, tagModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(tagModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t, resourceassert.TagResource(t, tagModel.ResourceReference()).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttr(tagModel.ResourceReference(), "allowed_values.#", "2")),
					assert.Check(resource.TestCheckTypeSetElemAttr(tagModel.ResourceReference(), "allowed_values.*", "foo")),
					assert.Check(resource.TestCheckTypeSetElemAttr(tagModel.ResourceReference(), "allowed_values.*", "bar")),
				),
			},
		},
	})
}

func tagV098(id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_tag" "test" {
	database				= "%[1]s"
	schema				    = "%[2]s"
	name					= "%[3]s"
	allowed_values			= ["bar", "foo"]
}
`, id.DatabaseName(), id.SchemaName(), id.Name())
}

func TestAcc_Tag_NoAllowedValues_WithoutExperimentFlag(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	basic := model.TagBase("test", id)
	withNoAllowedValues := model.TagBase("test", id).WithNoAllowedValues(true)
	withAllowedValues := model.TagBase("test", id).WithAllowedValues("value1", "value2")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tagsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Tag),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, withNoAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesNil(), // After initial creation, allowed values are null in Snowflake

					resourceassert.TagResource(t, withNoAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty(),
				),
			},
			// Update - set no_allowed_values to true, but experiment is off so it should be ignored
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesNil(),

					resourceassert.TagResource(t, basic.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty(),
				),
			},
			// Update - set no_allowed_values to true, but experiment is off so it should be ignored
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withNoAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesNil(),

					resourceassert.TagResource(t, withNoAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty(),
				),
			},
			// Update - add allowed_values (old behavior: add only)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withAllowedValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesUnordered("value1", "value2"),

					resourceassert.TagResource(t, withAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValues("value1", "value2"),
				),
			},
			// Update - remove allowed_values (old behavior: drop, not unset)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesEmpty(),

					resourceassert.TagResource(t, basic.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty(),
				),
			},
		},
	})
}

func TestAcc_Tag_AllowedValues_WithExperimentFlag(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	providerModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.TagAllowedValuesBehaviorChanges)

	basic := model.TagBase("test", id)
	withAllowedValues := model.TagBase("test", id).WithAllowedValues("value1", "value2")
	withDifferentAllowedValues := model.TagBase("test", id).WithAllowedValues("value2", "value3")
	withNoAllowedValues := model.TagBase("test", id).WithNoAllowedValues(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tagsWithExperimentFlagProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Tag),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: config.FromModels(t, providerModel, basic),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesNil(),

					resourceassert.TagResource(t, basic.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty().
						HasNoAllowedValuesString("false"),

					resourceshowoutputassert.TagShowOutput(t, basic.ResourceReference()).
						HasName(id.Name()).
						HasNoAllowedValues(),
				),
			},
			// Import - without optionals
			{
				Config:            config.FromModels(t, providerModel, basic),
				ResourceName:      basic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - set allowed_values
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withAllowedValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesUnordered("value1", "value2"),

					resourceassert.TagResource(t, withAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValues("value1", "value2").
						HasNoAllowedValuesString("false"),

					resourceshowoutputassert.TagShowOutput(t, withAllowedValues.ResourceReference()).
						HasName(id.Name()).
						HasAllowedValues("value1", "value2"),
				),
			},
			// Import - with allowed_values
			{
				Config:            config.FromModels(t, providerModel, withAllowedValues),
				ResourceName:      withAllowedValues.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - change allowed_values (partial overlap)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withDifferentAllowedValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withDifferentAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesUnordered("value2", "value3"),

					resourceassert.TagResource(t, withDifferentAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValues("value2", "value3").
						HasNoAllowedValuesString("false"),
				),
			},
			// Update - remove all allowed_values (uses UNSET; tag allows any value again)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, basic),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesNil(),

					resourceassert.TagResource(t, basic.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty().
						HasNoAllowedValuesString("false"),
				),
			},
			// Update - set no_allowed_values (tag blocks any value; empty non-nil AllowedValues in Snowflake)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withNoAllowedValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withNoAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesEmpty(),

					resourceassert.TagResource(t, withNoAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty().
						HasNoAllowedValuesString("true"),
				),
			},
			// Detect external change - someone adds allowed values externally
			{
				PreConfig: func() {
					testClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).WithAdd([]string{"external_value"}))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withNoAllowedValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withNoAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesEmpty(),

					resourceassert.TagResource(t, withNoAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty().
						HasNoAllowedValuesString("true"),
				),
			},
			// Destroy - ensure tag is destroyed before next step
			{
				Destroy: true,
				Config:  config.FromModels(t, providerModel, withNoAllowedValues),
			},
			// Create - with no_allowed_values (tag blocks any value from the start)
			{
				PreConfig: func() {
					// Assert tag does not exist before proceeding (check is here to avoid stale plan errors in the previous step)
					_, err := testClient().Tag.Show(t, id)
					require.ErrorIs(t, err, sdk.ErrObjectNotFound)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withNoAllowedValues.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, providerModel, withNoAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesEmpty(),

					resourceassert.TagResource(t, withNoAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty().
						HasNoAllowedValuesString("true"),
				),
			},
			// Destroy - ensure tag is destroyed before next step
			{
				Destroy: true,
				Config:  config.FromModels(t, providerModel, withNoAllowedValues),
			},
			// Create - with allowed_values
			{
				PreConfig: func() {
					// Assert tag does not exist before proceeding (check is here to avoid stale plan errors in the previous step)
					_, err := testClient().Tag.Show(t, id)
					require.ErrorIs(t, err, sdk.ErrObjectNotFound)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withAllowedValues.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, providerModel, withAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesUnordered("value1", "value2"),

					resourceassert.TagResource(t, withAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValues("value1", "value2").
						HasNoAllowedValuesString("false"),
				),
			},
			// Update - transition from allowed_values to no_allowed_values
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withNoAllowedValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withNoAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesEmpty(),

					resourceassert.TagResource(t, withNoAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty().
						HasNoAllowedValuesString("true"),
				),
			},
			// Detect external change - someone UNSET allowed values externally (null <-> empty transition)
			{
				PreConfig: func() {
					testClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).WithUnset(sdk.NewTagUnsetRequest().WithAllowedValues(true)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withNoAllowedValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withNoAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesEmpty(),

					resourceassert.TagResource(t, withNoAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty().
						HasNoAllowedValuesString("true"),
				),
			},
		},
	})
}

func TestAcc_Tag_TransitionToExperimentFlag_NullAllowedValues(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	providerModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.TagAllowedValuesBehaviorChanges)
	basic := model.TagBase("test", id)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create with external provider v2.14.0 - no allowed_values (null in Snowflake)
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.14.0"),
				Config:            config.FromModels(t, basic),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesNil(),

					resourceassert.TagResource(t, basic.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty(),
				),
			},
			// Transition to dev provider with experiment flag - noop
			{
				ProtoV6ProviderFactories: tagsWithExperimentFlagProviderFactory,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, providerModel, basic),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesNil(),

					resourceassert.TagResource(t, basic.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty().
						HasNoAllowedValuesString("false"),
				),
			},
		},
	})
}

func TestAcc_Tag_TransitionToExperimentFlag_EmptyAllowedValues(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	providerModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.TagAllowedValuesBehaviorChanges)
	basic := model.TagBase("test", id)
	withAllowedValues := model.TagBase("test", id).WithAllowedValues("v1", "v2")
	withNoAllowedValues := model.TagBase("test", id).WithNoAllowedValues(true)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create with external provider v2.14.0 - with allowed_values
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.14.0"),
				Config:            config.FromModels(t, withAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesUnordered("v1", "v2"),

					resourceassert.TagResource(t, withAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValues("v1", "v2"),
				),
			},
			// Remove allowed_values with external provider v2.14.0 - old behavior uses DROP, leaving empty allowed_values in Snowflake
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.14.0"),
				Config:            config.FromModels(t, basic),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesEmpty(),

					resourceassert.TagResource(t, basic.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty(),
				),
			},
			// Transition to dev provider with experiment flag - config uses no_allowed_values to match the empty Snowflake state
			{
				ProtoV6ProviderFactories: tagsWithExperimentFlagProviderFactory,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withNoAllowedValues.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, providerModel, withNoAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesEmpty(),

					resourceassert.TagResource(t, withNoAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValuesEmpty().
						HasNoAllowedValuesString("true"),
				),
			},
		},
	})
}

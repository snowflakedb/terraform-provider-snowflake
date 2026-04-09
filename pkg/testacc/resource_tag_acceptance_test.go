//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

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
		WithAllowedValues("value1", "value2", "FAIL").
		WithMaskingPolicies(maskingPolicy.ID()).
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithOnConflictCustomValue("FAIL")

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.Tag(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment("").
			HasAllowedValues().
			HasPropagateEnum(sdk.TagPropagationNone),

		resourceassert.TagResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCommentString("").
			HasAllowedValuesEmpty().
			HasMaskingPoliciesEmpty().
			HasPropagateEnum(sdk.TagPropagationNone).
			HasOnConflictEmpty(),

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
			HasAllowedValuesUnordered("value1", "value2", "FAIL").
			HasPropagateEnum(sdk.TagPropagationOnDependency),

		resourceassert.TagResource(t, complete.ResourceReference()).
			HasNameString(newId.Name()).
			HasFullyQualifiedNameString(newId.FullyQualifiedName()).
			HasDatabaseString(newId.DatabaseName()).
			HasSchemaString(newId.SchemaName()).
			HasCommentString(comment).
			HasAllowedValues("value1", "value2", "FAIL").
			HasMaskingPolicies(maskingPolicy.ID().FullyQualifiedName()).
			HasPropagateEnum(sdk.TagPropagationOnDependency).
			HasOnConflictCustomValue("FAIL"),

		resourceshowoutputassert.TagShowOutput(t, complete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasAllowedValues("FAIL", "value1", "value2"),
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
			// Import - with optionals (on_conflict is not readable from Snowflake)
			{
				Config:                  config.FromModels(t, complete),
				ResourceName:            complete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_conflict", "ordered_allowed_values", "allowed_values"},
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
					testClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).WithSet(*sdk.NewTagSetRequest().WithMaskingPolicies([]sdk.SchemaObjectIdentifier{maskingPolicy.ID()})))
					testClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).WithSet(*sdk.NewTagSetRequest().WithComment(comment)))
					testClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).WithSet(*sdk.NewTagSetRequest().WithPropagate(*sdk.NewTagPropagateRequest(sdk.TagPropagationOnDependency))))
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
			// Import - with allowed_values config (import populates ordered_allowed_values by default)
			{
				Config:                  config.FromModels(t, basic),
				ResourceName:            basic.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ordered_allowed_values", "allowed_values"},
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
			// Import - with different ordering (import populates ordered_allowed_values by default)
			{
				Config:                  config.FromModels(t, basicWithDifferentValues),
				ResourceName:            basicWithDifferentValues.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ordered_allowed_values", "allowed_values"},
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
			// Create with no_allowed_values=true and flag disabled - should error
			{
				Config:      config.FromModels(t, withNoAllowedValues),
				ExpectError: regexp.MustCompile("no_allowed_values is not supported"),
			},
			// Create basic tag (succeeds)
			{
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
			// Update - try to set no_allowed_values=true without flag - should error
			{
				Config:      config.FromModels(t, withNoAllowedValues),
				ExpectError: regexp.MustCompile("no_allowed_values is not supported"),
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

func TestAcc_Tag_ExternalChanges_WithoutExperimentFlag(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	allowedValues := []string{"value1", "value2"}
	basic := model.TagBase("test", id)
	withAllowedValues := model.TagBase("test", id).WithAllowedValues(allowedValues...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tagsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Tag),
		Steps: []resource.TestStep{
			// Create with allowed_values
			{
				Config: config.FromModels(t, withAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesUnordered(allowedValues...),

					resourceassert.TagResource(t, withAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValues(allowedValues...),
				),
			},
			// External change: drop all allowed values (tag enters blocking state).
			// Provider detects the drift (config still expects values) and restores them.
			{
				PreConfig: func() {
					testClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).WithDrop(allowedValues))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withAllowedValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasName(id.Name()).
						HasAllowedValuesUnordered(allowedValues...),

					resourceassert.TagResource(t, withAllowedValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAllowedValues(allowedValues...),
				),
			},
			// Destroy and recreate as basic tag to test the empty-config scenario
			{
				Destroy: true,
				Config:  config.FromModels(t, withAllowedValues),
			},
			{
				PreConfig: func() {
					_, err := testClient().Tag.Show(t, id)
					require.ErrorIs(t, err, sdk.ErrObjectNotFound)
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
			// External change: add and drop a temp value to make allowed_values empty (blocking).
			// Without the experiment flag, the provider cannot distinguish empty [] from null,
			// so no drift is detected and the tag silently stays in blocking state.
			{
				PreConfig: func() {
					testClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).WithAdd([]string{"temp_value"}))
					testClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).WithDrop([]string{"temp_value"}))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionNoop),
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
		WithExperimentalFeaturesEnabled(experimentalfeatures.TagsAllowEmptyAllowedValues)

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
			// Import - with allowed_values config (import populates ordered_allowed_values by default)
			{
				Config:                  config.FromModels(t, providerModel, withAllowedValues),
				ResourceName:            withAllowedValues.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ordered_allowed_values", "allowed_values"},
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
					testClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).WithUnset(*sdk.NewTagUnsetRequest().WithAllowedValues(true)))
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
			// Update - transition from no_allowed_values back to allowing any value (UNSET allowed_values in Snowflake)
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
		},
	})
}

func TestAcc_Tag_TransitionToExperimentFlag_NullAllowedValues(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	providerModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.TagsAllowEmptyAllowedValues)
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
		WithExperimentalFeaturesEnabled(experimentalfeatures.TagsAllowEmptyAllowedValues)
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

func TestAcc_Tag_TransitionToExperimentFlag_EmptyAllowedValuesWithBasicConfig(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	providerModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.TagsAllowEmptyAllowedValues)
	basic := model.TagBase("test", id)
	withAllowedValues := model.TagBase("test", id).WithAllowedValues("v1", "v2")

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
			// Remove allowed_values with v2.14.0 - old behavior uses DROP, leaving empty allowed_values (blocking) in Snowflake
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
			// Transition to the current provider version with experiment flag - stayed on basic config.
			// The new provider reads no_allowed_values=true from Snowflake (empty state; no values are allowed),
			// but config says no_allowed_values=false, so an update is planned to UNSET (nil state; all values are allowed).
			{
				ProtoV6ProviderFactories: tagsWithExperimentFlagProviderFactory,
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
		},
	})
}

func TestAcc_Tag_NoAllowedValues_DisableExperimentFlag(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	providerModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.TagsAllowEmptyAllowedValues)

	basic := model.TagBase("test", id)
	withNoAllowedValues := model.TagBase("test", id).WithNoAllowedValues(true)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create with experiment flag enabled and no_allowed_values=true
			{
				ProtoV6ProviderFactories: tagsWithExperimentFlagProviderFactory,
				Config:                   config.FromModels(t, providerModel, withNoAllowedValues),
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
			// Update the experimantal flag to disabled and no_allowed_values=false without modifying the allowed_values
			// The tag should stay in the blocking state and not go back to accepting any value.
			{
				ProtoV6ProviderFactories: tagsProviderFactory,
				Config:                   config.FromModels(t, basic),
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

func TestAcc_Tag_PropagateWithAllowedValuesSequence(t *testing.T) {
	tagId := testClient().Ids.RandomSchemaObjectIdentifier()

	var viewId sdk.SchemaObjectIdentifier

	withSeqInitial := model.TagBase("test", tagId).
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithOrderedAllowedValues("confidential", "internal", "public").
		WithOnConflictAllowedValuesSequence()

	withSeqReordered := model.TagBase("test", tagId).
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithOrderedAllowedValues("public", "internal", "confidential").
		WithOnConflictAllowedValuesSequence()

	withSeqAddedValue := model.TagBase("test", tagId).
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithOrderedAllowedValues("restricted", "public", "internal", "confidential").
		WithOnConflictAllowedValuesSequence()

	withSeqRemovedValue := model.TagBase("test", tagId).
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithOrderedAllowedValues("restricted", "internal", "confidential").
		WithOnConflictAllowedValuesSequence()

	withCustomConflict := model.TagBase("test", tagId).
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithOrderedAllowedValues("restricted", "internal", "confidential").
		WithOnConflictCustomValue("confidential")

	withSeqAgain := model.TagBase("test", tagId).
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithOrderedAllowedValues("restricted", "internal", "confidential").
		WithOnConflictAllowedValuesSequence()

	basic := model.TagBase("test", tagId).
		WithAllowedValues("restricted", "internal", "confidential")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tagsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Tag),
		Steps: []resource.TestStep{
			// Create with allowed_values_sequence and set up conflicting dependency to assert propagated value
			{
				Config: config.FromModels(t, withSeqInitial),
				Check: assertThat(t,
					objectassert.Tag(t, tagId).
						HasPropagateEnum(sdk.TagPropagationOnDependency).
						HasAllowedValuesUnordered("confidential", "internal", "public"),
					resourceassert.TagResource(t, withSeqInitial.ResourceReference()).
						HasOnConflictAllowedValuesSequence().
						HasOrderedAllowedValues("confidential", "internal", "public"),
				),
			},
			// Set up conflicting dependency (tag must exist first) and verify propagated value
			{
				PreConfig: func() {
					view := testClient().Tag.SetupTagPropagationConflictOnView(t, tagId, "internal", "public")
					viewId = view.ID()
				},
				Config: config.FromModels(t, withSeqInitial),
				Check: assertThat(t,
					// "confidential" is first in sequence but not assigned to either table;
					// "internal" is second and assigned to table1 -> first match -> propagated value
					invokeactionassert.TagValueOnObject(t, tagId, func() sdk.ObjectIdentifier { return viewId }, sdk.ObjectTypeView, "internal"),
				),
			},
			// Import (on_conflict not readable from Snowflake)
			{
				Config:                  config.FromModels(t, withSeqInitial),
				ResourceName:            withSeqInitial.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_conflict"},
			},
			// Update allowed_values order ONLY -> "public" is now first match -> view should get "public"
			{
				Config: config.FromModels(t, withSeqReordered),
				Check: assertThat(t,
					resourceassert.TagResource(t, withSeqReordered.ResourceReference()).
						HasOnConflictAllowedValuesSequence().
						HasOrderedAllowedValues("public", "internal", "confidential"),
					invokeactionassert.TagValueOnObject(t, tagId, func() sdk.ObjectIdentifier { return viewId }, sdk.ObjectTypeView, "public"),
				),
			},
			// Add new value "restricted" at front -> "public" still first match among conflicting ("internal","public")
			{
				Config: config.FromModels(t, withSeqAddedValue),
				Check: assertThat(t,
					resourceassert.TagResource(t, withSeqAddedValue.ResourceReference()).
						HasOrderedAllowedValues("restricted", "public", "internal", "confidential"),
					invokeactionassert.TagValueOnObject(t, tagId, func() sdk.ObjectIdentifier { return viewId }, sdk.ObjectTypeView, "public"),
				),
			},
			// Remove "public" -> only "internal" remains among conflicting values -> propagated = "internal"
			{
				Config: config.FromModels(t, withSeqRemovedValue),
				Check: assertThat(t,
					resourceassert.TagResource(t, withSeqRemovedValue.ResourceReference()).
						HasOrderedAllowedValues("restricted", "internal", "confidential"),
					invokeactionassert.TagValueOnObject(t, tagId, func() sdk.ObjectIdentifier { return viewId }, sdk.ObjectTypeView, "internal"),
				),
			},
			// Switch to custom_value -> view should get "confidential"
			{
				Config: config.FromModels(t, withCustomConflict),
				Check: assertThat(t,
					resourceassert.TagResource(t, withCustomConflict.ResourceReference()).
						HasOnConflictCustomValue("confidential"),
					invokeactionassert.TagValueOnObject(t, tagId, func() sdk.ObjectIdentifier { return viewId }, sdk.ObjectTypeView, "confidential"),
				),
			},
			// Switch back to allowed_values_sequence -> "internal" is now the only conflicting value in sequence
			{
				Config: config.FromModels(t, withSeqAgain),
				Check: assertThat(t,
					resourceassert.TagResource(t, withSeqAgain.ResourceReference()).
						HasOnConflictAllowedValuesSequence().
						HasOrderedAllowedValues("restricted", "internal", "confidential"),
					invokeactionassert.TagValueOnObject(t, tagId, func() sdk.ObjectIdentifier { return viewId }, sdk.ObjectTypeView, "internal"),
				),
			},
			// Remove propagate and on_conflict -> propagated tag value remains (Snowflake does not remove it), but state is clean
			{
				Config: config.FromModels(t, basic),
				Check: assertThat(t,
					objectassert.Tag(t, tagId).
						HasPropagateEnum(sdk.TagPropagationNone),
					resourceassert.TagResource(t, basic.ResourceReference()).
						HasPropagateEnum(sdk.TagPropagationNone).
						HasOnConflictEmpty(),
					// Snowflake does not remove already-propagated tags when UNSET PROPAGATE is issued
					invokeactionassert.TagValueOnObject(t, tagId, func() sdk.ObjectIdentifier { return viewId }, sdk.ObjectTypeView, "internal"),
				),
			},
		},
	})
}

func TestAcc_Tag_Propagation_CustomOnConflictValue(t *testing.T) {
	tagId := testClient().Ids.RandomSchemaObjectIdentifier()

	var viewId sdk.SchemaObjectIdentifier

	withCustomFail := model.TagBase("test", tagId).
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithOnConflictCustomValue("FAIL")

	withCustomFailAndAllowedValues := model.TagBase("test", tagId).
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithOnConflictCustomValue("FAIL").
		WithAllowedValues("FAIL", "alpha", "beta")

	withCustomRestricted := model.TagBase("test", tagId).
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithOnConflictCustomValue("RESTRICTED").
		WithAllowedValues("RESTRICTED", "FAIL", "alpha", "beta")

	propagateOnly := model.TagBase("test", tagId).
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithAllowedValues("RESTRICTED", "FAIL", "alpha", "beta")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tagsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Tag),
		Steps: []resource.TestStep{
			// Create the tag with propagate + on_conflict (tag must exist first)
			{
				Config: config.FromModels(t, withCustomFail),
				Check: assertThat(t,
					resourceassert.TagResource(t, withCustomFail.ResourceReference()).
						HasPropagateEnum(sdk.TagPropagationOnDependency).
						HasOnConflictCustomValue("FAIL"),
				),
			},
			// PreConfig sets up dependent objects (tables with conflicting tags, view), then verify
			{
				PreConfig: func() {
					view := testClient().Tag.SetupTagPropagationConflictOnView(t, tagId, "alpha", "beta")
					viewId = view.ID()
				},
				Config: config.FromModels(t, withCustomFail),
				Check:  assertThat(t, invokeactionassert.TagValueOnObject(t, tagId, func() sdk.ObjectIdentifier { return viewId }, sdk.ObjectTypeView, "FAIL")),
			},
			// Add allowed_values while keeping same propagate + on_conflict
			{
				Config: config.FromModels(t, withCustomFailAndAllowedValues),
				Check: assertThat(t,
					resourceassert.TagResource(t, withCustomFailAndAllowedValues.ResourceReference()).
						HasOnConflictCustomValue("FAIL").
						HasAllowedValues("FAIL", "alpha", "beta"),
					invokeactionassert.TagValueOnObject(t, tagId, func() sdk.ObjectIdentifier { return viewId }, sdk.ObjectTypeView, "FAIL"),
				),
			},
			// Update custom_value to "RESTRICTED"
			{
				Config: config.FromModels(t, withCustomRestricted),
				Check: assertThat(t,
					resourceassert.TagResource(t, withCustomRestricted.ResourceReference()).
						HasOnConflictCustomValue("RESTRICTED"),
					invokeactionassert.TagValueOnObject(t, tagId, func() sdk.ObjectIdentifier { return viewId }, sdk.ObjectTypeView, "RESTRICTED"),
				),
			},
			// Remove on_conflict (keep propagate)
			{
				Config: config.FromModels(t, propagateOnly),
				Check: assertThat(t,
					resourceassert.TagResource(t, propagateOnly.ResourceReference()).
						HasOnConflictEmpty(),
					invokeactionassert.TagValueOnObject(t, tagId, func() sdk.ObjectIdentifier { return viewId }, sdk.ObjectTypeView, "CONFLICT"),
				),
			},
		},
	})
}

func TestAcc_Tag_OrderedAllowedValues_FieldTransitions(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	providerModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.TagsAllowEmptyAllowedValues)

	withAllowedValues := model.TagBase("test", id).
		WithAllowedValues("a", "b", "c")
	withOrdered := model.TagBase("test", id).
		WithOrderedAllowedValues("c", "b", "a")
	withOrderedModified := model.TagBase("test", id).
		WithOrderedAllowedValues("x", "b", "c")
	withNoAllowedValues := model.TagBase("test", id).
		WithNoAllowedValues(true)
	withAllowedValuesAfterBlocking := model.TagBase("test", id).
		WithAllowedValues("x", "b", "c")

	ref := withAllowedValues.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tagsWithExperimentFlagProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Tag),
		Steps: []resource.TestStep{
			// Create with allowed_values (unordered).
			{
				Config: config.FromModels(t, providerModel, withAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).HasAllowedValuesUnordered("a", "b", "c"),
					resourceassert.TagResource(t, ref).
						HasAllowedValues("a", "b", "c").
						HasOrderedAllowedValuesEmpty(),
				),
			},
			// Switch to ordered_allowed_values - triggers Update, verify Snowflake order.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withOrdered),
				Check: assertThat(t,
					objectassert.Tag(t, id).HasAllowedValues("c", "b", "a"),
					resourceassert.TagResource(t, ref).
						HasOrderedAllowedValues("c", "b", "a").
						HasAllowedValuesEmpty(),
				),
			},
			// Modify ordered values.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withOrderedModified),
				Check: assertThat(t,
					objectassert.Tag(t, id).HasAllowedValues("x", "b", "c"),
					resourceassert.TagResource(t, ref).HasOrderedAllowedValues("x", "b", "c"),
				),
			},
			// Switch from ordered_allowed_values to no_allowed_values - blocks all values.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withNoAllowedValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).HasAllowedValuesEmpty(),
					resourceassert.TagResource(t, ref).
						HasNoAllowedValuesString("true").
						HasAllowedValuesEmpty().
						HasOrderedAllowedValuesEmpty(),
				),
			},
			// Switch from no_allowed_values back to ordered_allowed_values.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withOrderedModified),
				Check: assertThat(t,
					objectassert.Tag(t, id).HasAllowedValues("x", "b", "c"),
					resourceassert.TagResource(t, ref).HasOrderedAllowedValues("x", "b", "c"),
				),
			},
			// Switch to allowed_values (unordered).
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, providerModel, withAllowedValuesAfterBlocking),
				Check: assertThat(t,
					objectassert.Tag(t, id).HasAllowedValuesUnordered("x", "b", "c"),
					resourceassert.TagResource(t, ref).
						HasAllowedValues("x", "b", "c").
						HasOrderedAllowedValuesEmpty(),
				),
			},
			// Import with ordered_allowed_values config.
			{
				Config:            config.FromModels(t, providerModel, withOrderedModified),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Apply after import - values already in ordered_allowed_values, no reconciliation needed.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, providerModel, withOrderedModified),
				Check: assertThat(t,
					objectassert.Tag(t, id).HasAllowedValues("x", "b", "c"),
					resourceassert.TagResource(t, ref).HasOrderedAllowedValues("x", "b", "c"),
				),
			},
			// Import with allowed_values config (import populates ordered_allowed_values by default).
			{
				Config:                  config.FromModels(t, providerModel, withAllowedValuesAfterBlocking),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ordered_allowed_values", "allowed_values"},
			},
		},
	})
}

func TestAcc_Tag_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	invalidPropagate := model.TagBase("test", id).
		WithPropagate("INVALID")

	onConflictWithoutPropagate := model.TagBase("test", id).
		WithOnConflictCustomValue("FAIL")

	allowedValuesWithOrderedAllowedValues := model.TagBase("test", id).
		WithAllowedValues("a", "b").
		WithOrderedAllowedValues("a", "b")

	orderedAllowedValuesWithNoAllowedValues := model.TagBase("test", id).
		WithOrderedAllowedValues("a", "b").
		WithNoAllowedValues(true)

	allowedValuesSequenceWithAllowedValues := model.TagBase("test", id).
		WithAllowedValues("a", "b").
		WithPropagateEnum(sdk.TagPropagationOnDependency).
		WithOnConflictAllowedValuesSequence()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tagsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Tag),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidPropagate),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid tag propagation value`),
			},
			{
				Config:      config.FromModels(t, onConflictWithoutPropagate),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("\"on_conflict\": all of `on_conflict,propagate` must be specified"),
			},
			// allowed_values conflicts with ordered_allowed_values
			{
				Config:      config.FromModels(t, allowedValuesWithOrderedAllowedValues),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"allowed_values": conflicts with ordered_allowed_values`),
			},
			// ordered_allowed_values conflicts with no_allowed_values
			{
				Config:      config.FromModels(t, orderedAllowedValuesWithNoAllowedValues),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"ordered_allowed_values": conflicts with no_allowed_values`),
			},
			// allowed_values_sequence requires ordered_allowed_values
			{
				Config:      config.FromModels(t, allowedValuesSequenceWithAllowedValues),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"on_conflict.0.allowed_values_sequence": all of .+must be specified`),
			},
		},
	})
}

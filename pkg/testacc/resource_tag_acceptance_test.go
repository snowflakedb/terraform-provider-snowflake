//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
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

// TestAcc_Tag_NullableAllowedValues tests the distinction between null (field not specified)
// and empty set (allowed_values = []) using d.GetRawConfig() in the tag resource.
//
// Snowflake semantics:
// - UNSET ALLOWED_VALUES: any string value is allowed on the tag (permissive, no constraint).
// - DROP all ALLOWED_VALUES (resulting in zero allowed values): no value can be set on the tag (restrictive).
//
// Terraform config mapping:
// - allowed_values not specified (null) → calls UNSET ALLOWED_VALUES → permissive
// - allowed_values = [] (empty set)    → calls DROP on all existing values → restrictive
// - allowed_values = ["v1", "v2"]      → calls ADD/DROP to reach desired set
//
// The test verifies this by attempting to set the tag on a table after each state change.
func TestAcc_Tag_NullableAllowedValues(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	// Create a table to use as a target for tag value assignment verification.
	table, tableCleanup := testClient().Table.Create(t)
	t.Cleanup(tableCleanup)

	basic := model.TagBase("test", id)
	withValues := model.TagBase("test", id).WithAllowedValues("value1", "value2")
	withEmpty := model.TagBase("test", id).WithAllowedValuesEmpty()

	// trySetTagOnTable attempts to set the tag with a given value on the test table.
	// Returns nil on success, error on failure.
	trySetTagOnTable := func(tagValue string) error {
		return testClient().Tag.TrySetOnObject(t, sdk.ObjectTypeTable, table.ID(), []sdk.TagAssociation{
			{Name: id, Value: tagValue},
		})
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tagsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Tag),
		// Transitions
		// Step 1: values -> empty
		// Step 2: empty -> values
		// Step 3: values -> null
		// Step 4: null -> empty
		// Step 5: empty -> null
		// Step 6: null -> values
		Steps: []resource.TestStep{
			// Step 1: Create with allowed_values = ["value1", "value2"]
			// Only "value1" and "value2" should be assignable.
			{
				Config: config.FromModels(t, withValues),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						objectassert.Tag(t, id).
							HasAllowedValuesUnordered("value1", "value2"),
						resourceassert.TagResource(t, withValues.ResourceReference()).
							HasAllowedValues("value1", "value2"),
					),
					func(_ *terraform.State) error {
						// Allowed value should succeed
						if err := trySetTagOnTable("value1"); err != nil {
							return fmt.Errorf("expected setting tag to allowed value 'value1' to succeed, got: %w", err)
						}
						// Non-allowed value should fail
						if err := trySetTagOnTable("other_value"); err == nil {
							return fmt.Errorf("expected setting tag to non-allowed value 'other_value' to fail, but it succeeded")
						}
						return nil
					},
				),
			},
			// Step 2: Update to allowed_values = [] (empty set, not null)
			// Uses DROP to remove all values → restrictive: no value can be assigned.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withEmpty.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withEmpty),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						objectassert.Tag(t, id).
							HasAllowedValues(),
						resourceassert.TagResource(t, withEmpty.ResourceReference()).
							HasAllowedValuesEmpty(),
						resourceshowoutputassert.TagShowOutput(t, withEmpty.ResourceReference()).
							HasNoAllowedValues(),
					),
					func(_ *terraform.State) error {
						// After DROP all allowed values, no value should be assignable.
						if err := trySetTagOnTable("value1"); err == nil {
							return fmt.Errorf("expected setting tag value to fail after DROP all allowed values (restrictive), but it succeeded")
						}
						if err := trySetTagOnTable("any_value"); err == nil {
							return fmt.Errorf("expected setting tag value to fail after DROP all allowed values (restrictive), but it succeeded")
						}
						return nil
					},
				),
			},
			// Step 3: Set allowed_values again to prepare for the null transition.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withValues),
				Check: assertThat(t,
					objectassert.Tag(t, id).
						HasAllowedValuesUnordered("value1", "value2"),
					resourceassert.TagResource(t, withValues.ResourceReference()).
						HasAllowedValues("value1", "value2"),
				),
			},
			// Step 4: Remove allowed_values from config entirely (null)
			// Uses UNSET ALLOWED_VALUES → permissive: any value can be assigned.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						objectassert.Tag(t, id).
							HasAllowedValues(),
						resourceassert.TagResource(t, basic.ResourceReference()).
							HasAllowedValuesEmpty(),
						resourceshowoutputassert.TagShowOutput(t, basic.ResourceReference()).
							HasNoAllowedValues(),
					),
					func(_ *terraform.State) error {
						// After UNSET ALLOWED_VALUES, any value should be assignable.
						if err := trySetTagOnTable("value1"); err != nil {
							return fmt.Errorf("expected setting tag to 'value1' to succeed after UNSET (permissive), got: %w", err)
						}
						if err := trySetTagOnTable("any_arbitrary_value"); err != nil {
							return fmt.Errorf("expected setting tag to 'any_arbitrary_value' to succeed after UNSET (permissive), got: %w", err)
						}
						return nil
					},
				),
			},
			// Step 4: null → empty (Noop)
			// SDKv2 limitation: after Read, both null and empty configs produce [] in state.
			// Terraform sees state=[] and config=[] → no diff → Noop.
			// The tag remains in its previous UNSET/permissive state.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withEmpty.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, withEmpty),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						objectassert.Tag(t, id).
							HasAllowedValues(),
						resourceassert.TagResource(t, withEmpty.ResourceReference()).
							HasAllowedValuesEmpty(),
					),
					func(_ *terraform.State) error {
						// Tag remains permissive (UNSET was not overridden).
						if err := trySetTagOnTable("value1"); err != nil {
							return fmt.Errorf("expected tag to remain permissive after null→empty noop, got: %w", err)
						}
						if err := trySetTagOnTable("any_value"); err != nil {
							return fmt.Errorf("expected tag to remain permissive after null→empty noop, got: %w", err)
						}
						return nil
					},
				),
			},
			// Step 5: empty → null (Noop)
			// Same SDKv2 limitation: both produce [] in state → no diff → Noop.
			// The tag remains permissive (unchanged from the UNSET in step 3).
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, basic),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						objectassert.Tag(t, id).
							HasAllowedValues(),
						resourceassert.TagResource(t, basic.ResourceReference()).
							HasAllowedValuesEmpty(),
					),
					func(_ *terraform.State) error {
						// Tag is still permissive.
						if err := trySetTagOnTable("value1"); err != nil {
							return fmt.Errorf("expected tag to remain permissive after empty→null noop, got: %w", err)
						}
						if err := trySetTagOnTable("any_arbitrary_value"); err != nil {
							return fmt.Errorf("expected tag to remain permissive after empty→null noop, got: %w", err)
						}
						return nil
					},
				),
			},
			// Step 6: null → values (Update)
			// This transition works because state=[] differs from config=["value1","value2"].
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withValues),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						objectassert.Tag(t, id).
							HasAllowedValuesUnordered("value1", "value2"),
						resourceassert.TagResource(t, withValues.ResourceReference()).
							HasAllowedValues("value1", "value2"),
					),
					func(_ *terraform.State) error {
						if err := trySetTagOnTable("value1"); err != nil {
							return fmt.Errorf("expected setting tag to allowed value 'value1' to succeed, got: %w", err)
						}
						if err := trySetTagOnTable("other_value"); err == nil {
							return fmt.Errorf("expected setting tag to non-allowed value 'other_value' to fail, but it succeeded")
						}
						return nil
					},
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

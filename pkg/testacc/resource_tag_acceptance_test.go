//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tag_BasicUseCase(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	maskingPolicy, maskingPolicyCleanup := testClient().MaskingPolicy.CreateMaskingPolicy(t)
	t.Cleanup(maskingPolicyCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	newId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	comment := random.Comment()

	basic := model.TagBase("test", id)

	complete := model.TagBase("test", newId).
		WithComment(comment).
		WithAllowedValues("value1", "value2").
		WithMaskingPolicies(maskingPolicy.ID())

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.Tag(t, id).
			HasName(id.Name()).
			HasDatabaseName(database.ID().Name()).
			HasSchemaName(schema.ID().Name()).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment("").
			HasAllowedValues(),

		resourceassert.TagResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(database.ID().Name()).
			HasSchemaString(schema.ID().Name()).
			HasCommentString("").
			HasAllowedValuesEmpty().
			HasMaskingPoliciesEmpty(),

		resourceshowoutputassert.TagShowOutput(t, basic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(database.ID().Name()).
			HasSchemaName(schema.ID().Name()).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment("").
			HasOwnerRoleType("ROLE").
			HasNoAllowedValues(),
	}

	assertComplete := []assert.TestCheckFuncProvider{
		objectassert.Tag(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(database.ID().Name()).
			HasSchemaName(schema.ID().Name()).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment(comment).
			HasAllowedValuesSet("value1", "value2"),

		resourceassert.TagResource(t, complete.ResourceReference()).
			HasNameString(newId.Name()).
			HasFullyQualifiedNameString(newId.FullyQualifiedName()).
			HasDatabaseString(database.ID().Name()).
			HasSchemaString(schema.ID().Name()).
			HasCommentString(comment).
			HasAllowedValues("value1", "value2").
			HasMaskingPolicies(maskingPolicy.ID().FullyQualifiedName()),

		resourceshowoutputassert.TagShowOutput(t, complete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(newId.Name()).
			HasDatabaseName(database.ID().Name()).
			HasSchemaName(schema.ID().Name()).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasAllowedValues("value1", "value2"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
					testClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).WithSet(
						sdk.NewTagSetRequest().WithComment(comment),
					))
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
					objectassert.TagDoesNotExist(t, id),
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
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	basic := model.TagBase("test", id).WithAllowedValues("foo", "", "bar")
	basicWithDifferentValues := model.TagBase("test", id).WithAllowedValues("", "bar", "foo")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
						HasDatabaseName(database.ID().Name()).
						HasSchemaName(schema.ID().Name()).
						HasAllowedValuesSet("", "bar", "foo"),

					resourceassert.TagResource(t, basic.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasAllowedValues("", "bar", "foo"),

					resourceshowoutputassert.TagShowOutput(t, basic.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(database.ID().Name()).
						HasSchemaName(schema.ID().Name()).
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
						HasDatabaseName(database.ID().Name()).
						HasSchemaName(schema.ID().Name()).
						HasAllowedValuesSet("", "bar", "foo"),

					resourceassert.TagResource(t, basicWithDifferentValues.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasAllowedValues("", "bar", "foo"),

					resourceshowoutputassert.TagShowOutput(t, basicWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(database.ID().Name()).
						HasSchemaName(schema.ID().Name()).
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
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

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
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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

//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FileFormatOrc_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	renamedId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	basicModel := model.FileFormatOrc("test", id.DatabaseName(), id.SchemaName(), id.Name())
	ref := basicModel.ResourceReference()

	renamedModel := model.FileFormatOrc("test", id.DatabaseName(), id.SchemaName(), renamedId.Name())

	completeModel := model.FileFormatOrc("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithTrimSpace("true").
		WithNullIf("NULL_A", "NULL_B").
		WithComment(comment)

	alteredModel := model.FileFormatOrc("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithTrimSpace("false").
		WithReplaceInvalidCharacters("true").
		WithNullIf("NULL_C").
		WithComment(externalComment)

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatOrcResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasTrimSpace(r.BooleanDefault).
			HasReplaceInvalidCharacters(r.BooleanDefault).
			HasNullIfEmpty(),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeOrc).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(""),
		resourceshowoutputassert.FileFormatOrcDescribeOutput(t, ref).
			HasId(id).
			HasTrimSpace(false).
			HasReplaceInvalidCharacters(false).
			HasNoNullIf(),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatOrcResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasTrimSpace(r.BooleanTrue).
			HasReplaceInvalidCharacters(r.BooleanDefault).
			HasNullIf("NULL_A", "NULL_B").
			HasCommentString(comment),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeOrc).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment),
		resourceshowoutputassert.FileFormatOrcDescribeOutput(t, ref).
			HasId(id).
			HasTrimSpace(true).
			HasReplaceInvalidCharacters(false).
			HasNullIf("NULL_A", "NULL_B"),
	}

	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatOrcResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasTrimSpace(r.BooleanFalse).
			HasReplaceInvalidCharacters(r.BooleanTrue).
			HasNullIf("NULL_C").
			HasCommentString(externalComment),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeOrc).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(externalComment),
		resourceshowoutputassert.FileFormatOrcDescribeOutput(t, ref).
			HasId(id).
			HasTrimSpace(false).
			HasReplaceInvalidCharacters(true).
			HasNullIf("NULL_C"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatOrc),
		Steps: []resource.TestStep{
			// create with only required fields
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basicModel),
				Check:  assertThat(t, basicAssertions...),
			},
			// import
			{
				Config:            config.FromModels(t, basicModel),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// set all optional fields
			{
				Config: config.FromModels(t, completeModel),
				Check:  assertThat(t, completeAssertions...),
			},
			// import with all attributes
			{
				Config:            config.FromModels(t, completeModel),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// alter all optional fields (non-recreating change)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, alteredModel),
				Check:  assertThat(t, alteredAssertions...),
			},
			// external non-type change is detected and corrected back to the config value with an update (non-recreating change)
			{
				PreConfig: func() {
					testClient().FileFormat.CreateOrcWithRequest(t, id, sdk.NewCreateOrcFileFormatRequest(id).WithOrReplace(true))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, alteredModel),
				Check:  assertThat(t, alteredAssertions...),
			},
			// external type change is detected (the object is recreated in Snowflake as a different file format type)
			{
				PreConfig: func() {
					testClient().FileFormat.CreateCsvWithRequest(t, id, sdk.NewCreateCsvFileFormatRequest(id).WithOrReplace(true))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, alteredModel),
				Check:  assertThat(t, alteredAssertions...),
			},
			// unset optional fields
			{
				Config: config.FromModels(t, basicModel),
				Check:  assertThat(t, basicAssertions...),
			},
			// rename
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, renamedModel),
				Check: assertThat(
					t,
					resourceassert.FileFormatOrcResource(t, ref).
						HasNameString(renamedId.Name()).
						HasFullyQualifiedNameString(renamedId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_FileFormatOrc_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	completeModel := model.FileFormatOrcWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name()).
		WithTrimSpace("true").
		WithNullIf("NULL_A", "NULL_B").
		WithComment(comment)
	modelWithReplaceInvalidCharacters := model.FileFormatOrcWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name()).
		WithReplaceInvalidCharacters("true")
	ref := completeModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatOrc),
		Steps: []resource.TestStep{
			// create with all fields set
			{
				Config: config.FromModels(t, completeModel),
				Check: assertThat(
					t,
					resourceassert.FileFormatOrcResource(t, ref).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasTrimSpace("true").
						HasReplaceInvalidCharacters("default").
						HasNullIf("NULL_A", "NULL_B").
						HasCommentString(comment),
					resourceshowoutputassert.FileFormatShowOutput(t, ref).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.FileFormatTypeOrc).
						HasComment(comment),
					resourceshowoutputassert.FileFormatOrcDescribeOutput(t, ref).
						HasId(id).
						HasTrimSpace(true).
						HasReplaceInvalidCharacters(false).
						HasNullIf("NULL_A", "NULL_B"),
				),
			},
			// import
			{
				Config:            config.FromModels(t, completeModel),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  config.FromModels(t, completeModel),
				Destroy: true,
			},
			{
				Config: config.FromModels(t, modelWithReplaceInvalidCharacters),
				Check: assertThat(
					t,
					resourceassert.FileFormatOrcResource(t, ref).
						HasReplaceInvalidCharacters("true"),
					resourceshowoutputassert.FileFormatOrcDescribeOutput(t, ref).
						HasReplaceInvalidCharacters(true),
				),
			},
		},
	})
}

func TestAcc_FileFormatOrc_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	invalidTrimSpace := model.FileFormatOrc("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithTrimSpace("invalid")
	invalidReplaceInvalidCharacters := model.FileFormatOrc("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithReplaceInvalidCharacters("invalid")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatOrc),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidTrimSpace),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected .*trim_space.* to be one of \["true" "false"\], got invalid`),
			},
			{
				Config:      config.FromModels(t, invalidReplaceInvalidCharacters),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected .*replace_invalid_characters.* to be one of \["true" "false"\], got invalid`),
			},
		},
	})
}

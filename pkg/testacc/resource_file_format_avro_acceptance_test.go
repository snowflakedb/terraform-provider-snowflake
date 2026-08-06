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
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// avroImportedResourceAssert asserts the state produced by ImportFileFormatAvro. trim_space and
// replace_invalid_characters always land on their BooleanDefault sentinel, regardless of the live
// Snowflake value, to avoid manifesting a config-independent default as a Terraform-managed value - this
// holds for every scenario. Additional chained assertions cover the remaining, scenario-specific attributes.
func avroImportedResourceAssert(t *testing.T, resourceId string) *resourceassert.FileFormatAvroResourceAssert {
	t.Helper()
	return resourceassert.ImportedFileFormatAvroResource(t, resourceId).
		HasTrimSpace(r.BooleanDefault).
		HasReplaceInvalidCharacters(r.BooleanDefault)
}

func TestAcc_FileFormatAvro_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	resourceId := resourcehelpers.EncodeResourceIdentifier(id)
	renamedId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	basicModel := model.FileFormatAvro("test", id.DatabaseName(), id.SchemaName(), id.Name())
	ref := basicModel.ResourceReference()

	renamedModel := model.FileFormatAvro("test", id.DatabaseName(), id.SchemaName(), renamedId.Name())

	completeModel := model.FileFormatAvro("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("GZIP").
		WithTrimSpace("true").
		WithReplaceInvalidCharacters("true").
		WithNullIf("NULL_A", "NULL_B").
		WithComment(comment)

	alteredModel := model.FileFormatAvro("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("ZSTD").
		WithTrimSpace("false").
		WithReplaceInvalidCharacters("false").
		WithNullIf("NULL_C").
		WithComment(externalComment)

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatAvroResource(t, ref).
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
			HasType(sdk.FileFormatTypeAvro).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(""),
		resourceshowoutputassert.FileFormatAvroDescribeOutput(t, ref).
			HasId(id).
			HasCompression("AUTO").
			HasTrimSpace(false).
			HasReplaceInvalidCharacters(false).
			HasNoNullIf(),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatAvroResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCompressionString(string(sdk.AvroCompressionGzip)).
			HasTrimSpace(r.BooleanTrue).
			HasReplaceInvalidCharacters(r.BooleanTrue).
			HasNullIf("NULL_A", "NULL_B").
			HasCommentString(comment),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeAvro).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment),
		resourceshowoutputassert.FileFormatAvroDescribeOutput(t, ref).
			HasId(id).
			HasCompression("GZIP").
			HasTrimSpace(true).
			HasReplaceInvalidCharacters(true).
			HasNullIf("NULL_A", "NULL_B"),
	}

	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatAvroResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCompressionString("ZSTD").
			HasTrimSpace(r.BooleanFalse).
			HasReplaceInvalidCharacters(r.BooleanFalse).
			HasNullIf("NULL_C").
			HasCommentString(externalComment),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeAvro).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(externalComment),
		resourceshowoutputassert.FileFormatAvroDescribeOutput(t, ref).
			HasId(id).
			HasCompression("ZSTD").
			HasTrimSpace(false).
			HasReplaceInvalidCharacters(false).
			HasNullIf("NULL_C"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatAvro),
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
				Config:       config.FromModels(t, basicModel),
				ResourceName: ref,
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					avroImportedResourceAssert(t, resourceId).
						HasCompression("AUTO").
						HasNullIfEmpty(),
				),
			},
			// set all optional fields
			{
				Config: config.FromModels(t, completeModel),
				Check:  assertThat(t, completeAssertions...),
			},
			// import with all attributes
			{
				Config:       config.FromModels(t, completeModel),
				ResourceName: ref,
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					avroImportedResourceAssert(t, resourceId).
						HasCompression("GZIP").
						HasNullIf("NULL_A", "NULL_B"),
				),
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
					testClient().FileFormat.CreateAvroWithRequest(t, id, sdk.NewCreateAvroFileFormatRequest(id).WithOrReplace(true))
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
					resourceassert.FileFormatAvroResource(t, ref).
						HasNameString(renamedId.Name()).
						HasFullyQualifiedNameString(renamedId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_FileFormatAvro_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	resourceId := resourcehelpers.EncodeResourceIdentifier(id)
	comment := random.Comment()

	completeModel := model.FileFormatAvroWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression(string(sdk.AvroCompressionGzip)).
		WithTrimSpace("true").
		WithReplaceInvalidCharacters("true").
		WithNullIf("NULL_A", "NULL_B").
		WithComment(comment)
	ref := completeModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatAvro),
		Steps: []resource.TestStep{
			// create with all fields set
			{
				Config: config.FromModels(t, completeModel),
				Check: assertThat(
					t,
					resourceassert.FileFormatAvroResource(t, ref).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCompressionString("GZIP").
						HasTrimSpace("true").
						HasReplaceInvalidCharacters("true").
						HasNullIf("NULL_A", "NULL_B").
						HasCommentString(comment),
					resourceshowoutputassert.FileFormatShowOutput(t, ref).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.FileFormatTypeAvro).
						HasComment(comment),
					resourceshowoutputassert.FileFormatAvroDescribeOutput(t, ref).
						HasId(id).
						HasCompression("GZIP").
						HasTrimSpace(true).
						HasReplaceInvalidCharacters(true).
						HasNullIf("NULL_A", "NULL_B"),
				),
			},
			{
				Config:       config.FromModels(t, completeModel),
				ResourceName: ref,
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					avroImportedResourceAssert(t, resourceId).
						HasCompression("GZIP").
						HasNullIf("NULL_A", "NULL_B"),
				),
			},
		},
	})
}

func TestAcc_FileFormatAvro_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	invalidCompression := model.FileFormatAvro("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("INVALID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatAvro),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidCompression),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid avro compression: INVALID`),
			},
		},
	})
}

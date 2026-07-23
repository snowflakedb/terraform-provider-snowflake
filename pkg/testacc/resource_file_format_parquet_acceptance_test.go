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

func TestAcc_FileFormatParquet_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	renamedId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	basicModel := model.FileFormatParquet("test", id.DatabaseName(), id.SchemaName(), id.Name())
	ref := basicModel.ResourceReference()

	renamedModel := model.FileFormatParquet("test", id.DatabaseName(), id.SchemaName(), renamedId.Name())

	completeModel := model.FileFormatParquet("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("LZO").
		WithBinaryAsText("false").
		WithUseLogicalType("false").
		WithTrimSpace("true").
		WithUseVectorizedScanner("true").
		WithNullIf("NULL_A", "NULL_B").
		WithComment(comment)

	alteredModel := model.FileFormatParquet("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("NONE").
		WithBinaryAsText("true").
		WithUseLogicalType("true").
		WithTrimSpace("false").
		WithUseVectorizedScanner("false").
		WithNullIf("NULL_C").
		WithComment(externalComment)

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatParquetResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasSnappyCompressionString(r.BooleanDefault).
			HasBinaryAsText(r.BooleanDefault).
			HasUseLogicalType(r.BooleanDefault).
			HasTrimSpace(r.BooleanDefault).
			HasUseVectorizedScanner(r.BooleanDefault).
			HasReplaceInvalidCharacters(r.BooleanDefault).
			HasNullIfEmpty(),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeParquet).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(""),
		resourceshowoutputassert.FileFormatParquetDescribeOutput(t, ref).
			HasId(id).
			HasCompression("AUTO").
			HasBinaryAsText(true).
			HasUseLogicalType(false).
			HasTrimSpace(false).
			HasUseVectorizedScanner(false).
			HasReplaceInvalidCharacters(false).
			HasNoNullIf(),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatParquetResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCompressionString("LZO").
			HasSnappyCompressionString(r.BooleanDefault).
			HasBinaryAsText(r.BooleanFalse).
			HasUseLogicalType(r.BooleanFalse).
			HasTrimSpace(r.BooleanTrue).
			HasUseVectorizedScanner(r.BooleanTrue).
			HasReplaceInvalidCharacters(r.BooleanDefault).
			HasNullIf("NULL_A", "NULL_B").
			HasCommentString(comment),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeParquet).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment),
		resourceshowoutputassert.FileFormatParquetDescribeOutput(t, ref).
			HasId(id).
			HasCompression("LZO").
			HasBinaryAsText(false).
			HasUseLogicalType(false).
			HasTrimSpace(true).
			HasUseVectorizedScanner(true).
			HasReplaceInvalidCharacters(false).
			HasNullIf("NULL_A", "NULL_B"),
	}

	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatParquetResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCompressionString("NONE").
			HasSnappyCompressionString(r.BooleanDefault).
			HasBinaryAsText(r.BooleanTrue).
			HasUseLogicalType(r.BooleanTrue).
			HasTrimSpace(r.BooleanFalse).
			HasUseVectorizedScanner(r.BooleanFalse).
			HasReplaceInvalidCharacters(r.BooleanDefault).
			HasNullIf("NULL_C").
			HasCommentString(externalComment),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeParquet).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(externalComment),
		resourceshowoutputassert.FileFormatParquetDescribeOutput(t, ref).
			HasId(id).
			HasCompression("NONE").
			HasBinaryAsText(true).
			HasUseLogicalType(true).
			HasTrimSpace(false).
			HasUseVectorizedScanner(false).
			HasReplaceInvalidCharacters(false).
			HasNullIf("NULL_C"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatParquet),
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
					testClient().FileFormat.CreateParquetWithRequest(t, id, sdk.NewCreateParquetFileFormatRequest(id).WithOrReplace(true))
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
					resourceassert.FileFormatParquetResource(t, ref).
						HasNameString(renamedId.Name()).
						HasFullyQualifiedNameString(renamedId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_FileFormatParquet_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	completeModel := model.FileFormatParquetWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression(string(sdk.ParquetCompressionLzo)).
		WithBinaryAsText("false").
		WithUseLogicalType("true").
		WithTrimSpace("true").
		WithUseVectorizedScanner("true").
		WithNullIf("NULL_A", "NULL_B").
		WithComment(comment)
	// Compression and SnappyCompression are mutually exclusive; SnappyCompression is exercised separately here
	// because it cannot be read back from Snowflake (DESCRIBE FILE FORMAT folds it into COMPRESSION = SNAPPY).
	modelWithSnappyCompression := model.FileFormatParquetWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name()).
		WithSnappyCompression("true")
	modelWithReplaceInvalidCharacters := model.FileFormatParquetWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name()).
		WithReplaceInvalidCharacters("true")
	ref := completeModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatParquet),
		Steps: []resource.TestStep{
			// create with all fields set
			{
				Config: config.FromModels(t, completeModel),
				Check: assertThat(
					t,
					resourceassert.FileFormatParquetResource(t, ref).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCompressionString("LZO").
						HasBinaryAsText("false").
						HasUseLogicalType("true").
						HasTrimSpace("true").
						HasUseVectorizedScanner("true").
						HasReplaceInvalidCharacters("default").
						HasCommentString(comment),
					resourceshowoutputassert.FileFormatShowOutput(t, ref).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.FileFormatTypeParquet).
						HasComment(comment),
					resourceshowoutputassert.FileFormatParquetDescribeOutput(t, ref).
						HasId(id).
						HasCompression("LZO").
						HasBinaryAsText(false).
						HasUseLogicalType(true).
						HasTrimSpace(true).
						HasUseVectorizedScanner(true).
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
				Config: config.FromModels(t, modelWithSnappyCompression),
				Check: assertThat(
					t,
					resourceassert.FileFormatParquetResource(t, ref).
						HasSnappyCompressionString("true"),
					resourceshowoutputassert.FileFormatParquetDescribeOutput(t, ref).
						HasCompression("SNAPPY"),
				),
			},
			{
				Config:  config.FromModels(t, modelWithSnappyCompression),
				Destroy: true,
			},
			{
				Config: config.FromModels(t, modelWithReplaceInvalidCharacters),
				Check: assertThat(
					t,
					resourceassert.FileFormatParquetResource(t, ref).
						HasReplaceInvalidCharacters("true"),
					resourceshowoutputassert.FileFormatParquetDescribeOutput(t, ref).
						HasReplaceInvalidCharacters(true),
				),
			},
		},
	})
}

func TestAcc_FileFormatParquet_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	invalidCompression := model.FileFormatParquet("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("INVALID")
	compressionConflict := model.FileFormatParquet("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("LZO").
		WithSnappyCompression("true")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatParquet),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidCompression),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid parquet compression: INVALID`),
			},
			{
				Config:      config.FromModels(t, compressionConflict),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"compression": conflicts with snappy_compression|"snappy_compression": conflicts with compression`),
			},
		},
	})
}

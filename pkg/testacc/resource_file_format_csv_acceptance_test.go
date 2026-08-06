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

// csvImportedResourceAssert asserts the state produced by ImportFileFormatCsv. Every boolean attribute and
// skip_header always land on their BooleanDefault/IntDefault sentinel, regardless of the live Snowflake
// value, to avoid manifesting a config-independent default as a Terraform-managed value - this holds for
// every scenario. Additional chained assertions cover the remaining, scenario-specific attributes.
func csvImportedResourceAssert(t *testing.T, resourceId string) *resourceassert.FileFormatCsvResourceAssert {
	t.Helper()
	return resourceassert.ImportedFileFormatCsvResource(t, resourceId).
		HasParseHeader(r.BooleanDefault).
		HasTrimSpace(r.BooleanDefault).
		HasErrorOnColumnCountMismatch(r.BooleanDefault).
		HasSkipBlankLines(r.BooleanDefault).
		HasReplaceInvalidCharacters(r.BooleanDefault).
		HasEmptyFieldAsNull(r.BooleanDefault).
		HasSkipByteOrderMark(r.BooleanDefault).
		HasMultiLine(r.BooleanDefault).
		HasSkipHeader(r.IntDefault)
}

func TestAcc_FileFormatCsv_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	resourceId := resourcehelpers.EncodeResourceIdentifier(id)
	renamedId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	basicModel := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name())
	ref := basicModel.ResourceReference()

	renamedModel := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), renamedId.Name())

	completeModel := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression(string(sdk.CsvCompressionGzip)).
		WithRecordDelimiter(";").
		WithFieldDelimiter("|").
		WithMultiLine("false").
		WithFileExtension(".csv").
		// parse_header conflicts with skip_header
		WithSkipHeader(2).
		WithSkipBlankLines("true").
		WithDateFormat("YYYY-MM-DD").
		WithTimeFormat("HH24:MI:SS").
		WithTimestampFormat("YYYY-MM-DD HH24:MI:SS.FF3").
		WithBinaryFormat(string(sdk.BinaryFormatBase64)).
		WithEscape("NONE").
		WithEscapeUnenclosedField("NONE").
		WithTrimSpace("true").
		WithFieldOptionallyEnclosedBy(`"`).
		WithNullIf("NULL_A", "NULL_B").
		WithErrorOnColumnCountMismatch("false").
		WithReplaceInvalidCharacters("true").
		WithEmptyFieldAsNull("false").
		WithSkipByteOrderMark("false").
		WithEncoding(string(sdk.CsvEncodingIso88591)).
		WithComment(comment)

	alteredModel := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression(string(sdk.CsvCompressionBz2)).
		WithRecordDelimiter("#").
		WithFieldDelimiter(":").
		WithMultiLine("true").
		WithFileExtension(".tsv").
		// parse_header conflicts with skip_header, so only the latter is altered here
		WithSkipHeader(1).
		WithSkipBlankLines("false").
		WithDateFormat("MM-DD-YYYY").
		WithTimeFormat("HH24:MI").
		WithTimestampFormat("YYYY-MM-DD HH24:MI:SS.FF6").
		WithBinaryFormat(string(sdk.BinaryFormatUtf8)).
		WithEscape("\\").
		WithEscapeUnenclosedField("\\").
		WithTrimSpace("false").
		WithFieldOptionallyEnclosedBy("'").
		WithNullIf("NULL_C").
		WithErrorOnColumnCountMismatch("true").
		WithReplaceInvalidCharacters("false").
		WithEmptyFieldAsNull("true").
		WithSkipByteOrderMark("true").
		WithEncoding(string(sdk.CsvEncodingUtf8)).
		WithComment(externalComment)

	// parse_header and skip_header conflict with each other, so switching between them is verified with minimal configs
	parseHeaderModel := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithParseHeader("true")
	skipHeaderModel := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithSkipHeader(2)

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatCsvResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasMultiLine(r.BooleanDefault).
			HasParseHeader(r.BooleanDefault).
			HasSkipHeader(r.IntDefault).
			HasSkipBlankLines(r.BooleanDefault).
			HasTrimSpace(r.BooleanDefault).
			HasErrorOnColumnCountMismatch(r.BooleanDefault).
			HasReplaceInvalidCharacters(r.BooleanDefault).
			HasEmptyFieldAsNull(r.BooleanDefault).
			HasSkipByteOrderMark(r.BooleanDefault).
			HasNullIfEmpty(),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeCsv).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(""),
		resourceshowoutputassert.FileFormatCsvDescribeOutput(t, ref).
			HasId(id).
			HasCompression(string(sdk.CsvCompressionAuto)).
			HasRecordDelimiter(`\n`).
			HasFieldDelimiter(",").
			HasFileExtension("").
			HasSkipHeader(0).
			HasParseHeader(false).
			HasSkipBlankLines(false).
			HasDateFormat("AUTO").
			HasTimeFormat("AUTO").
			HasTimestampFormat("AUTO").
			HasBinaryFormat(string(sdk.BinaryFormatHex)).
			HasEscape("NONE").
			HasEscapeUnenclosedField(`\\`).
			HasTrimSpace(false).
			HasFieldOptionallyEnclosedBy("NONE").
			HasNullIf(`\\N`).
			HasErrorOnColumnCountMismatch(true).
			HasValidateUtf8(true).
			HasReplaceInvalidCharacters(false).
			HasEmptyFieldAsNull(true).
			HasSkipByteOrderMark(true).
			HasEncoding(string(sdk.CsvEncodingUtf8)).
			HasMultiLine(true),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatCsvResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCompressionString(string(sdk.CsvCompressionGzip)).
			HasRecordDelimiterString(";").
			HasFieldDelimiterString("|").
			HasMultiLine(r.BooleanFalse).
			HasFileExtensionString(".csv").
			HasParseHeader(r.BooleanDefault).
			HasSkipHeader(2).
			HasSkipBlankLines(r.BooleanTrue).
			HasDateFormatString("YYYY-MM-DD").
			HasTimeFormatString("HH24:MI:SS").
			HasTimestampFormatString("YYYY-MM-DD HH24:MI:SS.FF3").
			HasBinaryFormatString(string(sdk.BinaryFormatBase64)).
			HasEscapeString("NONE").
			HasEscapeUnenclosedFieldString("NONE").
			HasTrimSpace(r.BooleanTrue).
			HasFieldOptionallyEnclosedByString(`"`).
			HasNullIf("NULL_A", "NULL_B").
			HasErrorOnColumnCountMismatch(r.BooleanFalse).
			HasReplaceInvalidCharacters(r.BooleanTrue).
			HasEmptyFieldAsNull(r.BooleanFalse).
			HasSkipByteOrderMark(r.BooleanFalse).
			HasEncodingString(string(sdk.CsvEncodingIso88591)).
			HasCommentString(comment),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeCsv).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment),
		resourceshowoutputassert.FileFormatCsvDescribeOutput(t, ref).
			HasId(id).
			HasCompression(string(sdk.CsvCompressionGzip)).
			HasRecordDelimiter(";").
			HasFieldDelimiter("|").
			HasFileExtension(".csv").
			HasSkipHeader(2).
			HasParseHeader(false).
			HasSkipBlankLines(true).
			HasDateFormat("YYYY-MM-DD").
			HasTimeFormat("HH24:MI:SS").
			HasTimestampFormat("YYYY-MM-DD HH24:MI:SS.FF3").
			HasBinaryFormat(string(sdk.BinaryFormatBase64)).
			HasEscape("NONE").
			HasEscapeUnenclosedField("NONE").
			HasTrimSpace(true).
			HasFieldOptionallyEnclosedBy(`\"`).
			HasNullIf("NULL_A", "NULL_B").
			HasErrorOnColumnCountMismatch(false).
			HasReplaceInvalidCharacters(true).
			HasEmptyFieldAsNull(false).
			HasSkipByteOrderMark(false).
			HasEncoding(string(sdk.CsvEncodingIso88591)).
			HasMultiLine(false),
	}

	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatCsvResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCompressionString(string(sdk.CsvCompressionBz2)).
			HasRecordDelimiterString("#").
			HasFieldDelimiterString(":").
			HasMultiLine(r.BooleanTrue).
			HasFileExtensionString(".tsv").
			HasParseHeader(r.BooleanDefault).
			HasSkipHeader(1).
			HasSkipBlankLines(r.BooleanFalse).
			HasDateFormatString("MM-DD-YYYY").
			HasTimeFormatString("HH24:MI").
			HasTimestampFormatString("YYYY-MM-DD HH24:MI:SS.FF6").
			HasBinaryFormatString(string(sdk.BinaryFormatUtf8)).
			HasEscapeString("\\").
			HasEscapeUnenclosedFieldString("\\").
			HasTrimSpace(r.BooleanFalse).
			HasFieldOptionallyEnclosedByString("'").
			HasNullIf("NULL_C").
			HasErrorOnColumnCountMismatch(r.BooleanTrue).
			HasReplaceInvalidCharacters(r.BooleanFalse).
			HasEmptyFieldAsNull(r.BooleanTrue).
			HasSkipByteOrderMark(r.BooleanTrue).
			HasEncodingString(string(sdk.CsvEncodingUtf8)).
			HasCommentString(externalComment),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeCsv).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(externalComment),
		resourceshowoutputassert.FileFormatCsvDescribeOutput(t, ref).
			HasId(id).
			HasCompression(string(sdk.CsvCompressionBz2)).
			HasRecordDelimiter("#").
			HasFieldDelimiter(":").
			HasFileExtension(".tsv").
			HasSkipHeader(1).
			HasParseHeader(false).
			HasSkipBlankLines(false).
			HasDateFormat("MM-DD-YYYY").
			HasTimeFormat("HH24:MI").
			HasTimestampFormat("YYYY-MM-DD HH24:MI:SS.FF6").
			HasBinaryFormat(string(sdk.BinaryFormatUtf8)).
			HasEscape(`\\`).
			HasEscapeUnenclosedField(`\\`).
			HasTrimSpace(false).
			HasFieldOptionallyEnclosedBy("'").
			HasNullIf("NULL_C").
			HasErrorOnColumnCountMismatch(true).
			HasReplaceInvalidCharacters(false).
			HasEmptyFieldAsNull(true).
			HasSkipByteOrderMark(true).
			HasEncoding(string(sdk.CsvEncodingUtf8)).
			HasMultiLine(true),
	}

	parseHeaderAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatCsvResource(t, ref).
			HasNameString(id.Name()).
			HasParseHeader(r.BooleanTrue).
			HasSkipHeader(r.IntDefault),
		resourceshowoutputassert.FileFormatCsvDescribeOutput(t, ref).
			HasId(id).
			HasParseHeader(true).
			HasSkipHeader(0),
	}

	skipHeaderAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatCsvResource(t, ref).
			HasNameString(id.Name()).
			HasParseHeader(r.BooleanDefault).
			HasSkipHeader(2),
		resourceshowoutputassert.FileFormatCsvDescribeOutput(t, ref).
			HasId(id).
			HasParseHeader(false).
			HasSkipHeader(2),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatCsv),
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
					csvImportedResourceAssert(t, resourceId).
						HasCompression(string(sdk.CsvCompressionAuto)).
						HasRecordDelimiter(`\n`).
						HasFieldDelimiter(",").
						HasFileExtension("").
						HasDateFormat("AUTO").
						HasTimeFormat("AUTO").
						HasTimestampFormat("AUTO").
						HasBinaryFormat(string(sdk.BinaryFormatHex)).
						HasEscape("NONE").
						HasEscapeUnenclosedField(`\\`).
						HasFieldOptionallyEnclosedBy("NONE").
						HasEncoding(string(sdk.CsvEncodingUtf8)).
						HasNullIf(`\\N`),
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
					csvImportedResourceAssert(t, resourceId).
						HasCompression(string(sdk.CsvCompressionGzip)).
						HasRecordDelimiter(";").
						HasFieldDelimiter("|").
						HasFileExtension(".csv").
						HasDateFormat("YYYY-MM-DD").
						HasTimeFormat("HH24:MI:SS").
						HasTimestampFormat("YYYY-MM-DD HH24:MI:SS.FF3").
						HasBinaryFormat(string(sdk.BinaryFormatBase64)).
						HasEscape("NONE").
						HasEscapeUnenclosedField("NONE").
						HasFieldOptionallyEnclosedBy(`\"`).
						HasEncoding(string(sdk.CsvEncodingIso88591)).
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
					testClient().FileFormat.CreateCsvWithRequest(t, id, sdk.NewCreateCsvFileFormatRequest(id).WithOrReplace(true))
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
					testClient().FileFormat.CreateJsonWithRequest(t, id, sdk.NewCreateJsonFileFormatRequest(id).WithOrReplace(true))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, alteredModel),
				Check:  assertThat(t, alteredAssertions...),
			},
			// switch from skip_header to parse_header (non-recreating change)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, parseHeaderModel),
				Check:  assertThat(t, parseHeaderAssertions...),
			},
			// switch back from parse_header to skip_header (non-recreating change)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, skipHeaderModel),
				Check:  assertThat(t, skipHeaderAssertions...),
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
					resourceassert.FileFormatCsvResource(t, ref).
						HasNameString(renamedId.Name()).
						HasFullyQualifiedNameString(renamedId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_FileFormatCsv_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	resourceId := resourcehelpers.EncodeResourceIdentifier(id)
	comment := random.Comment()

	completeModel := model.FileFormatCsvWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression(string(sdk.CsvCompressionGzip)).
		WithRecordDelimiter("NONE").
		WithFieldDelimiter("NONE").
		WithMultiLine("false").
		WithFileExtension(".csv").
		// parse_header conflicts with skip_header
		WithParseHeader("true").
		WithSkipBlankLines("true").
		WithDateFormat("AUTO").
		WithTimeFormat("AUTO").
		WithTimestampFormat("AUTO").
		WithBinaryFormat(string(sdk.BinaryFormatBase64)).
		WithEscape("NONE").
		WithEscapeUnenclosedField("NONE").
		WithTrimSpace("true").
		WithFieldOptionallyEnclosedBy("NONE").
		WithNullIf("NULL").
		WithErrorOnColumnCountMismatch("false").
		WithReplaceInvalidCharacters("true").
		WithEmptyFieldAsNull("false").
		WithSkipByteOrderMark("false").
		WithEncoding(string(sdk.CsvEncodingUtf8)).
		WithComment(comment)
	ref := completeModel.ResourceReference()

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatCsvResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCompressionString(string(sdk.CsvCompressionGzip)).
			HasRecordDelimiterString("NONE").
			HasFieldDelimiterString("NONE").
			HasMultiLine(r.BooleanFalse).
			HasFileExtensionString(".csv").
			HasParseHeader(r.BooleanTrue).
			HasSkipHeader(r.IntDefault).
			HasSkipBlankLines(r.BooleanTrue).
			HasDateFormatString("AUTO").
			HasTimeFormatString("AUTO").
			HasTimestampFormatString("AUTO").
			HasBinaryFormatString(string(sdk.BinaryFormatBase64)).
			HasEscapeString("NONE").
			HasEscapeUnenclosedFieldString("NONE").
			HasTrimSpace(r.BooleanTrue).
			HasFieldOptionallyEnclosedByString("NONE").
			HasNullIf("NULL").
			HasErrorOnColumnCountMismatch(r.BooleanFalse).
			HasReplaceInvalidCharacters(r.BooleanTrue).
			HasEmptyFieldAsNull(r.BooleanFalse).
			HasSkipByteOrderMark(r.BooleanFalse).
			HasEncodingString(string(sdk.CsvEncodingUtf8)).
			HasCommentString(comment),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeCsv).
			HasComment(comment),
		resourceshowoutputassert.FileFormatCsvDescribeOutput(t, ref).
			HasId(id).
			HasCompression(string(sdk.CsvCompressionGzip)).
			HasRecordDelimiter("NONE").
			HasFieldDelimiter("NONE").
			HasFileExtension(".csv").
			HasParseHeader(true).
			HasSkipHeader(0).
			HasSkipBlankLines(true).
			HasDateFormat("AUTO").
			HasTimeFormat("AUTO").
			HasTimestampFormat("AUTO").
			HasBinaryFormat(string(sdk.BinaryFormatBase64)).
			HasEscape("NONE").
			HasEscapeUnenclosedField("NONE").
			HasTrimSpace(true).
			HasFieldOptionallyEnclosedBy("NONE").
			HasNullIf("NULL").
			HasErrorOnColumnCountMismatch(false).
			HasReplaceInvalidCharacters(true).
			HasEmptyFieldAsNull(false).
			HasSkipByteOrderMark(false).
			HasEncoding(string(sdk.CsvEncodingUtf8)).
			HasMultiLine(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatCsv),
		Steps: []resource.TestStep{
			// create with all fields set
			{
				Config: config.FromModels(t, completeModel),
				Check:  assertThat(t, completeAssertions...),
			},
			// import
			{
				Config:       config.FromModels(t, completeModel),
				ResourceName: ref,
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					csvImportedResourceAssert(t, resourceId).
						HasCompression(string(sdk.CsvCompressionGzip)).
						HasRecordDelimiter("NONE").
						HasFieldDelimiter("NONE").
						HasFileExtension(".csv").
						HasDateFormat("AUTO").
						HasTimeFormat("AUTO").
						HasTimestampFormat("AUTO").
						HasBinaryFormat(string(sdk.BinaryFormatBase64)).
						HasEscape("NONE").
						HasEscapeUnenclosedField("NONE").
						HasFieldOptionallyEnclosedBy("NONE").
						HasEncoding(string(sdk.CsvEncodingUtf8)).
						HasNullIf("NULL"),
				),
			},
		},
	})
}

func TestAcc_FileFormatCsv_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	invalidCompression := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("INVALID")
	invalidBinaryFormat := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithBinaryFormat("INVALID")
	invalidEncoding := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithEncoding("INVALID")
	invalidSkipHeader := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithSkipHeader(-2)
	conflictingHeaderFields := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithParseHeader("true").
		WithSkipHeader(1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatCsv),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidCompression),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid csv compression: INVALID`),
			},
			{
				Config:      config.FromModels(t, invalidBinaryFormat),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid binary format: INVALID`),
			},
			{
				Config:      config.FromModels(t, invalidEncoding),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid csv encoding: INVALID`),
			},
			{
				Config:      config.FromModels(t, invalidSkipHeader),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected skip_header to be at least \(0\)`),
			},
			{
				Config:      config.FromModels(t, conflictingHeaderFields),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`skip_header.*conflicts with\s+parse_header`),
			},
		},
	})
}

// proves https://github.com/snowflakedb/terraform-provider-snowflake/discussions/1950:
// an explicitly empty null_if list results in NULL_IF = () instead of Snowflake's default of (\\N).
func TestAcc_FileFormatCsv_EmptyNullIf(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	resourceId := resourcehelpers.EncodeResourceIdentifier(id)

	emptyModel := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithNullIf()
	nonEmptyModel := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithNullIf("NULL_A")
	ref := emptyModel.ResourceReference()

	emptyAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatCsvResource(t, ref).
			HasNameString(id.Name()).
			HasNullIfEmpty(),
		resourceshowoutputassert.FileFormatCsvDescribeOutput(t, ref).
			HasId(id).
			// an empty describe output means NULL_IF = () was set; without null_if in the config,
			// Snowflake's default is used instead and the describe output contains a single \\N entry
			HasNullIf(),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatCsv),
		Steps: []resource.TestStep{
			// create with an explicitly empty null_if
			{
				Config: config.FromModels(t, emptyModel),
				Check:  assertThat(t, emptyAssertions...),
			},
			// no drift
			{
				Config: config.FromModels(t, emptyModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply:             []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			// import
			{
				Config:       config.FromModels(t, emptyModel),
				ResourceName: ref,
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					csvImportedResourceAssert(t, resourceId).
						HasCompression(string(sdk.CsvCompressionAuto)).
						HasRecordDelimiter(`\n`).
						HasFieldDelimiter(",").
						HasFileExtension("").
						HasDateFormat("AUTO").
						HasTimeFormat("AUTO").
						HasTimestampFormat("AUTO").
						HasBinaryFormat(string(sdk.BinaryFormatHex)).
						HasEscape("NONE").
						HasEscapeUnenclosedField(`\\`).
						HasFieldOptionallyEnclosedBy("NONE").
						HasEncoding(string(sdk.CsvEncodingUtf8)).
						HasNullIfEmpty(),
				),
			},
			// change to a non-empty null_if
			{
				Config: config.FromModels(t, nonEmptyModel),
				Check: assertThat(
					t,
					resourceassert.FileFormatCsvResource(t, ref).
						HasNullIf("NULL_A"),
					resourceshowoutputassert.FileFormatCsvDescribeOutput(t, ref).
						HasNullIf("NULL_A"),
				),
			},
			// change back to an empty null_if
			{
				Config: config.FromModels(t, emptyModel),
				Check:  assertThat(t, emptyAssertions...),
			},
			// no drift after going back to an empty null_if
			{
				Config: config.FromModels(t, emptyModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply:             []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// proves https://github.com/snowflakedb/terraform-provider-snowflake/issues/3325:
// a null_if list containing an empty string results in NULL_IF with an empty string entry
// and does not cause a permanent plan.
// Note that DESCRIBE FILE FORMAT returns an empty list both for a null_if list with a single empty string
// and for an empty null_if list, so a lone empty string can not be read back from the describe output.
func TestAcc_FileFormatCsv_NullIfWithEmptyString(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	resourceId := resourcehelpers.EncodeResourceIdentifier(id)

	emptyStringModel := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithNullIf("")
	withEmptyStringModel := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithNullIf("NULL_A", "")
	ref := emptyStringModel.ResourceReference()

	emptyStringAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatCsvResource(t, ref).
			HasNameString(id.Name()).
			HasNullIf(""),
		resourceshowoutputassert.FileFormatCsvDescribeOutput(t, ref).
			HasId(id).
			// surprising, but correct: Snowflake reduces a single empty string to an empty list in the describe
			// output, so this is exactly what TestAcc_FileFormatCsv_EmptyNullIf asserts for an empty null_if
			HasNullIf(),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatCsv),
		Steps: []resource.TestStep{
			// create with null_if containing only an empty string
			{
				Config: config.FromModels(t, emptyStringModel),
				Check:  assertThat(t, emptyStringAssertions...),
			},
			// no drift
			{
				Config: config.FromModels(t, emptyStringModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply:             []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			// import; null_if is read back from the describe output, which does not contain the empty
			// string, so the imported value is an empty list rather than the configured [""]
			{
				Config:       config.FromModels(t, emptyStringModel),
				ResourceName: ref,
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					csvImportedResourceAssert(t, resourceId).
						HasCompression(string(sdk.CsvCompressionAuto)).
						HasRecordDelimiter(`\n`).
						HasFieldDelimiter(",").
						HasFileExtension("").
						HasDateFormat("AUTO").
						HasTimeFormat("AUTO").
						HasTimestampFormat("AUTO").
						HasBinaryFormat(string(sdk.BinaryFormatHex)).
						HasEscape("NONE").
						HasEscapeUnenclosedField(`\\`).
						HasFieldOptionallyEnclosedBy("NONE").
						HasEncoding(string(sdk.CsvEncodingUtf8)).
						HasNullIfEmpty(),
				),
			},
			// an empty string mixed with a regular value is visible in the describe output
			{
				Config: config.FromModels(t, withEmptyStringModel),
				Check: assertThat(
					t,
					resourceassert.FileFormatCsvResource(t, ref).
						HasNullIf("NULL_A", ""),
					resourceshowoutputassert.FileFormatCsvDescribeOutput(t, ref).
						HasNullIf("NULL_A", ""),
				),
			},
			// no drift
			{
				Config: config.FromModels(t, withEmptyStringModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply:             []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			// external change is detected and reverted; note that changing null_if to an empty list externally
			// would not be detected, because both values have the same describe output
			{
				PreConfig: func() {
					testClient().FileFormat.AlterCsv(t, sdk.NewAlterCsvFileFormatRequest(id).WithSet(
						*sdk.NewAlterCsvFileFormatSetRequest().WithNullIf(*sdk.NewNullIfListRequest().WithNullIf([]sdk.NullString{{S: "NULL_B"}})),
					))
				},
				Config: config.FromModels(t, emptyStringModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: assertThat(t, emptyStringAssertions...),
			},
		},
	})
}

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

// escape and escape_unenclosed_field are returned escaped by Snowflake, and skip_header uses a special
// "-1" default marker instead of Snowflake's 0, so all three are skipped during import verification.
var fileFormatCsvImportStateVerifyIgnore = []string{"escape", "escape_unenclosed_field", "skip_header"}

func TestAcc_FileFormatCsv_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
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
				Config:                  config.FromModels(t, basicModel),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: fileFormatCsvImportStateVerifyIgnore,
			},
			// set all optional fields
			{
				Config: config.FromModels(t, completeModel),
				Check:  assertThat(t, completeAssertions...),
			},
			// import with all attributes
			{
				Config:                  config.FromModels(t, completeModel),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: fileFormatCsvImportStateVerifyIgnore,
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
				Check: assertThat(
					t,
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
				),
			},
			// import
			{
				Config:                  config.FromModels(t, completeModel),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: fileFormatCsvImportStateVerifyIgnore,
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
				Config:      config.FromModels(t, conflictingHeaderFields),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`skip_header.*conflicts with\s+parse_header`),
			},
		},
	})
}

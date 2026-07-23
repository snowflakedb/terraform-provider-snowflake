//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_FileFormats(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("CreateCsv - minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateCsvFileFormatRequest(id)

		err := client.FileFormats.CreateCsv(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeCsv).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""))

		assertThatObject(t, objectassert.FileFormatCsv(t, id).
			HasId(id).
			HasType("CSV").
			HasCompression("AUTO").
			HasRecordDelimiter("\\n").
			HasFieldDelimiter(",").
			HasFileExtension("").
			HasSkipHeader(0).
			HasParseHeader(false).
			HasSkipBlankLines(false).
			HasDateFormat("AUTO").
			HasTimeFormat("AUTO").
			HasTimestampFormat("AUTO").
			HasBinaryFormat("HEX").
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
			HasEncoding("UTF8").
			HasMultiLine(true))
	})

	t.Run("CreateCsv - complete with SkipHeader", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateCsvFileFormatRequest(id).
			WithCompression(sdk.CsvCompressionGzip).
			WithRecordDelimiter(*sdk.NewStageFileFormatStringOrNoneRequest().WithValue("\\n")).
			WithFieldDelimiter(*sdk.NewStageFileFormatStringOrNoneRequest().WithValue(",")).
			WithMultiLine(false).
			WithFileExtension(".csv").
			WithSkipHeader(2).
			WithSkipBlankLines(true).
			WithDateFormat(*sdk.NewStageFileFormatStringOrAutoRequest().WithValue("YYYY-MM-DD")).
			WithTimeFormat(*sdk.NewStageFileFormatStringOrAutoRequest().WithValue("HH24:MI:SS")).
			WithTimestampFormat(*sdk.NewStageFileFormatStringOrAutoRequest().WithValue("YYYY-MM-DD HH24:MI:SS")).
			WithBinaryFormat(sdk.BinaryFormatBase64).
			WithEscape(*sdk.NewStageFileFormatStringOrNoneRequest().WithValue("!")).
			WithEscapeUnenclosedField(*sdk.NewStageFileFormatStringOrNoneRequest().WithValue("!")).
			WithTrimSpace(true).
			WithFieldOptionallyEnclosedBy(*sdk.NewStageFileFormatStringOrNoneRequest().WithValue("\"")).
			WithNullIf([]sdk.NullString{{S: "NULL"}, {S: ""}}).
			WithErrorOnColumnCountMismatch(false).
			WithReplaceInvalidCharacters(true).
			WithEmptyFieldAsNull(false).
			WithSkipByteOrderMark(false).
			WithEncoding(sdk.CsvEncodingUtf16).
			WithComment("csv complete with skip header")

		err := client.FileFormats.CreateCsv(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeCsv).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("csv complete with skip header"))

		assertThatObject(t, objectassert.FileFormatCsv(t, id).
			HasId(id).
			HasType("CSV").
			HasCompression(string(sdk.CsvCompressionGzip)).
			HasRecordDelimiter("\\n").
			HasFieldDelimiter(",").
			HasFileExtension(".csv").
			HasSkipHeader(2).
			HasParseHeader(false).
			HasSkipBlankLines(true).
			HasDateFormat("YYYY-MM-DD").
			HasTimeFormat("HH24:MI:SS").
			HasTimestampFormat("YYYY-MM-DD HH24:MI:SS").
			HasBinaryFormat(string(sdk.BinaryFormatBase64)).
			HasEscape("!").
			HasEscapeUnenclosedField("!").
			HasTrimSpace(true).
			HasFieldOptionallyEnclosedBy(`\"`).
			HasNullIf("NULL", "").
			HasErrorOnColumnCountMismatch(false).
			HasValidateUtf8(true).
			HasReplaceInvalidCharacters(true).
			HasEmptyFieldAsNull(false).
			HasSkipByteOrderMark(false).
			HasEncoding(string(sdk.CsvEncodingUtf16)).
			HasMultiLine(false))
	})

	t.Run("CreateCsv - complete with ParseHeader", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateCsvFileFormatRequest(id).
			WithParseHeader(true).
			WithCompression(sdk.CsvCompressionBz2).
			WithComment("csv complete with parse header")

		err := client.FileFormats.CreateCsv(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeCsv).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("csv complete with parse header"))

		assertThatObject(t, objectassert.FileFormatCsv(t, id).
			HasId(id).
			HasType("CSV").
			HasParseHeader(true).
			HasSkipHeader(0).
			HasCompression(string(sdk.CsvCompressionBz2)))
	})

	t.Run("AlterCsv - rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateCsv(ctx, sdk.NewCreateCsvFileFormatRequest(id))
		require.NoError(t, err)

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err = client.FileFormats.AlterCsv(ctx, sdk.NewAlterCsvFileFormatRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, newId))

		_, err = client.FileFormats.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)

		assertThatObject(t, objectassert.FileFormat(t, newId).HasName(newId.Name()))
	})

	t.Run("AlterCsv - set with SkipHeader", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateCsv(ctx, sdk.NewCreateCsvFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		err = client.FileFormats.AlterCsv(ctx, sdk.NewAlterCsvFileFormatRequest(id).
			WithSet(*sdk.NewAlterCsvFileFormatSetRequest().
				WithCompression(sdk.CsvCompressionGzip).
				WithRecordDelimiter(*sdk.NewStageFileFormatStringOrNoneRequest().WithValue("\\n")).
				WithFieldDelimiter(*sdk.NewStageFileFormatStringOrNoneRequest().WithValue(",")).
				WithMultiLine(false).
				WithFileExtension(".csv").
				WithSkipHeader(2).
				WithSkipBlankLines(true).
				WithDateFormat(*sdk.NewStageFileFormatStringOrAutoRequest().WithValue("YYYY-MM-DD")).
				WithTimeFormat(*sdk.NewStageFileFormatStringOrAutoRequest().WithValue("HH24:MI:SS")).
				WithTimestampFormat(*sdk.NewStageFileFormatStringOrAutoRequest().WithValue("YYYY-MM-DD HH24:MI:SS")).
				WithBinaryFormat(sdk.BinaryFormatBase64).
				WithEscape(*sdk.NewStageFileFormatStringOrNoneRequest().WithValue("!")).
				WithEscapeUnenclosedField(*sdk.NewStageFileFormatStringOrNoneRequest().WithValue("!")).
				WithTrimSpace(true).
				WithFieldOptionallyEnclosedBy(*sdk.NewStageFileFormatStringOrNoneRequest().WithValue("\"")).
				WithNullIf([]sdk.NullString{{S: "NULL"}}).
				WithErrorOnColumnCountMismatch(false).
				WithReplaceInvalidCharacters(true).
				WithEmptyFieldAsNull(false).
				WithSkipByteOrderMark(false).
				WithEncoding(sdk.CsvEncodingUtf16).
				WithComment("updated comment")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeCsv).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("updated comment"))

		assertThatObject(t, objectassert.FileFormatCsv(t, id).
			HasId(id).
			HasCompression(string(sdk.CsvCompressionGzip)).
			HasRecordDelimiter("\\n").
			HasFieldDelimiter(",").
			HasFileExtension(".csv").
			HasSkipHeader(2).
			HasParseHeader(false).
			HasSkipBlankLines(true).
			HasDateFormat("YYYY-MM-DD").
			HasTimeFormat("HH24:MI:SS").
			HasTimestampFormat("YYYY-MM-DD HH24:MI:SS").
			HasBinaryFormat(string(sdk.BinaryFormatBase64)).
			HasEscape("!").
			HasEscapeUnenclosedField("!").
			HasTrimSpace(true).
			HasFieldOptionallyEnclosedBy(`\"`).
			HasNullIf("NULL").
			HasErrorOnColumnCountMismatch(false).
			HasValidateUtf8(true).
			HasReplaceInvalidCharacters(true).
			HasEmptyFieldAsNull(false).
			HasSkipByteOrderMark(false).
			HasEncoding(string(sdk.CsvEncodingUtf16)).
			HasMultiLine(false))
	})

	t.Run("AlterCsv - set with ParseHeader", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateCsv(ctx, sdk.NewCreateCsvFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		err = client.FileFormats.AlterCsv(ctx, sdk.NewAlterCsvFileFormatRequest(id).
			WithSet(*sdk.NewAlterCsvFileFormatSetRequest().
				WithParseHeader(true).
				WithCompression(sdk.CsvCompressionBz2)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeCsv).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""))

		assertThatObject(t, objectassert.FileFormatCsv(t, id).
			HasId(id).
			HasParseHeader(true).
			HasSkipHeader(0).
			HasCompression(string(sdk.CsvCompressionBz2)))
	})

	t.Run("CreateJson - minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateJson(ctx, sdk.NewCreateJsonFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeJson).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""))

		assertThatObject(t, objectassert.FileFormatJson(t, id).
			HasId(id).
			HasType("JSON").
			HasCompression("AUTO").
			HasDateFormat("AUTO").
			HasTimeFormat("AUTO").
			HasTimestampFormat("AUTO").
			HasBinaryFormat("HEX").
			HasTrimSpace(false).
			HasMultiLine(true).
			HasNullIf().
			HasFileExtension("").
			HasEnableOctal(false).
			HasAllowDuplicate(false).
			HasStripOuterArray(false).
			HasStripNullValues(false).
			HasReplaceInvalidCharacters(false).
			HasIgnoreUtf8Errors(false).
			HasSkipByteOrderMark(true))
	})

	t.Run("CreateJson - complete with IgnoreUtf8Errors", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateJsonFileFormatRequest(id).
			WithCompression(sdk.JsonCompressionGzip).
			WithDateFormat(*sdk.NewStageFileFormatStringOrAutoRequest().WithValue("YYYY-MM-DD")).
			WithTimeFormat(*sdk.NewStageFileFormatStringOrAutoRequest().WithValue("HH24:MI:SS")).
			WithTimestampFormat(*sdk.NewStageFileFormatStringOrAutoRequest().WithValue("YYYY-MM-DD HH24:MI:SS")).
			WithBinaryFormat(sdk.BinaryFormatBase64).
			WithTrimSpace(true).
			WithMultiLine(false).
			WithNullIf(*sdk.NewNullIfListRequest().WithNullIf([]sdk.NullString{{S: "NULL"}})).
			WithFileExtension(".json").
			WithEnableOctal(true).
			WithAllowDuplicate(true).
			WithStripOuterArray(true).
			WithStripNullValues(true).
			WithIgnoreUtf8Errors(true).
			WithSkipByteOrderMark(false).
			WithComment("json complete with ignore utf8 errors")

		err := client.FileFormats.CreateJson(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeJson).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("json complete with ignore utf8 errors"))

		assertThatObject(t, objectassert.FileFormatJson(t, id).
			HasId(id).
			HasType("JSON").
			HasCompression(string(sdk.JsonCompressionGzip)).
			HasDateFormat("YYYY-MM-DD").
			HasTimeFormat("HH24:MI:SS").
			HasTimestampFormat("YYYY-MM-DD HH24:MI:SS").
			HasBinaryFormat(string(sdk.BinaryFormatBase64)).
			HasTrimSpace(true).
			HasMultiLine(false).
			HasNullIf("NULL").
			HasFileExtension(".json").
			HasEnableOctal(true).
			HasAllowDuplicate(true).
			HasStripOuterArray(true).
			HasStripNullValues(true).
			HasIgnoreUtf8Errors(true).
			HasSkipByteOrderMark(false))
	})

	t.Run("CreateJson - complete with ReplaceInvalidCharacters", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateJsonFileFormatRequest(id).
			WithReplaceInvalidCharacters(true).
			WithComment("json complete with replace invalid characters")

		err := client.FileFormats.CreateJson(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeJson).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("json complete with replace invalid characters"))

		assertThatObject(t, objectassert.FileFormatJson(t, id).
			HasId(id).
			HasReplaceInvalidCharacters(true).
			HasIgnoreUtf8Errors(false))
	})

	t.Run("AlterJson - rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateJson(ctx, sdk.NewCreateJsonFileFormatRequest(id))
		require.NoError(t, err)

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err = client.FileFormats.AlterJson(ctx, sdk.NewAlterJsonFileFormatRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, newId))

		assertThatObject(t, objectassert.FileFormat(t, newId).HasName(newId.Name()))
	})

	t.Run("AlterJson - set with IgnoreUtf8Errors", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateJson(ctx, sdk.NewCreateJsonFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		err = client.FileFormats.AlterJson(ctx, sdk.NewAlterJsonFileFormatRequest(id).
			WithSet(*sdk.NewAlterJsonFileFormatSetRequest().
				WithCompression(sdk.JsonCompressionGzip).
				WithDateFormat(*sdk.NewStageFileFormatStringOrAutoRequest().WithValue("YYYY-MM-DD")).
				WithTimeFormat(*sdk.NewStageFileFormatStringOrAutoRequest().WithValue("HH24:MI:SS")).
				WithTimestampFormat(*sdk.NewStageFileFormatStringOrAutoRequest().WithValue("YYYY-MM-DD HH24:MI:SS")).
				WithBinaryFormat(sdk.BinaryFormatBase64).
				WithTrimSpace(true).
				WithMultiLine(false).
				WithNullIf(*sdk.NewNullIfListRequest().WithNullIf([]sdk.NullString{{S: "NULL"}})).
				WithFileExtension(".json").
				WithEnableOctal(true).
				WithAllowDuplicate(true).
				WithStripOuterArray(true).
				WithStripNullValues(true).
				WithIgnoreUtf8Errors(true).
				WithSkipByteOrderMark(false).
				WithComment("updated comment")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeJson).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("updated comment"))

		assertThatObject(t, objectassert.FileFormatJson(t, id).
			HasId(id).
			HasCompression(string(sdk.JsonCompressionGzip)).
			HasDateFormat("YYYY-MM-DD").
			HasTimeFormat("HH24:MI:SS").
			HasTimestampFormat("YYYY-MM-DD HH24:MI:SS").
			HasBinaryFormat(string(sdk.BinaryFormatBase64)).
			HasTrimSpace(true).
			HasMultiLine(false).
			HasNullIf("NULL").
			HasFileExtension(".json").
			HasEnableOctal(true).
			HasAllowDuplicate(true).
			HasStripOuterArray(true).
			HasStripNullValues(true).
			HasIgnoreUtf8Errors(true).
			HasSkipByteOrderMark(false))
	})

	t.Run("AlterJson - set with ReplaceInvalidCharacters", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateJson(ctx, sdk.NewCreateJsonFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		err = client.FileFormats.AlterJson(ctx, sdk.NewAlterJsonFileFormatRequest(id).
			WithSet(*sdk.NewAlterJsonFileFormatSetRequest().WithReplaceInvalidCharacters(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeJson).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""))

		assertThatObject(t, objectassert.FileFormatJson(t, id).
			HasId(id).
			HasReplaceInvalidCharacters(true).
			HasIgnoreUtf8Errors(false))
	})

	t.Run("CreateAvro - minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateAvro(ctx, sdk.NewCreateAvroFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeAvro).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""))

		assertThatObject(t, objectassert.FileFormatAvro(t, id).
			HasId(id).
			HasType("AVRO").
			HasCompression("AUTO").
			HasTrimSpace(false).
			HasReplaceInvalidCharacters(false).
			HasNullIf())
	})

	t.Run("CreateAvro - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateAvroFileFormatRequest(id).
			WithCompression(sdk.AvroCompressionGzip).
			WithTrimSpace(true).
			WithReplaceInvalidCharacters(true).
			WithNullIf([]sdk.NullString{{S: "NULL"}, {S: ""}}).
			WithComment("avro complete")

		err := client.FileFormats.CreateAvro(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeAvro).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("avro complete"))

		assertThatObject(t, objectassert.FileFormatAvro(t, id).
			HasId(id).
			HasType("AVRO").
			HasCompression(string(sdk.AvroCompressionGzip)).
			HasTrimSpace(true).
			HasReplaceInvalidCharacters(true).
			HasNullIf("NULL", ""))
	})

	t.Run("AlterAvro - rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateAvro(ctx, sdk.NewCreateAvroFileFormatRequest(id))
		require.NoError(t, err)

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err = client.FileFormats.AlterAvro(ctx, sdk.NewAlterAvroFileFormatRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, newId))

		assertThatObject(t, objectassert.FileFormat(t, newId).HasName(newId.Name()))
	})

	t.Run("AlterAvro - set", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateAvro(ctx, sdk.NewCreateAvroFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		err = client.FileFormats.AlterAvro(ctx, sdk.NewAlterAvroFileFormatRequest(id).
			WithSet(*sdk.NewAlterAvroFileFormatSetRequest().
				WithCompression(sdk.AvroCompressionGzip).
				WithTrimSpace(true).
				WithReplaceInvalidCharacters(true).
				WithNullIf([]sdk.NullString{{S: "NULL"}}).
				WithComment("updated comment")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeAvro).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("updated comment"))

		assertThatObject(t, objectassert.FileFormatAvro(t, id).
			HasId(id).
			HasCompression(string(sdk.AvroCompressionGzip)).
			HasTrimSpace(true).
			HasReplaceInvalidCharacters(true).
			HasNullIf("NULL"))
	})

	t.Run("CreateOrc - minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateOrc(ctx, sdk.NewCreateOrcFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeOrc).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""))

		assertThatObject(t, objectassert.FileFormatOrc(t, id).
			HasId(id).
			HasType("ORC").
			HasTrimSpace(false).
			HasReplaceInvalidCharacters(false).
			HasNullIf())
	})

	t.Run("CreateOrc - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateOrcFileFormatRequest(id).
			WithTrimSpace(true).
			WithReplaceInvalidCharacters(true).
			WithNullIf([]sdk.NullString{{S: "NULL"}}).
			WithComment("orc complete")

		err := client.FileFormats.CreateOrc(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeOrc).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("orc complete"))

		assertThatObject(t, objectassert.FileFormatOrc(t, id).
			HasId(id).
			HasType("ORC").
			HasTrimSpace(true).
			HasReplaceInvalidCharacters(true).
			HasNullIf("NULL"))
	})

	t.Run("AlterOrc - rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateOrc(ctx, sdk.NewCreateOrcFileFormatRequest(id))
		require.NoError(t, err)

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err = client.FileFormats.AlterOrc(ctx, sdk.NewAlterOrcFileFormatRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, newId))

		assertThatObject(t, objectassert.FileFormat(t, newId).HasName(newId.Name()))
	})

	t.Run("AlterOrc - set", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateOrc(ctx, sdk.NewCreateOrcFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		err = client.FileFormats.AlterOrc(ctx, sdk.NewAlterOrcFileFormatRequest(id).
			WithSet(*sdk.NewAlterOrcFileFormatSetRequest().
				WithTrimSpace(true).
				WithReplaceInvalidCharacters(true).
				WithNullIf([]sdk.NullString{{S: "NULL"}}).
				WithComment("updated comment")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeOrc).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("updated comment"))

		assertThatObject(t, objectassert.FileFormatOrc(t, id).
			HasId(id).
			HasTrimSpace(true).
			HasReplaceInvalidCharacters(true).
			HasNullIf("NULL"))
	})

	t.Run("CreateParquet - minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateParquet(ctx, sdk.NewCreateParquetFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeParquet).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""))

		assertThatObject(t, objectassert.FileFormatParquet(t, id).
			HasId(id).
			HasType("PARQUET").
			HasCompression("AUTO").
			HasBinaryAsText(true).
			HasUseLogicalType(false).
			HasTrimSpace(false).
			HasUseVectorizedScanner(false).
			HasReplaceInvalidCharacters(false).
			HasNullIf())
	})

	t.Run("CreateParquet - complete with Compression", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateParquetFileFormatRequest(id).
			WithCompression(sdk.ParquetCompressionSnappy).
			WithBinaryAsText(false).
			WithUseLogicalType(true).
			WithTrimSpace(true).
			WithUseVectorizedScanner(true).
			WithReplaceInvalidCharacters(true).
			WithNullIf(*sdk.NewNullIfListRequest().WithNullIf([]sdk.NullString{{S: "NULL"}})).
			WithComment("parquet complete with compression")

		err := client.FileFormats.CreateParquet(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeParquet).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("parquet complete with compression"))

		assertThatObject(t, objectassert.FileFormatParquet(t, id).
			HasId(id).
			HasType("PARQUET").
			HasCompression(string(sdk.ParquetCompressionSnappy)).
			HasBinaryAsText(false).
			HasUseLogicalType(true).
			HasTrimSpace(true).
			HasUseVectorizedScanner(true).
			HasReplaceInvalidCharacters(true).
			HasNullIf("NULL"))
	})

	t.Run("CreateParquet - complete with SnappyCompression", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateParquetFileFormatRequest(id).
			WithSnappyCompression(true).
			WithBinaryAsText(false).
			WithUseLogicalType(true).
			WithTrimSpace(true).
			WithUseVectorizedScanner(true).
			WithReplaceInvalidCharacters(true).
			WithNullIf(*sdk.NewNullIfListRequest().WithNullIf([]sdk.NullString{{S: "NULL"}})).
			WithComment("parquet complete with snappy compression")

		err := client.FileFormats.CreateParquet(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeParquet).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("parquet complete with snappy compression"))

		assertThatObject(t, objectassert.FileFormatParquet(t, id).
			HasId(id).
			HasType("PARQUET").
			HasCompression("AUTO").
			HasBinaryAsText(false).
			HasUseLogicalType(true).
			HasTrimSpace(true).
			HasUseVectorizedScanner(true).
			HasReplaceInvalidCharacters(true).
			HasNullIf("NULL"))
	})

	t.Run("AlterParquet - rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateParquet(ctx, sdk.NewCreateParquetFileFormatRequest(id))
		require.NoError(t, err)

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err = client.FileFormats.AlterParquet(ctx, sdk.NewAlterParquetFileFormatRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, newId))

		assertThatObject(t, objectassert.FileFormat(t, newId).HasName(newId.Name()))
	})

	t.Run("AlterParquet - set with Compression", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateParquet(ctx, sdk.NewCreateParquetFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		err = client.FileFormats.AlterParquet(ctx, sdk.NewAlterParquetFileFormatRequest(id).
			WithSet(*sdk.NewAlterParquetFileFormatSetRequest().
				WithCompression(sdk.ParquetCompressionSnappy).
				WithBinaryAsText(false).
				WithUseLogicalType(true).
				WithTrimSpace(true).
				WithUseVectorizedScanner(true).
				WithReplaceInvalidCharacters(true).
				WithNullIf(*sdk.NewNullIfListRequest().WithNullIf([]sdk.NullString{{S: "NULL"}})).
				WithComment("updated comment")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeParquet).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("updated comment"))

		assertThatObject(t, objectassert.FileFormatParquet(t, id).
			HasId(id).
			HasCompression(string(sdk.ParquetCompressionSnappy)).
			HasBinaryAsText(false).
			HasUseLogicalType(true).
			HasTrimSpace(true).
			HasUseVectorizedScanner(true).
			HasReplaceInvalidCharacters(true).
			HasNullIf("NULL"))
	})

	t.Run("AlterParquet - set with SnappyCompression", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateParquet(ctx, sdk.NewCreateParquetFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		err = client.FileFormats.AlterParquet(ctx, sdk.NewAlterParquetFileFormatRequest(id).
			WithSet(*sdk.NewAlterParquetFileFormatSetRequest().WithSnappyCompression(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeParquet).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""))

		assertThatObject(t, objectassert.FileFormatParquet(t, id).
			HasId(id).
			HasCompression("AUTO"))
	})

	t.Run("CreateXml - minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateXml(ctx, sdk.NewCreateXmlFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeXml).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""))

		assertThatObject(t, objectassert.FileFormatXml(t, id).
			HasId(id).
			HasType("XML").
			HasCompression("AUTO").
			HasIgnoreUtf8Errors(false).
			HasPreserveSpace(false).
			HasStripOuterElement(false).
			HasDisableSnowflakeData(false).
			HasDisableAutoConvert(false).
			HasReplaceInvalidCharacters(false).
			HasSkipByteOrderMark(true))
	})

	t.Run("CreateXml - complete with IgnoreUtf8Errors", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateXmlFileFormatRequest(id).
			WithCompression(sdk.XmlCompressionGzip).
			WithIgnoreUtf8Errors(true).
			WithPreserveSpace(true).
			WithStripOuterElement(true).
			WithDisableSnowflakeData(true).
			WithDisableAutoConvert(true).
			WithSkipByteOrderMark(false).
			WithComment("xml complete with ignore utf8 errors")

		err := client.FileFormats.CreateXml(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeXml).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("xml complete with ignore utf8 errors"))

		assertThatObject(t, objectassert.FileFormatXml(t, id).
			HasId(id).
			HasType("XML").
			HasCompression(string(sdk.XmlCompressionGzip)).
			HasIgnoreUtf8Errors(true).
			HasPreserveSpace(true).
			HasStripOuterElement(true).
			HasDisableSnowflakeData(true).
			HasDisableAutoConvert(true).
			HasReplaceInvalidCharacters(false).
			HasSkipByteOrderMark(false))
	})

	t.Run("CreateXml - complete with ReplaceInvalidCharacters", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateXmlFileFormatRequest(id).
			WithCompression(sdk.XmlCompressionBz2).
			WithPreserveSpace(true).
			WithStripOuterElement(true).
			WithDisableSnowflakeData(true).
			WithDisableAutoConvert(true).
			WithReplaceInvalidCharacters(true).
			WithSkipByteOrderMark(false).
			WithComment("xml complete with replace invalid characters")

		err := client.FileFormats.CreateXml(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeXml).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("xml complete with replace invalid characters"))

		assertThatObject(t, objectassert.FileFormatXml(t, id).
			HasId(id).
			HasType("XML").
			HasCompression(string(sdk.XmlCompressionBz2)).
			HasIgnoreUtf8Errors(false).
			HasPreserveSpace(true).
			HasStripOuterElement(true).
			HasDisableSnowflakeData(true).
			HasDisableAutoConvert(true).
			HasReplaceInvalidCharacters(true).
			HasSkipByteOrderMark(false))
	})

	t.Run("AlterXml - rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateXml(ctx, sdk.NewCreateXmlFileFormatRequest(id))
		require.NoError(t, err)

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err = client.FileFormats.AlterXml(ctx, sdk.NewAlterXmlFileFormatRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, newId))

		assertThatObject(t, objectassert.FileFormat(t, newId).HasName(newId.Name()))
	})

	t.Run("AlterXml - set with IgnoreUtf8Errors", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateXml(ctx, sdk.NewCreateXmlFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		err = client.FileFormats.AlterXml(ctx, sdk.NewAlterXmlFileFormatRequest(id).
			WithSet(*sdk.NewAlterXmlFileFormatSetRequest().
				WithCompression(sdk.XmlCompressionGzip).
				WithIgnoreUtf8Errors(true).
				WithPreserveSpace(true).
				WithStripOuterElement(true).
				WithDisableSnowflakeData(true).
				WithDisableAutoConvert(true).
				WithSkipByteOrderMark(false).
				WithComment("updated comment")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeXml).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("updated comment"))

		assertThatObject(t, objectassert.FileFormatXml(t, id).
			HasId(id).
			HasCompression(string(sdk.XmlCompressionGzip)).
			HasIgnoreUtf8Errors(true).
			HasPreserveSpace(true).
			HasStripOuterElement(true).
			HasDisableSnowflakeData(true).
			HasDisableAutoConvert(true).
			HasSkipByteOrderMark(false))
	})

	t.Run("AlterXml - set with ReplaceInvalidCharacters", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateXml(ctx, sdk.NewCreateXmlFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		err = client.FileFormats.AlterXml(ctx, sdk.NewAlterXmlFileFormatRequest(id).
			WithSet(*sdk.NewAlterXmlFileFormatSetRequest().WithReplaceInvalidCharacters(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.FileFormat(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeXml).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""))

		assertThatObject(t, objectassert.FileFormatXml(t, id).
			HasId(id).
			HasReplaceInvalidCharacters(true).
			HasIgnoreUtf8Errors(false))
	})

	t.Run("Drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateCsv(ctx, sdk.NewCreateCsvFileFormatRequest(id))
		require.NoError(t, err)

		err = client.FileFormats.Drop(ctx, sdk.NewDropFileFormatRequest(id))
		require.NoError(t, err)

		_, err = client.FileFormats.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("Drop with IfExists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.Drop(ctx, sdk.NewDropFileFormatRequest(id).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("Show - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateCsv(ctx, sdk.NewCreateCsvFileFormatRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		fileFormats, err := client.FileFormats.Show(ctx, sdk.NewShowFileFormatRequest())
		require.NoError(t, err)
		require.NotEmpty(t, fileFormats)
	})

	t.Run("Show - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.CreateCsv(ctx, sdk.NewCreateCsvFileFormatRequest(id).WithComment("show test"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().FileFormat.DropFileFormatFunc(t, id))

		fileFormats, err := client.FileFormats.Show(ctx, sdk.NewShowFileFormatRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}).
			WithIn(sdk.In{Schema: id.SchemaId()}))
		require.NoError(t, err)
		require.Len(t, fileFormats, 1)
		assert.Equal(t, id.Name(), fileFormats[0].Name)
		assert.Equal(t, id.DatabaseName(), fileFormats[0].DatabaseName)
		assert.Equal(t, id.SchemaName(), fileFormats[0].SchemaName)
		assert.Equal(t, sdk.FileFormatTypeCsv, fileFormats[0].Type)
		assert.Equal(t, "show test", fileFormats[0].Comment)
		assert.NotEmpty(t, fileFormats[0].Owner)
		assert.NotEmpty(t, fileFormats[0].OwnerRoleType)
		assert.NotEmpty(t, fileFormats[0].FormatOptions)
	})
}

//go:build non_account_level_tests

package testint

import (
	"errors"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_FileFormatsCreateAndRead(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("CSV", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormatsLegacy.Create(ctx, id, &sdk.CreateFileFormatOptionsLegacy{
			Type: sdk.FileFormatTypeCsv,
			FileFormatTypeOptionsLegacy: sdk.FileFormatTypeOptionsLegacy{
				CSVCompression:                sdk.Pointer(sdk.CsvCompressionBz2),
				CSVRecordDelimiter:            sdk.String("\\123"),
				CSVFieldDelimiter:             sdk.String("0x42"),
				CSVFileExtension:              sdk.String("c"),
				CSVParseHeader:                sdk.Bool(true),
				CSVSkipBlankLines:             sdk.Bool(true),
				CSVDateFormat:                 sdk.String("d"),
				CSVTimeFormat:                 sdk.String("e"),
				CSVTimestampFormat:            sdk.String("f"),
				CSVBinaryFormat:               sdk.Pointer(sdk.BinaryFormatBase64),
				CSVEscape:                     sdk.String(`\`),
				CSVEscapeUnenclosedField:      sdk.String("h"),
				CSVTrimSpace:                  sdk.Bool(true),
				CSVFieldOptionallyEnclosedBy:  sdk.String("'"),
				CSVNullIf:                     &[]sdk.NullString{{S: "j"}, {S: "k"}},
				CSVErrorOnColumnCountMismatch: sdk.Bool(true),
				CSVReplaceInvalidCharacters:   sdk.Bool(true),
				CSVEmptyFieldAsNull:           sdk.Bool(true),
				CSVSkipByteOrderMark:          sdk.Bool(true),
				CSVEncoding:                   sdk.Pointer(sdk.CsvEncodingGb18030),
			},
			Comment: sdk.String("test comment"),
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormatsLegacy.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormatsLegacy.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, sdk.FileFormatTypeCsv, result.Type)
		assert.Equal(t, client.GetConfig().Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "ROLE", result.OwnerRoleType)
		assert.Equal(t, sdk.Pointer(sdk.CsvCompressionBz2), result.Options.CSVCompression)
		assert.Equal(t, "S", *result.Options.CSVRecordDelimiter) // o123 == 83 == 'S' (ASCII)
		assert.Equal(t, "B", *result.Options.CSVFieldDelimiter)  // 0x42 == 66 == 'B' (ASCII)
		assert.Equal(t, "c", *result.Options.CSVFileExtension)
		assert.True(t, *result.Options.CSVParseHeader)
		assert.True(t, *result.Options.CSVSkipBlankLines)
		assert.Equal(t, "d", *result.Options.CSVDateFormat)
		assert.Equal(t, "e", *result.Options.CSVTimeFormat)
		assert.Equal(t, "f", *result.Options.CSVTimestampFormat)
		assert.Equal(t, sdk.Pointer(sdk.BinaryFormatBase64), result.Options.CSVBinaryFormat)
		assert.Equal(t, `\`, *result.Options.CSVEscape)
		assert.Equal(t, "h", *result.Options.CSVEscapeUnenclosedField)
		assert.True(t, *result.Options.CSVTrimSpace)
		assert.Equal(t, "'", *result.Options.CSVFieldOptionallyEnclosedBy)
		assert.Equal(t, &[]sdk.NullString{{S: "j"}, {S: "k"}}, result.Options.CSVNullIf)
		assert.True(t, *result.Options.CSVErrorOnColumnCountMismatch)
		assert.True(t, *result.Options.CSVReplaceInvalidCharacters)
		assert.True(t, *result.Options.CSVEmptyFieldAsNull)
		assert.True(t, *result.Options.CSVSkipByteOrderMark)
		assert.Equal(t, sdk.Pointer(sdk.CsvEncodingGb18030), result.Options.CSVEncoding)

		describeResult, err := client.FileFormatsLegacy.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeCsv, describeResult.Type)
		assert.Equal(t, sdk.Pointer(sdk.CsvCompressionBz2), describeResult.Options.CSVCompression)
		assert.Equal(t, "S", *describeResult.Options.CSVRecordDelimiter) // o123 == 83 == 'S' (ASCII)
		assert.Equal(t, "B", *describeResult.Options.CSVFieldDelimiter)  // 0x42 == 66 == 'B' (ASCII)
		assert.Equal(t, "c", *describeResult.Options.CSVFileExtension)
		assert.True(t, *describeResult.Options.CSVParseHeader)
		assert.True(t, *describeResult.Options.CSVSkipBlankLines)
		assert.Equal(t, "d", *describeResult.Options.CSVDateFormat)
		assert.Equal(t, "e", *describeResult.Options.CSVTimeFormat)
		assert.Equal(t, "f", *describeResult.Options.CSVTimestampFormat)
		assert.Equal(t, sdk.Pointer(sdk.BinaryFormatBase64), describeResult.Options.CSVBinaryFormat)
		assert.Equal(t, `\\`, *describeResult.Options.CSVEscape) // Describe does not un-escape backslashes, but show does ....
		assert.Equal(t, "h", *describeResult.Options.CSVEscapeUnenclosedField)
		assert.True(t, *describeResult.Options.CSVTrimSpace)
		assert.Equal(t, "'", *describeResult.Options.CSVFieldOptionallyEnclosedBy)
		assert.Equal(t, &[]sdk.NullString{{S: "j"}, {S: "k"}}, describeResult.Options.CSVNullIf)
		assert.True(t, *describeResult.Options.CSVErrorOnColumnCountMismatch)
		assert.True(t, *describeResult.Options.CSVReplaceInvalidCharacters)
		assert.True(t, *describeResult.Options.CSVEmptyFieldAsNull)
		assert.True(t, *describeResult.Options.CSVSkipByteOrderMark)
		assert.Equal(t, sdk.Pointer(sdk.CsvEncodingGb18030), describeResult.Options.CSVEncoding)
	})

	// Check that field_optionally_enclosed_by can take the value NONE
	t.Run("CSV", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormatsLegacy.Create(ctx, id, &sdk.CreateFileFormatOptionsLegacy{
			Type: sdk.FileFormatTypeCsv,
			FileFormatTypeOptionsLegacy: sdk.FileFormatTypeOptionsLegacy{
				CSVFieldOptionallyEnclosedBy: sdk.String("NONE"),
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormatsLegacy.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormatsLegacy.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "NONE", *result.Options.CSVFieldOptionallyEnclosedBy)
	})
	t.Run("JSON", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormatsLegacy.Create(ctx, id, &sdk.CreateFileFormatOptionsLegacy{
			Type: sdk.FileFormatTypeJson,
			FileFormatTypeOptionsLegacy: sdk.FileFormatTypeOptionsLegacy{
				JSONCompression:       sdk.Pointer(sdk.JsonCompressionBrotli),
				JSONDateFormat:        sdk.String("a"),
				JSONTimeFormat:        sdk.String("b"),
				JSONTimestampFormat:   sdk.String("c"),
				JSONBinaryFormat:      sdk.Pointer(sdk.BinaryFormatHex),
				JSONTrimSpace:         sdk.Bool(true),
				JSONNullIf:            []sdk.NullString{{S: "d"}, {S: "e"}},
				JSONFileExtension:     sdk.String("f"),
				JSONEnableOctal:       sdk.Bool(true),
				JSONAllowDuplicate:    sdk.Bool(true),
				JSONStripOuterArray:   sdk.Bool(true),
				JSONStripNullValues:   sdk.Bool(true),
				JSONIgnoreUTF8Errors:  sdk.Bool(true),
				JSONSkipByteOrderMark: sdk.Bool(true),
			},
			Comment: sdk.String("test comment"),
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormatsLegacy.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormatsLegacy.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, sdk.FileFormatTypeJson, result.Type)
		assert.Equal(t, client.GetConfig().Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "ROLE", result.OwnerRoleType)

		assert.Equal(t, sdk.JsonCompressionBrotli, *result.Options.JSONCompression)
		assert.Equal(t, "a", *result.Options.JSONDateFormat)
		assert.Equal(t, "b", *result.Options.JSONTimeFormat)
		assert.Equal(t, "c", *result.Options.JSONTimestampFormat)
		assert.Equal(t, sdk.BinaryFormatHex, *result.Options.JSONBinaryFormat)
		assert.True(t, *result.Options.JSONTrimSpace)
		assert.Equal(t, []sdk.NullString{{S: "d"}, {S: "e"}}, result.Options.JSONNullIf)
		assert.Equal(t, "f", *result.Options.JSONFileExtension)
		assert.True(t, *result.Options.JSONEnableOctal)
		assert.True(t, *result.Options.JSONAllowDuplicate)
		assert.True(t, *result.Options.JSONStripOuterArray)
		assert.True(t, *result.Options.JSONStripNullValues)
		assert.True(t, *result.Options.JSONIgnoreUTF8Errors)
		assert.True(t, *result.Options.JSONSkipByteOrderMark)

		describeResult, err := client.FileFormatsLegacy.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeJson, describeResult.Type)
		assert.Equal(t, sdk.JsonCompressionBrotli, *describeResult.Options.JSONCompression)
		assert.Equal(t, "a", *describeResult.Options.JSONDateFormat)
		assert.Equal(t, "b", *describeResult.Options.JSONTimeFormat)
		assert.Equal(t, "c", *describeResult.Options.JSONTimestampFormat)
		assert.Equal(t, sdk.BinaryFormatHex, *describeResult.Options.JSONBinaryFormat)
		assert.True(t, *describeResult.Options.JSONTrimSpace)
		assert.Equal(t, []sdk.NullString{{S: "d"}, {S: "e"}}, describeResult.Options.JSONNullIf)
		assert.Equal(t, "f", *describeResult.Options.JSONFileExtension)
		assert.True(t, *describeResult.Options.JSONEnableOctal)
		assert.True(t, *describeResult.Options.JSONAllowDuplicate)
		assert.True(t, *describeResult.Options.JSONStripOuterArray)
		assert.True(t, *describeResult.Options.JSONStripNullValues)
		assert.True(t, *describeResult.Options.JSONIgnoreUTF8Errors)
		assert.True(t, *describeResult.Options.JSONSkipByteOrderMark)
	})
	t.Run("AVRO", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormatsLegacy.Create(ctx, id, &sdk.CreateFileFormatOptionsLegacy{
			Type: sdk.FileFormatTypeAvro,
			FileFormatTypeOptionsLegacy: sdk.FileFormatTypeOptionsLegacy{
				AvroCompression:              sdk.Pointer(sdk.AvroCompressionGzip),
				AvroTrimSpace:                sdk.Bool(true),
				AvroReplaceInvalidCharacters: sdk.Bool(true),
				AvroNullIf:                   &[]sdk.NullString{{S: "a"}, {S: "b"}},
			},
			Comment: sdk.String("test comment"),
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormatsLegacy.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormatsLegacy.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, sdk.FileFormatTypeAvro, result.Type)
		assert.Equal(t, client.GetConfig().Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "ROLE", result.OwnerRoleType)

		assert.Equal(t, sdk.AvroCompressionGzip, *result.Options.AvroCompression)
		assert.True(t, *result.Options.AvroTrimSpace)
		assert.True(t, *result.Options.AvroReplaceInvalidCharacters)
		assert.Equal(t, []sdk.NullString{{S: "a"}, {S: "b"}}, *result.Options.AvroNullIf)

		describeResult, err := client.FileFormatsLegacy.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeAvro, describeResult.Type)
		assert.Equal(t, sdk.AvroCompressionGzip, *describeResult.Options.AvroCompression)
		assert.True(t, *describeResult.Options.AvroTrimSpace)
		assert.True(t, *describeResult.Options.AvroReplaceInvalidCharacters)
		assert.Equal(t, []sdk.NullString{{S: "a"}, {S: "b"}}, *describeResult.Options.AvroNullIf)
	})
	t.Run("ORC", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormatsLegacy.Create(ctx, id, &sdk.CreateFileFormatOptionsLegacy{
			Type: sdk.FileFormatTypeOrc,
			FileFormatTypeOptionsLegacy: sdk.FileFormatTypeOptionsLegacy{
				ORCTrimSpace:                sdk.Bool(true),
				ORCReplaceInvalidCharacters: sdk.Bool(true),
				ORCNullIf:                   &[]sdk.NullString{{S: "a"}, {S: "b"}},
			},
			Comment: sdk.String("test comment"),
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormatsLegacy.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormatsLegacy.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, sdk.FileFormatTypeOrc, result.Type)
		assert.Equal(t, client.GetConfig().Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "ROLE", result.OwnerRoleType)

		assert.True(t, *result.Options.ORCTrimSpace)
		assert.True(t, *result.Options.ORCReplaceInvalidCharacters)
		assert.Equal(t, []sdk.NullString{{S: "a"}, {S: "b"}}, *result.Options.ORCNullIf)

		describeResult, err := client.FileFormatsLegacy.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeOrc, describeResult.Type)
		assert.True(t, *describeResult.Options.ORCTrimSpace)
		assert.True(t, *describeResult.Options.ORCReplaceInvalidCharacters)
		assert.Equal(t, []sdk.NullString{{S: "a"}, {S: "b"}}, *describeResult.Options.ORCNullIf)
	})
	t.Run("PARQUET", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormatsLegacy.Create(ctx, id, &sdk.CreateFileFormatOptionsLegacy{
			Type: sdk.FileFormatTypeParquet,
			FileFormatTypeOptionsLegacy: sdk.FileFormatTypeOptionsLegacy{
				ParquetCompression:              sdk.Pointer(sdk.ParquetCompressionLzo),
				ParquetBinaryAsText:             sdk.Bool(true),
				ParquetTrimSpace:                sdk.Bool(true),
				ParquetReplaceInvalidCharacters: sdk.Bool(true),
				ParquetNullIf:                   &[]sdk.NullString{{S: "a"}, {S: "b"}},
			},
			Comment: sdk.String("test comment"),
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormatsLegacy.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormatsLegacy.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, sdk.FileFormatTypeParquet, result.Type)
		assert.Equal(t, client.GetConfig().Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "ROLE", result.OwnerRoleType)

		assert.Equal(t, sdk.ParquetCompressionLzo, *result.Options.ParquetCompression)
		assert.True(t, *result.Options.ParquetBinaryAsText)
		assert.True(t, *result.Options.ParquetTrimSpace)
		assert.True(t, *result.Options.ParquetReplaceInvalidCharacters)
		assert.Equal(t, []sdk.NullString{{S: "a"}, {S: "b"}}, *result.Options.ParquetNullIf)

		describeResult, err := client.FileFormatsLegacy.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeParquet, describeResult.Type)
		assert.Equal(t, sdk.ParquetCompressionLzo, *describeResult.Options.ParquetCompression)
		assert.True(t, *describeResult.Options.ParquetBinaryAsText)
		assert.True(t, *describeResult.Options.ParquetTrimSpace)
		assert.True(t, *describeResult.Options.ParquetReplaceInvalidCharacters)
		assert.Equal(t, []sdk.NullString{{S: "a"}, {S: "b"}}, *describeResult.Options.ParquetNullIf)
	})
	t.Run("XML", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormatsLegacy.Create(ctx, id, &sdk.CreateFileFormatOptionsLegacy{
			Type: sdk.FileFormatTypeXml,
			FileFormatTypeOptionsLegacy: sdk.FileFormatTypeOptionsLegacy{
				XMLCompression:          sdk.Pointer(sdk.XmlCompressionDeflate),
				XMLIgnoreUTF8Errors:     sdk.Bool(true),
				XMLPreserveSpace:        sdk.Bool(true),
				XMLStripOuterElement:    sdk.Bool(true),
				XMLDisableSnowflakeData: sdk.Bool(true),
				XMLDisableAutoConvert:   sdk.Bool(true),
				XMLSkipByteOrderMark:    sdk.Bool(true),
			},
			Comment: sdk.String("test comment"),
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormatsLegacy.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormatsLegacy.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, sdk.FileFormatTypeXml, result.Type)
		assert.Equal(t, client.GetConfig().Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "ROLE", result.OwnerRoleType)

		assert.Equal(t, sdk.XmlCompressionDeflate, *result.Options.XMLCompression)
		assert.True(t, *result.Options.XMLIgnoreUTF8Errors)
		assert.True(t, *result.Options.XMLPreserveSpace)
		assert.True(t, *result.Options.XMLStripOuterElement)
		assert.True(t, *result.Options.XMLDisableSnowflakeData)
		assert.True(t, *result.Options.XMLDisableAutoConvert)
		assert.True(t, *result.Options.XMLSkipByteOrderMark)

		describeResult, err := client.FileFormatsLegacy.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeXml, describeResult.Type)
		assert.Equal(t, sdk.XmlCompressionDeflate, *describeResult.Options.XMLCompression)
		assert.True(t, *describeResult.Options.XMLIgnoreUTF8Errors)
		assert.True(t, *describeResult.Options.XMLPreserveSpace)
		assert.True(t, *describeResult.Options.XMLStripOuterElement)
		assert.True(t, *describeResult.Options.XMLDisableSnowflakeData)
		assert.True(t, *describeResult.Options.XMLDisableAutoConvert)
		assert.True(t, *describeResult.Options.XMLSkipByteOrderMark)
	})
}

func TestInt_FileFormatsAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("rename", func(t *testing.T) {
		fileFormat, fileFormatCleanup := testClientHelper().FileFormat.CreateFileFormat(t)
		t.Cleanup(fileFormatCleanup)
		oldId := fileFormat.ID()
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.FileFormatsLegacy.Alter(ctx, oldId, &sdk.AlterFileFormatOptionsLegacy{
			Rename: &sdk.AlterFileFormatRenameOptions{
				NewName: newId,
			},
		})
		require.NoError(t, err)

		_, err = client.FileFormatsLegacy.ShowByID(ctx, oldId)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)

		result, err := client.FileFormatsLegacy.ShowByID(ctx, newId)
		require.NoError(t, err)
		assert.Equal(t, newId, result.Name)

		// Undo rename so we can clean up
		err = client.FileFormatsLegacy.Alter(ctx, newId, &sdk.AlterFileFormatOptionsLegacy{
			Rename: &sdk.AlterFileFormatRenameOptions{
				NewName: oldId,
			},
		})
		require.NoError(t, err)
	})

	t.Run("set + set comment", func(t *testing.T) {
		fileFormat, fileFormatCleanup := testClientHelper().FileFormat.CreateFileFormatWithOptions(t, &sdk.CreateFileFormatOptionsLegacy{
			Type: sdk.FileFormatTypeCsv,
			FileFormatTypeOptionsLegacy: sdk.FileFormatTypeOptionsLegacy{
				CSVCompression: sdk.Pointer(sdk.CsvCompressionAuto),
				CSVParseHeader: sdk.Bool(false),
			},
		})
		t.Cleanup(fileFormatCleanup)

		err := client.FileFormatsLegacy.Alter(ctx, fileFormat.ID(), &sdk.AlterFileFormatOptionsLegacy{
			Set: &sdk.FileFormatTypeOptionsLegacy{
				Comment:        sdk.String("some comment"),
				CSVCompression: sdk.Pointer(sdk.CsvCompressionBz2),
				CSVParseHeader: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		result, err := client.FileFormatsLegacy.ShowByID(ctx, fileFormat.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.CsvCompressionBz2, *result.Options.CSVCompression)
		assert.True(t, *result.Options.CSVParseHeader)
		assert.Equal(t, "some comment", result.Comment)
	})
}

func TestInt_FileFormatsDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("no options", func(t *testing.T) {
		fileFormat, fileFormatCleanup := testClientHelper().FileFormat.CreateFileFormat(t)
		t.Cleanup(fileFormatCleanup)

		err := client.FileFormatsLegacy.Drop(ctx, fileFormat.ID(), nil)
		require.NoError(t, err)

		_, err = client.FileFormatsLegacy.ShowByID(ctx, fileFormat.ID())
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("with IfExists", func(t *testing.T) {
		fileFormat, fileFormatCleanup := testClientHelper().FileFormat.CreateFileFormat(t)
		t.Cleanup(fileFormatCleanup)

		err := client.FileFormatsLegacy.Drop(ctx, fileFormat.ID(), &sdk.DropFileFormatOptionsLegacy{
			IfExists: sdk.Bool(true),
		})
		require.NoError(t, err)

		_, err = client.FileFormatsLegacy.ShowByID(ctx, fileFormat.ID())
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})
}

func TestInt_FileFormatsShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	fileFormatTest, cleanupFileFormat := testClientHelper().FileFormat.CreateFileFormat(t)
	t.Cleanup(cleanupFileFormat)
	fileFormatTest2, cleanupFileFormat2 := testClientHelper().FileFormat.CreateFileFormat(t)
	t.Cleanup(cleanupFileFormat2)

	t.Run("without options", func(t *testing.T) {
		fileFormats, err := client.FileFormatsLegacy.Show(ctx, nil)
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(fileFormats))
		assert.Contains(t, fileFormats, *fileFormatTest)
		assert.Contains(t, fileFormats, *fileFormatTest2)
	})

	t.Run("LIKE", func(t *testing.T) {
		fileFormats, err := client.FileFormatsLegacy.Show(ctx, &sdk.ShowFileFormatsOptionsLegacy{
			Like: &sdk.Like{
				Pattern: sdk.String(fileFormatTest.Name.Name()),
			},
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, 1, len(fileFormats))
		assert.Contains(t, fileFormats, *fileFormatTest)
	})

	t.Run("IN", func(t *testing.T) {
		fileFormats, err := client.FileFormatsLegacy.Show(ctx, &sdk.ShowFileFormatsOptionsLegacy{
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(fileFormats))
		assert.Contains(t, fileFormats, *fileFormatTest)
		assert.Contains(t, fileFormats, *fileFormatTest2)
	})
}

func TestInt_FileFormatsShowById(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	fileFormatTest, cleanupFileFormat := testClientHelper().FileFormat.CreateFileFormat(t)
	t.Cleanup(cleanupFileFormat)

	// new database and schema created on purpose
	databaseTest2, cleanupDatabase2 := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase2)
	schemaTest2, cleanupSchema2 := testClientHelper().Schema.CreateSchemaInDatabase(t, databaseTest2.ID())
	t.Cleanup(cleanupSchema2)

	t.Run("show format in different schema", func(t *testing.T) {
		err := client.Sessions.UseDatabase(ctx, sdk.NewUseDatabaseSessionRequest(databaseTest2.ID()))
		require.NoError(t, err)
		err = client.Sessions.UseSchema(ctx, sdk.NewUseSchemaSessionRequest(schemaTest2.ID()))
		require.NoError(t, err)

		fileFormat, err := client.FileFormatsLegacy.ShowByID(ctx, fileFormatTest.ID())
		require.NoError(t, err)
		assert.Equal(t, testClientHelper().Ids.DatabaseId().Name(), fileFormat.Name.DatabaseName())
		assert.Equal(t, testClientHelper().Ids.SchemaId().Name(), fileFormat.Name.SchemaName())
		assert.Equal(t, fileFormatTest.Name.Name(), fileFormat.Name.Name())
	})
}

func TestInt_FileFormatsShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupFileFormatHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.FileFormatsLegacy.Drop(ctx, id, nil)
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createFileFormatHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		err := client.FileFormatsLegacy.Create(ctx, id, &sdk.CreateFileFormatOptionsLegacy{Type: sdk.FileFormatTypeCsv})
		require.NoError(t, err)
		t.Cleanup(cleanupFileFormatHandle(t, id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createFileFormatHandle(t, id1)
		createFileFormatHandle(t, id2)

		e1, err := client.FileFormatsLegacy.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.FileFormatsLegacy.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}

//go:build !account_level_tests

package testint

import (
	"errors"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_FileFormatsCreateAndRead(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("CSV", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.Create(ctx, id, &sdk.CreateFileFormatOptions{
			Type: sdk.FileFormatTypeCSV,
			FileFormatTypeOptions: sdk.FileFormatTypeOptions{
				CSVCompression:                &sdk.CSVCompressionBz2,
				CSVRecordDelimiter:            sdk.String("\\123"),
				CSVFieldDelimiter:             sdk.String("0x42"),
				CSVFileExtension:              sdk.String("c"),
				CSVParseHeader:                sdk.Bool(true),
				CSVSkipBlankLines:             sdk.Bool(true),
				CSVDateFormat:                 sdk.String("d"),
				CSVTimeFormat:                 sdk.String("e"),
				CSVTimestampFormat:            sdk.String("f"),
				CSVBinaryFormat:               &sdk.BinaryFormatBase64,
				CSVEscape:                     sdk.String(`\`),
				CSVEscapeUnenclosedField:      sdk.String("h"),
				CSVTrimSpace:                  sdk.Bool(true),
				CSVFieldOptionallyEnclosedBy:  sdk.String("'"),
				CSVNullIf:                     &[]sdk.NullString{{S: "j"}, {S: "k"}},
				CSVErrorOnColumnCountMismatch: sdk.Bool(true),
				CSVReplaceInvalidCharacters:   sdk.Bool(true),
				CSVEmptyFieldAsNull:           sdk.Bool(true),
				CSVSkipByteOrderMark:          sdk.Bool(true),
				CSVEncoding:                   &sdk.CSVEncodingGB18030,
			},
			Comment: sdk.String("test comment"),
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, sdk.FileFormatTypeCSV, result.Type)
		assert.Equal(t, client.GetConfig().Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "ROLE", result.OwnerRoleType)
		assert.Equal(t, &sdk.CSVCompressionBz2, result.Options.CSVCompression)
		assert.Equal(t, "S", *result.Options.CSVRecordDelimiter) // o123 == 83 == 'S' (ASCII)
		assert.Equal(t, "B", *result.Options.CSVFieldDelimiter)  // 0x42 == 66 == 'B' (ASCII)
		assert.Equal(t, "c", *result.Options.CSVFileExtension)
		assert.Equal(t, true, *result.Options.CSVParseHeader)
		assert.Equal(t, true, *result.Options.CSVSkipBlankLines)
		assert.Equal(t, "d", *result.Options.CSVDateFormat)
		assert.Equal(t, "e", *result.Options.CSVTimeFormat)
		assert.Equal(t, "f", *result.Options.CSVTimestampFormat)
		assert.Equal(t, &sdk.BinaryFormatBase64, result.Options.CSVBinaryFormat)
		assert.Equal(t, `\`, *result.Options.CSVEscape)
		assert.Equal(t, "h", *result.Options.CSVEscapeUnenclosedField)
		assert.Equal(t, true, *result.Options.CSVTrimSpace)
		assert.Equal(t, "'", *result.Options.CSVFieldOptionallyEnclosedBy)
		assert.Equal(t, &[]sdk.NullString{{S: "j"}, {S: "k"}}, result.Options.CSVNullIf)
		assert.Equal(t, true, *result.Options.CSVErrorOnColumnCountMismatch)
		assert.Equal(t, true, *result.Options.CSVReplaceInvalidCharacters)
		assert.Equal(t, true, *result.Options.CSVEmptyFieldAsNull)
		assert.Equal(t, true, *result.Options.CSVSkipByteOrderMark)
		assert.Equal(t, &sdk.CSVEncodingGB18030, result.Options.CSVEncoding)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeCSV, describeResult.Type)
		assert.Equal(t, &sdk.CSVCompressionBz2, describeResult.Options.CSVCompression)
		assert.Equal(t, "S", *describeResult.Options.CSVRecordDelimiter) // o123 == 83 == 'S' (ASCII)
		assert.Equal(t, "B", *describeResult.Options.CSVFieldDelimiter)  // 0x42 == 66 == 'B' (ASCII)
		assert.Equal(t, "c", *describeResult.Options.CSVFileExtension)
		assert.Equal(t, true, *describeResult.Options.CSVParseHeader)
		assert.Equal(t, true, *describeResult.Options.CSVSkipBlankLines)
		assert.Equal(t, "d", *describeResult.Options.CSVDateFormat)
		assert.Equal(t, "e", *describeResult.Options.CSVTimeFormat)
		assert.Equal(t, "f", *describeResult.Options.CSVTimestampFormat)
		assert.Equal(t, &sdk.BinaryFormatBase64, describeResult.Options.CSVBinaryFormat)
		assert.Equal(t, `\\`, *describeResult.Options.CSVEscape) // Describe does not un-escape backslashes, but show does ....
		assert.Equal(t, "h", *describeResult.Options.CSVEscapeUnenclosedField)
		assert.Equal(t, true, *describeResult.Options.CSVTrimSpace)
		assert.Equal(t, "'", *describeResult.Options.CSVFieldOptionallyEnclosedBy)
		assert.Equal(t, &[]sdk.NullString{{S: "j"}, {S: "k"}}, describeResult.Options.CSVNullIf)
		assert.Equal(t, true, *describeResult.Options.CSVErrorOnColumnCountMismatch)
		assert.Equal(t, true, *describeResult.Options.CSVReplaceInvalidCharacters)
		assert.Equal(t, true, *describeResult.Options.CSVEmptyFieldAsNull)
		assert.Equal(t, true, *describeResult.Options.CSVSkipByteOrderMark)
		assert.Equal(t, &sdk.CSVEncodingGB18030, describeResult.Options.CSVEncoding)
	})

	// Check that field_optionally_enclosed_by can take the value NONE
	t.Run("CSV", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.Create(ctx, id, &sdk.CreateFileFormatOptions{
			Type: sdk.FileFormatTypeCSV,
			FileFormatTypeOptions: sdk.FileFormatTypeOptions{
				CSVFieldOptionallyEnclosedBy: sdk.String("NONE"),
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "NONE", *result.Options.CSVFieldOptionallyEnclosedBy)
	})
	t.Run("JSON", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.Create(ctx, id, &sdk.CreateFileFormatOptions{
			Type: sdk.FileFormatTypeJSON,
			FileFormatTypeOptions: sdk.FileFormatTypeOptions{
				JSONCompression:       &sdk.JSONCompressionBrotli,
				JSONDateFormat:        sdk.String("a"),
				JSONTimeFormat:        sdk.String("b"),
				JSONTimestampFormat:   sdk.String("c"),
				JSONBinaryFormat:      &sdk.BinaryFormatHex,
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
			err := client.FileFormats.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, sdk.FileFormatTypeJSON, result.Type)
		assert.Equal(t, client.GetConfig().Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "ROLE", result.OwnerRoleType)

		assert.Equal(t, sdk.JSONCompressionBrotli, *result.Options.JSONCompression)
		assert.Equal(t, "a", *result.Options.JSONDateFormat)
		assert.Equal(t, "b", *result.Options.JSONTimeFormat)
		assert.Equal(t, "c", *result.Options.JSONTimestampFormat)
		assert.Equal(t, sdk.BinaryFormatHex, *result.Options.JSONBinaryFormat)
		assert.Equal(t, true, *result.Options.JSONTrimSpace)
		assert.Equal(t, []sdk.NullString{{S: "d"}, {S: "e"}}, result.Options.JSONNullIf)
		assert.Equal(t, "f", *result.Options.JSONFileExtension)
		assert.Equal(t, true, *result.Options.JSONEnableOctal)
		assert.Equal(t, true, *result.Options.JSONAllowDuplicate)
		assert.Equal(t, true, *result.Options.JSONStripOuterArray)
		assert.Equal(t, true, *result.Options.JSONStripNullValues)
		assert.Equal(t, true, *result.Options.JSONIgnoreUTF8Errors)
		assert.Equal(t, true, *result.Options.JSONSkipByteOrderMark)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeJSON, describeResult.Type)
		assert.Equal(t, sdk.JSONCompressionBrotli, *describeResult.Options.JSONCompression)
		assert.Equal(t, "a", *describeResult.Options.JSONDateFormat)
		assert.Equal(t, "b", *describeResult.Options.JSONTimeFormat)
		assert.Equal(t, "c", *describeResult.Options.JSONTimestampFormat)
		assert.Equal(t, sdk.BinaryFormatHex, *describeResult.Options.JSONBinaryFormat)
		assert.Equal(t, true, *describeResult.Options.JSONTrimSpace)
		assert.Equal(t, []sdk.NullString{{S: "d"}, {S: "e"}}, describeResult.Options.JSONNullIf)
		assert.Equal(t, "f", *describeResult.Options.JSONFileExtension)
		assert.Equal(t, true, *describeResult.Options.JSONEnableOctal)
		assert.Equal(t, true, *describeResult.Options.JSONAllowDuplicate)
		assert.Equal(t, true, *describeResult.Options.JSONStripOuterArray)
		assert.Equal(t, true, *describeResult.Options.JSONStripNullValues)
		assert.Equal(t, true, *describeResult.Options.JSONIgnoreUTF8Errors)
		assert.Equal(t, true, *describeResult.Options.JSONSkipByteOrderMark)
	})
	t.Run("AVRO", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.Create(ctx, id, &sdk.CreateFileFormatOptions{
			Type: sdk.FileFormatTypeAvro,
			FileFormatTypeOptions: sdk.FileFormatTypeOptions{
				AvroCompression:              &sdk.AvroCompressionGzip,
				AvroTrimSpace:                sdk.Bool(true),
				AvroReplaceInvalidCharacters: sdk.Bool(true),
				AvroNullIf:                   &[]sdk.NullString{{S: "a"}, {S: "b"}},
			},
			Comment: sdk.String("test comment"),
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, sdk.FileFormatTypeAvro, result.Type)
		assert.Equal(t, client.GetConfig().Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "ROLE", result.OwnerRoleType)

		assert.Equal(t, sdk.AvroCompressionGzip, *result.Options.AvroCompression)
		assert.Equal(t, true, *result.Options.AvroTrimSpace)
		assert.Equal(t, true, *result.Options.AvroReplaceInvalidCharacters)
		assert.Equal(t, []sdk.NullString{{S: "a"}, {S: "b"}}, *result.Options.AvroNullIf)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeAvro, describeResult.Type)
		assert.Equal(t, sdk.AvroCompressionGzip, *describeResult.Options.AvroCompression)
		assert.Equal(t, true, *describeResult.Options.AvroTrimSpace)
		assert.Equal(t, true, *describeResult.Options.AvroReplaceInvalidCharacters)
		assert.Equal(t, []sdk.NullString{{S: "a"}, {S: "b"}}, *describeResult.Options.AvroNullIf)
	})
	t.Run("ORC", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.Create(ctx, id, &sdk.CreateFileFormatOptions{
			Type: sdk.FileFormatTypeORC,
			FileFormatTypeOptions: sdk.FileFormatTypeOptions{
				ORCTrimSpace:                sdk.Bool(true),
				ORCReplaceInvalidCharacters: sdk.Bool(true),
				ORCNullIf:                   &[]sdk.NullString{{S: "a"}, {S: "b"}},
			},
			Comment: sdk.String("test comment"),
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, sdk.FileFormatTypeORC, result.Type)
		assert.Equal(t, client.GetConfig().Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "ROLE", result.OwnerRoleType)

		assert.Equal(t, true, *result.Options.ORCTrimSpace)
		assert.Equal(t, true, *result.Options.ORCReplaceInvalidCharacters)
		assert.Equal(t, []sdk.NullString{{S: "a"}, {S: "b"}}, *result.Options.ORCNullIf)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeORC, describeResult.Type)
		assert.Equal(t, true, *describeResult.Options.ORCTrimSpace)
		assert.Equal(t, true, *describeResult.Options.ORCReplaceInvalidCharacters)
		assert.Equal(t, []sdk.NullString{{S: "a"}, {S: "b"}}, *describeResult.Options.ORCNullIf)
	})
	t.Run("PARQUET", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.Create(ctx, id, &sdk.CreateFileFormatOptions{
			Type: sdk.FileFormatTypeParquet,
			FileFormatTypeOptions: sdk.FileFormatTypeOptions{
				ParquetCompression:              &sdk.ParquetCompressionLzo,
				ParquetBinaryAsText:             sdk.Bool(true),
				ParquetTrimSpace:                sdk.Bool(true),
				ParquetReplaceInvalidCharacters: sdk.Bool(true),
				ParquetNullIf:                   &[]sdk.NullString{{S: "a"}, {S: "b"}},
			},
			Comment: sdk.String("test comment"),
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, sdk.FileFormatTypeParquet, result.Type)
		assert.Equal(t, client.GetConfig().Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "ROLE", result.OwnerRoleType)

		assert.Equal(t, sdk.ParquetCompressionLzo, *result.Options.ParquetCompression)
		assert.Equal(t, true, *result.Options.ParquetBinaryAsText)
		assert.Equal(t, true, *result.Options.ParquetTrimSpace)
		assert.Equal(t, true, *result.Options.ParquetReplaceInvalidCharacters)
		assert.Equal(t, []sdk.NullString{{S: "a"}, {S: "b"}}, *result.Options.ParquetNullIf)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeParquet, describeResult.Type)
		assert.Equal(t, sdk.ParquetCompressionLzo, *describeResult.Options.ParquetCompression)
		assert.Equal(t, true, *describeResult.Options.ParquetBinaryAsText)
		assert.Equal(t, true, *describeResult.Options.ParquetTrimSpace)
		assert.Equal(t, true, *describeResult.Options.ParquetReplaceInvalidCharacters)
		assert.Equal(t, []sdk.NullString{{S: "a"}, {S: "b"}}, *describeResult.Options.ParquetNullIf)
	})
	t.Run("XML", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.FileFormats.Create(ctx, id, &sdk.CreateFileFormatOptions{
			Type: sdk.FileFormatTypeXML,
			FileFormatTypeOptions: sdk.FileFormatTypeOptions{
				XMLCompression:          &sdk.XMLCompressionDeflate,
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
			err := client.FileFormats.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, sdk.FileFormatTypeXML, result.Type)
		assert.Equal(t, client.GetConfig().Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "ROLE", result.OwnerRoleType)

		assert.Equal(t, sdk.XMLCompressionDeflate, *result.Options.XMLCompression)
		assert.Equal(t, true, *result.Options.XMLIgnoreUTF8Errors)
		assert.Equal(t, true, *result.Options.XMLPreserveSpace)
		assert.Equal(t, true, *result.Options.XMLStripOuterElement)
		assert.Equal(t, true, *result.Options.XMLDisableSnowflakeData)
		assert.Equal(t, true, *result.Options.XMLDisableAutoConvert)
		assert.Equal(t, true, *result.Options.XMLSkipByteOrderMark)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeXML, describeResult.Type)
		assert.Equal(t, sdk.XMLCompressionDeflate, *describeResult.Options.XMLCompression)
		assert.Equal(t, true, *describeResult.Options.XMLIgnoreUTF8Errors)
		assert.Equal(t, true, *describeResult.Options.XMLPreserveSpace)
		assert.Equal(t, true, *describeResult.Options.XMLStripOuterElement)
		assert.Equal(t, true, *describeResult.Options.XMLDisableSnowflakeData)
		assert.Equal(t, true, *describeResult.Options.XMLDisableAutoConvert)
		assert.Equal(t, true, *describeResult.Options.XMLSkipByteOrderMark)
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

		err := client.FileFormats.Alter(ctx, oldId, &sdk.AlterFileFormatOptions{
			Rename: &sdk.AlterFileFormatRenameOptions{
				NewName: newId,
			},
		})
		require.NoError(t, err)

		_, err = client.FileFormats.ShowByID(ctx, oldId)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)

		result, err := client.FileFormats.ShowByID(ctx, newId)
		require.NoError(t, err)
		assert.Equal(t, newId, result.Name)

		// Undo rename so we can clean up
		err = client.FileFormats.Alter(ctx, newId, &sdk.AlterFileFormatOptions{
			Rename: &sdk.AlterFileFormatRenameOptions{
				NewName: oldId,
			},
		})
		require.NoError(t, err)
	})

	t.Run("set + set comment", func(t *testing.T) {
		fileFormat, fileFormatCleanup := testClientHelper().FileFormat.CreateFileFormatWithOptions(t, &sdk.CreateFileFormatOptions{
			Type: sdk.FileFormatTypeCSV,
			FileFormatTypeOptions: sdk.FileFormatTypeOptions{
				CSVCompression: &sdk.CSVCompressionAuto,
				CSVParseHeader: sdk.Bool(false),
			},
		})
		t.Cleanup(fileFormatCleanup)

		err := client.FileFormats.Alter(ctx, fileFormat.ID(), &sdk.AlterFileFormatOptions{
			Set: &sdk.FileFormatTypeOptions{
				Comment:        sdk.String("some comment"),
				CSVCompression: &sdk.CSVCompressionBz2,
				CSVParseHeader: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		result, err := client.FileFormats.ShowByID(ctx, fileFormat.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.CSVCompressionBz2, *result.Options.CSVCompression)
		assert.Equal(t, true, *result.Options.CSVParseHeader)
		assert.Equal(t, "some comment", result.Comment)
	})
}

func TestInt_FileFormatsDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("no options", func(t *testing.T) {
		fileFormat, fileFormatCleanup := testClientHelper().FileFormat.CreateFileFormat(t)
		t.Cleanup(fileFormatCleanup)

		err := client.FileFormats.Drop(ctx, fileFormat.ID(), nil)
		require.NoError(t, err)

		_, err = client.FileFormats.ShowByID(ctx, fileFormat.ID())
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("with IfExists", func(t *testing.T) {
		fileFormat, fileFormatCleanup := testClientHelper().FileFormat.CreateFileFormat(t)
		t.Cleanup(fileFormatCleanup)

		err := client.FileFormats.Drop(ctx, fileFormat.ID(), &sdk.DropFileFormatOptions{
			IfExists: sdk.Bool(true),
		})
		require.NoError(t, err)

		_, err = client.FileFormats.ShowByID(ctx, fileFormat.ID())
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
		fileFormats, err := client.FileFormats.Show(ctx, nil)
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(fileFormats))
		assert.Contains(t, fileFormats, *fileFormatTest)
		assert.Contains(t, fileFormats, *fileFormatTest2)
	})

	t.Run("LIKE", func(t *testing.T) {
		fileFormats, err := client.FileFormats.Show(ctx, &sdk.ShowFileFormatsOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(fileFormatTest.Name.Name()),
			},
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, 1, len(fileFormats))
		assert.Contains(t, fileFormats, *fileFormatTest)
	})

	t.Run("IN", func(t *testing.T) {
		fileFormats, err := client.FileFormats.Show(ctx, &sdk.ShowFileFormatsOptions{
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
		err := client.Sessions.UseDatabase(ctx, databaseTest2.ID())
		require.NoError(t, err)
		err = client.Sessions.UseSchema(ctx, schemaTest2.ID())
		require.NoError(t, err)

		fileFormat, err := client.FileFormats.ShowByID(ctx, fileFormatTest.ID())
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
			err := client.FileFormats.Drop(ctx, id, nil)
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createFileFormatHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		err := client.FileFormats.Create(ctx, id, &sdk.CreateFileFormatOptions{Type: sdk.FileFormatTypeCSV})
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

		e1, err := client.FileFormats.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.FileFormats.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}

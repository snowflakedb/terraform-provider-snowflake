//go:build non_account_level_tests

package testint

import (
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_FileFormats(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("create, show, describe, alter, drop - CSV", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		createRequest := sdk.NewCreateFileFormatRequest(id, sdk.FileFormatTypeCsv).
			WithFileFormatObjectOptions(*sdk.NewFileFormatObjectOptionsRequest().
				WithCsvCompression(sdk.CsvCompressionBz2).
				WithCsvParseHeader(true).
				WithCsvFieldOptionallyEnclosedBy(*sdk.NewStageFileFormatStringOrNoneRequest().WithValue("\"")).
				WithCsvSkipHeader(1)).
			WithComment("test comment")
		err := client.FileFormats.Create(ctx, createRequest)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.DropSafely(ctx, id)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 1*time.Minute)
		assert.Equal(t, sdk.FileFormatTypeCsv, result.Type)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, sdk.Pointer(sdk.CsvCompressionBz2), result.Options.CsvCompression)
		assert.Equal(t, true, *result.Options.CsvParseHeader)
		assert.Equal(t, 1, *result.Options.CsvSkipHeader)

		details, err := client.FileFormats.DescribeDetails(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeCsv, details.Type)
		assert.Equal(t, sdk.Pointer(sdk.CsvCompressionBz2), details.Options.CsvCompression)
		assert.Equal(t, true, *details.Options.CsvParseHeader)
		assert.Equal(t, 1, *details.Options.CsvSkipHeader)

		alterRequest := sdk.NewAlterFileFormatRequest(id).
			WithSet(*sdk.NewFileFormatObjectOptionsRequest().
				WithComment("updated comment").
				WithCsvCompression(sdk.CsvCompressionGzip))
		err = client.FileFormats.Alter(ctx, alterRequest)
		require.NoError(t, err)

		result, err = client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "updated comment", result.Comment)
		assert.Equal(t, sdk.Pointer(sdk.CsvCompressionGzip), result.Options.CsvCompression)

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err = client.FileFormats.Alter(ctx, sdk.NewAlterFileFormatRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.DropSafely(ctx, newId)
			require.NoError(t, err)
		})

		_, err = client.FileFormats.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)

		result, err = client.FileFormats.ShowByID(ctx, newId)
		require.NoError(t, err)
		assert.Equal(t, newId.Name(), result.Name)
	})

	t.Run("create, show, describe, drop - Parquet", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		createRequest := sdk.NewCreateFileFormatRequest(id, sdk.FileFormatTypeParquet).
			WithFileFormatObjectOptions(*sdk.NewFileFormatObjectOptionsRequest().
				WithParquetCompression(sdk.ParquetCompressionLzo).
				WithParquetTrimSpace(true))
		err := client.FileFormats.Create(ctx, createRequest)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.DropSafely(ctx, id)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeParquet, result.Type)
		assert.Equal(t, sdk.Pointer(sdk.ParquetCompressionLzo), result.Options.ParquetCompression)
		assert.Equal(t, true, *result.Options.ParquetTrimSpace)

		details, err := client.FileFormats.DescribeDetails(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.FileFormatTypeParquet, details.Type)
		assert.Equal(t, sdk.Pointer(sdk.ParquetCompressionLzo), details.Options.ParquetCompression)
		assert.Equal(t, true, *details.Options.ParquetTrimSpace)

		err = client.FileFormats.Drop(ctx, sdk.NewDropFileFormatRequest(id))
		require.NoError(t, err)

		_, err = client.FileFormats.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})
}

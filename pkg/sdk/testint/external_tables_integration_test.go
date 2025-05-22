//go:build !account_level_tests

package testint

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ExternalTables(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	stage, stageCleanup := testClientHelper().Stage.CreateStageWithURL(t)
	t.Cleanup(stageCleanup)

	tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	defaultColumns := func() []*sdk.ExternalTableColumnRequest {
		return []*sdk.ExternalTableColumnRequest{
			sdk.NewExternalTableColumnRequest("filename", sdk.DataTypeString, "metadata$filename::string"),
			sdk.NewExternalTableColumnRequest("city", sdk.DataTypeString, "value:city:findname::string"),
			sdk.NewExternalTableColumnRequest("time", sdk.DataTypeTimestampLTZ, "to_timestamp_ltz(value:time::int)"),
			sdk.NewExternalTableColumnRequest("weather", sdk.DataTypeVariant, "value:weather::variant"),
		}
	}

	columns := defaultColumns()
	columnsWithPartition := append(defaultColumns(), []*sdk.ExternalTableColumnRequest{
		sdk.NewExternalTableColumnRequest("weather_date", sdk.DataTypeDate, "to_date(to_timestamp(value:time::int))"),
		sdk.NewExternalTableColumnRequest("part_date", sdk.DataTypeDate, "parse_json(metadata$external_table_partition):weather_date::date"),
	}...)

	minimalCreateExternalTableReq := func(id sdk.SchemaObjectIdentifier) *sdk.CreateExternalTableRequest {
		return sdk.NewCreateExternalTableRequest(
			id,
			stage.Location(),
		).WithFileFormat(*sdk.NewExternalTableFileFormatRequest().WithFileFormatType(sdk.ExternalTableFileFormatTypeJSON))
	}

	createExternalTableWithManualPartitioningReq := func(id sdk.SchemaObjectIdentifier) *sdk.CreateWithManualPartitioningExternalTableRequest {
		return sdk.NewCreateWithManualPartitioningExternalTableRequest(
			id,
			stage.Location(),
		).
			WithFileFormat(*sdk.NewExternalTableFileFormatRequest().WithFileFormatType(sdk.ExternalTableFileFormatTypeJSON)).
			WithOrReplace(true).
			WithColumns(columnsWithPartition).
			WithPartitionBy([]string{"part_date"}).
			WithCopyGrants(true).
			WithComment("some_comment").
			WithTag([]*sdk.TagAssociationRequest{sdk.NewTagAssociationRequest(tag.ID(), "tag-value")})
	}

	t.Run("Create: minimal", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
		require.NoError(t, err)
		assert.Equal(t, externalTableID.Name(), externalTable.Name)
	})

	t.Run("Create: with raw file format", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ExternalTables.Create(ctx, sdk.NewCreateExternalTableRequest(externalTableID, stage.Location()).WithRawFileFormat("TYPE = JSON"))
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
		require.NoError(t, err)
		assert.Equal(t, externalTableID.Name(), externalTable.Name)
	})

	t.Run("Create: complete", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := externalTableID.Name()
		err := client.ExternalTables.Create(
			ctx,
			sdk.NewCreateExternalTableRequest(
				externalTableID,
				stage.Location(),
			).
				WithFileFormat(*sdk.NewExternalTableFileFormatRequest().WithFileFormatType(sdk.ExternalTableFileFormatTypeJSON)).
				WithOrReplace(true).
				WithColumns(columns).
				WithPartitionBy([]string{"filename"}).
				WithRefreshOnCreate(false).
				WithAutoRefresh(false).
				WithPattern("weather-nyc/weather_2_3_0.json.gz").
				WithCopyGrants(true).
				WithComment("some_comment").
				WithTag([]*sdk.TagAssociationRequest{sdk.NewTagAssociationRequest(tag.ID(), "tag-value")}),
		)
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
		require.NoError(t, err)
		assert.Equal(t, name, externalTable.Name)
	})

	t.Run("Create: infer schema", func(t *testing.T) {
		fileFormat, fileFormatCleanup := testClientHelper().FileFormat.CreateFileFormat(t)
		t.Cleanup(fileFormatCleanup)

		err := client.Sessions.UseWarehouse(ctx, testClientHelper().Ids.WarehouseId())
		require.NoError(t, err)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		query := fmt.Sprintf(`SELECT ARRAY_AGG(OBJECT_CONSTRUCT(*)) WITHIN GROUP (ORDER BY order_id) FROM TABLE (INFER_SCHEMA(location => '%s', FILE_FORMAT=>'%s', ignore_case => true))`, stage.Location(), fileFormat.ID().FullyQualifiedName())
		err = client.ExternalTables.CreateUsingTemplate(
			ctx,
			sdk.NewCreateExternalTableUsingTemplateRequest(
				id,
				stage.Location(),
			).
				WithFileFormat(*sdk.NewExternalTableFileFormatRequest().WithName(fileFormat.ID().FullyQualifiedName())).
				WithQuery(query).
				WithAutoRefresh(false))
		require.NoError(t, err)

		_, err = client.ExternalTables.ShowByID(ctx, id)
		require.NoError(t, err)
	})

	t.Run("Create with manual partitioning: complete", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := externalTableID.Name()
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, createExternalTableWithManualPartitioningReq(externalTableID))
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
		require.NoError(t, err)
		assert.Equal(t, name, externalTable.Name)
	})

	t.Run("Create delta lake: complete", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := externalTableID.Name()
		err := client.ExternalTables.CreateDeltaLake(
			ctx,
			sdk.NewCreateDeltaLakeExternalTableRequest(
				externalTableID,
				stage.Location(),
			).
				WithFileFormat(*sdk.NewExternalTableFileFormatRequest().WithFileFormatType(sdk.ExternalTableFileFormatTypeParquet)).
				WithOrReplace(true).
				WithColumns(columnsWithPartition).
				WithPartitionBy([]string{"filename"}).
				WithAutoRefresh(false).
				WithRefreshOnCreate(false).
				WithCopyGrants(true).
				WithComment("some_comment").
				WithTag([]*sdk.TagAssociationRequest{sdk.NewTagAssociationRequest(tag.ID(), "tag-value")}),
		)
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
		require.NoError(t, err)
		assert.Equal(t, name, externalTable.Name)
	})

	t.Run("Alter: refresh", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			sdk.NewAlterExternalTableRequest(externalTableID).
				WithIfExists(true).
				WithRefresh(*sdk.NewRefreshExternalTableRequest("weather-nyc")),
		)
		require.NoError(t, err)
	})

	t.Run("Alter: add files", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ExternalTables.Create(
			ctx,
			minimalCreateExternalTableReq(externalTableID).
				WithPattern("weather-nyc/weather_2_3_0.json.gz"),
		)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			sdk.NewAlterExternalTableRequest(externalTableID).
				WithIfExists(true).
				WithAddFiles([]*sdk.ExternalTableFileRequest{sdk.NewExternalTableFileRequest("weather-nyc/weather_0_0_0.json.gz")}),
		)
		require.NoError(t, err)
	})

	t.Run("Alter: remove files", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ExternalTables.Create(
			ctx,
			minimalCreateExternalTableReq(externalTableID).
				WithPattern("weather-nyc/weather_2_3_0.json.gz"),
		)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			sdk.NewAlterExternalTableRequest(externalTableID).
				WithIfExists(true).
				WithAddFiles([]*sdk.ExternalTableFileRequest{sdk.NewExternalTableFileRequest("weather-nyc/weather_0_0_0.json.gz")}),
		)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			sdk.NewAlterExternalTableRequest(externalTableID).
				WithIfExists(true).
				WithRemoveFiles([]*sdk.ExternalTableFileRequest{sdk.NewExternalTableFileRequest("weather-nyc/weather_0_0_0.json.gz")}),
		)
		require.NoError(t, err)
	})

	t.Run("Alter: set auto refresh", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			sdk.NewAlterExternalTableRequest(externalTableID).
				WithIfExists(true).
				WithAutoRefresh(true),
		)
		require.NoError(t, err)
	})

	t.Run("Alter: add partitions", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, createExternalTableWithManualPartitioningReq(externalTableID))
		require.NoError(t, err)

		err = client.ExternalTables.AlterPartitions(
			ctx,
			sdk.NewAlterExternalTablePartitionRequest(externalTableID).
				WithIfExists(true).
				WithAddPartitions([]*sdk.PartitionRequest{sdk.NewPartitionRequest("part_date", "2019-06-25")}).
				WithLocation("2019/06"),
		)
		require.NoError(t, err)
	})

	t.Run("Alter: drop partitions", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, createExternalTableWithManualPartitioningReq(externalTableID))
		require.NoError(t, err)

		err = client.ExternalTables.AlterPartitions(
			ctx,
			sdk.NewAlterExternalTablePartitionRequest(externalTableID).
				WithIfExists(true).
				WithAddPartitions([]*sdk.PartitionRequest{sdk.NewPartitionRequest("part_date", "2019-06-25")}).
				WithLocation("2019/06"),
		)
		require.NoError(t, err)

		err = client.ExternalTables.AlterPartitions(
			ctx,
			sdk.NewAlterExternalTablePartitionRequest(externalTableID).
				WithIfExists(true).
				WithDropPartition(true).
				WithLocation("2019/06"),
		)
		require.NoError(t, err)
	})

	t.Run("Drop", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		err = client.ExternalTables.Drop(
			ctx,
			sdk.NewDropExternalTableRequest(externalTableID).
				WithIfExists(true).
				WithDropOption(*sdk.NewExternalTableDropOptionRequest().WithCascade(true)),
		)
		require.NoError(t, err)

		_, err = client.ExternalTables.ShowByID(ctx, externalTableID)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("Show", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := externalTableID.Name()
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		et, err := client.ExternalTables.Show(
			ctx,
			sdk.NewShowExternalTableRequest().
				WithTerse(true).
				WithLike(name).
				WithIn(*sdk.NewShowExternalTableInRequest().WithDatabase(testClientHelper().Ids.DatabaseId())).
				WithStartsWith(name).
				WithLimitFrom(*sdk.NewLimitFromRequest().WithRows(sdk.Int(1))),
		)
		require.NoError(t, err)
		assert.Equal(t, 1, len(et))
		assert.Equal(t, externalTableID, et[0].ID())
	})

	t.Run("Describe: columns", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := minimalCreateExternalTableReq(externalTableID)
		err := client.ExternalTables.Create(ctx, req)
		require.NoError(t, err)

		d, err := client.ExternalTables.DescribeColumns(ctx, sdk.NewDescribeExternalTableColumnsRequest(externalTableID))
		require.NoError(t, err)

		assert.Equal(t, len(req.GetColumns())+1, len(d)) // +1 because there's underlying Value column
		assert.Contains(t, d, sdk.ExternalTableColumnDetails{
			Name:       "VALUE",
			Type:       "VARIANT",
			Kind:       "COLUMN",
			IsNullable: true,
			Default:    nil,
			IsPrimary:  false,
			IsUnique:   false,
			Check:      nil,
			Expression: nil,
			Comment:    sdk.String("The value of this row"),
			PolicyName: nil,
		})
	})

	t.Run("Describe: stage", func(t *testing.T) {
		externalTableID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		d, err := client.ExternalTables.DescribeStage(ctx, sdk.NewDescribeExternalTableStageRequest(externalTableID))
		require.NoError(t, err)

		assert.Contains(t, d, sdk.ExternalTableStageDetails{
			ParentProperty:  "STAGE_FILE_FORMAT",
			Property:        "TIME_FORMAT",
			PropertyType:    "String",
			PropertyValue:   "AUTO",
			PropertyDefault: "AUTO",
		})
	})
}

func TestInt_ExternalTablesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	stage, stageCleanup := testClientHelper().Stage.CreateStageWithURL(t)
	t.Cleanup(stageCleanup)

	cleanupExternalTableHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.ExternalTables.Drop(ctx, sdk.NewDropExternalTableRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createExternalTableHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		request := sdk.NewCreateExternalTableRequest(id, stage.Location()).WithFileFormat(*sdk.NewExternalTableFileFormatRequest().WithFileFormatType(sdk.ExternalTableFileFormatTypeJSON))
		err := client.ExternalTables.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupExternalTableHandle(t, id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createExternalTableHandle(t, id1)
		createExternalTableHandle(t, id2)

		e1, err := client.ExternalTables.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.ExternalTables.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}

//go:build non_account_level_tests

package testint

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/ids"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_Stages(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)
	gcsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBucketUrl)
	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	azureSasToken := testenvs.GetOrSkipTest(t, testenvs.AzureExternalSasToken)

	s3StorageIntegration, err := client.StorageIntegrations.ShowByID(ctx, ids.PrecreatedS3StorageIntegration)
	require.NoError(t, err)
	gcpStorageIntegration, err := client.StorageIntegrations.ShowByID(ctx, ids.PrecreatedGcpStorageIntegration)
	require.NoError(t, err)
	azureStorageIntegration, err := client.StorageIntegrations.ShowByID(ctx, ids.PrecreatedAzureStorageIntegration)
	require.NoError(t, err)

	// ==================== INTERNAL STAGE TESTS ====================

	t.Run("CreateInternal - minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateInternalStageRequest(id)

		err := client.Stages.CreateInternal(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasType("INTERNAL").
			HasComment("").
			HasUrl("").
			HasDirectoryEnabled(false).
			HasHasCredentials(false).
			HasHasEncryptionKey(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))

		assertThatObject(t, objectassert.StageDetails(t, id).
			HasStageLocation(sdk.StageLocationDetails{
				Url: "",
			}).
			HasDirectoryTableEnable(false).
			HasDirectoryTableAutoRefresh(false).
			HasDirectoryTableNotificationChannelEmpty().
			HasDirectoryTableLastRefreshedOnNil(),
		)
	})

	t.Run("CreateInternal - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "test comment"

		fileFormat, fileFormatCleanup := testClientHelper().FileFormat.CreateFileFormat(t)
		t.Cleanup(fileFormatCleanup)

		request := sdk.NewCreateInternalStageRequest(id).
			WithIfNotExists(true).
			WithEncryption(*sdk.NewInternalStageEncryptionRequest().
				WithSnowflakeFull(*sdk.NewInternalStageEncryptionSnowflakeFullRequest())).
			WithDirectoryTableOptions(*sdk.NewInternalDirectoryTableOptionsRequest().
				WithEnable(true).
				WithAutoRefresh(true)).
			WithFileFormat(*sdk.NewStageFileFormatRequest().
				WithFormatName(fileFormat.ID())).
			WithComment(comment)

		err := client.Stages.CreateInternal(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		// Encryption type are not asserted because it's missing from Snowflake.
		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasType("INTERNAL").
			HasComment(comment).
			HasDirectoryEnabled(true).
			HasHasEncryptionKey(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))

		assertThatObject(t, objectassert.StageDetails(t, id).
			HasFileFormatName(fileFormat.ID()).
			HasDirectoryTableEnable(true).
			HasDirectoryTableAutoRefresh(true).
			HasDirectoryTableLastRefreshedOnNotEmpty())
	})

	t.Run("CreateInternal - temporary and or replace", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		request := sdk.NewCreateInternalStageRequest(id).
			WithTemporary(true).
			WithOrReplace(true)

		err := client.Stages.CreateInternal(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasType("INTERNAL TEMPORARY"),
		)
	})

	t.Run("CreateInternal - minimal with CSV file format", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		fileFormat := sdk.NewStageFileFormatRequest().
			WithFileFormatOptions(sdk.FileFormatOptions{
				CsvOptions: &sdk.FileFormatCsvOptions{},
			})

		request := sdk.NewCreateInternalStageRequest(id).
			WithFileFormat(*fileFormat)

		err := client.Stages.CreateInternal(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.StageDetails(t, id).
			HasFileFormatCsv(sdk.FileFormatCsv{
				Type:                       "CSV",
				RecordDelimiter:            "\\n",
				FieldDelimiter:             ",",
				FileExtension:              "",
				SkipHeader:                 0,
				ParseHeader:                false,
				DateFormat:                 "AUTO",
				TimeFormat:                 "AUTO",
				TimestampFormat:            "AUTO",
				BinaryFormat:               "HEX",
				Escape:                     "NONE",
				EscapeUnenclosedField:      `\\`,
				TrimSpace:                  false,
				FieldOptionallyEnclosedBy:  "NONE",
				NullIf:                     []string{`\\N`},
				Compression:                "AUTO",
				ErrorOnColumnCountMismatch: true,
				ValidateUtf8:               true,
				SkipBlankLines:             false,
				ReplaceInvalidCharacters:   false,
				EmptyFieldAsNull:           true,
				SkipByteOrderMark:          true,
				Encoding:                   "UTF8",
				MultiLine:                  true,
			}))
	})

	t.Run("CreateInternal - complete with CSV file format options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		gzipCompression := sdk.CSVCompressionGzip
		base64Format := sdk.BinaryFormatBase64
		utf8Encoding := sdk.CSVEncodingUTF8
		multiLine := false
		fileExtension := ".csv"
		skipHeader := 2
		skipBlankLines := true
		dateFormat := "YYYY-MM-DD"
		timeFormat := "HH24:MI:SS"
		timestampFormat := "YYYY-MM-DD HH24:MI:SS"
		escapeVal := `\\`
		trimSpace := true
		enclosedBy := `'`
		errorOnMismatch := false
		replaceInvalid := true
		emptyAsNull := false
		skipBom := false

		fileFormat := sdk.NewStageFileFormatRequest().
			WithFileFormatOptions(sdk.FileFormatOptions{
				CsvOptions: &sdk.FileFormatCsvOptions{
					Compression:                &gzipCompression,
					RecordDelimiter:            &sdk.StageFileFormatStringOrNone{Value: sdk.String("\\n")},
					FieldDelimiter:             &sdk.StageFileFormatStringOrNone{Value: sdk.String("|")},
					MultiLine:                  &multiLine,
					FileExtension:              &fileExtension,
					SkipHeader:                 &skipHeader,
					SkipBlankLines:             &skipBlankLines,
					DateFormat:                 &sdk.StageFileFormatStringOrAuto{Value: &dateFormat},
					TimeFormat:                 &sdk.StageFileFormatStringOrAuto{Value: &timeFormat},
					TimestampFormat:            &sdk.StageFileFormatStringOrAuto{Value: &timestampFormat},
					BinaryFormat:               &base64Format,
					Escape:                     &sdk.StageFileFormatStringOrNone{Value: &escapeVal},
					EscapeUnenclosedField:      &sdk.StageFileFormatStringOrNone{Value: &escapeVal},
					TrimSpace:                  &trimSpace,
					FieldOptionallyEnclosedBy:  &sdk.StageFileFormatStringOrNone{Value: &enclosedBy},
					NullIf:                     []sdk.NullString{{S: "NULL"}, {S: ""}},
					ErrorOnColumnCountMismatch: &errorOnMismatch,
					ReplaceInvalidCharacters:   &replaceInvalid,
					EmptyFieldAsNull:           &emptyAsNull,
					SkipByteOrderMark:          &skipBom,
					Encoding:                   &utf8Encoding,
				},
			})

		request := sdk.NewCreateInternalStageRequest(id).
			WithFileFormat(*fileFormat)

		err := client.Stages.CreateInternal(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.StageDetails(t, id).
			HasFileFormatCsv(sdk.FileFormatCsv{
				Type:                       "CSV",
				RecordDelimiter:            "\\n",
				FieldDelimiter:             "|",
				FileExtension:              ".csv",
				SkipHeader:                 2,
				ParseHeader:                false,
				DateFormat:                 "YYYY-MM-DD",
				TimeFormat:                 timeFormat,
				TimestampFormat:            timestampFormat,
				BinaryFormat:               "BASE64",
				Escape:                     escapeVal,
				EscapeUnenclosedField:      escapeVal,
				TrimSpace:                  true,
				FieldOptionallyEnclosedBy:  enclosedBy,
				NullIf:                     []string{"NULL", ""},
				Compression:                "GZIP",
				ErrorOnColumnCountMismatch: false,
				ValidateUtf8:               true,
				SkipBlankLines:             true,
				ReplaceInvalidCharacters:   true,
				EmptyFieldAsNull:           false,
				SkipByteOrderMark:          false,
				Encoding:                   "UTF8",
				MultiLine:                  false,
			}))
	})

	t.Run("CreateInternal - complete with CSV file format options; auto and none", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		fileFormat := sdk.NewStageFileFormatRequest().
			WithFileFormatOptions(sdk.FileFormatOptions{
				CsvOptions: &sdk.FileFormatCsvOptions{
					RecordDelimiter:           &sdk.StageFileFormatStringOrNone{None: sdk.Bool(true)},
					FieldDelimiter:            &sdk.StageFileFormatStringOrNone{None: sdk.Bool(true)},
					DateFormat:                &sdk.StageFileFormatStringOrAuto{Auto: sdk.Bool(true)},
					TimeFormat:                &sdk.StageFileFormatStringOrAuto{Auto: sdk.Bool(true)},
					TimestampFormat:           &sdk.StageFileFormatStringOrAuto{Auto: sdk.Bool(true)},
					Escape:                    &sdk.StageFileFormatStringOrNone{None: sdk.Bool(true)},
					EscapeUnenclosedField:     &sdk.StageFileFormatStringOrNone{None: sdk.Bool(true)},
					FieldOptionallyEnclosedBy: &sdk.StageFileFormatStringOrNone{None: sdk.Bool(true)},
				},
			})

		request := sdk.NewCreateInternalStageRequest(id).
			WithFileFormat(*fileFormat)

		err := client.Stages.CreateInternal(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.StageDetails(t, id).
			HasFileFormatCsv(sdk.FileFormatCsv{
				Type:                       "CSV",
				RecordDelimiter:            "NONE",
				FieldDelimiter:             "NONE",
				FileExtension:              "",
				SkipHeader:                 0,
				ParseHeader:                false,
				DateFormat:                 "AUTO",
				TimeFormat:                 "AUTO",
				TimestampFormat:            "AUTO",
				BinaryFormat:               "HEX",
				Escape:                     "NONE",
				EscapeUnenclosedField:      "NONE",
				TrimSpace:                  false,
				FieldOptionallyEnclosedBy:  "NONE",
				NullIf:                     []string{"\\\\N"},
				Compression:                "AUTO",
				ErrorOnColumnCountMismatch: true,
				ValidateUtf8:               true,
				SkipBlankLines:             false,
				ReplaceInvalidCharacters:   false,
				EmptyFieldAsNull:           true,
				SkipByteOrderMark:          true,
				Encoding:                   "UTF8",
				MultiLine:                  true,
			}))
	})
	t.Run("AlterInternalStage - complete", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasComment(""))

		err := client.Stages.AlterInternalStage(ctx, sdk.NewAlterInternalStageStageRequest(stage.ID()).
			WithIfExists(true).
			WithComment("altered comment"))
		require.NoError(t, err)

		stage, err = client.Stages.ShowByID(ctx, stage.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasComment("altered comment"))
	})

	// ==================== S3 STAGE TESTS ====================

	t.Run("CreateOnS3 - minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl)

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req)

		err := client.Stages.CreateOnS3(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasType("EXTERNAL").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasNoStorageIntegration().
			HasHasCredentials(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))
	})

	t.Run("CreateOnS3 - minimal with credentials", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithCredentials(*sdk.NewExternalStageS3CredentialsRequest().
				WithAwsKeyId(awsKeyId).
				WithAwsSecretKey(awsSecretKey).
				WithAwsToken("asdf"))

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req)

		err := client.Stages.CreateOnS3(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasHasCredentials(true).
			HasOwner(snowflakeroles.Accountadmin.Name()))
	})

	t.Run("CreateOnS3 - with credentials using AWS_ROLE", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithCredentials(*sdk.NewExternalStageS3CredentialsRequest().
				WithAwsRole("arn:aws:iam::123456789012:role/MyRole"))

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req)

		err := client.Stages.CreateOnS3(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		// AWS_ROLE is not returned by Snowflake in SHOW or DESCRIBE, so we can only verify the stage was created.
		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasHasCredentials(true).
			HasOwner(snowflakeroles.Accountadmin.Name()))
	})

	t.Run("CreateOnS3 - with AWS_CSE encryption", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		// It must be either 128 bits, 192 bits, or 256 bits long
		masterKey := random.AlphaN(256 / 8)

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithStorageIntegration(ids.PrecreatedS3StorageIntegration).
			WithEncryption(*sdk.NewExternalStageS3EncryptionRequest().
				WithAwsCse(*sdk.NewExternalStageS3EncryptionAwsCseRequest(masterKey)))

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req)

		err := client.Stages.CreateOnS3(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		// Encryption type and master key are not returned by Snowflake in SHOW or DESCRIBE.
		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasStorageIntegration(ids.PrecreatedS3StorageIntegration).
			HasHasEncryptionKey(true).
			HasOwner(snowflakeroles.Accountadmin.Name()))
	})

	t.Run("CreateOnS3 - minimal with private link", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithUsePrivatelinkEndpoint(true)

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req)

		err := client.Stages.CreateOnS3(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasHasCredentials(false).
			HasOwner(snowflakeroles.Accountadmin.Name()))

		assertThatObject(t, objectassert.StageDetails(t, id).
			HasPrivateLinkUsePrivatelinkEndpoint(true),
		)
	})

	t.Run("CreateOnS3 - complete with storage integration", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "complete s3 stage with credentials"

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithAwsAccessPointArn("arn:aws:s3:us-west-2:123456789012:accesspoint/my-data-ap").
			WithStorageIntegration(ids.PrecreatedS3StorageIntegration).
			WithEncryption(*sdk.NewExternalStageS3EncryptionRequest().
				WithAwsSseKms(*sdk.NewExternalStageS3EncryptionAwsSseKmsRequest().WithKmsKeyId(random.AlphaN(12))))

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req).
			WithIfNotExists(true).
			WithDirectoryTableOptions(*sdk.NewStageS3CommonDirectoryTableOptionsRequest().
				WithEnable(true).
				WithRefreshOnCreate(true).
				WithAutoRefresh(false)).
			WithComment(comment)

		err := client.Stages.CreateOnS3(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		// Encryption type, RefreshOnCreate, and other credentials fields are not asserted because it's missing from Snowflake.
		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasHasCredentials(false).
			HasStorageIntegration(ids.PrecreatedS3StorageIntegration).
			HasHasEncryptionKey(true).
			HasDirectoryEnabled(true).
			HasComment(comment))

		assertThatObject(t, objectassert.StageDetails(t, id).
			HasDirectoryTableEnable(true).
			HasDirectoryTableAutoRefresh(false).
			HasStageLocationAwsAccessPointArn("arn:aws:s3:us-west-2:123456789012:accesspoint/my-data-ap").
			HasPrivateLinkUsePrivatelinkEndpoint(false))
	})

	t.Run("CreateOnS3 - temporary and or replace", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithStorageIntegration(ids.PrecreatedS3StorageIntegration)

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req).
			WithOrReplace(true).
			WithTemporary(true).
			WithOrReplace(true)

		err := client.Stages.CreateOnS3(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasType("EXTERNAL TEMPORARY").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasStorageIntegration(s3StorageIntegration.ID()),
		)
	})

	t.Run("AlterExternalS3Stage - use privatelink", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStageOnS3WithCredentials(t)
		t.Cleanup(cleanup)

		require.Equal(t, "", stage.Comment)

		err := client.Stages.AlterExternalS3Stage(ctx, sdk.NewAlterExternalS3StageStageRequest(stage.ID()).
			WithExternalStageParams(*sdk.NewExternalS3StageParamsRequest(awsBucketUrl).WithUsePrivatelinkEndpoint(true)))
		require.NoError(t, err)

		stage, err = client.Stages.ShowByID(ctx, stage.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.StageDetails(t, stage.ID()).
			HasPrivateLinkUsePrivatelinkEndpoint(true),
		)
	})

	t.Run("AlterExternalS3Stage - complete", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStageOnS3WithCredentials(t)
		t.Cleanup(cleanup)

		require.Equal(t, "", stage.Comment)

		err := client.Stages.AlterExternalS3Stage(ctx, sdk.NewAlterExternalS3StageStageRequest(stage.ID()).
			WithIfExists(true).
			WithExternalStageParams(*sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
				WithStorageIntegration(ids.PrecreatedS3StorageIntegration).
				WithEncryption(*sdk.NewExternalStageS3EncryptionRequest().WithNone(*sdk.NewExternalStageS3EncryptionNoneRequest()))).
			WithComment("Updated comment"))
		require.NoError(t, err)

		stage, err = client.Stages.ShowByID(ctx, stage.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasType("EXTERNAL").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasStorageIntegration(s3StorageIntegration.ID()).
			HasComment("Updated comment"))
	})

	// ==================== GCS STAGE TESTS ====================

	t.Run("CreateOnGCS - minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		gcsReq := sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
			// We need to use a storage integration. Otherwise, we get:
			// Creation of stages with direct credentials, including accessing public stages, has been forbidden for GCS stages. See your account administrator for details.
			WithStorageIntegration(ids.PrecreatedGcpStorageIntegration)

		request := sdk.NewCreateOnGCSStageRequest(id, *gcsReq)

		err := client.Stages.CreateOnGCS(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasType("EXTERNAL").
			HasUrl(gcsBucketUrl).
			HasCloud("GCP").
			HasStorageIntegration(ids.PrecreatedGcpStorageIntegration).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))
	})

	t.Run("CreateOnGCS - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "complete gcs stage"
		kmsKeyId := random.AlphaN(12)

		gcsReq := sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
			WithStorageIntegration(ids.PrecreatedGcpStorageIntegration).
			WithEncryption(*sdk.NewExternalStageGCSEncryptionRequest().
				WithGcsSseKms(*sdk.NewExternalStageGCSEncryptionGcsSseKmsRequest().WithKmsKeyId(kmsKeyId)))

		request := sdk.NewCreateOnGCSStageRequest(id, *gcsReq).
			WithIfNotExists(true).
			WithDirectoryTableOptions(*sdk.NewExternalGCSDirectoryTableOptionsRequest().
				WithEnable(true).
				WithRefreshOnCreate(true).
				WithAutoRefresh(false)).
			WithComment(comment)

		err := client.Stages.CreateOnGCS(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		// Encryption type, RefreshOnCreate, and other credentials fields are not asserted because it's missing from Snowflake.
		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(gcsBucketUrl).
			HasCloud("GCP").
			HasStorageIntegration(gcpStorageIntegration.ID()).
			HasHasEncryptionKey(true).
			HasHasCredentials(false).
			HasDirectoryEnabled(true).
			HasComment(comment))

		assertThatObject(t, objectassert.StageDetails(t, id).
			HasDirectoryTableEnable(true).
			HasDirectoryTableAutoRefresh(false))
	})

	t.Run("CreateOnGCS - temporary", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "temporary gcs stage"

		gcsReq := sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
			WithStorageIntegration(ids.PrecreatedGcpStorageIntegration)

		request := sdk.NewCreateOnGCSStageRequest(id, *gcsReq).
			WithOrReplace(true).
			WithTemporary(true).
			WithComment(comment)

		err := client.Stages.CreateOnGCS(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasType("EXTERNAL TEMPORARY").
			HasUrl(gcsBucketUrl).
			HasCloud("GCP").
			HasStorageIntegration(gcpStorageIntegration.ID()).
			HasComment(comment))
	})

	t.Run("AlterExternalGCSStage - complete", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStageOnGCS(t)
		t.Cleanup(cleanup)

		require.Equal(t, "", stage.Comment)

		err := client.Stages.AlterExternalGCSStage(ctx, sdk.NewAlterExternalGCSStageStageRequest(stage.ID()).
			WithIfExists(true).
			WithExternalStageParams(*sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
				WithStorageIntegration(ids.PrecreatedGcpStorageIntegration).
				WithEncryption(*sdk.NewExternalStageGCSEncryptionRequest().WithNone(*sdk.NewExternalStageGCSEncryptionNoneRequest()))).
			WithComment("Updated comment"))
		require.NoError(t, err)

		stage, err = client.Stages.ShowByID(ctx, stage.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasType("EXTERNAL").
			HasUrl(gcsBucketUrl).
			HasCloud("GCP").
			HasStorageIntegration(gcpStorageIntegration.ID()).
			HasComment("Updated comment"))
	})

	// ==================== AZURE STAGE TESTS ====================

	t.Run("CreateOnAzure - minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		azureReq := sdk.NewExternalAzureStageParamsRequest(azureBucketUrl)

		request := sdk.NewCreateOnAzureStageRequest(id, *azureReq)

		err := client.Stages.CreateOnAzure(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(azureBucketUrl).
			HasCloud("AZURE").
			HasHasCredentials(false).
			HasOwner(snowflakeroles.Accountadmin.Name()))
	})

	t.Run("CreateOnAzure - minimal with storage integration", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		azureReq := sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
			WithStorageIntegration(ids.PrecreatedAzureStorageIntegration)

		request := sdk.NewCreateOnAzureStageRequest(id, *azureReq)

		err := client.Stages.CreateOnAzure(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasType("EXTERNAL").
			HasUrl(azureBucketUrl).
			HasCloud("AZURE").
			HasStorageIntegration(azureStorageIntegration.ID()).
			HasHasCredentials(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))
	})

	t.Run("CreateOnAzure - complete with credentials", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "complete azure stage with credentials"

		azureReq := sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
			WithCredentials(*sdk.NewExternalStageAzureCredentialsRequest(azureSasToken)).
			WithEncryption(*sdk.NewExternalStageAzureEncryptionRequest().
				WithAzureCse(*sdk.NewExternalStageAzureEncryptionAzureCseRequest(random.AlphaN(256 / 8))))

		request := sdk.NewCreateOnAzureStageRequest(id, *azureReq).
			WithIfNotExists(true).
			WithDirectoryTableOptions(*sdk.NewExternalAzureDirectoryTableOptionsRequest().
				WithEnable(true).
				WithRefreshOnCreate(false).
				WithAutoRefresh(false)).
			WithComment(comment)

		err := client.Stages.CreateOnAzure(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		// Encryption type, RefreshOnCreate, and other credentials fields are not asserted because it's missing from Snowflake.
		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(azureBucketUrl).
			HasCloud("AZURE").
			HasHasCredentials(true).
			HasHasEncryptionKey(true).
			HasDirectoryEnabled(true).
			HasComment(comment))

		assertThatObject(t, objectassert.StageDetails(t, id).
			HasDirectoryTableEnable(true).
			HasDirectoryTableAutoRefresh(false))
	})

	t.Run("CreateOnAzure - temporary", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "temporary azure stage"

		azureReq := sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
			WithStorageIntegration(ids.PrecreatedAzureStorageIntegration)

		request := sdk.NewCreateOnAzureStageRequest(id, *azureReq).
			WithOrReplace(true).
			WithTemporary(true).
			WithComment(comment)

		err := client.Stages.CreateOnAzure(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasType("EXTERNAL TEMPORARY").
			HasUrl(azureBucketUrl).
			HasCloud("AZURE").
			HasStorageIntegration(azureStorageIntegration.ID()).
			HasComment(comment))
	})

	t.Run("AlterExternalAzureStage - complete", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStageOnAzureWithCredentials(t)
		t.Cleanup(cleanup)

		require.Equal(t, "", stage.Comment)

		err := client.Stages.AlterExternalAzureStage(ctx, sdk.NewAlterExternalAzureStageStageRequest(stage.ID()).
			WithIfExists(true).
			WithExternalStageParams(*sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
				WithStorageIntegration(ids.PrecreatedAzureStorageIntegration).
				WithEncryption(*sdk.NewExternalStageAzureEncryptionRequest().WithNone(*sdk.NewExternalStageAzureEncryptionNoneRequest()))).
			WithComment("Updated comment"))
		require.NoError(t, err)

		stage, err = client.Stages.ShowByID(ctx, stage.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasType("EXTERNAL").
			HasUrl(azureBucketUrl).
			HasCloud("AZURE").
			HasStorageIntegration(azureStorageIntegration.ID()).
			HasComment("Updated comment"))
	})

	// ==================== S3-COMPATIBLE STAGE TESTS ====================

	t.Run("CreateOnS3Compatible - minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		compatibleBucketUrl := strings.Replace(awsBucketUrl, "s3://", "s3compat://", 1)
		endpoint := "s3.us-west-2.amazonaws.com"

		s3CompatReq := sdk.NewExternalS3CompatibleStageParamsRequest(compatibleBucketUrl, endpoint)

		request := sdk.NewCreateOnS3CompatibleStageRequest(id, *s3CompatReq)

		err := client.Stages.CreateOnS3Compatible(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasType("EXTERNAL").
			HasUrl(compatibleBucketUrl).
			HasCloud("AWS").
			HasEndpoint(endpoint).
			HasHasCredentials(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))
	})

	t.Run("CreateOnS3Compatible - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		compatibleBucketUrl := strings.Replace(awsBucketUrl, "s3://", "s3compat://", 1)
		endpoint := "s3.us-west-2.amazonaws.com"
		comment := "complete s3 compatible stage"

		s3CompatReq := sdk.NewExternalS3CompatibleStageParamsRequest(compatibleBucketUrl, endpoint).
			WithCredentials(*sdk.NewExternalStageS3CompatibleCredentialsRequest(awsKeyId, awsSecretKey))

		request := sdk.NewCreateOnS3CompatibleStageRequest(id, *s3CompatReq).
			WithIfNotExists(true).
			WithDirectoryTableOptions(*sdk.NewStageS3CommonDirectoryTableOptionsRequest().
				WithEnable(true).
				WithRefreshOnCreate(false).
				WithAutoRefresh(false)).
			WithComment(comment)

		err := client.Stages.CreateOnS3Compatible(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		// Encryption type, RefreshOnCreate, and other credentials fields are not asserted because it's missing from Snowflake.
		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(compatibleBucketUrl).
			HasCloud("AWS").
			HasEndpoint(endpoint).
			HasDirectoryEnabled(true).
			HasHasCredentials(true).
			HasHasEncryptionKey(false).
			HasComment(comment))

		assertThatObject(t, objectassert.StageDetails(t, id).
			HasDirectoryTableEnable(true).
			HasDirectoryTableAutoRefresh(false))
	})

	t.Run("CreateOnS3Compatible - temporary", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		compatibleBucketUrl := strings.Replace(awsBucketUrl, "s3://", "s3compat://", 1)
		endpoint := "s3.us-west-2.amazonaws.com"
		comment := "temporary s3 compatible stage"

		s3CompatReq := sdk.NewExternalS3CompatibleStageParamsRequest(compatibleBucketUrl, endpoint).
			WithCredentials(*sdk.NewExternalStageS3CompatibleCredentialsRequest(awsKeyId, awsSecretKey))

		request := sdk.NewCreateOnS3CompatibleStageRequest(id, *s3CompatReq).
			WithOrReplace(true).
			WithTemporary(true).
			WithComment(comment)

		err := client.Stages.CreateOnS3Compatible(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, id))

		assertThatObject(t, objectassert.Stage(t, id).
			HasName(id.Name()).
			HasType("EXTERNAL TEMPORARY").
			HasUrl(compatibleBucketUrl).
			HasCloud("AWS").
			HasEndpoint(endpoint).
			HasComment(comment))
	})

	// ==================== OTHER OPERATIONS ====================

	t.Run("Alter - rename", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStage(t)
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err = client.Stages.Alter(ctx, sdk.NewAlterStageRequest(stage.ID()).
			WithIfExists(true).
			WithRenameTo(newId))
		require.NoError(t, err)

		// Update cleanup to use new id
		t.Cleanup(func() {
			cleanup() // This will fail but we need to clean up with the new id
		})
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, newId))

		renamedStage, err := client.Stages.ShowByID(ctx, newId)
		require.NoError(t, err)
		require.NotNil(t, renamedStage)

		assertThatObject(t, objectassert.StageFromObject(t, renamedStage).
			HasName(newId.Name()))
	})

	t.Run("AlterDirectoryTable", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStageOnS3(t)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageDetails(t, stage.ID()).
			HasDirectoryTableEnable(false))

		err = client.Stages.AlterDirectoryTable(ctx, sdk.NewAlterDirectoryTableStageRequest(stage.ID()).
			WithSetDirectory(*sdk.NewDirectoryTableSetRequest(true)))
		require.NoError(t, err)

		err = client.Stages.AlterDirectoryTable(ctx, sdk.NewAlterDirectoryTableStageRequest(stage.ID()).
			WithRefresh(*sdk.NewDirectoryTableRefreshRequest().WithSubpath("/")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.StageDetails(t, stage.ID()).
			HasDirectoryTableEnable(true).
			HasDirectoryTableLastRefreshedOnNotEmpty())
	})

	t.Run("Drop", func(t *testing.T) {
		stage, _ := testClientHelper().Stage.CreateStage(t)

		foundStage, err := client.Stages.ShowByID(ctx, stage.ID())
		require.NotNil(t, foundStage)
		require.NoError(t, err)

		err = client.Stages.Drop(ctx, sdk.NewDropStageRequest(stage.ID()))
		require.NoError(t, err)

		foundStage, err = client.Stages.ShowByID(ctx, stage.ID())
		require.Nil(t, foundStage)
		require.Error(t, err)
	})

	t.Run("Describe internal", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageDetails(t, stage.ID()).
			HasDirectoryTableEnable(false).
			HasDirectoryTableAutoRefresh(false).
			HasStageLocationUrl(""))
	})

	t.Run("Describe external s3", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStageOnS3WithCredentials(t)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageDetails(t, stage.ID()).
			HasStageCredentialsAwsKeyId(awsKeyId).
			HasStageLocationUrl(fmt.Sprintf("[\"%s\"]", awsBucketUrl)))
	})

	t.Run("Describe external gcs", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStageOnGCS(t)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageDetails(t, stage.ID()).
			HasStageLocationUrl(fmt.Sprintf("[\"%s\"]", gcsBucketUrl)))
	})

	t.Run("Describe external azure", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStageOnAzureWithCredentials(t)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageDetails(t, stage.ID()).
			HasDirectoryTableEnable(false).
			HasStageLocationUrl(fmt.Sprintf("[\"%s\"]", azureBucketUrl)))
	})

	t.Run("Show internal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "show internal test"

		request := sdk.NewCreateInternalStageRequest(id).
			WithDirectoryTableOptions(*sdk.NewInternalDirectoryTableOptionsRequest().WithEnable(true)).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageWithRequest(t, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasUrl("").
			HasHasCredentials(false).
			HasHasEncryptionKey(false).
			HasComment(comment).
			HasType("INTERNAL").
			HasDirectoryEnabled(true).
			HasOwnerRoleType("ROLE"))

		assertThatObject(t, objectassert.StageDetails(t, id).
			HasDirectoryTableEnable(true))
	})
}

func TestInt_StagesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		stage1, cleanup1 := testClientHelper().Stage.CreateStageWithRequest(t, sdk.NewCreateInternalStageRequest(id1))
		t.Cleanup(cleanup1)
		stage2, cleanup2 := testClientHelper().Stage.CreateStageInSchema(t, schema.ID())
		t.Cleanup(cleanup2)

		// Re-create stage2 with the same name as stage1
		err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(stage2.ID()))
		require.NoError(t, err)
		stage2, cleanup2 = testClientHelper().Stage.CreateStageWithRequest(t, sdk.NewCreateInternalStageRequest(id2))
		t.Cleanup(cleanup2)

		e1, err := client.Stages.ShowByID(ctx, stage1.ID())
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Stages.ShowByID(ctx, stage2.ID())
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}

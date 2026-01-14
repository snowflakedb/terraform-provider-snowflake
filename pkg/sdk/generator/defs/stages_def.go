package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// TODO(SNOW-1019005): remove copy options
// TODO(SNOW-1019005): use a custom file format struct with a nice nesting
// TODO(SNOW-1019005): generate assertions
// TODO(SNOW-1019005): add parsers for DESC output and return a nice struct; use them in integration tests assertions
// TODO(SNOW-1019005): improve integration tests
func createStageOperation(structName string, apply func(qs *g.QueryStruct) *g.QueryStruct) *g.QueryStruct {
	qs := g.NewQueryStruct(structName).
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		SQL("STAGE").
		IfNotExists().
		Name()
	qs = apply(qs)
	return qs.
		OptionalQueryStructField("FileFormat", stageFileFormatDef, g.ListOptions().Parentheses().SQL("FILE_FORMAT =")).
		OptionalQueryStructField("CopyOptions", stageCopyOptionsDef(), g.ListOptions().Parentheses().NoComma().SQL("COPY_OPTIONS =")).
		OptionalComment().
		OptionalTags().
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists")
}

func alterStageOperation(structName string, apply func(qs *g.QueryStruct) *g.QueryStruct) *g.QueryStruct {
	qs := g.NewQueryStruct(structName).
		Alter().
		SQL("STAGE").
		IfExists().
		Name().
		SQL("SET")
	qs = apply(qs)
	return qs.
		OptionalQueryStructField("FileFormat", stageFileFormatDef, g.ListOptions().Parentheses().SQL("FILE_FORMAT =")).
		OptionalQueryStructField("CopyOptions", stageCopyOptionsDef(), g.ListOptions().Parentheses().NoComma().SQL("COPY_OPTIONS =")).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name")
}

var stageFileFormatDef = g.NewQueryStruct("StageFileFormat").
	OptionalTextAssignment("FORMAT_NAME", g.ParameterOptions().SingleQuotes()).
	OptionalAssignmentWithFieldName("TYPE", g.KindOfTPointer[sdkcommons.FileFormatType](), g.ParameterOptions(), "FileFormatType").
	PredefinedQueryStructField("Options", g.KindOfTPointer[sdkcommons.FileFormatTypeOptions](), g.ListOptions().NoComma())

var stageS3CommonDirectoryTableOptionsDef = func() *g.QueryStruct {
	return g.NewQueryStruct("StageS3CommonDirectoryTableOptions").
		BooleanAssignment("ENABLE", nil).
		OptionalBooleanAssignment("REFRESH_ON_CREATE", nil).
		OptionalBooleanAssignment("AUTO_REFRESH", nil)
}

var stageCopyOptionsDef = func() *g.QueryStruct {
	return g.NewQueryStruct("StageCopyOptions").
		OptionalQueryStructField(
			"OnError",
			g.NewQueryStruct("StageCopyOnErrorOptions").
				OptionalSQLWithCustomFieldName("Continue_", "CONTINUE").
				OptionalSQL("SKIP_FILE").
				// OptionalSQL("SKIP_FILE_n"). // TODO templated value - not even supported by structToSQL (could be keyword without space in-between)
				// OptionalSQL("SKIP_FILE_n%"). // TODO templated value with % - not even supported by structToSQL (could be keyword without space in-between)
				OptionalSQL("ABORT_STATEMENT"),
			g.ParameterOptions().SQL("ON_ERROR"),
		).
		OptionalNumberAssignment("SIZE_LIMIT", nil).
		OptionalBooleanAssignment("PURGE", nil).
		OptionalBooleanAssignment("RETURN_FAILED_ONLY", nil).
		OptionalAssignment("MATCH_BY_COLUMN_NAME", g.KindOfTPointer[sdkcommons.StageCopyColumnMapOption](), nil).
		OptionalBooleanAssignment("ENFORCE_LENGTH", nil).
		OptionalBooleanAssignment("TRUNCATECOLUMNS", nil).
		OptionalBooleanAssignment("FORCE", nil)
}

var externalS3StageParamsDef = func() *g.QueryStruct {
	return g.NewQueryStruct("ExternalS3StageParams").
		TextAssignment("URL", g.ParameterOptions().Required().SingleQuotes()).
		OptionalTextAssignment("AWS_ACCESS_POINT_ARN", g.ParameterOptions().SingleQuotes()).
		OptionalIdentifier("StorageIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("STORAGE_INTEGRATION")).
		OptionalQueryStructField(
			"Credentials",
			g.NewQueryStruct("ExternalStageS3Credentials").
				OptionalTextAssignment("AWS_KEY_ID", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("AWS_SECRET_KEY", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("AWS_TOKEN", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("AWS_ROLE", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.ConflictingFields, "AwsKeyId", "AwsRole").
				WithValidation(g.ConflictingFields, "AwsSecretKey", "AwsRole").
				WithValidation(g.ConflictingFields, "AwsToken", "AwsRole"),
			g.ListOptions().Parentheses().NoComma().SQL("CREDENTIALS ="),
		).
		OptionalQueryStructField("Encryption", g.NewQueryStruct("ExternalStageS3Encryption").
			OptionalAssignmentWithFieldName(
				"TYPE",
				g.KindOfT[sdkcommons.ExternalStageS3EncryptionOption](),
				g.ParameterOptions().SingleQuotes().Required(),
				"EncryptionType",
			).
			OptionalTextAssignment("MASTER_KEY", g.ParameterOptions().SingleQuotes()).
			OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
			g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
		).
		OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()).
		WithValidation(g.ConflictingFields, "StorageIntegration", "Credentials").
		WithValidation(g.ConflictingFields, "StorageIntegration", "UsePrivatelinkEndpoint")
}

var externalGCSStageParamsDef = func() *g.QueryStruct {
	return g.NewQueryStruct("ExternalGCSStageParams").
		TextAssignment("URL", g.ParameterOptions().SingleQuotes()).
		Identifier("StorageIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("STORAGE_INTEGRATION")).
		OptionalQueryStructField(
			"Encryption",
			g.NewQueryStruct("ExternalStageGCSEncryption").
				OptionalAssignmentWithFieldName(
					"TYPE",
					g.KindOfT[sdkcommons.ExternalStageGCSEncryptionOption](),
					g.ParameterOptions().SingleQuotes().Required(),
					"EncryptionType",
				).
				OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
			g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
		)
}

var externalAzureStageParamsDef = func() *g.QueryStruct {
	return g.NewQueryStruct("ExternalAzureStageParams").
		TextAssignment("URL", g.ParameterOptions().SingleQuotes()).
		OptionalIdentifier("StorageIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("STORAGE_INTEGRATION")).
		OptionalQueryStructField(
			"Credentials",
			g.NewQueryStruct("ExternalStageAzureCredentials").
				TextAssignment("AZURE_SAS_TOKEN", g.ParameterOptions().SingleQuotes()),
			g.ListOptions().Parentheses().NoComma().SQL("CREDENTIALS ="),
		).
		OptionalQueryStructField(
			"Encryption",
			g.NewQueryStruct("ExternalStageAzureEncryption").
				OptionalAssignmentWithFieldName(
					"TYPE",
					g.KindOfT[sdkcommons.ExternalStageAzureEncryptionOption](),
					g.ParameterOptions().SingleQuotes().Required(),
					"EncryptionType",
				).
				OptionalTextAssignment("MASTER_KEY", g.ParameterOptions().SingleQuotes()),
			g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
		).
		OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()).
		WithValidation(g.ConflictingFields, "StorageIntegration", "Credentials").
		WithValidation(g.ConflictingFields, "StorageIntegration", "UsePrivatelinkEndpoint")
}

var externalS3CompatibleStageParamsDef = func() *g.QueryStruct {
	return g.NewQueryStruct("ExternalS3CompatibleStageParams").
		TextAssignment("URL", g.ParameterOptions().Required().SingleQuotes()).
		TextAssignment("ENDPOINT", g.ParameterOptions().Required().SingleQuotes()).
		OptionalQueryStructField(
			"Credentials",
			g.NewQueryStruct("ExternalStageS3CompatibleCredentials").
				TextAssignment("AWS_KEY_ID", g.ParameterOptions().Required().SingleQuotes()).
				TextAssignment("AWS_SECRET_KEY", g.ParameterOptions().Required().SingleQuotes()),
			g.ListOptions().Parentheses().NoComma().SQL("CREDENTIALS ="),
		)
}

var stagesDef = g.NewInterface(
	"Stages",
	"Stage",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CustomOperation(
		"CreateInternal",
		"https://docs.snowflake.com/en/sql-reference/sql/create-stage",
		createStageOperation("CreateInternalStage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				OptionalQueryStructField(
					"Encryption",
					g.NewQueryStruct("InternalStageEncryption").
						AssignmentWithFieldName(
							"TYPE",
							g.KindOfT[sdkcommons.InternalStageEncryptionOption](),
							g.ParameterOptions().SingleQuotes().Required(),
							"EncryptionType",
						),
					g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
				).
				OptionalQueryStructField(
					"DirectoryTableOptions",
					g.NewQueryStruct("InternalDirectoryTableOptions").
						BooleanAssignment("ENABLE", nil).
						OptionalBooleanAssignment("AUTO_REFRESH", nil),
					g.ListOptions().Parentheses().NoComma().SQL("DIRECTORY ="),
				)
		}),
	).
	CustomOperation(
		"CreateOnS3",
		"https://docs.snowflake.com/en/sql-reference/sql/create-stage",
		createStageOperation("CreateExternalS3Stage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				QueryStructField("ExternalStageParams", externalS3StageParamsDef(), g.KeywordOptions().Required()).
				OptionalQueryStructField(
					"DirectoryTableOptions",
					stageS3CommonDirectoryTableOptionsDef(),
					g.ListOptions().Parentheses().NoComma().SQL("DIRECTORY ="),
				)
		}),
	).
	CustomOperation(
		"CreateOnGCS",
		"https://docs.snowflake.com/en/sql-reference/sql/create-stage",
		createStageOperation("CreateExternalGCSStage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				QueryStructField("ExternalStageParams", externalGCSStageParamsDef(), g.KeywordOptions().Required()).
				OptionalQueryStructField(
					"DirectoryTableOptions",
					g.NewQueryStruct("ExternalGCSDirectoryTableOptions").
						BooleanAssignment("ENABLE", nil).
						OptionalBooleanAssignment("REFRESH_ON_CREATE", nil).
						OptionalBooleanAssignment("AUTO_REFRESH", nil).
						OptionalTextAssignment("NOTIFICATION_INTEGRATION", g.ParameterOptions().SingleQuotes()),
					g.ListOptions().Parentheses().NoComma().SQL("DIRECTORY ="),
				)
		}),
	).
	CustomOperation(
		"CreateOnAzure",
		"https://docs.snowflake.com/en/sql-reference/sql/create-stage",
		createStageOperation("CreateExternalAzureStage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				QueryStructField("ExternalStageParams", externalAzureStageParamsDef(), g.KeywordOptions().Required()).
				OptionalQueryStructField(
					"DirectoryTableOptions",
					g.NewQueryStruct("ExternalAzureDirectoryTableOptions").
						BooleanAssignment("ENABLE", nil).
						OptionalBooleanAssignment("REFRESH_ON_CREATE", nil).
						OptionalBooleanAssignment("AUTO_REFRESH", nil).
						OptionalTextAssignment("NOTIFICATION_INTEGRATION", g.ParameterOptions().SingleQuotes()),
					g.ListOptions().Parentheses().NoComma().SQL("DIRECTORY ="),
				)
		}),
	).
	CustomOperation(
		"CreateOnS3Compatible",
		"https://docs.snowflake.com/en/sql-reference/sql/create-stage",
		createStageOperation("CreateExternalS3CompatibleStage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				QueryStructField("ExternalStageParams", externalS3CompatibleStageParamsDef(), g.KeywordOptions().Required()).
				OptionalQueryStructField(
					"DirectoryTableOptions",
					stageS3CommonDirectoryTableOptionsDef(),
					g.ListOptions().Parentheses().NoComma().SQL("DIRECTORY ="),
				)
		}),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-stage",
		g.NewQueryStruct("AlterStage").
			Alter().
			SQL("STAGE").
			IfExists().
			Name().
			OptionalIdentifier("RenameTo", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			OptionalSetTags().
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifierIfSet, "RenameTo").
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "SetTags", "UnsetTags").
			WithValidation(g.ValidIdentifier, "name"),
	).
	CustomOperation(
		"AlterInternalStage",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-stage",
		alterStageOperation("AlterInternalStage", func(qs *g.QueryStruct) *g.QueryStruct { return qs }),
	).
	CustomOperation(
		"AlterExternalS3Stage",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-stage",
		alterStageOperation("AlterExternalS3Stage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField("ExternalStageParams", externalS3StageParamsDef(), nil)
		}),
	).
	CustomOperation(
		"AlterExternalGCSStage",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-stage",
		alterStageOperation("AlterExternalGCSStage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField("ExternalStageParams", externalGCSStageParamsDef(), nil)
		}),
	).
	CustomOperation(
		"AlterExternalAzureStage",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-stage",
		alterStageOperation("AlterExternalAzureStage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField("ExternalStageParams", externalAzureStageParamsDef(), nil)
		}),
	).
	CustomOperation(
		"AlterDirectoryTable",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-stage",
		g.NewQueryStruct("AlterDirectoryTable").
			Alter().
			SQL("STAGE").
			IfExists().
			Name().
			OptionalQueryStructField(
				"SetDirectory",
				g.NewQueryStruct("DirectoryTableSet").BooleanAssignment("ENABLE", g.ParameterOptions().Required()),
				g.ListOptions().Parentheses().NoComma().SQL("SET DIRECTORY ="),
			).
			OptionalQueryStructField(
				"Refresh",
				g.NewQueryStruct("DirectoryTableRefresh").OptionalTextAssignment("SUBPATH", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions().SQL("REFRESH"),
			).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "SetDirectory", "Refresh"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-stage",
		g.NewQueryStruct("DropStage").
			Drop().
			SQL("STAGE").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-stage",
		g.DbStruct("stageDescRow").
			Field("parent_property", "string").
			Field("property", "string").
			Field("property_type", "string").
			Field("property_value", "string").
			Field("property_default", "string"),
		g.PlainStruct("StageProperty").
			Field("Parent", "string").
			Field("Name", "string").
			Field("Type", "string").
			Field("Value", "string").
			Field("Default", "string"),
		g.NewQueryStruct("DescStage").
			Describe().
			SQL("STAGE").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-stages",
		g.DbStruct("stageShowRow").
			Field("created_on", "time.Time").
			Field("name", "string").
			Field("database_name", "string").
			Field("schema_name", "string").
			Field("url", "string").
			Field("has_credentials", "string").
			Field("has_encryption_key", "string").
			Field("owner", "string").
			Field("comment", "string").
			Field("region", "sql.NullString").
			Field("type", "string").
			Field("cloud", "sql.NullString").
			// notification_channel is deprecated in Snowflake.
			Field("storage_integration", "sql.NullString").
			Field("endpoint", "sql.NullString").
			Field("owner_role_type", "sql.NullString").
			Field("directory_enabled", "string"),
		g.PlainStruct("Stage").
			Field("CreatedOn", "time.Time").
			Field("Name", "string").
			Field("DatabaseName", "string").
			Field("SchemaName", "string").
			Field("Url", "string").
			Field("HasCredentials", "bool").
			Field("HasEncryptionKey", "bool").
			Field("Owner", "string").
			Field("Comment", "string").
			Field("Region", "*string").
			Field("Type", "string").
			Field("Cloud", "*string").
			// notification_channel is deprecated in Snowflake.
			Field("StorageIntegration", "*string").
			Field("Endpoint", "*string").
			Field("OwnerRoleType", "*string").
			Field("DirectoryEnabled", "bool"),
		g.NewQueryStruct("ShowStages").
			Show().
			SQL("STAGES").
			OptionalLike().
			OptionalExtendedIn(),
	).
	ShowByIdOperationWithFiltering(
		g.ShowByIDLikeFiltering,
		g.ShowByIDExtendedInFiltering,
	)

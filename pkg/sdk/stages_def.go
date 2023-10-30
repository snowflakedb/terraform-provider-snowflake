package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

type InternalStageEncryptionOption string

const (
	InternalStageEncryptionFull InternalStageEncryptionOption = "SNOWFLAKE_FULL"
	InternalStageEncryptionSSE  InternalStageEncryptionOption = "SNOWFLAKE_SSE"
)

type ExternalStageS3EncryptionOption string

const (
	ExternalStageS3EncryptionCSE    ExternalStageS3EncryptionOption = "AWS_CSE"
	ExternalStageS3EncryptionSSES3  ExternalStageS3EncryptionOption = "AWS_SSE_S3"
	ExternalStageS3EncryptionSSEKMS ExternalStageS3EncryptionOption = "AWS_SSE_KMS"
	ExternalStageS3EncryptionNone   ExternalStageS3EncryptionOption = "NONE"
)

type ExternalStageGCSEncryptionOption string

const (
	ExternalStageGCSEncryptionSSEKMS ExternalStageGCSEncryptionOption = "GCS_SSE_KMS"
	ExternalStageGCSEncryptionNone   ExternalStageGCSEncryptionOption = "NONE"
)

type ExternalStageAzureEncryptionOption string

const (
	ExternalStageAzureEncryptionCSE  ExternalStageAzureEncryptionOption = "AZURE_CSE"
	ExternalStageAzureEncryptionNone ExternalStageAzureEncryptionOption = "NONE"
)

type StageCopyColumnMapOption string

const (
	StageCopyColumnMapCaseSensitive   StageCopyColumnMapOption = "CASE_SENSITIVE"
	StageCopyColumnMapCaseInsensitive StageCopyColumnMapOption = "CASE_INSENSITIVE"
	StageCopyColumnMapCaseNone        StageCopyColumnMapOption = "NONE"
)

// TODO PUT, GET, LS, etc. ???

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
		OptionalQueryStructField("FileFormat", stageFileFormatDef, g.ParameterOptions().Parentheses().SQL("FILE_FORMAT")).
		OptionalQueryStructField("CopyOptions", stageCopyOptionsDef, g.ParameterOptions().SQL("COPY_OPTIONS").Parentheses()).
		OptionalComment().
		OptionalTags()
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
		OptionalQueryStructField("FileFormat", stageFileFormatDef, g.ParameterOptions().SQL("FILE_FORMAT").Parentheses()).
		OptionalQueryStructField("CopyOptions", stageCopyOptionsDef, g.ParameterOptions().SQL("COPY_OPTIONS").Parentheses()).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name")
}

var stageFileFormatDef = g.NewQueryStruct("StageFileFormat").
	OptionalTextAssignment("FORMAT_NAME", g.ParameterOptions().SingleQuotes()).
	OptionalAssignment("TYPE", g.KindOfTPointer[FileFormatType](), g.ParameterOptions()).
	List("TYPE", g.KindOfT[FileFormatType](), g.ListOptions())

var stageCopyOptionsDef = g.NewQueryStruct("StageCopyOptions").
	OptionalQueryStructField(
		"OnError",
		g.NewQueryStruct("StageCopyOnErrorOptions").
			OptionalSQL("CONTINUE").
			OptionalSQL("SKIP_FILE").
			OptionalSQL("SKIP_FILE_").    // TODO templated value - not even supported by structToSQL (I think)
			OptionalSQL("SKIP_FILE_123"). // TODO templated value with % - not even supported by structToSQL (I think)
			OptionalSQL("ABORT_STATEMENT"),
		g.ParameterOptions().SQL("ON_ERROR"),
	).
	OptionalNumberAssignment("SIZE_LIMIT", nil).
	OptionalBooleanAssignment("PURGE", nil).
	OptionalBooleanAssignment("RETURN_FAILED_ONLY", nil).
	OptionalAssignment("MATCH_BY_COLUMN_NAME", g.KindOfTPointer[StageCopyColumnMapOption](), nil).
	OptionalBooleanAssignment("ENFORCE_LENGTH", nil).
	OptionalBooleanAssignment("TRUNCATECOLUMNS", nil).
	OptionalBooleanAssignment("FORCE", nil)

var externalS3StageParamsDef = g.NewQueryStruct("ExternalS3StageParams").
	TextAssignment("URL", g.ParameterOptions().SingleQuotes()).
	OptionalIdentifierAssignment("STORAGE_INTEGRATION", g.KindOfT[AccountObjectIdentifier](), g.ParameterOptions()).
	OptionalQueryStructField(
		"Credentials",
		g.NewQueryStruct("ExternalStageS3Credentials").
			OptionalTextAssignment("AWS_KEY_ID", g.ParameterOptions().SingleQuotes()).
			OptionalTextAssignment("AWS_SECRET_KEY", g.ParameterOptions().SingleQuotes()).
			OptionalTextAssignment("AWS_TOKEN", g.ParameterOptions().SingleQuotes()).
			OptionalTextAssignment("AWS_ROLE", g.ParameterOptions().SingleQuotes()),
		g.ParameterOptions().Parentheses().SQL("CREDENTIALS"),
	).
	OptionalQueryStructField("Encryption", g.NewQueryStruct("ExternalStageS3Encryption").
		OptionalAssignment(
			"TYPE",
			g.KindOfT[ExternalStageS3EncryptionOption](),
			g.ParameterOptions().SingleQuotes().Required(),
		).
		OptionalTextAssignment("MASTER_KEY", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
		g.ParameterOptions().Parentheses().SQL("ENCRYPTION"),
	)

var externalGCSStageParamsDef = g.NewQueryStruct("ExternalGCSStageParams").
	TextAssignment("URL", g.ParameterOptions().SingleQuotes()).
	OptionalIdentifierAssignment("STORAGE_INTEGRATION", g.KindOfT[AccountObjectIdentifier](), g.ParameterOptions()).
	OptionalQueryStructField(
		"Encryption",
		g.NewQueryStruct("ExternalStageGCSEncryption").
			OptionalAssignment(
				"TYPE",
				g.KindOfT[ExternalStageGCSEncryptionOption](),
				g.ParameterOptions().SingleQuotes().Required(),
			).
			OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
		g.ParameterOptions().Parentheses().SQL("ENCRYPTION"),
	)

var externalAzureStageParamsDef = g.NewQueryStruct("ExternalAzureStageParams").
	TextAssignment("URL", g.ParameterOptions().SingleQuotes()).
	OptionalIdentifierAssignment("STORAGE_INTEGRATION", g.KindOfT[AccountObjectIdentifier](), g.ParameterOptions()).
	OptionalQueryStructField(
		"Credentials",
		g.NewQueryStruct("ExternalStageAzureCredentials").
			OptionalTextAssignment("AZURE_SAS_TOKEN", g.ParameterOptions().SingleQuotes()),
		g.ParameterOptions().Parentheses(),
	).
	OptionalQueryStructField(
		"Encryption",
		g.NewQueryStruct("ExternalStageGCSEncryption").
			OptionalAssignment(
				"TYPE",
				g.KindOfT[ExternalStageAzureEncryptionOption](),
				g.ParameterOptions().SingleQuotes().Required(),
			).
			OptionalTextAssignment("MASTER_KEY", g.ParameterOptions().SingleQuotes()),
		g.ParameterOptions().Parentheses().SQL("ENCRYPTION"),
	)

var StagesDef = g.NewInterface(
	"Stages",
	"Stage",
	g.KindOfT[SchemaObjectIdentifier](),
).
	CustomOperation(
		"CreateInternal",
		"https://docs.snowflake.com/en/sql-reference/sql/create-stage",
		createStageOperation("CreateInternalStage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				OptionalQueryStructField(
					"Encryption",
					g.NewQueryStruct("InternalStageEncryption").
						OptionalAssignment(
							"TYPE",
							g.KindOfT[InternalStageEncryptionOption](),
							g.ParameterOptions().SingleQuotes(),
						),
					g.ParameterOptions().Parentheses().SQL("ENCRYPTION"),
				).
				OptionalQueryStructField(
					"DirectoryTableOptions",
					g.NewQueryStruct("InternalDirectoryTableOptions").
						OptionalBooleanAssignment("ENABLE", nil).
						OptionalBooleanAssignment("REFRESH_ON_CREATE", nil),
					g.ParameterOptions().Parentheses(),
				)
		}),
	).
	CustomOperation(
		"CreateOnS3",
		"https://docs.snowflake.com/en/sql-reference/sql/create-stage",
		createStageOperation("CreateExternalS3Stage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				OptionalQueryStructField("ExternalStageParams", externalS3StageParamsDef, nil).
				OptionalQueryStructField(
					"DirectoryTableOptions",
					g.NewQueryStruct("ExternalS3DirectoryTableOptions").
						OptionalBooleanAssignment("ENABLE", nil).
						OptionalBooleanAssignment("REFRESH_ON_CREATE", nil).
						OptionalBooleanAssignment("AUTO_REFRESH", nil),
					g.ParameterOptions().Parentheses(),
				)
		}),
	).
	CustomOperation(
		"CreateOnGCS",
		"https://docs.snowflake.com/en/sql-reference/sql/create-stage",
		createStageOperation("CreateExternalGCSStage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				OptionalQueryStructField("ExternalStageParams", externalGCSStageParamsDef, nil).
				OptionalQueryStructField(
					"DirectoryTableOptions",
					g.NewQueryStruct("ExternalGCSDirectoryTableOptions").
						OptionalBooleanAssignment("ENABLE", nil).
						OptionalBooleanAssignment("REFRESH_ON_CREATE", nil).
						OptionalBooleanAssignment("AUTO_REFRESH", nil).
						OptionalTextAssignment("NOTIFICATION_INTEGRATION", g.ParameterOptions().SingleQuotes()),
					g.ParameterOptions().Parentheses(),
				)
		}),
	).
	CustomOperation(
		"CreateOnAzure",
		"https://docs.snowflake.com/en/sql-reference/sql/create-stage",
		createStageOperation("CreateExternalAzureStage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				OptionalQueryStructField("ExternalStageParams", externalAzureStageParamsDef, nil).
				OptionalQueryStructField(
					"DirectoryTableOptions",
					g.NewQueryStruct("ExternalAzureDirectoryTableOptions").
						OptionalBooleanAssignment("ENABLE", nil).
						OptionalBooleanAssignment("REFRESH_ON_CREATE", nil).
						OptionalBooleanAssignment("AUTO_REFRESH", nil).
						OptionalTextAssignment("NOTIFICATION_INTEGRATION", g.ParameterOptions().SingleQuotes()),
					g.ParameterOptions().Parentheses(),
				)
		}),
	).
	CustomOperation(
		"CreateOnS3Compatible",
		"https://docs.snowflake.com/en/sql-reference/sql/create-stage",
		createStageOperation("CreateExternalS3CompatibleStage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				TextAssignment("URL", g.ParameterOptions().SingleQuotes().Required()).
				TextAssignment("ENDPOINT", g.ParameterOptions().SingleQuotes().Required()).
				OptionalQueryStructField(
					"Credentials",
					g.NewQueryStruct("ExternalStageS3CompatibleCredentials").
						OptionalTextAssignment("AWS_KEY_ID", g.ParameterOptions().SingleQuotes().Required()).
						OptionalTextAssignment("AWS_SECRET_KEY", g.ParameterOptions().SingleQuotes().Required()),
					g.ParameterOptions().Parentheses(),
				).
				// TODO: Can be used with compat ?
				OptionalQueryStructField(
					"DirectoryTableOptions",
					g.NewQueryStruct("ExternalS3DirectoryTableOptions").
						OptionalBooleanAssignment("ENABLE", nil).
						OptionalBooleanAssignment("REFRESH_ON_CREATE", nil).
						OptionalBooleanAssignment("AUTO_REFRESH", nil),
					g.ParameterOptions().Parentheses(),
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
			OptionalIdentifier("RenameTo", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().NoEquals()).
			List("SetTags", g.KindOfT[TagAssociation](), g.KeywordOptions().SQL("SET TAG")).
			List("UnsetTags", g.KindOfT[ObjectIdentifier](), g.KeywordOptions().SQL("UNSET TAG")).
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
			return qs.OptionalQueryStructField("ExternalStageParams", externalS3StageParamsDef, nil)
		}),
	).
	CustomOperation(
		"AlterExternalGCSStage",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-stage",
		alterStageOperation("AlterExternalGCSStage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField("ExternalStageParams", externalGCSStageParamsDef, nil)
		}),
	).
	CustomOperation(
		"AlterExternalAzureStage",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-stage",
		alterStageOperation("AlterExternalAzureStage", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField("ExternalStageParams", externalAzureStageParamsDef, nil)
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
				g.ParameterOptions().SQL("SET DIRECTORY").Parentheses(),
			).
			OptionalQueryStructField(
				"Refresh",
				g.NewQueryStruct("DirectoryTableRefresh").OptionalTextAssignment("SUBPATH", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions().SQL("REFRESH"),
			).
			WithValidation(g.ValidIdentifier, "name"),
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
			Field("property_value", "sql.NullString").
			Field("property_default", "sql.NullString"),
		g.PlainStruct("StageProperty").
			Field("Parent", "string").
			Field("Name", "string").
			Field("Type", "string").
			Field("Value", "*string").
			Field("Default", "*string"),
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
			Field("region", "string").
			Field("type", "string").
			Field("cloud", "sql.NullString").
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
			Field("Region", "string").
			Field("Type", "string").
			Field("Cloud", "*string").
			Field("StorageIntegration", "*string").
			Field("Endpoint", "*string").
			Field("OwnerRoleType", "*string").
			Field("DirectoryEnabled", "bool"),
		g.NewQueryStruct("ShowStages").
			Show().
			SQL("STAGES").
			OptionalLike().
			OptionalIn(),
	)

package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
	InternalStageEncryptionOptionEnumDef = g.NewEnum(
		"InternalStageEncryptionOption", "InternalStageEncryptionOptions",
		"SNOWFLAKE_FULL", "SNOWFLAKE_SSE",
	)
	ExternalStageS3EncryptionOptionEnumDef = g.NewEnum(
		"ExternalStageS3EncryptionOption", "ExternalStageS3EncryptionOptions",
		"AWS_CSE", "AWS_SSE_S3", "AWS_SSE_KMS", "NONE",
	)
	ExternalStageGCSEncryptionOptionEnumDef = g.NewEnum(
		"ExternalStageGCSEncryptionOption", "ExternalStageGCSEncryptionOptions",
		"GCS_SSE_KMS", "NONE",
	)
	ExternalStageAzureEncryptionOptionEnumDef = g.NewEnum(
		"ExternalStageAzureEncryptionOption", "ExternalStageAzureEncryptionOptions",
		"AZURE_CSE", "NONE",
	)
	StageCopyColumnMapOptionEnumDef = g.NewEnum(
		"StageCopyColumnMapOption", "StageCopyColumnMapOptions",
		"CASE_SENSITIVE", "CASE_INSENSITIVE", "NONE",
	)
	StageCloudEnumDef = g.NewEnum(
		"StageCloud", "StageClouds",
		"AZURE", "AWS", "GCP",
	)
	StageTypeEnumDef = g.NewEnum(
		"StageType", "StageTypes",
		"INTERNAL", "INTERNAL NO CSE", "INTERNAL TEMPORARY", "EXTERNAL", "EXTERNAL TEMPORARY",
	)
)

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
		OptionalQueryStructField(
			"FileFormat",
			g.NewQueryStruct("StageFileFormat").
				OptionalIdentifier("FormatName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("FORMAT_NAME =")).
				PredefinedQueryStructField("FileFormatOptions", "*FileFormatOptions", g.ListOptions()).
				WithValidation(g.ExactlyOneValueSet, "FormatName", "FileFormatOptions").
				WithAdditionalValidations(),
			g.ListOptions().Parentheses().NoComma().SQL("FILE_FORMAT ="),
		).
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
		OptionalQueryStructField(
			"FileFormat",
			g.NewQueryStruct("StageFileFormat").
				OptionalIdentifier("FormatName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("FORMAT_NAME =")).
				PredefinedQueryStructField("FileFormatOptions", "*FileFormatOptions", g.ListOptions()).
				WithValidation(g.ExactlyOneValueSet, "FormatName", "FileFormatOptions").
				WithAdditionalValidations(),
			g.ListOptions().Parentheses().NoComma().SQL("FILE_FORMAT ="),
		).
		// TODO(SNOW-3035788): use COMMENT in unset and here use OptionalComment
		OptionalAssignment("COMMENT", "StringAllowEmpty", g.ParameterOptions()).
		WithValidation(g.ValidIdentifier, "name")
}

var stageS3CommonDirectoryTableOptionsDef = func() *g.QueryStruct {
	return g.NewQueryStruct("StageS3CommonDirectoryTableOptions").
		BooleanAssignment("ENABLE", nil).
		OptionalBooleanAssignment("REFRESH_ON_CREATE", nil).
		OptionalBooleanAssignment("AUTO_REFRESH", nil)
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
		OptionalQueryStructField(
			"Encryption", g.NewQueryStruct("ExternalStageS3Encryption").
				OptionalQueryStructField(
					"AwsCse",
					g.NewQueryStruct("ExternalStageS3EncryptionAwsCse").
						SQLWithCustomFieldName("encryptionType", "TYPE = 'AWS_CSE'").
						TextAssignment("MASTER_KEY", g.ParameterOptions().Required().SingleQuotes()),
					g.KeywordOptions(),
				).
				OptionalQueryStructField(
					"AwsSseS3",
					g.NewQueryStruct("ExternalStageS3EncryptionAwsSseS3").
						SQLWithCustomFieldName("encryptionType", "TYPE = 'AWS_SSE_S3'"),
					g.KeywordOptions(),
				).
				OptionalQueryStructField(
					"AwsSseKms",
					g.NewQueryStruct("ExternalStageS3EncryptionAwsSseKms").
						SQLWithCustomFieldName("encryptionType", "TYPE = 'AWS_SSE_KMS'").
						OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
					g.KeywordOptions(),
				).
				OptionalQueryStructField(
					"None",
					g.NewQueryStruct("ExternalStageS3EncryptionNone").
						SQLWithCustomFieldName("encryptionType", "TYPE = 'NONE'"),
					g.KeywordOptions(),
				).
				WithValidation(g.ExactlyOneValueSet, "AwsCse", "AwsSseS3", "AwsSseKms", "None"),
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
				OptionalQueryStructField(
					"GcsSseKms",
					g.NewQueryStruct("ExternalStageGCSEncryptionGcsSseKms").
						SQLWithCustomFieldName("encryptionType", "TYPE = 'GCS_SSE_KMS'").
						OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
					g.KeywordOptions(),
				).
				OptionalQueryStructField(
					"None",
					g.NewQueryStruct("ExternalStageGCSEncryptionNone").
						SQLWithCustomFieldName("encryptionType", "TYPE = 'NONE'"),
					g.KeywordOptions(),
				).
				WithValidation(g.ExactlyOneValueSet, "GcsSseKms", "None"),
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
				OptionalQueryStructField(
					"AzureCse",
					g.NewQueryStruct("ExternalStageAzureEncryptionAzureCse").
						SQLWithCustomFieldName("encryptionType", "TYPE = 'AZURE_CSE'").
						TextAssignment("MASTER_KEY", g.ParameterOptions().Required().SingleQuotes()),
					g.KeywordOptions(),
				).
				OptionalQueryStructField(
					"None",
					g.NewQueryStruct("ExternalStageAzureEncryptionNone").
						SQLWithCustomFieldName("encryptionType", "TYPE = 'NONE'"),
					g.KeywordOptions(),
				).
				WithValidation(g.ExactlyOneValueSet, "AzureCse", "None"),
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
						OptionalQueryStructField(
							"SnowflakeFull",
							g.NewQueryStruct("InternalStageEncryptionSnowflakeFull").
								SQLWithCustomFieldName("encryptionType", "TYPE = 'SNOWFLAKE_FULL'"),
							g.KeywordOptions(),
						).
						OptionalQueryStructField(
							"SnowflakeSse",
							g.NewQueryStruct("InternalStageEncryptionSnowflakeSse").
								SQLWithCustomFieldName("encryptionType", "TYPE = 'SNOWFLAKE_SSE'"),
							g.KeywordOptions(),
						).
						WithValidation(g.ExactlyOneValueSet, "SnowflakeFull", "SnowflakeSse"),
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
			RenameTo().
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
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-stage",
		g.StructPair("stageDescRow", "StageProperty").
			Text("parent_property", g.WithPlainFieldName("Parent")).
			Text("property", g.WithPlainFieldName("Name")).
			Text("property_type", g.WithPlainFieldName("Type")).
			Text("property_value", g.WithPlainFieldName("Value")).
			Text("property_default", g.WithPlainFieldName("Default")),
		g.NewQueryStruct("DescStage").
			Describe().
			SQL("STAGE").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
		g.PlainStruct("StageDirectoryTable").
			Bool("Enable").
			Bool("AutoRefresh").
			OptionalText("DirectoryNotificationChannel").
			OptionalText("LastRefreshedOn"),
		g.PlainStruct("StagePrivateLink").
			Bool("UsePrivatelinkEndpoint"),
		g.PlainStruct("StageLocationDetails").
			StringList("Url").
			Text("AwsAccessPointArn"),
		g.PlainStruct("StageCredentials").
			Text("AwsKeyId"),
		g.PlainStruct("StageDetails").
			SchemaObjectIdentifier().
			OptionalField("FileFormatName", "SchemaObjectIdentifier").
			OptionalField("FileFormatCsv", "FileFormatCsv").
			OptionalField("FileFormatJson", "FileFormatJson").
			OptionalField("FileFormatAvro", "FileFormatAvro").
			OptionalField("FileFormatOrc", "FileFormatOrc").
			OptionalField("FileFormatParquet", "FileFormatParquet").
			OptionalField("FileFormatXml", "FileFormatXml").
			OptionalField("DirectoryTable", "StageDirectoryTable").
			OptionalField("PrivateLink", "StagePrivateLink").
			OptionalField("Location", "StageLocationDetails").
			OptionalField("Credentials", "StageCredentials"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-stages",
		g.StructPair("stageShowRow", "Stage").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("url").
			Field("has_credentials", "string", "bool").
			Field("has_encryption_key", "string", "bool").
			Text("owner").
			Text("comment").
			OptionalText("region").
			Enum("type", StageTypeEnumDef).
			OptionalEnum("cloud", StageCloudEnumDef).
			// notification_channel is deprecated in Snowflake.
			Field("storage_integration", "sql.NullString", "*AccountObjectIdentifier", g.WithPlainFieldName("StorageIntegration")).
			OptionalText("endpoint").
			Field("owner_role_type", "sql.NullString", "*string").
			Field("directory_enabled", "string", "bool"),
		g.NewQueryStruct("ShowStages").
			Show().
			SQL("STAGES").
			OptionalLike().
			OptionalExtendedIn(),
		g.ShowByIDLikeFiltering,
		g.ShowByIDExtendedInFiltering,
	).
	WithEnums(
		InternalStageEncryptionOptionEnumDef,
		ExternalStageS3EncryptionOptionEnumDef,
		ExternalStageGCSEncryptionOptionEnumDef,
		ExternalStageAzureEncryptionOptionEnumDef,
		StageCopyColumnMapOptionEnumDef,
		StageCloudEnumDef,
		StageTypeEnumDef,
	).
	WithCustomInterfaceMethod(
		"DescribeDetails",
		"DescribeDetails returns parsed describe output for stages.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
		"*StageDetails", "error",
	)

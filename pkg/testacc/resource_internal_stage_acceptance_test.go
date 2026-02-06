//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_InternalStage_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment, changedComment := random.Comment(), random.Comment()

	azureUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)

	modelBasic := model.InternalStageWithId(id)

	modelComplete := model.InternalStageWithId(id).
		WithComment(comment).
		WithDirectoryEnabled(r.BooleanTrue).
		WithEncryptionSnowflakeFull()

	modelUpdated := model.InternalStageWithId(id).
		WithComment(changedComment).
		WithDirectoryEnabled(r.BooleanFalse).
		WithEncryptionSnowflakeFull()

	modelSseEncryptionWithDirectoryTableAndComment := model.InternalStageWithId(id).
		WithComment(changedComment).
		WithDirectoryEnabled(r.BooleanTrue).
		WithEncryptionSnowflakeSse()

	modelSseEncryptionWithComment := model.InternalStageWithId(id).
		WithComment(changedComment).
		WithEncryptionSnowflakeSse()

	modelSseEncryption := model.InternalStageWithId(newId).
		WithEncryptionSnowflakeSse()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			// Create with empty optionals (basic)
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeInternal),
					resourceshowoutputassert.StageShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeInternal).
						HasComment("").
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// Import - without optionals
			{
				Config:                  accconfig.FromModels(t, modelBasic),
				ResourceName:            modelBasic.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"directory", "file_format"},
			},
			// Set optionals (complete)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(comment).
						HasDirectoryEnableString(r.BooleanTrue).
						HasEncryptionSnowflakeFull().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeInternal),
					resourceshowoutputassert.StageShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeInternal).
						HasComment(comment).
						HasDirectoryEnabled(true).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.directory_table.0.enable", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// Import - complete
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption", "directory", "file_format"},
			},
			// Alter (update comment, directory.enable)
			{
				Config: accconfig.FromModels(t, modelUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelUpdated.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelUpdated.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(changedComment).
						HasDirectoryEnableString(r.BooleanFalse).
						HasEncryptionSnowflakeFull().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeInternal),
					resourceshowoutputassert.StageShowOutput(t, modelUpdated.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeInternal).
						HasComment(changedComment).
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// External change detection
			{
				PreConfig: func() {
					testClient().Stage.AlterInternalStage(t, sdk.NewAlterInternalStageStageRequest(id).
						WithComment(sdk.StringAllowEmpty{Value: changedComment}))
				},
				Config: accconfig.FromModels(t, modelUpdated),
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelUpdated.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(changedComment).
						HasDirectoryEnableString(r.BooleanFalse).
						HasEncryptionSnowflakeFull().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeInternal),
					resourceshowoutputassert.StageShowOutput(t, modelUpdated.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeInternal).
						HasComment(changedComment).
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// ForceNew - encryption change (requires recreate)
			{
				Config: accconfig.FromModels(t, modelSseEncryptionWithDirectoryTableAndComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelSseEncryptionWithDirectoryTableAndComment.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelSseEncryptionWithDirectoryTableAndComment.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(changedComment).
						HasDirectoryEnableString(r.BooleanTrue).
						HasEncryptionSnowflakeSse().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeInternalNoCse),
					resourceshowoutputassert.StageShowOutput(t, modelSseEncryptionWithDirectoryTableAndComment.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeInternalNoCse).
						HasComment(changedComment).
						HasDirectoryEnabled(true).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelSseEncryptionWithDirectoryTableAndComment.ResourceReference(), "describe_output.0.directory_table.0.enable", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelSseEncryptionWithDirectoryTableAndComment.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// ForceNew - unset directorytable
			{
				Config: accconfig.FromModels(t, modelSseEncryptionWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelSseEncryptionWithComment.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelSseEncryptionWithComment.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(changedComment).
						HasDirectoryEmpty().
						HasEncryptionSnowflakeSse().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeInternalNoCse),
					resourceshowoutputassert.StageShowOutput(t, modelSseEncryptionWithComment.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeInternalNoCse).
						HasComment(changedComment).
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelSseEncryptionWithComment.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelSseEncryptionWithComment.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// Unset comment and rename
			{
				Config: accconfig.FromModels(t, modelSseEncryption),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelSseEncryption.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelSseEncryption.ResourceReference()).
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasSchemaString(newId.SchemaName()).
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionSnowflakeSse().
						HasFullyQualifiedNameString(newId.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeInternalNoCse),
					resourceshowoutputassert.StageShowOutput(t, modelSseEncryption.ResourceReference()).
						HasName(newId.Name()).
						HasDatabaseName(newId.DatabaseName()).
						HasSchemaName(newId.SchemaName()).
						HasType(sdk.StageTypeInternalNoCse).
						HasCommentEmpty().
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelSseEncryption.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelSseEncryption.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// External type change detection
			{
				PreConfig: func() {
					testClient().Stage.DropStageFunc(t, newId)()
					testClient().Stage.CreateStageOnAzureWithId(t, newId, azureUrl)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelSseEncryption.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, modelSseEncryption),
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelSseEncryption.ResourceReference()).
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasSchemaString(newId.SchemaName()).
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionSnowflakeSse().
						HasFullyQualifiedNameString(newId.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeInternalNoCse),
					resourceshowoutputassert.StageShowOutput(t, modelSseEncryption.ResourceReference()).
						HasName(newId.Name()).
						HasDatabaseName(newId.DatabaseName()).
						HasSchemaName(newId.SchemaName()).
						HasType(sdk.StageTypeInternalNoCse).
						HasCommentEmpty().
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelSseEncryption.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelSseEncryption.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
		},
	})
}

func TestAcc_InternalStage_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	modelComplete := model.InternalStageWithId(id).
		WithComment(comment).
		WithDirectoryEnabledAndAutoRefresh(true, r.BooleanFalse).
		WithEncryptionSnowflakeFull()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(comment).
						HasDirectoryEnableString(r.BooleanTrue).
						HasDirectoryAutoRefreshString(r.BooleanFalse).
						HasEncryptionSnowflakeFull().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeInternal),
					resourceshowoutputassert.StageShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeInternal).
						HasComment(comment).
						HasDirectoryEnabled(true).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.directory_table.0.enable", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption", "file_format"},
			},
		},
	})
}

func TestAcc_InternalStage_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelInvalidAutoRefresh := model.InternalStageWithId(id).
		WithDirectoryEnabledAndAutoRefresh(true, "invalid")

	modelBothEncryptionTypes := model.InternalStageWithId(id).
		WithEncryptionBothTypes()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelInvalidAutoRefresh),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected .*auto_refresh.* to be one of \["true" "false"], got invalid`),
			},
			{
				Config:      accconfig.FromModels(t, modelBothEncryptionTypes),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`encryption.0.snowflake_full,encryption.0.snowflake_sse.* can be specified`),
			},
		},
	})
}

func TestAcc_InternalStage_FileFormat_SwitchBetweenTypes(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	fileFormat, fileFormatCleanup := testClient().FileFormat.CreateFileFormat(t)
	t.Cleanup(fileFormatCleanup)

	modelBasic := model.InternalStageWithId(id)

	modelWithCsvFormat := model.InternalStageWithId(id).
		WithFileFormatCsv(model.CsvFileFormatOptions{})

	modelWithNamedFormat := model.InternalStageWithId(id).
		WithFileFormatName(fileFormat.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			// Start with inline CSV
			{
				Config: accconfig.FromModels(t, modelWithCsvFormat),
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelWithCsvFormat.ResourceReference()).
						HasFileFormatCsv(),
					assert.Check(resource.TestCheckResourceAttr(modelWithCsvFormat.ResourceReference(), "describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithCsvFormat.ResourceReference(), "describe_output.0.file_format.0.csv.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithCsvFormat.ResourceReference(), "describe_output.0.file_format.0.format_name", "")),
				),
			},
			// Switch to named format
			{
				Config: accconfig.FromModels(t, modelWithNamedFormat),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithNamedFormat.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelWithNamedFormat.ResourceReference()).
						HasFileFormatFormatName(fileFormat.ID().FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(modelWithNamedFormat.ResourceReference(), "describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithNamedFormat.ResourceReference(), "describe_output.0.file_format.0.csv.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelWithNamedFormat.ResourceReference(), "describe_output.0.file_format.0.format_name", fileFormat.ID().FullyQualifiedName())),
				),
			},
			// Detect external change
			{
				Config: accconfig.FromModels(t, modelWithNamedFormat),
				PreConfig: func() {
					testClient().Stage.AlterInternalStage(t, sdk.NewAlterInternalStageStageRequest(id).WithFileFormat(sdk.StageFileFormatRequest{FileFormatOptions: &sdk.FileFormatOptions{CsvOptions: &sdk.FileFormatCsvOptions{}}}))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithNamedFormat.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelWithNamedFormat.ResourceReference()).
						HasFileFormatFormatName(fileFormat.ID().FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(modelWithNamedFormat.ResourceReference(), "describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithNamedFormat.ResourceReference(), "describe_output.0.file_format.0.csv.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelWithNamedFormat.ResourceReference(), "describe_output.0.file_format.0.format_name", fileFormat.ID().FullyQualifiedName())),
				),
			},
			// Switch back to inline CSV
			{
				Config: accconfig.FromModels(t, modelWithCsvFormat),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithCsvFormat.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelWithCsvFormat.ResourceReference()).
						HasFileFormatCsv(),
					assert.Check(resource.TestCheckResourceAttr(modelWithCsvFormat.ResourceReference(), "describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithCsvFormat.ResourceReference(), "describe_output.0.file_format.0.csv.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithCsvFormat.ResourceReference(), "describe_output.0.file_format.0.format_name", "")),
				),
			},
			// Switch back to default
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelBasic.ResourceReference()).
						HasFileFormatEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.file_format.0.csv.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.file_format.0.format_name", "")),
				),
			},
		},
	})
}

func TestAcc_InternalStage_FileFormat_AllCsvOptions(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	multiLine := true
	parseHeader := false
	skipBlankLines := true
	trimSpace := true
	errorOnColumnCountMismatch := false
	replaceInvalidCharacters := true
	emptyFieldAsNull := true
	skipByteOrderMark := true

	modelWithoutFileFormat := model.InternalStageWithId(id).
		WithFileFormatCsv(model.CsvFileFormatOptions{})

	modelCompleteCsv := model.InternalStageWithId(id).
		WithFileFormatCsv(model.CsvFileFormatOptions{
			Compression:                sdk.CSVCompressionGzip,
			FieldDelimiter:             "|",
			MultiLine:                  &multiLine,
			FileExtension:              "csv",
			ParseHeader:                &parseHeader,
			SkipBlankLines:             &skipBlankLines,
			DateFormat:                 "AUTO",
			TimeFormat:                 "AUTO",
			TimestampFormat:            "AUTO",
			BinaryFormat:               sdk.BinaryFormatHex,
			Escape:                     `\`,
			EscapeUnenclosedField:      "NONE",
			TrimSpace:                  &trimSpace,
			FieldOptionallyEnclosedBy:  `"`,
			NullIf:                     []string{"NULL", ""},
			ErrorOnColumnCountMismatch: &errorOnColumnCountMismatch,
			ReplaceInvalidCharacters:   &replaceInvalidCharacters,
			EmptyFieldAsNull:           &emptyFieldAsNull,
			SkipByteOrderMark:          &skipByteOrderMark,
			Encoding:                   sdk.CSVEncodingUTF8,
			RecordDelimiter:            ";",
		})

	altMultiLine := false
	altParseHeader := true
	altSkipBlankLines := false
	altTrimSpace := false
	altErrorOnColumnCountMismatch := true
	altReplaceInvalidCharacters := false
	altEmptyFieldAsNull := false
	altSkipByteOrderMark := false

	modelAlteredCsv := model.InternalStageWithId(id).
		WithFileFormatCsv(model.CsvFileFormatOptions{
			Compression:                sdk.CSVCompressionZstd,
			FieldDelimiter:             ",",
			MultiLine:                  &altMultiLine,
			FileExtension:              "txt",
			ParseHeader:                &altParseHeader,
			SkipBlankLines:             &altSkipBlankLines,
			DateFormat:                 "YYYY",
			TimeFormat:                 "HH24:MI:SS",
			TimestampFormat:            "YYYY-MM-DD HH24:MI:SS",
			BinaryFormat:               sdk.BinaryFormatBase64,
			Escape:                     `\\`,
			EscapeUnenclosedField:      "NONE",
			TrimSpace:                  &altTrimSpace,
			FieldOptionallyEnclosedBy:  `"`,
			NullIf:                     []string{"NA"},
			ErrorOnColumnCountMismatch: &altErrorOnColumnCountMismatch,
			ReplaceInvalidCharacters:   &altReplaceInvalidCharacters,
			EmptyFieldAsNull:           &altEmptyFieldAsNull,
			SkipByteOrderMark:          &altSkipByteOrderMark,
			Encoding:                   sdk.CSVEncodingISO88591,
			RecordDelimiter:            ":",
		})
	defaultAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelCompleteCsv.ResourceReference()).
			HasFileFormatCsv(),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.type", string(sdk.FileFormatTypeCSV))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.record_delimiter", "\\n")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.field_delimiter", ",")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.file_extension", "")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_header", "0")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.parse_header", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.date_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.time_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.timestamp_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.binary_format", string(sdk.BinaryFormatHex))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.escape", "NONE")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.escape_unenclosed_field", "\\\\")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.trim_space", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.field_optionally_enclosed_by", "NONE")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.0", "\\\\N")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.compression", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.error_on_column_count_mismatch", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_blank_lines", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.replace_invalid_characters", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.empty_field_as_null", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_byte_order_mark", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.encoding", string(sdk.CSVEncodingUTF8))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.multi_line", "true")),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelCompleteCsv.ResourceReference()).
			HasFileFormatCsv().
			HasFileFormatCsvCompression(sdk.CSVCompressionGzip).
			HasFileFormatCsvFieldDelimiter("|").
			HasFileFormatCsvSkipHeader(-1).
			HasFileFormatCsvTrimSpace(true).
			HasFileFormatCsvParseHeader(false).
			HasFileFormatCsvMultiLine(true).
			HasFileFormatCsvSkipBlankLines(true).
			HasFileFormatCsvNullIfCount(2).
			HasFileFormatCsvRecordDelimiter(";").
			HasFileFormatCsvFileExtension("csv").
			HasFileFormatCsvDateFormat("AUTO").
			HasFileFormatCsvTimeFormat("AUTO").
			HasFileFormatCsvTimestampFormat("AUTO").
			HasFileFormatCsvBinaryFormat(sdk.BinaryFormatHex).
			HasFileFormatCsvEscape("\\").
			HasFileFormatCsvEscapeUnenclosedField("NONE").
			HasFileFormatCsvFieldOptionallyEnclosedBy("\"").
			HasFileFormatCsvErrorOnColumnCountMismatch(false).
			HasFileFormatCsvReplaceInvalidCharacters(true).
			HasFileFormatCsvEmptyFieldAsNull(true).
			HasFileFormatCsvSkipByteOrderMark(true).
			HasFileFormatCsvEncoding(sdk.CSVEncodingUTF8),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.type", string(sdk.FileFormatTypeCSV))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.record_delimiter", ";")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.field_delimiter", "|")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.file_extension", "csv")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_header", "0")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.parse_header", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.date_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.time_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.timestamp_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.binary_format", string(sdk.BinaryFormatHex))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.escape", "\\\\")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.escape_unenclosed_field", "NONE")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.trim_space", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.field_optionally_enclosed_by", "\\\"")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.#", "2")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.0", "NULL")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.1", "")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.compression", string(sdk.CSVCompressionGzip))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.error_on_column_count_mismatch", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_blank_lines", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.replace_invalid_characters", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.empty_field_as_null", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_byte_order_mark", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.encoding", string(sdk.CSVEncodingUTF8))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.multi_line", "true")),
	}
	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelAlteredCsv.ResourceReference()).
			HasFileFormatCsv().
			HasFileFormatCsvCompression(sdk.CSVCompressionZstd).
			HasFileFormatCsvFieldDelimiter(",").
			HasFileFormatCsvSkipHeader(-1).
			HasFileFormatCsvTrimSpace(false).
			HasFileFormatCsvParseHeader(true).
			HasFileFormatCsvMultiLine(false).
			HasFileFormatCsvSkipBlankLines(false).
			HasFileFormatCsvNullIfCount(1).
			HasFileFormatCsvRecordDelimiter(":").
			HasFileFormatCsvFileExtension("txt").
			HasFileFormatCsvDateFormat("YYYY").
			HasFileFormatCsvTimeFormat("HH24:MI:SS").
			HasFileFormatCsvTimestampFormat("YYYY-MM-DD HH24:MI:SS").
			HasFileFormatCsvBinaryFormat(sdk.BinaryFormatBase64).
			HasFileFormatCsvEscape("\\\\").
			HasFileFormatCsvEscapeUnenclosedField("NONE").
			HasFileFormatCsvFieldOptionallyEnclosedBy("\"").
			HasFileFormatCsvErrorOnColumnCountMismatch(true).
			HasFileFormatCsvReplaceInvalidCharacters(false).
			HasFileFormatCsvEmptyFieldAsNull(false).
			HasFileFormatCsvSkipByteOrderMark(false).
			HasFileFormatCsvEncoding(sdk.CSVEncodingISO88591),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.type", string(sdk.FileFormatTypeCSV))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.record_delimiter", ":")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.field_delimiter", ",")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.file_extension", "txt")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_header", "0")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.parse_header", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.date_format", "YYYY")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.time_format", "HH24:MI:SS")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.timestamp_format", "YYYY-MM-DD HH24:MI:SS")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.binary_format", string(sdk.BinaryFormatBase64))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.escape", "\\\\")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.escape_unenclosed_field", "NONE")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.trim_space", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.field_optionally_enclosed_by", "\\\"")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.0", "NA")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.compression", string(sdk.CSVCompressionZstd))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.error_on_column_count_mismatch", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_blank_lines", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.replace_invalid_characters", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.empty_field_as_null", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_byte_order_mark", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.encoding", string(sdk.CSVEncodingISO88591))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.multi_line", "false")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelCompleteCsv),
				Check: assertThat(t,
					completeAssertions...,
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelCompleteCsv),
				ResourceName:            modelCompleteCsv.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption", "directory", "file_format"},
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelWithoutFileFormat),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithoutFileFormat.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					defaultAssertions...,
				),
			},
			// Set all fields
			{
				Config: accconfig.FromModels(t, modelCompleteCsv),
				Check: assertThat(t,
					completeAssertions...,
				),
			},
			// alter values
			{
				Config: accconfig.FromModels(t, modelAlteredCsv),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelAlteredCsv.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					alteredAssertions...,
				),
			},
			// detect external changes
			{
				PreConfig: func() {
					testClient().Stage.AlterInternalStage(t, sdk.NewAlterInternalStageStageRequest(id).WithFileFormat(sdk.StageFileFormatRequest{
						FileFormatOptions: &sdk.FileFormatOptions{
							CsvOptions: &sdk.FileFormatCsvOptions{
								Compression:                sdk.Pointer(sdk.CSVCompressionZstd),
								FieldDelimiter:             &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer(",")},
								RecordDelimiter:            &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer("\n")},
								MultiLine:                  sdk.Bool(false),
								FileExtension:              sdk.Pointer("txt"),
								SkipBlankLines:             sdk.Bool(false),
								BinaryFormat:               sdk.Pointer(sdk.BinaryFormatBase64),
								TrimSpace:                  sdk.Bool(false),
								NullIf:                     []sdk.NullString{{S: "NA"}},
								ErrorOnColumnCountMismatch: sdk.Bool(true),
								ReplaceInvalidCharacters:   sdk.Bool(false),
								EmptyFieldAsNull:           sdk.Bool(false),
								SkipByteOrderMark:          sdk.Bool(false),
								Encoding:                   sdk.Pointer(sdk.CSVEncodingISO88591),
								ParseHeader:                new(bool),
								DateFormat:                 &sdk.StageFileFormatStringOrAuto{Auto: sdk.Bool(true)},
								TimeFormat:                 &sdk.StageFileFormatStringOrAuto{Auto: sdk.Bool(true)},
								TimestampFormat:            &sdk.StageFileFormatStringOrAuto{Auto: sdk.Bool(true)},
								Escape:                     &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer("\\")},
								EscapeUnenclosedField:      &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer("NONE")},
								FieldOptionallyEnclosedBy:  &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer("\"")},
							},
						},
					}))
				},
				Config: accconfig.FromModels(t, modelCompleteCsv),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteCsv.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					completeAssertions...,
				),
			},
		},
	})
}

func TestAcc_InternalStage_FileFormat_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelBothFormats := model.InternalStageWithId(id).
		WithFileFormatMultipleFormats()
	modelInvalidCompression := model.InternalStageWithId(id).
		WithFileFormatCsvInvalidCompression()
	modelInvalidBinaryFormat := model.InternalStageWithId(id).
		WithFileFormatCsvInvalidBinaryFormat()
	modelInvalidEncoding := model.InternalStageWithId(id).
		WithFileFormatCsvInvalidEncoding()
	modelInvalidBooleanString := model.InternalStageWithId(id).
		WithFileFormatCsvInvalidBooleanString()
	modelInvalidSkipHeader := model.InternalStageWithId(id).
		WithFileFormatCsvInvalidSkipHeader()
	modelInvalidFormatName := model.InternalStageWithId(id).
		WithFileFormatInvalidFormatName()
	modelConflictingOptions := model.InternalStageWithId(id).
		WithFileFormatCsvConflictingOptions()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelBothFormats),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`file_format.0.csv,file_format.0.format_name.* can be specified`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidCompression),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid csv compression: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidBinaryFormat),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid binary format: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidEncoding),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid csv encoding: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidBooleanString),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected .*multi_line.* to be one of \["true" "false"\], got invalid`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidSkipHeader),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected .* to be at least \(0\), got -1`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidFormatName),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Expected SchemaObjectIdentifier identifier type, but got:\nsdk.AccountObjectIdentifier."),
			},
			{
				Config:      accconfig.FromModels(t, modelConflictingOptions),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`file_format.0.csv.0.skip_header.*conflicts with\nfile_format.0.csv.0.parse_header`),
			},
		},
	})
}

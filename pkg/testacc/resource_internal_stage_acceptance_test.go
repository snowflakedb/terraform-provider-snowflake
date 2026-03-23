//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"strconv"
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
		WithFileFormatCsv(sdk.FileFormatCsvOptions{})

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
			// import named format
			{
				Config:                  accconfig.FromModels(t, modelWithNamedFormat),
				ResourceName:            modelWithNamedFormat.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption", "directory"},
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

	modelWithoutFileFormat := model.InternalStageWithId(id)

	modelCompleteCsv := model.InternalStageWithId(id).
		WithFileFormatCsv(sdk.FileFormatCsvOptions{
			Compression:                sdk.Pointer(sdk.CSVCompressionGzip),
			FieldDelimiter:             &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer("|")},
			MultiLine:                  &multiLine,
			FileExtension:              sdk.Pointer("csv"),
			ParseHeader:                &parseHeader,
			SkipBlankLines:             &skipBlankLines,
			DateFormat:                 &sdk.StageFileFormatStringOrAuto{Value: sdk.Pointer("AUTO")},
			TimeFormat:                 &sdk.StageFileFormatStringOrAuto{Value: sdk.Pointer("AUTO")},
			TimestampFormat:            &sdk.StageFileFormatStringOrAuto{Value: sdk.Pointer("AUTO")},
			BinaryFormat:               sdk.Pointer(sdk.BinaryFormatHex),
			Escape:                     &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer(`\`)},
			EscapeUnenclosedField:      &sdk.StageFileFormatStringOrNone{None: sdk.Pointer(true)},
			TrimSpace:                  &trimSpace,
			FieldOptionallyEnclosedBy:  &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer(`"`)},
			NullIf:                     []sdk.NullString{{S: "NULL"}, {S: ""}},
			ErrorOnColumnCountMismatch: &errorOnColumnCountMismatch,
			ReplaceInvalidCharacters:   &replaceInvalidCharacters,
			EmptyFieldAsNull:           &emptyFieldAsNull,
			SkipByteOrderMark:          &skipByteOrderMark,
			Encoding:                   sdk.Pointer(sdk.CSVEncodingUTF8),
			RecordDelimiter:            &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer(";")},
		})

	modelWithSkipHeader := model.InternalStageWithId(id).
		WithFileFormatCsv(sdk.FileFormatCsvOptions{SkipHeader: sdk.Pointer(1)})

	altMultiLine := false
	altParseHeader := true
	altSkipBlankLines := false
	altTrimSpace := false
	altErrorOnColumnCountMismatch := true
	altReplaceInvalidCharacters := false
	altEmptyFieldAsNull := false
	altSkipByteOrderMark := false

	modelAlteredCsv := model.InternalStageWithId(id).
		WithFileFormatCsv(sdk.FileFormatCsvOptions{
			Compression:                sdk.Pointer(sdk.CSVCompressionZstd),
			FieldDelimiter:             &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer(",")},
			MultiLine:                  &altMultiLine,
			FileExtension:              sdk.Pointer("txt"),
			ParseHeader:                &altParseHeader,
			SkipBlankLines:             &altSkipBlankLines,
			DateFormat:                 &sdk.StageFileFormatStringOrAuto{Value: sdk.Pointer("YYYY")},
			TimeFormat:                 &sdk.StageFileFormatStringOrAuto{Value: sdk.Pointer("HH24:MI:SS")},
			TimestampFormat:            &sdk.StageFileFormatStringOrAuto{Value: sdk.Pointer("YYYY-MM-DD HH24:MI:SS")},
			BinaryFormat:               sdk.Pointer(sdk.BinaryFormatBase64),
			Escape:                     &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer(`\\`)},
			EscapeUnenclosedField:      &sdk.StageFileFormatStringOrNone{None: sdk.Pointer(true)},
			TrimSpace:                  &altTrimSpace,
			FieldOptionallyEnclosedBy:  &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer(`"`)},
			NullIf:                     []sdk.NullString{{S: "NA"}},
			ErrorOnColumnCountMismatch: &altErrorOnColumnCountMismatch,
			ReplaceInvalidCharacters:   &altReplaceInvalidCharacters,
			EmptyFieldAsNull:           &altEmptyFieldAsNull,
			SkipByteOrderMark:          &altSkipByteOrderMark,
			Encoding:                   sdk.Pointer(sdk.CSVEncodingISO88591),
			RecordDelimiter:            &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer(":")},
		})
	defaultAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelWithoutFileFormat.ResourceReference()).
			HasFileFormatEmpty(),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.type", string(sdk.FileFormatTypeCSV))),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.record_delimiter", "\\n")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.field_delimiter", ",")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.file_extension", "")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_header", "0")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.parse_header", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.date_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.time_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.timestamp_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.binary_format", string(sdk.BinaryFormatHex))),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.escape", "NONE")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.escape_unenclosed_field", "\\\\")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.trim_space", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.field_optionally_enclosed_by", "NONE")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.0", "\\\\N")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.compression", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.error_on_column_count_mismatch", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_blank_lines", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.replace_invalid_characters", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.empty_field_as_null", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_byte_order_mark", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.encoding", string(sdk.CSVEncodingUTF8))),
		assert.Check(resource.TestCheckResourceAttr(modelWithoutFileFormat.ResourceReference(), "describe_output.0.file_format.0.csv.0.multi_line", "true")),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelCompleteCsv.ResourceReference()).
			HasFileFormatCsv().
			HasFileFormatCsvCompression(sdk.CSVCompressionGzip).
			HasFileFormatCsvFieldDelimiter("|").
			HasFileFormatCsvSkipHeader(-1).
			HasFileFormatCsvTrimSpace(trimSpace).
			HasFileFormatCsvParseHeader(parseHeader).
			HasFileFormatCsvMultiLine(multiLine).
			HasFileFormatCsvSkipBlankLines(skipBlankLines).
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
			HasFileFormatCsvErrorOnColumnCountMismatch(errorOnColumnCountMismatch).
			HasFileFormatCsvReplaceInvalidCharacters(replaceInvalidCharacters).
			HasFileFormatCsvEmptyFieldAsNull(emptyFieldAsNull).
			HasFileFormatCsvSkipByteOrderMark(skipByteOrderMark).
			HasFileFormatCsvEncoding(sdk.CSVEncodingUTF8),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.type", string(sdk.FileFormatTypeCSV))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.record_delimiter", ";")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.field_delimiter", "|")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.file_extension", "csv")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_header", "0")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.parse_header", strconv.FormatBool(parseHeader))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.date_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.time_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.timestamp_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.binary_format", string(sdk.BinaryFormatHex))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.escape", "\\\\")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.escape_unenclosed_field", "NONE")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.trim_space", strconv.FormatBool(trimSpace))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.field_optionally_enclosed_by", "\\\"")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.#", "2")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.0", "NULL")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.1", "")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.compression", string(sdk.CSVCompressionGzip))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.error_on_column_count_mismatch", strconv.FormatBool(errorOnColumnCountMismatch))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_blank_lines", strconv.FormatBool(skipBlankLines))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.replace_invalid_characters", strconv.FormatBool(replaceInvalidCharacters))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.empty_field_as_null", strconv.FormatBool(emptyFieldAsNull))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_byte_order_mark", strconv.FormatBool(skipByteOrderMark))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.encoding", string(sdk.CSVEncodingUTF8))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.multi_line", strconv.FormatBool(multiLine))),
	}
	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelAlteredCsv.ResourceReference()).
			HasFileFormatCsv().
			HasFileFormatCsvCompression(sdk.CSVCompressionZstd).
			HasFileFormatCsvFieldDelimiter(",").
			HasFileFormatCsvSkipHeader(-1).
			HasFileFormatCsvTrimSpace(altTrimSpace).
			HasFileFormatCsvParseHeader(altParseHeader).
			HasFileFormatCsvMultiLine(altMultiLine).
			HasFileFormatCsvSkipBlankLines(altSkipBlankLines).
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
			HasFileFormatCsvErrorOnColumnCountMismatch(altErrorOnColumnCountMismatch).
			HasFileFormatCsvReplaceInvalidCharacters(altReplaceInvalidCharacters).
			HasFileFormatCsvEmptyFieldAsNull(altEmptyFieldAsNull).
			HasFileFormatCsvSkipByteOrderMark(altSkipByteOrderMark).
			HasFileFormatCsvEncoding(sdk.CSVEncodingISO88591),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.type", string(sdk.FileFormatTypeCSV))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.record_delimiter", ":")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.field_delimiter", ",")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.file_extension", "txt")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_header", "0")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.parse_header", strconv.FormatBool(altParseHeader))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.date_format", "YYYY")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.time_format", "HH24:MI:SS")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.timestamp_format", "YYYY-MM-DD HH24:MI:SS")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.binary_format", string(sdk.BinaryFormatBase64))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.escape", "\\\\")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.escape_unenclosed_field", "NONE")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.trim_space", strconv.FormatBool(altTrimSpace))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.field_optionally_enclosed_by", "\\\"")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.null_if.0", "NA")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.compression", string(sdk.CSVCompressionZstd))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.error_on_column_count_mismatch", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_blank_lines", strconv.FormatBool(altSkipBlankLines))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.replace_invalid_characters", strconv.FormatBool(altReplaceInvalidCharacters))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.empty_field_as_null", strconv.FormatBool(altEmptyFieldAsNull))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_byte_order_mark", strconv.FormatBool(altSkipByteOrderMark))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.encoding", string(sdk.CSVEncodingISO88591))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredCsv.ResourceReference(), "describe_output.0.file_format.0.csv.0.multi_line", strconv.FormatBool(altMultiLine))),
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
				Config:            accconfig.FromModels(t, modelCompleteCsv),
				ResourceName:      modelCompleteCsv.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				// File format escape and escape_unenclosed_field are returned excessively escaped from Snowflake.
				// E.g., when escape is set to '\\' in SQL, Snowflake returns '\\\\' in the response.
				// This should be resolved with appropriate team. Alternatively, we can "unescape" such values in SDK.
				// skip_header is skipped due to "-1" default value.
				ImportStateVerifyIgnore: []string{"encryption", "directory", "file_format.0.csv.0.escape", "file_format.0.csv.0.field_optionally_enclosed_by", "file_format.0.csv.0.skip_header"},
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
								Compression:                sdk.Pointer(sdk.CSVCompressionGzip),
								FieldDelimiter:             &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer(",")},
								RecordDelimiter:            &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer("\n")},
								MultiLine:                  sdk.Bool(false),
								FileExtension:              sdk.Pointer("EXT"),
								SkipBlankLines:             sdk.Bool(false),
								BinaryFormat:               sdk.Pointer(sdk.BinaryFormatHex),
								TrimSpace:                  sdk.Bool(true),
								NullIf:                     []sdk.NullString{{S: "EXT"}},
								ErrorOnColumnCountMismatch: sdk.Bool(true),
								ReplaceInvalidCharacters:   sdk.Bool(true),
								EmptyFieldAsNull:           sdk.Bool(false),
								SkipByteOrderMark:          sdk.Bool(true),
								Encoding:                   sdk.Pointer(sdk.CSVEncodingUTF8),
								ParseHeader:                new(bool),
								DateFormat:                 &sdk.StageFileFormatStringOrAuto{Auto: sdk.Pointer(true)},
								TimeFormat:                 &sdk.StageFileFormatStringOrAuto{Auto: sdk.Pointer(true)},
								TimestampFormat:            &sdk.StageFileFormatStringOrAuto{Auto: sdk.Pointer(true)},
								Escape:                     &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer("\\")},
								EscapeUnenclosedField:      &sdk.StageFileFormatStringOrNone{None: sdk.Pointer(true)},
								FieldOptionallyEnclosedBy:  &sdk.StageFileFormatStringOrNone{Value: sdk.Pointer("\"")},
							},
						},
					}))
				},
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
			{
				Config: accconfig.FromModels(t, modelWithSkipHeader),
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelWithSkipHeader.ResourceReference()).
						HasFileFormatCsv().
						HasFileFormatCsvSkipHeader(1),
					assert.Check(resource.TestCheckResourceAttr(modelWithSkipHeader.ResourceReference(), "describe_output.0.file_format.0.csv.0.skip_header", "1")),
				),
			},
		},
	})
}

func TestAcc_InternalStage_FileFormat_AllJsonOptions(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	trimSpace := true
	multiLine := true
	enableOctal := true
	allowDuplicate := true
	stripOuterArray := true
	stripNullValues := true
	skipByteOrderMark := true
	ignoreUtf8Errors := false

	modelBasicJson := model.InternalStageWithId(id).
		WithFileFormatJson(sdk.FileFormatJsonOptions{})

	modelCompleteJson := model.InternalStageWithId(id).
		WithFileFormatJson(sdk.FileFormatJsonOptions{
			Compression:       sdk.Pointer(sdk.JSONCompressionGzip),
			DateFormat:        &sdk.StageFileFormatStringOrAuto{Value: sdk.Pointer("AUTO")},
			TimeFormat:        &sdk.StageFileFormatStringOrAuto{Value: sdk.Pointer("AUTO")},
			TimestampFormat:   &sdk.StageFileFormatStringOrAuto{Value: sdk.Pointer("AUTO")},
			BinaryFormat:      sdk.Pointer(sdk.BinaryFormatHex),
			TrimSpace:         &trimSpace,
			MultiLine:         &multiLine,
			NullIf:            []sdk.NullString{{S: "NULL"}, {S: ""}},
			FileExtension:     sdk.Pointer("json"),
			EnableOctal:       &enableOctal,
			AllowDuplicate:    &allowDuplicate,
			StripOuterArray:   &stripOuterArray,
			StripNullValues:   &stripNullValues,
			SkipByteOrderMark: &skipByteOrderMark,
			IgnoreUtf8Errors:  &ignoreUtf8Errors,
		})

	modelJsonWithReplaceInvalidCharacters := model.InternalStageWithId(id).
		WithFileFormatJson(sdk.FileFormatJsonOptions{
			ReplaceInvalidCharacters: sdk.Pointer(true),
		})

	altTrimSpace := false
	altMultiLine := false
	altEnableOctal := false
	altAllowDuplicate := false
	altStripOuterArray := false
	altStripNullValues := false
	altIgnoreUtf8Errors := true
	altSkipByteOrderMark := false

	modelAlteredJson := model.InternalStageWithId(id).
		WithFileFormatJson(sdk.FileFormatJsonOptions{
			Compression:       sdk.Pointer(sdk.JSONCompressionZstd),
			DateFormat:        &sdk.StageFileFormatStringOrAuto{Value: sdk.Pointer("YYYY")},
			TimeFormat:        &sdk.StageFileFormatStringOrAuto{Value: sdk.Pointer("HH24:MI:SS")},
			TimestampFormat:   &sdk.StageFileFormatStringOrAuto{Value: sdk.Pointer("YYYY-MM-DD HH24:MI:SS")},
			BinaryFormat:      sdk.Pointer(sdk.BinaryFormatBase64),
			TrimSpace:         &altTrimSpace,
			MultiLine:         &altMultiLine,
			NullIf:            []sdk.NullString{{S: "NA"}},
			FileExtension:     sdk.Pointer("txt"),
			EnableOctal:       &altEnableOctal,
			AllowDuplicate:    &altAllowDuplicate,
			StripOuterArray:   &altStripOuterArray,
			StripNullValues:   &altStripNullValues,
			IgnoreUtf8Errors:  &altIgnoreUtf8Errors,
			SkipByteOrderMark: &altSkipByteOrderMark,
		})

	defaultAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelBasicJson.ResourceReference()).
			HasFileFormatJson(),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.type", string(sdk.FileFormatTypeJSON))),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.compression", string(sdk.JSONCompressionAuto))),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.date_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.time_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.timestamp_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.binary_format", string(sdk.BinaryFormatHex))),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.trim_space", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.multi_line", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.null_if.#", "0")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.file_extension", "")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.enable_octal", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.allow_duplicate", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.strip_outer_array", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.strip_null_values", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.replace_invalid_characters", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.ignore_utf8_errors", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicJson.ResourceReference(), "describe_output.0.file_format.0.json.0.skip_byte_order_mark", "true")),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelCompleteJson.ResourceReference()).
			HasFileFormatJson().
			HasFileFormatJsonCompression(sdk.JSONCompressionGzip).
			HasFileFormatJsonDateFormat("AUTO").
			HasFileFormatJsonTimeFormat("AUTO").
			HasFileFormatJsonTimestampFormat("AUTO").
			HasFileFormatJsonBinaryFormat(sdk.BinaryFormatHex).
			HasFileFormatJsonTrimSpace(trimSpace).
			HasFileFormatJsonMultiLine(multiLine).
			HasFileFormatJsonNullIfCount(2).
			HasFileFormatJsonFileExtension("json").
			HasFileFormatJsonEnableOctal(enableOctal).
			HasFileFormatJsonAllowDuplicate(allowDuplicate).
			HasFileFormatJsonStripOuterArray(stripOuterArray).
			HasFileFormatJsonStripNullValues(stripNullValues).
			HasFileFormatJsonReplaceInvalidCharactersString(r.BooleanDefault).
			HasFileFormatJsonSkipByteOrderMark(skipByteOrderMark),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.type", string(sdk.FileFormatTypeJSON))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.compression", string(sdk.JSONCompressionGzip))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.date_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.time_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.timestamp_format", "AUTO")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.binary_format", string(sdk.BinaryFormatHex))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.trim_space", strconv.FormatBool(trimSpace))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.multi_line", strconv.FormatBool(multiLine))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.null_if.#", "2")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.null_if.0", "NULL")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.null_if.1", "")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.file_extension", "json")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.enable_octal", strconv.FormatBool(enableOctal))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.allow_duplicate", strconv.FormatBool(allowDuplicate))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.strip_outer_array", strconv.FormatBool(stripOuterArray))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.strip_null_values", strconv.FormatBool(stripNullValues))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.replace_invalid_characters", r.BooleanFalse)),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.ignore_utf8_errors", strconv.FormatBool(ignoreUtf8Errors))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteJson.ResourceReference(), "describe_output.0.file_format.0.json.0.skip_byte_order_mark", strconv.FormatBool(skipByteOrderMark))),
	}

	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelAlteredJson.ResourceReference()).
			HasFileFormatJson().
			HasFileFormatJsonCompression(sdk.JSONCompressionZstd).
			HasFileFormatJsonDateFormat("YYYY").
			HasFileFormatJsonTimeFormat("HH24:MI:SS").
			HasFileFormatJsonTimestampFormat("YYYY-MM-DD HH24:MI:SS").
			HasFileFormatJsonBinaryFormat(sdk.BinaryFormatBase64).
			HasFileFormatJsonTrimSpace(altTrimSpace).
			HasFileFormatJsonMultiLine(altMultiLine).
			HasFileFormatJsonNullIfCount(1).
			HasFileFormatJsonFileExtension("txt").
			HasFileFormatJsonEnableOctal(altEnableOctal).
			HasFileFormatJsonAllowDuplicate(altAllowDuplicate).
			HasFileFormatJsonStripOuterArray(altStripOuterArray).
			HasFileFormatJsonStripNullValues(altStripNullValues).
			HasFileFormatJsonIgnoreUtf8Errors(altIgnoreUtf8Errors).
			HasFileFormatJsonSkipByteOrderMark(altSkipByteOrderMark),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.type", string(sdk.FileFormatTypeJSON))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.compression", string(sdk.JSONCompressionZstd))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.date_format", "YYYY")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.time_format", "HH24:MI:SS")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.timestamp_format", "YYYY-MM-DD HH24:MI:SS")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.binary_format", string(sdk.BinaryFormatBase64))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.trim_space", strconv.FormatBool(altTrimSpace))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.multi_line", strconv.FormatBool(altMultiLine))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.null_if.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.null_if.0", "NA")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.file_extension", "txt")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.enable_octal", strconv.FormatBool(altEnableOctal))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.allow_duplicate", strconv.FormatBool(altAllowDuplicate))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.strip_outer_array", strconv.FormatBool(altStripOuterArray))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.strip_null_values", strconv.FormatBool(altStripNullValues))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.replace_invalid_characters", r.BooleanFalse)),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.ignore_utf8_errors", strconv.FormatBool(altIgnoreUtf8Errors))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredJson.ResourceReference(), "describe_output.0.file_format.0.json.0.skip_byte_order_mark", strconv.FormatBool(altSkipByteOrderMark))),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelCompleteJson),
				Check: assertThat(t,
					completeAssertions...,
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelCompleteJson),
				ResourceName:            modelCompleteJson.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption", "directory", "file_format.0.json.0.ignore_utf8_errors", "file_format.0.json.0.replace_invalid_characters"},
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasicJson),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasicJson.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					defaultAssertions...,
				),
			},
			// Set all fields
			{
				Config: accconfig.FromModels(t, modelCompleteJson),
				Check: assertThat(t,
					completeAssertions...,
				),
			},
			// alter values
			{
				Config: accconfig.FromModels(t, modelAlteredJson),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelAlteredJson.ResourceReference(), plancheck.ResourceActionUpdate),
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
							JsonOptions: &sdk.FileFormatJsonOptions{
								Compression:              sdk.Pointer(sdk.JSONCompressionGzip),
								DateFormat:               &sdk.StageFileFormatStringOrAuto{Auto: sdk.Bool(true)},
								TimeFormat:               &sdk.StageFileFormatStringOrAuto{Auto: sdk.Bool(true)},
								TimestampFormat:          &sdk.StageFileFormatStringOrAuto{Auto: sdk.Bool(true)},
								BinaryFormat:             sdk.Pointer(sdk.BinaryFormatHex),
								TrimSpace:                sdk.Bool(true),
								MultiLine:                sdk.Bool(true),
								NullIf:                   []sdk.NullString{{S: "EXT"}},
								FileExtension:            sdk.Pointer("EXT"),
								EnableOctal:              sdk.Bool(true),
								AllowDuplicate:           sdk.Bool(true),
								StripOuterArray:          sdk.Bool(true),
								StripNullValues:          sdk.Bool(true),
								ReplaceInvalidCharacters: sdk.Bool(true),
								SkipByteOrderMark:        sdk.Bool(true),
							},
						},
					}))
				},
				Config: accconfig.FromModels(t, modelAlteredJson),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelAlteredJson.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					alteredAssertions...,
				),
			},
			{
				Config: accconfig.FromModels(t, modelJsonWithReplaceInvalidCharacters),
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelJsonWithReplaceInvalidCharacters.ResourceReference()).
						HasFileFormatJson().
						HasFileFormatJsonReplaceInvalidCharacters(true),
					assert.Check(resource.TestCheckResourceAttr(modelJsonWithReplaceInvalidCharacters.ResourceReference(), "describe_output.0.file_format.0.json.0.replace_invalid_characters", "true")),
				),
			},
		},
	})
}

func TestAcc_InternalStage_FileFormat_AllAvroOptions(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	trimSpace := true
	replaceInvalidCharacters := true

	modelBasicAvro := model.InternalStageWithId(id).
		WithFileFormatAvro(sdk.FileFormatAvroOptions{})

	modelCompleteAvro := model.InternalStageWithId(id).
		WithFileFormatAvro(sdk.FileFormatAvroOptions{
			Compression:              sdk.Pointer(sdk.AvroCompressionGzip),
			TrimSpace:                &trimSpace,
			ReplaceInvalidCharacters: &replaceInvalidCharacters,
			NullIf:                   []sdk.NullString{{S: "NULL"}, {S: ""}},
		})

	altTrimSpace := false
	altReplaceInvalidCharacters := false

	modelAlteredAvro := model.InternalStageWithId(id).
		WithFileFormatAvro(sdk.FileFormatAvroOptions{
			Compression:              sdk.Pointer(sdk.AvroCompressionZstd),
			TrimSpace:                &altTrimSpace,
			ReplaceInvalidCharacters: &altReplaceInvalidCharacters,
			NullIf:                   []sdk.NullString{{S: "NA"}},
		})

	defaultAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelBasicAvro.ResourceReference()).
			HasFileFormatAvro().
			HasFileFormatAvroTrimSpaceString(r.BooleanDefault).
			HasFileFormatAvroReplaceInvalidCharactersString(r.BooleanDefault),
		assert.Check(resource.TestCheckResourceAttr(modelBasicAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.type", string(sdk.FileFormatTypeAvro))),
		assert.Check(resource.TestCheckResourceAttr(modelBasicAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.compression", string(sdk.AvroCompressionAuto))),
		assert.Check(resource.TestCheckResourceAttr(modelBasicAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.trim_space", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.replace_invalid_characters", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.null_if.#", "0")),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelCompleteAvro.ResourceReference()).
			HasFileFormatAvro().
			HasFileFormatAvroCompression(sdk.AvroCompressionGzip).
			HasFileFormatAvroTrimSpace(trimSpace).
			HasFileFormatAvroReplaceInvalidCharacters(replaceInvalidCharacters).
			HasFileFormatAvroNullIfCount(2),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.type", string(sdk.FileFormatTypeAvro))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.compression", string(sdk.AvroCompressionGzip))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.trim_space", strconv.FormatBool(trimSpace))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.replace_invalid_characters", strconv.FormatBool(replaceInvalidCharacters))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.null_if.#", "2")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.null_if.0", "NULL")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.null_if.1", "")),
	}

	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelAlteredAvro.ResourceReference()).
			HasFileFormatAvro().
			HasFileFormatAvroCompression(sdk.AvroCompressionZstd).
			HasFileFormatAvroTrimSpace(altTrimSpace).
			HasFileFormatAvroReplaceInvalidCharacters(altReplaceInvalidCharacters).
			HasFileFormatAvroNullIfCount(1),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.type", string(sdk.FileFormatTypeAvro))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.compression", string(sdk.AvroCompressionZstd))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.trim_space", strconv.FormatBool(altTrimSpace))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.replace_invalid_characters", strconv.FormatBool(altReplaceInvalidCharacters))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.null_if.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredAvro.ResourceReference(), "describe_output.0.file_format.0.avro.0.null_if.0", "NA")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelCompleteAvro),
				Check: assertThat(t,
					completeAssertions...,
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelCompleteAvro),
				ResourceName:            modelCompleteAvro.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption", "directory"},
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasicAvro),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasicAvro.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					defaultAssertions...,
				),
			},
			// Set all fields
			{
				Config: accconfig.FromModels(t, modelCompleteAvro),
				Check: assertThat(t,
					completeAssertions...,
				),
			},
			// alter values
			{
				Config: accconfig.FromModels(t, modelAlteredAvro),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelAlteredAvro.ResourceReference(), plancheck.ResourceActionUpdate),
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
							AvroOptions: &sdk.FileFormatAvroOptions{
								Compression:              sdk.Pointer(sdk.AvroCompressionGzip),
								TrimSpace:                sdk.Bool(true),
								ReplaceInvalidCharacters: sdk.Bool(true),
								NullIf:                   []sdk.NullString{{S: "EXT"}},
							},
						},
					}))
				},
				Config: accconfig.FromModels(t, modelAlteredAvro),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelAlteredAvro.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					alteredAssertions...,
				),
			},
		},
	})
}

func TestAcc_InternalStage_FileFormat_AllOrcOptions(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	trimSpace := true
	replaceInvalidCharacters := true

	modelBasicOrc := model.InternalStageWithId(id).
		WithFileFormatOrc(sdk.FileFormatOrcOptions{})

	modelCompleteOrc := model.InternalStageWithId(id).
		WithFileFormatOrc(sdk.FileFormatOrcOptions{
			TrimSpace:                &trimSpace,
			ReplaceInvalidCharacters: &replaceInvalidCharacters,
			NullIf:                   []sdk.NullString{{S: "NULL"}, {S: ""}},
		})

	altTrimSpace := false
	altReplaceInvalidCharacters := false

	modelAlteredOrc := model.InternalStageWithId(id).
		WithFileFormatOrc(sdk.FileFormatOrcOptions{
			TrimSpace:                &altTrimSpace,
			ReplaceInvalidCharacters: &altReplaceInvalidCharacters,
			NullIf:                   []sdk.NullString{{S: "NA"}},
		})

	defaultAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelBasicOrc.ResourceReference()).
			HasFileFormatOrc().
			HasFileFormatOrcTrimSpaceString(r.BooleanDefault).
			HasFileFormatOrcReplaceInvalidCharactersString(r.BooleanDefault),
		assert.Check(resource.TestCheckResourceAttr(modelBasicOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.type", string(sdk.FileFormatTypeORC))),
		assert.Check(resource.TestCheckResourceAttr(modelBasicOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.trim_space", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.replace_invalid_characters", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.null_if.#", "0")),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelCompleteOrc.ResourceReference()).
			HasFileFormatOrc().
			HasFileFormatOrcTrimSpace(trimSpace).
			HasFileFormatOrcReplaceInvalidCharacters(replaceInvalidCharacters).
			HasFileFormatOrcNullIfCount(2),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.type", string(sdk.FileFormatTypeORC))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.trim_space", strconv.FormatBool(trimSpace))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.replace_invalid_characters", strconv.FormatBool(replaceInvalidCharacters))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.null_if.#", "2")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.null_if.0", "NULL")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.null_if.1", "")),
	}

	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelAlteredOrc.ResourceReference()).
			HasFileFormatOrc().
			HasFileFormatOrcTrimSpace(altTrimSpace).
			HasFileFormatOrcReplaceInvalidCharacters(altReplaceInvalidCharacters).
			HasFileFormatOrcNullIfCount(1),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.type", string(sdk.FileFormatTypeORC))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.trim_space", strconv.FormatBool(altTrimSpace))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.replace_invalid_characters", strconv.FormatBool(altReplaceInvalidCharacters))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.null_if.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredOrc.ResourceReference(), "describe_output.0.file_format.0.orc.0.null_if.0", "NA")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelCompleteOrc),
				Check: assertThat(t,
					completeAssertions...,
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelCompleteOrc),
				ResourceName:            modelCompleteOrc.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption", "directory"},
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasicOrc),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasicOrc.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					defaultAssertions...,
				),
			},
			// Set all fields
			{
				Config: accconfig.FromModels(t, modelCompleteOrc),
				Check: assertThat(t,
					completeAssertions...,
				),
			},
			// alter values
			{
				Config: accconfig.FromModels(t, modelAlteredOrc),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelAlteredOrc.ResourceReference(), plancheck.ResourceActionUpdate),
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
							OrcOptions: &sdk.FileFormatOrcOptions{
								TrimSpace:                sdk.Bool(true),
								ReplaceInvalidCharacters: sdk.Bool(true),
								NullIf:                   []sdk.NullString{{S: "EXT"}},
							},
						},
					}))
				},
				Config: accconfig.FromModels(t, modelAlteredOrc),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelAlteredOrc.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					alteredAssertions...,
				),
			},
		},
	})
}

func TestAcc_InternalStage_FileFormat_AllParquetOptions(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	binaryAsText := true
	useLogicalType := true
	trimSpace := true
	useVectorizedScanner := true
	replaceInvalidCharacters := true

	modelBasicParquet := model.InternalStageWithId(id).
		WithFileFormatParquet(sdk.FileFormatParquetOptions{})

	modelCompleteParquet := model.InternalStageWithId(id).
		WithFileFormatParquet(sdk.FileFormatParquetOptions{
			Compression:              sdk.Pointer(sdk.ParquetCompressionSnappy),
			BinaryAsText:             &binaryAsText,
			UseLogicalType:           &useLogicalType,
			TrimSpace:                &trimSpace,
			UseVectorizedScanner:     &useVectorizedScanner,
			ReplaceInvalidCharacters: &replaceInvalidCharacters,
			NullIf:                   []sdk.NullString{{S: "NULL"}, {S: ""}},
		})

	altBinaryAsText := false
	altUseLogicalType := false
	altTrimSpace := false
	altUseVectorizedScanner := false
	altReplaceInvalidCharacters := false

	modelAlteredParquet := model.InternalStageWithId(id).
		WithFileFormatParquet(sdk.FileFormatParquetOptions{
			Compression:              sdk.Pointer(sdk.ParquetCompressionLzo),
			BinaryAsText:             &altBinaryAsText,
			UseLogicalType:           &altUseLogicalType,
			TrimSpace:                &altTrimSpace,
			UseVectorizedScanner:     &altUseVectorizedScanner,
			ReplaceInvalidCharacters: &altReplaceInvalidCharacters,
			NullIf:                   []sdk.NullString{{S: "NA"}},
		})

	defaultAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelBasicParquet.ResourceReference()).
			HasFileFormatParquet().
			HasFileFormatParquetTrimSpaceString(r.BooleanDefault).
			HasFileFormatParquetBinaryAsTextString(r.BooleanDefault).
			HasFileFormatParquetUseLogicalTypeString(r.BooleanDefault).
			HasFileFormatParquetUseVectorizedScannerString(r.BooleanDefault).
			HasFileFormatParquetReplaceInvalidCharactersString(r.BooleanDefault),
		assert.Check(resource.TestCheckResourceAttr(modelBasicParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.type", string(sdk.FileFormatTypeParquet))),
		assert.Check(resource.TestCheckResourceAttr(modelBasicParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.compression", string(sdk.ParquetCompressionAuto))),
		assert.Check(resource.TestCheckResourceAttr(modelBasicParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.binary_as_text", "true")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.use_logical_type", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.trim_space", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.use_vectorized_scanner", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.replace_invalid_characters", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.null_if.#", "0")),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelCompleteParquet.ResourceReference()).
			HasFileFormatParquet().
			HasFileFormatParquetCompression(sdk.ParquetCompressionSnappy).
			HasFileFormatParquetBinaryAsText(binaryAsText).
			HasFileFormatParquetUseLogicalType(useLogicalType).
			HasFileFormatParquetTrimSpace(trimSpace).
			HasFileFormatParquetUseVectorizedScanner(useVectorizedScanner).
			HasFileFormatParquetReplaceInvalidCharacters(replaceInvalidCharacters).
			HasFileFormatParquetNullIfCount(2),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.type", string(sdk.FileFormatTypeParquet))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.compression", string(sdk.ParquetCompressionSnappy))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.binary_as_text", strconv.FormatBool(binaryAsText))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.use_logical_type", strconv.FormatBool(useLogicalType))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.trim_space", strconv.FormatBool(trimSpace))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.use_vectorized_scanner", strconv.FormatBool(useVectorizedScanner))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.replace_invalid_characters", strconv.FormatBool(replaceInvalidCharacters))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.null_if.#", "2")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.null_if.0", "NULL")),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.null_if.1", "")),
	}

	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelAlteredParquet.ResourceReference()).
			HasFileFormatParquet().
			HasFileFormatParquetCompression(sdk.ParquetCompressionLzo).
			HasFileFormatParquetBinaryAsText(altBinaryAsText).
			HasFileFormatParquetUseLogicalType(altUseLogicalType).
			HasFileFormatParquetTrimSpace(altTrimSpace).
			HasFileFormatParquetUseVectorizedScanner(altUseVectorizedScanner).
			HasFileFormatParquetReplaceInvalidCharacters(altReplaceInvalidCharacters).
			HasFileFormatParquetNullIfCount(1),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.type", string(sdk.FileFormatTypeParquet))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.compression", string(sdk.ParquetCompressionLzo))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.binary_as_text", strconv.FormatBool(altBinaryAsText))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.use_logical_type", strconv.FormatBool(altUseLogicalType))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.trim_space", strconv.FormatBool(altTrimSpace))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.use_vectorized_scanner", strconv.FormatBool(altUseVectorizedScanner))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.replace_invalid_characters", strconv.FormatBool(altReplaceInvalidCharacters))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.null_if.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredParquet.ResourceReference(), "describe_output.0.file_format.0.parquet.0.null_if.0", "NA")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelCompleteParquet),
				Check: assertThat(t,
					completeAssertions...,
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelCompleteParquet),
				ResourceName:            modelCompleteParquet.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption", "directory"},
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasicParquet),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasicParquet.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					defaultAssertions...,
				),
			},
			// Set all fields
			{
				Config: accconfig.FromModels(t, modelCompleteParquet),
				Check: assertThat(t,
					completeAssertions...,
				),
			},
			// External change detection
			{
				PreConfig: func() {
					testClient().Stage.AlterInternalStage(t, sdk.NewAlterInternalStageStageRequest(id).WithFileFormat(sdk.StageFileFormatRequest{
						FileFormatOptions: &sdk.FileFormatOptions{
							ParquetOptions: &sdk.FileFormatParquetOptions{
								Compression:              sdk.Pointer(sdk.ParquetCompressionLzo),
								BinaryAsText:             sdk.Bool(true),
								UseLogicalType:           sdk.Bool(true),
								TrimSpace:                sdk.Bool(true),
								UseVectorizedScanner:     sdk.Bool(true),
								ReplaceInvalidCharacters: sdk.Bool(true),
								NullIf:                   []sdk.NullString{{S: "EXT"}},
							},
						},
					}))
				},
				Config: accconfig.FromModels(t, modelAlteredParquet),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelAlteredParquet.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					alteredAssertions...,
				),
			},
		},
	})
}

func TestAcc_InternalStage_FileFormat_AllXmlOptions(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	preserveSpace := true
	stripOuterElement := true
	disableAutoConvert := true
	skipByteOrderMark := true
	ignoreUtf8Errors := false

	modelBasicXml := model.InternalStageWithId(id).
		WithFileFormatXml(sdk.FileFormatXmlOptions{})

	modelCompleteXml := model.InternalStageWithId(id).
		WithFileFormatXml(sdk.FileFormatXmlOptions{
			Compression:        sdk.Pointer(sdk.XMLCompressionGzip),
			PreserveSpace:      &preserveSpace,
			StripOuterElement:  &stripOuterElement,
			DisableAutoConvert: &disableAutoConvert,
			SkipByteOrderMark:  &skipByteOrderMark,
			IgnoreUtf8Errors:   &ignoreUtf8Errors,
		})

	modelXmlWithReplaceInvalidCharacters := model.InternalStageWithId(id).
		WithFileFormatXml(sdk.FileFormatXmlOptions{
			ReplaceInvalidCharacters: sdk.Pointer(true),
		})

	altPreserveSpace := false
	altStripOuterElement := false
	altDisableAutoConvert := false
	altSkipByteOrderMark := false
	altIgnoreUtf8Errors := true

	modelAlteredXml := model.InternalStageWithId(id).
		WithFileFormatXml(sdk.FileFormatXmlOptions{
			Compression:        sdk.Pointer(sdk.XMLCompressionZstd),
			PreserveSpace:      &altPreserveSpace,
			StripOuterElement:  &altStripOuterElement,
			DisableAutoConvert: &altDisableAutoConvert,
			SkipByteOrderMark:  &altSkipByteOrderMark,
			IgnoreUtf8Errors:   &altIgnoreUtf8Errors,
		})

	defaultAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelBasicXml.ResourceReference()).
			HasFileFormatXml(),
		assert.Check(resource.TestCheckResourceAttr(modelBasicXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.type", string(sdk.FileFormatTypeXML))),
		assert.Check(resource.TestCheckResourceAttr(modelBasicXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.compression", string(sdk.XMLCompressionAuto))),
		assert.Check(resource.TestCheckResourceAttr(modelBasicXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.ignore_utf8_errors", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.preserve_space", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.strip_outer_element", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.disable_auto_convert", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.replace_invalid_characters", "false")),
		assert.Check(resource.TestCheckResourceAttr(modelBasicXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.skip_byte_order_mark", "true")),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelCompleteXml.ResourceReference()).
			HasFileFormatXml().
			HasFileFormatXmlCompression(sdk.XMLCompressionGzip).
			HasFileFormatXmlPreserveSpace(preserveSpace).
			HasFileFormatXmlStripOuterElement(stripOuterElement).
			HasFileFormatXmlDisableAutoConvert(disableAutoConvert).
			HasFileFormatXmlSkipByteOrderMark(skipByteOrderMark).
			HasFileFormatXmlReplaceInvalidCharactersString(r.BooleanDefault),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.type", string(sdk.FileFormatTypeXML))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.compression", string(sdk.XMLCompressionGzip))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.ignore_utf8_errors", strconv.FormatBool(ignoreUtf8Errors))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.preserve_space", strconv.FormatBool(preserveSpace))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.strip_outer_element", strconv.FormatBool(stripOuterElement))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.disable_auto_convert", strconv.FormatBool(disableAutoConvert))),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.replace_invalid_characters", r.BooleanFalse)),
		assert.Check(resource.TestCheckResourceAttr(modelCompleteXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.skip_byte_order_mark", strconv.FormatBool(skipByteOrderMark))),
	}

	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.InternalStageResource(t, modelAlteredXml.ResourceReference()).
			HasFileFormatXml().
			HasFileFormatXmlCompression(sdk.XMLCompressionZstd).
			HasFileFormatXmlPreserveSpace(altPreserveSpace).
			HasFileFormatXmlStripOuterElement(altStripOuterElement).
			HasFileFormatXmlDisableAutoConvert(altDisableAutoConvert).
			HasFileFormatXmlIgnoreUtf8Errors(altIgnoreUtf8Errors).
			HasFileFormatXmlSkipByteOrderMark(altSkipByteOrderMark),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.type", string(sdk.FileFormatTypeXML))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.compression", string(sdk.XMLCompressionZstd))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.ignore_utf8_errors", strconv.FormatBool(altIgnoreUtf8Errors))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.preserve_space", strconv.FormatBool(altPreserveSpace))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.strip_outer_element", strconv.FormatBool(altStripOuterElement))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.disable_auto_convert", strconv.FormatBool(altDisableAutoConvert))),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.replace_invalid_characters", r.BooleanFalse)),
		assert.Check(resource.TestCheckResourceAttr(modelAlteredXml.ResourceReference(), "describe_output.0.file_format.0.xml.0.skip_byte_order_mark", strconv.FormatBool(altSkipByteOrderMark))),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelCompleteXml),
				Check: assertThat(t,
					completeAssertions...,
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelCompleteXml),
				ResourceName:            modelCompleteXml.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption", "directory", "file_format.0.xml.0.ignore_utf8_errors", "file_format.0.xml.0.replace_invalid_characters"},
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasicXml),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasicXml.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					defaultAssertions...,
				),
			},
			// Set all fields
			{
				Config: accconfig.FromModels(t, modelCompleteXml),
				Check: assertThat(t,
					completeAssertions...,
				),
			},
			// alter values
			{
				Config: accconfig.FromModels(t, modelAlteredXml),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelAlteredXml.ResourceReference(), plancheck.ResourceActionUpdate),
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
							XmlOptions: &sdk.FileFormatXmlOptions{
								Compression:              sdk.Pointer(sdk.XMLCompressionGzip),
								PreserveSpace:            sdk.Bool(true),
								StripOuterElement:        sdk.Bool(true),
								DisableAutoConvert:       sdk.Bool(true),
								ReplaceInvalidCharacters: sdk.Bool(true),
								SkipByteOrderMark:        sdk.Bool(true),
							},
						},
					}))
				},
				Config: accconfig.FromModels(t, modelAlteredXml),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelAlteredXml.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					alteredAssertions...,
				),
			},
			{
				Config: accconfig.FromModels(t, modelXmlWithReplaceInvalidCharacters),
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelXmlWithReplaceInvalidCharacters.ResourceReference()).
						HasFileFormatXml().
						HasFileFormatXmlReplaceInvalidCharacters(true),
					assert.Check(resource.TestCheckResourceAttr(modelXmlWithReplaceInvalidCharacters.ResourceReference(), "describe_output.0.file_format.0.xml.0.replace_invalid_characters", "true")),
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
	modelJsonInvalidCompression := model.InternalStageWithId(id).
		WithFileFormatJsonInvalidCompression()
	modelJsonInvalidBinaryFormat := model.InternalStageWithId(id).
		WithFileFormatJsonInvalidBinaryFormat()
	modelJsonInvalidBooleanString := model.InternalStageWithId(id).
		WithFileFormatJsonInvalidBooleanString()
	modelJsonConflictingOptions := model.InternalStageWithId(id).
		WithFileFormatJsonConflictingOptions()
	modelParquetInvalidCompression := model.InternalStageWithId(id).
		WithFileFormatParquetInvalidCompression()
	modelParquetInvalidBooleanString := model.InternalStageWithId(id).
		WithFileFormatParquetInvalidBooleanString()
	modelXmlInvalidCompression := model.InternalStageWithId(id).
		WithFileFormatXmlInvalidCompression()
	modelXmlInvalidBooleanString := model.InternalStageWithId(id).
		WithFileFormatXmlInvalidBooleanString()
	modelXmlConflictingOptions := model.InternalStageWithId(id).
		WithFileFormatXmlConflictingOptions()

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
				ExpectError: regexp.MustCompile(`file_format.0.avro,file_format.0.csv,file_format.0.format_name,file_format.0.json,file_format.0.orc,file_format.0.parquet,file_format.0.xml`),
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
			{
				Config:      accconfig.FromModels(t, modelJsonInvalidCompression),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid json compression: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelJsonInvalidBinaryFormat),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid binary format: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelJsonInvalidBooleanString),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected .*multi_line.* to be one of \["true" "false"\], got invalid`),
			},
			{
				Config:      accconfig.FromModels(t, modelJsonConflictingOptions),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`file_format.0.json.0.replace_invalid_characters.*conflicts with\nfile_format.0.json.0.ignore_utf8_errors`),
			},
			{
				Config:      accconfig.FromModels(t, modelParquetInvalidCompression),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid parquet compression: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelParquetInvalidBooleanString),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected .*trim_space.* to be one of \["true" "false"\], got invalid`),
			},
			{
				Config:      accconfig.FromModels(t, modelXmlInvalidCompression),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid xml compression: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelXmlInvalidBooleanString),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected .*preserve_space.* to be one of \["true" "false"\], got invalid`),
			},
			{
				Config:      accconfig.FromModels(t, modelXmlConflictingOptions),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`file_format.0.xml.0.replace_invalid_characters.*conflicts with\nfile_format.0.xml.0.ignore_utf8_errors`),
			},
		},
	})
}

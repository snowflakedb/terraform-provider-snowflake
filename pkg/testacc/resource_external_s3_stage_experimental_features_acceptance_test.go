//go:build account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_Experimental_ExternalS3Stage_ImportJsonBooleanDefaults verifies that importing an external S3 stage
// without the IMPORT_BOOLEAN_DEFAULT experiment causes a permadiff (non-empty plan), and that enabling
// the experiment fixes it by setting BooleanDefault fields to "default" instead of the actual Snowflake value.
// Regression test for https://github.com/snowflakedb/terraform-provider-snowflake/issues/4549.
func TestAcc_Experimental_ExternalS3Stage_ImportJsonBooleanDefaults(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	awsUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	// storageIntegrationId := ids.PrecreatedS3StorageIntegration

	providerModelWithExperiment := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportBooleanDefault)

	experimentFactory := providerFactoryUsingCache("TestAcc_Experimental_ExternalS3Stage_ImportJsonBooleanDefaults")

	stageModel := model.ExternalS3StageWithId(id, awsUrl).
		WithDirectoryEnabledAndOptions(sdk.StageS3CommonDirectoryTableOptionsRequest{
			Enable:      false,
			AutoRefresh: sdk.Pointer(false),
		}).
		WithFileFormatJson(sdk.FileFormatJsonOptions{
			BinaryFormat: sdk.Pointer(sdk.BinaryFormatHex),
			Compression:  sdk.Pointer(sdk.JSONCompressionAuto),
			DateFormat: sdk.Pointer(sdk.StageFileFormatStringOrAuto{
				Value: sdk.Pointer("AUTO"),
			}),
			TimeFormat: sdk.Pointer(sdk.StageFileFormatStringOrAuto{
				Value: sdk.Pointer("AUTO"),
			}),
			TimestampFormat: sdk.Pointer(sdk.StageFileFormatStringOrAuto{
				Value: sdk.Pointer("AUTO"),
			}),
		})
	stageModelPure := model.ExternalS3StageWithId(id, awsUrl).
		WithFileFormatJson(sdk.FileFormatJsonOptions{
			BinaryFormat: sdk.Pointer(sdk.BinaryFormatHex),
			Compression:  sdk.Pointer(sdk.JSONCompressionAuto),
			DateFormat: sdk.Pointer(sdk.StageFileFormatStringOrAuto{
				Value: sdk.Pointer("AUTO"),
			}),
			TimeFormat: sdk.Pointer(sdk.StageFileFormatStringOrAuto{
				Value: sdk.Pointer("AUTO"),
			}),
			TimestampFormat: sdk.Pointer(sdk.StageFileFormatStringOrAuto{
				Value: sdk.Pointer("AUTO"),
			}),
		})
	createStage := func() {
		_, dropStage := testClient().Stage.CreateStageOnS3WithRequest(t, sdk.NewCreateOnS3StageRequest(id, *sdk.NewExternalS3StageParamsRequest(awsUrl)).
			// WithStorageIntegration(storageIntegrationId)
			WithFileFormat(sdk.StageFileFormatRequest{
				FileFormatOptions: &sdk.FileFormatOptions{
					JsonOptions: &sdk.FileFormatJsonOptions{
						BinaryFormat: sdk.Pointer(sdk.BinaryFormatHex),
						Compression:  sdk.Pointer(sdk.JSONCompressionAuto),
						DateFormat: sdk.Pointer(sdk.StageFileFormatStringOrAuto{
							Value: sdk.Pointer("AUTO"),
						}),
						TimeFormat: sdk.Pointer(sdk.StageFileFormatStringOrAuto{
							Value: sdk.Pointer("AUTO"),
						}),
						TimestampFormat: sdk.Pointer(sdk.StageFileFormatStringOrAuto{
							Value: sdk.Pointer("AUTO"),
						}),
					},
				},
			}))
		t.Cleanup(dropStage)
	}
	createStage()

	resourceId := resourcehelpers.EncodeResourceIdentifier(id)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Import WITHOUT experiment — booleans are imported as actual Snowflake values ("false"), causing a permadiff.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, stageModel),
				ResourceName:             stageModel.ResourceReference(),
				ImportState:              true,
				ImportStateId:            id.FullyQualifiedName(),
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.ignore_utf8_errors", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.trim_space", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.multi_line", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.enable_octal", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.allow_duplicate", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.strip_outer_array", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.strip_null_values", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.replace_invalid_characters", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.skip_byte_order_mark", "true"),
				),
				ImportStatePersist: true,
			},
			// Plan WITHOUT experiment — proves the bug: config has "default", state has "false" → permadiff
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, stageModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(stageModel.ResourceReference(), plancheck.ResourceActionUpdate),
						// planchecks.ExpectChange(stageModel.ResourceReference(), "file_format.0.json.0.ignore_utf8_errors", tfjson.ActionUpdate, sdk.String("false"), sdk.String("default")),
						planchecks.PrintPlanDetails(stageModel.ResourceReference(), "file_format", "directory", "use_privatelink_endpoint", "credentials"),
					},
				},
			},
			// Destroy to clear Terraform state before reimporting with the experiment enabled.
			// This also drops the stage, so we recreate it in the next step's PreConfig.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, stageModel),
				Destroy:                  true,
			},
			// Import WITH experiment — booleans are imported as "default" (matching schema defaults)
			{
				PreConfig:                createStage,
				ProtoV6ProviderFactories: experimentFactory,
				Config:                   accconfig.FromModels(t, providerModelWithExperiment, stageModel),
				ResourceName:             stageModel.ResourceReference(),
				ImportState:              true,
				ImportStatePersist:       true,
				ImportStateId:            id.FullyQualifiedName(),
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.ignore_utf8_errors", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.trim_space", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.multi_line", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.enable_octal", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.allow_duplicate", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.strip_outer_array", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.strip_null_values", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.replace_invalid_characters", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.skip_byte_order_mark", r.BooleanDefault),
				),
			},
			// Plan WITH experiment — proves the fix: config and state both have "default" -> no diff
			{
				ProtoV6ProviderFactories: experimentFactory,
				Config:                   accconfig.FromModels(t, providerModelWithExperiment, stageModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						planchecks.PrintPlanDetails(stageModel.ResourceReference(), "file_format", "directory", "use_privatelink_endpoint", "credentials"),
					},
				},
			},
			// Destroy to clear Terraform state before reimporting with the experiment enabled.
			// This also drops the stage, so we recreate it in the next step's PreConfig.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, stageModelPure),
				Destroy:                  true,
			},
			// Import WITH experiment — booleans are imported as "default" (matching schema defaults)
			{
				PreConfig:                createStage,
				ProtoV6ProviderFactories: experimentFactory,
				Config:                   accconfig.FromModels(t, providerModelWithExperiment, stageModelPure),
				ResourceName:             stageModelPure.ResourceReference(),
				ImportState:              true,
				ImportStatePersist:       true,
				ImportStateId:            id.FullyQualifiedName(),
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.ignore_utf8_errors", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.trim_space", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.multi_line", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.enable_octal", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.allow_duplicate", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.strip_outer_array", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.strip_null_values", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.replace_invalid_characters", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "file_format.0.json.0.skip_byte_order_mark", r.BooleanDefault),
				),
			},
			// Plan WITH experiment — proves the fix: config and state both have "default" -> no diff
			{
				ProtoV6ProviderFactories: experimentFactory,
				Config:                   accconfig.FromModels(t, providerModelWithExperiment, stageModelPure),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						planchecks.PrintPlanDetails(stageModel.ResourceReference(), "file_format", "directory", "use_privatelink_endpoint", "credentials"),
					},
				},
			},
		},
	})
}

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
// It covers all tri-value boolean attributes: JSON file format booleans, directory table booleans, and use_privatelink_endpoint.
// Regression test for https://github.com/snowflakedb/terraform-provider-snowflake/issues/4549.
func TestAcc_Experimental_ExternalS3Stage_ImportJsonBooleanDefaults(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	awsUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)

	providerModelWithExperiment := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportBooleanDefault)

	experimentFactory := providerFactoryUsingCache("TestAcc_Experimental_ExternalS3Stage_ImportJsonBooleanDefaults")

	// stageModel: directory block with enable=false but NO explicit auto_refresh,
	// so auto_refresh defaults to "default" in config. This lets us detect the permadiff
	// (import sets "false", config has "default") and verify the experiment fixes it.
	stageModel := model.ExternalS3StageWithId(id, awsUrl).
		WithDirectoryEnabledAndOptions(sdk.StageS3CommonDirectoryTableOptionsRequest{
			Enable: false,
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
	createStage := func() {
		_, dropStage := testClient().Stage.CreateStageOnS3WithRequest(t, sdk.NewCreateOnS3StageRequest(id, *sdk.NewExternalS3StageParamsRequest(awsUrl)).
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
			// Import WITHOUT experiment — booleans are imported as actual Snowflake values ("false"/"true"), causing a permadiff.
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
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "directory.0.auto_refresh", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "use_privatelink_endpoint", "false"),
				),
				ImportStatePersist: true,
			},
			// Plan WITHOUT experiment — proves the bug: config has "default", state has "false"/"true" → permadiff
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, stageModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(stageModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails(stageModel.ResourceReference(), "file_format", "directory", "use_privatelink_endpoint", "credentials"),
					},
				},
			},
			// Destroy to clear Terraform state before reimporting with the experiment enabled.
			// Unfortunately, one can't import a resource when it's already in the state,
			// and the framework doesn't support state removal.
			// This also drops the stage, so we recreate it in the next step's PreConfig.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, stageModel),
				Destroy:                  true,
			},
			// Import WITH experiment — all tri-value booleans are imported as "default"
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
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "directory.0.auto_refresh", r.BooleanDefault),
					importchecks.TestCheckResourceAttrInstanceState(resourceId, "use_privatelink_endpoint", r.BooleanDefault),
				),
			},
			// Plan WITH experiment — proves the fix: config and state both have "default" → no diff
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
		},
	})
}

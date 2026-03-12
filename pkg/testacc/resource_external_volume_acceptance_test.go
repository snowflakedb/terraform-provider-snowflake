//go:build non_account_level_tests

package testacc

// TODO Add test that includes Iceberg table creation, as this impacts the describe output (updates ACTIVE)

import (
	"fmt"
	"regexp"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// Note that generators currently don't handle lists of objects, which is required for storage locations
// Using the old approach of files for this reason

func getS3StorageLocation(
	locName string,
	provider string,
	baseUrl string,
	roleArn string,
	encryptionType string,
	s3EncryptionKmsKeyId string,
) config.Variable {
	m := map[string]config.Variable{
		"storage_location_name": config.StringVariable(locName),
		"storage_provider":      config.StringVariable(provider),
		"storage_base_url":      config.StringVariable(baseUrl),
		"storage_aws_role_arn":  config.StringVariable(roleArn),
		"encryption_type":       config.StringVariable(encryptionType),
	}
	if encryptionType == "AWS_SSE_KMS" {
		m["encryption_kms_key_id"] = config.StringVariable(s3EncryptionKmsKeyId)
	}
	return config.MapVariable(m)
}

func getS3StorageLocationWithExtras(
	locName string,
	provider string,
	baseUrl string,
	roleArn string,
	encryptionType string,
	s3EncryptionKmsKeyId string,
	accessPointArn string,
	usePrivatelinkEndpoint string,
	storageAwsExternalId string,
) config.Variable {
	m := map[string]config.Variable{
		"storage_location_name": config.StringVariable(locName),
		"storage_provider":      config.StringVariable(provider),
		"storage_base_url":      config.StringVariable(baseUrl),
		"storage_aws_role_arn":  config.StringVariable(roleArn),
		"encryption_type":       config.StringVariable(encryptionType),
	}
	if encryptionType == "AWS_SSE_KMS" {
		m["encryption_kms_key_id"] = config.StringVariable(s3EncryptionKmsKeyId)
	}
	if len(accessPointArn) > 0 {
		m["storage_aws_access_point_arn"] = config.StringVariable(accessPointArn)
	}
	if len(usePrivatelinkEndpoint) > 0 {
		m["use_privatelink_endpoint"] = config.StringVariable(usePrivatelinkEndpoint)
	}
	if len(storageAwsExternalId) > 0 {
		m["storage_aws_external_id"] = config.StringVariable(storageAwsExternalId)
	}
	return config.MapVariable(m)
}

func getS3CompatStorageLocation(
	locName string,
	baseUrl string,
	endpoint string,
	awsKeyId string,
	awsSecretKey string,
) config.Variable {
	return config.MapVariable(map[string]config.Variable{
		"storage_location_name":  config.StringVariable(locName),
		"storage_provider":       config.StringVariable("S3COMPAT"),
		"storage_base_url":       config.StringVariable(baseUrl),
		"storage_endpoint":       config.StringVariable(endpoint),
		"storage_aws_key_id":     config.StringVariable(awsKeyId),
		"storage_aws_secret_key": config.StringVariable(awsSecretKey),
	})
}

func getGcsStorageLocation(
	locName string,
	baseUrl string,
	encryptionType string,
	gcsEncryptionKmsKeyId string,
) config.Variable {
	gcsStorageProvider := "GCS"
	if encryptionType == "GCS_SSE_KMS" {
		return config.MapVariable(map[string]config.Variable{
			"storage_location_name": config.StringVariable(locName),
			"storage_provider":      config.StringVariable(gcsStorageProvider),
			"storage_base_url":      config.StringVariable(baseUrl),
			"encryption_type":       config.StringVariable(encryptionType),
			"encryption_kms_key_id": config.StringVariable(gcsEncryptionKmsKeyId),
		})
	} else {
		return config.MapVariable(map[string]config.Variable{
			"storage_location_name": config.StringVariable(locName),
			"storage_provider":      config.StringVariable(gcsStorageProvider),
			"storage_base_url":      config.StringVariable(baseUrl),
			"encryption_type":       config.StringVariable(encryptionType),
		})
	}
}

func getAzureStorageLocation(
	locName string,
	baseUrl string,
	azureTenantId string,
) config.Variable {
	azureStorageProvider := "AZURE"
	return config.MapVariable(map[string]config.Variable{
		"storage_location_name": config.StringVariable(locName),
		"storage_provider":      config.StringVariable(azureStorageProvider),
		"storage_base_url":      config.StringVariable(baseUrl),
		"azure_tenant_id":       config.StringVariable(azureTenantId),
	})
}

func externalVolume(storageLocations config.Variable, name string, comment string, allowWrites string) config.Variables {
	return config.Variables{
		"name":             config.StringVariable(name),
		"comment":          config.StringVariable(comment),
		"allow_writes":     config.StringVariable(allowWrites),
		"storage_location": storageLocations,
	}
}

func externalVolumeMultiple(s3StorageLocations config.Variable, gcsStorageLocations config.Variable, azureStorageLocations config.Variable, name string, comment string, allowWrites string) config.Variables {
	return config.Variables{
		"name":                    config.StringVariable(name),
		"comment":                 config.StringVariable(comment),
		"allow_writes":            config.StringVariable(allowWrites),
		"s3_storage_locations":    s3StorageLocations,
		"gcs_storage_locations":   gcsStorageLocations,
		"azure_storage_locations": azureStorageLocations,
	}
}

// Test volume with s3 storage locations
func TestAcc_ExternalVolume_BasicUseCase_S3(t *testing.T) {
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakePreProdGovEnvironment {
		t.Skip("Skipping test - Snowflake error 393962 (42601): External volumes in government deployments cannot use non-government S3 storage locations. Storage type S3 is not allowed. Please use S3GOV instead.")
	}
	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()

	ref := "snowflake_external_volume.complete"
	comment := random.Comment()
	comment2 := random.Comment()
	s3StorageLocationName := "s3Test"
	s3StorageProvider := "S3"
	s3StorageBaseUrl := "s3://my-example-bucket/"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionTypeNone := "NONE"
	s3EncryptionTypeSseKms := "AWS_SSE_KMS"
	s3EncryptionKmsKeyId := "123456789"
	s3EncryptionKmsKeyId2 := "987654321"
	awsExternalId := "123456789"
	awsExternalId2 := "987654321"
	s3AccessPointArn := "arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"
	s3AccessPointArnUpdated := "arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point-updated"
	s3StorageLocation := getS3StorageLocation(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeNone, "")
	s3StorageLocationComplete := getS3StorageLocationWithExtras(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeSseKms, s3EncryptionKmsKeyId, s3AccessPointArn, "true", awsExternalId)
	s3StorageLocationCompleteUpdated := getS3StorageLocationWithExtras(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeSseKms, s3EncryptionKmsKeyId2, s3AccessPointArnUpdated, "true", awsExternalId2)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			// create with a basic storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentEmpty().
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasCommentEmpty().
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasCommentEmpty().
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
						}),
				),
			},
			// import
			{
				ConfigDirectory:         ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables:         externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, "", ""),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"allow_writes", "storage_location.0.storage_aws_external_id"},
			},
			// update the location to have all optional fields
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocationComplete), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(1).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							awsExternalId,
							s3EncryptionTypeSseKms,
							s3EncryptionKmsKeyId,
							s3AccessPointArn,
							"true",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeSseKms,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn:        s3StorageAwsRoleArn,
									EncryptionKmsKeyId:       s3EncryptionKmsKeyId,
									StorageAwsAccessPointArn: s3AccessPointArn,
									UsePrivatelinkEndpoint:   sdk.Bool(true),
									StorageAwsExternalId:     awsExternalId,
								},
							},
						}),
				),
			},
			// import
			{
				ConfigDirectory:         ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:         externalVolume(config.ListVariable(s3StorageLocationComplete), externalVolumeName, comment, "true"),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_location.0.use_privatelink_endpoint", "storage_location.0.storage_aws_external_id"},
			},
			// update the location to have changed fields
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocationCompleteUpdated), externalVolumeName, comment2, "false"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString("false").
						HasStorageLocationLength(1).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							awsExternalId2,
							s3EncryptionTypeSseKms,
							s3EncryptionKmsKeyId2,
							s3AccessPointArnUpdated,
							"true",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(false),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasComment(comment2).
						HasAllowWrites("false").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeSseKms,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn:        s3StorageAwsRoleArn,
									EncryptionKmsKeyId:       s3EncryptionKmsKeyId2,
									StorageAwsAccessPointArn: s3AccessPointArnUpdated,
									UsePrivatelinkEndpoint:   sdk.Bool(true),
									StorageAwsExternalId:     awsExternalId2,
								},
							},
						}),
				),
			},
			// verify external changes
			{
				PreConfig: func() {
					testClient().ExternalVolume.Alter(t, sdk.NewAlterExternalVolumeRequest(id).WithSet(
						*sdk.NewAlterExternalVolumeSetRequest().WithComment("external comment"),
					))
					testClient().ExternalVolume.Alter(t, sdk.NewAlterExternalVolumeRequest(id).WithSet(
						*sdk.NewAlterExternalVolumeSetRequest().WithAllowWrites(true),
					))
					testClient().ExternalVolume.Alter(
						t,
						sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(
							*sdk.NewExternalVolumeStorageLocationItemRequest(
								*sdk.NewExternalVolumeStorageLocationRequest("externally-added-s3-storage-location").WithS3StorageLocationParams(
									*sdk.NewS3StorageLocationParamsRequest(
										"s3",
										"arn:aws:iam::123456789012:role/externally-added-role",
										"s3://externally-added-bucket",
									),
								),
							),
						),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocationCompleteUpdated), externalVolumeName, comment2, "false"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasCommentString(comment2).
						HasAllowWritesString("false").
						HasStorageLocationLength(1).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							awsExternalId2,
							s3EncryptionTypeSseKms,
							s3EncryptionKmsKeyId2,
							s3AccessPointArnUpdated,
							"true",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasComment(comment2).
						HasAllowWrites(false),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasComment(comment2).
						HasAllowWrites("false").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeSseKms,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn:        s3StorageAwsRoleArn,
									EncryptionKmsKeyId:       s3EncryptionKmsKeyId2,
									StorageAwsAccessPointArn: s3AccessPointArnUpdated,
									UsePrivatelinkEndpoint:   sdk.Bool(true),
									StorageAwsExternalId:     awsExternalId2,
								},
							},
						}),
				),
			},
			// unset the optional parameters
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentEmpty().
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasCommentEmpty().
						HasAllowWrites(false),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasCommentEmpty().
						HasAllowWrites("false").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
						}),
				),
			},
		},
	})
}

// Test volume with s3gov storage locations. It's a very simple smoke test, as the main functionalities are tested in the S3 test.
func TestAcc_ExternalVolume_BasicUseCase_S3Gov(t *testing.T) {
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() != testenvs.SnowflakePreProdGovEnvironment {
		t.Skip("Skipping S3Gov test, requires gov deployment")
	}

	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()
	ref := "snowflake_external_volume.complete"
	s3GovStorageLocationName := "s3GovTest"
	s3GovStorageProvider := "S3GOV"
	s3GovStorageBaseUrl := "s3gov://my-example-bucket"
	s3GovStorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3GovEncryptionTypeNone := "NONE"
	s3GovStorageLocation := getS3StorageLocation(s3GovStorageLocationName, s3GovStorageProvider, s3GovStorageBaseUrl, s3GovStorageAwsRoleArn, s3GovEncryptionTypeNone, "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			// create with a basic s3gov storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(s3GovStorageLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentEmpty().
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasS3StorageLocationAtIndex(
							0,
							s3GovStorageProvider,
							s3GovStorageLocationName,
							s3GovStorageBaseUrl,
							s3GovStorageAwsRoleArn,
							"",
							s3GovEncryptionTypeNone,
							"",
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasCommentEmpty().
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasCommentEmpty().
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3GovStorageLocationName,
								StorageProvider: s3GovStorageProvider,
								StorageBaseUrl:  s3GovStorageBaseUrl,
								EncryptionType:  s3GovEncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3GovStorageAwsRoleArn,
								},
							},
						}),
				),
			},
		},
	})
}

func TestAcc_ExternalVolume_BasicUseCase_S3Compat(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()

	ref := "snowflake_external_volume.complete"
	comment := random.Comment()
	comment2 := random.Comment()
	s3CompatLocationName := "s3CompatTest"
	s3CompatBaseUrl := "s3compat://my-example-bucket/"
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"
	s3CompatAwsKeyId := "AKIAIOSFODNN7EXAMPLE"
	s3CompatAwsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatLocation := getS3CompatStorageLocation(s3CompatLocationName, s3CompatBaseUrl, s3CompatEndpoint, s3CompatAwsKeyId, s3CompatAwsSecretKey)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			// create with a basic S3-compatible storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(s3CompatLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentEmpty().
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasS3CompatStorageLocationAtIndex(
							0,
							s3CompatLocationName,
							s3CompatBaseUrl,
							s3CompatEndpoint,
							s3CompatAwsKeyId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasCommentEmpty().
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasCommentEmpty().
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3CompatLocationName,
								StorageProvider: "S3COMPAT",
								StorageBaseUrl:  s3CompatBaseUrl,
								EncryptionType:  string(sdk.S3EncryptionNone),
								S3CompatStorageLocation: &sdk.StorageLocationS3CompatDetails{
									Endpoint:       s3CompatEndpoint,
									AwsAccessKeyId: s3CompatAwsKeyId,
								},
							},
						}),
				),
			},
			// import
			{
				ConfigDirectory:         ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables:         externalVolume(config.ListVariable(s3CompatLocation), externalVolumeName, "", ""),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"allow_writes", "storage_location.0.storage_aws_secret_key", "storage_location.0.use_privatelink_endpoint"},
			},
			// update to have all optional fields
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3CompatLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(1).
						HasS3CompatStorageLocationAtIndex(
							0,
							s3CompatLocationName,
							s3CompatBaseUrl,
							s3CompatEndpoint,
							s3CompatAwsKeyId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3CompatLocationName,
								StorageProvider: "S3COMPAT",
								StorageBaseUrl:  s3CompatBaseUrl,
								EncryptionType:  string(sdk.S3EncryptionNone),
								S3CompatStorageLocation: &sdk.StorageLocationS3CompatDetails{
									Endpoint:       s3CompatEndpoint,
									AwsAccessKeyId: s3CompatAwsKeyId,
								},
							},
						}),
				),
			},
			// import
			{
				ConfigDirectory:         ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:         externalVolume(config.ListVariable(s3CompatLocation), externalVolumeName, comment, "true"),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_location.0.storage_aws_secret_key", "storage_location.0.use_privatelink_endpoint"},
			},
			// update with changed fields
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3CompatLocation), externalVolumeName, comment2, "false"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString("false").
						HasStorageLocationLength(1).
						HasS3CompatStorageLocationAtIndex(
							0,
							s3CompatLocationName,
							s3CompatBaseUrl,
							s3CompatEndpoint,
							s3CompatAwsKeyId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(false),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasComment(comment2).
						HasAllowWrites("false").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3CompatLocationName,
								StorageProvider: "S3COMPAT",
								StorageBaseUrl:  s3CompatBaseUrl,
								EncryptionType:  string(sdk.S3EncryptionNone),
								S3CompatStorageLocation: &sdk.StorageLocationS3CompatDetails{
									Endpoint:       s3CompatEndpoint,
									AwsAccessKeyId: s3CompatAwsKeyId,
								},
							},
						}),
				),
			},
			// verify external changes
			{
				PreConfig: func() {
					testClient().ExternalVolume.Alter(t, sdk.NewAlterExternalVolumeRequest(id).WithSet(
						*sdk.NewAlterExternalVolumeSetRequest().WithAllowWrites(true).WithComment("external comment"),
					))
					testClient().ExternalVolume.Alter(
						t,
						sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(
							*sdk.NewExternalVolumeStorageLocationItemRequest(
								*sdk.NewExternalVolumeStorageLocationRequest("externally-added-s3compat-storage-location").WithS3CompatStorageLocationParams(
									*sdk.NewS3CompatStorageLocationParamsRequest(
										"s3compat://externally-added-bucket",
										"s3.us-east-2.amazonaws.com",
										*sdk.NewExternalVolumeS3CompatCredentialsRequest(
											"AKIAIOSFODNN7EXTERN",
											"externalSecretKey123456789012345678901",
										),
									),
								),
							),
						),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3CompatLocation), externalVolumeName, comment2, "false"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasCommentString(comment2).
						HasAllowWritesString("false").
						HasStorageLocationLength(1).
						HasS3CompatStorageLocationAtIndex(
							0,
							s3CompatLocationName,
							s3CompatBaseUrl,
							s3CompatEndpoint,
							s3CompatAwsKeyId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasComment(comment2).
						HasAllowWrites(false),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasComment(comment2).
						HasAllowWrites("false").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3CompatLocationName,
								StorageProvider: "S3COMPAT",
								StorageBaseUrl:  s3CompatBaseUrl,
								EncryptionType:  string(sdk.S3EncryptionNone),
								S3CompatStorageLocation: &sdk.StorageLocationS3CompatDetails{
									Endpoint:       s3CompatEndpoint,
									AwsAccessKeyId: s3CompatAwsKeyId,
								},
							},
						}),
				),
			},
			// unset the optional parameters
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(s3CompatLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentEmpty().
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasS3CompatStorageLocationAtIndex(
							0,
							s3CompatLocationName,
							s3CompatBaseUrl,
							s3CompatEndpoint,
							s3CompatAwsKeyId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasCommentEmpty().
						HasAllowWrites(false),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasCommentEmpty().
						HasAllowWrites("false").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3CompatLocationName,
								StorageProvider: "S3COMPAT",
								StorageBaseUrl:  s3CompatBaseUrl,
								EncryptionType:  string(sdk.S3EncryptionNone),
								S3CompatStorageLocation: &sdk.StorageLocationS3CompatDetails{
									Endpoint:       s3CompatEndpoint,
									AwsAccessKeyId: s3CompatAwsKeyId,
								},
							},
						}),
				),
			},
		},
	})
}

// Test volume with gcs storage locations
func TestAcc_ExternalVolume_BasicUseCase_GCS(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()

	ref := "snowflake_external_volume.complete"
	comment := random.Comment()
	comment2 := random.Comment()
	gcsStorageLocationName := "gcsTest"
	gcsStorageProvider := "GCS"
	gcsStorageBaseUrl := "gcs://my-example-bucket/"
	gcsEncryptionTypeNone := "NONE"
	gcsEncryptionTypeSseKms := "GCS_SSE_KMS"
	gcsEncryptionKmsKeyId := "123456789"
	gcsEncryptionKmsKeyId2 := "987654321"
	gcsStorageLocation := getGcsStorageLocation(gcsStorageLocationName, gcsStorageBaseUrl, gcsEncryptionTypeNone, "")
	gcsStorageLocationKmsEncryption := getGcsStorageLocation(gcsStorageLocationName, gcsStorageBaseUrl, gcsEncryptionTypeSseKms, gcsEncryptionKmsKeyId)
	gcsStorageLocationKmsEncryptionUpdatedKey := getGcsStorageLocation(gcsStorageLocationName, gcsStorageBaseUrl, gcsEncryptionTypeSseKms, gcsEncryptionKmsKeyId2)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			// create with a basic storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentEmpty().
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasGCSStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasCommentEmpty().
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasCommentEmpty().
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
						}),
				),
			},
			// import
			{
				ConfigDirectory:         ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables:         externalVolume(config.ListVariable(gcsStorageLocation), externalVolumeName, "", ""),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"allow_writes"},
			},
			// update the location to have all optional fields
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocationKmsEncryption), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(1).
						HasGCSStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeSseKms,
							gcsEncryptionKmsKeyId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeSseKms,
								GCSStorageLocation: &sdk.StorageLocationGcsDetails{
									EncryptionKmsKeyId: gcsEncryptionKmsKeyId,
								},
							},
						}),
				),
			},
			// import
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:   externalVolume(config.ListVariable(gcsStorageLocationKmsEncryption), externalVolumeName, comment, "true"),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// update the location to have changed fields
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocationKmsEncryptionUpdatedKey), externalVolumeName, comment2, "false"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString("false").
						HasStorageLocationLength(1).
						HasGCSStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeSseKms,
							gcsEncryptionKmsKeyId2,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(false),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasComment(comment2).
						HasAllowWrites("false").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeSseKms,
								GCSStorageLocation: &sdk.StorageLocationGcsDetails{
									EncryptionKmsKeyId: gcsEncryptionKmsKeyId2,
								},
							},
						}),
				),
			},
			// verify external changes
			{
				PreConfig: func() {
					testClient().ExternalVolume.Alter(t, sdk.NewAlterExternalVolumeRequest(id).WithSet(
						*sdk.NewAlterExternalVolumeSetRequest().WithAllowWrites(true).WithComment("external comment"),
					))
					testClient().ExternalVolume.Alter(
						t,
						sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(
							*sdk.NewExternalVolumeStorageLocationItemRequest(
								*sdk.NewExternalVolumeStorageLocationRequest("externally-added-gcs-storage-location").WithGCSStorageLocationParams(
									*sdk.NewGCSStorageLocationParamsRequest(
										"gcs://externally-added-bucket",
									),
								),
							),
						),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocationKmsEncryptionUpdatedKey), externalVolumeName, comment2, "false"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasCommentString(comment2).
						HasAllowWritesString("false").
						HasStorageLocationLength(1).
						HasGCSStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeSseKms,
							gcsEncryptionKmsKeyId2,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasComment(comment2).
						HasAllowWrites(false),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasComment(comment2).
						HasAllowWrites("false").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeSseKms,
								GCSStorageLocation: &sdk.StorageLocationGcsDetails{
									EncryptionKmsKeyId: gcsEncryptionKmsKeyId2,
								},
							},
						}),
				),
			},
			// unset the optional parameters
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentEmpty().
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasGCSStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasCommentEmpty().
						HasAllowWrites(false),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasCommentEmpty().
						HasAllowWrites("false").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
						}),
				),
			},
		},
	})
}

// Test volume with azure storage locations
func TestAcc_ExternalVolume_BasicUseCase_Azure(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()

	ref := "snowflake_external_volume.complete"
	comment := random.Comment()
	comment2 := random.Comment()
	azureStorageLocationName := "azureTest"
	azureStorageProvider := "AZURE"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container/"
	azureTenantId := "123456789"
	azureEncryptionTypeNone := "NONE"
	azureStorageLocation := getAzureStorageLocation(azureStorageLocationName, azureStorageBaseUrl, azureTenantId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			// create with a basic storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentEmpty().
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasAzureStorageLocationAtIndex(
							0,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasCommentEmpty().
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasCommentEmpty().
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// import
			{
				ConfigDirectory:         ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables:         externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, "", ""),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"allow_writes"},
			},
			// update the location to have all optional fields
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(1).
						HasAzureStorageLocationAtIndex(
							0,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// import
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:   externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// update the location to have changed fields
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, comment2, "false"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString("false").
						HasStorageLocationLength(1).
						HasAzureStorageLocationAtIndex(
							0,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(false),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasComment(comment2).
						HasAllowWrites("false").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// verify external changes
			{
				PreConfig: func() {
					testClient().ExternalVolume.Alter(t, sdk.NewAlterExternalVolumeRequest(id).WithSet(
						*sdk.NewAlterExternalVolumeSetRequest().WithComment("external comment"),
					))
					testClient().ExternalVolume.Alter(t, sdk.NewAlterExternalVolumeRequest(id).WithSet(
						*sdk.NewAlterExternalVolumeSetRequest().WithAllowWrites(true),
					))
					testClient().ExternalVolume.Alter(
						t,
						sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(
							*sdk.NewExternalVolumeStorageLocationItemRequest(
								*sdk.NewExternalVolumeStorageLocationRequest("externally-added-azure-storage-location").WithAzureStorageLocationParams(
									*sdk.NewAzureStorageLocationParamsRequest(
										azureTenantId,
										"azure://123456789.blob.core.windows.net/externally_added_container",
									),
								),
							),
						),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, comment2, "false"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasCommentString(comment2).
						HasAllowWritesString("false").
						HasStorageLocationLength(1).
						HasAzureStorageLocationAtIndex(
							0,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasComment(comment2).
						HasAllowWrites(false),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasComment(comment2).
						HasAllowWrites("false").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// unset the optional parameters
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, ref).
						HasNameString(externalVolumeName).
						HasCommentEmpty().
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasAzureStorageLocationAtIndex(
							0,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, ref).
						HasName(externalVolumeName).
						HasCommentEmpty().
						HasAllowWrites(false),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, ref).
						HasActiveEmpty().
						HasCommentEmpty().
						HasAllowWrites("false").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
		},
	})
}

// Test apply works when setting all optionals from the start
// Other tests start without setting all optionals
func TestAcc_ExternalVolume_CompleteUseCase(t *testing.T) {
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakePreProdGovEnvironment {
		t.Skip("Skipping test - Snowflake error 393962 (42601): External volumes in government deployments cannot use non-government S3 storage locations. Storage type S3 is not allowed. Please use S3GOV instead.")
	}
	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()
	comment := random.Comment()

	s3StorageLocationName := "s3Test"
	s3StorageProvider := "S3"
	s3StorageBaseUrl := "s3://my-example-bucket/"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionTypeSseKms := "AWS_SSE_KMS"
	s3EncryptionKmsKeyId := "123456789"
	awsExternalId := "123456789"
	s3StorageLocationKmsEncryption := getS3StorageLocationWithExtras(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeSseKms, s3EncryptionKmsKeyId, "", "", awsExternalId)

	gcsStorageLocationName := "gcsTest"
	gcsStorageProvider := "GCS"
	gcsStorageBaseUrl := "gcs://my-example-bucket/"
	gcsEncryptionTypeSseKms := "GCS_SSE_KMS"
	gcsEncryptionKmsKeyId := "123456789"
	gcsStorageLocationKmsEncryption := getGcsStorageLocation(gcsStorageLocationName, gcsStorageBaseUrl, gcsEncryptionTypeSseKms, gcsEncryptionKmsKeyId)

	azureStorageLocationName := "azureTest"
	azureStorageProvider := "AZURE"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container/"
	azureTenantId := "123456789"
	azureEncryptionTypeNone := "NONE"
	azureStorageLocation := getAzureStorageLocation(azureStorageLocationName, azureStorageBaseUrl, azureTenantId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocationKmsEncryption), config.ListVariable(gcsStorageLocationKmsEncryption), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "false"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("false").
						HasStorageLocationLength(3).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							awsExternalId,
							s3EncryptionTypeSseKms,
							s3EncryptionKmsKeyId,
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeSseKms,
							gcsEncryptionKmsKeyId,
						).
						HasAzureStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(false),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("false").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeSseKms,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn:    s3StorageAwsRoleArn,
									StorageAwsExternalId: awsExternalId,
									EncryptionKmsKeyId:   s3EncryptionKmsKeyId,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeSseKms,
								GCSStorageLocation: &sdk.StorageLocationGcsDetails{
									EncryptionKmsKeyId: gcsEncryptionKmsKeyId,
								},
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			{
				ConfigDirectory:         ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables:         externalVolumeMultiple(config.ListVariable(s3StorageLocationKmsEncryption), config.ListVariable(gcsStorageLocationKmsEncryption), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "false"),
				ResourceName:            "snowflake_external_volume.complete",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_location.0.storage_aws_external_id"},
			},
		},
	})
}

// Test volume with multiple storage locations that span multiple providers
// Test adding/removing storage locations at different positions in the storage_location list
func TestAcc_ExternalVolume_BasicUseCase_MultipleProviders(t *testing.T) {
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakePreProdGovEnvironment {
		t.Skip("Skipping test - Snowflake error 393962 (42601): External volumes in government deployments cannot use non-government S3 storage locations. Storage type S3 is not allowed. Please use S3GOV instead.")
	}
	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()
	comment := random.Comment()
	s3StorageLocationName := "s3Test"
	s3StorageLocationName2 := "s3Test2"
	s3StorageProvider := "S3"
	s3StorageBaseUrl := "s3://my-example-bucket/"
	s3StorageBaseUrl2 := "s3://my-example-bucket2/"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionTypeNone := "NONE"
	s3StorageLocation := getS3StorageLocation(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeNone, "")
	s3StorageLocationUpdatedBaseUrl := getS3StorageLocation(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl2, s3StorageAwsRoleArn, s3EncryptionTypeNone, "")
	s3StorageLocationUpdatedName := getS3StorageLocation(s3StorageLocationName2, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeNone, "")

	gcsStorageLocationName := "gcsTest"
	gcsStorageLocationName2 := "gcsTest2"
	gcsStorageProvider := "GCS"
	gcsStorageBaseUrl := "gcs://my-example-bucket/"
	gcsStorageBaseUrl2 := "gcs://my-example-bucket2/"
	gcsEncryptionTypeNone := "NONE"
	gcsStorageLocation := getGcsStorageLocation(gcsStorageLocationName, gcsStorageBaseUrl, gcsEncryptionTypeNone, "")
	gcsStorageLocationUpdatedName := getGcsStorageLocation(gcsStorageLocationName2, gcsStorageBaseUrl, gcsEncryptionTypeNone, "")
	gcsStorageLocationUpdatedBaseUrl := getGcsStorageLocation(gcsStorageLocationName, gcsStorageBaseUrl2, gcsEncryptionTypeNone, "")

	azureStorageLocationName := "azureTest"
	azureStorageLocationName2 := "azureTest2"
	azureStorageProvider := "AZURE"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container/"
	azureTenantId := "123456789"
	azureTenantId2 := "987654321"
	azureEncryptionTypeNone := "NONE"
	azureStorageLocation := getAzureStorageLocation(azureStorageLocationName, azureStorageBaseUrl, azureTenantId)
	azureStorageLocationUpdatedTenantId := getAzureStorageLocation(azureStorageLocationName, azureStorageBaseUrl, azureTenantId2)
	azureStorageLocationUpdatedName := getAzureStorageLocation(azureStorageLocationName2, azureStorageBaseUrl, azureTenantId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			// one location of each provider
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(3).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// import
			{
				ConfigDirectory:         ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables:         externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				ResourceName:            "snowflake_external_volume.complete",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_location.0.storage_aws_external_id"},
			},
			// change the s3 base url at position 0
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocationUpdatedBaseUrl), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(3).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl2,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl2,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// change back the s3 base url at position 0
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(3).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// add new s3 storage location to position 0
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocationUpdatedName, s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(4).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName2,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasS3StorageLocationAtIndex(
							1,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							2,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							3,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName2,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// remove s3 storage location at position 0
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(3).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// change the base url of the gcs storage location at position 1
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocationUpdatedBaseUrl), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(3).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl2,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl2,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// change back the encryption type of the gcs storage location at position 1
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(3).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// add new s3 storage location to position 1
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation, s3StorageLocationUpdatedName), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(4).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasS3StorageLocationAtIndex(
							1,
							s3StorageProvider,
							s3StorageLocationName2,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							2,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							3,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            s3StorageLocationName2,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// remove s3 storage location at position 1
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(3).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// change the tenant id of the azure storage location at position 2
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocationUpdatedTenantId), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(3).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId2,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId2,
								},
							},
						}),
				),
			},
			// change back the tenant id of the azure storage location at position 2
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(3).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// add new gcs storage location to position 2
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation, gcsStorageLocationUpdatedName), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(4).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasGCSStorageLocationAtIndex(
							2,
							gcsStorageLocationName2,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							3,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            gcsStorageLocationName2,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// remove gcs storage location at position 2
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(3).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// add new azure storage location to position 3
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation, azureStorageLocationUpdatedName), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(4).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						).
						HasAzureStorageLocationAtIndex(
							3,
							azureStorageLocationName2,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
							{
								Name:            azureStorageLocationName2,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
			// remove azure storage location from position 3
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, "true"),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString("true").
						HasStorageLocationLength(3).
						HasS3StorageLocationAtIndex(
							0,
							s3StorageProvider,
							s3StorageLocationName,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							"",
							s3EncryptionTypeNone,
							"",
							"",
							"",
						).
						HasGCSStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageBaseUrl,
							gcsEncryptionTypeNone,
							"",
						).
						HasAzureStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageBaseUrl,
							azureEncryptionTypeNone,
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(true),
					resourceshowoutputassert.ExternalVolumeDescribeOutput(t, "snowflake_external_volume.complete").
						HasComment(comment).
						HasAllowWrites("true").
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3StorageLocationName,
								StorageProvider: s3StorageProvider,
								StorageBaseUrl:  s3StorageBaseUrl,
								EncryptionType:  s3EncryptionTypeNone,
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3StorageAwsRoleArn,
								},
							},
							{
								Name:            gcsStorageLocationName,
								StorageProvider: gcsStorageProvider,
								StorageBaseUrl:  gcsStorageBaseUrl,
								EncryptionType:  gcsEncryptionTypeNone,
							},
							{
								Name:            azureStorageLocationName,
								StorageProvider: azureStorageProvider,
								StorageBaseUrl:  azureStorageBaseUrl,
								EncryptionType:  azureEncryptionTypeNone,
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
						}),
				),
			},
		},
	})
}

// Test invalid parameter combinations throw errors
func TestAcc_ExternalVolume_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	s3StorageLocationName := "s3Test"
	s3StorageProvider := "S3"
	s3StorageBaseUrl := "s3://my-example-bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionTypeNone := "NONE"
	s3EncryptionKmsKeyId := "123456789"

	gcsStorageLocationName := "gcsTest"
	gcsStorageProvider := "GCS"
	gcsStorageBaseUrl := "gcs://my-example-bucket"

	azureStorageLocationName := "azureTest"
	azureStorageProvider := "AZURE"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	externalVolumeName := id.Name()
	s3StorageLocationInvalidStorageProvider := getS3StorageLocation(s3StorageLocationName, "invalid-storage-provider", s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeNone, "")
	s3StorageLocationNoRoleArn := config.MapVariable(map[string]config.Variable{
		"storage_location_name": config.StringVariable(s3StorageLocationName),
		"storage_provider":      config.StringVariable(s3StorageProvider),
		"storage_base_url":      config.StringVariable(s3StorageBaseUrl),
	})
	s3StorageLocationWithTenantId := config.MapVariable(map[string]config.Variable{
		"storage_location_name": config.StringVariable(s3StorageLocationName),
		"storage_provider":      config.StringVariable(s3StorageProvider),
		"storage_base_url":      config.StringVariable(s3StorageBaseUrl),
		"azure_tenant_id":       config.StringVariable(azureTenantId),
	})
	gcsStorageLocationWithAwsRoleArn := config.MapVariable(map[string]config.Variable{
		"storage_location_name": config.StringVariable(gcsStorageLocationName),
		"storage_provider":      config.StringVariable(gcsStorageProvider),
		"storage_base_url":      config.StringVariable(gcsStorageBaseUrl),
		"storage_aws_role_arn":  config.StringVariable(s3StorageAwsRoleArn),
	})
	gcsStorageLocationWithTenantId := config.MapVariable(map[string]config.Variable{
		"storage_location_name": config.StringVariable(gcsStorageLocationName),
		"storage_provider":      config.StringVariable(gcsStorageProvider),
		"storage_base_url":      config.StringVariable(gcsStorageBaseUrl),
		"azure_tenant_id":       config.StringVariable(azureTenantId),
	})
	azureStorageLocationNoTenantId := config.MapVariable(map[string]config.Variable{
		"storage_location_name": config.StringVariable(azureStorageLocationName),
		"storage_provider":      config.StringVariable(azureStorageProvider),
		"storage_base_url":      config.StringVariable(azureStorageBaseUrl),
	})
	azureStorageLocationWithKmsKeyId := config.MapVariable(map[string]config.Variable{
		"storage_location_name": config.StringVariable(azureStorageLocationName),
		"storage_provider":      config.StringVariable(azureStorageProvider),
		"storage_base_url":      config.StringVariable(azureStorageBaseUrl),
		"azure_tenant_id":       config.StringVariable(azureTenantId),
		"encryption_kms_key_id": config.StringVariable(s3EncryptionKmsKeyId),
	})
	azureStorageLocationWithAwsRoleArn := config.MapVariable(map[string]config.Variable{
		"storage_location_name": config.StringVariable(azureStorageLocationName),
		"storage_provider":      config.StringVariable(azureStorageProvider),
		"storage_base_url":      config.StringVariable(azureStorageBaseUrl),
		"storage_aws_role_arn":  config.StringVariable(s3StorageAwsRoleArn),
		"azure_tenant_id":       config.StringVariable(azureTenantId),
	})
	azureStorageLocationWithEncryptionType := config.MapVariable(map[string]config.Variable{
		"storage_location_name": config.StringVariable(azureStorageLocationName),
		"storage_provider":      config.StringVariable(azureStorageProvider),
		"storage_base_url":      config.StringVariable(azureStorageBaseUrl),
		"azure_tenant_id":       config.StringVariable(azureTenantId),
		"encryption_type":       config.StringVariable(string(sdk.GCSEncryptionTypeSseKms)),
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			// invalid storage provider test
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocationInvalidStorageProvider), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("invalid storage provider: INVALID-STORAGE-PROVIDER"),
			},
			// no storage locations test
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("At least 1 \"storage_location\" blocks are required"),
			},
			// aws storage location doesn't specify storage_aws_role_arn
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocationNoRoleArn), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("unable to extract storage location, storage_aws_role_arn is required for s3 storage location"),
			},
			// azure storage location doesn't specify azure_tenant_id
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocationNoTenantId), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("unable to extract storage location, missing azure_tenant_id provider key in an azure storage location"),
			},
			// azure_tenant_id specified for s3 storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocationWithTenantId), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("unable to extract storage location, azure_tenant_id is not supported for s3 storage location"),
			},
			// storage_aws_role_arn specified for gcs storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocationWithAwsRoleArn), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("unable to extract storage location, storage_aws_role_arn is not supported for gcs storage location"),
			},
			// azure_tenant_id specified for gcs storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocationWithTenantId), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("unable to extract storage location, azure_tenant_id is not supported for gcs storage location"),
			},
			// storage_aws_role_arn specified for azure storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocationWithAwsRoleArn), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("unable to extract storage location, storage_aws_role_arn is not supported for azure storage location"),
			},
			// encryption_kms_key_id specified for azure storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocationWithKmsKeyId), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("unable to extract storage location, encryption_kms_key_id is not supported for azure storage location"),
			},
			// encryption_type is not supported for azure storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocationWithEncryptionType), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("unable to extract storage location, encryption_type is not supported for azure storage location"),
			},
		},
	})
}

func TestAcc_ExternalVolume_migrateFromVersion_2_14_0(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	configWithProvider := fmt.Sprintf(`
provider "snowflake" {
  preview_features_enabled = ["%s"]
}

resource "snowflake_external_volume" "complete" {
  name = "%s"
  storage_location {
    storage_location_name = "s3Test"
    storage_provider      = "S3"
    storage_base_url      = "s3://my-example-bucket"
    storage_aws_role_arn  = "arn:aws:iam::123456789012:role/myrole"
    encryption_type       = "NONE"
  }
}
`, previewfeatures.ExternalVolumeResource, id.Name())

	configWithoutProvider := fmt.Sprintf(`
resource "snowflake_external_volume" "complete" {
  name = "%s"
  storage_location {
    storage_location_name = "s3Test"
    storage_provider      = "S3"
    storage_base_url      = "s3://my-example-bucket"
    storage_aws_role_arn  = "arn:aws:iam::123456789012:role/myrole"
    encryption_type       = "NONE"
  }
}
`, id.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.14.0"),
				Config:            configWithProvider,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_volume.complete", "id", helpers.EncodeResourceIdentifier(id)),
					resource.TestCheckResourceAttr("snowflake_external_volume.complete", "name", id.Name()),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   configWithoutProvider,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_external_volume.complete", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_volume.complete", "id", helpers.EncodeResourceIdentifier(id)),
					resource.TestCheckResourceAttr("snowflake_external_volume.complete", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.0.storage_locations.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_external_volume.complete", "describe_output.0.allow_writes"),
				),
			},
		},
	})
}

func TestAcc_ExternalVolume_migrateFromVersion_2_14_0_externalId(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()

	s3StorageLocation := getS3StorageLocation("s3Test", "S3", "s3://my-example-bucket", "arn:aws:iam::123456789012:role/myrole", "NONE", "")

	// In 2.14.0, storage_aws_external_id was Computed: SF populated it automatically.
	// Now it's Optional: the state upgrader clears the old SF-generated value so there's no drift.
	configWithProvider := fmt.Sprintf(`
provider "snowflake" {
  preview_features_enabled = ["%s"]
}

resource "snowflake_external_volume" "complete" {
  name = "%s"
  storage_location {
    storage_location_name = "s3Test"
    storage_provider      = "S3"
    storage_base_url      = "s3://my-example-bucket"
    storage_aws_role_arn  = "arn:aws:iam::123456789012:role/myrole"
    encryption_type       = "NONE"
  }
}
`, previewfeatures.ExternalVolumeResource, externalVolumeName)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.14.0"),
				Config:            configWithProvider,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_volume.complete", "id", helpers.EncodeResourceIdentifier(id)),
					resource.TestCheckResourceAttr("snowflake_external_volume.complete", "name", externalVolumeName),
					// In 2.14.0, storage_aws_external_id is Computed and populated by SF.
					resource.TestCheckResourceAttrSet("snowflake_external_volume.complete", "storage_location.0.storage_aws_external_id"),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigDirectory:          ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables:          externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, "", ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_external_volume.complete", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_volume.complete", "id", helpers.EncodeResourceIdentifier(id)),
					resource.TestCheckResourceAttr("snowflake_external_volume.complete", "name", externalVolumeName),
					// After the state upgrader, storage_aws_external_id should be cleared.
					resource.TestCheckResourceAttr("snowflake_external_volume.complete", "storage_location.0.storage_aws_external_id", ""),
				),
			},
		},
	})
}

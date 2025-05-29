//go:build !account_level_tests

package testacc

// TODO Add test that includes Iceberg table creation, as this impacts the describe output (updates ACTIVE)
// TODO Add S3Gov tests

import (
	"regexp"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	if encryptionType == "AWS_SSE_KMS" {
		return config.MapVariable(map[string]config.Variable{
			"storage_location_name": config.StringVariable(locName),
			"storage_provider":      config.StringVariable(provider),
			"storage_base_url":      config.StringVariable(baseUrl),
			"storage_aws_role_arn":  config.StringVariable(roleArn),
			"encryption_type":       config.StringVariable(encryptionType),
			"encryption_kms_key_id": config.StringVariable(s3EncryptionKmsKeyId),
		})
	} else {
		return config.MapVariable(map[string]config.Variable{
			"storage_location_name": config.StringVariable(locName),
			"storage_provider":      config.StringVariable(provider),
			"storage_base_url":      config.StringVariable(baseUrl),
			"storage_aws_role_arn":  config.StringVariable(roleArn),
			"encryption_type":       config.StringVariable(encryptionType),
		})
	}
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
func TestAcc_External_Volume_S3(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()
	resourceId := helpers.EncodeResourceIdentifier(id)
	comment := random.Comment()
	comment2 := random.Comment()
	allowWritesTrue := "true"
	allowWritesFalse := "false"
	s3StorageLocationName := "s3Test"
	s3StorageLocationName2 := "s3Test2"
	s3StorageProvider := "S3"
	s3StorageBaseUrl := "s3://my-example-bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionTypeNone := "NONE"
	s3EncryptionTypeSseS3 := "AWS_SSE_S3"
	s3EncryptionTypeSseKms := "AWS_SSE_KMS"
	s3EncryptionKmsKeyId := "123456789"
	s3EncryptionKmsKeyId2 := "987654321"
	s3StorageLocation := getS3StorageLocation(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeNone, "")
	s3StorageLocationUpdatedName := getS3StorageLocation(s3StorageLocationName2, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeNone, "")
	s3StorageLocationSseEncryption := getS3StorageLocation(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeSseS3, "")
	s3StorageLocationKmsEncryption := getS3StorageLocation(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeSseKms, s3EncryptionKmsKeyId)
	s3StorageLocationKmsEncryptionUpdatedName := getS3StorageLocation(s3StorageLocationName2, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeSseKms, s3EncryptionKmsKeyId)
	s3StorageLocationKmsEncryptionUpdatedKey := getS3StorageLocation(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeSseKms, s3EncryptionKmsKeyId2)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			// without optionals
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString("").
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment("").
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "3")),
				),
			},
			// import - without optionals
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, "", ""),
				ResourceName:    "snowflake_external_volume.test",
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedExternalVolumeResource(t, resourceId).
						HasNameString(externalVolumeName).
						HasStorageLocationLength(1),
				),
			},
			// set external volume optionals
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// import - with external volume optionals
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:   externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, comment, allowWritesTrue),
				ResourceName:      "snowflake_external_volume.complete",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add second storage location without s3 optionals
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation, s3StorageLocationUpdatedName), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(2).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							s3StorageLocationName2,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "5")),
				),
			},
			// import - 2 storage locations without all s3 optionals set
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:   externalVolume(config.ListVariable(s3StorageLocation, s3StorageLocationUpdatedName), externalVolumeName, comment, allowWritesTrue),
				ResourceName:      "snowflake_external_volume.complete",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// update comment and change back to 1 storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, comment2, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// update allowWrites
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// add none encryption
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// add AWS_SSE_S3 encryption
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocationSseEncryption), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeSseS3,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// add AWS_SSE_KMS encryption
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocationKmsEncryption), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeSseKms,
							s3EncryptionKmsKeyId,
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// update kms key
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocationKmsEncryptionUpdatedKey), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeSseKms,
							s3EncryptionKmsKeyId2,
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// import - all optional s3 storage location optionals set
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:   externalVolume(config.ListVariable(s3StorageLocationKmsEncryption), externalVolumeName, comment2, allowWritesFalse),
				ResourceName:      "snowflake_external_volume.complete",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add second storage location with all s3 optionals set
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocationKmsEncryptionUpdatedKey, s3StorageLocationKmsEncryptionUpdatedName), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(2).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeSseKms,
							s3EncryptionKmsKeyId2,
							"",
						).
						HasStorageLocationAtIndex(
							1,
							s3StorageLocationName2,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeSseKms,
							s3EncryptionKmsKeyId,
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "5")),
				),
			},
			// import - 2 s3 storage locations, with all s3 optionals set
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:   externalVolume(config.ListVariable(s3StorageLocationKmsEncryptionUpdatedKey, s3StorageLocationKmsEncryptionUpdatedName), externalVolumeName, comment2, allowWritesFalse),
				ResourceName:      "snowflake_external_volume.complete",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// change back to AWS_SSE_S3 encryption with 1 s3 storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocationSseEncryption), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeSseS3,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// change back to none encryption
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// change back to no encryption
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// remove allow writes and comment from config
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString("").
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment("").
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "3")),
				),
			},
		},
	})
}

// Test volume with gcs storage locations
func TestAcc_External_Volume_GCS(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()
	resourceId := helpers.EncodeResourceIdentifier(id)
	comment := random.Comment()
	comment2 := random.Comment()
	allowWritesTrue := "true"
	allowWritesFalse := "false"
	gcsStorageLocationName := "gcsTest"
	gcsStorageLocationName2 := "gcsTest2"
	gcsStorageProvider := "GCS"
	gcsStorageBaseUrl := "gcs://my-example-bucket"
	gcsEncryptionTypeNone := "NONE"
	gcsEncryptionTypeSseKms := "GCS_SSE_KMS"
	gcsEncryptionKmsKeyId := "123456789"
	gcsEncryptionKmsKeyId2 := "987654321"
	gcsStorageLocation := getGcsStorageLocation(gcsStorageLocationName, gcsStorageBaseUrl, gcsEncryptionTypeNone, "")
	gcsStorageLocationUpdatedName := getGcsStorageLocation(gcsStorageLocationName2, gcsStorageBaseUrl, gcsEncryptionTypeNone, "")
	gcsStorageLocationKmsEncryption := getGcsStorageLocation(gcsStorageLocationName, gcsStorageBaseUrl, gcsEncryptionTypeSseKms, gcsEncryptionKmsKeyId)
	gcsStorageLocationKmsEncryptionUpdatedName := getGcsStorageLocation(gcsStorageLocationName2, gcsStorageBaseUrl, gcsEncryptionTypeSseKms, gcsEncryptionKmsKeyId)
	gcsStorageLocationKmsEncryptionUpdatedKey := getGcsStorageLocation(gcsStorageLocationName, gcsStorageBaseUrl, gcsEncryptionTypeSseKms, gcsEncryptionKmsKeyId2)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			// without optionals
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString("").
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment("").
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "3")),
				),
			},
			// import - without optionals
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocation), externalVolumeName, "", ""),
				ResourceName:    "snowflake_external_volume.test",
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedExternalVolumeResource(t, resourceId).
						HasNameString(externalVolumeName).
						HasStorageLocationLength(1),
				),
			},
			// set external volume optionals
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// import - with external volume optionals
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:   externalVolume(config.ListVariable(gcsStorageLocation), externalVolumeName, comment, allowWritesTrue),
				ResourceName:      "snowflake_external_volume.complete",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add second storage location without gcs optionals
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocation, gcsStorageLocationUpdatedName), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(2).
						HasStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName2,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "5")),
				),
			},
			// import - 2 storage locations without all gcs optionals set
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:   externalVolume(config.ListVariable(gcsStorageLocation, gcsStorageLocationUpdatedName), externalVolumeName, comment, allowWritesTrue),
				ResourceName:      "snowflake_external_volume.complete",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// update comment and change back to 1 storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocation), externalVolumeName, comment2, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// update allowWrites
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocation), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// add none encryption
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocation), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// add GCS_SSE_KMS encryption
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocationKmsEncryption), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeSseKms,
							gcsEncryptionKmsKeyId,
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// update kms key
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocationKmsEncryptionUpdatedKey), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeSseKms,
							gcsEncryptionKmsKeyId2,
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// import - all gcs storage location optionals set
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:   externalVolume(config.ListVariable(gcsStorageLocationKmsEncryption), externalVolumeName, comment2, allowWritesFalse),
				ResourceName:      "snowflake_external_volume.complete",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add second storage location with all gcs optionals set
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocationKmsEncryptionUpdatedKey, gcsStorageLocationKmsEncryptionUpdatedName), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(2).
						HasStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeSseKms,
							gcsEncryptionKmsKeyId2,
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName2,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeSseKms,
							gcsEncryptionKmsKeyId,
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "5")),
				),
			},
			// import - 2 gcs storage locations, with all gcs optionals set
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:   externalVolume(config.ListVariable(gcsStorageLocationKmsEncryptionUpdatedKey, gcsStorageLocationKmsEncryptionUpdatedName), externalVolumeName, comment2, allowWritesFalse),
				ResourceName:      "snowflake_external_volume.complete",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// change back to none encryption with 1 gcs storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocation), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// change back to no encryption
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocation), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// remove allow writes and comment from config
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString("").
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment("").
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "3")),
				),
			},
		},
	})
}

// Test volume with azure storage locations
func TestAcc_External_Volume_Azure(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()
	resourceId := helpers.EncodeResourceIdentifier(id)
	comment := random.Comment()
	comment2 := random.Comment()
	allowWritesTrue := "true"
	allowWritesFalse := "false"
	azureStorageLocationName := "azureTest"
	azureStorageLocationName2 := "azureTest2"
	azureStorageProvider := "AZURE"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"
	azureEncryptionTypeNone := "NONE"
	azureStorageLocation := getAzureStorageLocation(azureStorageLocationName, azureStorageBaseUrl, azureTenantId)
	azureStorageLocationUpdatedName := getAzureStorageLocation(azureStorageLocationName2, azureStorageBaseUrl, azureTenantId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			// without optionals
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString("").
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment("").
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "3")),
				),
			},
			// import - without optionals
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, "", ""),
				ResourceName:    "snowflake_external_volume.test",
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedExternalVolumeResource(t, resourceId).
						HasNameString(externalVolumeName).
						HasStorageLocationLength(1),
				),
			},
			// set external volume optionals
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// import - with external volume optionals
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:   externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				ResourceName:      "snowflake_external_volume.complete",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add second storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocation, azureStorageLocationUpdatedName), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(2).
						HasStorageLocationAtIndex(
							0,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						).
						HasStorageLocationAtIndex(
							1,
							azureStorageLocationName2,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "5")),
				),
			},
			// import - 2 storage locations
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables:   externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				ResourceName:      "snowflake_external_volume.complete",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// update comment and change back to 1 storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, comment2, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// update allowWrites
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, comment2, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment2).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment2).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "4")),
				),
			},
			// remove allow writes and comment from config
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocation), externalVolumeName, "", ""),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString("").
						HasAllowWritesString(r.BooleanDefault).
						HasStorageLocationLength(1).
						HasStorageLocationAtIndex(
							0,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment("").
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "3")),
				),
			},
		},
	})
}

// Test apply works when setting all optionals from the start
// Other tests start without setting all optionals
func TestAcc_External_Volume_All_Options(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()
	comment := random.Comment()
	allowWritesFalse := "false"

	s3StorageLocationName := "s3Test"
	s3StorageProvider := "S3"
	s3StorageBaseUrl := "s3://my-example-bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionTypeSseKms := "AWS_SSE_KMS"
	s3EncryptionKmsKeyId := "123456789"
	s3StorageLocationKmsEncryption := getS3StorageLocation(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeSseKms, s3EncryptionKmsKeyId)

	gcsStorageLocationName := "gcsTest"
	gcsStorageProvider := "GCS"
	gcsStorageBaseUrl := "gcs://my-example-bucket"
	gcsEncryptionTypeSseKms := "GCS_SSE_KMS"
	gcsEncryptionKmsKeyId := "123456789"
	gcsStorageLocationKmsEncryption := getGcsStorageLocation(gcsStorageLocationName, gcsStorageBaseUrl, gcsEncryptionTypeSseKms, gcsEncryptionKmsKeyId)

	azureStorageLocationName := "azureTest"
	azureStorageProvider := "AZURE"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"
	azureEncryptionTypeNone := "NONE"
	azureStorageLocation := getAzureStorageLocation(azureStorageLocationName, azureStorageBaseUrl, azureTenantId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocationKmsEncryption), config.ListVariable(gcsStorageLocationKmsEncryption), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(3).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeSseKms,
							s3EncryptionKmsKeyId,
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeSseKms,
							gcsEncryptionKmsKeyId,
							"",
						).
						HasStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "6")),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables:   externalVolumeMultiple(config.ListVariable(s3StorageLocationKmsEncryption), config.ListVariable(gcsStorageLocationKmsEncryption), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesFalse),
				ResourceName:      "snowflake_external_volume.complete",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test volume with multiple storage locations that span multiple providers
// Test adding/removing storage locations at different positions in the storage_location list
func TestAcc_External_Volume_Multiple(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()
	comment := random.Comment()
	allowWritesTrue := "true"
	s3StorageLocationName := "s3Test"
	s3StorageLocationName2 := "s3Test2"
	s3StorageProvider := "S3"
	s3StorageBaseUrl := "s3://my-example-bucket"
	s3StorageBaseUrl2 := "s3://my-example-bucket2"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionTypeNone := "NONE"
	s3StorageLocation := getS3StorageLocation(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeNone, "")
	s3StorageLocationUpdatedBaseUrl := getS3StorageLocation(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl2, s3StorageAwsRoleArn, s3EncryptionTypeNone, "")
	s3StorageLocationUpdatedName := getS3StorageLocation(s3StorageLocationName2, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeNone, "")

	gcsStorageLocationName := "gcsTest"
	gcsStorageLocationName2 := "gcsTest2"
	gcsStorageProvider := "GCS"
	gcsStorageBaseUrl := "gcs://my-example-bucket"
	gcsStorageBaseUrl2 := "gcs://my-example-bucket2"
	gcsEncryptionTypeNone := "NONE"
	gcsStorageLocation := getGcsStorageLocation(gcsStorageLocationName, gcsStorageBaseUrl, gcsEncryptionTypeNone, "")
	gcsStorageLocationUpdatedName := getGcsStorageLocation(gcsStorageLocationName2, gcsStorageBaseUrl, gcsEncryptionTypeNone, "")
	gcsStorageLocationUpdatedBaseUrl := getGcsStorageLocation(gcsStorageLocationName, gcsStorageBaseUrl2, gcsEncryptionTypeNone, "")

	azureStorageLocationName := "azureTest"
	azureStorageLocationName2 := "azureTest2"
	azureStorageProvider := "AZURE"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"
	azureTenantId2 := "987654321"
	azureEncryptionTypeNone := "NONE"
	azureStorageLocation := getAzureStorageLocation(azureStorageLocationName, azureStorageBaseUrl, azureTenantId)
	azureStorageLocationUpdatedTenantId := getAzureStorageLocation(azureStorageLocationName, azureStorageBaseUrl, azureTenantId2)
	azureStorageLocationUpdatedName := getAzureStorageLocation(azureStorageLocationName2, azureStorageBaseUrl, azureTenantId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			// one location of each provider
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(3).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "6")),
				),
			},
			// import
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables:   externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				ResourceName:      "snowflake_external_volume.complete",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// change the s3 base url at position 0
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocationUpdatedBaseUrl), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(3).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl2,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "6")),
				),
			},
			// change back the s3 base url at position 0
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(3).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "6")),
				),
			},
			// add new s3 storage location to position 0
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocationUpdatedName, s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(4).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName2,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							3,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "7")),
				),
			},
			// remove s3 storage location at position 0
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(3).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "6")),
				),
			},
			// change the base url of the gcs storage location at position 1
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocationUpdatedBaseUrl), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(3).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl2,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "6")),
				),
			},
			// change back the encryption type of the gcs storage location at position 1
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(3).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "6")),
				),
			},
			// add new s3 storage location to position 1
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation, s3StorageLocationUpdatedName), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(4).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							s3StorageLocationName2,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							3,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "7")),
				),
			},
			// remove s3 storage location at position 1
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(3).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "6")),
				),
			},
			// change the tenant id of the azure storage location at position 2
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocationUpdatedTenantId), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(3).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId2,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "6")),
				),
			},
			// change back the tenant id of the azure storage location at position 2
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(3).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "6")),
				),
			},
			// add new gcs storage location to position 2
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation, gcsStorageLocationUpdatedName), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(4).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							gcsStorageLocationName2,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							3,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "7")),
				),
			},
			// remove gcs storage location at position 2
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(3).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "6")),
				),
			},
			// add new azure storage location to position 3
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation, azureStorageLocationUpdatedName), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(4).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						).
						HasStorageLocationAtIndex(
							3,
							azureStorageLocationName2,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "7")),
				),
			},
			// remove azure storage location from position 3
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/multiple/complete"),
				ConfigVariables: externalVolumeMultiple(config.ListVariable(s3StorageLocation), config.ListVariable(gcsStorageLocation), config.ListVariable(azureStorageLocation), externalVolumeName, comment, allowWritesTrue),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesTrue).
						HasStorageLocationLength(3).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							gcsStorageLocationName,
							gcsStorageProvider,
							gcsStorageBaseUrl,
							"",
							gcsEncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							2,
							azureStorageLocationName,
							azureStorageProvider,
							azureStorageBaseUrl,
							"",
							azureEncryptionTypeNone,
							"",
							azureTenantId,
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesTrue),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "6")),
				),
			},
		},
	})
}

// Test that drifts are detected and fixed
func TestAcc_External_Volume_External_Changes(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeName := id.Name()
	comment := random.Comment()
	comment2 := random.Comment()
	allowWritesFalse := "false"
	s3StorageLocationName := "s3Test"
	s3StorageLocationName2 := "s3Test2"
	s3StorageProvider := "S3"
	s3StorageBaseUrl := "s3://my-example-bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionTypeNone := "NONE"
	s3StorageLocation := getS3StorageLocation(s3StorageLocationName, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeNone, "")
	s3StorageLocationUpdatedName := getS3StorageLocation(s3StorageLocationName2, s3StorageProvider, s3StorageBaseUrl, s3StorageAwsRoleArn, s3EncryptionTypeNone, "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalVolume),
		Steps: []resource.TestStep{
			// create volume with 2 s3 storage locations
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation, s3StorageLocationUpdatedName), externalVolumeName, comment, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(2).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							s3StorageLocationName2,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "5")),
				),
			},
			// externally remove storage location
			{
				PreConfig: func() {
					testClient().ExternalVolume.Alter(t, sdk.NewAlterExternalVolumeRequest(id).WithRemoveStorageLocation(s3StorageLocationName))
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation, s3StorageLocationUpdatedName), externalVolumeName, comment, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(2).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							s3StorageLocationName2,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "5")),
				),
			},
			// externally add storage location
			{
				PreConfig: func() {
					testClient().ExternalVolume.Alter(
						t,
						sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(
							*sdk.NewExternalVolumeStorageLocationRequest().WithS3StorageLocationParams(
								*sdk.NewS3StorageLocationParamsRequest(
									"externally-added-s3-storage-location",
									"s3",
									"arn:aws:iam::123456789012:role/externally-added-role",
									"s3://externally-added-bucket",
								),
							),
						),
					)
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation, s3StorageLocationUpdatedName), externalVolumeName, comment, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(2).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							s3StorageLocationName2,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "5")),
				),
			},
			// externally drop external volume
			{
				PreConfig: func() {
					testClient().ExternalVolume.DropFunc(t, id)
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation, s3StorageLocationUpdatedName), externalVolumeName, comment, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(2).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							s3StorageLocationName2,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "5")),
				),
			},
			// externally update comment
			{
				PreConfig: func() {
					testClient().ExternalVolume.Alter(t, sdk.NewAlterExternalVolumeRequest(id).WithSet(
						*sdk.NewAlterExternalVolumeSetRequest().WithComment(comment2),
					))
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation, s3StorageLocationUpdatedName), externalVolumeName, comment, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(2).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							s3StorageLocationName2,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "5")),
				),
			},
			// externally update allow writes
			{
				PreConfig: func() {
					testClient().ExternalVolume.Alter(t, sdk.NewAlterExternalVolumeRequest(id).WithSet(
						*sdk.NewAlterExternalVolumeSetRequest().WithAllowWrites(true),
					))
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/complete"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocation, s3StorageLocationUpdatedName), externalVolumeName, comment, allowWritesFalse),
				Check: assertThat(t,
					resourceassert.ExternalVolumeResource(t, "snowflake_external_volume.complete").
						HasNameString(externalVolumeName).
						HasCommentString(comment).
						HasAllowWritesString(allowWritesFalse).
						HasStorageLocationLength(2).
						HasStorageLocationAtIndex(
							0,
							s3StorageLocationName,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						).
						HasStorageLocationAtIndex(
							1,
							s3StorageLocationName2,
							s3StorageProvider,
							s3StorageBaseUrl,
							s3StorageAwsRoleArn,
							s3EncryptionTypeNone,
							"",
							"",
						),
					resourceshowoutputassert.ExternalVolumeShowOutput(t, "snowflake_external_volume.complete").
						HasName(externalVolumeName).
						HasComment(comment).
						HasAllowWrites(allowWritesFalse),
					assert.Check(resource.TestCheckResourceAttr("snowflake_external_volume.complete", "describe_output.#", "5")),
				),
			},
		},
	})
}

// Test invalid parameter combinations throw errors
func TestAcc_External_Volume_Invalid_Cases(t *testing.T) {
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			// invalid storage provider test
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(s3StorageLocationInvalidStorageProvider), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("invalid storage provider: invalid-storage-provider"),
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
				ExpectError:     regexp.MustCompile("unable to extract storage location, missing storage_aws_role_arn key in an s3 storage location"),
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
				ExpectError:     regexp.MustCompile("unable to extract storage location, azure_tenant_id provided for s3 storage location"),
			},
			// storage_aws_role_arn specified for gcs storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocationWithAwsRoleArn), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("unable to extract storage location, storage_aws_role_arn provided for gcs storage location"),
			},
			// azure_tenant_id specified for gcs storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(gcsStorageLocationWithTenantId), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("unable to extract storage location, azure_tenant_id provided for gcs storage location"),
			},
			// storage_aws_role_arn specified for azure storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocationWithAwsRoleArn), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("unable to extract storage location, storage_aws_role_arn provided for azure storage location"),
			},
			// encryption_kms_key_id specified for azure storage location
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalVolume/single/basic"),
				ConfigVariables: externalVolume(config.ListVariable(azureStorageLocationWithKmsKeyId), externalVolumeName, "", ""),
				ExpectError:     regexp.MustCompile("unable to extract storage location, encryption_kms_key_id provided for azure storage location"),
			},
			// TODO add test for encryption_type specified for azure storage location
		},
	})
}

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var s3StorageLocationDescribeSchema = map[string]*schema.Schema{
	"storage_aws_role_arn": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"storage_aws_iam_user_arn": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"storage_aws_external_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"storage_aws_access_point_arn": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"use_privatelink_endpoint": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"encryption_kms_key_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var gcsStorageLocationDescribeSchema = map[string]*schema.Schema{
	"storage_gcp_service_account": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"encryption_kms_key_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var azureStorageLocationDescribeSchema = map[string]*schema.Schema{
	"azure_tenant_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"azure_multi_tenant_app_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"azure_consent_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var s3CompatStorageLocationDescribeSchema = map[string]*schema.Schema{
	"endpoint": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"aws_access_key_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"encryption_kms_key_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var DescribeExternalVolumeSchema = map[string]*schema.Schema{
	"active": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"allow_writes": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"storage_locations": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"storage_provider": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"storage_base_url": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"storage_allowed_locations": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"encryption_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"s3_storage_location": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: s3StorageLocationDescribeSchema,
					},
				},
				"gcs_storage_location": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: gcsStorageLocationDescribeSchema,
					},
				},
				"azure_storage_location": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: azureStorageLocationDescribeSchema,
					},
				},
				"s3_compat_storage_location": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: s3CompatStorageLocationDescribeSchema,
					},
				},
			},
		},
	},
}

func ExternalVolumeDetailsToSchema(details sdk.ExternalVolumeDetails) map[string]any {
	result := map[string]any{
		"active":       details.Active,
		"comment":      details.Comment,
		"allow_writes": details.AllowWrites,
	}

	storageLocations := make([]map[string]any, len(details.StorageLocations))
	for i, loc := range details.StorageLocations {
		locMap := map[string]any{
			"name":                       loc.Name,
			"storage_provider":           loc.StorageProvider,
			"storage_base_url":           loc.StorageBaseUrl,
			"storage_allowed_locations":  loc.StorageAllowedLocations,
			"encryption_type":            loc.EncryptionType,
			"s3_storage_location":        []any{},
			"gcs_storage_location":       []any{},
			"azure_storage_location":     []any{},
			"s3_compat_storage_location": []any{},
		}

		switch {
		case loc.S3StorageLocation != nil:
			locMap["s3_storage_location"] = []any{map[string]any{
				"storage_aws_role_arn":         loc.S3StorageLocation.StorageAwsRoleArn,
				"storage_aws_iam_user_arn":     loc.S3StorageLocation.StorageAwsIamUserArn,
				"storage_aws_external_id":      loc.S3StorageLocation.StorageAwsExternalId,
				"storage_aws_access_point_arn": loc.S3StorageLocation.StorageAwsAccessPointArn,
				"use_privatelink_endpoint":     loc.S3StorageLocation.UsePrivatelinkEndpoint,
				"encryption_kms_key_id":        loc.S3StorageLocation.EncryptionKmsKeyId,
			}}
		case loc.GCSStorageLocation != nil:
			locMap["gcs_storage_location"] = []any{map[string]any{
				"storage_gcp_service_account": loc.GCSStorageLocation.StorageGcpServiceAccount,
				"encryption_kms_key_id":       loc.GCSStorageLocation.EncryptionKmsKeyId,
			}}
		case loc.AzureStorageLocation != nil:
			locMap["azure_storage_location"] = []any{map[string]any{
				"azure_tenant_id":             loc.AzureStorageLocation.AzureTenantId,
				"azure_multi_tenant_app_name": loc.AzureStorageLocation.AzureMultiTenantAppName,
				"azure_consent_url":           loc.AzureStorageLocation.AzureConsentUrl,
			}}
		case loc.S3CompatStorageLocation != nil:
			locMap["s3_compat_storage_location"] = []any{map[string]any{
				"endpoint":              loc.S3CompatStorageLocation.Endpoint,
				"aws_access_key_id":     loc.S3CompatStorageLocation.AwsAccessKeyId,
				"encryption_kms_key_id": loc.S3CompatStorageLocation.EncryptionKmsKeyId,
			}}
		}

		storageLocations[i] = locMap
	}

	result["storage_locations"] = storageLocations
	return result
}

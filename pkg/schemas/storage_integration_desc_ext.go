package schemas

import (
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var StorageIntegrationPropertiesNames = []string{
	"ENABLED",
	"STORAGE_PROVIDER",
	"STORAGE_ALLOWED_LOCATIONS",
	"STORAGE_BLOCKED_LOCATIONS",
	"STORAGE_AWS_IAM_USER_ARN",
	"STORAGE_AWS_OBJECT_ACL",
	"STORAGE_AWS_ROLE_ARN",
	"STORAGE_AWS_EXTERNAL_ID",
	"STORAGE_GCP_SERVICE_ACCOUNT",
	"AZURE_CONSENT_URL",
	"AZURE_MULTI_TENANT_APP_NAME",
	"USE_PRIVATELINK_ENDPOINT",
	"COMMENT",
}

// TODO [v3]: flatten nested DescribePropertyListSchema to typed scalars (or retire in favor of
// DescribeStorageIntegrationAllDetailsSchema). Breaking: describe_output.0.<key>.0.value → <key>.
func (storageIntegrationDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled":                     DescribePropertyListSchema,
		"storage_provider":            DescribePropertyListSchema,
		"storage_allowed_locations":   DescribePropertyListSchema,
		"storage_blocked_locations":   DescribePropertyListSchema,
		"storage_aws_iam_user_arn":    DescribePropertyListSchema,
		"storage_aws_object_acl":      DescribePropertyListSchema,
		"storage_aws_role_arn":        DescribePropertyListSchema,
		"storage_aws_external_id":     DescribePropertyListSchema,
		"storage_gcp_service_account": DescribePropertyListSchema,
		"azure_consent_url":           DescribePropertyListSchema,
		"azure_multi_tenant_app_name": DescribePropertyListSchema,
		"use_privatelink_endpoint":    DescribePropertyListSchema,
		"comment":                     DescribePropertyListSchema,
	}
}

func (storageIntegrationDetailsToSchemaMapper) additionalToSchema(src *sdk.StorageIntegrationDetails, dst map[string]any) {
	if src == nil {
		return
	}
	mapStorageIntegrationProperty(dst, "enabled", src.Enabled)
	mapStorageIntegrationProperty(dst, "storage_provider", src.StorageProvider)
	mapStorageIntegrationProperty(dst, "storage_allowed_locations", src.StorageAllowedLocations)
	mapStorageIntegrationProperty(dst, "storage_blocked_locations", src.StorageBlockedLocations)
	mapStorageIntegrationProperty(dst, "storage_aws_iam_user_arn", src.StorageAwsIamUserArn)
	mapStorageIntegrationProperty(dst, "storage_aws_object_acl", src.StorageAwsObjectAcl)
	mapStorageIntegrationProperty(dst, "storage_aws_role_arn", src.StorageAwsRoleArn)
	mapStorageIntegrationProperty(dst, "storage_aws_external_id", src.StorageAwsExternalId)
	mapStorageIntegrationProperty(dst, "storage_gcp_service_account", src.StorageGcpServiceAccount)
	mapStorageIntegrationProperty(dst, "azure_consent_url", src.AzureConsentUrl)
	mapStorageIntegrationProperty(dst, "azure_multi_tenant_app_name", src.AzureMultiTenantAppName)
	mapStorageIntegrationProperty(dst, "use_privatelink_endpoint", src.UsePrivatelinkEndpoint)
	mapStorageIntegrationProperty(dst, "comment", src.Comment)
}

func mapStorageIntegrationProperty(dst map[string]any, key string, property *sdk.StorageIntegrationProperty) {
	if property == nil {
		return
	}
	dst[key] = []map[string]any{StorageIntegrationPropertyToSchema(property)}
}

func DescribeStorageIntegrationToSchema(integrationProperties []sdk.StorageIntegrationProperty) map[string]any {
	for _, property := range integrationProperties {
		if !slices.Contains(StorageIntegrationPropertiesNames, property.Name) {
			log.Printf("[DEBUG] Unknown storage integration property %s", property.Name)
		}
	}
	return StorageIntegrationDetailsToSchema(sdk.AsStorageIntegrationDescribe(integrationProperties))
}

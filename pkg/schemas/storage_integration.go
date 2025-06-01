package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeStorageIntegrationSchema represents output of DESCRIBE query for the single StorageIntegration.
var DescribeStorageIntegrationSchema = map[string]*schema.Schema{
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
}

var _ = DescribeStorageIntegrationSchema

func DescribeStorageIntegrationToSchema(integrationProperties []sdk.StorageIntegrationProperty) map[string]any {
	propsSchema := make(map[string]any)
	for _, property := range integrationProperties {
		property := property
		// Convert property name to lowercase and add to schema
		propertyName := strings.ToLower(property.Name)
		propsSchema[propertyName] = []map[string]any{StorageIntegrationPropertyToSchema(&property)}
	}
	return propsSchema
}

// Helper function to convert StorageIntegrationProperty to schema
func StorageIntegrationPropertyToSchema(property *sdk.StorageIntegrationProperty) map[string]any {
	return map[string]any{
		"name":    property.Name,
		"type":    property.Type,
		"value":   property.Value,
		"default": property.Default,
	}
}

var _ = DescribeStorageIntegrationToSchema

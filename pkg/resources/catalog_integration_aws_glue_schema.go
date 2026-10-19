package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var catalogIntegrationAwsGlueSchema = collections.MergeMaps(
	catalogIntegrationAwsGlueIdentitySchemaExt,
	catalogIntegrationAwsGlueAttributesSchemaExt,
	map[string]*schema.Schema{
		FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	},
	catalogIntegrationAwsGlueShowOutputSchemaExt,
	catalogIntegrationAwsGlueDescribeOutputSchemaExt,
)

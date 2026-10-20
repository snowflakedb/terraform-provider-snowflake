package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Snowflake describes every catalog integration with DESCRIBE CATALOG INTEGRATION.
var catalogIntegrationObjectStorageDescribeOutputSchemaExt = map[string]*schema.Schema{
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE CATALOG INTEGRATION` for the given catalog integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeCatalogIntegrationObjectStorageDetailsSchema,
		},
	},
}

// catalogIntegrationObjectStorageAttributesModificationExt adds behavior the generator does not emit.
// catalog_source is describe-only. table_format is an SSOT attribute and stays in the generated schema.
func catalogIntegrationObjectStorageAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["refresh_interval_seconds"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInDescribe("refresh_interval_seconds")
	s["catalog_source"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: joinWithSpace("Specifies the type of catalog source.",
			"This field is used to detect when the catalog source was changed outside of Terraform and to recreate the resource when that happens."),
	}
	return s
}

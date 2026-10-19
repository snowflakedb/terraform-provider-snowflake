package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var catalogIntegrationObjectStorageIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier (i.e. name) of the catalog integration; must be unique in your account.",
	},
}

var catalogIntegrationObjectStorageAttributesSchemaExt = map[string]*schema.Schema{
	"enabled": {
		Type:     schema.TypeBool,
		Required: true,
		Description: joinWithSpace("Specifies whether the catalog integration is available for use for Iceberg tables.",
			"`true` allows users to create new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration function normally.",
			"`false` prevents users from creating new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration cannot access the catalog in the table definition."),
	},
	"refresh_interval_seconds": {
		Type:     schema.TypeInt,
		Optional: true,
		Description: joinWithSpace("Specifies the number of seconds to wait between attempts to poll the external Iceberg catalog for metadata updates for automated refresh.",
			"For Delta-based tables, specifies the number of seconds to wait between attempts to poll your external cloud storage for new metadata."),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("refresh_interval_seconds"),
		ValidateFunc:     validation.IntAtLeast(1),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "Specifies a comment for the catalog integration.",
	},
	"catalog_source": {
		Type:     schema.TypeString,
		Computed: true,
		Description: joinWithSpace("Specifies the type of catalog source.",
			"This field is used to detect when the catalog source was changed outside of Terraform and to recreate the resource when that happens."),
	},
	"table_format": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Specifies the table format. " + enumValuesDescription(sdk.AllCatalogIntegrationTableFormats),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToCatalogIntegrationTableFormat),
		ValidateDiagFunc: sdkValidation(sdk.ToCatalogIntegrationTableFormat),
	},
}

var catalogIntegrationObjectStorageShowOutputSchemaExt = map[string]*schema.Schema{
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW CATALOG INTEGRATIONS` for the given catalog integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowCatalogIntegrationSchema,
		},
	},
}

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

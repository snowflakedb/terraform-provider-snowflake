package schemas

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeStreamlitSchema represents output of DESCRIBE query for the single streamlit.
var DescribeStreamlitSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"title": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"root_location": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"main_file": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"query_warehouse": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"url_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"default_packages": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"user_packages": {
		Type:     schema.TypeSet,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	},
	"import_urls": {
		Type:     schema.TypeSet,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	},
	"external_access_integrations": {
		Type:     schema.TypeSet,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	},
	"external_access_secrets": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func StreamlitPropertiesToSchema(details sdk.StreamlitDetail) (map[string]any, error) {
	stageId, location, err := helpers.ParseRootLocation(details.RootLocation)
	if err != nil {
		return nil, err
	}
	rootLocation := fmt.Sprintf("@%s", stageId.FullyQualifiedName())
	if len(location) > 0 {
		rootLocation = fmt.Sprintf("%s/%s", rootLocation, location)
	}
	return map[string]any{
		"name":                         details.Name,
		"title":                        details.Title,
		"root_location":                rootLocation,
		"main_file":                    details.MainFile,
		"query_warehouse":              details.QueryWarehouse,
		"url_id":                       details.UrlId,
		"default_packages":             details.DefaultPackages,
		"user_packages":                details.UserPackages,
		"import_urls":                  details.ImportUrls,
		"external_access_integrations": details.ExternalAccessIntegrations,
		"external_access_secrets":      details.ExternalAccessSecrets,
	}, nil
}

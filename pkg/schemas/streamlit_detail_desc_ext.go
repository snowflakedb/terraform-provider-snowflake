package schemas

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (streamlitDetailToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"root_location": {
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
	}
}

func (streamlitDetailToSchemaMapper) additionalToSchema(src *sdk.StreamlitDetail, dst map[string]any) {
	rootLocation := src.RootLocation
	if stageId, location, err := helpers.ParseRootLocation(src.RootLocation); err == nil {
		rootLocation = fmt.Sprintf("@%s", stageId.FullyQualifiedName())
		if len(location) > 0 {
			rootLocation = fmt.Sprintf("%s/%s", rootLocation, location)
		}
	}
	dst["root_location"] = rootLocation
	dst["user_packages"] = src.UserPackages
	dst["import_urls"] = src.ImportUrls
	dst["external_access_integrations"] = src.ExternalAccessIntegrations
}

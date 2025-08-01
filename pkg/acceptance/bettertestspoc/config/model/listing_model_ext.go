package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func ListingWithStagedManifestWithOptionals(
	resourceName string,
	name string,
	stageId sdk.SchemaObjectIdentifier,
	versionName string,
	location string,
) *ListingModel {
	l := &ListingModel{ResourceModelMeta: config.Meta(resourceName, resources.Listing)}
	l.WithName(name)
	l.WithManifestValue(tfconfig.ListVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"from_stage": tfconfig.ListVariable(
				tfconfig.MapVariable(map[string]tfconfig.Variable{
					"stage":        tfconfig.StringVariable(stageId.FullyQualifiedName()),
					"version_name": tfconfig.StringVariable(versionName),
					"location":     tfconfig.StringVariable(location),
				}),
			),
		}),
	))
	return l
}

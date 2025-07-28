package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func ListingWithInlineManifest(
	resourceName string,
	name string,
	manifest string,
) *ListingModel {
	return Listing(resourceName, name, tfconfig.ListVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"from_string": tfconfig.StringVariable(manifest),
		}),
	))
}

func ListingWithStagedManifest(
	resourceName string,
	name string,
	stageId sdk.SchemaObjectIdentifier,
) *ListingModel {
	return Listing(resourceName, name, tfconfig.ListVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"from_stage": tfconfig.ListVariable(
				tfconfig.MapVariable(map[string]tfconfig.Variable{
					"stage": tfconfig.StringVariable(stageId.FullyQualifiedName()),
				}),
			),
		}),
	))
}

func ListingWithStagedManifestWithOptionals(
	resourceName string,
	name string,
	stageId sdk.SchemaObjectIdentifier,
	versionName string,
	location string,
) *ListingModel {
	return Listing(resourceName, name, tfconfig.ListVariable(
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
}

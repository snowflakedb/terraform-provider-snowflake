package model

import (
	"log"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TODO(SNOW-1501905): Remove after complex non-list type overrides are handled
func (l *ListingModel) WithManifest(locations []sdk.StageLocation) *ListingModel {
	if len(locations) != 1 {
		log.Fatalf("expected exactly one location for manifest, got %d", len(locations))
	}

	return l.WithManifestValue(tfconfig.ListVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"from_stage": tfconfig.ListVariable(
				tfconfig.MapVariable(map[string]tfconfig.Variable{
					"stage":    tfconfig.StringVariable(locations[0].GetStageId().FullyQualifiedName()),
					"location": tfconfig.StringVariable(locations[0].GetPath()),
				}),
			),
		}),
	))
}

func ListingWithInlineManifest(
	resourceName string,
	name string,
	manifest string,
) *ListingModel {
	l := &ListingModel{ResourceModelMeta: config.Meta(resourceName, resources.Listing)}
	l.WithName(name)
	l.WithManifestValue(tfconfig.ListVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"from_string": config.MultilineWrapperVariable(manifest),
		}),
	))
	return l
}

func ListingWithStagedManifest(
	resourceName string,
	name string,
	stageId sdk.SchemaObjectIdentifier,
) *ListingModel {
	l := &ListingModel{ResourceModelMeta: config.Meta(resourceName, resources.Listing)}
	l.WithName(name)
	l.WithManifestValue(tfconfig.ListVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"from_stage": tfconfig.ListVariable(
				tfconfig.MapVariable(map[string]tfconfig.Variable{
					"stage": tfconfig.StringVariable(stageId.FullyQualifiedName()),
				}),
			),
		}),
	))
	return l
}

func ListingWithStagedManifestWithLocation(
	resourceName string,
	name string,
	stageId sdk.SchemaObjectIdentifier,
	location string,
) *ListingModel {
	l := &ListingModel{ResourceModelMeta: config.Meta(resourceName, resources.Listing)}
	l.WithName(name)
	l.WithManifestValue(tfconfig.ListVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"from_stage": tfconfig.ListVariable(
				tfconfig.MapVariable(map[string]tfconfig.Variable{
					"stage":    tfconfig.StringVariable(stageId.FullyQualifiedName()),
					"location": tfconfig.StringVariable(location),
				}),
			),
		}),
	))
	return l
}

func ListingWithStagedManifestWithOptionals(
	resourceName string,
	name string,
	stageId sdk.SchemaObjectIdentifier,
	versionName string,
	versionComment string,
	location string,
) *ListingModel {
	l := &ListingModel{ResourceModelMeta: config.Meta(resourceName, resources.Listing)}
	l.WithName(name)
	l.WithManifestValue(tfconfig.ListVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"from_stage": tfconfig.ListVariable(
				tfconfig.MapVariable(map[string]tfconfig.Variable{
					"stage":           tfconfig.StringVariable(stageId.FullyQualifiedName()),
					"version_name":    tfconfig.StringVariable(versionName),
					"version_comment": tfconfig.StringVariable(versionComment),
					"location":        tfconfig.StringVariable(location),
				}),
			),
		}),
	))
	return l
}

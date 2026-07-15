package model

import (
	"log"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TODO(SNOW-1501905): Remove after complex non-list type overrides are handled
func (l *OrganizationListingModel) WithManifest(locations []sdk.StageLocation) *OrganizationListingModel {
	if len(locations) != 1 {
		log.Panicf("expected exactly one location for manifest, got %d", len(locations))
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

func OrganizationListingWithInlineManifest(
	resourceName string,
	name string,
	manifest string,
) *OrganizationListingModel {
	l := &OrganizationListingModel{ResourceModelMeta: config.Meta(resourceName, resources.OrganizationListing)}
	l.WithName(name)
	l.WithManifestValue(tfconfig.ListVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"from_string": config.MultilineWrapperVariable(manifest),
		}),
	))
	return l
}

func OrganizationListingWithStagedManifestWithLocation(
	resourceName string,
	name string,
	stageId sdk.SchemaObjectIdentifier,
	location string,
) *OrganizationListingModel {
	return OrganizationListing(resourceName, name, []sdk.StageLocation{sdk.NewStageLocation(stageId, location)})
}

func OrganizationListingWithStagedManifestWithOptionals(
	resourceName string,
	name string,
	stageId sdk.SchemaObjectIdentifier,
	versionName string,
	versionComment string,
	location string,
) *OrganizationListingModel {
	l := &OrganizationListingModel{ResourceModelMeta: config.Meta(resourceName, resources.OrganizationListing)}
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

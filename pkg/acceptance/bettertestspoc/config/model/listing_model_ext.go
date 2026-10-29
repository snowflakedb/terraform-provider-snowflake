package model

import (
	"log"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// listingAcceptanceTimeout is longer than the provider default of 20 minutes.
// CREATE LISTING on preprod has hit that deadline while the same create finishes
// in seconds on later runs. Acceptance configs wait 30 minutes so a slow
// Snowflake response is less likely to end as context deadline exceeded.
const listingAcceptanceTimeout = "30m"

func withListingAcceptanceTimeout(l *ListingModel) *ListingModel {
	return l.WithTimeout(config.Timeouts{
		Create: listingAcceptanceTimeout,
		Read:   listingAcceptanceTimeout,
		Update: listingAcceptanceTimeout,
		Delete: listingAcceptanceTimeout,
	})
}

// TODO(SNOW-1501905): Remove after complex non-list type overrides are handled
func (l *ListingModel) WithManifest(locations []sdk.StageLocation) *ListingModel {
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
	return withListingAcceptanceTimeout(l)
}

func ListingWithStagedManifestWithLocation(
	resourceName string,
	name string,
	stageId sdk.SchemaObjectIdentifier,
	location string,
) *ListingModel {
	return withListingAcceptanceTimeout(Listing(resourceName, name, []sdk.StageLocation{sdk.NewStageLocation(stageId, location)}))
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
	return withListingAcceptanceTimeout(l)
}

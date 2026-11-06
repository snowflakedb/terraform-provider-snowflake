package model

import (
	"log"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func withOrganizationListingAcceptanceTimeout(o *OrganizationListingModel) *OrganizationListingModel {
	return o.WithTimeout(config.Timeouts{
		Create: listingAcceptanceTimeout,
		Read:   listingAcceptanceTimeout,
		Update: listingAcceptanceTimeout,
		Delete: listingAcceptanceTimeout,
	})
}

// TODO(SNOW-1501905): Remove after complex non-list type overrides are handled
func (o *OrganizationListingModel) WithManifest(locations []sdk.StageLocation) *OrganizationListingModel {
	if len(locations) != 1 {
		log.Panicf("expected exactly one location for manifest, got %d", len(locations))
	}

	return o.WithManifestValue(tfconfig.ListVariable(
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
	o := &OrganizationListingModel{ResourceModelMeta: config.Meta(resourceName, resources.OrganizationListing)}
	o.WithName(name)
	o.WithManifestValue(tfconfig.ListVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"from_string": config.MultilineWrapperVariable(manifest),
		}),
	))
	return withOrganizationListingAcceptanceTimeout(o)
}

func OrganizationListingWithStagedManifestWithLocation(
	resourceName string,
	name string,
	stageId sdk.SchemaObjectIdentifier,
	location string,
) *OrganizationListingModel {
	return withOrganizationListingAcceptanceTimeout(OrganizationListing(resourceName, name, []sdk.StageLocation{sdk.NewStageLocation(stageId, location)}))
}

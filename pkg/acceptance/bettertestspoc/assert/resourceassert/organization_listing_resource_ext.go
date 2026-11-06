package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (o *OrganizationListingResourceAssert) HasManifestFromString(manifest string) *OrganizationListingResourceAssert {
	o.ValueSet("manifest.0.from_string", manifest)
	return o
}

func (o *OrganizationListingResourceAssert) HasManifestFromStringNotEmpty() *OrganizationListingResourceAssert {
	o.ValuePresent("manifest.0.from_string")
	return o
}

func (o *OrganizationListingResourceAssert) HasManifestFromStageNotEmpty() *OrganizationListingResourceAssert {
	o.ValuePresent("manifest.0.from_stage.0.stage")
	return o
}

func (o *OrganizationListingResourceAssert) HasManifestFromStageStageId(stageId sdk.SchemaObjectIdentifier) *OrganizationListingResourceAssert {
	o.ValueSet("manifest.0.from_stage.0.stage", stageId.FullyQualifiedName())
	return o
}

func (o *OrganizationListingResourceAssert) HasManifestFromStageLocation(location string) *OrganizationListingResourceAssert {
	o.ValueSet("manifest.0.from_stage.0.location", location)
	return o
}

func (o *OrganizationListingResourceAssert) HasNoManifest() *OrganizationListingResourceAssert {
	o.ValueNotSet("manifest")
	return o
}

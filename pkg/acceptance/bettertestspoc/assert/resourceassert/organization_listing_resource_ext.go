package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (l *OrganizationListingResourceAssert) HasManifestFromString(manifest string) *OrganizationListingResourceAssert {
	l.AddAssertion(assert.ValueSet("manifest.0.from_string", manifest))
	return l
}

func (l *OrganizationListingResourceAssert) HasManifestFromStringNotEmpty() *OrganizationListingResourceAssert {
	l.AddAssertion(assert.ValuePresent("manifest.0.from_string"))
	return l
}

func (l *OrganizationListingResourceAssert) HasManifestFromStageNotEmpty() *OrganizationListingResourceAssert {
	l.AddAssertion(assert.ValuePresent("manifest.0.from_stage.0.stage"))
	return l
}

func (l *OrganizationListingResourceAssert) HasManifestFromStageStageId(stageId sdk.SchemaObjectIdentifier) *OrganizationListingResourceAssert {
	l.AddAssertion(assert.ValueSet("manifest.0.from_stage.0.stage", stageId.FullyQualifiedName()))
	return l
}

func (l *OrganizationListingResourceAssert) HasManifestFromStageVersionName(versionName string) *OrganizationListingResourceAssert {
	l.AddAssertion(assert.ValueSet("manifest.0.from_stage.0.version_name", versionName))
	return l
}

func (l *OrganizationListingResourceAssert) HasManifestFromStageVersionComment(versionComment string) *OrganizationListingResourceAssert {
	l.AddAssertion(assert.ValueSet("manifest.0.from_stage.0.version_comment", versionComment))
	return l
}

func (l *OrganizationListingResourceAssert) HasManifestFromStageLocation(location string) *OrganizationListingResourceAssert {
	l.AddAssertion(assert.ValueSet("manifest.0.from_stage.0.location", location))
	return l
}

func (l *OrganizationListingResourceAssert) HasNoManifest() *OrganizationListingResourceAssert {
	l.AddAssertion(assert.ValueNotSet("manifest"))
	return l
}

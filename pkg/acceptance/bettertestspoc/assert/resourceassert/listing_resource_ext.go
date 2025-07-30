package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (l *ListingResourceAssert) HasManifestFromString(manifest string) *ListingResourceAssert {
	l.AddAssertion(assert.ValueSet("manifest.0.from_string", manifest))
	return l
}

func (l *ListingResourceAssert) HasManifestFromStringNotEmpty() *ListingResourceAssert {
	l.AddAssertion(assert.ValuePresent("manifest.0.from_string"))
	return l
}

func (l *ListingResourceAssert) HasManifestFromStageNotEmpty() *ListingResourceAssert {
	l.AddAssertion(assert.ValuePresent("manifest.0.from_stage.0.stage"))
	return l
}

func (l *ListingResourceAssert) HasManifestFromStageStageId(stageId sdk.SchemaObjectIdentifier) *ListingResourceAssert {
	l.AddAssertion(assert.ValueSet("manifest.0.from_stage.0.stage", stageId.FullyQualifiedName()))
	return l
}

func (l *ListingResourceAssert) HasManifestFromStageVersionName(versionName string) *ListingResourceAssert {
	l.AddAssertion(assert.ValueSet("manifest.0.from_stage.0.version_name", versionName))
	return l
}

func (l *ListingResourceAssert) HasManifestFromStageVersionComment(versionComment string) *ListingResourceAssert {
	l.AddAssertion(assert.ValueSet("manifest.0.from_stage.0.version_comment", versionComment))
	return l
}

func (l *ListingResourceAssert) HasManifestFromStageLocation(location string) *ListingResourceAssert {
	l.AddAssertion(assert.ValueSet("manifest.0.from_stage.0.location", location))
	return l
}

func (l *ListingResourceAssert) HasNoManifest() *ListingResourceAssert {
	l.AddAssertion(assert.ValueNotSet("manifest"))
	return l
}

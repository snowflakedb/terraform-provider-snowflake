package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
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

func (l *ListingResourceAssert) HasNoManifest() *ListingResourceAssert {
	l.AddAssertion(assert.ValueNotSet("manifest"))
	return l
}

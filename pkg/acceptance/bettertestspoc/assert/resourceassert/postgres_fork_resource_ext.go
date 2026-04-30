package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type PostgresForkResourceAssert struct {
	*assert.ResourceAssert
}

func PostgresForkResource(t *testing.T, name string) *PostgresForkResourceAssert {
	t.Helper()

	return &PostgresForkResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedPostgresForkResource(t *testing.T, id string) *PostgresForkResourceAssert {
	t.Helper()

	return &PostgresForkResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

func (p *PostgresForkResourceAssert) HasNameString(expected string) *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueSet("name", expected))
	return p
}

func (p *PostgresForkResourceAssert) HasForkFromString(expected string) *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueSet("fork_from", expected))
	return p
}

func (p *PostgresForkResourceAssert) HasAtTimestampString(expected string) *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueSet("at_timestamp", expected))
	return p
}

func (p *PostgresForkResourceAssert) HasAtOffsetString(expected string) *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueSet("at_offset", expected))
	return p
}

func (p *PostgresForkResourceAssert) HasBeforeTimestampString(expected string) *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueSet("before_timestamp", expected))
	return p
}

func (p *PostgresForkResourceAssert) HasBeforeOffsetString(expected string) *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueSet("before_offset", expected))
	return p
}

func (p *PostgresForkResourceAssert) HasComputeFamilyString(expected string) *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueSet("compute_family", expected))
	return p
}

func (p *PostgresForkResourceAssert) HasStorageSizeGbString(expected string) *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueSet("storage_size_gb", expected))
	return p
}

func (p *PostgresForkResourceAssert) HasAuthenticationAuthorityString(expected string) *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueSet("authentication_authority", expected))
	return p
}

func (p *PostgresForkResourceAssert) HasHighAvailabilityString(expected string) *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueSet("high_availability", expected))
	return p
}

func (p *PostgresForkResourceAssert) HasCommentString(expected string) *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueSet("comment", expected))
	return p
}

func (p *PostgresForkResourceAssert) HasFullyQualifiedNameString(expected string) *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return p
}

func (p *PostgresForkResourceAssert) HasNoComment() *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueNotSet("comment"))
	return p
}

func (p *PostgresForkResourceAssert) HasNoAtTimestamp() *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueNotSet("at_timestamp"))
	return p
}

func (p *PostgresForkResourceAssert) HasNoAtOffset() *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueNotSet("at_offset"))
	return p
}

func (p *PostgresForkResourceAssert) HasNoBeforeTimestamp() *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueNotSet("before_timestamp"))
	return p
}

func (p *PostgresForkResourceAssert) HasNoBeforeOffset() *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueNotSet("before_offset"))
	return p
}

func (p *PostgresForkResourceAssert) HasNoPostgresSettings() *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueNotSet("postgres_settings"))
	return p
}

// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type SecretWithBasicAuthenticationResourceAssert struct {
	*assert.ResourceAssert
}

func SecretWithBasicAuthenticationResource(t *testing.T, name string) *SecretWithBasicAuthenticationResourceAssert {
	t.Helper()

	return &SecretWithBasicAuthenticationResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedSecretWithBasicAuthenticationResource(t *testing.T, id string) *SecretWithBasicAuthenticationResourceAssert {
	t.Helper()

	return &SecretWithBasicAuthenticationResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (s *SecretWithBasicAuthenticationResourceAssert) HasCommentString(expected string) *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("comment", expected))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasDatabaseString(expected string) *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("database", expected))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasFullyQualifiedNameString(expected string) *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasNameString(expected string) *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("name", expected))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasPasswordString(expected string) *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("password", expected))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasSchemaString(expected string) *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("schema", expected))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasSecretTypeString(expected string) *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("secret_type", expected))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasUsernameString(expected string) *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("username", expected))
	return s
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (s *SecretWithBasicAuthenticationResourceAssert) HasNoComment() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueNotSet("comment"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasNoDatabase() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueNotSet("database"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasNoFullyQualifiedName() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasNoName() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueNotSet("name"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasNoPassword() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueNotSet("password"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasNoSchema() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueNotSet("schema"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasNoSecretType() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueNotSet("secret_type"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasNoUsername() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueNotSet("username"))
	return s
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (s *SecretWithBasicAuthenticationResourceAssert) HasCommentEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("comment", ""))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasDatabaseEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("database", ""))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasFullyQualifiedNameEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("fully_qualified_name", ""))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasNameEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("name", ""))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasPasswordEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("password", ""))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasSchemaEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("schema", ""))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasSecretTypeEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("secret_type", ""))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasUsernameEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValueSet("username", ""))
	return s
}

///////////////////////////////
// Attribute presence checks //
///////////////////////////////

func (s *SecretWithBasicAuthenticationResourceAssert) HasCommentNotEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValuePresent("comment"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasDatabaseNotEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValuePresent("database"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasFullyQualifiedNameNotEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValuePresent("fully_qualified_name"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasNameNotEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValuePresent("name"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasPasswordNotEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValuePresent("password"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasSchemaNotEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValuePresent("schema"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasSecretTypeNotEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValuePresent("secret_type"))
	return s
}

func (s *SecretWithBasicAuthenticationResourceAssert) HasUsernameNotEmpty() *SecretWithBasicAuthenticationResourceAssert {
	s.AddAssertion(assert.ValuePresent("username"))
	return s
}

// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type SecretWithClientCredentialsResourceAssert struct {
	*assert.ResourceAssert
}

func SecretWithClientCredentialsResource(t *testing.T, name string) *SecretWithClientCredentialsResourceAssert {
	t.Helper()

	return &SecretWithClientCredentialsResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedSecretWithClientCredentialsResource(t *testing.T, id string) *SecretWithClientCredentialsResourceAssert {
	t.Helper()

	return &SecretWithClientCredentialsResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (s *SecretWithClientCredentialsResourceAssert) HasApiAuthenticationString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("api_authentication", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasCommentString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("comment", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasDatabaseString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("database", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasFullyQualifiedNameString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNameString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("name", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasOauthScopesString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("oauth_scopes", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasSchemaString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("schema", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasSecretTypeString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("secret_type", expected))
	return s
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (s *SecretWithClientCredentialsResourceAssert) HasNoApiAuthentication() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("api_authentication"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoComment() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("comment"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoDatabase() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("database"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoFullyQualifiedName() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoName() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("name"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoOauthScopes() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("oauth_scopes.#", "0"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoSchema() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("schema"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoSecretType() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("secret_type"))
	return s
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (s *SecretWithClientCredentialsResourceAssert) HasApiAuthenticationEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("api_authentication", ""))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasCommentEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("comment", ""))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasDatabaseEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("database", ""))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasFullyQualifiedNameEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("fully_qualified_name", ""))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNameEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("name", ""))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasSchemaEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("schema", ""))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasSecretTypeEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("secret_type", ""))
	return s
}

///////////////////////////////
// Attribute presence checks //
///////////////////////////////

func (s *SecretWithClientCredentialsResourceAssert) HasApiAuthenticationNotEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValuePresent("api_authentication"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasCommentNotEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValuePresent("comment"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasDatabaseNotEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValuePresent("database"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasFullyQualifiedNameNotEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValuePresent("fully_qualified_name"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNameNotEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValuePresent("name"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasOauthScopesNotEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValuePresent("oauth_scopes"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasSchemaNotEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValuePresent("schema"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasSecretTypeNotEmpty() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValuePresent("secret_type"))
	return s
}

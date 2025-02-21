// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type TagAssociationResourceAssert struct {
	*assert.ResourceAssert
}

func TagAssociationResource(t *testing.T, name string) *TagAssociationResourceAssert {
	t.Helper()

	return &TagAssociationResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedTagAssociationResource(t *testing.T, id string) *TagAssociationResourceAssert {
	t.Helper()

	return &TagAssociationResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (t *TagAssociationResourceAssert) HasObjectIdentifiersString(expected string) *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueSet("object_identifiers", expected))
	return t
}

func (t *TagAssociationResourceAssert) HasObjectTypeString(expected string) *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueSet("object_type", expected))
	return t
}

func (t *TagAssociationResourceAssert) HasSkipValidationString(expected string) *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueSet("skip_validation", expected))
	return t
}

func (t *TagAssociationResourceAssert) HasTagIdString(expected string) *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueSet("tag_id", expected))
	return t
}

func (t *TagAssociationResourceAssert) HasTagValueString(expected string) *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueSet("tag_value", expected))
	return t
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (t *TagAssociationResourceAssert) HasNoObjectIdentifiers() *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueSet("object_identifiers.#", "0"))
	return t
}

func (t *TagAssociationResourceAssert) HasNoObjectType() *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueNotSet("object_type"))
	return t
}

func (t *TagAssociationResourceAssert) HasNoSkipValidation() *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueNotSet("skip_validation"))
	return t
}

func (t *TagAssociationResourceAssert) HasNoTagId() *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueNotSet("tag_id"))
	return t
}

func (t *TagAssociationResourceAssert) HasNoTagValue() *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueNotSet("tag_value"))
	return t
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (t *TagAssociationResourceAssert) HasObjectTypeEmpty() *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueSet("object_type", ""))
	return t
}

func (t *TagAssociationResourceAssert) HasTagIdEmpty() *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueSet("tag_id", ""))
	return t
}

func (t *TagAssociationResourceAssert) HasTagValueEmpty() *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueSet("tag_value", ""))
	return t
}

///////////////////////////////
// Attribute presence checks //
///////////////////////////////

func (t *TagAssociationResourceAssert) HasObjectIdentifiersNotEmpty() *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValuePresent("object_identifiers"))
	return t
}

func (t *TagAssociationResourceAssert) HasObjectTypeNotEmpty() *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValuePresent("object_type"))
	return t
}

func (t *TagAssociationResourceAssert) HasSkipValidationNotEmpty() *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValuePresent("skip_validation"))
	return t
}

func (t *TagAssociationResourceAssert) HasTagIdNotEmpty() *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValuePresent("tag_id"))
	return t
}

func (t *TagAssociationResourceAssert) HasTagValueNotEmpty() *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValuePresent("tag_value"))
	return t
}

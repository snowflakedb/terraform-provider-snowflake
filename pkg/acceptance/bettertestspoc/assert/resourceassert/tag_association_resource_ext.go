package resourceassert

import (
	"fmt"
)

func (t *TagAssociationResourceAssert) HasObjectIdentifiersLength(len int) *TagAssociationResourceAssert {
	t.ValueSet("object_identifiers.#", fmt.Sprintf("%d", len))
	return t
}

func (t *TagAssociationResourceAssert) HasIdString(expected string) *TagAssociationResourceAssert {
	t.ValueSet("id", expected)
	return t
}

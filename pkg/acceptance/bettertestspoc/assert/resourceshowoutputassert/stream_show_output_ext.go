package resourceshowoutputassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StreamShowOutputAssert) HasCreatedOnNotEmpty() *StreamShowOutputAssert {
	s.ValuePresent("created_on")
	return s
}

func (s *StreamShowOutputAssert) HasStaleAfterNotEmpty() *StreamShowOutputAssert {
	s.ValuePresent("stale_after")
	return s
}

func (s *StreamShowOutputAssert) HasBaseTables(ids ...sdk.SchemaObjectIdentifier) *StreamShowOutputAssert {
	s.StringValueSet("base_tables.#", strconv.FormatInt(int64(len(ids)), 10))
	for i := range ids {
		s.SetContainsElem("base_tables", ids[i].FullyQualifiedName())
	}
	return s
}

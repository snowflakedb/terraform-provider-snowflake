package resourceshowoutputassert

import (
	"strconv"
)

func (s *TagShowOutputAssert) HasCreatedOnNotEmpty() *TagShowOutputAssert {
	s.ValuePresent("created_on")
	return s
}

func (s *TagShowOutputAssert) HasAllowedValues(expected ...string) *TagShowOutputAssert {
	s.StringValueSet("allowed_values.#", strconv.FormatInt(int64(len(expected)), 10))
	for _, v := range expected {
		s.SetContainsElem("allowed_values", v)
	}
	return s
}

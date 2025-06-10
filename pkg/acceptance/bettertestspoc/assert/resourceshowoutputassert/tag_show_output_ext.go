package resourceshowoutputassert

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// TagsDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func TagsDatasourceShowOutput(t *testing.T, name string) *TagShowOutputAssert {
	t.Helper()

	s := TagShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "tags.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}

func (s *TagShowOutputAssert) HasCreatedOnNotEmpty() *TagShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return s
}

func (s *TagShowOutputAssert) HasAllowedValues(expected ...string) *TagShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("allowed_values.#", strconv.FormatInt(int64(len(expected)), 10)))
	for i := range expected {
		s.AddAssertion(assert.ResourceShowOutputValueSet(fmt.Sprintf("allowed_values.%d", i), expected[i]))
	}
	return s
}

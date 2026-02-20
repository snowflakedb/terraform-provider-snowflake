package resourceshowoutputassert

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type NetworkRuleDescOutputAssert struct {
	*assert.ResourceAssert
}

func NetworkRuleDescOutput(t *testing.T, name string) *NetworkRuleDescOutputAssert {
	t.Helper()

	networkRuleAssert := NetworkRuleDescOutputAssert{
		ResourceAssert: assert.NewResourceAssert(name, "describe_output"),
	}
	networkRuleAssert.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &networkRuleAssert
}

func ImportedNetworkRuleDescOutput(t *testing.T, id string) *NetworkRuleDescOutputAssert {
	t.Helper()

	networkRuleAssert := NetworkRuleDescOutputAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "describe_output"),
	}
	networkRuleAssert.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &networkRuleAssert
}

////////////////////////////
// Attribute value checks //
////////////////////////////

func (n *NetworkRuleDescOutputAssert) HasCreatedOn(expected time.Time) *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueSet("created_on", expected.String()))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasCreatedOnNotEmpty() *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValuePresent("created_on"))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasName(expected string) *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueSet("name", expected))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasDatabaseName(expected string) *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueSet("database_name", expected))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasSchemaName(expected string) *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueSet("schema_name", expected))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasOwner(expected string) *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueSet("owner", expected))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasComment(expected string) *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueSet("comment", expected))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasType(expected sdk.NetworkRuleType) *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputStringUnderlyingValueSet("type", expected))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasMode(expected sdk.NetworkRuleMode) *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputStringUnderlyingValueSet("mode", expected))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasValueList(expected []string) *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueSet("value_list.#", strconv.FormatInt(int64(len(expected)), 10)))
	for i, v := range expected {
		n.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("value_list.%d", i), v))
	}
	return n
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (n *NetworkRuleDescOutputAssert) HasNoCreatedOn() *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueNotSet("created_on"))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasNoName() *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueNotSet("name"))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasNoDatabaseName() *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueNotSet("database_name"))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasNoSchemaName() *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueNotSet("schema_name"))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasNoOwner() *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueNotSet("owner"))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasNoComment() *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueNotSet("comment"))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasNoType() *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputStringUnderlyingValueNotSet("type"))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasNoMode() *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputStringUnderlyingValueNotSet("mode"))
	return n
}

func (n *NetworkRuleDescOutputAssert) HasCommentEmpty() *NetworkRuleDescOutputAssert {
	n.AddAssertion(assert.ResourceDescribeOutputValueSet("comment", ""))
	return n
}

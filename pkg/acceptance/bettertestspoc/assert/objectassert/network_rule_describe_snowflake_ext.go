package objectassert

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type NetworkRuleDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.NetworkRuleDetails, sdk.SchemaObjectIdentifier]
}

func NetworkRuleDetails(t *testing.T, id sdk.SchemaObjectIdentifier) *NetworkRuleDetailsAssert {
	t.Helper()
	return &NetworkRuleDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(
			sdk.ObjectType("NETWORK_RULE_DETAILS"),
			id,
			func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.NetworkRuleDetails, sdk.SchemaObjectIdentifier] {
				return testClient.NetworkRule.Describe
			}),
	}
}

func (n *NetworkRuleDetailsAssert) HasCreatedOn(expected time.Time) *NetworkRuleDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NetworkRuleDetails) error {
		t.Helper()
		if o.CreatedOn != expected {
			return fmt.Errorf("expected created on: %v; got: %v", expected, o.CreatedOn)
		}
		return nil
	})
	return n
}

func (n *NetworkRuleDetailsAssert) HasCreatedOnNotEmpty() *NetworkRuleDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NetworkRuleDetails) error {
		t.Helper()
		if o.CreatedOn.IsZero() {
			return fmt.Errorf("expected created on to not be empty")
		}
		return nil
	})
	return n
}

func (n *NetworkRuleDetailsAssert) HasName(expected string) *NetworkRuleDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NetworkRuleDetails) error {
		t.Helper()
		if o.Name != expected {
			return fmt.Errorf("expected name: %v; got: %v", expected, o.Name)
		}
		return nil
	})
	return n
}

func (n *NetworkRuleDetailsAssert) HasDatabaseName(expected string) *NetworkRuleDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NetworkRuleDetails) error {
		t.Helper()
		if o.DatabaseName != expected {
			return fmt.Errorf("expected database name: %v; got: %v", expected, o.DatabaseName)
		}
		return nil
	})
	return n
}

func (n *NetworkRuleDetailsAssert) HasSchemaName(expected string) *NetworkRuleDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NetworkRuleDetails) error {
		t.Helper()
		if o.SchemaName != expected {
			return fmt.Errorf("expected schema name: %v; got: %v", expected, o.SchemaName)
		}
		return nil
	})
	return n
}

func (n *NetworkRuleDetailsAssert) HasOwner(expected string) *NetworkRuleDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NetworkRuleDetails) error {
		t.Helper()
		if o.Owner != expected {
			return fmt.Errorf("expected owner: %v; got: %v", expected, o.Owner)
		}
		return nil
	})
	return n
}

func (n *NetworkRuleDetailsAssert) HasComment(expected string) *NetworkRuleDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NetworkRuleDetails) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return n
}

func (n *NetworkRuleDetailsAssert) HasType(expected sdk.NetworkRuleType) *NetworkRuleDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NetworkRuleDetails) error {
		t.Helper()
		if o.Type != expected {
			return fmt.Errorf("expected type: %v; got: %v", expected, o.Type)
		}
		return nil
	})
	return n
}

func (n *NetworkRuleDetailsAssert) HasMode(expected sdk.NetworkRuleMode) *NetworkRuleDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NetworkRuleDetails) error {
		t.Helper()
		if o.Mode != expected {
			return fmt.Errorf("expected mode: %v; got: %v", expected, o.Mode)
		}
		return nil
	})
	return n
}

func (n *NetworkRuleDetailsAssert) HasValueList(expected []string) *NetworkRuleDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NetworkRuleDetails) error {
		t.Helper()
		if !reflect.DeepEqual(o.ValueList, expected) {
			return fmt.Errorf("expected value list: %v; got: %v", expected, o.ValueList)
		}
		return nil
	})
	return n
}

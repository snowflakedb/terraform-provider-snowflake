package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// HybridTableDetailsAssert wraps a slice of HybridTableDetails (one per column)
type HybridTableDetailsAssert struct {
	*assert.SnowflakeObjectAssert[[]sdk.HybridTableDetails, sdk.SchemaObjectIdentifier]
}

func HybridTableDetails(t *testing.T, id sdk.SchemaObjectIdentifier) *HybridTableDetailsAssert {
	t.Helper()
	return &HybridTableDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("HYBRID_TABLE_DETAILS"), id, func(testClient *helpers.TestClient) assert.ObjectProvider[[]sdk.HybridTableDetails, sdk.SchemaObjectIdentifier] {
			return testClient.HybridTable.DescribeDetails
		}),
	}
}

func HybridTableDetailsFromObject(t *testing.T, id sdk.SchemaObjectIdentifier, details []sdk.HybridTableDetails) *HybridTableDetailsAssert {
	t.Helper()
	return &HybridTableDetailsAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectType("HYBRID_TABLE_DETAILS"), id, &details),
	}
}

func (h *HybridTableDetailsAssert) HasColumnCount(expected int) *HybridTableDetailsAssert {
	h.AddAssertion(func(t *testing.T, details *[]sdk.HybridTableDetails) error {
		t.Helper()
		actual := len(*details)
		if actual != expected {
			return fmt.Errorf("expected column count: %d; got: %d", expected, actual)
		}
		return nil
	})
	return h
}

func (h *HybridTableDetailsAssert) HasColumn(name string, dataType string) *HybridTableDetailsAssert {
	h.AddAssertion(func(t *testing.T, details *[]sdk.HybridTableDetails) error {
		t.Helper()
		for _, col := range *details {
			if col.Name == name {
				if col.Type != dataType {
					return fmt.Errorf("column %s has type %s, expected %s", name, col.Type, dataType)
				}
				return nil
			}
		}
		return fmt.Errorf("column %s not found", name)
	})
	return h
}

func (h *HybridTableDetailsAssert) HasColumnWithKind(name string, kind string) *HybridTableDetailsAssert {
	h.AddAssertion(func(t *testing.T, details *[]sdk.HybridTableDetails) error {
		t.Helper()
		for _, col := range *details {
			if col.Name == name {
				if col.Kind != kind {
					return fmt.Errorf("column %s has kind %s, expected %s", name, col.Kind, kind)
				}
				return nil
			}
		}
		return fmt.Errorf("column %s not found", name)
	})
	return h
}

func (h *HybridTableDetailsAssert) HasPrimaryKey(columnName string) *HybridTableDetailsAssert {
	h.AddAssertion(func(t *testing.T, details *[]sdk.HybridTableDetails) error {
		t.Helper()
		for _, col := range *details {
			if col.Name == columnName {
				if col.PrimaryKey != "Y" {
					return fmt.Errorf("column %s is not a primary key (primary key = %s)", columnName, col.PrimaryKey)
				}
				return nil
			}
		}
		return fmt.Errorf("column %s not found", columnName)
	})
	return h
}

func (h *HybridTableDetailsAssert) HasNullableColumn(columnName string) *HybridTableDetailsAssert {
	h.AddAssertion(func(t *testing.T, details *[]sdk.HybridTableDetails) error {
		t.Helper()
		for _, col := range *details {
			if col.Name == columnName {
				if col.IsNullable != "Y" {
					return fmt.Errorf("column %s is not nullable (null? = %s)", columnName, col.IsNullable)
				}
				return nil
			}
		}
		return fmt.Errorf("column %s not found", columnName)
	})
	return h
}

func (h *HybridTableDetailsAssert) HasNotNullableColumn(columnName string) *HybridTableDetailsAssert {
	h.AddAssertion(func(t *testing.T, details *[]sdk.HybridTableDetails) error {
		t.Helper()
		for _, col := range *details {
			if col.Name == columnName {
				if col.IsNullable != "N" {
					return fmt.Errorf("column %s is nullable (null? = %s)", columnName, col.IsNullable)
				}
				return nil
			}
		}
		return fmt.Errorf("column %s not found", columnName)
	})
	return h
}

func (h *HybridTableDetailsAssert) HasColumnWithComment(columnName string, comment string) *HybridTableDetailsAssert {
	h.AddAssertion(func(t *testing.T, details *[]sdk.HybridTableDetails) error {
		t.Helper()
		for _, col := range *details {
			if col.Name == columnName {
				if col.Comment != comment {
					return fmt.Errorf("column %s has comment '%s', expected '%s'", columnName, col.Comment, comment)
				}
				return nil
			}
		}
		return fmt.Errorf("column %s not found", columnName)
	})
	return h
}

func (h *HybridTableDetailsAssert) HasColumnWithDefault(columnName string, defaultValue string) *HybridTableDetailsAssert {
	h.AddAssertion(func(t *testing.T, details *[]sdk.HybridTableDetails) error {
		t.Helper()
		for _, col := range *details {
			if col.Name == columnName {
				if col.Default != defaultValue {
					return fmt.Errorf("column %s has default '%s', expected '%s'", columnName, col.Default, defaultValue)
				}
				return nil
			}
		}
		return fmt.Errorf("column %s not found", columnName)
	})
	return h
}

func (h *HybridTableDetailsAssert) HasUniqueKey(columnName string) *HybridTableDetailsAssert {
	h.AddAssertion(func(t *testing.T, details *[]sdk.HybridTableDetails) error {
		t.Helper()
		for _, col := range *details {
			if col.Name == columnName {
				if col.UniqueKey != "Y" {
					return fmt.Errorf("column %s is not a unique key (unique key = %s)", columnName, col.UniqueKey)
				}
				return nil
			}
		}
		return fmt.Errorf("column %s not found", columnName)
	})
	return h
}

func (h *HybridTableDetailsAssert) HasColumnWithExpression(columnName string, expression string) *HybridTableDetailsAssert {
	h.AddAssertion(func(t *testing.T, details *[]sdk.HybridTableDetails) error {
		t.Helper()
		for _, col := range *details {
			if col.Name == columnName {
				if col.Expression != expression {
					return fmt.Errorf("column %s has expression '%s', expected '%s'", columnName, col.Expression, expression)
				}
				return nil
			}
		}
		return fmt.Errorf("column %s not found", columnName)
	})
	return h
}

func (h *HybridTableDetailsAssert) HasColumnWithPolicyName(columnName string, policyName string) *HybridTableDetailsAssert {
	h.AddAssertion(func(t *testing.T, details *[]sdk.HybridTableDetails) error {
		t.Helper()
		for _, col := range *details {
			if col.Name == columnName {
				if col.PolicyName != policyName {
					return fmt.Errorf("column %s has policy name '%s', expected '%s'", columnName, col.PolicyName, policyName)
				}
				return nil
			}
		}
		return fmt.Errorf("column %s not found", columnName)
	})
	return h
}

func (h *HybridTableDetailsAssert) HasColumnNames(expectedNames ...string) *HybridTableDetailsAssert {
	h.AddAssertion(func(t *testing.T, details *[]sdk.HybridTableDetails) error {
		t.Helper()
		if len(*details) != len(expectedNames) {
			return fmt.Errorf("expected %d columns, got %d", len(expectedNames), len(*details))
		}
		for i, expected := range expectedNames {
			if (*details)[i].Name != expected {
				return fmt.Errorf("column at index %d: expected name %s, got %s", i, expected, (*details)[i].Name)
			}
		}
		return nil
	})
	return h
}

func (h *HybridTableDetailsAssert) HasColumnAtIndex(index int, name string, dataType string) *HybridTableDetailsAssert {
	h.AddAssertion(func(t *testing.T, details *[]sdk.HybridTableDetails) error {
		t.Helper()
		if index >= len(*details) {
			return fmt.Errorf("index %d out of bounds (total columns: %d)", index, len(*details))
		}
		col := (*details)[index]
		if col.Name != name {
			return fmt.Errorf("column at index %d: expected name %s, got %s", index, name, col.Name)
		}
		if col.Type != dataType {
			return fmt.Errorf("column at index %d (%s): expected type %s, got %s", index, name, dataType, col.Type)
		}
		return nil
	})
	return h
}

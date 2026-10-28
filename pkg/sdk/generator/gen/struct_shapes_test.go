package gen

import (
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_requireConsistentStructShapes(t *testing.T) {
	policy := func(limitSql string) *QueryStruct {
		return NewQueryStruct("WidgetPolicy").OptionalNumberAssignment(limitSql, ParameterOptions())
	}
	interfaceWith := func(createPolicy, alterPolicy *QueryStruct) *Interface {
		i := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").
			CreateOperation("https://create", NewQueryStruct("CreateWidget").Create().SQL("WIDGET").Name().
				OptionalQueryStructField("Policy", createPolicy, ListOptions().SQL("POLICY =").Parentheses().NoComma())).
			AlterOperation("https://alter", NewQueryStruct("AlterWidget").Alter().SQL("WIDGET").Name().
				OptionalQueryStructField("Set", NewQueryStruct("WidgetSet").
					OptionalQueryStructField("Policy", alterPolicy, ListOptions().SQL("POLICY =").Parentheses().NoComma()), KeywordOptions().SQL("SET")))
		i.checkStructShapes = true
		return i
	}

	t.Run("same shape is shared", func(t *testing.T) {
		require.NotPanics(t, func() {
			preprocessDefinition(interfaceWith(policy("LIMIT"), policy("LIMIT")), log.New(io.Discard, "", 0))
		})
	})

	t.Run("different shape panics", func(t *testing.T) {
		require.PanicsWithValue(t, "Widget: struct WidgetPolicy is defined twice with a different shape:\nLimit *int ``ddl:\"parameter\" sql:\"LIMIT\"``\nOtherLimit *int ``ddl:\"parameter\" sql:\"OTHER_LIMIT\"``", func() {
			preprocessDefinition(interfaceWith(policy("LIMIT"), policy("OTHER_LIMIT")), log.New(io.Discard, "", 0))
		})
	})

	t.Run("not checked without the flag", func(t *testing.T) {
		i := interfaceWith(policy("LIMIT"), policy("OTHER_LIMIT"))
		i.checkStructShapes = false
		require.NotPanics(t, func() { preprocessDefinition(i, log.New(io.Discard, "", 0)) })
	})
}

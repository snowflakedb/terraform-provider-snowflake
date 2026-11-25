package gen

import (
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_toOptsOfSharedStructInSlice(t *testing.T) {
	constraint := func() *QueryStruct {
		return NewQueryStruct("WidgetConstraint").
			Text("Name", KeywordOptions()).
			OptionalQueryStructField("Inner", NewQueryStruct("WidgetInner").Text("Value", KeywordOptions()), KeywordOptions()).
			WithSharedToOpts()
	}
	// The constraint sits inside a column slice in CREATE and in a plain struct in ALTER, as iceberg tables use it.
	i := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").
		CreateOperation("https://create", NewQueryStruct("CreateWidget").Create().SQL("WIDGET").Name().
			ListQueryStructField("Columns", NewQueryStruct("WidgetColumn").
				Text("Name", KeywordOptions()).
				OptionalQueryStructField("Constraint", constraint(), KeywordOptions()), KeywordOptions())).
		AlterOperation("https://alter", NewQueryStruct("AlterWidget").Alter().SQL("WIDGET").Name().
			OptionalQueryStructField("AddColumn", NewQueryStruct("WidgetAdd").
				OptionalQueryStructField("Constraint", constraint(), KeywordOptions()), KeywordOptions()))
	preprocessDefinition(i, log.New(io.Discard, "", 0))

	t.Run("a struct used twice has one toOpts method", func(t *testing.T) {
		shared := i.SharedToOptsFields()
		require.Len(t, shared, 1)
		require.Equal(t, "WidgetConstraintRequest", shared[0].DtoDecl())
	})

	t.Run("the caller maps the struct inside its slice", func(t *testing.T) {
		columns := fieldNamed(t, i.Operations[0].OptsField, "Columns")
		inSlice := fieldNamed(t, &columns, "Constraint")
		require.Equal(t, "v", inSlice.ToOptsSourceRoot())
		require.Equal(t, "columns[i]", inSlice.ToOptsOptsRef())
	})

	t.Run("the standalone method does not use the slice of the caller", func(t *testing.T) {
		shared := i.SharedToOptsFields()[0]
		name := fieldNamed(t, shared, "Name")
		require.Equal(t, "r", name.ToOptsSourceRoot())
		require.Equal(t, "opts", name.ToOptsOptsRef())
		require.Equal(t, ".Name", name.ToOptsSourcePath())

		inner := fieldNamed(t, shared, "Inner")
		value := fieldNamed(t, &inner, "Value")
		require.Equal(t, "r", inner.ToOptsSourceRoot())
		require.Equal(t, "opts", inner.ToOptsOptsRef())
		require.Equal(t, ".Inner", inner.ToOptsSourcePath())
		require.Equal(t, "r", value.ToOptsSourceRoot())
		require.Equal(t, "opts", value.ToOptsOptsRef())
		require.Equal(t, ".Inner.Value", value.ToOptsSourcePath())
	})

	t.Run("the caller outside a slice uses the request", func(t *testing.T) {
		addColumn := fieldNamed(t, i.Operations[1].OptsField, "AddColumn")
		constraintField := fieldNamed(t, &addColumn, "Constraint")
		require.Equal(t, "r", constraintField.ToOptsSourceRoot())
		require.Equal(t, "opts", constraintField.ToOptsOptsRef())
	})
}

// fieldNamed returns the child field of parent with the given name. The returned pointer is the element of parent.Fields,
// which keeps its Parent link.
func fieldNamed(t *testing.T, parent *Field, name string) Field {
	t.Helper()
	for idx := range parent.Fields {
		if parent.Fields[idx].Name == name {
			return parent.Fields[idx]
		}
	}
	require.FailNowf(t, "field not found", "%s has no field %s", parent.Name, name)
	return Field{}
}

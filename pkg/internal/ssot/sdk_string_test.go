package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func sdkStringSchedule() Attribute {
	return Attribute{
		Name:           "schedule",
		Type:           AttrObject,
		SdkString:      true,
		SkipCreate:     true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Children: []Attribute{
			{Name: "hours", Type: AttrInt},
			{Name: "using_cron", Type: AttrString},
		},
	}
}

func Test_Attribute_SdkType(t *testing.T) {
	schedule := sdkStringSchedule()
	require.Equal(t, AttrString, schedule.SdkType())
	require.Equal(t, AttrObject, schedule.Type)
	require.Len(t, schedule.NestedFields(), 2, "the Terraform side keeps the block fields")

	require.Equal(t, AttrInt, Attribute{Type: AttrInt}.SdkType())
	require.Equal(t, AttrObject, Attribute{Type: AttrObject, Children: []Attribute{{Name: "a"}}}.SdkType())
}

func Test_SdkString_columnAndClauses(t *testing.T) {
	schedule := sdkStringSchedule()
	field := ColumnField(Object{Name: "Policy"}, schedule)
	require.Equal(t, "*string", field.ConcreteType)
	require.Equal(t, "Schedule", field.Name)
	require.Equal(t, StringParameter().Kind, schedule.EffectiveSqlClause(ClauseCreate).Kind)
}

func Test_ValidateObject_sdkString(t *testing.T) {
	t.Run("a block that is also a string column and clause", func(t *testing.T) {
		schedule := sdkStringSchedule()
		require.NotPanics(t, func() {
			ValidateObject(Object{
				Name:       "Policy",
				Attributes: []Attribute{schedule},
				Create:     CreateSpec{Clauses: []*Attribute{&schedule}},
				Set:        []*Attribute{&schedule},
				Show:       OutputSpec{Columns: []*Attribute{&schedule}},
			})
		})
	})

	t.Run("without SdkString the block is not a column", func(t *testing.T) {
		schedule := sdkStringSchedule()
		schedule.SdkString = false
		require.Panics(t, func() {
			ValidateObject(Object{
				Name:       "Policy",
				Attributes: []Attribute{schedule},
				Show:       OutputSpec{Columns: []*Attribute{&schedule}},
			})
		})
	})

	t.Run("must be an object", func(t *testing.T) {
		schedule := Attribute{Name: "schedule", Type: AttrString, SdkString: true, OptionalInShow: true}
		require.Panics(t, func() { ValidateObject(Object{Name: "Policy", Attributes: []Attribute{schedule}}) })
	})

	t.Run("must be optional", func(t *testing.T) {
		schedule := sdkStringSchedule()
		schedule.Required = true
		require.Panics(t, func() { ValidateObject(Object{Name: "Policy", Attributes: []Attribute{schedule}}) })
	})

	t.Run("needs the string column nullability", func(t *testing.T) {
		schedule := sdkStringSchedule()
		schedule.OptionalInShow = false
		schedule.OptionalInDb = false
		require.Panics(t, func() { ValidateObject(Object{Name: "Policy", Attributes: []Attribute{schedule}}) })
	})
}

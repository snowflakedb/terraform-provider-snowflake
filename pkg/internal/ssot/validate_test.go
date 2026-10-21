package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ValidateObject_writableMayBeOmittedFromColumns(t *testing.T) {
	minNodes := Attribute{Name: "min_nodes"}
	initiallySuspended := Attribute{Name: "initially_suspended"}
	name := Attribute{Name: "name", Readonly: true}
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:       "Pool",
			Attributes: []Attribute{minNodes, initiallySuspended},
			Show: OutputSpec{
				Columns: []*Attribute{&name, &minNodes},
			},
		})
	})
}

func Test_ValidateObject_readonlyInAttributes(t *testing.T) {
	require.Panics(t, func() {
		ValidateObject(Object{
			Name:       "Pool",
			Attributes: []Attribute{{Name: "name", Readonly: true}},
		})
	})
}

func Test_ValidateObject_writableColumnNotInAttributes(t *testing.T) {
	comment := Attribute{Name: "comment"}
	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Pool",
			Show: OutputSpec{Columns: []*Attribute{&comment}},
		})
	})
}

func Test_ValidateObject_conflictsWithUnknownAttribute(t *testing.T) {
	require.PanicsWithValue(t, `object "Warehouse": attribute "resource_constraint" conflicts with unknown attribute "generation"`, func() {
		ValidateObject(Object{
			Name: "Warehouse",
			Attributes: []Attribute{{
				Name:          "resource_constraint",
				ConflictsWith: []Attribute{{Name: "generation"}},
			}},
		})
	})
}

func Test_ValidateObject_conflictsWithItself(t *testing.T) {
	require.PanicsWithValue(t, `object "Warehouse": attribute "generation" conflicts with itself`, func() {
		ValidateObject(Object{
			Name: "Warehouse",
			Attributes: []Attribute{{
				Name:          "generation",
				ConflictsWith: []Attribute{{Name: "generation"}},
			}},
		})
	})
}

func Test_ValidateObject_conflictsWithKnownAttribute(t *testing.T) {
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name: "Warehouse",
			Attributes: []Attribute{
				{Name: "resource_constraint", ConflictsWith: []Attribute{{Name: "generation"}}},
				{Name: "generation", ConflictsWith: []Attribute{{Name: "resource_constraint"}}},
			},
		})
	})
}

func Test_ValidateObject_createOnlyMustBeOptionalBool(t *testing.T) {
	require.PanicsWithValue(t, `object "Warehouse": CreateOnly attribute "auto_suspend" must be an optional bool`, func() {
		ValidateObject(Object{
			Name:       "Warehouse",
			Attributes: []Attribute{{Name: "auto_suspend", Type: AttrInt, CreateOnly: true}},
		})
	})
	require.PanicsWithValue(t, `object "Warehouse": CreateOnly attribute "initially_suspended" must be an optional bool`, func() {
		ValidateObject(Object{
			Name:       "Warehouse",
			Attributes: []Attribute{{Name: "initially_suspended", Type: AttrBool, Required: true, CreateOnly: true}},
		})
	})
	require.PanicsWithValue(t, `object "Warehouse": NativeBool attribute "comment" must be a bool`, func() {
		ValidateObject(Object{
			Name:       "Warehouse",
			Attributes: []Attribute{{Name: "comment", Type: AttrString, NativeBool: true}},
		})
	})
}

func Test_ValidateObject_sqlClause(t *testing.T) {
	minNodes := Attribute{
		Name:      "min_nodes",
		Type:      AttrInt,
		CreateSql: IntParameter(),
		SetSql:    IntParameter(Optional()),
	}
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:       "Pool",
			Attributes: []Attribute{minNodes},
			Create:     []*Attribute{&minNodes},
			Set:        []*Attribute{&minNodes},
		})
	})

	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Pool",
			Attributes: []Attribute{{
				Name:      "created_on",
				Type:      AttrTime,
				CreateSql: IntParameter(),
			}},
		})
	})

	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Pool",
			Attributes: []Attribute{{
				Name:      "comment",
				Type:      AttrString,
				CreateSql: IdentifierKeyword("AccountObjectIdentifier"),
			}},
		})
	})

	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Pool",
			Attributes: []Attribute{{
				Name:      "instance_family",
				Type:      AttrEnum,
				CreateSql: EnumParameter(),
			}},
		})
	})

	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Pool",
			Attributes: []Attribute{{
				Name:      "backup",
				Type:      AttrList,
				CreateSql: ListParameter(""),
			}},
		})
	})

	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Pool",
			Attributes: []Attribute{{
				Name:      "for_application",
				Type:      AttrIdentifier,
				CreateSql: IdentifierKeyword(""),
			}},
		})
	})
}

func Test_ValidateObject_duplicateOutputKey(t *testing.T) {
	a := Attribute{Name: "for_application", ShowName: "application"}
	b := Attribute{Name: "application"}
	require.Panics(t, func() {
		ValidateObject(Object{
			Name:       "Pool",
			Attributes: []Attribute{a, b},
			Show:       OutputSpec{Columns: []*Attribute{&a, &b}},
		})
	})
}

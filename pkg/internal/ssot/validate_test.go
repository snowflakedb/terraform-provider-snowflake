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

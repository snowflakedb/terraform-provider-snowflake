package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Object_ShowOutputKeys_columnsAreMembershipAndOrder(t *testing.T) {
	minNodes := Attribute{Name: "min_nodes"}
	comment := Attribute{Name: "comment"}
	state := Attribute{Name: "state", Readonly: true}
	obj := Object{
		Attributes: []Attribute{minNodes, {Name: "initially_suspended"}, comment},
		Show: OutputSpec{
			Columns: []*Attribute{&state, &minNodes, &comment},
		},
	}
	require.Equal(t, []string{"state", "min_nodes", "comment"}, obj.ShowOutputKeys())
}

func Test_Object_DescribeOutputKeys_usesDescribeColumns(t *testing.T) {
	minNodes := Attribute{Name: "min_nodes"}
	state := Attribute{Name: "state", Readonly: true}
	errorCode := Attribute{Name: "error_code", Readonly: true}
	obj := Object{
		Attributes: []Attribute{minNodes},
		Show: OutputSpec{
			Columns: []*Attribute{&minNodes, &state},
		},
		Describe: OutputSpec{
			Columns: []*Attribute{&minNodes, &state, &errorCode},
		},
	}
	require.Equal(t, []string{"min_nodes", "state"}, obj.ShowOutputKeys())
	require.Equal(t, []string{"min_nodes", "state", "error_code"}, obj.DescribeOutputKeys())
}

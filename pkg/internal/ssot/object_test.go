package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AlterClause_BlockProjection(t *testing.T) {
	policy := Attribute{
		Name:      "policy",
		Type:      AttrObject,
		SqlOnly:   true,
		CreateSql: ObjectParameter(),
		SetSql:    ObjectParameter(FieldName("SetPolicy")),
		UnsetSql:  ObjectParameter(FieldName("UnsetPolicy")),
	}
	createOnly := Attribute{Name: "create_only", Type: AttrObject, SqlOnly: true, CreateSql: ObjectParameter()}

	testCases := []struct {
		name   string
		clause AlterClause
		want   ClauseUse
	}{
		{name: "set field", clause: AlterKeywordBlock("SetPolicy", "", &policy), want: ClauseSet},
		{name: "unset field", clause: AlterKeywordBlock("UnsetPolicy", "", &policy), want: ClauseUnset},
		{name: "other field", clause: AlterKeywordBlock("Policy", "", &policy), want: ClauseCreate},
		{name: "block without set or unset", clause: AlterKeywordBlock("SetPolicy", "", &createOnly), want: ClauseCreate},
		{name: "not a block", clause: AlterKeyword("RESUME"), want: ClauseCreate},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, tc.clause.BlockProjection())
		})
	}
}

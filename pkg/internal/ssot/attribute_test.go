package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Attribute_ChangesInPlace(t *testing.T) {
	tests := map[string]struct {
		attr Attribute
		want bool
	}{
		"plain":      {Attribute{Name: "comment"}, true},
		"ForceNew":   {Attribute{Name: "a", ForceNew: true}, false},
		"Readonly":   {Attribute{Name: "a", Readonly: true}, false},
		"CreateOnly": {Attribute{Name: "a", CreateOnly: true}, false},
		"Computed":   {Attribute{Name: "a", Computed: true}, false},
		"SqlOnly":    {Attribute{Name: "a", SqlOnly: true}, false},
		"SkipUpdate": {Attribute{Name: "a", SkipUpdate: true}, true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.attr.ChangesInPlace())
		})
	}
}

func Test_Attribute_CanUpdate(t *testing.T) {
	tests := map[string]struct {
		attr Attribute
		want bool
	}{
		"plain":      {Attribute{Name: "comment"}, true},
		"ForceNew":   {Attribute{Name: "a", ForceNew: true}, false},
		"Readonly":   {Attribute{Name: "a", Readonly: true}, false},
		"CreateOnly": {Attribute{Name: "a", CreateOnly: true}, false},
		"Computed":   {Attribute{Name: "a", Computed: true}, false},
		"SqlOnly":    {Attribute{Name: "a", SqlOnly: true}, false},
		// An attribute that a hook updates changes in place, but the generated ApplySetUnset leaves it out.
		"SkipUpdate": {Attribute{Name: "a", SkipUpdate: true}, false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.attr.CanUpdate())
		})
	}
}

func Test_Attribute_ClauseMarkerName(t *testing.T) {
	require.Equal(t, "tags", Tags.ClauseMarkerName())
	require.Equal(t, "parameters", Parameters.ClauseMarkerName())
	require.Empty(t, (&Attribute{Name: "comment", CreateSql: StringParameter()}).ClauseMarkerName())
	require.Empty(t, (*Attribute)(nil).ClauseMarkerName())
}

func Test_Attribute_NestedTagsClauseFollowsItsStatement(t *testing.T) {
	block := Attribute{
		Name: "alter_tags",
		Type: AttrObject,
		Children: []Attribute{
			{Name: "tags", CreateSql: TagsClause()},
			{Name: "set_tags", SetSql: TagsClause()},
			{Name: "unset_tags", UnsetSql: TagsClause()},
		},
	}
	names := func(use ClauseUse) []string {
		fields := block.NestedSqlFieldsFor(use)
		out := make([]string, 0, len(fields))
		for _, child := range fields {
			out = append(out, child.Name)
		}
		return out
	}
	require.Equal(t, []string{"tags"}, names(ClauseCreate))
	require.Equal(t, []string{"set_tags"}, names(ClauseSet))
	require.Equal(t, []string{"unset_tags"}, names(ClauseUnset))
}

func Test_validateSqlClause_parametersClauseIsNotANestedField(t *testing.T) {
	obj := Object{Name: "Widget"}
	require.NotPanics(t, func() {
		validateSqlClause(obj, Attribute{Name: "tags"}, "CREATE", TagsClause())
	})
	require.PanicsWithValue(t, `object "Widget": attribute "parameters" CREATE parameters clause is a clause list entry, not a nested field`, func() {
		validateSqlClause(obj, Attribute{Name: "parameters"}, "CREATE", ParametersClause())
	})
}

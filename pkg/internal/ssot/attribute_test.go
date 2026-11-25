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

package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToBehaviorChangeBundleStatus(t *testing.T) {
	type test struct {
		input string
		want  BehaviorChangeBundleStatus
	}

	valid := []test{
		{input: "enabled", want: BehaviorChangeBundleStatusEnabled},
		{input: "ENABLED", want: BehaviorChangeBundleStatusEnabled},
		{input: "DISABLED", want: BehaviorChangeBundleStatusDisabled},
		{input: "RELEASED", want: BehaviorChangeBundleStatusReleased},
	}

	invalid := []string{
		"",
		"foo",
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToBehaviorChangeBundleStatus(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, in := range invalid {
		t.Run(in, func(t *testing.T) {
			_, err := ToBehaviorChangeBundleStatus(in)
			require.Error(t, err)
		})
	}
}

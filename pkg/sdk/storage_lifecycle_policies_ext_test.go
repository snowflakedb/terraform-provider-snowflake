package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_normalizeStorageLifecyclePolicyArchiveTier(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected string
	}{
		{
			Name:     "NULL literal is normalized to empty string",
			Input:    "NULL",
			Expected: "",
		},
		{
			Name:     "lowercase null is normalized to empty string",
			Input:    "null",
			Expected: "",
		},
		{
			Name:     "empty string stays empty",
			Input:    "",
			Expected: "",
		},
		{
			Name:     "COLD tier is returned as is",
			Input:    "COLD",
			Expected: "COLD",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, normalizeStorageLifecyclePolicyArchiveTier(tc.Input))
		})
	}
}

package experimentalfeatures_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/stretchr/testify/require"
)

func Test_IsExperimentEnabled(t *testing.T) {
	type test struct {
		input       experimentalfeatures.ExperimentalFeature
		enabledList []string
		expected    bool
	}

	feature := experimentalfeatures.WarehouseShowImprovedPerformance
	lowercaseFeature := experimentalfeatures.ExperimentalFeature(strings.ToLower(string(experimentalfeatures.WarehouseShowImprovedPerformance)))

	listWithFeature := []string{string(feature)}
	listWithFeatureLowercase := []string{string(lowercaseFeature)}
	listWithOtherFeature := []string{"other"}

	valid := []test{
		{input: feature, enabledList: nil, expected: false},
		{input: feature, enabledList: []string{}, expected: false},
		{input: feature, enabledList: listWithOtherFeature, expected: false},
		{input: feature, enabledList: listWithFeature, expected: true},
		{input: feature, enabledList: listWithFeatureLowercase, expected: true},
		{input: lowercaseFeature, enabledList: listWithFeature, expected: true},
		{input: lowercaseFeature, enabledList: listWithFeatureLowercase, expected: true},
	}

	for _, tc := range valid {
		t.Run(fmt.Sprintf("List: %v, feature: %s", tc.enabledList, tc.input), func(t *testing.T) {
			got := experimentalfeatures.IsExperimentEnabled(tc.input, tc.enabledList)
			require.Equal(t, tc.expected, got)
		})
	}
}

package experimentalfeatures_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/stretchr/testify/require"
)

func Test_Experiments_IsEnabled(t *testing.T) {
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
			got := experimentalfeatures.New(tc.enabledList).IsEnabled(tc.input)
			require.Equal(t, tc.expected, got)
		})
	}

	t.Run("zero value is disabled", func(t *testing.T) {
		require.False(t, experimentalfeatures.Experiments{}.IsEnabled(feature))
	})
}

func Test_Experiments_UserEnabled(t *testing.T) {
	enabled := []string{string(experimentalfeatures.HierarchyRenames)}
	experiments := experimentalfeatures.New(enabled)
	got := experiments.UserEnabled()
	require.Equal(t, enabled, got)
	got[0] = "mutated"
	require.Equal(t, enabled, experiments.UserEnabled())
}

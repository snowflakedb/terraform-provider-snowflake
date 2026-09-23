package experimentalfeatures

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Experiments_IsEnabled(t *testing.T) {
	type test struct {
		input       ExperimentalFeature
		enabledList []string
		expected    bool
	}

	feature := WarehouseShowImprovedPerformance
	lowercaseFeature := ExperimentalFeature(strings.ToLower(string(WarehouseShowImprovedPerformance)))

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
			got := New(tc.enabledList, nil).IsEnabled(tc.input)
			require.Equal(t, tc.expected, got)
		})
	}

	t.Run("zero value is disabled", func(t *testing.T) {
		require.False(t, Experiments{}.IsEnabled(feature))
	})

	t.Run("default-on experiment is enabled without being listed", func(t *testing.T) {
		require.True(t, New(nil, nil).IsEnabled(InheritedGrants))
	})

	t.Run("default-on experiment can be disabled", func(t *testing.T) {
		require.False(t, New(nil, []string{string(InheritedGrants)}).IsEnabled(InheritedGrants))
	})

	t.Run("disabled wins over enabled for default-on experiment", func(t *testing.T) {
		name := string(InheritedGrants)
		require.False(t, New([]string{name}, []string{name}).IsEnabled(InheritedGrants))
	})
}

func Test_Experiments_UserLists(t *testing.T) {
	enabled := []string{string(HierarchyRenames)}
	disabled := []string{string(InheritedGrants)}
	experiments := New(enabled, disabled)

	gotEnabled := experiments.UserEnabled()
	require.Equal(t, enabled, gotEnabled)
	gotEnabled[0] = "mutated"
	require.Equal(t, enabled, experiments.UserEnabled())

	gotDisabled := experiments.UserDisabled()
	require.Equal(t, disabled, gotDisabled)
	gotDisabled[0] = "mutated"
	require.Equal(t, disabled, experiments.UserDisabled())

	gotEffective := experiments.EffectiveEnabled()
	require.Equal(t, enabled, gotEffective)
	gotEffective[0] = "mutated"
	require.Equal(t, enabled, experiments.EffectiveEnabled())
	require.NotContains(t, experiments.EffectiveEnabled(), string(InheritedGrants))
}

func Test_ConstructorsAssignState(t *testing.T) {
	active := NewActiveExperiment("ACTIVE_EXP", "description")
	require.Equal(t, ExperimentalFeature("ACTIVE_EXP"), active.Name())
	require.Equal(t, ExperimentalFeatureStateActive, active.state)
	require.Equal(t, "description", active.Description())
	require.Empty(t, active.removeInVersions)
	require.Empty(t, active.RemoveInVersionsPhrase())

	defaultOn := NewEnabledByDefaultExperiment("DEFAULT_EXP", []string{"v2.23.0", "v2.24.0"}, "part one", "part two")
	require.Equal(t, ExperimentalFeature("DEFAULT_EXP"), defaultOn.Name())
	require.Equal(t, ExperimentalFeatureStateEnabledByDefault, defaultOn.state)
	require.Equal(t, "part one\n\npart two\n\nIt will be promoted in v2.23.0 or v2.24.0.", defaultOn.Description())
	require.Equal(t, []string{"v2.23.0", "v2.24.0"}, defaultOn.removeInVersions)
	require.Equal(t, "v2.23.0 or v2.24.0", defaultOn.RemoveInVersionsPhrase())

	defaultOnNoVersions := NewEnabledByDefaultExperiment("DEFAULT_EXP_NO_VERSIONS", nil, "only body")
	require.Equal(t, "only body", defaultOnNoVersions.Description())

	oneVersion := NewEnabledByDefaultExperiment("ONE", []string{"v2.23.0"}, "body")
	require.Equal(t, "body\n\nIt will be promoted in v2.23.0.", oneVersion.Description())
	require.Equal(t, "v2.23.0", oneVersion.RemoveInVersionsPhrase())

	threeVersions := NewEnabledByDefaultExperiment("THREE", []string{"v2.23.0", "v2.24.0", "v2.25.0"}, "body")
	require.Equal(t, "body\n\nIt will be promoted in v2.23.0, v2.24.0, or v2.25.0.", threeVersions.Description())
	require.Equal(t, "v2.23.0, v2.24.0, or v2.25.0", threeVersions.RemoveInVersionsPhrase())

	promoted := NewPromotedExperiment("PROMOTED_EXP")
	require.Equal(t, ExperimentalFeature("PROMOTED_EXP"), promoted.Name())
	require.Equal(t, ExperimentalFeatureStatePromoted, promoted.state)
	require.Empty(t, promoted.Description())
	require.Empty(t, promoted.removeInVersions)

	discontinued := NewDiscontinuedExperiment("DISCONTINUED_EXP")
	require.Equal(t, ExperimentalFeature("DISCONTINUED_EXP"), discontinued.Name())
	require.Equal(t, ExperimentalFeatureStateDiscontinued, discontinued.state)
	require.Empty(t, discontinued.Description())
	require.Empty(t, discontinued.removeInVersions)
}

func Test_EnabledByDefaultExperimentsHaveRemoveInVersions(t *testing.T) {
	require.NotEmpty(t, EnabledByDefaultExperiments)
	for _, experiment := range EnabledByDefaultExperiments {
		require.NotEmptyf(t, experiment.removeInVersions, "ENABLED_BY_DEFAULT experiment %s must have removeInVersions", experiment.Name())
	}
}

func Test_ComputeEffectiveExperimentsEnabled_Public(t *testing.T) {
	active := string(WarehouseShowImprovedPerformance)
	defaultOn := string(InheritedGrants)

	t.Run("empty lists enable default-on experiments only", func(t *testing.T) {
		got := ComputeEffectiveExperimentsEnabled(nil, nil)
		require.Equal(t, []string{defaultOn}, got)
	})

	t.Run("active experiment is included when listed in enabled", func(t *testing.T) {
		got := ComputeEffectiveExperimentsEnabled([]string{active}, nil)
		require.Equal(t, []string{active, defaultOn}, got)
	})

	t.Run("default-on experiment is excluded when listed in disabled", func(t *testing.T) {
		got := ComputeEffectiveExperimentsEnabled(nil, []string{defaultOn})
		require.Empty(t, got)
	})

	t.Run("overlap: disabled wins for default-on experiment", func(t *testing.T) {
		got := ComputeEffectiveExperimentsEnabled([]string{defaultOn}, []string{defaultOn})
		require.Empty(t, got)
	})

	t.Run("matching is case-insensitive", func(t *testing.T) {
		got := ComputeEffectiveExperimentsEnabled([]string{strings.ToLower(active)}, []string{strings.ToLower(defaultOn)})
		require.Equal(t, []string{active}, got)
	})
}

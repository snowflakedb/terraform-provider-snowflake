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

	feature := HierarchyRenames
	lowercaseFeature := ExperimentalFeature(strings.ToLower(string(HierarchyRenames)))

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
		require.True(t, New(nil, nil).IsEnabled(WarehouseShowImprovedPerformance))
		require.True(t, New(nil, nil).IsEnabled(GrantAccountRoleSafePublicRole))
		require.True(t, New(nil, nil).IsEnabled(GrantsStrictPrivilegeManagement))
		require.True(t, New(nil, nil).IsEnabled(ImportBooleanDefault))
		require.True(t, New(nil, nil).IsEnabled(GrantsImportValidation))
		require.True(t, New(nil, nil).IsEnabled(ParametersIgnoreValueChangesIfNotOnObjectLevel))
		require.True(t, New(nil, nil).IsEnabled(GrantsSafeDestroy))
		require.True(t, New(nil, nil).IsEnabled(TagAssociationSafeDestroy))
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
	disabled := []string{string(InheritedGrants), string(WarehouseShowImprovedPerformance), string(GrantAccountRoleSafePublicRole), string(GrantsStrictPrivilegeManagement), string(ImportBooleanDefault), string(GrantsImportValidation), string(ParametersIgnoreValueChangesIfNotOnObjectLevel), string(GrantsSafeDestroy), string(TagAssociationSafeDestroy)}
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
	optIn := NewOptInExperiment("OPT_IN_EXP", "description")
	require.Equal(t, ExperimentalFeature("OPT_IN_EXP"), optIn.Name())
	require.Equal(t, ExperimentalFeatureStateOptIn, optIn.state)
	require.Equal(t, "description", optIn.Description())
	require.Empty(t, optIn.removeInVersions)
	require.Empty(t, optIn.RemoveInVersionsPhrase())

	defaultOn := NewEnabledByDefaultExperiment("DEFAULT_EXP", []string{"v2.23.0", "v2.24.0"}, "part one", "part two")
	require.Equal(t, ExperimentalFeature("DEFAULT_EXP"), defaultOn.Name())
	require.Equal(t, ExperimentalFeatureStateEnabledByDefault, defaultOn.state)
	require.Equal(t, "part one\n\npart two\n\n"+enabledByDefaultOptOutDescription+"\n\nIt will be promoted in v2.23.0 or v2.24.0.", defaultOn.Description())
	require.Equal(t, []string{"v2.23.0", "v2.24.0"}, defaultOn.removeInVersions)
	require.Equal(t, "v2.23.0 or v2.24.0", defaultOn.RemoveInVersionsPhrase())

	defaultOnNoVersions := NewEnabledByDefaultExperiment("DEFAULT_EXP_NO_VERSIONS", nil, "only body")
	require.Equal(t, "only body\n\n"+enabledByDefaultOptOutDescription, defaultOnNoVersions.Description())

	oneVersion := NewEnabledByDefaultExperiment("ONE", []string{"v2.23.0"}, "body")
	require.Equal(t, "body\n\n"+enabledByDefaultOptOutDescription+"\n\nIt will be promoted in v2.23.0.", oneVersion.Description())
	require.Equal(t, "v2.23.0", oneVersion.RemoveInVersionsPhrase())

	threeVersions := NewEnabledByDefaultExperiment("THREE", []string{"v2.23.0", "v2.24.0", "v2.25.0"}, "body")
	require.Equal(t, "body\n\n"+enabledByDefaultOptOutDescription+"\n\nIt will be promoted in v2.23.0, v2.24.0, or v2.25.0.", threeVersions.Description())
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
	optIn := string(HierarchyRenames)
	warehouse := string(WarehouseShowImprovedPerformance)
	grantsStrict := string(GrantsStrictPrivilegeManagement)
	parametersIgnore := string(ParametersIgnoreValueChangesIfNotOnObjectLevel)
	grantsImport := string(GrantsImportValidation)
	importBoolean := string(ImportBooleanDefault)
	grantsSafeDestroy := string(GrantsSafeDestroy)
	tagAssociationSafeDestroy := string(TagAssociationSafeDestroy)
	safePublic := string(GrantAccountRoleSafePublicRole)
	inherited := string(InheritedGrants)
	allDefaultOn := []string{warehouse, grantsStrict, parametersIgnore, grantsImport, importBoolean, grantsSafeDestroy, tagAssociationSafeDestroy, safePublic, inherited}

	t.Run("empty lists enable default-on experiments only", func(t *testing.T) {
		got := ComputeEffectiveExperimentsEnabled(nil, nil)
		require.Equal(t, allDefaultOn, got)
	})

	t.Run("opt-in experiment is included when listed in enabled", func(t *testing.T) {
		got := ComputeEffectiveExperimentsEnabled([]string{optIn}, nil)
		require.Equal(t, []string{warehouse, grantsStrict, parametersIgnore, grantsImport, importBoolean, grantsSafeDestroy, tagAssociationSafeDestroy, safePublic, optIn, inherited}, got)
	})

	t.Run("default-on experiment is excluded when listed in disabled", func(t *testing.T) {
		got := ComputeEffectiveExperimentsEnabled(nil, allDefaultOn)
		require.Empty(t, got)
	})

	t.Run("overlap: disabled wins for default-on experiment", func(t *testing.T) {
		got := ComputeEffectiveExperimentsEnabled(allDefaultOn, allDefaultOn)
		require.Empty(t, got)
	})

	t.Run("matching is case-insensitive", func(t *testing.T) {
		got := ComputeEffectiveExperimentsEnabled([]string{strings.ToLower(optIn)}, []string{strings.ToLower(inherited)})
		require.Equal(t, []string{warehouse, grantsStrict, parametersIgnore, grantsImport, importBoolean, grantsSafeDestroy, tagAssociationSafeDestroy, safePublic, optIn}, got)
	})
}

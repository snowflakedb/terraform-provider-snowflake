package experimentalfeatures

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_experimentListingWarnings(t *testing.T) {
	experiments := testLifecycleExperiments()

	tests := []struct {
		name             string
		userEnabled      []string
		userDisabled     []string
		wantSummaries    []string
		wantDetailSubstr []string
	}{
		{
			name: "empty lists produce no warnings",
		},
		{
			name:        "opt-in in enabled produces no warnings",
			userEnabled: []string{"OPT_IN_EXP"},
		},
		{
			name:          "default-on in enabled is redundant",
			userEnabled:   []string{"DEFAULT_EXP"},
			wantSummaries: []string{"Experiment already enabled by default."},
		},
		{
			name:          "default-on overlap: disabled wins then redundant",
			userEnabled:   []string{"DEFAULT_EXP"},
			userDisabled:  []string{"DEFAULT_EXP"},
			wantSummaries: []string{"Experiment listed in both enabled and disabled lists.", "Experiment already enabled by default."},
		},
		{
			name:             "promoted in enabled",
			userEnabled:      []string{"PROMOTED_EXP"},
			wantSummaries:    []string{"Promoted experiment listed in configuration."},
			wantDetailSubstr: []string{"PROMOTED_EXP"},
		},
		{
			name:             "promoted in disabled",
			userDisabled:     []string{"PROMOTED_EXP"},
			wantSummaries:    []string{"Promoted experiment listed in configuration."},
			wantDetailSubstr: []string{"PROMOTED_EXP"},
		},
		{
			name:          "promoted in both lists: overlap then promoted",
			userEnabled:   []string{"PROMOTED_EXP"},
			userDisabled:  []string{"PROMOTED_EXP"},
			wantSummaries: []string{"Experiment listed in both enabled and disabled lists.", "Promoted experiment listed in configuration."},
		},
		{
			name:             "discontinued in enabled",
			userEnabled:      []string{"DISCONTINUED_EXP"},
			wantSummaries:    []string{"Discontinued experiment listed in configuration."},
			wantDetailSubstr: []string{"next major version"},
		},
		{
			name:          "discontinued in disabled",
			userDisabled:  []string{"DISCONTINUED_EXP"},
			wantSummaries: []string{"Discontinued experiment listed in configuration."},
		},
		{
			name:          "discontinued in both lists: overlap then discontinued",
			userEnabled:   []string{"DISCONTINUED_EXP"},
			userDisabled:  []string{"DISCONTINUED_EXP"},
			wantSummaries: []string{"Experiment listed in both enabled and disabled lists.", "Discontinued experiment listed in configuration."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := experimentListingWarnings(experiments, tt.userEnabled, tt.userDisabled)
			if len(tt.wantSummaries) == 0 {
				require.Empty(t, got)
				return
			}
			require.Len(t, got, len(tt.wantSummaries))
			for i, summary := range tt.wantSummaries {
				require.Equal(t, summary, got[i].Summary)
			}
			for i, substr := range tt.wantDetailSubstr {
				require.Contains(t, got[i].Detail, substr)
			}
		})
	}
}

func Test_ExperimentListingWarnings_Public(t *testing.T) {
	optIn := string(HierarchyRenames)
	defaultOn := string(InheritedGrants)
	warehouse := string(WarehouseShowImprovedPerformance)

	t.Run("no warnings for empty lists", func(t *testing.T) {
		require.Empty(t, ExperimentListingWarnings(nil, nil))
	})

	t.Run("no warnings for enabling an opt-in experiment", func(t *testing.T) {
		require.Empty(t, ExperimentListingWarnings([]string{optIn}, nil))
	})

	t.Run("default-on warehouse listed in enabled is redundant", func(t *testing.T) {
		got := ExperimentListingWarnings([]string{warehouse}, nil)
		require.Len(t, got, 1)
		require.Equal(t, "Experiment already enabled by default.", got[0].Summary)
		require.Contains(t, got[0].Detail, warehouse)
	})

	t.Run("default-on listed in enabled is redundant", func(t *testing.T) {
		got := ExperimentListingWarnings([]string{defaultOn}, nil)
		require.Len(t, got, 1)
		require.Equal(t, "Experiment already enabled by default.", got[0].Summary)
		require.Contains(t, got[0].Detail, defaultOn)
	})

	t.Run("same name in both lists: disabled wins warning plus redundant for default-on", func(t *testing.T) {
		got := ExperimentListingWarnings([]string{defaultOn}, []string{defaultOn})
		require.Len(t, got, 2)
		require.Equal(t, "Experiment listed in both enabled and disabled lists.", got[0].Summary)
		require.Equal(t, "Experiment already enabled by default.", got[1].Summary)
	})
}

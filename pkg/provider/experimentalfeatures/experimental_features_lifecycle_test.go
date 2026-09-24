package experimentalfeatures

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testLifecycleExperiments() []Experiment {
	return []Experiment{
		NewOptInExperiment("OPT_IN_EXP", "opt-in"),
		NewEnabledByDefaultExperiment("DEFAULT_EXP", []string{"v1.0.0"}, "default on"),
		NewPromotedExperiment("PROMOTED_EXP"),
		NewDiscontinuedExperiment("DISCONTINUED_EXP"),
	}
}

func Test_computeEffectiveExperimentsEnabled(t *testing.T) {
	experiments := testLifecycleExperiments()

	tests := []struct {
		name         string
		userEnabled  []string
		userDisabled []string
		expected     []string
	}{
		{
			name:     "empty lists: only default-on",
			expected: []string{"DEFAULT_EXP"},
		},
		{
			name:        "opt-in in enabled",
			userEnabled: []string{"OPT_IN_EXP"},
			expected:    []string{"OPT_IN_EXP", "DEFAULT_EXP"},
		},
		{
			name:         "default-on in disabled",
			userDisabled: []string{"DEFAULT_EXP"},
			expected:     []string{},
		},
		{
			name:        "promoted in enabled is never effective",
			userEnabled: []string{"PROMOTED_EXP"},
			expected:    []string{"DEFAULT_EXP"},
		},
		{
			name:         "promoted in disabled is never effective",
			userDisabled: []string{"PROMOTED_EXP"},
			expected:     []string{"DEFAULT_EXP"},
		},
		{
			name:        "discontinued in enabled is never effective",
			userEnabled: []string{"DISCONTINUED_EXP"},
			expected:    []string{"DEFAULT_EXP"},
		},
		{
			name:         "discontinued in disabled is never effective",
			userDisabled: []string{"DISCONTINUED_EXP"},
			expected:     []string{"DEFAULT_EXP"},
		},
		{
			name:         "overlap on default-on: disabled wins",
			userEnabled:  []string{"DEFAULT_EXP"},
			userDisabled: []string{"DEFAULT_EXP"},
			expected:     []string{},
		},
		{
			name:         "overlap on opt-in: disabled wins",
			userEnabled:  []string{"OPT_IN_EXP"},
			userDisabled: []string{"OPT_IN_EXP"},
			expected:     []string{"DEFAULT_EXP"},
		},
		{
			name:         "overlap on promoted: never effective",
			userEnabled:  []string{"PROMOTED_EXP"},
			userDisabled: []string{"PROMOTED_EXP"},
			expected:     []string{"DEFAULT_EXP"},
		},
		{
			name:         "overlap on discontinued: never effective",
			userEnabled:  []string{"DISCONTINUED_EXP"},
			userDisabled: []string{"DISCONTINUED_EXP"},
			expected:     []string{"DEFAULT_EXP"},
		},
		{
			name:         "case-insensitive membership",
			userEnabled:  []string{"opt_in_exp"},
			userDisabled: []string{"default_exp"},
			expected:     []string{"OPT_IN_EXP"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeEffectiveExperimentsEnabled(experiments, tt.userEnabled, tt.userDisabled)
			if len(tt.expected) == 0 {
				require.Empty(t, got)
				return
			}
			require.Equal(t, tt.expected, got)
		})
	}
}

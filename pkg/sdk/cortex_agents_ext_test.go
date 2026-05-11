package sdk

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalCortexAgentProfile(t *testing.T) {
	validCases := []struct {
		name     string
		json     string
		expected *CortexAgentProfile
	}{
		{
			name: "all fields present",
			json: `{"display_name":"My Assistant","avatar":"assistant.png","color":"blue"}`,
			expected: &CortexAgentProfile{
				DisplayName: String("My Assistant"),
				Avatar:      String("assistant.png"),
				Color:       String("blue"),
			},
		},
		{
			name: "single field present",
			json: `{"color":"green"}`,
			expected: &CortexAgentProfile{
				Color: String("green"),
			},
		},
		{
			name:     "empty object",
			json:     `{}`,
			expected: &CortexAgentProfile{},
		},
	}

	for _, tc := range validCases {
		t.Run(tc.name, func(t *testing.T) {
			profile, err := UnmarshalCortexAgentProfile(tc.json)
			require.NoError(t, err)
			require.NotNil(t, profile)
			require.True(t, reflect.DeepEqual(tc.expected, profile), "expected %#v; got %#v", tc.expected, profile)
		})
	}

	invalidProfiles := []string{
		`{"broken"`,
		"",
	}

	for _, profile := range invalidProfiles {
		t.Run(profile, func(t *testing.T) {
			p, err := UnmarshalCortexAgentProfile(profile)
			require.Error(t, err)
			require.Nil(t, p)
		})
	}
}

func TestUnmarshalCortexAgentSpec(t *testing.T) {
	validCases := []struct {
		name     string
		json     string
		expected map[string]any
	}{
		{
			name:     "empty string",
			json:     "",
			expected: map[string]any{},
		},
		{
			name:     "whitespace only",
			json:     "  \n\t  ",
			expected: map[string]any{},
		},
		{
			name: "nested object",
			json: `{"orchestration":{"budget":{"seconds":30,"tokens":16000}}}`,
			expected: map[string]any{
				"orchestration": map[string]any{
					"budget": map[string]any{
						"seconds": float64(30),
						"tokens":  float64(16000),
					},
				},
			},
		},
	}

	for _, tc := range validCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := UnmarshalCortexAgentSpec(tc.json)
			require.NoError(t, err)
			require.NotNil(t, m)
			require.True(t, reflect.DeepEqual(tc.expected, m), "expected %#v; got %#v", tc.expected, m)
		})
	}

	invalidSpecs := []string{
		`{"broken"`,
		"[1,2]",
	}

	for _, spec := range invalidSpecs {
		t.Run(spec, func(t *testing.T) {
			m, err := UnmarshalCortexAgentSpec(spec)
			require.Error(t, err)
			require.Nil(t, m)
		})
	}
}

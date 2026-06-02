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
		expected CortexAgentProfile
	}{
		{
			name: "all fields present",
			json: `{"display_name":"My Assistant","avatar":"assistant.png","color":"blue"}`,
			expected: CortexAgentProfile{
				DisplayName: String("My Assistant"),
				Avatar:      String("assistant.png"),
				Color:       String("blue"),
			},
		},
		{
			name: "single field present",
			json: `{"color":"green"}`,
			expected: CortexAgentProfile{
				Color: String("green"),
			},
		},
		{
			name:     "empty object",
			json:     `{}`,
			expected: CortexAgentProfile{},
		},
	}

	for _, tc := range validCases {
		t.Run(tc.name, func(t *testing.T) {
			profile, err := UnmarshalCortexAgentProfile(tc.json)
			require.NoError(t, err)
			require.True(t, reflect.DeepEqual(tc.expected, profile), "expected %#v; got %#v", tc.expected, profile)
		})
	}

	invalidProfiles := []string{
		`{"broken"`,
		"",
		`[{"color":"blue"}]`,
	}

	for _, profile := range invalidProfiles {
		t.Run(profile, func(t *testing.T) {
			p, err := UnmarshalCortexAgentProfile(profile)
			require.Error(t, err)
			require.Equal(t, CortexAgentProfile{}, p)
		})
	}
}

func TestMarshalCortexAgentProfile(t *testing.T) {
	validCases := []struct {
		name    string
		profile CortexAgentProfile
		want    string
	}{
		{
			name:    "empty profile",
			profile: CortexAgentProfile{},
			want:    `{}`,
		},
		{
			name: "full profile",
			profile: CortexAgentProfile{
				DisplayName: String("My Assistant"),
				Avatar:      String("assistant.png"),
				Color:       String("blue"),
			},
			want: `{"display_name":"My Assistant","avatar":"assistant.png","color":"blue"}`,
		},
		{
			name: "partial profile",
			profile: CortexAgentProfile{
				Color: String("green"),
			},
			want: `{"color":"green"}`,
		},
	}

	for _, tc := range validCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := MarshalCortexAgentProfile(tc.profile)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestNormalizeCortexAgentSpecification(t *testing.T) {
	yamlAgentSpec := `orchestration:
  budget:
    seconds: 30
    tokens: 16000
instructions:
  response: "Basic acceptance tests"
`
	want, err := NormalizeCortexAgentSpecification(yamlAgentSpec)
	require.NoError(t, err)

	t.Run("equivalent agent specifications", func(t *testing.T) {
		equivalentAgentSpecifications := []string{
			`{"instructions":{"response":"Basic acceptance tests"},"orchestration":{"budget":{"seconds":30,"tokens":16000}}}`,
			`{"orchestration":{"budget":{"seconds":30,"tokens":16000}},"instructions":{"response":"Basic acceptance tests"}}`,
			`{"orchestration":{"budget":{"tokens":16000,"seconds":30}},"instructions":{"response":"Basic acceptance tests"}}`,
			`{  "instructions"  :  {  "response"  :  "Basic acceptance tests"  }  ,  "orchestration"  :  {  "budget"  :  {  "seconds"  :  30  ,  "tokens"  :  16000  }  }  }`,
			"{\n  \"instructions\": {\n    \"response\": \"Basic acceptance tests\"\n  },\n  \"orchestration\": {\n    \"budget\": {\n      \"seconds\": 30,\n      \"tokens\": 16000\n    }\n  }\n}",
			`{"orchestration" : { "budget" : { "seconds" : 30 , "tokens" : 16000 } } , "instructions" : { "response" : "Basic acceptance tests" } }`,
		}

		for _, spec := range equivalentAgentSpecifications {
			got, err := NormalizeCortexAgentSpecification(spec)
			require.NoError(t, err)
			require.Equal(t, want, got)
		}
	})

	t.Run("non-equivalent agent specifications", func(t *testing.T) {
		nonEquivalentAgentSpecifications := []struct {
			name string
			spec string
		}{
			{
				name: "different instructions.response string",
				spec: `{"instructions":{"response":"Different text"},"orchestration":{"budget":{"seconds":30,"tokens":16000}}}`,
			},
			{
				name: "different budget.seconds",
				spec: `{"instructions":{"response":"Basic acceptance tests"},"orchestration":{"budget":{"seconds":31,"tokens":16000}}}`,
			},
			{
				name: "different budget.tokens",
				spec: `{"instructions":{"response":"Basic acceptance tests"},"orchestration":{"budget":{"seconds":30,"tokens":8000}}}`,
			},
			{
				name: "extra orchestration field under instructions",
				spec: `{"instructions":{"response":"Basic acceptance tests","orchestration":"For any revenue question use Analyst"},"orchestration":{"budget":{"seconds":30,"tokens":16000}}}`,
			},
			{
				name: "missing budget.tokens",
				spec: `{"instructions":{"response":"Basic acceptance tests"},"orchestration":{"budget":{"seconds":30}}}`,
			},
			{
				name: "extra top-level key",
				spec: `{"instructions":{"response":"Basic acceptance tests"},"models":{"orchestration":"claude-4-sonnet"},"orchestration":{"budget":{"seconds":30,"tokens":16000}}}`,
			},
		}

		for _, tc := range nonEquivalentAgentSpecifications {
			t.Run(tc.name, func(t *testing.T) {
				got, err := NormalizeCortexAgentSpecification(tc.spec)
				require.NoError(t, err)
				require.NotEqual(t, want, got)
			})
		}
	})

	emptyLike := []string{
		"",
		"  \n\t ",
	}

	t.Run("empty and whitespace only", func(t *testing.T) {
		for _, spec := range emptyLike {
			got, err := NormalizeCortexAgentSpecification(spec)
			require.NoError(t, err)
			require.Equal(t, "{}", got)
		}
	})

	t.Run("invalid returns error", func(t *testing.T) {
		_, err := NormalizeCortexAgentSpecification("{broken")
		require.Error(t, err)
	})
}

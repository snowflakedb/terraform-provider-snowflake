package sdk

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	id := cortexAgentsTestIdSchemaObjectIdentifier
	schemaId := id.SchemaId()
	spec := `models:
	orchestration: claude-4-sonnet`
	profile := `{"display_name": "My Business Assistant", "avatar": "business-icon.png", "color": "blue"}`
	expectedProfile := strings.ReplaceAll(profile, `"`, `\"`)

	cortexAgentsTests.Create.
		withDefaultOpts(func() *CreateCortexAgentOptions {
			return &CreateCortexAgentOptions{
				name:              id,
				FromSpecification: spec,
			}
		}).
		withExpectedSqlf(
			case_CortexAgents_sql_Create_basic,
			"CREATE AGENT %s FROM SPECIFICATION $$%s$$", id.FullyQualifiedName(), spec,
		).
		withModifyAndExpectedSqlf(
			case_CortexAgents_sql_Create_all,
			func(opts *CreateCortexAgentOptions) {
				opts.IfNotExists = Bool(true)
				opts.Comment = String("some comment")
				opts.Profile = String(profile)
			},
			"CREATE AGENT IF NOT EXISTS %s COMMENT = 'some comment' PROFILE = '%s' FROM SPECIFICATION $$%s$$",
			id.FullyQualifiedName(), expectedProfile, spec,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateCortexAgentOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE AGENT %s FROM SPECIFICATION $$%s$$", id.FullyQualifiedName(), spec,
		)

	cortexAgentsTests.Alter.
		withModifyAndExpectedSqlf(
			case_CortexAgents_sql_Alter_Set,
			func(opts *AlterCortexAgentOptions) {
				opts.IfExists = Bool(true)
				opts.Set = &CortexAgentSet{
					Comment: &StringAllowEmpty{Value: "some comment"},
					Profile: String(profile),
				}
			},
			"ALTER AGENT IF EXISTS %s SET COMMENT = 'some comment', PROFILE = '%s'", id.FullyQualifiedName(), expectedProfile,
		).
		withModifyAndExpectedSqlf(
			case_CortexAgents_sql_Alter_ModifyLiveVersionSet,
			func(opts *AlterCortexAgentOptions) {
				opts.ModifyLiveVersionSet = &CortexAgentModifyLiveVersionSet{
					Specification: spec,
				}
			},
			"ALTER AGENT %s MODIFY LIVE VERSION SET SPECIFICATION = $$%s$$", id.FullyQualifiedName(), spec,
		)

	cortexAgentsTests.Drop.
		withExpectedSqlf(
			case_CortexAgents_sql_Drop_basic,
			"DROP AGENT %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexAgents_sql_Drop_all,
			func(opts *DropCortexAgentOptions) { opts.IfExists = Bool(true) },
			"DROP AGENT IF EXISTS %s", id.FullyQualifiedName(),
		)

	cortexAgentsTests.Show.
		withExpectedSql(case_CortexAgents_sql_Show_basic, "SHOW AGENTS").
		withModifyAndExpectedSqlf(
			case_CortexAgents_sql_Show_all,
			func(opts *ShowCortexAgentOptions) {
				opts.Like = &Like{Pattern: String("like-pattern")}
				opts.In = &ExtendedIn{In: In{Schema: schemaId}}
				opts.StartsWith = String("starts-with-pattern")
				opts.Limit = &LimitFrom{Rows: Int(10), From: String("limit-from")}
			},
			"SHOW AGENTS LIKE 'like-pattern' IN SCHEMA %s STARTS WITH 'starts-with-pattern' LIMIT 10 FROM 'limit-from'",
			schemaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexAgents_sql_Show_Like,
			func(opts *ShowCortexAgentOptions) { opts.Like = &Like{Pattern: String("like-pattern")} },
			"SHOW AGENTS LIKE 'like-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_CortexAgents_sql_Show_In,
			func(opts *ShowCortexAgentOptions) { opts.In = &ExtendedIn{In: In{Schema: schemaId}} },
			"SHOW AGENTS IN SCHEMA %s", schemaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexAgents_sql_Show_StartsWith,
			func(opts *ShowCortexAgentOptions) { opts.StartsWith = String("starts-with-pattern") },
			"SHOW AGENTS STARTS WITH 'starts-with-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_CortexAgents_sql_Show_Limit,
			func(opts *ShowCortexAgentOptions) { opts.Limit = &LimitFrom{Rows: Int(10), From: String("limit-from")} },
			"SHOW AGENTS LIMIT 10 FROM 'limit-from'",
		)

	cortexAgentsTests.Describe.
		withExpectedSqlf(
			case_CortexAgents_sql_Describe_basic,
			"DESCRIBE AGENT %s", id.FullyQualifiedName(),
		)
}

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
				DisplayName: new("My Assistant"),
				Avatar:      new("assistant.png"),
				Color:       new("blue"),
			},
		},
		{
			name: "single field present",
			json: `{"color":"green"}`,
			expected: CortexAgentProfile{
				Color: new("green"),
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
				DisplayName: new("My Assistant"),
				Avatar:      new("assistant.png"),
				Color:       new("blue"),
			},
			want: `{"display_name":"My Assistant","avatar":"assistant.png","color":"blue"}`,
		},
		{
			name: "partial profile",
			profile: CortexAgentProfile{
				Color: new("green"),
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

package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	openflowConnectorDefinitionsTests.Show.
		withExpectedSql(
			case_OpenflowConnectorDefinitions_sql_Show_basic,
			"SHOW OPENFLOW CONNECTOR DEFINITIONS",
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectorDefinitions_sql_Show_Like,
			func(opts *ShowOpenflowConnectorDefinitionOptions) {
				opts.Like = &Like{Pattern: new("%_NO_CORTEX")}
			},
			"SHOW OPENFLOW CONNECTOR DEFINITIONS LIKE '%%_NO_CORTEX'",
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectorDefinitions_sql_Show_StartsWith,
			func(opts *ShowOpenflowConnectorDefinitionOptions) {
				opts.StartsWith = new("OPENFLOW_")
			},
			"SHOW OPENFLOW CONNECTOR DEFINITIONS STARTS WITH 'OPENFLOW_'",
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectorDefinitions_sql_Show_Limit,
			func(opts *ShowOpenflowConnectorDefinitionOptions) {
				opts.Limit = &LimitFrom{Rows: new(2)}
			},
			"SHOW OPENFLOW CONNECTOR DEFINITIONS LIMIT 2",
		).
		// Definitions are account-wide, so there is no IN scope. Clause order is LIKE, STARTS WITH, LIMIT.
		withModifyAndExpectedSqlf(
			case_OpenflowConnectorDefinitions_sql_Show_all,
			func(opts *ShowOpenflowConnectorDefinitionOptions) {
				opts.Like = &Like{Pattern: new("%POSTGRES%")}
				opts.StartsWith = new("OPENFLOW_")
				opts.Limit = &LimitFrom{Rows: new(5), From: new("OPENFLOW_A")}
			},
			"SHOW OPENFLOW CONNECTOR DEFINITIONS LIKE '%%POSTGRES%%' STARTS WITH 'OPENFLOW_' LIMIT 5 FROM 'OPENFLOW_A'",
		)
}

func TestParseOpenflowConnectorDefinitionCategories(t *testing.T) {
	testCases := []struct {
		Name     string
		Value    string
		Expected []string
	}{
		{
			Name:     "empty string",
			Value:    "",
			Expected: []string{},
		},
		{
			Name:     "blank string",
			Value:    "   ",
			Expected: []string{},
		},
		{
			Name:     "empty JSON array",
			Value:    "[]",
			Expected: []string{},
		},
		{
			// Not theoretical: the runtime's external_access_integrations column comes back as the literal
			// null once its last value is removed, so a values column really can report emptiness this way.
			Name:     "JSON null",
			Value:    "null",
			Expected: []string{},
		},
		{
			Name:     "empty array with spaces",
			Value:    " [] ",
			Expected: []string{},
		},
		{
			Name:     "single category",
			Value:    `["Databases"]`,
			Expected: []string{"Databases"},
		},
		{
			Name:     "multiple categories",
			Value:    `["Databases","Analytics"]`,
			Expected: []string{"Databases", "Analytics"},
		},
		{
			Name:     "spaces after the separator",
			Value:    `["Databases", "Analytics"]`,
			Expected: []string{"Databases", "Analytics"},
		},
		{
			Name:     "category containing a comma is not split",
			Value:    `["Databases, relational","Analytics"]`,
			Expected: []string{"Databases, relational", "Analytics"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			categories, err := ParseOpenflowConnectorDefinitionCategories(testCase.Value)
			require.NoError(t, err)
			assert.Equal(t, testCase.Expected, categories)
			// Stripping the quotes and brackets is the whole point.
			for _, category := range categories {
				assert.NotContains(t, category, `"`)
				assert.NotContains(t, category, "[")
			}
		})
	}

	// Unquoted elements are accepted rather than rejected, since the shared parser splits on commas
	// rather than unmarshalling.
	t.Run("unquoted elements", func(t *testing.T) {
		categories, err := ParseOpenflowConnectorDefinitionCategories("[Databases, Analytics]")
		require.NoError(t, err)
		assert.Equal(t, []string{"Databases", "Analytics"}, categories)
	})
}

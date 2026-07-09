package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestV098TaskStateUpgrader_sessionParametersDoNotHijackIdentificationFields verifies the fix for SNOW-3649833:
// the v0.98 state upgrader must not let arbitrary keys from the legacy free-form session_parameters map
// overwrite the reserved task identification fields ("id", "schema", "name", "database").
//
// This behavior is only exercisable as a unit test, not an acceptance/migration test. The
// key/value pairs would have to live in the persisted state under session_parameters WITHOUT appearing
// in the configuration first. An acceptance test creates the resource through a real provider from config, so
// these keys would either be rejected by Snowflake (they are not valid session parameters) during the
// old-provider create step.
func TestV098TaskStateUpgrader_sessionParametersDoNotHijackIdentificationFields(t *testing.T) {
	rawState := map[string]any{
		"id":                          "MYDB|MYSCHEMA|MYTASK",
		"database":                    "MYDB",
		"schema":                      "MYSCHEMA",
		"name":                        "MYTASK",
		"when":                        "",
		"enabled":                     false,
		"allow_overlapping_execution": false,
		"session_parameters": map[string]any{
			// Malicious entries attempting to clobber identification fields.
			"id":       "A|B|C",
			"database": "D",
			"schema":   "E",
			"name":     "F",
			// A legitimate session parameter that must still be carried over to the top level.
			"LOG_LEVEL": "INFO",
		},
	}

	result, err := v098TaskStateUpgrader(context.Background(), rawState, nil)
	require.NoError(t, err)

	// Identification fields must keep the original task's values, not the hijacked ones.
	assert.Equal(t, `"MYDB"."MYSCHEMA"."MYTASK"`, result["id"])
	assert.Equal(t, "MYDB", result["database"])
	assert.Equal(t, "MYSCHEMA", result["schema"])
	assert.Equal(t, "MYTASK", result["name"])

	// Legitimate session parameters are still copied up, and the original map is removed.
	assert.Equal(t, "INFO", result["LOG_LEVEL"])
	assert.NotContains(t, result, "session_parameters")
}

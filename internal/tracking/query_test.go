package tracking

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/stretchr/testify/require"
)

func TestTrimMetadata(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected string
	}{
		{
			Input:    "select 1",
			Expected: "select 1",
		},
		{
			Input:    "select 1; --some comment",
			Expected: "select 1; --some comment",
		},
		{
			Input:    fmt.Sprintf("select 1; --%s", MetadataPrefix),
			Expected: "select 1;",
		},
		{
			Input:    fmt.Sprintf("select 1; --%s ", MetadataPrefix),
			Expected: "select 1;",
		},
		{
			Input:    fmt.Sprintf("select 1; --%s some text after", MetadataPrefix),
			Expected: "select 1;",
		},
	}

	for _, tc := range testCases {
		t.Run("TrimMetadata: "+tc.Input, func(t *testing.T) {
			trimmedInput := TrimMetadata(tc.Input)
			assert.Equal(t, tc.Expected, trimmedInput)
		})
	}
}

func TestAppendMetadata(t *testing.T) {
	metadata := newTestMetadata("123", resources.Account, CreateOperation)
	sql := "SELECT 1"

	bytes, err := json.Marshal(metadata)
	require.NoError(t, err)

	expectedSql := fmt.Sprintf("%s --%s %s", sql, MetadataPrefix, string(bytes))

	newSql, err := AppendMetadata(sql, metadata)
	require.NoError(t, err)
	require.Equal(t, expectedSql, newSql)
}

func TestParseMetadata(t *testing.T) {
	metadata := newTestMetadata("123", resources.Account, CreateOperation)
	bytes, err := json.Marshal(metadata)
	require.NoError(t, err)
	sql := fmt.Sprintf("SELECT 1 --%s %s", MetadataPrefix, string(bytes))

	parsedMetadata, err := ParseMetadata(sql)
	require.NoError(t, err)
	require.Equal(t, metadata, parsedMetadata)
}

func TestParseInvalidMetadataKeys(t *testing.T) {
	sql := fmt.Sprintf(`SELECT 1 --%s {"key": "value"}`, MetadataPrefix)

	parsedMetadata, err := ParseMetadata(sql)
	require.ErrorContains(t, err, "schema version for metadata should not be empty")
	require.ErrorContains(t, err, "provider version for metadata should not be empty")
	require.ErrorContains(t, err, "either resource or data source name for metadata should be specified")
	require.ErrorContains(t, err, "operation for metadata should not be empty")
	require.Equal(t, Metadata{}, parsedMetadata)
}

func TestParseInvalidMetadataJson(t *testing.T) {
	sql := fmt.Sprintf(`SELECT 1 --%s "key": "value"`, MetadataPrefix)

	parsedMetadata, err := ParseMetadata(sql)
	require.ErrorContains(t, err, "failed to unmarshal metadata from sql")
	require.Equal(t, Metadata{}, parsedMetadata)
}

func TestParseMetadataFromInvalidSqlCommentPrefix(t *testing.T) {
	metadata := newTestMetadata("123", resources.Account, CreateOperation)
	sql := "SELECT 1"

	bytes, err := json.Marshal(metadata)
	require.NoError(t, err)

	parsedMetadata, err := ParseMetadata(fmt.Sprintf("%s --invalid_prefix %s", sql, string(bytes)))
	require.ErrorContains(t, err, "failed to parse metadata from sql")
	require.Equal(t, Metadata{}, parsedMetadata)
}

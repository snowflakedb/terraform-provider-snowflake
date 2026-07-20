package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseNullIfProperty(t *testing.T) {
	require.Equal(t, []NullString{{"a"}, {"b"}}, parseNullIfProperty("[a, b]"))
	require.Equal(t, []NullString{}, parseNullIfProperty("[]"))
}

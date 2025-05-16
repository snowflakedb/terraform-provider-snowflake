package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestEscapeString(t *testing.T) {
	r := require.New(t)

	r.Equal(`\'`, snowflake.EscapeString(`'`))
	r.Equal(`\\\'`, snowflake.EscapeString(`\'`))
}

func TestEscapeSnowflakeString(t *testing.T) {
	r := require.New(t)
	r.Equal(`'table''s quoted'`, snowflake.EscapeSnowflakeString(`table's quoted`))
}

func TestUnescapeSnowflakeString(t *testing.T) {
	r := require.New(t)
	r.Equal(`table's quoted`, snowflake.UnescapeSnowflakeString(`'table''s quoted'`))
}

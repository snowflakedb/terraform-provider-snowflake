package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestEscapeString(t *testing.T) {
	a := require.New(t)

	a.Equal(`\'`, snowflake.EscapeString(`'`))
	a.Equal(`\\\'`, snowflake.EscapeString(`\'`))
}

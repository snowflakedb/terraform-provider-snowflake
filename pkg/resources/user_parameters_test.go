package resources

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_UserParametersSchema(t *testing.T) {
	t.Run("description references parameter docs correctly", func(t *testing.T) {
		require.True(t, strings.HasSuffix(userParametersSchema["abort_detached_query"].Description, "For more information, check [ABORT_DETACHED_QUERY docs](https://docs.snowflake.com/en/sql-reference/parameters#abort-detached-query)."))
	})

	t.Run("preserves pre-catalog integer validation behavior", func(t *testing.T) {
		for name, parameterSchema := range userParametersSchema {
			if parameterSchema.Type == schema.TypeInt {
				require.Nil(t, parameterSchema.ValidateDiagFunc, name)
			}
		}
	})
}

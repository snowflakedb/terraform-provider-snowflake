package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_ResourceModel_idEncoding(t *testing.T) {
	t.Run("current encoding", func(t *testing.T) {
		model := lifecycleModel()
		require.Equal(t, "helpers.EncodeResourceIdentifier[sdk.SchemaObjectIdentifier]", model.EncodeIdFn())
	})

	t.Run("legacy encoding", func(t *testing.T) {
		model := lifecycleModel(func(obj *ssot.Object) { obj.LegacyIdEncoding = true })
		require.Equal(t, "func(renamed ...sdk.SchemaObjectIdentifier) string { return helpers.EncodeSnowflakeID(renamed[0]) }", model.EncodeIdFn())
	})
}

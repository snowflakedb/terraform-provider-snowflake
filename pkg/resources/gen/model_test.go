package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_ResourceModel_idEncoding(t *testing.T) {
	t.Run("current encoding", func(t *testing.T) {
		model := lifecycleModel()
		require.Equal(t, "helpers.EncodeResourceIdentifier(newId)", model.EncodeId("newId"))
		require.Equal(t, "helpers.EncodeResourceIdentifier[sdk.SchemaObjectIdentifier]", model.EncodeIdFn())
		require.Equal(t, "sdk.ParseSchemaObjectIdentifier", model.DecodeId())
	})

	t.Run("legacy encoding", func(t *testing.T) {
		model := lifecycleModel(func(obj *ssot.Object) { obj.LegacyIdEncoding = true })
		require.Equal(t, "helpers.EncodeSnowflakeID(newId)", model.EncodeId("newId"))
		require.Equal(t, "func(renamed ...sdk.SchemaObjectIdentifier) string { return helpers.EncodeSnowflakeID(renamed[0]) }", model.EncodeIdFn())
		require.Equal(t, "helpers.DecodeSnowflakeIDErrLegacy[sdk.SchemaObjectIdentifier]", model.DecodeId())
	})
}

func Test_ResourceModel_ReadCollectsDiags(t *testing.T) {
	require.False(t, lifecycleModel().ReadCollectsDiags())
	require.True(t, lifecycleModel(withHooks(ssot.WithAfterShowById(ssot.ExtCall()))).ReadCollectsDiags())
}

package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_ResourceModel_ReadCollectsDiags(t *testing.T) {
	require.False(t, lifecycleModel().ReadCollectsDiags())
	require.True(t, lifecycleModel(withHooks(ssot.WithAfterShowById(ssot.ExtCall()))).ReadCollectsDiags())
}

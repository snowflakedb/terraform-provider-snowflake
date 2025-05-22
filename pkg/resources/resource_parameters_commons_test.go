package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func Test_enrichWithReferenceToParameterDocs(t *testing.T) {
	t.Run("formats the docs reference correctly", func(t *testing.T) {
		description := random.Comment()

		enrichedDescription := enrichWithReferenceToParameterDocs(sdk.UserParameterAbortDetachedQuery, description)

		require.Equal(t, description+" "+"For more information, check [ABORT_DETACHED_QUERY docs](https://docs.snowflake.com/en/sql-reference/parameters#abort-detached-query).", enrichedDescription)
	})
}

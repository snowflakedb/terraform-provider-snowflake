package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_projectionPolicyApplySetUnset_sendsTheBodyInItsOwnRequest(t *testing.T) {
	d := schema.TestResourceDataRaw(t, projectionPolicySchema, map[string]any{
		"body":    "projection_constraint(allow => true)",
		"comment": "c",
	})
	requests, err := projectionPolicyApplySetUnset(d)
	require.NoError(t, err)

	// SET BODY cannot be combined with SET COMMENT.
	require.Equal(t, sdk.NewProjectionPolicySetRequest().WithSetBody("projection_constraint(allow => true)"), requests.setBody)
	require.Equal(t, sdk.NewProjectionPolicySetRequest().WithComment("c"), requests.set)
	require.Equal(t, sdk.NewProjectionPolicyUnsetRequest(), requests.unset)
}

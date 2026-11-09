package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_externalVolumeApplySetUnset_hasOnlyASetRequest(t *testing.T) {
	d := schema.TestResourceDataRaw(t, externalVolumeSchema, map[string]any{
		"comment":      "c",
		"allow_writes": "false",
	})
	requests, err := externalVolumeApplySetUnset(d)
	require.NoError(t, err)

	// The object has no UNSET, so a removed allow_writes is set to its default instead.
	require.Equal(t, sdk.NewExternalVolumeSetRequest().WithComment("c").WithAllowWrites(false), requests.set)
}

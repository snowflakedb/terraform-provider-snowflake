//go:build !account_level_tests

package testint

import (
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ToTimestampLTZ(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
		Set: &sdk.AccountSet{
			LegacyParameters: &sdk.AccountLevelParameters{
				SessionParameters: &sdk.SessionParameters{
					TimestampTypeMapping: sdk.Pointer(sdk.TimestampTypeMappingLtz),
					Timezone:             sdk.String("UTC"),
				},
			},
		},
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				LegacyParameters: &sdk.AccountLevelParametersUnset{
					SessionParameters: &sdk.SessionParametersUnset{
						TimestampTypeMapping: sdk.Bool(true),
						Timezone:             sdk.Bool(true),
					},
				},
			},
		})
		require.NoError(t, err)
	})
	// new warehouse created on purpose
	warehouseTest, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	err = client.Sessions.UseWarehouse(ctx, warehouseTest.ID())
	require.NoError(t, err)
	now := time.Now()
	actual, err := client.ConversionFunctions.ToTimestampLTZ(ctx, now)
	require.NoError(t, err)
	expected := now.UTC()
	assert.Equal(t, expected, actual)
}

func TestInt_ToTimestampNTZ(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
		Set: &sdk.AccountSet{
			LegacyParameters: &sdk.AccountLevelParameters{
				SessionParameters: &sdk.SessionParameters{
					TimestampTypeMapping: sdk.Pointer(sdk.TimestampTypeMappingLtz),
					Timezone:             sdk.String("UTC"),
				},
			},
		},
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				LegacyParameters: &sdk.AccountLevelParametersUnset{
					SessionParameters: &sdk.SessionParametersUnset{
						TimestampTypeMapping: sdk.Bool(true),
						Timezone:             sdk.Bool(true),
					},
				},
			},
		})
		require.NoError(t, err)
	})
	// new warehouse created on purpose
	warehouseTest, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	err = client.Sessions.UseWarehouse(ctx, warehouseTest.ID())
	require.NoError(t, err)
	now := time.Now()
	actual, err := client.ConversionFunctions.ToTimestampNTZ(ctx, now)
	require.NoError(t, err)
	assert.Equal(t, time.UTC, actual.Location())
}

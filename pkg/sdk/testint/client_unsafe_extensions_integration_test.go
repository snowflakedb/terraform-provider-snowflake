//go:build !account_level_tests

package testint

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Client_UnsafeQuery(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("test show databases", func(t *testing.T) {
		sql := fmt.Sprintf("SHOW DATABASES LIKE '%%%s%%'", testClientHelper().Ids.DatabaseId().Name())
		results, err := client.QueryUnsafe(ctx, sql)
		require.NoError(t, err)

		assert.Len(t, results, 1)
		row := results[0]

		require.NotNil(t, row["name"])
		require.NotNil(t, row["created_on"])
		require.NotNil(t, row["owner"])
		require.NotNil(t, row["options"])
		require.NotNil(t, row["comment"])
		require.NotNil(t, row["is_default"])

		assert.Equal(t, testClientHelper().Ids.DatabaseId().Name(), *row["name"])
		assert.NotEmpty(t, *row["created_on"])
		assert.Equal(t, "STANDARD", *row["kind"])
		assert.Equal(t, "ACCOUNTADMIN", *row["owner"])
		assert.Equal(t, "", *row["options"])
		assert.Equal(t, "", *row["comment"])
		assert.Equal(t, "N", *row["is_default"])
	})

	t.Run("test more results", func(t *testing.T) {
		db1, db1Cleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(db1Cleanup)
		db2, db2Cleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(db2Cleanup)
		db3, db3Cleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(db3Cleanup)

		sql := "SHOW DATABASES"
		results, err := client.QueryUnsafe(ctx, sql)
		require.NoError(t, err)

		require.GreaterOrEqual(t, len(results), 4)
		names := make([]any, len(results))
		for i, r := range results {
			names[i] = *r["name"]
		}
		assert.Contains(t, names, testClientHelper().Ids.DatabaseId().Name())
		assert.Contains(t, names, db1.Name)
		assert.Contains(t, names, db2.Name)
		assert.Contains(t, names, db3.Name)
	})
}

//go:build !account_level_tests

package testint

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
)

// Research for basic object tracking done as part of SNOW-1737787

// https://docs.snowflake.com/en/sql-reference/parameters#query-tag
func TestInt_ContextQueryTags(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	// set query_tag on user level
	userQueryTag := "user query tag"
	testClientHelper().User.AlterCurrentUser(t, &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			SessionParameters: &sdk.SessionParameters{
				QueryTag: sdk.String(userQueryTag),
			},
		},
	})
	t.Cleanup(func() {
		testClientHelper().User.AlterCurrentUser(t, &sdk.AlterUserOptions{
			Unset: &sdk.UserUnset{
				SessionParameters: &sdk.SessionParametersUnset{
					QueryTag: sdk.Bool(true),
				},
			},
		})
	})
	queryId := executeQueryAndReturnQueryId(t, context.Background(), client)
	queryTagResult := testClientHelper().InformationSchema.GetQueryHistoryByQueryId(t, 20, queryId)
	require.Equal(t, userQueryTag, queryTagResult.QueryTag)

	// set query_tag on session level
	sessionQueryTag := "session query tag"
	require.NoError(t, client.Sessions.AlterSession(ctx, &sdk.AlterSessionOptions{
		Set: &sdk.SessionSet{
			SessionParameters: &sdk.SessionParameters{
				QueryTag: sdk.String(sessionQueryTag),
			},
		},
	}))
	t.Cleanup(func() {
		require.NoError(t, client.Sessions.AlterSession(ctx, &sdk.AlterSessionOptions{
			Unset: &sdk.SessionUnset{
				SessionParametersUnset: &sdk.SessionParametersUnset{
					QueryTag: sdk.Bool(true),
				},
			},
		}))
	})
	queryId = executeQueryAndReturnQueryId(t, context.Background(), client)
	queryTagResult = testClientHelper().InformationSchema.GetQueryHistoryByQueryId(t, 20, queryId)
	require.Equal(t, sessionQueryTag, queryTagResult.QueryTag)

	// set query_tag on query level
	perQueryQueryTag := "per-query query tag"
	ctxWithQueryTag := gosnowflake.WithQueryTag(context.Background(), perQueryQueryTag)
	queryId = executeQueryAndReturnQueryId(t, ctxWithQueryTag, client)
	queryTagResult = testClientHelper().InformationSchema.GetQueryHistoryByQueryId(t, 20, queryId)
	require.Equal(t, perQueryQueryTag, queryTagResult.QueryTag)
}

func executeQueryAndReturnQueryId(t *testing.T, ctx context.Context, client *sdk.Client) string {
	t.Helper()
	queryIdChan := make(chan string, 1)
	ctx = gosnowflake.WithQueryIDChan(ctx, queryIdChan)

	_, err := client.QueryUnsafe(ctx, "SELECT 1")
	require.NoError(t, err)

	return <-queryIdChan
}

// https://select.dev/posts/snowflake-query-tags#using-query-comments-instead-of-query-tags
func TestInt_QueryComment(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	queryIdChan := make(chan string, 1)
	metadata := `{"comment": "some comment"}`
	_, err := client.QueryUnsafe(gosnowflake.WithQueryIDChan(ctx, queryIdChan), fmt.Sprintf(`SELECT 1; --%s`, metadata))
	require.NoError(t, err)
	queryId := <-queryIdChan

	queryText := testClientHelper().InformationSchema.GetQueryHistoryByQueryId(t, 20, queryId).QueryText
	require.Equal(t, metadata, strings.Split(queryText, "--")[1])
}

// TODO [SNOW-2054366]: Use dedicated users for these tests.
func TestInt_AppName(t *testing.T) {
	// https://community.snowflake.com/s/article/How-to-see-application-name-added-in-the-connection-string-in-Snowsight
	t.Skip("there is no way to check client application name by querying Snowflake's")

	version := "v0.99.0"
	config := sdk.DefaultConfig()
	config.Application = fmt.Sprintf("terraform-provider-snowflake:%s", version)
	client, err := sdk.NewClient(config)
	require.NoError(t, err)

	_, err = client.QueryUnsafe(context.Background(), "SELECT 1")
	require.NoError(t, err)
}

// TODO(SNOW-1805150): Document potential usage of connection string

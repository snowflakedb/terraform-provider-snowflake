//go:build non_account_level_tests

package testint

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_CortexAgents(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	expectedSpecAsMap := func(response string) map[string]any {
		return map[string]any{
			"orchestration": map[string]any{
				"budget": map[string]any{
					"seconds": float64(30),
					"tokens":  float64(16000),
				},
			},
			"instructions": map[string]any{
				"response": response,
			},
		}
	}

	createCortexAgent := func(t *testing.T) sdk.SchemaObjectIdentifier {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		cleanup := testClientHelper().CortexAgent.CreateWithId(t, id)
		t.Cleanup(cleanup)

		return id
	}

	t.Run("create cortex agent: complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		profile := `{"display_name":"My Business Assistant","avatar":"business-icon.png","color":"blue"}`
		response := "Complete integration test"
		spec := testClientHelper().CortexAgent.SampleSpecWithResponse(t, response)

		cleanup := testClientHelper().CortexAgent.CreateWithRequest(t, sdk.NewCreateCortexAgentRequest(id, spec).
			WithIfNotExists(true).
			WithComment(comment).
			WithProfile(profile))
		t.Cleanup(cleanup)

		expectedProfile := sdk.CortexAgentProfile{
			DisplayName: sdk.String("My Business Assistant"),
			Avatar:      sdk.String("business-icon.png"),
			Color:       sdk.String("blue"),
		}
		assertThatObject(t, objectassert.CortexAgent(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasCortexAgentProfile(expectedProfile),
		)
		assertThatObject(t, objectassert.CortexAgentDetails(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasCortexAgentProfile(expectedProfile).
			HasCortexAgentSpec(expectedSpecAsMap(response)).
			HasCreatedOnNotEmpty().
			HasDefaultVersionName("LAST").
			HasVersions(`["VERSION$1"]`).
			HasAliases(`{"DEFAULT":"VERSION$1","FIRST":"VERSION$1","LAST":"VERSION$1"}`),
		)
	})

	t.Run("create cortex agent: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		response := "Without optionals"
		spec := testClientHelper().CortexAgent.SampleSpecWithResponse(t, response)

		cleanup := testClientHelper().CortexAgent.CreateWithRequest(t, sdk.NewCreateCortexAgentRequest(id, spec))
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.CortexAgent(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasNoComment().
			HasNoProfile(),
		)
		assertThatObject(t, objectassert.CortexAgentDetails(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasNoComment().
			HasNoProfile().
			HasCortexAgentSpec(expectedSpecAsMap(response)).
			HasCreatedOnNotEmpty().
			HasDefaultVersionName("LAST").
			HasVersions(`["VERSION$1"]`).
			HasAliases(`{"DEFAULT":"VERSION$1","FIRST":"VERSION$1","LAST":"VERSION$1"}`),
		)
	})

	t.Run("alter cortex agent: set comment and profile", func(t *testing.T) {
		id := createCortexAgent(t)
		comment := random.Comment()
		profile := `{"display_name":"Renamed Assistant"}`

		err := client.CortexAgents.Alter(ctx, sdk.NewAlterCortexAgentRequest(id).
			WithSet(*sdk.NewCortexAgentSetRequest().
				WithComment(sdk.StringAllowEmpty{Value: comment}).
				WithProfile(profile)))
		require.NoError(t, err)

		expectedProfile := sdk.CortexAgentProfile{
			DisplayName: sdk.String("Renamed Assistant"),
		}
		assertThatObject(t, objectassert.CortexAgent(t, id).
			HasComment(comment).
			HasCortexAgentProfile(expectedProfile),
		)
		assertThatObject(t, objectassert.CortexAgentDetails(t, id).
			HasComment(comment).
			HasCortexAgentProfile(expectedProfile))

		err = client.CortexAgents.Alter(ctx, sdk.NewAlterCortexAgentRequest(id).
			WithSet(*sdk.NewCortexAgentSetRequest().
				WithComment(sdk.StringAllowEmpty{}).
				WithProfile("{}")))
		require.NoError(t, err)

		expectedEmptyProfile := sdk.CortexAgentProfile{}
		assertThatObject(t, objectassert.CortexAgent(t, id).
			HasComment("").
			HasCortexAgentProfile(expectedEmptyProfile),
		)
		assertThatObject(t, objectassert.CortexAgentDetails(t, id).
			HasComment("").
			HasCortexAgentProfile(expectedEmptyProfile))
	})

	t.Run("alter cortex agent: modify live version set", func(t *testing.T) {
		id := createCortexAgent(t)
		newResponse := "Updated live version"
		newSpec := testClientHelper().CortexAgent.SampleSpecWithResponse(t, newResponse)

		err := client.CortexAgents.Alter(ctx, sdk.NewAlterCortexAgentRequest(id).
			WithModifyLiveVersionSet(*sdk.NewCortexAgentModifyLiveVersionSetRequest(newSpec)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.CortexAgentDetails(t, id).
			HasCortexAgentSpec(expectedSpecAsMap(newResponse)))
	})

	t.Run("drop cortex agent: existing", func(t *testing.T) {
		id := createCortexAgent(t)

		err := client.CortexAgents.Drop(ctx, sdk.NewDropCortexAgentRequest(id))
		require.NoError(t, err)

		_, err = client.CortexAgents.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop cortex agent: non-existing", func(t *testing.T) {
		err := client.CortexAgents.Drop(ctx, sdk.NewDropCortexAgentRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("show cortex agents", func(t *testing.T) {
		db, dbCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(dbCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_cortex_agentzzz")
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_cortex_agent_2_")
		id3 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_cortex_agent_3_")
		id4 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(sdk.NewDatabaseObjectIdentifier(db.Name, "PUBLIC"))
		ids := []sdk.SchemaObjectIdentifier{id1, id2, id3, id4}
		for _, id := range ids {
			cleanup := testClientHelper().CortexAgent.CreateWithId(t, id)
			t.Cleanup(cleanup)
		}

		t.Run("like", func(t *testing.T) {
			cortexAgents, err := client.CortexAgents.Show(ctx, sdk.NewShowCortexAgentRequest().
				WithLike(sdk.Like{Pattern: sdk.String("test_cortex_agent_2_%")}).
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id1.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, cortexAgents, 1)
		})

		t.Run("starts_with", func(t *testing.T) {
			cortexAgents, err := client.CortexAgents.Show(ctx, sdk.NewShowCortexAgentRequest().
				WithStartsWith("test_cortex_agent_").
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id1.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, cortexAgents, 2)
		})

		t.Run("in_account", func(t *testing.T) {
			cortexAgents, err := client.CortexAgents.Show(ctx, sdk.NewShowCortexAgentRequest().
				WithIn(sdk.ExtendedIn{In: sdk.In{Account: sdk.Bool(true)}}))
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(cortexAgents), 4)
		})

		t.Run("in_database", func(t *testing.T) {
			cortexAgents, err := client.CortexAgents.Show(ctx, sdk.NewShowCortexAgentRequest().
				WithIn(sdk.ExtendedIn{In: sdk.In{Database: id1.DatabaseId()}}))
			require.NoError(t, err)
			assert.Len(t, cortexAgents, 3)
		})

		t.Run("in_schema", func(t *testing.T) {
			cortexAgents, err := client.CortexAgents.Show(ctx, sdk.NewShowCortexAgentRequest().
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id1.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, cortexAgents, 3)
		})

		t.Run("limit", func(t *testing.T) {
			cortexAgents, err := client.CortexAgents.Show(ctx, sdk.NewShowCortexAgentRequest().
				WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}).
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id1.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, cortexAgents, 1)
		})

		t.Run("limit from", func(t *testing.T) {
			cortexAgents, err := client.CortexAgents.Show(ctx, sdk.NewShowCortexAgentRequest().
				WithLimit(sdk.LimitFrom{Rows: sdk.Int(1), From: sdk.String("test_cortex_agent_")}).
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id1.SchemaId()}}))
			require.NoError(t, err)
			require.Len(t, cortexAgents, 1)
			require.True(t, strings.HasPrefix(cortexAgents[0].Name, "test_cortex_agent_2"))
		})
	})

	t.Run("describe cortex agent: non-existing", func(t *testing.T) {
		_, err := client.CortexAgents.Describe(ctx, NonExistingSchemaObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

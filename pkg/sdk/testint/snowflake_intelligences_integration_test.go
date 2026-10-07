//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SnowflakeIntelligences(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertSnowflakeIntelligence := func(t *testing.T, intelligence *sdk.SnowflakeIntelligence, id sdk.AccountObjectIdentifier) {
		t.Helper()
		assertThatObject(t, objectassert.SnowflakeIntelligenceFromObject(t, intelligence).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(""))
	}

	assertSnowflakeIntelligenceDetails := func(t *testing.T, id sdk.AccountObjectIdentifier) {
		t.Helper()
		assertThatObject(t, objectassert.SnowflakeIntelligenceDetails(t, id).
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasCreatedOnNotEmpty().
			HasNoBrandName().
			HasNoVanityUrl().
			HasNoWelcomeMessage().
			HasNoIconLightPath().
			HasNoIconDarkPath().
			HasNoLogoLightPath().
			HasNoLogoDarkPath().
			HasNoFaviconPath().
			HasNoAccentColorLight().
			HasNoAccentColorDark().
			HasNoComment())
	}

	createSnowflakeIntelligenceWithRequest := func(t *testing.T, request *sdk.CreateSnowflakeIntelligenceRequest) *sdk.SnowflakeIntelligence {
		t.Helper()
		id := request.GetName()

		err := client.SnowflakeIntelligences.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SnowflakeIntelligence.DropFunc(t, id))

		intelligence, err := client.SnowflakeIntelligences.ShowByID(ctx, id)
		require.NoError(t, err)

		return intelligence
	}

	createSnowflakeIntelligenceRequest := func(t *testing.T) *sdk.CreateSnowflakeIntelligenceRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateSnowflakeIntelligenceRequest(id)
	}

	createSnowflakeIntelligence := func(t *testing.T) *sdk.SnowflakeIntelligence {
		t.Helper()
		return createSnowflakeIntelligenceWithRequest(t, createSnowflakeIntelligenceRequest(t))
	}

	createAgent := func(t *testing.T) sdk.SchemaObjectIdentifier {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		t.Cleanup(testClientHelper().CortexAgent.CreateWithId(t, id))
		return id
	}

	showAgents := func(t *testing.T, id sdk.AccountObjectIdentifier) []sdk.SnowflakeIntelligenceAgent {
		t.Helper()
		agents, err := client.SnowflakeIntelligences.ShowAgents(ctx, sdk.NewShowAgentsSnowflakeIntelligenceRequest(id))
		require.NoError(t, err)
		return agents
	}

	detachAgentOnCleanup := func(t *testing.T, id sdk.AccountObjectIdentifier, agentId sdk.SchemaObjectIdentifier) {
		t.Helper()
		t.Cleanup(func() {
			agents, err := client.SnowflakeIntelligences.ShowAgents(ctx, sdk.NewShowAgentsSnowflakeIntelligenceRequest(id))
			if err != nil {
				return
			}
			_, err = collections.FindFirst(agents, func(agent sdk.SnowflakeIntelligenceAgent) bool {
				return agent.ID().FullyQualifiedName() == agentId.FullyQualifiedName()
			})
			if err != nil {
				return
			}
			_ = client.SnowflakeIntelligences.Alter(ctx, sdk.NewAlterSnowflakeIntelligenceRequest(id).WithDropAgent(agentId))
		})
	}

	t.Run("create: basic", func(t *testing.T) {
		request := createSnowflakeIntelligenceRequest(t)

		intelligence := createSnowflakeIntelligenceWithRequest(t, request)

		assertSnowflakeIntelligence(t, intelligence, request.GetName())
		assertSnowflakeIntelligenceDetails(t, intelligence.ID())
		require.Empty(t, showAgents(t, intelligence.ID()))
	})

	t.Run("create: if not exists", func(t *testing.T) {
		request := createSnowflakeIntelligenceRequest(t).WithIfNotExists(true)

		intelligence := createSnowflakeIntelligenceWithRequest(t, request)

		assertSnowflakeIntelligence(t, intelligence, request.GetName())
		assertSnowflakeIntelligenceDetails(t, intelligence.ID())

		err := client.SnowflakeIntelligences.Create(ctx, sdk.NewCreateSnowflakeIntelligenceRequest(intelligence.ID()).WithIfNotExists(true))
		require.NoError(t, err)
	})

	t.Run("create: or replace", func(t *testing.T) {
		id := createSnowflakeIntelligence(t).ID()
		agentId := createAgent(t)
		detachAgentOnCleanup(t, id, agentId)

		err := client.SnowflakeIntelligences.Alter(ctx, sdk.NewAlterSnowflakeIntelligenceRequest(id).WithAddAgent(agentId))
		require.NoError(t, err)
		require.Len(t, showAgents(t, id), 1)

		err = client.SnowflakeIntelligences.Create(ctx, sdk.NewCreateSnowflakeIntelligenceRequest(id).WithOrReplace(true))
		require.NoError(t, err)

		intelligence, err := client.SnowflakeIntelligences.ShowByID(ctx, id)
		require.NoError(t, err)
		assertSnowflakeIntelligence(t, intelligence, id)
		assertSnowflakeIntelligenceDetails(t, id)
		require.Empty(t, showAgents(t, id))
	})

	t.Run("create: second object", func(t *testing.T) {
		_ = createSnowflakeIntelligence(t)

		err := client.SnowflakeIntelligences.Create(ctx, createSnowflakeIntelligenceRequest(t))
		require.ErrorContains(t, err, "Multiple SNOWFLAKE INTELLIGENCE objects are not allowed in the same account.")
	})

	t.Run("alter: add and drop agent", func(t *testing.T) {
		id := createSnowflakeIntelligence(t).ID()
		agentId := createAgent(t)
		detachAgentOnCleanup(t, id, agentId)

		agent, err := client.CortexAgents.ShowByID(ctx, agentId)
		require.NoError(t, err)
		require.Empty(t, showAgents(t, id))

		err = client.SnowflakeIntelligences.Alter(ctx, sdk.NewAlterSnowflakeIntelligenceRequest(id).WithAddAgent(agentId))
		require.NoError(t, err)

		agents := showAgents(t, id)
		require.Len(t, agents, 1)
		assertThatObject(t, objectassert.SnowflakeIntelligenceAgentFromObject(t, &agents[0]).
			HasCreatedOnNotEmpty().
			HasName(agentId.Name()).
			HasDatabaseName(agentId.DatabaseName()).
			HasSchemaName(agentId.SchemaName()).
			HasOwner(agent.Owner).
			HasComment(agent.Comment).
			HasProfile(agent.Profile).
			HasIsSecure(false))

		err = client.SnowflakeIntelligences.Alter(ctx, sdk.NewAlterSnowflakeIntelligenceRequest(id).WithIfExists(true).WithDropAgent(agentId))
		require.NoError(t, err)
		require.Empty(t, showAgents(t, id))
	})

	t.Run("drop safely: existing", func(t *testing.T) {
		id := createSnowflakeIntelligence(t).ID()

		err := client.SnowflakeIntelligences.DropSafely(ctx, id)
		require.NoError(t, err)

		_, err = client.SnowflakeIntelligences.ShowByIDSafely(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop safely: non-existing", func(t *testing.T) {
		err := client.SnowflakeIntelligences.DropSafely(ctx, NonExistingAccountObjectIdentifier)
		require.NoError(t, err)
	})

	t.Run("show: default", func(t *testing.T) {
		intelligence := createSnowflakeIntelligence(t)

		returnedIntelligences, err := client.SnowflakeIntelligences.Show(ctx, sdk.NewShowSnowflakeIntelligenceRequest())
		require.NoError(t, err)

		assert.Contains(t, returnedIntelligences, *intelligence)
	})

	t.Run("describe: basic", func(t *testing.T) {
		id := createSnowflakeIntelligence(t).ID()

		assertSnowflakeIntelligenceDetails(t, id)
	})
}

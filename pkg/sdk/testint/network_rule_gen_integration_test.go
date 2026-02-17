//go:build non_account_level_tests

package testint

import (
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_NetworkRules(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("Create", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.NetworkRules.Create(ctx, sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{}, sdk.NetworkRuleModeIngress))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().NetworkRule.DropFunc(t, id))

		assertThatObject(t, objectassert.NetworkRule(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment("").
			HasType(sdk.NetworkRuleTypeIpv4).
			HasMode(sdk.NetworkRuleModeIngress).
			HasEntriesInValueList(0).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))

		assertThatObject(t, objectassert.NetworkRuleDetails(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasType(sdk.NetworkRuleTypeIpv4).
			HasMode(sdk.NetworkRuleModeIngress).
			HasValueList([]string{}))
	})

	t.Run("Create with all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "test comment"
		values := []sdk.NetworkRuleValue{
			{Value: "0.0.0.0"},
			{Value: "1.1.1.1"},
		}

		err := client.NetworkRules.Create(ctx, sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, values, sdk.NetworkRuleModeIngress).
			WithOrReplace(true).
			WithComment(comment))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().NetworkRule.DropFunc(t, id))

		assertThatObject(t, objectassert.NetworkRule(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment(comment).
			HasType(sdk.NetworkRuleTypeIpv4).
			HasMode(sdk.NetworkRuleModeIngress).
			HasEntriesInValueList(2).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))

		assertThatObject(t, objectassert.NetworkRuleDetails(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasType(sdk.NetworkRuleTypeIpv4).
			HasMode(sdk.NetworkRuleModeIngress).
			HasValueList([]string{"0.0.0.0", "1.1.1.1"}))
	})

	t.Run("Alter: set and unset", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.NetworkRules.Create(ctx, sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{}, sdk.NetworkRuleModeIngress))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().NetworkRule.DropFunc(t, id))

		setReq := sdk.NewNetworkRuleSetRequest([]sdk.NetworkRuleValue{
			{Value: "0.0.0.0"},
			{Value: "1.1.1.1"},
		}).WithComment("some comment")
		err = client.NetworkRules.Alter(ctx, sdk.NewAlterNetworkRuleRequest(id).WithSet(*setReq))
		require.NoError(t, err)

		assertThatObject(t, objectassert.NetworkRule(t, id).
			HasEntriesInValueList(2).
			HasComment("some comment"))
		assertThatObject(t, objectassert.NetworkRuleDetails(t, id).
			HasValueList([]string{"0.0.0.0", "1.1.1.1"}).
			HasComment("some comment"))

		unsetReq := sdk.NewNetworkRuleUnsetRequest().
			WithValueList(true).
			WithComment(true)
		err = client.NetworkRules.Alter(ctx, sdk.NewAlterNetworkRuleRequest(id).WithUnset(*unsetReq))
		require.NoError(t, err)

		assertThatObject(t, objectassert.NetworkRule(t, id).
			HasEntriesInValueList(0).
			HasComment(""))
		assertThatObject(t, objectassert.NetworkRuleDetails(t, id).
			HasValueList([]string{}).
			HasComment(""))
	})

	t.Run("Drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.NetworkRules.Create(ctx, sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{}, sdk.NetworkRuleModeIngress))
		require.NoError(t, err)

		_, err = client.NetworkRules.ShowByID(ctx, id)
		require.NoError(t, err)

		err = client.NetworkRules.Drop(ctx, sdk.NewDropNetworkRuleRequest(id))
		require.NoError(t, err)

		_, err = client.NetworkRules.ShowByID(ctx, id)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("Show", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.NetworkRules.Create(ctx, sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{}, sdk.NetworkRuleModeIngress).WithComment("some comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().NetworkRule.DropFunc(t, id))

		networkRules, err := client.NetworkRules.Show(ctx, sdk.NewShowNetworkRuleRequest().WithIn(sdk.In{
			Schema: id.SchemaId(),
		}).WithLike(sdk.Like{
			Pattern: sdk.String(id.Name()),
		}))
		require.NoError(t, err)
		require.Len(t, networkRules, 1)

		assertThatObject(t, objectassert.NetworkRuleFromObject(t, &networkRules[0]).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("some comment").
			HasType(sdk.NetworkRuleTypeIpv4).
			HasMode(sdk.NetworkRuleModeIngress).
			HasEntriesInValueList(0).
			HasOwnerRoleType("ROLE"))
	})

	t.Run("Describe", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.NetworkRules.Create(ctx, sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{}, sdk.NetworkRuleModeIngress).WithComment("some comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().NetworkRule.DropFunc(t, id))

		assertThatObject(t, objectassert.NetworkRuleDetails(t, id).
			HasCreatedOnNotEmpty().
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("some comment").
			HasValueList([]string{}).
			HasMode(sdk.NetworkRuleModeIngress).
			HasType(sdk.NetworkRuleTypeIpv4))
	})
}

func TestInt_NetworkRulesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupNetworkRuleHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.NetworkRules.Drop(ctx, sdk.NewDropNetworkRuleRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createNetworkRuleHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		request := sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{}, sdk.NetworkRuleModeIngress)
		err := client.NetworkRules.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupNetworkRuleHandle(t, id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createNetworkRuleHandle(t, id1)
		createNetworkRuleHandle(t, id2)

		e1, err := client.NetworkRules.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.NetworkRules.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}

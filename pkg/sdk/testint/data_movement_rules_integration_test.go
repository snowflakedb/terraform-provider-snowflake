//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// Data movement rules need Enterprise Edition.
func TestInt_DataMovementRules(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	const defaultMaxRows = "(100)"
	const otherMaxRows = "(500)"

	createBasic := func(t *testing.T) sdk.SchemaObjectIdentifier {
		t.Helper()
		rule, cleanup := testClientHelper().DataMovementRule.Create(t)
		t.Cleanup(cleanup)
		return rule.ID()
	}

	t.Run("create: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.DataMovementRules.Create(ctx, sdk.NewCreateDataMovementRuleRequest(id, sdk.DataMovementRuleTypeCopyIntoExternalStage, defaultMaxRows))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().DataMovementRule.DropFunc(t, id))

		assertThatObject(t, objectassert.DataMovementRule(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment("").
			HasType(sdk.DataMovementRuleTypeCopyIntoExternalStage).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"),
		)
		assertThatObject(t, objectassert.DataMovementRuleDetails(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasType(sdk.DataMovementRuleTypeCopyIntoExternalStage).
			HasMaxRowsSignature("()").
			HasMaxRowsBody(defaultMaxRows),
		)
	})

	t.Run("create: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		maxRows := "CASE WHEN SYS_CONTEXT('SNOWFLAKE$SESSION', 'ROLE') = 'COMP_TEAM_ROLE' THEN NULL ELSE 0 END"

		err := client.DataMovementRules.Create(ctx, sdk.NewCreateDataMovementRuleRequest(id, sdk.DataMovementRuleTypeSnowsightUi, maxRows).
			WithIfNotExists(true).
			WithComment(comment))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().DataMovementRule.DropFunc(t, id))

		assertThatObject(t, objectassert.DataMovementRule(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment(comment).
			HasType(sdk.DataMovementRuleTypeSnowsightUi).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"),
		)
		assertThatObject(t, objectassert.DataMovementRuleDetails(t, id).
			HasName(id.Name()).
			HasComment(comment).
			HasType(sdk.DataMovementRuleTypeSnowsightUi).
			HasMaxRowsSignature("()").
			HasMaxRowsBody(maxRows),
		)

		err = client.DataMovementRules.Create(ctx, sdk.NewCreateDataMovementRuleRequest(id, sdk.DataMovementRuleTypeSnowsightUi, otherMaxRows).
			WithIfNotExists(true))
		require.NoError(t, err)
		assertThatObject(t, objectassert.DataMovementRuleDetails(t, id).
			HasMaxRowsBody(maxRows),
		)
	})

	t.Run("create: or replace", func(t *testing.T) {
		id := createBasic(t)

		comment := random.Comment()
		err := client.DataMovementRules.Create(ctx, sdk.NewCreateDataMovementRuleRequest(id, sdk.DataMovementRuleTypeUiDownload, otherMaxRows).
			WithOrReplace(true).
			WithComment(comment))
		require.NoError(t, err)

		assertThatObject(t, objectassert.DataMovementRule(t, id).
			HasComment(comment).
			HasType(sdk.DataMovementRuleTypeUiDownload),
		)
		assertThatObject(t, objectassert.DataMovementRuleDetails(t, id).
			HasMaxRowsBody(otherMaxRows),
		)
	})

	t.Run("show: like, in, and limit", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		prefix := random.AlphaN(6)
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID())
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID())
		id3 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		for _, id := range []sdk.SchemaObjectIdentifier{id1, id2, id3} {
			_, cleanup := testClientHelper().DataMovementRule.CreateWithRequest(t, sdk.NewCreateDataMovementRuleRequest(id, sdk.DataMovementRuleTypeProgrammaticFetch, defaultMaxRows))
			t.Cleanup(cleanup)
		}

		byLike, err := client.DataMovementRules.Show(ctx, sdk.NewShowDataMovementRuleRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id1.Name())}).
			WithIn(sdk.In{Schema: schema.ID()}))
		require.NoError(t, err)
		require.Len(t, byLike, 1)
		assertThatObject(t, objectassert.DataMovementRuleFromObject(t, &byLike[0]).
			HasName(id1.Name()).
			HasDatabaseName(id1.DatabaseName()).
			HasSchemaName(id1.SchemaName()).
			HasType(sdk.DataMovementRuleTypeProgrammaticFetch),
		)

		byLimit, err := client.DataMovementRules.Show(ctx, sdk.NewShowDataMovementRuleRequest().
			WithIn(sdk.In{Schema: schema.ID()}).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		require.Len(t, byLimit, 1)

		inSchema, err := client.DataMovementRules.Show(ctx, sdk.NewShowDataMovementRuleRequest().
			WithIn(sdk.In{Schema: schema.ID()}))
		require.NoError(t, err)
		require.Len(t, inSchema, 3)

		inDatabase, err := client.DataMovementRules.Show(ctx, sdk.NewShowDataMovementRuleRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id3.Name())}).
			WithIn(sdk.In{Database: id3.DatabaseId()}))
		require.NoError(t, err)
		require.Len(t, inDatabase, 1)

		inAccount, err := client.DataMovementRules.Show(ctx, sdk.NewShowDataMovementRuleRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id3.Name())}).
			WithIn(sdk.In{Account: sdk.Bool(true)}))
		require.NoError(t, err)
		require.Len(t, inAccount, 1)
	})

	t.Run("show by id: same name in two schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())
		_, cleanup1 := testClientHelper().DataMovementRule.CreateWithRequest(t, sdk.NewCreateDataMovementRuleRequest(id1, sdk.DataMovementRuleTypeSnowsightUi, defaultMaxRows))
		t.Cleanup(cleanup1)
		_, cleanup2 := testClientHelper().DataMovementRule.CreateWithRequest(t, sdk.NewCreateDataMovementRuleRequest(id2, sdk.DataMovementRuleTypeAgentAccess, defaultMaxRows).WithComment("second"))
		t.Cleanup(cleanup2)

		assertThatObject(t, objectassert.DataMovementRule(t, id1).
			HasName(id1.Name()).
			HasSchemaName(id1.SchemaName()).
			HasType(sdk.DataMovementRuleTypeSnowsightUi).
			HasComment(""),
		)
		assertThatObject(t, objectassert.DataMovementRule(t, id2).
			HasName(id2.Name()).
			HasSchemaName(id2.SchemaName()).
			HasType(sdk.DataMovementRuleTypeAgentAccess).
			HasComment("second"),
		)
	})

	t.Run("show by id safely: existing and non-existing", func(t *testing.T) {
		id := createBasic(t)

		rule, err := client.DataMovementRules.ShowByIDSafely(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectassert.DataMovementRuleFromObject(t, rule).
			HasName(id.Name()),
		)

		_, err = client.DataMovementRules.ShowByIDSafely(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("describe: existing", func(t *testing.T) {
		id := createBasic(t)

		assertThatObject(t, objectassert.DataMovementRuleDetails(t, id).
			HasName(id.Name()).
			HasType(sdk.DataMovementRuleTypeCopyIntoExternalStage).
			HasMaxRowsSignature("()").
			HasMaxRowsBody(defaultMaxRows),
		)
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		_, err := client.DataMovementRules.Describe(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter: set max rows", func(t *testing.T) {
		id := createBasic(t)

		err := client.DataMovementRules.Alter(ctx, sdk.NewAlterDataMovementRuleRequest(id).
			WithIfExists(true).
			WithSet(*sdk.NewDataMovementRuleSetRequest().WithSetMaxRows(otherMaxRows)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.DataMovementRuleDetails(t, id).
			HasMaxRowsBody(otherMaxRows),
		)
	})

	t.Run("alter: set and unset comment", func(t *testing.T) {
		id := createBasic(t)
		comment := random.Comment()

		err := client.DataMovementRules.Alter(ctx, sdk.NewAlterDataMovementRuleRequest(id).
			WithSet(*sdk.NewDataMovementRuleSetRequest().WithComment(comment)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.DataMovementRule(t, id).
			HasComment(comment),
		)

		err = client.DataMovementRules.Alter(ctx, sdk.NewAlterDataMovementRuleRequest(id).
			WithUnset(*sdk.NewDataMovementRuleUnsetRequest().WithComment(true)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.DataMovementRule(t, id).
			HasComment(""),
		)
	})

	t.Run("alter: if exists on non-existing", func(t *testing.T) {
		err := client.DataMovementRules.Alter(ctx, sdk.NewAlterDataMovementRuleRequest(NonExistingSchemaObjectIdentifier).
			WithIfExists(true).
			WithSet(*sdk.NewDataMovementRuleSetRequest().WithComment("x")))
		require.NoError(t, err)
	})

	t.Run("alter: rename", func(t *testing.T) {
		id := createBasic(t)
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.DataMovementRules.Alter(ctx, sdk.NewAlterDataMovementRuleRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().DataMovementRule.DropFunc(t, newId))

		_, err = client.DataMovementRules.ShowByID(ctx, id)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)

		assertThatObject(t, objectassert.DataMovementRule(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasType(sdk.DataMovementRuleTypeCopyIntoExternalStage),
		)
	})

	t.Run("alter: rename to another schema", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)
		id := createBasic(t)
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

		err := client.DataMovementRules.Alter(ctx, sdk.NewAlterDataMovementRuleRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().DataMovementRule.DropFunc(t, newId))

		assertThatObject(t, objectassert.DataMovementRule(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()),
		)
	})

	t.Run("drop: existing", func(t *testing.T) {
		id := createBasic(t)

		err := client.DataMovementRules.Drop(ctx, sdk.NewDropDataMovementRuleRequest(id))
		require.NoError(t, err)
		_, err = client.DataMovementRules.ShowByID(ctx, id)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		err := client.DataMovementRules.Drop(ctx, sdk.NewDropDataMovementRuleRequest(NonExistingSchemaObjectIdentifier))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists on non-existing", func(t *testing.T) {
		err := client.DataMovementRules.Drop(ctx, sdk.NewDropDataMovementRuleRequest(NonExistingSchemaObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop safely: existing", func(t *testing.T) {
		id := createBasic(t)

		err := client.DataMovementRules.DropSafely(ctx, id)
		require.NoError(t, err)
		_, err = client.DataMovementRules.ShowByIDSafely(ctx, id)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop safely: non-existing", func(t *testing.T) {
		err := client.DataMovementRules.DropSafely(ctx, NonExistingSchemaObjectIdentifier)
		require.NoError(t, err)
	})
}

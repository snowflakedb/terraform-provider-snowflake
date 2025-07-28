//go:build !account_level_tests

package testint

import (
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
 * todo: `ALTER SEQUENCE [ IF EXISTS ] <name> UNSET COMMENT` not works, and error: Syntax error: unexpected 'COMMENT'. (line 39)
 */

func TestInt_Sequences(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupSequenceHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.Sequences.Drop(ctx, sdk.NewDropSequenceRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createSequenceHandle := func(t *testing.T) *sdk.Sequence {
		t.Helper()

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		sr := sdk.NewCreateSequenceRequest(id).WithStart(sdk.Int(1)).WithIncrement(sdk.Int(1))
		err := client.Sequences.Create(ctx, sr)
		require.NoError(t, err)
		t.Cleanup(cleanupSequenceHandle(t, id))

		s, err := client.Sequences.ShowByID(ctx, id)
		require.NoError(t, err)
		return s
	}

	assertSequence := func(t *testing.T, id sdk.SchemaObjectIdentifier, interval int, ordered bool, comment string) {
		t.Helper()

		e, err := client.Sequences.ShowByID(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, e.CreatedOn)
		require.Equal(t, id.Name(), e.Name)
		require.Equal(t, id.DatabaseName(), e.DatabaseName)
		require.Equal(t, id.SchemaName(), e.SchemaName)
		require.Equal(t, 1, e.NextValue)
		require.Equal(t, interval, e.Interval)
		require.Equal(t, "ACCOUNTADMIN", e.Owner)
		require.Equal(t, "ROLE", e.OwnerRoleType)
		require.Equal(t, comment, e.Comment)
		require.Equal(t, ordered, e.Ordered)
	}

	t.Run("create sequence", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		comment := random.Comment()
		request := sdk.NewCreateSequenceRequest(id).
			WithStart(sdk.Int(1)).
			WithIncrement(sdk.Int(1)).
			WithIfNotExists(sdk.Bool(true)).
			WithValuesBehavior(sdk.ValuesBehaviorPointer(sdk.ValuesBehaviorOrder)).
			WithComment(&comment)
		err := client.Sequences.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSequenceHandle(t, id))
		assertSequence(t, id, 1, true, comment)
	})

	t.Run("show event table: without like", func(t *testing.T) {
		e1 := createSequenceHandle(t)
		e2 := createSequenceHandle(t)

		sequences, err := client.Sequences.Show(ctx, sdk.NewShowSequenceRequest())
		require.NoError(t, err)
		require.Len(t, sequences, 2)
		require.Contains(t, sequences, *e1)
		require.Contains(t, sequences, *e2)
	})

	t.Run("show sequence: with like", func(t *testing.T) {
		e1 := createSequenceHandle(t)
		e2 := createSequenceHandle(t)

		sequences, err := client.Sequences.Show(ctx, sdk.NewShowSequenceRequest().WithLike(sdk.Like{Pattern: &e1.Name}))
		require.NoError(t, err)
		require.Len(t, sequences, 1)
		require.Contains(t, sequences, *e1)
		require.NotContains(t, sequences, *e2)
	})

	t.Run("show sequence: no matches", func(t *testing.T) {
		sequences, err := client.Sequences.Show(ctx, sdk.NewShowSequenceRequest().WithLike(sdk.Like{Pattern: sdk.String("non-existent")}))
		require.NoError(t, err)
		require.Empty(t, sequences)
	})

	t.Run("describe sequence", func(t *testing.T) {
		e := createSequenceHandle(t)
		id := e.ID()

		details, err := client.Sequences.Describe(ctx, id)
		require.NoError(t, err)
		require.Equal(t, e.CreatedOn, details.CreatedOn)
		require.Equal(t, e.Name, details.Name)
		require.Equal(t, e.SchemaName, details.SchemaName)
		require.Equal(t, e.DatabaseName, details.DatabaseName)
		require.Equal(t, e.NextValue, details.NextValue)
		require.Equal(t, e.Interval, details.Interval)
		require.Equal(t, e.Owner, details.Owner)
		require.Equal(t, e.OwnerRoleType, details.OwnerRoleType)
		require.Equal(t, e.Comment, details.Comment)
		require.Equal(t, e.Ordered, details.Ordered)
	})

	t.Run("alter sequence: set options", func(t *testing.T) {
		e := createSequenceHandle(t)
		id := e.ID()

		comment := random.Comment()
		set := sdk.NewSequenceSetRequest().WithComment(&comment).WithValuesBehavior(sdk.ValuesBehaviorPointer(sdk.ValuesBehaviorNoOrder))
		err := client.Sequences.Alter(ctx, sdk.NewAlterSequenceRequest(id).WithSet(set))
		require.NoError(t, err)

		assertSequence(t, id, 1, false, comment)
	})

	t.Run("alter sequence: set increment", func(t *testing.T) {
		e := createSequenceHandle(t)
		id := e.ID()

		increment := 2
		err := client.Sequences.Alter(ctx, sdk.NewAlterSequenceRequest(id).WithSetIncrement(&increment))
		require.NoError(t, err)
		assertSequence(t, id, 2, false, "")
	})

	t.Run("alter sequence: rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Sequences.Create(ctx, sdk.NewCreateSequenceRequest(id))
		require.NoError(t, err)
		nid := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err = client.Sequences.Alter(ctx, sdk.NewAlterSequenceRequest(id).WithRenameTo(&nid))
		if err != nil {
			t.Cleanup(cleanupSequenceHandle(t, id))
		} else {
			t.Cleanup(cleanupSequenceHandle(t, nid))
		}
		require.NoError(t, err)

		_, err = client.Sequences.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
		_, err = client.Sequences.ShowByID(ctx, nid)
		require.NoError(t, err)
	})
}

func TestInt_SequencesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupSequenceHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.Sequences.Drop(ctx, sdk.NewDropSequenceRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createSequenceHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		sr := sdk.NewCreateSequenceRequest(id).WithStart(sdk.Int(1)).WithIncrement(sdk.Int(1))
		err := client.Sequences.Create(ctx, sr)
		require.NoError(t, err)
		t.Cleanup(cleanupSequenceHandle(t, id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createSequenceHandle(t, id1)
		createSequenceHandle(t, id2)

		e1, err := client.Sequences.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Sequences.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}

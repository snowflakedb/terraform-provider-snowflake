//go:build non_account_level_tests

package testint

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_MaskingPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	defaultBody := "REPLACE('X', 1, 2)"

	singleVarcharSignature := []sdk.CreateMaskingPolicySignatureRequest{
		*sdk.NewCreateMaskingPolicySignatureRequest("col1", testdatatypes.DataTypeVarchar),
	}
	expectedSingleVarcharSignature := []sdk.TableColumnSignature{
		{Name: "col1", Type: testdatatypes.DataTypeVarchar},
	}

	twoVarcharSignature := []sdk.CreateMaskingPolicySignatureRequest{
		*sdk.NewCreateMaskingPolicySignatureRequest("col1", testdatatypes.DataTypeVarchar),
		*sdk.NewCreateMaskingPolicySignatureRequest("col2", testdatatypes.DataTypeVarchar),
	}
	expectedTwoVarcharSignature := []sdk.TableColumnSignature{
		{Name: "col1", Type: testdatatypes.DataTypeVarchar},
		{Name: "col2", Type: testdatatypes.DataTypeVarchar},
	}

	// creates a minimal masking policy with the given id and registers its cleanup
	createMaskingPolicyWithId := func(t *testing.T, id sdk.SchemaObjectIdentifier) *sdk.MaskingPolicy {
		t.Helper()
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateOrReplaceMaskingPolicyWithRequest(t, id, singleVarcharSignature, testdatatypes.DataTypeVarchar, defaultBody)
		t.Cleanup(maskingPolicyCleanup)
		return maskingPolicy
	}

	// creates a minimal masking policy in the test schema and registers its cleanup
	createMaskingPolicy := func(t *testing.T) *sdk.MaskingPolicy {
		t.Helper()
		return createMaskingPolicyWithId(t, testClientHelper().Ids.RandomSchemaObjectIdentifier())
	}

	t.Run("create: complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		exemptOtherPolicies := random.Bool()

		err := client.MaskingPolicies.Create(ctx, sdk.NewCreateMaskingPolicyRequest(id, twoVarcharSignature, testdatatypes.DataTypeVarchar, defaultBody).
			WithOrReplace(true).
			WithComment(comment).
			WithExemptOtherPolicies(exemptOtherPolicies))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().MaskingPolicy.DropMaskingPolicyFunc(t, id))

		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), maskingPolicyDetails.Name)
		assert.Equal(t, expectedTwoVarcharSignature, maskingPolicyDetails.Signature)
		assert.Equal(t, testdatatypes.DefaultVarcharAsString, maskingPolicyDetails.ReturnType.ToSql())
		assert.Equal(t, defaultBody, maskingPolicyDetails.Body)

		maskingPolicy, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
		require.NoError(t, err)
		assert.Len(t, maskingPolicy, 1)
		assert.Equal(t, id.Name(), maskingPolicy[0].Name)
		assert.Equal(t, comment, maskingPolicy[0].Comment)
		assert.Equal(t, exemptOtherPolicies, maskingPolicy[0].ExemptOtherPolicies)
	})

	t.Run("create: if not exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()

		err := client.MaskingPolicies.Create(ctx, sdk.NewCreateMaskingPolicyRequest(id, twoVarcharSignature, testdatatypes.DataTypeVarchar, defaultBody).
			WithIfNotExists(true).
			WithComment(comment).
			WithExemptOtherPolicies(true))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().MaskingPolicy.DropMaskingPolicyFunc(t, id))

		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), maskingPolicyDetails.Name)
		assert.Equal(t, expectedTwoVarcharSignature, maskingPolicyDetails.Signature)
		assert.Equal(t, testdatatypes.DefaultVarcharAsString, maskingPolicyDetails.ReturnType.ToSql())
		assert.Equal(t, defaultBody, maskingPolicyDetails.Body)

		maskingPolicy, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
		require.NoError(t, err)
		assert.Len(t, maskingPolicy, 1)
		assert.Equal(t, id.Name(), maskingPolicy[0].Name)
		assert.Equal(t, comment, maskingPolicy[0].Comment)
		assert.True(t, maskingPolicy[0].ExemptOtherPolicies)
	})

	t.Run("create: no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.MaskingPolicies.Create(ctx, sdk.NewCreateMaskingPolicyRequest(id, singleVarcharSignature, testdatatypes.DataTypeVarchar, defaultBody))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().MaskingPolicy.DropMaskingPolicyFunc(t, id))

		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), maskingPolicyDetails.Name)
		assert.Equal(t, expectedSingleVarcharSignature, maskingPolicyDetails.Signature)
		assert.Equal(t, testdatatypes.DefaultVarcharAsString, maskingPolicyDetails.ReturnType.ToSql())
		assert.Equal(t, defaultBody, maskingPolicyDetails.Body)

		maskingPolicy, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
		require.NoError(t, err)
		assert.Len(t, maskingPolicy, 1)
		assert.Equal(t, id.Name(), maskingPolicy[0].Name)
		assert.Equal(t, "", maskingPolicy[0].Comment)
		assert.False(t, maskingPolicy[0].ExemptOtherPolicies)
	})

	t.Run("create: DECFLOAT", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		signature := []sdk.CreateMaskingPolicySignatureRequest{
			*sdk.NewCreateMaskingPolicySignatureRequest("col1", testdatatypes.DataTypeDecfloat),
		}
		expectedSignature := []sdk.TableColumnSignature{
			{Name: "col1", Type: testdatatypes.DataTypeDecfloat},
		}
		expression := "REPLACE('X', 1, 2)::DECFLOAT"

		err := client.MaskingPolicies.Create(ctx, sdk.NewCreateMaskingPolicyRequest(id, signature, testdatatypes.DataTypeDecfloat, expression))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().MaskingPolicy.DropMaskingPolicyFunc(t, id))

		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), maskingPolicyDetails.Name)
		assert.Equal(t, expectedSignature, maskingPolicyDetails.Signature)
		assert.Equal(t, "DECFLOAT(38)", maskingPolicyDetails.ReturnType.ToSqlWithoutUnknowns())
		assert.Equal(t, expression, maskingPolicyDetails.Body)
	})

	t.Run("create: multiline expression", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		signature := []sdk.CreateMaskingPolicySignatureRequest{
			*sdk.NewCreateMaskingPolicySignatureRequest("val", testdatatypes.DataTypeVarchar),
		}
		expectedSignature := []sdk.TableColumnSignature{
			{Name: "val", Type: testdatatypes.DataTypeVarchar},
		}
		expression := `
		case
			when current_role() in ('ROLE_A') then
				val
			when is_role_in_session( 'ROLE_B' ) then
				'ABC123'
			else
				'******'
		end
		`

		err := client.MaskingPolicies.Create(ctx, sdk.NewCreateMaskingPolicyRequest(id, signature, testdatatypes.DataTypeVarchar, expression))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().MaskingPolicy.DropMaskingPolicyFunc(t, id))

		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), maskingPolicyDetails.Name)
		assert.Equal(t, expectedSignature, maskingPolicyDetails.Signature)
		assert.Equal(t, testdatatypes.DefaultVarcharAsString, maskingPolicyDetails.ReturnType.ToSql())
		assert.Equal(t, strings.TrimSpace(expression), maskingPolicyDetails.Body)
	})

	t.Run("alter: set and unset comment", func(t *testing.T) {
		maskingPolicy := createMaskingPolicy(t)
		id := maskingPolicy.ID()
		comment := random.Comment()

		err := client.MaskingPolicies.Alter(ctx, sdk.NewAlterMaskingPolicyRequest(id).
			WithSetComment(comment))
		require.NoError(t, err)
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
		require.NoError(t, err)
		assert.Len(t, maskingPolicies, 1)
		assert.Equal(t, comment, maskingPolicies[0].Comment)

		err = client.MaskingPolicies.Alter(ctx, sdk.NewAlterMaskingPolicyRequest(id).
			WithSetComment(comment))
		require.NoError(t, err)
		err = client.MaskingPolicies.Alter(ctx, sdk.NewAlterMaskingPolicyRequest(id).
			WithUnsetComment(true))
		require.NoError(t, err)
		maskingPolicies, err = client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
		require.NoError(t, err)
		assert.Len(t, maskingPolicies, 1)
		assert.Equal(t, "", maskingPolicies[0].Comment)
	})

	t.Run("alter: rename", func(t *testing.T) {
		maskingPolicy := createMaskingPolicy(t)
		oldId := maskingPolicy.ID()
		// a masking policy can be renamed only within the same schema
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(oldId.SchemaId())

		err := client.MaskingPolicies.Alter(ctx, sdk.NewAlterMaskingPolicyRequest(oldId).WithRenameTo(newId))
		if err == nil {
			t.Cleanup(testClientHelper().MaskingPolicy.DropMaskingPolicyFunc(t, newId))
		}
		require.NoError(t, err)

		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, newId)
		require.NoError(t, err)
		assert.Equal(t, newId.Name(), maskingPolicyDetails.Name)
	})

	t.Run("alter: set body", func(t *testing.T) {
		maskingPolicy := createMaskingPolicy(t)
		id := maskingPolicy.ID()
		newBody := "'***'"

		err := client.MaskingPolicies.Alter(ctx, sdk.NewAlterMaskingPolicyRequest(id).
			WithSetBody(newBody))
		require.NoError(t, err)

		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, newBody, maskingPolicyDetails.Body)
	})

	t.Run("drop: existing", func(t *testing.T) {
		maskingPolicy := createMaskingPolicy(t)
		id := maskingPolicy.ID()

		err := client.MaskingPolicies.Drop(ctx, sdk.NewDropMaskingPolicyRequest(id))
		require.NoError(t, err)

		_, err = client.MaskingPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		err := client.MaskingPolicies.Drop(ctx, sdk.NewDropMaskingPolicyRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	// SHOW MASKING POLICIES without IN is scoped to the schema currently in use, so this case has to be run
	// before the "show" subtest below, which changes the current schema by creating a new database and schema.
	t.Run("show: without options", func(t *testing.T) {
		maskingPolicy := createMaskingPolicy(t)
		maskingPolicy2 := createMaskingPolicy(t)

		maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest())
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(maskingPolicies), 2)
		assert.Contains(t, maskingPolicies, *maskingPolicy)
		assert.Contains(t, maskingPolicies, *maskingPolicy2)
	})

	t.Run("show", func(t *testing.T) {
		// the fixtures are created in a freshly created database, so that the exact counts below are not affected
		// by masking policies created in the test schema by other tests
		db, dbCleanup := testClientHelper().Database.CreateDatabaseWithParametersSet(t)
		t.Cleanup(dbCleanup)
		schemaId := testClientHelper().Ids.NewDatabaseObjectIdentifierInDatabase("PUBLIC", db.ID())
		otherSchema, otherSchemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
		t.Cleanup(otherSchemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix("test_masking_policy_1_", schemaId)
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix("test_masking_policy_2_", schemaId)
		id3 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix("test_masking_policy_3_", schemaId)
		id4 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(otherSchema.ID())

		maskingPolicy1 := createMaskingPolicyWithId(t, id1)
		maskingPolicy2 := createMaskingPolicyWithId(t, id2)
		maskingPolicy3 := createMaskingPolicyWithId(t, id3)
		maskingPolicy4 := createMaskingPolicyWithId(t, id4)

		t.Run("like", func(t *testing.T) {
			maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
				WithLike(sdk.Like{Pattern: sdk.String("test_masking_policy_2_%")}).
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: schemaId}}))
			require.NoError(t, err)
			assert.Len(t, maskingPolicies, 1)
			assert.Contains(t, maskingPolicies, *maskingPolicy2)
		})

		t.Run("like: exact name", func(t *testing.T) {
			maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
				WithLike(sdk.Like{Pattern: sdk.String(id1.Name())}).
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: schemaId}}))
			require.NoError(t, err)
			assert.Len(t, maskingPolicies, 1)
			assert.Contains(t, maskingPolicies, *maskingPolicy1)
		})

		t.Run("like: non-existing", func(t *testing.T) {
			maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
				WithLike(sdk.Like{Pattern: sdk.String("non-existent")}))
			require.NoError(t, err)
			assert.Empty(t, maskingPolicies)
		})

		t.Run("in_account", func(t *testing.T) {
			maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
				WithIn(sdk.ExtendedIn{In: sdk.In{Account: sdk.Bool(true)}}))
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(maskingPolicies), 4)
		})

		t.Run("in_database", func(t *testing.T) {
			maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
				WithIn(sdk.ExtendedIn{In: sdk.In{Database: db.ID()}}))
			require.NoError(t, err)
			assert.Len(t, maskingPolicies, 4)
			assert.Contains(t, maskingPolicies, *maskingPolicy4)
		})

		t.Run("in_schema", func(t *testing.T) {
			maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: schemaId}}))
			require.NoError(t, err)
			assert.Len(t, maskingPolicies, 3)
			assert.Contains(t, maskingPolicies, *maskingPolicy1)
			assert.Contains(t, maskingPolicies, *maskingPolicy2)
			assert.Contains(t, maskingPolicies, *maskingPolicy3)
		})

		t.Run("limit", func(t *testing.T) {
			maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: schemaId}}).
				WithLimit(sdk.LimitFrom{Rows: sdk.Pointer(1)}))
			require.NoError(t, err)
			assert.Len(t, maskingPolicies, 1)
		})
	})

	t.Run("describe: existing", func(t *testing.T) {
		maskingPolicy := createMaskingPolicy(t)

		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, maskingPolicy.ID())
		require.NoError(t, err)
		assert.Equal(t, maskingPolicy.Name, maskingPolicyDetails.Name)
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		_, err := client.MaskingPolicies.Describe(ctx, NonExistingSchemaObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("show by id: same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createMaskingPolicyWithId(t, id1)
		createMaskingPolicyWithId(t, id2)

		e1, err := client.MaskingPolicies.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.MaskingPolicies.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})

	t.Run("show by id: check fields", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createMaskingPolicyWithId(t, id)

		mp, err := client.MaskingPolicies.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, mp.ID())
		assert.NotEmpty(t, mp.CreatedOn)
		assert.Equal(t, id.Name(), mp.Name)
		assert.Equal(t, testClientHelper().Ids.DatabaseId().Name(), mp.DatabaseName)
		assert.Equal(t, testClientHelper().Ids.SchemaId().Name(), mp.SchemaName)
		assert.Equal(t, "MASKING_POLICY", mp.Kind)
		assert.Equal(t, "ACCOUNTADMIN", mp.Owner)
		assert.Equal(t, "", mp.Comment)
		assert.False(t, mp.ExemptOtherPolicies)
		assert.Equal(t, "ROLE", mp.OwnerRoleType)
	})
}

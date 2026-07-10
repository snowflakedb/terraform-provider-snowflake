//go:build non_account_level_tests

package testint

import (
	"errors"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-2298256]: merge these tests
func TestInt_MaskingPoliciesShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	maskingPolicyTest, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
	t.Cleanup(maskingPolicyCleanup)

	maskingPolicy2Test, maskingPolicy2Cleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
	t.Cleanup(maskingPolicy2Cleanup)

	t.Run("without show options", func(t *testing.T) {
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest())
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(maskingPolicies), 2)
	})

	t.Run("with show options", func(t *testing.T) {
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}))
		require.NoError(t, err)
		assert.Contains(t, maskingPolicies, *maskingPolicyTest)
		assert.Contains(t, maskingPolicies, *maskingPolicy2Test)
		assert.Len(t, maskingPolicies, 2)
	})

	t.Run("with show options and like", func(t *testing.T) {
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(maskingPolicyTest.Name)}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}))
		require.NoError(t, err)
		assert.Contains(t, maskingPolicies, *maskingPolicyTest)
		assert.Len(t, maskingPolicies, 1)
	})

	t.Run("when searching a non-existent masking policy", func(t *testing.T) {
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String("non-existent")}))
		require.NoError(t, err)
		assert.Empty(t, maskingPolicies)
	})

	t.Run("when limiting the number of results", func(t *testing.T) {
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}).
			WithLimit(sdk.LimitFrom{Rows: sdk.Pointer(1)}))
		require.NoError(t, err)
		assert.Len(t, maskingPolicies, 1)
	})
}

func TestInt_MaskingPolicyCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("test complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := id.Name()
		signature := []sdk.CreateMaskingPolicySignatureRequest{
			*sdk.NewCreateMaskingPolicySignatureRequest("col1", testdatatypes.DataTypeVarchar),
			*sdk.NewCreateMaskingPolicySignatureRequest("col2", testdatatypes.DataTypeVarchar),
		}
		expectedSignature := []sdk.TableColumnSignature{
			{Name: "col1", Type: testdatatypes.DataTypeVarchar},
			{Name: "col2", Type: testdatatypes.DataTypeVarchar},
		}
		expression := "REPLACE('X', 1, 2)"
		comment := random.Comment()
		exemptOtherPolicies := random.Bool()
		err := client.MaskingPolicies.Create(ctx, sdk.NewCreateMaskingPolicyRequest(id, signature, testdatatypes.DataTypeVarchar, expression).
			WithOrReplace(true).
			WithComment(comment).
			WithExemptOtherPolicies(exemptOtherPolicies))
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, maskingPolicyDetails.Name)
		assert.Equal(t, expectedSignature, maskingPolicyDetails.Signature)
		assert.Equal(t, testdatatypes.DefaultVarcharAsString, maskingPolicyDetails.ReturnType.ToSql())
		assert.Equal(t, expression, maskingPolicyDetails.Body)

		maskingPolicy, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(name)}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}))
		require.NoError(t, err)
		assert.Len(t, maskingPolicy, 1)
		assert.Equal(t, name, maskingPolicy[0].Name)
		assert.Equal(t, comment, maskingPolicy[0].Comment)
		assert.Equal(t, exemptOtherPolicies, maskingPolicy[0].ExemptOtherPolicies)
	})

	t.Run("test if_not_exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := id.Name()
		signature := []sdk.CreateMaskingPolicySignatureRequest{
			*sdk.NewCreateMaskingPolicySignatureRequest("col1", testdatatypes.DataTypeVarchar),
			*sdk.NewCreateMaskingPolicySignatureRequest("col2", testdatatypes.DataTypeVarchar),
		}
		expectedSignature := []sdk.TableColumnSignature{
			{Name: "col1", Type: testdatatypes.DataTypeVarchar},
			{Name: "col2", Type: testdatatypes.DataTypeVarchar},
		}
		expression := "REPLACE('X', 1, 2)"
		comment := random.Comment()
		err := client.MaskingPolicies.Create(ctx, sdk.NewCreateMaskingPolicyRequest(id, signature, testdatatypes.DataTypeVarchar, expression).
			WithIfNotExists(true).
			WithComment(comment).
			WithExemptOtherPolicies(true))
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, maskingPolicyDetails.Name)
		assert.Equal(t, expectedSignature, maskingPolicyDetails.Signature)
		assert.Equal(t, testdatatypes.DefaultVarcharAsString, maskingPolicyDetails.ReturnType.ToSql())
		assert.Equal(t, expression, maskingPolicyDetails.Body)

		maskingPolicy, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(name)}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}))
		require.NoError(t, err)
		assert.Len(t, maskingPolicy, 1)
		assert.Equal(t, name, maskingPolicy[0].Name)
		assert.Equal(t, comment, maskingPolicy[0].Comment)
		assert.True(t, maskingPolicy[0].ExemptOtherPolicies)
	})

	t.Run("test no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := id.Name()
		signature := []sdk.CreateMaskingPolicySignatureRequest{
			*sdk.NewCreateMaskingPolicySignatureRequest("col1", testdatatypes.DataTypeVarchar),
		}
		expectedSignature := []sdk.TableColumnSignature{
			{Name: "col1", Type: testdatatypes.DataTypeVarchar},
		}
		expression := "REPLACE('X', 1, 2)"
		err := client.MaskingPolicies.Create(ctx, sdk.NewCreateMaskingPolicyRequest(id, signature, testdatatypes.DataTypeVarchar, expression))
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, maskingPolicyDetails.Name)
		assert.Equal(t, expectedSignature, maskingPolicyDetails.Signature)
		assert.Equal(t, testdatatypes.DefaultVarcharAsString, maskingPolicyDetails.ReturnType.ToSql())
		assert.Equal(t, expression, maskingPolicyDetails.Body)

		maskingPolicy, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(name)}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}))
		require.NoError(t, err)
		assert.Len(t, maskingPolicy, 1)
		assert.Equal(t, name, maskingPolicy[0].Name)
		assert.Equal(t, "", maskingPolicy[0].Comment)
		assert.False(t, maskingPolicy[0].ExemptOtherPolicies)
	})

	t.Run("create: DECFLOAT", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := id.Name()
		signature := []sdk.CreateMaskingPolicySignatureRequest{
			*sdk.NewCreateMaskingPolicySignatureRequest("col1", testdatatypes.DataTypeDecfloat),
		}
		expectedSignature := []sdk.TableColumnSignature{
			{Name: "col1", Type: testdatatypes.DataTypeDecfloat},
		}
		expression := "REPLACE('X', 1, 2)::DECFLOAT"
		err := client.MaskingPolicies.Create(ctx, sdk.NewCreateMaskingPolicyRequest(id, signature, testdatatypes.DataTypeDecfloat, expression))
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, maskingPolicyDetails.Name)
		assert.Equal(t, expectedSignature, maskingPolicyDetails.Signature)
		assert.Equal(t, "DECFLOAT(38)", maskingPolicyDetails.ReturnType.ToSqlWithoutUnknowns())
		assert.Equal(t, expression, maskingPolicyDetails.Body)
	})

	t.Run("test multiline expression", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := id.Name()
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
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, maskingPolicyDetails.Name)
		assert.Equal(t, expectedSignature, maskingPolicyDetails.Signature)
		assert.Equal(t, testdatatypes.DefaultVarcharAsString, maskingPolicyDetails.ReturnType.ToSql())
		assert.Equal(t, strings.TrimSpace(expression), maskingPolicyDetails.Body)
	})
}

func TestInt_MaskingPolicyDescribe(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
	t.Cleanup(maskingPolicyCleanup)

	t.Run("when masking policy exists", func(t *testing.T) {
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, maskingPolicy.ID())
		require.NoError(t, err)
		assert.Equal(t, maskingPolicy.Name, maskingPolicyDetails.Name)
	})

	t.Run("when masking policy does not exist", func(t *testing.T) {
		_, err := client.MaskingPolicies.Describe(ctx, NonExistingSchemaObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_MaskingPolicyAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when setting and unsetting a value", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		t.Cleanup(maskingPolicyCleanup)
		comment := random.Comment()
		err := client.MaskingPolicies.Alter(ctx, sdk.NewAlterMaskingPolicyRequest(maskingPolicy.ID()).
			WithSetComment(comment))
		require.NoError(t, err)
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(maskingPolicy.Name)}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}))
		require.NoError(t, err)
		assert.Len(t, maskingPolicies, 1)
		assert.Equal(t, comment, maskingPolicies[0].Comment)

		err = client.MaskingPolicies.Alter(ctx, sdk.NewAlterMaskingPolicyRequest(maskingPolicy.ID()).
			WithSetComment(comment))
		require.NoError(t, err)
		err = client.MaskingPolicies.Alter(ctx, sdk.NewAlterMaskingPolicyRequest(maskingPolicy.ID()).
			WithUnsetComment(true))
		require.NoError(t, err)
		maskingPolicies, err = client.MaskingPolicies.Show(ctx, sdk.NewShowMaskingPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(maskingPolicy.Name)}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}))
		require.NoError(t, err)
		assert.Len(t, maskingPolicies, 1)
		assert.Equal(t, "", maskingPolicies[0].Comment)
	})

	t.Run("when renaming", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		oldID := maskingPolicy.ID()
		t.Cleanup(maskingPolicyCleanup)
		newID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.MaskingPolicies.Alter(ctx, sdk.NewAlterMaskingPolicyRequest(oldID).WithNewName(newID))
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newID.Name(), maskingPolicyDetails.Name)
		// rename back to original name, so it can be cleaned up
		err = client.MaskingPolicies.Alter(ctx, sdk.NewAlterMaskingPolicyRequest(newID).WithNewName(oldID))
		require.NoError(t, err)
	})

	t.Run("set body", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		id := maskingPolicy.ID()
		newBody := "'***'"
		t.Cleanup(maskingPolicyCleanup)

		err := client.MaskingPolicies.Alter(ctx, sdk.NewAlterMaskingPolicyRequest(id).
			WithSetBody(newBody))
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, newBody, maskingPolicyDetails.Body)
	})
}

func TestInt_MaskingPolicyDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when masking policy exists", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		t.Cleanup(maskingPolicyCleanup)
		id := maskingPolicy.ID()
		err := client.MaskingPolicies.Drop(ctx, sdk.NewDropMaskingPolicyRequest(id))
		require.NoError(t, err)
		_, err = client.MaskingPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("when masking policy does not exist", func(t *testing.T) {
		err := client.MaskingPolicies.Drop(ctx, sdk.NewDropMaskingPolicyRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_MaskingPoliciesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupMaskingPolicyHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.MaskingPolicies.Drop(ctx, sdk.NewDropMaskingPolicyRequest(id).WithIfExists(true))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createMaskingPolicyHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		signature := []sdk.CreateMaskingPolicySignatureRequest{
			*sdk.NewCreateMaskingPolicySignatureRequest(testClientHelper().Ids.Alpha(), testdatatypes.DataTypeVarchar),
		}
		expression := "REPLACE('X', 1, 2)"
		err := client.MaskingPolicies.Create(ctx, sdk.NewCreateMaskingPolicyRequest(id, signature, testdatatypes.DataTypeVarchar, expression))
		require.NoError(t, err)
		t.Cleanup(cleanupMaskingPolicyHandle(t, id))
	}

	assertMaskingPolicy := func(t *testing.T, mp *sdk.MaskingPolicy, id sdk.SchemaObjectIdentifier) {
		t.Helper()
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
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createMaskingPolicyHandle(t, id1)
		createMaskingPolicyHandle(t, id2)

		e1, err := client.MaskingPolicies.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.MaskingPolicies.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})

	t.Run("show by id: check fields", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		createMaskingPolicyHandle(t, id)

		mp, err := client.MaskingPolicies.ShowByID(ctx, id)
		require.NoError(t, err)
		assertMaskingPolicy(t, mp, id)
	})
}

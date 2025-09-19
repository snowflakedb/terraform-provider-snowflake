//go:build !account_level_tests

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
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, nil)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(maskingPolicies), 2)
	})

	t.Run("with show options", func(t *testing.T) {
		showOptions := &sdk.ShowMaskingPolicyOptions{
			In: &sdk.ExtendedIn{
				In: sdk.In{
					Schema: testClientHelper().Ids.SchemaId(),
				},
			},
		}
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, maskingPolicies, *maskingPolicyTest)
		assert.Contains(t, maskingPolicies, *maskingPolicy2Test)
		assert.Len(t, maskingPolicies, 2)
	})

	t.Run("with show options and like", func(t *testing.T) {
		showOptions := &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(maskingPolicyTest.Name),
			},
			In: &sdk.ExtendedIn{
				In: sdk.In{
					Schema: testClientHelper().Ids.SchemaId(),
				},
			},
		}
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, maskingPolicies, *maskingPolicyTest)
		assert.Len(t, maskingPolicies, 1)
	})

	t.Run("when searching a non-existent masking policy", func(t *testing.T) {
		showOptions := &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Empty(t, maskingPolicies)
	})

	t.Run("when limiting the number of results", func(t *testing.T) {
		showOptions := &sdk.ShowMaskingPolicyOptions{
			In: &sdk.ExtendedIn{
				In: sdk.In{
					Schema: testClientHelper().Ids.SchemaId(),
				},
			},
			Limit: &sdk.LimitFrom{
				Rows: sdk.Pointer(1),
			},
		}
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, showOptions)
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
		signature := []sdk.TableColumnSignature{
			{
				Name: "col1",
				Type: testdatatypes.DataTypeVarchar,
			},
			{
				Name: "col2",
				Type: testdatatypes.DataTypeVarchar,
			},
		}
		expression := "REPLACE('X', 1, 2)"
		comment := random.Comment()
		exemptOtherPolicies := random.Bool()
		err := client.MaskingPolicies.Create(ctx, id, signature, testdatatypes.DataTypeVarchar, expression, &sdk.CreateMaskingPolicyOptions{
			OrReplace:           sdk.Bool(true),
			IfNotExists:         sdk.Bool(false),
			Comment:             sdk.String(comment),
			ExemptOtherPolicies: sdk.Bool(exemptOtherPolicies),
		})
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, maskingPolicyDetails.Name)
		assert.Equal(t, signature, maskingPolicyDetails.Signature)
		assert.Equal(t, testdatatypes.DefaultVarcharAsString, maskingPolicyDetails.ReturnType.ToSql())
		assert.Equal(t, expression, maskingPolicyDetails.Body)

		maskingPolicy, err := client.MaskingPolicies.Show(ctx, &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.ExtendedIn{
				In: sdk.In{
					Schema: testClientHelper().Ids.SchemaId(),
				},
			},
		})
		require.NoError(t, err)
		assert.Len(t, maskingPolicy, 1)
		assert.Equal(t, name, maskingPolicy[0].Name)
		assert.Equal(t, comment, maskingPolicy[0].Comment)
		assert.Equal(t, exemptOtherPolicies, maskingPolicy[0].ExemptOtherPolicies)
	})

	t.Run("test if_not_exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := id.Name()
		signature := []sdk.TableColumnSignature{
			{
				Name: "col1",
				Type: testdatatypes.DataTypeVarchar,
			},
			{
				Name: "col2",
				Type: testdatatypes.DataTypeVarchar,
			},
		}
		expression := "REPLACE('X', 1, 2)"
		comment := random.Comment()
		err := client.MaskingPolicies.Create(ctx, id, signature, testdatatypes.DataTypeVarchar, expression, &sdk.CreateMaskingPolicyOptions{
			OrReplace:           sdk.Bool(false),
			IfNotExists:         sdk.Bool(true),
			Comment:             sdk.String(comment),
			ExemptOtherPolicies: sdk.Bool(true),
		})
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, maskingPolicyDetails.Name)
		assert.Equal(t, signature, maskingPolicyDetails.Signature)
		assert.Equal(t, testdatatypes.DefaultVarcharAsString, maskingPolicyDetails.ReturnType.ToSql())
		assert.Equal(t, expression, maskingPolicyDetails.Body)

		maskingPolicy, err := client.MaskingPolicies.Show(ctx, &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.ExtendedIn{
				In: sdk.In{
					Schema: testClientHelper().Ids.SchemaId(),
				},
			},
		})
		require.NoError(t, err)
		assert.Len(t, maskingPolicy, 1)
		assert.Equal(t, name, maskingPolicy[0].Name)
		assert.Equal(t, comment, maskingPolicy[0].Comment)
		assert.True(t, maskingPolicy[0].ExemptOtherPolicies)
	})

	t.Run("test no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := id.Name()
		signature := []sdk.TableColumnSignature{
			{
				Name: "col1",
				Type: testdatatypes.DataTypeVarchar,
			},
		}
		expression := "REPLACE('X', 1, 2)"
		err := client.MaskingPolicies.Create(ctx, id, signature, testdatatypes.DataTypeVarchar, expression, nil)
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, maskingPolicyDetails.Name)
		assert.Equal(t, signature, maskingPolicyDetails.Signature)
		assert.Equal(t, testdatatypes.DefaultVarcharAsString, maskingPolicyDetails.ReturnType.ToSql())
		assert.Equal(t, expression, maskingPolicyDetails.Body)

		maskingPolicy, err := client.MaskingPolicies.Show(ctx, &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.ExtendedIn{
				In: sdk.In{
					Schema: testClientHelper().Ids.SchemaId(),
				},
			},
		})
		require.NoError(t, err)
		assert.Len(t, maskingPolicy, 1)
		assert.Equal(t, name, maskingPolicy[0].Name)
		assert.Equal(t, "", maskingPolicy[0].Comment)
		assert.False(t, maskingPolicy[0].ExemptOtherPolicies)
	})

	t.Run("test multiline expression", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := id.Name()
		signature := []sdk.TableColumnSignature{
			{
				Name: "val",
				Type: testdatatypes.DataTypeVarchar,
			},
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
		err := client.MaskingPolicies.Create(ctx, id, signature, testdatatypes.DataTypeVarchar, expression, nil)
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, maskingPolicyDetails.Name)
		assert.Equal(t, signature, maskingPolicyDetails.Signature)
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
		alterOptions := &sdk.AlterMaskingPolicyOptions{
			Set: &sdk.MaskingPolicySet{
				Comment: sdk.String(comment),
			},
		}
		err := client.MaskingPolicies.Alter(ctx, maskingPolicy.ID(), alterOptions)
		require.NoError(t, err)
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(maskingPolicy.Name),
			},
			In: &sdk.ExtendedIn{
				In: sdk.In{
					Schema: testClientHelper().Ids.SchemaId(),
				},
			},
		})
		require.NoError(t, err)
		assert.Len(t, maskingPolicies, 1)
		assert.Equal(t, comment, maskingPolicies[0].Comment)

		err = client.MaskingPolicies.Alter(ctx, maskingPolicy.ID(), alterOptions)
		require.NoError(t, err)
		alterOptions = &sdk.AlterMaskingPolicyOptions{
			Unset: &sdk.MaskingPolicyUnset{
				Comment: sdk.Bool(true),
			},
		}
		err = client.MaskingPolicies.Alter(ctx, maskingPolicy.ID(), alterOptions)
		require.NoError(t, err)
		maskingPolicies, err = client.MaskingPolicies.Show(ctx, &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(maskingPolicy.Name),
			},
			In: &sdk.ExtendedIn{
				In: sdk.In{
					Schema: testClientHelper().Ids.SchemaId(),
				},
			},
		})
		require.NoError(t, err)
		assert.Len(t, maskingPolicies, 1)
		assert.Equal(t, "", maskingPolicies[0].Comment)
	})

	t.Run("when renaming", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		oldID := maskingPolicy.ID()
		t.Cleanup(maskingPolicyCleanup)
		newID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		alterOptions := &sdk.AlterMaskingPolicyOptions{
			NewName: &newID,
		}
		err := client.MaskingPolicies.Alter(ctx, oldID, alterOptions)
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newID.Name(), maskingPolicyDetails.Name)
		// rename back to original name, so it can be cleaned up
		alterOptions = &sdk.AlterMaskingPolicyOptions{
			NewName: &oldID,
		}
		err = client.MaskingPolicies.Alter(ctx, newID, alterOptions)
		require.NoError(t, err)
	})

	t.Run("set body", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		id := maskingPolicy.ID()
		newBody := "'***'"
		t.Cleanup(maskingPolicyCleanup)

		alterOptions := &sdk.AlterMaskingPolicyOptions{
			Set: &sdk.MaskingPolicySet{
				Body: sdk.Pointer(newBody),
			},
		}
		err := client.MaskingPolicies.Alter(ctx, id, alterOptions)
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
		err := client.MaskingPolicies.Drop(ctx, id, nil)
		require.NoError(t, err)
		_, err = client.MaskingPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("when masking policy does not exist", func(t *testing.T) {
		err := client.MaskingPolicies.Drop(ctx, NonExistingSchemaObjectIdentifier, nil)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_MaskingPoliciesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupMaskingPolicyHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.MaskingPolicies.Drop(ctx, id, &sdk.DropMaskingPolicyOptions{IfExists: sdk.Bool(true)})
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createMaskingPolicyHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		signature := []sdk.TableColumnSignature{
			{
				Name: testClientHelper().Ids.Alpha(),
				Type: testdatatypes.DataTypeVarchar,
			},
		}
		expression := "REPLACE('X', 1, 2)"
		err := client.MaskingPolicies.Create(ctx, id, signature, testdatatypes.DataTypeVarchar, expression, &sdk.CreateMaskingPolicyOptions{})
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

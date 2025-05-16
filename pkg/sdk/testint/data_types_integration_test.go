//go:build !account_level_tests

package testint

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DataTypes(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	incorrectBooleanDatatypes := []string{
		"BOOLEAN()",
		"BOOLEAN(1)",
		"BOOL",
	}
	incorrectFloatDatatypes := []string{
		"DOUBLE()",
		"DOUBLE(1)",
		"DOUBLE PRECISION(1)",
	}
	incorrectlyCorrectFloatDatatypes := []string{
		"FLOAT()",
		"FLOAT(20)",
		"FLOAT4(20)",
		"FLOAT8(20)",
		"REAL(20)",
	}
	incorrectNumberDatatypes := []string{
		"NUMBER()",
		"NUMBER(x)",
		"INT()",
		"NUMBER(36, 5, 7)",
	}
	incorrectTextDatatypes := []string{
		"VARCHAR()",
		"VARCHAR(x)",
		"VARCHAR(36, 5)",
	}
	vectorInnerTypesSynonyms := slices.Concat(datatypes.AllNumberDataTypes, datatypes.FloatDataTypeSynonyms)
	vectorInnerTypeSynonymsThatWork := []string{
		"INTEGER",
		"INT",
		"FLOAT8",
		"FLOAT4",
		"FLOAT",
	}

	for _, c := range datatypes.ArrayDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of array datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT []::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT []::%s(36)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.ErrorContains(t, err, "SQL compilation error")
			assert.ErrorContains(t, err, "unexpected '36'")
		})
	}

	for _, c := range datatypes.BinaryDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of binary datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT TO_BINARY('AB')::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT TO_BINARY('AB')::%s(36)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT TO_BINARY('AB')::%s(36, 2)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.ErrorContains(t, err, "SQL compilation error")
			assert.ErrorContains(t, err, "','")
			assert.ErrorContains(t, err, "')'")
		})
	}

	for _, c := range datatypes.BooleanDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of boolean datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT TRUE::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
		})
	}

	for _, c := range incorrectBooleanDatatypes {
		t.Run(fmt.Sprintf("check behavior of boolean datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT TRUE::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			require.Error(t, err)
		})
	}

	for _, c := range datatypes.DateDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of date datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT '2024-12-02'::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
		})
	}

	for _, c := range datatypes.FloatDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of float datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1.1::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
		})
	}

	for _, c := range incorrectFloatDatatypes {
		t.Run(fmt.Sprintf("check behavior of float datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1.1::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			require.Error(t, err)
		})
	}

	// There is no attribute documented for float numbers: https://docs.snowflake.com/en/sql-reference/data-types-numeric#float-float4-float8.
	// However, adding it succeeds for FLOAT, FLOAT4, FLOAT8, and REAL, but ift fails both for DOUBLE and DOUBLE PRECISION.
	for _, c := range incorrectlyCorrectFloatDatatypes {
		t.Run(fmt.Sprintf("document incorrect behavior of float datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1.1::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			require.NoError(t, err)
		})
	}

	// Testing on table creation here because casting (::GEOGRAPHY) was ending with errors (even for the "correct" cases).
	for _, c := range datatypes.GeographyDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of geography datatype: %s", c), func(t *testing.T) {
			tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			sql := fmt.Sprintf("CREATE TABLE %s (i %s)", tableId.FullyQualifiedName(), c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
			t.Cleanup(testClientHelper().Table.DropFunc(t, tableId))

			tableId = testClientHelper().Ids.RandomSchemaObjectIdentifier()
			sql = fmt.Sprintf("CREATE TABLE %s (i %s())", tableId.FullyQualifiedName(), c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.ErrorContains(t, err, "SQL compilation error")
			assert.ErrorContains(t, err, "unexpected '('")
			t.Cleanup(testClientHelper().Table.DropFunc(t, tableId))
		})
	}

	// Testing on table creation here because casting (::GEOMETRY) was ending with errors (even for the "correct" cases).
	for _, c := range datatypes.GeometryDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of geometry datatype: %s", c), func(t *testing.T) {
			tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			sql := fmt.Sprintf("CREATE TABLE %s (i %s)", tableId.FullyQualifiedName(), c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
			t.Cleanup(testClientHelper().Table.DropFunc(t, tableId))

			tableId = testClientHelper().Ids.RandomSchemaObjectIdentifier()
			sql = fmt.Sprintf("CREATE TABLE %s (i %s())", tableId.FullyQualifiedName(), c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.ErrorContains(t, err, "SQL compilation error")
			assert.ErrorContains(t, err, "unexpected '('")
			t.Cleanup(testClientHelper().Table.DropFunc(t, tableId))
		})
	}

	for _, c := range datatypes.NumberDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of number datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT 1::%s(36)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT 1::%s(36, 5)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
		})
	}

	for _, c := range datatypes.NumberDataTypeSubTypes {
		t.Run(fmt.Sprintf("check behavior of number data type subtype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT 1::%s(36)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.ErrorContains(t, err, "SQL compilation error")
			assert.ErrorContains(t, err, "unexpected '36'")
		})
	}

	for _, c := range incorrectNumberDatatypes {
		t.Run(fmt.Sprintf("check behavior of number datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			require.Error(t, err)
		})
	}

	for _, c := range datatypes.ObjectDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of object data type: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT {}::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT {}::%s(36)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.ErrorContains(t, err, "SQL compilation error")
			assert.ErrorContains(t, err, "unexpected '36'")
		})
	}

	for _, c := range datatypes.AllTextDataTypes {
		t.Run(fmt.Sprintf("check behavior of text data type: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 'A'::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT 'ABC'::%s(36)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
		})
	}

	for _, c := range incorrectTextDatatypes {
		t.Run(fmt.Sprintf("check behavior of text datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT ABC::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			require.Error(t, err)
		})
	}

	for _, c := range datatypes.TimeDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of time data type: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT '00:00:00'::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT '00:00:00'::%s(5)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
		})
	}

	for _, c := range datatypes.TimestampLtzDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of timestamp ltz data types: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT '2024-12-02 00:00:00 +0000'::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT '2024-12-02 00:00:00 +0000'::%s(3)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
		})
	}

	for _, c := range datatypes.TimestampNtzDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of timestamp ntz data types: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT '2024-12-02 00:00:00 +0000'::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT '2024-12-02 00:00:00 +0000'::%s(3)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
		})
	}

	for _, c := range datatypes.TimestampTzDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of timestamp tz data types: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT '2024-12-02 00:00:00 +0000'::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT '2024-12-02 00:00:00 +0000'::%s(3)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
		})
	}

	for _, c := range datatypes.VariantDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of variant data type: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT TO_VARIANT(1)::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT TO_VARIANT(1)::%s(36)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.ErrorContains(t, err, "SQL compilation error")
			assert.ErrorContains(t, err, "unexpected '36'")
		})
	}

	// Testing on table creation here because apparently VECTOR is not supported as query in the gosnowflake driver.
	// It ends with "unsupported data type" from https://github.com/snowflakedb/gosnowflake/blob/171ddf2540f3a24f2a990e8453dc425ea864a4a0/converter.go#L1599.
	for _, c := range datatypes.VectorDataTypeSynonyms {
		for _, inner := range datatypes.VectorAllowedInnerTypes {
			t.Run(fmt.Sprintf("check behavior of vector data type: %s, %s", c, inner), func(t *testing.T) {
				tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
				sql := fmt.Sprintf("CREATE TABLE %s (i %s(%s, 2))", tableId.FullyQualifiedName(), c, inner)
				_, err := client.QueryUnsafe(ctx, sql)
				assert.NoError(t, err)
				t.Cleanup(testClientHelper().Table.DropFunc(t, tableId))

				tableId = testClientHelper().Ids.RandomSchemaObjectIdentifier()
				sql = fmt.Sprintf("CREATE TABLE %s (i %s(%s))", tableId.FullyQualifiedName(), c, inner)
				_, err = client.QueryUnsafe(ctx, sql)
				assert.ErrorContains(t, err, "SQL compilation error")
				assert.ErrorContains(t, err, "unexpected ')'")
				t.Cleanup(testClientHelper().Table.DropFunc(t, tableId))
			})
		}
	}

	// Testing on table creation here because apparently VECTOR is not supported as query in the gosnowflake driver.
	// It ends with "unsupported data type" from https://github.com/snowflakedb/gosnowflake/blob/171ddf2540f3a24f2a990e8453dc425ea864a4a0/converter.go#L1599.
	for _, c := range vectorInnerTypesSynonyms {
		t.Run(fmt.Sprintf("document behavior of vector data type synonyms: %s", c), func(t *testing.T) {
			tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			sql := fmt.Sprintf("CREATE TABLE %s (i VECTOR(%s, 3))", tableId.FullyQualifiedName(), c)
			_, err := client.QueryUnsafe(ctx, sql)
			if slices.Contains(vectorInnerTypeSynonymsThatWork, c) {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "SQL compilation error")
				switch {
				case slices.Contains(datatypes.NumberDataTypeSynonyms, c):
					assert.ErrorContains(t, err, fmt.Sprintf("unexpected '%s'", c))
				case slices.Contains(datatypes.NumberDataTypeSubTypes, c):
					assert.ErrorContains(t, err, "Unsupported vector element type 'NUMBER(38,0)'")
				case slices.Contains(datatypes.FloatDataTypeSynonyms, c):
					assert.ErrorContains(t, err, "Unsupported vector element type 'FLOAT'")
				default:
					t.Fail()
				}
			}
			t.Cleanup(testClientHelper().Table.DropFunc(t, tableId))
		})
	}
}

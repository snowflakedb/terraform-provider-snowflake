//go:build non_account_level_tests

package resources_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagIdentifierAndObjectIdentifier(t *testing.T) {
	tagId := sdk.NewSchemaObjectIdentifier("test_db", "test_schema", "test_tag")
	t.Run("account identifier", func(t *testing.T) {
		in := map[string]any{
			"tag_id":      tagId.FullyQualifiedName(),
			"object_type": "ACCOUNT",
			"object_identifiers": []any{
				"orgname.accountname",
			},
		}
		d := schema.TestResourceDataRaw(t, resources.TagAssociation().Schema, in)
		tid, identifiers, objectType, err := resources.TagIdentifierAndObjectIdentifier(d)
		require.NoError(t, err)
		assert.Equal(t, tagId, tid)
		assert.Len(t, identifiers, 1)
		assert.Equal(t, `"orgname"."accountname"`, identifiers[0].FullyQualifiedName())
		assert.Equal(t, sdk.ObjectTypeAccount, objectType)
	})
	t.Run("account object identifier", func(t *testing.T) {
		in := map[string]any{
			"tag_id":      tagId.FullyQualifiedName(),
			"object_type": "DATABASE",
			"object_identifiers": []any{
				"test_db",
			},
		}
		d := schema.TestResourceDataRaw(t, resources.TagAssociation().Schema, in)
		tid, identifiers, objectType, err := resources.TagIdentifierAndObjectIdentifier(d)
		require.NoError(t, err)
		assert.Equal(t, tagId, tid)
		assert.Len(t, identifiers, 1)
		assert.Equal(t, "\"test_db\"", identifiers[0].FullyQualifiedName())
		assert.Equal(t, sdk.ObjectTypeDatabase, objectType)
	})

	t.Run("database object identifier", func(t *testing.T) {
		in := map[string]any{
			"tag_id":      tagId.FullyQualifiedName(),
			"object_type": "SCHEMA",
			"object_identifiers": []any{
				"test_db.test_schema",
			},
		}
		d := schema.TestResourceDataRaw(t, resources.TagAssociation().Schema, in)
		tid, identifiers, objectType, err := resources.TagIdentifierAndObjectIdentifier(d)
		require.NoError(t, err)
		assert.Equal(t, tagId, tid)
		assert.Len(t, identifiers, 1)
		assert.Equal(t, "\"test_db\".\"test_schema\"", identifiers[0].FullyQualifiedName())
		assert.Equal(t, sdk.ObjectTypeSchema, objectType)
	})

	t.Run("schema object identifier", func(t *testing.T) {
		in := map[string]any{
			"tag_id":      tagId.FullyQualifiedName(),
			"object_type": "TABLE",
			"object_identifiers": []any{
				"test_db.test_schema.test_table",
			},
		}
		d := schema.TestResourceDataRaw(t, resources.TagAssociation().Schema, in)
		tid, identifiers, objectType, err := resources.TagIdentifierAndObjectIdentifier(d)
		require.NoError(t, err)
		assert.Equal(t, tagId, tid)
		assert.Len(t, identifiers, 1)
		assert.Equal(t, "\"test_db\".\"test_schema\".\"test_table\"", identifiers[0].FullyQualifiedName())
		assert.Equal(t, sdk.ObjectTypeTable, objectType)
	})

	t.Run("column object identifier", func(t *testing.T) {
		in := map[string]any{
			"tag_id":      "\"test_db\".\"test_schema\".\"test_tag\"",
			"object_type": "COLUMN",
			"object_identifiers": []any{
				"test_db.test_schema.test_table.test_column",
			},
		}
		d := schema.TestResourceDataRaw(t, resources.TagAssociation().Schema, in)
		tid, identifiers, objectType, err := resources.TagIdentifierAndObjectIdentifier(d)
		require.NoError(t, err)
		assert.Equal(t, sdk.NewSchemaObjectIdentifier("test_db", "test_schema", "test_tag"), tid)
		assert.Len(t, identifiers, 1)
		assert.Equal(t, "\"test_db\".\"test_schema\".\"test_table\".\"test_column\"", identifiers[0].FullyQualifiedName())
		assert.Equal(t, sdk.ObjectTypeColumn, objectType)
	})

	t.Run("invalid object identifier", func(t *testing.T) {
		in := map[string]any{
			"tag_id":      tagId.FullyQualifiedName(),
			"object_type": "COLUMN",
			"object_identifiers": []any{
				"\"",
			},
		}
		d := schema.TestResourceDataRaw(t, resources.TagAssociation().Schema, in)
		_, _, _, err := resources.TagIdentifierAndObjectIdentifier(d)
		require.ErrorContains(t, err, `unable to read identifier: ", err = parse error on line 1, column 2: extraneous or missing " in quoted-field`)
	})

	t.Run("invalid tag identifier", func(t *testing.T) {
		in := map[string]any{
			"tag_id":      "\"test_schema\".\"test_tag\"",
			"object_type": "DATABASE",
			"object_identifiers": []any{
				"test_db",
			},
		}
		d := schema.TestResourceDataRaw(t, resources.TagAssociation().Schema, in)
		_, _, _, err := resources.TagIdentifierAndObjectIdentifier(d)
		require.ErrorContains(t, err, `unexpected number of parts 2 in identifier "test_schema"."test_tag", expected 3 in a form of "<database_name>.<schema_name>.<schema_object_name>"`)
	})
}

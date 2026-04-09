package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagReferencesGetForEntity(t *testing.T) {
	t.Run("validation: missing parameters", func(t *testing.T) {
		opts := &getForEntityTagReferenceOptions{}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("getForEntityTagReferenceOptions", "parameters"))
	})

	t.Run("validation: missing arguments", func(t *testing.T) {
		opts := &getForEntityTagReferenceOptions{
			parameters: &tagReferenceParameters{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("tagReferenceParameters", "arguments"))
	})

	t.Run("validation: missing objectName", func(t *testing.T) {
		opts := &getForEntityTagReferenceOptions{
			parameters: &tagReferenceParameters{
				arguments: &tagReferenceFunctionArguments{
					objectDomain: Pointer(TagReferenceObjectDomainTable),
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("tagReferenceFunctionArguments", "objectName"))
	})

	t.Run("validation: missing objectDomain", func(t *testing.T) {
		opts := &getForEntityTagReferenceOptions{
			parameters: &tagReferenceParameters{
				arguments: &tagReferenceFunctionArguments{
					objectName: Pointer("some_name"),
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("tagReferenceFunctionArguments", "objectDomain"))
	})

	t.Run("table domain", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &getForEntityTagReferenceOptions{
			parameters: &tagReferenceParameters{
				arguments: &tagReferenceFunctionArguments{
					objectName:   Pointer(id.FullyQualifiedName()),
					objectDomain: Pointer(TagReferenceObjectDomainTable),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.TAG_REFERENCES ('%s', 'TABLE'))`, temporaryReplace(id))
	})

	t.Run("warehouse domain", func(t *testing.T) {
		opts := &getForEntityTagReferenceOptions{
			parameters: &tagReferenceParameters{
				arguments: &tagReferenceFunctionArguments{
					objectName:   Pointer(NewAccountObjectIdentifier("my_warehouse").FullyQualifiedName()),
					objectDomain: Pointer(TagReferenceObjectDomainWarehouse),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.TAG_REFERENCES ('\"my_warehouse\"', 'WAREHOUSE'))`)
	})

	t.Run("via request builder", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		request := NewGetForEntityTagReferenceRequest(id, TagReferenceObjectDomainTable)
		opts := request.toOpts()
		assertOptsValidAndSQLEquals(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.TAG_REFERENCES ('%s', 'TABLE'))`, temporaryReplace(id))
	})
}

func TestToTagReferenceObjectDomain(t *testing.T) {
	t.Run("valid values", func(t *testing.T) {
		for _, d := range AllTagReferenceObjectDomains {
			result, err := ToTagReferenceObjectDomain(string(d))
			require.NoError(t, err)
			assert.Equal(t, d, result)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		result, err := ToTagReferenceObjectDomain("table")
		require.NoError(t, err)
		assert.Equal(t, TagReferenceObjectDomainTable, result)
	})

	t.Run("invalid value", func(t *testing.T) {
		_, err := ToTagReferenceObjectDomain("INVALID")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid TagReferenceObjectDomain")
	})
}

func TestToTagReferenceApplyMethod(t *testing.T) {
	t.Run("valid values", func(t *testing.T) {
		for _, m := range AllTagReferenceApplyMethods {
			result, err := ToTagReferenceApplyMethod(string(m))
			require.NoError(t, err)
			assert.Equal(t, m, result)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		result, err := ToTagReferenceApplyMethod("manual")
		require.NoError(t, err)
		assert.Equal(t, TagReferenceApplyMethodManual, result)
	})

	t.Run("invalid value", func(t *testing.T) {
		_, err := ToTagReferenceApplyMethod("INVALID")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid TagReferenceApplyMethod")
	})
}

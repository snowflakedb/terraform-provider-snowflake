package sdk

import (
	"testing"
)

func TestTagReferencesGetForEntity(t *testing.T) {
	t.Run("validation: missing parameters", func(t *testing.T) {
		opts := &GetForEntityTagReferenceOptions{}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GetForEntityTagReferenceOptions", "parameters"))
	})

	t.Run("validation: missing arguments", func(t *testing.T) {
		opts := &GetForEntityTagReferenceOptions{
			parameters: &tagReferenceParameters{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GetForEntityTagReferenceOptions.parameters", "arguments"))
	})

	t.Run("validation: missing objectName", func(t *testing.T) {
		opts := &GetForEntityTagReferenceOptions{
			parameters: &tagReferenceParameters{
				arguments: &tagReferenceFunctionArguments{
					ObjectDomain: Pointer(TagReferenceObjectDomainTable),
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GetForEntityTagReferenceOptions.parameters.arguments", "ObjectName"))
	})

	t.Run("validation: missing objectDomain", func(t *testing.T) {
		opts := &GetForEntityTagReferenceOptions{
			parameters: &tagReferenceParameters{
				arguments: &tagReferenceFunctionArguments{
					ObjectName: Pointer("some_name"),
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GetForEntityTagReferenceOptions.parameters.arguments", "ObjectDomain"))
	})

	t.Run("table domain", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &GetForEntityTagReferenceOptions{
			parameters: &tagReferenceParameters{
				arguments: &tagReferenceFunctionArguments{
					ObjectName:   Pointer(id.FullyQualifiedName()),
					ObjectDomain: Pointer(TagReferenceObjectDomainTable),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.TAG_REFERENCES ('%s', 'TABLE'))`, temporaryReplace(id))
	})

	t.Run("warehouse domain", func(t *testing.T) {
		opts := &GetForEntityTagReferenceOptions{
			parameters: &tagReferenceParameters{
				arguments: &tagReferenceFunctionArguments{
					ObjectName:   Pointer(NewAccountObjectIdentifier("my_warehouse").FullyQualifiedName()),
					ObjectDomain: Pointer(TagReferenceObjectDomainWarehouse),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.TAG_REFERENCES ('\"my_warehouse\"', 'WAREHOUSE'))`)
	})
}

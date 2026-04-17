package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagReferencesGetForEntity(t *testing.T) {
	t.Run("validation: missing parameters", func(t *testing.T) {
		opts := &GetForEntityTagReferenceOptions{}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("getForEntityTagReferenceOptions", "parameters"))
	})

	t.Run("validation: missing arguments", func(t *testing.T) {
		opts := &GetForEntityTagReferenceOptions{
			parameters: &tagReferenceParameters{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("tagReferenceParameters", "arguments"))
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

	t.Run("via request builder", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		request := NewGetForEntityTagReferenceRequest(
			NewtagReferenceParametersRequest(
				NewtagReferenceFunctionArgumentsRequest(
					Pointer(id.FullyQualifiedName()),
					Pointer(TagReferenceObjectDomainTable),
				),
			),
		)
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

	t.Run("all declared constants are covered", func(t *testing.T) {
		assert.ElementsMatch(t, []TagReferenceObjectDomain{
			TagReferenceObjectDomainAccount,
			TagReferenceObjectDomainAlert,
			TagReferenceObjectDomainColumn,
			TagReferenceObjectDomainComputePool,
			TagReferenceObjectDomainDatabase,
			TagReferenceObjectDomainDatabaseRole,
			TagReferenceObjectDomainFailoverGroup,
			TagReferenceObjectDomainFunction,
			TagReferenceObjectDomainIntegration,
			TagReferenceObjectDomainNetworkPolicy,
			TagReferenceObjectDomainProcedure,
			TagReferenceObjectDomainReplicationGroup,
			TagReferenceObjectDomainRole,
			TagReferenceObjectDomainSchema,
			TagReferenceObjectDomainShare,
			TagReferenceObjectDomainStage,
			TagReferenceObjectDomainStream,
			TagReferenceObjectDomainTable,
			TagReferenceObjectDomainTask,
			TagReferenceObjectDomainUser,
			TagReferenceObjectDomainWarehouse,
		}, AllTagReferenceObjectDomains)
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

	t.Run("all declared constants are covered", func(t *testing.T) {
		assert.ElementsMatch(t, []TagReferenceApplyMethod{
			TagReferenceApplyMethodClassified,
			TagReferenceApplyMethodInherited,
			TagReferenceApplyMethodManual,
			TagReferenceApplyMethodPropagated,
		}, AllTagReferenceApplyMethods)
	})

	t.Run("convert returns joined mapping errors", func(t *testing.T) {
		row := tagReferenceDBRow{
			TagDatabase: "db",
			TagSchema:   "schema",
			TagName:     "tag",
			TagValue:    "value",
			Level:       "invalid-level",
			ObjectName:  "obj",
			Domain:      "invalid-domain",
			ApplyMethod: "invalid-apply-method",
		}

		_, err := row.convert()
		require.Error(t, err)
		assert.ErrorContains(t, err, fmt.Sprintf("invalid TagReferenceObjectDomain: %s", "INVALID-LEVEL"))
		assert.ErrorContains(t, err, fmt.Sprintf("invalid TagReferenceObjectDomain: %s", "INVALID-DOMAIN"))
		assert.ErrorContains(t, err, fmt.Sprintf("invalid TagReferenceApplyMethod: %s", "INVALID-APPLY-METHOD"))
	})
}

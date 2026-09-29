package sdk

import (
	"strings"
	"testing"
)

func TestPolicyReferencesGetForEntity(t *testing.T) {
	t.Run("validation: missing parameters", func(t *testing.T) {
		opts := &GetForEntityPolicyReferenceOptions{}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GetForEntityPolicyReferenceOptions", "parameters"))
	})

	t.Run("validation: missing arguments", func(t *testing.T) {
		opts := &GetForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GetForEntityPolicyReferenceOptions.parameters", "arguments"))
	})

	t.Run("validation: missing refEntityName", func(t *testing.T) {
		opts := &GetForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					RefEntityDomain: PolicyEntityDomainUser,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GetForEntityPolicyReferenceOptions.parameters.arguments", "refEntityName"))
	})

	t.Run("validation: missing refEntityDomain", func(t *testing.T) {
		opts := &GetForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName: []ObjectIdentifier{NewAccountObjectIdentifierFromFullyQualifiedName("user_name")},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GetForEntityPolicyReferenceOptions.parameters.arguments", "RefEntityDomain"))
	})

	t.Run("user domain", func(t *testing.T) {
		opts := &GetForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{NewAccountObjectIdentifier("user_name")},
					RefEntityDomain: PolicyEntityDomainUser,
				},
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '\"user_name\"', REF_ENTITY_DOMAIN => 'USER'))`)
	})

	t.Run("table domain", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &GetForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{id},
					RefEntityDomain: PolicyEntityDomainTable,
				},
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '%s', REF_ENTITY_DOMAIN => 'TABLE'))`, temporaryReplace(id))
	})

	t.Run("account domain", func(t *testing.T) {
		opts := &GetForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{NewAccountObjectIdentifier("account_name")},
					RefEntityDomain: PolicyEntityDomainAccount,
				},
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '\"account_name\"', REF_ENTITY_DOMAIN => 'ACCOUNT'))`)
	})

	t.Run("integration domain", func(t *testing.T) {
		opts := &GetForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{NewAccountObjectIdentifier("integration_name")},
					RefEntityDomain: PolicyEntityDomainIntegration,
				},
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '\"integration_name\"', REF_ENTITY_DOMAIN => 'INTEGRATION'))`)
	})

	t.Run("tag domain", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &GetForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{id},
					RefEntityDomain: PolicyEntityDomainTag,
				},
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '%s', REF_ENTITY_DOMAIN => 'TAG'))`, temporaryReplace(id))
	})

	t.Run("view domain", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &GetForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{id},
					RefEntityDomain: PolicyEntityDomainView,
				},
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '%s', REF_ENTITY_DOMAIN => 'VIEW'))`, temporaryReplace(id))
	})
}

// TODO [SNOW-1569516]: make nicer during the identifiers rework follow up
func temporaryReplace(id SchemaObjectIdentifier) string {
	return strings.ReplaceAll(id.FullyQualifiedName(), `"`, `\"`)
}

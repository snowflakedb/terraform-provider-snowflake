package sdk

import "strings"

func init() {
	id := randomSchemaObjectIdentifier()
	userId := NewAccountObjectIdentifier("user_name")
	accountId := NewAccountObjectIdentifier("account_name")
	integrationId := NewAccountObjectIdentifier("integration_name")

	policyReferencesTests.GetForEntity.
		withDefaultOpts(func() *GetForEntityPolicyReferenceOptions {
			return &GetForEntityPolicyReferenceOptions{
				parameters: &policyReferenceParameters{
					arguments: &policyReferenceFunctionArguments{
						refEntityName:   []ObjectIdentifier{id},
						RefEntityDomain: PolicyEntityDomainTable,
					},
				},
			}
		}).
		withModify(case_PolicyReferences_validation_GetForEntity_parameters_arguments_refEntityName_ValidateValueSet, func(opts *GetForEntityPolicyReferenceOptions) {
			opts.parameters = &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					RefEntityDomain: PolicyEntityDomainUser,
				},
			}
		}).
		withModify(case_PolicyReferences_validation_GetForEntity_parameters_arguments_RefEntityDomain_ValidateValueSet, func(opts *GetForEntityPolicyReferenceOptions) {
			opts.parameters = &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName: []ObjectIdentifier{userId},
				},
			}
		}).
		withExpectedSqlf(
			case_PolicyReferences_sql_GetForEntity_basic,
			`SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '%s', REF_ENTITY_DOMAIN => 'TABLE'))`,
			temporaryReplace(id),
		).
		withAdditionalSqlCasef(
			"sql_GetForEntity_userDomain",
			func(opts *GetForEntityPolicyReferenceOptions) {
				opts.parameters.arguments.refEntityName = []ObjectIdentifier{userId}
				opts.parameters.arguments.RefEntityDomain = PolicyEntityDomainUser
			},
			`SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '\"user_name\"', REF_ENTITY_DOMAIN => 'USER'))`,
		).
		withAdditionalSqlCasef(
			"sql_GetForEntity_accountDomain",
			func(opts *GetForEntityPolicyReferenceOptions) {
				opts.parameters.arguments.refEntityName = []ObjectIdentifier{accountId}
				opts.parameters.arguments.RefEntityDomain = PolicyEntityDomainAccount
			},
			`SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '\"account_name\"', REF_ENTITY_DOMAIN => 'ACCOUNT'))`,
		).
		withAdditionalSqlCasef(
			"sql_GetForEntity_integrationDomain",
			func(opts *GetForEntityPolicyReferenceOptions) {
				opts.parameters.arguments.refEntityName = []ObjectIdentifier{integrationId}
				opts.parameters.arguments.RefEntityDomain = PolicyEntityDomainIntegration
			},
			`SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '\"integration_name\"', REF_ENTITY_DOMAIN => 'INTEGRATION'))`,
		).
		withAdditionalSqlCasef(
			"sql_GetForEntity_tagDomain",
			func(opts *GetForEntityPolicyReferenceOptions) {
				opts.parameters.arguments.RefEntityDomain = PolicyEntityDomainTag
			},
			`SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '%s', REF_ENTITY_DOMAIN => 'TAG'))`,
			temporaryReplace(id),
		).
		withAdditionalSqlCasef(
			"sql_GetForEntity_viewDomain",
			func(opts *GetForEntityPolicyReferenceOptions) {
				opts.parameters.arguments.RefEntityDomain = PolicyEntityDomainView
			},
			`SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '%s', REF_ENTITY_DOMAIN => 'VIEW'))`,
			temporaryReplace(id),
		)
}

// TODO [SNOW-1569516]: make nicer during the identifiers rework follow up
func temporaryReplace(id SchemaObjectIdentifier) string {
	return strings.ReplaceAll(id.FullyQualifiedName(), `"`, `\"`)
}

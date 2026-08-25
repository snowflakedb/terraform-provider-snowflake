package sdk

func init() {
	dataMetricFunctionReferencesTests.GetForEntity.
		withDefaultOpts(func() *GetForEntityDataMetricFunctionReferenceOptions {
			return &GetForEntityDataMetricFunctionReferenceOptions{
				parameters: &dataMetricFunctionReferenceParameters{
					arguments: &dataMetricFunctionReferenceFunctionArguments{
						refEntityName:   []ObjectIdentifier{NewSchemaObjectIdentifier("a", "b", "c")},
						RefEntityDomain: Pointer(DataMetricFunctionRefEntityDomainOptionView),
					},
				},
			}
		}).
		withExpectedSql(
			case_DataMetricFunctionReferences_sql_GetForEntity_basic,
			`SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.DATA_METRIC_FUNCTION_REFERENCES (REF_ENTITY_NAME => '\"a\".\"b\".\"c\"', REF_ENTITY_DOMAIN => 'VIEW'))`,
		)
}

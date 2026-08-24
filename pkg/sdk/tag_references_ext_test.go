package sdk

func init() {
	id := randomSchemaObjectIdentifier()
	warehouseId := NewAccountObjectIdentifier("my_warehouse")

	tagReferencesTests.GetForEntity.
		withDefaultOpts(func() *GetForEntityTagReferenceOptions {
			return &GetForEntityTagReferenceOptions{
				parameters: &TagReferenceParameters{
					arguments: &TagReferenceFunctionArguments{
						ObjectName:   id.FullyQualifiedName(),
						ObjectDomain: TagReferenceObjectDomainTable,
					},
				},
			}
		}).
		withModify(case_TagReferences_validation_GetForEntity_parameters_arguments_ObjectName_ValidateValueSet, func(opts *GetForEntityTagReferenceOptions) {
			opts.parameters = &TagReferenceParameters{
				arguments: &TagReferenceFunctionArguments{
					ObjectDomain: TagReferenceObjectDomainTable,
				},
			}
		}).
		withModify(case_TagReferences_validation_GetForEntity_parameters_arguments_ObjectDomain_ValidateValueSet, func(opts *GetForEntityTagReferenceOptions) {
			opts.parameters = &TagReferenceParameters{
				arguments: &TagReferenceFunctionArguments{
					ObjectName: "some_name",
				},
			}
		}).
		withExpectedSqlf(
			case_TagReferences_sql_GetForEntity_basic,
			`SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.TAG_REFERENCES ('%s', 'TABLE'))`,
			temporaryReplace(id),
		).
		withAdditionalSqlCasef(
			"sql_GetForEntity_all",
			func(opts *GetForEntityTagReferenceOptions) {
				opts.parameters.arguments.ObjectName = warehouseId.FullyQualifiedName()
				opts.parameters.arguments.ObjectDomain = TagReferenceObjectDomainWarehouse
			},
			`SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.TAG_REFERENCES ('\"my_warehouse\"', 'WAREHOUSE'))`,
		)
}

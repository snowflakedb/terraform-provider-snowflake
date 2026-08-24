package sdk

func init() {
	id := externalFunctionsTestIdSchemaObjectIdentifier
	idWithArgs := externalFunctionsTestIdSchemaObjectIdentifierWithArguments
	noArgsId := randomSchemaObjectIdentifierWithArguments()
	apiIntegration := NewAccountObjectIdentifier("api_integration")
	as := "https://xyz.execute-api.us-west-2.amazonaws.com/prod/remote_echo"
	requestTranslatorId := randomSchemaObjectIdentifier()
	responseTranslatorId := randomSchemaObjectIdentifier()
	headers := []ExternalFunctionHeader{
		{Name: "header1", Value: "value1"},
		{Name: "header2", Value: "value2"},
	}
	contextHeaders := []ExternalFunctionContextHeader{
		{ContextFunction: "CURRENT_ACCOUNT"},
		{ContextFunction: "CURRENT_USER"},
	}
	schemaId := randomDatabaseObjectIdentifier()

	externalFunctionsTests.Create.
		withDefaultOpts(func() *CreateExternalFunctionOptions {
			return &CreateExternalFunctionOptions{
				name:           id,
				ResultDataType: DataTypeVARCHAR,
				ApiIntegration: &apiIntegration,
				As:             as,
			}
		}).
		withExpectedSqlf(
			case_ExternalFunctions_sql_Create_basic,
			`CREATE EXTERNAL FUNCTION %s () RETURNS VARCHAR API_INTEGRATION = "api_integration" AS '%s'`,
			id.FullyQualifiedName(), as,
		).
		withModifyAndExpectedSqlf(
			case_ExternalFunctions_sql_Create_all,
			func(opts *CreateExternalFunctionOptions) {
				opts.OrReplace = new(true)
				opts.Secure = new(true)
				opts.Arguments = []ExternalFunctionArgument{
					{ArgName: "id", ArgDataType: DataTypeNumber},
					{ArgName: "name", ArgDataType: DataTypeVARCHAR},
				}
				opts.ReturnNullValues = new(ReturnNullValuesNotNull)
				opts.NullInputBehavior = new(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Comment = new("comment")
				opts.Headers = headers
				opts.ContextHeaders = contextHeaders
				opts.MaxBatchRows = new(100)
				opts.Compression = new("GZIP")
				opts.RequestTranslator = &requestTranslatorId
				opts.ResponseTranslator = &responseTranslatorId
			},
			`CREATE OR REPLACE SECURE EXTERNAL FUNCTION %s (id NUMBER, name VARCHAR) RETURNS VARCHAR NOT NULL CALLED ON NULL INPUT IMMUTABLE COMMENT = 'comment' API_INTEGRATION = "api_integration" HEADERS = ('header1' = 'value1', 'header2' = 'value2') CONTEXT_HEADERS = (CURRENT_ACCOUNT, CURRENT_USER) MAX_BATCH_ROWS = 100 COMPRESSION = GZIP REQUEST_TRANSLATOR = %s RESPONSE_TRANSLATOR = %s AS '%s'`,
			id.FullyQualifiedName(), requestTranslatorId.FullyQualifiedName(), responseTranslatorId.FullyQualifiedName(), as,
		)

	externalFunctionsTests.Alter.
		withDefaultOpts(func() *AlterExternalFunctionOptions {
			return &AlterExternalFunctionOptions{
				name:     idWithArgs,
				IfExists: new(true),
			}
		}).
		withModify(case_ExternalFunctions_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet, func(opts *AlterExternalFunctionOptions) {
			opts.Set = &ExternalFunctionSet{}
			opts.Unset = &ExternalFunctionUnset{}
		}).
		withModify(case_ExternalFunctions_validation_Alter_opts_Set_ExactlyOneValueSet_MoreThanOneSet, func(opts *AlterExternalFunctionOptions) {
			opts.Set = &ExternalFunctionSet{
				MaxBatchRows: new(100),
				Headers:      headers,
			}
		}).
		withModifyAndExpectedSqlf(
			case_ExternalFunctions_sql_Alter_Set,
			func(opts *AlterExternalFunctionOptions) {
				opts.Set = &ExternalFunctionSet{ApiIntegration: &apiIntegration}
			},
			`ALTER FUNCTION IF EXISTS %s SET API_INTEGRATION = "api_integration"`,
			idWithArgs.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalFunctions_sql_Alter_Unset,
			func(opts *AlterExternalFunctionOptions) {
				opts.Unset = &ExternalFunctionUnset{
					Comment:            new(true),
					Headers:            new(true),
					ContextHeaders:     new(true),
					MaxBatchRows:       new(true),
					Compression:        new(true),
					Secure:             new(true),
					RequestTranslator:  new(true),
					ResponseTranslator: new(true),
				}
			},
			`ALTER FUNCTION IF EXISTS %s UNSET COMMENT, HEADERS, CONTEXT_HEADERS, MAX_BATCH_ROWS, COMPRESSION, SECURE, REQUEST_TRANSLATOR, RESPONSE_TRANSLATOR`,
			idWithArgs.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_headers",
			func(opts *AlterExternalFunctionOptions) {
				opts.Set = &ExternalFunctionSet{Headers: headers}
			},
			`ALTER FUNCTION IF EXISTS %s SET HEADERS = ('header1' = 'value1', 'header2' = 'value2')`,
			idWithArgs.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_maxBatchRows",
			func(opts *AlterExternalFunctionOptions) {
				opts.Set = &ExternalFunctionSet{MaxBatchRows: new(100)}
			},
			`ALTER FUNCTION IF EXISTS %s SET MAX_BATCH_ROWS = 100`,
			idWithArgs.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_compression",
			func(opts *AlterExternalFunctionOptions) {
				opts.Set = &ExternalFunctionSet{Compression: new("GZIP")}
			},
			`ALTER FUNCTION IF EXISTS %s SET COMPRESSION = GZIP`,
			idWithArgs.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_contextHeaders",
			func(opts *AlterExternalFunctionOptions) {
				opts.Set = &ExternalFunctionSet{ContextHeaders: contextHeaders}
			},
			`ALTER FUNCTION IF EXISTS %s SET CONTEXT_HEADERS = (CURRENT_ACCOUNT, CURRENT_USER)`,
			idWithArgs.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_requestTranslator",
			func(opts *AlterExternalFunctionOptions) {
				opts.Set = &ExternalFunctionSet{RequestTranslator: &requestTranslatorId}
			},
			`ALTER FUNCTION IF EXISTS %s SET REQUEST_TRANSLATOR = %s`,
			idWithArgs.FullyQualifiedName(), requestTranslatorId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_responseTranslator",
			func(opts *AlterExternalFunctionOptions) {
				opts.Set = &ExternalFunctionSet{ResponseTranslator: &responseTranslatorId}
			},
			`ALTER FUNCTION IF EXISTS %s SET RESPONSE_TRANSLATOR = %s`,
			idWithArgs.FullyQualifiedName(), responseTranslatorId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_noArguments",
			func(opts *AlterExternalFunctionOptions) {
				opts.name = noArgsId
				opts.Unset = &ExternalFunctionUnset{
					Comment:            new(true),
					Headers:            new(true),
					ContextHeaders:     new(true),
					MaxBatchRows:       new(true),
					Compression:        new(true),
					Secure:             new(true),
					RequestTranslator:  new(true),
					ResponseTranslator: new(true),
				}
			},
			`ALTER FUNCTION IF EXISTS %s UNSET COMMENT, HEADERS, CONTEXT_HEADERS, MAX_BATCH_ROWS, COMPRESSION, SECURE, REQUEST_TRANSLATOR, RESPONSE_TRANSLATOR`,
			noArgsId.FullyQualifiedName(),
		)

	externalFunctionsTests.Show.
		withExpectedSql(case_ExternalFunctions_sql_Show_basic, `SHOW EXTERNAL FUNCTIONS`).
		withModifyAndExpectedSqlf(
			case_ExternalFunctions_sql_Show_all,
			func(opts *ShowExternalFunctionOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.In = &In{Schema: schemaId}
			},
			`SHOW EXTERNAL FUNCTIONS LIKE 'pattern' IN SCHEMA %s`,
			schemaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalFunctions_sql_Show_Like,
			func(opts *ShowExternalFunctionOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			`SHOW EXTERNAL FUNCTIONS LIKE 'pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_ExternalFunctions_sql_Show_In,
			func(opts *ShowExternalFunctionOptions) { opts.In = &In{Schema: schemaId} },
			`SHOW EXTERNAL FUNCTIONS IN SCHEMA %s`,
			schemaId.FullyQualifiedName(),
		)

	externalFunctionsTests.Describe.
		withExpectedSqlf(
			case_ExternalFunctions_sql_Describe_basic,
			`DESCRIBE FUNCTION %s`,
			idWithArgs.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Describe_noArguments",
			func(opts *DescribeExternalFunctionOptions) { opts.name = noArgsId },
			`DESCRIBE FUNCTION %s`,
			noArgsId.FullyQualifiedName(),
		)
}

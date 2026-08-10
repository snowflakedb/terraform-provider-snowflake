package sdk

func init() {
	newId := randomSchemaObjectIdentifierInSchema(proceduresTestIdSchemaObjectIdentifierWithArguments.SchemaId())
	noArgsId := randomSchemaObjectIdentifierWithArguments()
	secretId := randomSchemaObjectIdentifier()
	secretId2 := randomSchemaObjectIdentifier()
	cteId := NewAccountObjectIdentifier("album_info_1976")

	proceduresTests.CreateForJava.
		withDefaultOpts(func() *CreateForJavaProcedureOptions {
			return &CreateForJavaProcedureOptions{
				name:           proceduresTestIdSchemaObjectIdentifier,
				Returns:        ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataType: dataTypeVarchar_100}},
				Handler:        "TestFunc.echoVarchar",
				Packages:       []ProcedurePackage{{ProcedurePackage: "com.snowflake:snowpark:1.2.0"}},
				RuntimeVersion: "1.8",
			}
		}).
		withModify(case_Procedures_validation_CreateForJava_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForJavaProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat}}
		}).
		withModify(case_Procedures_validation_CreateForJava_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForJavaProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat}, {ArgName: "arg2"}}
		}).
		withModify(case_Procedures_validation_CreateForJava_opts_Returns_ResultDataType_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForJavaProcedureOptions) {
			opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{ResultDataTypeOld: DataTypeFloat, ResultDataType: dataTypeFloat}
		}).
		withModify(case_Procedures_validation_CreateForJava_opts_Returns_Table_Columns_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForJavaProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR, ColumnDataType: dataTypeVarchar_100}}}
		}).
		withModify(case_Procedures_validation_CreateForJava_opts_Returns_Table_Columns_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForJavaProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "col2"}}}
		}).
		withAdditionalValidationCase(
			"validation_CreateForJava_procedureDefinition",
			func(opts *CreateForJavaProcedureOptions) { opts.TargetPath = new("@~/testfunc.jar") },
			NewError("TARGET_PATH must be nil when AS is nil"),
		).
		withExpectedSqlf(
			case_Procedures_sql_CreateForJava_basic,
			`CREATE PROCEDURE %s () RETURNS VARCHAR(100) LANGUAGE JAVA RUNTIME_VERSION = '1.8' PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'TestFunc.echoVarchar'`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_CreateForJava_all,
			func(opts *CreateForJavaProcedureOptions) {
				opts.OrReplace = new(true)
				opts.Secure = new(true)
				opts.Arguments = []ProcedureArgument{
					{ArgName: "id", ArgDataType: dataTypeNumber_36_2},
					{ArgName: "name", ArgDataType: dataTypeVarchar_100, DefaultValue: new("'test'")},
				}
				opts.CopyGrants = new(true)
				opts.Returns = ProcedureReturns{
					Table: &ProcedureReturnsTable{
						Columns: []ProcedureColumn{{ColumnName: "country_code", ColumnDataType: dataTypeVarchar_100}},
					},
				}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Imports = []ProcedureImport{{ProcedureImport: "test_jar.jar"}}
				opts.ExternalAccessIntegrations = []AccountObjectIdentifier{NewAccountObjectIdentifier("ext_integration")}
				opts.Secrets = []SecretReference{{VariableName: "variable1", Name: secretId}, {VariableName: "variable2", Name: secretId2}}
				opts.TargetPath = new("@~/testfunc.jar")
				opts.Comment = new("test comment")
				opts.ExecuteAs = new(ExecuteAsCaller)
				opts.ProcedureDefinition = new("return id + name;")
			},
			`CREATE OR REPLACE SECURE PROCEDURE %s ("id" NUMBER(36, 2), "name" VARCHAR(100) DEFAULT 'test') COPY GRANTS RETURNS TABLE ("country_code" VARCHAR(100)) LANGUAGE JAVA STRICT IMMUTABLE RUNTIME_VERSION = '1.8' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('test_jar.jar') HANDLER = 'TestFunc.echoVarchar' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) TARGET_PATH = '@~/testfunc.jar' COMMENT = 'test comment' EXECUTE AS CALLER AS return id + name;`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		).
		// TODO [SNOW-1348106]: remove with old procedure removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForJava_allOldDataTypes",
			func(opts *CreateForJavaProcedureOptions) {
				opts.OrReplace = new(true)
				opts.Secure = new(true)
				opts.Arguments = []ProcedureArgument{
					{ArgName: "id", ArgDataTypeOld: DataTypeNumber},
					{ArgName: "name", ArgDataTypeOld: DataTypeVARCHAR, DefaultValue: new("'test'")},
				}
				opts.CopyGrants = new(true)
				opts.Returns = ProcedureReturns{
					Table: &ProcedureReturnsTable{
						Columns: []ProcedureColumn{{ColumnName: "country_code", ColumnDataTypeOld: DataTypeVARCHAR}},
					},
				}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Imports = []ProcedureImport{{ProcedureImport: "test_jar.jar"}}
				opts.ExternalAccessIntegrations = []AccountObjectIdentifier{NewAccountObjectIdentifier("ext_integration")}
				opts.Secrets = []SecretReference{{VariableName: "variable1", Name: secretId}, {VariableName: "variable2", Name: secretId2}}
				opts.TargetPath = new("@~/testfunc.jar")
				opts.Comment = new("test comment")
				opts.ExecuteAs = new(ExecuteAsCaller)
				opts.ProcedureDefinition = new("return id + name;")
			},
			`CREATE OR REPLACE SECURE PROCEDURE %s ("id" NUMBER, "name" VARCHAR DEFAULT 'test') COPY GRANTS RETURNS TABLE ("country_code" VARCHAR) LANGUAGE JAVA STRICT IMMUTABLE RUNTIME_VERSION = '1.8' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('test_jar.jar') HANDLER = 'TestFunc.echoVarchar' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) TARGET_PATH = '@~/testfunc.jar' COMMENT = 'test comment' EXECUTE AS CALLER AS return id + name;`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		)

	proceduresTests.CreateForJavaScript.
		withDefaultOpts(func() *CreateForJavaScriptProcedureOptions {
			return &CreateForJavaScriptProcedureOptions{
				name:                proceduresTestIdSchemaObjectIdentifier,
				ResultDataType:      dataTypeFloat,
				ProcedureDefinition: "return 1;",
			}
		}).
		withModify(case_Procedures_validation_CreateForJavaScript_opts_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForJavaScriptProcedureOptions) {
			opts.ResultDataTypeOld = DataTypeFloat
		}).
		withModify(case_Procedures_validation_CreateForJavaScript_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForJavaScriptProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "d", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat}}
		}).
		withModify(case_Procedures_validation_CreateForJavaScript_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForJavaScriptProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "d", ArgDataTypeOld: DataTypeFloat}, {ArgName: "d2"}}
		}).
		withExpectedSqlf(
			case_Procedures_sql_CreateForJavaScript_basic,
			`CREATE PROCEDURE %s () RETURNS FLOAT LANGUAGE JAVASCRIPT AS return 1;`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_CreateForJavaScript_all,
			func(opts *CreateForJavaScriptProcedureOptions) {
				opts.OrReplace = new(true)
				opts.Secure = new(true)
				opts.Arguments = []ProcedureArgument{{ArgName: "d", ArgDataType: dataTypeFloat, DefaultValue: new("1.0")}}
				opts.CopyGrants = new(true)
				opts.NotNull = new(true)
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Comment = new("test comment")
				opts.ExecuteAs = new(ExecuteAsCaller)
				opts.ProcedureDefinition = "return 1;"
			},
			`CREATE OR REPLACE SECURE PROCEDURE %s ("d" FLOAT DEFAULT 1.0) COPY GRANTS RETURNS FLOAT NOT NULL LANGUAGE JAVASCRIPT STRICT IMMUTABLE COMMENT = 'test comment' EXECUTE AS CALLER AS return 1;`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		// TODO [SNOW-1348106]: remove with old procedure removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForJavaScript_allOldDataTypes",
			func(opts *CreateForJavaScriptProcedureOptions) {
				opts.OrReplace = new(true)
				opts.Secure = new(true)
				opts.Arguments = []ProcedureArgument{{ArgName: "d", ArgDataTypeOld: "DOUBLE", DefaultValue: new("1.0")}}
				opts.CopyGrants = new(true)
				opts.ResultDataType = nil
				opts.ResultDataTypeOld = "DOUBLE"
				opts.NotNull = new(true)
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Comment = new("test comment")
				opts.ExecuteAs = new(ExecuteAsCaller)
				opts.ProcedureDefinition = "return 1;"
			},
			`CREATE OR REPLACE SECURE PROCEDURE %s ("d" DOUBLE DEFAULT 1.0) COPY GRANTS RETURNS DOUBLE NOT NULL LANGUAGE JAVASCRIPT STRICT IMMUTABLE COMMENT = 'test comment' EXECUTE AS CALLER AS return 1;`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	proceduresTests.CreateForPython.
		withDefaultOpts(func() *CreateForPythonProcedureOptions {
			return &CreateForPythonProcedureOptions{
				name:           proceduresTestIdSchemaObjectIdentifier,
				Returns:        ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataType: dataTypeVarchar_100}},
				Handler:        "udf",
				Packages:       []ProcedurePackage{{ProcedurePackage: "snowflake-snowpark-python"}},
				RuntimeVersion: "3.9",
			}
		}).
		withModify(case_Procedures_validation_CreateForPython_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForPythonProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "i", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat}}
		}).
		withModify(case_Procedures_validation_CreateForPython_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForPythonProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "i", ArgDataTypeOld: DataTypeFloat}, {ArgName: "i2"}}
		}).
		withModify(case_Procedures_validation_CreateForPython_opts_Returns_ResultDataType_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForPythonProcedureOptions) {
			opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{ResultDataTypeOld: DataTypeFloat, ResultDataType: dataTypeFloat}
		}).
		withModify(case_Procedures_validation_CreateForPython_opts_Returns_Table_Columns_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForPythonProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR, ColumnDataType: dataTypeVarchar_100}}}
		}).
		withModify(case_Procedures_validation_CreateForPython_opts_Returns_Table_Columns_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForPythonProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "col2"}}}
		}).
		withExpectedSqlf(
			case_Procedures_sql_CreateForPython_basic,
			`CREATE PROCEDURE %s () RETURNS VARCHAR(100) LANGUAGE PYTHON RUNTIME_VERSION = '3.9' PACKAGES = ('snowflake-snowpark-python') HANDLER = 'udf'`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_CreateForPython_all,
			func(opts *CreateForPythonProcedureOptions) {
				opts.OrReplace = new(true)
				opts.Secure = new(true)
				opts.Arguments = []ProcedureArgument{{ArgName: "i", ArgDataType: dataTypeNumber_36_2, DefaultValue: new("1")}}
				opts.CopyGrants = new(true)
				opts.Returns = ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataType: dataTypeVariant, Null: new(true)}}
				opts.Packages = []ProcedurePackage{{ProcedurePackage: "numpy"}, {ProcedurePackage: "pandas"}}
				opts.Imports = []ProcedureImport{{ProcedureImport: "numpy"}, {ProcedureImport: "pandas"}}
				opts.ExternalAccessIntegrations = []AccountObjectIdentifier{NewAccountObjectIdentifier("ext_integration")}
				opts.Secrets = []SecretReference{{VariableName: "variable1", Name: secretId}, {VariableName: "variable2", Name: secretId2}}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Comment = new("test comment")
				opts.ExecuteAs = new(ExecuteAsCaller)
				opts.ProcedureDefinition = new("import numpy as np")
			},
			`CREATE OR REPLACE SECURE PROCEDURE %s ("i" NUMBER(36, 2) DEFAULT 1) COPY GRANTS RETURNS VARIANT NULL LANGUAGE PYTHON STRICT IMMUTABLE RUNTIME_VERSION = '3.9' PACKAGES = ('numpy', 'pandas') IMPORTS = ('numpy', 'pandas') HANDLER = 'udf' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) COMMENT = 'test comment' EXECUTE AS CALLER AS import numpy as np`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		).
		// TODO [SNOW-1348106]: remove with old procedure removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForPython_allOldDataTypes",
			func(opts *CreateForPythonProcedureOptions) {
				opts.OrReplace = new(true)
				opts.Secure = new(true)
				opts.Arguments = []ProcedureArgument{{ArgName: "i", ArgDataTypeOld: "int", DefaultValue: new("1")}}
				opts.CopyGrants = new(true)
				opts.Returns = ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataTypeOld: "VARIANT", Null: new(true)}}
				opts.Packages = []ProcedurePackage{{ProcedurePackage: "numpy"}, {ProcedurePackage: "pandas"}}
				opts.Imports = []ProcedureImport{{ProcedureImport: "numpy"}, {ProcedureImport: "pandas"}}
				opts.ExternalAccessIntegrations = []AccountObjectIdentifier{NewAccountObjectIdentifier("ext_integration")}
				opts.Secrets = []SecretReference{{VariableName: "variable1", Name: secretId}, {VariableName: "variable2", Name: secretId2}}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Comment = new("test comment")
				opts.ExecuteAs = new(ExecuteAsCaller)
				opts.ProcedureDefinition = new("import numpy as np")
			},
			`CREATE OR REPLACE SECURE PROCEDURE %s ("i" int DEFAULT 1) COPY GRANTS RETURNS VARIANT NULL LANGUAGE PYTHON STRICT IMMUTABLE RUNTIME_VERSION = '3.9' PACKAGES = ('numpy', 'pandas') IMPORTS = ('numpy', 'pandas') HANDLER = 'udf' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) COMMENT = 'test comment' EXECUTE AS CALLER AS import numpy as np`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		)

	proceduresTests.CreateForScala.
		withDefaultOpts(func() *CreateForScalaProcedureOptions {
			return &CreateForScalaProcedureOptions{
				name:           proceduresTestIdSchemaObjectIdentifier,
				Returns:        ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataType: dataTypeVarchar_100}},
				Handler:        "Echo.echoVarchar",
				Packages:       []ProcedurePackage{{ProcedurePackage: "com.snowflake:snowpark:1.2.0"}},
				RuntimeVersion: "2.0",
			}
		}).
		withModify(case_Procedures_validation_CreateForScala_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForScalaProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat}}
		}).
		withModify(case_Procedures_validation_CreateForScala_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForScalaProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat}, {ArgName: "arg2"}}
		}).
		withModify(case_Procedures_validation_CreateForScala_opts_Returns_ResultDataType_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForScalaProcedureOptions) {
			opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{ResultDataTypeOld: DataTypeFloat, ResultDataType: dataTypeFloat}
		}).
		withModify(case_Procedures_validation_CreateForScala_opts_Returns_Table_Columns_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForScalaProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR, ColumnDataType: dataTypeVarchar_100}}}
		}).
		withModify(case_Procedures_validation_CreateForScala_opts_Returns_Table_Columns_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForScalaProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "col2"}}}
		}).
		withAdditionalValidationCase(
			"validation_CreateForScala_procedureDefinition",
			func(opts *CreateForScalaProcedureOptions) { opts.TargetPath = new("@~/testfunc.jar") },
			NewError("TARGET_PATH must be nil when AS is nil"),
		).
		withExpectedSqlf(
			case_Procedures_sql_CreateForScala_basic,
			`CREATE PROCEDURE %s () RETURNS VARCHAR(100) LANGUAGE SCALA RUNTIME_VERSION = '2.0' PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'Echo.echoVarchar'`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_CreateForScala_all,
			func(opts *CreateForScalaProcedureOptions) {
				opts.OrReplace = new(true)
				opts.Secure = new(true)
				opts.Arguments = []ProcedureArgument{{ArgName: "x", ArgDataType: dataTypeVarchar_100, DefaultValue: new("'test'")}}
				opts.CopyGrants = new(true)
				opts.Returns = ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataType: dataTypeVarchar_100, NotNull: new(true)}}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Imports = []ProcedureImport{{ProcedureImport: "@udf_libs/echohandler.jar"}}
				opts.TargetPath = new("@~/testfunc.jar")
				opts.Comment = new("test comment")
				opts.ExecuteAs = new(ExecuteAsCaller)
				opts.ProcedureDefinition = new("return x")
			},
			`CREATE OR REPLACE SECURE PROCEDURE %s ("x" VARCHAR(100) DEFAULT 'test') COPY GRANTS RETURNS VARCHAR(100) NOT NULL LANGUAGE SCALA STRICT IMMUTABLE RUNTIME_VERSION = '2.0' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar' TARGET_PATH = '@~/testfunc.jar' COMMENT = 'test comment' EXECUTE AS CALLER AS return x`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		// TODO [SNOW-1348106]: remove with old procedure removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForScala_allOldDataTypes",
			func(opts *CreateForScalaProcedureOptions) {
				opts.OrReplace = new(true)
				opts.Secure = new(true)
				opts.Arguments = []ProcedureArgument{{ArgName: "x", ArgDataTypeOld: "VARCHAR", DefaultValue: new("'test'")}}
				opts.CopyGrants = new(true)
				opts.Returns = ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataTypeOld: "VARCHAR", NotNull: new(true)}}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Imports = []ProcedureImport{{ProcedureImport: "@udf_libs/echohandler.jar"}}
				opts.TargetPath = new("@~/testfunc.jar")
				opts.Comment = new("test comment")
				opts.ExecuteAs = new(ExecuteAsCaller)
				opts.ProcedureDefinition = new("return x")
			},
			`CREATE OR REPLACE SECURE PROCEDURE %s ("x" VARCHAR DEFAULT 'test') COPY GRANTS RETURNS VARCHAR NOT NULL LANGUAGE SCALA STRICT IMMUTABLE RUNTIME_VERSION = '2.0' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar' TARGET_PATH = '@~/testfunc.jar' COMMENT = 'test comment' EXECUTE AS CALLER AS return x`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	proceduresTests.CreateForSQL.
		withDefaultOpts(func() *CreateForSQLProcedureOptions {
			return &CreateForSQLProcedureOptions{
				name:                proceduresTestIdSchemaObjectIdentifier,
				Returns:             ProcedureSQLReturns{ResultDataType: &ProcedureSQLReturnsResultDataType{ResultDataType: dataTypeVarchar_100}},
				ProcedureDefinition: "3.141592654::FLOAT",
			}
		}).
		withModify(case_Procedures_validation_CreateForSQL_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForSQLProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat}}
		}).
		withModify(case_Procedures_validation_CreateForSQL_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForSQLProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat}, {ArgName: "arg2"}}
		}).
		withModify(case_Procedures_validation_CreateForSQL_opts_Returns_ResultDataType_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForSQLProcedureOptions) {
			opts.Returns.ResultDataType = &ProcedureSQLReturnsResultDataType{ResultDataTypeOld: DataTypeFloat, ResultDataType: dataTypeFloat}
		}).
		withModify(case_Procedures_validation_CreateForSQL_opts_Returns_Table_Columns_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForSQLProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR, ColumnDataType: dataTypeVarchar_100}}}
		}).
		withModify(case_Procedures_validation_CreateForSQL_opts_Returns_Table_Columns_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForSQLProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "col2"}}}
		}).
		withExpectedSqlf(
			case_Procedures_sql_CreateForSQL_basic,
			`CREATE PROCEDURE %s () RETURNS VARCHAR(100) LANGUAGE SQL AS 3.141592654::FLOAT`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_CreateForSQL_all,
			func(opts *CreateForSQLProcedureOptions) {
				opts.OrReplace = new(true)
				opts.Secure = new(true)
				opts.Arguments = []ProcedureArgument{{ArgName: "message", ArgDataType: dataTypeVarchar_100, DefaultValue: new("'test'")}}
				opts.CopyGrants = new(true)
				opts.Returns = ProcedureSQLReturns{
					ResultDataType: &ProcedureSQLReturnsResultDataType{ResultDataType: dataTypeVarchar_100},
					NotNull:        new(true),
				}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Comment = new("test comment")
				opts.ExecuteAs = new(ExecuteAsCaller)
				opts.ProcedureDefinition = "3.141592654::FLOAT"
			},
			`CREATE OR REPLACE SECURE PROCEDURE %s ("message" VARCHAR(100) DEFAULT 'test') COPY GRANTS RETURNS VARCHAR(100) NOT NULL LANGUAGE SQL STRICT IMMUTABLE COMMENT = 'test comment' EXECUTE AS CALLER AS 3.141592654::FLOAT`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		// TODO [SNOW-1348106]: remove with old procedure removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForSQL_allOldDataTypes",
			func(opts *CreateForSQLProcedureOptions) {
				opts.OrReplace = new(true)
				opts.Secure = new(true)
				opts.Arguments = []ProcedureArgument{{ArgName: "message", ArgDataTypeOld: "VARCHAR", DefaultValue: new("'test'")}}
				opts.CopyGrants = new(true)
				opts.Returns = ProcedureSQLReturns{
					ResultDataType: &ProcedureSQLReturnsResultDataType{ResultDataTypeOld: "VARCHAR"},
					NotNull:        new(true),
				}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Comment = new("test comment")
				opts.ExecuteAs = new(ExecuteAsCaller)
				opts.ProcedureDefinition = "3.141592654::FLOAT"
			},
			`CREATE OR REPLACE SECURE PROCEDURE %s ("message" VARCHAR DEFAULT 'test') COPY GRANTS RETURNS VARCHAR NOT NULL LANGUAGE SQL STRICT IMMUTABLE COMMENT = 'test comment' EXECUTE AS CALLER AS 3.141592654::FLOAT`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateForSQL_createWithNoArguments",
			func(opts *CreateForSQLProcedureOptions) {
				opts.Returns = ProcedureSQLReturns{ResultDataType: &ProcedureSQLReturnsResultDataType{ResultDataType: dataTypeFloat}}
			},
			`CREATE PROCEDURE %s () RETURNS FLOAT LANGUAGE SQL AS 3.141592654::FLOAT`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	proceduresTests.Alter.
		withDefaultOpts(func() *AlterProcedureOptions {
			return &AlterProcedureOptions{
				name:     proceduresTestIdSchemaObjectIdentifierWithArguments,
				IfExists: new(true),
			}
		}).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_Alter_RenameTo,
			func(opts *AlterProcedureOptions) { opts.RenameTo = &newId },
			`ALTER PROCEDURE IF EXISTS %s RENAME TO %s`,
			proceduresTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(), newId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_Alter_ExecuteAs,
			func(opts *AlterProcedureOptions) { opts.ExecuteAs = new(ExecuteAsCaller) },
			`ALTER PROCEDURE IF EXISTS %s EXECUTE AS CALLER`,
			proceduresTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_Alter_Set,
			func(opts *AlterProcedureOptions) {
				opts.Set = &ProcedureSet{Comment: new("comment"), TraceLevel: new(TraceLevelOff)}
			},
			`ALTER PROCEDURE IF EXISTS %s SET COMMENT = 'comment', TRACE_LEVEL = 'OFF'`,
			proceduresTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_Alter_Unset,
			func(opts *AlterProcedureOptions) {
				opts.Unset = &ProcedureUnset{Comment: new(true), TraceLevel: new(true)}
			},
			`ALTER PROCEDURE IF EXISTS %s UNSET COMMENT, TRACE_LEVEL`,
			proceduresTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_Alter_SetTags,
			func(opts *AlterProcedureOptions) {
				opts.SetTags = []TagAssociation{{Name: NewAccountObjectIdentifier("tag1"), Value: "value1"}}
			},
			`ALTER PROCEDURE IF EXISTS %s SET TAG "tag1" = 'value1'`,
			proceduresTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_Alter_UnsetTags,
			func(opts *AlterProcedureOptions) {
				opts.UnsetTags = []ObjectIdentifier{NewAccountObjectIdentifier("tag1"), NewAccountObjectIdentifier("tag2")}
			},
			`ALTER PROCEDURE IF EXISTS %s UNSET TAG "tag1", "tag2"`,
			proceduresTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_emptySecrets",
			func(opts *AlterProcedureOptions) { opts.Set = &ProcedureSet{SecretsList: &SecretsList{}} },
			`ALTER PROCEDURE IF EXISTS %s SET SECRETS = ()`,
			proceduresTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_nonEmptySecrets",
			func(opts *AlterProcedureOptions) {
				opts.Set = &ProcedureSet{SecretsList: &SecretsList{SecretsList: []SecretReference{{VariableName: "abc", Name: secretId}}}}
			},
			`ALTER PROCEDURE IF EXISTS %s SET SECRETS = ('abc' = %s)`,
			proceduresTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(), secretId.FullyQualifiedName(),
		)

	proceduresTests.Drop.
		withExpectedSqlf(
			case_Procedures_sql_Drop_basic,
			`DROP PROCEDURE %s`,
			proceduresTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_Drop_all,
			func(opts *DropProcedureOptions) { opts.IfExists = new(true) },
			`DROP PROCEDURE IF EXISTS %s`,
			proceduresTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Drop_noArguments",
			func(opts *DropProcedureOptions) { opts.name = noArgsId },
			`DROP PROCEDURE %s`,
			noArgsId.FullyQualifiedName(),
		)

	proceduresTests.Show.
		withExpectedSql(case_Procedures_sql_Show_basic, "SHOW PROCEDURES").
		withModifyAndExpectedSqlf(
			case_Procedures_sql_Show_all,
			func(opts *ShowProcedureOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
			},
			"SHOW PROCEDURES LIKE 'pattern' IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_Show_Like,
			func(opts *ShowProcedureOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			"SHOW PROCEDURES LIKE 'pattern'",
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_Show_In,
			func(opts *ShowProcedureOptions) { opts.In = &ExtendedIn{In: In{Account: new(true)}} },
			"SHOW PROCEDURES IN ACCOUNT",
		)

	proceduresTests.Describe.
		withExpectedSqlf(
			case_Procedures_sql_Describe_basic,
			`DESCRIBE PROCEDURE %s`,
			proceduresTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		)

	proceduresTests.Call.
		withExpectedSqlf(
			case_Procedures_sql_Call_basic,
			`CALL %s ()`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Call_allOptions_namedArgs",
			func(opts *CallProcedureOptions) {
				opts.CallArguments = []string{"province => 'Manitoba'", "amount => 127.4"}
				opts.ScriptingVariable = new(":ret")
			},
			`CALL %s (province => 'Manitoba', amount => 127.4) INTO :ret`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Call_allOptions_positionalArgs",
			func(opts *CallProcedureOptions) {
				opts.CallArguments = []string{"'Manitoba'", "127.4"}
				opts.ScriptingVariable = new(":ret")
			},
			`CALL %s ('Manitoba', 127.4) INTO :ret`,
			proceduresTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	proceduresTests.CreateAndCallForJava.
		withDefaultOpts(func() *CreateAndCallForJavaProcedureOptions {
			return &CreateAndCallForJavaProcedureOptions{
				Name:           proceduresTestIdAccountObjectIdentifier,
				ProcedureName:  proceduresTestIdAccountObjectIdentifier,
				Returns:        ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataType: dataTypeVarchar_100}},
				Handler:        "TestFunc.echoVarchar",
				Packages:       []ProcedurePackage{{ProcedurePackage: "com.snowflake:snowpark:1.2.0"}},
				RuntimeVersion: "1.8",
			}
		}).
		withModify(case_Procedures_validation_CreateAndCallForJava_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForJavaProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat}}
		}).
		withModify(case_Procedures_validation_CreateAndCallForJava_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateAndCallForJavaProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat}, {ArgName: "arg2"}}
		}).
		withModify(case_Procedures_validation_CreateAndCallForJava_opts_Returns_ResultDataType_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForJavaProcedureOptions) {
			opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{ResultDataTypeOld: DataTypeFloat, ResultDataType: dataTypeFloat}
		}).
		withModify(case_Procedures_validation_CreateAndCallForJava_opts_Returns_Table_Columns_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForJavaProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR, ColumnDataType: dataTypeVarchar_100}}}
		}).
		withModify(case_Procedures_validation_CreateAndCallForJava_opts_Returns_Table_Columns_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateAndCallForJavaProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "col2"}}}
		}).
		withExpectedSqlf(
			case_Procedures_sql_CreateAndCallForJava_basic,
			`WITH %s AS PROCEDURE () RETURNS VARCHAR(100) LANGUAGE JAVA RUNTIME_VERSION = '1.8' PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'TestFunc.echoVarchar' CALL %s ()`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_CreateAndCallForJava_all,
			func(opts *CreateAndCallForJavaProcedureOptions) {
				opts.Arguments = []ProcedureArgument{{ArgName: "id", ArgDataType: dataTypeNumber_36_2}, {ArgName: "name", ArgDataType: dataTypeVarchar_100}}
				opts.Returns = ProcedureReturns{Table: &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "country_code", ColumnDataType: dataTypeVarchar_100}}}}
				opts.Imports = []ProcedureImport{{ProcedureImport: "test_jar.jar"}}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ProcedureDefinition = new("return id + name;")
				opts.WithClause = &ProcedureWithClause{CteName: cteId, CteColumns: []string{"x", "y"}, Statement: "(select m.album_ID, m.album_name, b.band_name from music_albums)"}
				opts.ScriptingVariable = new(":ret")
				opts.CallArguments = []string{"1", "rnd"}
			},
			`WITH %s AS PROCEDURE ("id" NUMBER(36, 2), "name" VARCHAR(100)) RETURNS TABLE ("country_code" VARCHAR(100)) LANGUAGE JAVA STRICT RUNTIME_VERSION = '1.8' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('test_jar.jar') HANDLER = 'TestFunc.echoVarchar' AS 'return id + name;' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1, rnd) INTO :ret`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), cteId.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateAndCallForJava_noArguments",
			func(opts *CreateAndCallForJavaProcedureOptions) {
				opts.Returns = ProcedureReturns{Table: &ProcedureReturnsTable{}}
				opts.Packages = []ProcedurePackage{{ProcedurePackage: "com.snowflake:snowpark:latest"}}
			},
			`WITH %s AS PROCEDURE () RETURNS TABLE () LANGUAGE JAVA RUNTIME_VERSION = '1.8' PACKAGES = ('com.snowflake:snowpark:latest') HANDLER = 'TestFunc.echoVarchar' CALL %s ()`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		// TODO [SNOW-1348106]: remove with old procedure removal for V1
		withAdditionalSqlCasef(
			"sql_CreateAndCallForJava_allOldDataTypes",
			func(opts *CreateAndCallForJavaProcedureOptions) {
				opts.Arguments = []ProcedureArgument{{ArgName: "id", ArgDataTypeOld: DataTypeNumber}, {ArgName: "name", ArgDataTypeOld: DataTypeVARCHAR}}
				opts.Returns = ProcedureReturns{Table: &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "country_code", ColumnDataTypeOld: DataTypeVARCHAR}}}}
				opts.Imports = []ProcedureImport{{ProcedureImport: "test_jar.jar"}}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ProcedureDefinition = new("return id + name;")
				opts.WithClause = &ProcedureWithClause{CteName: cteId, CteColumns: []string{"x", "y"}, Statement: "(select m.album_ID, m.album_name, b.band_name from music_albums)"}
				opts.ScriptingVariable = new(":ret")
				opts.CallArguments = []string{"1", "rnd"}
			},
			`WITH %s AS PROCEDURE ("id" NUMBER, "name" VARCHAR) RETURNS TABLE ("country_code" VARCHAR) LANGUAGE JAVA STRICT RUNTIME_VERSION = '1.8' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('test_jar.jar') HANDLER = 'TestFunc.echoVarchar' AS 'return id + name;' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1, rnd) INTO :ret`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), cteId.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		)

	proceduresTests.CreateAndCallForScala.
		withDefaultOpts(func() *CreateAndCallForScalaProcedureOptions {
			return &CreateAndCallForScalaProcedureOptions{
				Name:           proceduresTestIdAccountObjectIdentifier,
				ProcedureName:  proceduresTestIdAccountObjectIdentifier,
				Returns:        ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataType: dataTypeVarchar_100}},
				Handler:        "TestFunc.echoVarchar",
				Packages:       []ProcedurePackage{{ProcedurePackage: "com.snowflake:snowpark:1.2.0"}},
				RuntimeVersion: "2.12",
			}
		}).
		withModify(case_Procedures_validation_CreateAndCallForScala_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForScalaProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat}}
		}).
		withModify(case_Procedures_validation_CreateAndCallForScala_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateAndCallForScalaProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat}, {ArgName: "arg2"}}
		}).
		withModify(case_Procedures_validation_CreateAndCallForScala_opts_Returns_ResultDataType_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForScalaProcedureOptions) {
			opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{ResultDataTypeOld: DataTypeFloat, ResultDataType: dataTypeFloat}
		}).
		withModify(case_Procedures_validation_CreateAndCallForScala_opts_Returns_Table_Columns_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForScalaProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR, ColumnDataType: dataTypeVarchar_100}}}
		}).
		withModify(case_Procedures_validation_CreateAndCallForScala_opts_Returns_Table_Columns_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateAndCallForScalaProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "col2"}}}
		}).
		withExpectedSqlf(
			case_Procedures_sql_CreateAndCallForScala_basic,
			`WITH %s AS PROCEDURE () RETURNS VARCHAR(100) LANGUAGE SCALA RUNTIME_VERSION = '2.12' PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'TestFunc.echoVarchar' CALL %s ()`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_CreateAndCallForScala_all,
			func(opts *CreateAndCallForScalaProcedureOptions) {
				opts.Arguments = []ProcedureArgument{{ArgName: "id", ArgDataType: dataTypeNumber_36_2}, {ArgName: "name", ArgDataType: dataTypeVarchar_100}}
				opts.Returns = ProcedureReturns{Table: &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "country_code", ColumnDataType: dataTypeVarchar_100}}}}
				opts.Imports = []ProcedureImport{{ProcedureImport: "test_jar.jar"}}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ProcedureDefinition = new("return id + name;")
				opts.WithClauses = []ProcedureWithClause{{CteName: cteId, CteColumns: []string{"x", "y"}, Statement: "(select m.album_ID, m.album_name, b.band_name from music_albums)"}}
				opts.ScriptingVariable = new(":ret")
				opts.CallArguments = []string{"1", "rnd"}
			},
			`WITH %s AS PROCEDURE ("id" NUMBER(36, 2), "name" VARCHAR(100)) RETURNS TABLE ("country_code" VARCHAR(100)) LANGUAGE SCALA STRICT RUNTIME_VERSION = '2.12' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('test_jar.jar') HANDLER = 'TestFunc.echoVarchar' AS 'return id + name;' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1, rnd) INTO :ret`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), cteId.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateAndCallForScala_noArguments",
			func(opts *CreateAndCallForScalaProcedureOptions) {
				opts.Returns = ProcedureReturns{Table: &ProcedureReturnsTable{}}
			},
			`WITH %s AS PROCEDURE () RETURNS TABLE () LANGUAGE SCALA RUNTIME_VERSION = '2.12' PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'TestFunc.echoVarchar' CALL %s ()`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		// TODO [SNOW-1348106]: remove with old procedure removal for V1
		withAdditionalSqlCasef(
			"sql_CreateAndCallForScala_allOldDataTypes",
			func(opts *CreateAndCallForScalaProcedureOptions) {
				opts.Arguments = []ProcedureArgument{{ArgName: "id", ArgDataTypeOld: DataTypeNumber}, {ArgName: "name", ArgDataTypeOld: DataTypeVARCHAR}}
				opts.Returns = ProcedureReturns{Table: &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "country_code", ColumnDataTypeOld: DataTypeVARCHAR}}}}
				opts.Imports = []ProcedureImport{{ProcedureImport: "test_jar.jar"}}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ProcedureDefinition = new("return id + name;")
				opts.WithClauses = []ProcedureWithClause{{CteName: cteId, CteColumns: []string{"x", "y"}, Statement: "(select m.album_ID, m.album_name, b.band_name from music_albums)"}}
				opts.ScriptingVariable = new(":ret")
				opts.CallArguments = []string{"1", "rnd"}
			},
			`WITH %s AS PROCEDURE ("id" NUMBER, "name" VARCHAR) RETURNS TABLE ("country_code" VARCHAR) LANGUAGE SCALA STRICT RUNTIME_VERSION = '2.12' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('test_jar.jar') HANDLER = 'TestFunc.echoVarchar' AS 'return id + name;' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1, rnd) INTO :ret`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), cteId.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		)

	proceduresTests.CreateAndCallForJavaScript.
		withDefaultOpts(func() *CreateAndCallForJavaScriptProcedureOptions {
			return &CreateAndCallForJavaScriptProcedureOptions{
				Name:                proceduresTestIdAccountObjectIdentifier,
				ProcedureName:       proceduresTestIdAccountObjectIdentifier,
				ResultDataType:      dataTypeFloat,
				ProcedureDefinition: "return 1;",
			}
		}).
		withModify(case_Procedures_validation_CreateAndCallForJavaScript_opts_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForJavaScriptProcedureOptions) {
			opts.ResultDataTypeOld = DataTypeFloat
		}).
		withModify(case_Procedures_validation_CreateAndCallForJavaScript_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForJavaScriptProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "d", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat}}
		}).
		withModify(case_Procedures_validation_CreateAndCallForJavaScript_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateAndCallForJavaScriptProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "d", ArgDataTypeOld: DataTypeFloat}, {ArgName: "d2"}}
		}).
		withExpectedSqlf(
			case_Procedures_sql_CreateAndCallForJavaScript_basic,
			`WITH %s AS PROCEDURE () RETURNS FLOAT LANGUAGE JAVASCRIPT AS 'return 1;' CALL %s ()`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_CreateAndCallForJavaScript_all,
			func(opts *CreateAndCallForJavaScriptProcedureOptions) {
				opts.Arguments = []ProcedureArgument{{ArgName: "d", ArgDataType: dataTypeFloat, DefaultValue: new("1.0")}}
				opts.NotNull = new(true)
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ProcedureDefinition = "return 1;"
				opts.WithClauses = []ProcedureWithClause{{CteName: cteId, CteColumns: []string{"x", "y"}, Statement: "(select m.album_ID, m.album_name, b.band_name from music_albums)"}}
				opts.ScriptingVariable = new(":ret")
				opts.CallArguments = []string{"1"}
			},
			`WITH %s AS PROCEDURE ("d" FLOAT DEFAULT 1.0) RETURNS FLOAT NOT NULL LANGUAGE JAVASCRIPT STRICT AS 'return 1;' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1) INTO :ret`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), cteId.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		// TODO [SNOW-1348106]: remove with old procedure removal for V1
		withAdditionalSqlCasef(
			"sql_CreateAndCallForJavaScript_allOldDataTypes",
			func(opts *CreateAndCallForJavaScriptProcedureOptions) {
				opts.Arguments = []ProcedureArgument{{ArgName: "d", ArgDataTypeOld: "DOUBLE", DefaultValue: new("1.0")}}
				opts.ResultDataType = nil
				opts.ResultDataTypeOld = "DOUBLE"
				opts.NotNull = new(true)
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ProcedureDefinition = "return 1;"
				opts.WithClauses = []ProcedureWithClause{{CteName: cteId, CteColumns: []string{"x", "y"}, Statement: "(select m.album_ID, m.album_name, b.band_name from music_albums)"}}
				opts.ScriptingVariable = new(":ret")
				opts.CallArguments = []string{"1"}
			},
			`WITH %s AS PROCEDURE ("d" DOUBLE DEFAULT 1.0) RETURNS DOUBLE NOT NULL LANGUAGE JAVASCRIPT STRICT AS 'return 1;' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1) INTO :ret`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), cteId.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		)

	proceduresTests.CreateAndCallForPython.
		withDefaultOpts(func() *CreateAndCallForPythonProcedureOptions {
			return &CreateAndCallForPythonProcedureOptions{
				Name:           proceduresTestIdAccountObjectIdentifier,
				ProcedureName:  proceduresTestIdAccountObjectIdentifier,
				Returns:        ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataType: dataTypeVarchar_100}},
				Handler:        "udf",
				Packages:       []ProcedurePackage{{ProcedurePackage: "snowflake-snowpark-python"}},
				RuntimeVersion: "3.9",
			}
		}).
		withModify(case_Procedures_validation_CreateAndCallForPython_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForPythonProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "i", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat}}
		}).
		withModify(case_Procedures_validation_CreateAndCallForPython_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateAndCallForPythonProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "i", ArgDataTypeOld: DataTypeFloat}, {ArgName: "i2"}}
		}).
		withModify(case_Procedures_validation_CreateAndCallForPython_opts_Returns_ResultDataType_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForPythonProcedureOptions) {
			opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{ResultDataTypeOld: DataTypeFloat, ResultDataType: dataTypeFloat}
		}).
		withModify(case_Procedures_validation_CreateAndCallForPython_opts_Returns_Table_Columns_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForPythonProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR, ColumnDataType: dataTypeVarchar_100}}}
		}).
		withModify(case_Procedures_validation_CreateAndCallForPython_opts_Returns_Table_Columns_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateAndCallForPythonProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "col2"}}}
		}).
		withExpectedSqlf(
			case_Procedures_sql_CreateAndCallForPython_basic,
			`WITH %s AS PROCEDURE () RETURNS VARCHAR(100) LANGUAGE PYTHON RUNTIME_VERSION = '3.9' PACKAGES = ('snowflake-snowpark-python') HANDLER = 'udf' CALL %s ()`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_CreateAndCallForPython_all,
			func(opts *CreateAndCallForPythonProcedureOptions) {
				opts.Arguments = []ProcedureArgument{{ArgName: "i", ArgDataType: dataTypeNumber_36_2, DefaultValue: new("1")}}
				opts.Returns = ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataType: dataTypeVariant, Null: new(true)}}
				opts.Packages = []ProcedurePackage{{ProcedurePackage: "numpy"}, {ProcedurePackage: "pandas"}}
				opts.Imports = []ProcedureImport{{ProcedureImport: "numpy"}, {ProcedureImport: "pandas"}}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ProcedureDefinition = new("import numpy as np")
				opts.WithClauses = []ProcedureWithClause{{CteName: cteId, CteColumns: []string{"x", "y"}, Statement: "(select m.album_ID, m.album_name, b.band_name from music_albums)"}}
				opts.ScriptingVariable = new(":ret")
				opts.CallArguments = []string{"1"}
			},
			`WITH %s AS PROCEDURE ("i" NUMBER(36, 2) DEFAULT 1) RETURNS VARIANT NULL LANGUAGE PYTHON STRICT RUNTIME_VERSION = '3.9' PACKAGES = ('numpy', 'pandas') IMPORTS = ('numpy', 'pandas') HANDLER = 'udf' AS 'import numpy as np' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1) INTO :ret`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), cteId.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateAndCallForPython_noArguments",
			func(opts *CreateAndCallForPythonProcedureOptions) {
				opts.Returns = ProcedureReturns{Table: &ProcedureReturnsTable{}}
			},
			`WITH %s AS PROCEDURE () RETURNS TABLE () LANGUAGE PYTHON RUNTIME_VERSION = '3.9' PACKAGES = ('snowflake-snowpark-python') HANDLER = 'udf' CALL %s ()`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		// TODO [SNOW-1348106]: remove with old procedure removal for V1
		withAdditionalSqlCasef(
			"sql_CreateAndCallForPython_allOldDataTypes",
			func(opts *CreateAndCallForPythonProcedureOptions) {
				opts.Arguments = []ProcedureArgument{{ArgName: "i", ArgDataTypeOld: "int", DefaultValue: new("1")}}
				opts.Returns = ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataTypeOld: "VARIANT", Null: new(true)}}
				opts.Packages = []ProcedurePackage{{ProcedurePackage: "numpy"}, {ProcedurePackage: "pandas"}}
				opts.Imports = []ProcedureImport{{ProcedureImport: "numpy"}, {ProcedureImport: "pandas"}}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ProcedureDefinition = new("import numpy as np")
				opts.WithClauses = []ProcedureWithClause{{CteName: cteId, CteColumns: []string{"x", "y"}, Statement: "(select m.album_ID, m.album_name, b.band_name from music_albums)"}}
				opts.ScriptingVariable = new(":ret")
				opts.CallArguments = []string{"1"}
			},
			`WITH %s AS PROCEDURE ("i" int DEFAULT 1) RETURNS VARIANT NULL LANGUAGE PYTHON STRICT RUNTIME_VERSION = '3.9' PACKAGES = ('numpy', 'pandas') IMPORTS = ('numpy', 'pandas') HANDLER = 'udf' AS 'import numpy as np' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1) INTO :ret`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), cteId.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		)

	proceduresTests.CreateAndCallForSQL.
		withDefaultOpts(func() *CreateAndCallForSQLProcedureOptions {
			return &CreateAndCallForSQLProcedureOptions{
				Name:                proceduresTestIdAccountObjectIdentifier,
				ProcedureName:       proceduresTestIdAccountObjectIdentifier,
				Returns:             ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataType: dataTypeFloat}},
				ProcedureDefinition: "3.141592654::FLOAT",
			}
		}).
		withModify(case_Procedures_validation_CreateAndCallForSQL_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForSQLProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat}}
		}).
		withModify(case_Procedures_validation_CreateAndCallForSQL_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateAndCallForSQLProcedureOptions) {
			opts.Arguments = []ProcedureArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat}, {ArgName: "arg2"}}
		}).
		withModify(case_Procedures_validation_CreateAndCallForSQL_opts_Returns_ResultDataType_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForSQLProcedureOptions) {
			opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{ResultDataTypeOld: DataTypeFloat, ResultDataType: dataTypeFloat}
		}).
		withModify(case_Procedures_validation_CreateAndCallForSQL_opts_Returns_Table_Columns_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateAndCallForSQLProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR, ColumnDataType: dataTypeVarchar_100}}}
		}).
		withModify(case_Procedures_validation_CreateAndCallForSQL_opts_Returns_Table_Columns_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateAndCallForSQLProcedureOptions) {
			opts.Returns.ResultDataType = nil
			opts.Returns.Table = &ProcedureReturnsTable{Columns: []ProcedureColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "col2"}}}
		}).
		withExpectedSqlf(
			case_Procedures_sql_CreateAndCallForSQL_basic,
			`WITH %s AS PROCEDURE () RETURNS FLOAT LANGUAGE SQL AS '3.141592654::FLOAT' CALL %s ()`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Procedures_sql_CreateAndCallForSQL_all,
			func(opts *CreateAndCallForSQLProcedureOptions) {
				opts.Arguments = []ProcedureArgument{{ArgName: "message", ArgDataType: dataTypeVarchar_100, DefaultValue: new("'test'")}}
				opts.Returns = ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataType: dataTypeFloat}}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ProcedureDefinition = "3.141592654::FLOAT"
				opts.WithClauses = []ProcedureWithClause{{CteName: cteId, CteColumns: []string{"x", "y"}, Statement: "(select m.album_ID, m.album_name, b.band_name from music_albums)"}}
				opts.ScriptingVariable = new(":ret")
				opts.CallArguments = []string{"1"}
			},
			`WITH %s AS PROCEDURE ("message" VARCHAR(100) DEFAULT 'test') RETURNS FLOAT LANGUAGE SQL STRICT AS '3.141592654::FLOAT' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1) INTO :ret`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), cteId.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateAndCallForSQL_noArguments",
			func(opts *CreateAndCallForSQLProcedureOptions) {
				opts.Returns = ProcedureReturns{Table: &ProcedureReturnsTable{}}
			},
			`WITH %s AS PROCEDURE () RETURNS TABLE () LANGUAGE SQL AS '3.141592654::FLOAT' CALL %s ()`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		// TODO [SNOW-1348106]: remove with old procedure removal for V1
		withAdditionalSqlCasef(
			"sql_CreateAndCallForSQL_allOldDataTypes",
			func(opts *CreateAndCallForSQLProcedureOptions) {
				opts.Arguments = []ProcedureArgument{{ArgName: "message", ArgDataTypeOld: "VARCHAR", DefaultValue: new("'test'")}}
				opts.Returns = ProcedureReturns{ResultDataType: &ProcedureReturnsResultDataType{ResultDataTypeOld: DataTypeFloat}}
				opts.NullInputBehavior = new(NullInputBehaviorStrict)
				opts.ProcedureDefinition = "3.141592654::FLOAT"
				opts.WithClauses = []ProcedureWithClause{{CteName: cteId, CteColumns: []string{"x", "y"}, Statement: "(select m.album_ID, m.album_name, b.band_name from music_albums)"}}
				opts.ScriptingVariable = new(":ret")
				opts.CallArguments = []string{"1"}
			},
			`WITH %s AS PROCEDURE ("message" VARCHAR DEFAULT 'test') RETURNS FLOAT LANGUAGE SQL STRICT AS '3.141592654::FLOAT' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1) INTO :ret`,
			proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(), cteId.FullyQualifiedName(), proceduresTestIdAccountObjectIdentifier.FullyQualifiedName(),
		)
}

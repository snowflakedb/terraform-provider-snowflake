package sdk

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/require"
)

func wrapFunctionDefinition(def string) string {
	return fmt.Sprintf(`$$%s$$`, def)
}

func init() {
	javaId := functionsTestIdSchemaObjectIdentifier
	alterDescribeId := functionsTestIdSchemaObjectIdentifierWithArguments
	renameTarget := randomSchemaObjectIdentifier()
	secretId := randomSchemaObjectIdentifier()
	secretId2 := randomSchemaObjectIdentifier()

	functionsTests.CreateForJava.
		withDefaultOpts(func() *CreateForJavaFunctionOptions {
			return &CreateForJavaFunctionOptions{
				name:    javaId,
				Returns: FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataType: dataTypeVarchar_100}},
				Handler: "TestFunc.echoVarchar",
				// additionalValidations() requires Imports to be set when FunctionDefinition is nil.
				Imports: []FunctionImport{{FunctionImport: "@~/my_lib.jar"}},
			}
		}).
		withModify(case_Functions_validation_CreateForJava_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForJavaFunctionOptions) {
			opts.Arguments = []FunctionArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat}}
		}).
		withModify(case_Functions_validation_CreateForJava_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForJavaFunctionOptions) {
			opts.Arguments = []FunctionArgument{{ArgName: "arg", ArgDataTypeOld: DataTypeFloat}, {ArgName: "arg"}}
		}).
		withModify(case_Functions_validation_CreateForJava_opts_Returns_ResultDataType_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForJavaFunctionOptions) {
			opts.Returns.ResultDataType = &FunctionReturnsResultDataType{ResultDataTypeOld: DataTypeFloat, ResultDataType: dataTypeFloat}
		}).
		withModify(case_Functions_validation_CreateForJava_opts_Returns_Table_Columns_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForJavaFunctionOptions) {
			opts.Returns.Table = &FunctionReturnsTable{Columns: []FunctionColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR, ColumnDataType: dataTypeVarchar_100}}}
		}).
		withModify(case_Functions_validation_CreateForJava_opts_Returns_Table_Columns_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForJavaFunctionOptions) {
			opts.Returns.Table = &FunctionReturnsTable{Columns: []FunctionColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "col2"}}}
		}).
		withAdditionalValidationCase(
			"validation_CreateForJava_targetPath",
			func(opts *CreateForJavaFunctionOptions) {
				opts.TargetPath = String("@~/testfunc.jar")
			},
			NewError("TARGET_PATH must be nil when AS is nil"),
		).
		withAdditionalValidationCase(
			"validation_CreateForJava_imports",
			func(opts *CreateForJavaFunctionOptions) {
				opts.Imports = nil
			},
			NewError("IMPORTS must not be empty when AS is nil"),
		).
		withExpectedSqlf(
			case_Functions_sql_CreateForJava_basic,
			`CREATE FUNCTION %s () RETURNS VARCHAR(100) LANGUAGE JAVA IMPORTS = ('@~/my_lib.jar') HANDLER = 'TestFunc.echoVarchar'`,
			javaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_CreateForJava_all,
			func(opts *CreateForJavaFunctionOptions) {
				opts.OrReplace = Bool(true)
				opts.Temporary = Bool(true)
				opts.Secure = Bool(true)
				opts.Arguments = []FunctionArgument{{ArgName: "id", ArgDataType: dataTypeNumber_36_2}, {ArgName: "name", ArgDataType: dataTypeVarchar_100, DefaultValue: String("'test'")}}
				opts.CopyGrants = Bool(true)
				opts.Returns = FunctionReturns{Table: &FunctionReturnsTable{Columns: []FunctionColumn{{ColumnName: "country_code", ColumnDataType: dataTypeVarchar_100}, {ColumnName: "country_name", ColumnDataType: dataTypeVarchar_100}}}}
				opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
				opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
				opts.RuntimeVersion = String("2.0")
				opts.Comment = String("comment")
				opts.Imports = []FunctionImport{{FunctionImport: "@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar"}}
				opts.Packages = []FunctionPackage{{FunctionPackage: "com.snowflake:snowpark:1.2.0"}}
				opts.ExternalAccessIntegrations = []AccountObjectIdentifier{NewAccountObjectIdentifier("ext_integration")}
				opts.Secrets = []SecretReference{{VariableName: "variable1", Name: secretId}, {VariableName: "variable2", Name: secretId2}}
				opts.TargetPath = String("@~/testfunc.jar")
				opts.FunctionDefinition = String(wrapFunctionDefinition("return id + name;"))
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("id" NUMBER(36, 2), "name" VARCHAR(100) DEFAULT 'test') COPY GRANTS RETURNS TABLE ("country_code" VARCHAR(100), "country_name" VARCHAR(100)) NOT NULL LANGUAGE JAVA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar') PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'TestFunc.echoVarchar' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) TARGET_PATH = '@~/testfunc.jar' AS $$return id + name;$$`,
			javaId.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		).
		// TODO [SNOW-1348103]: remove with old function removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForJava_allOldDataTypes",
			func(opts *CreateForJavaFunctionOptions) {
				opts.OrReplace = Bool(true)
				opts.Temporary = Bool(true)
				opts.Secure = Bool(true)
				opts.Arguments = []FunctionArgument{{ArgName: "id", ArgDataTypeOld: DataTypeNumber}, {ArgName: "name", ArgDataTypeOld: DataTypeVARCHAR, DefaultValue: String("'test'")}}
				opts.CopyGrants = Bool(true)
				opts.Returns = FunctionReturns{Table: &FunctionReturnsTable{Columns: []FunctionColumn{{ColumnName: "country_code", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "country_name", ColumnDataTypeOld: DataTypeVARCHAR}}}}
				opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
				opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
				opts.RuntimeVersion = String("2.0")
				opts.Comment = String("comment")
				opts.Imports = []FunctionImport{{FunctionImport: "@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar"}}
				opts.Packages = []FunctionPackage{{FunctionPackage: "com.snowflake:snowpark:1.2.0"}}
				opts.ExternalAccessIntegrations = []AccountObjectIdentifier{NewAccountObjectIdentifier("ext_integration")}
				opts.Secrets = []SecretReference{{VariableName: "variable1", Name: secretId}, {VariableName: "variable2", Name: secretId2}}
				opts.TargetPath = String("@~/testfunc.jar")
				opts.FunctionDefinition = String(wrapFunctionDefinition("return id + name;"))
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("id" NUMBER, "name" VARCHAR DEFAULT 'test') COPY GRANTS RETURNS TABLE ("country_code" VARCHAR, "country_name" VARCHAR) NOT NULL LANGUAGE JAVA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar') PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'TestFunc.echoVarchar' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) TARGET_PATH = '@~/testfunc.jar' AS $$return id + name;$$`,
			javaId.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		)

	functionsTests.CreateForJavascript.
		withDefaultOpts(func() *CreateForJavascriptFunctionOptions {
			return &CreateForJavascriptFunctionOptions{
				name: javaId,
				Returns: FunctionReturns{
					ResultDataType: &FunctionReturnsResultDataType{ResultDataType: dataTypeFloat},
				},
				FunctionDefinition: wrapFunctionDefinition("return 1;"),
			}
		}).
		withModify(case_Functions_validation_CreateForJavascript_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForJavascriptFunctionOptions) {
			opts.Arguments = []FunctionArgument{{ArgName: "d", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat}}
		}).
		withModify(case_Functions_validation_CreateForJavascript_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForJavascriptFunctionOptions) {
			opts.Arguments = []FunctionArgument{{ArgName: "d", ArgDataTypeOld: DataTypeFloat}, {ArgName: "d2"}}
		}).
		withModify(case_Functions_validation_CreateForJavascript_opts_Returns_ResultDataType_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForJavascriptFunctionOptions) {
			opts.Returns.ResultDataType = &FunctionReturnsResultDataType{ResultDataTypeOld: DataTypeFloat, ResultDataType: dataTypeFloat}
		}).
		withModify(case_Functions_validation_CreateForJavascript_opts_Returns_Table_Columns_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForJavascriptFunctionOptions) {
			opts.Returns.Table = &FunctionReturnsTable{Columns: []FunctionColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR, ColumnDataType: dataTypeVarchar_100}}}
		}).
		withModify(case_Functions_validation_CreateForJavascript_opts_Returns_Table_Columns_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForJavascriptFunctionOptions) {
			opts.Returns.Table = &FunctionReturnsTable{Columns: []FunctionColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "col2"}}}
		}).
		withExpectedSqlf(
			case_Functions_sql_CreateForJavascript_basic,
			`CREATE FUNCTION %s () RETURNS FLOAT LANGUAGE JAVASCRIPT AS $$return 1;$$`,
			javaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_CreateForJavascript_all,
			func(opts *CreateForJavascriptFunctionOptions) {
				opts.OrReplace = Bool(true)
				opts.Temporary = Bool(true)
				opts.Secure = Bool(true)
				opts.Arguments = []FunctionArgument{{ArgName: "d", ArgDataType: dataTypeFloat, DefaultValue: String("1.0")}}
				opts.CopyGrants = Bool(true)
				opts.Returns = FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataType: dataTypeFloat}}
				opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
				opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
				opts.Comment = String("comment")
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("d" FLOAT DEFAULT 1.0) COPY GRANTS RETURNS FLOAT NOT NULL LANGUAGE JAVASCRIPT CALLED ON NULL INPUT IMMUTABLE COMMENT = 'comment' AS $$return 1;$$`,
			javaId.FullyQualifiedName(),
		).
		// TODO [SNOW-1348103]: remove with old function removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForJavascript_allOldDataTypes",
			func(opts *CreateForJavascriptFunctionOptions) {
				opts.OrReplace = Bool(true)
				opts.Temporary = Bool(true)
				opts.Secure = Bool(true)
				opts.Arguments = []FunctionArgument{{ArgName: "d", ArgDataTypeOld: DataTypeFloat, DefaultValue: String("1.0")}}
				opts.CopyGrants = Bool(true)
				opts.Returns = FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataTypeOld: DataTypeFloat}}
				opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
				opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
				opts.Comment = String("comment")
				opts.FunctionDefinition = wrapFunctionDefinition("return 1;")
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("d" FLOAT DEFAULT 1.0) COPY GRANTS RETURNS FLOAT NOT NULL LANGUAGE JAVASCRIPT CALLED ON NULL INPUT IMMUTABLE COMMENT = 'comment' AS $$return 1;$$`,
			javaId.FullyQualifiedName(),
		)

	functionsTests.CreateForPython.
		withDefaultOpts(func() *CreateForPythonFunctionOptions {
			return &CreateForPythonFunctionOptions{
				name:           javaId,
				Returns:        FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataType: dataTypeVariant}},
				RuntimeVersion: "3.9",
				Handler:        "udf",
				// additionalValidations() requires Imports to be set when FunctionDefinition is nil.
				Imports: []FunctionImport{{FunctionImport: "numpy"}},
			}
		}).
		withModify(case_Functions_validation_CreateForPython_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForPythonFunctionOptions) {
			opts.Arguments = []FunctionArgument{{ArgName: "i", ArgDataTypeOld: DataTypeNumber, ArgDataType: dataTypeNumber_36_2}}
		}).
		withModify(case_Functions_validation_CreateForPython_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForPythonFunctionOptions) {
			opts.Arguments = []FunctionArgument{{ArgName: "i", ArgDataTypeOld: DataTypeNumber}, {ArgName: "i2"}}
		}).
		withModify(case_Functions_validation_CreateForPython_opts_Returns_ResultDataType_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForPythonFunctionOptions) {
			opts.Returns.ResultDataType = &FunctionReturnsResultDataType{ResultDataTypeOld: DataTypeVariant, ResultDataType: dataTypeVariant}
		}).
		withModify(case_Functions_validation_CreateForPython_opts_Returns_Table_Columns_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForPythonFunctionOptions) {
			opts.Returns.Table = &FunctionReturnsTable{Columns: []FunctionColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR, ColumnDataType: dataTypeVarchar_100}}}
		}).
		withModify(case_Functions_validation_CreateForPython_opts_Returns_Table_Columns_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForPythonFunctionOptions) {
			opts.Returns.Table = &FunctionReturnsTable{Columns: []FunctionColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "col2"}}}
		}).
		withAdditionalValidationCase(
			"validation_CreateForPython_functionDefinition",
			func(opts *CreateForPythonFunctionOptions) {
				opts.Packages = []FunctionPackage{{FunctionPackage: "numpy"}}
				opts.Imports = nil // clear default imports so additionalValidations error fires
			},
			NewError("IMPORTS must not be empty when AS is nil"),
		).
		withExpectedSqlf(
			case_Functions_sql_CreateForPython_basic,
			`CREATE FUNCTION %s () RETURNS VARIANT LANGUAGE PYTHON RUNTIME_VERSION = '3.9' IMPORTS = ('numpy') HANDLER = 'udf'`,
			javaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_CreateForPython_all,
			func(opts *CreateForPythonFunctionOptions) {
				opts.OrReplace = Bool(true)
				opts.Temporary = Bool(true)
				opts.Secure = Bool(true)
				opts.Arguments = []FunctionArgument{{ArgName: "i", ArgDataType: dataTypeNumber_36_2, DefaultValue: String("1")}}
				opts.CopyGrants = Bool(true)
				opts.Returns = FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataType: dataTypeVariant}}
				opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
				opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
				opts.Comment = String("comment")
				opts.Imports = []FunctionImport{{FunctionImport: "numpy"}, {FunctionImport: "pandas"}}
				opts.Packages = []FunctionPackage{{FunctionPackage: "numpy"}, {FunctionPackage: "pandas"}}
				opts.ExternalAccessIntegrations = []AccountObjectIdentifier{NewAccountObjectIdentifier("ext_integration")}
				opts.Secrets = []SecretReference{{VariableName: "variable1", Name: secretId}, {VariableName: "variable2", Name: secretId2}}
				opts.FunctionDefinition = String(wrapFunctionDefinition("import numpy as np"))
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("i" NUMBER(36, 2) DEFAULT 1) COPY GRANTS RETURNS VARIANT NOT NULL LANGUAGE PYTHON CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '3.9' COMMENT = 'comment' IMPORTS = ('numpy', 'pandas') PACKAGES = ('numpy', 'pandas') HANDLER = 'udf' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) AS $$import numpy as np$$`,
			javaId.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		).
		// TODO [SNOW-1348103]: remove with old function removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForPython_allOldDataTypes",
			func(opts *CreateForPythonFunctionOptions) {
				opts.OrReplace = Bool(true)
				opts.Temporary = Bool(true)
				opts.Secure = Bool(true)
				opts.Arguments = []FunctionArgument{{ArgName: "i", ArgDataTypeOld: DataTypeNumber, DefaultValue: String("1")}}
				opts.CopyGrants = Bool(true)
				opts.Returns = FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataTypeOld: DataTypeVariant}}
				opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
				opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
				opts.RuntimeVersion = "3.9"
				opts.Comment = String("comment")
				opts.Imports = []FunctionImport{{FunctionImport: "numpy"}, {FunctionImport: "pandas"}}
				opts.Packages = []FunctionPackage{{FunctionPackage: "numpy"}, {FunctionPackage: "pandas"}}
				opts.ExternalAccessIntegrations = []AccountObjectIdentifier{NewAccountObjectIdentifier("ext_integration")}
				opts.Secrets = []SecretReference{{VariableName: "variable1", Name: secretId}, {VariableName: "variable2", Name: secretId2}}
				opts.FunctionDefinition = String(wrapFunctionDefinition("import numpy as np"))
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("i" NUMBER DEFAULT 1) COPY GRANTS RETURNS VARIANT NOT NULL LANGUAGE PYTHON CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '3.9' COMMENT = 'comment' IMPORTS = ('numpy', 'pandas') PACKAGES = ('numpy', 'pandas') HANDLER = 'udf' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) AS $$import numpy as np$$`,
			javaId.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		)

	functionsTests.CreateForScala.
		withDefaultOpts(func() *CreateForScalaFunctionOptions {
			return &CreateForScalaFunctionOptions{
				name:           javaId,
				ResultDataType: dataTypeVarchar_100,
				Handler:        "Echo.echoVarchar",
				// additionalValidations() requires Imports to be set when FunctionDefinition is nil.
				Imports: []FunctionImport{{FunctionImport: "@udf_libs/echohandler.jar"}},
			}
		}).
		withModify(case_Functions_validation_CreateForScala_opts_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForScalaFunctionOptions) {
			opts.ResultDataTypeOld = DataTypeFloat
			opts.ResultDataType = dataTypeFloat
		}).
		withModify(case_Functions_validation_CreateForScala_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForScalaFunctionOptions) {
			opts.Arguments = []FunctionArgument{{ArgName: "x", ArgDataTypeOld: DataTypeVARCHAR, ArgDataType: dataTypeVarchar_100}}
		}).
		withModify(case_Functions_validation_CreateForScala_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForScalaFunctionOptions) {
			opts.Arguments = []FunctionArgument{{ArgName: "x", ArgDataTypeOld: DataTypeVARCHAR}, {ArgName: "y"}}
		}).
		withAdditionalValidationCase(
			"validation_CreateForScala_targetPath",
			func(opts *CreateForScalaFunctionOptions) {
				opts.TargetPath = String("@~/testfunc.jar")
			},
			NewError("TARGET_PATH must be nil when AS is nil"),
		).
		withAdditionalValidationCase(
			"validation_CreateForScala_imports",
			func(opts *CreateForScalaFunctionOptions) {
				opts.Imports = nil
				// TargetPath is nil — only IMPORTS error fires
			},
			NewError("IMPORTS must not be empty when AS is nil"),
		).
		withExpectedSqlf(
			case_Functions_sql_CreateForScala_basic,
			`CREATE FUNCTION %s () RETURNS VARCHAR(100) LANGUAGE SCALA RUNTIME_VERSION = '' IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar'`,
			javaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_CreateForScala_all,
			func(opts *CreateForScalaFunctionOptions) {
				opts.OrReplace = Bool(true)
				opts.Temporary = Bool(true)
				opts.Secure = Bool(true)
				opts.Arguments = []FunctionArgument{{ArgName: "x", ArgDataType: dataTypeVarchar_100, DefaultValue: String("'test'")}}
				opts.CopyGrants = Bool(true)
				opts.ResultDataType = dataTypeVarchar_100
				opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
				opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
				opts.RuntimeVersion = "2.0"
				opts.Comment = String("comment")
				opts.Imports = []FunctionImport{{FunctionImport: "@udf_libs/echohandler.jar"}}
				opts.FunctionDefinition = String(wrapFunctionDefinition("return x"))
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("x" VARCHAR(100) DEFAULT 'test') COPY GRANTS RETURNS VARCHAR(100) NOT NULL LANGUAGE SCALA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar' AS $$return x$$`,
			javaId.FullyQualifiedName(),
		).
		// TODO [SNOW-1348103]: remove with old function removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForScala_allOldDataTypes",
			func(opts *CreateForScalaFunctionOptions) {
				opts.OrReplace = Bool(true)
				opts.Temporary = Bool(true)
				opts.Secure = Bool(true)
				opts.Arguments = []FunctionArgument{{ArgName: "x", ArgDataTypeOld: DataTypeVARCHAR, DefaultValue: String("'test'")}}
				opts.CopyGrants = Bool(true)
				opts.ResultDataType = nil // clear new-style data type set by default opts
				opts.ResultDataTypeOld = DataTypeVARCHAR
				opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
				opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
				opts.RuntimeVersion = "2.0"
				opts.Comment = String("comment")
				opts.Imports = []FunctionImport{{FunctionImport: "@udf_libs/echohandler.jar"}}
				opts.FunctionDefinition = String(wrapFunctionDefinition("return x"))
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("x" VARCHAR DEFAULT 'test') COPY GRANTS RETURNS VARCHAR NOT NULL LANGUAGE SCALA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar' AS $$return x$$`,
			javaId.FullyQualifiedName(),
		)

	functionsTests.CreateForSQL.
		withDefaultOpts(func() *CreateForSQLFunctionOptions {
			return &CreateForSQLFunctionOptions{
				name:               javaId,
				Returns:            FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataType: dataTypeFloat}},
				FunctionDefinition: wrapFunctionDefinition("3.141592654::FLOAT"),
			}
		}).
		withModify(case_Functions_validation_CreateForSQL_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForSQLFunctionOptions) {
			opts.Arguments = []FunctionArgument{{ArgName: "message", ArgDataTypeOld: "VARCHAR", ArgDataType: dataTypeVarchar_100}}
		}).
		withModify(case_Functions_validation_CreateForSQL_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForSQLFunctionOptions) {
			opts.Arguments = []FunctionArgument{{ArgName: "message", ArgDataTypeOld: "VARCHAR"}, {ArgName: "msg2"}}
		}).
		withModify(case_Functions_validation_CreateForSQL_opts_Returns_ResultDataType_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForSQLFunctionOptions) {
			opts.Returns.ResultDataType = &FunctionReturnsResultDataType{ResultDataTypeOld: DataTypeFloat, ResultDataType: dataTypeFloat}
		}).
		withModify(case_Functions_validation_CreateForSQL_opts_Returns_Table_Columns_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateForSQLFunctionOptions) {
			opts.Returns.Table = &FunctionReturnsTable{Columns: []FunctionColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR, ColumnDataType: dataTypeVarchar_100}}}
		}).
		withModify(case_Functions_validation_CreateForSQL_opts_Returns_Table_Columns_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateForSQLFunctionOptions) {
			opts.Returns.Table = &FunctionReturnsTable{Columns: []FunctionColumn{{ColumnName: "col", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "col2"}}}
		}).
		withExpectedSqlf(
			case_Functions_sql_CreateForSQL_basic,
			`CREATE FUNCTION %s () RETURNS FLOAT AS $$3.141592654::FLOAT$$`,
			javaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_CreateForSQL_all,
			func(opts *CreateForSQLFunctionOptions) {
				opts.OrReplace = Bool(true)
				opts.Temporary = Bool(true)
				opts.Secure = Bool(true)
				opts.Arguments = []FunctionArgument{{ArgName: "message", ArgDataType: dataTypeVarchar_100, DefaultValue: String("'test'")}}
				opts.CopyGrants = Bool(true)
				opts.Returns = FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataType: dataTypeFloat}}
				opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
				opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
				opts.Memoizable = Bool(true)
				opts.Comment = String("comment")
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("message" VARCHAR(100) DEFAULT 'test') COPY GRANTS RETURNS FLOAT NOT NULL IMMUTABLE MEMOIZABLE COMMENT = 'comment' AS $$3.141592654::FLOAT$$`,
			javaId.FullyQualifiedName(),
		).
		// TODO [SNOW-1348103]: remove with old function removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForSQL_allOldDataTypes",
			func(opts *CreateForSQLFunctionOptions) {
				opts.OrReplace = Bool(true)
				opts.Temporary = Bool(true)
				opts.Secure = Bool(true)
				opts.Arguments = []FunctionArgument{{ArgName: "message", ArgDataTypeOld: "VARCHAR", DefaultValue: String("'test'")}}
				opts.CopyGrants = Bool(true)
				opts.Returns = FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataTypeOld: DataTypeFloat}}
				opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
				opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
				opts.Memoizable = Bool(true)
				opts.Comment = String("comment")
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("message" VARCHAR DEFAULT 'test') COPY GRANTS RETURNS FLOAT NOT NULL IMMUTABLE MEMOIZABLE COMMENT = 'comment' AS $$3.141592654::FLOAT$$`,
			javaId.FullyQualifiedName(),
		)

	functionsTests.Alter.
		withDefaultOpts(func() *AlterFunctionOptions {
			return &AlterFunctionOptions{name: alterDescribeId, IfExists: Bool(true)}
		}).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_RenameTo,
			func(opts *AlterFunctionOptions) { opts.RenameTo = &renameTarget },
			`ALTER FUNCTION IF EXISTS %s RENAME TO %s`,
			alterDescribeId.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_Set,
			func(opts *AlterFunctionOptions) {
				opts.Set = &FunctionSet{Comment: String("comment"), TraceLevel: Pointer(TraceLevelOff)}
			},
			`ALTER FUNCTION IF EXISTS %s SET COMMENT = 'comment', TRACE_LEVEL = 'OFF'`,
			alterDescribeId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_emptySecrets",
			func(opts *AlterFunctionOptions) {
				opts.Set = &FunctionSet{SecretsList: &SecretsList{}}
			},
			`ALTER FUNCTION IF EXISTS %s SET SECRETS = ()`, alterDescribeId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_nonEmptySecrets",
			func(opts *AlterFunctionOptions) {
				opts.Set = &FunctionSet{SecretsList: &SecretsList{[]SecretReference{{VariableName: "abc", Name: secretId}}}}
			},
			`ALTER FUNCTION IF EXISTS %s SET SECRETS = ('abc' = %s)`, alterDescribeId.FullyQualifiedName(), secretId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_Unset,
			func(opts *AlterFunctionOptions) {
				opts.Unset = &FunctionUnset{Comment: Bool(true), TraceLevel: Bool(true)}
			},
			`ALTER FUNCTION IF EXISTS %s UNSET COMMENT, TRACE_LEVEL`,
			alterDescribeId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_SetSecure,
			func(opts *AlterFunctionOptions) { opts.SetSecure = new(true) },
			`ALTER FUNCTION IF EXISTS %s SET SECURE`, alterDescribeId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_UnsetSecure,
			func(opts *AlterFunctionOptions) { opts.UnsetSecure = new(true) },
			`ALTER FUNCTION IF EXISTS %s UNSET SECURE`, alterDescribeId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_SetTags,
			func(opts *AlterFunctionOptions) {
				opts.SetTags = []TagAssociation{{Name: NewAccountObjectIdentifier("tag1"), Value: "value1"}}
			},
			`ALTER FUNCTION IF EXISTS %s SET TAG "tag1" = 'value1'`, alterDescribeId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_UnsetTags,
			func(opts *AlterFunctionOptions) {
				opts.UnsetTags = []ObjectIdentifier{NewAccountObjectIdentifier("tag1"), NewAccountObjectIdentifier("tag2")}
			},
			`ALTER FUNCTION IF EXISTS %s UNSET TAG "tag1", "tag2"`, alterDescribeId.FullyQualifiedName(),
		)

	functionsTests.Drop.
		withExpectedSqlf(case_Functions_sql_Drop_basic,
			`DROP FUNCTION %s`, alterDescribeId.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Drop_all,
			func(opts *DropFunctionOptions) { opts.IfExists = Bool(true) },
			`DROP FUNCTION IF EXISTS %s`, alterDescribeId.FullyQualifiedName(),
		)

	functionsTests.Show.
		withExpectedSql(case_Functions_sql_Show_basic, `SHOW USER FUNCTIONS`).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Show_all,
			func(opts *ShowFunctionOptions) {
				opts.Like = &Like{Pattern: String("pattern")}
				opts.In = &ExtendedIn{In: In{Account: Bool(true)}}
			},
			`SHOW USER FUNCTIONS LIKE 'pattern' IN ACCOUNT`,
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Show_Like,
			func(opts *ShowFunctionOptions) { opts.Like = &Like{Pattern: String("pattern")} },
			`SHOW USER FUNCTIONS LIKE 'pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Show_In,
			func(opts *ShowFunctionOptions) { opts.In = &ExtendedIn{In: In{Account: Bool(true)}} },
			`SHOW USER FUNCTIONS IN ACCOUNT`,
		)

	functionsTests.Describe.
		withDefaultOpts(func() *DescribeFunctionOptions {
			return &DescribeFunctionOptions{name: alterDescribeId}
		}).
		withExpectedSqlf(case_Functions_sql_Describe_basic,
			`DESCRIBE FUNCTION %s`, alterDescribeId.FullyQualifiedName())
}

// TODO [SNOW-1850370]: test parsing single
func Test_parseFunctionOrProcedureImports(t *testing.T) {
	inputs := []struct {
		rawInput string
		expected []NormalizedPath
	}{
		{"", []NormalizedPath{}},
		{`[]`, []NormalizedPath{}},
		{`[@~/abc]`, []NormalizedPath{{"~", "abc"}}},
		{`[@~/abc/def]`, []NormalizedPath{{"~", "abc/def"}}},
		{`[@"db"."sc"."st"/abc/def]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}}},
		{`[@db.sc.st/abc/def]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}}},
		{`[db.sc.st/abc/def]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}}},
		{`[@"db"."sc".st/abc/def]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}}},
		{`[@"db"."sc".st/abc/def, db."sc".st/abc]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}, {`"db"."sc"."st"`, "abc"}}},
	}

	badInputs := []struct {
		rawInput          string
		expectedErrorPart string
	}{
		{"[", "wrapping brackets not found"},
		{"]", "wrapping brackets not found"},
		{`[@~/]`, "contains empty path"},
		{`[@~]`, "cannot be split into stage and path"},
		{`[@"db"."sc"/abc]`, "contains incorrect stage location"},
		{`[@"db"/abc]`, "contains incorrect stage location"},
		{`[@"db"."sc"."st"."smth"/abc]`, "contains incorrect stage location"},
		{`[@"db/a"."sc"."st"/abc]`, "contains incorrect stage location"},
		{`[@"db"."sc"."st"/abc], @"db"."sc"/abc]`, "contains incorrect stage location"},
	}

	for _, tc := range inputs {
		t.Run(fmt.Sprintf("Snowflake raw imports: %s", tc.rawInput), func(t *testing.T) {
			results, err := parseFunctionOrProcedureImports(&tc.rawInput)
			require.NoError(t, err)
			require.Equal(t, tc.expected, results)
		})
	}

	for _, tc := range badInputs {
		t.Run(fmt.Sprintf("incorrect Snowflake input: %s, expecting error with: %s", tc.rawInput, tc.expectedErrorPart), func(t *testing.T) {
			_, err := parseFunctionOrProcedureImports(&tc.rawInput)
			require.Error(t, err)
			require.ErrorContains(t, err, "could not parse imports from Snowflake")
			require.ErrorContains(t, err, tc.expectedErrorPart)
		})
	}

	t.Run("Snowflake raw imports nil", func(t *testing.T) {
		results, err := parseFunctionOrProcedureImports(nil)
		require.NoError(t, err)
		require.Equal(t, []NormalizedPath{}, results)
	})
}

func Test_parseFunctionOrProcedureReturns(t *testing.T) {
	inputs := []struct {
		rawInput              string
		expectedRawDataType   string
		expectedReturnNotNull bool
	}{
		{"CHAR", "CHAR(1)", false},
		{"CHAR(1)", "CHAR(1)", false},
		{"NUMBER(30, 2)", "NUMBER(30, 2)", false},
		{"NUMBER(30,2)", "NUMBER(30, 2)", false},
		{"NUMBER(30,2) NOT NULL", "NUMBER(30, 2)", true},
		{"CHAR NOT NULL", "CHAR(1)", true},
		{"  CHAR   NOT NULL  ", "CHAR(1)", true},
		{"OBJECT", "OBJECT", false},
		{"OBJECT NOT NULL", "OBJECT", true},
	}

	badInputs := []struct {
		rawInput          string
		expectedErrorPart string
	}{
		{"", "invalid data type"},
		{"NOT NULL", "invalid data type"},
		{"CHA NOT NULL", "invalid data type"},
		{"CHA NOT NULLS", "invalid data type"},
	}

	for _, tc := range inputs {
		t.Run(fmt.Sprintf("return data type raw: %s", tc.rawInput), func(t *testing.T) {
			dt, returnNotNull, err := parseFunctionOrProcedureReturns(tc.rawInput)
			require.NoError(t, err)
			require.Equal(t, tc.expectedRawDataType, dt.ToSql())
			require.Equal(t, tc.expectedReturnNotNull, returnNotNull)
		})
	}

	for _, tc := range badInputs {
		t.Run(fmt.Sprintf("incorrect return data type raw: %s, expecting error with: %s", tc.rawInput, tc.expectedErrorPart), func(t *testing.T) {
			_, _, err := parseFunctionOrProcedureReturns(tc.rawInput)
			require.Error(t, err)
			require.ErrorContains(t, err, tc.expectedErrorPart)
		})
	}
}

func Test_parseFunctionOrProcedureSignature(t *testing.T) {
	inputs := []struct {
		rawInput     string
		expectedArgs []NormalizedArgument
	}{
		{"()", []NormalizedArgument{}},
		{"(abc CHAR)", []NormalizedArgument{{"abc", dataTypeChar}}},
		{"(abc CHAR(1))", []NormalizedArgument{{"abc", dataTypeChar}}},
		{"(abc CHAR(100))", []NormalizedArgument{{"abc", dataTypeChar_100}}},
		{"  (   abc CHAR(100  )  )", []NormalizedArgument{{"abc", dataTypeChar_100}}},
		{"(  abc   CHAR  )", []NormalizedArgument{{"abc", dataTypeChar}}},
		{"(abc DOUBLE PRECISION)", []NormalizedArgument{{"abc", dataTypeDoublePrecision}}},
		{"(abc double precision)", []NormalizedArgument{{"abc", dataTypeDoublePrecision}}},
		{"(abc TIMESTAMP WITHOUT TIME ZONE(5))", []NormalizedArgument{{"abc", dataTypeTimestampWithoutTimeZone_5}}},
	}

	badInputs := []struct {
		rawInput          string
		expectedErrorPart string
	}{
		{"", "can't be empty"},
		{"(abc CHAR", "wrapping parentheses not found"},
		{"abc CHAR)", "wrapping parentheses not found"},
		{"(abc)", "cannot be split into arg name, data type, and default"},
		{"(CHAR)", "cannot be split into arg name, data type, and default"},
		{"(abc CHA)", "invalid data type"},
		{"(abc CHA(123))", "invalid data type"},
		{"(abc CHAR(1) DEFAULT)", "cannot be parsed"},
		{"(abc CHAR(1) DEFAULT 'a')", "cannot be parsed"},
		// TODO [SNOW-1850370]: Snowflake currently does not return concrete data types so we can fail on them currently but it should be improved in the future
		{"(abc NUMBER(30,2))", "cannot be parsed"},
		{"(abc NUMBER(30, 2))", "cannot be parsed"},
	}

	for _, tc := range inputs {
		t.Run(fmt.Sprintf("return data type raw: %s", tc.rawInput), func(t *testing.T) {
			args, err := parseFunctionOrProcedureSignature(tc.rawInput)

			require.NoError(t, err)
			require.Len(t, args, len(tc.expectedArgs))
			for i, arg := range args {
				require.Equal(t, tc.expectedArgs[i].Name, arg.Name)
				require.True(t, datatypes.AreTheSame(tc.expectedArgs[i].DataType, arg.DataType))
			}
		})
	}

	for _, tc := range badInputs {
		t.Run(fmt.Sprintf("incorrect signature raw: %s, expecting error with: %s", tc.rawInput, tc.expectedErrorPart), func(t *testing.T) {
			_, err := parseFunctionOrProcedureSignature(tc.rawInput)
			require.Error(t, err)
			require.ErrorContains(t, err, "could not parse signature from Snowflake")
			require.ErrorContains(t, err, tc.expectedErrorPart)
		})
	}
}

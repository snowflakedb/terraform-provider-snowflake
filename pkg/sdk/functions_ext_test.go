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
	newId := randomSchemaObjectIdentifier()
	secretId := randomSchemaObjectIdentifier()
	secretId2 := randomSchemaObjectIdentifier()

	functionsTests.CreateForJava.
		withDefaultOpts(func() *CreateForJavaFunctionOptions {
			return &CreateForJavaFunctionOptions{
				name:    functionsTestIdSchemaObjectIdentifier,
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
				opts.TargetPath = new("@~/testfunc.jar")
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
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_CreateForJava_all,
			func(opts *CreateForJavaFunctionOptions) {
				opts.IfNotExists = new(true)
				opts.Temporary = new(true)
				opts.Secure = new(true)
				opts.Arguments = []FunctionArgument{{ArgName: "id", ArgDataType: dataTypeNumber_36_2}, {ArgName: "name", ArgDataType: dataTypeVarchar_100, DefaultValue: new("'test'")}}
				opts.Returns = FunctionReturns{Table: &FunctionReturnsTable{Columns: []FunctionColumn{{ColumnName: "country_code", ColumnDataType: dataTypeVarchar_100}, {ColumnName: "country_name", ColumnDataType: dataTypeVarchar_100}}}}
				opts.ReturnNullValues = new(ReturnNullValuesNotNull)
				opts.NullInputBehavior = new(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.RuntimeVersion = new("2.0")
				opts.Comment = new("comment")
				opts.Imports = []FunctionImport{{FunctionImport: "@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar"}}
				opts.Packages = []FunctionPackage{{FunctionPackage: "com.snowflake:snowpark:1.2.0"}}
				opts.ExternalAccessIntegrations = []AccountObjectIdentifier{NewAccountObjectIdentifier("ext_integration")}
				opts.Secrets = []SecretReference{{VariableName: "variable1", Name: secretId}, {VariableName: "variable2", Name: secretId2}}
				opts.TargetPath = new("@~/testfunc.jar")
				opts.FunctionDefinition = new(wrapFunctionDefinition("return id + name;"))
			},
			`CREATE TEMPORARY SECURE FUNCTION IF NOT EXISTS %s ("id" NUMBER(36, 2), "name" VARCHAR(100) DEFAULT 'test') RETURNS TABLE ("country_code" VARCHAR(100), "country_name" VARCHAR(100)) NOT NULL LANGUAGE JAVA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar') PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'TestFunc.echoVarchar' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) TARGET_PATH = '@~/testfunc.jar' AS $$return id + name;$$`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		).
		// TODO [SNOW-1348103]: remove with old function removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForJava_allOldDataTypes",
			func(opts *CreateForJavaFunctionOptions) {
				opts.OrReplace = new(true)
				opts.Temporary = new(true)
				opts.Secure = new(true)
				opts.Arguments = []FunctionArgument{{ArgName: "id", ArgDataTypeOld: DataTypeNumber}, {ArgName: "name", ArgDataTypeOld: DataTypeVARCHAR, DefaultValue: new("'test'")}}
				opts.CopyGrants = new(true)
				opts.Returns = FunctionReturns{Table: &FunctionReturnsTable{Columns: []FunctionColumn{{ColumnName: "country_code", ColumnDataTypeOld: DataTypeVARCHAR}, {ColumnName: "country_name", ColumnDataTypeOld: DataTypeVARCHAR}}}}
				opts.ReturnNullValues = new(ReturnNullValuesNotNull)
				opts.NullInputBehavior = new(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.RuntimeVersion = new("2.0")
				opts.Comment = new("comment")
				opts.Imports = []FunctionImport{{FunctionImport: "@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar"}}
				opts.Packages = []FunctionPackage{{FunctionPackage: "com.snowflake:snowpark:1.2.0"}}
				opts.ExternalAccessIntegrations = []AccountObjectIdentifier{NewAccountObjectIdentifier("ext_integration")}
				opts.Secrets = []SecretReference{{VariableName: "variable1", Name: secretId}, {VariableName: "variable2", Name: secretId2}}
				opts.TargetPath = new("@~/testfunc.jar")
				opts.FunctionDefinition = new(wrapFunctionDefinition("return id + name;"))
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("id" NUMBER, "name" VARCHAR DEFAULT 'test') COPY GRANTS RETURNS TABLE ("country_code" VARCHAR, "country_name" VARCHAR) NOT NULL LANGUAGE JAVA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar') PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'TestFunc.echoVarchar' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) TARGET_PATH = '@~/testfunc.jar' AS $$return id + name;$$`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateForJava_orReplace",
			func(opts *CreateForJavaFunctionOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			`CREATE OR REPLACE FUNCTION %s () COPY GRANTS RETURNS VARCHAR(100) LANGUAGE JAVA IMPORTS = ('@~/my_lib.jar') HANDLER = 'TestFunc.echoVarchar'`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	functionsTests.CreateForJavascript.
		withDefaultOpts(func() *CreateForJavascriptFunctionOptions {
			return &CreateForJavascriptFunctionOptions{
				name: functionsTestIdSchemaObjectIdentifier,
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
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_CreateForJavascript_all,
			func(opts *CreateForJavascriptFunctionOptions) {
				opts.IfNotExists = new(true)
				opts.Temporary = new(true)
				opts.Secure = new(true)
				opts.Arguments = []FunctionArgument{{ArgName: "d", ArgDataType: dataTypeFloat, DefaultValue: new("1.0")}}
				opts.Returns = FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataType: dataTypeFloat}}
				opts.ReturnNullValues = new(ReturnNullValuesNotNull)
				opts.NullInputBehavior = new(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Comment = new("comment")
			},
			`CREATE TEMPORARY SECURE FUNCTION IF NOT EXISTS %s ("d" FLOAT DEFAULT 1.0) RETURNS FLOAT NOT NULL LANGUAGE JAVASCRIPT CALLED ON NULL INPUT IMMUTABLE COMMENT = 'comment' AS $$return 1;$$`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		// TODO [SNOW-1348103]: remove with old function removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForJavascript_allOldDataTypes",
			func(opts *CreateForJavascriptFunctionOptions) {
				opts.OrReplace = new(true)
				opts.Temporary = new(true)
				opts.Secure = new(true)
				opts.Arguments = []FunctionArgument{{ArgName: "d", ArgDataTypeOld: DataTypeFloat, DefaultValue: new("1.0")}}
				opts.CopyGrants = new(true)
				opts.Returns = FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataTypeOld: DataTypeFloat}}
				opts.ReturnNullValues = new(ReturnNullValuesNotNull)
				opts.NullInputBehavior = new(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Comment = new("comment")
				opts.FunctionDefinition = wrapFunctionDefinition("return 1;")
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("d" FLOAT DEFAULT 1.0) COPY GRANTS RETURNS FLOAT NOT NULL LANGUAGE JAVASCRIPT CALLED ON NULL INPUT IMMUTABLE COMMENT = 'comment' AS $$return 1;$$`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateForJavascript_orReplace",
			func(opts *CreateForJavascriptFunctionOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			`CREATE OR REPLACE FUNCTION %s () COPY GRANTS RETURNS FLOAT LANGUAGE JAVASCRIPT AS $$return 1;$$`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	functionsTests.CreateForPython.
		withDefaultOpts(func() *CreateForPythonFunctionOptions {
			return &CreateForPythonFunctionOptions{
				name:           functionsTestIdSchemaObjectIdentifier,
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
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_CreateForPython_all,
			func(opts *CreateForPythonFunctionOptions) {
				opts.IfNotExists = new(true)
				opts.Temporary = new(true)
				opts.Secure = new(true)
				opts.Arguments = []FunctionArgument{{ArgName: "i", ArgDataType: dataTypeNumber_36_2, DefaultValue: new("1")}}
				opts.Returns = FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataType: dataTypeVariant}}
				opts.ReturnNullValues = new(ReturnNullValuesNotNull)
				opts.NullInputBehavior = new(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Comment = new("comment")
				opts.Imports = []FunctionImport{{FunctionImport: "numpy"}, {FunctionImport: "pandas"}}
				opts.Packages = []FunctionPackage{{FunctionPackage: "numpy"}, {FunctionPackage: "pandas"}}
				opts.ExternalAccessIntegrations = []AccountObjectIdentifier{NewAccountObjectIdentifier("ext_integration")}
				opts.Secrets = []SecretReference{{VariableName: "variable1", Name: secretId}, {VariableName: "variable2", Name: secretId2}}
				opts.FunctionDefinition = new(wrapFunctionDefinition("import numpy as np"))
			},
			`CREATE TEMPORARY SECURE FUNCTION IF NOT EXISTS %s ("i" NUMBER(36, 2) DEFAULT 1) RETURNS VARIANT NOT NULL LANGUAGE PYTHON CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '3.9' COMMENT = 'comment' IMPORTS = ('numpy', 'pandas') PACKAGES = ('numpy', 'pandas') HANDLER = 'udf' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) AS $$import numpy as np$$`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		).
		// TODO [SNOW-1348103]: remove with old function removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForPython_allOldDataTypes",
			func(opts *CreateForPythonFunctionOptions) {
				opts.OrReplace = new(true)
				opts.Temporary = new(true)
				opts.Secure = new(true)
				opts.Arguments = []FunctionArgument{{ArgName: "i", ArgDataTypeOld: DataTypeNumber, DefaultValue: new("1")}}
				opts.CopyGrants = new(true)
				opts.Returns = FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataTypeOld: DataTypeVariant}}
				opts.ReturnNullValues = new(ReturnNullValuesNotNull)
				opts.NullInputBehavior = new(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.RuntimeVersion = "3.9"
				opts.Comment = new("comment")
				opts.Imports = []FunctionImport{{FunctionImport: "numpy"}, {FunctionImport: "pandas"}}
				opts.Packages = []FunctionPackage{{FunctionPackage: "numpy"}, {FunctionPackage: "pandas"}}
				opts.ExternalAccessIntegrations = []AccountObjectIdentifier{NewAccountObjectIdentifier("ext_integration")}
				opts.Secrets = []SecretReference{{VariableName: "variable1", Name: secretId}, {VariableName: "variable2", Name: secretId2}}
				opts.FunctionDefinition = new(wrapFunctionDefinition("import numpy as np"))
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("i" NUMBER DEFAULT 1) COPY GRANTS RETURNS VARIANT NOT NULL LANGUAGE PYTHON CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '3.9' COMMENT = 'comment' IMPORTS = ('numpy', 'pandas') PACKAGES = ('numpy', 'pandas') HANDLER = 'udf' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) AS $$import numpy as np$$`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateForPython_orReplace",
			func(opts *CreateForPythonFunctionOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			`CREATE OR REPLACE FUNCTION %s () COPY GRANTS RETURNS VARIANT LANGUAGE PYTHON RUNTIME_VERSION = '3.9' IMPORTS = ('numpy') HANDLER = 'udf'`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	functionsTests.CreateForScala.
		withDefaultOpts(func() *CreateForScalaFunctionOptions {
			return &CreateForScalaFunctionOptions{
				name:           functionsTestIdSchemaObjectIdentifier,
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
				opts.TargetPath = new("@~/testfunc.jar")
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
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_CreateForScala_all,
			func(opts *CreateForScalaFunctionOptions) {
				opts.IfNotExists = new(true)
				opts.Temporary = new(true)
				opts.Secure = new(true)
				opts.Arguments = []FunctionArgument{{ArgName: "x", ArgDataType: dataTypeVarchar_100, DefaultValue: new("'test'")}}
				opts.ResultDataType = dataTypeVarchar_100
				opts.ReturnNullValues = new(ReturnNullValuesNotNull)
				opts.NullInputBehavior = new(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.RuntimeVersion = "2.0"
				opts.Comment = new("comment")
				opts.Imports = []FunctionImport{{FunctionImport: "@udf_libs/echohandler.jar"}}
				opts.FunctionDefinition = new(wrapFunctionDefinition("return x"))
			},
			`CREATE TEMPORARY SECURE FUNCTION IF NOT EXISTS %s ("x" VARCHAR(100) DEFAULT 'test') RETURNS VARCHAR(100) NOT NULL LANGUAGE SCALA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar' AS $$return x$$`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		// TODO [SNOW-1348103]: remove with old function removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForScala_allOldDataTypes",
			func(opts *CreateForScalaFunctionOptions) {
				opts.OrReplace = new(true)
				opts.Temporary = new(true)
				opts.Secure = new(true)
				opts.Arguments = []FunctionArgument{{ArgName: "x", ArgDataTypeOld: DataTypeVARCHAR, DefaultValue: new("'test'")}}
				opts.CopyGrants = new(true)
				opts.ResultDataType = nil // clear new-style data type set by default opts
				opts.ResultDataTypeOld = DataTypeVARCHAR
				opts.ReturnNullValues = new(ReturnNullValuesNotNull)
				opts.NullInputBehavior = new(NullInputBehaviorCalledOnNullInput)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.RuntimeVersion = "2.0"
				opts.Comment = new("comment")
				opts.Imports = []FunctionImport{{FunctionImport: "@udf_libs/echohandler.jar"}}
				opts.FunctionDefinition = new(wrapFunctionDefinition("return x"))
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("x" VARCHAR DEFAULT 'test') COPY GRANTS RETURNS VARCHAR NOT NULL LANGUAGE SCALA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar' AS $$return x$$`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateForScala_orReplace",
			func(opts *CreateForScalaFunctionOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			`CREATE OR REPLACE FUNCTION %s () COPY GRANTS RETURNS VARCHAR(100) LANGUAGE SCALA RUNTIME_VERSION = '' IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar'`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	functionsTests.CreateForSQL.
		withDefaultOpts(func() *CreateForSQLFunctionOptions {
			return &CreateForSQLFunctionOptions{
				name:               functionsTestIdSchemaObjectIdentifier,
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
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_CreateForSQL_all,
			func(opts *CreateForSQLFunctionOptions) {
				opts.IfNotExists = new(true)
				opts.Temporary = new(true)
				opts.Secure = new(true)
				opts.Arguments = []FunctionArgument{{ArgName: "message", ArgDataType: dataTypeVarchar_100, DefaultValue: new("'test'")}}
				opts.Returns = FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataType: dataTypeFloat}}
				opts.ReturnNullValues = new(ReturnNullValuesNotNull)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Memoizable = new(true)
				opts.Comment = new("comment")
			},
			`CREATE TEMPORARY SECURE FUNCTION IF NOT EXISTS %s ("message" VARCHAR(100) DEFAULT 'test') RETURNS FLOAT NOT NULL IMMUTABLE MEMOIZABLE COMMENT = 'comment' AS $$3.141592654::FLOAT$$`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		// TODO [SNOW-1348103]: remove with old function removal for V1
		withAdditionalSqlCasef(
			"sql_CreateForSQL_allOldDataTypes",
			func(opts *CreateForSQLFunctionOptions) {
				opts.OrReplace = new(true)
				opts.Temporary = new(true)
				opts.Secure = new(true)
				opts.Arguments = []FunctionArgument{{ArgName: "message", ArgDataTypeOld: "VARCHAR", DefaultValue: new("'test'")}}
				opts.CopyGrants = new(true)
				opts.Returns = FunctionReturns{ResultDataType: &FunctionReturnsResultDataType{ResultDataTypeOld: DataTypeFloat}}
				opts.ReturnNullValues = new(ReturnNullValuesNotNull)
				opts.ReturnResultsBehavior = new(ReturnResultsBehaviorImmutable)
				opts.Memoizable = new(true)
				opts.Comment = new("comment")
			},
			`CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("message" VARCHAR DEFAULT 'test') COPY GRANTS RETURNS FLOAT NOT NULL IMMUTABLE MEMOIZABLE COMMENT = 'comment' AS $$3.141592654::FLOAT$$`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateForSQL_orReplace",
			func(opts *CreateForSQLFunctionOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			`CREATE OR REPLACE FUNCTION %s () COPY GRANTS RETURNS FLOAT AS $$3.141592654::FLOAT$$`,
			functionsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	functionsTests.Alter.
		withDefaultOpts(func() *AlterFunctionOptions {
			return &AlterFunctionOptions{name: functionsTestIdSchemaObjectIdentifierWithArguments, IfExists: new(true)}
		}).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_RenameTo,
			func(opts *AlterFunctionOptions) { opts.RenameTo = &newId },
			`ALTER FUNCTION IF EXISTS %s RENAME TO %s`,
			functionsTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(), newId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_Set,
			func(opts *AlterFunctionOptions) {
				opts.Set = &FunctionSet{Comment: new("comment"), TraceLevel: new(TraceLevelOff)}
			},
			`ALTER FUNCTION IF EXISTS %s SET COMMENT = 'comment', TRACE_LEVEL = 'OFF'`,
			functionsTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_emptySecrets",
			func(opts *AlterFunctionOptions) {
				opts.Set = &FunctionSet{SecretsList: &SecretsList{}}
			},
			`ALTER FUNCTION IF EXISTS %s SET SECRETS = ()`, functionsTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_nonEmptySecrets",
			func(opts *AlterFunctionOptions) {
				opts.Set = &FunctionSet{SecretsList: &SecretsList{[]SecretReference{{VariableName: "abc", Name: secretId}}}}
			},
			`ALTER FUNCTION IF EXISTS %s SET SECRETS = ('abc' = %s)`, functionsTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(), secretId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_Unset,
			func(opts *AlterFunctionOptions) {
				opts.Unset = &FunctionUnset{Comment: new(true), TraceLevel: new(true)}
			},
			`ALTER FUNCTION IF EXISTS %s UNSET COMMENT, TRACE_LEVEL`,
			functionsTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_SetSecure,
			func(opts *AlterFunctionOptions) { opts.SetSecure = new(true) },
			`ALTER FUNCTION IF EXISTS %s SET SECURE`, functionsTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_UnsetSecure,
			func(opts *AlterFunctionOptions) { opts.UnsetSecure = new(true) },
			`ALTER FUNCTION IF EXISTS %s UNSET SECURE`, functionsTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_SetTags,
			func(opts *AlterFunctionOptions) {
				opts.SetTags = []TagAssociation{{Name: NewAccountObjectIdentifier("tag1"), Value: "value1"}}
			},
			`ALTER FUNCTION IF EXISTS %s SET TAG "tag1" = 'value1'`, functionsTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Alter_UnsetTags,
			func(opts *AlterFunctionOptions) {
				opts.UnsetTags = []ObjectIdentifier{NewAccountObjectIdentifier("tag1"), NewAccountObjectIdentifier("tag2")}
			},
			`ALTER FUNCTION IF EXISTS %s UNSET TAG "tag1", "tag2"`, functionsTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		)

	functionsTests.Drop.
		withExpectedSqlf(case_Functions_sql_Drop_basic,
			`DROP FUNCTION %s`, functionsTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Drop_all,
			func(opts *DropFunctionOptions) { opts.IfExists = new(true) },
			`DROP FUNCTION IF EXISTS %s`, functionsTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName(),
		)

	functionsTests.Show.
		withExpectedSql(case_Functions_sql_Show_basic, `SHOW USER FUNCTIONS`).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Show_all,
			func(opts *ShowFunctionOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
			},
			`SHOW USER FUNCTIONS LIKE 'pattern' IN ACCOUNT`,
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Show_Like,
			func(opts *ShowFunctionOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			`SHOW USER FUNCTIONS LIKE 'pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_Functions_sql_Show_In,
			func(opts *ShowFunctionOptions) { opts.In = &ExtendedIn{In: In{Account: new(true)}} },
			`SHOW USER FUNCTIONS IN ACCOUNT`,
		)

	functionsTests.Describe.
		withDefaultOpts(func() *DescribeFunctionOptions {
			return &DescribeFunctionOptions{name: functionsTestIdSchemaObjectIdentifierWithArguments}
		}).
		withExpectedSqlf(case_Functions_sql_Describe_basic,
			`DESCRIBE FUNCTION %s`, functionsTestIdSchemaObjectIdentifierWithArguments.FullyQualifiedName())
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

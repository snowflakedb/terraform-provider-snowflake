package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	javaProcedureSchema = setUpProcedureSchema(javaProcedureSchemaDefinition)
	javascriptProcedureSchema = setUpProcedureSchema(javascriptProcedureSchemaDefinition)
	pythonProcedureSchema = setUpProcedureSchema(pythonProcedureSchemaDefinition)
	scalaProcedureSchema = setUpProcedureSchema(scalaProcedureSchemaDefinition)
	sqlProcedureSchema = setUpProcedureSchema(sqlProcedureSchemaDefinition)
}

type procedureSchemaDef struct {
	additionalArguments            []string
	procedureDefinitionDescription string
	procedureDefinitionRequired    bool
	returnTypeLinkName             string
	returnTypeLinkUrl              string
	runtimeVersionDescription      string
	importsDescription             string
	handlerDescription             string
	targetPathDescription          string
}

func setUpProcedureSchema(definition procedureSchemaDef) map[string]*schema.Schema {
	currentSchema := make(map[string]*schema.Schema)
	for k, v := range procedureBaseSchema() {
		v := v
		if slices.Contains(definition.additionalArguments, k) || slices.Contains(commonProcedureArguments, k) {
			currentSchema[k] = &v
		}
	}
	if v, ok := currentSchema["procedure_definition"]; ok && v != nil {
		v.Description = diffSuppressStatementFieldDescription(definition.procedureDefinitionDescription)
		if definition.procedureDefinitionRequired {
			v.Required = true
		} else {
			v.Optional = true
		}
	}
	if v, ok := currentSchema["return_type"]; ok && v != nil {
		v.Description = procedureReturnsTemplate(definition.returnTypeLinkName, definition.returnTypeLinkUrl)
	}
	if v, ok := currentSchema["runtime_version"]; ok && v != nil {
		v.Description = definition.runtimeVersionDescription
	}
	if v, ok := currentSchema["imports"]; ok && v != nil {
		v.Description = definition.importsDescription
	}
	if v, ok := currentSchema["handler"]; ok && v != nil {
		v.Description = definition.handlerDescription
	}
	if v, ok := currentSchema["target_path"]; ok && v != nil {
		v.Description = definition.handlerDescription
	}
	return currentSchema
}

func procedureDefinitionTemplate(language string, linkName string, linkUrl string) string {
	return fmt.Sprintf("Defines the code executed by the stored procedure. The definition can consist of any valid code. Wrapping `$$` signs are added by the provider automatically; do not include them. The `procedure_definition` value must be %[1]s source code. For more information, see [%[2]s](%[3]s).", language, linkName, linkUrl)
}

func procedureReturnsTemplate(linkName string, linkUrl string) string {
	return fmt.Sprintf("Specifies the type of the result returned by the stored procedure. For `<result_data_type>`, use the Snowflake data type that corresponds to the type of the language that you are using (see [%s](%s)). For `RETURNS TABLE ( [ col_name col_data_type [ , ... ] ] )`, if you know the Snowflake data types of the columns in the returned table, specify the column names and types. Otherwise (e.g. if you are determining the column types during run time), you can omit the column names and types (i.e. `TABLE ()`).", linkName, linkUrl)
}

var (
	commonProcedureArguments = []string{
		"database",
		"schema",
		"name",
		"is_secure",
		"arguments",
		"return_type",
		"null_input_behavior",
		"comment",
		"execute_as",
		"procedure_definition",
		"procedure_language",
		ShowOutputAttributeName,
		ParametersAttributeName,
		FullyQualifiedNameAttributeName,
	}
	javaProcedureSchemaDefinition = procedureSchemaDef{
		additionalArguments: []string{
			"runtime_version",
			"imports",
			"snowpark_package",
			"packages",
			"handler",
			"external_access_integrations",
			"secrets",
			"target_path",
		},
		procedureDefinitionDescription: procedureDefinitionTemplate("Java", "Java (using Snowpark)", "https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-java"),
		returnTypeLinkName:             "SQL-Java Data Type Mappings",
		returnTypeLinkUrl:              "https://docs.snowflake.com/en/developer-guide/udf-stored-procedure-data-type-mapping.html#label-sql-java-data-type-mappings",
		runtimeVersionDescription:      "The language runtime version to use. Currently, the supported versions are: 11.",
		importsDescription:             "The location (stage), path, and name of the file(s) to import. You must set the IMPORTS clause to include any files that your stored procedure depends on. If you are writing an in-line stored procedure, you can omit this clause, unless your code depends on classes defined outside the stored procedure or resource files. If you are writing a stored procedure with a staged handler, you must also include a path to the JAR file containing the stored procedure’s handler code. The IMPORTS definition cannot reference variables from arguments that are passed into the stored procedure. Each file in the IMPORTS clause must have a unique name, even if the files are in different subdirectories or different stages.",
		handlerDescription:             "Use the fully qualified name of the method or function for the stored procedure. This is typically in the following form `com.my_company.my_package.MyClass.myMethod` where `com.my_company.my_package` corresponds to the package containing the object or class: `package com.my_company.my_package;`.",
		targetPathDescription:          "For stored procedures with inline handler code, specifies the location to which Snowflake should write the compiled code (JAR file) after compiling the source code specified in the `<procedure_definition>`. If this clause is omitted, Snowflake re-compiles the source code each time the code is needed. If you specify this clause uou cannot set this to an existing file. Snowflake returns an error if the TARGET_PATH points to an existing file. If you specify both the IMPORTS and TARGET_PATH clauses, the file name in the TARGET_PATH clause must be different from each file name in the IMPORTS clause, even if the files are in different subdirectories or different stages. If you no longer need to use the stored procedure (e.g. if you drop the stored procedure), you must manually remove this JAR file.",
	}
	javascriptProcedureSchemaDefinition = procedureSchemaDef{
		additionalArguments:            []string{},
		returnTypeLinkName:             "SQL and JavaScript data type mapping",
		returnTypeLinkUrl:              "https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-javascript.html#label-stored-procedure-data-type-mapping",
		procedureDefinitionDescription: procedureDefinitionTemplate("JavaScript", "JavaScript", "https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-javascript"),
		procedureDefinitionRequired:    true,
	}
	pythonProcedureSchemaDefinition = procedureSchemaDef{
		additionalArguments: []string{
			"runtime_version",
			"imports",
			"snowpark_package",
			"packages",
			"handler",
			"external_access_integrations",
			"secrets",
		},
		procedureDefinitionDescription: procedureDefinitionTemplate("Python", "Python (using Snowpark)", "https://docs.snowflake.com/en/developer-guide/stored-procedure/python/procedure-python-overview"),
		returnTypeLinkName:             "SQL-Python Data Type Mappings",
		returnTypeLinkUrl:              "https://docs.snowflake.com/en/developer-guide/udf-stored-procedure-data-type-mapping.html#label-sql-python-data-type-mappings",
		runtimeVersionDescription:      "The language runtime version to use. Currently, the supported versions are: 3.9, 3.10, and 3.11.",
		importsDescription:             "The location (stage), path, and name of the file(s) to import. You must set the IMPORTS clause to include any files that your stored procedure depends on. If you are writing an in-line stored procedure, you can omit this clause, unless your code depends on classes defined outside the stored procedure or resource files. If your stored procedure’s code will be on a stage, you must also include a path to the module file your code is in. The IMPORTS definition cannot reference variables from arguments that are passed into the stored procedure. Each file in the IMPORTS clause must have a unique name, even if the files are in different subdirectories or different stages.",
		handlerDescription:             "Use the name of the stored procedure’s function or method. This can differ depending on whether the code is in-line or referenced at a stage. When the code is in-line, you can specify just the function name. When the code is imported from a stage, specify the fully-qualified handler function name as `<module_name>.<function_name>`.",
	}
	scalaProcedureSchemaDefinition = procedureSchemaDef{
		additionalArguments: []string{
			"runtime_version",
			"imports",
			"snowpark_package",
			"packages",
			"handler",
			"external_access_integrations",
			"secrets",
			"target_path",
		},
		procedureDefinitionDescription: procedureDefinitionTemplate("Scala", "Scala (using Snowpark)", "https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-scala"),
		returnTypeLinkName:             "SQL-Scala Data Type Mappings",
		returnTypeLinkUrl:              "https://docs.snowflake.com/en/developer-guide/udf-stored-procedure-data-type-mapping.html#label-sql-types-to-scala-types",
		runtimeVersionDescription:      "The language runtime version to use. Currently, the supported versions are: 2.12.",
		importsDescription:             "The location (stage), path, and name of the file(s) to import. You must set the IMPORTS clause to include any files that your stored procedure depends on. If you are writing an in-line stored procedure, you can omit this clause, unless your code depends on classes defined outside the stored procedure or resource files. If you are writing a stored procedure with a staged handler, you must also include a path to the JAR file containing the stored procedure’s handler code. The IMPORTS definition cannot reference variables from arguments that are passed into the stored procedure. Each file in the IMPORTS clause must have a unique name, even if the files are in different subdirectories or different stages.",
		handlerDescription:             "Use the fully qualified name of the method or function for the stored procedure. This is typically in the following form: `com.my_company.my_package.MyClass.myMethod` where `com.my_company.my_package` corresponds to the package containing the object or class: `package com.my_company.my_package;`.",
		targetPathDescription:          "For stored procedures with inline handler code, specifies the location to which Snowflake should write the compiled code (JAR file) after compiling the source code specified in the procedure_definition. If this clause is omitted, Snowflake re-compiles the source code each time the code is needed. If you specify this clause you cannot set this to an existing file. Snowflake returns an error if the TARGET_PATH points to an existing file. If you specify both the IMPORTS and TARGET_PATH clauses, the file name in the TARGET_PATH clause must be different from each file name in the IMPORTS clause, even if the files are in different subdirectories or different stages. If you no longer need to use the stored procedure (e.g. if you drop the stored procedure), you must manually remove this JAR file.",
	}
	sqlProcedureSchemaDefinition = procedureSchemaDef{
		additionalArguments:            []string{},
		procedureDefinitionDescription: procedureDefinitionTemplate("SQL", "Snowflake Scripting", "https://docs.snowflake.com/en/developer-guide/snowflake-scripting/index"),
		procedureDefinitionRequired:    true,
		returnTypeLinkName:             "SQL data type",
		returnTypeLinkUrl:              "https://docs.snowflake.com/en/sql-reference-data-types",
	}
)

var (
	javaProcedureSchema       map[string]*schema.Schema
	javascriptProcedureSchema map[string]*schema.Schema
	pythonProcedureSchema     map[string]*schema.Schema
	scalaProcedureSchema      map[string]*schema.Schema
	sqlProcedureSchema        map[string]*schema.Schema
)

// TODO [SNOW-1348103]: add null/not null
// TODO [SNOW-1348103]: currently database and schema are ForceNew but based on the docs it is possible to rename with moving to different db/schema
// TODO [SNOW-1348103]: copyGrants and orReplace logic omitted for now, will be added to the limitations docs
func procedureBaseSchema() map[string]schema.Schema {
	return map[string]schema.Schema{
		"database": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppressIdentifierQuoting,
			Description:      blocklistedCharactersFieldDescription("The database in which to create the procedure."),
		},
		"schema": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppressIdentifierQuoting,
			Description:      blocklistedCharactersFieldDescription("The schema in which to create the procedure."),
		},
		"name": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      blocklistedCharactersFieldDescription("The name of the procedure; the identifier does not need to be unique for the schema in which the procedure is created because stored procedures are [identified and resolved by the combination of the name and argument types](https://docs.snowflake.com/en/developer-guide/udf-stored-procedure-naming-conventions.html#label-procedure-function-name-overloading)."),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"is_secure": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("is_secure"),
			Description:      booleanStringFieldDescription("Specifies that the procedure is secure. For more information about secure procedures, see [Protecting Sensitive Information with Secure UDFs and Stored Procedures](https://docs.snowflake.com/en/developer-guide/secure-udf-procedure)."),
		},
		"arguments": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"arg_name": {
						Type:     schema.TypeString,
						Required: true,
						// TODO [SNOW-1348103]: adjust diff suppression accordingly.
						Description: "The argument name. The provider wraps it in double quotes by default, so be aware of that while referencing the argument in the procedure definition.",
					},
					"arg_data_type": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: IsDataTypeValid,
						DiffSuppressFunc: DiffSuppressDataTypes,
						Description:      "The argument type.",
					},
					"arg_default_value": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: externalChangesNotDetectedFieldDescription("Optional default value for the argument. For text values use single quotes. Numeric values can be unquoted."),
					},
				},
			},
			Optional:    true,
			ForceNew:    true,
			Description: "List of the arguments for the procedure. Consult the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-procedure#all-languages) for more details.",
		},
		"return_type": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: IsDataTypeValid,
			DiffSuppressFunc: DiffSuppressDataTypes,
		},
		"null_input_behavior": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: sdkValidation(sdk.ToNullInputBehavior),
			DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToNullInputBehavior)), // IgnoreChangeToCurrentSnowflakeValueInShow("null_input_behavior")),
			Description:      fmt.Sprintf("Specifies the behavior of the procedure when called with null inputs. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllAllowedNullInputBehaviors)),
		},
		// "return_behavior" removed because it is deprecated in the docs: https://docs.snowflake.com/en/sql-reference/sql/create-procedure#id1
		"runtime_version": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"comment": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "user-defined procedure",
			Description: "Specifies a comment for the procedure.",
		},
		"imports": {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"stage_location": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Stage location without leading `@`. To use your user's stage set this to `~`, otherwise pass fully qualified name of the stage (with every part contained in double quotes or use `snowflake_stage.<your stage's resource name>.fully_qualified_name` if you manage this stage through terraform).",
					},
					"path_on_stage": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Path for import on stage, without the leading `/`.",
					},
				},
			},
		},
		"snowpark_package": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The Snowpark package is required for stored procedures, so it must always be present. For more information about Snowpark, see [Snowpark API](https://docs.snowflake.com/en/developer-guide/snowpark/index).",
		},
		// TODO [SNOW-1850370]: what do we do with the version "latest".
		"packages": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			ForceNew:    true,
			Description: "List of the names of packages deployed in Snowflake that should be included in the handler code’s execution environment. The Snowpark package is required for stored procedures, but is specified in the `snowpark_package` attribute. For more information about Snowpark, see [Snowpark API](https://docs.snowflake.com/en/developer-guide/snowpark/index).",
		},
		"handler": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"external_access_integrations": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			},
			Optional:    true,
			Description: "The names of [external access integrations](https://docs.snowflake.com/en/sql-reference/sql/create-external-access-integration) needed in order for this procedure’s handler code to access external networks. An external access integration specifies [network rules](https://docs.snowflake.com/en/sql-reference/sql/create-network-rule) and [secrets](https://docs.snowflake.com/en/sql-reference/sql/create-secret) that specify external locations and credentials (if any) allowed for use by handler code when making requests of an external network, such as an external REST API.",
		},
		"secrets": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"secret_variable_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The variable that will be used in handler code when retrieving information from the secret.",
					},
					"secret_id": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Fully qualified name of the allowed [secret](https://docs.snowflake.com/en/sql-reference/sql/create-secret). You will receive an error if you specify a SECRETS value whose secret isn’t also included in an integration specified by the EXTERNAL_ACCESS_INTEGRATIONS parameter.",
						DiffSuppressFunc: suppressIdentifierQuoting,
					},
				},
			},
			Description: "Assigns the names of [secrets](https://docs.snowflake.com/en/sql-reference/sql/create-secret) to variables so that you can use the variables to reference the secrets when retrieving information from secrets in handler code. Secrets you specify here must be allowed by the [external access integration](https://docs.snowflake.com/en/sql-reference/sql/create-external-access-integration) specified as a value of this CREATE FUNCTION command’s EXTERNAL_ACCESS_INTEGRATIONS parameter.",
		},
		"target_path": {
			Type:     schema.TypeSet,
			MaxItems: 1,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"stage_location": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Stage location without leading `@`. To use your user's stage set this to `~`, otherwise pass fully qualified name of the stage (with every part contained in double quotes or use `snowflake_stage.<your stage's resource name>.fully_qualified_name` if you manage this stage through terraform).",
					},
					"path_on_stage": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Path for import on stage, without the leading `/`.",
					},
				},
			},
		},
		"execute_as": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: sdkValidation(sdk.ToExecuteAs),
			DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToExecuteAs), IgnoreChangeToCurrentSnowflakeValueInShow("execute_as")),
			Description:      fmt.Sprintf("Specifies whether the stored procedure executes with the privileges of the owner (an “owner’s rights” stored procedure) or with the privileges of the caller (a “caller’s rights” stored procedure). If you execute the statement CREATE PROCEDURE … EXECUTE AS CALLER, then in the future the procedure will execute as a caller’s rights procedure. If you execute CREATE PROCEDURE … EXECUTE AS OWNER, then the procedure will execute as an owner’s rights procedure. For more information, see [Understanding caller’s rights and owner’s rights stored procedures](https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-rights). Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllAllowedExecuteAs)),
		},
		"procedure_definition": {
			Type:             schema.TypeString,
			ForceNew:         true,
			DiffSuppressFunc: DiffSuppressStatement,
		},
		"procedure_language": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Specifies language for the procedure. Used to detect external changes.",
		},
		ShowOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW PROCEDURE` for the given procedure.",
			Elem: &schema.Resource{
				Schema: schemas.ShowProcedureSchema,
			},
		},
		ParametersAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW PARAMETERS IN PROCEDURE` for the given procedure.",
			Elem: &schema.Resource{
				Schema: schemas.ShowProcedureParametersSchema,
			},
		},
		FullyQualifiedNameAttributeName: *schemas.FullyQualifiedNameSchema,
	}
}

var DeleteProcedure = ResourceDeleteContextFunc(
	sdk.ParseSchemaObjectIdentifierWithArguments,
	func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifierWithArguments] {
		return client.Procedures.DropSafely
	},
)

func UpdateProcedure(language string, readFunc func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics) func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifierWithArguments(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		if d.HasChange("name") {
			newId := sdk.NewSchemaObjectIdentifierWithArgumentsInSchema(id.SchemaId(), d.Get("name").(string), id.ArgumentDataTypes()...)

			err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithRenameTo(newId.SchemaObjectId()))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error renaming procedure %v err = %w", d.Id(), err))
			}

			d.SetId(helpers.EncodeResourceIdentifier(newId))
			id = newId
		}

		// Batch SET operations and UNSET operations
		setRequest := sdk.NewProcedureSetRequest()
		unsetRequest := sdk.NewProcedureUnsetRequest()

		_ = stringAttributeUpdate(d, "comment", &setRequest.Comment, &unsetRequest.Comment)

		switch language {
		case "JAVA", "SCALA", "PYTHON":
			err = errors.Join(
				func() error {
					if d.HasChange("secrets") {
						return setSecretsInBuilder(d, func(references []sdk.SecretReference) *sdk.ProcedureSetRequest {
							return setRequest.WithSecretsList(sdk.SecretsListRequest{SecretsList: references})
						})
					}
					return nil
				}(),
				func() error {
					if d.HasChange("external_access_integrations") {
						return setExternalAccessIntegrationsInBuilder(d, func(references []sdk.AccountObjectIdentifier) any {
							if len(references) == 0 {
								return unsetRequest.WithExternalAccessIntegrations(true)
							} else {
								return setRequest.WithExternalAccessIntegrations(references)
							}
						})
					}
					return nil
				}(),
			)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if updateParamDiags := handleProcedureParametersUpdate(d, setRequest, unsetRequest); len(updateParamDiags) > 0 {
			return updateParamDiags
		}

		// Apply SET and UNSET changes
		if !reflect.DeepEqual(*setRequest, *sdk.NewProcedureSetRequest()) {
			err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithSet(*setRequest))
			if err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}
		}
		if !reflect.DeepEqual(*unsetRequest, *sdk.NewProcedureUnsetRequest()) {
			err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithUnset(*unsetRequest))
			if err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}
		}

		// has to be handled separately
		if d.HasChange("execute_as") {
			var value sdk.ExecuteAs
			if v, ok := d.GetOk("execute_as"); ok {
				value, err = sdk.ToExecuteAs(v.(string))
				if err != nil {
					return diag.FromErr(err)
				}
			} else {
				// there is no UNSET, so we need to set it manually
				value = sdk.ExecuteAsOwner
			}
			err = client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithExecuteAs(value))
			if err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}
		}

		return readFunc(ctx, d, meta)
	}
}

func ImportProcedure(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifierWithArguments(d.Id())
	if err != nil {
		return nil, err
	}

	procedureDetails, err := client.Procedures.DescribeDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	procedure, err := client.Procedures.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = errors.Join(
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
		d.Set("name", id.Name()),
		d.Set("is_secure", booleanStringFromBool(procedure.IsSecure)),
		setOptionalFromStringPtr(d, "null_input_behavior", procedureDetails.NullHandling),
		d.Set("execute_as", procedureDetails.ExecuteAs),
		importFunctionOrProcedureArguments(d, procedureDetails.NormalizedArguments),
		// all others are set in read
	)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func queryAllProcedureDetailsCommon(ctx context.Context, d *schema.ResourceData, client *sdk.Client, id sdk.SchemaObjectIdentifierWithArguments) (*allProcedureDetailsCommon, diag.Diagnostics) {
	procedure, err := client.Procedures.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query procedure. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Procedure: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}

	procedureDetails, err := client.Procedures.DescribeDetails(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			log.Printf("[DEBUG] procedure (%s) not found or we are not authorized. Err: %s", d.Id(), err)
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query procedure. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Procedure: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}

	procedureParameters, err := client.Procedures.ShowParameters(ctx, id)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return &allProcedureDetailsCommon{
		procedure:           procedure,
		procedureDetails:    procedureDetails,
		procedureParameters: procedureParameters,
	}, nil
}

type allProcedureDetailsCommon struct {
	procedure           *sdk.Procedure
	procedureDetails    *sdk.ProcedureDetails
	procedureParameters []*sdk.Parameter
}

// TODO [SNOW-1850370]: Make the rest of the functions in this file generic (for reuse with functions)
// These were copy-pasted for now.
func parseProcedureArgumentsCommon(d *schema.ResourceData) ([]sdk.ProcedureArgumentRequest, error) {
	args := make([]sdk.ProcedureArgumentRequest, 0)
	if v, ok := d.GetOk("arguments"); ok {
		for _, arg := range v.([]any) {
			argName := arg.(map[string]any)["arg_name"].(string)
			argDataType := arg.(map[string]any)["arg_data_type"].(string)
			dataType, err := datatypes.ParseDataType(argDataType)
			if err != nil {
				return nil, err
			}
			request := sdk.NewProcedureArgumentRequest(argName, dataType)

			if argDefaultValue, defaultValuePresent := arg.(map[string]any)["arg_default_value"]; defaultValuePresent && argDefaultValue.(string) != "" {
				request.WithDefaultValue(argDefaultValue.(string))
			}

			args = append(args, *request)
		}
	}
	return args, nil
}

func parseProcedureImportsCommon(d *schema.ResourceData) ([]sdk.ProcedureImportRequest, error) {
	imports := make([]sdk.ProcedureImportRequest, 0)
	if v, ok := d.GetOk("imports"); ok {
		for _, imp := range v.(*schema.Set).List() {
			stageLocation := imp.(map[string]any)["stage_location"].(string)
			pathOnStage := imp.(map[string]any)["path_on_stage"].(string)
			imports = append(imports, *sdk.NewProcedureImportRequest(fmt.Sprintf("@%s/%s", stageLocation, pathOnStage)))
		}
	}
	return imports, nil
}

func parseProceduresPackagesCommon(d *schema.ResourceData) ([]sdk.ProcedurePackageRequest, error) {
	packages := make([]sdk.ProcedurePackageRequest, 0)
	if v, ok := d.GetOk("packages"); ok {
		for _, pkg := range v.(*schema.Set).List() {
			packages = append(packages, *sdk.NewProcedurePackageRequest(pkg.(string)))
		}
	}
	return packages, nil
}

func parseProcedureTargetPathCommon(d *schema.ResourceData) (string, error) {
	var tp string
	if v, ok := d.GetOk("target_path"); ok {
		for _, p := range v.(*schema.Set).List() {
			stageLocation := p.(map[string]any)["stage_location"].(string)
			pathOnStage := p.(map[string]any)["path_on_stage"].(string)
			tp = fmt.Sprintf("@%s/%s", stageLocation, pathOnStage)
		}
	}
	return tp, nil
}

func parseProcedureReturnsCommon(d *schema.ResourceData) (*sdk.ProcedureReturnsRequest, error) {
	returnTypeRaw := d.Get("return_type").(string)
	dataType, err := datatypes.ParseDataType(returnTypeRaw)
	if err != nil {
		return nil, err
	}
	returns := sdk.NewProcedureReturnsRequest()
	switch v := dataType.(type) {
	case *datatypes.TableDataType:
		var cr []sdk.ProcedureColumnRequest
		for _, c := range v.Columns() {
			cr = append(cr, *sdk.NewProcedureColumnRequest(c.ColumnName(), c.ColumnType()))
		}
		returns.WithTable(*sdk.NewProcedureReturnsTableRequest().WithColumns(cr))
	default:
		returns.WithResultDataType(*sdk.NewProcedureReturnsResultDataTypeRequest(dataType))
	}
	return returns, nil
}

func parseProcedureSqlReturns(d *schema.ResourceData) (*sdk.ProcedureSQLReturnsRequest, error) {
	returnTypeRaw := d.Get("return_type").(string)
	dataType, err := datatypes.ParseDataType(returnTypeRaw)
	if err != nil {
		return nil, err
	}
	returns := sdk.NewProcedureSQLReturnsRequest()
	switch v := dataType.(type) {
	case *datatypes.TableDataType:
		var cr []sdk.ProcedureColumnRequest
		for _, c := range v.Columns() {
			cr = append(cr, *sdk.NewProcedureColumnRequest(c.ColumnName(), c.ColumnType()))
		}
		returns.WithTable(*sdk.NewProcedureReturnsTableRequest().WithColumns(cr))
	default:
		returns.WithResultDataType(*sdk.NewProcedureReturnsResultDataTypeRequest(dataType))
	}
	return returns, nil
}

func setProcedureImportsInBuilder[T any](d *schema.ResourceData, setImports func([]sdk.ProcedureImportRequest) T) error {
	imports, err := parseProcedureImportsCommon(d)
	if err != nil {
		return err
	}
	setImports(imports)
	return nil
}

func setProcedureTargetPathInBuilder[T any](d *schema.ResourceData, setTargetPath func(string) T) error {
	tp, err := parseProcedureTargetPathCommon(d)
	if err != nil {
		return err
	}
	if tp != "" {
		setTargetPath(tp)
	}
	return nil
}

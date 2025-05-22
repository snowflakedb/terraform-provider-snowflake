package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	javaFunctionSchema = setUpFunctionSchema(javaFunctionSchemaDefinition)
	javascriptFunctionSchema = setUpFunctionSchema(javascriptFunctionSchemaDefinition)
	pythonFunctionSchema = setUpFunctionSchema(pythonFunctionSchemaDefinition)
	scalaFunctionSchema = setUpFunctionSchema(scalaFunctionSchemaDefinition)
	sqlFunctionSchema = setUpFunctionSchema(sqlFunctionSchemaDefinition)
}

type functionSchemaDef struct {
	additionalArguments           []string
	functionDefinitionDescription string
	functionDefinitionRequired    bool
	runtimeVersionRequired        bool
	runtimeVersionDescription     string
	importsDescription            string
	packagesDescription           string
	handlerDescription            string
	targetPathDescription         string
}

func setUpFunctionSchema(definition functionSchemaDef) map[string]*schema.Schema {
	currentSchema := make(map[string]*schema.Schema)
	for k, v := range functionBaseSchema() {
		v := v
		if slices.Contains(definition.additionalArguments, k) || slices.Contains(commonFunctionArguments, k) {
			currentSchema[k] = &v
		}
	}
	if v, ok := currentSchema["function_definition"]; ok && v != nil {
		v.Description = diffSuppressStatementFieldDescription(definition.functionDefinitionDescription)
		if definition.functionDefinitionRequired {
			v.Required = true
		} else {
			v.Optional = true
		}
	}
	if v, ok := currentSchema["runtime_version"]; ok && v != nil {
		if definition.runtimeVersionRequired {
			v.Required = true
		} else {
			v.Optional = true
		}
		v.Description = definition.runtimeVersionDescription
	}
	if v, ok := currentSchema["imports"]; ok && v != nil {
		v.Description = definition.importsDescription
	}
	if v, ok := currentSchema["packages"]; ok && v != nil {
		v.Description = definition.packagesDescription
	}
	if v, ok := currentSchema["handler"]; ok && v != nil {
		v.Description = definition.handlerDescription
	}
	if v, ok := currentSchema["target_path"]; ok && v != nil {
		v.Description = definition.handlerDescription
	}
	return currentSchema
}

func functionDefinitionTemplate(language string, linkUrl string) string {
	return fmt.Sprintf("Defines the handler code executed when the UDF is called. Wrapping `$$` signs are added by the provider automatically; do not include them. The `function_definition` value must be %[1]s source code. For more information, see [Introduction to %[1]s UDFs](%[2]s).", language, linkUrl)
}

var (
	commonFunctionArguments = []string{
		"database",
		"schema",
		"name",
		"is_secure",
		"arguments",
		"return_type",
		"return_results_behavior",
		"comment",
		"function_definition",
		"function_language",
		ShowOutputAttributeName,
		ParametersAttributeName,
		FullyQualifiedNameAttributeName,
	}
	javaFunctionSchemaDefinition = functionSchemaDef{
		additionalArguments: []string{
			"runtime_version",
			"null_input_behavior",
			"imports",
			"packages",
			"handler",
			"external_access_integrations",
			"secrets",
			"target_path",
		},
		functionDefinitionDescription: functionDefinitionTemplate("Java", "https://docs.snowflake.com/en/developer-guide/udf/java/udf-java-introduction"),
		// May be optional for java because if it is not set, describe return empty version.
		runtimeVersionRequired:    false,
		runtimeVersionDescription: "Specifies the Java JDK runtime version to use. The supported versions of Java are 11.x and 17.x. If RUNTIME_VERSION is not set, Java JDK 11 is used.",
		importsDescription:        "The location (stage), path, and name of the file(s) to import. A file can be a JAR file or another type of file. If the file is a JAR file, it can contain one or more .class files and zero or more resource files. JNI (Java Native Interface) is not supported. Snowflake prohibits loading libraries that contain native code (as opposed to Java bytecode). Java UDFs can also read non-JAR files. For an example, see [Reading a file specified statically in IMPORTS](https://docs.snowflake.com/en/developer-guide/udf/java/udf-java-cookbook.html#label-reading-file-from-java-udf-imports). Consult the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#java).",
		packagesDescription:       "The name and version number of Snowflake system packages required as dependencies. The value should be of the form `package_name:version_number`, where `package_name` is `snowflake_domain:package`.",
		handlerDescription:        "The name of the handler method or class. If the handler is for a scalar UDF, returning a non-tabular value, the HANDLER value should be a method name, as in the following form: `MyClass.myMethod`. If the handler is for a tabular UDF, the HANDLER value should be the name of a handler class.",
		targetPathDescription:     "The TARGET_PATH clause specifies the location to which Snowflake should write the compiled code (JAR file) after compiling the source code specified in the `function_definition`. If this clause is included, the user should manually remove the JAR file when it is no longer needed (typically when the Java UDF is dropped). If this clause is omitted, Snowflake re-compiles the source code each time the code is needed. The JAR file is not stored permanently, and the user does not need to clean up the JAR file. Snowflake returns an error if the TARGET_PATH matches an existing file; you cannot use TARGET_PATH to overwrite an existing file.",
	}
	javascriptFunctionSchemaDefinition = functionSchemaDef{
		additionalArguments: []string{
			"null_input_behavior",
		},
		functionDefinitionDescription: functionDefinitionTemplate("JavaScript", "https://docs.snowflake.com/en/developer-guide/udf/javascript/udf-javascript-introduction"),
		functionDefinitionRequired:    true,
	}
	pythonFunctionSchemaDefinition = functionSchemaDef{
		additionalArguments: []string{
			"is_aggregate",
			"runtime_version",
			"null_input_behavior",
			"imports",
			"packages",
			"handler",
			"external_access_integrations",
			"secrets",
		},
		functionDefinitionDescription: functionDefinitionTemplate("Python", "https://docs.snowflake.com/en/developer-guide/udf/python/udf-python-introduction"),
		runtimeVersionRequired:        true,
		runtimeVersionDescription:     "Specifies the Python version to use. The supported versions of Python are: 3.9, 3.10, and 3.11.",
		importsDescription:            "The location (stage), path, and name of the file(s) to import. A file can be a `.py` file or another type of file. Python UDFs can also read non-Python files, such as text files. For an example, see [Reading a file](https://docs.snowflake.com/en/developer-guide/udf/python/udf-python-examples.html#label-udf-python-read-files). Consult the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#python).",
		packagesDescription:           "The name and version number of packages required as dependencies. The value should be of the form `package_name==version_number`.",
		handlerDescription:            "The name of the handler function or class. If the handler is for a scalar UDF, returning a non-tabular value, the HANDLER value should be a function name. If the handler code is in-line with the CREATE FUNCTION statement, you can use the function name alone. When the handler code is referenced at a stage, this value should be qualified with the module name, as in the following form: `my_module.my_function`. If the handler is for a tabular UDF, the HANDLER value should be the name of a handler class.",
	}
	scalaFunctionSchemaDefinition = functionSchemaDef{
		additionalArguments: []string{
			"runtime_version",
			"null_input_behavior",
			"imports",
			"packages",
			"handler",
			"external_access_integrations",
			"secrets",
			"target_path",
		},
		functionDefinitionDescription: functionDefinitionTemplate("Scala", "https://docs.snowflake.com/en/developer-guide/udf/scala/udf-scala-introduction"),
		runtimeVersionRequired:        true,
		runtimeVersionDescription:     "Specifies the Scala runtime version to use. The supported versions of Scala are: 2.12.",
		importsDescription:            "The location (stage), path, and name of the file(s) to import, such as a JAR or other kind of file. The JAR file might contain handler dependency libraries. It can contain one or more .class files and zero or more resource files. JNI (Java Native Interface) is not supported. Snowflake prohibits loading libraries that contain native code (as opposed to Java bytecode). A non-JAR file might a file read by handler code. For an example, see [Reading a file specified statically in IMPORTS](https://docs.snowflake.com/en/developer-guide/udf/java/udf-java-cookbook.html#label-reading-file-from-java-udf-imports). Consult the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#scala).",
		packagesDescription:           "The name and version number of Snowflake system packages required as dependencies. The value should be of the form `package_name:version_number`, where `package_name` is `snowflake_domain:package`.",
		handlerDescription:            "The name of the handler method or class. If the handler is for a scalar UDF, returning a non-tabular value, the HANDLER value should be a method name, as in the following form: `MyClass.myMethod`.",
		targetPathDescription:         "The TARGET_PATH clause specifies the location to which Snowflake should write the compiled code (JAR file) after compiling the source code specified in the `function_definition`. If this clause is included, you should manually remove the JAR file when it is no longer needed (typically when the UDF is dropped). If this clause is omitted, Snowflake re-compiles the source code each time the code is needed. The JAR file is not stored permanently, and you do not need to clean up the JAR file. Snowflake returns an error if the TARGET_PATH matches an existing file; you cannot use TARGET_PATH to overwrite an existing file.",
	}
	sqlFunctionSchemaDefinition = functionSchemaDef{
		additionalArguments:           []string{},
		functionDefinitionDescription: functionDefinitionTemplate("SQL", "https://docs.snowflake.com/en/developer-guide/udf/sql/udf-sql-introduction"),
		functionDefinitionRequired:    true,
	}
)

var (
	javaFunctionSchema       map[string]*schema.Schema
	javascriptFunctionSchema map[string]*schema.Schema
	pythonFunctionSchema     map[string]*schema.Schema
	scalaFunctionSchema      map[string]*schema.Schema
	sqlFunctionSchema        map[string]*schema.Schema
)

// TODO [SNOW-1348103]: add null/not null
// TODO [SNOW-1348103]: currently database and schema are ForceNew but based on the docs it is possible to rename with moving to different db/schema
// TODO [SNOW-1348103]: copyGrants and orReplace logic omitted for now, will be added to the limitations docs
// TODO [SNOW-1348103]: temporary is not supported because it creates a per-session object; add to limitations/design decisions
func functionBaseSchema() map[string]schema.Schema {
	return map[string]schema.Schema{
		"database": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppressIdentifierQuoting,
			Description:      blocklistedCharactersFieldDescription("The database in which to create the function."),
		},
		"schema": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppressIdentifierQuoting,
			Description:      blocklistedCharactersFieldDescription("The schema in which to create the function."),
		},
		"name": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      blocklistedCharactersFieldDescription("The name of the function; the identifier does not need to be unique for the schema in which the function is created because UDFs are identified and resolved by the combination of the name and argument types. Check the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#all-languages)."),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"is_secure": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("is_secure"),
			Description:      booleanStringFieldDescription("Specifies that the function is secure. By design, the Snowflake's `SHOW FUNCTIONS` command does not provide information about secure functions (consult [function docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#id1) and [Protecting Sensitive Information with Secure UDFs and Stored Procedures](https://docs.snowflake.com/en/developer-guide/secure-udf-procedure)) which is essential to manage/import function with Terraform. Use the role owning the function while managing secure functions."),
		},
		"is_aggregate": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("is_aggregate"),
			Description:      booleanStringFieldDescription("Specifies that the function is an aggregate function. For more information about user-defined aggregate functions, see [Python user-defined aggregate functions](https://docs.snowflake.com/en/developer-guide/udf/python/udf-python-aggregate-functions)."),
		},
		"arguments": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"arg_name": {
						Type:     schema.TypeString,
						Required: true,
						// TODO [SNOW-1348103]: adjust diff suppression accordingly.
						Description: "The argument name. The provider wraps it in double quotes by default, so be aware of that while referencing the argument in the function definition.",
					},
					// TODO [SNOW-1348103]: after testing weird names add limitations to the docs and add validation here
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
			Description: "List of the arguments for the function. Consult the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#all-languages) for more details.",
		},
		"return_type": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: IsDataTypeValid,
			DiffSuppressFunc: DiffSuppressDataTypes,
			Description:      "Specifies the results returned by the UDF, which determines the UDF type. Use `<result_data_type>` to create a scalar UDF that returns a single value with the specified data type. Use `TABLE (col_name col_data_type, ...)` to creates a table UDF that returns tabular results with the specified table column(s) and column type(s). For the details, consult the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#all-languages).",
		},
		"null_input_behavior": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: sdkValidation(sdk.ToNullInputBehavior),
			DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToNullInputBehavior)), // TODO [SNOW-1348103]: IgnoreChangeToCurrentSnowflakeValueInShow("null_input_behavior") but not in show
			Description:      fmt.Sprintf("Specifies the behavior of the function when called with null inputs. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllAllowedNullInputBehaviors)),
		},
		"return_results_behavior": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: sdkValidation(sdk.ToReturnResultsBehavior),
			DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToReturnResultsBehavior)), // TODO [SNOW-1348103]: IgnoreChangeToCurrentSnowflakeValueInShow("return_results_behavior") but not in show
			Description:      fmt.Sprintf("Specifies the behavior of the function when returning results. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllAllowedReturnResultsBehaviors)),
		},
		"runtime_version": {
			Type:     schema.TypeString,
			ForceNew: true,
		},
		"comment": {
			Type:     schema.TypeString,
			Optional: true,
			// TODO [SNOW-1348103]: handle dynamic comment - this is a workaround for now
			Default:     "user-defined function",
			Description: "Specifies a comment for the function.",
		},
		// split into two because of https://docs.snowflake.com/en/sql-reference/sql/create-function#id6
		// TODO [SNOW-1348103]: add validations preventing setting improper stage and path
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
		// TODO [SNOW-1348103]: what do we do with the version "latest".
		"packages": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			ForceNew: true,
		},
		"handler": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		// TODO [SNOW-1348103]: use suppress from network policies when adding logic
		"external_access_integrations": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			},
			Optional:    true,
			Description: "The names of [external access integrations](https://docs.snowflake.com/en/sql-reference/sql/create-external-access-integration) needed in order for this function’s handler code to access external networks. An external access integration specifies [network rules](https://docs.snowflake.com/en/sql-reference/sql/create-network-rule) and [secrets](https://docs.snowflake.com/en/sql-reference/sql/create-secret) that specify external locations and credentials (if any) allowed for use by handler code when making requests of an external network, such as an external REST API.",
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
		"function_definition": {
			Type:             schema.TypeString,
			ForceNew:         true,
			DiffSuppressFunc: DiffSuppressStatement,
		},
		"function_language": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Specifies language for the user. Used to detect external changes.",
		},
		ShowOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW FUNCTION` for the given function.",
			Elem: &schema.Resource{
				Schema: schemas.ShowFunctionSchema,
			},
		},
		ParametersAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW PARAMETERS IN FUNCTION` for the given function.",
			Elem: &schema.Resource{
				Schema: schemas.ShowFunctionParametersSchema,
			},
		},
		FullyQualifiedNameAttributeName: *schemas.FullyQualifiedNameSchema,
	}
}

var DeleteFunction = ResourceDeleteContextFunc(
	sdk.ParseSchemaObjectIdentifierWithArguments,
	func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifierWithArguments] {
		return client.Functions.DropSafely
	},
)

func UpdateFunction(language string, readFunc func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics) func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifierWithArguments(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		if d.HasChange("name") {
			newId := sdk.NewSchemaObjectIdentifierWithArgumentsInSchema(id.SchemaId(), d.Get("name").(string), id.ArgumentDataTypes()...)

			err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithRenameTo(newId.SchemaObjectId()))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error renaming function %v err = %w", d.Id(), err))
			}

			d.SetId(helpers.EncodeResourceIdentifier(newId))
			id = newId
		}

		// Batch SET operations and UNSET operations
		setRequest := sdk.NewFunctionSetRequest()
		unsetRequest := sdk.NewFunctionUnsetRequest()

		_ = stringAttributeUpdate(d, "comment", &setRequest.Comment, &unsetRequest.Comment)

		switch language {
		case "JAVA", "SCALA", "PYTHON":
			err = errors.Join(
				func() error {
					if d.HasChange("secrets") {
						return setSecretsInBuilder(d, func(references []sdk.SecretReference) *sdk.FunctionSetRequest {
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

		if updateParamDiags := handleFunctionParametersUpdate(d, setRequest, unsetRequest); len(updateParamDiags) > 0 {
			return updateParamDiags
		}

		// Apply SET and UNSET changes
		if !reflect.DeepEqual(*setRequest, *sdk.NewFunctionSetRequest()) {
			err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSet(*setRequest))
			if err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}
		}
		if !reflect.DeepEqual(*unsetRequest, *sdk.NewFunctionUnsetRequest()) {
			err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnset(*unsetRequest))
			if err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}
		}

		// has to be handled separately
		if d.HasChange("is_secure") {
			if v := d.Get("is_secure").(string); v != BooleanDefault {
				parsed, err := booleanStringToBool(v)
				if err != nil {
					return diag.FromErr(err)
				}
				err = client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetSecure(parsed))
				if err != nil {
					d.Partial(true)
					return diag.FromErr(err)
				}
			} else {
				err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnsetSecure(true))
				if err != nil {
					d.Partial(true)
					return diag.FromErr(err)
				}
			}
		}

		return readFunc(ctx, d, meta)
	}
}

func ImportFunction(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifierWithArguments(d.Id())
	if err != nil {
		return nil, err
	}

	functionDetails, err := client.Functions.DescribeDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	function, err := client.Functions.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = errors.Join(
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
		d.Set("name", id.Name()),
		d.Set("is_secure", booleanStringFromBool(function.IsSecure)),
		setOptionalFromStringPtr(d, "null_input_behavior", functionDetails.NullHandling),
		setOptionalFromStringPtr(d, "return_results_behavior", functionDetails.Volatility),
		importFunctionOrProcedureArguments(d, functionDetails.NormalizedArguments),
		// all others are set in read
	)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

// TODO [SNOW-1850370]: Make the rest of the functions in this file generic (for reuse with procedures)
func parseFunctionArgumentsCommon(d *schema.ResourceData) ([]sdk.FunctionArgumentRequest, error) {
	args := make([]sdk.FunctionArgumentRequest, 0)
	if v, ok := d.GetOk("arguments"); ok {
		for _, arg := range v.([]any) {
			argName := arg.(map[string]any)["arg_name"].(string)
			argDataType := arg.(map[string]any)["arg_data_type"].(string)
			dataType, err := datatypes.ParseDataType(argDataType)
			if err != nil {
				return nil, err
			}
			request := sdk.NewFunctionArgumentRequest(argName, dataType)

			if argDefaultValue, defaultValuePresent := arg.(map[string]any)["arg_default_value"]; defaultValuePresent && argDefaultValue.(string) != "" {
				request.WithDefaultValue(argDefaultValue.(string))
			}

			args = append(args, *request)
		}
	}
	return args, nil
}

func parseFunctionImportsCommon(d *schema.ResourceData) ([]sdk.FunctionImportRequest, error) {
	imports := make([]sdk.FunctionImportRequest, 0)
	if v, ok := d.GetOk("imports"); ok {
		for _, imp := range v.(*schema.Set).List() {
			stageLocation := imp.(map[string]any)["stage_location"].(string)
			pathOnStage := imp.(map[string]any)["path_on_stage"].(string)
			imports = append(imports, *sdk.NewFunctionImportRequest().WithImport(fmt.Sprintf("@%s/%s", stageLocation, pathOnStage)))
		}
	}
	return imports, nil
}

func parseFunctionPackagesCommon(d *schema.ResourceData) ([]sdk.FunctionPackageRequest, error) {
	packages := make([]sdk.FunctionPackageRequest, 0)
	if v, ok := d.GetOk("packages"); ok {
		for _, pkg := range v.(*schema.Set).List() {
			packages = append(packages, *sdk.NewFunctionPackageRequest().WithPackage(pkg.(string)))
		}
	}
	return packages, nil
}

func parseFunctionTargetPathCommon(d *schema.ResourceData) (string, error) {
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

func parseFunctionReturnsCommon(d *schema.ResourceData) (*sdk.FunctionReturnsRequest, error) {
	returnTypeRaw := d.Get("return_type").(string)
	dataType, err := datatypes.ParseDataType(returnTypeRaw)
	if err != nil {
		return nil, err
	}
	returns := sdk.NewFunctionReturnsRequest()
	switch v := dataType.(type) {
	case *datatypes.TableDataType:
		var cr []sdk.FunctionColumnRequest
		for _, c := range v.Columns() {
			cr = append(cr, *sdk.NewFunctionColumnRequest(c.ColumnName(), c.ColumnType()))
		}
		returns.WithTable(*sdk.NewFunctionReturnsTableRequest().WithColumns(cr))
	default:
		returns.WithResultDataType(*sdk.NewFunctionReturnsResultDataTypeRequest(dataType))
	}
	return returns, nil
}

func setFunctionImportsInBuilder[T any](d *schema.ResourceData, setImports func([]sdk.FunctionImportRequest) T) error {
	imports, err := parseFunctionImportsCommon(d)
	if err != nil {
		return err
	}
	setImports(imports)
	return nil
}

func setFunctionPackagesInBuilder[T any](d *schema.ResourceData, setPackages func([]sdk.FunctionPackageRequest) T) error {
	packages, err := parseFunctionPackagesCommon(d)
	if err != nil {
		return err
	}
	setPackages(packages)
	return nil
}

func setFunctionTargetPathInBuilder[T any](d *schema.ResourceData, setTargetPath func(string) T) error {
	tp, err := parseFunctionTargetPathCommon(d)
	if err != nil {
		return err
	}
	if tp != "" {
		setTargetPath(tp)
	}
	return nil
}

func queryAllFunctionDetailsCommon(ctx context.Context, d *schema.ResourceData, client *sdk.Client, id sdk.SchemaObjectIdentifierWithArguments) (*allFunctionDetailsCommon, diag.Diagnostics) {
	function, err := client.Functions.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query function. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Function: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}

	functionDetails, err := client.Functions.DescribeDetails(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			log.Printf("[DEBUG] function (%s) not found or we are not authorized. Err: %s", d.Id(), err)
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query function. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Function: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}

	functionParameters, err := client.Functions.ShowParameters(ctx, id)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return &allFunctionDetailsCommon{
		function:           function,
		functionDetails:    functionDetails,
		functionParameters: functionParameters,
	}, nil
}

type allFunctionDetailsCommon struct {
	function           *sdk.Function
	functionDetails    *sdk.FunctionDetails
	functionParameters []*sdk.Parameter
}

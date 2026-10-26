package gen

import (
	"strings"
	"text/template"

	_ "embed"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

var (
	//go:embed templates/lifecycle.tmpl
	lifecycleTemplateContent string
	LifecycleTemplate        *template.Template

	//go:embed templates/schema.tmpl
	schemaTemplateContent string
	SchemaTemplate        *template.Template

	//go:embed templates/parameters.tmpl
	parametersTemplateContent string
	ParametersTemplate        *template.Template

	//go:embed templates/sub_templates/constructor.tmpl
	constructorTemplateContent string

	//go:embed templates/sub_templates/parse_id.tmpl
	parseIdTemplateContent string

	//go:embed templates/sub_templates/generated_parse_id.tmpl
	generatedParseIdTemplateContent string

	//go:embed templates/sub_templates/parse_id_from_config.tmpl
	parseIdFromConfigTemplateContent string

	//go:embed templates/sub_templates/generated_parse_id_from_config.tmpl
	generatedParseIdFromConfigTemplateContent string

	//go:embed templates/sub_templates/before_create.tmpl
	beforeCreateTemplateContent string

	//go:embed templates/sub_templates/new_create_request.tmpl
	newCreateRequestTemplateContent string

	//go:embed templates/sub_templates/generated_create_request.tmpl
	generatedCreateRequestTemplateContent string

	//go:embed templates/sub_templates/apply_create_optionals.tmpl
	applyCreateOptionalsTemplateContent string

	//go:embed templates/sub_templates/generated_apply_create_optionals.tmpl
	generatedApplyCreateOptionalsTemplateContent string

	//go:embed templates/sub_templates/create_constructor_args.tmpl
	createConstructorArgsTemplateContent string

	//go:embed templates/sub_templates/create_optional_calls.tmpl
	createOptionalCallsTemplateContent string

	//go:embed templates/sub_templates/create_in_sdk.tmpl
	createInSdkTemplateContent string

	//go:embed templates/sub_templates/after_create.tmpl
	afterCreateTemplateContent string

	//go:embed templates/sub_templates/after_alter.tmpl
	afterAlterTemplateContent string

	//go:embed templates/sub_templates/show_by_id_safely_in_sdk.tmpl
	showByIdSafelyInSdkTemplateContent string

	//go:embed templates/sub_templates/describe_in_sdk.tmpl
	describeInSdkTemplateContent string

	//go:embed templates/sub_templates/handle_external_changes.tmpl
	handleExternalChangesTemplateContent string

	//go:embed templates/sub_templates/handle_external_changes_in_describe.tmpl
	handleExternalChangesInDescribeTemplateContent string

	//go:embed templates/sub_templates/set_state_to_values_from_config.tmpl
	setStateToValuesFromConfigTemplateContent string

	//go:embed templates/sub_templates/show_output_set.tmpl
	showOutputSetTemplateContent string

	//go:embed templates/sub_templates/describe_output_set.tmpl
	describeOutputSetTemplateContent string

	//go:embed templates/sub_templates/set_describe_output.tmpl
	setDescribeOutputTemplateContent string

	//go:embed templates/sub_templates/fqn_set.tmpl
	fqnSetTemplateContent string

	//go:embed templates/sub_templates/set_config_fields.tmpl
	setConfigFieldsTemplateContent string

	//go:embed templates/sub_templates/hierarchy_rename.tmpl
	hierarchyRenameTemplateContent string

	//go:embed templates/sub_templates/before_alter.tmpl
	beforeAlterTemplateContent string

	//go:embed templates/sub_templates/apply_set_unset.tmpl
	applySetUnsetTemplateContent string

	//go:embed templates/sub_templates/generated_apply_set_unset.tmpl
	generatedApplySetUnsetTemplateContent string

	//go:embed templates/sub_templates/alter_set_in_sdk.tmpl
	alterSetInSdkTemplateContent string

	//go:embed templates/sub_templates/alter_unset_in_sdk.tmpl
	alterUnsetInSdkTemplateContent string

	//go:embed templates/sub_templates/set_fields_not_set_by_read.tmpl
	setFieldsNotSetByReadTemplateContent string

	//go:embed templates/sub_templates/show_by_id_in_sdk.tmpl
	showByIdInSdkTemplateContent string

	//go:embed templates/sub_templates/drop_safely_in_sdk.tmpl
	dropSafelyInSdkTemplateContent string

	//go:embed templates/sub_templates/generated_delete.tmpl
	generatedDeleteTemplateContent string

	//go:embed templates/sub_templates/identity_schema.tmpl
	identitySchemaTemplateContent string

	//go:embed templates/sub_templates/attributes_schema.tmpl
	attributesSchemaTemplateContent string

	//go:embed templates/sub_templates/attributes_schema_definition.tmpl
	attributesSchemaDefinitionTemplateContent string

	//go:embed templates/sub_templates/fqn_schema.tmpl
	fqnSchemaTemplateContent string

	//go:embed templates/sub_templates/show_output_schema.tmpl
	showOutputSchemaTemplateContent string

	//go:embed templates/sub_templates/describe_output_schema.tmpl
	describeOutputSchemaTemplateContent string

	//go:embed templates/sub_templates/apply_parameters_create.tmpl
	applyParametersCreateTemplateContent string

	//go:embed templates/sub_templates/show_parameters_in_sdk.tmpl
	showParametersInSdkTemplateContent string

	//go:embed templates/sub_templates/parameters_details_from_raw.tmpl
	parametersDetailsFromRawTemplateContent string

	//go:embed templates/sub_templates/set_parameters_fields.tmpl
	setParametersFieldsTemplateContent string

	//go:embed templates/sub_templates/parameters_output_set.tmpl
	parametersOutputSetTemplateContent string

	//go:embed templates/sub_templates/apply_parameters_changes.tmpl
	applyParametersChangesTemplateContent string

	//go:embed templates/sub_templates/parameters_attributes_schema.tmpl
	parametersAttributesSchemaTemplateContent string

	//go:embed templates/sub_templates/parameters_output_schema.tmpl
	parametersOutputSchemaTemplateContent string
)

func init() {
	subTemplates := template.New("subTemplates").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.Unimplemented,
		genhelpers.CamelToWords,
		genhelpers.FirstLetterLowercase,
		genhelpers.SnakeCaseToCamel,
		strings.ToUpper,
	))
	subTemplates, _ = subTemplates.New("constructor").Parse(constructorTemplateContent)
	subTemplates, _ = subTemplates.New("parseId").Parse(parseIdTemplateContent)
	subTemplates, _ = subTemplates.New("generatedParseId").Parse(generatedParseIdTemplateContent)
	subTemplates, _ = subTemplates.New("parseIdFromConfig").Parse(parseIdFromConfigTemplateContent)
	subTemplates, _ = subTemplates.New("generatedParseIdFromConfig").Parse(generatedParseIdFromConfigTemplateContent)
	subTemplates, _ = subTemplates.New("beforeCreate").Parse(beforeCreateTemplateContent)
	subTemplates, _ = subTemplates.New("newCreateRequest").Parse(newCreateRequestTemplateContent)
	subTemplates, _ = subTemplates.New("createConstructorArgs").Parse(createConstructorArgsTemplateContent)
	subTemplates, _ = subTemplates.New("createOptionalCalls").Parse(createOptionalCallsTemplateContent)
	subTemplates, _ = subTemplates.New("generatedCreateRequest").Parse(generatedCreateRequestTemplateContent)
	subTemplates, _ = subTemplates.New("applyCreateOptionals").Parse(applyCreateOptionalsTemplateContent)
	subTemplates, _ = subTemplates.New("generatedApplyCreateOptionals").Parse(generatedApplyCreateOptionalsTemplateContent)
	subTemplates, _ = subTemplates.New("createInSdk").Parse(createInSdkTemplateContent)
	subTemplates, _ = subTemplates.New("afterCreate").Parse(afterCreateTemplateContent)
	subTemplates, _ = subTemplates.New("afterAlter").Parse(afterAlterTemplateContent)
	subTemplates, _ = subTemplates.New("showByIdSafelyInSdk").Parse(showByIdSafelyInSdkTemplateContent)
	subTemplates, _ = subTemplates.New("describeInSdk").Parse(describeInSdkTemplateContent)
	subTemplates, _ = subTemplates.New("handleExternalChanges").Parse(handleExternalChangesTemplateContent)
	subTemplates, _ = subTemplates.New("handleExternalChangesInDescribe").Parse(handleExternalChangesInDescribeTemplateContent)
	subTemplates, _ = subTemplates.New("setStateToValuesFromConfig").Parse(setStateToValuesFromConfigTemplateContent)
	subTemplates, _ = subTemplates.New("showOutputSet").Parse(showOutputSetTemplateContent)
	subTemplates, _ = subTemplates.New("describeOutputSet").Parse(describeOutputSetTemplateContent)
	subTemplates, _ = subTemplates.New("setDescribeOutput").Parse(setDescribeOutputTemplateContent)
	subTemplates, _ = subTemplates.New("fqnSet").Parse(fqnSetTemplateContent)
	subTemplates, _ = subTemplates.New("setConfigFields").Parse(setConfigFieldsTemplateContent)
	subTemplates, _ = subTemplates.New("hierarchyRename").Parse(hierarchyRenameTemplateContent)
	subTemplates, _ = subTemplates.New("beforeAlter").Parse(beforeAlterTemplateContent)
	subTemplates, _ = subTemplates.New("applySetUnset").Parse(applySetUnsetTemplateContent)
	subTemplates, _ = subTemplates.New("generatedApplySetUnset").Parse(generatedApplySetUnsetTemplateContent)
	subTemplates, _ = subTemplates.New("alterSetInSdk").Parse(alterSetInSdkTemplateContent)
	subTemplates, _ = subTemplates.New("alterUnsetInSdk").Parse(alterUnsetInSdkTemplateContent)
	subTemplates, _ = subTemplates.New("setFieldsNotSetByRead").Parse(setFieldsNotSetByReadTemplateContent)
	subTemplates, _ = subTemplates.New("showByIdInSdk").Parse(showByIdInSdkTemplateContent)
	subTemplates, _ = subTemplates.New("dropSafelyInSdk").Parse(dropSafelyInSdkTemplateContent)
	subTemplates, _ = subTemplates.New("generatedDelete").Parse(generatedDeleteTemplateContent)
	subTemplates, _ = subTemplates.New("identitySchema").Parse(identitySchemaTemplateContent)
	subTemplates, _ = subTemplates.New("attributesSchema").Parse(attributesSchemaTemplateContent)
	subTemplates, _ = subTemplates.New("attributesSchemaDefinition").Parse(attributesSchemaDefinitionTemplateContent)
	subTemplates, _ = subTemplates.New("fqnSchema").Parse(fqnSchemaTemplateContent)
	subTemplates, _ = subTemplates.New("showOutputSchema").Parse(showOutputSchemaTemplateContent)
	subTemplates, _ = subTemplates.New("describeOutputSchema").Parse(describeOutputSchemaTemplateContent)
	subTemplates, _ = subTemplates.New("applyParametersCreate").Parse(applyParametersCreateTemplateContent)
	subTemplates, _ = subTemplates.New("showParametersInSdk").Parse(showParametersInSdkTemplateContent)
	subTemplates, _ = subTemplates.New("parametersDetailsFromRaw").Parse(parametersDetailsFromRawTemplateContent)
	subTemplates, _ = subTemplates.New("setParametersFields").Parse(setParametersFieldsTemplateContent)
	subTemplates, _ = subTemplates.New("parametersOutputSet").Parse(parametersOutputSetTemplateContent)
	subTemplates, _ = subTemplates.New("applyParametersChanges").Parse(applyParametersChangesTemplateContent)
	subTemplates, _ = subTemplates.New("parametersAttributesSchema").Parse(parametersAttributesSchemaTemplateContent)
	subTemplates, _ = subTemplates.New("parametersOutputSchema").Parse(parametersOutputSchemaTemplateContent)
	LifecycleTemplate, _ = subTemplates.New("lifecycleTemplate").Parse(lifecycleTemplateContent)
	SchemaTemplate, _ = subTemplates.New("schemaTemplate").Parse(schemaTemplateContent)
	ParametersTemplate, _ = subTemplates.New("parametersTemplate").Parse(parametersTemplateContent)
}

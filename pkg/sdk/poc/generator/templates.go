package generator

import (
	_ "embed"
	"strings"
	"text/template"
)

var (
	//go:embed templates/package.tmpl
	packageTemplateContent string
	PackageTemplate, _     = template.New("packageTemplates").Parse(packageTemplateContent)

	//go:embed templates/interface.tmpl
	interfaceTemplateContent string
	InterfaceTemplate, _     = template.New("interfaceTemplate").Funcs(template.FuncMap{
		"deref":     func(p *DescriptionMappingKind) string { return string(*p) },
		"hasPrefix": func(str string, prefix string) bool { return strings.HasPrefix(str, prefix) },
	}).Parse(interfaceTemplateContent)

	//go:embed templates/operation_struct.tmpl
	operationStructTemplateContent string
	OperationStructTemplate, _     = template.New("optionsTemplate").Parse(operationStructTemplateContent)

	//go:embed templates/struct.tmpl
	structTemplateContent string
	StructTemplate, _     = template.New("structTemplate").Parse(structTemplateContent)

	//go:embed templates/show_object_id_method.tmpl
	showObjectIdMethodTemplateContent string
	ShowObjectIdMethodTemplate, _     = template.New("showObjectIdMethodTemplate").Parse(showObjectIdMethodTemplateContent)

	//go:embed templates/show_object_type_method.tmpl
	showObjectTypeMethodTemplateContent string
	ShowObjectTypeMethodTemplate, _     = template.New("showObjectTypeMethodTemplate").Parse(showObjectTypeMethodTemplateContent)

	//go:embed templates/dto_declarations.tmpl
	dtoDeclarationsTemplateContent string
	DtoTemplate, _                 = template.New("dtoTemplate").Parse(dtoDeclarationsTemplateContent)

	//go:embed templates/dto_structs.tmpl
	dtoStructsTemplateContent string
	DtoDeclTemplate, _        = template.New("dtoTemplate").Parse(dtoStructsTemplateContent)

	//go:embed templates/implementation.tmpl
	implementationTemplateContent string
	ImplementationTemplate        *template.Template

	//go:embed templates/unit_tests.tmpl
	unitTestTemplateContent string
	UnitTestsTemplate       *template.Template

	//go:embed templates/validations.tmpl
	validationTemplateContent string
	ValidationsTemplate       *template.Template

	//go:embed templates/sub_templates/to_opts_mapping.tmpl
	toOptsMappingTemplateContent string

	//go:embed templates/sub_templates/convert.tmpl
	convertTemplateContent string

	//go:embed templates/sub_templates/implementation_mappings.tmpl
	implementationMappingsTemplateContent string

	//go:embed templates/sub_templates/implementation_functions.tmpl
	implementationFunctionsTemplateContent string

	//go:embed templates/sub_templates/validation_test.tmpl
	validationTestTemplateContent string

	//go:embed templates/sub_templates/validation_tests.tmpl
	validationTestsTemplateContent string

	//go:embed templates/sub_templates/validation_implementation.tmpl
	validationImplementationTemplateContent string
)

func init() {
	subTemplates := template.New("subTemplates").Funcs(template.FuncMap{
		"deref":     func(p *DescriptionMappingKind) string { return string(*p) },
		"hasPrefix": func(str string, prefix string) bool { return strings.HasPrefix(str, prefix) },
	})
	subTemplates, _ = subTemplates.New("toOptsMapping").Parse(toOptsMappingTemplateContent)
	subTemplates, _ = subTemplates.New("convert").Parse(convertTemplateContent)
	subTemplates, _ = subTemplates.New("implementationMappings").Parse(implementationMappingsTemplateContent)
	subTemplates, _ = subTemplates.New("implementationFunctions").Parse(implementationFunctionsTemplateContent)
	subTemplates, _ = subTemplates.New("validationTest").Parse(validationTestTemplateContent)
	subTemplates, _ = subTemplates.New("validationTests").Parse(validationTestsTemplateContent)
	subTemplates, _ = subTemplates.New("validationImplementation").Parse(validationImplementationTemplateContent)

	ImplementationTemplate, _ = subTemplates.New("implementationTemplate").Parse(implementationTemplateContent)
	UnitTestsTemplate, _ = subTemplates.New("unitTestsTemplate").Parse(unitTestTemplateContent)
	ValidationsTemplate, _ = subTemplates.New("validationsTemplate").Parse(validationTemplateContent)
}

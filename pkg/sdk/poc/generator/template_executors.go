package generator

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"text/template"
)

var (
	generatedStructs []string
	generatedDtos    []string
)

func GenerateInterface(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, InterfaceTemplate, def)
	// TODO [this PR]: replace this with the template logic
	for _, o := range def.Operations {
		if o.OptsField != nil {
			generateOptionsStruct(writer, o)
		}

		if o.Name == string(OperationKindShow) {
			idKind, err := toObjectIdentifierKind(def.IdentifierKind)
			if err != nil {
				log.Printf("[WARN] for showObjectIdMethod: %v", err)
			}
			if checkRequiredFieldsForIdMethod(def.NameSingular, o.HelperStructs, idKind) {
				generateShowObjectIdMethod(writer, newShowObjectIDMethod(def.NameSingular, idKind))
			}

			generateShowObjectTypeMethod(writer, newShowObjectTypeMethod(def.NameSingular))
		}
	}
}

func generateShowObjectIdMethod(writer io.Writer, m *ShowObjectIdMethod) {
	printTo(writer, ShowObjectIdMethodTemplate, m)
}

func generateShowObjectTypeMethod(writer io.Writer, m *ShowObjectTypeMethod) {
	printTo(writer, ShowObjectTypeMethodTemplate, m)
}

// TODO [next PR]: remove this method
func generateOptionsStruct(writer io.Writer, operation *Operation) {
	printTo(writer, OperationStructTemplate, operation)

	for _, f := range operation.HelperStructs {
		if !slices.Contains(generatedStructs, f.KindNoPtr()) {
			generateStruct(writer, f)
		}
	}

	for _, f := range operation.OptsField.Fields {
		if len(f.Fields) > 0 && !slices.Contains(generatedStructs, f.KindNoPtr()) {
			generateStruct(writer, f)
		}
	}
}

// TODO [next PR]: remove this method
func generateStruct(writer io.Writer, field *Field) {
	if !slices.Contains(generatedStructs, field.KindNoPtr()) {
		fmt.Println("Generating: " + field.KindNoPtr())
		printTo(writer, StructTemplate, field)
		generatedStructs = append(generatedStructs, field.KindNoPtr())
	}

	for _, f := range field.Fields {
		if len(f.Fields) > 0 && !slices.Contains(generatedStructs, f.Name) {
			generateStruct(writer, f)
		}
	}
}

func GenerateDtos(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, DtoTemplate, def)
	for _, o := range def.Operations {
		if o.OptsField != nil {
			generateDtoDecls(writer, o.OptsField)
		}
	}
}

func generateDtoDecls(writer io.Writer, field *Field) {
	if !slices.Contains(generatedDtos, field.DtoDecl()) {
		printTo(writer, DtoDeclTemplate, field)
		generatedDtos = append(generatedDtos, field.DtoDecl())

		for _, f := range field.Fields {
			if f.IsStruct() {
				generateDtoDecls(writer, f)
			}
		}
	}
}

func GenerateImplementation(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, ImplementationTemplate, def)
}

func GenerateUnitTests(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, UnitTestsTemplate, def)
}

func GenerateValidations(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, ValidationsTemplate, def)
}

func generatePackageDirective(writer io.Writer) {
	printTo(writer, PackageTemplate, os.Getenv("GOPACKAGE"))
}

func printTo(writer io.Writer, template *template.Template, model any) {
	err := template.Execute(writer, model)
	if err != nil {
		log.Panicln(err)
	}
}

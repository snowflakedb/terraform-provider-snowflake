package gen

import (
	"fmt"
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

var AllSdkObjectDefinitions = make([]*Interface, 0)

func GetSdkDefinitions() []*Interface {
	allDefinitions := AllSdkObjectDefinitions
	interfaces := make([]*Interface, len(allDefinitions))
	for idx, def := range allDefinitions {
		preprocessDefinition(def)
		interfaces[idx] = def
	}
	return interfaces
}

// preprocessDefinition is needed because current simple builder is not ideal, should be removed later
func preprocessDefinition(definition *Interface) {
	generatedStructs := make([]string, 0)
	generatedDtos := make([]string, 0)

	for _, o := range definition.Operations {
		o.ObjectInterface = definition
		if o.OptsField != nil {
			o.OptsField = deepCopyFieldHierarchy(o.OptsField)

			o.OptsField.Name = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
			o.OptsField.Kind = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
			setParent(o.OptsField)

			// TODO [SNOW-2324252]: this logic is currently the old logic adjusted. Let's clean it after new generation is working.
			// fill out StructsToGenerate; it replaces the old generateOptionsStruct and generateStruct
			structsToGenerate := make([]*Field, 0)
			for _, f := range o.HelperStructs {
				if !slices.Contains(generatedStructs, f.KindNoPtr()) {
					structsToGenerate, generatedStructs = addStructToGenerate(f, structsToGenerate, generatedStructs)
				}
			}
			for idx, f := range o.OptsField.Fields {
				if len(f.Fields) > 0 && !slices.Contains(generatedStructs, f.KindNoPtr()) {
					structsToGenerate, generatedStructs = addStructToGenerate(&(o.OptsField.Fields[idx]), structsToGenerate, generatedStructs)
				}
			}
			log.Printf("[DEBUG] Structs to generate (length: %d): %v", len(structsToGenerate), structsToGenerate)
			o.StructsToGenerate = structsToGenerate

			// TODO [SNOW-2324252]: this logic is currently the old logic adjusted. Let's clean it after new generation is working.
			// fill out ObjectIdMethod and ObjectIdType; it replaces the old template executors logic
			if o.Name == string(OperationKindShow) {
				// TODO [SNOW-2324252]: do we really conversion logic? The definition file should handle this.
				idKind, err := ToObjectIdentifierKind(definition.IdentifierKind)
				if err != nil {
					log.Printf("[WARN] for showObjectIdMethod: %v", err)
				}
				if CheckRequiredFieldsForIdMethod(definition.NameSingular, o.HelperStructs, idKind) {
					o.ObjectIdMethod = NewShowObjectIDMethod(definition.NameSingular, idKind)
				}

				o.ObjectTypeMethod = NewShowObjectTypeMethod(definition.NameSingular)
			}

			// TODO [SNOW-2324252]: this logic is currently the old logic adjusted. Let's clean it after new generation is working.
			// fill out DtosToGenerate; it replaces the old GenerateDtos and generateDtoDecls logic
			dtosToGenerate := make([]*Field, 0)
			dtosToGenerate, generatedDtos = addDtoToGenerate(o.OptsField, dtosToGenerate, generatedDtos)
			log.Printf("[DEBUG] Dtos to generate (length: %d): %v", len(dtosToGenerate), dtosToGenerate)
			o.DtosToGenerate = dtosToGenerate
		}
	}
}

func deepCopyFieldHierarchy(field *Field) *Field {
	newField := *field
	if field.Fields != nil && len(field.Fields) > 0 {
		children := make([]Field, len(field.Fields))
		for idx, f := range field.Fields {
			newF := deepCopyFieldHierarchy(&f)
			children[idx] = *newF
		}
		newField.Fields = children
	}
	return &newField
}

func setParent(field *Field) {
	for idx, f := range field.Fields {
		if f.Parent != nil {
			log.Panicf("Field %s already has a parent\nold parent: %s (path: %s)\nnew parent: %s (path: %s);\n\nit is caused by the current incorrect implementation of nested fields;\nreuse the common definition by wrapping it in function invocation", f.Name, f.Parent.KindNoPtr(), f.Parent.PathWithRoot(), field.Name, field.PathWithRoot())
		}
		(&(field.Fields[idx])).Parent = field
		setParent(&(field.Fields[idx]))
	}
}

func addStructToGenerate(field *Field, structsToGenerate []*Field, generatedStructs []string) ([]*Field, []string) {
	if !slices.Contains(generatedStructs, field.KindNoPtr()) {
		log.Printf("[DEBUG] Adding %s (path: %s) to structs to be generated", field.KindNoPtr(), field.PathWithRoot())
		structsToGenerate = append(structsToGenerate, field)
		generatedStructs = append(generatedStructs, field.KindNoPtr())
	} else {
		log.Printf("[DEBUG] Struct %s (path: %s) already queued for generation", field.KindNoPtr(), field.PathWithRoot())
	}

	for idx, f := range field.Fields {
		if len(f.Fields) > 0 && !slices.Contains(generatedStructs, f.Name) {
			structsToGenerate, generatedStructs = addStructToGenerate(&(field.Fields[idx]), structsToGenerate, generatedStructs)
		}
	}
	return structsToGenerate, generatedStructs
}

func addDtoToGenerate(field *Field, dtosToGenerate []*Field, generatedDtos []string) ([]*Field, []string) {
	if !slices.Contains(generatedDtos, field.DtoDecl()) {
		log.Printf("[DEBUG] Adding %s (path: %s) to structs to be generated", field.DtoDecl(), field.PathWithRoot())
		dtosToGenerate = append(dtosToGenerate, field)
		generatedDtos = append(generatedDtos, field.DtoDecl())

		for idx, f := range field.Fields {
			if f.IsStruct() {
				dtosToGenerate, generatedDtos = addDtoToGenerate(&(field.Fields[idx]), dtosToGenerate, generatedDtos)
			}
		}
	} else {
		log.Printf("[DEBUG] Struct %s (path: %s) already queued for generation", field.DtoDecl(), field.PathWithRoot())
	}
	return dtosToGenerate, generatedDtos
}

func ExtendInterface() func(*Interface, *genhelpers.PreambleModel) *Interface {
	return func(i *Interface, preamble *genhelpers.PreambleModel) *Interface {
		i.PreambleModel = preamble
		return i
	}
}

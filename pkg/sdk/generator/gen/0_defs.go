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
			o.OptsField.Name = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
			o.OptsField.Kind = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
			setParent(o.OptsField)
			relocateIdentifierElementSliceValidations(o.OptsField)
			resolveInterfaceIdentifierKinds(o.OptsField, definition.IdentifierKind)

			// TODO [SNOW-2324252]: this logic is currently the old logic adjusted. Let's clean it after new generation is working.
			// fill out StructsToGenerate; it replaces the old generateOptionsStruct and generateStruct
			structsToGenerate := make([]*Field, 0)
			for _, f := range o.HelperStructs {
				if !slices.Contains(generatedStructs, f.KindNoPtr()) {
					structsToGenerate, generatedStructs = addStructToGenerate(f, structsToGenerate, generatedStructs)
				}
			}
			for idx, f := range o.OptsField.Fields {
				if f.IsShared {
					continue
				}
				if f.IsStruct() && !slices.Contains(generatedStructs, f.KindNoPtr()) {
					structsToGenerate, generatedStructs = addStructToGenerate(&o.OptsField.Fields[idx], structsToGenerate, generatedStructs)
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

				typeName := definition.NameSingular
				if definition.ShowObjectTypeName != "" {
					typeName = definition.ShowObjectTypeName
				}
				o.ObjectTypeMethod = NewShowObjectTypeMethod(definition.ShowObjectName, typeName)
			}

			if o.NoRequest {
				if hasDtoField(o.OptsField) {
					log.Panicf("operation %s: WithNoRequest is set but the query struct has request fields", o.Name)
				}
				o.DtosToGenerate = nil
			} else {
				// TODO [SNOW-2324252]: this logic is currently the old logic adjusted. Let's clean it after new generation is working.
				// fill out DtosToGenerate; it replaces the old GenerateDtos and generateDtoDecls logic
				dtosToGenerate := make([]*Field, 0)
				dtosToGenerate, generatedDtos = addDtoToGenerate(o.OptsField, dtosToGenerate, generatedDtos)
				log.Printf("[DEBUG] Dtos to generate (length: %d): %v", len(dtosToGenerate), dtosToGenerate)
				o.DtosToGenerate = dtosToGenerate
			}
		}
	}

	// Deduplicate convert() and convertibleRow guard emissions.
	// A convert() method can only be declared once per receiver type per package.
	seenMappingReceivers := make([]string, 0)
	for _, o := range definition.Operations {
		for _, mapping := range []*Mapping{o.ShowMapping, o.DescribeMapping, o.InstanceMethodMapping} {
			if mapping == nil {
				continue
			}
			if slices.Contains(seenMappingReceivers, mapping.From.Name) {
				mapping.SkipConvert = true
			} else {
				seenMappingReceivers = append(seenMappingReceivers, mapping.From.Name)
			}
		}
	}

	// Built for every object regardless of opt-in; main.go's generation part decides whether to render it.
	definition.UnitTests = definition.buildUnitTestsModel()
}

// resolveInterfaceIdentifierKinds walks all fields recursively and replaces sentinel kind values
// with the actual identifier type from the interface definition.
func resolveInterfaceIdentifierKinds(field *Field, identifierKind string) {
	switch field.Kind {
	case InterfaceIdentifierKind:
		field.Kind = identifierKind
	case InterfaceIdentifierPointerKind:
		field.Kind = KindOfPointer(identifierKind)
	}
	for idx := range field.Fields {
		resolveInterfaceIdentifierKinds(&field.Fields[idx], identifierKind)
	}
}

// relocateIdentifierElementSliceValidations moves WithValidation(ValidIdentifier, "SliceField")
// from the container onto the slice field when that child is a flat identifier list or a
// []TagAssociation. The validations template only opens a per-element loop when the current
// field is the slice. Mutate through f.Fields[i] — FindChild returns a copy.
func relocateIdentifierElementSliceValidations(f *Field) {
	remaining := make([]*Validation, 0, len(f.Validations))
	for _, v := range f.Validations {
		moved := false
		if v.Type == ValidIdentifier && len(v.FieldNames) == 1 {
			for i := range f.Fields {
				if f.Fields[i].Name == v.FieldNames[0] && f.Fields[i].acceptsRelocatedValidIdentifier() {
					f.Fields[i].Validations = append(f.Fields[i].Validations, v)
					moved = true
					break
				}
			}
		}
		if !moved {
			remaining = append(remaining, v)
		}
	}
	f.Validations = remaining
	for i := range f.Fields {
		relocateIdentifierElementSliceValidations(&f.Fields[i])
	}
}

func setParent(field *Field) {
	for idx, f := range field.Fields {
		if f.Parent != nil {
			log.Panicf("Field %s already has a parent\nold parent: %s (path: %s)\nnew parent: %s (path: %s);\n\nit is caused by the current incorrect implementation of nested fields;\nreuse the common definition by wrapping it in function invocation", f.Name, f.Parent.KindNoPtr(), f.Parent.PathWithRoot(), field.Name, field.PathWithRoot())
		}
		(&field.Fields[idx]).Parent = field
		setParent(&field.Fields[idx])
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
		if f.IsShared {
			continue
		}
		if f.IsStruct() && !slices.Contains(generatedStructs, f.KindNoPtr()) {
			structsToGenerate, generatedStructs = addStructToGenerate(&field.Fields[idx], structsToGenerate, generatedStructs)
		}
	}
	return structsToGenerate, generatedStructs
}

// hasDtoField reports whether any field would appear on a request DTO.
func hasDtoField(field *Field) bool {
	for i := range field.Fields {
		child := &field.Fields[i]
		if child.ShouldBeInDto() || hasDtoField(child) {
			return true
		}
	}
	return false
}

func addDtoToGenerate(field *Field, dtosToGenerate []*Field, generatedDtos []string) ([]*Field, []string) {
	if !slices.Contains(generatedDtos, field.DtoDecl()) {
		log.Printf("[DEBUG] Adding %s (path: %s) to structs to be generated", field.DtoDecl(), field.PathWithRoot())
		dtosToGenerate = append(dtosToGenerate, field)
		generatedDtos = append(generatedDtos, field.DtoDecl())

		for idx, f := range field.Fields {
			if f.IsShared {
				continue
			}
			if f.HasAnyFields() {
				dtosToGenerate, generatedDtos = addDtoToGenerate(&field.Fields[idx], dtosToGenerate, generatedDtos)
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

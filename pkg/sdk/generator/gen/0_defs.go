package gen

import (
	"fmt"
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

var AllSdkObjectDefinitions = make([]*Interface, 0)

func GetSdkDefinitions() []*Interface {
	interfaces := make([]*Interface, len(AllSdkObjectDefinitions))
	copy(interfaces, AllSdkObjectDefinitions)
	return interfaces
}

// MergeSdkDefinitions keeps old interfaces whose names ssot does not emit, then appends the ssot interfaces.
func MergeSdkDefinitions(old, fromSsot []*Interface) []*Interface {
	plugged := make(map[string]struct{}, len(fromSsot))
	for _, iface := range fromSsot {
		plugged[iface.Name] = struct{}{}
	}
	merged := make([]*Interface, 0, len(old)+len(fromSsot))
	for _, iface := range old {
		if _, ok := plugged[iface.Name]; ok {
			continue
		}
		merged = append(merged, iface)
	}
	return append(merged, fromSsot...)
}

// PrepareSdkDefinition prepares one definition on the worker that generates it.
// It is the SDK generator's WithObjectPrepare hook.
func PrepareSdkDefinition(definition *Interface, logger *log.Logger) {
	preprocessDefinition(definition, logger)
}

// preprocessDefinition is needed because current simple builder is not ideal, should be removed later
func preprocessDefinition(definition *Interface, logger *log.Logger) {
	generatedStructs := make([]string, 0)
	generatedDtos := make([]string, 0)

	for _, o := range definition.Operations {
		o.ObjectInterface = definition
		if o.OptsField != nil {
			o.OptsField.Name = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
			o.OptsField.Kind = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
			setParent(o.OptsField, logger)
			relocateIdentifierElementSliceValidations(o.OptsField)
			resolveInterfaceIdentifierKinds(o.OptsField, definition.IdentifierKind)

			// TODO [SNOW-2324252]: this logic is currently the old logic adjusted. Let's clean it after new generation is working.
			// fill out StructsToGenerate; it replaces the old generateOptionsStruct and generateStruct
			structsToGenerate := make([]*Field, 0)
			for _, f := range o.HelperStructs {
				if !slices.Contains(generatedStructs, f.KindNoPtr()) {
					structsToGenerate, generatedStructs = addStructToGenerate(f, structsToGenerate, generatedStructs, logger)
				}
			}
			for idx, f := range o.OptsField.Fields {
				if f.IsShared {
					continue
				}
				if f.IsStruct() && !slices.Contains(generatedStructs, f.KindNoPtr()) {
					structsToGenerate, generatedStructs = addStructToGenerate(&o.OptsField.Fields[idx], structsToGenerate, generatedStructs, logger)
				}
			}
			logger.Printf("[DEBUG] Structs to generate (length: %d): %v", len(structsToGenerate), structsToGenerate)
			o.StructsToGenerate = structsToGenerate

			// TODO [SNOW-2324252]: this logic is currently the old logic adjusted. Let's clean it after new generation is working.
			// fill out ObjectIdMethod and ObjectIdType; it replaces the old template executors logic
			if o.Name == string(OperationKindShow) {
				// TODO [SNOW-2324252]: do we really conversion logic? The definition file should handle this.
				idKind, err := ToObjectIdentifierKind(definition.IdentifierKind)
				if err != nil {
					logger.Printf("[WARN] for showObjectIdMethod: %v", err)
				}
				if CheckRequiredFieldsForIdMethod(definition.NameSingular, o.HelperStructs, idKind, logger) {
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
					logger.Panicf("operation %s: WithNoRequest is set but the query struct has request fields", o.Name)
				}
				o.DtosToGenerate = nil
			} else {
				// TODO [SNOW-2324252]: this logic is currently the old logic adjusted. Let's clean it after new generation is working.
				// fill out DtosToGenerate; it replaces the old GenerateDtos and generateDtoDecls logic
				dtosToGenerate := make([]*Field, 0)
				dtosToGenerate, generatedDtos = addDtoToGenerate(o.OptsField, dtosToGenerate, generatedDtos, logger)
				logger.Printf("[DEBUG] Dtos to generate (length: %d): %v", len(dtosToGenerate), dtosToGenerate)
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

	if definition.checkStructShapes {
		requireConsistentStructShapes(definition)
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

func setParent(field *Field, logger *log.Logger) {
	for idx, f := range field.Fields {
		if f.Parent != nil {
			logger.Panicf("Field %s already has a parent\nold parent: %s (path: %s)\nnew parent: %s (path: %s);\n\nit is caused by the current incorrect implementation of nested fields;\nreuse the common definition by wrapping it in function invocation", f.Name, f.Parent.KindNoPtr(), f.Parent.PathWithRoot(), field.Name, field.PathWithRoot())
		}
		(&field.Fields[idx]).Parent = field
		setParent(&field.Fields[idx], logger)
	}
}

func addStructToGenerate(field *Field, structsToGenerate []*Field, generatedStructs []string, logger *log.Logger) ([]*Field, []string) {
	if !slices.Contains(generatedStructs, field.KindNoPtr()) {
		logger.Printf("[DEBUG] Adding %s (path: %s) to structs to be generated", field.KindNoPtr(), field.PathWithRoot())
		structsToGenerate = append(structsToGenerate, field)
		generatedStructs = append(generatedStructs, field.KindNoPtr())
	} else {
		logger.Printf("[DEBUG] Struct %s (path: %s) already queued for generation", field.KindNoPtr(), field.PathWithRoot())
	}

	for idx, f := range field.Fields {
		if f.IsShared {
			continue
		}
		if f.IsStruct() && !slices.Contains(generatedStructs, f.KindNoPtr()) {
			structsToGenerate, generatedStructs = addStructToGenerate(&field.Fields[idx], structsToGenerate, generatedStructs, logger)
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

func addDtoToGenerate(field *Field, dtosToGenerate []*Field, generatedDtos []string, logger *log.Logger) ([]*Field, []string) {
	if !slices.Contains(generatedDtos, field.DtoDecl()) {
		logger.Printf("[DEBUG] Adding %s (path: %s) to structs to be generated", field.DtoDecl(), field.PathWithRoot())
		dtosToGenerate = append(dtosToGenerate, field)
		generatedDtos = append(generatedDtos, field.DtoDecl())

		for idx, f := range field.Fields {
			if f.IsShared {
				continue
			}
			if f.HasAnyFields() {
				dtosToGenerate, generatedDtos = addDtoToGenerate(&field.Fields[idx], dtosToGenerate, generatedDtos, logger)
			}
		}
	} else {
		logger.Printf("[DEBUG] Struct %s (path: %s) already queued for generation", field.DtoDecl(), field.PathWithRoot())
	}
	return dtosToGenerate, generatedDtos
}

func ExtendInterface() func(*Interface, *genhelpers.PreambleModel, *log.Logger) *Interface {
	return func(i *Interface, preamble *genhelpers.PreambleModel, _ *log.Logger) *Interface {
		i.PreambleModel = preamble
		return i
	}
}

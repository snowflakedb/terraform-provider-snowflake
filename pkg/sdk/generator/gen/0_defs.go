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

			// TODO [SNOW-2324252]: this logic is currently the old logic adjusted. Let's clean it after new generation is working.
			// fill out StructsToGenerate; it replaces the old generateOptionsStruct and generateStruct
			structsToGenerate := make([]*Field, 0)
			for _, f := range o.HelperStructs {
				if !slices.Contains(generatedStructs, f.KindNoPtr()) {
					structsToGenerate, generatedStructs = addStructToGenerate(f, structsToGenerate, generatedStructs)
				}
			}
			for _, f := range o.OptsField.Fields {
				if len(f.Fields) > 0 && !slices.Contains(generatedStructs, f.KindNoPtr()) {
					structsToGenerate, generatedStructs = addStructToGenerate(f, structsToGenerate, generatedStructs)
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

func setParent(field *Field) {
	for _, f := range field.Fields {
		f.Parent = field
		setParent(f)
	}
}

func addStructToGenerate(field *Field, structsToGenerate []*Field, generatedStructs []string) ([]*Field, []string) {
	if !slices.Contains(generatedStructs, field.KindNoPtr()) {
		log.Printf("[DEBUG] Adding %s to structs to be generated", field.KindNoPtr())
		structsToGenerate = append(structsToGenerate, field)
		generatedStructs = append(generatedStructs, field.KindNoPtr())
	}

	for _, f := range field.Fields {
		if len(f.Fields) > 0 && !slices.Contains(generatedStructs, f.Name) {
			structsToGenerate, generatedStructs = addStructToGenerate(f, structsToGenerate, generatedStructs)
		}
	}
	return structsToGenerate, generatedStructs
}

func addDtoToGenerate(field *Field, dtosToGenerate []*Field, generatedDtos []string) ([]*Field, []string) {
	if !slices.Contains(generatedDtos, field.DtoDecl()) {
		log.Printf("[DEBUG] Adding %s to structs to be generated", field.DtoDecl())
		dtosToGenerate = append(dtosToGenerate, field)
		generatedDtos = append(generatedDtos, field.DtoDecl())

		for _, f := range field.Fields {
			if f.IsStruct() {
				dtosToGenerate, generatedDtos = addDtoToGenerate(f, dtosToGenerate, generatedDtos)
			}
		}
	}
	return dtosToGenerate, generatedDtos
}

func ExtendInterface() func(*Interface, *genhelpers.PreambleModel) *Interface {
	return func(i *Interface, preamble *genhelpers.PreambleModel) *Interface {
		i.PreambleModel = preamble
		return i
	}
}

// TODO [next PRs]: move these definitions
// var allSdkObjectDefinitions = []*generator.Interface{
//	sdk.ExternalVolumesDef,
//	sdk.NotificationIntegrationsDef,
//	sdk.OrganizationAccountsDef,
//	sdk.RowAccessPoliciesDef,
//	sdk.SecretsDef,
//	sdk.SecurityIntegrationsDef,
//	sdk.ServicesDef,
//	sdk.SessionPoliciesDef,
//	sdk.StagesDef,
//	sdk.StorageIntegrationDef,
//	sdk.StreamlitsDef,
//	sdk.StreamsDef,
//	sdk.TasksDef,
//	sdk.UserProgrammaticAccessTokensDef,
//	sdk.ViewsDef,
// }

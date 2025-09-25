package poc

import (
	"fmt"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

type SdkObjectDef struct {
	name string
	// TODO [next PR]: can be removed?
	file       string
	definition *generator.Interface
}

func GetSdkDefinitions() []*generator.Interface {
	allDefinitions := allSdkObjectDefinitions
	interfaces := make([]*generator.Interface, len(allDefinitions))
	for idx, s := range allDefinitions {
		preprocessDefinition(s.definition)
		interfaces[idx] = s.definition
	}
	return interfaces
}

// preprocessDefinition is needed because current simple builder is not ideal, should be removed later
func preprocessDefinition(definition *generator.Interface) {
	for _, o := range definition.Operations {
		o.ObjectInterface = definition
		if o.OptsField != nil {
			o.OptsField.Name = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
			o.OptsField.Kind = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
			setParent(o.OptsField)

			// TODO [this PR]: this logic is currently the old logic adjusted. Let's clean it after new generation is working.
			// fill out StructsToGenerate; it replaces the old generateOptionsStruct and generateStruct
			structsToGenerate := make([]*generator.Field, 0)
			generatedStructs := make([]string, 0)
			for _, f := range o.HelperStructs {
				if !slices.Contains(generatedStructs, f.KindNoPtr()) {
					structsToGenerate = addStructToGenerate(f, structsToGenerate, generatedStructs)
				}
			}
			for _, f := range o.OptsField.Fields {
				if len(f.Fields) > 0 && !slices.Contains(generatedStructs, f.KindNoPtr()) {
					structsToGenerate = addStructToGenerate(f, structsToGenerate, generatedStructs)
				}
			}
			// TODO [this PR]: replace with log or remove
			fmt.Printf("Structs to generate length: %d\n", len(structsToGenerate))
			o.StructsToGenerate = structsToGenerate
		}
	}
}

func setParent(field *generator.Field) {
	for _, f := range field.Fields {
		f.Parent = field
		setParent(f)
	}
}

func addStructToGenerate(field *generator.Field, structsToGenerate []*generator.Field, generatedStructs []string) []*generator.Field {
	if !slices.Contains(generatedStructs, field.KindNoPtr()) {
		// TODO [this PR]: replace with log or remove
		fmt.Printf("Adding: %s\n", field.KindNoPtr())
		structsToGenerate = append(structsToGenerate, field)
		generatedStructs = append(generatedStructs, field.KindNoPtr())
	}

	for _, f := range field.Fields {
		if len(f.Fields) > 0 && !slices.Contains(generatedStructs, f.Name) {
			structsToGenerate = addStructToGenerate(f, structsToGenerate, generatedStructs)
		}
	}
	return structsToGenerate
}

func WithPreamble(i *generator.Interface, preamble *genhelpers.PreambleModel) *generator.Interface {
	i.PreambleModel = preamble
	return i
}

var allSdkObjectDefinitions = []SdkObjectDef{
	{
		name:       "Sequences",
		file:       "sequences_def.go",
		definition: sdk.SequencesDef,
	},
}

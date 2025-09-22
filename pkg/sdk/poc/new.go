package poc

import (
	"fmt"

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
		}
	}
}

func setParent(field *generator.Field) {
	for _, f := range field.Fields {
		f.Parent = field
		setParent(f)
	}
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

package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// ParameterHookModel is one catalog parameter as the parameter hook templates see it.
// Conversion supplies the mapper and read expression. DiffFunc names the custom-diff helper.
type ParameterHookModel struct {
	Parameter  ssot.ParameterDef
	Classified ssot.ClassifiedField
	Conversion ValueConversionModel
}

func newParameterHookModel(parameter ssot.ParameterDef) ParameterHookModel {
	classified := ssot.ClassifyParameter(parameter.FieldName(), parameter.Kind)
	return ParameterHookModel{
		Parameter:  parameter,
		Classified: classified,
		Conversion: parameterValueConversion(parameter, classified),
	}
}

// NeedsMapping reports whether create and update pass a value mapper.
func (h ParameterHookModel) NeedsMapping() bool {
	return h.Conversion.NeedsDecoding()
}

// Mapper is the Go expression that converts the Terraform value to the SDK type.
// It is empty for bool, int, and string.
func (h ParameterHookModel) Mapper() string {
	return h.Conversion.DecodeFunc
}

// Read is the Go expression passed to setResourceData.
func (h ParameterHookModel) Read() string {
	value := "parameters." + h.Parameter.StructFieldName() + ".Value"
	return h.Conversion.StateExpression(value)
}

// DiffFunc is the custom-diff helper name for this parameter.
func (h ParameterHookModel) DiffFunc() string {
	switch h.Parameter.Kind {
	case string(ssot.AttrBool):
		return "BoolTypedParameterValueComputedIf"
	case string(ssot.AttrInt):
		return "IntTypedParameterValueComputedIf"
	case string(ssot.AttrFloat):
		return "FloatTypedParameterValueComputedIf"
	case "StringAllowEmpty":
		return "StringAllowEmptyTypedParameterValueComputedIf"
	case "AccountObjectIdentifier":
		return "IdentifierTypedParameterValueComputedIf"
	}
	if h.Classified.SdkKind == "identifier" {
		return "IdentifierTypedParameterValueComputedIf"
	}
	return "StringTypedParameterValueComputedIf"
}

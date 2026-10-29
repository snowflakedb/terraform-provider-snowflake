package gen

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// ValueConversionModel describes how one value moves between Terraform storage and its SDK type.
// DecodeFunc is empty when the stored value needs no conversion.
type ValueConversionModel struct {
	DecodeFunc string
	Encode     genhelpers.Mapper
}

// NeedsDecoding reports whether create and update pass DecodeFunc.
func (m ValueConversionModel) NeedsDecoding() bool {
	return m.DecodeFunc != ""
}

// StateExpression converts an SDK value expression into the Terraform state expression.
func (m ValueConversionModel) StateExpression(value string) string {
	return genhelpers.RunMapper(m.Encode, value)
}

func valueConversionFromClassified(field ssot.ClassifiedField) ValueConversionModel {
	switch field.SdkKind {
	case "bool":
		return ValueConversionModel{
			DecodeFunc: "booleanStringToBool",
			Encode: func(value string) string {
				return "booleanStringFromBool(" + value + ")"
			},
		}
	case "identifier":
		sdkType := strings.TrimPrefix(field.SdkType, "sdk.")
		return ValueConversionModel{
			DecodeFunc: "sdk.Parse" + sdkType,
			Encode:     genhelpers.FullyQualifiedName,
		}
	case "enum":
		sdkType := strings.TrimPrefix(field.SdkType, "sdk.")
		return ValueConversionModel{
			DecodeFunc: "sdk.To" + sdkType,
			Encode:     genhelpers.Identity,
		}
	default:
		return ValueConversionModel{Encode: genhelpers.Identity}
	}
}

func parameterValueConversion(parameter ssot.ParameterDef, classified ssot.ClassifiedField) ValueConversionModel {
	conversion := valueConversionFromClassified(classified)
	switch parameter.Kind {
	case "StringAllowEmpty":
		conversion.DecodeFunc = "stringToStringAllowEmpty"
	case string(ssot.AttrFloat):
		conversion.DecodeFunc = "sdk.ToFloat64"
		conversion.Encode = func(value string) string {
			return "formatFloatParameter(" + value + ")"
		}
	case "AccountObjectIdentifier":
		conversion.DecodeFunc = "stringToAccountObjectIdentifier"
	}
	return conversion
}

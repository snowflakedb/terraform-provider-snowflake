package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
)

// ParameterDetailField is one field of a generated <Object>ParametersDetails struct.
type ParameterDetailField struct {
	FieldName string // Go field name, e.g. "DataRetentionTimeInDays"
	Key       string // SQL name / SHOW PARAMETERS key, e.g. "DATA_RETENTION_TIME_IN_DAYS"
	GoType    string // typed value, e.g. "int", "LogLevel", "AccountObjectIdentifier"
	Parser    string // string->GoType parser, e.g. "strconv.Atoi", "ToLogLevel", "ParseAccountObjectIdentifier"
}

// ParametersDetailsConfig holds the data needed to generate the typed parameters-details accessor.
type ParametersDetailsConfig struct {
	Fields []ParameterDetailField
}

// ShowParametersDetails declares a typed ShowParametersDetails accessor for the given catalog parameters.
func (i *Interface) ShowParametersDetails(params ...parameterdefs.ParameterDef) *Interface {
	fields := make([]ParameterDetailField, 0, len(params))
	for _, p := range params {
		f := ParameterDetailField{FieldName: sqlToFieldName(p.SqlName, true), Key: p.SqlName}
		if p.Kind == KindOfT[sdkcommons.StringAllowEmpty]() {
			f.GoType, f.Parser = "string", "identityParse"
			fields = append(fields, f)
			continue
		}
		switch p.Kind {
		case KindBool:
			f.GoType, f.Parser = "bool", "strconv.ParseBool"
		case KindInt:
			f.GoType, f.Parser = "int", "strconv.Atoi"
		case KindString:
			f.GoType, f.Parser = "string", "identityParse"
		default:
			f.GoType = p.Kind
			if _, err := ToObjectIdentifierKind(p.Kind); err == nil {
				f.Parser = "Parse" + p.Kind
			} else {
				f.Parser = "To" + p.Kind
			}
		}
		fields = append(fields, f)
	}
	i.ParametersDetails = &ParametersDetailsConfig{Fields: fields}
	return i.WithCustomInterfaceMethod(
		"ShowParametersDetails",
		"",
		[]*MethodParameter{NewMethodParameter("id", i.IdentifierKind)},
		"*"+i.NameSingular+"ParametersDetails", "error",
	)
}

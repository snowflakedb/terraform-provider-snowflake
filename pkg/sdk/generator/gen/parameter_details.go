package gen

import (
	"fmt"

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
	Fields            []ParameterDetailField
	WithoutIdentifier bool
}

// ShowParametersDetails declares a typed ShowParametersDetails accessor for the given catalog parameters.
func (i *Interface) ShowParametersDetails(params ...parameterdefs.ParameterDef) *Interface {
	return i.showParametersDetails(false, params...)
}

// ShowParametersDetailsWithoutIdentifier declares a typed ShowParametersDetails accessor for
// account-scoped SHOW PARAMETERS, which does not accept an object identifier.
func (i *Interface) ShowParametersDetailsWithoutIdentifier(params ...parameterdefs.ParameterDef) *Interface {
	return i.showParametersDetails(true, params...)
}

func (i *Interface) showParametersDetails(withoutIdentifier bool, params ...parameterdefs.ParameterDef) *Interface {
	fields := make([]ParameterDetailField, 0, len(params))
	for _, p := range params {
		info, err := InfoForKind(p.Kind)
		if err != nil {
			panic(fmt.Sprintf("parameter %s: %v", p.SqlName, err))
		}
		fields = append(fields, ParameterDetailField{
			FieldName: sqlToFieldName(p.SqlName, true),
			Key:       p.SqlName,
			GoType:    info.GoType,
			Parser:    info.ReadParser,
		})
	}
	i.ParametersDetails = &ParametersDetailsConfig{Fields: fields, WithoutIdentifier: withoutIdentifier}
	parameters := []*MethodParameter{NewMethodParameter("id", i.IdentifierKind)}
	if withoutIdentifier {
		parameters = nil
	}
	return i.WithCustomInterfaceMethod(
		"ShowParametersDetails",
		"",
		parameters,
		"*"+i.NameSingular+"ParametersDetails", "error",
	)
}

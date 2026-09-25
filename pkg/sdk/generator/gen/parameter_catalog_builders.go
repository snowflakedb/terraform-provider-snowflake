package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
)

// WithParameters expands catalog parameters into CREATE/ALTER-SET assignment fields. String and enum
// values are single-quoted — Snowflake accepts quoted and unquoted forms uniformly for all
// string/enum parameters.
func (v *QueryStruct) WithParameters(params ...parameterdefs.ParameterDef) *QueryStruct {
	for _, p := range params {
		if p.Kind == KindOfT[sdkcommons.StringAllowEmpty]() {
			v.OptionalAssignment(p.SqlName, p.Kind, ParameterOptions())
			continue
		}
		switch p.Kind {
		case KindBool:
			v.OptionalBooleanAssignment(p.SqlName, nil)
		case KindInt:
			v.OptionalNumberAssignment(p.SqlName, ParameterOptions())
		case KindString:
			v.OptionalTextAssignment(p.SqlName, ParameterOptions().SingleQuotes())
		default:
			if _, err := ToObjectIdentifierKind(p.Kind); err == nil {
				v.OptionalIdentifier(sqlToFieldName(p.SqlName, true), p.Kind, IdentifierOptions().SQL(p.SqlName).Equals())
			} else {
				v.OptionalAssignment(p.SqlName, p.Kind, ParameterOptions().SingleQuotes())
			}
		}
	}
	return v
}

// WithParametersUnset expands catalog parameters into ALTER-UNSET keyword fields.
func (v *QueryStruct) WithParametersUnset(params ...parameterdefs.ParameterDef) *QueryStruct {
	for _, p := range params {
		v.OptionalSQL(p.SqlName)
	}
	return v
}

package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
)

// WithParameters expands catalog parameters into CREATE/ALTER-SET assignment fields.
func (v *QueryStruct) WithParameters(params ...parameterdefs.ParameterDef) *QueryStruct {
	for _, p := range params {
		v.setParameters = append(v.setParameters, newParameterField(p))
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
		case KindFloat:
			v.OptionalAssignment(p.SqlName, KindFloat, ParameterOptions())
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
		v.unsetParameters = append(v.unsetParameters, newParameterField(p))
		v.OptionalSQL(p.SqlName)
	}
	return v
}

func newParameterField(p parameterdefs.ParameterDef) ParameterField {
	info, err := InfoForKind(p.Kind)
	if err != nil {
		panic(fmt.Sprintf("parameter %s: %v", p.SqlName, err))
	}
	return ParameterField{
		SqlName:   p.SqlName,
		FieldName: sqlToFieldName(p.SqlName, true),
		Parser:    info.WriteParser,
	}
}

// ParameterSqlToFieldName maps a catalog parameter's SQL name to the exported Go field
// name the generator uses for it. Exported so definitions can derive validation
// field-name lists that are guaranteed to match generated fields.
func ParameterSqlToFieldName(p parameterdefs.ParameterDef) string {
	return sqlToFieldName(p.SqlName, true)
}

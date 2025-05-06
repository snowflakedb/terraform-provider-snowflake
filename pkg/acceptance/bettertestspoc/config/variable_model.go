package config

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

type VariableModel struct {
	Type    tfconfig.Variable `json:"type,omitempty"`
	Default tfconfig.Variable `json:"default,omitempty"`

	name string
}

func (v *VariableModel) CommonTfName() string {
	return v.name
}

func (v *VariableModel) CommonTfType() string {
	return "variable"
}

func Variable(
	variableName string,
	type_ string,
) *VariableModel {
	v := &VariableModel{
		name: variableName,
	}
	v.WithType(type_)
	return v
}

func StringVariable(
	variableName string,
) *VariableModel {
	return Variable(variableName, "string")
}

func NumberVariable(
	variableName string,
) *VariableModel {
	return Variable(variableName, "number")
}

func SetMapStringVariable(
	variableName string,
) *VariableModel {
	return Variable(variableName, "set(map(string))")
}

func (v *VariableModel) WithType(type_ string) *VariableModel {
	v.Type = UnquotedWrapperVariable(type_)
	return v
}

func (v *VariableModel) WithStringDefault(default_ string) *VariableModel {
	v.Default = tfconfig.StringVariable(default_)
	return v
}

func (v *VariableModel) WithUnquotedDefault(default_ string) *VariableModel {
	v.Default = UnquotedWrapperVariable(default_)
	return v
}

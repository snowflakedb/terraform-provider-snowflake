package config

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

// TODO [this PR]: use types instead of tfconfig.Variable?
type VariableModel struct {
	Type      tfconfig.Variable `json:"type,omitempty"`
	Default   tfconfig.Variable `json:"default,omitempty"`
	Sensitive tfconfig.Variable `json:"sensitive,omitempty"`

	name string
}

func (v *VariableModel) BlockName() string {
	return v.name
}

func (v *VariableModel) BlockType() string {
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

func (v *VariableModel) WithDefault(variable tfconfig.Variable) *VariableModel {
	v.Default = variable
	return v
}

func (v *VariableModel) WithUnquotedDefault(default_ string) *VariableModel {
	v.Default = UnquotedWrapperVariable(default_)
	return v
}

func (v *VariableModel) WithSensitive(sensitive bool) *VariableModel {
	v.Sensitive = tfconfig.BoolVariable(sensitive)
	return v
}

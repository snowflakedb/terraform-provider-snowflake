package config

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

// VariableModel allows to use predefined variable block.
// Check: https://developer.hashicorp.com/terraform/language/values/variables#declaring-an-input-variable.
// TODO [SNOW-1501905]: Consider using types instead of tfconfig.Variable
//   - The reason to use tfconfig.Variable in the earlier models was to be able to put any value to also test invalid configs.
//   - Maybe there is no sense to do that for all the fields in the Terraform predefined blocks like locals, variables, etc.
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

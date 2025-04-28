package config

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/config"
)

type DynamicBlock map[string]DynamicBlockContent

type DynamicBlockContent struct {
	ForEach string            `json:"for_each"`
	Content map[string]string `json:"content"`
}

// NewDynamicBlock is a quick and dirty implementation to add dynamic blocks to our config builders.
// Dynamic blocks look like this:
//
//	dynamic "<label_name>" {
//		for_each = var.<variable_name>
//		content {
//			name = argument.value["name"]
//			type = argument.value["type"]
//		}
//	}
//
// which in JSON would look like (<...> mark the fields with dynamic names):
//
//	{
//		"dynamic": {
//			"<label_name>": {
//				"for_each": "var.<variable_name>",
//				"content": {
//					"<arg1_name>": "<label_name>.value[\"<arg1_name>\"]",
//					"<arg2_name>": "<label_name>.value[\"<arg2_name>\"]"
//					...
//					"<argN_name>": "<label_name>.value[\"<argN_name>\"]"
//				}
//			}
//		}
//	}
//
// The main complexity with our struct -> json -> hcl -> formatted hcl approach is that in JSON all string values are quoted.
// We need to unquote them and unescape the escaped quotes inside the quoted value.
// We use the SnowflakeProviderConfigUnquoteMarker and SnowflakeProviderConfigQuoteMarker to apply this formatting.
func NewDynamicBlock(label string, variableName string, values []string) *DynamicBlock {
	args := make(map[string]string)
	for _, v := range values {
		quotedValue := fmt.Sprintf(`%[2]s%[1]s%[2]s`, v, SnowflakeProviderConfigQuoteMarker)
		argumentReference := fmt.Sprintf(`%[1]s.value[%[2]s]`, label, quotedValue)
		unquotedArgumentReference := fmt.Sprintf(`%[2]s%[1]s%[2]s`, argumentReference, SnowflakeProviderConfigUnquoteMarker)
		args[v] = unquotedArgumentReference
	}

	return &DynamicBlock{
		label: {
			ForEach: fmt.Sprintf(`%[2]svar.%[1]s%[2]s`, variableName, SnowflakeProviderConfigUnquoteMarker),
			Content: args,
		},
	}
}

func TempSecretStringVariableConfig(variableName string, value string) (string, config.Variables) {
	temporaryVariableDefinition := fmt.Sprintf(`
	variable "%s" {
		type = string
		sensitive = true
	}
`, variableName)
	configVariables := config.Variables{
		variableName: config.StringVariable(value),
	}

	return temporaryVariableDefinition, configVariables
}

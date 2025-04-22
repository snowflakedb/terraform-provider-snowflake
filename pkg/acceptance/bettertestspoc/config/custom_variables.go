package config

import (
	"encoding/json"
	"fmt"
)

type emptyListVariable struct{}

// MarshalJSON returns the JSON encoding of emptyListVariable.
func (v emptyListVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{})
}

// EmptyListVariable returns Variable representing an empty list. This is because the current hcl parser handles empty SetVariable incorrectly.
func EmptyListVariable() emptyListVariable {
	return emptyListVariable{}
}

type replacementPlaceholderVariable struct {
	placeholder ReplacementPlaceholder
}

// MarshalJSON returns the JSON encoding of replacementPlaceholderVariable.
func (v replacementPlaceholderVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.placeholder)
}

// ReplacementPlaceholderVariable returns Variable containing one of the ReplacementPlaceholder which is later replaced by HclFormatter.
func ReplacementPlaceholderVariable(placeholder ReplacementPlaceholder) replacementPlaceholderVariable {
	return replacementPlaceholderVariable{placeholder}
}

type multilineWrapperVariable struct {
	content string
}

// MarshalJSON returns the JSON encoding of multilineWrapperVariable.
func (v multilineWrapperVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf(`%[1]s%[2]s%[1]s`, SnowflakeProviderConfigMultilineMarker, v.content))
}

// MultilineWrapperVariable returns Variable containing multiline content wrapped with SnowflakeProviderConfigMultilineMarker later replaced by HclFormatter.
func MultilineWrapperVariable(content string) multilineWrapperVariable {
	return multilineWrapperVariable{content}
}

type secretVariable struct {
	content string
}

// MarshalJSON returns the JSON encoding of multilineWrapperVariable.
func (v secretVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf(`%[1]s%[2]s%[1]s`, SnowflakeProviderConfigMultilineMarker, v.content))
}

// SecretVariable returns Variable containing multiline content wrapped with SnowflakeProviderConfigMultilineMarker later replaced by HclFormatter.
func SecretVariable(content string) secretVariable {
	return secretVariable{content}
}

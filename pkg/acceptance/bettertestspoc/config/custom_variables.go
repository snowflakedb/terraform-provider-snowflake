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

type wrapperVariable struct {
	wrapper ReplacementPlaceholder
	content string
}

// MarshalJSON returns the JSON encoding of multilineWrapperVariable.
func (v wrapperVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf(`%[1]s%[2]s%[1]s`, v.wrapper, v.content))
}

// MultilineWrapperVariable returns Variable containing multiline content wrapped with SnowflakeProviderConfigMultilineMarker later replaced by HclFormatter.
func MultilineWrapperVariable(content string) wrapperVariable {
	return wrapperVariable{SnowflakeProviderConfigMultilineMarker, content}
}

// QuotedWrapperVariable returns Variable containing quoted content wrapped with SnowflakeProviderConfigQuoteMarker later replaced by HclFormatter.
func QuotedWrapperVariable(content string) wrapperVariable {
	return wrapperVariable{SnowflakeProviderConfigQuoteMarker, content}
}

// UnquotedWrapperVariable returns Variable containing unquoted content wrapped with SnowflakeProviderConfigUnquoteMarker later replaced by HclFormatter.
func UnquotedWrapperVariable(content string) wrapperVariable {
	return wrapperVariable{SnowflakeProviderConfigUnquoteMarker, content}
}

func VariableReference(variableName string) wrapperVariable {
	return UnquotedWrapperVariable(fmt.Sprintf("var.%s", variableName))
}

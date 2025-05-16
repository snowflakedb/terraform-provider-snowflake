package helpers

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

const ResourceIdDelimiter = '|'

func ParseResourceIdentifier(identifier string) []string {
	if identifier == "" {
		return make([]string, 0)
	}
	return strings.Split(identifier, string(ResourceIdDelimiter))
}

func EncodeResourceIdentifier[T sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.SchemaObjectIdentifierWithArguments | sdk.TableColumnIdentifier | sdk.ExternalObjectIdentifier | sdk.AccountIdentifier | string](parts ...T) string {
	result := make([]string, len(parts))
	for i, part := range parts {
		switch typedPart := any(part).(type) {
		case sdk.AccountObjectIdentifier:
			result[i] = typedPart.Name()
		case sdk.DatabaseObjectIdentifier, sdk.SchemaObjectIdentifier, sdk.SchemaObjectIdentifierWithArguments, sdk.TableColumnIdentifier, sdk.ExternalObjectIdentifier, sdk.AccountIdentifier:
			result[i] = typedPart.(sdk.ObjectIdentifier).FullyQualifiedName()
		case string:
			result[i] = typedPart
		}
	}
	return strings.Join(result, string(ResourceIdDelimiter))
}

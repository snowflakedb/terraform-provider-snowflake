package genhelpers

// sdkIdentifierTypes is the set of known SDK identifier type names (without package prefix or pointer).
var sdkIdentifierTypes = map[string]struct{}{
	"AccountObjectIdentifier":             {},
	"DatabaseObjectIdentifier":            {},
	"SchemaObjectIdentifier":              {},
	"SchemaObjectIdentifierWithArguments": {},
	"ExternalObjectIdentifier":            {},
	"AccountIdentifier":                   {},
	"TableColumnIdentifier":               {},
}

// IsIdentifierType reports whether kind (possibly pointer-prefixed) is a known SDK identifier type.
func IsIdentifierType(kind string) bool {
	_, ok := sdkIdentifierTypes[TypeWithoutPointer(kind)]
	return ok
}

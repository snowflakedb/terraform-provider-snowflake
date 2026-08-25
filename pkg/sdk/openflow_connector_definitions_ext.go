package sdk

import "strings"

// ParseOpenflowConnectorDefinitionCategories parses the categories column of
// SHOW OPENFLOW CONNECTOR DEFINITIONS, which Snowflake returns as a JSON array such as
// ["Databases","Analytics"]. It exists only to adapt the shared parser to the signature the generator
// expects for a custom parser, which is (string) (T, error).
func ParseOpenflowConnectorDefinitionCategories(value string) ([]string, error) {
	// A column with no values means no categories. Snowflake spells that as SQL NULL, the literal string
	// null, or an empty array depending on how the row got there; the shared parser would turn the first
	// into one empty element and the second into a single category named "null". The runtime's
	// external_access_integrations column has the same three shapes.
	switch strings.TrimSpace(value) {
	case "", "null", "[]":
		return []string{}, nil
	}
	return ParseCommaSeparatedStringArray(value, true), nil
}

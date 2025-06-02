package genhelpers

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var (
	splitOnTheWordsBeginnings = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	splitRemainingWordBreaks  = regexp.MustCompile("([a-z0-9])([A-Z]+)")

	forbiddenAttributeNames = []string{"type"}
)

const forbiddenAttributeNameSuffix = "_"

var PredefinedImports = map[string]string{
	"sdk":         "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk",
	"collections": "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections",
}

// ToSnakeCase allows converting a CamelCase text to camel_case one (needed for schema attribute names). Examples:
// - CamelCase -> camel_case
// - ACamelCase -> a_camel_case
// - URLParser -> url_parser
// - Camel1Case -> camel1_case
// - camelCase -> camel_case
// - camelURL -> camel_url
// - camelURLSomething -> camel_url_something
func ToSnakeCase(str string) string {
	wordsSplit := splitOnTheWordsBeginnings.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(splitRemainingWordBreaks.ReplaceAllString(wordsSplit, "${1}_${2}"))
}

// ColumnOutput is a helper to align a tabular output with the given columnWidth, e.g. (for 20 spaces column width):
// Name                string              string
// State               sdk.WarehouseState  string
func ColumnOutput(columnWidth int, columns ...string) string {
	var sb strings.Builder
	for i, column := range columns {
		sb.WriteString(column)
		if i != len(columns)-1 {
			spaces := max(columnWidth-len(column), 1)
			sb.WriteString(strings.Repeat(" ", spaces))
		}
	}
	return sb.String()
}

// WriteCodeToFile formats and saves content from the given buffer into file relative to the current working directory.
// TODO [SNOW-1501905]: test
func WriteCodeToFile(buffer *bytes.Buffer, fileName string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("writing code to file %s failed with err: %w", fileName, err)
	}
	outputPath := filepath.Join(wd, fileName)
	src, err := format.Source(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("writing code to file %s failed with err: %w", fileName, err)
	}
	if err := os.WriteFile(outputPath, src, 0o600); err != nil {
		return fmt.Errorf("writing code to file %s failed with err: %w", fileName, err)
	}
	return nil
}

func SanitizeAttributeName(name string) string {
	if slices.Contains(forbiddenAttributeNames, name) {
		return name + forbiddenAttributeNameSuffix
	}
	return name
}

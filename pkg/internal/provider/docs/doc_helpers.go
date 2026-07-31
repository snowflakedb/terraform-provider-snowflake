package docs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// Replacement describes a single resource/datasource that replaces a deprecated one.
type Replacement struct {
	// Name is the full object name, e.g. `snowflake_storage_integration_aws`.
	Name string
}

// Page returns the name of the object's documentation page, i.e. Name without the `snowflake_` prefix.
func (r Replacement) Page() string {
	return strings.TrimPrefix(r.Name, "snowflake_")
}

// Quoted returns the replacement exactly as it appears in the deprecation message (check deprecatedResourceDescription).
// Replacing the quoted form (instead of just Name) keeps the replacements delimited, so that one of them can not be
// matched inside another one it is a prefix of, and leaves no quotes around the resulting link.
func (r Replacement) Quoted() string {
	return fmt.Sprintf("`%s`", r.Name)
}

// listedReplacementsRegex matches the replacements listed in the deprecation message, e.g.
// "Please use one of the new resources instead: `snowflake_storage_integration_aws` | `snowflake_storage_integration_azure`."
var listedReplacementsRegex = regexp.MustCompile(`Please use one of the new (?:resources|data sources) instead: ([^.]+)\.`)

// quotedObjectNameRegex matches a single quoted object name from the list of replacements.
var quotedObjectNameRegex = regexp.MustCompile("`(snowflake_\\w+)`")

// GetDeprecatedObjectReplacements returns all replacements listed in the given resource/datasource DeprecationMessage.
// It returns nil when the message does not list any (i.e. it was not built with deprecatedResourceDescription).
func GetDeprecatedObjectReplacements(deprecationMessage string) []Replacement {
	listMatch := listedReplacementsRegex.FindStringSubmatch(deprecationMessage)
	if listMatch == nil {
		return nil
	}
	return collections.Map(quotedObjectNameRegex.FindAllStringSubmatch(listMatch[1], -1), func(nameMatch []string) Replacement {
		return Replacement{Name: nameMatch[1]}
	})
}

// RelativeLink allows us to get relative link to the resource/datasource in the same subtree. Will have to change when we introduce subcategories.
func RelativeLink(title string, path string) string {
	return fmt.Sprintf(`[%s](./%s)`, title, path)
}

func PossibleValuesListed[T ~string | ~int](values []T) string {
	valuesWrapped := make([]string, len(values))
	for i, value := range values {
		valuesWrapped[i] = fmt.Sprintf("`%v`", value)
	}
	return strings.Join(valuesWrapped, " | ")
}

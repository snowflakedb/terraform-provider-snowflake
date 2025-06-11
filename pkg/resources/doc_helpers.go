package resources

import (
	"fmt"
	"strings"

	providerresources "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
)

func possibleValuesListed[T ~string | ~int](values []T) string {
	return docs.PossibleValuesListed(values)
}

func characterList(values []rune) string {
	valuesWrapped := make([]string, len(values))
	for i, value := range values {
		valuesWrapped[i] = fmt.Sprintf("`%c`", value)
	}
	return strings.Join(valuesWrapped, ", ")
}

func booleanStringFieldDescription(description string) string {
	return fmt.Sprintf(`%s Available options are: "%s" or "%s". When the value is not set in the configuration the provider will put "%s" there which means to use the Snowflake default for this value.`, description, BooleanTrue, BooleanFalse, BooleanDefault)
}

func externalChangesNotDetectedFieldDescription(description string) string {
	return fmt.Sprintf(`%s External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".`, description)
}

func withPrivilegedRolesDescription(description, paramName string) string {
	return fmt.Sprintf(`%s By default, this list includes the ACCOUNTADMIN, ORGADMIN and SECURITYADMIN roles. To remove these privileged roles from the list, use the ALTER ACCOUNT command to set the %s account parameter to FALSE.`, description, paramName)
}

func blocklistedCharactersFieldDescription(description string) string {
	return fmt.Sprintf(`%s Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: %s.`, description, characterList([]rune{'|', '.', '"'}))
}

func diffSuppressStatementFieldDescription(description string) string {
	return fmt.Sprintf(`%s To mitigate permadiff on this field, the provider replaces blank characters with a space. This can lead to false positives in cases where a change in case or run of whitespace is semantically significant.`, description)
}

func dataTypeFieldDescription(description string) string {
	return fmt.Sprintf(`%s For more information about data types, check [Snowflake docs](https://docs.snowflake.com/en/sql-reference/intro-summary-data-types).`, description)
}

func deprecatedResourceDescription(alternatives ...string) string {
	return fmt.Sprintf(`This resource is deprecated and will be removed in a future major version release. Please use one of the new resources instead: %s.`, possibleValuesListed(alternatives))
}

func copyGrantsDescription(description string) string {
	return fmt.Sprintf("%s This is used when the provider detects changes for fields that can not be changed by ALTER. This value will not have any effect during creating a new object with Terraform.", description)
}

func relatedResourceDescription(description string, resource providerresources.Resource) string {
	return fmt.Sprintf(`%s For more information about this resource, see [docs](./%s).`, description, strings.TrimPrefix(resource.String(), "snowflake_"))
}

func joinWithSpace(parts ...string) string {
	return strings.Join(parts, " ")
}

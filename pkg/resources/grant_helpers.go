package resources

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// toPrivileges validates and normalizes a list of privilege names using sdk.ToPrivilege.
// It is used when parsing untrusted input (e.g. resource identifiers during import).
func toPrivileges(privileges []string) ([]string, error) {
	return collections.MapErr(privileges, sdk.ToPrivilege)
}

func isNotOwnershipGrant() func(value any, path cty.Path) diag.Diagnostics {
	return func(value any, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		if privilege, ok := value.(string); ok && strings.ToUpper(privilege) == "OWNERSHIP" {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Unsupported privilege 'OWNERSHIP'",
				Detail:        "Granting ownership is only allowed in snowflake_grant_ownership resource.",
				AttributePath: nil,
			})
		}
		return diags
	}
}

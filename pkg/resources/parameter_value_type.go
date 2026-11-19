package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/validators"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// parameterValueTypeSchema is the value_type field of the parameter resources. supportedListLocation tells where
// the resource lists the parameters that the provider supports, and experiment is the one that turns the field on.
func parameterValueTypeSchema(supportedListLocation string, experiment experimentalfeatures.ExperimentalFeature, diffSuppress schema.SchemaDiffSuppressFunc) *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validators.NormalizeValidation(sdk.ToParameterValueKind),
		DiffSuppressFunc: SuppressIfAny(diffSuppress, NormalizeAndCompare(sdk.ToParameterValueKind)),
		Description: fmt.Sprintf("Type of `value`, used only for parameters that the provider does not support yet (forward compatibility). "+
			"Valid values are (case-insensitive): %[1]s. `value` is rendered the same way as for the parameters that the provider supports: "+
			"`STRING` (also use it for enum parameters) is sent as a single-quoted string, `IDENTIFIER` as a quoted identifier (it can be fully qualified, e.g. `DB.SCHEMA.NAME`), "+
			"and the others as a number or a boolean. The field is required for a parameter that the provider does not support, and it is ignored for the supported parameters "+
			"(they are listed in %[2]s), so it can stay in the configuration when a later provider version starts to support the parameter. Changes of this field are ignored for the supported parameters, "+
			"but setting it for a supported parameter logs a warning. The field works only with the %[3]s experiment (enabled by default). "+
			"The transition does not need configuration changes only if `value_type` matches the type of the parameter in Snowflake. Some types are easy to confuse: "+
			"for example, an identifier (like a network policy name) is also accepted by Snowflake as a `STRING`, but once the provider supports the parameter, it is sent as a quoted, case-sensitive identifier, "+
			"which may show a plan difference or fail. Check the parameter type in the [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/parameters) before choosing `value_type`.",
			possibleValuesListed(sdk.AsStringList(sdk.AllParameterValueKinds)), supportedListLocation, experiment),
	}
}

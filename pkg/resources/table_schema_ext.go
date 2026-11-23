package resources

import (
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-1031688]: move data manipulation logic to the SDK - SQL generation or builders part (e.g. different default types/identity)

// tableAttributesModificationExt restores schema chrome the generator does not emit.
func tableAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	column := s["column"].Elem.(*schema.Resource).Schema
	column["type"].ValidateDiagFunc = IsDataTypeValid
	column["type"].DiffSuppressFunc = DiffSuppressDataTypes
	column["nullable"].Default = true
	column["comment"].Default = ""
	column["masking_policy"].Default = ""
	column["collate"].Default = ""
	// A column without an explicit collation inherits the table's default_ddl_collation at the moment it is created,
	// and DESCRIBE keeps reporting it after the table default is changed or unset. An omitted collate is therefore
	// treated as "not managed", so it must not plan a change to the collation Snowflake assigned. An explicit
	// `collate = ""` is managed: it is a change that clears the collation, so it is not suppressed.
	column["collate"].DiffSuppressFunc = func(k, old, new string, d *schema.ResourceData) bool {
		return new == "" && old != "" && isColumnCollateOmitted(k, d)
	}

	defaultExpression := column["default"].Elem.(*schema.Resource).Schema["expression"]
	defaultExpression.ValidateDiagFunc = validateTableColumnDefaultExpression
	defaultExpression.DiffSuppressFunc = suppressColumnDefaultExpressionDifferences

	identity := column["identity"].Elem.(*schema.Resource).Schema
	identity["start_num"].Default = 1
	identity["step_num"].Default = 1

	s["primary_key"].Deprecated = "Use snowflake_table_constraint instead"

	// Identifier validation is not generated for nested-block fields (identifierSchemaType skips InBlock).
	for _, key := range []string{"row_access_policy", "aggregation_policy", "join_policy"} {
		s[key].Elem.(*schema.Resource).Schema["policy_name"].ValidateDiagFunc = IsValidIdentifier[sdk.SchemaObjectIdentifier]()
	}

	return s
}

// isColumnCollateOmitted reports whether the collate of the column under the given schema key (column.<index>.collate)
// is absent from the configuration. When the configuration cannot be read, collate is treated as omitted,
// which keeps the previous behavior of not planning a change to the collation Snowflake assigned.
func isColumnCollateOmitted(key string, d *schema.ResourceData) bool {
	parts := strings.Split(key, ".")
	if len(parts) != 3 || parts[0] != "column" || parts[2] != "collate" {
		return true
	}
	index, err := strconv.Atoi(parts[1])
	if err != nil {
		return true
	}
	value, diags := d.GetRawConfigAt(cty.GetAttrPath("column").IndexInt(index).GetAttr("collate"))
	if diags.HasError() {
		return true
	}
	return value.IsNull()
}

// validateTableColumnDefaultExpression rejects an expression that is only whitespace. An empty expression ("") is
// valid and means an empty string default (see getColumnDefaultExpression), but "  " is most likely a mistake and
// would not be read back from DESCRIBE TABLE as the same value.
func validateTableColumnDefaultExpression(value any, path cty.Path) diag.Diagnostics {
	if expression, ok := value.(string); ok && expression != "" && strings.TrimSpace(expression) == "" {
		return diag.Diagnostics{{
			Severity:      diag.Error,
			Summary:       "Invalid column default expression",
			Detail:        "The default expression cannot be only whitespace. Use \"\" or \"''\" for an empty string default.",
			AttributePath: path,
		}}
	}
	return nil
}

// suppressColumnDefaultExpressionDifferences hides the differences between the configured default expression and
// the one DESCRIBE TABLE returns that do not change the default:
//   - an empty expression ("") is created as the empty string literal, which DESCRIBE TABLE returns as two single quotes;
//   - the quoting of the sequence in <sequence>.NEXTVAL (see suppressColumnDefaultSequenceQuoting).
func suppressColumnDefaultExpressionDifferences(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if isEmptyStringDefault(oldValue) && isEmptyStringDefault(newValue) {
		return true
	}
	return suppressColumnDefaultSequenceQuoting(k, oldValue, newValue, d)
}

func isEmptyStringDefault(expression string) bool {
	return expression == "" || expression == columnDefaultEmptyStringLiteral
}

// suppressColumnDefaultSequenceQuoting compares <sequence>.NEXTVAL defaults by the sequence identifier, because
// DESCRIBE TABLE returns the sequence name without the double quotes a fully qualified name has.
func suppressColumnDefaultSequenceQuoting(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	oldSequence, oldIsSequence := cutNextval(oldValue)
	newSequence, newIsSequence := cutNextval(newValue)
	if !oldIsSequence || !newIsSequence {
		return false
	}
	oldId, err := sdk.ParseSchemaObjectIdentifier(oldSequence)
	if err != nil {
		return false
	}
	newId, err := sdk.ParseSchemaObjectIdentifier(newSequence)
	if err != nil {
		return false
	}
	return oldId.FullyQualifiedName() == newId.FullyQualifiedName()
}

func cutNextval(expression string) (string, bool) {
	const nextval = ".NEXTVAL"
	expression = strings.TrimSpace(expression)
	if len(expression) <= len(nextval) || !strings.EqualFold(expression[len(expression)-len(nextval):], nextval) {
		return "", false
	}
	return expression[:len(expression)-len(nextval)], true
}

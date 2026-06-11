package sdk

import (
	"context"
	"regexp"
	"strings"
)

func (r hybridTableDetailsRow) additionalConvert(result *HybridTableDetails) error {
	type_, collation := r.splitTypeAndCollation()
	result.Type = type_
	result.Collation = collation
	return nil
}

// ShowParameters returns the parameters visible at the TABLE level for the given hybrid table.
// Mirrors pkg/sdk/functions_ext.go:155 (ParametersIn.Function) with ParametersIn.Table.
func (v *hybridTables) ShowParameters(ctx context.Context, id SchemaObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Table: id,
		},
	})
}

// splitTypeAndCollation extracts the COLLATE 'X' suffix that DESCRIBE TABLE returns for
// collated columns (e.g. "VARCHAR(200) COLLATE 'en-ci'") and returns the bare type and
// the collation as separate values. Mirrors pkg/sdk/tables.go:736 — the same shape exists
// for classic tables but the generator cannot derive a Collation field from the raw Type
// column, so this lives in the _ext.go file and is invoked manually from convert().
func (r hybridTableDetailsRow) splitTypeAndCollation() (string, *string) {
	collateRegexp := regexp.MustCompile(`COLLATE +'([a-zA-Z0-9_-]*)'`)
	matches := collateRegexp.FindStringSubmatch(r.Type)
	if len(matches) == 2 {
		collation := matches[1]
		type_ := strings.TrimSpace(collateRegexp.ReplaceAllString(r.Type, ""))
		return type_, &collation
	}
	return r.Type, nil
}

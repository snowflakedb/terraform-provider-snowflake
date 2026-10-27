package gen

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// ConfigFieldSet is one d.Set inside the generated SetConfigFields function.
type ConfigFieldSet struct {
	Name string
	Expr string
}

// ConfigFieldSets is the SHOW values copied into state.
// A field is included when it is a resource attribute in ShowColumns, not in the
// set-state list, not CreateOnly, and not a nested block. Names follow Attributes
// order. Show membership matches Name, because Attributes is a copy of the defs.
func (m ResourceModel) ConfigFieldSets() []ConfigFieldSet {
	excluded := map[string]struct{}{}
	for _, name := range m.SetStateToValueFromConfigFields() {
		excluded[name] = struct{}{}
	}
	columns := map[string]ssot.Attribute{}
	for _, col := range m.ShowColumns {
		if col == nil {
			continue
		}
		columns[col.Name] = *col
	}
	sets := make([]ConfigFieldSet, 0)
	for _, attr := range m.Attributes {
		col, ok := columns[attr.Name]
		if !ok || col.CreateOnly || len(attr.Children) > 0 || len(col.Children) > 0 || strings.Contains(attr.DiffSuppressExpr, "IgnoreAfterCreation") {
			continue
		}
		if _, skip := excluded[attr.Name]; skip {
			continue
		}
		sets = append(sets, ConfigFieldSet{Name: attr.Name, Expr: configFieldCode(ssot.Object{Name: m.Name, IdType: m.IdType}, col)})
	}
	if len(sets) == 0 {
		return nil
	}
	return sets
}

func configFieldCode(obj ssot.Object, col ssot.Attribute) string {
	field := "obj." + genhelpers.SnakeCaseToCamel(col.OutputKey())
	if col.OptionalInShow {
		return field
	}
	return valueConversionFromClassified(ssot.ClassifyAttribute(obj, col)).StateExpression(field)
}

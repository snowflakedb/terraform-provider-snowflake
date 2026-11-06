package gen

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// showComputedAttributes returns the computed attributes that are SHOW columns, in ComputedAttributes order.
// Membership matches Name, like ConfigFieldSets: the SHOW column and the computed attribute are the same attribute reused.
// One that is not a SHOW column (for example a DESCRIBE-only field) stays hand-written in the schema ext.
func showComputedAttributes(obj ssot.Object) (attrs []ssot.Attribute, columns []*ssot.Attribute) {
	for _, attr := range obj.ComputedAttributes {
		idx := slices.IndexFunc(obj.Show.Columns, func(col *ssot.Attribute) bool { return col != nil && col.Name == attr.Name })
		if idx < 0 {
			continue
		}
		attrs = append(attrs, attr)
		columns = append(columns, obj.Show.Columns[idx])
	}
	return attrs, columns
}

// computedResourceAttribute is the schema chrome of a computed-only attribute: a read-only field with its type and
// description. It is never Required, Optional, ForceNew, or part of CREATE and ALTER.
func computedResourceAttribute(objectName string, attr ssot.Attribute) ResourceAttribute {
	if attr.Description == "" {
		panic(fmt.Sprintf("object %q: computed attribute %q needs a Description", objectName, attr.Name))
	}
	if attr.Type == ssot.AttrList || attr.Type == ssot.AttrSet || attr.Type == ssot.AttrObject || len(attr.Children) > 0 || attr.Elem != nil {
		panic(fmt.Sprintf("object %q: computed attribute %q must be a scalar", objectName, attr.Name))
	}
	return ResourceAttribute{
		Name:            attr.Name,
		TfType:          tfTypeFor(attr.Type),
		Computed:        true,
		DescriptionExpr: strconv.Quote(attr.Description),
	}
}

// ComputedFieldSets is the SHOW values copied into the computed attributes. Names follow ComputedAttributes order.
func (m ResourceModel) ComputedFieldSets() []ConfigFieldSet {
	sets := make([]ConfigFieldSet, 0, len(m.ComputedColumns))
	for _, col := range m.ComputedColumns {
		sets = append(sets, ConfigFieldSet{Name: col.Name, Expr: configFieldCode(ssot.Object{Name: m.Name, IdType: m.IdType}, *col)})
	}
	return sets
}

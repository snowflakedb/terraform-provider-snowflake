package datasourcemodel

import (
    tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (l *ListingsModel) WithLimit(rows int) *ListingsModel {
    return l.WithLimitValue(
        tfconfig.ObjectVariable(map[string]tfconfig.Variable{
            "rows": tfconfig.IntegerVariable(rows),
        }),
    )
}

func (l *ListingsModel) WithRowsAndFrom(rows int, from string) *ListingsModel {
    return l.WithLimitValue(
        tfconfig.ObjectVariable(map[string]tfconfig.Variable{
            "rows": tfconfig.IntegerVariable(rows),
            "from": tfconfig.StringVariable(from),
        }),
    )
}



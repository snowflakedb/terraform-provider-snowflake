package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (f *FileFormatParquetModel) WithNullIf(nullIf ...string) *FileFormatParquetModel {
	nullIfList := collections.Map(nullIf, func(v string) tfconfig.Variable { return tfconfig.StringVariable(v) })
	f.NullIf = tfconfig.ListVariable(nullIfList...)
	return f
}

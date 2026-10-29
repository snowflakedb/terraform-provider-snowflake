package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

func ssotShowResultSchemaDetails() []ShowResultSchemaDetails {
	var details []ShowResultSchemaDetails
	for _, obj := range defs.All() {
		g := ssot.ObjectGenerators(obj.Generators...)
		if g.Schema.Show {
			if len(obj.ShowOutputKeys()) == 0 {
				panic(fmt.Sprintf("object %q is enabled for schema SHOW but has no SHOW fields", obj.Name))
			}
			details = append(details, showResultSchemaDetailsFromSsot(obj, false))
		}
		if g.Schema.Describe {
			if len(obj.DescribeOutputKeys()) == 0 {
				panic(fmt.Sprintf("object %q is enabled for schema DESCRIBE but has no DESCRIBE fields", obj.Name))
			}
			details = append(details, showResultSchemaDetailsFromSsot(obj, true))
		}
	}
	return details
}

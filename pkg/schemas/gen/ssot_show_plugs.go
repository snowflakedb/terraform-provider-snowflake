package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

func ssotShowResultSchemaDetails() []ShowResultSchemaDetails {
	var details []ShowResultSchemaDetails
	for _, obj := range defs.All() {
		g := ssot.ObjectGenerators(obj.Generators...)
		if g.Schema.Show && selectsForSchemaShow(obj) {
			details = append(details, showResultSchemaDetailsFromSsot(obj, false))
		}
		if g.Schema.Describe && selectsForSchemaDescribe(obj) {
			details = append(details, showResultSchemaDetailsFromSsot(obj, true))
		}
	}
	return details
}

// selectsForSchemaShow is the one filter of the SHOW schema generator: the object has SHOW output columns.
func selectsForSchemaShow(obj ssot.Object) bool {
	return obj.HasShowOutput()
}

// selectsForSchemaDescribe is the one filter of the DESCRIBE schema generator: the object has DESCRIBE output columns.
func selectsForSchemaDescribe(obj ssot.Object) bool {
	return obj.HasDescribeOutput()
}

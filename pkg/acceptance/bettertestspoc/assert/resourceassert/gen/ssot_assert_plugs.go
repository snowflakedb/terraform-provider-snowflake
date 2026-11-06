package gen

import (
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

func ResourceAssertionModels() []ResourceAssertionsModel {
	schemaDetails := GetResourceSchemaDetails()
	ssotModels := ssotResourceAssertionModels()
	plugged := make(map[string]struct{}, len(ssotModels))
	for _, model := range ssotModels {
		plugged[model.Name] = struct{}{}
	}
	all := make([]ResourceAssertionsModel, 0, len(schemaDetails)+len(ssotModels))
	for _, details := range schemaDetails {
		if _, ok := plugged[details.ObjectName()]; ok {
			continue
		}
		all = append(all, ModelFromResourceSchemaDetails(details, nil, nil))
	}
	return append(all, ssotModels...)
}

func ssotResourceAssertionModels() []ResourceAssertionsModel {
	return collections.Map(defs.SelectObjects(func(g ssot.Generators) bool { return g.ResourceAssert.Enabled }, selectsForResourceAssert), AssertionsFromSsotObject)
}

// selectsForResourceAssert is the one filter of the resource assertion generator: the object is a resource and has attributes.
func selectsForResourceAssert(obj ssot.Object) bool {
	return !obj.VariantBaseOnly && len(obj.Attributes) > 0
}

func WithPreamble(model ResourceAssertionsModel, preamble *genhelpers.PreambleModel, _ *log.Logger) ResourceAssertionsModel {
	model.PreambleModel = preamble
	return model
}

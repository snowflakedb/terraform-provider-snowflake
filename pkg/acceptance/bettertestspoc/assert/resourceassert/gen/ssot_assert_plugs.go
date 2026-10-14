package gen

import (
	"log"

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
	var models []ResourceAssertionsModel
	for _, obj := range defs.All() {
		if !ssot.ObjectGenerators(obj.Generators...).ResourceAssert.Enabled {
			continue
		}
		models = append(models, AssertionsFromSsotObject(obj))
	}
	return models
}

func WithPreamble(model ResourceAssertionsModel, preamble *genhelpers.PreambleModel, _ *log.Logger) ResourceAssertionsModel {
	model.PreambleModel = preamble
	return model
}

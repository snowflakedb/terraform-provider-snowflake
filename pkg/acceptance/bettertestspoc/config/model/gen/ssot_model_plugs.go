package gen

import (
	"log"

	resourceassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

func ModelBuilderModels() []ResourceConfigBuilderModel {
	schemaDetails := resourceassertgen.GetResourceSchemaDetails()
	ssotModels := ssotModelBuilderModels()
	plugged := make(map[string]struct{}, len(ssotModels))
	for _, model := range ssotModels {
		plugged[model.Name] = struct{}{}
	}
	all := make([]ResourceConfigBuilderModel, 0, len(schemaDetails)+len(ssotModels))
	for _, details := range schemaDetails {
		if _, ok := plugged[details.ObjectName()]; ok {
			continue
		}
		all = append(all, ModelFromResourceSchemaDetails(details, nil, log.Default()))
	}
	return append(all, ssotModels...)
}

func ssotModelBuilderModels() []ResourceConfigBuilderModel {
	var models []ResourceConfigBuilderModel
	for _, obj := range defs.All() {
		if !ssot.ObjectGenerators(obj.Generators...).Model.Enabled {
			continue
		}
		models = append(models, ModelFromSsotObject(obj))
	}
	return models
}

func WithPreamble(model ResourceConfigBuilderModel, preamble *genhelpers.PreambleModel, _ *log.Logger) ResourceConfigBuilderModel {
	model.PreambleModel = preamble
	return model
}

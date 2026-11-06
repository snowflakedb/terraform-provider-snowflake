package gen

import (
	"log"

	resourceassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
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
	return collections.Map(defs.SelectObjects(func(g ssot.Generators) bool { return g.Model.Enabled }, selectsForModel), ModelFromSsotObject)
}

// selectsForModel is the one filter of the config model builder generator: the object is a resource and has attributes.
func selectsForModel(obj ssot.Object) bool {
	return !obj.VariantBaseOnly && len(obj.Attributes) > 0
}

func WithPreamble(model ResourceConfigBuilderModel, preamble *genhelpers.PreambleModel, _ *log.Logger) ResourceConfigBuilderModel {
	model.PreambleModel = preamble
	return model
}

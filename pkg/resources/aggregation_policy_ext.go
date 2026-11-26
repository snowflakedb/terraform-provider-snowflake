package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	aggregationPolicyCreateExt      = PreviewFeatureCreateContextWrapper(string(previewfeatures.AggregationPolicyResource), TrackingCreateWrapper(resources.AggregationPolicy, CreateAggregationPolicy))
	aggregationPolicyReadExt        = PreviewFeatureReadContextWrapper(string(previewfeatures.AggregationPolicyResource), TrackingReadWrapper(resources.AggregationPolicy, ReadAggregationPolicyFunc(true)))
	aggregationPolicyUpdateExt      = PreviewFeatureUpdateContextWrapper(string(previewfeatures.AggregationPolicyResource), TrackingUpdateWrapper(resources.AggregationPolicy, UpdateAggregationPolicy))
	aggregationPolicyDescriptionExt = docs.WithPreviewFeatureMark() + "Resource used to manage aggregation policies. For more information, check [aggregation policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-aggregation-policy)."
	aggregationPolicyImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.AggregationPolicy, ImportAggregationPolicy),
	}
)

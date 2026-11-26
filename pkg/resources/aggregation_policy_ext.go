package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	aggregationPolicyCreateExt      = TrackingCreateWrapper(resources.AggregationPolicy, CreateAggregationPolicy)
	aggregationPolicyReadExt        = TrackingReadWrapper(resources.AggregationPolicy, ReadAggregationPolicyFunc(true))
	aggregationPolicyUpdateExt      = TrackingUpdateWrapper(resources.AggregationPolicy, UpdateAggregationPolicy)
	aggregationPolicyDescriptionExt = "Resource used to manage aggregation policies. For more information, check [aggregation policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-aggregation-policy)."
	aggregationPolicyImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.AggregationPolicy, ImportAggregationPolicy),
	}
)

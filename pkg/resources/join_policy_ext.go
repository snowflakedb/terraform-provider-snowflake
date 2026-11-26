package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	joinPolicyCreateExt      = TrackingCreateWrapper(resources.JoinPolicy, CreateJoinPolicy)
	joinPolicyReadExt        = TrackingReadWrapper(resources.JoinPolicy, ReadJoinPolicyFunc(true))
	joinPolicyUpdateExt      = TrackingUpdateWrapper(resources.JoinPolicy, UpdateJoinPolicy)
	joinPolicyDescriptionExt = "Resource used to manage join policies. For more information, check [join policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-join-policy)."
	joinPolicyImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.JoinPolicy, ImportJoinPolicy),
	}
)

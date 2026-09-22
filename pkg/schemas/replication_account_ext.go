package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (replicationAccountToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"comment": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (replicationAccountToSchemaMapper) additionalToSchema(src *sdk.ReplicationAccount, dst map[string]any) {
	// comment is sql.NullString; skip when unset.
	if src.Comment.Valid {
		dst["comment"] = src.Comment.String
	}
}

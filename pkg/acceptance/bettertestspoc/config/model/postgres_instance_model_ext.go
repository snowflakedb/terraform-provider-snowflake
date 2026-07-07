package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func PostgresInstanceFromId(
	resourceName string,
	id sdk.AccountObjectIdentifier,
) *PostgresInstanceModel {
	p := &PostgresInstanceModel{ResourceModelMeta: config.Meta(resourceName, resources.PostgresInstance)}
	p.WithName(id.Name())
	return p
}

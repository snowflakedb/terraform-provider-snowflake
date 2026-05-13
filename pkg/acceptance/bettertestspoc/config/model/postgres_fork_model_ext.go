package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func PostgresForkFromId(
	resourceName string,
	id sdk.AccountObjectIdentifier,
	forkFrom sdk.AccountObjectIdentifier,
) *PostgresForkModel {
	p := &PostgresForkModel{ResourceModelMeta: config.Meta(resourceName, resources.PostgresFork)}
	p.WithName(id.Name())
	p.WithForkFrom(forkFrom.FullyQualifiedName())
	return p
}

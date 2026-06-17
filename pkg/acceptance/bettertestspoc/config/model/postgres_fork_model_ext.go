package model

import (
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
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

func (p *PostgresForkModel) WithAtTimestamp(timestamp time.Time) *PostgresForkModel {
	p.At = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"timestamp": tfconfig.StringVariable(timestamp.UTC().Format("2006-01-02 15:04:05")),
	})
	return p
}

func (p *PostgresForkModel) WithAtOffset(offset string) *PostgresForkModel {
	p.At = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"offset": tfconfig.StringVariable(offset),
	})
	return p
}

func (p *PostgresForkModel) WithBeforeTimestamp(timestamp time.Time) *PostgresForkModel {
	p.Before = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"timestamp": tfconfig.StringVariable(timestamp.UTC().Format("2006-01-02 15:04:05")),
	})
	return p
}

func (p *PostgresForkModel) WithBeforeOffset(offset string) *PostgresForkModel {
	p.Before = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"offset": tfconfig.StringVariable(offset),
	})
	return p
}

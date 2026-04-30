package model

import (
	"encoding/json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

type PostgresForkModel struct {
	Name               tfconfig.Variable `json:"name,omitempty"`
	ForkFrom           tfconfig.Variable `json:"fork_from,omitempty"`
	AtTimestamp        tfconfig.Variable `json:"at_timestamp,omitempty"`
	AtOffset           tfconfig.Variable `json:"at_offset,omitempty"`
	BeforeTimestamp    tfconfig.Variable `json:"before_timestamp,omitempty"`
	BeforeOffset       tfconfig.Variable `json:"before_offset,omitempty"`
	ComputeFamily      tfconfig.Variable `json:"compute_family,omitempty"`
	StorageSizeGb      tfconfig.Variable `json:"storage_size_gb,omitempty"`
	HighAvailability   tfconfig.Variable `json:"high_availability,omitempty"`
	PostgresSettings   tfconfig.Variable `json:"postgres_settings,omitempty"`
	Comment            tfconfig.Variable `json:"comment,omitempty"`
	FullyQualifiedName tfconfig.Variable `json:"fully_qualified_name,omitempty"`

	*config.ResourceModelMeta
}

func PostgresFork(
	resourceName string,
	name string,
	forkFrom string,
) *PostgresForkModel {
	p := &PostgresForkModel{ResourceModelMeta: config.Meta(resourceName, resources.PostgresFork)}
	p.WithName(name)
	p.WithForkFrom(forkFrom)
	return p
}

func PostgresForkWithDefaultMeta(
	name string,
	forkFrom string,
) *PostgresForkModel {
	p := &PostgresForkModel{ResourceModelMeta: config.DefaultMeta(resources.PostgresFork)}
	p.WithName(name)
	p.WithForkFrom(forkFrom)
	return p
}

func (p *PostgresForkModel) MarshalJSON() ([]byte, error) {
	type Alias PostgresForkModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(p),
		DependsOn: p.DependsOn(),
	})
}

func (p *PostgresForkModel) WithDependsOn(values ...string) *PostgresForkModel {
	p.SetDependsOn(values...)
	return p
}

func (p *PostgresForkModel) WithName(name string) *PostgresForkModel {
	p.Name = tfconfig.StringVariable(name)
	return p
}

func (p *PostgresForkModel) WithForkFrom(forkFrom string) *PostgresForkModel {
	p.ForkFrom = tfconfig.StringVariable(forkFrom)
	return p
}

func (p *PostgresForkModel) WithAtTimestamp(atTimestamp string) *PostgresForkModel {
	p.AtTimestamp = tfconfig.StringVariable(atTimestamp)
	return p
}

func (p *PostgresForkModel) WithAtOffset(atOffset string) *PostgresForkModel {
	p.AtOffset = tfconfig.StringVariable(atOffset)
	return p
}

func (p *PostgresForkModel) WithBeforeTimestamp(beforeTimestamp string) *PostgresForkModel {
	p.BeforeTimestamp = tfconfig.StringVariable(beforeTimestamp)
	return p
}

func (p *PostgresForkModel) WithBeforeOffset(beforeOffset string) *PostgresForkModel {
	p.BeforeOffset = tfconfig.StringVariable(beforeOffset)
	return p
}

func (p *PostgresForkModel) WithComputeFamily(computeFamily string) *PostgresForkModel {
	p.ComputeFamily = tfconfig.StringVariable(computeFamily)
	return p
}

func (p *PostgresForkModel) WithStorageSizeGb(storageSizeGb int) *PostgresForkModel {
	p.StorageSizeGb = tfconfig.IntegerVariable(storageSizeGb)
	return p
}

func (p *PostgresForkModel) WithHighAvailability(highAvailability bool) *PostgresForkModel {
	p.HighAvailability = tfconfig.BoolVariable(highAvailability)
	return p
}

func (p *PostgresForkModel) WithPostgresSettings(postgresSettings string) *PostgresForkModel {
	p.PostgresSettings = tfconfig.StringVariable(postgresSettings)
	return p
}

func (p *PostgresForkModel) WithComment(comment string) *PostgresForkModel {
	p.Comment = tfconfig.StringVariable(comment)
	return p
}

func (p *PostgresForkModel) WithFullyQualifiedName(fullyQualifiedName string) *PostgresForkModel {
	p.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return p
}

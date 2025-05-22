package config

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
)

// TODO [SNOW-1501905]: add possibility to have reference to another object (e.g. WithResourceMonitorReference); new config.Variable impl?
// TODO [SNOW-1501905]: generate With/SetDependsOn for the resources to preserve builder pattern
// TODO [SNOW-1501905]: add a convenience method to use multiple configs from multiple models
// TODO [SNOW-1501905]: add provider to resource/datasource models (use in the grant_ownership_acceptance_test)

// ResourceModel is the base interface all of our resource config models will implement.
// To allow easy implementation, ResourceModelMeta can be embedded inside the struct (and the struct will automatically implement it).
type ResourceModel interface {
	Resource() resources.Resource
	ResourceName() string
	SetResourceName(name string)
	ResourceReference() string
	DependsOn() []string
	SetDependsOn(values ...string)
}

type ResourceModelMeta struct {
	name      string
	resource  resources.Resource
	dependsOn []string
}

func (m *ResourceModelMeta) Resource() resources.Resource {
	return m.resource
}

func (m *ResourceModelMeta) ResourceName() string {
	return m.name
}

func (m *ResourceModelMeta) SetResourceName(name string) {
	m.name = name
}

func (m *ResourceModelMeta) ResourceReference() string {
	return fmt.Sprintf(`%s.%s`, m.resource, m.name)
}

func (m *ResourceModelMeta) DependsOn() []string {
	return m.dependsOn
}

func (m *ResourceModelMeta) SetDependsOn(values ...string) {
	m.dependsOn = values
}

// DefaultResourceName is exported to allow assertions against the resources using the default name.
const DefaultResourceName = "test"

func DefaultMeta(resource resources.Resource) *ResourceModelMeta {
	return &ResourceModelMeta{name: DefaultResourceName, resource: resource}
}

func Meta(resourceName string, resource resources.Resource) *ResourceModelMeta {
	return &ResourceModelMeta{name: resourceName, resource: resource}
}

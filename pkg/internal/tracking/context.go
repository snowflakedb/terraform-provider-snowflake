package tracking

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

const (
	CurrentSchemaVersion string = "1"
	ProviderVersion      string = "v1.0.6" // TODO(SNOW-1814934): Currently hardcoded, make it computed
	MetadataPrefix       string = "terraform_provider_usage_tracking"
)

type key struct{}

var metadataContextKey key

type Operation string

const (
	CreateOperation     Operation = "create"
	ReadOperation       Operation = "read"
	UpdateOperation     Operation = "update"
	DeleteOperation     Operation = "delete"
	ImportOperation     Operation = "import"
	CustomDiffOperation Operation = "custom_diff"
)

type Metadata struct {
	SchemaVersion string    `json:"json_schema_version,omitempty"`
	Version       string    `json:"version,omitempty"`
	Resource      string    `json:"resource,omitempty"`
	Datasource    string    `json:"datasource,omitempty"`
	Operation     Operation `json:"operation,omitempty"`
}

func (m Metadata) validate() error {
	errs := make([]error, 0)
	if m.SchemaVersion == "" {
		errs = append(errs, errors.New("schema version for metadata should not be empty"))
	}
	if m.Version == "" {
		errs = append(errs, errors.New("provider version for metadata should not be empty"))
	}
	if m.Resource == "" && m.Datasource == "" {
		errs = append(errs, errors.New("either resource or data source name for metadata should be specified"))
	}
	if m.Operation == "" {
		errs = append(errs, errors.New("operation for metadata should not be empty"))
	}
	return errors.Join(errs...)
}

// newTestMetadata is a helper constructor that is used only for testing purposes
func newTestMetadata(version string, resource resources.Resource, operation Operation) Metadata {
	return Metadata{
		SchemaVersion: CurrentSchemaVersion,
		Version:       version,
		Resource:      resource.String(),
		Operation:     operation,
	}
}

func NewVersionedResourceMetadata(resource resources.Resource, operation Operation) Metadata {
	return Metadata{
		SchemaVersion: CurrentSchemaVersion,
		Version:       ProviderVersion,
		Resource:      resource.String(),
		Operation:     operation,
	}
}

func NewVersionedDatasourceMetadata(datasource datasources.Datasource) Metadata {
	return Metadata{
		SchemaVersion: CurrentSchemaVersion,
		Version:       ProviderVersion,
		Datasource:    datasource.String(),
		Operation:     ReadOperation,
	}
}

func NewContext(ctx context.Context, metadata Metadata) context.Context {
	return context.WithValue(ctx, metadataContextKey, metadata)
}

func FromContext(ctx context.Context) (Metadata, bool) {
	metadata, ok := ctx.Value(metadataContextKey).(Metadata)
	return metadata, ok
}

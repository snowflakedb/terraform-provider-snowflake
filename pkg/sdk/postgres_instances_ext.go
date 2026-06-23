package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

func (r *CreatePostgresInstanceRequest) GetName() AccountObjectIdentifier {
	return r.name
}

// PostgresInstanceDetails represents the parsed result of DESCRIBE POSTGRES INSTANCE
type PostgresInstanceDetails struct {
	Name                         string
	Owner                        string
	OwnerRoleType                string
	CreatedOn                    string
	UpdatedOn                    string
	Type                         string
	Host                         string
	Origin                       *string
	PrivatelinkServiceIdentifier *string
	ComputeFamily                string
	StorageSizeGb                int
	PostgresVersion              int
	HighAvailability             bool
	AuthenticationAuthority      string
	State                        string
	RetentionTime                int
	MaintenanceWindowStart       int
	Comment                      *string
	NetworkPolicy                *AccountObjectIdentifier
	PostgresSettings             *string
	StorageIntegration           *AccountObjectIdentifier
}

// ParsePostgresInstanceDetails parses []PostgresInstanceProperty into PostgresInstanceDetails
func ParsePostgresInstanceDetails(properties []PostgresInstanceProperty) (*PostgresInstanceDetails, error) {
	details := &PostgresInstanceDetails{}
	var errs []error
	for _, prop := range properties {
		switch strings.ToLower(prop.Property) {
		case "name":
			details.Name = prop.Value
		case "owner":
			details.Owner = prop.Value
		case "owner_role_type":
			details.OwnerRoleType = prop.Value
		case "created_on":
			details.CreatedOn = prop.Value
		case "updated_on":
			details.UpdatedOn = prop.Value
		case "type":
			details.Type = prop.Value
		case "host":
			details.Host = prop.Value
		case "origin":
			details.Origin = String(prop.Value)
		case "privatelink_service_identifier":
			details.PrivatelinkServiceIdentifier = String(prop.Value)
		case "compute_family":
			details.ComputeFamily = prop.Value
		case "storage_size_gb":
			if prop.Value != "" {
				if val, err := strconv.Atoi(prop.Value); err != nil {
					errs = append(errs, err)
				} else {
					details.StorageSizeGb = val
				}
			}
		case "postgres_version":
			if prop.Value != "" {
				if val, err := strconv.Atoi(prop.Value); err != nil {
					errs = append(errs, err)
				} else {
					details.PostgresVersion = val
				}
			}
		case "high_availability":
			details.HighAvailability = prop.Value == "true"
		case "authentication_authority":
			details.AuthenticationAuthority = prop.Value
		case "state":
			details.State = prop.Value
		case "retention_time":
			if prop.Value != "" {
				if val, err := strconv.Atoi(prop.Value); err != nil {
					errs = append(errs, err)
				} else {
					details.RetentionTime = val
				}
			}
		case "maintenance_window_start":
			if prop.Value != "" {
				if val, err := strconv.Atoi(prop.Value); err != nil {
					errs = append(errs, err)
				} else {
					details.MaintenanceWindowStart = val
				}
			}
		case "comment":
			details.Comment = String(prop.Value)
		case "network_policy":
			if prop.Value != "" {
				details.NetworkPolicy = Pointer(NewAccountObjectIdentifier(prop.Value))
			}
		case "postgres_settings":
			details.PostgresSettings = String(prop.Value)
		case "storage_integration":
			if prop.Value != "" {
				details.StorageIntegration = Pointer(NewAccountObjectIdentifier(prop.Value))
			}
		}
	}
	return details, errors.Join(errs...)
}

func (v *postgresInstances) DescribeDetails(ctx context.Context, id AccountObjectIdentifier) (*PostgresInstanceDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return ParsePostgresInstanceDetails(properties)
}

// NormalizePostgresSettings parses a postgres_settings JSON string into a canonical
// form so Terraform can compare user input with Snowflake responses without spurious
// diffs due to key ordering or whitespace. An empty string or an empty JSON object
// ("{}") is normalized to "" to represent "not set".
func NormalizePostgresSettings(s string) (string, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return "", nil
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(trimmed), &m); err != nil {
		return "", err
	}

	if len(m) == 0 {
		return "", nil
	}

	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// NormalizePostgresSettingsPtr is a pointer-safe variant of NormalizePostgresSettings
// for use on the read path. Returns nil for nil input, empty/"{}" JSON, or parse errors.
func NormalizePostgresSettingsPtr(s *string) *string {
	if s == nil {
		return nil
	}
	normalized, err := NormalizePostgresSettings(*s)
	if err != nil || normalized == "" {
		return nil
	}
	return &normalized
}

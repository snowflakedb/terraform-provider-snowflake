package sdk

import (
	"context"
	"errors"
	"strconv"
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
	PostgresVersion              string
	HighAvailability             bool
	AuthenticationAuthority      string
	State                        string
	RetentionTime                int
	MaintenanceWindowStart       int
	Comment                      *string
	NetworkPolicy                *string
	PostgresSettings             *string
	StorageIntegration           *string
}

// ParsePostgresInstanceDetails parses []PostgresInstanceProperty into PostgresInstanceDetails
func ParsePostgresInstanceDetails(properties []PostgresInstanceProperty) (*PostgresInstanceDetails, error) {
	details := &PostgresInstanceDetails{}
	var errs []error
	for _, prop := range properties {
		switch prop.Property {
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
			if val, err := strconv.Atoi(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.StorageSizeGb = val
			}
		case "postgres_version":
			details.PostgresVersion = prop.Value
		case "high_availability":
			details.HighAvailability = prop.Value == "true"
		case "authentication_authority":
			details.AuthenticationAuthority = prop.Value
		case "state":
			details.State = prop.Value
		case "retention_time":
			if val, err := strconv.Atoi(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.RetentionTime = val
			}
		case "maintenance_window_start":
			if val, err := strconv.Atoi(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.MaintenanceWindowStart = val
			}
		case "comment":
			details.Comment = String(prop.Value)
		case "network_policy":
			details.NetworkPolicy = String(prop.Value)
		case "postgres_settings":
			details.PostgresSettings = String(prop.Value)
		case "storage_integration":
			details.StorageIntegration = String(prop.Value)
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

func (r postgresInstancesRow) convert() (*PostgresInstance, error) {
	pi := &PostgresInstance{
		Name:                    r.Name,
		Owner:                   r.Owner,
		OwnerRoleType:           r.OwnerRoleType,
		CreatedOn:               r.CreatedOn,
		UpdatedOn:               r.UpdatedOn,
		Type:                    r.Type,
		ComputeFamily:           r.ComputeFamily,
		AuthenticationAuthority: r.AuthenticationAuthority,
		StorageSize:             r.StorageSize,
		PostgresVersion:         r.PostgresVersion,
		IsHa:                    r.IsHa == "true",
		RetentionTime:           r.RetentionTime,
	}
	mapNullString(&pi.Origin, r.Origin)
	mapNullString(&pi.Host, r.Host)
	mapNullString(&pi.PrivatelinkServiceIdentifier, r.PrivatelinkServiceIdentifier)
	mapNullString(&pi.PostgresSettings, r.PostgresSettings)
	mapNullString(&pi.Comment, r.Comment)
	mapStringWithMapping(&pi.State, r.State, ToPostgresInstanceState)
	return pi, nil
}

func (r postgresInstanceDetailsRow) convert() (*PostgresInstanceProperty, error) {
	return &PostgresInstanceProperty{
		Property: r.Property,
		Value:    r.Value,
	}, nil
}

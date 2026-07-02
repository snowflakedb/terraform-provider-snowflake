package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (r *CreatePostgresInstanceRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *AlterPostgresInstanceRequest) GetName() AccountObjectIdentifier {
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
			if prop.Value != "" {
				details.Comment = String(prop.Value)
			}
		case "network_policy":
			if prop.Value != "" {
				details.NetworkPolicy = Pointer(NewAccountObjectIdentifier(prop.Value))
			}
		case "postgres_settings":
			if prop.Value != "" {
				details.PostgresSettings = String(prop.Value)
			}
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

// CreateSafely creates a Postgres instance and polls ShowByID every 3 seconds until the
// instance reaches READY state. The caller controls the wait budget via ctx — use
// context.WithTimeout to set a deadline. Returns the ready instance or an error.
func (v *postgresInstances) CreateSafely(ctx context.Context, req *CreatePostgresInstanceRequest) (*PostgresInstance, error) {
	return createSafelyPolling(
		ctx,
		func() error { return v.Create(ctx, req) },
		func() (*PostgresInstance, error) { return v.ShowByID(ctx, req.GetName()) },
	)
}

// createSafelyPolling is the polling loop shared between CreateSafely and its unit tests.
func createSafelyPolling(ctx context.Context, doCreate func() error, doShowByID func() (*PostgresInstance, error)) (*PostgresInstance, error) {
	if err := doCreate(); err != nil {
		return nil, err
	}
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("postgres instance did not reach READY state: %w", ctx.Err())
		default:
		}
		instance, err := doShowByID()
		if err != nil {
			return nil, err
		}
		if instance.State == PostgresInstanceStateReady {
			return instance, nil
		}
		time.Sleep(3 * time.Second)
	}
}

func (v *postgresInstances) AlterSafely(ctx context.Context, req *AlterPostgresInstanceRequest) error {
	return updateSafelyPolling(
		ctx,
		func() error { return v.Alter(ctx, req) },
		func() (*PostgresInstance, error) { return v.ShowByID(ctx, req.GetName()) },
	)
}

func updateSafelyPolling(ctx context.Context, doUpdate func() error, doShowByID func() (*PostgresInstance, error)) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("postgres instance did not reach READY state: %w", ctx.Err())
		default:
		}
		instance, err := doShowByID()
		if err != nil {
			return err
		}
		if instance.State != PostgresInstanceStateReady {
			time.Sleep(3 * time.Second)
			continue
		}
		if err := doUpdate(); err != nil {
			if strings.Contains(err.Error(), "must be complete before issuing ALTER") {
				time.Sleep(3 * time.Second)
				continue
			}
			return err
		}
		// ALTER accepted; wait for READY to ensure the backend has committed the
		// change before the caller issues a read (Snowflake applies some mutations
		// asynchronously even after returning success from ALTER).
		for {
			select {
			case <-ctx.Done():
				return fmt.Errorf("postgres instance did not settle after alter: %w", ctx.Err())
			default:
			}
			time.Sleep(3 * time.Second)
			instance, err = doShowByID()
			if err != nil {
				return err
			}
			if instance.State == PostgresInstanceStateReady {
				return nil
			}
		}
	}
}

// NormalizePostgresSettings parses a postgres_settings JSON string into a canonical
// form so Terraform can compare user input with Snowflake responses without spurious
// diffs due to key ordering or whitespace. An empty string or an empty JSON object
// ("{}") is normalized to "" to represent "not set".
func NormalizePostgresSettings(s string) (string, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" || trimmed == "{}" {
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

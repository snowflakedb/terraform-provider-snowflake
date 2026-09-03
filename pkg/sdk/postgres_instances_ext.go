package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
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
	MaintenanceWindowStart       *int
	Comment                      *string
	NetworkPolicy                *AccountObjectIdentifier
	PostgresSettings             *string
	StorageIntegration           *AccountObjectIdentifier
	// HasAnyRunningOperations is true when DESCRIBE operations JSON column contains any items,
	// indicating some work (e.g. altering some field on the postgres instance) is running in the background on the server side.
	HasAnyRunningOperations bool
	// OperationErrors is the joined set of FAILED DESCRIBE operation records.
	OperationErrors error
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
					details.MaintenanceWindowStart = &val
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
		case "operations":
			trimmed := strings.TrimSpace(prop.Value)
			if trimmed != "" && !strings.EqualFold(trimmed, "none") {
				var entries map[string]struct {
					State string `json:"state"`
					Error string `json:"error"`
				}
				if err := json.Unmarshal([]byte(trimmed), &entries); err != nil {
					details.HasAnyRunningOperations = true
				} else {
					var failedNames []string
					failedReasons := make(map[string]string)
					for name, entry := range entries {
						state := strings.ToUpper(strings.TrimSpace(entry.State))
						switch state {
						case "FAILED":
							reason := entry.Error
							if reason == "" {
								reason = "no error detail reported"
							}
							failedNames = append(failedNames, name)
							failedReasons[name] = reason
						case string(PostgresInstanceStateReady):
							// Leftover create snapshot; not running.
						default:
							details.HasAnyRunningOperations = true
						}
					}
					slices.Sort(failedNames)
					opErrs := make([]error, 0, len(failedNames))
					for _, name := range failedNames {
						opErrs = append(opErrs, fmt.Errorf("postgres instance %s operation failed: %s", name, failedReasons[name]))
					}
					details.OperationErrors = errors.Join(opErrs...)
				}
			}
		}
	}
	return details, errors.Join(errs...)
}

func (d *PostgresInstanceDetails) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(d.Name)
}

func (v *postgresInstances) DescribeDetails(ctx context.Context, id AccountObjectIdentifier) (*PostgresInstanceDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return ParsePostgresInstanceDetails(properties)
}

// CreateSafely creates a Postgres instance and polls ShowByID every 3 seconds until the
// instance reaches READY state. CREATE actually changes SHOW state, unlike ALTER, so this
// does not wait on HasAnyRunningOperations. The caller controls the wait budget via ctx —
// use context.WithTimeout to set a deadline. Returns the ready instance or an error.
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
	return pollUntilStateOneOf(ctx, doShowByID, PostgresInstanceStateReady)
}

const postgresInstancePollInterval = 3 * time.Second

func (v *postgresInstances) AlterSafely(ctx context.Context, req *AlterPostgresInstanceRequest) error {
	return updateSafelyPolling(
		ctx,
		req,
		func() error { return v.Alter(ctx, req) },
		func() (*PostgresInstanceDetails, error) { return v.DescribeDetails(ctx, req.GetName()) },
		func() (*PostgresInstance, error) { return v.ShowByID(ctx, req.GetName()) },
	)
}

func updateSafelyPolling(
	ctx context.Context,
	req *AlterPostgresInstanceRequest,
	doUpdate func() error,
	doDescribe func() (*PostgresInstanceDetails, error),
	doShowByID func() (*PostgresInstance, error),
) error {
	// SHOW state is not a completion signal: ALTER does not change state, and several
	// async properties leave the instance READY for the whole change.
	if _, err := pollUntilNoRunningOperations(ctx, doDescribe); err != nil {
		return err
	}
	for {
		if err := doUpdate(); err != nil {
			if strings.Contains(err.Error(), ErrPostgresOperationMustBeComplete.Error()) {
				if _, err := pollUntilNoRunningOperations(ctx, doDescribe); err != nil {
					return err
				}
				continue
			}
			return err
		}
		break
	}
	details, err := pollUntilNoRunningOperations(ctx, doDescribe)
	if err != nil {
		return err
	}
	if details != nil && details.OperationErrors != nil {
		return details.OperationErrors
	}
	switch {
	case req != nil && req.Suspend != nil && *req.Suspend:
		_, err := pollUntilStateOneOf(ctx, doShowByID, PostgresInstanceStateSuspending, PostgresInstanceStateSuspended)
		return err
	case req != nil && req.Resume != nil && *req.Resume:
		_, err := pollUntilStateOneOf(ctx, doShowByID, PostgresInstanceStateReady)
		return err
	}
	return pollUntilNetworkPolicyMatches(ctx, req, doDescribe)
}

func pollUntilNoRunningOperations(ctx context.Context, doDescribe func() (*PostgresInstanceDetails, error)) (*PostgresInstanceDetails, error) {
	return util.PollUntil(ctx, postgresInstancePollInterval, doDescribe, func(d *PostgresInstanceDetails) bool {
		return !d.HasAnyRunningOperations
	}, "postgres instance still has running operations")
}

func pollUntilNetworkPolicyMatches(ctx context.Context, req *AlterPostgresInstanceRequest, doDescribe func() (*PostgresInstanceDetails, error)) error {
	var want string
	switch {
	case req.Set != nil && req.Set.NetworkPolicy != nil:
		want = req.Set.NetworkPolicy.Name()
	case req.Unset != nil && req.Unset.NetworkPolicy != nil && *req.Unset.NetworkPolicy:
		want = ""
	default:
		return nil
	}

	_, err := util.PollUntil(ctx, postgresInstancePollInterval, doDescribe, func(d *PostgresInstanceDetails) bool {
		got := ""
		if d.NetworkPolicy != nil && !strings.EqualFold(d.NetworkPolicy.Name(), "None") {
			got = d.NetworkPolicy.Name()
		}
		return strings.EqualFold(got, want)
	}, "postgres instance network_policy did not converge")
	return err
}

func pollUntilStateOneOf(ctx context.Context, doShowByID func() (*PostgresInstance, error), states ...PostgresInstanceState) (*PostgresInstance, error) {
	return util.PollUntil(ctx, postgresInstancePollInterval, doShowByID, func(i *PostgresInstance) bool {
		return slices.Contains(states, i.State)
	}, "postgres instance did not reach expected state")
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

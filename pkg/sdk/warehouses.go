package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
)

var (
	_ validatable = new(CreateWarehouseOptions)
	_ validatable = new(AlterWarehouseOptions)
	_ validatable = new(DropWarehouseOptions)
	_ validatable = new(ShowWarehouseOptions)
	_ validatable = new(describeWarehouseOptions)

	_ convertibleRow[Warehouse] = new(warehouseDBRow)
)

type Warehouses interface {
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateWarehouseOptions) error
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterWarehouseOptions) error
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropWarehouseOptions) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, opts *ShowWarehouseOptions) ([]Warehouse, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Warehouse, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Warehouse, error)
	Describe(ctx context.Context, id AccountObjectIdentifier) (*WarehouseDetails, error)
	ShowParameters(ctx context.Context, id AccountObjectIdentifier) ([]*Parameter, error)
}

var _ Warehouses = (*warehouses)(nil)

type warehouses struct {
	client *Client
}

type WarehouseType string

const (
	WarehouseTypeStandard          WarehouseType = "STANDARD"
	WarehouseTypeSnowparkOptimized WarehouseType = "SNOWPARK-OPTIMIZED"
)

func ToWarehouseType(s string) (WarehouseType, error) {
	switch strings.ToUpper(s) {
	case string(WarehouseTypeStandard):
		return WarehouseTypeStandard, nil
	case string(WarehouseTypeSnowparkOptimized):
		return WarehouseTypeSnowparkOptimized, nil
	default:
		return "", fmt.Errorf("invalid warehouse type: %s", s)
	}
}

func (e WarehouseType) FromString(s string) (WarehouseType, error) {
	return ToWarehouseType(s)
}

type WarehouseSize string

const (
	WarehouseSizeXSmall   WarehouseSize = "XSMALL"
	WarehouseSizeSmall    WarehouseSize = "SMALL"
	WarehouseSizeMedium   WarehouseSize = "MEDIUM"
	WarehouseSizeLarge    WarehouseSize = "LARGE"
	WarehouseSizeXLarge   WarehouseSize = "XLARGE"
	WarehouseSizeXXLarge  WarehouseSize = "XXLARGE"
	WarehouseSizeXXXLarge WarehouseSize = "XXXLARGE"
	WarehouseSizeX4Large  WarehouseSize = "X4LARGE"
	WarehouseSizeX5Large  WarehouseSize = "X5LARGE"
	WarehouseSizeX6Large  WarehouseSize = "X6LARGE"
)

func ToWarehouseSize(s string) (WarehouseSize, error) {
	switch strings.ToUpper(s) {
	case string(WarehouseSizeXSmall), "X-SMALL":
		return WarehouseSizeXSmall, nil
	case string(WarehouseSizeSmall):
		return WarehouseSizeSmall, nil
	case string(WarehouseSizeMedium):
		return WarehouseSizeMedium, nil
	case string(WarehouseSizeLarge):
		return WarehouseSizeLarge, nil
	case string(WarehouseSizeXLarge), "X-LARGE":
		return WarehouseSizeXLarge, nil
	case string(WarehouseSizeXXLarge), "X2LARGE", "2X-LARGE":
		return WarehouseSizeXXLarge, nil
	case string(WarehouseSizeXXXLarge), "X3LARGE", "3X-LARGE":
		return WarehouseSizeXXXLarge, nil
	case string(WarehouseSizeX4Large), "4X-LARGE":
		return WarehouseSizeX4Large, nil
	case string(WarehouseSizeX5Large), "5X-LARGE":
		return WarehouseSizeX5Large, nil
	case string(WarehouseSizeX6Large), "6X-LARGE":
		return WarehouseSizeX6Large, nil
	default:
		return "", fmt.Errorf("invalid warehouse size: %s", s)
	}
}

func (e WarehouseSize) FromString(s string) (WarehouseSize, error) {
	return ToWarehouseSize(s)
}

type ScalingPolicy string

const (
	ScalingPolicyStandard ScalingPolicy = "STANDARD"
	ScalingPolicyEconomy  ScalingPolicy = "ECONOMY"
)

func ToScalingPolicy(s string) (ScalingPolicy, error) {
	switch strings.ToUpper(s) {
	case string(ScalingPolicyStandard):
		return ScalingPolicyStandard, nil
	case string(ScalingPolicyEconomy):
		return ScalingPolicyEconomy, nil
	default:
		return "", fmt.Errorf("invalid scaling policy: %s", s)
	}
}

func (e ScalingPolicy) FromString(s string) (ScalingPolicy, error) {
	return ToScalingPolicy(s)
}

// CreateWarehouseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-warehouse.
type CreateWarehouseOptions struct {
	create      bool                    `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	warehouse   bool                    `ddl:"static" sql:"WAREHOUSE"`
	IfNotExists *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        AccountObjectIdentifier `ddl:"identifier"`

	// Object properties
	WarehouseType                   *WarehouseType           `ddl:"parameter,single_quotes" sql:"WAREHOUSE_TYPE"`
	WarehouseSize                   *WarehouseSize           `ddl:"parameter,single_quotes" sql:"WAREHOUSE_SIZE"`
	MaxClusterCount                 *int                     `ddl:"parameter" sql:"MAX_CLUSTER_COUNT"`
	MinClusterCount                 *int                     `ddl:"parameter" sql:"MIN_CLUSTER_COUNT"`
	ScalingPolicy                   *ScalingPolicy           `ddl:"parameter,single_quotes" sql:"SCALING_POLICY"`
	AutoSuspend                     *int                     `ddl:"parameter" sql:"AUTO_SUSPEND"`
	AutoResume                      *bool                    `ddl:"parameter" sql:"AUTO_RESUME"`
	InitiallySuspended              *bool                    `ddl:"parameter" sql:"INITIALLY_SUSPENDED"`
	ResourceMonitor                 *AccountObjectIdentifier `ddl:"identifier,equals" sql:"RESOURCE_MONITOR"`
	Comment                         *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
	EnableQueryAcceleration         *bool                    `ddl:"parameter" sql:"ENABLE_QUERY_ACCELERATION"`
	QueryAccelerationMaxScaleFactor *int                     `ddl:"parameter" sql:"QUERY_ACCELERATION_MAX_SCALE_FACTOR"`

	// Object params
	MaxConcurrencyLevel             *int             `ddl:"parameter" sql:"MAX_CONCURRENCY_LEVEL"`
	StatementQueuedTimeoutInSeconds *int             `ddl:"parameter" sql:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds       *int             `ddl:"parameter" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`
	Tag                             []TagAssociation `ddl:"keyword,parentheses" sql:"TAG"`
}

func (opts *CreateWarehouseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.MinClusterCount) && valueSet(opts.MaxClusterCount) && !validateIntGreaterThanOrEqual(*opts.MaxClusterCount, *opts.MinClusterCount) {
		errs = append(errs, fmt.Errorf("MinClusterCount must be less than or equal to MaxClusterCount"))
	}
	if valueSet(opts.QueryAccelerationMaxScaleFactor) && !validateIntInRangeInclusive(*opts.QueryAccelerationMaxScaleFactor, 0, 100) {
		errs = append(errs, errIntBetween("CreateWarehouseOptions", "QueryAccelerationMaxScaleFactor", 0, 100))
	}
	return errors.Join(errs...)
}

func (c *warehouses) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateWarehouseOptions) error {
	if opts == nil {
		opts = &CreateWarehouseOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	stmt, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, stmt)
	return err
}

// AlterWarehouseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse.
type AlterWarehouseOptions struct {
	alter     bool                    `ddl:"static" sql:"ALTER"`
	warehouse bool                    `ddl:"static" sql:"WAREHOUSE"`
	IfExists  *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name      AccountObjectIdentifier `ddl:"identifier"`

	Suspend         *bool                    `ddl:"keyword" sql:"SUSPEND"`
	Resume          *bool                    `ddl:"keyword" sql:"RESUME"`
	IfSuspended     *bool                    `ddl:"keyword" sql:"IF SUSPENDED"`
	AbortAllQueries *bool                    `ddl:"keyword" sql:"ABORT ALL QUERIES"`
	NewName         *AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`

	Set      *WarehouseSet      `ddl:"keyword" sql:"SET"`
	Unset    *WarehouseUnset    `ddl:"list,no_parentheses" sql:"UNSET"`
	SetTag   []TagAssociation   `ddl:"keyword" sql:"SET TAG"`
	UnsetTag []ObjectIdentifier `ddl:"keyword" sql:"UNSET TAG"`
}

func (opts *AlterWarehouseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Suspend, opts.Resume, opts.AbortAllQueries, opts.NewName, opts.Set, opts.Unset, opts.SetTag, opts.UnsetTag) {
		errs = append(errs, errExactlyOneOf("AlterWarehouseOptions", "Suspend", "Resume", "AbortAllQueries", "NewName", "Set", "Unset", "SetTag", "UnsetTag"))
	}
	if everyValueSet(opts.Suspend, opts.Resume) && (*opts.Suspend && *opts.Resume) {
		errs = append(errs, errOneOf("AlterWarehouseOptions", "Suspend", "Resume"))
	}
	if (valueSet(opts.IfSuspended) && *opts.IfSuspended) && (!valueSet(opts.Resume) || !*opts.Resume) {
		errs = append(errs, fmt.Errorf(`"Resume" has to be set when using "IfSuspended"`))
	}
	if valueSet(opts.Set) {
		if err := opts.Set.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Unset) {
		if err := opts.Unset.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type WarehouseSet struct {
	// Object properties
	WarehouseType                   *WarehouseType          `ddl:"parameter,single_quotes" sql:"WAREHOUSE_TYPE"`
	WarehouseSize                   *WarehouseSize          `ddl:"parameter,single_quotes" sql:"WAREHOUSE_SIZE"`
	WaitForCompletion               *bool                   `ddl:"parameter" sql:"WAIT_FOR_COMPLETION"`
	MaxClusterCount                 *int                    `ddl:"parameter" sql:"MAX_CLUSTER_COUNT"`
	MinClusterCount                 *int                    `ddl:"parameter" sql:"MIN_CLUSTER_COUNT"`
	ScalingPolicy                   *ScalingPolicy          `ddl:"parameter,single_quotes" sql:"SCALING_POLICY"`
	AutoSuspend                     *int                    `ddl:"parameter" sql:"AUTO_SUSPEND"`
	AutoResume                      *bool                   `ddl:"parameter" sql:"AUTO_RESUME"`
	ResourceMonitor                 AccountObjectIdentifier `ddl:"identifier,equals" sql:"RESOURCE_MONITOR"`
	Comment                         *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
	EnableQueryAcceleration         *bool                   `ddl:"parameter" sql:"ENABLE_QUERY_ACCELERATION"`
	QueryAccelerationMaxScaleFactor *int                    `ddl:"parameter" sql:"QUERY_ACCELERATION_MAX_SCALE_FACTOR"`

	// Object params
	MaxConcurrencyLevel             *int `ddl:"parameter" sql:"MAX_CONCURRENCY_LEVEL"`
	StatementQueuedTimeoutInSeconds *int `ddl:"parameter" sql:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds       *int `ddl:"parameter" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`
}

func (v *WarehouseSet) validate() error {
	// we validate only the case then both are set together, if only MinClusterCount is set, we leave it for Snowflake to validate
	if v.MinClusterCount != nil && valueSet(v.MaxClusterCount) {
		if ok := validateIntInRangeInclusive(*v.MinClusterCount, 1, *v.MaxClusterCount); !ok {
			return fmt.Errorf("MinClusterCount must be less than or equal to MaxClusterCount")
		}
	}
	if v.AutoSuspend != nil {
		if ok := validateIntGreaterThanOrEqual(*v.AutoSuspend, 0); !ok {
			return fmt.Errorf("AutoSuspend must be greater than or equal to 0")
		}
	}
	if v.QueryAccelerationMaxScaleFactor != nil {
		if ok := validateIntInRangeInclusive(*v.QueryAccelerationMaxScaleFactor, 0, 100); !ok {
			return fmt.Errorf("QueryAccelerationMaxScaleFactor must be between 0 and 100")
		}
	}
	if everyValueNil(v.WarehouseType, v.WarehouseSize, v.WaitForCompletion, v.MaxClusterCount, v.MinClusterCount, v.ScalingPolicy, v.AutoSuspend, v.AutoResume, v.ResourceMonitor, v.Comment, v.EnableQueryAcceleration, v.QueryAccelerationMaxScaleFactor, v.MaxConcurrencyLevel, v.StatementQueuedTimeoutInSeconds, v.StatementTimeoutInSeconds) {
		return errAtLeastOneOf("WarehouseSet", "WarehouseType", "WarehouseSize", "WaitForCompletion", "MaxClusterCount", "MinClusterCount", "ScalingPolicy", "AutoSuspend", "AutoResume", "ResourceMonitor", "Comment", "EnableQueryAcceleration", "QueryAccelerationMaxScaleFactor", "MaxConcurrencyLevel", "StatementQueuedTimeoutInSeconds", "StatementTimeoutInSeconds")
	}
	return nil
}

type WarehouseUnset struct {
	// Object properties
	WarehouseType                   *bool `ddl:"keyword" sql:"WAREHOUSE_TYPE"`
	WaitForCompletion               *bool `ddl:"keyword" sql:"WAIT_FOR_COMPLETION"`
	MaxClusterCount                 *bool `ddl:"keyword" sql:"MAX_CLUSTER_COUNT"`
	MinClusterCount                 *bool `ddl:"keyword" sql:"MIN_CLUSTER_COUNT"`
	ScalingPolicy                   *bool `ddl:"keyword" sql:"SCALING_POLICY"`
	AutoSuspend                     *bool `ddl:"keyword" sql:"AUTO_SUSPEND"`
	AutoResume                      *bool `ddl:"keyword" sql:"AUTO_RESUME"`
	ResourceMonitor                 *bool `ddl:"keyword" sql:"RESOURCE_MONITOR"`
	Comment                         *bool `ddl:"keyword" sql:"COMMENT"`
	EnableQueryAcceleration         *bool `ddl:"keyword" sql:"ENABLE_QUERY_ACCELERATION"`
	QueryAccelerationMaxScaleFactor *bool `ddl:"keyword" sql:"QUERY_ACCELERATION_MAX_SCALE_FACTOR"`

	// Object params
	MaxConcurrencyLevel             *bool `ddl:"keyword" sql:"MAX_CONCURRENCY_LEVEL"`
	StatementQueuedTimeoutInSeconds *bool `ddl:"keyword" sql:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds       *bool `ddl:"keyword" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`
}

func (v *WarehouseUnset) validate() error {
	if everyValueNil(v.WarehouseType, v.WaitForCompletion, v.MaxClusterCount, v.MinClusterCount, v.ScalingPolicy, v.AutoSuspend, v.AutoResume, v.ResourceMonitor, v.Comment, v.EnableQueryAcceleration, v.QueryAccelerationMaxScaleFactor, v.MaxConcurrencyLevel, v.StatementQueuedTimeoutInSeconds, v.StatementTimeoutInSeconds) {
		return errAtLeastOneOf("WarehouseUnset", "WarehouseType", "WaitForCompletion", "MaxClusterCount", "MinClusterCount", "ScalingPolicy", "AutoSuspend", "AutoResume", "ResourceMonitor", "Comment", "EnableQueryAcceleration", "QueryAccelerationMaxScaleFactor", "MaxConcurrencyLevel", "StatementQueuedTimeoutInSeconds", "StatementTimeoutInSeconds")
	}
	return nil
}

func (c *warehouses) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterWarehouseOptions) error {
	if opts == nil {
		opts = &AlterWarehouseOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}

	// Warehouse needs to be suspended to change its type.
	if opts.warehouseTypeIsChanged() {
		warehouse, err := c.ShowByID(ctx, id)
		if err != nil {
			return err
		}
		if warehouse.State == WarehouseStateStarted {
			err := c.Alter(ctx, id, &AlterWarehouseOptions{Suspend: Bool(true)})
			if err != nil {
				return err
			}
			defer func() {
				err := c.Alter(ctx, id, &AlterWarehouseOptions{Resume: Bool(true), IfSuspended: Bool(true)})
				if err != nil {
					log.Printf("[DEBUG] error occurred during warehouse resumption, err=%v", err)
				}
			}()

			// needed to make sure that warehouse is suspended
			var warehouseSuspensionErrs []error
			err = util.Retry(3, 1*time.Second, func() (error, bool) {
				warehouse, err = c.ShowByID(ctx, id)
				if err != nil {
					warehouseSuspensionErrs = append(warehouseSuspensionErrs, err)
					return nil, false
				}
				if warehouse.State != WarehouseStateSuspended {
					return nil, false
				}
				return nil, true
			})
			if err != nil {
				return fmt.Errorf("warehouse suspension failed, err: %w, original errors: %w", err, errors.Join(warehouseSuspensionErrs...))
			}
		}
	}

	_, err = c.client.exec(ctx, sql)
	return err
}

func (opts *AlterWarehouseOptions) warehouseTypeIsChanged() bool {
	return opts.Set != nil && opts.Set.WarehouseType != nil
}

// DropWarehouseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-warehouse.
type DropWarehouseOptions struct {
	drop      bool                    `ddl:"static" sql:"DROP"`
	warehouse bool                    `ddl:"static" sql:"WAREHOUSE"`
	IfExists  *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name      AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *DropWarehouseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (c *warehouses) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropWarehouseOptions) error {
	if opts == nil {
		opts = &DropWarehouseOptions{
			name: id,
		}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, sql)
	if err != nil {
		return err
	}
	return err
}

func (c *warehouses) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(c.client, func() error { return c.Drop(ctx, id, &DropWarehouseOptions{IfExists: Bool(true)}) }, ctx, id)
}

// ShowWarehouseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-warehouses.
type ShowWarehouseOptions struct {
	show       bool  `ddl:"static" sql:"SHOW"`
	warehouses bool  `ddl:"static" sql:"WAREHOUSES"`
	Like       *Like `ddl:"keyword" sql:"LIKE"`
}

func (opts *ShowWarehouseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

type WarehouseState string

const (
	WarehouseStateSuspended  WarehouseState = "SUSPENDED"
	WarehouseStateSuspending WarehouseState = "SUSPENDING"
	WarehouseStateStarted    WarehouseState = "STARTED"
	WarehouseStateResizing   WarehouseState = "RESIZING"
	WarehouseStateResuming   WarehouseState = "RESUMING"
)

type Warehouse struct {
	Name                            string
	State                           WarehouseState
	Type                            WarehouseType
	Size                            WarehouseSize
	MinClusterCount                 int
	MaxClusterCount                 int
	StartedClusters                 int
	Running                         int
	Queued                          int
	IsDefault                       bool
	IsCurrent                       bool
	AutoSuspend                     int
	AutoResume                      bool
	Available                       float64
	Provisioning                    float64
	Quiescing                       float64
	Other                           float64
	CreatedOn                       time.Time
	ResumedOn                       time.Time
	UpdatedOn                       time.Time
	Owner                           string
	Comment                         string
	EnableQueryAcceleration         bool
	QueryAccelerationMaxScaleFactor int
	ResourceMonitor                 AccountObjectIdentifier
	ScalingPolicy                   ScalingPolicy
	OwnerRoleType                   string
}

type warehouseDBRow struct {
	Name                            string         `db:"name"`
	State                           string         `db:"state"`
	Type                            string         `db:"type"`
	Size                            string         `db:"size"`
	MinClusterCount                 int            `db:"min_cluster_count"`
	MaxClusterCount                 int            `db:"max_cluster_count"`
	StartedClusters                 int            `db:"started_clusters"`
	Running                         int            `db:"running"`
	Queued                          int            `db:"queued"`
	IsDefault                       string         `db:"is_default"`
	IsCurrent                       string         `db:"is_current"`
	AutoSuspend                     sql.NullInt64  `db:"auto_suspend"`
	AutoResume                      bool           `db:"auto_resume"`
	Available                       string         `db:"available"`
	Provisioning                    string         `db:"provisioning"`
	Quiescing                       string         `db:"quiescing"`
	Other                           string         `db:"other"`
	CreatedOn                       time.Time      `db:"created_on"`
	ResumedOn                       time.Time      `db:"resumed_on"`
	UpdatedOn                       time.Time      `db:"updated_on"`
	Owner                           string         `db:"owner"`
	Comment                         string         `db:"comment"`
	EnableQueryAcceleration         bool           `db:"enable_query_acceleration"`
	QueryAccelerationMaxScaleFactor int            `db:"query_acceleration_max_scale_factor"`
	ResourceMonitor                 string         `db:"resource_monitor"`
	Actives                         string         `db:"actives"`
	Pendings                        string         `db:"pendings"`
	Failed                          string         `db:"failed"`
	Suspended                       string         `db:"suspended"`
	UUID                            string         `db:"uuid"`
	ScalingPolicy                   string         `db:"scaling_policy"`
	OwnerRoleType                   sql.NullString `db:"owner_role_type"`
}

func (row warehouseDBRow) convertErr() (*Warehouse, error) {
	wh := &Warehouse{
		Name:                            row.Name,
		State:                           WarehouseState(row.State),
		Type:                            WarehouseType(row.Type),
		MinClusterCount:                 row.MinClusterCount,
		MaxClusterCount:                 row.MaxClusterCount,
		StartedClusters:                 row.StartedClusters,
		Running:                         row.Running,
		Queued:                          row.Queued,
		IsDefault:                       row.IsDefault == "Y",
		IsCurrent:                       row.IsCurrent == "Y",
		AutoResume:                      row.AutoResume,
		CreatedOn:                       row.CreatedOn,
		ResumedOn:                       row.ResumedOn,
		UpdatedOn:                       row.UpdatedOn,
		Owner:                           row.Owner,
		Comment:                         row.Comment,
		EnableQueryAcceleration:         row.EnableQueryAcceleration,
		QueryAccelerationMaxScaleFactor: row.QueryAccelerationMaxScaleFactor,
		ScalingPolicy:                   ScalingPolicy(row.ScalingPolicy),
	}
	if size, err := ToWarehouseSize(row.Size); err != nil {
		return nil, err
	} else {
		wh.Size = size
	}
	if available := strings.TrimSpace(row.Available); available != "" {
		if val, err := strconv.ParseFloat(available, 64); err != nil {
			return nil, fmt.Errorf(`row 'available' has incorrect value '%s', %w`, available, err)
		} else {
			wh.Available = val
		}
	}
	if provisioning := strings.TrimSpace(row.Provisioning); provisioning != "" {
		if val, err := strconv.ParseFloat(provisioning, 64); err != nil {
			return nil, fmt.Errorf(`row 'provisioning' has incorrect value '%s', %w`, provisioning, err)
		} else {
			wh.Provisioning = val
		}
	}
	if quiescing := strings.TrimSpace(row.Quiescing); quiescing != "" {
		if val, err := strconv.ParseFloat(quiescing, 64); err != nil {
			return nil, fmt.Errorf(`row 'quiescing' has incorrect value '%s', %w`, quiescing, err)
		} else {
			wh.Quiescing = val
		}
	}
	if other := strings.TrimSpace(row.Other); other != "" {
		if val, err := strconv.ParseFloat(other, 64); err != nil {
			return nil, fmt.Errorf(`row 'other' has incorrect value '%s', %w`, other, err)
		} else {
			wh.Other = val
		}
	}
	if row.AutoSuspend.Valid {
		wh.AutoSuspend = int(row.AutoSuspend.Int64)
	}
	if row.OwnerRoleType.Valid {
		wh.OwnerRoleType = row.OwnerRoleType.String
	}
	if row.ResourceMonitor != "null" {
		wh.ResourceMonitor = NewAccountObjectIdentifierFromFullyQualifiedName(row.ResourceMonitor)
	}
	return wh, nil
}

func (c *warehouses) Show(ctx context.Context, opts *ShowWarehouseOptions) ([]Warehouse, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[warehouseDBRow](c.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRowsErr[warehouseDBRow, Warehouse](dbRows)
}

func (c *warehouses) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Warehouse, error) {
	warehouses, err := c.Show(ctx, &ShowWarehouseOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}

	return collections.FindFirst(warehouses, func(warehouse Warehouse) bool {
		return warehouse.ID().FullyQualifiedName() == id.FullyQualifiedName()
	})
}

func (c *warehouses) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Warehouse, error) {
	return SafeShowById(c.client, c.ShowByID, ctx, id)
}

// describeWarehouseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-warehouse.
type describeWarehouseOptions struct {
	describe  bool                    `ddl:"static" sql:"DESCRIBE"`
	warehouse bool                    `ddl:"static" sql:"WAREHOUSE"`
	name      AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *describeWarehouseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

type warehouseDetailsRow struct {
	CreatedOn time.Time `db:"created_on"`
	Name      string    `db:"name"`
	Kind      string    `db:"kind"`
}

func (row *warehouseDetailsRow) toWarehouseDetails() *WarehouseDetails {
	return &WarehouseDetails{
		CreatedOn: row.CreatedOn,
		Name:      row.Name,
		Kind:      row.Kind,
	}
}

type WarehouseDetails struct {
	CreatedOn time.Time
	Name      string
	Kind      string
}

func (c *warehouses) Describe(ctx context.Context, id AccountObjectIdentifier) (*WarehouseDetails, error) {
	opts := &describeWarehouseOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}

	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := warehouseDetailsRow{}
	err = c.client.queryOne(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}

	return dest.toWarehouseDetails(), nil
}

func (v *Warehouse) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *Warehouse) ObjectType() ObjectType {
	return ObjectTypeWarehouse
}

func (c *warehouses) ShowParameters(ctx context.Context, id AccountObjectIdentifier) ([]*Parameter, error) {
	return c.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Warehouse: id,
		},
	})
}

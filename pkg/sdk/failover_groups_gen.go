package sdk

import (
	"context"
	"database/sql"
	"time"
)

var (
	_ FailoverGroups                = (*failoverGroups)(nil)
	_ convertibleRow[FailoverGroup] = new(failoverGroupDBRow)
)

type FailoverGroups interface {
	Create(ctx context.Context, id AccountObjectIdentifier, objectTypes []PluralObjectType, allowedAccounts []AccountIdentifier, opts *CreateFailoverGroupOptions) error
	CreateSecondaryReplicationGroup(ctx context.Context, id AccountObjectIdentifier, primaryFailoverGroupID ExternalObjectIdentifier, opts *CreateSecondaryReplicationGroupOptions) error
	AlterSource(ctx context.Context, id AccountObjectIdentifier, opts *AlterSourceFailoverGroupOptions) error
	AlterTarget(ctx context.Context, id AccountObjectIdentifier, opts *AlterTargetFailoverGroupOptions) error
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropFailoverGroupOptions) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, opts *ShowFailoverGroupOptions) ([]FailoverGroup, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*FailoverGroup, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*FailoverGroup, error)
	ShowDatabases(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error)
	ShowShares(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error)
}

// failoverGroups implements FailoverGroups.
type failoverGroups struct {
	client *Client
}

// CreateFailoverGroupOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-failover-group.
type CreateFailoverGroupOptions struct {
	create        bool                    `ddl:"static" sql:"CREATE"`
	failoverGroup bool                    `ddl:"static" sql:"FAILOVER GROUP"`
	IfNotExists   *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`

	objectTypes             []PluralObjectType        `ddl:"parameter" sql:"OBJECT_TYPES"`
	AllowedDatabases        []AccountObjectIdentifier `ddl:"parameter" sql:"ALLOWED_DATABASES"`
	AllowedShares           []AccountObjectIdentifier `ddl:"parameter" sql:"ALLOWED_SHARES"`
	AllowedIntegrationTypes []IntegrationType         `ddl:"parameter" sql:"ALLOWED_INTEGRATION_TYPES"`
	allowedAccounts         []AccountIdentifier       `ddl:"parameter" sql:"ALLOWED_ACCOUNTS"`
	IgnoreEditionCheck      *bool                     `ddl:"keyword" sql:"IGNORE EDITION CHECK"`
	ReplicationSchedule     *string                   `ddl:"parameter,single_quotes" sql:"REPLICATION_SCHEDULE"`
}

// CreateSecondaryReplicationGroupOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-failover-group.
type CreateSecondaryReplicationGroupOptions struct {
	create               bool                     `ddl:"static" sql:"CREATE"`
	failoverGroup        bool                     `ddl:"static" sql:"FAILOVER GROUP"`
	IfNotExists          *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                 AccountObjectIdentifier  `ddl:"identifier"`
	primaryFailoverGroup ExternalObjectIdentifier `ddl:"identifier" sql:"AS REPLICA OF"`
}

// AlterSourceFailoverGroupOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-failover-group.
type AlterSourceFailoverGroupOptions struct {
	alter         bool                    `ddl:"static" sql:"ALTER"`
	failoverGroup bool                    `ddl:"static" sql:"FAILOVER GROUP"`
	IfExists      *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`
	NewName       AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	Set           *FailoverGroupSet       `ddl:"keyword" sql:"SET"`
	Unset         *FailoverGroupUnset     `ddl:"list,no_parentheses" sql:"UNSET"`
	Add           *FailoverGroupAdd       `ddl:"keyword" sql:"ADD"`
	Move          *FailoverGroupMove      `ddl:"keyword" sql:"MOVE"`
	Remove        *FailoverGroupRemove    `ddl:"keyword" sql:"REMOVE"`
}

type FailoverGroupSet struct {
	ObjectTypes             []PluralObjectType `ddl:"parameter" sql:"OBJECT_TYPES"`
	AllowedIntegrationTypes []IntegrationType  `ddl:"parameter" sql:"ALLOWED_INTEGRATION_TYPES"`
	ReplicationSchedule     *string            `ddl:"parameter,single_quotes" sql:"REPLICATION_SCHEDULE"`
}

type FailoverGroupUnset struct {
	ReplicationSchedule *bool `ddl:"keyword" sql:"REPLICATION_SCHEDULE"`
}

type FailoverGroupAdd struct {
	AllowedDatabases   []AccountObjectIdentifier `ddl:"parameter,reverse" sql:"TO ALLOWED_DATABASES"`
	AllowedShares      []AccountObjectIdentifier `ddl:"parameter,reverse" sql:"TO ALLOWED_SHARES"`
	AllowedAccounts    []AccountIdentifier       `ddl:"parameter,reverse" sql:"TO ALLOWED_ACCOUNTS"`
	IgnoreEditionCheck *bool                     `ddl:"keyword" sql:"IGNORE_EDITION_CHECK"`
}

type FailoverGroupMove struct {
	Databases []AccountObjectIdentifier `ddl:"parameter,no_equals" sql:"DATABASES"`
	Shares    []AccountObjectIdentifier `ddl:"parameter,no_equals" sql:"SHARES"`
	To        AccountObjectIdentifier   `ddl:"identifier" sql:"TO FAILOVER GROUP"`
}

type FailoverGroupRemove struct {
	AllowedDatabases []AccountObjectIdentifier `ddl:"parameter,reverse" sql:"FROM ALLOWED_DATABASES"`
	AllowedShares    []AccountObjectIdentifier `ddl:"parameter,reverse" sql:"FROM ALLOWED_SHARES"`
	AllowedAccounts  []AccountIdentifier       `ddl:"parameter,reverse" sql:"FROM ALLOWED_ACCOUNTS"`
}

// AlterTargetFailoverGroupOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-failover-group.
type AlterTargetFailoverGroupOptions struct {
	alter         bool                    `ddl:"static" sql:"ALTER"`
	failoverGroup bool                    `ddl:"static" sql:"FAILOVER GROUP"`
	IfExists      *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`
	Refresh       *bool                   `ddl:"keyword" sql:"REFRESH"`
	Primary       *bool                   `ddl:"keyword" sql:"PRIMARY"`
	Suspend       *bool                   `ddl:"keyword" sql:"SUSPEND"`
	Resume        *bool                   `ddl:"keyword" sql:"RESUME"`
}

// DropFailoverGroupOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-failover-group.
type DropFailoverGroupOptions struct {
	drop          bool                    `ddl:"static" sql:"DROP"`
	failoverGroup bool                    `ddl:"static" sql:"FAILOVER GROUP"`
	IfExists      *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`
}

// ShowFailoverGroupOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-failover-groups.
type ShowFailoverGroupOptions struct {
	show           bool              `ddl:"static" sql:"SHOW"`
	failoverGroups bool              `ddl:"static" sql:"FAILOVER GROUPS"`
	InAccount      AccountIdentifier `ddl:"identifier" sql:"IN ACCOUNT"`
}

// FailoverGroup is a user friendly result for a CREATE FAILOVER GROUP query.
type FailoverGroup struct {
	RegionGroup             string
	SnowflakeRegion         string
	CreatedOn               time.Time
	AccountName             string
	Name                    string
	Type                    string
	Comment                 string
	IsPrimary               bool
	Primary                 ExternalObjectIdentifier
	ObjectTypes             []PluralObjectType
	AllowedIntegrationTypes []IntegrationType
	AllowedAccounts         []AccountIdentifier
	OrganizationName        string
	AccountLocator          string
	ReplicationSchedule     string
	SecondaryState          FailoverGroupSecondaryState
	NextScheduledRefresh    string
	Owner                   string
}

func (v *FailoverGroup) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *FailoverGroup) ObjectType() ObjectType {
	return ObjectTypeFailoverGroup
}

// failoverGroupDBRow is used to decode the result of a CREATE FAILOVER GROUP query.
type failoverGroupDBRow struct {
	RegionGroup             string         `db:"region_group"`
	SnowflakeRegion         string         `db:"snowflake_region"`
	CreatedOn               time.Time      `db:"created_on"`
	AccountName             string         `db:"account_name"`
	Name                    string         `db:"name"`
	Type                    string         `db:"type"`
	Comment                 sql.NullString `db:"comment"`
	IsPrimary               bool           `db:"is_primary"`
	Primary                 string         `db:"primary"`
	ObjectTypes             string         `db:"object_types"`
	AllowedIntegrationTypes string         `db:"allowed_integration_types"`
	AllowedAccounts         string         `db:"allowed_accounts"`
	OrganizationName        string         `db:"organization_name"`
	AccountLocator          string         `db:"account_locator"`
	ReplicationSchedule     sql.NullString `db:"replication_schedule"`
	SecondaryState          sql.NullString `db:"secondary_state"`
	NextScheduledRefresh    sql.NullString `db:"next_scheduled_refresh"`
	Owner                   sql.NullString `db:"owner"`
}

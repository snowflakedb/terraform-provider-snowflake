package sdk

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ FailoverGroups = (*failoverGroups)(nil)
var _ convertibleRow[FailoverGroup] = new(failoverGroupDBRow)

var (
	_ validatable = new(CreateDatabaseOptions)
	_ validatable = new(CreateSharedDatabaseOptions)
	_ validatable = new(CreateSecondaryDatabaseOptions)
	_ validatable = new(AlterDatabaseOptions)
	_ validatable = new(AlterDatabaseReplicationOptions)
	_ validatable = new(AlterDatabaseFailoverOptions)
	_ validatable = new(DropDatabaseOptions)
	_ validatable = new(undropDatabaseOptions)
	_ validatable = new(ShowDatabasesOptions)
	_ validatable = new(describeDatabaseOptions)
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

// FailoverGroups implements FailoverGroups.
type failoverGroups struct {
	client *Client
}

// IntegrationType is the type of integration.
type IntegrationType string

const (
	IntegrationTypeSecurityIntegrations       IntegrationType = "SECURITY INTEGRATIONS"
	IntegrationTypeAPIIntegrations            IntegrationType = "API INTEGRATIONS"
	IntegrationTypeStorageIntegrations        IntegrationType = "STORAGE INTEGRATIONS"
	IntegrationTypeExternalAccessIntegrations IntegrationType = "EXTERNAL ACCESS INTEGRATIONS"
	IntegrationTypeNotificationIntegrations   IntegrationType = "NOTIFICATION INTEGRATIONS"
)

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

func (opts *CreateFailoverGroupOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *failoverGroups) Create(ctx context.Context, id AccountObjectIdentifier, objectTypes []PluralObjectType, allowedAccounts []AccountIdentifier, opts *CreateFailoverGroupOptions) error {
	if opts == nil {
		opts = &CreateFailoverGroupOptions{}
	}
	opts.name = id
	opts.allowedAccounts = allowedAccounts
	opts.objectTypes = objectTypes
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

// CreateSecondaryReplicationGroupOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-failover-group.
type CreateSecondaryReplicationGroupOptions struct {
	create               bool                     `ddl:"static" sql:"CREATE"`
	failoverGroup        bool                     `ddl:"static" sql:"FAILOVER GROUP"`
	IfNotExists          *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                 AccountObjectIdentifier  `ddl:"identifier"`
	primaryFailoverGroup ExternalObjectIdentifier `ddl:"identifier" sql:"AS REPLICA OF"`
}

func (opts *CreateSecondaryReplicationGroupOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.primaryFailoverGroup) {
		errs = append(errs, errInvalidIdentifier("CreateSecondaryReplicationGroupOptions", "primaryFailoverGroup"))
	}
	return errors.Join(errs...)
}

func (v *failoverGroups) CreateSecondaryReplicationGroup(ctx context.Context, id AccountObjectIdentifier, primaryFailoverGroupID ExternalObjectIdentifier, opts *CreateSecondaryReplicationGroupOptions) error {
	if opts == nil {
		opts = &CreateSecondaryReplicationGroupOptions{}
	}
	opts.name = id
	opts.primaryFailoverGroup = primaryFailoverGroupID
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
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

func (opts *AlterSourceFailoverGroupOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.Add, opts.Move, opts.Remove, opts.NewName) {
		errs = append(errs, errExactlyOneOf("AlterSourceFailoverGroupOptions", "Set", "Unset", "Add", "Move", "Remove", "NewName"))
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
	if valueSet(opts.Add) {
		if err := opts.Add.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Move) {
		if err := opts.Move.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Remove) {
		if err := opts.Remove.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type FailoverGroupSet struct {
	ObjectTypes             []PluralObjectType `ddl:"parameter" sql:"OBJECT_TYPES"`
	AllowedIntegrationTypes []IntegrationType  `ddl:"parameter" sql:"ALLOWED_INTEGRATION_TYPES"`
	ReplicationSchedule     *string            `ddl:"parameter,single_quotes" sql:"REPLICATION_SCHEDULE"`
}

func (v *FailoverGroupSet) validate() error {
	if len(v.AllowedIntegrationTypes) > 0 {
		// INTEGRATIONS must be set in object types
		if !slices.Contains(v.ObjectTypes, PluralObjectTypeIntegrations) {
			return errors.New("INTEGRATIONS must be set in OBJECT_TYPES when setting allowed integration types")
		}
	}
	return nil
}

type FailoverGroupUnset struct {
	ReplicationSchedule *bool `ddl:"keyword" sql:"REPLICATION_SCHEDULE"`
}

func (v *FailoverGroupUnset) validate() error {
	if everyValueNil(v.ReplicationSchedule) {
		return errAtLeastOneOf("FailoverGroupUnset", "ReplicationSchedule")
	}
	return nil
}

type FailoverGroupAdd struct {
	AllowedDatabases   []AccountObjectIdentifier `ddl:"parameter,reverse" sql:"TO ALLOWED_DATABASES"`
	AllowedShares      []AccountObjectIdentifier `ddl:"parameter,reverse" sql:"TO ALLOWED_SHARES"`
	AllowedAccounts    []AccountIdentifier       `ddl:"parameter,reverse" sql:"TO ALLOWED_ACCOUNTS"`
	IgnoreEditionCheck *bool                     `ddl:"keyword" sql:"IGNORE_EDITION_CHECK"`
}

func (v *FailoverGroupAdd) validate() error {
	return nil
}

type FailoverGroupMove struct {
	Databases []AccountObjectIdentifier `ddl:"parameter,no_equals" sql:"DATABASES"`
	Shares    []AccountObjectIdentifier `ddl:"parameter,no_equals" sql:"SHARES"`
	To        AccountObjectIdentifier   `ddl:"identifier" sql:"TO FAILOVER GROUP"`
}

func (v *FailoverGroupMove) validate() error {
	return nil
}

type FailoverGroupRemove struct {
	AllowedDatabases []AccountObjectIdentifier `ddl:"parameter,reverse" sql:"FROM ALLOWED_DATABASES"`
	AllowedShares    []AccountObjectIdentifier `ddl:"parameter,reverse" sql:"FROM ALLOWED_SHARES"`
	AllowedAccounts  []AccountIdentifier       `ddl:"parameter,reverse" sql:"FROM ALLOWED_ACCOUNTS"`
}

func (v *FailoverGroupRemove) validate() error {
	return nil
}

func (v *failoverGroups) AlterSource(ctx context.Context, id AccountObjectIdentifier, opts *AlterSourceFailoverGroupOptions) error {
	if opts == nil {
		opts = &AlterSourceFailoverGroupOptions{}
	}
	opts.name = id

	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
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

func (opts *AlterTargetFailoverGroupOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Refresh, opts.Primary, opts.Suspend, opts.Resume) {
		errs = append(errs, errExactlyOneOf("AlterTargetFailoverGroupOptions", "Refresh", "Primary", "Suspend", "Resume"))
	}
	return errors.Join(errs...)
}

func (v *failoverGroups) AlterTarget(ctx context.Context, id AccountObjectIdentifier, opts *AlterTargetFailoverGroupOptions) error {
	if opts == nil {
		opts = &AlterTargetFailoverGroupOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

// DropFailoverGroupOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-failover-group.
type DropFailoverGroupOptions struct {
	drop          bool                    `ddl:"static" sql:"DROP"`
	failoverGroup bool                    `ddl:"static" sql:"FAILOVER GROUP"`
	IfExists      *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *DropFailoverGroupOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *failoverGroups) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropFailoverGroupOptions) error {
	if opts == nil {
		opts = &DropFailoverGroupOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	if err != nil {
		return err
	}
	return err
}

func (v *failoverGroups) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropFailoverGroupOptions{IfExists: Bool(true)}) }, ctx, id)
}

// ShowFailoverGroupOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-failover-groups.
type ShowFailoverGroupOptions struct {
	show           bool              `ddl:"static" sql:"SHOW"`
	failoverGroups bool              `ddl:"static" sql:"FAILOVER GROUPS"`
	InAccount      AccountIdentifier `ddl:"identifier" sql:"IN ACCOUNT"`
}

func (opts *ShowFailoverGroupOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

type FailoverGroupSecondaryState string

const (
	FailoverGroupSecondaryStateSuspended FailoverGroupSecondaryState = "SUSPENDED"
	FailoverGroupSecondaryStateStarted   FailoverGroupSecondaryState = "STARTED"
	FailoverGroupSecondaryStateNull      FailoverGroupSecondaryState = "NULL"
)

// FailoverGroups is a user friendly result for a CREATE FAILOVER GROUP query.
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

func (v *FailoverGroup) ExternalID() ExternalObjectIdentifier {
	return NewExternalObjectIdentifier(AccountIdentifier{
		organizationName: v.OrganizationName,
		accountName:      v.AccountName,
		accountLocator:   v.AccountLocator,
	}, v.ID())
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

func (row failoverGroupDBRow) convertErr() (*FailoverGroup, error) {
	ots := strings.Split(row.ObjectTypes, ",")
	pluralObjectTypes := make([]PluralObjectType, 0, len(ots))
	for _, ot := range ots {
		pot := PluralObjectType(strings.TrimSpace(ot))
		if pot == PluralObjectTypeParameters {
			pluralObjectTypes = append(pluralObjectTypes, PluralObjectType("ACCOUNT PARAMETERS"))
		} else {
			pluralObjectTypes = append(pluralObjectTypes, pot)
		}
	}
	its := strings.Split(row.AllowedIntegrationTypes, ",")
	allowedIntegrationTypes := make([]IntegrationType, 0, len(its))
	for _, it := range its {
		if it == "" {
			continue
		}
		allowedIntegrationTypes = append(allowedIntegrationTypes, IntegrationType(strings.ReplaceAll(strings.TrimSpace(it), "_", " ")+" INTEGRATIONS"))
	}
	aas := strings.Split(row.AllowedAccounts, ",")
	allowedAccounts := make([]AccountIdentifier, 0, len(aas))
	for _, aa := range aas {
		s := strings.TrimSpace(aa)
		p := strings.Split(s, ".")
		if len(p) != 2 {
			continue
		}
		orgName := p[0]
		accountName := p[1]
		allowedAccounts = append(allowedAccounts, NewAccountIdentifier(orgName, accountName))
	}
	var comment string
	if row.Comment.Valid {
		comment = row.Comment.String
	}
	var replicationSchedule string
	if row.ReplicationSchedule.Valid {
		replicationSchedule = row.ReplicationSchedule.String
	}

	secondaryState := FailoverGroupSecondaryStateNull
	if row.SecondaryState.Valid {
		secondaryState = FailoverGroupSecondaryState(row.SecondaryState.String)
	}
	nextScheduledRefresh := ""
	if row.NextScheduledRefresh.Valid {
		nextScheduledRefresh = row.NextScheduledRefresh.String
	}
	return &FailoverGroup{
		RegionGroup:             row.RegionGroup,
		SnowflakeRegion:         row.SnowflakeRegion,
		CreatedOn:               row.CreatedOn,
		AccountName:             row.AccountName,
		OrganizationName:        row.OrganizationName,
		AccountLocator:          row.AccountLocator,
		Name:                    row.Name,
		Comment:                 comment,
		IsPrimary:               row.IsPrimary,
		Primary:                 NewExternalObjectIdentifierFromFullyQualifiedName(row.Primary),
		ObjectTypes:             pluralObjectTypes,
		AllowedIntegrationTypes: allowedIntegrationTypes,
		AllowedAccounts:         allowedAccounts,
		ReplicationSchedule:     replicationSchedule,
		SecondaryState:          secondaryState,
		NextScheduledRefresh:    nextScheduledRefresh,
		Owner:                   row.Owner.String,
		Type:                    row.Type,
	}, nil
}

func (v *failoverGroups) Show(ctx context.Context, opts *ShowFailoverGroupOptions) ([]FailoverGroup, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[failoverGroupDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRowsErr[failoverGroupDBRow, FailoverGroup](dbRows)
}

func (v *failoverGroups) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*FailoverGroup, error) {
	currentAccount, err := v.client.ContextFunctions.CurrentAccount(ctx)
	if err != nil {
		return nil, err
	}

	failoverGroups, err := v.Show(ctx, nil)
	if err != nil {
		return nil, err
	}

	return collections.FindFirst(failoverGroups, func(group FailoverGroup) bool {
		return group.ID().FullyQualifiedName() == id.FullyQualifiedName() && group.AccountLocator == currentAccount
	})
}

func (v *failoverGroups) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*FailoverGroup, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

// showFailoverGroupDatabasesOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-databases-in-failover-group.
type showFailoverGroupDatabasesOptions struct {
	show      bool                    `ddl:"static" sql:"SHOW"`
	databases bool                    `ddl:"static" sql:"DATABASES"`
	in        AccountObjectIdentifier `ddl:"identifier" sql:"IN FAILOVER GROUP"`
}

func (opts *showFailoverGroupDatabasesOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.in) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *failoverGroups) ShowDatabases(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error) {
	opts := &showFailoverGroupDatabasesOptions{
		in: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []struct {
		Name string `db:"name"`
	}{}
	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	resultList := make([]AccountObjectIdentifier, len(dest))
	for i, row := range dest {
		resultList[i] = NewAccountObjectIdentifier(row.Name)
	}
	return resultList, nil
}

// showFailoverGroupSharesOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-shares-in-failover-group.
type showFailoverGroupSharesOptions struct {
	show      bool                    `ddl:"static" sql:"SHOW"`
	databases bool                    `ddl:"static" sql:"SHARES"`
	in        AccountObjectIdentifier `ddl:"identifier" sql:"IN FAILOVER GROUP"`
}

func (opts *showFailoverGroupSharesOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.in) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *failoverGroups) ShowShares(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error) {
	opts := &showFailoverGroupSharesOptions{
		in: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []struct {
		Name         string `db:"name"`
		OwnerAccount string `db:"owner_account"`
	}{}
	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	resultList := make([]AccountObjectIdentifier, len(dest))
	for i, r := range dest {
		// TODO [SNOW-1348343]: change during failover groups rework; this was not working correctly with identifiers containing `.` character
		resultList[i] = NewExternalObjectIdentifier(NewAccountIdentifierFromFullyQualifiedName(r.OwnerAccount), NewAccountObjectIdentifier(r.Name)).objectIdentifier.(AccountObjectIdentifier)
	}
	return resultList, nil
}

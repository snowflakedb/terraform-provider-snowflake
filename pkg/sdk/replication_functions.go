package sdk

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

var _ ReplicationFunctions = (*replicationFunctions)(nil)

var (
	_ validatable = new(ShowRegionsOptions)
	_ validatable = new(ShowReplicationDatabasesOptions)
)

var _ convertibleRowDeprecated[ReplicationDatabase] = new(replicationDatabaseRow)

type ReplicationFunctions interface {
	ShowReplicationAccounts(ctx context.Context) ([]*ReplicationAccount, error)
	ShowReplicationDatabases(ctx context.Context, opts *ShowReplicationDatabasesOptions) ([]ReplicationDatabase, error)
	ShowRegions(ctx context.Context, opts *ShowRegionsOptions) ([]*Region, error)
}

type replicationFunctions struct {
	client *Client
}

type ReplicationAccount struct {
	SnowflakeRegion  string    `db:"snowflake_region"`
	CreatedOn        time.Time `db:"created_on"`
	AccountName      string    `db:"account_name"`
	AccountLocator   string    `db:"account_locator"`
	Comment          string    `db:"comment"`
	OrganizationName string    `db:"organization_name"`
	IsOrgAdmin       bool      `db:"is_org_admin"`
}

func (v *ReplicationAccount) ID() AccountIdentifier {
	return AccountIdentifier{
		organizationName: v.OrganizationName,
		accountName:      v.AccountName,
		accountLocator:   v.AccountLocator,
	}
}

// ShowReplicationAccounts is based on https://docs.snowflake.com/en/sql-reference/sql/show-replication-accounts.
func (c *replicationFunctions) ShowReplicationAccounts(ctx context.Context) ([]*ReplicationAccount, error) {
	rows := []ReplicationAccount{}
	sql := "SHOW REPLICATION ACCOUNTS"
	err := c.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	replicationAccounts := make([]*ReplicationAccount, len(rows))
	for i, row := range rows {
		replicationAccount := row
		replicationAccounts[i] = &replicationAccount
	}
	return replicationAccounts, nil
}

type replicationDatabaseRow struct {
	RegionGroup                  sql.NullString `db:"region_group"`
	SnowflakeRegion              string         `db:"snowflake_region"`
	CreatedOn                    string         `db:"created_on"`
	AccountName                  string         `db:"account_name"`
	Name                         string         `db:"name"`
	Comment                      sql.NullString `db:"comment"`
	IsPrimary                    bool           `db:"is_primary"`
	PrimaryDatabase              string         `db:"primary"`
	ReplicationAllowedToAccounts sql.NullString `db:"replication_allowed_to_accounts"`
	FailoverAllowedToAccounts    sql.NullString `db:"failover_allowed_to_accounts"`
	OrganizationName             string         `db:"organization_name"`
	AccountLocator               string         `db:"account_locator"`
}

type ReplicationDatabase struct {
	RegionGroup                  string
	SnowflakeRegion              string
	CreatedOn                    string
	AccountName                  string
	Name                         string
	Comment                      string
	IsPrimary                    bool
	PrimaryDatabase              *ExternalObjectIdentifier
	ReplicationAllowedToAccounts string
	FailoverAllowedToAccounts    string
	OrganizationName             string
	AccountLocator               string
}

func (row replicationDatabaseRow) convert() *ReplicationDatabase {
	db := &ReplicationDatabase{
		SnowflakeRegion:  row.SnowflakeRegion,
		CreatedOn:        row.CreatedOn,
		AccountName:      row.AccountName,
		Name:             row.Name,
		IsPrimary:        row.IsPrimary,
		OrganizationName: row.OrganizationName,
		AccountLocator:   row.AccountLocator,
	}
	if row.PrimaryDatabase != "" {
		primaryDatabaseId, err := ParseExternalObjectIdentifier(row.PrimaryDatabase)
		if err != nil {
			log.Printf("[DEBUG] unable to parse primary database identifier: %v, err = %s", row.PrimaryDatabase, err)
		} else {
			db.PrimaryDatabase = &primaryDatabaseId
		}
	}
	if row.RegionGroup.Valid {
		db.RegionGroup = row.RegionGroup.String
	}
	if row.Comment.Valid {
		db.Comment = row.Comment.String
	}
	if row.ReplicationAllowedToAccounts.Valid {
		db.ReplicationAllowedToAccounts = row.ReplicationAllowedToAccounts.String
	}
	if row.FailoverAllowedToAccounts.Valid {
		db.FailoverAllowedToAccounts = row.FailoverAllowedToAccounts.String
	}
	return db
}

// ShowReplicationDatabasesOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-replication-databases.
type ShowReplicationDatabasesOptions struct {
	show                 bool                      `ddl:"static" sql:"SHOW"`
	replicationDatabases bool                      `ddl:"static" sql:"REPLICATION DATABASES"`
	Like                 *Like                     `ddl:"keyword" sql:"LIKE"`
	WithPrimary          *ExternalObjectIdentifier `ddl:"identifier" sql:"WITH PRIMARY"`
}

func (opts *ShowReplicationDatabasesOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if opts.WithPrimary != nil && !ValidObjectIdentifier(opts.WithPrimary) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (c *replicationFunctions) ShowReplicationDatabases(ctx context.Context, opts *ShowReplicationDatabasesOptions) ([]ReplicationDatabase, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[replicationDatabaseRow](c.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRowsErr[replicationDatabaseRow, ReplicationDatabase](dbRows)
}

type CloudType string

const (
	CloudTypeAWS   CloudType = "aws"
	CloudTypeAzure CloudType = "azure"
	CloudTypeGCP   CloudType = "gcp"
)

type Region struct {
	RegionGroup     string
	SnowflakeRegion string
	CloudType       CloudType
	Region          string
	DisplayName     string
}

type regionRow struct {
	RegionGroup     string `db:"region_group"`
	SnowflakeRegion string `db:"snowflake_region"`
	Cloud           string `db:"cloud"`
	Region          string `db:"region"`
	DisplayName     string `db:"display_name"`
}

func (row *regionRow) toRegion() *Region {
	return &Region{
		RegionGroup:     row.RegionGroup,
		SnowflakeRegion: row.SnowflakeRegion,
		CloudType:       CloudType(row.Cloud),
		Region:          row.Region,
		DisplayName:     row.DisplayName,
	}
}

type ShowRegionsOptions struct {
	show    bool  `ddl:"static" sql:"SHOW"`
	regions bool  `ddl:"static" sql:"REGIONS"`
	Like    *Like `ddl:"keyword" sql:"LIKE"`
}

func (opts *ShowRegionsOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

// ShowRegions is based on https://docs.snowflake.com/en/sql-reference/sql/show-regions.
func (c *replicationFunctions) ShowRegions(ctx context.Context, opts *ShowRegionsOptions) ([]*Region, error) {
	if opts == nil {
		opts = &ShowRegionsOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	rows := []regionRow{}
	err = c.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	regions := make([]*Region, len(rows))
	for i, row := range rows {
		regions[i] = row.toRegion()
	}
	return regions, nil
}

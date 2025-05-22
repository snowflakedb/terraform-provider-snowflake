package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"

	"github.com/jmoiron/sqlx"
	"github.com/snowflakedb/gosnowflake"
)

type Client struct {
	config         *gosnowflake.Config
	db             *sqlx.DB
	sessionID      string
	accountLocator string
	dryRun         bool
	traceLogs      []string

	// System-Defined Functions
	ContextFunctions     ContextFunctions
	ConversionFunctions  ConversionFunctions
	SystemFunctions      SystemFunctions
	ReplicationFunctions ReplicationFunctions

	// DDL Commands
	Accounts                     Accounts
	Alerts                       Alerts
	ApiIntegrations              ApiIntegrations
	ApplicationPackages          ApplicationPackages
	ApplicationRoles             ApplicationRoles
	Applications                 Applications
	AuthenticationPolicies       AuthenticationPolicies
	Comments                     Comments
	ComputePools                 ComputePools
	Connections                  Connections
	CortexSearchServices         CortexSearchServices
	DatabaseRoles                DatabaseRoles
	Databases                    Databases
	DataMetricFunctionReferences DataMetricFunctionReferences
	DynamicTables                DynamicTables
	ExternalFunctions            ExternalFunctions
	ExternalVolumes              ExternalVolumes
	ExternalTables               ExternalTables
	EventTables                  EventTables
	FailoverGroups               FailoverGroups
	FileFormats                  FileFormats
	Functions                    Functions
	Grants                       Grants
	ImageRepositories            ImageRepositories
	ManagedAccounts              ManagedAccounts
	MaskingPolicies              MaskingPolicies
	MaterializedViews            MaterializedViews
	NetworkPolicies              NetworkPolicies
	NetworkRules                 NetworkRules
	NotificationIntegrations     NotificationIntegrations
	Parameters                   Parameters
	PasswordPolicies             PasswordPolicies
	Pipes                        Pipes
	PolicyReferences             PolicyReferences
	Procedures                   Procedures
	ResourceMonitors             ResourceMonitors
	Roles                        Roles
	RowAccessPolicies            RowAccessPolicies
	Schemas                      Schemas
	Secrets                      Secrets
	SecurityIntegrations         SecurityIntegrations
	Services                     Services
	Sequences                    Sequences
	SessionPolicies              SessionPolicies
	Sessions                     Sessions
	Shares                       Shares
	Stages                       Stages
	StorageIntegrations          StorageIntegrations
	Streamlits                   Streamlits
	Streams                      Streams
	Tables                       Tables
	Tags                         Tags
	Tasks                        Tasks
	Users                        Users
	Views                        Views
	Warehouses                   Warehouses
}

func (c *Client) GetAccountLocator() string {
	return c.accountLocator
}

func (c *Client) GetConfig() *gosnowflake.Config {
	return c.config
}

func (c *Client) GetConn() *sqlx.DB {
	return c.db
}

func NewDefaultClient(opts ...func(*FileReaderConfig)) (*Client, error) {
	return NewClient(nil, opts...)
}

func NewDryRunClient() *Client {
	client := &Client{
		dryRun:    true,
		traceLogs: []string{},
	}
	client.initialize()
	return client
}

func NewClient(cfg *gosnowflake.Config, opts ...func(*FileReaderConfig)) (*Client, error) {
	var err error
	if cfg == nil {
		log.Printf("[DEBUG] Searching for default config in credentials chain...")
		cfg = DefaultConfig(opts...)
	}

	dsn, err := gosnowflake.DSN(cfg)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("snowflake", dsn)
	if err != nil {
		return nil, fmt.Errorf("open snowflake connection: %w", err)
	}

	client := &Client{
		// snowflake does not adhere to the normal sql driver interface, so we have to use unsafe
		db:     db.Unsafe(),
		config: cfg,
	}
	client.initialize()

	err = client.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping snowflake: %w", err)
	}
	ctx := context.Background()
	currentAccount, err := client.ContextFunctions.CurrentAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current account: %w", err)
	}
	client.accountLocator = currentAccount

	sessionID, err := client.ContextFunctions.CurrentSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current session: %w", err)
	}
	client.sessionID = sessionID
	log.Printf("[DEBUG] connection success! Account: %s", currentAccount)

	return client, nil
}

func (c *Client) initialize() {
	c.Accounts = &accounts{client: c}
	c.Alerts = &alerts{client: c}
	c.ApiIntegrations = &apiIntegrations{client: c}
	c.ApplicationPackages = &applicationPackages{client: c}
	c.ApplicationRoles = &applicationRoles{client: c}
	c.Applications = &applications{client: c}
	c.AuthenticationPolicies = &authenticationPolicies{client: c}
	c.Comments = &comments{client: c}
	c.ComputePools = &computePools{client: c}
	c.Connections = &connections{client: c}
	c.ContextFunctions = &contextFunctions{client: c}
	c.ConversionFunctions = &conversionFunctions{client: c}
	c.CortexSearchServices = &cortexSearchServices{client: c}
	c.DatabaseRoles = &databaseRoles{client: c}
	c.Databases = &databases{client: c}
	c.DataMetricFunctionReferences = &dataMetricFunctionReferences{client: c}
	c.DynamicTables = &dynamicTables{client: c}
	c.ExternalFunctions = &externalFunctions{client: c}
	c.ExternalVolumes = &externalVolumes{client: c}
	c.ExternalTables = &externalTables{client: c}
	c.EventTables = &eventTables{client: c}
	c.FailoverGroups = &failoverGroups{client: c}
	c.FileFormats = &fileFormats{client: c}
	c.Functions = &functions{client: c}
	c.Grants = &grants{client: c}
	c.ImageRepositories = &imageRepositories{client: c}
	c.ManagedAccounts = &managedAccounts{client: c}
	c.MaskingPolicies = &maskingPolicies{client: c}
	c.MaterializedViews = &materializedViews{client: c}
	c.NetworkPolicies = &networkPolicies{client: c}
	c.NetworkRules = &networkRules{client: c}
	c.NotificationIntegrations = &notificationIntegrations{client: c}
	c.Parameters = &parameters{client: c}
	c.PasswordPolicies = &passwordPolicies{client: c}
	c.Pipes = &pipes{client: c}
	c.PolicyReferences = &policyReference{client: c}
	c.Procedures = &procedures{client: c}
	c.ReplicationFunctions = &replicationFunctions{client: c}
	c.ResourceMonitors = &resourceMonitors{client: c}
	c.Roles = &roles{client: c}
	c.RowAccessPolicies = &rowAccessPolicies{client: c}
	c.Schemas = &schemas{client: c}
	c.Secrets = &secrets{client: c}
	c.SecurityIntegrations = &securityIntegrations{client: c}
	c.Sequences = &sequences{client: c}
	c.Services = &services{client: c}
	c.SessionPolicies = &sessionPolicies{client: c}
	c.Sessions = &sessions{client: c}
	c.Shares = &shares{client: c}
	c.Stages = &stages{client: c}
	c.StorageIntegrations = &storageIntegrations{client: c}
	c.Streamlits = &streamlits{client: c}
	c.Streams = &streams{client: c}
	c.SystemFunctions = &systemFunctions{client: c}
	c.Tables = &tables{client: c}
	c.Tags = &tags{client: c}
	c.Tasks = &tasks{client: c}
	c.Users = &users{client: c}
	c.Views = &views{client: c}
	c.Warehouses = &warehouses{client: c}
}

func (c *Client) TraceLogs() []string {
	return c.traceLogs
}

func (c *Client) Ping() error {
	return c.db.Ping()
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

type accountLocatorContextKey struct{}

var snowflakeAccountLocatorContextKey accountLocatorContextKey

// Exec executes a query that does not return rows.
func (c *Client) exec(ctx context.Context, sql string) (sql.Result, error) {
	if c.dryRun {
		c.traceLogs = append(c.traceLogs, sql)
		// TODO(SNOW-926146): Decide what to do with logs during plugin framework poc
		// log.Printf("[DEBUG] sql-conn-exec-dry: %v", sql)
		return nil, nil
	}
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	sql = appendQueryMetadata(ctx, sql)
	result, err := c.db.ExecContext(ctx, sql)
	return result, decodeDriverError(err)
}

// query runs a query and returns the rows. dest is expected to be a slice of structs.
func (c *Client) query(ctx context.Context, dest interface{}, sql string) error {
	if c.dryRun {
		c.traceLogs = append(c.traceLogs, sql)
		// TODO(SNOW-926146): Decide what to do with logs during plugin framework poc
		// log.Printf("[DEBUG] sql-conn-query-dry: %v", sql)
		return nil
	}
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	sql = appendQueryMetadata(ctx, sql)
	return decodeDriverError(c.db.SelectContext(ctx, dest, sql))
}

// queryOne runs a query and returns one row. dest is expected to be a pointer to a struct.
func (c *Client) queryOne(ctx context.Context, dest interface{}, sql string) error {
	if c.dryRun {
		c.traceLogs = append(c.traceLogs, sql)
		// TODO(SNOW-926146): Decide what to do with logs during plugin framework poc
		// log.Printf("[DEBUG] sql-conn-query-one-dry: %v", sql)
		return nil
	}
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	sql = appendQueryMetadata(ctx, sql)
	return decodeDriverError(c.db.GetContext(ctx, dest, sql))
}

func appendQueryMetadata(ctx context.Context, sql string) string {
	if metadata, ok := tracking.FromContext(ctx); ok {
		newSql, err := tracking.AppendMetadata(sql, metadata)
		if err != nil {
			log.Printf("[ERROR] failed to append metadata tracking: %v", err)
			return sql
		}
		return newSql
	}
	return sql
}

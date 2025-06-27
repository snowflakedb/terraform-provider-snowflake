package sdk

import (
	"github.com/jmoiron/sqlx"
	"github.com/snowflakedb/gosnowflake"
)

// TODO: Extend dto builder to support generating a new builder type from a struct (e.g. below code could be entirely generated out of Client struct).
//go:generate go run ./dto-builder-generator/main.go

type ClientBuilder struct {
	Config         *gosnowflake.Config
	Db             *sqlx.DB
	SessionID      string
	AccountLocator string
	DryRun         bool
	TraceLogs      []string

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
	GitRepositories              GitRepositories
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

func (s *ClientBuilder) Build() *Client {
	return &Client{
		config:         s.Config,
		db:             s.Db,
		sessionID:      s.SessionID,
		accountLocator: s.AccountLocator,
		dryRun:         s.DryRun,
		traceLogs:      s.TraceLogs,

		// System-Defined Functions
		ContextFunctions:     s.ContextFunctions,
		ConversionFunctions:  s.ConversionFunctions,
		SystemFunctions:      s.SystemFunctions,
		ReplicationFunctions: s.ReplicationFunctions,

		// DDL Commands
		Accounts:                     s.Accounts,
		Alerts:                       s.Alerts,
		ApiIntegrations:              s.ApiIntegrations,
		ApplicationPackages:          s.ApplicationPackages,
		ApplicationRoles:             s.ApplicationRoles,
		Applications:                 s.Applications,
		AuthenticationPolicies:       s.AuthenticationPolicies,
		Comments:                     s.Comments,
		ComputePools:                 s.ComputePools,
		Connections:                  s.Connections,
		CortexSearchServices:         s.CortexSearchServices,
		DatabaseRoles:                s.DatabaseRoles,
		Databases:                    s.Databases,
		DataMetricFunctionReferences: s.DataMetricFunctionReferences,
		DynamicTables:                s.DynamicTables,
		ExternalFunctions:            s.ExternalFunctions,
		ExternalVolumes:              s.ExternalVolumes,
		ExternalTables:               s.ExternalTables,
		EventTables:                  s.EventTables,
		FailoverGroups:               s.FailoverGroups,
		FileFormats:                  s.FileFormats,
		Functions:                    s.Functions,
		GitRepositories:              s.GitRepositories,
		Grants:                       s.Grants,
		ImageRepositories:            s.ImageRepositories,
		ManagedAccounts:              s.ManagedAccounts,
		MaskingPolicies:              s.MaskingPolicies,
		MaterializedViews:            s.MaterializedViews,
		NetworkPolicies:              s.NetworkPolicies,
		NetworkRules:                 s.NetworkRules,
		NotificationIntegrations:     s.NotificationIntegrations,
		Parameters:                   s.Parameters,
		PasswordPolicies:             s.PasswordPolicies,
		Pipes:                        s.Pipes,
		PolicyReferences:             s.PolicyReferences,
		Procedures:                   s.Procedures,
		ResourceMonitors:             s.ResourceMonitors,
		Roles:                        s.Roles,
		RowAccessPolicies:            s.RowAccessPolicies,
		Schemas:                      s.Schemas,
		Secrets:                      s.Secrets,
		SecurityIntegrations:         s.SecurityIntegrations,
		Services:                     s.Services,
		Sequences:                    s.Sequences,
		SessionPolicies:              s.SessionPolicies,
		Sessions:                     s.Sessions,
		Shares:                       s.Shares,
		Stages:                       s.Stages,
		StorageIntegrations:          s.StorageIntegrations,
		Streamlits:                   s.Streamlits,
		Streams:                      s.Streams,
		Tables:                       s.Tables,
		Tags:                         s.Tags,
		Tasks:                        s.Tasks,
		Users:                        s.Users,
		Views:                        s.Views,
		Warehouses:                   s.Warehouses,
	}
}

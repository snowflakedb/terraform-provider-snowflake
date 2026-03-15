package sdk

import (
	"context"
	"database/sql"
	"time"
)

type ExternalAccessIntegrations interface {
	Create(ctx context.Context, request *CreateExternalAccessIntegrationRequest) error
	Alter(ctx context.Context, request *AlterExternalAccessIntegrationRequest) error
	Drop(ctx context.Context, request *DropExternalAccessIntegrationRequest) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, request *ShowExternalAccessIntegrationRequest) ([]ExternalAccessIntegration, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ExternalAccessIntegration, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*ExternalAccessIntegration, error)
	Describe(ctx context.Context, id AccountObjectIdentifier) ([]ExternalAccessIntegrationProperty, error)
}

// CreateExternalAccessIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-external-access-integration.
type CreateExternalAccessIntegrationOptions struct {
	create                       bool                     `ddl:"static" sql:"CREATE"`
	OrReplace                    *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	externalAccessIntegration    bool                     `ddl:"static" sql:"EXTERNAL ACCESS INTEGRATION"`
	IfNotExists                  *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                         AccountObjectIdentifier  `ddl:"identifier"`
	AllowedNetworkRules          []SchemaObjectIdentifier `ddl:"parameter,parentheses" sql:"ALLOWED_NETWORK_RULES"`
	AllowedAuthenticationSecrets []SchemaObjectIdentifier `ddl:"parameter,parentheses" sql:"ALLOWED_AUTHENTICATION_SECRETS"`
	Enabled                      bool                     `ddl:"parameter" sql:"ENABLED"`
	Comment                      *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterExternalAccessIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-external-access-integration.
type AlterExternalAccessIntegrationOptions struct {
	alter                     bool                            `ddl:"static" sql:"ALTER"`
	externalAccessIntegration bool                            `ddl:"static" sql:"EXTERNAL ACCESS INTEGRATION"`
	IfExists                  *bool                           `ddl:"keyword" sql:"IF EXISTS"`
	name                      AccountObjectIdentifier         `ddl:"identifier"`
	Set                       *ExternalAccessIntegrationSet   `ddl:"keyword" sql:"SET"`
	Unset                     *ExternalAccessIntegrationUnset `ddl:"list,no_parentheses" sql:"UNSET"`
}

type ExternalAccessIntegrationSet struct {
	AllowedNetworkRules          []SchemaObjectIdentifier `ddl:"parameter,parentheses" sql:"ALLOWED_NETWORK_RULES"`
	AllowedAuthenticationSecrets []SchemaObjectIdentifier `ddl:"parameter,parentheses" sql:"ALLOWED_AUTHENTICATION_SECRETS"`
	Enabled                      *bool                    `ddl:"parameter" sql:"ENABLED"`
	Comment                      *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ExternalAccessIntegrationUnset struct {
	AllowedAuthenticationSecrets *bool `ddl:"keyword" sql:"ALLOWED_AUTHENTICATION_SECRETS"`
	Comment                      *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropExternalAccessIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-integration.
type DropExternalAccessIntegrationOptions struct {
	drop                      bool                    `ddl:"static" sql:"DROP"`
	externalAccessIntegration bool                    `ddl:"static" sql:"EXTERNAL ACCESS INTEGRATION"`
	IfExists                  *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name                      AccountObjectIdentifier `ddl:"identifier"`
}

// ShowExternalAccessIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-integrations.
type ShowExternalAccessIntegrationOptions struct {
	show                       bool  `ddl:"static" sql:"SHOW"`
	externalAccessIntegrations bool  `ddl:"static" sql:"EXTERNAL ACCESS INTEGRATIONS"`
	Like                       *Like `ddl:"keyword" sql:"LIKE"`
}

type showExternalAccessIntegrationsDbRow struct {
	Name      string         `db:"name"`
	Type      string         `db:"type"`
	Category  string         `db:"category"`
	Enabled   bool           `db:"enabled"`
	Comment   sql.NullString `db:"comment"`
	CreatedOn time.Time      `db:"created_on"`
}

type ExternalAccessIntegration struct {
	Name      string
	IntType   string
	Category  string
	Enabled   bool
	Comment   string
	CreatedOn time.Time
}

func (v *ExternalAccessIntegration) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *ExternalAccessIntegration) ObjectType() ObjectType {
	return ObjectTypeIntegration
}

// DescribeExternalAccessIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-integration.
type DescribeExternalAccessIntegrationOptions struct {
	describe                  bool                    `ddl:"static" sql:"DESCRIBE"`
	externalAccessIntegration bool                    `ddl:"static" sql:"EXTERNAL ACCESS INTEGRATION"`
	name                      AccountObjectIdentifier `ddl:"identifier"`
}

type descExternalAccessIntegrationsDbRow struct {
	Property        string `db:"property"`
	PropertyType    string `db:"property_type"`
	PropertyValue   string `db:"property_value"`
	PropertyDefault string `db:"property_default"`
}

type ExternalAccessIntegrationProperty struct {
	Name    string
	Type    string
	Value   string
	Default string
}

package sdk

import (
	"context"
	"database/sql"
	"time"
)

type ApiIntegrations interface {
	Create(ctx context.Context, request *CreateApiIntegrationRequest) error
	Alter(ctx context.Context, request *AlterApiIntegrationRequest) error
	Drop(ctx context.Context, request *DropApiIntegrationRequest) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, request *ShowApiIntegrationRequest) ([]ApiIntegration, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ApiIntegration, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*ApiIntegration, error)
	Describe(ctx context.Context, id AccountObjectIdentifier) ([]ApiIntegrationProperty, error)
}

// CreateApiIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-api-integration.
type CreateApiIntegrationOptions struct {
	create                  bool                           `ddl:"static" sql:"CREATE"`
	OrReplace               *bool                          `ddl:"keyword" sql:"OR REPLACE"`
	apiIntegration          bool                           `ddl:"static" sql:"API INTEGRATION"`
	IfNotExists             *bool                          `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                    AccountObjectIdentifier        `ddl:"identifier"`
	AwsApiProviderParams    *AwsApiParams                  `ddl:"keyword"`
	AzureApiProviderParams  *AzureApiParams                `ddl:"keyword"`
	GoogleApiProviderParams *GoogleApiParams               `ddl:"keyword"`
	ApiAllowedPrefixes      []ApiIntegrationEndpointPrefix `ddl:"parameter,parentheses" sql:"API_ALLOWED_PREFIXES"`
	ApiBlockedPrefixes      []ApiIntegrationEndpointPrefix `ddl:"parameter,parentheses" sql:"API_BLOCKED_PREFIXES"`
	Enabled                 bool                           `ddl:"parameter" sql:"ENABLED"`
	Comment                 *string                        `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ApiIntegrationEndpointPrefix struct {
	Path string `ddl:"keyword,single_quotes"`
}

type AwsApiParams struct {
	ApiProvider   ApiIntegrationAwsApiProviderType `ddl:"parameter,no_quotes" sql:"API_PROVIDER"`
	ApiAwsRoleArn string                           `ddl:"parameter,single_quotes" sql:"API_AWS_ROLE_ARN"`
	ApiKey        *string                          `ddl:"parameter,single_quotes" sql:"API_KEY"`
}

type AzureApiParams struct {
	apiProvider          string  `ddl:"static" sql:"API_PROVIDER = azure_api_management"`
	AzureTenantId        string  `ddl:"parameter,single_quotes" sql:"AZURE_TENANT_ID"`
	AzureAdApplicationId string  `ddl:"parameter,single_quotes" sql:"AZURE_AD_APPLICATION_ID"`
	ApiKey               *string `ddl:"parameter,single_quotes" sql:"API_KEY"`
}

type GoogleApiParams struct {
	apiProvider    string `ddl:"static" sql:"API_PROVIDER = google_api_gateway"`
	GoogleAudience string `ddl:"parameter,single_quotes" sql:"GOOGLE_AUDIENCE"`
}

// AlterApiIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-api-integration.
type AlterApiIntegrationOptions struct {
	alter          bool                    `ddl:"static" sql:"ALTER"`
	apiIntegration bool                    `ddl:"static" sql:"API INTEGRATION"`
	IfExists       *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name           AccountObjectIdentifier `ddl:"identifier"`
	Set            *ApiIntegrationSet      `ddl:"keyword" sql:"SET"`
	Unset          *ApiIntegrationUnset    `ddl:"list,no_parentheses" sql:"UNSET"`
	SetTags        []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTags      []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
}

type ApiIntegrationSet struct {
	AwsParams          *SetAwsApiParams               `ddl:"keyword"`
	AzureParams        *SetAzureApiParams             `ddl:"keyword"`
	GoogleParams       *SetGoogleApiParams            `ddl:"keyword"`
	Enabled            *bool                          `ddl:"parameter" sql:"ENABLED"`
	ApiAllowedPrefixes []ApiIntegrationEndpointPrefix `ddl:"parameter,parentheses" sql:"API_ALLOWED_PREFIXES"`
	ApiBlockedPrefixes []ApiIntegrationEndpointPrefix `ddl:"parameter,parentheses" sql:"API_BLOCKED_PREFIXES"`
	Comment            *string                        `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type SetAwsApiParams struct {
	ApiAwsRoleArn *string `ddl:"parameter,single_quotes" sql:"API_AWS_ROLE_ARN"`
	ApiKey        *string `ddl:"parameter,single_quotes" sql:"API_KEY"`
}

type SetAzureApiParams struct {
	AzureTenantId        *string `ddl:"parameter,single_quotes" sql:"AZURE_TENANT_ID"`
	AzureAdApplicationId *string `ddl:"parameter,single_quotes" sql:"AZURE_AD_APPLICATION_ID"`
	ApiKey               *string `ddl:"parameter,single_quotes" sql:"API_KEY"`
}

type SetGoogleApiParams struct {
	GoogleAudience string `ddl:"parameter,single_quotes" sql:"GOOGLE_AUDIENCE"`
}

type ApiIntegrationUnset struct {
	ApiKey             *bool `ddl:"keyword" sql:"API_KEY"`
	Enabled            *bool `ddl:"keyword" sql:"ENABLED"`
	ApiBlockedPrefixes *bool `ddl:"keyword" sql:"API_BLOCKED_PREFIXES"`
	Comment            *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropApiIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-integration.
type DropApiIntegrationOptions struct {
	drop           bool                    `ddl:"static" sql:"DROP"`
	apiIntegration bool                    `ddl:"static" sql:"API INTEGRATION"`
	IfExists       *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name           AccountObjectIdentifier `ddl:"identifier"`
}

// ShowApiIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-integrations.
type ShowApiIntegrationOptions struct {
	show            bool  `ddl:"static" sql:"SHOW"`
	apiIntegrations bool  `ddl:"static" sql:"API INTEGRATIONS"`
	Like            *Like `ddl:"keyword" sql:"LIKE"`
}

type showApiIntegrationsDbRow struct {
	Name      string         `db:"name"`
	Type      string         `db:"type"`
	Category  string         `db:"category"`
	Enabled   bool           `db:"enabled"`
	Comment   sql.NullString `db:"comment"`
	CreatedOn time.Time      `db:"created_on"`
}

type ApiIntegration struct {
	Name      string
	ApiType   string
	Category  string
	Enabled   bool
	Comment   string
	CreatedOn time.Time
}

func (v *ApiIntegration) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *ApiIntegration) ObjectType() ObjectType {
	return ObjectTypeIntegration
}

// DescribeApiIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-integration.
type DescribeApiIntegrationOptions struct {
	describe       bool                    `ddl:"static" sql:"DESCRIBE"`
	apiIntegration bool                    `ddl:"static" sql:"API INTEGRATION"`
	name           AccountObjectIdentifier `ddl:"identifier"`
}

type descApiIntegrationsDbRow struct {
	Property        string `db:"property"`
	PropertyType    string `db:"property_type"`
	PropertyValue   string `db:"property_value"`
	PropertyDefault string `db:"property_default"`
}

type ApiIntegrationProperty struct {
	Name    string
	Type    string
	Value   string
	Default string
}

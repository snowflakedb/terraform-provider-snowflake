package sdk

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/snowflakedb/gosnowflake"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ validatable = new(CreateAccountOptions)
	_ validatable = new(AlterAccountOptions)
	_ validatable = new(ShowAccountOptions)
)

type Accounts interface {
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateAccountOptions) (*AccountCreateResponse, error)
	Alter(ctx context.Context, opts *AlterAccountOptions) error
	Show(ctx context.Context, opts *ShowAccountOptions) ([]Account, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Account, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Account, error)
	Drop(ctx context.Context, id AccountObjectIdentifier, gracePeriodInDays int, opts *DropAccountOptions) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier, gracePeriodInDays int) error
	Undrop(ctx context.Context, id AccountObjectIdentifier) error
	ShowParameters(ctx context.Context) ([]*Parameter, error)
	UnsetAllParameters(ctx context.Context) error
	UnsetAllPoliciesSafely(ctx context.Context) error
	UnsetPolicySafely(ctx context.Context, kind PolicyKind) error
	UnsetAllAttachments(ctx context.Context) error
}

var _ Accounts = (*accounts)(nil)

type accounts struct {
	client *Client
}

type AccountEdition string

var (
	EditionStandard         AccountEdition = "STANDARD"
	EditionEnterprise       AccountEdition = "ENTERPRISE"
	EditionBusinessCritical AccountEdition = "BUSINESS_CRITICAL"
)

var AllAccountEditions = []AccountEdition{
	EditionStandard,
	EditionEnterprise,
	EditionBusinessCritical,
}

func ToAccountEdition(edition string) (AccountEdition, error) {
	switch typedEdition := AccountEdition(strings.ToUpper(edition)); typedEdition {
	case EditionStandard, EditionEnterprise, EditionBusinessCritical:
		return typedEdition, nil
	default:
		return "", fmt.Errorf("unknown account edition: %s", edition)
	}
}

// CreateAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-account.
type CreateAccountOptions struct {
	create  bool                    `ddl:"static" sql:"CREATE"`
	account bool                    `ddl:"static" sql:"ACCOUNT"`
	name    AccountObjectIdentifier `ddl:"identifier"`

	// Object properties
	AdminName                string         `ddl:"parameter,single_quotes" sql:"ADMIN_NAME"`
	AdminPassword            *string        `ddl:"parameter,single_quotes" sql:"ADMIN_PASSWORD"`
	AdminRSAPublicKey        *string        `ddl:"parameter,single_quotes" sql:"ADMIN_RSA_PUBLIC_KEY"`
	AdminUserType            *UserType      `ddl:"parameter" sql:"ADMIN_USER_TYPE"`
	FirstName                *string        `ddl:"parameter,single_quotes" sql:"FIRST_NAME"`
	LastName                 *string        `ddl:"parameter,single_quotes" sql:"LAST_NAME"`
	Email                    string         `ddl:"parameter,single_quotes" sql:"EMAIL"`
	MustChangePassword       *bool          `ddl:"parameter" sql:"MUST_CHANGE_PASSWORD"`
	Edition                  AccountEdition `ddl:"parameter" sql:"EDITION"`
	RegionGroup              *string        `ddl:"parameter" sql:"REGION_GROUP"`
	Region                   *string        `ddl:"parameter" sql:"REGION"`
	Comment                  *string        `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ConsumptionBillingEntity *string        `ddl:"parameter,double_quotes" sql:"CONSUMPTION_BILLING_ENTITY"`
	Polaris                  *bool          `ddl:"parameter" sql:"POLARIS"`
}

func (opts *CreateAccountOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if opts.AdminName == "" {
		errs = append(errs, errNotSet("CreateAccountOptions", "AdminName"))
	}
	if !anyValueSet(opts.AdminPassword, opts.AdminRSAPublicKey) {
		errs = append(errs, errAtLeastOneOf("CreateAccountOptions", "AdminPassword", "AdminRSAPublicKey"))
	}
	if opts.Email == "" {
		errs = append(errs, errNotSet("CreateAccountOptions", "Email"))
	}
	if opts.Edition == "" {
		errs = append(errs, errNotSet("CreateAccountOptions", "Edition"))
	}
	return errors.Join(errs...)
}

type AccountCreateResponse struct {
	AccountLocator    string `json:"accountLocator,omitempty"`
	AccountLocatorUrl string `json:"accountLocatorUrl,omitempty"`
	OrganizationName  string
	AccountName       string         `json:"accountName,omitempty"`
	Url               string         `json:"url,omitempty"`
	Edition           AccountEdition `json:"edition,omitempty"`
	RegionGroup       string         `json:"regionGroup,omitempty"`
	Cloud             string         `json:"cloud,omitempty"`
	Region            string         `json:"region,omitempty"`
}

func ToAccountCreateResponse(v string) (*AccountCreateResponse, error) {
	var res AccountCreateResponse
	err := json.Unmarshal([]byte(v), &res)
	if err != nil {
		return nil, err
	}
	if len(res.Url) > 0 {
		url := strings.TrimPrefix(res.Url, `https://`)
		url = strings.TrimPrefix(url, `http://`)
		parts := strings.SplitN(url, "-", 2)
		if len(parts) == 2 {
			res.OrganizationName = strings.ToUpper(parts[0])
		}
	}
	return &res, nil
}

func (c *accounts) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateAccountOptions) (*AccountCreateResponse, error) {
	if opts == nil {
		opts = &CreateAccountOptions{}
	}
	opts.name = id
	queryChanId := make(chan string, 1)
	err := validateAndExec(c.client, gosnowflake.WithQueryIDChan(ctx, queryChanId), opts)
	if err != nil {
		return nil, err
	}

	queryId := <-queryChanId
	rows, err := c.client.QueryUnsafe(gosnowflake.WithFetchResultByID(ctx, queryId), "")
	if err != nil {
		log.Printf("[WARN] Unable to retrieve create account output, err = %v", err)
	}

	if len(rows) == 1 && rows[0]["status"] != nil {
		if status, ok := (*rows[0]["status"]).(string); ok {
			return ToAccountCreateResponse(status)
		}
	}

	return nil, nil
}

// AlterAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-account.
type AlterAccountOptions struct {
	alter   bool                     `ddl:"static" sql:"ALTER"`
	account bool                     `ddl:"static" sql:"ACCOUNT"`
	Name    *AccountObjectIdentifier `ddl:"identifier"`

	Set      *AccountSet        `ddl:"keyword" sql:"SET"`
	Unset    *AccountUnset      `ddl:"list,no_parentheses" sql:"UNSET"`
	SetTag   []TagAssociation   `ddl:"keyword" sql:"SET TAG"`
	UnsetTag []ObjectIdentifier `ddl:"keyword" sql:"UNSET TAG"`
	Rename   *AccountRename     `ddl:"-"`
	Drop     *AccountDrop       `ddl:"-"`
}

func (opts *AlterAccountOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTag, opts.UnsetTag, opts.Drop, opts.Rename) {
		errs = append(errs, errExactlyOneOf("AlterAccountOptions", "Set", "Unset", "SetTag", "UnsetTag", "Drop", "Rename"))
	}
	if valueSet(opts.Set) {
		if valueSet(opts.Set.ConsumptionBillingEntity) {
			if !valueSet(opts.Name) || !ValidObjectIdentifier(opts.Name) {
				errs = append(errs, ErrInvalidObjectIdentifier)
			}
		}
		if err := opts.Set.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Unset) {
		if valueSet(opts.Unset.ConsumptionBillingEntity) {
			if !valueSet(opts.Name) || !ValidObjectIdentifier(opts.Name) {
				errs = append(errs, ErrInvalidObjectIdentifier)
			}
		}
		if err := opts.Unset.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Drop) {
		if err := opts.Drop.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Rename) {
		if err := opts.Rename.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Drop) || valueSet(opts.Rename) {
		if !valueSet(opts.Name) || !ValidObjectIdentifier(opts.Name) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	return errors.Join(errs...)
}

type AccountLevelParameters struct {
	AccountParameters *LegacyAccountParameters `ddl:"list,no_parentheses"`
	SessionParameters *SessionParameters       `ddl:"list,no_parentheses"`
	ObjectParameters  *ObjectParameters        `ddl:"list,no_parentheses"`
	UserParameters    *UserParameters          `ddl:"list,no_parentheses"`
}

func (opts *AccountLevelParameters) validate() error {
	var errs []error
	if valueSet(opts.AccountParameters) {
		if err := opts.AccountParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.SessionParameters) {
		if err := opts.SessionParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.ObjectParameters) {
		if err := opts.ObjectParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.UserParameters) {
		if err := opts.UserParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type AccountSet struct {
	Parameters               *AccountParameters       `ddl:"list,no_parentheses"`
	LegacyParameters         *AccountLevelParameters  `ddl:"list,no_parentheses"`
	ResourceMonitor          *AccountObjectIdentifier `ddl:"identifier,equals" sql:"RESOURCE_MONITOR"`
	PackagesPolicy           *SchemaObjectIdentifier  `ddl:"identifier" sql:"PACKAGES POLICY"`
	PasswordPolicy           *SchemaObjectIdentifier  `ddl:"identifier" sql:"PASSWORD POLICY"`
	SessionPolicy            *SchemaObjectIdentifier  `ddl:"identifier" sql:"SESSION POLICY"`
	AuthenticationPolicy     *SchemaObjectIdentifier  `ddl:"identifier" sql:"AUTHENTICATION POLICY"`
	FeaturePolicySet         *AccountFeaturePolicySet `ddl:"keyword"`
	ConsumptionBillingEntity *string                  `ddl:"parameter,double_quotes" sql:"CONSUMPTION_BILLING_ENTITY"`
	OrgAdmin                 *bool                    `ddl:"parameter" sql:"IS_ORG_ADMIN"`
	Force                    *bool                    `ddl:"keyword" sql:"FORCE"`
}

type AccountFeaturePolicySet struct {
	FeaturePolicy      *SchemaObjectIdentifier `ddl:"identifier" sql:"FEATURE POLICY"`
	forAllApplications bool                    `ddl:"static" sql:"FOR ALL APPLICATIONS"`
}

func (opts *AccountSet) validate() error {
	var errs []error
	if !exactlyOneValueSet(opts.Parameters, opts.LegacyParameters, opts.ResourceMonitor, opts.PackagesPolicy, opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy, opts.FeaturePolicySet, opts.OrgAdmin, opts.ConsumptionBillingEntity) {
		errs = append(errs, errExactlyOneOf("AccountSet", "Parameters", "LegacyParameters", "ResourceMonitor", "PackagesPolicy", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "FeaturePolicySet", "OrgAdmin", "ConsumptionBillingEntity"))
	}
	if valueSet(opts.Force) && !valueSet(opts.PackagesPolicy) && !valueSet(opts.FeaturePolicySet) {
		errs = append(errs, NewError("force can only be set with PackagesPolicy and FeaturePolicy"))
	}
	if valueSet(opts.LegacyParameters) {
		if err := opts.LegacyParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type AccountLevelParametersUnset struct {
	AccountParameters *LegacyAccountParametersUnset `ddl:"list,no_parentheses"`
	SessionParameters *SessionParametersUnset       `ddl:"list,no_parentheses"`
	ObjectParameters  *ObjectParametersUnset        `ddl:"list,no_parentheses"`
	UserParameters    *UserParametersUnset          `ddl:"list,no_parentheses"`
}

func (opts *AccountLevelParametersUnset) validate() error {
	if !anyValueSet(opts.AccountParameters, opts.SessionParameters, opts.ObjectParameters, opts.UserParameters) {
		return errAtLeastOneOf("AccountLevelParametersUnset", "LegacyAccountParameters", "SessionParameters", "ObjectParameters", "UserParameters")
	}
	return nil
}

type AccountUnset struct {
	Parameters               *AccountParametersUnset      `ddl:"list,no_parentheses"`
	LegacyParameters         *AccountLevelParametersUnset `ddl:"list,no_parentheses"`
	AuthenticationPolicy     *bool                        `ddl:"keyword" sql:"AUTHENTICATION POLICY"`
	FeaturePolicyUnset       *AccountFeaturePolicyUnset   `ddl:"keyword"`
	PackagesPolicy           *bool                        `ddl:"keyword" sql:"PACKAGES POLICY"`
	PasswordPolicy           *bool                        `ddl:"keyword" sql:"PASSWORD POLICY"`
	SessionPolicy            *bool                        `ddl:"keyword" sql:"SESSION POLICY"`
	ResourceMonitor          *bool                        `ddl:"keyword" sql:"RESOURCE_MONITOR"`
	ConsumptionBillingEntity *bool                        `ddl:"keyword" sql:"CONSUMPTION_BILLING_ENTITY"`
}

type AccountFeaturePolicyUnset struct {
	FeaturePolicy      *bool `ddl:"keyword" sql:"FEATURE POLICY"`
	forAllApplications bool  `ddl:"static" sql:"FOR ALL APPLICATIONS"`
}

func (opts *AccountUnset) validate() error {
	var errs []error
	if !exactlyOneValueSet(opts.LegacyParameters, opts.Parameters, opts.PackagesPolicy, opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy, opts.ResourceMonitor, opts.FeaturePolicyUnset, opts.ConsumptionBillingEntity) {
		errs = append(errs, errExactlyOneOf("AccountUnset", "Parameters", "LegacyParameters", "PackagesPolicy", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "ResourceMonitor", "FeaturePolicyUnset", "ConsumptionBillingEntity"))
	}
	if valueSet(opts.LegacyParameters) {
		if err := opts.LegacyParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type AccountRename struct {
	NewName    AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SaveOldURL *bool                   `ddl:"parameter" sql:"SAVE_OLD_URL"`
}

func (opts *AccountRename) validate() error {
	var errs []error
	if !ValidObjectIdentifier(opts.NewName) {
		errs = append(errs, errInvalidIdentifier("AccountRename", "NewName"))
	}
	return errors.Join(errs...)
}

type AccountDrop struct {
	OldUrl             *bool `ddl:"keyword" sql:"DROP OLD URL"`
	OldOrganizationUrl *bool `ddl:"keyword" sql:"DROP OLD ORGANIZATION URL"`
}

func (opts *AccountDrop) validate() error {
	var errs []error
	if !exactlyOneValueSet(opts.OldUrl, opts.OldOrganizationUrl) {
		errs = append(errs, errExactlyOneOf("AccountDrop", "OldUrl", "OldOrganizationUrl"))
	}
	return errors.Join(errs...)
}

func (c *accounts) Alter(ctx context.Context, opts *AlterAccountOptions) error {
	if opts == nil {
		opts = &AlterAccountOptions{}
	}
	return validateAndExec(c.client, ctx, opts)
}

// ShowAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-organisation-accounts.
type ShowAccountOptions struct {
	show     bool  `ddl:"static" sql:"SHOW"`
	accounts bool  `ddl:"static" sql:"ACCOUNTS"`
	History  *bool `ddl:"keyword" sql:"HISTORY"`
	Like     *Like `ddl:"keyword" sql:"LIKE"`
}

func (opts *ShowAccountOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

type Account struct {
	OrganizationName                     string
	AccountName                          string
	SnowflakeRegion                      string
	RegionGroup                          *string // shows only for organizations that span multiple region groups
	Edition                              *AccountEdition
	AccountURL                           *string
	CreatedOn                            *time.Time
	Comment                              *string
	AccountLocator                       string
	AccountLocatorUrl                    *string
	ManagedAccounts                      *int
	ConsumptionBillingEntityName         *string
	MarketplaceConsumerBillingEntityName *string
	MarketplaceProviderBillingEntityName *string
	OldAccountURL                        *string
	IsOrgAdmin                           *bool
	AccountOldUrlSavedOn                 *time.Time
	AccountOldUrlLastUsed                *time.Time
	OrganizationOldUrl                   *string
	OrganizationOldUrlSavedOn            *time.Time
	OrganizationOldUrlLastUsed           *time.Time
	IsEventsAccount                      *bool
	IsOrganizationAccount                bool
	// Available only with the History keyword set
	DroppedOn                   *time.Time
	ScheduledDeletionTime       *time.Time
	RestoredOn                  *time.Time
	MovedToOrganization         *string
	MovedOn                     *string
	OrganizationUrlExpirationOn *time.Time
}

func (v *Account) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.AccountName)
}

func (v *Account) AccountID() AccountIdentifier {
	return NewAccountIdentifier(v.OrganizationName, v.AccountName)
}

type accountDBRow struct {
	OrganizationName                     string         `db:"organization_name"`
	AccountName                          string         `db:"account_name"`
	RegionGroup                          sql.NullString `db:"region_group"`
	SnowflakeRegion                      string         `db:"snowflake_region"`
	Edition                              sql.NullString `db:"edition"`
	AccountURL                           sql.NullString `db:"account_url"`
	CreatedOn                            sql.NullTime   `db:"created_on"`
	Comment                              sql.NullString `db:"comment"`
	AccountLocator                       string         `db:"account_locator"`
	AccountLocatorURL                    sql.NullString `db:"account_locator_url"`
	ManagedAccounts                      sql.NullInt32  `db:"managed_accounts"`
	ConsumptionBillingEntityName         sql.NullString `db:"consumption_billing_entity_name"`
	MarketplaceConsumerBillingEntityName sql.NullString `db:"marketplace_consumer_billing_entity_name"`
	MarketplaceProviderBillingEntityName sql.NullString `db:"marketplace_provider_billing_entity_name"`
	OldAccountURL                        sql.NullString `db:"old_account_url"`
	IsOrgAdmin                           sql.NullBool   `db:"is_org_admin"`
	AccountOldUrlSavedOn                 sql.NullTime   `db:"account_old_url_saved_on"`
	AccountOldUrlLastUsed                sql.NullTime   `db:"account_old_url_last_used"`
	OrganizationOldUrl                   sql.NullString `db:"organization_old_url"`
	OrganizationOldUrlSavedOn            sql.NullTime   `db:"organization_old_url_saved_on"`
	OrganizationOldUrlLastUsed           sql.NullTime   `db:"organization_old_url_last_used"`
	IsEventsAccount                      sql.NullBool   `db:"is_events_account"`
	IsOrganizationAccount                bool           `db:"is_organization_account"`
	// Available only with the History keyword set
	DroppedOn                   sql.NullTime   `db:"dropped_on"`
	ScheduledDeletionTime       sql.NullTime   `db:"scheduled_deletion_time"`
	RestoredOn                  sql.NullTime   `db:"restored_on"`
	MovedToOrganization         sql.NullString `db:"moved_to_organization"`
	MovedOn                     sql.NullString `db:"moved_on"`
	OrganizationUrlExpirationOn sql.NullTime   `db:"organization_URL_expiration_on"`
}

func (row accountDBRow) convert() *Account {
	acc := &Account{
		OrganizationName:      row.OrganizationName,
		AccountName:           row.AccountName,
		SnowflakeRegion:       row.SnowflakeRegion,
		AccountLocator:        row.AccountLocator,
		IsOrganizationAccount: row.IsOrganizationAccount,
	}
	if row.RegionGroup.Valid {
		acc.RegionGroup = &row.RegionGroup.String
	}
	if row.Edition.Valid {
		acc.Edition = Pointer(AccountEdition(row.Edition.String))
	}
	if row.AccountURL.Valid {
		acc.AccountURL = &row.AccountURL.String
	}
	if row.CreatedOn.Valid {
		acc.CreatedOn = &row.CreatedOn.Time
	}
	if row.Comment.Valid {
		acc.Comment = &row.Comment.String
	}
	if row.AccountLocatorURL.Valid {
		acc.AccountLocatorUrl = &row.AccountLocatorURL.String
	}
	if row.ManagedAccounts.Valid {
		acc.ManagedAccounts = Int(int(row.ManagedAccounts.Int32))
	}
	if row.ConsumptionBillingEntityName.Valid {
		acc.ConsumptionBillingEntityName = &row.ConsumptionBillingEntityName.String
	}
	if row.OldAccountURL.Valid {
		acc.OldAccountURL = &row.OldAccountURL.String
	}
	if row.IsOrgAdmin.Valid {
		acc.IsOrgAdmin = &row.IsOrgAdmin.Bool
	}
	if row.OrganizationOldUrl.Valid {
		acc.OrganizationOldUrl = &row.OrganizationOldUrl.String
	}
	if row.IsEventsAccount.Valid {
		acc.IsEventsAccount = &row.IsEventsAccount.Bool
	}
	if row.MarketplaceConsumerBillingEntityName.Valid {
		acc.MarketplaceConsumerBillingEntityName = &row.MarketplaceConsumerBillingEntityName.String
	}
	if row.MarketplaceProviderBillingEntityName.Valid {
		acc.MarketplaceProviderBillingEntityName = &row.MarketplaceProviderBillingEntityName.String
	}
	if row.AccountOldUrlSavedOn.Valid {
		acc.AccountOldUrlSavedOn = &row.AccountOldUrlSavedOn.Time
	}
	if row.AccountOldUrlLastUsed.Valid {
		acc.AccountOldUrlLastUsed = &row.AccountOldUrlLastUsed.Time
	}
	if row.OrganizationOldUrlSavedOn.Valid {
		acc.OrganizationOldUrlSavedOn = &row.OrganizationOldUrlSavedOn.Time
	}
	if row.OrganizationOldUrlLastUsed.Valid {
		acc.OrganizationOldUrlLastUsed = &row.OrganizationOldUrlLastUsed.Time
	}
	if row.DroppedOn.Valid {
		acc.DroppedOn = &row.DroppedOn.Time
	}
	if row.ScheduledDeletionTime.Valid {
		acc.ScheduledDeletionTime = &row.ScheduledDeletionTime.Time
	}
	if row.RestoredOn.Valid {
		acc.RestoredOn = &row.RestoredOn.Time
	}
	if row.MovedToOrganization.Valid {
		acc.MovedToOrganization = &row.MovedToOrganization.String
	}
	if row.MovedOn.Valid {
		acc.MovedOn = &row.MovedOn.String
	}
	if row.OrganizationUrlExpirationOn.Valid {
		acc.OrganizationUrlExpirationOn = &row.OrganizationUrlExpirationOn.Time
	}
	return acc
}

func (c *accounts) Show(ctx context.Context, opts *ShowAccountOptions) ([]Account, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[accountDBRow](c.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[accountDBRow, Account](dbRows)
	return resultList, nil
}

func (c *accounts) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Account, error) {
	accounts, err := c.Show(ctx, &ShowAccountOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(accounts, func(account Account) bool {
		return account.AccountName == id.Name()
	})
}

func (c *accounts) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Account, error) {
	return SafeShowById(c.client, c.ShowByID, ctx, id)
}

// DropAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-account.
type DropAccountOptions struct {
	drop              bool                    `ddl:"static" sql:"DROP"`
	account           bool                    `ddl:"static" sql:"ACCOUNT"`
	IfExists          *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name              AccountObjectIdentifier `ddl:"identifier"`
	gracePeriodInDays int                     `ddl:"parameter" sql:"GRACE_PERIOD_IN_DAYS"`
}

func (opts *DropAccountOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (c *accounts) Drop(ctx context.Context, id AccountObjectIdentifier, gracePeriodInDays int, opts *DropAccountOptions) error {
	if opts == nil {
		opts = &DropAccountOptions{}
	}
	opts.name = id
	opts.gracePeriodInDays = gracePeriodInDays
	return validateAndExec(c.client, ctx, opts)
}

func (c *accounts) DropSafely(ctx context.Context, id AccountObjectIdentifier, gracePeriodInDays int) error {
	return SafeDrop(c.client, func() error { return c.Drop(ctx, id, gracePeriodInDays, &DropAccountOptions{IfExists: Bool(true)}) }, ctx, id)
}

// undropAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/undrop-account.
type undropAccountOptions struct {
	undrop  bool                    `ddl:"static" sql:"UNDROP"`
	account bool                    `ddl:"static" sql:"ACCOUNT"`
	name    AccountObjectIdentifier `ddl:"identifier"`
}

func (c *accounts) Undrop(ctx context.Context, id AccountObjectIdentifier) error {
	opts := &undropAccountOptions{
		name: id,
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, sql)
	return err
}

func (c *accounts) ShowParameters(ctx context.Context) ([]*Parameter, error) {
	return c.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Account: Bool(true),
		},
	})
}

func (c *accounts) UnsetAllParameters(ctx context.Context) error {
	return c.client.Accounts.Alter(ctx, &AlterAccountOptions{Unset: &AccountUnset{
		Parameters: &AccountParametersUnset{
			AbortDetachedQuery:                               Bool(true),
			ActivePythonProfiler:                             Bool(true),
			AllowClientMFACaching:                            Bool(true),
			AllowIDToken:                                     Bool(true),
			Autocommit:                                       Bool(true),
			BaseLocationPrefix:                               Bool(true),
			BinaryInputFormat:                                Bool(true),
			BinaryOutputFormat:                               Bool(true),
			Catalog:                                          Bool(true),
			CatalogSync:                                      Bool(true),
			ClientEnableLogInfoStatementParameters:           Bool(true),
			ClientEncryptionKeySize:                          Bool(true),
			ClientMemoryLimit:                                Bool(true),
			ClientMetadataRequestUseConnectionCtx:            Bool(true),
			ClientMetadataUseSessionDatabase:                 Bool(true),
			ClientPrefetchThreads:                            Bool(true),
			ClientResultChunkSize:                            Bool(true),
			ClientResultColumnCaseInsensitive:                Bool(true),
			ClientSessionKeepAlive:                           Bool(true),
			ClientSessionKeepAliveHeartbeatFrequency:         Bool(true),
			ClientTimestampTypeMapping:                       Bool(true),
			CortexEnabledCrossRegion:                         Bool(true),
			CortexModelsAllowlist:                            Bool(true),
			CsvTimestampFormat:                               Bool(true),
			DataRetentionTimeInDays:                          Bool(true),
			DateInputFormat:                                  Bool(true),
			DateOutputFormat:                                 Bool(true),
			DefaultDDLCollation:                              Bool(true),
			DefaultNotebookComputePoolCpu:                    Bool(true),
			DefaultNotebookComputePoolGpu:                    Bool(true),
			DefaultNullOrdering:                              Bool(true),
			DefaultStreamlitNotebookWarehouse:                Bool(true),
			DisableUiDownloadButton:                          Bool(true),
			DisableUserPrivilegeGrants:                       Bool(true),
			EnableAutomaticSensitiveDataClassificationLog:    Bool(true),
			EnableEgressCostOptimizer:                        Bool(true),
			EnableIdentifierFirstLogin:                       Bool(true),
			EnableInternalStagesPrivatelink:                  Bool(true),
			EnableTriSecretAndRekeyOptOutForImageRepository:  Bool(true),
			EnableTriSecretAndRekeyOptOutForSpcsBlockStorage: Bool(true),
			EnableUnhandledExceptionsReporting:               Bool(true),
			EnableUnloadPhysicalTypeOptimization:             Bool(true),
			EnableUnredactedQuerySyntaxError:                 Bool(true),
			EnableUnredactedSecureObjectError:                Bool(true),
			EnforceNetworkRulesForInternalStages:             Bool(true),
			ErrorOnNondeterministicMerge:                     Bool(true),
			ErrorOnNondeterministicUpdate:                    Bool(true),
			EventTable:                                       Bool(true),
			ExternalOAuthAddPrivilegedRolesToBlockedList:     Bool(true),
			ExternalVolume:                                   Bool(true),
			GeographyOutputFormat:                            Bool(true),
			GeometryOutputFormat:                             Bool(true),
			HybridTableLockTimeout:                           Bool(true),
			InitialReplicationSizeLimitInTB:                  Bool(true),
			JdbcTreatDecimalAsInt:                            Bool(true),
			JdbcTreatTimestampNtzAsUtc:                       Bool(true),
			JdbcUseSessionTimezone:                           Bool(true),
			JsonIndent:                                       Bool(true),
			JsTreatIntegerAsBigInt:                           Bool(true),
			ListingAutoFulfillmentReplicationRefreshSchedule: Bool(true),
			LockTimeout:                                      Bool(true),
			LogLevel:                                         Bool(true),
			MaxConcurrencyLevel:                              Bool(true),
			MaxDataExtensionTimeInDays:                       Bool(true),
			MetricLevel:                                      Bool(true),
			MinDataRetentionTimeInDays:                       Bool(true),
			MultiStatementCount:                              Bool(true),
			NetworkPolicy:                                    Bool(true),
			NoorderSequenceAsDefault:                         Bool(true),
			OAuthAddPrivilegedRolesToBlockedList:             Bool(true),
			OdbcTreatDecimalAsInt:                            Bool(true),
			PeriodicDataRekeying:                             Bool(true),
			PipeExecutionPaused:                              Bool(true),
			PreventUnloadToInlineURL:                         Bool(true),
			PreventUnloadToInternalStages:                    Bool(true),
			PythonProfilerModules:                            Bool(true),
			PythonProfilerTargetStage:                        Bool(true),
			QueryTag:                                         Bool(true),
			QuotedIdentifiersIgnoreCase:                      Bool(true),
			ReplaceInvalidCharacters:                         Bool(true),
			RequireStorageIntegrationForStageCreation:        Bool(true),
			RequireStorageIntegrationForStageOperation:       Bool(true),
			RowsPerResultset:                                 Bool(true),
			S3StageVpceDnsName:                               Bool(true),
			SamlIdentityProvider:                             Bool(true),
			SearchPath:                                       Bool(true),
			ServerlessTaskMaxStatementSize:                   Bool(true),
			ServerlessTaskMinStatementSize:                   Bool(true),
			SimulatedDataSharingConsumer:                     Bool(true),
			SsoLoginPage:                                     Bool(true),
			StatementQueuedTimeoutInSeconds:                  Bool(true),
			StatementTimeoutInSeconds:                        Bool(true),
			StorageSerializationPolicy:                       Bool(true),
			StrictJsonOutput:                                 Bool(true),
			SuspendTaskAfterNumFailures:                      Bool(true),
			TaskAutoRetryAttempts:                            Bool(true),
			TimestampDayIsAlways24h:                          Bool(true),
			TimestampInputFormat:                             Bool(true),
			TimestampLtzOutputFormat:                         Bool(true),
			TimestampNtzOutputFormat:                         Bool(true),
			TimestampOutputFormat:                            Bool(true),
			TimestampTypeMapping:                             Bool(true),
			TimestampTzOutputFormat:                          Bool(true),
			Timezone:                                         Bool(true),
			TimeInputFormat:                                  Bool(true),
			TimeOutputFormat:                                 Bool(true),
			TraceLevel:                                       Bool(true),
			TransactionAbortOnError:                          Bool(true),
			TransactionDefaultIsolationLevel:                 Bool(true),
			TwoDigitCenturyStart:                             Bool(true),
			UnsupportedDdlAction:                             Bool(true),
			UserTaskManagedInitialWarehouseSize:              Bool(true),
			UserTaskMinimumTriggerIntervalInSeconds:          Bool(true),
			UserTaskTimeoutMs:                                Bool(true),
			UseCachedResult:                                  Bool(true),
			WeekOfYearPolicy:                                 Bool(true),
			WeekStart:                                        Bool(true),
		},
	}})
}

func (c *accounts) UnsetAllPoliciesSafely(ctx context.Context) error {
	return errors.Join(
		c.UnsetPolicySafely(ctx, PolicyKindAuthenticationPolicy),
		c.UnsetPolicySafely(ctx, PolicyKindFeaturePolicy),
		c.UnsetPolicySafely(ctx, PolicyKindPackagesPolicy),
		c.UnsetPolicySafely(ctx, PolicyKindPasswordPolicy),
		c.UnsetPolicySafely(ctx, PolicyKindSessionPolicy),
	)
}

func (c *accounts) UnsetPolicySafely(ctx context.Context, kind PolicyKind) error {
	var unset *AccountUnset
	switch kind {
	case PolicyKindAuthenticationPolicy:
		unset = &AccountUnset{AuthenticationPolicy: Bool(true)}
	case PolicyKindFeaturePolicy:
		unset = &AccountUnset{FeaturePolicyUnset: &AccountFeaturePolicyUnset{FeaturePolicy: Bool(true)}}
	case PolicyKindPackagesPolicy:
		unset = &AccountUnset{PackagesPolicy: Bool(true)}
	case PolicyKindPasswordPolicy:
		unset = &AccountUnset{PasswordPolicy: Bool(true)}
	case PolicyKindSessionPolicy:
		unset = &AccountUnset{SessionPolicy: Bool(true)}
	default:
		return fmt.Errorf("policy kind %s is not supported for account policies", kind)
	}
	err := c.client.Accounts.Alter(ctx, &AlterAccountOptions{Unset: unset})
	// If the policy is not attached to the account, Snowflake returns an error.
	if err != nil && strings.Contains(err.Error(), fmt.Sprintf("Any policy of kind %s is not attached to ACCOUNT", kind)) {
		return nil
	}
	return err
}

func (c *accounts) UnsetAllAttachments(ctx context.Context) error {
	return errors.Join(
		c.UnsetAllParameters(ctx),
		c.UnsetAllPoliciesSafely(ctx),
		c.Alter(ctx, &AlterAccountOptions{Unset: &AccountUnset{ResourceMonitor: Bool(true)}}),
	)
}

package sdk

import (
	"context"
	"errors"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

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

func (v *failoverGroups) Show(ctx context.Context, opts *ShowFailoverGroupOptions) ([]FailoverGroup, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[failoverGroupDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[failoverGroupDBRow, FailoverGroup](dbRows)
}

func (v *failoverGroups) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*FailoverGroup, error) {
	failoverGroups, err := v.Show(ctx, &ShowFailoverGroupOptions{})
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(failoverGroups, func(r FailoverGroup) bool {
		return r.Name == id.Name() && r.AccountLocator == v.client.GetAccountLocator()
	})
}

func (v *failoverGroups) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*FailoverGroup, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (row failoverGroupDBRow) convert() (*FailoverGroup, error) {
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

// FailoverGroupDatabase and FailoverGroupShare, their db row types, Request types, and
// ShowFailoverGroupDatabases/ShowFailoverGroupShares methods are generated in Step 3.
// These stubs let the package compile until the generator runs.

// The types below are stubs that will be replaced by generator output in Step 3.

type FailoverGroupDatabase struct {
	Name string
}

type FailoverGroupShare struct {
	Name         string
	OwnerAccount string
}

// ShowFailoverGroupDatabasesOptions and ShowFailoverGroupSharesOptions mirror
// what the generator will produce from the def (public In field, value type).
type ShowFailoverGroupDatabasesOptions struct {
	show      bool                    `ddl:"static" sql:"SHOW"`
	databases bool                    `ddl:"static" sql:"DATABASES"`
	In        AccountObjectIdentifier `ddl:"identifier" sql:"IN FAILOVER GROUP"`
}

func (opts *ShowFailoverGroupDatabasesOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.In) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

type ShowFailoverGroupSharesOptions struct {
	show   bool                    `ddl:"static" sql:"SHOW"`
	shares bool                    `ddl:"static" sql:"SHARES"`
	In     AccountObjectIdentifier `ddl:"identifier" sql:"IN FAILOVER GROUP"`
}

func (opts *ShowFailoverGroupSharesOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.In) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

type ShowFailoverGroupDatabasesRequest struct {
	in AccountObjectIdentifier
}

func NewShowFailoverGroupDatabasesRequest(in AccountObjectIdentifier) *ShowFailoverGroupDatabasesRequest {
	return &ShowFailoverGroupDatabasesRequest{in: in}
}

type ShowFailoverGroupSharesRequest struct {
	in AccountObjectIdentifier
}

func NewShowFailoverGroupSharesRequest(in AccountObjectIdentifier) *ShowFailoverGroupSharesRequest {
	return &ShowFailoverGroupSharesRequest{in: in}
}

func (v *failoverGroups) ShowFailoverGroupDatabases(ctx context.Context, req *ShowFailoverGroupDatabasesRequest) ([]FailoverGroupDatabase, error) {
	opts := &ShowFailoverGroupDatabasesOptions{In: req.in}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dbRows := []failoverGroupDatabaseDBRow{}
	if err := v.client.query(ctx, &dbRows, sql); err != nil {
		return nil, err
	}
	return convertRows[failoverGroupDatabaseDBRow, FailoverGroupDatabase](dbRows)
}

type failoverGroupDatabaseDBRow struct {
	Name string `db:"name"`
}

func (r failoverGroupDatabaseDBRow) convert() (*FailoverGroupDatabase, error) {
	return &FailoverGroupDatabase{Name: r.Name}, nil
}

func (v *failoverGroups) ShowFailoverGroupShares(ctx context.Context, req *ShowFailoverGroupSharesRequest) ([]FailoverGroupShare, error) {
	opts := &ShowFailoverGroupSharesOptions{In: req.in}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dbRows := []failoverGroupShareDBRow{}
	if err := v.client.query(ctx, &dbRows, sql); err != nil {
		return nil, err
	}
	return convertRows[failoverGroupShareDBRow, FailoverGroupShare](dbRows)
}

type failoverGroupShareDBRow struct {
	Name         string `db:"name"`
	OwnerAccount string `db:"owner_account"`
}

func (r failoverGroupShareDBRow) convert() (*FailoverGroupShare, error) {
	return &FailoverGroupShare{Name: r.Name, OwnerAccount: r.OwnerAccount}, nil
}

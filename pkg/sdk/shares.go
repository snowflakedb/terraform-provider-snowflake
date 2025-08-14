package sdk

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ validatable = new(CreateShareOptions)
	_ validatable = new(AlterShareOptions)
	_ validatable = new(DropShareOptions)
	_ validatable = new(ShowShareOptions)
	_ validatable = new(describeShareOptions)

	_ convertibleRow[Share] = new(shareRow)
)

type Shares interface {
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateShareOptions) error
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterShareOptions) error
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropShareOptions) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, opts *ShowShareOptions) ([]Share, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Share, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Share, error)
	DescribeProvider(ctx context.Context, id AccountObjectIdentifier) (*ShareDetails, error)
	DescribeConsumer(ctx context.Context, id ExternalObjectIdentifier) (*ShareDetails, error)
}

var _ Shares = (*shares)(nil)

type shares struct {
	client *Client
}

type ShareKind string

const (
	ShareKindInbound  ShareKind = "INBOUND"
	ShareKindOutbound ShareKind = "OUTBOUND"
)

type Share struct {
	CreatedOn    time.Time
	Kind         ShareKind
	Name         ExternalObjectIdentifier
	DatabaseName AccountObjectIdentifier
	To           []AccountIdentifier
	Owner        string
	Comment      string
}

func (v *Share) ID() AccountObjectIdentifier {
	return v.Name.objectIdentifier.(AccountObjectIdentifier)
}

func (v *Share) ExternalID() ExternalObjectIdentifier {
	return v.Name
}

func (v *Share) ObjectType() ObjectType {
	return ObjectTypeShare
}

type shareRow struct {
	CreatedOn    time.Time `db:"created_on"`
	Kind         string    `db:"kind"`
	OwnerAccount string    `db:"owner_account"`
	Name         string    `db:"name"`
	DatabaseName string    `db:"database_name"`
	To           string    `db:"to"`
	Owner        string    `db:"owner"`
	Comment      string    `db:"comment"`
}

func (r shareRow) convert() (*Share, error) {
	toAccounts := strings.Split(r.To, ",")
	var to []AccountIdentifier
	if len(toAccounts) != 0 {
		for _, a := range toAccounts {
			if a == "" {
				continue
			}
			parts := strings.Split(a, ".")
			if len(parts) == 1 {
				accountLocator := parts[0]
				to = append(to, NewAccountIdentifierFromAccountLocator(accountLocator))
				continue
			}
			orgName := parts[0]
			accountName := strings.Join(parts[1:], ".")
			to = append(to, NewAccountIdentifier(orgName, accountName))
		}
	}
	return &Share{
		CreatedOn:    r.CreatedOn,
		Kind:         ShareKind(r.Kind),
		Name:         NewExternalObjectIdentifier(NewAccountIdentifierFromFullyQualifiedName(r.OwnerAccount), NewAccountObjectIdentifier(r.Name)),
		DatabaseName: NewAccountObjectIdentifier(r.DatabaseName),
		To:           to,
		Owner:        r.Owner,
		Comment:      r.Comment,
	}, nil
}

// CreateShareOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-share.
type CreateShareOptions struct {
	create    bool                    `ddl:"static" sql:"CREATE"`
	OrReplace *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	share     bool                    `ddl:"static" sql:"SHARE"`
	name      AccountObjectIdentifier `ddl:"identifier"`
	Comment   *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (opts *CreateShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (s *shares) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateShareOptions) error {
	if opts == nil {
		opts = &CreateShareOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = s.client.exec(ctx, sql)
	return err
}

// DropShareOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-share.
type DropShareOptions struct {
	drop     bool                    `ddl:"static" sql:"DROP"`
	share    bool                    `ddl:"static" sql:"SHARE"`
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *DropShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (s *shares) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropShareOptions) error {
	opts = createIfNil(opts)
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = s.client.exec(ctx, sql)
	return err
}

func (s *shares) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(s.client, func() error { return s.Drop(ctx, id, &DropShareOptions{IfExists: Bool(true)}) }, ctx, id)
}

// AlterShareOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-share.
type AlterShareOptions struct {
	alter    bool                    `ddl:"static" sql:"ALTER"`
	share    bool                    `ddl:"static" sql:"SHARE"`
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`
	Add      *ShareAdd               `ddl:"keyword" sql:"ADD"`
	Remove   *ShareRemove            `ddl:"keyword" sql:"REMOVE"`
	Set      *ShareSet               `ddl:"keyword" sql:"SET"`
	Unset    *ShareUnset             `ddl:"keyword" sql:"UNSET"`
	SetTag   []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTag []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
}

func (opts *AlterShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Add, opts.Remove, opts.Set, opts.Unset, opts.SetTag, opts.UnsetTag) {
		errs = append(errs, errExactlyOneOf("AlterShareOptions", "Add", "Remove", "Set", "Unset", "SetTag", "UnsetTag"))
	}
	if valueSet(opts.Add) {
		if err := opts.Add.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Remove) {
		if err := opts.Remove.validate(); err != nil {
			errs = append(errs, err)
		}
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

type ShareAdd struct {
	Accounts          []AccountIdentifier `ddl:"parameter" sql:"ACCOUNTS"`
	ShareRestrictions *bool               `ddl:"parameter" sql:"SHARE_RESTRICTIONS"`
}

func (v *ShareAdd) validate() error {
	if len(v.Accounts) == 0 {
		return fmt.Errorf("at least one account must be specified")
	}
	return nil
}

type ShareRemove struct {
	Accounts []AccountIdentifier `ddl:"parameter" sql:"ACCOUNTS"`
}

func (v *ShareRemove) validate() error {
	if len(v.Accounts) == 0 {
		return fmt.Errorf("at least one account must be specified")
	}
	return nil
}

type ShareSet struct {
	Accounts []AccountIdentifier `ddl:"parameter" sql:"ACCOUNTS"`
	Comment  *string             `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (v *ShareSet) validate() error {
	if !anyValueSet(v.Accounts, v.Comment) {
		return errAtLeastOneOf("ShareSet", "Accounts", "Comment")
	}
	return nil
}

type ShareUnset struct {
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}

func (v *ShareUnset) validate() error {
	if !exactlyOneValueSet(v.Comment) {
		return errExactlyOneOf("ShareUnset", "Comment")
	}
	return nil
}

func (s *shares) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterShareOptions) error {
	if opts == nil {
		opts = &AlterShareOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = s.client.exec(ctx, sql)
	return err
}

// ShowShareOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-shares.
type ShowShareOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"`
	shares     bool       `ddl:"static" sql:"SHARES"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	StartsWith *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit      *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

func (opts *ShowShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (s *shares) Show(ctx context.Context, opts *ShowShareOptions) ([]Share, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[shareRow](s.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[shareRow, Share](dbRows)
}

func (s *shares) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Share, error) {
	shares, err := s.Show(ctx, &ShowShareOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(shares, func(share Share) bool {
		return share.ID().FullyQualifiedName() == id.FullyQualifiedName()
	})
}

func (s *shares) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Share, error) {
	return SafeShowById(s.client, s.ShowByID, ctx, id)
}

type ShareDetails struct {
	SharedObjects []ShareInfo
}

type ShareInfo struct {
	Kind     ObjectType
	Name     ObjectIdentifier
	SharedOn time.Time
}

type shareDetailsRow struct {
	Kind     string    `db:"kind"`
	Name     string    `db:"name"`
	SharedOn time.Time `db:"shared_on"`
}

func (row *shareDetailsRow) toShareInfo() *ShareInfo {
	objectType := ObjectType(row.Kind)
	trimmedS := strings.Trim(row.Name, "\"")
	// TODO(SNOW-1229218): Use a common mapper to get object id.
	id := objectType.GetObjectIdentifier(trimmedS)
	return &ShareInfo{
		Kind:     objectType,
		Name:     id,
		SharedOn: row.SharedOn,
	}
}

func shareDetailsFromRows(rows []shareDetailsRow) *ShareDetails {
	v := &ShareDetails{}
	for _, row := range rows {
		v.SharedObjects = append(v.SharedObjects, *row.toShareInfo())
	}
	return v
}

// describeShareOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-share.
type describeShareOptions struct {
	describe bool             `ddl:"static" sql:"DESCRIBE"`
	share    bool             `ddl:"static" sql:"SHARE"`
	name     ObjectIdentifier `ddl:"identifier"`
}

func (opts *describeShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (s *shares) DescribeProvider(ctx context.Context, id AccountObjectIdentifier) (*ShareDetails, error) {
	opts := &describeShareOptions{
		name: id,
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []shareDetailsRow
	err = s.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	return shareDetailsFromRows(rows), nil
}

func (s *shares) DescribeConsumer(ctx context.Context, id ExternalObjectIdentifier) (*ShareDetails, error) {
	opts := &describeShareOptions{
		name: id,
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []shareDetailsRow
	err = s.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	return shareDetailsFromRows(rows), nil
}

package sdk

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ Shares                    = (*shares)(nil)
	_ convertibleRow[Share]     = new(shareRow)
	_ convertibleRow[ShareInfo] = new(shareDetailsRow)
)

type shares struct {
	client *Client
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

func (r shareDetailsRow) convert() (*ShareInfo, error) {
	result := &ShareInfo{
		Kind:     ObjectType(r.Kind),
		SharedOn: r.SharedOn,
	}
	if err := r.additionalConvert(result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *shares) Describe(ctx context.Context, id AccountObjectIdentifier) ([]ShareInfo, error) {
	opts := &describeShareOptions{name: id}
	dbRows, err := validateAndQuery[shareDetailsRow](s.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[shareDetailsRow, ShareInfo](dbRows)
}

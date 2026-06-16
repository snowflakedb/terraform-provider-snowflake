package sdk

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Databases = (*databases)(nil)

type databases struct {
	client *Client
}

func (v *databases) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateDatabaseOptions) error {
	if opts == nil {
		opts = &CreateDatabaseOptions{}
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

func (v *databases) CreateShared(ctx context.Context, id AccountObjectIdentifier, shareID ExternalObjectIdentifier, opts *CreateSharedDatabaseOptions) error {
	if opts == nil {
		opts = &CreateSharedDatabaseOptions{}
	}

	opts.name = id
	opts.fromShare = shareID

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

func (v *databases) CreateSecondary(ctx context.Context, id AccountObjectIdentifier, primaryID ExternalObjectIdentifier, opts *CreateSecondaryDatabaseOptions) error {
	if opts == nil {
		opts = &CreateSecondaryDatabaseOptions{}
	}
	opts.name = id
	opts.primaryDatabase = primaryID
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

func (v *databases) CreateFromListing(ctx context.Context, id AccountObjectIdentifier, listingGlobalName string, opts *CreateDatabaseFromListingOptions) error {
	if opts == nil {
		opts = &CreateDatabaseFromListingOptions{}
	}
	opts.name = id
	opts.fromListing = listingGlobalName
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

func (v *databases) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseOptions) error {
	if opts == nil {
		opts = &AlterDatabaseOptions{}
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

func (v *databases) AlterReplication(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseReplicationOptions) error {
	if opts == nil {
		opts = &AlterDatabaseReplicationOptions{}
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

func (v *databases) AlterFailover(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseFailoverOptions) error {
	if opts == nil {
		opts = &AlterDatabaseFailoverOptions{}
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

func (v *databases) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropDatabaseOptions) error {
	if opts == nil {
		opts = &DropDatabaseOptions{}
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

func (v *databases) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropDatabaseOptions{IfExists: Bool(true)}) }, ctx, id)
}

func (v *databases) Undrop(ctx context.Context, id AccountObjectIdentifier) error {
	opts := &undropDatabaseOptions{
		name: id,
	}
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

func (v *databases) Show(ctx context.Context, opts *ShowDatabasesOptions) ([]Database, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[databaseRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[databaseRow, Database](dbRows)
}

func (v *databases) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Database, error) {
	databases, err := v.Show(ctx, &ShowDatabasesOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}

	return collections.FindFirst(databases, func(r Database) bool { return r.Name == id.Name() })
}

func (v *databases) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Database, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (row databaseRow) convert() (*Database, error) {
	database := &Database{
		CreatedOn: row.CreatedOn,
		Name:      row.Name,
	}
	if row.IsDefault.Valid {
		database.IsDefault = row.IsDefault.String == "Y"
	}
	if row.IsCurrent.Valid {
		database.IsCurrent = row.IsCurrent.String == "Y"
	}
	if row.Origin.Valid && row.Origin.String != "" && row.Origin.String != "<revoked>" {
		originId, err := ParseObjectIdentifierString(row.Origin.String)
		if err != nil {
			return nil, fmt.Errorf("unable to parse origin ID: %w", err)
		} else {
			database.Origin = originId
		}
	}
	if row.Owner.Valid {
		database.Owner = row.Owner.String
	}
	if row.Comment.Valid {
		database.Comment = row.Comment.String
	}
	if row.Options.Valid {
		database.Options = row.Options.String
	}
	if row.RetentionTime.Valid {
		retentionTimeInt, err := strconv.Atoi(row.RetentionTime.String)
		if err != nil {
			return nil, fmt.Errorf("unable to parse retention time: %w", err)
		}
		database.RetentionTime = retentionTimeInt
	}
	if row.ResourceGroup.Valid {
		database.ResourceGroup = row.ResourceGroup.String
	}
	if row.DroppedOn.Valid {
		database.DroppedOn = row.DroppedOn.Time
	}
	if row.Options.Valid {
		database.SetTransient(row.Options.String)
	}
	if row.Kind.Valid {
		database.Kind = row.Kind.String
	}
	if row.OwnerRoleType.Valid {
		database.OwnerRoleType = row.OwnerRoleType.String
	}
	return database, nil
}

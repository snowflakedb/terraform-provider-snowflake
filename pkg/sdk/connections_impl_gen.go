package sdk

import (
	"context"
	"log"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Connections = (*connections)(nil)

type connections struct {
	client *Client
}

func (v *connections) Create(ctx context.Context, request *CreateConnectionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *connections) Alter(ctx context.Context, request *AlterConnectionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *connections) Drop(ctx context.Context, request *DropConnectionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *connections) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropConnectionRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *connections) Show(ctx context.Context, request *ShowConnectionRequest) ([]Connection, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[connectionRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[connectionRow, Connection](dbRows)
	return resultList, nil
}

func (v *connections) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Connection, error) {
	request := NewShowConnectionRequest().
		WithLike(Like{Pattern: String(id.Name())})
	connections, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(connections, func(r Connection) bool {
		// manually adjusted
		return r.Name == id.Name() && r.AccountLocator == v.client.GetAccountLocator()
	})
}

func (v *connections) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Connection, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (r *CreateConnectionRequest) toOpts() *CreateConnectionOptions {
	opts := &CreateConnectionOptions{
		IfNotExists: r.IfNotExists,
		name:        r.name,
		AsReplicaOf: r.AsReplicaOf,
		Comment:     r.Comment,
	}
	return opts
}

func (r *AlterConnectionRequest) toOpts() *AlterConnectionOptions {
	opts := &AlterConnectionOptions{
		IfExists: r.IfExists,
		name:     r.name,
		Primary:  r.Primary,
	}

	if r.EnableConnectionFailover != nil {
		opts.EnableConnectionFailover = &EnableConnectionFailover{
			ToAccounts: r.EnableConnectionFailover.ToAccounts,
		}
	}

	if r.DisableConnectionFailover != nil {
		opts.DisableConnectionFailover = &DisableConnectionFailover{}

		if r.DisableConnectionFailover.ToAccounts != nil {
			opts.DisableConnectionFailover.ToAccounts = &ToAccounts{
				Accounts: r.DisableConnectionFailover.ToAccounts.Accounts,
			}
		}
	}

	if r.Set != nil {
		opts.Set = &ConnectionSet{
			Comment: r.Set.Comment,
		}
	}

	if r.Unset != nil {
		opts.Unset = &ConnectionUnset{
			Comment: r.Unset.Comment,
		}
	}

	return opts
}

func (r *DropConnectionRequest) toOpts() *DropConnectionOptions {
	opts := &DropConnectionOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowConnectionRequest) toOpts() *ShowConnectionOptions {
	opts := &ShowConnectionOptions{
		Like: r.Like,
	}
	return opts
}

func (r connectionRow) convert() *Connection {
	c := &Connection{
		SnowflakeRegion:  r.SnowflakeRegion,
		CreatedOn:        r.CreatedOn,
		AccountName:      r.AccountName,
		Name:             r.Name,
		ConnectionUrl:    r.ConnectionUrl,
		OrganizationName: r.OrganizationName,
		AccountLocator:   r.AccountLocator,
	}

	parsedIsPrimary, err := strconv.ParseBool(r.IsPrimary)
	if err != nil {
		log.Printf("[DEBUG] Unable to parse bool is_primary for connection: %v, err = %s", r.IsPrimary, err)
	} else {
		c.IsPrimary = parsedIsPrimary
	}

	primaryExternalId, err := ParseExternalObjectIdentifier(r.Primary)
	if err != nil {
		log.Printf("[DEBUG] Unable to parse primary connection external identifier: %v, err = %s", r.Primary, err)
	} else {
		c.Primary = primaryExternalId
	}

	if allowedToAccounts, err := ParseCommaSeparatedAccountIdentifierArray(r.FailoverAllowedToAccounts); err != nil {
		log.Printf("[DEBUG] Unable to parse account identifier list for enable failover to accounts, err = %v", err)
	} else {
		c.FailoverAllowedToAccounts = allowedToAccounts
	}

	if r.Comment.Valid {
		c.Comment = String(r.Comment.String)
	}

	return c
}

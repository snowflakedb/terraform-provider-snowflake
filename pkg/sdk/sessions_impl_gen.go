package sdk

import "context"

var _ Sessions = (*sessions)(nil)

type sessions struct {
	client *Client
}

func (v *sessions) AlterSession(ctx context.Context, opts *AlterSessionOptions) error {
	if opts == nil {
		opts = &AlterSessionOptions{}
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

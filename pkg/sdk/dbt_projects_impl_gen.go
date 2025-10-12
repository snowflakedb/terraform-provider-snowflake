package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ DbtProjects = (*dbtProjects)(nil)

var (
	_ convertibleRow[DbtProject]        = new(dbtProjectDBRow)
	_ convertibleRow[DbtProjectDetails] = new(dbtProjectDetailsRow)
)

type dbtProjects struct {
	client *Client
}

func (v *dbtProjects) Create(ctx context.Context, request *CreateDbtProjectRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *dbtProjects) Alter(ctx context.Context, request *AlterDbtProjectRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *dbtProjects) Drop(ctx context.Context, request *DropDbtProjectRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *dbtProjects) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropDbtProjectRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *dbtProjects) Show(ctx context.Context, request *ShowDbtProjectRequest) ([]DbtProject, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[dbtProjectDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[dbtProjectDBRow, DbtProject](dbRows)
}

func (v *dbtProjects) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*DbtProject, error) {
	request := NewShowDbtProjectRequest().
		WithIn(In{Schema: id.SchemaId()}).
		WithLike(Like{Pattern: String(id.Name())})
	dbtProjects, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(dbtProjects, func(r DbtProject) bool { return r.Name == id.Name() })
}

func (v *dbtProjects) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*DbtProject, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *dbtProjects) Describe(ctx context.Context, id SchemaObjectIdentifier) (*DbtProjectDetails, error) {
	opts := &DescribeDbtProjectOptions{
		name: id,
	}
	result, err := validateAndQueryOne[dbtProjectDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert()
}

func (r *CreateDbtProjectRequest) toOpts() *CreateDbtProjectOptions {
	opts := &CreateDbtProjectOptions{
		create:         true, // Set static field to true
		OrReplace:      r.OrReplace,
		dbtProject:     true, // Set static field to true
		IfNotExists:    r.IfNotExists,
		name:           r.name,
		From:           r.From,
		DefaultArgs:    r.DefaultArgs,
		DefaultVersion: r.DefaultVersion,
		Comment:        r.Comment,
	}
	return opts
}

func (r *AlterDbtProjectRequest) toOpts() *AlterDbtProjectOptions {
	opts := &AlterDbtProjectOptions{
		alter:      true, // Set static field to true
		dbtProject: true, // Set static field to true
		IfExists:   r.IfExists,
		name:       r.name,
	}
	if r.Set != nil {
		opts.Set = &DbtProjectSet{
			DefaultArgs:    r.Set.DefaultArgs,
			DefaultVersion: r.Set.DefaultVersion,
			Comment:        r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &DbtProjectUnset{
			DefaultArgs:    r.Unset.DefaultArgs,
			DefaultVersion: r.Unset.DefaultVersion,
			Comment:        r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropDbtProjectRequest) toOpts() *DropDbtProjectOptions {
	opts := &DropDbtProjectOptions{
		drop:       true, // Set static field to true
		dbtProject: true, // Set static field to true
		IfExists:   r.IfExists,
		name:       r.name,
	}
	return opts
}

func (r *ShowDbtProjectRequest) toOpts() *ShowDbtProjectOptions {
	opts := &ShowDbtProjectOptions{
		show:        true, // Set static field to true
		dbtProjects: true, // Set static field to true
		Like:        r.Like,
		In:          r.In,
	}
	return opts
}

func (r dbtProjectDBRow) convert() (*DbtProject, error) {
	dbtProject := &DbtProject{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		DatabaseName:  r.DatabaseName,
		SchemaName:    r.SchemaName,
		Owner:         r.Owner,
		OwnerRoleType: r.OwnerRoleType,
	}

	// Convert sql.NullString to *string
	if r.SourceLocation.Valid {
		dbtProject.SourceLocation = &r.SourceLocation.String
	}
	if r.DefaultArgs.Valid {
		dbtProject.DefaultArgs = &r.DefaultArgs.String
	}
	if r.DefaultVersion.Valid {
		dbtProject.DefaultVersion = &r.DefaultVersion.String
	}
	if r.Comment.Valid {
		dbtProject.Comment = &r.Comment.String
	}

	return dbtProject, nil
}

func (r *DescribeDbtProjectRequest) toOpts() *DescribeDbtProjectOptions {
	opts := &DescribeDbtProjectOptions{
		describe:   true, // Set static field to true
		dbtProject: true, // Set static field to true
		name:       r.name,
	}
	return opts
}

func (r dbtProjectDetailsRow) convert() (*DbtProjectDetails, error) {
	return &DbtProjectDetails{
		Property: r.Property,
		Value:    r.Value,
	}, nil
}

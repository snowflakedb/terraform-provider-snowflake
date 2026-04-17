package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ Tags                = (*tags)(nil)
	_ convertibleRow[Tag] = new(tagRow)
)

type tags struct {
	client *Client
}

func (v *tags) Create(ctx context.Context, request *CreateTagRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tags) Alter(ctx context.Context, request *AlterTagRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tags) Show(ctx context.Context, request *ShowTagRequest) ([]Tag, error) {
	opts := request.toOpts()
	rows, err := validateAndQuery[tagRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[tagRow, Tag](rows)
}

func (v *tags) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Tag, error) {
	request := NewShowTagRequest().WithIn(ExtendedIn{
		In: In{
			Schema: id.SchemaId(),
		},
	}).WithLike(Like{Pattern: String(id.Name())})

	tags, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(tags, func(r Tag) bool { return r.Name == id.Name() })
}

func (v *tags) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Tag, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *tags) Drop(ctx context.Context, request *DropTagRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tags) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropTagRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *tags) Undrop(ctx context.Context, request *UndropTagRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tags) Set(ctx context.Context, request *SetTagRequest) error {
	objectType, err := normalizeGetTagObjectType(request.objectType)
	if err != nil {
		return err
	}
	request.objectType = objectType

	// TODO [SNOW-1022645]: use query from resource sdk - similarly to https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/0e88e082282adf35f605c323569908a99bd406f9/pkg/acceptance/check_destroy.go#L67
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tags) Unset(ctx context.Context, request *UnsetTagRequest) error {
	objectType, err := normalizeGetTagObjectType(request.objectType)
	if err != nil {
		return err
	}
	request.objectType = objectType

	// TODO [SNOW-1022645]: use query from resource sdk - similarly to https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/0e88e082282adf35f605c323569908a99bd406f9/pkg/acceptance/check_destroy.go#L67
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tags) UnsetSafely(ctx context.Context, request *UnsetTagRequest) error {
	return SafeUnsetTag(func() error {
		return v.Unset(ctx, request.WithIfExists(true))
	})
}

func (v *tags) SetOnCurrentAccount(ctx context.Context, request *SetTagOnCurrentAccountRequest) error {
	return v.client.Accounts.Alter(ctx, &AlterAccountOptions{
		SetTag: request.SetTags,
	})
}

func (v *tags) UnsetOnCurrentAccount(ctx context.Context, request *UnsetTagOnCurrentAccountRequest) error {
	return v.client.Accounts.Alter(ctx, &AlterAccountOptions{
		UnsetTag: request.UnsetTags,
	})
}

func (s *CreateTagRequest) toOpts() *createTagOptions {
	opts := &createTagOptions{
		OrReplace:   s.orReplace,
		IfNotExists: s.ifNotExists,
		name:        s.name,
		Propagate:   s.propagate.toTagPropagate(),
		Comment:     s.comment,
	}
	if len(s.allowedValues) > 0 {
		opts.AllowedValues = createAllowedValues(s.allowedValues)
	}
	return opts
}

func (s *AlterTagRequest) toOpts() *alterTagOptions {
	opts := &alterTagOptions{
		name:     s.name,
		ifExists: s.ifExists,
	}
	if len(s.add) > 0 {
		opts.Add = &TagAdd{AllowedValues: createAllowedValues(s.add)}
	}
	if len(s.drop) > 0 {
		opts.Drop = &TagDrop{AllowedValues: createAllowedValues(s.drop)}
	}
	if s.set != nil {
		set := &TagSet{
			Propagate: s.set.propagate.toTagPropagate(),
			Comment:   s.set.comment,
		}
		if len(s.set.allowedValues) > 0 {
			set.AllowedValues = createAllowedValues(s.set.allowedValues)
		}
		if len(s.set.maskingPolicies) > 0 {
			set.MaskingPolicies = &TagSetMaskingPolicies{
				MaskingPolicies: createTagMaskingPolicies(s.set.maskingPolicies),
				Force:           s.set.force,
			}
		}
		opts.Set = set
	}
	if s.unset != nil {
		unset := &TagUnset{
			AllowedValues: s.unset.allowedValues,
			Propagate:     s.unset.propagate,
			OnConflict:    s.unset.onConflict,
			Comment:       s.unset.comment,
		}
		if len(s.unset.maskingPolicies) > 0 {
			unset.MaskingPolicies = &TagUnsetMaskingPolicies{
				MaskingPolicies: createTagMaskingPolicies(s.unset.maskingPolicies),
			}
		}
		opts.Unset = unset
	}
	if s.rename != nil {
		opts.Rename = &TagRename{Name: *s.rename}
	}
	return opts
}

func (s *ShowTagRequest) toOpts() *showTagOptions {
	return &showTagOptions{
		Like: s.Like,
		In:   s.In,
	}
}

func (s *DropTagRequest) toOpts() *dropTagOptions {
	return &dropTagOptions{
		IfExists: Bool(s.ifExists),
		name:     s.name,
	}
}

func (s *UndropTagRequest) toOpts() *undropTagOptions {
	return &undropTagOptions{
		name: s.name,
	}
}

func (s *SetTagRequest) toOpts() *setTagOptions {
	o := &setTagOptions{
		objectType: s.objectType,
		objectName: s.objectName,
		SetTags:    s.SetTags,
	}
	// TODO [SNOW-1022645]: discuss how we handle situation like this in the SDK
	if o.objectType == ObjectTypeColumn {
		id, ok := o.objectName.(TableColumnIdentifier)
		if ok {
			o.objectType = ObjectTypeTable
			o.objectName = id.SchemaObjectId()
			o.column = String(id.Name())
		}
	}
	// TODO(SNOW-1818976): Remove this workaround. Currently ALTER "ORGNAME"."ACCOUNTNAME" SET TAG does not work, but ALTER "ACCOUNTNAME" does.
	if o.objectType == ObjectTypeAccount {
		id, ok := o.objectName.(AccountIdentifier)
		if ok {
			o.objectName = NewAccountIdentifierFromFullyQualifiedName(id.AccountName())
		}
	}

	return o
}

func (s *UnsetTagRequest) toOpts() *unsetTagOptions {
	o := &unsetTagOptions{
		objectType: s.objectType,
		IfExists:   s.IfExists,
		objectName: s.objectName,
		UnsetTags:  s.UnsetTags,
	}
	// TODO [SNOW-1022645]: discuss how we handle situation like this in the SDK
	if o.objectType == ObjectTypeColumn {
		id, ok := o.objectName.(TableColumnIdentifier)
		if ok {
			o.objectType = ObjectTypeTable
			o.objectName = id.SchemaObjectId()
			o.column = String(id.Name())
		}
	}

	// TODO(SNOW-1818976): Remove this workaround. Currently ALTER "ORGNAME"."ACCOUNTNAME" SET TAG does not work, but ALTER "ACCOUNTNAME" does.
	if o.objectType == ObjectTypeAccount {
		id, ok := o.objectName.(AccountIdentifier)
		if ok {
			o.objectName = NewAccountIdentifierFromFullyQualifiedName(id.AccountName())
		}
	}

	return o
}

func (r *TagPropagateRequest) toTagPropagate() *TagPropagate {
	if r == nil {
		return nil
	}
	return &TagPropagate{
		PropagationMethod: &r.propagationMethod,
		OnConflict:        r.onConflict,
	}
}

func createAllowedValues(values []string) *AllowedValues {
	items := make([]StringAllowEmpty, 0, len(values))
	for _, value := range values {
		items = append(items, StringAllowEmpty{
			Value: value,
		})
	}
	return &AllowedValues{
		Values: items,
	}
}

func createTagMaskingPolicies(maskingPolicies []SchemaObjectIdentifier) []TagMaskingPolicy {
	items := make([]TagMaskingPolicy, 0, len(maskingPolicies))
	for _, value := range maskingPolicies {
		items = append(items, TagMaskingPolicy{
			Name: value,
		})
	}
	return items
}

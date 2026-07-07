package sdk

import (
	"context"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

var (
	_ MaskingPolicies                      = (*maskingPolicies)(nil)
	_ convertibleRow[MaskingPolicy]        = new(maskingPolicyDBRow)
	_ convertibleRow[MaskingPolicyDetails] = new(maskingPolicyDetailsRow)
)

type maskingPolicies struct {
	client *Client
}

func (v *maskingPolicies) Create(ctx context.Context, id SchemaObjectIdentifier, signature []TableColumnSignature, returns datatypes.DataType, body string, opts *CreateMaskingPolicyOptions) error {
	if opts == nil {
		opts = &CreateMaskingPolicyOptions{}
	}
	opts.name = id
	opts.signature = signature
	opts.returns = returns
	opts.body = body
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

func (v *maskingPolicies) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterMaskingPolicyOptions) error {
	if opts == nil {
		opts = &AlterMaskingPolicyOptions{}
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

func (v *maskingPolicies) Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropMaskingPolicyOptions) error {
	opts = createIfNil(opts)
	opts.name = id
	if err := opts.validate(); err != nil {
		return fmt.Errorf("validate drop options: %w", err)
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

func (v *maskingPolicies) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropMaskingPolicyOptions{IfExists: Bool(true)}) }, ctx, id)
}

func (v *maskingPolicies) Show(ctx context.Context, opts *ShowMaskingPolicyOptions) ([]MaskingPolicy, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[maskingPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[maskingPolicyDBRow, MaskingPolicy](dbRows)
}

func (v *maskingPolicies) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*MaskingPolicy, error) {
	maskingPolicies, err := v.Show(ctx, &ShowMaskingPolicyOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
		In: &ExtendedIn{
			In: In{
				Schema: id.SchemaId(),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(maskingPolicies, func(r MaskingPolicy) bool { return r.Name == id.Name() })
}

func (v *maskingPolicies) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*MaskingPolicy, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *maskingPolicies) Describe(ctx context.Context, id SchemaObjectIdentifier) (*MaskingPolicyDetails, error) {
	opts := &describeMaskingPolicyOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}

	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := maskingPolicyDetailsRow{}
	err = v.client.queryOne(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}

	return dest.convert()
}

func (row maskingPolicyDetailsRow) convert() (*MaskingPolicyDetails, error) {
	result := &MaskingPolicyDetails{
		Name:      row.Name,
		Signature: []TableColumnSignature{},
		Body:      row.Body,
	}
	if v, err := ParseTableColumnSignature(row.Signature); err != nil {
		log.Printf("[DEBUG] parsing table column signature: %v", err)
	} else {
		result.Signature = v
	}
	dataType, err := datatypes.ParseDataType(row.ReturnType)
	if err != nil {
		return nil, fmt.Errorf("parsing return type: %w", err)
	}
	result.ReturnType = dataType
	return result, nil
}

func (row maskingPolicyDBRow) convert() (*MaskingPolicy, error) {
	maskingPolicy := &MaskingPolicy{
		CreatedOn:     row.CreatedOn,
		Name:          row.Name,
		DatabaseName:  row.DatabaseName,
		SchemaName:    row.SchemaName,
		Kind:          row.Kind,
		Owner:         row.Owner,
		Comment:       row.Comment,
		OwnerRoleType: row.OwnerRoleType,
	}

	if row.Options != "" {
		options, err := ParseMaskingPolicyOptions(row.Options)
		if err != nil {
			return nil, fmt.Errorf("converting masking policy row: error unmarshaling options: %w", err)
		}
		maskingPolicy.ExemptOtherPolicies = options.ExemptOtherPolicies
	}

	return maskingPolicy, nil
}

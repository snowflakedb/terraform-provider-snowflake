package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

var _ MaskingPolicies = (*maskingPolicies)(nil)

var (
	_ validatable = new(CreateMaskingPolicyOptions)
	_ validatable = new(AlterMaskingPolicyOptions)
	_ validatable = new(DropMaskingPolicyOptions)
	_ validatable = new(ShowMaskingPolicyOptions)
	_ validatable = new(describeMaskingPolicyOptions)
)

type MaskingPolicies interface {
	Create(ctx context.Context, id SchemaObjectIdentifier, signature []TableColumnSignature, returns datatypes.DataType, expression string, opts *CreateMaskingPolicyOptions) error
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterMaskingPolicyOptions) error
	Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropMaskingPolicyOptions) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, opts *ShowMaskingPolicyOptions) ([]MaskingPolicy, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*MaskingPolicy, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*MaskingPolicy, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*MaskingPolicyDetails, error)
}

type maskingPolicies struct {
	client *Client
}

// CreateMaskingPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-masking-policy.
type CreateMaskingPolicyOptions struct {
	create        bool                   `ddl:"static" sql:"CREATE"`
	OrReplace     *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"`
	IfNotExists   *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`

	// required
	signature []TableColumnSignature `ddl:"keyword,parentheses" sql:"AS"`
	returns   datatypes.DataType     `ddl:"parameter,no_equals" sql:"RETURNS"`
	body      string                 `ddl:"parameter,no_equals" sql:"->"`

	// optional
	Comment             *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExemptOtherPolicies *bool   `ddl:"parameter" sql:"EXEMPT_OTHER_POLICIES"`
}

func (opts *CreateMaskingPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) && *opts.OrReplace && *opts.IfNotExists {
		errs = append(errs, errOneOf("CreateMaskingPolicyOptions", "OrReplace", "IfNotExists"))
	}
	if !valueSet(opts.signature) {
		errs = append(errs, errNotSet("CreateMaskingPolicyOptions", "signature"))
	}
	if !valueSet(opts.returns) {
		errs = append(errs, errNotSet("CreateMaskingPolicyOptions", "returns"))
	}
	if !valueSet(opts.body) {
		errs = append(errs, errNotSet("CreateMaskingPolicyOptions", "body"))
	}
	return errors.Join(errs...)
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

// AlterMaskingPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-masking-policy.
type AlterMaskingPolicyOptions struct {
	alter         bool                    `ddl:"static" sql:"ALTER"`
	maskingPolicy bool                    `ddl:"static" sql:"MASKING POLICY"`
	IfExists      *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier  `ddl:"identifier"`
	NewName       *SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	Set           *MaskingPolicySet       `ddl:"keyword" sql:"SET"`
	Unset         *MaskingPolicyUnset     `ddl:"keyword" sql:"UNSET"`
	SetTag        []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTag      []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
}

func (opts *AlterMaskingPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.NewName != nil {
		if !ValidObjectIdentifier(opts.NewName) {
			errs = append(errs, errInvalidIdentifier("AlterMaskingPolicyOptions", "NewName"))
		}
		if opts.name.DatabaseName() != opts.NewName.DatabaseName() {
			errs = append(errs, ErrDifferentDatabase)
		}
		if opts.name.SchemaName() != opts.NewName.SchemaName() {
			errs = append(errs, ErrDifferentSchema)
		}
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTag, opts.UnsetTag, opts.NewName) {
		errs = append(errs, errExactlyOneOf("AlterMaskingPolicyOptions", "Set", "Unset", "SetTag", "UnsetTag", "NewName"))
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

type MaskingPolicySet struct {
	Body    *string `ddl:"parameter,no_equals" sql:"BODY ->"`
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (v *MaskingPolicySet) validate() error {
	if !exactlyOneValueSet(v.Body, v.Comment) {
		return errExactlyOneOf("MaskingPolicySet", "Body", "Comment")
	}
	return nil
}

type MaskingPolicyUnset struct {
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}

func (v *MaskingPolicyUnset) validate() error {
	if !exactlyOneValueSet(v.Comment) {
		return errExactlyOneOf("MaskingPolicyUnset", "Comment")
	}
	return nil
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

// DropMaskingPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-masking-policy.
type DropMaskingPolicyOptions struct {
	drop          bool                   `ddl:"static" sql:"DROP"`
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"`
	IfExists      *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *DropMaskingPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
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

// ShowMaskingPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-masking-policies.
type ShowMaskingPolicyOptions struct {
	show            bool        `ddl:"static" sql:"SHOW"`
	maskingPolicies bool        `ddl:"static" sql:"MASKING POLICIES"`
	Like            *Like       `ddl:"keyword" sql:"LIKE"`
	In              *ExtendedIn `ddl:"keyword" sql:"IN"`
	Limit           *LimitFrom  `ddl:"keyword" sql:"LIMIT"`
}

func (opts *ShowMaskingPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

// MaskingPolicys is a user friendly result for a CREATE MASKING POLICY query.
type MaskingPolicy struct {
	CreatedOn           time.Time
	Name                string
	DatabaseName        string
	SchemaName          string
	Kind                string
	Owner               string
	Comment             string
	ExemptOtherPolicies bool
	OwnerRoleType       string
}

type MaskingPolicyOptions struct {
	ExemptOtherPolicies bool `json:"EXEMPT_OTHER_POLICIES"`
}

func ParseMaskingPolicyOptions(s string) (MaskingPolicyOptions, error) {
	var options MaskingPolicyOptions
	err := json.Unmarshal([]byte(s), &options)
	if err != nil {
		return MaskingPolicyOptions{}, err
	}

	return options, nil
}

func (v *MaskingPolicy) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *MaskingPolicy) ObjectType() ObjectType {
	return ObjectTypeMaskingPolicy
}

// maskingPolicyDBRow is used to decode the result of a CREATE MASKING POLICY query.
type maskingPolicyDBRow struct {
	CreatedOn     time.Time `db:"created_on"`
	Name          string    `db:"name"`
	DatabaseName  string    `db:"database_name"`
	SchemaName    string    `db:"schema_name"`
	Kind          string    `db:"kind"`
	Owner         string    `db:"owner"`
	Comment       string    `db:"comment"`
	OwnerRoleType string    `db:"owner_role_type"`
	Options       string    `db:"options"`
}

func (row maskingPolicyDBRow) convert() *MaskingPolicy {
	options, err := ParseMaskingPolicyOptions(row.Options)
	if err != nil {
		log.Printf("[DEBUG] converting masking policy row: error unmarshaling options: %v", err)
		log.Printf("[DEBUG] setting exempt_other_policies = false")
		options.ExemptOtherPolicies = false
	}
	return &MaskingPolicy{
		CreatedOn:           row.CreatedOn,
		Name:                row.Name,
		DatabaseName:        row.DatabaseName,
		SchemaName:          row.SchemaName,
		Kind:                row.Kind,
		Owner:               row.Owner,
		Comment:             row.Comment,
		ExemptOtherPolicies: options.ExemptOtherPolicies,
		OwnerRoleType:       row.OwnerRoleType,
	}
}

// List all the masking policies by pattern.
func (v *maskingPolicies) Show(ctx context.Context, opts *ShowMaskingPolicyOptions) ([]MaskingPolicy, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[maskingPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[maskingPolicyDBRow, MaskingPolicy](dbRows)
	return resultList, nil
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

// describeMaskingPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-masking-policy.
type describeMaskingPolicyOptions struct {
	describe      bool                   `ddl:"static" sql:"DESCRIBE"`
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *describeMaskingPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

type MaskingPolicyDetails struct {
	Name       string
	Signature  []TableColumnSignature
	ReturnType datatypes.DataType
	Body       string
}

type maskingPolicyDetailsRow struct {
	Name       string `db:"name"`
	Signature  string `db:"signature"`
	ReturnType string `db:"return_type"`
	Body       string `db:"body"`
}

func (row maskingPolicyDetailsRow) toMaskingPolicyDetails() *MaskingPolicyDetails {
	dataType, err := datatypes.ParseDataType(row.ReturnType)
	if err != nil {
		return nil
	}
	v := &MaskingPolicyDetails{
		Name:       row.Name,
		Signature:  []TableColumnSignature{},
		ReturnType: dataType,
		Body:       row.Body,
	}

	signature, err := ParseTableColumnSignature(row.Signature)
	if err != nil {
		log.Printf("[DEBUG] parsing table column signature: %v", err)
	} else {
		v.Signature = signature
	}
	return v
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

	return dest.toMaskingPolicyDetails(), nil
}

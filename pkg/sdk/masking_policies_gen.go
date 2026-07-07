package sdk

import (
	"context"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
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

type MaskingPolicySet struct {
	Body    *string `ddl:"parameter,no_equals" sql:"BODY ->"`
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type MaskingPolicyUnset struct {
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropMaskingPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-masking-policy.
type DropMaskingPolicyOptions struct {
	drop          bool                   `ddl:"static" sql:"DROP"`
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"`
	IfExists      *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowMaskingPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-masking-policies.
type ShowMaskingPolicyOptions struct {
	show            bool        `ddl:"static" sql:"SHOW"`
	maskingPolicies bool        `ddl:"static" sql:"MASKING POLICIES"`
	Like            *Like       `ddl:"keyword" sql:"LIKE"`
	In              *ExtendedIn `ddl:"keyword" sql:"IN"`
	Limit           *LimitFrom  `ddl:"keyword" sql:"LIMIT"`
}

// MaskingPolicy is a user friendly result for a CREATE MASKING POLICY query.
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

// describeMaskingPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-masking-policy.
type describeMaskingPolicyOptions struct {
	describe      bool                   `ddl:"static" sql:"DESCRIBE"`
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
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

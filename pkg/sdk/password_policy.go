package sdk

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
)

var _ PasswordPolicies = (*passwordPolicies)(nil)

var (
	_ validatable = new(CreatePasswordPolicyOptions)
	_ validatable = new(AlterPasswordPolicyOptions)
	_ validatable = new(DropPasswordPolicyOptions)
	_ validatable = new(ShowPasswordPolicyOptions)
	_ validatable = new(describePasswordPolicyOptions)
)

type PasswordPolicies interface {
	Create(ctx context.Context, id SchemaObjectIdentifier, opts *CreatePasswordPolicyOptions) error
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterPasswordPolicyOptions) error
	Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropPasswordPolicyOptions) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, opts *ShowPasswordPolicyOptions) ([]PasswordPolicy, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*PasswordPolicy, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*PasswordPolicy, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*PasswordPolicyDetails, error)
}

// passwordPolicies implements PasswordPolicies.
type passwordPolicies struct {
	client *Client
}

// CreatePasswordPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-password-policy.
type CreatePasswordPolicyOptions struct {
	create         bool                   `ddl:"static" sql:"CREATE"`
	OrReplace      *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	passwordPolicy bool                   `ddl:"static" sql:"PASSWORD POLICY"`
	IfNotExists    *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name           SchemaObjectIdentifier `ddl:"identifier"`

	PasswordMinLength         *int `ddl:"parameter" sql:"PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         *int `ddl:"parameter" sql:"PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars *int `ddl:"parameter" sql:"PASSWORD_MIN_UPPER_CASE_CHARS"`
	PasswordMinLowerCaseChars *int `ddl:"parameter" sql:"PASSWORD_MIN_LOWER_CASE_CHARS"`
	PasswordMinNumericChars   *int `ddl:"parameter" sql:"PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   *int `ddl:"parameter" sql:"PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMinAgeDays        *int `ddl:"parameter" sql:"PASSWORD_MIN_AGE_DAYS"`
	PasswordMaxAgeDays        *int `ddl:"parameter" sql:"PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        *int `ddl:"parameter" sql:"PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   *int `ddl:"parameter" sql:"PASSWORD_LOCKOUT_TIME_MINS"`
	PasswordHistory           *int `ddl:"parameter" sql:"PASSWORD_HISTORY"`

	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (opts *CreatePasswordPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *passwordPolicies) Create(ctx context.Context, id SchemaObjectIdentifier, opts *CreatePasswordPolicyOptions) error {
	if opts == nil {
		opts = &CreatePasswordPolicyOptions{}
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

// AlterPasswordPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-password-policy.
type AlterPasswordPolicyOptions struct {
	alter          bool                    `ddl:"static" sql:"ALTER"`
	passwordPolicy bool                    `ddl:"static" sql:"PASSWORD POLICY"`
	IfExists       *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name           SchemaObjectIdentifier  `ddl:"identifier"`
	NewName        *SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	Set            *PasswordPolicySet      `ddl:"list,no_parentheses" sql:"SET"`
	Unset          *PasswordPolicyUnset    `ddl:"list,no_parentheses" sql:"UNSET"`
}

func (opts *AlterPasswordPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.NewName) {
		errs = append(errs, errExactlyOneOf("Set", "Unset", "NewName"))
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

type PasswordPolicySet struct {
	PasswordMinLength         *int    `ddl:"parameter" sql:"PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         *int    `ddl:"parameter" sql:"PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars *int    `ddl:"parameter" sql:"PASSWORD_MIN_UPPER_CASE_CHARS"`
	PasswordMinLowerCaseChars *int    `ddl:"parameter" sql:"PASSWORD_MIN_LOWER_CASE_CHARS"`
	PasswordMinNumericChars   *int    `ddl:"parameter" sql:"PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   *int    `ddl:"parameter" sql:"PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMinAgeDays        *int    `ddl:"parameter" sql:"PASSWORD_MIN_AGE_DAYS"`
	PasswordMaxAgeDays        *int    `ddl:"parameter" sql:"PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        *int    `ddl:"parameter" sql:"PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   *int    `ddl:"parameter" sql:"PASSWORD_LOCKOUT_TIME_MINS"`
	PasswordHistory           *int    `ddl:"parameter" sql:"PASSWORD_HISTORY"`
	Comment                   *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (v *PasswordPolicySet) validate() error {
	if everyValueNil(
		v.PasswordMinLength,
		v.PasswordMaxLength,
		v.PasswordMinUpperCaseChars,
		v.PasswordMinLowerCaseChars,
		v.PasswordMinNumericChars,
		v.PasswordMinSpecialChars,
		v.PasswordMinAgeDays,
		v.PasswordMaxAgeDays,
		v.PasswordMaxRetries,
		v.PasswordLockoutTimeMins,
		v.PasswordHistory,
		v.Comment) {
		return errAtLeastOneOf("PasswordPolicySet", "PasswordMinLength", "PasswordMaxLength", "PasswordMinUpperCaseChars", "PasswordMinLowerCaseChars", "PasswordMinNumericChars", "PasswordMinSpecialChars", "PasswordMinAgeDays", "PasswordMaxAgeDays", "PasswordMaxRetries", "PasswordLockoutTimeMins", "PasswordHistory", "Comment")
	}
	return nil
}

type PasswordPolicyUnset struct {
	PasswordMinLength         *bool `ddl:"keyword" sql:"PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         *bool `ddl:"keyword" sql:"PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars *bool `ddl:"keyword" sql:"PASSWORD_MIN_UPPER_CASE_CHARS"`
	PasswordMinLowerCaseChars *bool `ddl:"keyword" sql:"PASSWORD_MIN_LOWER_CASE_CHARS"`
	PasswordMinNumericChars   *bool `ddl:"keyword" sql:"PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   *bool `ddl:"keyword" sql:"PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMinAgeDays        *bool `ddl:"keyword" sql:"PASSWORD_MIN_AGE_DAYS"`
	PasswordMaxAgeDays        *bool `ddl:"keyword" sql:"PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        *bool `ddl:"keyword" sql:"PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   *bool `ddl:"keyword" sql:"PASSWORD_LOCKOUT_TIME_MINS"`
	PasswordHistory           *bool `ddl:"keyword" sql:"PASSWORD_HISTORY"`
	Comment                   *bool `ddl:"keyword" sql:"COMMENT"`
}

func (v *PasswordPolicyUnset) validate() error {
	if everyValueNil(
		v.PasswordMinLength,
		v.PasswordMaxLength,
		v.PasswordMinUpperCaseChars,
		v.PasswordMinLowerCaseChars,
		v.PasswordMinNumericChars,
		v.PasswordMinSpecialChars,
		v.PasswordMinAgeDays,
		v.PasswordMaxAgeDays,
		v.PasswordMaxRetries,
		v.PasswordLockoutTimeMins,
		v.PasswordHistory,
		v.Comment) {
		return errAtLeastOneOf("PasswordPolicyUnset", "PasswordMinLength", "PasswordMaxLength", "PasswordMinUpperCaseChars", "PasswordMinLowerCaseChars", "PasswordMinNumericChars", "PasswordMinSpecialChars", "PasswordMinAgeDays", "PasswordMaxAgeDays", "PasswordMaxRetries", "PasswordLockoutTimeMins", "PasswordHistory", "Comment")
	}
	return nil
}

func (v *passwordPolicies) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterPasswordPolicyOptions) error {
	if opts == nil {
		opts = &AlterPasswordPolicyOptions{}
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

// DropPasswordPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-password-policy.
type DropPasswordPolicyOptions struct {
	drop           bool                   `ddl:"static" sql:"DROP"`
	passwordPolicy bool                   `ddl:"static" sql:"PASSWORD POLICY"`
	IfExists       *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name           SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *DropPasswordPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *passwordPolicies) Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropPasswordPolicyOptions) error {
	if opts == nil {
		opts = &DropPasswordPolicyOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return fmt.Errorf("validate drop options: %w", err)
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

func (v *passwordPolicies) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropPasswordPolicyOptions{IfExists: Bool(true)}) }, ctx, id)
}

// ShowPasswordPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-password-policies.
type ShowPasswordPolicyOptions struct {
	show             bool  `ddl:"static" sql:"SHOW"`
	passwordPolicies bool  `ddl:"static" sql:"PASSWORD POLICIES"`
	Like             *Like `ddl:"keyword" sql:"LIKE"`
	In               *In   `ddl:"keyword" sql:"IN"`
	Limit            *int  `ddl:"parameter,no_equals" sql:"LIMIT"`
}

func (opts *ShowPasswordPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

// PasswordPolicy is a user-friendly result for a CREATE PASSWORD POLICY query.
type PasswordPolicy struct {
	CreatedOn     time.Time
	Name          string
	DatabaseName  string
	SchemaName    string
	Kind          string
	Owner         string
	Comment       string
	OwnerRoleType string
}

func (v *PasswordPolicy) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *PasswordPolicy) ObjectType() ObjectType {
	return ObjectTypePasswordPolicy
}

// passwordPolicyDBRow is used to decode the result of a CREATE PASSWORD POLICY query.
type passwordPolicyDBRow struct {
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

func (row passwordPolicyDBRow) convert() PasswordPolicy {
	return PasswordPolicy{
		CreatedOn:     row.CreatedOn,
		Name:          row.Name,
		DatabaseName:  row.DatabaseName,
		SchemaName:    row.SchemaName,
		Kind:          row.Kind,
		Owner:         row.Owner,
		Comment:       row.Comment,
		OwnerRoleType: row.OwnerRoleType,
	}
}

// List all the password policies by pattern.
func (v *passwordPolicies) Show(ctx context.Context, opts *ShowPasswordPolicyOptions) ([]PasswordPolicy, error) {
	opts = createIfNil(opts)
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var dest []passwordPolicyDBRow
	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	resultList := make([]PasswordPolicy, len(dest))
	for i, row := range dest {
		resultList[i] = row.convert()
	}

	return resultList, nil
}

func (v *passwordPolicies) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*PasswordPolicy, error) {
	passwordPolicies, err := v.Show(ctx, &ShowPasswordPolicyOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
		In: &In{
			Schema: id.SchemaId(),
		},
	})
	if err != nil {
		return nil, err
	}

	return collections.FindFirst(passwordPolicies, func(policy PasswordPolicy) bool {
		return policy.ID().FullyQualifiedName() == id.FullyQualifiedName()
	})
}

func (v *passwordPolicies) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*PasswordPolicy, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

// describePasswordPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-password-policy.
type describePasswordPolicyOptions struct {
	describe       bool                   `ddl:"static" sql:"DESCRIBE"`
	passwordPolicy bool                   `ddl:"static" sql:"PASSWORD POLICY"`
	name           SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *describePasswordPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

type PasswordPolicyDetails struct {
	Name                      *StringProperty
	Owner                     *StringProperty
	Comment                   *StringProperty
	PasswordMinLength         *IntProperty
	PasswordMaxLength         *IntProperty
	PasswordMinUpperCaseChars *IntProperty
	PasswordMinLowerCaseChars *IntProperty
	PasswordMinNumericChars   *IntProperty
	PasswordMinSpecialChars   *IntProperty
	PasswordMinAgeDays        *IntProperty
	PasswordMaxAgeDays        *IntProperty
	PasswordMaxRetries        *IntProperty
	PasswordLockoutTimeMins   *IntProperty
	PasswordHistory           *IntProperty
}

func passwordPolicyDetailsFromRows(rows []propertyRow) *PasswordPolicyDetails {
	v := &PasswordPolicyDetails{}
	for _, row := range rows {
		switch row.Property {
		case "NAME":
			v.Name = row.toStringProperty()
		case "OWNER":
			v.Owner = row.toStringProperty()
		case "COMMENT":
			v.Comment = row.toStringProperty()
		case "PASSWORD_MIN_LENGTH":
			v.PasswordMinLength = row.toIntProperty()
		case "PASSWORD_MAX_LENGTH":
			v.PasswordMaxLength = row.toIntProperty()
		case "PASSWORD_MIN_UPPER_CASE_CHARS":
			v.PasswordMinUpperCaseChars = row.toIntProperty()
		case "PASSWORD_MIN_LOWER_CASE_CHARS":
			v.PasswordMinLowerCaseChars = row.toIntProperty()
		case "PASSWORD_MIN_NUMERIC_CHARS":
			v.PasswordMinNumericChars = row.toIntProperty()
		case "PASSWORD_MIN_SPECIAL_CHARS":
			v.PasswordMinSpecialChars = row.toIntProperty()
		case "PASSWORD_MIN_AGE_DAYS":
			v.PasswordMinAgeDays = row.toIntProperty()
		case "PASSWORD_MAX_AGE_DAYS":
			v.PasswordMaxAgeDays = row.toIntProperty()
		case "PASSWORD_MAX_RETRIES":
			v.PasswordMaxRetries = row.toIntProperty()
		case "PASSWORD_LOCKOUT_TIME_MINS":
			v.PasswordLockoutTimeMins = row.toIntProperty()
		case "PASSWORD_HISTORY":
			v.PasswordHistory = row.toIntProperty()
		}
	}
	return v
}

func (v *passwordPolicies) Describe(ctx context.Context, id SchemaObjectIdentifier) (*PasswordPolicyDetails, error) {
	opts := &describePasswordPolicyOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}

	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []propertyRow{}
	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}

	return passwordPolicyDetailsFromRows(dest), nil
}

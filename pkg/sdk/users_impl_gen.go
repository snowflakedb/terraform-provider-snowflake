package sdk

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ convertibleRow[User] = new(userDBRow)

func (v *users) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateUserOptions) error {
	if opts == nil {
		opts = &CreateUserOptions{}
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

func (v *users) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterUserOptions) error {
	if opts == nil {
		opts = &AlterUserOptions{}
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

func (v *users) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropUserOptions) error {
	if opts == nil {
		opts = &DropUserOptions{}
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
	if err != nil {
		return err
	}
	return err
}

func (v *users) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropUserOptions{IfExists: Bool(true)}) }, ctx, id)
}

// TODO [Step 3]: Describe ([]UserProperty) and ShowUserWorkloadIdentityAuthenticationMethodOptions pre-gen stubs are at the bottom of this file.

func (v *users) Show(ctx context.Context, opts *ShowUserOptions) ([]User, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[userDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[userDBRow, User](dbRows)
}

func (v *users) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*User, error) {
	users, err := v.Show(ctx, &ShowUserOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(users, func(user User) bool {
		return user.ID().Name() == id.Name()
	})
}

func (v *users) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*User, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *users) ShowParameters(ctx context.Context, id AccountObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			User: id,
		},
	})
}

// AddProgrammaticAccessToken, ModifyProgrammaticAccessToken, RotateProgrammaticAccessToken,
// RemoveProgrammaticAccessToken, RemoveProgrammaticAccessTokenSafely, ShowProgrammaticAccessTokens,
// ShowProgrammaticAccessTokenByName, ShowProgrammaticAccessTokenByNameSafely are in users_ext.go — PAT delegation methods.
// ShowUserWorkloadIdentityAuthenticationMethodOptions is generated from wifMethodsPairs + showWifMethodsQueryStruct.

func (row userDBRow) convert() (*User, error) {
	user := &User{
		Name:      row.Name,
		CreatedOn: row.CreatedOn,
		Owner:     row.Owner,
	}
	if row.LoginName.Valid {
		user.LoginName = row.LoginName.String
	}
	if row.DisplayName.Valid {
		user.DisplayName = row.DisplayName.String
	}
	if row.FirstName.Valid {
		user.FirstName = row.FirstName.String
	}
	if row.LastName.Valid {
		user.LastName = row.LastName.String
	}
	if row.Email.Valid {
		user.Email = row.Email.String
	}
	if row.MinsToUnlock.Valid {
		user.MinsToUnlock = row.MinsToUnlock.String
	}
	if row.DaysToExpiry.Valid {
		user.DaysToExpiry = row.DaysToExpiry.String
	}
	if row.Comment.Valid {
		user.Comment = row.Comment.String
	}
	if err := handleNullableBoolString(row.Disabled, &user.Disabled); err != nil {
		return nil, fmt.Errorf("error parsing disabled: %w", err)
	}
	if err := handleNullableBoolString(row.MustChangePassword, &user.MustChangePassword); err != nil {
		return nil, fmt.Errorf("error parsing must change password: %w", err)
	}
	if err := handleNullableBoolString(row.SnowflakeLock, &user.SnowflakeLock); err != nil {
		return nil, fmt.Errorf("error parsing snowflake lock: %w", err)
	}
	if err := handleNullableBoolString(row.ExtAuthnDuo, &user.ExtAuthnDuo); err != nil {
		return nil, fmt.Errorf("error parsing ext authn duo: %w", err)
	}
	if row.ExtAuthnUid.Valid {
		user.ExtAuthnUid = row.ExtAuthnUid.String
	}
	if row.MinsToBypassMfa.Valid {
		user.MinsToBypassMfa = row.MinsToBypassMfa.String
	}
	if row.DefaultWarehouse.Valid {
		user.DefaultWarehouse = row.DefaultWarehouse.String
	}
	if row.DefaultNamespace.Valid {
		user.DefaultNamespace = row.DefaultNamespace.String
	}
	if row.DefaultRole.Valid {
		user.DefaultRole = row.DefaultRole.String
	}
	if row.DefaultSecondaryRoles.Valid {
		user.DefaultSecondaryRoles = row.DefaultSecondaryRoles.String
	}
	if row.LastSuccessLogin.Valid {
		user.LastSuccessLogin = row.LastSuccessLogin.Time
	}
	if row.ExpiresAtTime.Valid {
		user.ExpiresAtTime = row.ExpiresAtTime.Time
	}
	if row.LockedUntilTime.Valid {
		user.LockedUntilTime = row.LockedUntilTime.Time
	}
	if row.HasPassword.Valid {
		user.HasPassword = row.HasPassword.Bool
	}
	if row.HasRsaPublicKey.Valid {
		user.HasRsaPublicKey = row.HasRsaPublicKey.Bool
	}
	if row.Type.Valid {
		user.Type = row.Type.String
	}
	if row.HasMfa.Valid {
		user.HasMfa = row.HasMfa.Bool
	}
	if row.HasWorkloadIdentity.Valid {
		user.HasWorkloadIdentity = row.HasWorkloadIdentity.Bool
	}
	return user, nil
}

// TODO [Step 3]: Remove — generated by wifMethodsPairs.
var _ convertibleRow[UserWorkloadIdentityAuthenticationMethod] = new(userWorkloadIdentityAuthenticationMethodsDBRow)

// TODO [Step 3]: Remove — generated by wifMethodsPairs (calls additionalConvert() in users_ext.go).
func (row userWorkloadIdentityAuthenticationMethodsDBRow) convert() (*UserWorkloadIdentityAuthenticationMethod, error) {
	result := &UserWorkloadIdentityAuthenticationMethod{
		Name:      row.Name,
		CreatedOn: row.CreatedOn,
	}
	if row.Comment.Valid {
		result.Comment = row.Comment.String
	}
	if row.LastUsed.Valid {
		result.LastUsed = row.LastUsed.Time
	}
	if err := row.additionalConvert(result); err != nil {
		return nil, err
	}
	return result, nil
}

// TODO [Step 3]: Remove — generated by describeUserPropertyPairs.
var _ convertibleRow[UserProperty] = new(describeUserPropertyRow)

// TODO [Step 3]: Remove — generated by describeUserPropertyPairs.
func (row describeUserPropertyRow) convert() (*UserProperty, error) {
	return &UserProperty{
		Property:     row.Property,
		Value:        row.Value,
		DefaultValue: row.DefaultValue,
		Description:  row.Description,
	}, nil
}

// TODO [Step 3]: Remove — replaced by generated Describe() from describeUserQueryStruct() + describeUserPropertyPairs.
func (v *users) Describe(ctx context.Context, id AccountObjectIdentifier) ([]UserProperty, error) {
	opts := &describeUserOptions{name: id}
	dbRows, err := validateAndQuery[describeUserPropertyRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[describeUserPropertyRow, UserProperty](dbRows)
}

// TODO [Step 3]: Remove — replaced by generated ShowUserWorkloadIdentityAuthenticationMethodOptions() from showWifMethodsQueryStruct().
func (v *users) ShowUserWorkloadIdentityAuthenticationMethodOptions(ctx context.Context, userId AccountObjectIdentifier) ([]UserWorkloadIdentityAuthenticationMethod, error) {
	opts := &showUserWorkloadIdentityAuthenticationMethodOptionsOptions{
		ForUser: userId,
	}
	dbRows, err := validateAndQuery[userWorkloadIdentityAuthenticationMethodsDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[userWorkloadIdentityAuthenticationMethodsDBRow, UserWorkloadIdentityAuthenticationMethod](dbRows)
}

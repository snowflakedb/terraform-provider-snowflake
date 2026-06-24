package sdk

import (
	"context"
	"encoding/json"
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

func (v *users) Describe(ctx context.Context, id AccountObjectIdentifier) (*UserDetails, error) {
	opts := &describeUserOptions{
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
	return userDetailsFromRows(dest), nil
}

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

func (v *users) AddProgrammaticAccessToken(ctx context.Context, request *AddUserProgrammaticAccessTokenRequest) (*AddProgrammaticAccessTokenResult, error) {
	return v.client.UserProgrammaticAccessTokens.Add(ctx, request)
}

func (v *users) ModifyProgrammaticAccessToken(ctx context.Context, request *ModifyUserProgrammaticAccessTokenRequest) error {
	return v.client.UserProgrammaticAccessTokens.Modify(ctx, request)
}

func (v *users) RotateProgrammaticAccessToken(ctx context.Context, request *RotateUserProgrammaticAccessTokenRequest) (*RotateProgrammaticAccessTokenResult, error) {
	return v.client.UserProgrammaticAccessTokens.Rotate(ctx, request)
}

func (v *users) RemoveProgrammaticAccessToken(ctx context.Context, request *RemoveUserProgrammaticAccessTokenRequest) error {
	return v.client.UserProgrammaticAccessTokens.Remove(ctx, request)
}

func (v *users) RemoveProgrammaticAccessTokenSafely(ctx context.Context, request *RemoveUserProgrammaticAccessTokenRequest) error {
	return v.client.UserProgrammaticAccessTokens.RemoveByIDSafely(ctx, request)
}

func (v *users) ShowProgrammaticAccessTokens(ctx context.Context, request *ShowUserProgrammaticAccessTokenRequest) ([]ProgrammaticAccessToken, error) {
	return v.client.UserProgrammaticAccessTokens.Show(ctx, request)
}

func (v *users) ShowProgrammaticAccessTokenByName(ctx context.Context, userId AccountObjectIdentifier, tokenName AccountObjectIdentifier) (*ProgrammaticAccessToken, error) {
	return v.client.UserProgrammaticAccessTokens.ShowByID(ctx, userId, tokenName)
}

func (v *users) ShowProgrammaticAccessTokenByNameSafely(ctx context.Context, userId AccountObjectIdentifier, tokenName AccountObjectIdentifier) (*ProgrammaticAccessToken, error) {
	return v.client.UserProgrammaticAccessTokens.ShowByIDSafely(ctx, userId, tokenName)
}

func (v *users) ShowUserWorkloadIdentityAuthenticationMethodOptions(ctx context.Context, userId AccountObjectIdentifier) ([]UserWorkloadIdentityAuthenticationMethod, error) {
	opts := &showUserAuthenticationMethodOptions{
		ForUser: userId,
	}
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[userWorkloadIdentityAuthenticationMethodsDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[userWorkloadIdentityAuthenticationMethodsDBRow, UserWorkloadIdentityAuthenticationMethod](dbRows)
}

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

func (row userWorkloadIdentityAuthenticationMethodsDBRow) convert() (*UserWorkloadIdentityAuthenticationMethod, error) {
	methods := &UserWorkloadIdentityAuthenticationMethod{
		Name:      row.Name,
		CreatedOn: row.CreatedOn,
	}
	wifType, err := ToWIFTypeType(row.Type)
	if err != nil {
		return nil, err
	}
	methods.Type = wifType
	switch wifType {
	case WIFTypeAWS:
		additionalInfo := &UserWorkloadIdentityAuthenticationMethodsAwsAdditionalInfo{}
		if err := json.Unmarshal([]byte(row.AdditionalInfo.String), additionalInfo); err != nil {
			return nil, err
		}
		methods.AwsAdditionalInfo = additionalInfo
	case WIFTypeAzure:
		additionalInfo := &UserWorkloadIdentityAuthenticationMethodsAzureAdditionalInfo{}
		if err := json.Unmarshal([]byte(row.AdditionalInfo.String), additionalInfo); err != nil {
			return nil, err
		}
		methods.AzureAdditionalInfo = additionalInfo
	case WIFTypeGCP:
		additionalInfo := &UserWorkloadIdentityAuthenticationMethodsGcpAdditionalInfo{}
		if err := json.Unmarshal([]byte(row.AdditionalInfo.String), additionalInfo); err != nil {
			return nil, err
		}
		methods.GcpAdditionalInfo = additionalInfo
	case WIFTypeOIDC:
		additionalInfo := &UserWorkloadIdentityAuthenticationMethodsOidcAdditionalInfo{}
		if err := json.Unmarshal([]byte(row.AdditionalInfo.String), additionalInfo); err != nil {
			return nil, err
		}
		methods.OidcAdditionalInfo = additionalInfo
	}
	if row.LastUsed.Valid {
		methods.LastUsed = row.LastUsed.Time
	}
	if row.Comment.Valid {
		methods.Comment = row.Comment.String
	}
	return methods, nil
}

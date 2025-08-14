package sdk

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ UserProgrammaticAccessTokens = (*userProgrammaticAccessTokens)(nil)
var _ convertibleRow[ProgrammaticAccessToken] = new(programmaticAccessTokenResultDBRow)
var _ convertibleRow[AddProgrammaticAccessTokenResult] = new(addProgrammaticAccessTokenResultDBRow)
var _ convertibleRow[RotateProgrammaticAccessTokenResult] = new(rotateProgrammaticAccessTokenResultDBRow)

type userProgrammaticAccessTokens struct {
	client *Client
}

func (v *userProgrammaticAccessTokens) Add(ctx context.Context, request *AddUserProgrammaticAccessTokenRequest) (*AddProgrammaticAccessTokenResult, error) {
	opts := request.toOpts()
	result, err := validateAndQueryOne[addProgrammaticAccessTokenResultDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convertErr()
}

func (v *userProgrammaticAccessTokens) Modify(ctx context.Context, request *ModifyUserProgrammaticAccessTokenRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *userProgrammaticAccessTokens) Rotate(ctx context.Context, request *RotateUserProgrammaticAccessTokenRequest) (*RotateProgrammaticAccessTokenResult, error) {
	opts := request.toOpts()
	result, err := validateAndQueryOne[rotateProgrammaticAccessTokenResultDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convertErr()
}

func (v *userProgrammaticAccessTokens) Remove(ctx context.Context, request *RemoveUserProgrammaticAccessTokenRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

// Adjusted manually to include the user id in the request.
func (v *userProgrammaticAccessTokens) RemoveByIDSafely(ctx context.Context, request *RemoveUserProgrammaticAccessTokenRequest) error {
	return SafeRemoveProgrammaticAccessToken(v.client, ctx, request)
}

func (v *userProgrammaticAccessTokens) Show(ctx context.Context, request *ShowUserProgrammaticAccessTokenRequest) ([]ProgrammaticAccessToken, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[programmaticAccessTokenResultDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRowsErr[programmaticAccessTokenResultDBRow, ProgrammaticAccessToken](dbRows)
}

// Adjusted manually to include the user id in the request.
func (v *userProgrammaticAccessTokens) ShowByID(ctx context.Context, userId, id AccountObjectIdentifier) (*ProgrammaticAccessToken, error) {
	request := NewShowUserProgrammaticAccessTokenRequest().WithUserName(userId)
	userProgrammaticAccessTokens, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(userProgrammaticAccessTokens, func(r ProgrammaticAccessToken) bool { return r.Name == id.Name() })
}

// Adjusted manually to include the user id in the request.
func (v *userProgrammaticAccessTokens) ShowByIDSafely(ctx context.Context, userId, id AccountObjectIdentifier) (*ProgrammaticAccessToken, error) {
	return SafeShowProgrammaticAccessTokenByName(v.client, ctx, userId, id)
}

func (r *AddUserProgrammaticAccessTokenRequest) toOpts() *AddUserProgrammaticAccessTokenOptions {
	opts := &AddUserProgrammaticAccessTokenOptions{
		IfExists:                             r.IfExists,
		UserName:                             r.UserName,
		name:                                 r.name,
		RoleRestriction:                      r.RoleRestriction,
		DaysToExpiry:                         r.DaysToExpiry,
		MinsToBypassNetworkPolicyRequirement: r.MinsToBypassNetworkPolicyRequirement,
		Comment:                              r.Comment,
	}
	return opts
}

func (r addProgrammaticAccessTokenResultDBRow) convertErr() (*AddProgrammaticAccessTokenResult, error) {
	return &AddProgrammaticAccessTokenResult{
		TokenName:   r.TokenName,
		TokenSecret: r.TokenSecret,
	}, nil
}

func (r *ModifyUserProgrammaticAccessTokenRequest) toOpts() *ModifyUserProgrammaticAccessTokenOptions {
	opts := &ModifyUserProgrammaticAccessTokenOptions{
		IfExists: r.IfExists,
		UserName: r.UserName,
		name:     r.name,

		RenameTo: r.RenameTo,
	}
	if r.Set != nil {
		opts.Set = &ModifyProgrammaticAccessTokenSet{
			Disabled:                             r.Set.Disabled,
			MinsToBypassNetworkPolicyRequirement: r.Set.MinsToBypassNetworkPolicyRequirement,
			Comment:                              r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &ModifyProgrammaticAccessTokenUnset{
			Disabled:                             r.Unset.Disabled,
			MinsToBypassNetworkPolicyRequirement: r.Unset.MinsToBypassNetworkPolicyRequirement,
			Comment:                              r.Unset.Comment,
		}
	}
	return opts
}

func (r *RotateUserProgrammaticAccessTokenRequest) toOpts() *RotateUserProgrammaticAccessTokenOptions {
	opts := &RotateUserProgrammaticAccessTokenOptions{
		IfExists:                     r.IfExists,
		UserName:                     r.UserName,
		name:                         r.name,
		ExpireRotatedTokenAfterHours: r.ExpireRotatedTokenAfterHours,
	}
	return opts
}

func (r rotateProgrammaticAccessTokenResultDBRow) convertErr() (*RotateProgrammaticAccessTokenResult, error) {
	return &RotateProgrammaticAccessTokenResult{
		TokenName:        r.TokenName,
		TokenSecret:      r.TokenSecret,
		RotatedTokenName: r.RotatedTokenName,
	}, nil
}

func (r *RemoveUserProgrammaticAccessTokenRequest) toOpts() *RemoveUserProgrammaticAccessTokenOptions {
	opts := &RemoveUserProgrammaticAccessTokenOptions{
		IfExists: r.IfExists,
		UserName: r.UserName,
		name:     r.name,
	}
	return opts
}

func (r *ShowUserProgrammaticAccessTokenRequest) toOpts() *ShowUserProgrammaticAccessTokenOptions {
	opts := &ShowUserProgrammaticAccessTokenOptions{
		UserName: r.UserName,
	}
	return opts
}

func (r programmaticAccessTokenResultDBRow) convertErr() (*ProgrammaticAccessToken, error) {
	token := &ProgrammaticAccessToken{
		Name:      r.Name,
		ExpiresAt: r.ExpiresAt,
		CreatedOn: r.CreatedOn,
		CreatedBy: r.CreatedBy,
	}
	userName, err := ParseAccountObjectIdentifier(r.UserName)
	if err != nil {
		return nil, fmt.Errorf("error parsing user name: %w", err)
	} else {
		token.UserName = userName
	}
	if r.RoleRestriction != "" {
		roleRestriction, err := ParseAccountObjectIdentifier(r.RoleRestriction)
		if err != nil {
			return nil, fmt.Errorf("error parsing role restriction: %w", err)
		}
		token.RoleRestriction = &roleRestriction
	}
	status, err := toProgrammaticAccessTokenStatus(r.Status)
	if err != nil {
		return nil, fmt.Errorf("error parsing programmatic access token status: %w", err)
	} else {
		token.Status = status
	}
	if r.Comment.Valid {
		token.Comment = &r.Comment.String
	}
	if r.MinsToBypassNetworkPolicyRequirement.Valid {
		token.MinsToBypassNetworkPolicyRequirement = Pointer(int(r.MinsToBypassNetworkPolicyRequirement.Int64))
	}
	if r.RotatedTo.Valid {
		token.RotatedTo = &r.RotatedTo.String
	}
	return token, nil
}

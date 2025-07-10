package sdk

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

func (v *organizationAccounts) ShowParameters(ctx context.Context) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Account: Bool(true),
		},
	})
}

func (v *organizationAccounts) UnsetAllParameters(ctx context.Context) error {
	return v.client.Accounts.UnsetAllParameters(ctx)
}

func (v *organizationAccounts) UnsetPolicySafely(ctx context.Context, kind PolicyKind) error {
	unset := NewOrganizationAccountUnsetRequest()
	switch kind {
	case PolicyKindPasswordPolicy:
		unset.WithPasswordPolicy(true)
	case PolicyKindSessionPolicy:
		unset.WithSessionPolicy(true)
	default:
		return fmt.Errorf("policy kind %s is not supported for account policies", kind)
	}
	err := v.client.OrganizationAccounts.Alter(ctx, NewAlterOrganizationAccountRequest().WithUnset(*unset))
	// If the policy is not attached to the account, Snowflake returns an error.
	if err != nil && strings.Contains(err.Error(), fmt.Sprintf("Any policy of kind %s is not attached to ACCOUNT", kind)) {
		return nil
	}
	return err
}

func (v *organizationAccounts) SetPolicySafely(ctx context.Context, kind PolicyKind, id SchemaObjectIdentifier) error {
	set := NewOrganizationAccountSetRequest()
	switch kind {
	case PolicyKindPasswordPolicy:
		set.WithPasswordPolicy(id)
	case PolicyKindSessionPolicy:
		set.WithSessionPolicy(id)
	}
	return errors.Join(
		v.UnsetPolicySafely(ctx, kind),
		v.client.OrganizationAccounts.Alter(ctx, NewAlterOrganizationAccountRequest().WithSet(*set)),
	)
}

func (v *organizationAccounts) UnsetAll(ctx context.Context) error {
	return errors.Join(
		v.client.OrganizationAccounts.UnsetAllParameters(ctx),
		v.client.OrganizationAccounts.UnsetPolicySafely(ctx, PolicyKindPasswordPolicy),
		v.client.OrganizationAccounts.UnsetPolicySafely(ctx, PolicyKindSessionPolicy),
	)
}

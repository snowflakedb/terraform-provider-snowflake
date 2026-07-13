package sdk

import (
	"context"
	"errors"
	"fmt"
)

func (v *Account) AccountID() AccountIdentifier {
	return NewAccountIdentifier(v.OrganizationName, v.AccountName)
}

func (c *accounts) UnsetAllPoliciesSafely(ctx context.Context) error {
	return errors.Join(
		c.UnsetPolicySafely(ctx, PolicyKindAuthenticationPolicy),
		c.UnsetPolicySafely(ctx, PolicyKindFeaturePolicy),
		c.UnsetPolicySafely(ctx, PolicyKindPackagesPolicy),
		c.UnsetPolicySafely(ctx, PolicyKindPasswordPolicy),
		c.UnsetPolicySafely(ctx, PolicyKindSessionPolicy),
	)
}

func (c *accounts) UnsetPolicySafely(ctx context.Context, kind PolicyKind) error {
	var unset *AccountUnset
	switch kind {
	case PolicyKindAuthenticationPolicy:
		unset = &AccountUnset{AuthenticationPolicy: Bool(true)}
	case PolicyKindFeaturePolicy:
		unset = &AccountUnset{FeaturePolicyUnset: &AccountFeaturePolicyUnset{FeaturePolicy: Bool(true)}}
	case PolicyKindPackagesPolicy:
		unset = &AccountUnset{PackagesPolicy: Bool(true)}
	case PolicyKindPasswordPolicy:
		unset = &AccountUnset{PasswordPolicy: Bool(true)}
	case PolicyKindSessionPolicy:
		unset = &AccountUnset{SessionPolicy: Bool(true)}
	default:
		return fmt.Errorf("policy kind %s is not supported for account policies", kind)
	}
	err := c.client.Accounts.Alter(ctx, &AlterAccountOptions{Unset: unset})
	// If the policy is not attached to the account, Snowflake returns an error.
	if errors.Is(err, ErrPolicyNotAttachedToAccount) {
		return nil
	}
	return err
}

func (c *accounts) UnsetAll(ctx context.Context) error {
	return errors.Join(
		c.UnsetAllParameters(ctx),
		c.UnsetAllPoliciesSafely(ctx),
		c.Alter(ctx, &AlterAccountOptions{Unset: &AccountUnset{ResourceMonitor: Bool(true)}}),
	)
}

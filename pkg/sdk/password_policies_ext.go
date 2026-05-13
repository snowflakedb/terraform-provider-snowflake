package sdk

import (
	"context"
	"errors"
	"strconv"
)

func (v *passwordPolicies) DescribeDetails(ctx context.Context, id SchemaObjectIdentifier) (*PasswordPolicyDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parsePasswordPolicyProperties(properties)
}

func parsePasswordPolicyProperties(properties []PasswordPolicyProperty) (*PasswordPolicyDetails, error) {
	details := &PasswordPolicyDetails{}
	var errs []error
	for _, prop := range properties {
		switch prop.Property {
		case "NAME":
			details.Name = prop.Value
		case "OWNER":
			details.Owner = prop.Value
		case "COMMENT":
			details.Comment = emptyIfNull(prop.Value)
		case "PASSWORD_MIN_LENGTH":
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.PasswordMinLength = int(val)
			}
		case "PASSWORD_MAX_LENGTH":
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.PasswordMaxLength = int(val)
			}
		case "PASSWORD_MIN_UPPER_CASE_CHARS":
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.PasswordMinUpperCaseChars = int(val)
			}
		case "PASSWORD_MIN_LOWER_CASE_CHARS":
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.PasswordMinLowerCaseChars = int(val)
			}
		case "PASSWORD_MIN_NUMERIC_CHARS":
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.PasswordMinNumericChars = int(val)
			}
		case "PASSWORD_MIN_SPECIAL_CHARS":
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.PasswordMinSpecialChars = int(val)
			}
		case "PASSWORD_MIN_AGE_DAYS":
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.PasswordMinAgeDays = int(val)
			}
		case "PASSWORD_MAX_AGE_DAYS":
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.PasswordMaxAgeDays = int(val)
			}
		case "PASSWORD_MAX_RETRIES":
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.PasswordMaxRetries = int(val)
			}
		case "PASSWORD_LOCKOUT_TIME_MINS":
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.PasswordLockoutTimeMins = int(val)
			}
		case "PASSWORD_HISTORY":
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.PasswordHistory = int(val)
			}
		}
	}
	return details, errors.Join(errs...)
}

func (r *CreatePasswordPolicyRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

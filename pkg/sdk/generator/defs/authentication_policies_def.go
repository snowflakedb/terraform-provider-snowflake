package defs

import (
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var authenticationPoliciesDef = g.NewInterface(
	"AuthenticationPolicies",
	"AuthenticationPolicy",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperationFromObject(ssotdefs.AuthenticationPolicy).
	AlterOperationFromObject(ssotdefs.AuthenticationPolicy).
	DropOperationFromObject(ssotdefs.AuthenticationPolicy).
	ShowOperationFromObject(ssotdefs.AuthenticationPolicy).
	DescribeOperationFromObject(ssotdefs.AuthenticationPolicy).
	WithEnums(
		ssotdefs.AuthenticationMethodsOption,
		ssotdefs.MfaAuthenticationMethodsOption,
		ssotdefs.MfaAuthenticationMethodsReadOption,
		ssotdefs.MfaEnrollmentOption,
		ssotdefs.MfaEnrollmentReadOption,
		ssotdefs.ClientTypesOption,
		ssotdefs.ClientPolicyDriverType,
		ssotdefs.AllowedProviderOption,
		ssotdefs.NetworkPolicyEvaluationOption,
		ssotdefs.MfaPolicyAllowedMethodsOption,
		ssotdefs.EnforceMfaOnExternalAuthenticationOption,
		ssotdefs.AuthenticationPolicyTargetScope,
	)

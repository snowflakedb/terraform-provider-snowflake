package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (u *LegacyServiceUserModel) WithBinaryInputFormatEnum(binaryInputFormat sdk.BinaryInputFormat) *LegacyServiceUserModel {
	u.BinaryInputFormat = tfconfig.StringVariable(string(binaryInputFormat))
	return u
}

func (u *LegacyServiceUserModel) WithBinaryOutputFormatEnum(binaryOutputFormat sdk.BinaryOutputFormat) *LegacyServiceUserModel {
	u.BinaryOutputFormat = tfconfig.StringVariable(string(binaryOutputFormat))
	return u
}

func (u *LegacyServiceUserModel) WithClientTimestampTypeMappingEnum(clientTimestampTypeMapping sdk.ClientTimestampTypeMapping) *LegacyServiceUserModel {
	u.ClientTimestampTypeMapping = tfconfig.StringVariable(string(clientTimestampTypeMapping))
	return u
}

func (u *LegacyServiceUserModel) WithGeographyOutputFormatEnum(geographyOutputFormat sdk.GeographyOutputFormat) *LegacyServiceUserModel {
	u.GeographyOutputFormat = tfconfig.StringVariable(string(geographyOutputFormat))
	return u
}

func (u *LegacyServiceUserModel) WithGeometryOutputFormatEnum(geometryOutputFormat sdk.GeometryOutputFormat) *LegacyServiceUserModel {
	u.GeometryOutputFormat = tfconfig.StringVariable(string(geometryOutputFormat))
	return u
}

func (u *LegacyServiceUserModel) WithLogLevelEnum(logLevel sdk.LogLevel) *LegacyServiceUserModel {
	u.LogLevel = tfconfig.StringVariable(string(logLevel))
	return u
}

func (u *LegacyServiceUserModel) WithTimestampTypeMappingEnum(timestampTypeMapping sdk.TimestampTypeMapping) *LegacyServiceUserModel {
	u.TimestampTypeMapping = tfconfig.StringVariable(string(timestampTypeMapping))
	return u
}

func (u *LegacyServiceUserModel) WithTraceLevelEnum(traceLevel sdk.TraceLevel) *LegacyServiceUserModel {
	u.TraceLevel = tfconfig.StringVariable(string(traceLevel))
	return u
}

func (u *LegacyServiceUserModel) WithTransactionDefaultIsolationLevelEnum(transactionDefaultIsolationLevel sdk.TransactionDefaultIsolationLevel) *LegacyServiceUserModel {
	u.TransactionDefaultIsolationLevel = tfconfig.StringVariable(string(transactionDefaultIsolationLevel))
	return u
}

func (u *LegacyServiceUserModel) WithUnsupportedDdlActionEnum(unsupportedDdlAction sdk.UnsupportedDDLAction) *LegacyServiceUserModel {
	u.UnsupportedDdlAction = tfconfig.StringVariable(string(unsupportedDdlAction))
	return u
}

func (u *LegacyServiceUserModel) WithNetworkPolicyId(networkPolicy sdk.AccountObjectIdentifier) *LegacyServiceUserModel {
	u.NetworkPolicy = tfconfig.StringVariable(networkPolicy.Name())
	return u
}

func (u *LegacyServiceUserModel) WithDefaultSecondaryRolesOptionEnum(option sdk.SecondaryRolesOption) *LegacyServiceUserModel {
	return u.WithDefaultSecondaryRolesOption(string(option))
}

// WIF (Workload Identity Federation) helper methods

// WithDefaultWorkloadIdentityAws sets the default workload identity to use AWS federation.
func (u *LegacyServiceUserModel) WithDefaultWorkloadIdentityAws(arn string) *LegacyServiceUserModel {
	u.DefaultWorkloadIdentity = tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"aws": tfconfig.ListVariable(tfconfig.ObjectVariable(
				map[string]tfconfig.Variable{
					"arn": tfconfig.StringVariable(arn),
				},
			)),
		},
	)
	return u
}

// WithDefaultWorkloadIdentityGcp sets the default workload identity to use GCP federation.
func (u *LegacyServiceUserModel) WithDefaultWorkloadIdentityGcp(subject string) *LegacyServiceUserModel {
	u.DefaultWorkloadIdentity = tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"gcp": tfconfig.ListVariable(tfconfig.ObjectVariable(
				map[string]tfconfig.Variable{
					"subject": tfconfig.StringVariable(subject),
				},
			)),
		},
	)
	return u
}

// WithDefaultWorkloadIdentityAzure sets the default workload identity to use Azure federation.
func (u *LegacyServiceUserModel) WithDefaultWorkloadIdentityAzure(issuer, subject string) *LegacyServiceUserModel {
	u.DefaultWorkloadIdentity = tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"azure": tfconfig.ListVariable(tfconfig.ObjectVariable(
				map[string]tfconfig.Variable{
					"issuer":  tfconfig.StringVariable(issuer),
					"subject": tfconfig.StringVariable(subject),
				},
			)),
		},
	)
	return u
}

// WithDefaultWorkloadIdentityOidc sets the default workload identity to use generic OIDC federation.
func (u *LegacyServiceUserModel) WithDefaultWorkloadIdentityOidc(issuer, subject string, audienceList ...string) *LegacyServiceUserModel {
	m := map[string]tfconfig.Variable{
		"issuer":  tfconfig.StringVariable(issuer),
		"subject": tfconfig.StringVariable(subject),
	}
	if len(audienceList) > 0 {
		audiences := make([]tfconfig.Variable, len(audienceList))
		for i, a := range audienceList {
			audiences[i] = tfconfig.StringVariable(a)
		}
		m["oidc_audience_list"] = tfconfig.ListVariable(audiences...)
	}

	u.DefaultWorkloadIdentity = tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"oidc": tfconfig.ListVariable(tfconfig.ObjectVariable(m)),
		},
	)
	return u
}

// WithoutDefaultWorkloadIdentity removes any default workload identity configuration.
func (u *LegacyServiceUserModel) WithoutDefaultWorkloadIdentity() *LegacyServiceUserModel {
	u.DefaultWorkloadIdentity = nil
	return u
}

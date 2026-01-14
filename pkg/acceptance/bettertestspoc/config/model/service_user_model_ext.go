package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (u *ServiceUserModel) WithBinaryInputFormatEnum(binaryInputFormat sdk.BinaryInputFormat) *ServiceUserModel {
	u.BinaryInputFormat = tfconfig.StringVariable(string(binaryInputFormat))
	return u
}

func (u *ServiceUserModel) WithBinaryOutputFormatEnum(binaryOutputFormat sdk.BinaryOutputFormat) *ServiceUserModel {
	u.BinaryOutputFormat = tfconfig.StringVariable(string(binaryOutputFormat))
	return u
}

func (u *ServiceUserModel) WithClientTimestampTypeMappingEnum(clientTimestampTypeMapping sdk.ClientTimestampTypeMapping) *ServiceUserModel {
	u.ClientTimestampTypeMapping = tfconfig.StringVariable(string(clientTimestampTypeMapping))
	return u
}

func (u *ServiceUserModel) WithGeographyOutputFormatEnum(geographyOutputFormat sdk.GeographyOutputFormat) *ServiceUserModel {
	u.GeographyOutputFormat = tfconfig.StringVariable(string(geographyOutputFormat))
	return u
}

func (u *ServiceUserModel) WithGeometryOutputFormatEnum(geometryOutputFormat sdk.GeometryOutputFormat) *ServiceUserModel {
	u.GeometryOutputFormat = tfconfig.StringVariable(string(geometryOutputFormat))
	return u
}

func (u *ServiceUserModel) WithLogLevelEnum(logLevel sdk.LogLevel) *ServiceUserModel {
	u.LogLevel = tfconfig.StringVariable(string(logLevel))
	return u
}

func (u *ServiceUserModel) WithTimestampTypeMappingEnum(timestampTypeMapping sdk.TimestampTypeMapping) *ServiceUserModel {
	u.TimestampTypeMapping = tfconfig.StringVariable(string(timestampTypeMapping))
	return u
}

func (u *ServiceUserModel) WithTraceLevelEnum(traceLevel sdk.TraceLevel) *ServiceUserModel {
	u.TraceLevel = tfconfig.StringVariable(string(traceLevel))
	return u
}

func (u *ServiceUserModel) WithTransactionDefaultIsolationLevelEnum(transactionDefaultIsolationLevel sdk.TransactionDefaultIsolationLevel) *ServiceUserModel {
	u.TransactionDefaultIsolationLevel = tfconfig.StringVariable(string(transactionDefaultIsolationLevel))
	return u
}

func (u *ServiceUserModel) WithUnsupportedDdlActionEnum(unsupportedDdlAction sdk.UnsupportedDDLAction) *ServiceUserModel {
	u.UnsupportedDdlAction = tfconfig.StringVariable(string(unsupportedDdlAction))
	return u
}

func (u *ServiceUserModel) WithNetworkPolicyId(networkPolicy sdk.AccountObjectIdentifier) *ServiceUserModel {
	u.NetworkPolicy = tfconfig.StringVariable(networkPolicy.Name())
	return u
}

func (u *ServiceUserModel) WithDefaultSecondaryRolesOptionEnum(option sdk.SecondaryRolesOption) *ServiceUserModel {
	return u.WithDefaultSecondaryRolesOption(string(option))
}

// WIF (Workload Identity Federation) helper methods

// WithDefaultWorkloadIdentityAws sets the default workload identity to use AWS federation.
func (u *ServiceUserModel) WithDefaultWorkloadIdentityAws(arn string) *ServiceUserModel {
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
func (u *ServiceUserModel) WithDefaultWorkloadIdentityGcp(subject string) *ServiceUserModel {
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
func (u *ServiceUserModel) WithDefaultWorkloadIdentityAzure(issuer, subject string) *ServiceUserModel {
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
func (u *ServiceUserModel) WithDefaultWorkloadIdentityOidc(issuer, subject string, audienceList []string) *ServiceUserModel {
	audiences := make([]tfconfig.Variable, len(audienceList))
	for i, a := range audienceList {
		audiences[i] = tfconfig.StringVariable(a)
	}

	u.DefaultWorkloadIdentity = tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"oidc": tfconfig.ListVariable(tfconfig.ObjectVariable(
				map[string]tfconfig.Variable{
					"issuer":             tfconfig.StringVariable(issuer),
					"subject":            tfconfig.StringVariable(subject),
					"oidc_audience_list": tfconfig.ListVariable(audiences...),
				},
			)),
		},
	)
	return u
}

// WithoutDefaultWorkloadIdentity removes any default workload identity configuration.
func (u *ServiceUserModel) WithoutDefaultWorkloadIdentity() *ServiceUserModel {
	u.DefaultWorkloadIdentity = nil
	return u
}

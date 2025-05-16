package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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

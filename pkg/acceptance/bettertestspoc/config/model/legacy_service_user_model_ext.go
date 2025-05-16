package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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

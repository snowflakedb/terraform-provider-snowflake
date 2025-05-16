package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (u *UserModel) WithBinaryInputFormatEnum(binaryInputFormat sdk.BinaryInputFormat) *UserModel {
	u.BinaryInputFormat = tfconfig.StringVariable(string(binaryInputFormat))
	return u
}

func (u *UserModel) WithBinaryOutputFormatEnum(binaryOutputFormat sdk.BinaryOutputFormat) *UserModel {
	u.BinaryOutputFormat = tfconfig.StringVariable(string(binaryOutputFormat))
	return u
}

func (u *UserModel) WithClientTimestampTypeMappingEnum(clientTimestampTypeMapping sdk.ClientTimestampTypeMapping) *UserModel {
	u.ClientTimestampTypeMapping = tfconfig.StringVariable(string(clientTimestampTypeMapping))
	return u
}

func (u *UserModel) WithGeographyOutputFormatEnum(geographyOutputFormat sdk.GeographyOutputFormat) *UserModel {
	u.GeographyOutputFormat = tfconfig.StringVariable(string(geographyOutputFormat))
	return u
}

func (u *UserModel) WithGeometryOutputFormatEnum(geometryOutputFormat sdk.GeometryOutputFormat) *UserModel {
	u.GeometryOutputFormat = tfconfig.StringVariable(string(geometryOutputFormat))
	return u
}

func (u *UserModel) WithLogLevelEnum(logLevel sdk.LogLevel) *UserModel {
	u.LogLevel = tfconfig.StringVariable(string(logLevel))
	return u
}

func (u *UserModel) WithTimestampTypeMappingEnum(timestampTypeMapping sdk.TimestampTypeMapping) *UserModel {
	u.TimestampTypeMapping = tfconfig.StringVariable(string(timestampTypeMapping))
	return u
}

func (u *UserModel) WithTraceLevelEnum(traceLevel sdk.TraceLevel) *UserModel {
	u.TraceLevel = tfconfig.StringVariable(string(traceLevel))
	return u
}

func (u *UserModel) WithTransactionDefaultIsolationLevelEnum(transactionDefaultIsolationLevel sdk.TransactionDefaultIsolationLevel) *UserModel {
	u.TransactionDefaultIsolationLevel = tfconfig.StringVariable(string(transactionDefaultIsolationLevel))
	return u
}

func (u *UserModel) WithUnsupportedDdlActionEnum(unsupportedDdlAction sdk.UnsupportedDDLAction) *UserModel {
	u.UnsupportedDdlAction = tfconfig.StringVariable(string(unsupportedDdlAction))
	return u
}

func (u *UserModel) WithNetworkPolicyId(networkPolicy sdk.AccountObjectIdentifier) *UserModel {
	u.NetworkPolicy = tfconfig.StringVariable(networkPolicy.Name())
	return u
}

func (u *UserModel) WithNullPassword() *UserModel {
	return u.WithPasswordValue(config.ReplacementPlaceholderVariable(config.SnowflakeProviderConfigNull))
}

func (u *UserModel) WithDefaultSecondaryRolesOptionEnum(option sdk.SecondaryRolesOption) *UserModel {
	return u.WithDefaultSecondaryRolesOption(string(option))
}

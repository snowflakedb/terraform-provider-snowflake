package helpers

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SnowflakeDefaultsClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewSnowflakeDefaultsClient(context *TestClientContext) *SnowflakeDefaultsClient {
	return &SnowflakeDefaultsClient{
		context: context,
	}
}

func (c *SnowflakeDefaultsClient) WarehouseGenerationEmptyByDefault(t *testing.T) bool {
	t.Helper()
	if c.context.snowflakeEnvironment == testenvs.SnowflakePreProdGovEnvironment {
		return true
	}
	return false
}

// DefaultQuotedIdentifiersIgnoreCaseLevel returns the expected parameter level for
// QUOTED_IDENTIFIERS_IGNORE_CASE on a newly created object. Prod accounts surface the
// Snowflake default (empty level); non-prod accounts have it set at ACCOUNT level.
func (c *SnowflakeDefaultsClient) DefaultQuotedIdentifiersIgnoreCaseLevel(t *testing.T) sdk.ParameterType {
	t.Helper()
	if c.context.snowflakeEnvironment != testenvs.SnowflakeProdEnvironment {
		return sdk.ParameterTypeAccount
	}
	return sdk.ParameterTypeSnowflakeDefault
}

// DefaultAutoSuspend returns the Snowflake default for AUTO_SUSPEND.
// Prod accounts default to 600; preprod still defaults to 34.
func (c *SnowflakeDefaultsClient) DefaultAutoSuspend(t *testing.T) int {
	t.Helper()
	if c.context.snowflakeEnvironment != testenvs.SnowflakeProdEnvironment {
		return 34
	}
	return 600
}

// DefaultQueryAccelerationMaxScaleFactor returns the Snowflake default for
// QUERY_ACCELERATION_MAX_SCALE_FACTOR. Prod accounts default to 2; preprod still defaults to 8.
func (c *SnowflakeDefaultsClient) DefaultQueryAccelerationMaxScaleFactor(t *testing.T) int {
	t.Helper()
	if c.context.snowflakeEnvironment != testenvs.SnowflakeProdEnvironment {
		return 8
	}
	return 2
}

// DefaultStatementTimeoutInSecondsLevel returns the expected parameter level for
// STATEMENT_TIMEOUT_IN_SECONDS on a newly created warehouse. Prod inherits the
// Snowflake default (empty level); non-prod accounts currently surface ACCOUNT because they need to be overridden.
func (c *SnowflakeDefaultsClient) DefaultStatementTimeoutInSecondsLevel(t *testing.T) sdk.ParameterType {
	t.Helper()
	if c.context.snowflakeEnvironment != testenvs.SnowflakeProdEnvironment {
		return sdk.ParameterTypeAccount
	}
	return sdk.ParameterTypeSnowflakeDefault
}

// DefaultMfaEnrollment returns the expected mfa_enrollment value for a newly created
// authentication policy without an explicit mfa_enrollment setting.
// Prod accounts still use the legacy OPTIONAL default; non-prod accounts have
// already moved to the REQUIRED default (Snowflake is deprecating OPTIONAL).
func (c *SnowflakeDefaultsClient) DefaultMfaEnrollment(t *testing.T) string {
	t.Helper()
	if c.context.snowflakeEnvironment != testenvs.SnowflakeProdEnvironment {
		return string(sdk.MfaEnrollmentReadOptionRequired)
	}
	return string(sdk.MfaEnrollmentOptionOptional)
}

// DefaultEnableCortexAnalystLevel returns the expected parameter level for
// ENABLE_CORTEX_ANALYST after unset. Prod accounts surface SYSTEM; non-prod
// accounts inherit the Snowflake default (empty level).
func (c *SnowflakeDefaultsClient) DefaultEnableCortexAnalystLevel(t *testing.T) sdk.ParameterType {
	t.Helper()
	if c.context.snowflakeEnvironment == testenvs.SnowflakeProdEnvironment {
		return sdk.ParameterTypeSystem
	}
	return sdk.ParameterTypeSnowflakeDefault
}

// DefaultEnableNotebookCreationInPersonalDbLevel returns the expected parameter
// level for ENABLE_NOTEBOOK_CREATION_IN_PERSONAL_DB after unset. Prod accounts
// surface SYSTEM; non-prod accounts inherit the Snowflake default (empty level).
func (c *SnowflakeDefaultsClient) DefaultEnableNotebookCreationInPersonalDbLevel(t *testing.T) sdk.ParameterType {
	t.Helper()
	if c.context.snowflakeEnvironment == testenvs.SnowflakeProdEnvironment {
		return sdk.ParameterTypeSystem
	}
	return sdk.ParameterTypeSnowflakeDefault
}

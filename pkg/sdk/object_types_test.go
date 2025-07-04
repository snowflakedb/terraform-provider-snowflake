package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToObjectType(t *testing.T) {
	type test struct {
		input string
		want  ObjectType
	}

	valid := []test{
		// Case insensitive.
		{input: "schema", want: ObjectTypeSchema},

		// Supported Values.
		{input: "ACCOUNT", want: ObjectTypeAccount},
		{input: "MANAGED ACCOUNT", want: ObjectTypeManagedAccount},
		{input: "USER", want: ObjectTypeUser},
		{input: "DATABASE ROLE", want: ObjectTypeDatabaseRole},
		{input: "DATASET", want: ObjectTypeDataset},
		{input: "ROLE", want: ObjectTypeRole},
		{input: "INTEGRATION", want: ObjectTypeIntegration},
		{input: "NETWORK POLICY", want: ObjectTypeNetworkPolicy},
		{input: "PASSWORD POLICY", want: ObjectTypePasswordPolicy},
		{input: "SESSION POLICY", want: ObjectTypeSessionPolicy},
		{input: "PRIVACY POLICY", want: ObjectTypePrivacyPolicy},
		{input: "REPLICATION GROUP", want: ObjectTypeReplicationGroup},
		{input: "FAILOVER GROUP", want: ObjectTypeFailoverGroup},
		{input: "CONNECTION", want: ObjectTypeConnection},
		{input: "PARAMETER", want: ObjectTypeParameter},
		{input: "WAREHOUSE", want: ObjectTypeWarehouse},
		{input: "RESOURCE MONITOR", want: ObjectTypeResourceMonitor},
		{input: "DATABASE", want: ObjectTypeDatabase},
		{input: "SCHEMA", want: ObjectTypeSchema},
		{input: "SHARE", want: ObjectTypeShare},
		{input: "TABLE", want: ObjectTypeTable},
		{input: "DYNAMIC TABLE", want: ObjectTypeDynamicTable},
		{input: "CORTEX SEARCH SERVICE", want: ObjectTypeCortexSearchService},
		{input: "EXTERNAL TABLE", want: ObjectTypeExternalTable},
		{input: "EVENT TABLE", want: ObjectTypeEventTable},
		{input: "VIEW", want: ObjectTypeView},
		{input: "MATERIALIZED VIEW", want: ObjectTypeMaterializedView},
		{input: "SEQUENCE", want: ObjectTypeSequence},
		{input: "SNAPSHOT", want: ObjectTypeSnapshot},
		{input: "FUNCTION", want: ObjectTypeFunction},
		{input: "EXTERNAL FUNCTION", want: ObjectTypeExternalFunction},
		{input: "PROCEDURE", want: ObjectTypeProcedure},
		{input: "STREAM", want: ObjectTypeStream},
		{input: "TASK", want: ObjectTypeTask},
		{input: "MASKING POLICY", want: ObjectTypeMaskingPolicy},
		{input: "ROW ACCESS POLICY", want: ObjectTypeRowAccessPolicy},
		{input: "TAG", want: ObjectTypeTag},
		{input: "SECRET", want: ObjectTypeSecret},
		{input: "STAGE", want: ObjectTypeStage},
		{input: "FILE FORMAT", want: ObjectTypeFileFormat},
		{input: "PIPE", want: ObjectTypePipe},
		{input: "ALERT", want: ObjectTypeAlert},
		{input: "SNOWFLAKE.CORE.BUDGET", want: ObjectTypeBudget},
		{input: "SNOWFLAKE.ML.CLASSIFICATION", want: ObjectTypeClassification},
		{input: "APPLICATION", want: ObjectTypeApplication},
		{input: "APPLICATION PACKAGE", want: ObjectTypeApplicationPackage},
		{input: "APPLICATION ROLE", want: ObjectTypeApplicationRole},
		{input: "STREAMLIT", want: ObjectTypeStreamlit},
		{input: "COLUMN", want: ObjectTypeColumn},
		{input: "ICEBERG TABLE", want: ObjectTypeIcebergTable},
		{input: "EXTERNAL VOLUME", want: ObjectTypeExternalVolume},
		{input: "NETWORK RULE", want: ObjectTypeNetworkRule},
		{input: "NOTEBOOK", want: ObjectTypeNotebook},
		{input: "PACKAGES POLICY", want: ObjectTypePackagesPolicy},
		{input: "COMPUTE POOL", want: ObjectTypeComputePool},
		{input: "AGGREGATION POLICY", want: ObjectTypeAggregationPolicy},
		{input: "AUTHENTICATION POLICY", want: ObjectTypeAuthenticationPolicy},
		{input: "HYBRID TABLE", want: ObjectTypeHybridTable},
		{input: "IMAGE REPOSITORY", want: ObjectTypeImageRepository},
		{input: "PROJECTION POLICY", want: ObjectTypeProjectionPolicy},
		{input: "DATA METRIC FUNCTION", want: ObjectTypeDataMetricFunction},
		{input: "GIT REPOSITORY", want: ObjectTypeGitRepository},
		{input: "MODEL", want: ObjectTypeModel},
		{input: "SERVICE", want: ObjectTypeService},
		{input: "PROGRAMMATIC ACCESS TOKEN", want: ObjectTypeProgrammaticAccessToken},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToObjectType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToObjectType(tc.input)
			require.Error(t, err)
		})
	}
}

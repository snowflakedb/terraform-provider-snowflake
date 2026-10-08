package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type enumParameterMetadata struct {
	validate           schema.SchemaValidateDiagFunc
	diffSuppress       schema.SchemaDiffSuppressFunc
	optionsDescription string
}

// enumParameterValidators maps a catalog Kind naming an enum type to its derived Terraform behaviors.
var enumParameterValidators = map[string]enumParameterMetadata{
	"ActivePythonProfiler":                   {sdkValidation(sdk.ToActivePythonProfiler), NormalizeAndCompare(sdk.ToActivePythonProfiler), enumValuesDescription(sdk.AllActivePythonProfilers)},
	"BinaryInputFormat":                      {sdkValidation(sdk.ToBinaryInputFormat), NormalizeAndCompare(sdk.ToBinaryInputFormat), enumValuesDescription(sdk.AllBinaryInputFormats)},
	"BinaryOutputFormat":                     {sdkValidation(sdk.ToBinaryOutputFormat), NormalizeAndCompare(sdk.ToBinaryOutputFormat), enumValuesDescription(sdk.AllBinaryOutputFormats)},
	"ClientTimestampTypeMapping":             {sdkValidation(sdk.ToClientTimestampTypeMapping), NormalizeAndCompare(sdk.ToClientTimestampTypeMapping), enumValuesDescription(sdk.AllClientTimestampTypeMappings)},
	"DefaultNullOrdering":                    {sdkValidation(sdk.ToDefaultNullOrdering), NormalizeAndCompare(sdk.ToDefaultNullOrdering), enumValuesDescription(sdk.AllDefaultNullOrderings)},
	"GeographyOutputFormat":                  {sdkValidation(sdk.ToGeographyOutputFormat), NormalizeAndCompare(sdk.ToGeographyOutputFormat), enumValuesDescription(sdk.AllGeographyOutputFormats)},
	"GeometryOutputFormat":                   {sdkValidation(sdk.ToGeometryOutputFormat), NormalizeAndCompare(sdk.ToGeometryOutputFormat), enumValuesDescription(sdk.AllGeometryOutputFormats)},
	"IcebergTableIcebergMergeOnReadBehavior": {sdkValidation(sdk.ToIcebergTableIcebergMergeOnReadBehavior), NormalizeAndCompare(sdk.ToIcebergTableIcebergMergeOnReadBehavior), enumValuesDescription(sdk.AllIcebergTableIcebergMergeOnReadBehaviors)},
	"IcebergTableTargetFileSize":             {sdkValidation(sdk.ToIcebergTableTargetFileSize), NormalizeAndCompare(sdk.ToIcebergTableTargetFileSize), enumValuesDescription(sdk.AllIcebergTableTargetFileSizes)},
	"LogLevel":                               {sdkValidation(sdk.ToLogLevel), NormalizeAndCompare(sdk.ToLogLevel), enumValuesDescription(sdk.AllLogLevels)},
	"MetricLevel":                            {sdkValidation(sdk.ToMetricLevel), NormalizeAndCompare(sdk.ToMetricLevel), enumValuesDescription(sdk.AllMetricLevels)},
	"StorageSerializationPolicy":             {sdkValidation(sdk.ToStorageSerializationPolicy), NormalizeAndCompare(sdk.ToStorageSerializationPolicy), enumValuesDescription(sdk.AllStorageSerializationPolicies)},
	"TimestampTypeMapping":                   {sdkValidation(sdk.ToTimestampTypeMapping), NormalizeAndCompare(sdk.ToTimestampTypeMapping), enumValuesDescription(sdk.AllTimestampTypeMappings)},
	"TraceLevel":                             {sdkValidation(sdk.ToTraceLevel), NormalizeAndCompare(sdk.ToTraceLevel), enumValuesDescription(sdk.AllTraceLevels)},
	"TransactionDefaultIsolationLevel":       {sdkValidation(sdk.ToTransactionDefaultIsolationLevel), NormalizeAndCompare(sdk.ToTransactionDefaultIsolationLevel), enumValuesDescription(sdk.AllTransactionDefaultIsolationLevels)},
	"UnsupportedDDLAction":                   {sdkValidation(sdk.ToUnsupportedDDLAction), NormalizeAndCompare(sdk.ToUnsupportedDDLAction), enumValuesDescription(sdk.AllUnsupportedDDLActions)},
	"WarehouseSize":                          {sdkValidation(sdk.ToWarehouseSize), NormalizeAndCompare(sdk.ToWarehouseSize), enumValuesDescription(sdk.AllWarehouseSizes)},
}

// identifierParameterValidators maps a catalog Kind naming an identifier type to its validator.
var identifierParameterValidators = map[string]schema.SchemaValidateDiagFunc{
	"AccountObjectIdentifier": IsValidIdentifier[sdk.AccountObjectIdentifier](),
	"SchemaObjectIdentifier":  IsValidIdentifier[sdk.SchemaObjectIdentifier](),
}

// KindFloat ("float64") stays a Terraform string. The SDK field is float64 so
// ALTER renders an unquoted decimal literal, while the attribute keeps scale ("10.0").
var primitiveParameterValueTypes = map[string]schema.ValueType{
	"bool":             schema.TypeBool,
	"int":              schema.TypeInt,
	"float64":          schema.TypeString,
	"string":           schema.TypeString,
	"StringAllowEmpty": schema.TypeString,
}

type parameterSchemaModifier func(*schema.Schema)

func withIntAtLeast(minimum int) parameterSchemaModifier {
	return func(s *schema.Schema) {
		s.ValidateDiagFunc = validation.ToDiagFunc(validation.IntAtLeast(minimum))
	}
}

func withIntBetween(minimum, maximum int) parameterSchemaModifier {
	return func(s *schema.Schema) {
		s.ValidateDiagFunc = validation.ToDiagFunc(validation.IntBetween(minimum, maximum))
	}
}

func withFloatValidation() parameterSchemaModifier {
	return func(s *schema.Schema) {
		s.ValidateDiagFunc = sdkValidation(sdk.ToFloat64)
		s.DiffSuppressFunc = NormalizeAndCompare(sdk.ToFloat64)
	}
}

func withAccountObjectIdentifierValidation() parameterSchemaModifier {
	return func(s *schema.Schema) {
		s.ValidateDiagFunc = IsValidIdentifier[sdk.AccountObjectIdentifier]()
		s.DiffSuppressFunc = suppressIdentifierQuoting
	}
}

var parameterSchemaModifiers = buildParameterSchemaModifiers()

func buildParameterSchemaModifiers() map[string][]parameterSchemaModifier {
	result := make(map[string][]parameterSchemaModifier)
	register := func(modifier parameterSchemaModifier, parameters ...parameterdefs.ParameterDef) {
		for _, parameter := range parameters {
			result[parameter.SqlName] = append(result[parameter.SqlName], modifier)
		}
	}

	register(withIntAtLeast(0),
		defs.DataRetentionTimeInDays,
		defs.HybridTableLockTimeout,
		defs.IcebergVersionDefault,
		defs.JsonIndent,
		defs.LockTimeout,
		defs.MaxDataExtensionTimeInDays,
		defs.MinDataRetentionTimeInDays,
		defs.MultiStatementCount,
		defs.RowsPerResultset,
		defs.StatementQueuedTimeoutInSeconds,
		defs.StatementTimeoutInSeconds,
		defs.SuspendTaskAfterNumFailures,
		defs.TaskAutoRetryAttempts,
		defs.UserTaskMinimumTriggerIntervalInSeconds,
		defs.UserTaskTimeoutMs,
	)
	register(withIntAtLeast(1),
		defs.ClientMemoryLimit,
		defs.ClientPrefetchThreads,
		defs.MaxConcurrencyLevel,
	)
	register(withIntAtLeast(16), defs.ClientResultChunkSize)
	register(withIntAtLeast(900), defs.ClientSessionKeepAliveHeartbeatFrequency)
	register(withIntAtLeast(1900), defs.TwoDigitCenturyStart)
	register(withFloatValidation(), defs.InitialReplicationSizeLimitInTb)
	register(withAccountObjectIdentifierValidation(),
		defs.CatalogSync,
		defs.DefaultNotebookComputePoolCpu,
		defs.DefaultNotebookComputePoolGpu,
		defs.DefaultStreamlitComputePool,
	)

	return result
}

// parameterSchema derives a resource schema entry from a catalog parameter.
func parameterSchema(p parameterdefs.ParameterDef) *schema.Schema {
	return parameterSchemaWithModifiers(p, true)
}

func parameterSchemaWithoutModifiers(p parameterdefs.ParameterDef) *schema.Schema {
	return parameterSchemaWithModifiers(p, false)
}

func parameterSchemaWithModifiers(p parameterdefs.ParameterDef, applyModifiers bool) *schema.Schema {
	valueType, isPrimitive := primitiveParameterValueTypes[p.Kind]
	enumMetadata, isEnum := enumParameterValidators[p.Kind]
	identifierValidator, isIdentifier := identifierParameterValidators[p.Kind]

	if !isPrimitive && !isEnum && !isIdentifier {
		panic(fmt.Sprintf("parameter %s has unhandled kind %q; register it as a primitive, enum, or identifier kind", p.SqlName, p.Kind))
	}
	if !isPrimitive {
		valueType = schema.TypeString
	}

	s := &schema.Schema{
		Type:        valueType,
		Description: enrichWithReferenceToParameterDocs(p.SqlName, p.Description),
		Optional:    true,
		Computed:    true,
	}

	switch {
	case isEnum:
		s.ValidateDiagFunc = enumMetadata.validate
		s.DiffSuppressFunc = enumMetadata.diffSuppress
		s.Description = enrichWithReferenceToParameterDocs(p.SqlName, joinWithSpace(p.Description, enumMetadata.optionsDescription))
	case isIdentifier:
		s.ValidateDiagFunc = identifierValidator
		s.DiffSuppressFunc = suppressIdentifierQuoting
	}
	if applyModifiers {
		for _, modify := range parameterSchemaModifiers[p.SqlName] {
			modify(s)
		}
	}

	return s
}

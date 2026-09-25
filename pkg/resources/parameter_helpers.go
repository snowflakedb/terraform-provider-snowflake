package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type enumParameterMetadata struct {
	validate           schema.SchemaValidateDiagFunc
	diffSuppress       schema.SchemaDiffSuppressFunc
	optionsDescription string
}

// enumParameterValidators maps a catalog Kind naming an enum type to its derived Terraform behaviors.
var enumParameterValidators = map[string]enumParameterMetadata{
	"LogLevel":                   {sdkValidation(sdk.ToLogLevel), NormalizeAndCompare(sdk.ToLogLevel), enumValuesDescription(sdk.AllLogLevels)},
	"TraceLevel":                 {sdkValidation(sdk.ToTraceLevel), NormalizeAndCompare(sdk.ToTraceLevel), enumValuesDescription(sdk.AllTraceLevels)},
	"StorageSerializationPolicy": {sdkValidation(sdk.ToStorageSerializationPolicy), NormalizeAndCompare(sdk.ToStorageSerializationPolicy), enumValuesDescription(sdk.AllStorageSerializationPolicies)},
	"WarehouseSize":              {sdkValidation(sdk.ToWarehouseSize), NormalizeAndCompare(sdk.ToWarehouseSize), enumValuesDescription(sdk.AllWarehouseSizes)},
}

// identifierParameterValidators maps a catalog Kind naming an identifier type to its validator.
var identifierParameterValidators = map[string]schema.SchemaValidateDiagFunc{
	"AccountObjectIdentifier": IsValidIdentifier[sdk.AccountObjectIdentifier](),
}

var primitiveParameterValueTypes = map[string]schema.ValueType{
	"bool":             schema.TypeBool,
	"int":              schema.TypeInt,
	"string":           schema.TypeString,
	"StringAllowEmpty": schema.TypeString,
}

// parameterSchema derives a resource schema entry from a catalog parameter.
func parameterSchema(p parameterdefs.ParameterDef) *schema.Schema {
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

	return s
}

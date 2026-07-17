package gen

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type SnowflakeObjectAssertionsModel struct {
	Name                 string
	SdkType              string
	IdType               string
	IsDataSourceOutput   bool
	IsSubStruct          bool
	ObjectTypeName       string
	NoShowById           bool
	NoIdentifiableObject bool
	ShowByParentId       *genhelpers.ShowByParentIdDef
	DescribeOverride     *genhelpers.DescribeOverrideDef
	NestedAssertFields   []string
	Fields               []SnowflakeObjectFieldAssertion

	*genhelpers.PreambleModel
}

// ComparisonMethod is the strategy used to compare the expected and actual field values in a Has* assertion.
type ComparisonMethod string

const (
	// ComparisonMethodIdentity uses != for primitives, strings, enums, time.Time, and identifier types mapped to strings.
	ComparisonMethodIdentity ComparisonMethod = "identity"
	// ComparisonMethodDeepEqual uses reflect.DeepEqual for struct types that have no special mapper.
	ComparisonMethodDeepEqual ComparisonMethod = "deepEqual"
	// ComparisonMethodDataTypeCompare uses datatypes.AreTheSame for datatypes.DataType interface fields.
	ComparisonMethodDataTypeCompare ComparisonMethod = "dataTypeCompare"
)

type SnowflakeObjectFieldAssertion struct {
	Name                     string
	ConcreteType             string
	IsOriginalTypePointer    bool
	IsOriginalTypeSlice      bool
	Mapper                   genhelpers.Mapper
	ExpectedValueMapper      genhelpers.Mapper
	ErrorMapper              genhelpers.Mapper // mapper used in fmt.Errorf; may differ from Mapper (e.g. .ToSql() for datatypes.DataType)
	ErrorExpectedValueMapper genhelpers.Mapper
	ComparisonMethod         ComparisonMethod
	GenerateNestedAdapter    bool // true for NestedAssertFields — generates nested assertion adapter instead of Has*
	SkipGeneration           bool // true for SkipFields — suppresses Has* generation entirely
}

// ComparisonFunc returns the function to use in a `!ComparisonFunc(actual, expected)` comparison,
// or an empty string when identity (`!=`) should be used instead.
func (f SnowflakeObjectFieldAssertion) ComparisonFunc() string {
	switch f.ComparisonMethod {
	case ComparisonMethodDeepEqual:
		return "reflect.DeepEqual"
	case ComparisonMethodDataTypeCompare:
		return "datatypes.AreTheSame"
	default:
		return ""
	}
}

func (m SnowflakeObjectAssertionsModel) PlaceholderIdentifier() string {
	switch m.IdType {
	case "sdk.AccountObjectIdentifier":
		return `sdk.NewAccountObjectIdentifier("")`
	case "sdk.SchemaObjectIdentifier":
		return `sdk.NewSchemaObjectIdentifier("", "", "")`
	case "sdk.DatabaseObjectIdentifier":
		return `sdk.NewDatabaseObjectIdentifier("", "")`
	default:
		return `sdk.NewAccountObjectIdentifier("")`
	}
}

func ModelFromSdkObjectDetails(sdkObject genhelpers.SdkObjectDetails, preamble *genhelpers.PreambleModel) SnowflakeObjectAssertionsModel {
	name, _ := strings.CutPrefix(sdkObject.Name, "sdk.")
	fields := make([]SnowflakeObjectFieldAssertion, len(sdkObject.Fields))
	for idx, field := range sdkObject.Fields {
		fieldAssertion := MapToSnowflakeObjectFieldAssertion(field)
		if slices.Contains(sdkObject.NestedAssertFields, field.Name) {
			fieldAssertion.GenerateNestedAdapter = true
		}
		if slices.Contains(sdkObject.SkipFields, field.Name) {
			fieldAssertion.SkipGeneration = true
		}
		fields[idx] = fieldAssertion
	}

	objectTypeName := name
	if sdkObject.ObjectTypeName != "" {
		objectTypeName = sdkObject.ObjectTypeName
	}

	return SnowflakeObjectAssertionsModel{
		Name:                 name,
		SdkType:              sdkObject.Name,
		IdType:               sdkObject.IdType,
		IsDataSourceOutput:   sdkObject.IsDataSourceOutput,
		IsSubStruct:          sdkObject.IsSubStruct,
		ObjectTypeName:       objectTypeName,
		NoShowById:           sdkObject.NoShowById,
		NoIdentifiableObject: sdkObject.NoIdentifiableObject,
		ShowByParentId:       sdkObject.ShowByParentId,
		DescribeOverride:     sdkObject.DescribeOverride,
		Fields:               fields,
		NestedAssertFields:   sdkObject.NestedAssertFields,
		PreambleModel:        preamble,
	}
}

func MapToSnowflakeObjectFieldAssertion(field genhelpers.Field) SnowflakeObjectFieldAssertion {
	concreteBase := field.ConcreteTypeNoPointerNoArray()
	underlyingKind := strings.TrimPrefix(field.UnderlyingType, "*")

	mapper := genhelpers.Identity
	if field.IsPointer() {
		mapper = genhelpers.Dereference
	}
	expectedValueMapper := genhelpers.Identity

	comparisonMethod := ComparisonMethodIdentity
	errorMapper := mapper
	errorExpectedValueMapper := expectedValueMapper
	switch {
	case concreteBase == "datatypes.DataType":
		comparisonMethod = ComparisonMethodDataTypeCompare
		errorMapper = genhelpers.ToSql
		errorExpectedValueMapper = genhelpers.ToSql
	case concreteBase == "sdk.AccountObjectIdentifier", concreteBase == "sdk.SchemaObjectIdentifier":
		comparisonMethod = ComparisonMethodIdentity
	case underlyingKind == "struct" && concreteBase != "time.Time":
		comparisonMethod = ComparisonMethodDeepEqual
	}

	// TODO [SNOW-1501905]: handle other mappings if needed
	if concreteBase == "sdk.AccountObjectIdentifier" {
		mapper = genhelpers.Name
		if field.IsPointer() {
			mapper = func(s string) string {
				return genhelpers.Name(genhelpers.Parentheses(genhelpers.Dereference(s)))
			}
		}
		expectedValueMapper = genhelpers.Name
		errorMapper = genhelpers.Name
		errorExpectedValueMapper = genhelpers.Name
	}
	if concreteBase == "sdk.SchemaObjectIdentifier" {
		mapper = genhelpers.FullyQualifiedName
		if field.IsPointer() {
			mapper = func(s string) string {
				return genhelpers.FullyQualifiedName(genhelpers.Parentheses(genhelpers.Dereference(s)))
			}
		}
		expectedValueMapper = genhelpers.FullyQualifiedName
		errorMapper = genhelpers.FullyQualifiedName
		errorExpectedValueMapper = genhelpers.FullyQualifiedName
	}

	return SnowflakeObjectFieldAssertion{
		Name:                     field.Name,
		ConcreteType:             field.ConcreteType,
		IsOriginalTypePointer:    field.IsPointer(),
		IsOriginalTypeSlice:      field.IsSlice(),
		Mapper:                   mapper,
		ExpectedValueMapper:      expectedValueMapper,
		ErrorMapper:              errorMapper,
		ErrorExpectedValueMapper: errorExpectedValueMapper,
		ComparisonMethod:         comparisonMethod,
	}
}

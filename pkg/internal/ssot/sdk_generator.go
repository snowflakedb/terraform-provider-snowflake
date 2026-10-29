package ssot

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

const (
	showDocURLPrefix     = "https://docs.snowflake.com/en/sql-reference/sql/show-"
	describeDocURLPrefix = "https://docs.snowflake.com/en/sql-reference/sql/desc-"
	dropDocURLPrefix     = "https://docs.snowflake.com/en/sql-reference/sql/drop-"
	createDocURLPrefix   = "https://docs.snowflake.com/en/sql-reference/sql/create-"
	alterDocURLPrefix    = "https://docs.snowflake.com/en/sql-reference/sql/alter-"
	// customOperationDocURLPrefix has no verb. The kind slug is the verb.
	customOperationDocURLPrefix = "https://docs.snowflake.com/en/sql-reference/sql/"
)

// SdkShowByIDFilteringKind is one ShowByID filter, in call-argument order.
// An empty list emits LIKE. SdkShowByIDNone emits no filter and cannot be combined with another kind.
type SdkShowByIDFilteringKind string

const (
	SdkShowByIDLike       SdkShowByIDFilteringKind = "like"
	SdkShowByIDIn         SdkShowByIDFilteringKind = "in"
	SdkShowByIDExtendedIn SdkShowByIDFilteringKind = "extended_in"
	SdkShowByIDNone       SdkShowByIDFilteringKind = "none"
)

// SdkShowByIDPredicateKind selects the ShowByID row match.
// Empty keeps the generator default: name, or full id for a schema object with arguments.
type SdkShowByIDPredicateKind string

const (
	SdkShowByIDPredicateName           SdkShowByIDPredicateKind = "name"
	SdkShowByIDPredicateFullID         SdkShowByIDPredicateKind = "full_id"
	SdkShowByIDPredicateAccountName    SdkShowByIDPredicateKind = "account_name"
	SdkShowByIDPredicateNameAndLocator SdkShowByIDPredicateKind = "name_and_locator"
)

// SdkDescriptionMapping is the DESCRIBE result shape the SDK method emits.
type SdkDescriptionMapping string

const (
	SdkDescriptionMappingSingleValue SdkDescriptionMapping = "single_value"
	SdkDescriptionMappingSlice       SdkDescriptionMapping = "slice"
)

// SdkShowDoc overrides the default SHOW doc URL.
// The default is https://docs.snowflake.com/en/sql-reference/sql/show- plus the hyphenated NamePlural.
func SdkShowDoc(url string) SdkShowOption {
	return func(s *SdkGenerator) { s.ShowDoc = url }
}

// SdkShowSql overrides the default SHOW command words.
// The default is the uppercase words of NamePlural, so ComputePools is COMPUTE POOLS.
func SdkShowSql(sql string) SdkShowOption {
	return func(s *SdkGenerator) { s.ShowSql = sql }
}

func SdkShowTerse() SdkShowOption {
	return func(s *SdkGenerator) { s.ShowTerse = true }
}

// SdkShowOptionalSql appends optional SQL clauses, in order, after the SHOW object words.
func SdkShowOptionalSql(parts ...string) SdkShowOption {
	return func(s *SdkGenerator) { s.ShowOptionalSql = append(s.ShowOptionalSql, parts...) }
}

func SdkShowOptionalLike() SdkShowOption {
	return func(s *SdkGenerator) { s.ShowOptionalLike = true }
}

// SdkShowOptionalIn emits IN for account, database, and schema.
func SdkShowOptionalIn() SdkShowOption {
	return func(s *SdkGenerator) { s.ShowOptionalIn = true }
}

// SdkShowOptionalExtendedIn emits IN for account, database, schema, application, and application package.
func SdkShowOptionalExtendedIn() SdkShowOption {
	return func(s *SdkGenerator) { s.ShowOptionalExtendedIn = true }
}

// SdkShowOptionalOn adds the optional ON clause, after IN and before STARTS WITH.
func SdkShowOptionalOn() SdkShowOption {
	return func(s *SdkGenerator) { s.ShowOptionalOn = true }
}

// SdkShowResultFilterHook calls the excludeFromShow method of the SHOW row, which the SDK ext file implements,
// for each row. A row for which it returns true is left out of the result.
func SdkShowResultFilterHook() SdkShowOption {
	return func(s *SdkGenerator) { s.ShowResultFilterHook = true }
}

// SdkShowPlainOnlyField is a SHOW column that only the plain struct has, with Go type Kind (without a package prefix).
type SdkShowPlainOnlyField struct {
	Column *Attribute
	Kind   string
}

// SdkShowPlainOnly marks a SHOW column as plain-only. column is the same *Attribute that sits in Show.Columns.
// The column stays in Columns for the schema and the assertions, but the SHOW row does not read it:
// convert skips it, and additionalConvert in the SDK ext file fills the plain field. kind is the plain Go type.
func SdkShowPlainOnly(column *Attribute, kind string) SdkShowOption {
	return func(s *SdkGenerator) {
		s.ShowPlainOnlyFields = append(s.ShowPlainOnlyFields, SdkShowPlainOnlyField{Column: column, Kind: kind})
	}
}

func SdkShowOptionalStartsWith() SdkShowOption {
	return func(s *SdkGenerator) { s.ShowOptionalStartsWith = true }
}

func SdkShowOptionalLimitFrom() SdkShowOption {
	return func(s *SdkGenerator) { s.ShowOptionalLimitFrom = true }
}

func SdkShowOptionalLimit() SdkShowOption {
	return func(s *SdkGenerator) { s.ShowOptionalLimit = true }
}

func SdkShowByIDFiltering(kinds ...SdkShowByIDFilteringKind) SdkShowOption {
	return func(s *SdkGenerator) { s.ShowByIDFiltering = kinds }
}

// SdkShowByIDNoFiltering emits ShowByID with no LIKE, IN, or other filter.
func SdkShowByIDNoFiltering() SdkShowOption {
	return SdkShowByIDFiltering(SdkShowByIDNone)
}

// SdkShowByIDPredicate sets the ShowByID FindFirst match.
func SdkShowByIDPredicate(kind SdkShowByIDPredicateKind) SdkShowOption {
	return func(s *SdkGenerator) { s.ShowByIDPredicate = kind }
}

// SdkDescribeDoc overrides the default DESCRIBE doc URL.
// The default is https://docs.snowflake.com/en/sql-reference/sql/desc- plus the hyphenated Name.
func SdkDescribeDoc(url string) SdkDescribeOption {
	return func(s *SdkGenerator) { s.DescribeDoc = url }
}

// SdkDescribeSql overrides ObjectType.
// CREATE, ALTER, DROP, and DESCRIBE use that value unless the statement sets its own SQL.
func SdkDescribeSql(sql string) SdkDescribeOption {
	return func(s *SdkGenerator) { s.DescribeSql = sql }
}

// SdkDescribeMapping overrides the default DESCRIBE result shape.
// The default is a single value, or a slice when the object sets DescribeRows.
// Pass SdkDescriptionMappingSlice when DESCRIBE returns a slice of another shape.
func SdkDescribeMapping(kind SdkDescriptionMapping) SdkDescribeOption {
	return func(s *SdkGenerator) { s.DescribeMapping = kind }
}

// SdkDescribeQueryStructName overrides the default DESCRIBE query struct name.
// The default is Desc plus Name, so ComputePool is DescComputePool.
func SdkDescribeQueryStructName(name string) SdkDescribeOption {
	return func(s *SdkGenerator) { s.DescribeQueryStructName = name }
}

// SdkDescribeIdentifierOverride makes the DESCRIBE name a required identifier of this Go type, such as ObjectIdentifier,
// instead of the identifier type of the object. Use it when DESCRIBE accepts more than one kind of identifier,
// such as an outbound share name and an inbound share <provider_account>.<share_name>.
func SdkDescribeIdentifierOverride(kind string) SdkDescribeOption {
	return func(s *SdkGenerator) { s.DescribeIdentifierOverride = kind }
}

// SdkPlainOnlyField is a field on the DESCRIBE plain struct that DESCRIBE does not return.
// Name is the Go field name. Kind is the plain Go type, without a package prefix.
type SdkPlainOnlyField struct {
	Name string
	Kind string
}

// SdkDescribePlainOnlyField adds a plain-only field to the DESCRIBE paired struct.
// It is the ssot form of PlainOnlyField. The field is not a column and convert skips it.
// additionalConvert in the SDK ext file fills it, unless SdkDescribeDetailsFromIdentifier copies every plain-only field.
func SdkDescribePlainOnlyField(name, kind string) SdkDescribeOption {
	return func(s *SdkGenerator) {
		s.DescribePlainOnlyFields = append(s.DescribePlainOnlyFields, SdkPlainOnlyField{Name: name, Kind: kind})
	}
}

// SdkConfig returns the SDK settings for obj with defaults filled from Name and NamePlural.
// A functional option that sets a value wins over the default.
func (obj Object) SdkConfig() SdkGenerator {
	cfg := ObjectGenerators(obj.Generators...).Sdk
	if cfg.Show {
		plural := obj.namePluralForShow()
		if cfg.ShowDoc == "" {
			cfg.ShowDoc = showDocURLPrefix + docsObjectSlug(plural)
		}
		if cfg.ShowSql == "" {
			cfg.ShowSql = strings.ToUpper(genhelpers.CamelToWords(plural))
		}
		if cfg.ShowQueryStructName == "" {
			cfg.ShowQueryStructName = "Show" + plural
		}
	}
	if cfg.Describe {
		name := obj.nameForDescribe()
		if cfg.DescribeDoc == "" {
			cfg.DescribeDoc = describeDocURLPrefix + docsObjectSlug(name)
		}
		if cfg.DescribeSql == "" {
			cfg.DescribeSql = obj.ObjectType()
		}
		if cfg.DescribeQueryStructName == "" {
			cfg.DescribeQueryStructName = "Desc" + name
		}
		if cfg.DescribeMapping == "" {
			// A raw row list is a slice. Every property-list DESCRIBE in the SDK defs maps to a slice.
			cfg.DescribeMapping = SdkDescriptionMappingSingleValue
			if len(obj.DescribeRows.Columns) > 0 {
				cfg.DescribeMapping = SdkDescriptionMappingSlice
			}
		}
	}
	if cfg.Drop {
		name := obj.nameForDrop()
		if cfg.DropDoc == "" {
			cfg.DropDoc = dropDocURLPrefix + docsObjectSlug(name)
		}
		if cfg.DropSql == "" {
			cfg.DropSql = obj.ObjectType()
		}
		if cfg.DropQueryStructName == "" {
			cfg.DropQueryStructName = "Drop" + name
		}
	}
	if cfg.Create {
		name := obj.nameForCreate()
		if cfg.CreateDoc == "" {
			cfg.CreateDoc = createDocURLPrefix + docsObjectSlug(name)
		}
		if cfg.CreateSql == "" {
			cfg.CreateSql = obj.ObjectType()
		}
		if cfg.CreateQueryStructName == "" {
			cfg.CreateQueryStructName = "Create" + name
		}
	}
	if cfg.Alter {
		name := obj.nameForAlter()
		if cfg.AlterDoc == "" {
			cfg.AlterDoc = alterDocURLPrefix + docsObjectSlug(name)
		}
		if cfg.AlterSql == "" {
			cfg.AlterSql = obj.ObjectType()
		}
		if cfg.AlterQueryStructName == "" {
			cfg.AlterQueryStructName = "Alter" + name
		}
	}
	// TODO [next PRs]: These empty-name loops repeat ValidateObject. The SDK emitter copies them too
	// because make generate-sdk does not call ValidateObject. See CustomOperationsFromObject.
	if cfg.CustomOperations {
		obj.singularName("custom operation")
		for _, spec := range obj.CustomOperations {
			if spec.Kind == "" {
				panic(fmt.Sprintf("object %q: custom operation kind is empty", obj.Name))
			}
		}
	}
	if cfg.CustomInterfaceMethods {
		obj.singularName("custom interface method")
		for _, spec := range obj.CustomInterfaceMethods {
			if spec.Name == "" {
				panic(fmt.Sprintf("object %q: custom interface method name is empty", obj.Name))
			}
		}
	}
	return cfg
}

// CustomOperationDoc is the operation doc URL.
// An explicit Doc wins. Otherwise it is the kind slug, a hyphen, and the name slug.
func (obj Object) CustomOperationDoc(spec CustomOperationSpec) string {
	if spec.Doc != "" {
		return spec.Doc
	}
	if spec.Kind == "" {
		panic(fmt.Sprintf("object %q: custom operation kind is empty", obj.Name))
	}
	name := obj.singularName("custom operation")
	return customOperationDocURLPrefix + docsObjectSlug(spec.Kind) + "-" + docsObjectSlug(name)
}

// CustomOperationSQL is the uppercase words of Name.
func (obj Object) CustomOperationSQL() string {
	return strings.ToUpper(genhelpers.CamelToWords(obj.singularName("custom operation")))
}

// CustomOperationQueryStructName is the kind plus Name.
func (obj Object) CustomOperationQueryStructName(kind string) string {
	if kind == "" {
		panic(fmt.Sprintf("object %q: custom operation kind is empty", obj.Name))
	}
	return kind + obj.singularName("custom operation")
}

func (obj Object) namePluralForShow() string {
	if obj.NamePlural == "" {
		panic(fmt.Sprintf("object %q needs NamePlural to build the SHOW defaults", obj.Name))
	}
	return obj.NamePlural
}

// ObjectType is the Snowflake SQL object type words, so ComputePool is COMPUTE POOL.
// Sdk.DescribeSql overrides the words. The default is the uppercase words of Name.
func (obj Object) ObjectType() string {
	if sql := ObjectGenerators(obj.Generators...).Sdk.DescribeSql; sql != "" {
		return sql
	}
	if obj.Name == "" {
		panic(fmt.Sprintf("object %q needs Name to build the object type", obj.Name))
	}
	return strings.ToUpper(genhelpers.CamelToWords(obj.Name))
}

func (obj Object) nameForDescribe() string {
	if obj.Name == "" {
		panic(fmt.Sprintf("object %q needs Name to build the DESCRIBE defaults", obj.Name))
	}
	return obj.Name
}

func (obj Object) nameForDrop() string {
	return obj.singularName("DROP")
}

func (obj Object) nameForCreate() string {
	return obj.singularName("CREATE")
}

// VariantStem is the variant name with baseName removed. WarehouseAdaptive with base Warehouse is Adaptive.
func (obj Object) VariantStem(baseName string) string {
	if baseName == "" || !strings.HasPrefix(obj.Name, baseName) || obj.Name == baseName {
		panic(fmt.Sprintf("object %q: cannot build a variant stem from base %q", obj.Name, baseName))
	}
	return strings.TrimPrefix(obj.Name, baseName)
}

// VariantCreateMethod is Create plus the variant name with baseName removed.
// WarehouseAdaptive with base Warehouse is CreateAdaptive.
func (obj Object) VariantCreateMethod(baseName string) string {
	return "Create" + obj.VariantStem(baseName)
}

// VariantAlterMethod is Alter plus the variant name with baseName removed.
// WarehouseAdaptive with base Warehouse is AlterAdaptive.
func (obj Object) VariantAlterMethod(baseName string) string {
	if baseName == "" || !strings.HasPrefix(obj.Name, baseName) || obj.Name == baseName {
		panic(fmt.Sprintf("object %q: cannot build a variant ALTER method from base %q", obj.Name, baseName))
	}
	return "Alter" + strings.TrimPrefix(obj.Name, baseName)
}

// VariantAlterConfig is the ALTER config for emitting this variant as a custom method.
// It does not read WithSdk. Doc and AdditionalValidations come from Alter.
// SQL words are VariantEntityAlterName, or VariantOf.ObjectType when that override is empty.
// The query struct name is Alter plus SdkCreateRequestObjectName, or Alter plus Name.
func (obj Object) VariantAlterConfig() SdkGenerator {
	if !obj.HasAlterSpec() {
		panic(fmt.Sprintf("object %q: variant ALTER has no spec", obj.Name))
	}
	stem := obj.SdkCreateRequestObjectName
	if stem == "" {
		stem = obj.Name
	}
	return SdkGenerator{
		AlterDoc:                      obj.Alter.Doc,
		AlterSql:                      obj.variantEntityName(obj.VariantEntityAlterName),
		AlterQueryStructName:          "Alter" + stem,
		AlterAdditionalValidations:    obj.Alter.AdditionalValidations,
		AlterSetAdditionalValidations: obj.Alter.SetAdditionalValidations,
	}
}

// VariantCreateConfig is the CREATE config for emitting this variant as a custom method.
// It does not read WithSdk. Doc, AdditionalValidations, and WithoutOrReplace come from Create.
// SQL words are VariantEntityCreateName, or VariantOf.ObjectType when that override is empty.
// The query struct name is Create plus SdkCreateRequestObjectName, or Create plus Name.
func (obj Object) VariantCreateConfig() SdkGenerator {
	if !obj.HasCreateSpec() {
		panic(fmt.Sprintf("object %q: variant CREATE has no spec", obj.Name))
	}
	stem := obj.SdkCreateRequestObjectName
	if stem == "" {
		stem = obj.Name
	}
	return SdkGenerator{
		CreateDoc:                   obj.Create.Doc,
		CreateSql:                   obj.variantEntityName(obj.VariantEntityCreateName),
		CreateQueryStructName:       "Create" + stem,
		CreateAdditionalValidations: obj.Create.AdditionalValidations,
		CreateWithoutOrReplace:      obj.Create.WithoutOrReplace,
	}
}

// variantEntityName is the SQL object words for a variant statement.
// An override wins. Empty uses the base object's type words.
func (obj Object) variantEntityName(override string) string {
	if override != "" {
		return override
	}
	if obj.VariantOf == nil {
		panic(fmt.Sprintf("object %q: variant entity name needs VariantOf", obj.Name))
	}
	return obj.VariantOf.ObjectType()
}

func (obj Object) nameForAlter() string {
	return obj.singularName("ALTER")
}

func (obj Object) singularName(statement string) string {
	if obj.Name == "" {
		panic(fmt.Sprintf("object %q needs Name to build the %s defaults", obj.Name, statement))
	}
	return obj.Name
}

// docsObjectSlug is the docs path segment: snake case with hyphens, so ComputePools is compute-pools.
func docsObjectSlug(name string) string {
	return strings.ReplaceAll(genhelpers.ToSnakeCase(name), "_", "-")
}

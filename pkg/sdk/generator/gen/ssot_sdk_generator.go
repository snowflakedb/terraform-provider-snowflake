package gen

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

// parametersForLevel is the catalog a Parameters entry spreads. Tests replace it with a short list.
var parametersForLevel = ssotdefs.ParameterDefsForLevel

// variantsOf is the catalog a variant CREATE reads. Tests replace it with a short list.
var variantsOf = ssotdefs.VariantsOf

// ShowOperationFromObject adds SHOW and ShowByID from obj.SdkConfig.
// The field chain comes from ShowPairedStructs. Clauses are emitted in a fixed order:
// Show, Terse, SQL, OptionalSQL, Like, ExtendedIn, StartsWith, LimitFrom, Limit.
// An empty ShowByID filtering list passes no filtering argument, so ShowByID stays Like.
func (i *Interface) ShowOperationFromObject(obj ssot.Object) *Interface {
	cfg := obj.SdkConfig()
	if !cfg.Show {
		panic(fmt.Sprintf("object %q is not enabled for SDK SHOW generation", obj.Name))
	}
	return i.ShowOperationWithPairedStructs(
		cfg.ShowDoc,
		ShowPairedStructs(obj),
		showQueryStructFromConfig(cfg),
		showByIDFilteringFromObject(obj, cfg)...,
	)
}

// DescribeOperationFromObject adds DESCRIBE from obj.SdkConfig.
// The field chain comes from DescribePairedStructs. The query struct is Describe, SQL,
// Name, and ValidIdentifier on name.
// helperStructs are extra structs the describe operation emits.
// ShowPlainName and ShowVariants are emitted here too. Each one gets a projection on the paired SHOW struct.
// When ShowByID is already registered, each plain also gets ShowByID and ShowByIDSafely for that suffix.
// TODO [next PRs]: those plain structs are SHOW structs. Emit them from ShowOperationFromObject. Show() still returns the union row.
func (i *Interface) DescribeOperationFromObject(obj ssot.Object, helperStructs ...IntoField) *Interface {
	cfg := obj.SdkConfig()
	if !cfg.Describe {
		panic(fmt.Sprintf("object %q is not enabled for SDK DESCRIBE generation", obj.Name))
	}
	plains := plainShowHelpers(obj)
	helpers := make([]IntoField, 0, len(helperStructs)+len(plains))
	helpers = append(helpers, helperStructs...)
	for _, plain := range plains {
		helpers = append(helpers, plain)
	}
	if cfg.DescribeDetails {
		helpers = append(helpers, DescribeDetailsPlain(obj))
	}
	i = i.DescribeOperationWithPairedStructs(
		describeMappingFromObject(obj, cfg),
		cfg.DescribeDoc,
		DescribePairedStructs(obj),
		describeQueryStructFromConfig(cfg),
		helpers...,
	).showProjectionsFromPlains(obj, plains)
	if len(cfg.DescribeDetailsParse) > 0 {
		i = i.WithDescribeParse(describeParseConfig(obj, cfg))
	}
	if cfg.DescribeDetailsMethod {
		i = i.WithDescribeDetails(obj.Name+"Details", "Parse"+obj.Name+"Described")
	}
	return i
}

func describeParseConfig(obj ssot.Object, cfg ssot.SdkGenerator) DescribeParseConfig {
	config := DescribeParseConfig{
		RowsType:     obj.DescribePlainName(),
		DetailsType:  obj.Name + "Details",
		ErrorLabel:   genhelpers.CamelToWords(obj.Name),
		ValidateFunc: "validate" + obj.Name + "Described",
	}
	for _, mapping := range cfg.DescribeDetailsParse {
		field := sqlToFieldName(mapping.Attr.OutputKey(), true)
		if !mapping.Prefix {
			config.Properties = append(config.Properties, DescribeParseProperty{Property: mapping.Property, Field: field})
			continue
		}
		element := ssot.NestedDetailsElem(obj, mapping.Attr)
		config.Prefixes = append(config.Prefixes, DescribeParsePrefix{
			Prefix:        mapping.Property,
			Field:         field,
			Var:           unexported(field),
			ElementType:   element,
			ElementParser: "Parse" + strings.TrimSuffix(element, "Details") + "Described",
		})
	}
	return config
}

func plainShowHelpers(obj ssot.Object) []*Plain {
	var out []*Plain
	if obj.ShowPlainName != "" {
		if len(obj.Show.Columns) == 0 {
			panic(fmt.Sprintf("object %q: ShowPlainName %q has no SHOW columns", obj.Name, obj.ShowPlainName))
		}
		out = append(out, plainFromColumns(obj, obj.ShowPlainName, obj.Show.Columns))
	}
	for _, variant := range obj.ShowVariants {
		if variant == nil || variant.Name == "" || len(variant.Show.Columns) == 0 {
			panic(fmt.Sprintf("object %q: show variant needs a name and SHOW columns", obj.Name))
		}
		out = append(out, plainFromColumns(*variant, variant.Name, variant.Show.Columns))
	}
	return out
}

// showProjectionsFromPlains attaches one projection per plain onto the paired SHOW struct.
// The method is As plus the plain name with the object name removed.
// When ShowByID is already registered, the same suffix is passed to ShowByIDProjection.
func (i *Interface) showProjectionsFromPlains(obj ssot.Object, targets []*Plain) *Interface {
	if len(targets) == 0 {
		return i
	}
	columns := obj.PairedShowColumns()
	if len(columns) == 0 {
		panic(fmt.Sprintf("object %q: SHOW projections have no paired SHOW columns", obj.Name))
	}
	source := plainFromColumns(obj, obj.Name, columns)
	for _, target := range targets {
		method := showProjectionMethod(obj.Name, target.name)
		i.ShowProjection(source, method, target)
		if i.hasShowByID() {
			i.ShowByIDProjection(strings.TrimPrefix(method, "As"))
		}
	}
	return i
}

// showProjectionMethod is As plus plainName with baseName removed.
// WarehouseRegular on Warehouse is AsRegular. WarehouseAdaptive is AsAdaptive.
func showProjectionMethod(baseName, plainName string) string {
	if baseName == "" || !strings.HasPrefix(plainName, baseName) || plainName == baseName {
		panic(fmt.Sprintf("object %q: cannot build a SHOW projection method from %q", baseName, plainName))
	}
	return "As" + strings.TrimPrefix(plainName, baseName)
}

func showQueryStructFromConfig(cfg ssot.SdkGenerator) *QueryStruct {
	q := NewQueryStruct(cfg.ShowQueryStructName).Show()
	if cfg.ShowTerse {
		q = q.Terse()
	}
	q = q.SQL(cfg.ShowSql)
	for _, part := range cfg.ShowOptionalSql {
		q = q.OptionalSQL(part)
	}
	if cfg.ShowOptionalLike {
		q = q.OptionalLike()
	}
	if cfg.ShowOptionalExtendedIn {
		q = q.OptionalExtendedIn()
	}
	if cfg.ShowOptionalStartsWith {
		q = q.OptionalStartsWith()
	}
	if cfg.ShowOptionalLimitFrom {
		q = q.OptionalLimitFrom()
	}
	if cfg.ShowOptionalLimit {
		q = q.OptionalLimit()
	}
	return q
}

func showByIDFilteringFromObject(obj ssot.Object, cfg ssot.SdkGenerator) []ShowByIDFilteringKind {
	if len(cfg.ShowByIDFiltering) == 0 {
		return nil
	}
	filtering := make([]ShowByIDFilteringKind, len(cfg.ShowByIDFiltering))
	for i, kind := range cfg.ShowByIDFiltering {
		switch kind {
		case ssot.SdkShowByIDLike:
			filtering[i] = ShowByIDLikeFiltering
		case ssot.SdkShowByIDExtendedIn:
			filtering[i] = ShowByIDExtendedInFiltering
		default:
			panic(fmt.Sprintf("object %q: unknown ShowByID filtering %q", obj.Name, kind))
		}
	}
	return filtering
}

func describeMappingFromObject(obj ssot.Object, cfg ssot.SdkGenerator) DescriptionMappingKind {
	switch cfg.DescribeMapping {
	case ssot.SdkDescriptionMappingSingleValue:
		return DescriptionMappingKindSingleValue
	case ssot.SdkDescriptionMappingSlice:
		return DescriptionMappingKindSlice
	default:
		panic(fmt.Sprintf("object %q: unknown description mapping %q", obj.Name, cfg.DescribeMapping))
	}
}

func describeQueryStructFromConfig(cfg ssot.SdkGenerator) *QueryStruct {
	return NewQueryStruct(cfg.DescribeQueryStructName).
		Describe().
		SQL(cfg.DescribeSql).
		Name().
		WithValidation(ValidIdentifier, "name")
}

// DropOperationFromObject adds DROP from obj.SdkConfig.
// Clauses are emitted in a fixed order: Drop, SQL, IfExists, Name, ValidIdentifier.
// SdkDropCascadeRestrict adds CASCADE, RESTRICT, and their conflict check after that.
// DropSafely stays the template default unless a drop option sets the hook or force.
func (i *Interface) DropOperationFromObject(obj ssot.Object) *Interface {
	cfg := obj.SdkConfig()
	if !cfg.Drop {
		panic(fmt.Sprintf("object %q is not enabled for SDK DROP generation", obj.Name))
	}
	var opts []DropOperationOption
	if cfg.DropSafelyHook {
		opts = append(opts, WithDropSafelyHook())
	}
	if cfg.DropSafelyForce {
		opts = append(opts, WithDropSafelyForce())
	}
	return i.DropOperation(cfg.DropDoc, dropQueryStructFromConfig(cfg), opts...)
}

func dropQueryStructFromConfig(cfg ssot.SdkGenerator) *QueryStruct {
	q := NewQueryStruct(cfg.DropQueryStructName).
		Drop().
		SQL(cfg.DropSql).
		IfExists().
		Name().
		WithValidation(ValidIdentifier, "name")
	if cfg.DropCascadeRestrict {
		q = q.OptionalSQL("CASCADE").
			OptionalSQL("RESTRICT").
			WithValidation(ConflictingFields, "Cascade", "Restrict")
	}
	return q
}

// CreateOperationFromObject adds CREATE from obj.SdkConfig.
// Clauses are Create, optional OR REPLACE, Create.KeywordsBefore, SQL, optional IF NOT EXISTS, Name, then Create.Clauses, then ValidIdentifier.
// helperStructs are query structs the attribute clauses name, other than an enum list item.
// A list of enums adds its item query struct here. The struct name is the enum name plus ListItem, and its field is Value.
func (i *Interface) CreateOperationFromObject(obj ssot.Object, helperStructs ...IntoField) *Interface {
	if obj.VariantOf != nil {
		panic(fmt.Sprintf("object %q: a variant cannot enable SDK generation", obj.Name))
	}
	cfg := obj.SdkConfig()
	if !cfg.Create {
		panic(fmt.Sprintf("object %q is not enabled for SDK CREATE generation", obj.Name))
	}
	return i.CreateOperation(cfg.CreateDoc, createQueryStructFromConfig(obj, cfg), append(helperStructs, enumListItemStructs(obj)...)...)
}

// CreateVariantOperationsFromObject emits each variant CREATE as a custom method.
// Variants come from the catalog for base. The variant supplies the CREATE spec. The base must turn on SdkCreateForVariants.
// TODO [next PRs]: the same switch for variant SHOW/DESCRIBE.
func (i *Interface) CreateVariantOperationsFromObject(base ssot.Object) *Interface {
	if !base.SdkConfig().CreateForVariants {
		panic(fmt.Sprintf("object %q is not enabled for variant CREATE generation", base.Name))
	}
	variants := variantsOf(base)
	if len(variants) == 0 {
		panic(fmt.Sprintf("object %q: variant CREATE generation has no variants", base.Name))
	}
	for _, variant := range variants {
		if variant.VariantOf == nil || variant.VariantOf.Name != base.Name {
			panic(fmt.Sprintf("object %q: %q is not a variant", base.Name, variant.Name))
		}
		if !variant.HasCreateSpec() {
			panic(fmt.Sprintf("object %q: variant CREATE has no spec", variant.Name))
		}
		cfg := variant.VariantCreateConfig()
		q := createQueryStructFromConfig(variant, cfg)
		i = i.CustomOperation(variant.VariantCreateMethod(base.Name), cfg.CreateDoc, q, enumListItemStructs(variant)...)
	}
	return i
}

func createQueryStructFromConfig(obj ssot.Object, cfg ssot.SdkGenerator) *QueryStruct {
	q := NewQueryStruct(cfg.CreateQueryStructName).Create()
	q.owner = obj.Name
	if !cfg.CreateWithoutOrReplace {
		q = q.OrReplace()
	}
	for _, field := range obj.Create.KeywordsBefore {
		if field != nil && field.IsTags() {
			panic("KeywordsBefore list contains tags")
		}
		if field != nil && field.IsParameters() {
			panic("KeywordsBefore list contains parameters")
		}
	}
	createParams := parametersForStatement(obj, obj.Create.Clauses, "CREATE")
	q = q.appendSqlSequence(obj.Create.KeywordsBefore, attributeSqlCreate).SQL(cfg.CreateSql)
	if !cfg.CreateWithoutIfNotExists {
		q = q.IfNotExists()
	}
	q = q.Name()
	q = q.appendSqlSequenceParams(obj.Create.Clauses, attributeSqlCreate, createParams).WithValidation(ValidIdentifier, "name")
	if !cfg.CreateWithoutOrReplace && !cfg.CreateWithoutIfNotExists {
		q = q.WithValidation(ConflictingFields, "OrReplace", "IfNotExists")
	}
	q = q.withIdentifierParameterValidations(createParams)
	if cfg.CreateAdditionalValidations {
		q = q.WithAdditionalValidations()
	}
	return q
}

// TODO [next PRs]: Delete the repeated definition checks in CustomOperationsFromObject,
// CustomInterfaceMethodsFromObject, customInterfaceMethodParameters, and Object.SdkConfig.
// ValidateObject already rejects an empty name, a duplicate name, a missing IdType, and a parameter
// with no type. The copies exist because make generate-sdk does not call ValidateObject. It reads the
// object directly. The generator tests also call these emitters on a fake object without ValidateObject.
// Run ValidateObject from SDK generation, or share one checker, and delete the copies. Keep the opt-in
// panic in the emitter. That checks whether the call is enabled, not whether the spec is well formed.

// CustomOperationsFromObject adds each custom operation from obj.CustomOperations, in slice order.
// The operations are emitted together at this call. Clause order is the emission order.
// The query struct name is the kind plus Name.
// A name clause adds ValidIdentifier. OR REPLACE and IF NOT EXISTS together are marked conflicting.
func (i *Interface) CustomOperationsFromObject(obj ssot.Object) *Interface {
	requireCustomOperations(obj)
	seenKind := make(map[string]struct{})
	for _, spec := range obj.CustomOperations {
		if spec.Kind == "" {
			panic(fmt.Sprintf("object %q: custom operation kind is empty", obj.Name))
		}
		if _, dup := seenKind[spec.Kind]; dup {
			panic(fmt.Sprintf("object %q: duplicate custom operation %q", obj.Name, spec.Kind))
		}
		seenKind[spec.Kind] = struct{}{}
		i = i.customOperationFromSpec(obj, spec)
	}
	return i
}

func requireCustomOperations(obj ssot.Object) {
	if !obj.SdkConfig().CustomOperations {
		panic(fmt.Sprintf("object %q is not enabled for SDK custom operation generation", obj.Name))
	}
}

// CustomInterfaceMethodsFromObject adds each custom interface method from obj.CustomInterfaceMethods, in slice order.
// The declaration is generated. The body stays in the SDK ext file.
// A nil parameter list is id of IdType. An empty return list is error.
func (i *Interface) CustomInterfaceMethodsFromObject(obj ssot.Object) *Interface {
	if !obj.SdkConfig().CustomInterfaceMethods {
		panic(fmt.Sprintf("object %q is not enabled for SDK custom interface method generation", obj.Name))
	}
	seen := make(map[string]struct{})
	for _, spec := range obj.CustomInterfaceMethods {
		if spec.Name == "" {
			panic(fmt.Sprintf("object %q: custom interface method name is empty", obj.Name))
		}
		if _, dup := seen[spec.Name]; dup {
			panic(fmt.Sprintf("object %q: duplicate custom interface method %q", obj.Name, spec.Name))
		}
		seen[spec.Name] = struct{}{}
		returns := spec.Returns
		if len(returns) == 0 {
			returns = []string{"error"}
		}
		i = i.WithCustomInterfaceMethod(spec.Name, spec.Doc, customInterfaceMethodParameters(obj, spec), returns...)
	}
	return i
}

func customInterfaceMethodParameters(obj ssot.Object, spec ssot.CustomInterfaceMethodSpec) []*MethodParameter {
	if spec.Parameters == nil {
		if obj.IdType == "" {
			panic(fmt.Sprintf("object %q: custom interface method %q needs IdType", obj.Name, spec.Name))
		}
		return []*MethodParameter{NewMethodParameter("id", obj.IdType)}
	}
	params := make([]*MethodParameter, len(spec.Parameters))
	for i, param := range spec.Parameters {
		if param.Name == "" {
			panic(fmt.Sprintf("object %q: custom interface method %q parameter name is empty", obj.Name, spec.Name))
		}
		kind := param.Kind
		if kind == "" {
			kind = obj.IdType
		}
		if kind == "" {
			panic(fmt.Sprintf("object %q: custom interface method %q parameter %q needs a type", obj.Name, spec.Name, param.Name))
		}
		params[i] = NewMethodParameter(param.Name, kind)
	}
	return params
}

// ShowParametersFromObject adds the ShowParameters declaration.
// It is the custom interface method ShowParameters: empty doc, id of IdType, []*Parameter and error.
// The body stays in the SDK ext file.
func (i *Interface) ShowParametersFromObject(obj ssot.Object) *Interface {
	if !obj.SdkConfig().ShowParameters {
		panic(fmt.Sprintf("object %q is not enabled for SDK ShowParameters generation", obj.Name))
	}
	return i.ShowParameters(obj.IdType)
}

// ShowParametersDetailsFromObject adds ShowParametersDetails and the typed details struct.
// The fields are ParameterLevel() plus ShowParametersDetailsAdditionalLevels.
// ALTER SET and UNSET stay ParameterLevel(). The body is generated and calls ShowParameters.
func (i *Interface) ShowParametersDetailsFromObject(obj ssot.Object) *Interface {
	if !obj.SdkConfig().ShowParametersDetails {
		panic(fmt.Sprintf("object %q is not enabled for SDK ShowParametersDetails generation", obj.Name))
	}
	levels := append([]ssot.ParameterLevel{obj.ParameterLevel()}, obj.ShowParametersDetailsAdditionalLevels...)
	return i.ShowParametersDetails(parametersForLevels(levels)...)
}

// InterfaceFromObject builds one SDK interface from the operations whose flags are on.
// A flag that is off is absent. The existing FromObject helpers still panic if called while their flag is off.
// Order is Create, Alter, Drop, Show, Describe, ShowParameters, ShowParametersDetails, custom operations, then custom interface methods.
func InterfaceFromObject(obj ssot.Object) *Interface {
	if obj.NamePlural == "" || obj.Name == "" || obj.IdType == "" {
		panic(fmt.Sprintf("object %q needs NamePlural, Name, and IdType to build an SDK interface", obj.Name))
	}
	cfg := obj.SdkConfig()
	iface := NewInterface(obj.NamePlural, obj.Name, obj.IdType)
	if cfg.Create {
		iface = iface.CreateOperationFromObject(obj)
	}
	if cfg.CreateForVariants {
		iface = iface.CreateVariantOperationsFromObject(obj)
	}
	if cfg.Alter {
		iface = iface.AlterOperationFromObject(obj)
	}
	if cfg.AlterForVariants {
		iface = iface.AlterVariantOperationsFromObject(obj)
	}
	if cfg.Drop {
		iface = iface.DropOperationFromObject(obj)
	}
	if cfg.Show {
		iface = iface.ShowOperationFromObject(obj)
	}
	if cfg.Describe {
		iface = iface.DescribeOperationFromObject(obj)
	}
	if cfg.ShowParameters {
		iface = iface.ShowParametersFromObject(obj)
	}
	if cfg.ShowParametersDetails {
		iface = iface.ShowParametersDetailsFromObject(obj)
	}
	if cfg.CustomOperations {
		iface = iface.CustomOperationsFromObject(obj)
	}
	if cfg.CustomInterfaceMethods {
		iface = iface.CustomInterfaceMethodsFromObject(obj)
	}
	if enums := referencedEnums(obj); len(enums) > 0 {
		iface = iface.WithEnums(enums...)
	}
	return iface
}

// SdkInterfaceOverride changes one ssot-built SDK interface. The key is the ssot object Name.
type SdkInterfaceOverride func(*Interface)

// SdkInterfacesFromSsot builds SDK interfaces for objects that turn on SdkDefineFromSsot.
// NamePlural must be unique among objects that are not variants. A variant shares the base interface.
// A shared plural among non-variants, such as the catalog integration variants, stays on the old def.
// The interface contains the operation parts whose flags are on.
// overrides is keyed by the ssot object Name. Each function receives the built interface and may add to it.
// A key that does not match a built interface panics.
func SdkInterfacesFromSsot(objects []ssot.Object, overrides map[string]SdkInterfaceOverride) []*Interface {
	counts := make(map[string]int, len(objects))
	for _, obj := range objects {
		if obj.NamePlural == "" || obj.VariantOf != nil {
			continue
		}
		counts[obj.NamePlural]++
	}
	applied := make(map[string]struct{}, len(overrides))
	var out []*Interface
	for _, obj := range objects {
		if obj.NamePlural == "" || counts[obj.NamePlural] != 1 || !obj.SdkConfig().DefineFromSsot {
			continue
		}
		iface := InterfaceFromObject(obj)
		if override, ok := overrides[obj.Name]; ok {
			override(iface)
			applied[obj.Name] = struct{}{}
		}
		out = append(out, iface)
	}
	for name := range overrides {
		if _, ok := applied[name]; !ok {
			panic(fmt.Sprintf("sdk interface override %q did not match an ssot interface", name))
		}
	}
	return out
}

// referencedEnums is the enums on the object's columns and attributes, first reference first.
// SHOW columns come first. ColumnEnum prefers ShowEnum, so the SHOW type is the generated SDK enum.
// An enum with SdkType is a parser name for another SDK type and is not generated.
func referencedEnums(obj ssot.Object) []*ssot.Enum {
	var out []*ssot.Enum
	seen := map[*ssot.Enum]struct{}{}
	addEnum := func(enum *ssot.Enum) {
		if enum == nil || enum.SdkType != "" {
			return
		}
		if _, ok := seen[enum]; ok {
			return
		}
		seen[enum] = struct{}{}
		out = append(out, enum)
	}
	var addAttr func(attr *ssot.Attribute)
	addAttr = func(attr *ssot.Attribute) {
		if attr == nil {
			return
		}
		addEnum(attr.ColumnEnum())
		if attr.Elem != nil {
			addEnum(attr.Elem.ColumnEnum())
		}
		addEnum(attr.CreateSql.Enum)
		for i := range attr.Children {
			addAttr(&attr.Children[i])
		}
		if attr.Elem != nil && attr.Elem.Type == ssot.AttrObject {
			for i := range attr.Elem.Children {
				addAttr(&attr.Elem.Children[i])
			}
			for _, variant := range attr.Elem.CreateSql.Variants {
				for _, field := range variant.Fields {
					addAttr(field)
				}
			}
		}
	}
	addList := func(list []*ssot.Attribute) {
		for _, attr := range list {
			addAttr(attr)
		}
	}
	addList(obj.PairedShowColumns())
	addList(obj.PairedDescribeColumns())
	for i := range obj.Attributes {
		addAttr(&obj.Attributes[i])
	}
	addList(obj.Create.KeywordsBefore)
	addList(obj.Create.Clauses)
	addList(obj.Set)
	addList(obj.Unset)
	for _, enum := range obj.AdditionalEnums {
		addEnum(enum)
	}
	for _, spec := range obj.CustomOperations {
		for _, clause := range spec.Clauses {
			if clause.Kind == ssot.CustomOperationClauseAttribute {
				addAttr(clause.Attribute)
			}
		}
	}
	return out
}

func (i *Interface) customOperationFromSpec(obj ssot.Object, spec ssot.CustomOperationSpec) *Interface {
	return i.CustomOperation(spec.Kind, obj.CustomOperationDoc(spec), customOperationQueryStruct(obj, spec))
}

func customOperationQueryStruct(obj ssot.Object, spec ssot.CustomOperationSpec) *QueryStruct {
	q := NewQueryStruct(obj.CustomOperationQueryStructName(spec.Kind))
	sawOrReplace := false
	sawIfNotExists := false
	sawName := false
	seen := make(map[ssot.CustomOperationClauseKind]struct{})
	seenAttr := make(map[string]struct{})
	seenPredefined := make(map[string]struct{})
	singleton := func(kind ssot.CustomOperationClauseKind) {
		if _, dup := seen[kind]; dup {
			panic(fmt.Sprintf("object %q: custom operation %q has duplicate %s", obj.Name, spec.Kind, kind))
		}
		seen[kind] = struct{}{}
	}
	for _, clause := range spec.Clauses {
		switch clause.Kind {
		case ssot.CustomOperationClauseStaticSQL:
			if clause.SQL == "" {
				panic(fmt.Sprintf("object %q: custom operation %q static SQL is empty", obj.Name, spec.Kind))
			}
			q = q.SQL(clause.SQL)
		case ssot.CustomOperationClauseOrReplace:
			singleton(clause.Kind)
			sawOrReplace = true
			q = q.OrReplace()
		case ssot.CustomOperationClauseIfNotExists:
			singleton(clause.Kind)
			sawIfNotExists = true
			q = q.IfNotExists()
		case ssot.CustomOperationClauseAttribute:
			attr := customOperationAttribute(obj, spec, clause, seenAttr)
			q = q.appendSqlSequence([]*ssot.Attribute{attr}, attributeSqlCreate)
		case ssot.CustomOperationClauseObjectWords:
			singleton(clause.Kind)
			q = q.SQL(obj.CustomOperationSQL())
		case ssot.CustomOperationClauseName:
			singleton(clause.Kind)
			sawName = true
			q = q.Name()
		case ssot.CustomOperationClausePredefined:
			if clause.Predefined == "" {
				panic(fmt.Sprintf("object %q: custom operation %q predefined type is empty", obj.Name, spec.Kind))
			}
			if _, dup := seenPredefined[clause.Predefined]; dup {
				panic(fmt.Sprintf("object %q: custom operation %q has duplicate predefined %q", obj.Name, spec.Kind, clause.Predefined))
			}
			seenPredefined[clause.Predefined] = struct{}{}
			q = q.PredefinedQueryStructField(clause.Predefined, clause.Predefined, KeywordOptions())
		default:
			panic(fmt.Sprintf("object %q: custom operation %q clause has kind %q", obj.Name, spec.Kind, clause.Kind))
		}
	}
	if sawName {
		q = q.WithValidation(ValidIdentifier, "name")
	}
	if sawOrReplace && sawIfNotExists {
		q = q.WithValidation(ConflictingFields, "OrReplace", "IfNotExists")
	}
	if spec.AdditionalValidations {
		q = q.WithAdditionalValidations()
	}
	return q
}

func customOperationAttribute(obj ssot.Object, spec ssot.CustomOperationSpec, clause ssot.CustomOperationClause, seenAttr map[string]struct{}) *ssot.Attribute {
	if clause.Attribute == nil {
		panic(fmt.Sprintf("object %q: custom operation %q attribute is empty", obj.Name, spec.Kind))
	}
	if clause.Attribute.IsTags() {
		panic(fmt.Sprintf("object %q: custom operation %q attribute is tags", obj.Name, spec.Kind))
	}
	if clause.Attribute.IsParameters() {
		panic(fmt.Sprintf("object %q: custom operation %q attribute is parameters", obj.Name, spec.Kind))
	}
	name := clause.Attribute.Name
	if name == "" {
		panic(fmt.Sprintf("object %q: custom operation %q attribute has empty Name", obj.Name, spec.Kind))
	}
	if _, dup := seenAttr[name]; dup {
		panic(fmt.Sprintf("object %q: custom operation %q has duplicate attribute %q", obj.Name, spec.Kind, name))
	}
	seenAttr[name] = struct{}{}
	return clause.Attribute
}

// AlterOperationFromObject adds ALTER from obj.SdkConfig.
// Clauses are Alter, SQL, optional IF EXISTS, Name, then obj.Alter in order, then ValidIdentifier.
// AtLeastOneValueSet on SET and UNSET lists the fields those structs contain.
// ExactlyOneValueSet follows obj.Alter.
func (i *Interface) AlterOperationFromObject(obj ssot.Object) *Interface {
	if obj.VariantOf != nil {
		panic(fmt.Sprintf("object %q: a variant cannot enable SDK generation", obj.Name))
	}
	cfg := obj.SdkConfig()
	if !cfg.Alter {
		panic(fmt.Sprintf("object %q is not enabled for SDK ALTER generation", obj.Name))
	}
	return i.AlterOperation(cfg.AlterDoc, alterQueryStructFromConfig(obj, cfg)).withEnumListItems(obj)
}

// AlterVariantOperationsFromObject emits each variant ALTER as a custom method.
// Variants come from the catalog for base. The variant supplies the ALTER spec. The base must turn on SdkAlterForVariants.
// TODO [next PRs]: the same switch for variant SHOW/DESCRIBE.
func (i *Interface) AlterVariantOperationsFromObject(base ssot.Object) *Interface {
	if !base.SdkConfig().AlterForVariants {
		panic(fmt.Sprintf("object %q is not enabled for variant ALTER generation", base.Name))
	}
	variants := variantsOf(base)
	if len(variants) == 0 {
		panic(fmt.Sprintf("object %q: variant ALTER generation has no variants", base.Name))
	}
	for _, variant := range variants {
		if variant.VariantOf == nil || variant.VariantOf.Name != base.Name {
			panic(fmt.Sprintf("object %q: %q is not a variant", base.Name, variant.Name))
		}
		if !variant.HasAlterSpec() {
			panic(fmt.Sprintf("object %q: variant ALTER has no spec", variant.Name))
		}
		cfg := variant.VariantAlterConfig()
		q := alterQueryStructFromConfig(variant, cfg)
		i = i.CustomOperation(variant.VariantAlterMethod(base.Name), cfg.AlterDoc, q)
		i = i.withEnumListItems(variant)
	}
	return i
}

func (i *Interface) withEnumListItems(obj ssot.Object) *Interface {
	if len(i.Operations) == 0 {
		return i
	}
	i.Operations[len(i.Operations)-1].withHelperStructs(enumListItemFields(obj)...)
	return i
}

func enumListItemQueryStructs(obj ssot.Object) []*QueryStruct {
	// TODO [next PRs]: Move this shape into ssot as a description (struct name, field Value, the enum, single quotes, required). This function should only translate that description into a query struct. Do that together with more complex list items, not for this enum case alone.
	seen := map[string]struct{}{}
	var out []*QueryStruct
	add := func(list []*ssot.Attribute) {
		for _, attr := range list {
			if attr == nil || attr.Type != ssot.AttrList || attr.Elem == nil || attr.Elem.Type != ssot.AttrEnum || attr.Elem.Enum == nil {
				continue
			}
			name := ssot.EnumListItemName(attr.Elem.Enum)
			if _, ok := seen[name]; ok {
				continue
			}
			seen[name] = struct{}{}
			out = append(out, NewQueryStruct(name).Enum("Value", attr.Elem.Enum, KeywordOptions().SingleQuotes().Required()))
		}
	}
	add(obj.Create.KeywordsBefore)
	add(obj.Create.Clauses)
	add(obj.Set)
	add(obj.Unset)
	return out
}

func enumListItemStructs(obj ssot.Object) []IntoField {
	qs := enumListItemQueryStructs(obj)
	out := make([]IntoField, len(qs))
	for i := range qs {
		out[i] = qs[i]
	}
	return out
}

func enumListItemFields(obj ssot.Object) []*Field {
	qs := enumListItemQueryStructs(obj)
	out := make([]*Field, len(qs))
	for i := range qs {
		out[i] = qs[i].IntoField()
	}
	return out
}

func alterQueryStructFromConfig(obj ssot.Object, cfg ssot.SdkGenerator) *QueryStruct {
	q := NewQueryStruct(cfg.AlterQueryStructName).Alter().SQL(cfg.AlterSql)
	q.owner = obj.Name
	if !cfg.AlterWithoutIfExists {
		q = q.IfExists()
	}
	q = q.Name()
	seen := make(map[ssot.AlterClauseKind]struct{})
	for _, clause := range obj.Alter.Clauses {
		switch clause.Kind {
		case ssot.AlterClauseKeyword:
			if clause.SQL == "" {
				panic(fmt.Sprintf("object %q: ALTER keyword is empty", obj.Name))
			}
			q = q.OptionalSQL(clause.SQL)
		case ssot.AlterClauseRenameTo:
			if _, dup := seen[clause.Kind]; dup {
				panic(fmt.Sprintf("object %q: duplicate ALTER rename_to", obj.Name))
			}
			seen[clause.Kind] = struct{}{}
			q = q.RenameTo()
		case ssot.AlterClauseSwapWith:
			if _, dup := seen[clause.Kind]; dup {
				panic(fmt.Sprintf("object %q: duplicate ALTER swap_with", obj.Name))
			}
			seen[clause.Kind] = struct{}{}
			q = q.Identifier("SwapWith", InterfaceIdentifierPointerKind, IdentifierOptions().SQL("SWAP WITH"))
		case ssot.AlterClauseSet:
			if _, dup := seen[clause.Kind]; dup {
				panic(fmt.Sprintf("object %q: duplicate ALTER set", obj.Name))
			}
			seen[clause.Kind] = struct{}{}
			q = appendAlterSet(q, obj, cfg)
		case ssot.AlterClauseUnset:
			if _, dup := seen[clause.Kind]; dup {
				panic(fmt.Sprintf("object %q: duplicate ALTER unset", obj.Name))
			}
			seen[clause.Kind] = struct{}{}
			q = appendAlterUnset(q, obj)
		case ssot.AlterClauseSetTags:
			if _, dup := seen[clause.Kind]; dup {
				panic(fmt.Sprintf("object %q: duplicate ALTER set_tags", obj.Name))
			}
			seen[clause.Kind] = struct{}{}
			q = q.OptionalSetTags()
		case ssot.AlterClauseUnsetTags:
			if _, dup := seen[clause.Kind]; dup {
				panic(fmt.Sprintf("object %q: duplicate ALTER unset_tags", obj.Name))
			}
			seen[clause.Kind] = struct{}{}
			q = q.OptionalUnsetTags()
		case ssot.AlterClauseList:
			if clause.SQL == "" || clause.ListItem == "" {
				panic(fmt.Sprintf("object %q: ALTER list %q has no SQL or item type", obj.Name, clause.SQL))
			}
			q = q.ListAssignment(clause.SQL, clause.ListItem, ParameterOptions().Parentheses().NoEquals())
		case ssot.AlterClauseText:
			if clause.SQL == "" {
				panic(fmt.Sprintf("object %q: ALTER text is empty", obj.Name))
			}
			q = q.OptionalTextAssignment(clause.SQL, ParameterOptions().SingleQuotes().NoEquals())
		case ssot.AlterClauseBlock:
			block := queryStructFromBlockAttribute(obj.Name, *clause.Block)
			var transformer FieldTransformer = ParameterOptions().SQL(clause.SQL)
			if clause.BlockKeyword {
				transformer = KeywordOptions().SQL(clause.SQL)
			}
			q = q.OptionalQueryStructField(clause.Field, block, transformer)
		default:
			panic(fmt.Sprintf("object %q: ALTER clause has kind %q", obj.Name, clause.Kind))
		}
	}
	q = q.WithValidation(ValidIdentifier, "name")
	if !cfg.AlterWithoutExactlyOneValueSet {
		q = q.WithValidation(ExactlyOneValueSet, alterExactlyOneFields(obj)...)
	}
	q = q.withAlterIdentifierValidations(obj)
	for _, clause := range obj.Alter.Clauses {
		if clause.Kind == ssot.AlterClauseList {
			q = q.WithValidation(ValidIdentifier, sqlToFieldName(clause.SQL, true))
		}
	}
	if cfg.AlterAdditionalValidations {
		q = q.WithAdditionalValidations()
	}
	return q
}

func appendAlterSet(q *QueryStruct, obj ssot.Object, cfg ssot.SdkGenerator) *QueryStruct {
	setParams := parametersForStatement(obj, obj.Set, "SET")
	setStruct := NewQueryStruct(obj.Name+"Set").appendSqlSequenceParams(obj.Set, attributeSqlSet, setParams)
	setStruct = setStruct.withIdentifierParameterValidations(setParams)
	setStruct = setStruct.WithValidation(AtLeastOneValueSet, atLeastOneFields(setStruct, setParams)...)
	if cfg.AlterSetAdditionalValidations {
		setStruct = setStruct.WithAdditionalValidations()
	}
	var setOptions FieldTransformer = KeywordOptions().SQL("SET")
	if cfg.AlterSetList {
		setOptions = ListOptions().NoParentheses().SQL("SET")
	}
	return q.OptionalQueryStructField("Set", setStruct, setOptions)
}

func appendAlterUnset(q *QueryStruct, obj ssot.Object) *QueryStruct {
	unsetParams := parametersForStatement(obj, obj.Unset, "UNSET")
	unsetStruct := NewQueryStruct(obj.Name+"Unset").appendSqlSequenceParams(obj.Unset, attributeSqlUnset, unsetParams)
	unsetStruct = unsetStruct.WithValidation(AtLeastOneValueSet, atLeastOneFields(unsetStruct, unsetParams)...)
	return q.OptionalQueryStructField("Unset", unsetStruct, ListOptions().NoParentheses().SQL("UNSET"))
}

func parametersForStatement(obj ssot.Object, list []*ssot.Attribute, statement string) []ssot.ParameterDef {
	for _, field := range list {
		if field != nil && field.IsParameters() {
			return parametersForLevels(obj.ParameterLevelsFor(statement))
		}
	}
	return nil
}

// parametersForLevels is the union of levels in catalog order.
// One level uses parametersForLevel so tests can replace that list.
func parametersForLevels(levels []ssot.ParameterLevel) []ssot.ParameterDef {
	if len(levels) == 1 {
		return parametersForLevel(levels[0])
	}
	return ssotdefs.ParameterDefsForLevels(levels...)
}

// atLeastOneFields lists parameter fields first, then the other struct fields.
// Struct order stays the statement order. The validation argument order does not.
func atLeastOneFields(q *QueryStruct, params []ssot.ParameterDef) []string {
	paramNames := make(map[string]struct{}, len(params))
	fields := make([]string, 0, len(params)+len(q.fields))
	for _, param := range params {
		name := param.StructFieldName()
		if _, ok := paramNames[name]; ok {
			continue
		}
		paramNames[name] = struct{}{}
		fields = append(fields, name)
	}
	for _, name := range queryStructFieldNames(q) {
		if _, ok := paramNames[name]; ok {
			continue
		}
		fields = append(fields, name)
	}
	return fields
}

func (v *QueryStruct) withIdentifierParameterValidations(params []ssot.ParameterDef) *QueryStruct {
	for _, p := range params {
		if _, err := ToObjectIdentifierKind(p.Kind); err != nil {
			continue
		}
		v = v.WithValidation(ValidIdentifierIfSet, sqlToFieldName(p.SqlName, true))
	}
	return v
}

func (v *QueryStruct) withAlterIdentifierValidations(obj ssot.Object) *QueryStruct {
	for _, clause := range obj.Alter.Clauses {
		switch clause.Kind {
		case ssot.AlterClauseRenameTo:
			v = v.WithValidation(ValidIdentifierIfSet, "RenameTo")
		case ssot.AlterClauseSwapWith:
			v = v.WithValidation(ValidIdentifierIfSet, "SwapWith")
		}
	}
	return v
}

func queryStructFieldNames(q *QueryStruct) []string {
	return collections.Map(q.fields, func(field *Field) string { return field.Name })
}

func alterExactlyOneFields(obj ssot.Object) []string {
	var fields []string
	for _, clause := range obj.Alter.Clauses {
		switch clause.Kind {
		case ssot.AlterClauseKeyword:
			if clause.NotExactlyOne {
				continue
			}
			fields = append(fields, sqlToFieldName(clause.SQL, true))
		case ssot.AlterClauseList, ssot.AlterClauseText:
			fields = append(fields, sqlToFieldName(clause.SQL, true))
		case ssot.AlterClauseBlock:
			fields = append(fields, clause.Field)
		case ssot.AlterClauseRenameTo:
			fields = append(fields, "RenameTo")
		case ssot.AlterClauseSwapWith:
			fields = append(fields, "SwapWith")
		case ssot.AlterClauseSet:
			fields = append(fields, "Set")
		case ssot.AlterClauseUnset:
			fields = append(fields, "Unset")
		case ssot.AlterClauseSetTags:
			fields = append(fields, "SetTags")
		case ssot.AlterClauseUnsetTags:
			fields = append(fields, "UnsetTags")
		}
	}
	return fields
}

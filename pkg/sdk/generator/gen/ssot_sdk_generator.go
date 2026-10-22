package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

// parametersForObject is the catalog a Parameters entry spreads. Tests replace it with a short list.
var parametersForObject = func(obj ssot.Object) []ssot.ParameterDef {
	return ssotdefs.ParameterDefsForLevel(obj.ParameterLevel())
}

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
// helperStructs are extra query structs the describe operation names, such as a warehouse type variant.
func (i *Interface) DescribeOperationFromObject(obj ssot.Object, helperStructs ...IntoField) *Interface {
	cfg := obj.SdkConfig()
	if !cfg.Describe {
		panic(fmt.Sprintf("object %q is not enabled for SDK DESCRIBE generation", obj.Name))
	}
	return i.DescribeOperationWithPairedStructs(
		describeMappingFromObject(obj, cfg),
		cfg.DescribeDoc,
		DescribePairedStructs(obj),
		describeQueryStructFromConfig(cfg),
		helperStructs...,
	)
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
	cfg := obj.SdkConfig()
	if !cfg.Create {
		panic(fmt.Sprintf("object %q is not enabled for SDK CREATE generation", obj.Name))
	}
	return i.CreateOperation(cfg.CreateDoc, createQueryStructFromConfig(obj, cfg), append(helperStructs, enumListItemStructs(obj)...)...)
}

func createQueryStructFromConfig(obj ssot.Object, cfg ssot.SdkGenerator) *QueryStruct {
	q := NewQueryStruct(cfg.CreateQueryStructName).Create()
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
	createParams := parametersForStatement(obj, obj.Create.Clauses)
	q = q.appendSqlSequence(obj.Create.KeywordsBefore, attributeSqlCreate).SQL(cfg.CreateSql)
	if !cfg.CreateWithoutIfNotExists {
		q = q.IfNotExists()
	}
	q = q.Name().appendSqlSequenceParams(obj.Create.Clauses, attributeSqlCreate, createParams).WithValidation(ValidIdentifier, "name")
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
		i = i.WithCustomInterfaceMethod(spec.Name, obj.CustomInterfaceMethodDoc(spec), customInterfaceMethodParameters(obj, spec), returns...)
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
	cfg := obj.SdkConfig()
	if !cfg.Alter {
		panic(fmt.Sprintf("object %q is not enabled for SDK ALTER generation", obj.Name))
	}
	return i.AlterOperation(cfg.AlterDoc, alterQueryStructFromConfig(obj, cfg)).withEnumListItems(obj)
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
		default:
			panic(fmt.Sprintf("object %q: ALTER clause has kind %q", obj.Name, clause.Kind))
		}
	}
	q = q.WithValidation(ValidIdentifier, "name")
	if !cfg.AlterWithoutExactlyOneValueSet {
		q = q.WithValidation(ExactlyOneValueSet, alterExactlyOneFields(obj)...)
	}
	q = q.withAlterIdentifierValidations(obj)
	return q
}

func appendAlterSet(q *QueryStruct, obj ssot.Object, cfg ssot.SdkGenerator) *QueryStruct {
	setParams := parametersForStatement(obj, obj.Set)
	setStruct := NewQueryStruct(obj.Name+"Set").appendSqlSequenceParams(obj.Set, attributeSqlSet, setParams)
	setStruct = setStruct.withIdentifierParameterValidations(setParams)
	if cfg.AlterSetAdditionalValidations {
		setStruct = setStruct.WithAdditionalValidations()
	}
	setStruct = setStruct.WithValidation(AtLeastOneValueSet, queryStructFieldNames(setStruct)...)
	var setOptions FieldTransformer = KeywordOptions().SQL("SET")
	if cfg.AlterSetList {
		setOptions = ListOptions().NoParentheses().SQL("SET")
	}
	return q.OptionalQueryStructField("Set", setStruct, setOptions)
}

func appendAlterUnset(q *QueryStruct, obj ssot.Object) *QueryStruct {
	unsetParams := parametersForStatement(obj, obj.Unset)
	unsetStruct := NewQueryStruct(obj.Name+"Unset").appendSqlSequenceParams(obj.Unset, attributeSqlUnset, unsetParams)
	unsetStruct = unsetStruct.WithValidation(AtLeastOneValueSet, queryStructFieldNames(unsetStruct)...)
	return q.OptionalQueryStructField("Unset", unsetStruct, ListOptions().NoParentheses().SQL("UNSET"))
}

func parametersForStatement(obj ssot.Object, list []*ssot.Attribute) []ssot.ParameterDef {
	for _, field := range list {
		if field != nil && field.IsParameters() {
			return parametersForObject(obj)
		}
	}
	return nil
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
			fields = append(fields, sqlToFieldName(clause.SQL, true))
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

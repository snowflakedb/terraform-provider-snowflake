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
func (i *Interface) DescribeOperationFromObject(obj ssot.Object) *Interface {
	cfg := obj.SdkConfig()
	if !cfg.Describe {
		panic(fmt.Sprintf("object %q is not enabled for SDK DESCRIBE generation", obj.Name))
	}
	return i.DescribeOperationWithPairedStructs(
		describeMappingFromObject(obj, cfg),
		cfg.DescribeDoc,
		DescribePairedStructs(obj),
		describeQueryStructFromConfig(cfg),
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
// helperStructs are query structs the attribute clauses name, such as a list item.
func (i *Interface) CreateOperationFromObject(obj ssot.Object, helperStructs ...IntoField) *Interface {
	cfg := obj.SdkConfig()
	if !cfg.Create {
		panic(fmt.Sprintf("object %q is not enabled for SDK CREATE generation", obj.Name))
	}
	return i.CreateOperation(cfg.CreateDoc, createQueryStructFromConfig(obj, cfg), helperStructs...)
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

// AlterOperationFromObject adds ALTER from obj.SdkConfig.
// Clauses are Alter, SQL, optional IF EXISTS, Name, then obj.Alter in order, then ValidIdentifier.
// AtLeastOneValueSet on SET and UNSET lists the fields those structs contain.
// ExactlyOneValueSet follows obj.Alter.
func (i *Interface) AlterOperationFromObject(obj ssot.Object) *Interface {
	cfg := obj.SdkConfig()
	if !cfg.Alter {
		panic(fmt.Sprintf("object %q is not enabled for SDK ALTER generation", obj.Name))
	}
	return i.AlterOperation(cfg.AlterDoc, alterQueryStructFromConfig(obj, cfg))
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

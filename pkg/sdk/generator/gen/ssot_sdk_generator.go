package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

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
// Clauses are Create, SQL, optional OR REPLACE, optional IF NOT EXISTS, Name, then obj.Create, then ValidIdentifier.
// helperStructs are query structs the attribute clauses name, such as a list item.
func (i *Interface) CreateOperationFromObject(obj ssot.Object, helperStructs ...IntoField) *Interface {
	cfg := obj.SdkConfig()
	if !cfg.Create {
		panic(fmt.Sprintf("object %q is not enabled for SDK CREATE generation", obj.Name))
	}
	return i.CreateOperation(cfg.CreateDoc, createQueryStructFromConfig(obj, cfg), helperStructs...)
}

func createQueryStructFromConfig(obj ssot.Object, cfg ssot.SdkGenerator) *QueryStruct {
	q := NewQueryStruct(cfg.CreateQueryStructName).Create().SQL(cfg.CreateSql)
	if !cfg.CreateWithoutOrReplace {
		q = q.OrReplace()
	}
	if !cfg.CreateWithoutIfNotExists {
		q = q.IfNotExists()
	}
	q = q.Name().appendSqlSequence(obj.Create, attributeSqlCreate).WithValidation(ValidIdentifier, "name")
	if !cfg.CreateWithoutOrReplace && !cfg.CreateWithoutIfNotExists {
		q = q.WithValidation(ConflictingFields, "OrReplace", "IfNotExists")
	}
	if cfg.CreateAdditionalValidations {
		q = q.WithAdditionalValidations()
	}
	return q
}

// AlterOperationFromObject adds ALTER from obj.SdkConfig.
// Clauses are Alter, SQL, optional IF EXISTS, Name, action keywords, optional SET, optional UNSET, optional tag clauses, ValidIdentifier.
// AtLeastOneValueSet on SET and UNSET lists the fields those structs contain.
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
	// TODO [next PRs]: These keywords are emitted before SET and UNSET. That order matters, and some objects place action keywords after those blocks.
	for _, action := range cfg.AlterActions {
		q = q.OptionalSQL(action)
	}
	if !cfg.AlterWithoutSet {
		setStruct := NewQueryStruct(obj.Name+"Set").appendSqlSequence(obj.Set, attributeSqlSet)
		setFields := queryStructFieldNames(setStruct)
		if cfg.AlterSetAdditionalValidations {
			setStruct = setStruct.WithAdditionalValidations()
		}
		setStruct = setStruct.WithValidation(AtLeastOneValueSet, setFields...)
		q = q.OptionalQueryStructField("Set", setStruct, KeywordOptions().SQL("SET"))
	}
	if !cfg.AlterWithoutUnset {
		unsetStruct := NewQueryStruct(obj.Name+"Unset").appendSqlSequence(obj.Unset, attributeSqlUnset)
		unsetStruct = unsetStruct.WithValidation(AtLeastOneValueSet, queryStructFieldNames(unsetStruct)...)
		q = q.OptionalQueryStructField("Unset", unsetStruct, ListOptions().NoParentheses().SQL("UNSET"))
	}
	if !cfg.AlterWithoutSetTags {
		q = q.OptionalSetTags()
	}
	if !cfg.AlterWithoutUnsetTags {
		q = q.OptionalUnsetTags()
	}
	q = q.WithValidation(ValidIdentifier, "name")
	if !cfg.AlterWithoutExactlyOneValueSet {
		q = q.WithValidation(ExactlyOneValueSet, alterExactlyOneFields(cfg)...)
	}
	return q
}

func queryStructFieldNames(q *QueryStruct) []string {
	return collections.Map(q.fields, func(field *Field) string { return field.Name })
}

func alterExactlyOneFields(cfg ssot.SdkGenerator) []string {
	var fields []string
	for _, action := range cfg.AlterActions {
		fields = append(fields, sqlToFieldName(action, true))
	}
	if !cfg.AlterWithoutSet {
		fields = append(fields, "Set")
	}
	if !cfg.AlterWithoutUnset {
		fields = append(fields, "Unset")
	}
	if !cfg.AlterWithoutSetTags {
		fields = append(fields, "SetTags")
	}
	if !cfg.AlterWithoutUnsetTags {
		fields = append(fields, "UnsetTags")
	}
	return fields
}

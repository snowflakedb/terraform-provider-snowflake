package gen

import (
	"fmt"

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

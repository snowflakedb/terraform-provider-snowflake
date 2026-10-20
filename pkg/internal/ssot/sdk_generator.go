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
)

// SdkShowByIDFilteringKind is one ShowByID filter, in call-argument order.
// Omitted filters mean the SHOW call passes no filtering argument.
type SdkShowByIDFilteringKind string

const (
	SdkShowByIDLike       SdkShowByIDFilteringKind = "like"
	SdkShowByIDExtendedIn SdkShowByIDFilteringKind = "extended_in"
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

func SdkShowOptionalExtendedIn() SdkShowOption {
	return func(s *SdkGenerator) { s.ShowOptionalExtendedIn = true }
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

// SdkDescribeDoc overrides the default DESCRIBE doc URL.
// The default is https://docs.snowflake.com/en/sql-reference/sql/desc- plus the hyphenated Name.
func SdkDescribeDoc(url string) SdkDescribeOption {
	return func(s *SdkGenerator) { s.DescribeDoc = url }
}

// SdkDescribeSql overrides the default DESCRIBE command words.
// The default is the uppercase words of Name, so ComputePool is COMPUTE POOL.
func SdkDescribeSql(sql string) SdkDescribeOption {
	return func(s *SdkGenerator) { s.DescribeSql = sql }
}

// SdkDescribeMapping overrides the default DESCRIBE result shape.
// The default is a single value. Pass SdkDescriptionMappingSlice when DESCRIBE returns a slice.
func SdkDescribeMapping(kind SdkDescriptionMapping) SdkDescribeOption {
	return func(s *SdkGenerator) { s.DescribeMapping = kind }
}

// SdkDescribeQueryStructName overrides the default DESCRIBE query struct name.
// The default is Desc plus Name, so ComputePool is DescComputePool.
func SdkDescribeQueryStructName(name string) SdkDescribeOption {
	return func(s *SdkGenerator) { s.DescribeQueryStructName = name }
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
			cfg.DescribeSql = strings.ToUpper(genhelpers.CamelToWords(name))
		}
		if cfg.DescribeQueryStructName == "" {
			cfg.DescribeQueryStructName = "Desc" + name
		}
		if cfg.DescribeMapping == "" {
			cfg.DescribeMapping = SdkDescriptionMappingSingleValue
		}
	}
	if cfg.Drop {
		name := obj.nameForDrop()
		if cfg.DropDoc == "" {
			cfg.DropDoc = dropDocURLPrefix + docsObjectSlug(name)
		}
		if cfg.DropSql == "" {
			cfg.DropSql = strings.ToUpper(genhelpers.CamelToWords(name))
		}
		if cfg.DropQueryStructName == "" {
			cfg.DropQueryStructName = "Drop" + name
		}
	}
	return cfg
}

func (obj Object) namePluralForShow() string {
	if obj.NamePlural == "" {
		panic(fmt.Sprintf("object %q needs NamePlural to build the SHOW defaults", obj.Name))
	}
	return obj.NamePlural
}

func (obj Object) nameForDescribe() string {
	if obj.Name == "" {
		panic(fmt.Sprintf("object %q needs Name to build the DESCRIBE defaults", obj.Name))
	}
	return obj.Name
}

func (obj Object) nameForDrop() string {
	if obj.Name == "" {
		panic(fmt.Sprintf("object %q needs Name to build the DROP defaults", obj.Name))
	}
	return obj.Name
}

// docsObjectSlug is the docs path segment: snake case with hyphens, so ComputePools is compute-pools.
func docsObjectSlug(name string) string {
	return strings.ReplaceAll(genhelpers.ToSnakeCase(name), "_", "-")
}

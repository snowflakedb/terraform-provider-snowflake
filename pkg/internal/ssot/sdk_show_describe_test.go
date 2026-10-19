package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SdkConfig_fillsDefaults(t *testing.T) {
	obj := Object{
		Name:       "ComputePool",
		NamePlural: "ComputePools",
		Generators: []GeneratorOption{WithSdk(SdkShow(), SdkDescribe())},
	}

	require.Equal(t, SdkGenerator{
		Show:                    true,
		Describe:                true,
		ShowDoc:                 "https://docs.snowflake.com/en/sql-reference/sql/show-compute-pools",
		ShowSql:                 "COMPUTE POOLS",
		ShowQueryStructName:     "ShowComputePools",
		DescribeDoc:             "https://docs.snowflake.com/en/sql-reference/sql/desc-compute-pool",
		DescribeSql:             "COMPUTE POOL",
		DescribeQueryStructName: "DescComputePool",
		DescribeMapping:         SdkDescriptionMappingSingleValue,
	}, obj.SdkConfig())

	raw := ObjectGenerators(obj.Generators...).Sdk
	require.Equal(t, SdkGenerator{Show: true, Describe: true}, raw)
}

func Test_SdkConfig_optionsWin(t *testing.T) {
	obj := Object{
		Name:       "Schema",
		NamePlural: "Schemas",
		Generators: []GeneratorOption{WithSdk(
			SdkShow(
				SdkShowDoc("https://custom/show"),
				SdkShowSql("CUSTOM OBJECTS"),
			),
			SdkDescribe(
				SdkDescribeDoc("https://custom/desc"),
				SdkDescribeSql("CUSTOM OBJECT"),
				SdkDescribeMapping(SdkDescriptionMappingSlice),
				SdkDescribeQueryStructName("DescribeSchema"),
			),
		)},
	}

	cfg := obj.SdkConfig()
	require.Equal(t, "https://custom/show", cfg.ShowDoc)
	require.Equal(t, "CUSTOM OBJECTS", cfg.ShowSql)
	require.Equal(t, "ShowSchemas", cfg.ShowQueryStructName)
	require.Equal(t, "https://custom/desc", cfg.DescribeDoc)
	require.Equal(t, "CUSTOM OBJECT", cfg.DescribeSql)
	require.Equal(t, SdkDescriptionMappingSlice, cfg.DescribeMapping)
	require.Equal(t, "DescribeSchema", cfg.DescribeQueryStructName)
}

func Test_SdkConfig_panicsWithoutName(t *testing.T) {
	require.Panics(t, func() {
		Object{Generators: []GeneratorOption{WithSdk(SdkShow())}}.SdkConfig()
	})
	require.Panics(t, func() {
		Object{Generators: []GeneratorOption{WithSdk(SdkDescribe())}}.SdkConfig()
	})
}

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

func Test_SdkConfig_dropDefaults(t *testing.T) {
	obj := Object{
		Name:       "ComputePool",
		NamePlural: "ComputePools",
		Generators: []GeneratorOption{WithSdk(SdkDrop(SdkDropSafelyHook()))},
	}
	require.Equal(t, SdkGenerator{
		Drop:                true,
		DropDoc:             "https://docs.snowflake.com/en/sql-reference/sql/drop-compute-pool",
		DropSql:             "COMPUTE POOL",
		DropQueryStructName: "DropComputePool",
		DropSafelyHook:      true,
	}, obj.SdkConfig())

	raw := ObjectGenerators(obj.Generators...).Sdk
	require.Equal(t, SdkGenerator{Drop: true, DropSafelyHook: true}, raw)
}

func Test_ObjectType_matchesSdkConfig(t *testing.T) {
	computePool := Object{
		Name:       "ComputePool",
		NamePlural: "ComputePools",
		Generators: []GeneratorOption{WithSdk(SdkShow(), SdkDescribe(), SdkCreate(), SdkAlter(), SdkDrop())},
	}
	require.Equal(t, "COMPUTE POOL", computePool.ObjectType())
	cfg := computePool.SdkConfig()
	require.Equal(t, computePool.ObjectType(), cfg.DescribeSql)
	require.Equal(t, computePool.ObjectType(), cfg.CreateSql)
	require.Equal(t, computePool.ObjectType(), cfg.AlterSql)
	require.Equal(t, computePool.ObjectType(), cfg.DropSql)

	overridden := Object{
		Name: "CatalogIntegrationObjectStorage",
		Generators: []GeneratorOption{WithSdk(
			SdkDescribe(SdkDescribeSql("CATALOG INTEGRATION")),
			SdkCreate(),
			SdkAlter(),
			SdkDrop(),
		)},
	}
	require.Equal(t, "CATALOG INTEGRATION", overridden.ObjectType())
	cfg = overridden.SdkConfig()
	require.Equal(t, overridden.ObjectType(), cfg.DescribeSql)
	require.Equal(t, overridden.ObjectType(), cfg.CreateSql)
	require.Equal(t, overridden.ObjectType(), cfg.AlterSql)
	require.Equal(t, overridden.ObjectType(), cfg.DropSql)
}

func Test_SdkConfig_createAlterDefaults(t *testing.T) {
	obj := Object{
		Name: "ComputePool",
		Generators: []GeneratorOption{WithSdk(
			SdkCreate(SdkCreateWithoutOrReplace(), SdkCreateAdditionalValidations()),
			SdkAlter(),
		)},
	}
	cfg := obj.SdkConfig()
	require.Equal(t, "https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool", cfg.CreateDoc)
	require.Equal(t, "COMPUTE POOL", cfg.CreateSql)
	require.Equal(t, "CreateComputePool", cfg.CreateQueryStructName)
	require.False(t, cfg.CreateWithoutIfNotExists)
	require.True(t, cfg.CreateWithoutOrReplace)
	require.False(t, cfg.AlterWithoutIfExists)
	require.Equal(t, "https://docs.snowflake.com/en/sql-reference/sql/alter-compute-pool", cfg.AlterDoc)
	require.Equal(t, "COMPUTE POOL", cfg.AlterSql)
	require.Equal(t, "AlterComputePool", cfg.AlterQueryStructName)

	overridden := Object{
		Name: "ComputePool",
		Generators: []GeneratorOption{WithSdk(
			SdkCreate(SdkCreateDoc("https://custom/create"), SdkCreateSql("CUSTOM"), SdkCreateQueryStructName("MakePool")),
			SdkAlter(SdkAlterDoc("https://custom/alter"), SdkAlterSql("CUSTOM"), SdkAlterQueryStructName("ChangePool")),
		)},
	}.SdkConfig()
	require.Equal(t, "https://custom/create", overridden.CreateDoc)
	require.Equal(t, "CUSTOM", overridden.CreateSql)
	require.Equal(t, "MakePool", overridden.CreateQueryStructName)
	require.Equal(t, "https://custom/alter", overridden.AlterDoc)
	require.Equal(t, "CUSTOM", overridden.AlterSql)
	require.Equal(t, "ChangePool", overridden.AlterQueryStructName)
}

func Test_SdkConfig_panicsWithoutName(t *testing.T) {
	require.Panics(t, func() {
		Object{Generators: []GeneratorOption{WithSdk(SdkShow())}}.SdkConfig()
	})
	require.Panics(t, func() {
		Object{Generators: []GeneratorOption{WithSdk(SdkDescribe())}}.SdkConfig()
	})
	require.Panics(t, func() {
		Object{Generators: []GeneratorOption{WithSdk(SdkDrop())}}.SdkConfig()
	})
	require.Panics(t, func() {
		Object{Generators: []GeneratorOption{WithSdk(SdkCreate())}}.SdkConfig()
	})
	require.Panics(t, func() {
		Object{Generators: []GeneratorOption{WithSdk(SdkAlter())}}.SdkConfig()
	})
	require.Panics(t, func() {
		Object{Generators: []GeneratorOption{WithSdk(SdkCustomOperations())}}.SdkConfig()
	})
	require.Panics(t, func() {
		Object{Generators: []GeneratorOption{WithSdk(SdkCustomInterfaceMethods())}}.SdkConfig()
	})
}

func Test_SdkConfig_customOperations(t *testing.T) {
	obj := Object{
		Name:       "Fake",
		Generators: []GeneratorOption{WithSdk(SdkCustomOperations())},
		CustomOperations: []CustomOperationSpec{{
			Kind: "Undrop",
		}},
	}
	require.True(t, obj.SdkConfig().CustomOperations)
	require.Equal(t, SdkGenerator{CustomOperations: true}, ObjectGenerators(obj.Generators...).Sdk)
	require.Equal(t, "https://docs.snowflake.com/en/sql-reference/sql/undrop-fake", obj.CustomOperationDoc(obj.CustomOperations[0]))
	require.Equal(t, "FAKE", obj.CustomOperationSQL())
	require.Equal(t, "UndropFake", obj.CustomOperationQueryStructName("Undrop"))
	require.Equal(t, "https://docs.snowflake.com/en/sql-reference/sql/create-fake", obj.CustomOperationDoc(CustomOperationSpec{
		Kind: "Clone",
		Doc:  "https://docs.snowflake.com/en/sql-reference/sql/create-fake",
	}))

	require.Panics(t, func() {
		Object{
			Name:             "Fake",
			Generators:       []GeneratorOption{WithSdk(SdkCustomOperations())},
			CustomOperations: []CustomOperationSpec{{}},
		}.SdkConfig()
	})
}

func Test_SdkConfig_customInterfaceMethods(t *testing.T) {
	obj := Object{
		Name:       "Fake",
		IdType:     "DatabaseObjectIdentifier",
		Generators: []GeneratorOption{WithSdk(SdkCustomInterfaceMethods())},
		CustomInterfaceMethods: []CustomInterfaceMethodSpec{{
			Name: "Use",
		}},
	}
	require.True(t, obj.SdkConfig().CustomInterfaceMethods)
	require.Equal(t, SdkGenerator{CustomInterfaceMethods: true}, ObjectGenerators(obj.Generators...).Sdk)
	require.Empty(t, obj.CustomInterfaceMethods[0].Doc)
}

func Test_VariantCreateConfig_usesEntityName(t *testing.T) {
	base := Object{Name: "Warehouse"}
	comment := Attribute{Name: "comment", Type: AttrString}
	obj := Object{
		Name:                       "WarehouseAdaptive",
		VariantOf:                  &base,
		SdkCreateRequestObjectName: "AdaptiveWarehouse",
		Create: CreateSpec{
			Doc:     "https://docs.example.com/adaptive",
			Clauses: []*Attribute{&comment},
		},
	}
	require.Equal(t, "WAREHOUSE", obj.VariantCreateConfig().CreateSql)

	obj.VariantEntityCreateName = "INTERACTIVE WAREHOUSE"
	require.Equal(t, SdkGenerator{
		CreateDoc:             "https://docs.example.com/adaptive",
		CreateSql:             "INTERACTIVE WAREHOUSE",
		CreateQueryStructName: "CreateAdaptiveWarehouse",
	}, obj.VariantCreateConfig())
}

func Test_VariantAlterConfig_usesEntityName(t *testing.T) {
	base := Object{Name: "Warehouse"}
	obj := Object{
		Name:                       "WarehouseAdaptive",
		VariantOf:                  &base,
		SdkCreateRequestObjectName: "AdaptiveWarehouse",
		Alter: AlterSpec{
			Doc:     "https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse",
			Clauses: []AlterClause{AlterSet()},
		},
	}
	require.Equal(t, SdkGenerator{
		AlterDoc:             "https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse",
		AlterSql:             "WAREHOUSE",
		AlterQueryStructName: "AlterAdaptiveWarehouse",
	}, obj.VariantAlterConfig())

	obj.VariantEntityAlterName = "CUSTOM WAREHOUSE"
	require.Equal(t, "CUSTOM WAREHOUSE", obj.VariantAlterConfig().AlterSql)
}

func Test_SdkConfig_describeMappingDefault(t *testing.T) {
	row := Attribute{Name: "property", Type: AttrString, Readonly: true}
	obj := Object{Name: "Widget", NamePlural: "Widgets", Generators: []GeneratorOption{WithSdk(SdkDescribe())}}
	require.Equal(t, SdkDescriptionMappingSingleValue, obj.SdkConfig().DescribeMapping)

	obj.DescribeRows = OutputSpec{Columns: []*Attribute{&row}}
	require.Equal(t, SdkDescriptionMappingSlice, obj.SdkConfig().DescribeMapping)

	obj.Generators = []GeneratorOption{WithSdk(SdkDescribe(SdkDescribeMapping(SdkDescriptionMappingSingleValue)))}
	require.Equal(t, SdkDescriptionMappingSingleValue, obj.SdkConfig().DescribeMapping)
}

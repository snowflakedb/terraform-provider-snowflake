package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ValidateObject_writableMayBeOmittedFromColumns(t *testing.T) {
	minNodes := Attribute{Name: "min_nodes"}
	initiallySuspended := Attribute{Name: "initially_suspended"}
	name := Attribute{Name: "name", Readonly: true}
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:       "Pool",
			Attributes: []Attribute{minNodes, initiallySuspended},
			Show: OutputSpec{
				Columns: []*Attribute{&name, &minNodes},
			},
		})
	})
}

func Test_ValidateObject_readonlyInAttributes(t *testing.T) {
	require.Panics(t, func() {
		ValidateObject(Object{
			Name:       "Pool",
			Attributes: []Attribute{{Name: "name", Readonly: true}},
		})
	})
}

func Test_ValidateObject_writableColumnNotInAttributes(t *testing.T) {
	comment := Attribute{Name: "comment"}
	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Pool",
			Show: OutputSpec{Columns: []*Attribute{&comment}},
		})
	})
}

func Test_ValidateObject_identifierRequiresType(t *testing.T) {
	owner := Attribute{Name: "owner", Type: AttrIdentifier}
	require.Panics(t, func() {
		ValidateObject(Object{
			Name:       "Pool",
			Attributes: []Attribute{owner},
		})
	})
	refs := Attribute{Name: "refs", Type: AttrList, Elem: &Attribute{Type: AttrIdentifier}}
	require.Panics(t, func() {
		ValidateObject(Object{
			Name:       "Pool",
			Attributes: []Attribute{refs},
		})
	})
	state := Attribute{Name: "state", Type: AttrEnum, Readonly: true}
	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Pool",
			Show: OutputSpec{Columns: []*Attribute{&state}},
		})
	})
}

func Test_ValidateObject_conflictsWithUnknownAttribute(t *testing.T) {
	require.PanicsWithValue(t, `object "Warehouse": attribute "resource_constraint" conflicts with unknown attribute "generation"`, func() {
		ValidateObject(Object{
			Name: "Warehouse",
			Attributes: []Attribute{{
				Name:          "resource_constraint",
				ConflictsWith: []Attribute{{Name: "generation"}},
			}},
		})
	})
}

func Test_ValidateObject_objectAssertSkipMustBeColumn(t *testing.T) {
	tables := Attribute{Name: "tables", Type: AttrString, Readonly: true}
	other := Attribute{Name: "other", Type: AttrString, Readonly: true}
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name: "Warehouse",
			Generators: []GeneratorOption{
				WithObjectAssert(ObjectAssertShow(ObjectAssertShowSkipFields(&tables))),
			},
			Show: OutputSpec{Columns: []*Attribute{&tables}},
		})
	})
	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Warehouse",
			Generators: []GeneratorOption{
				WithObjectAssert(ObjectAssertShow(ObjectAssertShowSkipFields(&other))),
			},
			Show: OutputSpec{Columns: []*Attribute{&tables}},
		})
	})
	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Warehouse",
			Generators: []GeneratorOption{
				WithObjectAssert(ObjectAssertDescribe(ObjectAssertDescribeSkipFields(&tables))),
			},
			Show: OutputSpec{Columns: []*Attribute{&tables}},
		})
	})
}

func Test_ValidateObject_conflictsWithItself(t *testing.T) {
	require.PanicsWithValue(t, `object "Warehouse": attribute "generation" conflicts with itself`, func() {
		ValidateObject(Object{
			Name: "Warehouse",
			Attributes: []Attribute{{
				Name:          "generation",
				ConflictsWith: []Attribute{{Name: "generation"}},
			}},
		})
	})
}

func Test_ValidateObject_conflictsWithKnownAttribute(t *testing.T) {
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name: "Warehouse",
			Attributes: []Attribute{
				{Name: "resource_constraint", ConflictsWith: []Attribute{{Name: "generation"}}},
				{Name: "generation", ConflictsWith: []Attribute{{Name: "resource_constraint"}}},
			},
		})
	})
}

func Test_ValidateObject_skipCreateBlockMustBeListBlock(t *testing.T) {
	require.PanicsWithValue(t, `object "Alert": SkipCreateBlock attribute "comment" must be a list block`, func() {
		ValidateObject(Object{
			Name:       "Alert",
			Attributes: []Attribute{{Name: "comment", Type: AttrString, SkipCreateBlock: true}},
		})
	})
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name: "Alert",
			Attributes: []Attribute{{
				Name:            "alert_schedule",
				Type:            AttrList,
				SkipCreateBlock: true,
				Children:        []Attribute{{Name: "interval", Type: AttrInt}},
			}},
		})
	})
}

func Test_ValidateObject_createOnlyMustBeOptionalBool(t *testing.T) {
	require.PanicsWithValue(t, `object "Warehouse": CreateOnly attribute "auto_suspend" must be an optional bool`, func() {
		ValidateObject(Object{
			Name:       "Warehouse",
			Attributes: []Attribute{{Name: "auto_suspend", Type: AttrInt, CreateOnly: true}},
		})
	})
	require.PanicsWithValue(t, `object "Warehouse": CreateOnly attribute "initially_suspended" must be an optional bool`, func() {
		ValidateObject(Object{
			Name:       "Warehouse",
			Attributes: []Attribute{{Name: "initially_suspended", Type: AttrBool, Required: true, CreateOnly: true}},
		})
	})
	require.PanicsWithValue(t, `object "Warehouse": NativeBool attribute "comment" must be a bool`, func() {
		ValidateObject(Object{
			Name:       "Warehouse",
			Attributes: []Attribute{{Name: "comment", Type: AttrString, NativeBool: true}},
		})
	})
}

func Test_ValidateObject_additionalCustomizeDiffs(t *testing.T) {
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name: "Widget",
			AdditionalCustomizeDiffs: []CustomizeDiffFunc{
				"widgetForceNewIfSizeRemoved",
			},
			ResourceHooks: []ResourceHookOption{
				WithCustomizeDiff(Generated()),
			},
		})
	})
}

func Test_ValidateObject_additionalCustomizeDiffsEmptyEntry(t *testing.T) {
	require.PanicsWithValue(t, `object "Widget": AdditionalCustomizeDiffs entry is empty`, func() {
		ValidateObject(Object{
			Name: "Widget",
			AdditionalCustomizeDiffs: []CustomizeDiffFunc{
				"",
			},
			ResourceHooks: []ResourceHookOption{
				WithCustomizeDiff(Generated()),
			},
		})
	})
}

func Test_ValidateObject_additionalCustomizeDiffsNotAFunctionName(t *testing.T) {
	require.PanicsWithValue(t, `object "Widget": AdditionalCustomizeDiffs entry "not a func" is not a function name`, func() {
		ValidateObject(Object{
			Name: "Widget",
			AdditionalCustomizeDiffs: []CustomizeDiffFunc{
				"not a func",
			},
			ResourceHooks: []ResourceHookOption{
				WithCustomizeDiff(Generated()),
			},
		})
	})
}

func Test_ValidateObject_additionalCustomizeDiffsRequireGenerated(t *testing.T) {
	require.PanicsWithValue(t, `object "Widget": AdditionalCustomizeDiffs require WithCustomizeDiff(Generated())`, func() {
		ValidateObject(Object{
			Name: "Widget",
			AdditionalCustomizeDiffs: []CustomizeDiffFunc{
				"widgetForceNewIfSizeRemoved",
			},
		})
	})
}

func Test_ValidateObject_schemaStateUpgraders(t *testing.T) {
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name: "Widget",
			SchemaStateUpgraders: []StateUpgraderFunc{
				"v092WidgetStateUpgrader",
				"v2_6_0_WidgetStateUpgrader",
			},
			ResourceHooks: []ResourceHookOption{
				WithStateUpgraders(Generated()),
			},
		})
	})
}

func Test_ValidateObject_schemaStateUpgradersEmptyEntry(t *testing.T) {
	require.PanicsWithValue(t, `object "Widget": SchemaStateUpgraders entry is empty`, func() {
		ValidateObject(Object{
			Name: "Widget",
			SchemaStateUpgraders: []StateUpgraderFunc{
				"",
			},
			ResourceHooks: []ResourceHookOption{
				WithStateUpgraders(Generated()),
			},
		})
	})
}

func Test_ValidateObject_schemaStateUpgradersNotAFunctionName(t *testing.T) {
	require.PanicsWithValue(t, `object "Widget": SchemaStateUpgraders entry "not a func" is not a function name`, func() {
		ValidateObject(Object{
			Name: "Widget",
			SchemaStateUpgraders: []StateUpgraderFunc{
				"not a func",
			},
			ResourceHooks: []ResourceHookOption{
				WithStateUpgraders(Generated()),
			},
		})
	})
}

func Test_ValidateObject_schemaStateUpgradersRequireGenerated(t *testing.T) {
	require.PanicsWithValue(t, `object "Widget": SchemaStateUpgraders require WithStateUpgraders(Generated())`, func() {
		ValidateObject(Object{
			Name: "Widget",
			SchemaStateUpgraders: []StateUpgraderFunc{
				"v092WidgetStateUpgrader",
			},
		})
	})
}

func Test_ValidateObject_stateUpgradersGenerationNeedsFunctions(t *testing.T) {
	require.PanicsWithValue(t, `object "Widget": WithStateUpgraders(Generated()) needs upgrade functions`, func() {
		ValidateObject(Object{
			Name: "Widget",
			ResourceHooks: []ResourceHookOption{
				WithStateUpgraders(Generated()),
			},
		})
	})
}

func Test_ValidateObject_sqlClause(t *testing.T) {
	minNodes := Attribute{
		Name:      "min_nodes",
		Type:      AttrInt,
		CreateSql: IntParameter(),
		SetSql:    IntParameter(Optional()),
	}
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:       "Pool",
			Attributes: []Attribute{minNodes},
			Create:     CreateSpec{Clauses: []*Attribute{&minNodes}},
			Set:        []*Attribute{&minNodes},
		})
	})

	plain := Attribute{Name: "count", Type: AttrInt, Required: true}
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:       "Widget",
			Attributes: []Attribute{plain},
			Create:     CreateSpec{Clauses: []*Attribute{&plain}},
			Set:        []*Attribute{&plain},
			Unset:      []*Attribute{&plain},
		})
	})
	byValue := plain
	byValue.UnsetSql = IntParameter()
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:       "Widget",
			Attributes: []Attribute{byValue},
			Create:     CreateSpec{Clauses: []*Attribute{&byValue}},
			Unset:      []*Attribute{&byValue},
		})
	})

	transient := Attribute{Name: "is_transient", Type: AttrBool, CreateSql: BoolKeyword(Sql("TRANSIENT"))}
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:       "Widget",
			Attributes: []Attribute{transient, minNodes},
			Create: CreateSpec{
				KeywordsBefore: []*Attribute{&transient},
				Clauses:        []*Attribute{&minNodes},
			},
			Set: []*Attribute{&minNodes},
		})
	})
	require.PanicsWithValue(t, `object "Widget": KeywordsBefore list contains tags`, func() {
		ValidateObject(Object{
			Name:       "Widget",
			Attributes: []Attribute{minNodes},
			Create: CreateSpec{
				KeywordsBefore: []*Attribute{Tags},
				Clauses:        []*Attribute{&minNodes},
			},
			Set: []*Attribute{&minNodes},
		})
	})
	comment := Attribute{
		Name:      "comment",
		Type:      AttrString,
		CreateSql: StringParameter(),
		SetSql:    StringParameter(),
		UnsetSql:  BoolKeyword(),
	}
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:       "Widget",
			Attributes: []Attribute{comment},
			Create:     CreateSpec{Clauses: []*Attribute{Parameters, &comment}},
			Set:        []*Attribute{Parameters, &comment},
			Unset:      []*Attribute{Parameters, &comment},
		})
	})
	require.PanicsWithValue(t, `object "Widget": parameters is not an attribute`, func() {
		ValidateObject(Object{
			Name:       "Widget",
			Attributes: []Attribute{*Parameters},
		})
	})
	require.PanicsWithValue(t, `object "Widget": KeywordsBefore list contains parameters`, func() {
		ValidateObject(Object{
			Name:       "Widget",
			Attributes: []Attribute{comment},
			Create:     CreateSpec{KeywordsBefore: []*Attribute{Parameters}, Clauses: []*Attribute{&comment}},
		})
	})
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:  "Widget",
			Alter: AlterSpec{Clauses: []AlterClause{AlterSet(), AlterUnsetTags(), AlterKeyword("ENABLE MANAGED ACCESS")}},
		})
	})
	require.PanicsWithValue(t, `object "Widget": ALTER keyword is empty`, func() {
		ValidateObject(Object{
			Name:  "Widget",
			Alter: AlterSpec{Clauses: []AlterClause{{Kind: AlterClauseKeyword}}},
		})
	})
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:  "Widget",
			Alter: AlterSpec{Clauses: []AlterClause{AlterRenameTo(), AlterSwapWith()}},
		})
	})
	require.PanicsWithValue(t, `object "Widget": duplicate ALTER rename_to`, func() {
		ValidateObject(Object{
			Name:  "Widget",
			Alter: AlterSpec{Clauses: []AlterClause{AlterRenameTo(), AlterRenameTo()}},
		})
	})
	require.PanicsWithValue(t, `object "Widget": duplicate ALTER set`, func() {
		ValidateObject(Object{
			Name:  "Widget",
			Alter: AlterSpec{Clauses: []AlterClause{AlterSet(), AlterSet()}},
		})
	})
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name: "Widget",
			Alter: AlterSpec{Clauses: []AlterClause{
				AlterKeyword("RESUME"),
				AlterKeywordNotExactlyOne("IF SUSPENDED"),
				AlterList("ADD TABLES", "SchemaObjectIdentifier"),
			}},
		})
	})
	require.PanicsWithValue(t, `object "Widget": ALTER list "ADD TABLES" has no item type`, func() {
		ValidateObject(Object{
			Name:  "Widget",
			Alter: AlterSpec{Clauses: []AlterClause{AlterList("ADD TABLES", "")}},
		})
	})
	require.PanicsWithValue(t, `object "Widget": duplicate ALTER list "ADD TABLES"`, func() {
		ValidateObject(Object{
			Name: "Widget",
			Alter: AlterSpec{Clauses: []AlterClause{
				AlterList("ADD TABLES", "SchemaObjectIdentifier"),
				AlterList("ADD TABLES", "SchemaObjectIdentifier"),
			}},
		})
	})
	require.PanicsWithValue(t, `object "Widget": ALTER list "ADD TABLES" is an exactly-one branch`, func() {
		ValidateObject(Object{
			Name:  "Widget",
			Alter: AlterSpec{Clauses: []AlterClause{{Kind: AlterClauseList, SQL: "ADD TABLES", ListItem: "SchemaObjectIdentifier", NotExactlyOne: true}}},
		})
	})

	require.PanicsWithValue(t, `object "Widget": parameters is not a SHOW column`, func() {
		ValidateObject(Object{
			Name: "Widget",
			Show: OutputSpec{Columns: []*Attribute{Parameters}},
		})
	})

	require.PanicsWithValue(t, `object "Widget": duplicate CREATE field "min_nodes"`, func() {
		ValidateObject(Object{
			Name:       "Widget",
			Attributes: []Attribute{minNodes},
			Create: CreateSpec{
				KeywordsBefore: []*Attribute{&minNodes},
				Clauses:        []*Attribute{&minNodes},
			},
			Set: []*Attribute{&minNodes},
		})
	})

	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Pool",
			Attributes: []Attribute{{
				Name:      "created_on",
				Type:      AttrTime,
				CreateSql: IntParameter(),
			}},
		})
	})

	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Pool",
			Attributes: []Attribute{{
				Name:      "comment",
				Type:      AttrString,
				CreateSql: IdentifierKeyword("AccountObjectIdentifier"),
			}},
		})
	})

	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Pool",
			Attributes: []Attribute{{
				Name:      "instance_family",
				Type:      AttrEnum,
				CreateSql: EnumParameter(),
			}},
		})
	})

	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Pool",
			Attributes: []Attribute{{
				Name:      "backup",
				Type:      AttrList,
				CreateSql: ListParameter(""),
			}},
		})
	})

	require.Panics(t, func() {
		ValidateObject(Object{
			Name: "Pool",
			Attributes: []Attribute{{
				Name:      "for_application",
				Type:      AttrIdentifier,
				CreateSql: IdentifierKeyword(""),
			}},
		})
	})
}

func Test_ValidateObject_sqlOnlyField(t *testing.T) {
	wait := Attribute{Name: "wait_for_completion", Type: AttrBool}
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:  "Widget",
			Set:   []*Attribute{&wait},
			Unset: []*Attribute{&wait},
		})
	})
}

func Test_ValidateObject_nestedAttribute(t *testing.T) {
	leaf := Attribute{Name: "leaf", Type: AttrString}
	inner := Attribute{Name: "inner", Type: AttrList, Children: []Attribute{leaf}}
	outer := Attribute{Name: "outer", Type: AttrSet, Children: []Attribute{inner}}
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:       "Widget",
			Attributes: []Attribute{outer},
		})
	})

	require.PanicsWithValue(t, `object "Widget": attribute "label" Nested requires a list or set`, func() {
		ValidateObject(Object{
			Name:       "Widget",
			Attributes: []Attribute{{Name: "label", Type: AttrString, Children: []Attribute{leaf}}},
		})
	})

	require.PanicsWithValue(t, `object "Widget": attribute "items" has both Elem and Nested`, func() {
		ValidateObject(Object{
			Name: "Widget",
			Attributes: []Attribute{{
				Name:     "items",
				Type:     AttrSet,
				Elem:     &Attribute{Type: AttrString},
				Children: []Attribute{leaf},
			}},
		})
	})

	require.PanicsWithValue(t, `object "Widget": attribute "outer" has duplicate nested attribute "inner"`, func() {
		ValidateObject(Object{
			Name: "Widget",
			Attributes: []Attribute{{
				Name:     "outer",
				Type:     AttrSet,
				Children: []Attribute{inner, {Name: "inner", Type: AttrString}},
			}},
		})
	})

	require.PanicsWithValue(t, `object "Widget": nested attribute "outer"."inner" has a CREATE clause`, func() {
		ValidateObject(Object{
			Name: "Widget",
			Attributes: []Attribute{{
				Name: "outer",
				Type: AttrSet,
				Children: []Attribute{{
					Name:      "inner",
					Type:      AttrString,
					CreateSql: StringParameter(),
				}},
			}},
		})
	})

	require.PanicsWithValue(t, `object "Widget": nested attribute "inner"."leaf" has a SET clause`, func() {
		ValidateObject(Object{
			Name: "Widget",
			Attributes: []Attribute{{
				Name: "outer",
				Type: AttrSet,
				Children: []Attribute{{
					Name: "inner",
					Type: AttrList,
					Children: []Attribute{{
						Name:   "leaf",
						Type:   AttrString,
						SetSql: StringParameter(),
					}},
				}},
			}},
		})
	})

	column := outer
	require.PanicsWithValue(t, `object "Widget": nested attribute "outer" is not a SHOW column`, func() {
		ValidateObject(Object{
			Name:       "Widget",
			Attributes: []Attribute{outer},
			Show:       OutputSpec{Columns: []*Attribute{&column}},
		})
	})
}

func Test_ValidateObject_collectionItemBounds(t *testing.T) {
	locations := Attribute{
		Name:     "storage_location",
		Type:     AttrList,
		Required: true,
		MinItems: new(1),
		Children: []Attribute{
			{Name: "storage_location_name", Type: AttrString, Required: true},
			{Name: "storage_base_url", Type: AttrString, Required: true},
			{Name: "allowed_locations", Type: AttrList, MinItems: new(1), Elem: &Attribute{Type: AttrString}},
		},
	}
	require.NotPanics(t, func() {
		ValidateObject(Object{Name: "ExternalVolume", Attributes: []Attribute{locations}})
	})

	require.PanicsWithValue(t, `object "ExternalVolume": attribute "comment" item bounds require a list or set`, func() {
		ValidateObject(Object{
			Name:       "ExternalVolume",
			Attributes: []Attribute{{Name: "comment", Type: AttrString, MinItems: new(1)}},
		})
	})

	require.PanicsWithValue(t, `object "ExternalVolume": attribute "storage_location" MinItems 2 is greater than MaxItems 1`, func() {
		ValidateObject(Object{
			Name: "ExternalVolume",
			Attributes: []Attribute{{
				Name:     "storage_location",
				Type:     AttrList,
				MinItems: new(2),
				MaxItems: new(1),
			}},
		})
	})

	require.PanicsWithValue(t, `object "ExternalVolume": attribute "storage_base_url" item bounds require a list or set`, func() {
		ValidateObject(Object{
			Name: "ExternalVolume",
			Attributes: []Attribute{{
				Name: "storage_location",
				Type: AttrSet,
				Children: []Attribute{{
					Name:     "storage_base_url",
					Type:     AttrString,
					MinItems: new(1),
				}},
			}},
		})
	})
}

func Test_ValidateObject_customOperations(t *testing.T) {
	transient := Attribute{Name: "is_transient", Type: AttrBool, CreateSql: BoolKeyword(Sql("TRANSIENT"))}
	valid := Object{
		Name:       "Fake",
		Attributes: []Attribute{transient},
		Create:     CreateSpec{KeywordsBefore: []*Attribute{&transient}},
		CustomOperations: []CustomOperationSpec{
			{
				Kind:                  "Clone",
				Doc:                   "https://docs.snowflake.com/en/sql-reference/sql/create-fake",
				AdditionalValidations: true,
				Clauses: []CustomOperationClause{
					CustomStaticSQL("CREATE"),
					CustomOrReplace(),
					CustomAttribute(&transient),
					CustomObjectWords(),
					CustomIfNotExists(),
					CustomName(),
					CustomPredefined("Clone"),
				},
			},
			{
				Kind: "Undrop",
				Clauses: []CustomOperationClause{
					CustomStaticSQL("UNDROP"),
					CustomObjectWords(),
					CustomName(),
				},
			},
		},
	}
	require.NotPanics(t, func() { ValidateObject(valid) })

	missing := valid
	other := Attribute{Name: "other", Type: AttrBool, CreateSql: BoolKeyword(Sql("OTHER"))}
	missing.CustomOperations = []CustomOperationSpec{{
		Kind:    "Clone",
		Clauses: []CustomOperationClause{CustomAttribute(&other)},
	}}
	require.Panics(t, func() { ValidateObject(missing) })

	duplicated := valid
	duplicated.CustomOperations = append(duplicated.CustomOperations, CustomOperationSpec{Kind: "Clone"})
	require.Panics(t, func() { ValidateObject(duplicated) })
}

func Test_ValidateObject_customInterfaceMethods(t *testing.T) {
	valid := Object{
		Name:   "Fake",
		IdType: "DatabaseObjectIdentifier",
		CustomInterfaceMethods: []CustomInterfaceMethodSpec{{
			Name: "Use",
		}},
	}
	require.NotPanics(t, func() { ValidateObject(valid) })

	missingID := valid
	missingID.IdType = ""
	require.Panics(t, func() { ValidateObject(missingID) })

	duplicated := valid
	duplicated.CustomInterfaceMethods = append(duplicated.CustomInterfaceMethods, CustomInterfaceMethodSpec{Name: "Use"})
	require.Panics(t, func() { ValidateObject(duplicated) })
}

func Test_ValidateObject_duplicateOutputKey(t *testing.T) {
	a := Attribute{Name: "for_application", ShowName: "application"}
	b := Attribute{Name: "application"}
	require.Panics(t, func() {
		ValidateObject(Object{
			Name:       "Pool",
			Attributes: []Attribute{a, b},
			Show:       OutputSpec{Columns: []*Attribute{&a, &b}},
		})
	})
}

func Test_ValidateObject_diffSuppressStrategies(t *testing.T) {
	for _, attr := range []Attribute{
		{Name: "text", Type: AttrString, DiffSuppress: DiffSuppressKind("unknown")},
		{Name: "text", Type: AttrString, DiffSuppress: SuppressStatement.WithKey("other")},
		{Name: "list", Type: AttrList, Elem: &Attribute{Type: AttrString, DiffSuppress: SuppressShow}},
	} {
		require.Panics(t, func() { ValidateObject(Object{Name: "Test", Attributes: []Attribute{attr}}) })
	}
}

func Test_ValidateObject_diffSuppressComposition(t *testing.T) {
	for _, attr := range []Attribute{
		{Name: "text", Type: AttrString, DiffSuppress: SuppressIfAny(SuppressStatement, SuppressDescribe.WithKey("definition"))},
		{Name: "text", Type: AttrString, DiffSuppress: SuppressIfAny()},
		{Name: "list", Type: AttrList, Elem: &Attribute{Type: AttrString, DiffSuppress: SuppressStatement}},
	} {
		require.NotPanics(t, func() { ValidateObject(Object{Name: "Test", Attributes: []Attribute{attr}}) })
	}
	for _, attr := range []Attribute{
		{Name: "text", Type: AttrString, DiffSuppress: SuppressIfAny(SuppressShow, DiffSuppressKind("unknown"))},
		{Name: "count", Type: AttrInt, DiffSuppress: SuppressIfAny(SuppressShow, SuppressStatement)},
		{Name: "list", Type: AttrList, Elem: &Attribute{Type: AttrString, DiffSuppress: SuppressIfAny(SuppressStatement, SuppressIfAny(SuppressDescribe))}},
	} {
		require.Panics(t, func() { ValidateObject(Object{Name: "Test", Attributes: []Attribute{attr}}) })
	}
}

func Test_ValidateObject_sqlStatement(t *testing.T) {
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:       "Alert",
			Attributes: []Attribute{{Name: "condition", Type: AttrString, DiffSuppress: SuppressStatement}},
		})
	})
	require.PanicsWithValue(t, `object "Alert": invalid statement suppression for attribute "interval"`, func() {
		ValidateObject(Object{
			Name:       "Alert",
			Attributes: []Attribute{{Name: "interval", Type: AttrInt, DiffSuppress: SuppressStatement}},
		})
	})
}

func Test_ValidateObject_variantEntityCreateName(t *testing.T) {
	base := Object{Name: "Warehouse"}
	kind := Attribute{
		Name:      "warehouse_type_static",
		SqlOnly:   true,
		CreateSql: StaticSQL("warehouseType", "WAREHOUSE_TYPE = 'ADAPTIVE'"),
	}
	comment := Attribute{Name: "comment", Type: AttrString}
	obj := Object{
		Name:                    "WarehouseAdaptive",
		VariantOf:               &base,
		VariantEntityCreateName: "WAREHOUSE",
		Attributes:              []Attribute{kind, comment},
		Create:                  CreateSpec{Clauses: []*Attribute{&kind, &comment}},
	}
	require.NotPanics(t, func() { ValidateObject(obj) })
	require.Equal(t, []Attribute{comment}, obj.ResourceAttributes())

	missing := obj
	missing.VariantEntityCreateName = ""
	require.PanicsWithValue(t, `object "WarehouseAdaptive": variant CREATE needs VariantEntityCreateName`, func() {
		ValidateObject(missing)
	})

	withSql := obj
	withSql.Create.Sql = "WAREHOUSE"
	require.PanicsWithValue(t, `object "WarehouseAdaptive": variant CREATE words are VariantEntityCreateName`, func() {
		ValidateObject(withSql)
	})

	require.PanicsWithValue(t, `object "Widget": VariantEntityCreateName is for a variant`, func() {
		ValidateObject(Object{Name: "Widget", VariantEntityCreateName: "WIDGET"})
	})
}

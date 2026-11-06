package ssot

import (
	"fmt"
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

func Test_ValidateObject_atLeastOneOfAndNestedRelations(t *testing.T) {
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name: "AuthenticationPolicy",
			Attributes: []Attribute{{
				Name:     "mfa_policy",
				Type:     AttrList,
				MaxItems: new(1),
				Elem: ObjectElem([]Attribute{
					{Name: "enforce_mfa_on_external_authentication", AtLeastOneOf: []Attribute{{Name: "enforce_mfa_on_external_authentication"}, {Name: "allowed_methods"}}},
					{Name: "allowed_methods", AtLeastOneOf: []Attribute{{Name: "enforce_mfa_on_external_authentication"}, {Name: "allowed_methods"}}},
				}),
			}},
		})
	})
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name: "AuthenticationPolicy",
			Attributes: []Attribute{{
				Name: "mfa_policy",
				Type: AttrObject,
				Children: []Attribute{
					{Name: "enforce_mfa_on_external_authentication", AtLeastOneOf: []Attribute{{Name: "enforce_mfa_on_external_authentication"}, {Name: "allowed_methods"}}},
					{Name: "allowed_methods", AtLeastOneOf: []Attribute{{Name: "enforce_mfa_on_external_authentication"}, {Name: "allowed_methods"}}},
				},
			}},
		})
	})
	require.PanicsWithValue(t, `object "AuthenticationPolicy": attribute "allowed_methods" at least one of unknown attribute "missing"`, func() {
		ValidateObject(Object{
			Name: "AuthenticationPolicy",
			Attributes: []Attribute{{
				Name:     "mfa_policy",
				Type:     AttrList,
				MaxItems: new(1),
				Elem: ObjectElem([]Attribute{
					{Name: "allowed_methods", AtLeastOneOf: []Attribute{{Name: "missing"}}},
				}),
			}},
		})
	})
	require.PanicsWithValue(t, `object "Widget": attribute "items" relations require a single block`, func() {
		ValidateObject(Object{
			Name: "Widget",
			Attributes: []Attribute{{
				Name: "items",
				Type: AttrSet,
				Elem: ObjectElem([]Attribute{
					{Name: "left", ConflictsWith: []Attribute{{Name: "right"}}},
					{Name: "right"},
				}),
			}},
		})
	})
	require.PanicsWithValue(t, `object "Widget": ConflictsWith attribute "left" is required`, func() {
		ValidateObject(Object{
			Name: "Widget",
			Attributes: []Attribute{
				{Name: "left", Required: true, ConflictsWith: []Attribute{{Name: "right"}}},
				{Name: "right"},
			},
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

func Test_ValidateObject_skipCreate(t *testing.T) {
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name: "Alert",
			Attributes: []Attribute{{
				Name:       "enabled",
				Type:       AttrBool,
				NativeBool: true,
				SkipCreate: true,
			}, {
				Name:       "alert_schedule",
				Type:       AttrList,
				SkipCreate: true,
				Elem:       ObjectElem([]Attribute{{Name: "interval", Type: AttrInt}}),
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
	require.PanicsWithValue(t, `object "Widget": ALTER text is empty`, func() {
		ValidateObject(Object{Name: "Widget", Alter: AlterSpec{Clauses: []AlterClause{AlterText("")}}})
	})
	require.PanicsWithValue(t, `object "Widget": ALTER block needs a field, SQL, and an attribute`, func() {
		ValidateObject(Object{Name: "Widget", Alter: AlterSpec{Clauses: []AlterClause{AlterParameterBlock("Add", "ADD X", nil)}}})
	})
	require.PanicsWithValue(t, `object "Widget": ALTER block "ADD X" needs a list of objects or an object with Children`, func() {
		scalar := Attribute{Name: "x", Type: AttrString}
		ValidateObject(Object{Name: "Widget", Alter: AlterSpec{Clauses: []AlterClause{AlterParameterBlock("Add", "ADD X", &scalar)}}})
	})
	require.PanicsWithValue(t, `object "Widget": duplicate ALTER clause "ADD X"`, func() {
		block := Attribute{Name: "x", Type: AttrObject, SqlOnly: true, Children: []Attribute{{Name: "y", Type: AttrString}}}
		ValidateObject(Object{
			Name: "Widget",
			Alter: AlterSpec{Clauses: []AlterClause{
				AlterParameterBlock("Add", "ADD X", &block),
				AlterKeywordBlock("AddAgain", "ADD X", &block),
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
	inner := Attribute{Name: "inner", Type: AttrList, Elem: ObjectElem([]Attribute{leaf})}
	outer := Attribute{Name: "outer", Type: AttrSet, Elem: ObjectElem([]Attribute{inner})}
	require.NotPanics(t, func() {
		ValidateObject(Object{
			Name:       "Widget",
			Attributes: []Attribute{outer},
		})
	})

	legacyMessage := `object "Widget": attribute "%s" has Children but is not an object; declare the fields of a list or set on an object Elem (ObjectElem)`
	for name, attr := range map[string]Attribute{
		"label":  {Name: "label", Type: AttrString, Children: []Attribute{leaf}},
		"legacy": {Name: "legacy", Type: AttrList, Children: []Attribute{leaf}},
		"items":  {Name: "items", Type: AttrSet, Elem: &Attribute{Type: AttrString}, Children: []Attribute{leaf}},
	} {
		require.PanicsWithValue(t, fmt.Sprintf(legacyMessage, name), func() {
			ValidateObject(Object{Name: "Widget", Attributes: []Attribute{attr}})
		})
	}
	require.PanicsWithValue(t, fmt.Sprintf(legacyMessage, "nested"), func() {
		ValidateObject(Object{Name: "Widget", Attributes: []Attribute{{
			Name: "outer",
			Type: AttrList,
			Elem: ObjectElem([]Attribute{{Name: "nested", Type: AttrList, Children: []Attribute{leaf}}}),
		}}})
	})

	require.PanicsWithValue(t, `object "Widget": attribute "label" has an object Elem and needs to be a list or set`, func() {
		ValidateObject(Object{Name: "Widget", Attributes: []Attribute{{Name: "label", Type: AttrString, Elem: ObjectElem([]Attribute{leaf})}}})
	})

	require.PanicsWithValue(t, `object "Widget": attribute "empty" has an object Elem with no Children`, func() {
		ValidateObject(Object{Name: "Widget", Attributes: []Attribute{{Name: "empty", Type: AttrList, Elem: ObjectElem(nil)}}})
	})

	require.NotPanics(t, func() {
		ValidateObject(Object{Name: "Widget", Attributes: []Attribute{{Name: "object", Type: AttrObject, Children: []Attribute{leaf}}}})
	})

	require.PanicsWithValue(t, `object "Widget": attribute "object" is an object with no Children`, func() {
		ValidateObject(Object{Name: "Widget", Attributes: []Attribute{{Name: "object", Type: AttrObject}}})
	})

	require.PanicsWithValue(t, `object "Widget": attribute "object" is an object and has an Elem; the fields of an object are its Children`, func() {
		ValidateObject(Object{Name: "Widget", Attributes: []Attribute{{Name: "object", Type: AttrObject, Children: []Attribute{leaf}, Elem: ObjectElem([]Attribute{leaf})}}})
	})

	for name, attr := range map[string]Attribute{
		"min": {Name: "min", Type: AttrObject, Children: []Attribute{leaf}, MinItems: new(1)},
		"max": {Name: "max", Type: AttrObject, Children: []Attribute{leaf}, MaxItems: new(1)},
	} {
		require.PanicsWithValue(t, fmt.Sprintf(`object "Widget": attribute %q item bounds require a list or set`, name), func() {
			ValidateObject(Object{Name: "Widget", Attributes: []Attribute{attr}})
		})
	}

	// The fields of a single block carry their CREATE clause when the block has an object clause.
	block := Attribute{
		Name:      "block",
		Type:      AttrObject,
		CreateSql: ObjectParameter(Sql("BLOCK ="), Parentheses()),
		Children:  []Attribute{{Name: "inner", Type: AttrString, CreateSql: StringParameter()}},
	}
	require.NotPanics(t, func() {
		ValidateObject(Object{Name: "Widget", Attributes: []Attribute{block}, Create: CreateSpec{Clauses: []*Attribute{&block}}})
	})

	require.PanicsWithValue(t, `object "Widget": attribute "outer" has duplicate nested attribute "inner"`, func() {
		ValidateObject(Object{
			Name: "Widget",
			Attributes: []Attribute{{
				Name: "outer",
				Type: AttrSet,
				Elem: ObjectElem([]Attribute{inner, {Name: "inner", Type: AttrString}}),
			}},
		})
	})

	require.PanicsWithValue(t, `object "Widget": nested attribute "outer"."inner" has a CREATE clause`, func() {
		ValidateObject(Object{
			Name: "Widget",
			Attributes: []Attribute{{
				Name: "outer",
				Type: AttrSet,
				Elem: ObjectElem([]Attribute{{
					Name:      "inner",
					Type:      AttrString,
					CreateSql: StringParameter(),
				}}),
			}},
		})
	})

	require.PanicsWithValue(t, `object "Widget": nested attribute "inner"."leaf" has a SET clause`, func() {
		ValidateObject(Object{
			Name: "Widget",
			Attributes: []Attribute{{
				Name: "outer",
				Type: AttrSet,
				Elem: ObjectElem([]Attribute{{
					Name: "inner",
					Type: AttrList,
					Elem: ObjectElem([]Attribute{{
						Name:   "leaf",
						Type:   AttrString,
						SetSql: StringParameter(),
					}}),
				}}),
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
		Elem: ObjectElem([]Attribute{
			{Name: "storage_location_name", Type: AttrString, Required: true},
			{Name: "storage_base_url", Type: AttrString, Required: true},
			{Name: "allowed_locations", Type: AttrList, MinItems: new(1), Elem: &Attribute{Type: AttrString}},
		}),
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
				Elem: ObjectElem([]Attribute{{
					Name:     "storage_base_url",
					Type:     AttrString,
					MinItems: new(1),
				}}),
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

	paired := valid
	column := Attribute{Name: "name", Type: AttrString}
	paired.CustomOperations = append(paired.CustomOperations, CustomOperationSpec{
		Kind: "DescribeColumns",
		Clauses: []CustomOperationClause{
			CustomStaticSQL("DESCRIBE"),
			CustomStaticField("columnsType", "TYPE = COLUMNS"),
			CustomIdentifier("Database", "AccountObjectIdentifier", true),
			CustomText("TABLE_NAME"),
		},
		Paired: &CustomPairedShow{DbRow: "fakeRow", Plain: "Fake", Columns: []*Attribute{&column}},
	})
	require.NotPanics(t, func() { ValidateObject(paired) })

	emptyPaired := paired
	emptyPaired.CustomOperations[len(emptyPaired.CustomOperations)-1].Paired = &CustomPairedShow{Plain: "Fake"}
	require.Panics(t, func() { ValidateObject(emptyPaired) })
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
		Name:       "WarehouseAdaptive",
		VariantOf:  &base,
		Attributes: []Attribute{kind, comment},
		Create:     CreateSpec{Clauses: []*Attribute{&kind, &comment}},
	}
	require.NotPanics(t, func() { ValidateObject(obj) })
	require.Equal(t, []Attribute{comment}, obj.ResourceAttributes())
	require.Equal(t, "WAREHOUSE", obj.VariantCreateConfig().CreateSql)

	withSql := obj
	withSql.Create.Sql = "WAREHOUSE"
	require.PanicsWithValue(t, `object "WarehouseAdaptive": variant CREATE words are VariantEntityCreateName`, func() {
		ValidateObject(withSql)
	})

	require.PanicsWithValue(t, `object "Widget": VariantEntityCreateName is for a variant`, func() {
		ValidateObject(Object{Name: "Widget", VariantEntityCreateName: "WIDGET"})
	})
	require.PanicsWithValue(t, `object "Widget": Create.WithoutOrReplace is for a variant`, func() {
		ValidateObject(Object{Name: "Widget", Create: CreateSpec{WithoutOrReplace: true}})
	})
}

func Test_ValidateObject_describeGeneric(t *testing.T) {
	property := Attribute{Name: "property", Type: AttrString, Readonly: true}
	obj := func(column Attribute, generators ...GeneratorOption) Object {
		return Object{
			Name:         "Widget",
			DescribeRows: OutputSpec{Columns: []*Attribute{&column}},
			Generators:   generators,
		}
	}
	require.NotPanics(t, func() { ValidateObject(obj(property, WithSdk(SdkDescribe()))) })
	require.PanicsWithValue(t, `object "Widget": DescribeRows needs SdkDescribe`, func() {
		ValidateObject(obj(property))
	})
	writable := Attribute{Name: "property", Type: AttrString}
	require.PanicsWithValue(t, `object "Widget": DESCRIBE rows entry "property" is writable but not in Attributes`, func() {
		ValidateObject(obj(writable, WithSdk(SdkDescribe())))
	})
}

func Test_ValidateObject_showPlainOnly(t *testing.T) {
	name := Attribute{Name: "name", Type: AttrString, Readonly: true}
	other := Attribute{Name: "other", Type: AttrString, Readonly: true}
	obj := func(generators ...GeneratorOption) Object {
		return Object{
			Name:       "Widget",
			NamePlural: "Widgets",
			Show:       OutputSpec{Columns: []*Attribute{&name}},
			Generators: generators,
		}
	}
	require.NotPanics(t, func() { ValidateObject(obj(WithSdk(SdkShow(SdkShowPlainOnly(&name, "string"))))) })
	require.PanicsWithValue(t, `object "Widget": SdkShowPlainOnly needs SdkShow`, func() {
		ValidateObject(obj(WithSdk(func(s *SdkGenerator) {
			s.ShowPlainOnlyFields = []SdkShowPlainOnlyField{{Column: &name, Kind: "string"}}
		})))
	})
	require.PanicsWithValue(t, `object "Widget": SdkShowPlainOnly needs a column and a kind`, func() {
		ValidateObject(obj(WithSdk(SdkShow(SdkShowPlainOnly(&name, "")))))
	})
	require.PanicsWithValue(t, `object "Widget": SdkShowPlainOnly column "other" is not in the SHOW columns`, func() {
		ValidateObject(obj(WithSdk(SdkShow(SdkShowPlainOnly(&other, "string")))))
	})
	require.PanicsWithValue(t, `object "Widget": SdkShowPlainOnly column "name" is repeated`, func() {
		ValidateObject(obj(WithSdk(SdkShow(SdkShowPlainOnly(&name, "string"), SdkShowPlainOnly(&name, "string")))))
	})
}

func Test_ValidateObject_describePlainOnlyField(t *testing.T) {
	name := Attribute{Name: "name", Type: AttrString, Readonly: true}
	obj := func(generators ...GeneratorOption) Object {
		return Object{
			Name:       "Widget",
			Describe:   OutputSpec{Columns: []*Attribute{&name}},
			Generators: generators,
		}
	}
	require.NotPanics(t, func() {
		ValidateObject(obj(WithSdk(SdkDescribe(
			SdkDescribePlainOnlyField("DatabaseName", "string"),
			SdkDescribePlainOnlyField("SchemaName", "string"),
		))))
	})
	require.PanicsWithValue(t, `object "Widget": SdkDescribePlainOnlyField needs SdkDescribe`, func() {
		ValidateObject(obj(WithSdk(func(s *SdkGenerator) {
			s.DescribePlainOnlyFields = []SdkPlainOnlyField{{Name: "DatabaseName", Kind: "string"}}
		})))
	})
	require.PanicsWithValue(t, `object "Widget": SdkDescribePlainOnlyField needs a name and a kind`, func() {
		ValidateObject(obj(WithSdk(SdkDescribe(SdkDescribePlainOnlyField("", "string")))))
	})
	require.PanicsWithValue(t, `object "Widget": SdkDescribePlainOnlyField "Name" duplicates a DESCRIBE field`, func() {
		ValidateObject(obj(WithSdk(SdkDescribe(SdkDescribePlainOnlyField("Name", "string")))))
	})
	row := Attribute{Name: "property", Type: AttrString, Readonly: true}
	withRows := obj(WithSdk(SdkDescribe(SdkDescribePlainOnlyField("DatabaseName", "string"))))
	withRows.DescribeRows = OutputSpec{Columns: []*Attribute{&row}}
	require.PanicsWithValue(t, `object "Widget": SdkDescribePlainOnlyField cannot be used with DescribeRows`, func() {
		ValidateObject(withRows)
	})
}

func Test_ValidateObject_describeDetailsFromIdentifier(t *testing.T) {
	name := Attribute{Name: "name", Type: AttrString, Readonly: true}
	row := Attribute{Name: "property", Type: AttrString, Readonly: true}
	obj := func(idType string, generators ...GeneratorOption) Object {
		return Object{
			Name:       "Widget",
			IdType:     idType,
			Describe:   OutputSpec{Columns: []*Attribute{&name}},
			Generators: generators,
		}
	}
	both := SdkDescribe(
		SdkDescribePlainOnlyField("DatabaseName", "string"),
		SdkDescribePlainOnlyField("SchemaName", "string"),
		SdkDescribeDetailsFromIdentifier(),
	)
	require.NotPanics(t, func() { ValidateObject(obj("SchemaObjectIdentifier", WithSdk(both))) })
	require.NotPanics(t, func() {
		ValidateObject(obj("SchemaObjectIdentifier", WithSdk(SdkDescribe(
			SdkDescribePlainOnlyField("DatabaseName", "string"),
			SdkDescribeDetailsFromIdentifier(),
		))))
	})
	require.NotPanics(t, func() {
		ValidateObject(obj("SchemaObjectIdentifier", WithSdk(SdkDescribe(
			SdkDescribePlainOnlyField("SchemaName", "string"),
			SdkDescribeDetailsFromIdentifier(),
		))))
	})
	withOtherMethod := obj("SchemaObjectIdentifier", WithSdk(both))
	withOtherMethod.CustomInterfaceMethods = []CustomInterfaceMethodSpec{{Name: "Use"}}
	require.NotPanics(t, func() { ValidateObject(withOtherMethod) })
	require.NotPanics(t, func() {
		ValidateObject(obj("DatabaseObjectIdentifier", WithSdk(SdkDescribe(
			SdkDescribePlainOnlyField("DatabaseName", "string"),
			SdkDescribeDetailsFromIdentifier(),
		))))
	})

	require.PanicsWithValue(t, `object "Widget": SdkDescribeDetailsFromIdentifier needs SdkDescribe`, func() {
		ValidateObject(obj("SchemaObjectIdentifier", WithSdk(func(s *SdkGenerator) { s.DescribeDetailsFromIdentifier = true })))
	})
	require.PanicsWithValue(t, `object "Widget": SdkDescribeDetailsFromIdentifier needs IdType SchemaObjectIdentifier or DatabaseObjectIdentifier`, func() {
		ValidateObject(obj("AccountObjectIdentifier", WithSdk(both)))
	})
	require.PanicsWithValue(t, `object "Widget": SdkDescribeDetailsFromIdentifier plain-only field "SchemaName" must be DatabaseName`, func() {
		ValidateObject(obj("DatabaseObjectIdentifier", WithSdk(both)))
	})
	require.PanicsWithValue(t, `object "Widget": SdkDescribeDetailsFromIdentifier needs a plain-only DatabaseName`, func() {
		ValidateObject(obj("DatabaseObjectIdentifier", WithSdk(SdkDescribe(SdkDescribeDetailsFromIdentifier()))))
	})
	require.PanicsWithValue(t, `object "Widget": SdkDescribeDetailsFromIdentifier needs a plain-only DatabaseName or SchemaName`, func() {
		ValidateObject(obj("SchemaObjectIdentifier", WithSdk(SdkDescribe(SdkDescribeDetailsFromIdentifier()))))
	})
	require.PanicsWithValue(t, `object "Widget": SdkDescribeDetailsFromIdentifier plain-only field "Comment" must be DatabaseName or SchemaName`, func() {
		ValidateObject(obj("SchemaObjectIdentifier", WithSdk(SdkDescribe(
			SdkDescribePlainOnlyField("Comment", "string"),
			SdkDescribeDetailsFromIdentifier(),
		))))
	})
	require.PanicsWithValue(t, `object "Widget": SdkDescribeDetailsFromIdentifier plain-only field "DatabaseName" must be a string`, func() {
		ValidateObject(obj("SchemaObjectIdentifier", WithSdk(SdkDescribe(
			SdkDescribePlainOnlyField("DatabaseName", "int"),
			SdkDescribeDetailsFromIdentifier(),
		))))
	})
	withRows := obj("SchemaObjectIdentifier", WithSdk(SdkDescribe(SdkDescribeDetailsFromIdentifier())))
	withRows.DescribeRows = OutputSpec{Columns: []*Attribute{&row}}
	require.PanicsWithValue(t, `object "Widget": SdkDescribeDetailsFromIdentifier cannot be used with DescribeRows`, func() {
		ValidateObject(withRows)
	})
	withDetails := obj("SchemaObjectIdentifier", WithSdk(SdkDescribe(
		SdkDescribeDetails(),
		SdkDescribeDetailsFromIdentifier(),
	)))
	withDetails.DescribeRows = OutputSpec{Columns: []*Attribute{&row}}
	require.PanicsWithValue(t, `object "Widget": SdkDescribeDetailsFromIdentifier cannot be used with SdkDescribeDetails`, func() {
		ValidateObject(withDetails)
	})
	withCustom := obj("SchemaObjectIdentifier", WithSdk(both))
	withCustom.CustomInterfaceMethods = []CustomInterfaceMethodSpec{{Name: "DescribeDetails"}}
	require.PanicsWithValue(t, `object "Widget": SdkDescribeDetailsFromIdentifier conflicts with custom interface method DescribeDetails`, func() {
		ValidateObject(withCustom)
	})
}

func Test_Object_PairedDescribeColumns(t *testing.T) {
	details := Attribute{Name: "active", Type: AttrString, Readonly: true}
	row := Attribute{Name: "property", Type: AttrString, Readonly: true}
	plain := Object{Describe: OutputSpec{Columns: []*Attribute{&details}}}
	require.Equal(t, []*Attribute{&details}, plain.PairedDescribeColumns())
	generic := Object{Describe: plain.Describe, DescribeRows: OutputSpec{Columns: []*Attribute{&row}}}
	require.Equal(t, []*Attribute{&row}, generic.PairedDescribeColumns())
	require.Equal(t, []*Attribute{&details}, generic.Describe.Columns)
}

func Test_Object_DescribePlainName(t *testing.T) {
	row := Attribute{Name: "property", Type: AttrString, Readonly: true}
	require.Equal(t, "WidgetDetails", Object{Name: "Widget"}.DescribePlainName())
	require.Equal(t, "WidgetProperty", Object{Name: "Widget", DescribeRows: OutputSpec{Columns: []*Attribute{&row}}}.DescribePlainName())
}

func Test_ValidateObject_describeDetails(t *testing.T) {
	row := Attribute{Name: "property", Type: AttrString, Readonly: true}
	active := Attribute{Name: "active", Type: AttrString, Readonly: true}
	other := Attribute{Name: "other", Type: AttrString, Readonly: true}
	obj := func(rows []*Attribute, opts ...SdkDescribeOption) Object {
		return Object{
			Name:         "Widget",
			Describe:     OutputSpec{Columns: []*Attribute{&active}},
			DescribeRows: OutputSpec{Columns: rows},
			Generators:   []GeneratorOption{WithSdk(SdkDescribe(opts...))},
		}
	}
	require.NotPanics(t, func() {
		ValidateObject(obj([]*Attribute{&row}, SdkDescribeDetails(SdkDescribeDetailsNestedField(&active))))
	})
	require.PanicsWithValue(t, `object "Widget": SdkDescribeDetails needs SdkDescribe and DescribeRows`, func() {
		ValidateObject(obj(nil, SdkDescribeDetails()))
	})
	require.Panics(t, func() {
		ValidateObject(obj([]*Attribute{&row}, SdkDescribeDetails(SdkDescribeDetailsNestedField(&other))))
	})
	require.PanicsWithValue(t, `object "Widget": SdkDescribeDetailsMethod needs SdkDescribeDetails`, func() {
		o := obj([]*Attribute{&row})
		o.Generators = []GeneratorOption{WithSdk(func(s *SdkGenerator) { s.Describe = true; s.DescribeDetailsMethod = true })}
		ValidateObject(o)
	})
}

func Test_NestedDetailsElem(t *testing.T) {
	locations := Attribute{Name: "storage_locations", Type: AttrList}
	require.Equal(t, "WidgetStorageLocationDetails", NestedDetailsElem(Object{Name: "Widget"}, &locations))
	require.Panics(t, func() { NestedDetailsElem(Object{Name: "Widget"}, &Attribute{Name: "x", Type: AttrString}) })
}

func Test_ValidateObject_describeDetailsParse(t *testing.T) {
	row := Attribute{Name: "property", Type: AttrString, Readonly: true}
	active := Attribute{Name: "active", Type: AttrString, Readonly: true}
	size := Attribute{Name: "size", Type: AttrInt, Readonly: true}
	locations := Attribute{Name: "locations", Type: AttrList, Readonly: true, Elem: &Attribute{Type: AttrString}}
	obj := func(columns []*Attribute, mappings ...DescribeDetailsMapping) Object {
		return Object{
			Name:         "Widget",
			Describe:     OutputSpec{Columns: columns},
			DescribeRows: OutputSpec{Columns: []*Attribute{&row}},
			Generators: []GeneratorOption{WithSdk(SdkDescribe(SdkDescribeDetails(
				SdkDescribeDetailsNestedField(&locations),
				SdkDescribeDetailsParse(mappings...),
			)))},
		}
	}
	require.NotPanics(t, func() {
		ValidateObject(obj([]*Attribute{&active, &locations}, DescribeProperty(&active, "ACTIVE"), DescribePropertyPrefix(&locations, "LOC_")))
	})
	require.PanicsWithValue(t, `object "Widget": DESCRIBE column "active" has no details mapping`, func() {
		ValidateObject(obj([]*Attribute{&active, &locations}, DescribePropertyPrefix(&locations, "LOC_")))
	})
	require.PanicsWithValue(t, `object "Widget": DESCRIBE details mapping "active" is set twice`, func() {
		ValidateObject(obj([]*Attribute{&active}, DescribeProperty(&active, "A"), DescribeProperty(&active, "B")))
	})
	require.PanicsWithValue(t, `object "Widget": DESCRIBE details prefix mapping "active" must be a nested field`, func() {
		ValidateObject(obj([]*Attribute{&active}, DescribePropertyPrefix(&active, "A_")))
	})
	require.PanicsWithValue(t, `object "Widget": DESCRIBE details property mapping "size" must be a required string column`, func() {
		ValidateObject(obj([]*Attribute{&size}, DescribeProperty(&size, "SIZE")))
	})
	require.PanicsWithValue(t, `object "Widget": DESCRIBE details property mapping "locations" must be a required string column`, func() {
		ValidateObject(obj([]*Attribute{&locations}, DescribeProperty(&locations, "LOC")))
	})
	require.PanicsWithValue(t, `object "Widget": DESCRIBE details mapping "size" is not a DESCRIBE column`, func() {
		ValidateObject(obj([]*Attribute{&active}, DescribeProperty(&active, "A"), DescribeProperty(&size, "SIZE")))
	})
}

func Test_ValidateObject_objectListVariants(t *testing.T) {
	baseUrl := Attribute{Name: "base_url", Type: AttrString, CreateSql: Clause()}
	keyId := Attribute{Name: "key_id", Type: AttrString, CreateSql: Clause()}
	static := Attribute{Name: "kind", SqlOnly: true, CreateSql: StaticSQL("kind", "KIND = 'X'")}
	list := func(elem SqlClause, children ...Attribute) Object {
		attr := Attribute{
			Name:      "location",
			Type:      AttrList,
			CreateSql: ListOfObjects(Sql("LOCATIONS")),
			Elem:      &Attribute{Type: AttrObject, CreateSql: elem, Children: children},
		}
		return Object{Name: "Widget", Attributes: []Attribute{attr}, Create: CreateSpec{Clauses: []*Attribute{&attr}}}
	}

	require.NotPanics(t, func() {
		ValidateObject(list(ObjectParameter(Parentheses(), WithVariants(Variant("First", &baseUrl), Variant("Second", &static, &baseUrl))), baseUrl, keyId))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "location" has an empty or duplicate variant "First"`, func() {
		ValidateObject(list(ObjectParameter(WithVariants(Variant("First", &baseUrl), Variant("First", &keyId))), baseUrl, keyId))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "location" variant "First" has no fields`, func() {
		ValidateObject(list(ObjectParameter(WithVariants(Variant("First"))), baseUrl))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "location" variant "First" has duplicate field "base_url"`, func() {
		ValidateObject(list(ObjectParameter(WithVariants(Variant("First", &baseUrl, &baseUrl))), baseUrl))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "location" variant "First" field "key_id" is not a child of the element`, func() {
		ValidateObject(list(ObjectParameter(WithVariants(Variant("First", &keyId))), baseUrl))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "location" variant "First" field "key_id" has no CREATE clause`, func() {
		noClause := Attribute{Name: "key_id", Type: AttrString}
		ValidateObject(list(ObjectParameter(WithVariants(Variant("First", &noClause))), baseUrl))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "location" variant "First" field "block" needs Children`, func() {
		block := Attribute{Name: "block", Type: AttrObject, SqlOnly: true, CreateSql: ObjectParameter()}
		ValidateObject(list(ObjectParameter(WithVariants(Variant("First", &block))), baseUrl))
	})
	require.PanicsWithValue(t, `object "Widget": nested attribute "location"."base_url" has a CREATE clause`, func() {
		attr := Attribute{Name: "location", Type: AttrList, Elem: &Attribute{Type: AttrObject, Children: []Attribute{baseUrl}}}
		ValidateObject(Object{Name: "Widget", Attributes: []Attribute{attr}})
	})
}

func Test_Attribute_nestedFields(t *testing.T) {
	child := Attribute{Name: "child"}
	list := Attribute{Type: AttrList, Elem: ObjectElem([]Attribute{child})}
	require.Equal(t, []Attribute{child}, list.NestedFields())
	require.True(t, list.IsObjectList())
	set := Attribute{Type: AttrSet, Elem: ObjectElem([]Attribute{child})}
	require.Equal(t, []Attribute{child}, set.NestedFields())
	require.Equal(t, []Attribute{child}, Attribute{Type: AttrObject, Children: []Attribute{child}}.NestedFields())
	scalars := Attribute{Type: AttrList, Elem: &Attribute{Type: AttrString}}
	require.Empty(t, scalars.NestedFields())
	require.False(t, scalars.IsObjectList())
	require.Empty(t, Attribute{Type: AttrList, Children: []Attribute{child}}.NestedFields(), "Children on a list is not a block")
}

func Test_Attribute_sqlStructNames(t *testing.T) {
	attr := Attribute{Name: "storage_location"}
	require.Equal(t, "ExternalVolumeStorageLocation", attr.SqlStructName("ExternalVolume"))
	require.Equal(t, "GCSStorageLocationParams", attr.SqlVariantStructName("GCS"))
}

func Test_ValidateObject_skipCreateOnObjectLists(t *testing.T) {
	leaf := Attribute{Name: "leaf", Type: AttrString}
	inner := Attribute{Name: "inner", Type: AttrList, Elem: ObjectElem([]Attribute{leaf})}
	outer := Attribute{Name: "outer", Type: AttrList, SkipCreate: true, Elem: ObjectElem([]Attribute{inner})}
	require.NotPanics(t, func() { ValidateObject(Object{Name: "Widget", Attributes: []Attribute{outer}}) })
}

func Test_ObjectElem(t *testing.T) {
	child := Attribute{Name: "child", Type: AttrString}
	list := Attribute{Name: "block", Type: AttrList, Elem: ObjectElem([]Attribute{child})}
	require.True(t, list.IsObjectList())
	require.Equal(t, []Attribute{child}, list.NestedFields())
	require.Equal(t, AttrObject, list.Elem.Type)
}

func Test_Attribute_IsSingleBlock(t *testing.T) {
	leaf := Attribute{Name: "leaf", Type: AttrString}
	require.True(t, Attribute{Type: AttrObject, Children: []Attribute{leaf}}.IsSingleBlock())
	require.True(t, Attribute{Type: AttrList, MaxItems: new(1), Elem: ObjectElem([]Attribute{leaf})}.IsSingleBlock())
	require.False(t, Attribute{Type: AttrList, Elem: ObjectElem([]Attribute{leaf})}.IsSingleBlock())
	require.False(t, Attribute{Type: AttrList, MaxItems: new(1), Elem: &Attribute{Type: AttrString}}.IsSingleBlock())
}

func Test_Attribute_nestedSqlFields(t *testing.T) {
	typeOauth := Attribute{Name: "type_oauth", SqlOnly: true, CreateSql: StaticSQL("restAuthType", "TYPE = OAUTH")}
	uri := Attribute{Name: "uri", Type: AttrString}

	object := Attribute{Type: AttrObject, Children: []Attribute{typeOauth, uri}}
	require.Equal(t, []Attribute{uri}, object.NestedFields())
	require.Equal(t, []Attribute{typeOauth, uri}, object.NestedSqlFields())

	list := Attribute{Type: AttrList, Elem: ObjectElem([]Attribute{typeOauth, uri})}
	require.Equal(t, []Attribute{uri}, list.NestedFields())
	require.Equal(t, []Attribute{typeOauth, uri}, list.NestedSqlFields())

	require.Empty(t, Attribute{Type: AttrString}.NestedSqlFields())
}

func Test_ValidateObject_nestedStaticFields(t *testing.T) {
	typeOauth := Attribute{Name: "type_oauth", SqlOnly: true, CreateSql: StaticSQL("restAuthType", "TYPE = OAUTH")}
	uri := Attribute{Name: "uri", Type: AttrString, CreateSql: StringParameter()}
	validate := func(block Attribute) {
		ValidateObject(Object{Name: "Widget", Attributes: []Attribute{block}, Create: CreateSpec{Clauses: []*Attribute{&block}}})
	}

	t.Run("static field of a single block with an object clause", func(t *testing.T) {
		require.NotPanics(t, func() {
			validate(Attribute{Name: "block", Type: AttrObject, CreateSql: ObjectParameter(Sql("BLOCK ="), Parentheses()), Children: []Attribute{typeOauth, uri}})
		})
	})

	t.Run("static field of a list element with an object clause", func(t *testing.T) {
		require.NotPanics(t, func() {
			validate(Attribute{
				Name:      "blocks",
				Type:      AttrList,
				CreateSql: ListOfObjects(Sql("BLOCKS"), Parentheses()),
				Elem: &Attribute{
					Type:      AttrObject,
					CreateSql: ObjectParameter(Parentheses()),
					Children:  []Attribute{typeOauth, uri},
				},
			})
		})
	})

	t.Run("without an object clause", func(t *testing.T) {
		require.PanicsWithValue(t, `object "Widget": nested attribute "block"."type_oauth" is SQL-only but not a static clause, a keyword, an identifier, a parameter, or a nested struct of an object clause`, func() {
			ValidateObject(Object{Name: "Widget", Attributes: []Attribute{{Name: "block", Type: AttrObject, Children: []Attribute{typeOauth, uri}}}})
		})
	})

	t.Run("sql-only without a clause", func(t *testing.T) {
		bare := Attribute{Name: "other", SqlOnly: true}
		require.PanicsWithValue(t, `object "Widget": nested attribute "block"."other" is SQL-only but not a static clause, a keyword, an identifier, a parameter, or a nested struct of an object clause`, func() {
			validate(Attribute{Name: "block", Type: AttrObject, CreateSql: ObjectParameter(Sql("BLOCK =")), Children: []Attribute{bare, uri}})
		})
	})

	t.Run("parameter child", func(t *testing.T) {
		where := Attribute{Name: "other", Type: AttrString, SqlOnly: true, CreateSql: PlainParameter(Sql("WHERE"), NoEquals())}
		require.NotPanics(t, func() {
			validate(Attribute{Name: "block", Type: AttrObject, CreateSql: ObjectParameter(Sql("BLOCK =")), Children: []Attribute{where, uri}})
		})
	})

	t.Run("static field is private", func(t *testing.T) {
		exported := Attribute{Name: "exported", SqlOnly: true, CreateSql: StaticSQL("RestAuthType", "TYPE = OAUTH")}
		require.Panics(t, func() {
			validate(Attribute{Name: "block", Type: AttrObject, CreateSql: ObjectParameter(Sql("BLOCK =")), Children: []Attribute{exported, uri}})
		})
	})

	t.Run("only static fields", func(t *testing.T) {
		require.PanicsWithValue(t, `object "Widget": attribute "block" is an object with no Children`, func() {
			validate(Attribute{Name: "block", Type: AttrObject, CreateSql: ObjectParameter(Sql("BLOCK =")), Children: []Attribute{typeOauth}})
		})
	})
}

func Test_ValidateObject_exactlyOneOf(t *testing.T) {
	one := Attribute{Name: "one", Type: AttrString}
	two := Attribute{Name: "two", Type: AttrString}
	three := Attribute{Name: "three", Type: AttrString}
	required := Attribute{Name: "required", Type: AttrString, Required: true}
	static := Attribute{Name: "static", SqlOnly: true, CreateSql: StaticSQL("static", "A = B")}
	outside := Attribute{Name: "outside", Type: AttrString}
	newObject := func(groups ...[]*Attribute) Object {
		return Object{
			Name:         "Widget",
			Attributes:   []Attribute{one, two, three, required, static},
			ExactlyOneOf: groups,
			Create:       CreateSpec{Clauses: []*Attribute{&static}},
		}
	}

	require.NotPanics(t, func() { ValidateObject(newObject()) }, "no group")
	require.NotPanics(t, func() { ValidateObject(newObject([]*Attribute{&one, &two})) })

	require.PanicsWithValue(t, `object "Widget": ExactlyOneOf group 0 needs at least two attributes`, func() {
		ValidateObject(newObject([]*Attribute{&one}))
	})
	require.PanicsWithValue(t, `object "Widget": ExactlyOneOf group 0 has a nil attribute`, func() {
		ValidateObject(newObject([]*Attribute{&one, nil}))
	})
	require.PanicsWithValue(t, `object "Widget": ExactlyOneOf group 0 attribute "outside" is not an attribute of the object`, func() {
		ValidateObject(newObject([]*Attribute{&one, &outside}))
	})
	require.PanicsWithValue(t, `object "Widget": ExactlyOneOf group 0 attribute "static" is not an attribute of the object`, func() {
		ValidateObject(newObject([]*Attribute{&one, &static}))
	})
	require.PanicsWithValue(t, `object "Widget": ExactlyOneOf group 0 attribute "required" is required`, func() {
		ValidateObject(newObject([]*Attribute{&one, &required}))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "one" is in more than one ExactlyOneOf group, or twice in one`, func() {
		ValidateObject(newObject([]*Attribute{&one, &two}, []*Attribute{&one, &three}))
	})
}

func Test_ValidateObject_allowEmpty(t *testing.T) {
	validate := func(attr Attribute) {
		ValidateObject(Object{Name: "Widget", Attributes: []Attribute{attr}, Set: []*Attribute{&attr}})
	}
	require.NotPanics(t, func() {
		validate(Attribute{Name: "comment", Type: AttrString, SetSql: StringParameter(AllowEmpty())})
	})
	require.PanicsWithValue(t, `object "Widget": attribute "count" SET clause allows an empty value but is not a string parameter`, func() {
		validate(Attribute{Name: "count", Type: AttrInt, SetSql: IntParameter(AllowEmpty())})
	})
	require.PanicsWithValue(t, `object "Widget": attribute "flag" SET clause allows an empty value but is not a string parameter`, func() {
		validate(Attribute{Name: "flag", Type: AttrString, SetSql: BoolKeyword(AllowEmpty())})
	})
}

func Test_ValidateObject_setConflictingFields(t *testing.T) {
	one := Attribute{Name: "one", Type: AttrString}
	two := Attribute{Name: "two", Type: AttrString}
	three := Attribute{Name: "three", Type: AttrString}
	newObject := func(pairs ...[]*Attribute) Object {
		return Object{
			Name:       "Widget",
			Attributes: []Attribute{one, two, three},
			Set:        []*Attribute{&one, &two},
			Alter:      AlterSpec{SetConflictingFields: pairs},
		}
	}
	require.NotPanics(t, func() { ValidateObject(newObject([]*Attribute{&one, &two})) })
	require.PanicsWithValue(t, `object "Widget": SetConflictingFields entry 0 needs exactly two attributes`, func() {
		ValidateObject(newObject([]*Attribute{&one, &two, &three}))
	})
	require.PanicsWithValue(t, `object "Widget": SetConflictingFields entry 0 has an attribute that is not in Set`, func() {
		ValidateObject(newObject([]*Attribute{&one, &three}))
	})
	require.PanicsWithValue(t, `object "Widget": SetConflictingFields entry 0 has an attribute that is not in Set`, func() {
		ValidateObject(newObject([]*Attribute{&one, nil}))
	})
}

func Test_ValidateObject_describeDetailsObjectColumnsAndVariants(t *testing.T) {
	active := Attribute{Name: "active", Type: AttrString, Readonly: true}
	config := Attribute{Name: "config", Type: AttrObject, Required: true, Children: []Attribute{{Name: "uri", Type: AttrString}}}
	base := Object{Name: "Widget", NamePlural: "Widgets", IdType: "AccountObjectIdentifier"}

	newVariant := func(withDetails bool) Object {
		obj := Object{
			Name:       "WidgetMini",
			VariantOf:  &base,
			Attributes: []Attribute{config},
			Describe:   OutputSpec{Columns: []*Attribute{&active, &config}},
		}
		if withDetails {
			obj.Generators = []GeneratorOption{WithSdk(SdkDescribe(SdkDescribeDetails()))}
		}
		return obj
	}

	t.Run("variant details use the rows of the base and an attribute that is an object is a nested struct", func(t *testing.T) {
		require.NotPanics(t, func() { ValidateObject(newVariant(true)) })
	})

	t.Run("an object column needs the details", func(t *testing.T) {
		require.PanicsWithValue(t, `object "WidgetMini": nested attribute "config" is not a DESCRIBE column`, func() {
			ValidateObject(newVariant(false))
		})
	})

	t.Run("an object without a variant needs rows", func(t *testing.T) {
		obj := newVariant(true)
		obj.VariantOf = nil
		require.PanicsWithValue(t, `object "WidgetMini": SdkDescribeDetails needs SdkDescribe and DescribeRows`, func() { ValidateObject(obj) })
	})
}

func Test_Attribute_IsOptionalColumn(t *testing.T) {
	require.True(t, Attribute{Type: AttrString, OptionalInShow: true}.IsOptionalColumn())
	require.False(t, Attribute{Type: AttrString}.IsOptionalColumn())
	require.True(t, Attribute{Type: AttrObject}.IsOptionalColumn(), "an object that is not required is a pointer")
	require.False(t, Attribute{Type: AttrObject, Required: true}.IsOptionalColumn())
}

func Test_ValidateObject_childrenAtLeastOneValueSet(t *testing.T) {
	obj := func(policy Attribute) Object {
		return Object{Name: "Widget", Attributes: []Attribute{policy}}
	}
	children := []Attribute{{Name: "limit", Type: AttrInt}, {Name: "size", Type: AttrInt}}
	require.NotPanics(t, func() {
		ValidateObject(obj(Attribute{Name: "policy", Type: AttrObject, ChildrenRule: ChildrenAtLeastOne, Children: children}))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "policy" sets ChildrenRule and needs at least two Children that are fields`, func() {
		ValidateObject(obj(Attribute{Name: "policy", Type: AttrObject, ChildrenRule: ChildrenAtLeastOne, Children: children[:1]}))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "policy" sets ChildrenRule but is not an object`, func() {
		ValidateObject(obj(Attribute{Name: "policy", Type: AttrString, ChildrenRule: ChildrenAtLeastOne}))
	})
}

func Test_ValidateObject_boolTrueValueNeedsStoredAsText(t *testing.T) {
	require.PanicsWithValue(t, `object "Widget": attribute "flag" sets BoolTrueValue without StoredAsText`, func() {
		ValidateObject(Object{Name: "Widget", Attributes: []Attribute{{Name: "flag", Type: AttrBool, BoolTrueValue: "ON"}}})
	})
}

func Test_ValidateObject_enumKeywordAndPreexistingStruct(t *testing.T) {
	driver := NewEnum("WidgetDriver", "WidgetDrivers", "JDBC")
	version := Attribute{Name: "minimum_version", Type: AttrString}
	params := Attribute{
		Name:      "entry_params",
		Type:      AttrObject,
		SqlOnly:   true,
		CreateSql: ObjectParameter(Sql("="), Parentheses()),
		Children:  []Attribute{version},
	}
	entries := func(driverAttr Attribute) Attribute {
		return Attribute{
			Name:      "entries",
			Type:      AttrSet,
			CreateSql: ListOfObjects(Sql("ENTRIES"), Parentheses()),
			Elem:      &Attribute{Type: AttrObject, CreateSql: ObjectParameter(), Children: []Attribute{driverAttr, version, params}},
		}
	}
	obj := func(list Attribute) Object {
		return Object{Name: "Widget", Attributes: []Attribute{list}, Create: CreateSpec{Clauses: []*Attribute{&list}}}
	}
	require.NotPanics(t, func() {
		ValidateObject(obj(entries(Attribute{Name: "driver", Type: AttrEnum, Enum: driver, CreateSql: EnumKeyword(NoQuotes())})))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "driver" CREATE value keyword needs an enum attribute with an Enum, a string, a list of identifiers, a list of strings, or a list of objects`, func() {
		ValidateObject(obj(entries(Attribute{Name: "driver", Type: AttrString, CreateSql: EnumKeyword(NoQuotes())})))
	})
}

func Test_ValidateObject_stringListAndStructClauses(t *testing.T) {
	objectWith := func(attr Attribute) Object {
		return Object{Name: "Widget", Create: CreateSpec{Clauses: []*Attribute{&attr}}}
	}
	clusterBy := Attribute{
		Name:      "cluster_by",
		Type:      AttrList,
		SqlOnly:   true,
		Elem:      &Attribute{Type: AttrString},
		CreateSql: StringListKeyword(Parentheses(), Sql("CLUSTER BY")),
	}
	policy := Attribute{
		Name:      "row_access_policy",
		Type:      AttrObject,
		SqlOnly:   true,
		CreateSql: ObjectParameter(),
		Children: []Attribute{
			{
				Name:           "name",
				Type:           AttrIdentifier,
				IdentifierType: "SchemaObjectIdentifier",
				Required:       true,
				SqlOnly:        true,
				CreateSql:      IdentifierKeyword("SchemaObjectIdentifier", Sql("ROW ACCESS POLICY"), FieldName("Name")),
			},
			{
				Name:      "on",
				Type:      AttrList,
				Required:  true,
				SqlOnly:   true,
				Elem:      &Attribute{Type: AttrString},
				CreateSql: StringListKeyword(Parentheses(), Sql("ON")),
			},
		},
	}

	require.NotPanics(t, func() {
		ValidateObject(objectWith(clusterBy))
	})
	require.NotPanics(t, func() {
		ValidateObject(objectWith(policy))
	})
	require.NotPanics(t, func() {
		ValidateObject(objectWith(Attribute{Name: "format", Type: AttrObject, SqlOnly: true, CreateSql: PreexistingStruct("WidgetFormat", Parentheses(), NoComma())}))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "cluster_by" CREATE value keyword needs an enum attribute with an Enum, a string, a list of identifiers, a list of strings, or a list of objects`, func() {
		ValidateObject(objectWith(Attribute{Name: "cluster_by", Type: AttrList, SqlOnly: true, CreateSql: StringListKeyword(Parentheses(), Sql("CLUSTER BY"))}))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "cluster_by" CREATE string list keyword has quotes or equals`, func() {
		ValidateObject(objectWith(Attribute{Name: "cluster_by", Type: AttrList, SqlOnly: true, Elem: &Attribute{Type: AttrString}, CreateSql: StringListKeyword(Parentheses(), Sql("CLUSTER BY"), SingleQuotes())}))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "format" preexisting struct clause has no struct name`, func() {
		ValidateObject(objectWith(Attribute{Name: "format", Type: AttrObject, SqlOnly: true, CreateSql: PreexistingStruct("")}))
	})
	require.PanicsWithValue(t, `object "Widget": attribute "format" preexisting_struct clause is CREATE only`, func() {
		attr := Attribute{Name: "format", Type: AttrObject, SqlOnly: true, CreateSql: PreexistingStruct("WidgetFormat"), SetSql: PreexistingStruct("WidgetFormat")}
		ValidateObject(Object{Name: "Widget", Create: CreateSpec{Clauses: []*Attribute{&attr}}, Set: []*Attribute{&attr}})
	})
}

func Test_ValidateObject_schemaName(t *testing.T) {
	nullable := Attribute{Name: "null?", Type: AttrBool, Readonly: true, SchemaName: "is_nullable"}
	same := Attribute{Name: "is_nullable", Type: AttrString, Readonly: true}
	require.PanicsWithValue(t, `object "Widget": TableColumnDetails columns "null?" and "is_nullable" share schema key "is_nullable"`, func() {
		ValidateObject(widgetWithColumns(&nullable, &same))
	})

	note := Attribute{Name: "comment", Type: AttrString, Readonly: true, SchemaName: "note"}
	renamed := Attribute{Name: "null?", Type: AttrBool, Readonly: true, SchemaName: "comment"}
	require.PanicsWithValue(t, `object "Widget": TableColumnDetails column "null?" schema name "comment" collides with the generated key of "comment"`, func() {
		ValidateObject(widgetWithColumns(&renamed, &note))
	})
}

func Test_SchemaNameIndex(t *testing.T) {
	nullable := Attribute{Name: "null?", Type: AttrBool, Readonly: true, SchemaName: "is_nullable"}
	name := Attribute{Name: "name", Type: AttrString, Readonly: true}
	index := SchemaNameIndex([]Object{widgetWithColumns(&nullable, &name)})
	require.Equal(t, map[string]map[string]string{
		"sdk.TableColumnDetails": {"Null": "is_nullable"},
	}, index)
}

func widgetWithColumns(columns ...*Attribute) Object {
	return Object{
		Name: "Widget",
		CustomOperations: []CustomOperationSpec{{
			Kind:    "DescribeColumns",
			Clauses: []CustomOperationClause{CustomStaticSQL("DESCRIBE")},
			Paired: &CustomPairedShow{
				DbRow:   "row",
				Plain:   "TableColumnDetails",
				Columns: columns,
			},
		}},
	}
}

func identityNamesCreatedBySnowflakeObject(change func(*Object)) Object {
	obj := Object{
		Name:                            "Widget",
		IdentityNamesCreatedBySnowflake: []string{"PUBLIC"},
		Alter:                           AlterSpec{Clauses: []AlterClause{AlterSet()}},
		ResourceHooks:                   []ResourceHookOption{WithBeforeCreate(Generated())},
	}
	change(&obj)
	return obj
}

func Test_ValidateObject_identityNamesCreatedBySnowflake(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		require.NotPanics(t, func() { ValidateObject(identityNamesCreatedBySnowflakeObject(func(*Object) {})) })
	})

	for name, tc := range map[string]struct {
		change func(*Object)
		panics string
	}{
		"empty name": {
			func(o *Object) { o.IdentityNamesCreatedBySnowflake = []string{" "} },
			`object "Widget": IdentityNamesCreatedBySnowflake name " " must be set and must not contain quotes or backslashes`,
		},
		"name that breaks the Go string": {
			func(o *Object) { o.IdentityNamesCreatedBySnowflake = []string{`PUB"LIC`} },
			`object "Widget": IdentityNamesCreatedBySnowflake name "PUB\"LIC" must be set and must not contain quotes or backslashes`,
		},
		"names without the generated hook": {
			func(o *Object) { o.ResourceHooks = nil },
			`object "Widget": IdentityNamesCreatedBySnowflake requires WithBeforeCreate(Generated())`,
		},
		"generated hook without names": {
			func(o *Object) { o.IdentityNamesCreatedBySnowflake = nil },
			`object "Widget": WithBeforeCreate(Generated()) needs IdentityNamesCreatedBySnowflake`,
		},
		"identity schema is missing": {
			func(o *Object) { o.ResourceHooks = append(o.ResourceHooks, WithIdentitySchema(Omitted())) },
			`object "Widget": IdentityNamesCreatedBySnowflake reads the identity name and needs the identity schema`,
		},
		"update is missing": {
			func(o *Object) { o.Alter = AlterSpec{} },
			`object "Widget": IdentityNamesCreatedBySnowflake runs Update and needs ALTER to update the object`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			require.PanicsWithValue(t, tc.panics, func() { ValidateObject(identityNamesCreatedBySnowflakeObject(tc.change)) })
		})
	}
}

func showReuseBaseObject(opts ...GeneratorOption) *Object {
	return &Object{Name: "Base", Generators: opts}
}

func Test_ValidateObject_showReuse(t *testing.T) {
	baseWithShow := func() *Object {
		return showReuseBaseObject(WithSchema(SchemaShow()), WithObjectAssert(ObjectAssertShow()))
	}
	variant := func(base *Object, opts ...GeneratorOption) Object {
		return Object{Name: "Variant", VariantOf: base, Generators: opts}
	}
	reuseBoth := []GeneratorOption{
		WithSchema(SchemaShow(SchemaShowFromBase())),
		WithObjectAssert(ObjectAssertShow(ObjectAssertShowRowFromBase())),
	}

	t.Run("valid reuse resolves to the base", func(t *testing.T) {
		obj := variant(baseWithShow(), reuseBoth...)
		require.NotPanics(t, func() { ValidateObject(obj) })
		require.True(t, obj.ReusesBaseSchemaShow())
		require.True(t, obj.ReusesBaseObjectAssertShow())
		require.Equal(t, "Base", obj.SchemaShowOwnerName())
		require.Equal(t, "Base", obj.ObjectAssertShowOwnerName())
	})
	t.Run("without reuse the object owns its artifacts", func(t *testing.T) {
		obj := variant(baseWithShow(), WithSchema(SchemaShow()), WithObjectAssert(ObjectAssertShow()))
		require.False(t, obj.ReusesBaseSchemaShow())
		require.False(t, obj.ReusesBaseObjectAssertShow())
		require.Equal(t, "Variant", obj.SchemaShowOwnerName())
		require.Equal(t, "Variant", obj.ObjectAssertShowOwnerName())
	})
	t.Run("reuse is independent per generator", func(t *testing.T) {
		obj := variant(baseWithShow(), WithSchema(SchemaShow(SchemaShowFromBase())), WithObjectAssert(ObjectAssertShow()))
		require.Equal(t, "Base", obj.SchemaShowOwnerName())
		require.Equal(t, "Variant", obj.ObjectAssertShowOwnerName())
	})
	t.Run("requires VariantOf", func(t *testing.T) {
		require.Panics(t, func() {
			ValidateObject(Object{Name: "Lone", Generators: []GeneratorOption{WithSchema(SchemaShow(SchemaShowFromBase()))}})
		})
		require.Panics(t, func() {
			ValidateObject(Object{Name: "Lone", Generators: []GeneratorOption{WithObjectAssert(ObjectAssertShow(ObjectAssertShowRowFromBase()))}})
		})
	})
	t.Run("the base must generate the artifact", func(t *testing.T) {
		require.Panics(t, func() {
			ValidateObject(variant(showReuseBaseObject(WithObjectAssert(ObjectAssertShow())), WithSchema(SchemaShow(SchemaShowFromBase()))))
		})
		require.Panics(t, func() {
			ValidateObject(variant(showReuseBaseObject(WithSchema(SchemaShow())), WithObjectAssert(ObjectAssertShow(ObjectAssertShowRowFromBase()))))
		})
	})
	t.Run("a reusing variant cannot override show-side options", func(t *testing.T) {
		tables := Attribute{Name: "tables"}
		obj := func(opts ...GeneratorOption) Object {
			return Object{
				Name:       "Variant",
				VariantOf:  baseWithShow(),
				Attributes: []Attribute{tables},
				Show:       OutputSpec{Columns: []*Attribute{&tables}},
				Generators: opts,
			}
		}
		require.Panics(t, func() {
			ValidateObject(obj(WithSchema(SchemaShow(SchemaShowFromBase(), SchemaShowManualFields(&tables)))))
		})
		require.Panics(t, func() {
			ValidateObject(obj(WithSchema(SchemaShow(SchemaShowFromBase(), SchemaShowListFields(&tables)))))
		})
		require.Panics(t, func() {
			ValidateObject(obj(WithObjectAssert(ObjectAssertShow(ObjectAssertShowRowFromBase(), ObjectAssertShowSkipFields(&tables)))))
		})
		require.Panics(t, func() {
			ValidateObject(obj(WithObjectAssert(ObjectAssertShow(ObjectAssertShowRowFromBase(), ObjectAssertShowNestedField(&tables)))))
		})
	})
}

func Test_ValidateObject_showOutputAssertReuse(t *testing.T) {
	base := func(opts ...GeneratorOption) *Object { return &Object{Name: "Base", Generators: opts} }
	variant := func(b *Object, opts ...ShowOutputShowOption) Object {
		return Object{Name: "Variant", VariantOf: b, Generators: []GeneratorOption{WithShowOutputAssert(ShowOutputShow(opts...))}}
	}
	baseWithShowOutput := func() *Object { return base(WithShowOutputAssert(ShowOutputShow())) }

	t.Run("valid reuse resolves to the base", func(t *testing.T) {
		obj := variant(baseWithShowOutput(), ShowOutputShowRowFromBase())
		require.NotPanics(t, func() { ValidateObject(obj) })
		require.True(t, obj.ReusesBaseShowOutputAssertShow())
		require.Equal(t, "Base", obj.ShowOutputAssertShowOwnerName())
	})
	t.Run("without reuse the object owns its assertion", func(t *testing.T) {
		obj := variant(baseWithShowOutput())
		require.False(t, obj.ReusesBaseShowOutputAssertShow())
		require.Equal(t, "Variant", obj.ShowOutputAssertShowOwnerName())
	})
	t.Run("requires VariantOf", func(t *testing.T) {
		require.Panics(t, func() { ValidateObject(variant(nil, ShowOutputShowRowFromBase())) })
	})
	t.Run("the base must generate the assertion", func(t *testing.T) {
		require.Panics(t, func() { ValidateObject(variant(base(), ShowOutputShowRowFromBase())) })
	})
	t.Run("cannot set skip fields", func(t *testing.T) {
		tables := Attribute{Name: "tables"}
		obj := variant(baseWithShowOutput(), ShowOutputShowRowFromBase(), ShowOutputShowSkipFields(&tables))
		obj.Attributes = []Attribute{tables}
		obj.Show = OutputSpec{Columns: []*Attribute{&tables}}
		require.Panics(t, func() { ValidateObject(obj) })
	})
}

var widgetEnabled = Attribute{Name: "enabled", Type: AttrBool, NativeBool: true, SkipCreate: true, SetAfterCreate: true}

func setAfterCreateObject(change func(*Object)) Object {
	obj := Object{
		Name:       "Widget",
		Attributes: []Attribute{widgetEnabled},
		Alter: AlterSpec{Clauses: []AlterClause{
			AlterKeywordFor(&widgetEnabled, true, "RESUME"),
			AlterKeywordFor(&widgetEnabled, false, "SUSPEND"),
		}},
		ResourceHooks: []ResourceHookOption{WithAfterCreate(Generated())},
	}
	change(&obj)
	return obj
}

func Test_ValidateObject_setAfterCreate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		require.NotPanics(t, func() { ValidateObject(setAfterCreateObject(func(*Object) {})) })
	})

	t.Run("the false link is optional", func(t *testing.T) {
		require.NotPanics(t, func() {
			ValidateObject(setAfterCreateObject(func(o *Object) { o.Alter.Clauses = o.Alter.Clauses[:1] }))
		})
	})

	for name, tc := range map[string]struct {
		change func(*Object)
		panics string
	}{
		"SetAfterCreate without SkipCreate": {
			func(o *Object) { o.Attributes[0].SkipCreate = false },
			`object "Widget": attribute "enabled" is SetAfterCreate and needs SkipCreate; CREATE would set it`,
		},
		"SetAfterCreate on an attribute that is not a native bool": {
			func(o *Object) {
				o.Attributes = append(o.Attributes, Attribute{Name: "size", Type: AttrInt, SkipCreate: true, SetAfterCreate: true})
			},
			`object "Widget": attribute "size" is SetAfterCreate, which supports only NativeBool attributes for now; another type needs a new kind of AfterCreate step`,
		},
		"link to an attribute that is not a native bool": {
			func(o *Object) { o.Attributes[0].NativeBool = false },
			`object "Widget": ALTER keyword "RESUME" links "enabled", which is not a NativeBool attribute of the object`,
		},
		"SetAfterCreate without a link to true": {
			func(o *Object) { o.Alter.Clauses = []AlterClause{AlterKeywordFor(&widgetEnabled, false, "SUSPEND")} },
			`object "Widget": attribute "enabled" is SetAfterCreate and needs an ALTER keyword linked to true (AlterKeywordFor)`,
		},
		"link to an attribute that the object does not have": {
			func(o *Object) {
				o.Alter.Clauses = append(o.Alter.Clauses, AlterKeywordFor(&Attribute{Name: "paused", Type: AttrBool, NativeBool: true}, true, "PAUSE"))
			},
			`object "Widget": ALTER keyword "PAUSE" links "paused", which is not a NativeBool attribute of the object`,
		},
		"link on a clause that is not a keyword": {
			func(o *Object) {
				o.Alter.Clauses = append(o.Alter.Clauses, AlterClause{Kind: AlterClauseText, SQL: "REMOVE ITEM", Attribute: &widgetEnabled})
			},
			`object "Widget": ALTER clause "REMOVE ITEM" links an attribute but is not a keyword`,
		},
		"SetAfterCreate without the generated hook": {
			func(o *Object) { o.ResourceHooks = nil },
			`object "Widget": SetAfterCreate attributes require WithAfterCreate(Generated())`,
		},
		"generated hook without SetAfterCreate": {
			func(o *Object) { o.Attributes[0].SetAfterCreate = false },
			`object "Widget": WithAfterCreate(Generated()) needs a SetAfterCreate attribute`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			require.PanicsWithValue(t, tc.panics, func() { ValidateObject(setAfterCreateObject(tc.change)) })
		})
	}
}

func Test_Object_LinkedAlterKeyword(t *testing.T) {
	obj := setAfterCreateObject(func(*Object) {})
	resume, ok := obj.LinkedAlterKeyword("enabled", true)
	require.True(t, ok)
	require.Equal(t, "RESUME", resume.SQL)
	suspend, ok := obj.LinkedAlterKeyword("enabled", false)
	require.True(t, ok)
	require.Equal(t, "SUSPEND", suspend.SQL)
	_, ok = obj.LinkedAlterKeyword("comment", true)
	require.False(t, ok)
}

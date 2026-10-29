package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_PairedStructsFromColumns_matchesHandwrittenChain(t *testing.T) {
	state := NewEnum("FakeState", "FakeStates", "STARTED")
	family := NewEnum("FakeFamily", "FakeFamilies", "CPU")

	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	minNodes := ssot.Attribute{Name: "min_nodes", Type: ssot.AttrInt}
	autoResume := ssot.Attribute{Name: "auto_resume", Type: ssot.AttrBool}
	isDefault := ssot.Attribute{Name: "is_default", Type: ssot.AttrBool, StoredAsText: true}
	createdOn := ssot.Attribute{Name: "created_on", Type: ssot.AttrTime}
	stateAttr := ssot.Attribute{Name: "state", Type: ssot.AttrEnum, Enum: state}
	forApplication := ssot.Attribute{
		Name:           "for_application",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "AccountObjectIdentifier",
		OptionalInShow: true,
		OptionalInDb:   true,
		ShowName:       "application",
	}
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString, OptionalInDb: true}
	droppedOn := ssot.Attribute{Name: "dropped_on", Type: ssot.AttrTime, OptionalInDb: true}
	options := ssot.Attribute{Name: "options", Type: ssot.AttrString, OptionalInShow: true, OptionalInDb: true}
	backup := ssot.Attribute{
		Name: "backup_instance_families",
		Type: ssot.AttrList,
		Elem: &ssot.Attribute{Type: ssot.AttrEnum, Enum: family},
	}
	kind := ssot.Attribute{Name: "kind", Type: ssot.AttrString}

	obj := ssot.Object{
		Name:       "Fake",
		NamePlural: "Fakes",
		IdType:     "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithSdk(ssot.SdkShow(), ssot.SdkDescribe()),
		},
		Show: ssot.OutputSpec{Columns: []*ssot.Attribute{
			&name,
			&minNodes,
			&autoResume,
			&isDefault,
			&createdOn,
			&stateAttr,
			&forApplication,
			&comment,
			&droppedOn,
			&options,
			&backup,
		}},
		Describe: ssot.OutputSpec{Columns: []*ssot.Attribute{
			&createdOn,
			&name,
			&kind,
		}},
	}

	showGot := ShowPairedStructs(obj)
	showWant := StructPair("fakesRow", "Fake").
		Text("name").
		Number("min_nodes").
		Bool("auto_resume").
		BoolFromText("is_default").
		Time("created_on").
		Enum("state", state).
		Field("application", "sql.NullString", "*AccountObjectIdentifier", WithPlainFieldName("Application")).
		OptionalText("comment", WithRequiredInPlain()).
		OptionalTime("dropped_on", WithRequiredInPlain()).
		OptionalText("options").
		StringList("backup_instance_families")
	require.Equal(t, showWant, showGot)
	require.True(t, pairedFieldByColumn(showGot, "state").isEnum)
	gotAutoResume := pairedFieldByColumn(showGot, "auto_resume")
	require.Equal(t, StructPair("fakesRow", "Fake").Bool("auto_resume").fields[0], gotAutoResume)
	require.Equal(t, "bool", gotAutoResume.dbKind)
	require.Equal(t, "bool", gotAutoResume.plainKind)
	require.NotEqual(t, StructPair("fakesRow", "Fake").BoolFromText("auto_resume").fields[0], gotAutoResume)
	gotDefault := pairedFieldByColumn(showGot, "is_default")
	require.Equal(t, StructPair("fakesRow", "Fake").BoolFromText("is_default").fields[0], gotDefault)
	require.Equal(t, "string", gotDefault.dbKind)
	require.Equal(t, "bool", gotDefault.plainKind)
	require.NotEqual(t, StructPair("fakesRow", "Fake").Bool("is_default").fields[0], gotDefault)

	storedText := ssot.Attribute{Name: "database_name", Type: ssot.AttrString, StoredAsText: true}
	require.Equal(t, StructPair("fakesRow", "Fake").Text("database_name"), ShowPairedStructs(showPairedObject("AccountObjectIdentifier", &storedText)))

	describeGot := DescribePairedStructs(obj)
	describeWant := StructPair("fakeDescRow", "FakeDetails").
		Time("created_on").
		Text("name").
		Text("kind")
	require.Equal(t, describeWant, describeGot)

	withPlainOnly := obj
	withPlainOnly.Generators = []ssot.GeneratorOption{ssot.WithSdk(
		ssot.SdkShow(),
		ssot.SdkDescribe(
			ssot.SdkDescribePlainOnlyField("DatabaseName", "string"),
			ssot.SdkDescribePlainOnlyField("SchemaName", "string"),
		),
	)}
	require.Equal(t, describeWant.
		PlainOnlyField("DatabaseName", "string").
		PlainOnlyField("SchemaName", "string"),
		DescribePairedStructs(withPlainOnly),
	)
}

func Test_ShowPairedStructs_plainOnlyHookAndManualConvert(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	database := ssot.Attribute{Name: "database_name", Type: ssot.AttrString, OptionalInDb: true, ManualConvert: true}
	scopes := ssot.Attribute{Name: "target_scopes", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrEnum}}
	obj := ssot.Object{
		Name:       "Fake",
		NamePlural: "Fakes",
		Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{&name, &database, &scopes}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkShow(
			ssot.SdkShowResultFilterHook(),
			ssot.SdkShowPlainOnly(&scopes, "[]FakeScope"),
		))},
	}
	want := StructPair("fakesRow", "Fake").
		Text("name").
		OptionalText("database_name", WithRequiredInPlain(), WithManualConvert()).
		PlainOnlyField("TargetScopes", "[]FakeScope").
		WithShowResultFilterHook()
	require.Equal(t, want, ShowPairedStructs(obj))
}

func Test_PairedStructsFromColumns_panics(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	state := NewEnum("FakeState", "FakeStates", "STARTED")

	tests := []struct {
		name string
		call func()
	}{
		{
			name: "show option off",
			call: func() {
				obj := showPairedObject("AccountObjectIdentifier", &name)
				obj.Generators = []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe())}
				ShowPairedStructs(obj)
			},
		},
		{
			name: "show missing NamePlural",
			call: func() {
				obj := showPairedObject("AccountObjectIdentifier", &name)
				obj.NamePlural = ""
				ShowPairedStructs(obj)
			},
		},
		{
			name: "describe option off",
			call: func() {
				DescribePairedStructs(showPairedObject("AccountObjectIdentifier", &name))
			},
		},
		{
			name: "describe missing Name",
			call: func() {
				obj := showPairedObject("AccountObjectIdentifier", &name)
				obj.Name = ""
				obj.Generators = []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe())}
				obj.Describe = ssot.OutputSpec{Columns: []*ssot.Attribute{&name}}
				DescribePairedStructs(obj)
			},
		},
		{
			name: "empty columns",
			call: func() {
				ShowPairedStructs(showPairedObject("AccountObjectIdentifier"))
			},
		},
		{
			name: "nil column",
			call: func() {
				ShowPairedStructs(showPairedObject("AccountObjectIdentifier", &name, nil))
			},
		},
		{
			name: "OptionalInShow without OptionalInDb",
			call: func() {
				comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString, OptionalInShow: true}
				ShowPairedStructs(showPairedObject("AccountObjectIdentifier", &comment))
			},
		},
		{
			name: "list whose element is not an enum",
			call: func() {
				names := ssot.Attribute{Name: "names", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrString}}
				ShowPairedStructs(showPairedObject("AccountObjectIdentifier", &names))
			},
		},
		{
			name: "set of enum",
			call: func() {
				families := ssot.Attribute{
					Name: "families",
					Type: ssot.AttrSet,
					Elem: &ssot.Attribute{Type: ssot.AttrEnum, Enum: state},
				}
				ShowPairedStructs(showPairedObject("AccountObjectIdentifier", &families))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Panics(t, tt.call)
		})
	}
}

func Test_PairedStructsFromColumns_mapsWarehouseColumnShapes(t *testing.T) {
	size := NewEnum("WarehouseSize", "WarehouseSizes", "XSMALL")
	state := NewEnum("WarehouseState", "WarehouseStates", "STARTED")
	optionalInt := ssot.Attribute{Name: "min_cluster_count", Type: ssot.AttrInt, OptionalInShow: true, OptionalInDb: true}
	dbOnlyInt := ssot.Attribute{Name: "auto_suspend", Type: ssot.AttrInt, OptionalInDb: true}
	optionalBool := ssot.Attribute{Name: "enable_query_acceleration", Type: ssot.AttrBool, OptionalInShow: true, OptionalInDb: true}
	optionalBoolText := ssot.Attribute{Name: "flag", Type: ssot.AttrBool, OptionalInShow: true, OptionalInDb: true, StoredAsText: true}
	optionalEnum := ssot.Attribute{Name: "size", Type: ssot.AttrEnum, Enum: size, OptionalInShow: true, OptionalInDb: true}
	manualEnum := ssot.Attribute{Name: "state", Type: ssot.AttrEnum, Enum: state, ManualConvert: true}
	manualOptionalEnum := ssot.Attribute{Name: "resource_constraint", Type: ssot.AttrEnum, Enum: size, OptionalInShow: true, OptionalInDb: true, ManualConvert: true}
	available := ssot.Attribute{Name: "available", Type: ssot.AttrFloat}
	monitor := ssot.Attribute{Name: "resource_monitor", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier", CustomParser: "ParseAccountObjectIdentifierExcludingExplicitNullString"}
	tables := ssot.Attribute{
		Name:         "tables",
		Type:         ssot.AttrList,
		OptionalInDb: true,
		CustomParser: "ParseCommaSeparatedSchemaObjectIdentifierArray",
		Elem:         &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "SchemaObjectIdentifier"},
	}
	objectTypes := ssot.Attribute{
		Name:         "object_types",
		Type:         ssot.AttrList,
		OptionalInDb: true,
		CustomParser: "ParseCommaSeparatedWarehouseStateArray",
		Elem:         &ssot.Attribute{Type: ssot.AttrEnum, Enum: state},
	}

	got := ShowPairedStructs(showPairedObject("AccountObjectIdentifier",
		&optionalInt, &dbOnlyInt, &optionalBool, &optionalBoolText, &optionalEnum, &manualEnum, &manualOptionalEnum, &available, &monitor, &tables, &objectTypes,
	))
	want := StructPair("fakesRow", "Fake").
		OptionalNumber("min_cluster_count").
		OptionalNumber("auto_suspend", WithRequiredInPlain()).
		OptionalBool("enable_query_acceleration").
		OptionalBoolFromText("flag").
		OptionalEnum("size", size).
		Field("state", "string", "WarehouseState", WithManualConvert()).
		OptionalEnum("resource_constraint", size, WithManualConvert()).
		Field("available", "string", "float64", WithManualConvert()).
		Field("resource_monitor", "string", "AccountObjectIdentifier", WithCustomParser("ParseAccountObjectIdentifierExcludingExplicitNullString")).
		Field("tables", "sql.NullString", "[]SchemaObjectIdentifier", WithCustomParser("ParseCommaSeparatedSchemaObjectIdentifierArray")).
		Field("object_types", "sql.NullString", "[]WarehouseState", WithCustomParser("ParseCommaSeparatedWarehouseStateArray"))
	require.Equal(t, want, got)

	requiredID := ssot.Attribute{Name: "database_name", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"}
	manualID := ssot.Attribute{Name: "name", Type: ssot.AttrIdentifier, IdentifierType: "ObjectIdentifier", ManualConvert: true}
	manualIDs := ssot.Attribute{
		Name:          "to",
		Type:          ssot.AttrList,
		ManualConvert: true,
		Elem:          &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountIdentifier"},
	}
	require.Equal(t,
		StructPair("fakesRow", "Fake").
			Field("database_name", "string", "AccountObjectIdentifier").
			Field("name", "string", "ObjectIdentifier", WithManualConvert()).
			Field("to", "string", "[]AccountIdentifier", WithManualConvert()),
		ShowPairedStructs(showPairedObject("AccountObjectIdentifier", &requiredID, &manualID, &manualIDs)),
	)
}

func showPairedObject(idType string, columns ...*ssot.Attribute) ssot.Object {
	return ssot.Object{
		Name:       "Fake",
		NamePlural: "Fakes",
		IdType:     idType,
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkShow())},
		Show:       ssot.OutputSpec{Columns: columns},
	}
}

func Test_PairedStructsFromColumns_boolParsed(t *testing.T) {
	lock := ssot.Attribute{Name: "has_retention_lock", Type: ssot.AttrBool, StoredAsText: true, BoolParsed: true}
	got := ShowPairedStructs(showPairedObject("SchemaObjectIdentifier", &lock))
	want := StructPair("fakesRow", "Fake").BoolFromText("has_retention_lock", WithBoolParsed())
	require.Equal(t, want, got)
}

func pairedFieldByColumn(pairs *PairedStructs, column string) pairedField {
	for _, field := range pairs.fields {
		if field.dbColumnName == column {
			return field
		}
	}
	panic("missing column " + column)
}

func Test_detailsGoType(t *testing.T) {
	family := ssot.NewEnum("WidgetFamily", "WidgetFamilies", "SMALL")
	obj := ssot.Object{Name: "Widget"}
	require.Equal(t, "[]WidgetFamily", detailsGoType(obj, ssot.Attribute{Name: "families", Type: ssot.AttrSet, Elem: &ssot.Attribute{Type: ssot.AttrEnum, Enum: family}}))
	require.Equal(t, "[]WidgetFamily", detailsGoType(obj, ssot.Attribute{Name: "families", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrEnum, Enum: family}}))
	require.Equal(t, "[]string", detailsGoType(obj, ssot.Attribute{Name: "names", Type: ssot.AttrSet, Elem: &ssot.Attribute{Type: ssot.AttrString}}))
	require.Equal(t, "WidgetFamily", detailsGoType(obj, ssot.Attribute{Name: "family", Type: ssot.AttrEnum, Enum: family}))
}

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
			name: "optional int",
			call: func() {
				minNodes := ssot.Attribute{Name: "min_nodes", Type: ssot.AttrInt, OptionalInDb: true}
				ShowPairedStructs(showPairedObject("AccountObjectIdentifier", &minNodes))
			},
		},
		{
			name: "optional bool",
			call: func() {
				autoResume := ssot.Attribute{Name: "auto_resume", Type: ssot.AttrBool, OptionalInShow: true, OptionalInDb: true}
				ShowPairedStructs(showPairedObject("AccountObjectIdentifier", &autoResume))
			},
		},
		{
			name: "optional bool with StoredAsText",
			call: func() {
				autoResume := ssot.Attribute{Name: "auto_resume", Type: ssot.AttrBool, OptionalInShow: true, OptionalInDb: true, StoredAsText: true}
				ShowPairedStructs(showPairedObject("AccountObjectIdentifier", &autoResume))
			},
		},
		{
			name: "optional enum",
			call: func() {
				stateAttr := ssot.Attribute{Name: "state", Type: ssot.AttrEnum, Enum: state, OptionalInDb: true}
				ShowPairedStructs(showPairedObject("AccountObjectIdentifier", &stateAttr))
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
		{
			name: "identifier missing IdType",
			call: func() {
				application := ssot.Attribute{
					Name:           "for_application",
					Type:           ssot.AttrIdentifier,
					OptionalInShow: true,
					OptionalInDb:   true,
					ShowName:       "application",
				}
				ShowPairedStructs(showPairedObject("", &application))
			},
		},
		{
			name: "enum missing Enum",
			call: func() {
				stateAttr := ssot.Attribute{Name: "state", Type: ssot.AttrEnum}
				ShowPairedStructs(showPairedObject("AccountObjectIdentifier", &stateAttr))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Panics(t, tt.call)
		})
	}
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

func pairedFieldByColumn(pairs *PairedStructs, column string) pairedField {
	for _, field := range pairs.fields {
		if field.dbColumnName == column {
			return field
		}
	}
	panic("missing column " + column)
}

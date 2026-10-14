package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_ShowOutputAssertionsFromSsot_rendersColumnAssertions(t *testing.T) {
	family := ssot.NewEnum("PoolFamily", "PoolFamilies", "SMALL")
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	minNodes := ssot.Attribute{Name: "min_nodes", Type: ssot.AttrInt}
	autoResume := ssot.Attribute{Name: "auto_resume", Type: ssot.AttrBool}
	createdOn := ssot.Attribute{Name: "created_on", Type: ssot.AttrTime}
	instanceFamily := ssot.Attribute{Name: "instance_family", Type: ssot.AttrEnum, Enum: family}
	owner := ssot.Attribute{Name: "owner", Type: ssot.AttrIdentifier}
	forApplication := ssot.Attribute{Name: "for_application", Type: ssot.AttrIdentifier, OptionalInShow: true, ShowName: "application"}
	backupFamilies := ssot.Attribute{Name: "backup_instance_families", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrEnum, Enum: family}}
	columns := []*ssot.Attribute{
		&comment, &minNodes, &autoResume, &createdOn, &instanceFamily, &owner, &forApplication, &backupFamilies,
	}
	obj := ssot.Object{
		Name:   "Pool",
		IdType: "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithShowOutputAssert(),
			ssot.WithDescribeOutputAssert(),
		},
		Show:     ssot.OutputSpec{Columns: columns},
		Describe: ssot.OutputSpec{Columns: columns},
	}

	showDetails := showOutputDetailsFromSsot(obj, false)
	require.Equal(t, "sdk.Pool", showDetails.ObjectName())
	require.False(t, showDetails.IsDataSourceOutput)

	show := ShowOutputAssertionsFromSsot(obj, false)
	require.Equal(t, "Pool", show.Name)
	require.False(t, show.IsDescribeOutput)
	assertions := renderShowOutputAssertions(t, show)
	require.Contains(t, assertions, "func (p *PoolShowOutputAssert) HasComment(expected string) *PoolShowOutputAssert {")
	require.Contains(t, assertions, `p.StringValueSet("comment", expected)`)
	require.Contains(t, assertions, "func (p *PoolShowOutputAssert) HasMinNodes(expected int) *PoolShowOutputAssert {")
	require.Contains(t, assertions, `p.IntValueSet("min_nodes", expected)`)
	require.Contains(t, assertions, "func (p *PoolShowOutputAssert) HasAutoResume(expected bool) *PoolShowOutputAssert {")
	require.Contains(t, assertions, `p.BoolValueSet("auto_resume", expected)`)
	require.NotContains(t, assertions, "HasAutoResumeBool")
	require.Contains(t, assertions, "func (p *PoolShowOutputAssert) HasCreatedOn(expected time.Time) *PoolShowOutputAssert {")
	require.Contains(t, assertions, `p.StringValueSet("created_on", expected.String())`)
	require.Contains(t, assertions, "func (p *PoolShowOutputAssert) HasInstanceFamily(expected sdk.PoolFamily) *PoolShowOutputAssert {")
	require.Contains(t, assertions, `p.StringValueSet("instance_family", string(expected))`)
	require.NotContains(t, assertions, "HasInstanceFamilySdk")
	require.NotContains(t, assertions, "HasInstanceFamilyString")
	require.Contains(t, assertions, "func (p *PoolShowOutputAssert) HasOwner(expected sdk.AccountObjectIdentifier) *PoolShowOutputAssert {")
	require.Contains(t, assertions, `p.StringValueSet("owner", expected.Name())`)
	require.Contains(t, assertions, "func (p *PoolShowOutputAssert) HasApplication(expected sdk.AccountObjectIdentifier) *PoolShowOutputAssert {")
	require.Contains(t, assertions, `p.StringValueSet("application", expected.Name())`)
	require.NotContains(t, assertions, "HasApplicationSdk")
	require.Contains(t, assertions, "func (p *PoolShowOutputAssert) HasNoBackupInstanceFamilies() *PoolShowOutputAssert {")
	require.Contains(t, assertions, `p.ValueSet("backup_instance_families.#", "0")`)
	require.NotContains(t, assertions, "HasBackupInstanceFamilies(")

	descDetails := showOutputDetailsFromSsot(obj, true)
	require.Equal(t, "sdk.PoolDetails", descDetails.ObjectName())
	require.True(t, descDetails.IsDataSourceOutput)
	desc := ShowOutputAssertionsFromSsot(obj, true)
	require.Equal(t, "PoolDetails", desc.Name)
	require.True(t, desc.IsDescribeOutput)
	descAssertions := renderShowOutputAssertions(t, desc)
	require.Contains(t, descAssertions, "func (p *PoolDescribeOutputAssert) HasInstanceFamily(expected sdk.PoolFamily) *PoolDescribeOutputAssert {")
	require.Contains(t, descAssertions, `p.StringValueSet("instance_family", string(expected))`)
	require.NotContains(t, descAssertions, "HasBackupInstanceFamilies(")
}

func Test_ShowOutputAssertionsFromSsot_refusesDisabledOrEmpty(t *testing.T) {
	column := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	require.Panics(t, func() {
		ShowOutputAssertionsFromSsot(ssot.Object{
			Name:   "Pool",
			IdType: "AccountObjectIdentifier",
			Show:   ssot.OutputSpec{Columns: []*ssot.Attribute{&column}},
		}, false)
	})
	require.Panics(t, func() {
		ShowOutputAssertionsFromSsot(ssot.Object{
			Name:       "Pool",
			IdType:     "AccountObjectIdentifier",
			Generators: []ssot.GeneratorOption{ssot.WithShowOutputAssert()},
		}, false)
	})
	require.Panics(t, func() {
		ShowOutputAssertionsFromSsot(ssot.Object{
			Name:       "Pool",
			IdType:     "AccountObjectIdentifier",
			Generators: []ssot.GeneratorOption{ssot.WithShowOutputAssert()},
			Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{nil}},
		}, false)
	})
	require.Panics(t, func() {
		ShowOutputAssertionsFromSsot(ssot.Object{
			Name:       "Pool",
			IdType:     "AccountObjectIdentifier",
			Generators: []ssot.GeneratorOption{ssot.WithSchemaShow()},
			Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{&column}},
		}, false)
	})
}

func Test_MapToResourceShowOutputAssertion_reflectedField(t *testing.T) {
	field := genhelpers.Field{Name: "WarehouseSize", ConcreteType: "string", UnderlyingType: "string"}
	model := ResourceShowOutputAssertionsModel{
		Name:       "Warehouse",
		Attributes: []ResourceShowOutputAssertionModel{MapToResourceShowOutputAssertion(field)},
	}
	assertions := renderShowOutputAssertions(t, model)
	require.Contains(t, assertions, "func (w *WarehouseShowOutputAssert) HasWarehouseSize(expected string) *WarehouseShowOutputAssert {")
	require.Contains(t, assertions, `w.StringValueSet("warehouse_size", expected)`)
	require.NotContains(t, assertions, "HasWarehouseSizeSdk")
	require.NotContains(t, assertions, "HasWarehouseSizeBool")
}

func renderShowOutputAssertions(t *testing.T, model ResourceShowOutputAssertionsModel) string {
	t.Helper()
	var buf strings.Builder
	require.NoError(t, AssertionsTemplate.Execute(&buf, model))
	return buf.String()
}

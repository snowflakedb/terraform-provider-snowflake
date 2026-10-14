package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_showResultSchemaDetailsFromSsot_synthesizesFields(t *testing.T) {
	family := ssot.NewEnum("PoolFamily", "PoolFamilies", "SMALL")
	stateEnum := ssot.NewEnum("PoolState", "PoolStates", "IDLE")
	forApplication := ssot.Attribute{Name: "for_application", Type: ssot.AttrIdentifier, OptionalInShow: true, OptionalInDb: true, ShowName: "application"}
	minNodes := ssot.Attribute{Name: "min_nodes", Type: ssot.AttrInt}
	instanceFamily := ssot.Attribute{Name: "instance_family", Type: ssot.AttrEnum, Enum: family}
	backupFamilies := ssot.Attribute{Name: "backup_families", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrEnum, Enum: family}}
	autoResume := ssot.Attribute{Name: "auto_resume", Type: ssot.AttrBool}
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString, OptionalInShow: true, OptionalInDb: true}
	tag := ssot.Attribute{Name: "tag", Type: ssot.AttrString, OptionalInDb: true}
	state := ssot.Attribute{Name: "state", Type: ssot.AttrEnum, Enum: stateEnum, Readonly: true}
	createdOn := ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true}
	updatedOn := ssot.Attribute{Name: "updated_on", Type: ssot.AttrTime, OptionalInShow: true, Readonly: true}
	owner := ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true}
	errorCode := ssot.Attribute{Name: "error_code", Type: ssot.AttrString, Readonly: true}
	obj := ssot.Object{
		Name:   "Pool",
		IdType: "AccountObjectIdentifier",
		Attributes: []ssot.Attribute{
			forApplication, minNodes, instanceFamily, backupFamilies, autoResume, comment, tag,
		},
		Show: ssot.OutputSpec{
			Columns: []*ssot.Attribute{
				&forApplication, &minNodes, &instanceFamily, &backupFamilies, &autoResume, &comment, &tag,
				&state, &createdOn, &updatedOn, &owner,
			},
		},
		Describe: ssot.OutputSpec{
			Columns: []*ssot.Attribute{
				&forApplication, &minNodes, &instanceFamily, &backupFamilies, &autoResume, &comment, &tag,
				&state, &createdOn, &updatedOn, &owner, &errorCode,
			},
		},
	}

	show := showResultSchemaDetailsFromSsot(obj, false)
	require.Equal(t, "sdk.Pool", show.Name)
	require.False(t, show.IsDescribe)
	require.Equal(t, []genhelpers.Field{
		{Name: "Application", ConcreteType: "*sdk.AccountObjectIdentifier", UnderlyingType: "*struct"},
		{Name: "MinNodes", ConcreteType: "int", UnderlyingType: "int"},
		{Name: "InstanceFamily", ConcreteType: "sdk.PoolFamily", UnderlyingType: "string"},
		{Name: "BackupFamilies", ConcreteType: "[]string", UnderlyingType: "slice"},
		{Name: "AutoResume", ConcreteType: "bool", UnderlyingType: "bool"},
		{Name: "Comment", ConcreteType: "*string", UnderlyingType: "*string"},
		{Name: "Tag", ConcreteType: "string", UnderlyingType: "string"},
		{Name: "State", ConcreteType: "sdk.PoolState", UnderlyingType: "string"},
		{Name: "CreatedOn", ConcreteType: "time.Time", UnderlyingType: "struct"},
		{Name: "UpdatedOn", ConcreteType: "*time.Time", UnderlyingType: "*struct"},
		{Name: "Owner", ConcreteType: "string", UnderlyingType: "string"},
	}, show.Fields)

	desc := showResultSchemaDetailsFromSsot(obj, true)
	require.Equal(t, "sdk.PoolDetails", desc.Name)
	require.True(t, desc.IsDescribe)
	require.Equal(t, append(append([]genhelpers.Field{}, show.Fields...), genhelpers.Field{
		Name: "ErrorCode", ConcreteType: "string", UnderlyingType: "string",
	}), desc.Fields)
}

func Test_showResultSchemaDetailsFromSsot_collectionsAndDescribeColumns(t *testing.T) {
	names := ssot.Attribute{Name: "names", Type: ssot.AttrSet, Elem: &ssot.Attribute{Type: ssot.AttrString}}
	refs := ssot.Attribute{Name: "refs", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrIdentifier}}
	counts := ssot.Attribute{Name: "counts", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrInt}}
	status := ssot.Attribute{Name: "status", Type: ssot.AttrString, Readonly: true}
	obj := ssot.Object{
		Name:       "Item",
		IdType:     "SchemaObjectIdentifier",
		Attributes: []ssot.Attribute{names, refs, counts},
		Show: ssot.OutputSpec{
			Columns: []*ssot.Attribute{&names, &refs, &counts},
		},
		Describe: ssot.OutputSpec{
			Columns: []*ssot.Attribute{&refs, &status},
		},
	}

	show := showResultSchemaDetailsFromSsot(obj, false)
	require.Equal(t, []genhelpers.Field{
		{Name: "Names", ConcreteType: "[]string", UnderlyingType: "slice"},
		{Name: "Refs", ConcreteType: "[]sdk.SchemaObjectIdentifier", UnderlyingType: "slice"},
		{Name: "Counts", ConcreteType: "[]int", UnderlyingType: "slice"},
	}, show.Fields)

	desc := showResultSchemaDetailsFromSsot(obj, true)
	require.Equal(t, "sdk.ItemDetails", desc.Name)
	require.Equal(t, []genhelpers.Field{
		{Name: "Refs", ConcreteType: "[]sdk.SchemaObjectIdentifier", UnderlyingType: "slice"},
		{Name: "Status", ConcreteType: "string", UnderlyingType: "string"},
	}, desc.Fields)
}

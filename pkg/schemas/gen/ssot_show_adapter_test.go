package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_showResultSchemaDetailsFromSsot_usesColumnsAndOptionalInShow(t *testing.T) {
	forApplication := ssot.Attribute{Name: "for_application", Type: ssot.AttrIdentifier, OptionalInShow: true, ShowName: "application"}
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString, OptionalInShow: true, OptionalInDb: true}
	tag := ssot.Attribute{Name: "tag", Type: ssot.AttrString, OptionalInDb: true}
	updatedOn := ssot.Attribute{Name: "updated_on", Type: ssot.AttrTime, OptionalInShow: true}
	owner := ssot.Attribute{Name: "owner", Type: ssot.AttrString}
	errorCode := ssot.Attribute{Name: "error_code", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:   "Pool",
		IdType: "AccountObjectIdentifier",
		Show: ssot.OutputSpec{
			Columns: []*ssot.Attribute{&forApplication, &comment, &tag, &updatedOn, &owner},
		},
		Describe: ssot.OutputSpec{
			Columns: []*ssot.Attribute{&comment, &errorCode},
		},
	}

	show := showResultSchemaDetailsFromSsot(obj, false)
	require.Equal(t, "sdk.Pool", show.Name)
	require.False(t, show.IsDescribe)
	require.Equal(t, []string{"Application", "Comment", "Tag", "UpdatedOn", "Owner"}, fieldNames(show.Fields))
	requirePointer(t, fieldByName(t, show.Fields, "Application"))
	requirePointer(t, fieldByName(t, show.Fields, "Comment"))
	requirePlain(t, fieldByName(t, show.Fields, "Tag"))
	requirePointer(t, fieldByName(t, show.Fields, "UpdatedOn"))
	requirePlain(t, fieldByName(t, show.Fields, "Owner"))

	desc := showResultSchemaDetailsFromSsot(obj, true)
	require.Equal(t, "sdk.PoolDetails", desc.Name)
	require.True(t, desc.IsDescribe)
	require.Equal(t, []string{"Comment", "ErrorCode"}, fieldNames(desc.Fields))
	requirePointer(t, fieldByName(t, desc.Fields, "Comment"))
	requirePlain(t, fieldByName(t, desc.Fields, "ErrorCode"))
}

func fieldNames(fields []genhelpers.Field) []string {
	names := make([]string, len(fields))
	for i, field := range fields {
		names[i] = field.Name
	}
	return names
}

func fieldByName(t *testing.T, fields []genhelpers.Field, name string) genhelpers.Field {
	t.Helper()
	for _, field := range fields {
		if field.Name == name {
			return field
		}
	}
	t.Fatalf("field %q not found", name)
	return genhelpers.Field{}
}

func requirePointer(t *testing.T, field genhelpers.Field) {
	t.Helper()
	require.Truef(t, strings.HasPrefix(field.ConcreteType, "*") && strings.HasPrefix(field.UnderlyingType, "*"), "%s concrete %q underlying %q", field.Name, field.ConcreteType, field.UnderlyingType)
}

func requirePlain(t *testing.T, field genhelpers.Field) {
	t.Helper()
	require.Falsef(t, strings.HasPrefix(field.ConcreteType, "*") || strings.HasPrefix(field.UnderlyingType, "*"), "%s concrete %q underlying %q", field.Name, field.ConcreteType, field.UnderlyingType)
}

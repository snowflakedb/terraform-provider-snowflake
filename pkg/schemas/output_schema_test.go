package schemas

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_OutputListSchema(t *testing.T) {
	elem := map[string]*schema.Schema{"name": {Type: schema.TypeString, Computed: true}}

	got := OutputListSchema("show_output", "Holds the output of SHOW THINGS.", elem)

	require.Len(t, got, 1)
	attribute := got["show_output"]
	require.NotNil(t, attribute)
	assert.Equal(t, schema.TypeList, attribute.Type)
	assert.True(t, attribute.Computed)
	assert.Equal(t, "Holds the output of SHOW THINGS.", attribute.Description)
	assert.Equal(t, elem, attribute.Elem.(*schema.Resource).Schema)
}

func Test_ShowOutputListSchema(t *testing.T) {
	elem := map[string]*schema.Schema{"name": {Type: schema.TypeString, Computed: true}}

	got := ShowOutputListSchema(sdk.ObjectTypeBackupPolicy, elem)

	require.Len(t, got, 1)
	attribute := got[showOutputAttributeName]
	require.NotNil(t, attribute)
	assert.Equal(t, schema.TypeList, attribute.Type)
	assert.True(t, attribute.Computed)
	assert.Equal(t, "Holds the output of `SHOW BACKUP POLICIES` for the given backup policy.", attribute.Description)
	assert.Equal(t, elem, attribute.Elem.(*schema.Resource).Schema)

	withExtra := ShowOutputListSchema(sdk.ObjectTypeBackupPolicy, elem, "Extra.")
	assert.Equal(t, "Holds the output of `SHOW BACKUP POLICIES` for the given backup policy. Extra.", withExtra[showOutputAttributeName].Description)
}

func Test_DescribeOutputListSchema(t *testing.T) {
	elem := map[string]*schema.Schema{"name": {Type: schema.TypeString, Computed: true}}

	got := DescribeOutputListSchema(sdk.ObjectTypeBackupPolicy, elem)

	attribute := got[describeOutputAttributeName]
	require.NotNil(t, attribute)
	assert.Equal(t, "Holds the output of `DESCRIBE BACKUP POLICY` for the given backup policy.", attribute.Description)
}

func Test_ParametersOutputListSchema(t *testing.T) {
	elem := map[string]*schema.Schema{"name": {Type: schema.TypeString, Computed: true}}

	got := ParametersOutputListSchema(sdk.ObjectTypeWarehouse, elem)

	attribute := got[parametersAttributeName]
	require.NotNil(t, attribute)
	assert.Equal(t, "Holds the output of `SHOW PARAMETERS FOR WAREHOUSE` for the given warehouse.", attribute.Description)
}

func Test_FullyQualifiedNameOutputSchema(t *testing.T) {
	got := FullyQualifiedNameOutputSchema()

	require.Len(t, got, 1)
	attribute := got[fullyQualifiedNameAttributeName]
	require.NotNil(t, attribute)
	assert.Equal(t, FullyQualifiedNameSchema, attribute)
}

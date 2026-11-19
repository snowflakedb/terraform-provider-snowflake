package datasources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_outputListSchema(t *testing.T) {
	elem := map[string]*schema.Schema{"name": {Type: schema.TypeString, Computed: true}}

	got := outputListSchema("show_output", "Holds the output of SHOW THINGS.", elem)

	require.Len(t, got, 1)
	attribute := got["show_output"]
	require.NotNil(t, attribute)
	assert.Equal(t, schema.TypeList, attribute.Type)
	assert.True(t, attribute.Computed)
	assert.Equal(t, "Holds the output of SHOW THINGS.", attribute.Description)
	assert.Equal(t, elem, attribute.Elem.(*schema.Resource).Schema)
}

func Test_showOutputListSchema(t *testing.T) {
	elem := map[string]*schema.Schema{"name": {Type: schema.TypeString, Computed: true}}

	got := showOutputListSchema(sdk.ObjectTypeBackupPolicy.Plural(), elem)

	require.Len(t, got, 1)
	attribute := got[resources.ShowOutputAttributeName]
	require.NotNil(t, attribute)
	assert.Equal(t, schema.TypeList, attribute.Type)
	assert.True(t, attribute.Computed)
	assert.Equal(t, "Holds the output of SHOW BACKUP POLICIES.", attribute.Description)
	assert.Equal(t, elem, attribute.Elem.(*schema.Resource).Schema)

	withExtra := showOutputListSchema(sdk.ObjectTypeBackupPolicy.Plural(), elem, "Extra.")
	assert.Equal(t, "Holds the output of SHOW BACKUP POLICIES. Extra.", withExtra[resources.ShowOutputAttributeName].Description)
}

func Test_describeOutputListSchema(t *testing.T) {
	elem := map[string]*schema.Schema{"name": {Type: schema.TypeString, Computed: true}}

	got := describeOutputListSchema(sdk.ObjectTypeBackupPolicy, elem)

	attribute := got[resources.DescribeOutputAttributeName]
	require.NotNil(t, attribute)
	assert.Equal(t, "Holds the output of DESCRIBE BACKUP POLICY.", attribute.Description)
}

func Test_parametersOutputListSchema(t *testing.T) {
	elem := map[string]*schema.Schema{"name": {Type: schema.TypeString, Computed: true}}

	got := parametersOutputListSchema(sdk.ObjectTypeWarehouse, elem)

	attribute := got[resources.ParametersAttributeName]
	require.NotNil(t, attribute)
	assert.Equal(t, "Holds the output of SHOW PARAMETERS FOR WAREHOUSE.", attribute.Description)
}

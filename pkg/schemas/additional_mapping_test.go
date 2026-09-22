package schemas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MergeSchema(t *testing.T) {
	generated := map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Computed: true},
	}
	additional := map[string]*schema.Schema{
		"is_ha": {Type: schema.TypeBool, Computed: true},
	}

	merged := mergeSchema(generated, additional)

	require.Len(t, merged, 2)
	assert.Equal(t, schema.TypeString, merged["name"].Type)
	assert.Equal(t, schema.TypeBool, merged["is_ha"].Type)
}

func Test_MergeSchema_CollisionPanics(t *testing.T) {
	generated := map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Computed: true},
	}
	additional := map[string]*schema.Schema{
		"name": {Type: schema.TypeBool, Computed: true},
	}

	assert.PanicsWithValue(t, `additionalSchema key "name" collides with generated schema`, func() {
		mergeSchema(generated, additional)
	})
}

package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlattenNotificationConfiguration(t *testing.T) {
	t.Run("full configuration", func(t *testing.T) {
		jsonStr := `{"NOTIFY_ADMINS":true,"SEVERITY_THRESHOLD":"High","USERS":["USER1","USER2"]}`
		result, err := flattenNotificationConfiguration(jsonStr)
		require.NoError(t, err)
		require.Len(t, result, 1)
		m := result[0].(map[string]interface{})
		assert.Equal(t, true, m["notify_admins"])
		assert.Equal(t, "High", m["severity_threshold"])
		assert.Equal(t, []string{"USER1", "USER2"}, m["users"])
	})

	t.Run("partial configuration", func(t *testing.T) {
		jsonStr := `{"SEVERITY_THRESHOLD":"Critical"}`
		result, err := flattenNotificationConfiguration(jsonStr)
		require.NoError(t, err)
		require.Len(t, result, 1)
		m := result[0].(map[string]interface{})
		assert.Equal(t, false, m["notify_admins"])
		assert.Equal(t, "Critical", m["severity_threshold"])
	})

	t.Run("empty string", func(t *testing.T) {
		result, err := flattenNotificationConfiguration("")
		require.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := flattenNotificationConfiguration("not valid json")
		assert.Error(t, err)
	})
}

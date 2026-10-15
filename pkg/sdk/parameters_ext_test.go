package sdk

import "testing"

func TestParameters_ShowParameters(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		assertOptsValidAndSqlEquals(t, &ShowParametersOptions{}, "SHOW PARAMETERS")
	})

	t.Run("like", func(t *testing.T) {
		opts := &ShowParametersOptions{
			Like: &Like{Pattern: String("test")},
		}
		assertOptsValidAndSqlEquals(t, opts, "SHOW PARAMETERS LIKE 'test'")
	})

	t.Run("in account", func(t *testing.T) {
		opts := &ShowParametersOptions{
			In: &ParametersIn{Account: Bool(true)},
		}
		assertOptsValidAndSqlEquals(t, opts, "SHOW PARAMETERS IN ACCOUNT")
	})

	t.Run("like and in account", func(t *testing.T) {
		opts := &ShowParametersOptions{
			Like: &Like{Pattern: String("test")},
			In:   &ParametersIn{Account: Bool(true)},
		}
		assertOptsValidAndSqlEquals(t, opts, "SHOW PARAMETERS LIKE 'test' IN ACCOUNT")
	})
}

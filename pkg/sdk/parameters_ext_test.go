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

	t.Run("in pipe", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &ShowParametersOptions{
			In: &ParametersIn{Pipe: id},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW PARAMETERS IN PIPE %s", id.FullyQualifiedName())
	})

	t.Run("like log event level in pipe", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &ShowParametersOptions{
			Like: &Like{Pattern: String("LOG_EVENT_LEVEL")},
			In:   &ParametersIn{Pipe: id},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW PARAMETERS LIKE 'LOG_EVENT_LEVEL' IN PIPE %s", id.FullyQualifiedName())
	})

	t.Run("no scope", func(t *testing.T) {
		opts := &ShowParametersOptions{
			In: &ParametersIn{},
		}
		assertOptsInvalid(t, opts, errAtLeastOneOf("Session", "Account", "User", "Warehouse", "Database", "Schema", "Task", "Table", "Function", "Procedure", "OpenflowDeployment", "Service", "Pipe"))
	})
}

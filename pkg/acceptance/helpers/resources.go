package helpers

func DummyResource() string {
	return `
resource snowflake_execute "t" {
    execute = "SELECT 1"
    query   = "SELECT 1"
    revert  = "SELECT 1"
}`
}

package snowflake

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func Exec(db *sql.DB, query string) error {
	_, err := db.Exec(query)
	return err
}

// QueryRow will run stmt against the db and return the row. We use
// [DB.Unsafe](https://godoc.org/github.com/jmoiron/sqlx#DB.Unsafe) so that we can scan to structs
// without worrying about newly introduced columns.
func QueryRow(db *sql.DB, stmt string) *sqlx.Row {
	return sqlx.NewDb(db, "snowflake").Unsafe().QueryRowx(stmt)
}

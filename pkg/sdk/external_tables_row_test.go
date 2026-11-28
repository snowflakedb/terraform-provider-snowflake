package sdk

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type externalTableRowTestDriver struct{ region driver.Value }

func (d externalTableRowTestDriver) Open(string) (driver.Conn, error)           { return d, nil }
func (d externalTableRowTestDriver) Prepare(string) (driver.Stmt, error)        { return d, nil }
func (d externalTableRowTestDriver) Close() error                               { return nil }
func (d externalTableRowTestDriver) Begin() (driver.Tx, error)                  { return nil, io.EOF }
func (d externalTableRowTestDriver) NumInput() int                              { return 0 }
func (d externalTableRowTestDriver) Exec([]driver.Value) (driver.Result, error) { return nil, io.EOF }
func (d externalTableRowTestDriver) Query([]driver.Value) (driver.Rows, error) {
	return &externalTableRowTestRows{region: d.region}, nil
}

type externalTableRowTestRows struct {
	region driver.Value
	done   bool
}

func (r *externalTableRowTestRows) Columns() []string {
	return []string{
		"created_on", "name", "database_name", "schema_name", "invalid", "invalid_reason", "owner", "comment", "stage",
		"location", "file_format_name", "file_format_type", "cloud", "region", "notification_channel",
		"last_refreshed_on", "table_format", "last_refresh_details", "owner_role_type",
	}
}

func (r *externalTableRowTestRows) Close() error { return nil }

func (r *externalTableRowTestRows) Next(dest []driver.Value) error {
	if r.done {
		return io.EOF
	}
	r.done = true
	copy(dest, []driver.Value{
		time.Now(), "ET", "DB", "SC", false, nil, "OWNER", "", "STAGE",
		"@stage/", "", nil, "AWS", r.region, nil, nil, "DELTA", nil, "ROLE",
	})
	return nil
}

func TestExternalTableRow_convert_region(t *testing.T) {
	for name, tc := range map[string]struct {
		region   driver.Value
		expected string
	}{
		"null region":  {region: nil, expected: ""},
		"set region":   {region: "US_EAST_1", expected: "US_EAST_1"},
		"empty region": {region: "", expected: ""},
	} {
		t.Run(name, func(t *testing.T) {
			driverName := "external_table_row_test_" + name
			sql.Register(driverName, externalTableRowTestDriver{region: tc.region})
			db, err := sqlx.Open(driverName, "")
			require.NoError(t, err)
			defer db.Close()

			var rows []externalTableRow
			require.NoError(t, db.Select(&rows, "SHOW EXTERNAL TABLES"))
			require.Len(t, rows, 1)

			externalTable, err := rows[0].convert()
			require.NoError(t, err)
			assert.Equal(t, tc.expected, externalTable.Region)
		})
	}
}

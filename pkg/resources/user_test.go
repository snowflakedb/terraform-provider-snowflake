package resources_test

import (
	"database/sql"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	r := require.New(t)
	err := resources.User().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestUserCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                 "good_name",
		"comment":              "great comment",
		"password":             "awesomepassword",
		"login_name":           "gname",
		"display_name":         "Display Name",
		"first_name":           "Marcin",
		"last_name":            "Zukowski",
		"email":                "fake@email.com",
		"disabled":             true,
		"default_warehouse":    "mywarehouse",
		"default_namespace":    "mynamespace",
		"default_role":         "bestrole",
		"rsa_public_key":       "asdf",
		"rsa_public_key_2":     "asdf2",
		"must_change_password": true,
	}
	d := schema.TestResourceDataRaw(t, resources.User().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		name := "good_name"
		q := fmt.Sprintf(`^CREATE USER "%s" COMMENT='great comment' DEFAULT_NAMESPACE='mynamespace' DEFAULT_ROLE='bestrole' DEFAULT_WAREHOUSE='mywarehouse' DISPLAY_NAME='Display Name' EMAIL='fake@email.com' FIRST_NAME='Marcin' LAST_NAME='Zukowski' LOGIN_NAME='gname' PASSWORD='awesomepassword' RSA_PUBLIC_KEY='asdf' RSA_PUBLIC_KEY_2='asdf2' DISABLED=true MUST_CHANGE_PASSWORD=true$`, name)
		mock.ExpectExec(q).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadUser(mock, name)
		err := resources.CreateUser(d, db)
		r.NoError(err)
	})
}

func expectReadUser(mock sqlmock.Sqlmock, name string) {
	rows := sqlmock.NewRows(
		[]string{"property", "value", "default", "description"},
	).AddRow(
		"NAME", name, "", "",
	).AddRow(
		"CREATED_ON", "created_on", "", "",
	).AddRow(
		"LOGIN_NAME", "myloginname", "", "",
	).AddRow(
		"DISPLAY_NAME", "display_name", "", "",
	).AddRow(
		"FIRST_NAME", "first_name", "", "",
	).AddRow(
		"LAST_NAME", "last_name", "", "",
	).AddRow(
		"EMAIL", "email", "", "",
	).AddRow(
		"MINS_TO_UNLOCK", "mins_to_unlock", "", "",
	).AddRow(
		"DAYS_TO_EXPIRY", "days_to_expiry", "", "",
	).AddRow(
		"COMMENT", "mock comment", "", "",
	).AddRow(
		"DISABLED", "false", "", "",
	).AddRow(
		"MUST_CHANGE_PASSWORD", "true", "", "",
	).AddRow(
		"SNOWFLAKE_LOCK", "snowflake_lock", "", "",
	).AddRow(
		"DEFAULT_WAREHOUSE", "default_warehouse", "", "",
	).AddRow(
		"DEFAULT_NAMESPACE", "default_namespace", "", "",
	).AddRow(
		"DEFAULT_ROLE", "default_role", "", "",
	).AddRow(
		"EXT_AUTHN_DUO", "ext_authn_duo", "", "",
	).AddRow(
		"EXT_AUTHN_UID", "ext_authn_uid", "", "",
	).AddRow(
		"MINS_TO_BYPASS_MFA", "mins_to_bypass_mfa", "", "",
	).AddRow(
		"OWNER", "owner", "", "",
	).AddRow(
		"LAST_SUCCESS_LOGIN", "last_success_login", "", "",
	).AddRow(
		"EXPIRES_AT_TIME", "expires_at_time", "", "",
	).AddRow(
		"LOCKED_UNTIL_TIME", "locked_until_time", "", "",
	).AddRow(
		"HAS_PASSWORD", "has_password", "", "",
	).AddRow(
		"HAS_RSA_PUBLIC_KEY", "false", "", "",
	)

	q := fmt.Sprintf(`^DESCRIBE USER "%s"$`, name)
	mock.ExpectQuery(q).WillReturnRows(rows)
}

func TestUserRead(t *testing.T) {
	r := require.New(t)
	name := "good_name"
	d := user(t, name, map[string]interface{}{"name": name})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadUser(mock, name)
		err := resources.ReadUser(d, db)
		r.NoError(err)
		r.Equal("mock comment", d.Get("comment").(string))
		r.Equal("myloginname", d.Get("login_name").(string))
		r.Equal(false, d.Get("disabled").(bool))

		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		q := snowflake.User(d.Id()).Describe()
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		err2 := resources.ReadUser(d, db)
		r.Empty(d.State())
		r.Nil(err2)
	})
}

func TestUserExists(t *testing.T) {
	r := require.New(t)
	name := "good_name"
	d := user(t, name, map[string]interface{}{"name": name})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadUser(mock, name)
		b, err := resources.UserExists(d, db)
		r.NoError(err)
		r.True(b)
	})
}

func TestUserDelete(t *testing.T) {
	r := require.New(t)

	d := user(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^DROP USER "drop_it"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteUser(d, db)
		r.NoError(err)
	})
}

package sdk

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	id := postgresInstancesTestIdAccountObjectIdentifier
	forkSourceId := randomAccountObjectIdentifier()
	renameTarget := randomAccountObjectIdentifier()
	networkPolicyId := randomAccountObjectIdentifier()
	storageIntegrationId := randomAccountObjectIdentifier()
	tagId := NewAccountObjectIdentifier("tag1")
	tagId2 := NewAccountObjectIdentifier("tag2")
	authAuthority := PostgresInstanceAuthenticationAuthorityPostgres

	postgresInstancesTests.Create.
		withDefaultOpts(func() *CreatePostgresInstanceOptions {
			return &CreatePostgresInstanceOptions{
				name:                    id,
				ComputeFamily:           "STANDARD_M",
				StorageSizeGb:           10,
				AuthenticationAuthority: PostgresInstanceAuthenticationAuthorityPostgres,
			}
		}).
		withExpectedSqlf(
			case_PostgresInstances_sql_Create_basic,
			"CREATE POSTGRES INSTANCE %s COMPUTE_FAMILY = 'STANDARD_M' STORAGE_SIZE_GB = 10 AUTHENTICATION_AUTHORITY = POSTGRES",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Create_all,
			func(opts *CreatePostgresInstanceOptions) {
				opts.PostgresVersion = new(15)
				opts.NetworkPolicy = &networkPolicyId
				opts.HighAvailability = new(true)
				opts.StorageIntegration = &storageIntegrationId
				opts.PostgresSettings = new(`{"max_connections":"100"}`)
				opts.Comment = new("my comment")
			},
			`CREATE POSTGRES INSTANCE %s COMPUTE_FAMILY = 'STANDARD_M' STORAGE_SIZE_GB = 10 AUTHENTICATION_AUTHORITY = POSTGRES POSTGRES_VERSION = 15 NETWORK_POLICY = %s HIGH_AVAILABILITY = true STORAGE_INTEGRATION = %s POSTGRES_SETTINGS = '{\"max_connections\":\"100\"}' COMMENT = 'my comment'`,
			id.FullyQualifiedName(), networkPolicyId.FullyQualifiedName(), storageIntegrationId.FullyQualifiedName(),
		)

	postgresInstancesTests.Fork.
		withDefaultOpts(func() *ForkPostgresInstanceOptions {
			return &ForkPostgresInstanceOptions{
				name: id,
				Fork: forkSourceId,
			}
		}).
		withModify(
			case_PostgresInstances_validation_Fork_opts_ConflictingFields,
			func(opts *ForkPostgresInstanceOptions) {
				opts.At = &PostgresInstanceForkAt{Timestamp: new("2023-01-01 00:00:00")}
				opts.Before = &PostgresInstanceForkBefore{Timestamp: new("2023-01-01 00:00:00")}
			},
		).
		withExpectedSqlf(
			case_PostgresInstances_sql_Fork_basic,
			"CREATE POSTGRES INSTANCE %s FORK %s",
			id.FullyQualifiedName(), forkSourceId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Fork_all",
			func(opts *ForkPostgresInstanceOptions) {
				opts.At = &PostgresInstanceForkAt{Timestamp: new("2023-01-01 00:00:00")}
				opts.ComputeFamily = new("STANDARD_M")
				opts.StorageSizeGb = new(10)
				opts.HighAvailability = new(true)
				opts.Comment = new("my fork")
			},
			"CREATE POSTGRES INSTANCE %s FORK %s AT (TIMESTAMP => '2023-01-01 00:00:00') COMPUTE_FAMILY = 'STANDARD_M' STORAGE_SIZE_GB = 10 HIGH_AVAILABILITY = true COMMENT = 'my fork'",
			id.FullyQualifiedName(), forkSourceId.FullyQualifiedName(),
		)

	postgresInstancesTests.Alter.
		withModify(
			case_PostgresInstances_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterPostgresInstanceOptions) {
				opts.RenameTo = &renameTarget
				opts.Set = &PostgresInstanceSet{Comment: new("my comment")}
			},
		).
		withModify(
			case_PostgresInstances_validation_Alter_opts_Set_Apply_ExactlyOneValueSet_NoneSet,
			func(opts *AlterPostgresInstanceOptions) {
				opts.Set = &PostgresInstanceSet{
					Comment: new("x"),
					Apply:   &PostgresInstanceApply{},
				}
			},
		).
		withModify(
			case_PostgresInstances_validation_Alter_opts_Set_Apply_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterPostgresInstanceOptions) {
				opts.Set = &PostgresInstanceSet{
					Comment: new("x"),
					Apply:   &PostgresInstanceApply{Immediately: new(true), On: new("foo")},
				}
			},
		).
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Alter_RenameTo,
			func(opts *AlterPostgresInstanceOptions) { opts.RenameTo = &renameTarget },
			"ALTER POSTGRES INSTANCE %s RENAME TO %s", id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Alter_Set,
			func(opts *AlterPostgresInstanceOptions) {
				opts.IfExists = new(true)
				opts.Set = &PostgresInstanceSet{
					NetworkPolicy:           &networkPolicyId,
					AuthenticationAuthority: &authAuthority,
					Comment:                 new("my comment"),
					HighAvailability:        new(true),
					ComputeFamily:           new("STANDARD_M"),
					StorageSizeGb:           new(10),
					StorageIntegration:      &storageIntegrationId,
					PostgresVersion:         new(15),
					MaintenanceWindowStart:  new(3),
					PostgresSettings:        new(`{"max_connections":"100"}`),
					Apply:                   &PostgresInstanceApply{Immediately: new(true)},
				}
			},
			`ALTER POSTGRES INSTANCE IF EXISTS %s SET NETWORK_POLICY = %s AUTHENTICATION_AUTHORITY = POSTGRES COMMENT = 'my comment' HIGH_AVAILABILITY = true COMPUTE_FAMILY = 'STANDARD_M' STORAGE_SIZE_GB = 10 STORAGE_INTEGRATION = %s POSTGRES_VERSION = 15 MAINTENANCE_WINDOW_START = 3 POSTGRES_SETTINGS = '{\"max_connections\":\"100\"}' APPLY IMMEDIATELY`,
			id.FullyQualifiedName(), networkPolicyId.FullyQualifiedName(), storageIntegrationId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Alter_Unset,
			func(opts *AlterPostgresInstanceOptions) {
				opts.Unset = &PostgresInstanceUnset{
					Comment:                new(true),
					PostgresSettings:       new(true),
					NetworkPolicy:          new(true),
					MaintenanceWindowStart: new(true),
					StorageIntegration:     new(true),
				}
			},
			"ALTER POSTGRES INSTANCE %s UNSET COMMENT, POSTGRES_SETTINGS, NETWORK_POLICY, MAINTENANCE_WINDOW_START, STORAGE_INTEGRATION",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Alter_Suspend,
			func(opts *AlterPostgresInstanceOptions) { opts.Suspend = new(true) },
			"ALTER POSTGRES INSTANCE %s SUSPEND", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Alter_Resume,
			func(opts *AlterPostgresInstanceOptions) { opts.Resume = new(true) },
			"ALTER POSTGRES INSTANCE %s RESUME", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Alter_ResetAccess,
			func(opts *AlterPostgresInstanceOptions) {
				opts.ResetAccess = &PostgresInstanceResetAccess{ForRole: PostgresInstanceResetAccessRoleSnowflakeAdmin}
			},
			"ALTER POSTGRES INSTANCE %s RESET ACCESS FOR 'snowflake_admin'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Alter_SetTags,
			func(opts *AlterPostgresInstanceOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tagId, Value: "value-123"},
					{Name: tagId2, Value: "value-123"},
				}
			},
			"ALTER POSTGRES INSTANCE %s SET TAG %s = 'value-123', %s = 'value-123'",
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Alter_UnsetTags,
			func(opts *AlterPostgresInstanceOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId, tagId2}
			},
			"ALTER POSTGRES INSTANCE %s UNSET TAG %s, %s",
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		)

	postgresInstancesTests.Drop.
		withExpectedSqlf(
			case_PostgresInstances_sql_Drop_basic,
			"DROP POSTGRES INSTANCE %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Drop_all,
			func(opts *DropPostgresInstanceOptions) { opts.IfExists = new(true) },
			"DROP POSTGRES INSTANCE IF EXISTS %s", id.FullyQualifiedName(),
		)

	postgresInstancesTests.Show.
		withExpectedSql(case_PostgresInstances_sql_Show_basic, "SHOW POSTGRES INSTANCES").
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Show_all,
			func(opts *ShowPostgresInstanceOptions) {
				opts.Like = &Like{Pattern: new("my-pattern")}
				opts.StartsWith = new("my-prefix")
				opts.Limit = &LimitFrom{Rows: new(10), From: new("my-from")}
			},
			"SHOW POSTGRES INSTANCES LIKE 'my-pattern' STARTS WITH 'my-prefix' LIMIT 10 FROM 'my-from'",
		).
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Show_Like,
			func(opts *ShowPostgresInstanceOptions) { opts.Like = &Like{Pattern: new("my-pattern")} },
			"SHOW POSTGRES INSTANCES LIKE 'my-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Show_StartsWith,
			func(opts *ShowPostgresInstanceOptions) { opts.StartsWith = new("my-prefix") },
			"SHOW POSTGRES INSTANCES STARTS WITH 'my-prefix'",
		).
		withModifyAndExpectedSqlf(
			case_PostgresInstances_sql_Show_Limit,
			func(opts *ShowPostgresInstanceOptions) {
				opts.Limit = &LimitFrom{Rows: new(10), From: new("my-from")}
			},
			"SHOW POSTGRES INSTANCES LIMIT 10 FROM 'my-from'",
		)

	postgresInstancesTests.Describe.
		withExpectedSqlf(
			case_PostgresInstances_sql_Describe_basic,
			"DESCRIBE POSTGRES INSTANCE %s", id.FullyQualifiedName(),
		)
}

func TestPostgresInstances_ParseDetails(t *testing.T) {
	t.Run("optional string fields: empty becomes nil, non-empty becomes pointer", func(t *testing.T) {
		tests := []struct {
			name      string
			property  string
			value     string
			wantNil   bool
			wantValue string
			getField  func(*PostgresInstanceDetails) *string
		}{
			{
				name:     "empty comment yields nil",
				property: "comment", value: "",
				wantNil:  true,
				getField: func(d *PostgresInstanceDetails) *string { return d.Comment },
			},
			{
				name:     "non-empty comment yields pointer to value",
				property: "comment", value: "my comment",
				wantValue: "my comment",
				getField:  func(d *PostgresInstanceDetails) *string { return d.Comment },
			},
			{
				name:     "empty postgres_settings yields nil",
				property: "postgres_settings", value: "",
				wantNil:  true,
				getField: func(d *PostgresInstanceDetails) *string { return d.PostgresSettings },
			},
			{
				name:     "non-empty postgres_settings yields pointer to value",
				property: "postgres_settings", value: `{"work_mem":"64KB"}`,
				wantValue: `{"work_mem":"64KB"}`,
				getField:  func(d *PostgresInstanceDetails) *string { return d.PostgresSettings },
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				properties := []PostgresInstanceProperty{
					{Property: "name", Value: "test_instance"},
					{Property: tc.property, Value: tc.value},
				}
				details, err := ParsePostgresInstanceDetails(properties)
				require.NoError(t, err)
				if tc.wantNil {
					require.Nil(t, tc.getField(details))
				} else {
					require.NotNil(t, tc.getField(details))
					assert.Equal(t, tc.wantValue, *tc.getField(details))
				}
			})
		}
	})

	t.Run("parse network policy into AccountObjectIdentifier", func(t *testing.T) {
		properties := []PostgresInstanceProperty{
			{Property: "name", Value: "test_instance"},
			{Property: "network_policy", Value: "my_network_policy"},
		}
		details, err := ParsePostgresInstanceDetails(properties)
		require.NoError(t, err)
		require.NotNil(t, details.NetworkPolicy)
		assert.Equal(t, NewAccountObjectIdentifier("my_network_policy"), *details.NetworkPolicy)
	})

	t.Run("parse storage integration into AccountObjectIdentifier", func(t *testing.T) {
		properties := []PostgresInstanceProperty{
			{Property: "name", Value: "test_instance"},
			{Property: "storage_integration", Value: "my_storage_integration"},
		}
		details, err := ParsePostgresInstanceDetails(properties)
		require.NoError(t, err)
		require.NotNil(t, details.StorageIntegration)
		assert.Equal(t, NewAccountObjectIdentifier("my_storage_integration"), *details.StorageIntegration)
	})

	t.Run("parse mixed-case property keys", func(t *testing.T) {
		properties := []PostgresInstanceProperty{
			{Property: "Name", Value: "test_instance"},
			{Property: "COMPUTE_FAMILY", Value: "STANDARD_M"},
			{Property: "Storage_Size_Gb", Value: "100"},
		}
		details, err := ParsePostgresInstanceDetails(properties)
		require.NoError(t, err)
		assert.Equal(t, "test_instance", details.Name)
		assert.Equal(t, "STANDARD_M", details.ComputeFamily)
		assert.Equal(t, 100, details.StorageSizeGb)
	})
}

func TestNormalizePostgresSettings(t *testing.T) {
	t.Run("empty and whitespace only", func(t *testing.T) {
		for _, s := range []string{"", "  ", "\t\n"} {
			got, err := NormalizePostgresSettings(s)
			require.NoError(t, err)
			require.Equal(t, "", got)
		}
	})

	t.Run("empty JSON object", func(t *testing.T) {
		got, err := NormalizePostgresSettings("{}")
		require.NoError(t, err)
		require.Equal(t, "", got)
	})

	t.Run("equivalent JSON with different formatting", func(t *testing.T) {
		want, err := NormalizePostgresSettings(`{"max_connections":"100","shared_buffers":"256MB"}`)
		require.NoError(t, err)

		equivalentForms := []string{
			`{"shared_buffers":"256MB","max_connections":"100"}`,
			`{  "max_connections"  :  "100"  ,  "shared_buffers"  :  "256MB"  }`,
			"{\n  \"max_connections\": \"100\",\n  \"shared_buffers\": \"256MB\"\n}",
		}
		for _, s := range equivalentForms {
			got, err := NormalizePostgresSettings(s)
			require.NoError(t, err)
			require.Equal(t, want, got)
		}
	})

	t.Run("non-equivalent JSON", func(t *testing.T) {
		want, err := NormalizePostgresSettings(`{"max_connections":"100"}`)
		require.NoError(t, err)

		got, err := NormalizePostgresSettings(`{"max_connections":"200"}`)
		require.NoError(t, err)
		require.NotEqual(t, want, got)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		_, err := NormalizePostgresSettings("{broken")
		require.Error(t, err)
	})
}

func TestNormalizePostgresSettingsPtr(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		require.Nil(t, NormalizePostgresSettingsPtr(nil))
	})

	t.Run("empty string returns nil", func(t *testing.T) {
		s := ""
		require.Nil(t, NormalizePostgresSettingsPtr(&s))
	})

	t.Run("empty JSON object returns nil", func(t *testing.T) {
		s := "{}"
		require.Nil(t, NormalizePostgresSettingsPtr(&s))
	})

	t.Run("valid JSON returns normalized pointer", func(t *testing.T) {
		s := `{"shared_buffers":"256MB","max_connections":"100"}`
		got := NormalizePostgresSettingsPtr(&s)
		require.NotNil(t, got)
		want, err := NormalizePostgresSettings(s)
		require.NoError(t, err)
		require.Equal(t, want, *got)
	})

	t.Run("invalid JSON returns nil", func(t *testing.T) {
		s := "{broken"
		require.Nil(t, NormalizePostgresSettingsPtr(&s))
	})
}

// stubPostgresInstances is a minimal test double for testing CreateSafely / updateSafely
// polling logic without a live SDK client.
type stubPostgresInstances struct {
	createErr  error
	showStates []PostgresInstanceState // sequence of states returned by successive ShowByID calls
	showIdx    int
	showErr    error

	// For updateSafelyPolling tests
	updateErr    error
	updateCalled int
}

func (s *stubPostgresInstances) showByID() (*PostgresInstance, error) {
	if s.showErr != nil {
		return nil, s.showErr
	}
	if s.showIdx >= len(s.showStates) {
		return &PostgresInstance{Name: "test", State: PostgresInstanceStateReady}, nil
	}
	state := s.showStates[s.showIdx]
	s.showIdx++
	return &PostgresInstance{Name: "test", State: state}, nil
}

func (s *stubPostgresInstances) update() error {
	s.updateCalled++
	return s.updateErr
}

func TestCreateSafely(t *testing.T) {
	t.Run("returns error when Create fails", func(t *testing.T) {
		createErr := errors.New("create failed")
		_, err := createSafelyPolling(context.Background(), func() error { return createErr }, nil)
		require.ErrorIs(t, err, createErr)
	})

	t.Run("returns instance when immediately READY", func(t *testing.T) {
		stub := &stubPostgresInstances{
			showStates: []PostgresInstanceState{PostgresInstanceStateReady},
		}
		instance, err := createSafelyPolling(context.Background(), func() error { return nil }, stub.showByID)
		require.NoError(t, err)
		assert.Equal(t, PostgresInstanceStateReady, instance.State)
	})

	t.Run("returns instance after polling through non-READY states", func(t *testing.T) {
		stub := &stubPostgresInstances{
			showStates: []PostgresInstanceState{
				PostgresInstanceStateCreating,
				PostgresInstanceStateCreating,
				PostgresInstanceStateReady,
			},
		}
		instance, err := createSafelyPolling(context.Background(), func() error { return nil }, stub.showByID)
		require.NoError(t, err)
		assert.Equal(t, PostgresInstanceStateReady, instance.State)
		assert.Equal(t, 3, stub.showIdx)
	})

	t.Run("returns error when context is canceled before READY", func(t *testing.T) {
		stub := &stubPostgresInstances{
			showStates: []PostgresInstanceState{PostgresInstanceStateCreating},
		}
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		time.Sleep(5 * time.Millisecond) // ensure deadline is already exceeded
		_, err := createSafelyPolling(ctx, func() error { return nil }, stub.showByID)
		require.Error(t, err)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("propagates ShowByID error", func(t *testing.T) {
		showErr := errors.New("show failed")
		stub := &stubPostgresInstances{showErr: showErr}
		_, err := createSafelyPolling(context.Background(), func() error { return nil }, stub.showByID)
		require.ErrorIs(t, err, showErr)
	})
}

func TestUpdateSafely(t *testing.T) {
	t.Run("calls update when instance is immediately READY", func(t *testing.T) {
		stub := &stubPostgresInstances{}
		err := updateSafelyPolling(context.Background(), stub.update, stub.showByID)
		require.NoError(t, err)
		assert.Equal(t, 1, stub.updateCalled)
	})

	t.Run("waits for READY state before calling update", func(t *testing.T) {
		stub := &stubPostgresInstances{
			showStates: []PostgresInstanceState{
				PostgresInstanceStateCreating,
				PostgresInstanceStateCreating,
				PostgresInstanceStateReady,
			},
		}
		err := updateSafelyPolling(context.Background(), stub.update, stub.showByID)
		require.NoError(t, err)
		assert.Equal(t, 1, stub.updateCalled)
	})

	t.Run("retries update on must-be-complete error", func(t *testing.T) {
		calls := 0
		doUpdate := func() error {
			calls++
			if calls < 3 {
				return fmt.Errorf("604009 (03000): Running operation CREATE POSTGRES SERVICE on X %s SET POSTGRES_SETTINGS", ErrPostgresOperationMustBeComplete.Error())
			}
			return nil
		}
		stub := &stubPostgresInstances{}
		err := updateSafelyPolling(context.Background(), doUpdate, stub.showByID)
		require.NoError(t, err)
		assert.Equal(t, 3, calls)
	})

	t.Run("returns non-retryable update error", func(t *testing.T) {
		updateErr := errors.New("unexpected error")
		stub := &stubPostgresInstances{updateErr: updateErr}
		err := updateSafelyPolling(context.Background(), stub.update, stub.showByID)
		require.ErrorIs(t, err, updateErr)
	})

	t.Run("returns error when context is canceled before READY", func(t *testing.T) {
		stub := &stubPostgresInstances{
			showStates: []PostgresInstanceState{PostgresInstanceStateCreating},
		}
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		time.Sleep(5 * time.Millisecond)
		err := updateSafelyPolling(ctx, stub.update, stub.showByID)
		require.Error(t, err)
		assert.Equal(t, 0, stub.updateCalled)
	})
}

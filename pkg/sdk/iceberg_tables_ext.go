package sdk

import "context"

func (r *CreateIcebergTableRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

func (v *icebergTables) ShowParameters(ctx context.Context, id SchemaObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Table: id,
		},
	})
}

// TypeString returns the best available textual representation of the column's data type:
// the parsed form when available, otherwise the raw DESC output.
func (d IcebergTableDetails) TypeString() string {
	if d.Type != nil {
		return d.Type.ToSql()
	}
	return d.DataTypeRaw
}

type IcebergTablePartitionSpec struct {
	SpecId int                              `json:"spec-id"`
	Fields []IcebergTablePartitionSpecField `json:"fields"`
}

type IcebergTablePartitionSpecField struct {
	Name      string `json:"name"`
	Transform string `json:"transform"`
	SourceId  int    `json:"source-id"`
	FieldId   int    `json:"field-id"`
}

// IcebergTableAutoRefreshStatus is the parsed form of the auto_refresh_status JSON column returned by
// SHOW ICEBERG TABLES. It is empty (the column is "") for tables without auto-refresh configured.
type IcebergTableAutoRefreshStatus struct {
	CurrentSnapshotId    int     `json:"currentSnapshotId"`
	LastSnapshotTime     *string `json:"lastSnapshotTime"`
	PendingSnapshotCount int     `json:"pendingSnapshotCount"`
	ExecutionState       string  `json:"executionState"`
	LastUpdatedTime      string  `json:"lastUpdatedTime"`
}

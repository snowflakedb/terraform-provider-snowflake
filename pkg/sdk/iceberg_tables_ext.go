package sdk

import "fmt"

func (opts *CreateIcebergTableOptions) additionalValidations() error {
	var errs []error
	// PartitionBy is a slice, validate each element
	for i, p := range opts.PartitionBy {
		if !exactlyOneValueSet(p.Identity, p.Bucket, p.Truncate, p.Year, p.Month, p.Day, p.Hour) {
			errs = append(errs, errExactlyOneOf(fmt.Sprintf("CreateIcebergTableOptions.PartitionBy[%d]", i), "Identity", "Bucket", "Truncate", "Year", "Month", "Day", "Hour"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterIcebergTableOptions) additionalValidations() error {
	var errs []error
	// AlterColumnAction is a slice, validate each element
	for i, col := range opts.AlterColumnAction {
		if !exactlyOneValueSet(col.SetNotNull, col.DropNotNull, col.DataType, col.Comment, col.UnsetComment, col.SetWriteDefault, col.DropWriteDefault) {
			errs = append(errs, errExactlyOneOf(fmt.Sprintf("AlterIcebergTableOptions.AlterColumnAction[%d]", i), "SetNotNull", "DropNotNull", "DataType", "Comment", "UnsetComment", "SetWriteDefault", "DropWriteDefault"))
		}
	}
	return JoinErrors(errs...)
}

func (d *TableDropSearchOptimization) additionalValidations() error {
	var errs []error
	// each Drop.On entry must have exactly one of SearchMethodWithTarget, ColumnName, or ExpressionId set.
	for _, on := range d.On {
		if !exactlyOneValueSet(on.SearchMethodWithTarget, on.ColumnName, on.ExpressionId) {
			errs = append(errs, errExactlyOneOf("AlterIcebergTableOptions.SearchOptimizationAction.Drop.On", "SearchMethodWithTarget", "ColumnName", "ExpressionId"))
		}
	}
	return JoinErrors(errs...)
}

// icebergTableExternalVolumeQuoted formats an AccountObjectIdentifier for the
// EXTERNAL_VOLUME clause of CREATE ICEBERG TABLE, which expects a single-quoted
// string literal whose content is the double-quoted volume name (e.g. '"vol1"').
//
// TODO(SNOW-2236323): Use a proper generation option instead.
// We need to use a custom parsing here, see SNOW-1833593 for more details.
func icebergTableExternalVolumeQuoted(id *AccountObjectIdentifier) *string {
	if id == nil {
		return nil
	}
	return Pointer(fmt.Sprintf("'%s'", id.FullyQualifiedName()))
}

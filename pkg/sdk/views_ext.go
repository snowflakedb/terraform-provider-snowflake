package sdk

import (
	"fmt"
	"strings"
)

func (p *ViewRowAccessPolicy) additionalValidations() error {
	if !valueSet(p.On) {
		return errNotSet("CreateViewOptions.RowAccessPolicy", "On")
	}
	return nil
}

func (p *ViewAddRowAccessPolicy) additionalValidations() error {
	if !valueSet(p.On) {
		return errNotSet("AlterViewOptions.AddRowAccessPolicy", "On")
	}
	return nil
}

func (opts *CreateViewOptions) additionalValidations() error {
	var errs []error
	if valueSet(opts.Columns) {
		for i, columnOption := range opts.Columns {
			if valueSet(columnOption.MaskingPolicy) {
				if !ValidObjectIdentifier(columnOption.MaskingPolicy.MaskingPolicy) {
					errs = append(errs, errInvalidIdentifier(fmt.Sprintf("CreateViewOptions.Columns[%d]", i), "MaskingPolicy"))
				}
			}
			if valueSet(columnOption.ProjectionPolicy) {
				if !ValidObjectIdentifier(columnOption.ProjectionPolicy.ProjectionPolicy) {
					errs = append(errs, errInvalidIdentifier(fmt.Sprintf("CreateViewOptions.Columns[%d]", i), "ProjectionPolicy"))
				}
			}
		}
	}
	return JoinErrors(errs...)
}

var AllViewDataMetricScheduleMinutes = []int{5, 15, 30, 60, 720, 1440}

// TODO(SNOW-1636212): remove
func (v *View) HasCopyGrants() bool {
	return strings.Contains(v.Text, " COPY GRANTS ")
}

func (v *View) IsTemporary() bool {
	return strings.Contains(v.Text, "TEMPORARY")
}

func (v *View) IsRecursive() bool {
	return strings.Contains(v.Text, "RECURSIVE")
}

func (v *View) IsChangeTracking() bool {
	return v.ChangeTracking == "ON"
}

func (r *CreateViewRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

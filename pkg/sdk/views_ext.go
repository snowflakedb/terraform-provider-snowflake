package sdk

import "strings"

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

package sdk

import "fmt"

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

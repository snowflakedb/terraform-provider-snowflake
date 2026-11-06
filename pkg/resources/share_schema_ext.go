package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// shareAttributesModificationExt keeps the case-insensitive comparison of the accounts, which the generated one does not do.
// The generated suppressor stays, so quoting differences are still suppressed.
func shareAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["accounts"].DiffSuppressFunc = SuppressIfAny(s["accounts"].DiffSuppressFunc, IgnoreCaseInSet("accounts"))
	return s
}

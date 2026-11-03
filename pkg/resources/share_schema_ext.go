package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var shareIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the identifier for the share; must be unique for the account in which the share is created.",
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

// shareAttributesModificationExt keeps the case-insensitive comparison of the accounts, which the generated one does not do.
// The generated suppressor stays, so quoting differences are still suppressed.
func shareAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["accounts"].DiffSuppressFunc = SuppressIfAny(s["accounts"].DiffSuppressFunc, IgnoreCaseInSet("accounts"))
	return s
}

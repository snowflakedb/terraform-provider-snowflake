package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var shareSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the identifier for the share; must be unique for the account in which the share is created.",
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the managed account.",
	},
	"accounts": {
		// Changed from Set to List to use DiffSuppressFunc: https://github.com/hashicorp/terraform-plugin-sdk/issues/160
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidAccountIdentifier(),
		},
		Optional: true,
		Description: "A list of accounts to be added to the share. Values should not be the account locator, but " +
			"in the form of 'organization_name.account_name",
		DiffSuppressFunc: ignoreCaseSuppressFunc,
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

package resources

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// Snowflake has no command that lists the external volumes of a failover group.
func failoverGroupAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["allowed_external_volumes"].Description = externalChangesNotDetectedFieldDescription(s["allowed_external_volumes"].Description)
	return s
}

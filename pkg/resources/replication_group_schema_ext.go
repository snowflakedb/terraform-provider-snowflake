package resources

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	replicationGroupIsPrimaryDescription          = "Indicates if the replication group is primary. When Terraform detects that the replication group is not primary, the resource is recreated."
	replicationGroupSecondaryIsPrimaryDescription = "Indicates if the replication group primary status has been changed. If change is detected, resource will be recreated."
)

// replicationGroupAttributesModificationExt adds the computed is_primary flag.
// The generator does not emit a computed attribute, and this field is not a CREATE or ALTER clause.
func replicationGroupAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["is_primary"] = replicationGroupIsPrimarySchema(replicationGroupIsPrimaryDescription)
	return s
}

func replicationGroupIsPrimarySchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: description,
	}
}

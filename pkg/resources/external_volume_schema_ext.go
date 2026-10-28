package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// A list is used as the order of storage locations matter. Storage location position in the list is used to select
// the active storage location - https://docs.snowflake.com/en/user-guide/tables-iceberg-storage#active-storage-location
// This is also why it has been left as one list with optional cloud dependent parameters, rather than splitting into
// one list per cloud provider. The block itself is generated. These adjustments are schema behavior the generator does not emit.
func externalVolumeAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	// allow_writes is an optional boolean string, but the existing schema does not validate it.
	s["allow_writes"].ValidateDiagFunc = nil

	fields := s["storage_location"].Elem.(*schema.Resource).Schema
	fields["storage_location_name"].Description = blocklistedCharactersFieldDescription("Name of the storage location. Must be unique for the external volume. Do not use the name `terraform_provider_sentinel_storage_location` - this is reserved for the provider for performing update operations.")
	fields["encryption_type"].DiffSuppressFunc = func(_, oldValue, newValue string, _ *schema.ResourceData) bool {
		return oldValue == "NONE" && newValue == ""
	}
	fields["use_privatelink_endpoint"].DiffSuppressFunc = func(_, oldValue, newValue string, _ *schema.ResourceData) bool {
		return oldValue == "" && (newValue == BooleanDefault || newValue == BooleanFalse)
	}
	return s
}

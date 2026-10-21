package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// ExternalVolumeAttrs holds the external volume fields that SHOW reuses.
// DESCRIBE is a property list parsed into a different shape, so those columns
// are not recorded here.
var ExternalVolumeAttrs = struct {
	Comment     ssot.Attribute
	AllowWrites ssot.Attribute
	Name        ssot.Attribute
}{
	Comment: ssot.Attribute{
		Name:          "comment",
		Type:          ssot.AttrString,
		Description:   "Specifies a comment for the external volume.",
		OptionalInDb:  true,
		CreateMapping: ssot.StringAttributeCreateBuilder,
	},
	AllowWrites: ssot.Attribute{
		Name:          "allow_writes",
		Type:          ssot.AttrBool,
		Description:   "Specifies whether write operations are allowed for the external volume; must be set to TRUE for Iceberg tables that use Snowflake as the catalog.",
		CreateMapping: ssot.BooleanStringAttributeCreateBuilder,
	},
	Name: ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
}

// ExternalVolumeObject is not named ExternalVolume because that identifier is the
// EXTERNAL_VOLUME parameter in the catalog.
var ExternalVolumeObject = ssot.Object{
	Name:                                "ExternalVolume",
	NamePlural:                          "ExternalVolumes",
	DescribeFailRead:                    true,
	IdType:                              "AccountObjectIdentifier",
	IdentityNameForceNew:                true,
	AttributesModification:              true,
	AdditionalDescribeOutputDescription: "Because of Terraform limitations, the changes on storage_location field do not mark this field as computed.",
	Attributes: []ssot.Attribute{
		ExternalVolumeAttrs.Comment,
		ExternalVolumeAttrs.AllowWrites,
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&ExternalVolumeAttrs.Name,
			&ExternalVolumeAttrs.AllowWrites,
			&ExternalVolumeAttrs.Comment,
		},
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithDescribeOutputSet(ssot.Omitted()),
		ssot.WithSetDescribeOutput(ssot.ExtCall()),
		ssot.WithStateUpgraders(ssot.Generated()),
	},
	SchemaStateUpgraders: []ssot.StateUpgraderFunc{
		"v2_14_0_ExternalVolumeStateUpgrader",
	},
}

package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

var Schema = ssot.Object{
	Name:                                "Schema",
	NamePlural:                          "Schemas",
	DescribeUsesShowId:                  true,
	HasPartialOnUpdate:                  true,
	AdditionalDescribeOutputDescription: "In order to handle this output, one must grant sufficient privileges, e.g. [grant_ownership](./grant_ownership) on all objects in the schema.",
	IdType:                              "DatabaseObjectIdentifier",
	AdditionalIdentityNameDescription:   "When the name is `PUBLIC`, during creation the provider checks if this schema has already been created and, in such case, `ALTER` is used to match the desired state.",
	AttributesModification:              true,
	Attributes: []ssot.Attribute{
		{
			Name:        "with_managed_access",
			Type:        ssot.AttrBool,
			Description: "Specifies a managed schema. Managed access schemas centralize privilege management with the schema owner.",
		},
		{
			Name:        "is_transient",
			Type:        ssot.AttrBool,
			ForceNew:    true,
			Description: "Specifies the schema as transient. Transient schemas do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
		},
		{
			Name:        "comment",
			Type:        ssot.AttrString,
			Description: "Specifies a comment for the schema.",
		},
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithBeforeCreate(ssot.ExtCall()),
		ssot.WithHierarchyRename(ssot.ExtCall()),
		ssot.WithBeforeAlter(ssot.ExtCall()),
		ssot.WithDescribeOutputSet(ssot.Omitted()),
		ssot.WithSetDescribeOutput(ssot.ExtCall()),
		ssot.WithSchemaVersion(ssot.ExtCall()),
		ssot.WithApplyParametersCreate(ssot.ExtCall()),
		ssot.WithShowParametersInSdk(ssot.ExtCall()),
		ssot.WithParametersDetailsFromRaw(ssot.ExtCall()),
		ssot.WithSetParametersFields(ssot.ExtCall()),
		ssot.WithParametersOutputSet(ssot.ExtCall()),
		ssot.WithApplyParametersChanges(ssot.ExtCall()),
		ssot.WithParametersAttributesSchema(ssot.Generated()),
		ssot.WithParametersOutputSchema(ssot.Generated()),
		ssot.WithParametersFieldNames(ssot.Generated()),
		ssot.WithParametersCustomDiff(ssot.Generated()),
	},
}

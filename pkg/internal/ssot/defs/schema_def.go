package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

var SchemaAttrs = struct {
	WithManagedAccess ssot.Attribute
	IsTransient       ssot.Attribute
	Comment           ssot.Attribute
	CreatedOn         ssot.Attribute
	DroppedOn         ssot.Attribute
	Name              ssot.Attribute
	IsDefault         ssot.Attribute
	IsCurrent         ssot.Attribute
	DatabaseName      ssot.Attribute
	Owner             ssot.Attribute
	Options           ssot.Attribute
	RetentionTime     ssot.Attribute
	OwnerRoleType     ssot.Attribute
	Kind              ssot.Attribute
}{
	WithManagedAccess: ssot.Attribute{
		Name:        "with_managed_access",
		Type:        ssot.AttrBool,
		Description: "Specifies a managed schema. Managed access schemas centralize privilege management with the schema owner.",
	},
	IsTransient: ssot.Attribute{
		Name:        "is_transient",
		Type:        ssot.AttrBool,
		ForceNew:    true,
		Description: "Specifies the schema as transient. Transient schemas do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
	},
	Comment: ssot.Attribute{
		Name:         "comment",
		Type:         ssot.AttrString,
		Description:  "Specifies a comment for the schema.",
		OptionalInDb: true,
	},
	CreatedOn:     ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	DroppedOn:     ssot.Attribute{Name: "dropped_on", Type: ssot.AttrTime, OptionalInDb: true, Readonly: true},
	Name:          ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	IsDefault:     ssot.Attribute{Name: "is_default", Type: ssot.AttrBool, Readonly: true},
	IsCurrent:     ssot.Attribute{Name: "is_current", Type: ssot.AttrBool, Readonly: true},
	DatabaseName:  ssot.Attribute{Name: "database_name", Type: ssot.AttrString, Readonly: true},
	Owner:         ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true},
	Options:       ssot.Attribute{Name: "options", Type: ssot.AttrString, OptionalInShow: true, OptionalInDb: true, Readonly: true},
	RetentionTime: ssot.Attribute{Name: "retention_time", Type: ssot.AttrString, Readonly: true},
	OwnerRoleType: ssot.Attribute{Name: "owner_role_type", Type: ssot.AttrString, Readonly: true},
	Kind:          ssot.Attribute{Name: "kind", Type: ssot.AttrString, Readonly: true},
}

var Schema = ssot.Object{
	Name:                                "Schema",
	NamePlural:                          "Schemas",
	DescribeUsesShowId:                  true,
	HasPartialOnUpdate:                  true,
	AdditionalDescribeOutputDescription: "In order to handle this output, one must grant sufficient privileges, e.g. [grant_ownership](./grant_ownership) on all objects in the schema.",
	IdType:                              "DatabaseObjectIdentifier",
	AdditionalIdentityNameDescription:   "When the name is `PUBLIC`, during creation the provider checks if this schema has already been created and, in such case, `ALTER` is used to match the desired state.",
	AttributesModification:              true,
	Generators: []ssot.GeneratorOption{
		ssot.WithSchemaShow(),
		ssot.WithSchemaDescribe(),
		ssot.WithModel(),
		ssot.WithResourceAssert(),
	},
	Attributes: []ssot.Attribute{
		SchemaAttrs.WithManagedAccess,
		SchemaAttrs.IsTransient,
		SchemaAttrs.Comment,
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&SchemaAttrs.CreatedOn,
			&SchemaAttrs.DroppedOn,
			&SchemaAttrs.Name,
			&SchemaAttrs.IsDefault,
			&SchemaAttrs.IsCurrent,
			&SchemaAttrs.DatabaseName,
			&SchemaAttrs.Owner,
			&SchemaAttrs.Comment,
			&SchemaAttrs.Options,
			&SchemaAttrs.RetentionTime,
			&SchemaAttrs.OwnerRoleType,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&SchemaAttrs.CreatedOn,
			&SchemaAttrs.Name,
			&SchemaAttrs.Kind,
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

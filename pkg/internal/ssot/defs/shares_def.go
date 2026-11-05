package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

var ShareKind = ssot.NewEnum("ShareKind", "ShareKinds", "INBOUND", "OUTBOUND")

// shareDescribeKindEnum names sdk.ObjectType for the DESCRIBE kind column. The enum is not generated (SdkType is set):
// the SDK ext converts the column into sdk.ObjectType.
var shareDescribeKindEnum = ssot.NewEnum("ObjectType", "ObjectTypes").WithSdkType("ObjectType")

// ShareAttrs holds the share fields that CREATE, ALTER, SHOW and DESCRIBE reuse.
//
// SHOW and DESCRIBE both return a name column, but of different types: SHOW a plain string, DESCRIBE the name of an
// object in the share, which is an outbound or an inbound (external) identifier. They are separate attributes.
// The to column of a share linked to a listing holds "ALL ACCOUNTS IN DATA EXCHANGE", which is not an account: the manual
// conversion in pkg/sdk/shares_ext.go skips it.
// to_total is not a Snowflake column: it is the N of the "... (N TOTAL)" marker that SHOW SHARES appends to the to column.
var ShareAttrs = struct {
	SecureObjectsOnly ssot.Attribute
	Comment           ssot.Attribute
	Accounts          ssot.Attribute
	CreatedOn         ssot.Attribute
	Kind              ssot.Attribute
	OwnerAccount      ssot.Attribute
	Name              ssot.Attribute
	DatabaseName      ssot.Attribute
	To                ssot.Attribute
	Owner             ssot.Attribute
	ListingGlobalName ssot.Attribute
	ToTotal           ssot.Attribute
	DescribeKind      ssot.Attribute
	DescribeName      ssot.Attribute
	SharedOn          ssot.Attribute
}{
	SecureObjectsOnly: ssot.Attribute{
		Name:     "secure_objects_only",
		Type:     ssot.AttrBool,
		ForceNew: true,
		// The same attribute is the SHOW column. SHOW SHARES returns "true" or "false" as text, and an empty value for an inbound share.
		StoredAsText:   true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Specifies whether only secure objects can be added to the share. Once set to false, it cannot be set back to true.",
	},
	Comment: ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the share.",
	},
	Accounts: ssot.Attribute{
		Name: "accounts",
		Type: ssot.AttrSet,
		Elem: &ssot.Attribute{
			Type:           ssot.AttrIdentifier,
			IdentifierType: "AccountIdentifier",
		},
		Description: "A set of accounts the share is available to, in the form `organization_name.account_name` (not the account locator). " +
			"The provider reads the recipients from Snowflake, so accounts added outside of Terraform are removed on the next apply.",
		// TODO [SNOW-4428377]: Add an optional share_restrictions field (SHARE_RESTRICTIONS of ALTER SHARE ... ADD ACCOUNTS, available in the SDK).
		// It is a statement-time modifier of ADD (default TRUE) that needs the OVERRIDE SHARE RESTRICTIONS privilege and cannot be read back from
		// SHOW, DESCRIBE or ACCOUNT_USAGE, so it would be applied on add only and its external changes would not be detected.
		// The ACCOUNTS list has no parentheses and no tag for them. The ADD and REMOVE blocks set the same clause.
		SetSql: ssot.PlainParameter(),
		// CREATE SHARE has no ACCOUNTS clause: the accounts are added after the share is created (AfterCreate).
		SkipCreate: true,
	},
	CreatedOn:    ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	Kind:         ssot.Attribute{Name: "kind", Type: ssot.AttrEnum, Enum: ShareKind, Readonly: true},
	OwnerAccount: ssot.Attribute{Name: "owner_account", Type: ssot.AttrString, Readonly: true},
	Name:         ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	DatabaseName: ssot.Attribute{Name: "database_name", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier", Readonly: true},
	To: ssot.Attribute{
		Name:          "to",
		Type:          ssot.AttrList,
		Elem:          &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountIdentifier"},
		ManualConvert: true,
		Readonly:      true,
	},
	Owner:             ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true},
	ListingGlobalName: ssot.Attribute{Name: "listing_global_name", Type: ssot.AttrString, Readonly: true},
	ToTotal:           ssot.Attribute{Name: "to_total", Type: ssot.AttrInt, Readonly: true},
	DescribeKind:      ssot.Attribute{Name: "kind", Type: ssot.AttrEnum, Enum: shareDescribeKindEnum, ManualConvert: true, Readonly: true},
	DescribeName:      ssot.Attribute{Name: "name", Type: ssot.AttrIdentifier, IdentifierType: "ObjectIdentifier", ManualConvert: true, Readonly: true},
	SharedOn:          ssot.Attribute{Name: "shared_on", Type: ssot.AttrTime, Readonly: true},
}

// shareAdd and shareRemove are the ADD ACCOUNTS and REMOVE ACCOUNTS blocks of ALTER. They are not Terraform fields.
var (
	shareAddAccounts = ssot.Attribute{
		Name:      "add",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(),
		Children: []ssot.Attribute{
			ssot.WithCreateSql(&ShareAttrs.Accounts, ssot.PlainParameter(ssot.Required())),
			{Name: "share_restrictions", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolParameter()},
		},
	}
	shareRemoveAccounts = ssot.Attribute{
		Name:      "remove",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(),
		Children: []ssot.Attribute{
			ssot.WithCreateSql(&ShareAttrs.Accounts, ssot.PlainParameter(ssot.Required())),
		},
	}
)

// Share generates the whole SDK interface, the resource (a preview feature), and their test helpers. The ext parts of the
// resource are about the accounts (see the comment on ResourceHooks) and the import.
//
// Checked against create-share, alter-share, drop-share, show-shares, desc-share, and parameters.
//
// SHOW SHARES returns a secure_objects_only column ("true" for an outbound share, "" for an inbound one, observed on a live
// account) that the show-shares page does not list. It is an optional bool, because the inbound value is empty, and it is the
// same attribute as the secure_objects_only CREATE attribute.
//
// CREATE OR ALTER SHARE is documented (COMMENT only). It is never added, because our resources do not use CREATE OR ALTER.
//
// The docs show DROP SHARE <name> without IF EXISTS. The SDK keeps IF EXISTS.
//
// Object parameters: none is documented for a share.
var Share = ssot.Object{
	Name:                                "Share",
	NamePlural:                          "Shares",
	IdType:                              "AccountObjectIdentifier",
	PreviewFeature:                      true,
	AdditionalDescribeOutputDescription: "DESCRIBE SHARE returns a row for every object granted to the share, including the system objects that Snowflake adds to the shared database. In order to see the rows, the role must have a privilege on the share.",
	IdentityNameForceNew:                true,
	AttributesModification:              true,
	SchemaStateUpgraders: []ssot.StateUpgraderFunc{
		"v2_23_0_ShareStateUpgrader",
	},
	// The accounts are the reason for the ext hooks: they need a temporary database when no database is granted yet,
	// the update makes the recipients equal to the configured accounts (BeforeAlter; the generated SET and UNSET handle the
	// comment), and the read merges the recipients that SHOW SHARES does not list (AfterShowById).
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithParseId(ssot.Generated()),
		ssot.WithNewCreateRequest(ssot.Generated()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithCreateInSdk(ssot.Generated()),
		ssot.WithAfterCreate(ssot.ExtCall()),
		ssot.WithShowByIdSafelyInSdk(ssot.Generated()),
		ssot.WithAfterShowById(ssot.ExtCall()),
		ssot.WithDescribeInSdk(ssot.ExtCall()),
		ssot.WithHandleExternalChanges(ssot.Generated()),
		ssot.WithSetStateToValuesFromConfig(ssot.Generated()),
		ssot.WithSetConfigFields(ssot.ExtCall()),
		ssot.WithApplySetUnset(ssot.Generated()),
		ssot.WithAlterSetInSdk(ssot.Generated()),
		ssot.WithAlterUnsetInSdk(ssot.Generated()),
		ssot.WithBeforeAlter(ssot.ExtCall()),
		ssot.WithSetFieldsNotSetByRead(ssot.ExtCall()),
		ssot.WithShowByIdInSdk(ssot.Generated()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
		ssot.WithIdentitySchema(ssot.Generated()),
		ssot.WithDescribeOutputSchema(ssot.Generated()),
		// Describe is best-effort and its output is a list of rows, so the ext sets it.
		ssot.WithDescribeOutputSet(ssot.Omitted()),
		ssot.WithSetDescribeOutput(ssot.ExtCall()),
		ssot.WithCustomizeDiff(ssot.Generated()),
		ssot.WithStateUpgraders(ssot.Generated()),
	},
	Generators: []ssot.GeneratorOption{
		ssot.WithAllGenerators(
			ssot.WithSchema(ssot.SchemaShow(), ssot.SchemaDescribe()),
			ssot.WithObjectAssert(
				ssot.ObjectAssertShow(),
				ssot.ObjectAssertDescribe(ssot.WithNoIdentifiableObject()),
			),
		),
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkShow(
				ssot.SdkShowOptionalLike(),
				ssot.SdkShowOptionalStartsWith(),
				ssot.SdkShowOptionalLimit(),
				ssot.SdkShowByIDFiltering(ssot.SdkShowByIDLike),
				ssot.SdkShowPlainOnly(&ShareAttrs.ToTotal, "int"),
			),
			ssot.SdkCreate(),
			ssot.SdkDescribe(
				ssot.SdkDescribeMapping(ssot.SdkDescriptionMappingSlice),
				ssot.SdkDescribeIdentifierOverride("ObjectIdentifier"),
			),
			ssot.SdkAlter(),
			ssot.SdkCustomInterfaceMethods(),
			ssot.SdkDrop(),
		),
	},
	// The bodies stay in pkg/sdk/shares_ext.go.
	CustomInterfaceMethods: []ssot.CustomInterfaceMethodSpec{
		{
			Name:    "DescribeProvider",
			Doc:     "DescribeProvider returns the share contents for a provider (outbound) share.",
			Returns: []string{"[]ShareDetails", "error"},
		},
		{
			Name:       "DescribeConsumer",
			Doc:        "DescribeConsumer returns the share contents for a consumer (inbound) share.",
			Parameters: []ssot.CustomInterfaceMethodParameter{{Name: "id", Kind: "ExternalObjectIdentifier"}},
			Returns:    []string{"[]ShareDetails", "error"},
		},
		{
			Name:    "ShowTargetAccounts",
			Doc:     "ShowTargetAccounts returns all target accounts of an outbound share from SNOWFLAKE.ACCOUNT_USAGE.SHARES, because SHOW SHARES lists at most three of them.",
			Returns: []string{"[]AccountIdentifier", "error"},
		},
	},
	Attributes: []ssot.Attribute{
		ShareAttrs.SecureObjectsOnly,
		ShareAttrs.Comment,
		ShareAttrs.Accounts,
	},
	Create: ssot.CreateSpec{
		Clauses: []*ssot.Attribute{
			&ShareAttrs.SecureObjectsOnly,
			&ShareAttrs.Comment,
		},
	},
	Set: []*ssot.Attribute{
		&ShareAttrs.Accounts,
		&ShareAttrs.Comment,
	},
	Unset: []*ssot.Attribute{
		&ShareAttrs.Comment,
	},
	Alter: ssot.AlterSpec{
		Clauses: []ssot.AlterClause{
			ssot.AlterKeywordBlock("Add", "ADD", &shareAddAccounts),
			ssot.AlterKeywordBlock("Remove", "REMOVE", &shareRemoveAccounts),
			ssot.AlterSet(),
			ssot.AlterUnset(),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
		},
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&ShareAttrs.CreatedOn,
			&ShareAttrs.Kind,
			&ShareAttrs.OwnerAccount,
			&ShareAttrs.Name,
			&ShareAttrs.DatabaseName,
			&ShareAttrs.To,
			&ShareAttrs.Owner,
			&ShareAttrs.Comment,
			&ShareAttrs.ListingGlobalName,
			&ShareAttrs.SecureObjectsOnly,
			&ShareAttrs.ToTotal,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&ShareAttrs.DescribeKind,
			&ShareAttrs.DescribeName,
			&ShareAttrs.SharedOn,
		},
	},
}

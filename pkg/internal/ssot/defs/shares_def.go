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
// to_total is not a Snowflake column: it is the N of the "... (N TOTAL)" marker that SHOW SHARES appends to the to column.
var ShareAttrs = struct {
	SecureObjectsOnly     ssot.Attribute
	Comment               ssot.Attribute
	Accounts              ssot.Attribute
	CreatedOn             ssot.Attribute
	Kind                  ssot.Attribute
	OwnerAccount          ssot.Attribute
	Name                  ssot.Attribute
	DatabaseName          ssot.Attribute
	To                    ssot.Attribute
	Owner                 ssot.Attribute
	ListingGlobalName     ssot.Attribute
	ShowSecureObjectsOnly ssot.Attribute
	ToTotal               ssot.Attribute
	DescribeKind          ssot.Attribute
	DescribeName          ssot.Attribute
	SharedOn              ssot.Attribute
}{
	SecureObjectsOnly: ssot.Attribute{
		Name:        "secure_objects_only",
		Type:        ssot.AttrBool,
		Description: "Specifies whether only secure objects can be added to the share. Once set to false, it cannot be set back to true.",
	},
	Comment: ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the share.",
	},
	Accounts: ssot.Attribute{
		Name: "accounts",
		Type: ssot.AttrList,
		Elem: &ssot.Attribute{
			Type:           ssot.AttrIdentifier,
			IdentifierType: "AccountIdentifier",
		},
		Description: "Specifies the accounts that can consume the share.",
		// The ACCOUNTS list has no parentheses and no tag for them. The ADD and REMOVE blocks set the same clause.
		SetSql: ssot.PlainParameter(),
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
	Owner:                 ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true},
	ListingGlobalName:     ssot.Attribute{Name: "listing_global_name", Type: ssot.AttrString, Readonly: true},
	ShowSecureObjectsOnly: ssot.Attribute{Name: "secure_objects_only_output", ShowName: "secure_objects_only", Type: ssot.AttrString, Readonly: true},
	ToTotal:               ssot.Attribute{Name: "to_total", Type: ssot.AttrInt, Readonly: true},
	DescribeKind:          ssot.Attribute{Name: "kind", Type: ssot.AttrEnum, Enum: shareDescribeKindEnum, ManualConvert: true, Readonly: true},
	DescribeName:          ssot.Attribute{Name: "name", Type: ssot.AttrIdentifier, IdentifierType: "ObjectIdentifier", ManualConvert: true, Readonly: true},
	SharedOn:              ssot.Attribute{Name: "shared_on", Type: ssot.AttrTime, Readonly: true},
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

// Share is defined for the SDK only: no Generators entry sets WithResource. The whole SDK interface is built from here.
//
// Checked against create-share, alter-share, drop-share, show-shares, desc-share, and parameters.
//
// SHOW SHARES returns a secure_objects_only column ("true" for an outbound share, "" for an inbound one, observed on a live
// account) that the show-shares page does not list. It is a string, because the inbound value is empty. Separate from the
// secure_objects_only CREATE attribute, which is a bool, so it has the ShowName.
//
// CREATE OR ALTER SHARE is documented (COMMENT only) and is not added, because the existing interface did not have it.
//
// The docs show DROP SHARE <name> without IF EXISTS. The SDK keeps IF EXISTS.
//
// Object parameters: none is documented for a share.
var Share = ssot.Object{
	Name:       "Share",
	NamePlural: "Shares",
	IdType:     "AccountObjectIdentifier",
	Generators: []ssot.GeneratorOption{
		ssot.WithSchema(ssot.SchemaShow(ssot.SchemaShowManualFields(&ShareAttrs.Name))),
		ssot.WithObjectAssert(
			ssot.ObjectAssertShow(),
			ssot.ObjectAssertDescribe(ssot.WithNoIdentifiableObject()),
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
			&ShareAttrs.ShowSecureObjectsOnly,
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

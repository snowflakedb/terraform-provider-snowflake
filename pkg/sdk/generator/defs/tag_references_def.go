package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var TagReferenceObjectDomainDef = g.NewEnum(
	"TagReferenceObjectDomain",
	"TagReferenceObjectDomains",
	"ACCOUNT",
	"ALERT",
	"COLUMN",
	"COMPUTE POOL",
	"DATABASE",
	"DATABASE ROLE",
	"FAILOVER GROUP",
	"FUNCTION",
	"INTEGRATION",
	"NETWORK POLICY",
	"PROCEDURE",
	"REPLICATION GROUP",
	"ROLE",
	"SCHEMA",
	"SHARE",
	"STAGE",
	"STREAM",
	"TABLE",
	"TASK",
	"USER",
	"WAREHOUSE",
)

var TagReferenceApplyMethodDef = g.NewEnum(
	"TagReferenceApplyMethod",
	"TagReferenceApplyMethods",
	"CLASSIFIED",
	"INHERITED",
	"MANUAL",
	"PROPAGATED",
)

var tagReferenceParametersDef = g.NewQueryStruct("tagReferenceParameters").
	SQLWithCustomFieldName("functionFullyQualifiedName", "SNOWFLAKE.INFORMATION_SCHEMA.TAG_REFERENCES").
	OptionalQueryStructField(
		"arguments",
		tagReferenceFunctionArgumentsDef,
		g.ListOptions().Parentheses().Required(),
	).WithValidation(g.ValidateValueSet, "arguments")

var tagReferenceFunctionArgumentsDef = g.NewQueryStruct("tagReferenceFunctionArguments").
	OptionalText(
		"ObjectName",
		g.KeywordOptions().SingleQuotes().Required(),
	).
	OptionalEnum(
		"ObjectDomain",
		"TagReferenceObjectDomain",
		g.KeywordOptions().SingleQuotes().Required(),
	).
	WithValidation(g.ValidateValueSet, "ObjectName").
	WithValidation(g.ValidateValueSet, "ObjectDomain")

var tagReferencesDef = g.NewInterface(
	"TagReferences",
	"TagReference",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CustomShowOperation(
	"GetForEntity",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/functions/tag_references",
	g.DbStruct("tagReferenceDBRow").
		Text("TAG_DATABASE").
		Text("TAG_SCHEMA").
		Text("TAG_NAME").
		Text("TAG_VALUE").
		Text("LEVEL").
		OptionalText("OBJECT_DATABASE").
		OptionalText("OBJECT_SCHEMA").
		Text("OBJECT_NAME").
		Text("DOMAIN").
		OptionalText("COLUMN_NAME").
		Text("APPLY_METHOD"),
	g.PlainStruct("TagReference").
		Text("TagDatabase").
		Text("TagSchema").
		Text("TagName").
		Text("TagValue").
		Field("Level", "TagReferenceObjectDomain").
		Field("ObjectDatabase", "*string").
		Field("ObjectSchema", "*string").
		Text("ObjectName").
		Field("Domain", "TagReferenceObjectDomain").
		Field("ColumnName", "*string").
		Field("ApplyMethod", "TagReferenceApplyMethod"),
	g.NewQueryStruct("GetForEntity").
		SQLWithCustomFieldName("selectEverythingFrom", "SELECT * FROM TABLE").
		OptionalQueryStructField(
			"parameters",
			tagReferenceParametersDef,
			g.ListOptions().Parentheses().NoComma().Required(),
		).WithValidation(g.ValidateValueSet, "parameters"),
	tagReferenceParametersDef,
	tagReferenceFunctionArgumentsDef,
).WithEnums(TagReferenceObjectDomainDef, TagReferenceApplyMethodDef)

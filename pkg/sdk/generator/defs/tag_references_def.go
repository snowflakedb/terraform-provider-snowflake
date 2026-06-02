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

var tagReferenceParametersDef = g.NewQueryStruct("TagReferenceParameters").
	SQLWithCustomFieldName("functionFullyQualifiedName", "SNOWFLAKE.INFORMATION_SCHEMA.TAG_REFERENCES").
	OptionalQueryStructField(
		"arguments",
		tagReferenceFunctionArgumentsDef,
		g.ListOptions().Parentheses().Required(),
	).WithValidation(g.ValidateValueSet, "arguments")

var tagReferenceFunctionArgumentsDef = g.NewQueryStruct("TagReferenceFunctionArguments").
	Text(
		"ObjectName",
		g.KeywordOptions().SingleQuotes().Required(),
	).
	PredefinedQueryStructField(
		"ObjectDomain",
		TagReferenceObjectDomainDef.Kind(),
		g.KeywordOptions().SingleQuotes().Required(),
	).
	WithValidation(g.ValidateValueSet, "ObjectName").
	WithValidation(g.ValidateValueSet, "ObjectDomain")

var tagReferencesDef = g.NewInterface(
	"TagReferences",
	"TagReference",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CustomShowOperationWithPairedStructs(
	"GetForEntity",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/functions/tag_references",
	g.StructPair("tagReferenceDBRow", "TagReference").
		Text("TAG_DATABASE").
		Text("TAG_SCHEMA").
		Text("TAG_NAME").
		Text("TAG_VALUE").
		Enum("LEVEL", TagReferenceObjectDomainDef).
		OptionalText("OBJECT_DATABASE").
		OptionalText("OBJECT_SCHEMA").
		Text("OBJECT_NAME").
		Enum("DOMAIN", TagReferenceObjectDomainDef).
		OptionalText("COLUMN_NAME").
		Enum("APPLY_METHOD", TagReferenceApplyMethodDef),
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

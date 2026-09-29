package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var PolicyEntityDomainEnumDef = g.NewEnum(
	"PolicyEntityDomain", "PolicyEntityDomains",
	"ACCOUNT",
	"DYNAMIC_TABLE",
	"ICEBERG_TABLE",
	"INTEGRATION",
	"TABLE",
	"TAG",
	"USER",
	"VIEW",
)

var PolicyKindEnumDef = g.NewEnum(
	"PolicyKind", "PolicyKinds",
	"AGGREGATION_POLICY",
	"AUTHENTICATION_POLICY",
	"FEATURE_POLICY",
	"MASKING_POLICY",
	"PACKAGES_POLICY",
	"PASSWORD_POLICY",
	"PROJECTION_POLICY",
	"ROW_ACCESS_POLICY",
	"SESSION_POLICY",
	"STORAGE_LIFECYCLE_POLICY",
)

var policyReferenceParametersDef = g.NewQueryStruct("policyReferenceParameters").
	SQLWithCustomFieldName("functionFullyQualifiedName", "SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES").
	OptionalQueryStructField(
		"arguments",
		policyReferenceFunctionArgumentsDef,
		g.ListOptions().Parentheses().Required(),
	).WithValidation(g.ValidateValueSet, "arguments")

var policyReferenceFunctionArgumentsDef = g.NewQueryStruct("policyReferenceFunctionArguments").
	PredefinedQueryStructField("refEntityName", "[]ObjectIdentifier", g.ParameterOptions().ArrowEquals().SingleQuotes().SQL("REF_ENTITY_NAME").Required()).
	EnumAssignment("REF_ENTITY_DOMAIN", PolicyEntityDomainEnumDef, g.ParameterOptions().ArrowEquals().SingleQuotes().Required()).
	WithValidation(g.ValidateValueSet, "RefEntityDomain").
	WithValidation(g.ValidateValueSet, "refEntityName")

var policyReferencePairs = g.StructPair("policyReferenceDBRow", "PolicyReference").
	OptionalText("POLICY_DB").
	OptionalText("POLICY_SCHEMA").
	Text("POLICY_NAME").
	Enum("POLICY_KIND", PolicyKindEnumDef).
	OptionalText("REF_DATABASE_NAME").
	OptionalText("REF_SCHEMA_NAME").
	Text("REF_ENTITY_NAME").
	Text("REF_ENTITY_DOMAIN").
	OptionalText("REF_COLUMN_NAME").
	OptionalText("REF_ARG_COLUMN_NAMES").
	OptionalText("TAG_DATABASE").
	OptionalText("TAG_SCHEMA").
	OptionalText("TAG_NAME").
	OptionalText("POLICY_STATUS")

var policyReferencesDef = g.NewInterface(
	"PolicyReferences",
	"PolicyReference",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CustomShowOperationWithPairedStructs(
	"GetForEntity",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/functions/policy_references",
	policyReferencePairs,
	g.NewQueryStruct("GetForEntity").
		SQLWithCustomFieldName("selectEverythingFrom", "SELECT * FROM TABLE").
		OptionalQueryStructField(
			"parameters",
			policyReferenceParametersDef,
			g.ListOptions().Parentheses().NoComma().Required(),
		).WithValidation(g.ValidateValueSet, "parameters"),
	policyReferenceParametersDef,
	policyReferenceFunctionArgumentsDef,
).WithEnums(PolicyEntityDomainEnumDef, PolicyKindEnumDef)

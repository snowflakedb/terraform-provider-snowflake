package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

func pairedStructExampleAllOptions(dbName, plainName string) *g.PairedStructs {
	return g.StructPair(dbName, plainName).
		// Non-nullable string in both; auto-derived field names
		Text("object_name").
		// Both non-nullable strings; plain name explicitly overridden
		Text("both_non_nullable_strings", g.WithPlainFieldName("OverriddenNonNullableStringPlainName")).
		// Both db and plain names overridden
		Text("type", g.WithDbFieldName("StorageTypeDb"), g.WithPlainFieldName("StorageTypePlain")).
		// db sql.NullString, plain *string — both nullable; auto-derived field names
		OptionalText("both_nullable_strings").
		// db sql.NullString, plain string — WithRequiredInPlain strips the pointer
		OptionalText("organization_name", g.WithRequiredInPlain()).
		// Both non-nullable bool; auto-derived field names
		Bool("is_primary").
		// db sql.NullBool, plain *bool — both nullable
		OptionalBool("is_default").
		// db sql.NullBool, plain bool — WithRequiredInPlain strips the pointer
		OptionalBool("enabled", g.WithRequiredInPlain()).
		// Both non-nullable int; auto-derived field names
		Number("next_value").
		// db sql.NullInt64, plain *int — both nullable; auto-derived field names
		OptionalNumber("port").
		// db sql.NullInt64, plain int — WithRequiredInPlain strips the pointer
		OptionalNumber("retry_limit", g.WithRequiredInPlain()).
		// Both time.Time; auto-derived field names
		Time("created_on").
		// db sql.NullTime, plain *time.Time — both nullable; auto-derived field names
		OptionalTime("updated_at").
		// db string, plain ExternalObjectIdentifier — fully explicit custom types
		Field("primary", "string", "ExternalObjectIdentifier").
		// db string (raw CSV), plain []AccountIdentifier — custom plain type
		PlainField("failover_allowed_to_accounts", "[]AccountIdentifier").
		// db string, plain []string; auto-derived field names
		StringList("tags").
		// db string, plain AccountObjectIdentifier; plain name defaults to "Id"
		AccountObjectIdentifier("account_id").
		// db string, plain AccountObjectIdentifier; plain name overridden
		AccountObjectIdentifier("second_account_id", g.WithPlainFieldName("OverriddenSecondAccountId"))
}

// PairedStructExample demonstrates the PairedStructs single-definition approach.
// It is the functional equivalent of calling ShowOperation/DescribeOperation with separate
// DbStruct and PlainStruct builders, but defines both structs in a single field-by-field chain.
// Every supported method and should be added here.
var PairedStructExample = g.NewInterface(
	"PairedStructExamples",
	"PairedStructExample",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	WithAllowedGenerationParts("default", "dto", "dto_builders", "impl", "validations").
	ShowOperationWithPairedStructs(
		"https://example.com",
		pairedStructExampleAllOptions("pairedStructExampleRow", "PairedStructExample"),
		g.NewQueryStruct("ShowPairedStructExamples").
			Show().
			SQL("WHATEVER").
			OptionalLike(),
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSingleValue,
		"https://example.com",
		pairedStructExampleAllOptions("pairedStructExampleDetailRow", "PairedStructExampleDetail"),
		g.NewQueryStruct("DescribePairedStructExamples").
			Describe().
			SQL("WHATEVER").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)

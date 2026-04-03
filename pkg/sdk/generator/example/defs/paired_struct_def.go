package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// PairedStructExample demonstrates the PairedStructs single-definition approach.
// It is the functional equivalent of calling ShowOperation/DescribeOperation with separate
// DbStruct and PlainStruct builders, but defines both structs in a single field-by-field chain.
//
// The example covers every supported variation:
//
//	Text          - both non-nullable strings; plain name auto-derived from db column name
//	Text+option   - both non-nullable strings; plain name explicitly overridden
//	OptionalText  - db sql.NullString, plain *string (both nullable)
//	OptionalText  - db sql.NullString, plain string (WithRequiredInPlain strips the pointer)
//	Time          - both time.Time
//	Bool          - both non-nullable bool
//	Field         - fully explicit db/plain types (e.g. db string → plain ExternalObjectIdentifier)
//	PlainField    - db string, custom plain type (e.g. []AccountIdentifier)
var PairedStructExample = g.NewInterface(
	"PairedStructExamples",
	"PairedStructExample",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	WithAllowedGenerationParts("default", "dto", "dto_builders", "impl", "validations").
	ShowOperationWithPairedStructs(
		"https://example.com",
		g.StructPair("pairedStructExampleRow", "PairedStructExample").
			// Both non-nullable strings; plain field name is auto-derived: "snowflake_region" → "SnowflakeRegion"
			Text("snowflake_region").
			// Both non-nullable strings; plain name overridden to "AccountName" instead of auto "AccountName"
			// (here the override matches the auto-derived result, but it shows the mechanism)
			Text("account_name", g.WithPlainFieldName("AccountName")).
			// Non-nullable string in both; auto-derived plain name "Name"
			Text("name").
			// db sql.NullString, plain *string — both nullable
			OptionalText("region_group").
			// db sql.NullString, plain string — WithRequiredInPlain strips the pointer
			OptionalText("organization_name", g.WithRequiredInPlain()).
			// Both time.Time
			Time("created_on").
			// Both non-nullable bool
			Bool("is_primary").
			// db string, plain ExternalObjectIdentifier — fully explicit custom types
			Field("primary", "string", "ExternalObjectIdentifier").
			// db string (raw CSV), plain []AccountIdentifier — custom plain type
			PlainField("failover_allowed_to_accounts", "[]AccountIdentifier").
			// Simple non-nullable strings
			Text("connection_url").
			Text("account_locator"),
		g.NewQueryStruct("ShowPairedStructExamples").
			Show().
			SQL("WHATEVER").
			OptionalLike(),
	).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://example.com",
	g.StructPair("pairedStructExampleDetailRow", "PairedStructExampleDetail").
		Text("name").
		Text("snowflake_region").
		// db sql.NullBool, plain bool — WithRequiredInPlain strips the pointer from *bool
		OptionalBool("is_primary", g.WithRequiredInPlain()).
		OptionalText("comment"),
	g.NewQueryStruct("DescribePairedStructExamples").
		Describe().
		SQL("WHATEVER").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)

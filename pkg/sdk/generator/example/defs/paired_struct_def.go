package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// PairedStructExample demonstrates the PairedStructs single-definition approach.
// It is the functional equivalent of calling ShowOperation/DescribeOperation with separate
// DbStruct and PlainStruct builders, but defines both structs in a single field-by-field chain.
//
// Every supported method and option is exercised here:
//
//	Text                   - both non-nullable strings; plain name auto-derived from db column name
//	Text + WithPlainFieldName - plain name explicitly overridden
//	Text + WithDbFieldName + WithPlainFieldName - both db and plain Go names overridden
//	OptionalText           - db sql.NullString, plain *string (both nullable)
//	OptionalText + WithRequiredInPlain - db nullable, plain non-nullable (pointer stripped)
//	Bool                   - both non-nullable bool
//	OptionalBool           - db sql.NullBool, plain *bool (both nullable)
//	OptionalBool + WithRequiredInPlain - db nullable, plain non-nullable
//	Number                 - both non-nullable int
//	OptionalNumber         - db sql.NullInt64, plain *int (both nullable)
//	OptionalNumber + WithRequiredInPlain - db nullable, plain non-nullable
//	Time                   - both time.Time
//	OptionalTime           - db sql.NullTime, plain *time.Time (both nullable)
//	Field                  - fully explicit db/plain types (e.g. db string → plain ExternalObjectIdentifier)
//	PlainField             - db string, custom plain type (e.g. []AccountIdentifier)
//	StringList             - db string, plain []string
//	AccountObjectIdentifier - db string, plain AccountObjectIdentifier, plain name defaults to "Id"
var PairedStructExample = g.NewInterface(
	"PairedStructExamples",
	"PairedStructExample",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	WithAllowedGenerationParts("default", "dto", "dto_builders", "impl", "validations").
	ShowOperationWithPairedStructs(
		"https://example.com",
		g.StructPair("pairedStructExampleRow", "PairedStructExample").
			// Both non-nullable strings; plain field name auto-derived: "snowflake_region" → "SnowflakeRegion"
			Text("snowflake_region").
			// Both non-nullable strings; plain name explicitly overridden
			Text("account_name", g.WithPlainFieldName("AccountName")).
			// Both db and plain Go names overridden: column "type" → Go field "StorageType" in both structs
			Text("type", g.WithDbFieldName("StorageType"), g.WithPlainFieldName("StorageType")).
			// Non-nullable string in both; auto-derived plain name "Name"
			Text("name").
			// db sql.NullString, plain *string — both nullable
			OptionalText("region_group").
			// db sql.NullString, plain string — WithRequiredInPlain strips the pointer
			OptionalText("organization_name", g.WithRequiredInPlain()).
			// Both non-nullable bool
			Bool("is_primary").
			// db sql.NullBool, plain *bool — both nullable
			OptionalBool("is_default").
			// db sql.NullBool, plain bool — WithRequiredInPlain strips the pointer
			OptionalBool("enabled", g.WithRequiredInPlain()).
			// Both non-nullable int
			Number("next_value").
			// db sql.NullInt64, plain *int — both nullable
			OptionalNumber("port").
			// db sql.NullInt64, plain int — WithRequiredInPlain strips the pointer
			OptionalNumber("retry_limit", g.WithRequiredInPlain()).
			// Both time.Time
			Time("created_on").
			// db sql.NullTime, plain *time.Time — both nullable
			OptionalTime("updated_at").
			// db string, plain ExternalObjectIdentifier — fully explicit custom types
			Field("primary", "string", "ExternalObjectIdentifier").
			// db string (raw CSV), plain []AccountIdentifier — custom plain type
			PlainField("failover_allowed_to_accounts", "[]AccountIdentifier").
			// db string, plain []string
			StringList("tags").
			// db string, plain AccountObjectIdentifier; plain name defaults to "Id"
			AccountObjectIdentifier("account_id"),
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

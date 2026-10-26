package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// ShowProjectionExampleMini is the subset struct. ShowProjection emits the method only, so the struct is a describe helper.
var showProjectionExampleMini = g.PlainStruct("ShowProjectionExampleMini").
	Text("Name").
	Number("Running")

// ShowProjectionExample demonstrates ShowProjection and ShowByIDProjection.
// AsMini copies Name and Running. ShowByIDMini and ShowByIDSafelyMini return that subset.
var ShowProjectionExampleDef = g.NewInterface(
	"ShowProjectionExamples",
	"ShowProjectionExample",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	WithAllowedGenerationParts(g.PartDefault, g.PartDto, g.PartDtoBuilders, g.PartImpl, g.PartValidations).
	ShowOperationWithPairedStructs(
		"https://example.com",
		g.StructPair("showProjectionExampleRow", "ShowProjectionExample").
			Text("name").
			Number("running").
			Text("comment"),
		g.NewQueryStruct("ShowShowProjectionExamples").
			Show().
			SQL("WHATEVER").
			OptionalLike(),
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSingleValue,
		"https://example.com",
		g.StructPair("showProjectionExampleDescRow", "ShowProjectionExampleDetails").
			Text("name"),
		g.NewQueryStruct("DescribeShowProjectionExample").
			Describe().
			SQL("WHATEVER").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
		showProjectionExampleMini,
	).
	ShowProjection(
		g.PlainStruct("ShowProjectionExample").
			Text("Name").
			Number("Running").
			Text("Comment"),
		"AsMini",
		showProjectionExampleMini,
	).
	ShowByIDProjection("Mini")

package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// DropSafelyHookExampleDef demonstrates DropOperation with WithDropSafelyHook.
// The generated DropSafely implementation calls v.dropSafelyHook(ctx, id) before dropping.
var DropSafelyHookExampleDef = g.NewInterface(
	"DropSafelyHookExamples",
	"DropSafelyHookExample",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).WithAllowedGenerationParts(g.PartDefault, g.PartImpl, g.PartDto, g.PartDtoBuilders, g.PartValidations).
	DropOperation(
		"https://example.com",
		g.NewQueryStruct("DropDropSafelyHookExample").
			Drop().
			SQL("WHATEVER").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
		g.WithDropSafelyHook(),
	)

// DropSafelyForceExampleDef demonstrates DropOperation with WithDropSafelyForce.
// The generated DropSafely appends .WithForce(true) to the Drop request.
var DropSafelyForceExampleDef = g.NewInterface(
	"DropSafelyForceExamples",
	"DropSafelyForceExample",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).WithAllowedGenerationParts(g.PartDefault, g.PartImpl, g.PartDto, g.PartDtoBuilders, g.PartValidations).
	DropOperation(
		"https://example.com",
		g.NewQueryStruct("DropDropSafelyForceExample").
			Drop().
			SQL("WHATEVER").
			IfExists().
			OptionalSQL("FORCE").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
		g.WithDropSafelyForce(),
	)

// DropSafelySuppressedExampleDef demonstrates DropOperation with WithDropSafelySuppressed.
// No DropSafely is generated; the one declared with WithCustomInterfaceMethod is implemented in the _ext.go file.
var DropSafelySuppressedExampleDef = g.NewInterface(
	"DropSafelySuppressedExamples",
	"DropSafelySuppressedExample",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).WithAllowedGenerationParts(g.PartDefault, g.PartImpl, g.PartDto, g.PartDtoBuilders, g.PartValidations).
	DropOperation(
		"https://example.com",
		g.NewQueryStruct("DropDropSafelySuppressedExample").
			Drop().
			SQL("WHATEVER").
			Name().
			OptionalSQL("CASCADE").
			WithValidation(g.ValidIdentifier, "name"),
		g.WithDropSafelySuppressed(),
	).
	WithCustomInterfaceMethod(
		"DropSafely",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("request", "*DropDropSafelySuppressedExampleRequest")},
		"error",
	)

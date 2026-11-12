package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var suspendIfStartedExampleStateDef = g.NewEnum(
	"SuspendIfStartedExampleState", "SuspendIfStartedExampleStates",
	"STARTED",
	"SUSPENDED",
)

// SuspendIfStartedExampleDef demonstrates SuspendIfStarted. The generated method reads the state with ShowByID, sends
// SUSPEND when it is STARTED, and waits for SUSPENDED. The returned function sends RESUME IF SUSPENDED.
var SuspendIfStartedExampleDef = g.NewInterface(
	"SuspendIfStartedExamples",
	"SuspendIfStartedExample",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	WithAllowedGenerationParts(g.PartDefault, g.PartDto, g.PartDtoBuilders, g.PartImpl, g.PartValidations, g.PartEnums).
	AlterOperation(
		"https://example.com",
		g.NewQueryStruct("AlterSuspendIfStartedExample").
			Alter().
			SQL("WHATEVER").
			Name().
			OptionalSQL("SUSPEND").
			OptionalSQL("RESUME").
			OptionalSQL("IF SUSPENDED").
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Suspend", "Resume", "Comment"),
	).
	ShowOperationWithPairedStructs(
		"https://example.com",
		g.StructPair("suspendIfStartedExampleRow", "SuspendIfStartedExample").
			Text("name").
			Enum("state", suspendIfStartedExampleStateDef),
		g.NewQueryStruct("ShowSuspendIfStartedExamples").
			Show().
			SQL("WHATEVERS").
			OptionalLike(),
		g.ShowByIDLikeFiltering,
	).
	SuspendIfStarted(g.SuspendIfStartedConfig{
		StateField:        "State",
		StartedConstant:   "SuspendIfStartedExampleStateStarted",
		SuspendedConstant: "SuspendIfStartedExampleStateSuspended",
		ResumeIfSuspended: true,
	}).
	WithEnums(suspendIfStartedExampleStateDef)

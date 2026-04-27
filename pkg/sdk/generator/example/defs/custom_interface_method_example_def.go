package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var CustomInterfaceMethodExamplesDef = g.NewInterface(
	"CustomInterfaceMethodExamples",
	"CustomInterfaceMethodExample",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	WithAllowedGenerationParts("default", "dto", "impl", "validations").
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-database-role",
		g.NewQueryStruct("CreateCustomInterfaceMethodExample").
			Create().
			SQL("CUSTOM EXAMPLE").
			Name(),
	).
	// WithCustomInterfaceMethod with no extra parameters and a single return type.
	WithCustomInterfaceMethod(
		"UnsetAll",
		"UnsetAll allows unsetting all parameters simultaneously",
		nil,
		"error",
	).
	// WithCustomInterfaceMethod with one parameter and two return types.
	WithCustomInterfaceMethod(
		"ShowParameters",
		"ShowParameters redirects the invocation to the common Parameters command",
		[]*g.MethodParameter{
			g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]()),
		},
		"[]*Parameter", "error",
	).
	// WithCustomInterfaceMethod with multiple parameters and two return types.
	WithCustomInterfaceMethod(
		"SuspendRootTasks",
		"SuspendRootTasks is added to iterate over tasks to find all roots and suspend them, as no such method is available on the Snowflake side",
		[]*g.MethodParameter{
			g.NewMethodParameter("taskId", g.KindOfT[sdkcommons.SchemaObjectIdentifier]()),
			g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]()),
		},
		"[]SchemaObjectIdentifier", "error",
	).
	// WithCustomInterfaceMethod with no return values and no comment.
	WithCustomInterfaceMethod(
		"Refresh",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]()),
		},
	).
	// WithCustomInterfaceMethod with a multiline doc comment.
	WithCustomInterfaceMethod(
		"ResumeTasks",
		"ResumeTasks is added manually;\nit resumes all given tasks in the correct order",
		[]*g.MethodParameter{
			g.NewMethodParameter("ids", "[]SchemaObjectIdentifier"),
		},
		"error",
	)

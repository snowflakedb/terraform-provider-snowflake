package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

const parametersDoc = "https://docs.snowflake.com/en/sql-reference/sql/show-parameters"

var parametersDef = g.NewInterface(
	"Parameters",
	"Parameter",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	WithAllowedGenerationParts(g.PartDefault).
	WithCustomInterfaceMethod(
		"SetAccountParameter",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", "AccountParameter"),
			g.NewMethodParameter("value", g.KindString),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"UnsetAccountParameter",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("parameter", "AccountParameter")},
		"error",
	).
	WithCustomInterfaceMethod(
		"SetSessionParameterOnAccount",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", "SessionParameter"),
			g.NewMethodParameter("value", g.KindString),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"SetSessionParameterOnUser",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("userID", g.KindOfT[sdkcommons.AccountObjectIdentifier]()),
			g.NewMethodParameter("parameter", "SessionParameter"),
			g.NewMethodParameter("value", g.KindString),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"SetObjectParameterOnAccount",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", "ObjectParameter"),
			g.NewMethodParameter("value", g.KindString),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"SetObjectParameterOnObject",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("object", "Object"),
			g.NewMethodParameter("parameter", "ObjectParameter"),
			g.NewMethodParameter("value", g.KindString),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"UnsetObjectParameterOnAccount",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("parameter", "ObjectParameter")},
		"error",
	).
	WithCustomInterfaceMethod(
		"UnsetObjectParameterOnObject",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("object", "Object"),
			g.NewMethodParameter("parameter", "ObjectParameter"),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"ShowParameters",
		parametersDoc,
		[]*g.MethodParameter{g.NewMethodParameter("opts", "*ShowParametersOptions")},
		"[]*Parameter", "error",
	).
	WithCustomInterfaceMethod(
		"ShowAccountParameter",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("parameter", "AccountParameter")},
		"*Parameter", "error",
	).
	WithCustomInterfaceMethod(
		"ShowSessionParameter",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("parameter", "SessionParameter")},
		"*Parameter", "error",
	).
	WithCustomInterfaceMethod(
		"ShowUserParameter",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", "UserParameter"),
			g.NewMethodParameter("user", g.KindOfT[sdkcommons.AccountObjectIdentifier]()),
		},
		"*Parameter", "error",
	).
	WithCustomInterfaceMethod(
		"ShowObjectParameter",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", "ObjectParameter"),
			g.NewMethodParameter("object", "Object"),
		},
		"*Parameter", "error",
	)

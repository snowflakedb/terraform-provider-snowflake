package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

func All() []ssot.Object {
	return []ssot.Object{
		ComputePool,
		Schema,
	}
}

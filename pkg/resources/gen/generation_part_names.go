package gen

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"

type generationPartName string

func (g generationPartName) GenerationPartName() string { return string(g) }

const (
	PartDefault generationPartName = "default"
	PartSchema  generationPartName = "schema"
)

var (
	_ genhelpers.GenerationPartNamer = PartDefault
	_ genhelpers.GenerationPartNamer = PartSchema
)

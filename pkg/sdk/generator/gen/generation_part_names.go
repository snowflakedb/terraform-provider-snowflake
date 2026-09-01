package gen

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"

type generationPartName string

// GenerationPartName restricts which values can be passed to SDK generator part APIs.
// Only constants defined in this package satisfy this interface.
type GenerationPartName interface {
	genhelpers.GenerationPartNamer
	xxxProtected()
}

func (g generationPartName) GenerationPartName() string { return string(g) }
func (g generationPartName) xxxProtected()              {}

const (
	PartDefault     generationPartName = "default"
	PartDto         generationPartName = "dto"
	PartDtoBuilders generationPartName = "dto_builders"
	PartImpl        generationPartName = "impl"
	PartUnitTests   generationPartName = "unit_tests"
	PartValidations generationPartName = "validations"
	PartEnums       generationPartName = "enums"
)

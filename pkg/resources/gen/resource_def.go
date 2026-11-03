package gen

import (
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

// ResourceObject is an ssot object plus this generator's per-object part settings.
type ResourceObject struct {
	ssot.Object
	*genhelpers.ObjectGenerationSettings
}

var resourceGenerationSettings = map[string]*genhelpers.ObjectGenerationSettings{
	"Pipe": {
		AllowedGenerationParts: []genhelpers.GenerationPartNamer{PartDefault, PartSchema, PartParameters},
		EnabledGenerationParts: []genhelpers.GenerationPartNamer{PartParameters},
	},
	"Schema": {
		AllowedGenerationParts: []genhelpers.GenerationPartNamer{PartDefault, PartSchema, PartParameters},
		EnabledGenerationParts: []genhelpers.GenerationPartNamer{PartParameters},
	},
	"Warehouse": {
		AllowedGenerationParts: []genhelpers.GenerationPartNamer{PartDefault, PartSchema, PartParameters},
		EnabledGenerationParts: []genhelpers.GenerationPartNamer{PartParameters},
	},
	"WarehouseAdaptive": {
		AllowedGenerationParts: []genhelpers.GenerationPartNamer{PartDefault, PartSchema, PartParameters},
		EnabledGenerationParts: []genhelpers.GenerationPartNamer{PartParameters},
	},
	"WarehouseInteractive": {
		AllowedGenerationParts: []genhelpers.GenerationPartNamer{PartDefault, PartSchema, PartParameters},
		EnabledGenerationParts: []genhelpers.GenerationPartNamer{PartParameters},
	},
}

func GetAllObjects() []ResourceObject {
	if len(defs.All()) == 0 {
		log.Fatal("no object definitions registered")
	}
	fmt.Println("resource object definitions:")
	objects := defs.SelectObjects(func(g ssot.Generators) bool { return g.Resource.Enabled }, selectsForResource)
	out := make([]ResourceObject, 0, len(objects))
	for _, def := range objects {
		fmt.Printf(" - %s\n", def.Name)
		out = append(out, ResourceObject{
			Object:                   def,
			ObjectGenerationSettings: resourceGenerationSettings[def.Name],
		})
	}
	return out
}

// selectsForResource is the one filter of the resource generator: the object is a resource and has attributes.
func selectsForResource(obj ssot.Object) bool {
	return !obj.VariantBaseOnly && len(obj.Attributes) > 0
}

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
	// TODO [next PRs]: generate the table lifecycle and remove this entry. The template encodes
	// ids with EncodeResourceIdentifier, has no preview-feature wrappers, and imports through SHOW.
	// Table still uses the legacy pipe-delimited id, preview wrappers, and ImportStatePassthrough.
	// Read and Update apply column diffs rather than SET/UNSET.
	"Table": {
		AllowedGenerationParts: []genhelpers.GenerationPartNamer{PartSchema},
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
	objects := defs.All()
	if len(objects) == 0 {
		log.Fatal("no object definitions registered")
	}
	fmt.Println("resource object definitions:")
	out := make([]ResourceObject, 0, len(objects))
	for _, def := range objects {
		if !ssot.ObjectGenerators(def.Generators...).Resource.Enabled {
			continue
		}
		if def.VariantBaseOnly {
			continue
		}
		fmt.Printf(" - %s\n", def.Name)
		out = append(out, ResourceObject{
			Object:                   def,
			ObjectGenerationSettings: resourceGenerationSettings[def.Name],
		})
	}
	return out
}

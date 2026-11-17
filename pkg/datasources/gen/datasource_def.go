package gen

import (
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

type DataSourceObject struct {
	ssot.Object
}

func GetAllObjects() []DataSourceObject {
	if len(defs.All()) == 0 {
		log.Fatal("no object definitions registered")
	}
	fmt.Println("data source object definitions:")
	objects := defs.SelectObjects(func(g ssot.Generators) bool { return g.DataSource.Enabled }, selectsForDataSource)
	out := make([]DataSourceObject, 0, len(objects))
	for _, def := range objects {
		fmt.Printf(" - %s\n", def.Name)
		out = append(out, DataSourceObject{Object: def})
	}
	return out
}

// selectsForDataSource takes the base of a plural that has an SDK SHOW. A variant is covered by its base.
func selectsForDataSource(obj ssot.Object) bool {
	return obj.VariantOf == nil && len(obj.PairedShowColumns()) > 0 && obj.SdkConfig().Show
}

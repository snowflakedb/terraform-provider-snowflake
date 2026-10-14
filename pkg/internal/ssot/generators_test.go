package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ObjectGenerators_options(t *testing.T) {
	g := ObjectGenerators(WithSchemaShow(), WithSchemaDescribe())
	require.Equal(t, Generators{Schema: SchemaGenerator{Show: true, Describe: true}}, g)

	g = ObjectGenerators(WithSchemaShow())
	require.Equal(t, Generators{Schema: SchemaGenerator{Show: true}}, g)

	g = ObjectGenerators(WithSchemaDescribe())
	require.Equal(t, Generators{Schema: SchemaGenerator{Describe: true}}, g)

	g = ObjectGenerators(WithModel())
	require.Equal(t, Generators{Model: ModelGenerator{Enabled: true}}, g)

	g = ObjectGenerators(WithResourceAssert())
	require.Equal(t, Generators{ResourceAssert: ResourceAssertGenerator{Enabled: true}}, g)

	g = ObjectGenerators(WithShowOutputAssert(), WithDescribeOutputAssert())
	require.Equal(t, Generators{ShowOutputAssert: ShowOutputAssertGenerator{Show: true, Describe: true}}, g)

	g = ObjectGenerators(WithShowOutputAssert())
	require.Equal(t, Generators{ShowOutputAssert: ShowOutputAssertGenerator{Show: true}}, g)

	g = ObjectGenerators(WithDescribeOutputAssert())
	require.Equal(t, Generators{ShowOutputAssert: ShowOutputAssertGenerator{Describe: true}}, g)

	g = ObjectGenerators(WithObjectAssertShow(), WithObjectAssertDescribe())
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{Show: true, Describe: true}}, g)

	g = ObjectGenerators(WithObjectAssertShow())
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{Show: true}}, g)

	g = ObjectGenerators(WithObjectAssertDescribe())
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{Describe: true}}, g)

	g = ObjectGenerators(WithObjectAssertDescribe(WithNoIdentifiableObject()))
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{Describe: true, NoIdentifiableObject: true}}, g)

	g = ObjectGenerators(WithSdkShow(), WithSdkDescribe())
	require.Equal(t, Generators{Sdk: SdkGenerator{Show: true, Describe: true}}, g)

	g = ObjectGenerators(WithSdkShow())
	require.Equal(t, Generators{Sdk: SdkGenerator{Show: true}}, g)

	g = ObjectGenerators(WithSdkDescribe())
	require.Equal(t, Generators{Sdk: SdkGenerator{Describe: true}}, g)

	g = ObjectGenerators()
	require.Equal(t, Generators{}, g)
}

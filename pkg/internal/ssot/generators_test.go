package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ObjectGenerators_options(t *testing.T) {
	g := ObjectGenerators(WithSchema(SchemaShow(), SchemaDescribe()))
	require.Equal(t, Generators{Schema: SchemaGenerator{Show: true, Describe: true}}, g)

	g = ObjectGenerators(WithSchema(SchemaShow()))
	require.Equal(t, Generators{Schema: SchemaGenerator{Show: true}}, g)

	g = ObjectGenerators(WithSchema(SchemaDescribe()))
	require.Equal(t, Generators{Schema: SchemaGenerator{Describe: true}}, g)

	g = ObjectGenerators(WithSchema())
	require.Equal(t, Generators{}, g)

	g = ObjectGenerators(WithModel())
	require.Equal(t, Generators{Model: ModelGenerator{Enabled: true}}, g)

	g = ObjectGenerators(WithResourceAssert())
	require.Equal(t, Generators{ResourceAssert: ResourceAssertGenerator{Enabled: true}}, g)

	g = ObjectGenerators(WithShowOutputAssert(ShowOutputShow(), ShowOutputDescribe()))
	require.Equal(t, Generators{ShowOutputAssert: ShowOutputAssertGenerator{Show: true, Describe: true}}, g)

	g = ObjectGenerators(WithShowOutputAssert(ShowOutputShow()))
	require.Equal(t, Generators{ShowOutputAssert: ShowOutputAssertGenerator{Show: true}}, g)

	g = ObjectGenerators(WithShowOutputAssert(ShowOutputDescribe()))
	require.Equal(t, Generators{ShowOutputAssert: ShowOutputAssertGenerator{Describe: true}}, g)

	g = ObjectGenerators(WithShowOutputAssert())
	require.Equal(t, Generators{}, g)

	g = ObjectGenerators(WithObjectAssert(ObjectAssertShow(), ObjectAssertDescribe()))
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{Show: true, Describe: true}}, g)

	g = ObjectGenerators(WithObjectAssert(ObjectAssertShow()))
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{Show: true}}, g)

	g = ObjectGenerators(WithObjectAssert(ObjectAssertDescribe()))
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{Describe: true}}, g)

	g = ObjectGenerators(WithObjectAssert(ObjectAssertDescribe(WithNoIdentifiableObject())))
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{Describe: true, NoIdentifiableObject: true}}, g)

	g = ObjectGenerators(WithObjectAssert())
	require.Equal(t, Generators{}, g)

	g = ObjectGenerators(WithSdk(SdkShow(), SdkDescribe()))
	require.Equal(t, Generators{Sdk: SdkGenerator{Show: true, Describe: true}}, g)

	g = ObjectGenerators(WithSdk(SdkShow()))
	require.Equal(t, Generators{Sdk: SdkGenerator{Show: true}}, g)

	g = ObjectGenerators(WithSdk(SdkDescribe()))
	require.Equal(t, Generators{Sdk: SdkGenerator{Describe: true}}, g)

	g = ObjectGenerators(WithSdk())
	require.Equal(t, Generators{}, g)

	g = ObjectGenerators()
	require.Equal(t, Generators{}, g)
}

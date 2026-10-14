package ssot

type SchemaGenerator struct {
	Show     bool
	Describe bool
}

type SdkGenerator struct{}

type Generators struct {
	Schema SchemaGenerator
	Sdk    SdkGenerator
}

type GeneratorOption func(*Generators)

func WithSchemaShow() GeneratorOption {
	return func(g *Generators) { g.Schema.Show = true }
}

func WithSchemaDescribe() GeneratorOption {
	return func(g *Generators) { g.Schema.Describe = true }
}

func ObjectGenerators(opts ...GeneratorOption) Generators {
	var g Generators
	for _, opt := range opts {
		opt(&g)
	}
	return g
}

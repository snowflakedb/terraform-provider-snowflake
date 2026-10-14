package ssot

type SchemaGenerator struct {
	Show     bool
	Describe bool
}

type SdkGenerator struct {
	Show     bool
	Describe bool
}

type ModelGenerator struct {
	Enabled bool
}

type ResourceAssertGenerator struct {
	Enabled bool
}

type ShowOutputAssertGenerator struct {
	Show     bool
	Describe bool
}

type ObjectAssertGenerator struct {
	Show                 bool
	Describe             bool
	NoIdentifiableObject bool
}

type Generators struct {
	Schema           SchemaGenerator
	Sdk              SdkGenerator
	Model            ModelGenerator
	ResourceAssert   ResourceAssertGenerator
	ShowOutputAssert ShowOutputAssertGenerator
	ObjectAssert     ObjectAssertGenerator
}

type GeneratorOption func(*Generators)

func WithSchemaShow() GeneratorOption {
	return func(g *Generators) { g.Schema.Show = true }
}

func WithSchemaDescribe() GeneratorOption {
	return func(g *Generators) { g.Schema.Describe = true }
}

func WithModel() GeneratorOption {
	return func(g *Generators) { g.Model.Enabled = true }
}

func WithResourceAssert() GeneratorOption {
	return func(g *Generators) { g.ResourceAssert.Enabled = true }
}

func WithShowOutputAssert() GeneratorOption {
	return func(g *Generators) { g.ShowOutputAssert.Show = true }
}

func WithDescribeOutputAssert() GeneratorOption {
	return func(g *Generators) { g.ShowOutputAssert.Describe = true }
}

func WithObjectAssertShow() GeneratorOption {
	return func(g *Generators) { g.ObjectAssert.Show = true }
}

type ObjectAssertDescribeOption func(*ObjectAssertGenerator)

func WithObjectAssertDescribe(opts ...ObjectAssertDescribeOption) GeneratorOption {
	return func(g *Generators) {
		g.ObjectAssert.Describe = true
		for _, opt := range opts {
			opt(&g.ObjectAssert)
		}
	}
}

func WithNoIdentifiableObject() ObjectAssertDescribeOption {
	return func(a *ObjectAssertGenerator) { a.NoIdentifiableObject = true }
}

func WithSdkShow() GeneratorOption {
	return func(g *Generators) { g.Sdk.Show = true }
}

func WithSdkDescribe() GeneratorOption {
	return func(g *Generators) { g.Sdk.Describe = true }
}

func ObjectGenerators(opts ...GeneratorOption) Generators {
	var g Generators
	for _, opt := range opts {
		opt(&g)
	}
	return g
}

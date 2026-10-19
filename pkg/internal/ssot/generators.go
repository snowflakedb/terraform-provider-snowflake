package ssot

type SchemaGenerator struct {
	Show     bool
	Describe bool
}

type SdkGenerator struct {
	Show     bool
	Describe bool

	ShowDoc                string
	ShowSql                string
	ShowTerse              bool
	ShowOptionalSql        []string
	ShowOptionalLike       bool
	ShowOptionalExtendedIn bool
	ShowOptionalStartsWith bool
	ShowOptionalLimitFrom  bool
	ShowOptionalLimit      bool
	ShowByIDFiltering      []SdkShowByIDFilteringKind
	ShowQueryStructName    string

	DescribeDoc             string
	DescribeSql             string
	DescribeMapping         SdkDescriptionMapping
	DescribeQueryStructName string
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

// ObjectParametersAssertGenerator emits Snowflake object parameter assertions from the parameter catalog.
// Constants are sdk.<Name>Parameter.
type ObjectParametersAssertGenerator struct {
	Enabled bool
}

// ResourceParametersAssertGenerator emits resource parameter assertions from the parameter catalog.
// Constants are sdk.<Name>Parameter.
type ResourceParametersAssertGenerator struct {
	Enabled bool
}

type Generators struct {
	Schema                   SchemaGenerator
	Sdk                      SdkGenerator
	Model                    ModelGenerator
	ResourceAssert           ResourceAssertGenerator
	ShowOutputAssert         ShowOutputAssertGenerator
	ObjectAssert             ObjectAssertGenerator
	ObjectParametersAssert   ObjectParametersAssertGenerator
	ResourceParametersAssert ResourceParametersAssertGenerator
}

type GeneratorOption func(*Generators)

type SchemaOption func(*SchemaGenerator)

func SchemaShow() SchemaOption {
	return func(s *SchemaGenerator) { s.Show = true }
}

func SchemaDescribe() SchemaOption {
	return func(s *SchemaGenerator) { s.Describe = true }
}

func WithSchema(opts ...SchemaOption) GeneratorOption {
	return func(g *Generators) {
		for _, opt := range opts {
			opt(&g.Schema)
		}
	}
}

type SdkOption func(*SdkGenerator)

type SdkShowOption func(*SdkGenerator)

type SdkDescribeOption func(*SdkGenerator)

func SdkShow(opts ...SdkShowOption) SdkOption {
	return func(s *SdkGenerator) {
		s.Show = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

func SdkDescribe(opts ...SdkDescribeOption) SdkOption {
	return func(s *SdkGenerator) {
		s.Describe = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

func WithSdk(opts ...SdkOption) GeneratorOption {
	return func(g *Generators) {
		for _, opt := range opts {
			opt(&g.Sdk)
		}
	}
}

type ModelOption func(*ModelGenerator)

func WithModel(opts ...ModelOption) GeneratorOption {
	return func(g *Generators) {
		g.Model.Enabled = true
		for _, opt := range opts {
			opt(&g.Model)
		}
	}
}

type ResourceAssertOption func(*ResourceAssertGenerator)

func WithResourceAssert(opts ...ResourceAssertOption) GeneratorOption {
	return func(g *Generators) {
		g.ResourceAssert.Enabled = true
		for _, opt := range opts {
			opt(&g.ResourceAssert)
		}
	}
}

type ShowOutputAssertOption func(*ShowOutputAssertGenerator)

func ShowOutputShow() ShowOutputAssertOption {
	return func(s *ShowOutputAssertGenerator) { s.Show = true }
}

func ShowOutputDescribe() ShowOutputAssertOption {
	return func(s *ShowOutputAssertGenerator) { s.Describe = true }
}

func WithShowOutputAssert(opts ...ShowOutputAssertOption) GeneratorOption {
	return func(g *Generators) {
		for _, opt := range opts {
			opt(&g.ShowOutputAssert)
		}
	}
}

type ObjectAssertOption func(*ObjectAssertGenerator)

type ObjectAssertDescribeOption func(*ObjectAssertGenerator)

func ObjectAssertShow() ObjectAssertOption {
	return func(a *ObjectAssertGenerator) { a.Show = true }
}

func ObjectAssertDescribe(opts ...ObjectAssertDescribeOption) ObjectAssertOption {
	return func(a *ObjectAssertGenerator) {
		a.Describe = true
		for _, opt := range opts {
			opt(a)
		}
	}
}

func WithNoIdentifiableObject() ObjectAssertDescribeOption {
	return func(a *ObjectAssertGenerator) { a.NoIdentifiableObject = true }
}

func WithObjectAssert(opts ...ObjectAssertOption) GeneratorOption {
	return func(g *Generators) {
		for _, opt := range opts {
			opt(&g.ObjectAssert)
		}
	}
}

type ObjectParametersAssertOption func(*ObjectParametersAssertGenerator)

func WithObjectParametersAssert(opts ...ObjectParametersAssertOption) GeneratorOption {
	return func(g *Generators) {
		g.ObjectParametersAssert.Enabled = true
		for _, opt := range opts {
			opt(&g.ObjectParametersAssert)
		}
	}
}

type ResourceParametersAssertOption func(*ResourceParametersAssertGenerator)

func WithResourceParametersAssert(opts ...ResourceParametersAssertOption) GeneratorOption {
	return func(g *Generators) {
		g.ResourceParametersAssert.Enabled = true
		for _, opt := range opts {
			opt(&g.ResourceParametersAssert)
		}
	}
}

func ObjectGenerators(opts ...GeneratorOption) Generators {
	var g Generators
	for _, opt := range opts {
		opt(&g)
	}
	return g
}

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

	Drop                bool
	DropDoc             string
	DropSql             string
	DropQueryStructName string
	DropSafelyHook      bool
	DropSafelyForce     bool
	DropCascadeRestrict bool

	Create                      bool
	CreateDoc                   string
	CreateSql                   string
	CreateQueryStructName       string
	CreateWithoutIfNotExists    bool
	CreateWithoutOrReplace      bool
	CreateAdditionalValidations bool

	Alter                          bool
	AlterDoc                       string
	AlterSql                       string
	AlterQueryStructName           string
	AlterWithoutIfExists           bool
	AlterSetList                   bool
	AlterSetAdditionalValidations  bool
	AlterWithoutExactlyOneValueSet bool

	CustomOperations bool
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

type SdkDropOption func(*SdkGenerator)

func SdkDrop(opts ...SdkDropOption) SdkOption {
	return func(s *SdkGenerator) {
		s.Drop = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

func SdkDropSafelyHook() SdkDropOption {
	return func(s *SdkGenerator) { s.DropSafelyHook = true }
}

func SdkDropSafelyForce() SdkDropOption {
	return func(s *SdkGenerator) { s.DropSafelyForce = true }
}

// SdkDropCascadeRestrict adds optional CASCADE and RESTRICT, and marks them as conflicting.
func SdkDropCascadeRestrict() SdkDropOption {
	return func(s *SdkGenerator) { s.DropCascadeRestrict = true }
}

type SdkCreateOption func(*SdkGenerator)

func SdkCreate(opts ...SdkCreateOption) SdkOption {
	return func(s *SdkGenerator) {
		s.Create = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

// SdkCreateDoc overrides the default CREATE doc URL.
func SdkCreateDoc(url string) SdkCreateOption {
	return func(s *SdkGenerator) { s.CreateDoc = url }
}

// SdkCreateSql overrides the default CREATE command words.
func SdkCreateSql(sql string) SdkCreateOption {
	return func(s *SdkGenerator) { s.CreateSql = sql }
}

// SdkCreateQueryStructName overrides the default CREATE query struct name.
func SdkCreateQueryStructName(name string) SdkCreateOption {
	return func(s *SdkGenerator) { s.CreateQueryStructName = name }
}

// SdkCreateWithoutIfNotExists leaves IF NOT EXISTS out. CREATE includes it unless this option is set.
func SdkCreateWithoutIfNotExists() SdkCreateOption {
	return func(s *SdkGenerator) { s.CreateWithoutIfNotExists = true }
}

// SdkCreateWithoutOrReplace leaves OR REPLACE out. CREATE includes it unless this option is set.
func SdkCreateWithoutOrReplace() SdkCreateOption {
	return func(s *SdkGenerator) { s.CreateWithoutOrReplace = true }
}

func SdkCreateAdditionalValidations() SdkCreateOption {
	return func(s *SdkGenerator) { s.CreateAdditionalValidations = true }
}

type SdkAlterOption func(*SdkGenerator)

func SdkAlter(opts ...SdkAlterOption) SdkOption {
	return func(s *SdkGenerator) {
		s.Alter = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

// SdkAlterDoc overrides the default ALTER doc URL.
func SdkAlterDoc(url string) SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterDoc = url }
}

// SdkAlterSql overrides the default ALTER command words.
func SdkAlterSql(sql string) SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterSql = sql }
}

// SdkAlterQueryStructName overrides the default ALTER query struct name.
func SdkAlterQueryStructName(name string) SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterQueryStructName = name }
}

// SdkAlterWithoutIfExists leaves IF EXISTS out. ALTER includes it unless this option is set.
func SdkAlterWithoutIfExists() SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterWithoutIfExists = true }
}

// SdkAlterSetList makes SET a list with no parentheses. SET is a keyword unless this option is set.
// TODO [next PRs]: Check why SET is a keyword for some objects and a list for others.
func SdkAlterSetList() SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterSetList = true }
}

// SdkAlterSetAdditionalValidations marks the SET struct with additional validations.
func SdkAlterSetAdditionalValidations() SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterSetAdditionalValidations = true }
}

// SdkAlterWithoutExactlyOneValueSet leaves out the check that exactly one alter branch is set.
// ALTER includes that check unless this option is set.
func SdkAlterWithoutExactlyOneValueSet() SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterWithoutExactlyOneValueSet = true }
}

// SdkCustomOperations emits Object.CustomOperations. Each entry is one custom operation.
// The query struct name is the kind plus Name. The object words are the uppercase words of Name.
// The doc URL is the kind slug, a hyphen, and the name slug, unless the spec sets Doc.
func SdkCustomOperations() SdkOption {
	return func(s *SdkGenerator) { s.CustomOperations = true }
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

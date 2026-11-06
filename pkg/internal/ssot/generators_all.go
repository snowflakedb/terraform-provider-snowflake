package ssot

// WithAllGenerators enables every generator except the SDK, then applies overrides in order.
//
// Enabled: SHOW and DESCRIBE of the schema, object assertion, and show output assertion generators, the resource, the
// config model builders, the resource assertions, and the object and resource parameter assertions.
// The SDK (WithSdk) is not part of it: pass WithSdk next to it in Generators.
//
// An enabled generator that has no input for the object is skipped by its selector, so enabling all of them is safe for an
// object that has, for example, no parameters or no DESCRIBE columns. An object that is still being migrated and cannot
// generate everything should not use it.
//
// Overrides modify what is enabled. A With* option adds its settings on top of the defaults, for example
// WithSchema(SchemaShow(SchemaShowListFields(&attrs.Tables))). A Without* option turns a generator, or one side of it, off
// and drops the settings that belong to it.
func WithAllGenerators(overrides ...GeneratorOption) GeneratorOption {
	return func(g *Generators) {
		g.Schema.Show = true
		g.Schema.Describe = true
		g.ObjectAssert.Show = true
		g.ObjectAssert.Describe = true
		g.ShowOutputAssert.Show = true
		g.ShowOutputAssert.Describe = true
		g.Model.Enabled = true
		g.Resource.Enabled = true
		g.ResourceAssert.Enabled = true
		g.ObjectParametersAssert.Enabled = true
		g.ResourceParametersAssert.Enabled = true
		for _, override := range overrides {
			override(g)
		}
	}
}

// WithoutSchema turns off the SHOW and DESCRIBE schema generator and drops its settings.
func WithoutSchema() GeneratorOption {
	return func(g *Generators) { g.Schema = SchemaGenerator{} }
}

// WithoutSchemaShow turns off the SHOW schema generator and drops its settings.
func WithoutSchemaShow() GeneratorOption {
	return func(g *Generators) {
		g.Schema.Show = false
		g.Schema.ShowFromBase = false
		g.Schema.ShowManualFields = nil
		g.Schema.ShowListFields = nil
	}
}

// WithoutSchemaDescribe turns off the DESCRIBE schema generator and drops its settings.
func WithoutSchemaDescribe() GeneratorOption {
	return func(g *Generators) {
		g.Schema.Describe = false
		g.Schema.DescribeManualFields = nil
		g.Schema.DescribeListFields = nil
	}
}

// WithoutObjectAssert turns off the SHOW and DESCRIBE object assertion generator and drops its settings.
func WithoutObjectAssert() GeneratorOption {
	return func(g *Generators) { g.ObjectAssert = ObjectAssertGenerator{} }
}

// WithoutObjectAssertShow turns off the SHOW object assertion generator and drops its settings.
func WithoutObjectAssertShow() GeneratorOption {
	return func(g *Generators) {
		g.ObjectAssert.Show = false
		g.ObjectAssert.ShowSkipFields = nil
		g.ObjectAssert.ShowProjection = ""
		g.ObjectAssert.ShowRowFromBase = false
		g.ObjectAssert.ShowNestedFields = nil
	}
}

// WithoutObjectAssertDescribe turns off the DESCRIBE object assertion generator and drops its settings.
func WithoutObjectAssertDescribe() GeneratorOption {
	return func(g *Generators) {
		g.ObjectAssert.Describe = false
		g.ObjectAssert.NoIdentifiableObject = false
		g.ObjectAssert.DescribeSkipFields = nil
		g.ObjectAssert.DescribeMethod = ""
		g.ObjectAssert.DescribeOnBase = false
		g.ObjectAssert.DescribeNestedFields = nil
	}
}

// WithoutShowOutputAssert turns off the SHOW and DESCRIBE show output assertion generator and drops its settings.
func WithoutShowOutputAssert() GeneratorOption {
	return func(g *Generators) { g.ShowOutputAssert = ShowOutputAssertGenerator{} }
}

// WithoutShowOutputAssertShow turns off the SHOW show output assertion generator and drops its settings.
func WithoutShowOutputAssertShow() GeneratorOption {
	return func(g *Generators) {
		g.ShowOutputAssert.Show = false
		g.ShowOutputAssert.ShowRowFromBase = false
		g.ShowOutputAssert.ShowSkipFields = nil
	}
}

// WithoutShowOutputAssertDescribe turns off the DESCRIBE show output assertion generator and drops its settings.
func WithoutShowOutputAssertDescribe() GeneratorOption {
	return func(g *Generators) {
		g.ShowOutputAssert.Describe = false
		g.ShowOutputAssert.DescribeSkipFields = nil
	}
}

// WithoutModel turns off the config model builder generator and drops its settings.
func WithoutModel() GeneratorOption {
	return func(g *Generators) { g.Model = ModelGenerator{} }
}

// WithoutResource turns off the resource generator.
func WithoutResource() GeneratorOption {
	return func(g *Generators) { g.Resource = ResourceGenerator{} }
}

// WithoutResourceAssert turns off the resource assertion generator.
func WithoutResourceAssert() GeneratorOption {
	return func(g *Generators) { g.ResourceAssert = ResourceAssertGenerator{} }
}

// WithoutObjectParametersAssert turns off the object parameter assertion generator.
func WithoutObjectParametersAssert() GeneratorOption {
	return func(g *Generators) { g.ObjectParametersAssert = ObjectParametersAssertGenerator{} }
}

// WithoutResourceParametersAssert turns off the resource parameter assertion generator.
func WithoutResourceParametersAssert() GeneratorOption {
	return func(g *Generators) { g.ResourceParametersAssert = ResourceParametersAssertGenerator{} }
}

package ssot

type Object struct {
	Name                                string
	NamePlural                          string
	DescribeFailRead                    bool
	DescribeUsesShowId                  bool
	HasPartialOnUpdate                  bool
	AdditionalDescribeOutputDescription string
	IdType                              string
	IdentityNameForceNew                bool
	AdditionalIdentityNameDescription   string
	Attributes                          []Attribute
	AttributesModification              bool
	ResourceHooks                       []ResourceHookOption
	Generators                          []GeneratorOption
	Show                                OutputSpec
	Describe                            OutputSpec
}

func (d Object) ObjectName() string { return d.Name }

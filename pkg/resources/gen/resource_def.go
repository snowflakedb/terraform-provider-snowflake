package gen

type ResourceDef struct {
	name               string
	describeFailRead   bool
	hasPartialOnUpdate bool
	hooks              []HookOption
}

func (d ResourceDef) ObjectName() string { return d.name }

func GetAllObjects() []ResourceDef {
	return []ResourceDef{
		computePoolDef,
	}
}

var computePoolDef = ResourceDef{
	name:             "ComputePool",
	describeFailRead: true,
}

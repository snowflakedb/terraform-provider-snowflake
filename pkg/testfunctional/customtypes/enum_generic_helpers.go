package customtypes

type EnumCreator[T ~string] interface {
	~string
	FromString(string) (T, error)
}

type dummyEnumType string

func (e dummyEnumType) FromString(s string) (dummyEnumType, error) {
	return dummyEnumType(s), nil
}

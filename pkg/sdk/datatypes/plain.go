package datatypes

// plainDataType is the common implementation of the data types without any attributes (e.g. BOOLEAN).
// underlyingType keeps the used synonym (e.g. DOUBLE), legacyType keeps the base type (e.g. FLOAT).
type plainDataType struct {
	underlyingType string
	legacyType     string
}

func (t *plainDataType) ToSql() string {
	return t.underlyingType
}

func (t *plainDataType) ToLegacyDataTypeSql() string {
	return t.legacyType
}

func (t *plainDataType) Canonical() string {
	return t.legacyType
}

func (t *plainDataType) ToSqlWithoutUnknowns() string {
	return t.underlyingType
}

func (t *plainDataType) isPlain() {}

type plain interface {
	isPlain()
}

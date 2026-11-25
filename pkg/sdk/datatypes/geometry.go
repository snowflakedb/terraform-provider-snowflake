package datatypes

// GeometryDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-geospatial#geometry-data-type
// It does not have synonyms. It does not have any attributes.
type GeometryDataType struct {
	plainDataType
}

func NewGeometryDataType() DataType {
	return newGeometryDataType(GeometryLegacyDataType)
}

func newGeometryDataType(underlyingType string) *GeometryDataType {
	return &GeometryDataType{plainDataType{underlyingType, GeometryLegacyDataType}}
}

var GeometryDataTypeSynonyms = []string{GeometryLegacyDataType}

func parseGeometryDataTypeRaw(raw sanitizedDataTypeRaw) (*GeometryDataType, error) {
	return newGeometryDataType(raw.matchedByType), nil
}

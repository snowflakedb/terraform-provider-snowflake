package datatypes

// GeographyDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-geospatial#geography-data-type
// It does not have synonyms. It does not have any attributes.
type GeographyDataType struct {
	plainDataType
}

func NewGeographyDataType() DataType {
	return newGeographyDataType(GeographyLegacyDataType)
}

func newGeographyDataType(underlyingType string) *GeographyDataType {
	return &GeographyDataType{plainDataType{underlyingType, GeographyLegacyDataType}}
}

var GeographyDataTypeSynonyms = []string{GeographyLegacyDataType}

func parseGeographyDataTypeRaw(raw sanitizedDataTypeRaw) (*GeographyDataType, error) {
	return newGeographyDataType(raw.matchedByType), nil
}

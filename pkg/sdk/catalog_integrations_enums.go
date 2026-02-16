package sdk

type CatalogSource string

var (
	CatalogSourceObjectStore CatalogSource = "OBJECT_STORE"
	CatalogSourceGlue        CatalogSource = "GLUE"
	CatalogSourceIcebergRest CatalogSource = "ICEBERG_REST"
	CatalogSourcePolaris     CatalogSource = "POLARIS"
	CatalogSourceSapBdc      CatalogSource = "SAP_BDC"
)

type TableFormat string

var (
	TableFormatIceberg TableFormat = "ICEBERG"
	TableFormatDelta   TableFormat = "DELTA"
)

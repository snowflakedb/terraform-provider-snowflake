package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FileFormat_ToFileFormatType(t *testing.T) {
	type test struct {
		input string
		want  FileFormatType
	}

	valid := []test{
		// case insensitive.
		{input: "csv", want: FileFormatTypeCsv},

		// Supported Values
		{input: "CSV", want: FileFormatTypeCsv},
		{input: "JSON", want: FileFormatTypeJson},
		{input: "AVRO", want: FileFormatTypeAvro},
		{input: "ORC", want: FileFormatTypeOrc},
		{input: "PARQUET", want: FileFormatTypeParquet},
		{input: "XML", want: FileFormatTypeXml},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToFileFormatType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToFileFormatType(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_FileFormat_ToBinaryFormat(t *testing.T) {
	type test struct {
		input string
		want  BinaryFormat
	}

	valid := []test{
		// case insensitive.
		{input: "hex", want: BinaryFormatHex},

		// Supported Values
		{input: "HEX", want: BinaryFormatHex},
		{input: "BASE64", want: BinaryFormatBase64},
		{input: "UTF8", want: BinaryFormatUtf8},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToBinaryFormat(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToBinaryFormat(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_FileFormat_ToCsvCompression(t *testing.T) {
	type test struct {
		input string
		want  CsvCompression
	}

	valid := []test{
		// case insensitive.
		{input: "gzip", want: CsvCompressionGzip},

		// Supported Values
		{input: "AUTO", want: CsvCompressionAuto},
		{input: "GZIP", want: CsvCompressionGzip},
		{input: "BZ2", want: CsvCompressionBz2},
		{input: "BROTLI", want: CsvCompressionBrotli},
		{input: "ZSTD", want: CsvCompressionZstd},
		{input: "DEFLATE", want: CsvCompressionDeflate},
		{input: "RAW_DEFLATE", want: CsvCompressionRawDeflate},
		{input: "NONE", want: CsvCompressionNone},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToCsvCompression(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToCsvCompression(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_FileFormat_ToCsvEncoding(t *testing.T) {
	type test struct {
		input string
		want  CsvEncoding
	}

	valid := []test{
		// case insensitive.
		{input: "utf8", want: CsvEncodingUtf8},

		// Supported Values
		{input: "BIG5", want: CsvEncodingBig5},
		{input: "EUCJP", want: CsvEncodingEucjp},
		{input: "EUCKR", want: CsvEncodingEuckr},
		{input: "GB18030", want: CsvEncodingGb18030},
		{input: "IBM420", want: CsvEncodingIbm420},
		{input: "IBM424", want: CsvEncodingIbm424},
		{input: "ISO2022CN", want: CsvEncodingIso2022cn},
		{input: "ISO2022JP", want: CsvEncodingIso2022jp},
		{input: "ISO2022KR", want: CsvEncodingIso2022kr},
		{input: "ISO88591", want: CsvEncodingIso88591},
		{input: "ISO88592", want: CsvEncodingIso88592},
		{input: "ISO88595", want: CsvEncodingIso88595},
		{input: "ISO88596", want: CsvEncodingIso88596},
		{input: "ISO88597", want: CsvEncodingIso88597},
		{input: "ISO88598", want: CsvEncodingIso88598},
		{input: "ISO88599", want: CsvEncodingIso88599},
		{input: "ISO885915", want: CsvEncodingIso885915},
		{input: "KOI8R", want: CsvEncodingKoi8r},
		{input: "SHIFTJIS", want: CsvEncodingShiftjis},
		{input: "UTF8", want: CsvEncodingUtf8},
		{input: "UTF16", want: CsvEncodingUtf16},
		{input: "UTF16BE", want: CsvEncodingUtf16be},
		{input: "UTF16LE", want: CsvEncodingUtf16le},
		{input: "UTF32", want: CsvEncodingUtf32},
		{input: "UTF32BE", want: CsvEncodingUtf32be},
		{input: "UTF32LE", want: CsvEncodingUtf32le},
		{input: "WINDOWS1250", want: CsvEncodingWindows1250},
		{input: "WINDOWS1251", want: CsvEncodingWindows1251},
		{input: "WINDOWS1252", want: CsvEncodingWindows1252},
		{input: "WINDOWS1253", want: CsvEncodingWindows1253},
		{input: "WINDOWS1254", want: CsvEncodingWindows1254},
		{input: "WINDOWS1255", want: CsvEncodingWindows1255},
		{input: "WINDOWS1256", want: CsvEncodingWindows1256},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToCsvEncoding(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToCsvEncoding(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_FileFormat_ToJsonCompression(t *testing.T) {
	type test struct {
		input string
		want  JsonCompression
	}

	valid := []test{
		// case insensitive.
		{input: "gzip", want: JsonCompressionGzip},

		// Supported Values
		{input: "AUTO", want: JsonCompressionAuto},
		{input: "GZIP", want: JsonCompressionGzip},
		{input: "BZ2", want: JsonCompressionBz2},
		{input: "BROTLI", want: JsonCompressionBrotli},
		{input: "ZSTD", want: JsonCompressionZstd},
		{input: "DEFLATE", want: JsonCompressionDeflate},
		{input: "RAW_DEFLATE", want: JsonCompressionRawDeflate},
		{input: "NONE", want: JsonCompressionNone},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToJsonCompression(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToJsonCompression(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_FileFormat_ToAvroCompression(t *testing.T) {
	type test struct {
		input string
		want  AvroCompression
	}

	valid := []test{
		// case insensitive.
		{input: "gzip", want: AvroCompressionGzip},

		// Supported Values
		{input: "AUTO", want: AvroCompressionAuto},
		{input: "GZIP", want: AvroCompressionGzip},
		{input: "BROTLI", want: AvroCompressionBrotli},
		{input: "ZSTD", want: AvroCompressionZstd},
		{input: "DEFLATE", want: AvroCompressionDeflate},
		{input: "RAW_DEFLATE", want: AvroCompressionRawDeflate},
		{input: "NONE", want: AvroCompressionNone},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToAvroCompression(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToAvroCompression(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_FileFormat_ToParquetCompression(t *testing.T) {
	type test struct {
		input string
		want  ParquetCompression
	}

	valid := []test{
		// case insensitive.
		{input: "snappy", want: ParquetCompressionSnappy},

		// Supported Values
		{input: "AUTO", want: ParquetCompressionAuto},
		{input: "LZO", want: ParquetCompressionLzo},
		{input: "SNAPPY", want: ParquetCompressionSnappy},
		{input: "NONE", want: ParquetCompressionNone},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToParquetCompression(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToParquetCompression(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_FileFormat_ToXmlCompression(t *testing.T) {
	type test struct {
		input string
		want  XmlCompression
	}

	valid := []test{
		// case insensitive.
		{input: "gzip", want: XmlCompressionGzip},

		// Supported Values
		{input: "AUTO", want: XmlCompressionAuto},
		{input: "GZIP", want: XmlCompressionGzip},
		{input: "BZ2", want: XmlCompressionBz2},
		{input: "BROTLI", want: XmlCompressionBrotli},
		{input: "ZSTD", want: XmlCompressionZstd},
		{input: "DEFLATE", want: XmlCompressionDeflate},
		{input: "RAW_DEFLATE", want: XmlCompressionRawDeflate},
		{input: "NONE", want: XmlCompressionNone},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToXmlCompression(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToXmlCompression(tc.input)
			require.Error(t, err)
		})
	}
}

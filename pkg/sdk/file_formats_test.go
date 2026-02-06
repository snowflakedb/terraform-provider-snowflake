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
		{input: "csv", want: FileFormatTypeCSV},

		// Supported Values
		{input: "CSV", want: FileFormatTypeCSV},
		{input: "JSON", want: FileFormatTypeJSON},
		{input: "AVRO", want: FileFormatTypeAvro},
		{input: "ORC", want: FileFormatTypeORC},
		{input: "PARQUET", want: FileFormatTypeParquet},
		{input: "XML", want: FileFormatTypeXML},
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
		{input: "UTF8", want: BinaryFormatUTF8},
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
		{input: "gzip", want: CSVCompressionGzip},

		// Supported Values
		{input: "AUTO", want: CSVCompressionAuto},
		{input: "GZIP", want: CSVCompressionGzip},
		{input: "BZ2", want: CSVCompressionBz2},
		{input: "BROTLI", want: CSVCompressionBrotli},
		{input: "ZSTD", want: CSVCompressionZstd},
		{input: "DEFLATE", want: CSVCompressionDeflate},
		{input: "RAW_DEFLATE", want: CSVCompressionRawDeflate},
		{input: "NONE", want: CSVCompressionNone},
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
		{input: "utf8", want: CSVEncodingUTF8},

		// Supported Values
		{input: "BIG5", want: CSVEncodingBIG5},
		{input: "EUCJP", want: CSVEncodingEUCJP},
		{input: "EUCKR", want: CSVEncodingEUCKR},
		{input: "GB18030", want: CSVEncodingGB18030},
		{input: "IBM420", want: CSVEncodingIBM420},
		{input: "IBM424", want: CSVEncodingIBM424},
		{input: "ISO2022CN", want: CSVEncodingISO2022CN},
		{input: "ISO2022JP", want: CSVEncodingISO2022JP},
		{input: "ISO2022KR", want: CSVEncodingISO2022KR},
		{input: "ISO88591", want: CSVEncodingISO88591},
		{input: "ISO88592", want: CSVEncodingISO88592},
		{input: "ISO88595", want: CSVEncodingISO88595},
		{input: "ISO88596", want: CSVEncodingISO88596},
		{input: "ISO88597", want: CSVEncodingISO88597},
		{input: "ISO88598", want: CSVEncodingISO88598},
		{input: "ISO88599", want: CSVEncodingISO88599},
		{input: "ISO885915", want: CSVEncodingISO885915},
		{input: "KOI8R", want: CSVEncodingKOI8R},
		{input: "SHIFTJIS", want: CSVEncodingSHIFTJIS},
		{input: "UTF8", want: CSVEncodingUTF8},
		{input: "UTF16", want: CSVEncodingUTF16},
		{input: "UTF16BE", want: CSVEncodingUTF16BE},
		{input: "UTF16LE", want: CSVEncodingUTF16LE},
		{input: "UTF32", want: CSVEncodingUTF32},
		{input: "UTF32BE", want: CSVEncodingUTF32BE},
		{input: "UTF32LE", want: CSVEncodingUTF32LE},
		{input: "WINDOWS1250", want: CSVEncodingWINDOWS1250},
		{input: "WINDOWS1251", want: CSVEncodingWINDOWS1251},
		{input: "WINDOWS1252", want: CSVEncodingWINDOWS1252},
		{input: "WINDOWS1253", want: CSVEncodingWINDOWS1253},
		{input: "WINDOWS1254", want: CSVEncodingWINDOWS1254},
		{input: "WINDOWS1255", want: CSVEncodingWINDOWS1255},
		{input: "WINDOWS1256", want: CSVEncodingWINDOWS1256},
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
		{input: "gzip", want: JSONCompressionGzip},

		// Supported Values
		{input: "AUTO", want: JSONCompressionAuto},
		{input: "GZIP", want: JSONCompressionGzip},
		{input: "BZ2", want: JSONCompressionBz2},
		{input: "BROTLI", want: JSONCompressionBrotli},
		{input: "ZSTD", want: JSONCompressionZstd},
		{input: "DEFLATE", want: JSONCompressionDeflate},
		{input: "RAW_DEFLATE", want: JSONCompressionRawDeflate},
		{input: "NONE", want: JSONCompressionNone},
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

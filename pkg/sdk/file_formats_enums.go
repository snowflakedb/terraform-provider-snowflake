package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type FileFormatType string

var (
	FileFormatTypeCSV     FileFormatType = "CSV"
	FileFormatTypeJSON    FileFormatType = "JSON"
	FileFormatTypeAvro    FileFormatType = "AVRO"
	FileFormatTypeORC     FileFormatType = "ORC"
	FileFormatTypeParquet FileFormatType = "PARQUET"
	FileFormatTypeXML     FileFormatType = "XML"
)

var AllFileFormatTypes = []FileFormatType{
	FileFormatTypeCSV,
	FileFormatTypeJSON,
	FileFormatTypeAvro,
	FileFormatTypeORC,
	FileFormatTypeParquet,
	FileFormatTypeXML,
}

func ToFileFormatType(s string) (FileFormatType, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllFileFormatTypes, FileFormatType(s)) {
		return "", fmt.Errorf("invalid file format type: %s", s)
	}
	return FileFormatType(s), nil
}

type BinaryFormat string

var (
	BinaryFormatHex    BinaryFormat = "HEX"
	BinaryFormatBase64 BinaryFormat = "BASE64"
	BinaryFormatUTF8   BinaryFormat = "UTF8"
)

var AllBinaryFormats = []BinaryFormat{
	BinaryFormatHex,
	BinaryFormatBase64,
	BinaryFormatUTF8,
}

func ToBinaryFormat(s string) (BinaryFormat, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllBinaryFormats, BinaryFormat(s)) {
		return "", fmt.Errorf("invalid binary format: %s", s)
	}
	return BinaryFormat(s), nil
}

type CsvCompression string

var (
	CSVCompressionAuto       CsvCompression = "AUTO"
	CSVCompressionGzip       CsvCompression = "GZIP"
	CSVCompressionBz2        CsvCompression = "BZ2"
	CSVCompressionBrotli     CsvCompression = "BROTLI"
	CSVCompressionZstd       CsvCompression = "ZSTD"
	CSVCompressionDeflate    CsvCompression = "DEFLATE"
	CSVCompressionRawDeflate CsvCompression = "RAW_DEFLATE"
	CSVCompressionNone       CsvCompression = "NONE"
)

var AllCsvCompressions = []CsvCompression{
	CSVCompressionAuto,
	CSVCompressionGzip,
	CSVCompressionBz2,
	CSVCompressionBrotli,
	CSVCompressionZstd,
	CSVCompressionDeflate,
	CSVCompressionRawDeflate,
	CSVCompressionNone,
}

func ToCsvCompression(s string) (CsvCompression, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllCsvCompressions, CsvCompression(s)) {
		return "", fmt.Errorf("invalid csv compression: %s", s)
	}
	return CsvCompression(s), nil
}

type CsvEncoding string

var (
	CSVEncodingBIG5        CsvEncoding = "BIG5"
	CSVEncodingEUCJP       CsvEncoding = "EUCJP"
	CSVEncodingEUCKR       CsvEncoding = "EUCKR"
	CSVEncodingGB18030     CsvEncoding = "GB18030"
	CSVEncodingIBM420      CsvEncoding = "IBM420"
	CSVEncodingIBM424      CsvEncoding = "IBM424"
	CSVEncodingISO2022CN   CsvEncoding = "ISO2022CN"
	CSVEncodingISO2022JP   CsvEncoding = "ISO2022JP"
	CSVEncodingISO2022KR   CsvEncoding = "ISO2022KR"
	CSVEncodingISO88591    CsvEncoding = "ISO88591"
	CSVEncodingISO88592    CsvEncoding = "ISO88592"
	CSVEncodingISO88595    CsvEncoding = "ISO88595"
	CSVEncodingISO88596    CsvEncoding = "ISO88596"
	CSVEncodingISO88597    CsvEncoding = "ISO88597"
	CSVEncodingISO88598    CsvEncoding = "ISO88598"
	CSVEncodingISO88599    CsvEncoding = "ISO88599"
	CSVEncodingISO885915   CsvEncoding = "ISO885915"
	CSVEncodingKOI8R       CsvEncoding = "KOI8R"
	CSVEncodingSHIFTJIS    CsvEncoding = "SHIFTJIS"
	CSVEncodingUTF8        CsvEncoding = "UTF8"
	CSVEncodingUTF16       CsvEncoding = "UTF16"
	CSVEncodingUTF16BE     CsvEncoding = "UTF16BE"
	CSVEncodingUTF16LE     CsvEncoding = "UTF16LE"
	CSVEncodingUTF32       CsvEncoding = "UTF32"
	CSVEncodingUTF32BE     CsvEncoding = "UTF32BE"
	CSVEncodingUTF32LE     CsvEncoding = "UTF32LE"
	CSVEncodingWINDOWS1250 CsvEncoding = "WINDOWS1250"
	CSVEncodingWINDOWS1251 CsvEncoding = "WINDOWS1251"
	CSVEncodingWINDOWS1252 CsvEncoding = "WINDOWS1252"
	CSVEncodingWINDOWS1253 CsvEncoding = "WINDOWS1253"
	CSVEncodingWINDOWS1254 CsvEncoding = "WINDOWS1254"
	CSVEncodingWINDOWS1255 CsvEncoding = "WINDOWS1255"
	CSVEncodingWINDOWS1256 CsvEncoding = "WINDOWS1256"
)

var AllCsvEncodings = []CsvEncoding{
	CSVEncodingBIG5,
	CSVEncodingEUCJP,
	CSVEncodingEUCKR,
	CSVEncodingGB18030,
	CSVEncodingIBM420,
	CSVEncodingIBM424,
	CSVEncodingISO2022CN,
	CSVEncodingISO2022JP,
	CSVEncodingISO2022KR,
	CSVEncodingISO88591,
	CSVEncodingISO88592,
	CSVEncodingISO88595,
	CSVEncodingISO88596,
	CSVEncodingISO88597,
	CSVEncodingISO88598,
	CSVEncodingISO88599,
	CSVEncodingISO885915,
	CSVEncodingKOI8R,
	CSVEncodingSHIFTJIS,
	CSVEncodingUTF8,
	CSVEncodingUTF16,
	CSVEncodingUTF16BE,
	CSVEncodingUTF16LE,
	CSVEncodingUTF32,
	CSVEncodingUTF32BE,
	CSVEncodingUTF32LE,
	CSVEncodingWINDOWS1250,
	CSVEncodingWINDOWS1251,
	CSVEncodingWINDOWS1252,
	CSVEncodingWINDOWS1253,
	CSVEncodingWINDOWS1254,
	CSVEncodingWINDOWS1255,
	CSVEncodingWINDOWS1256,
}

func ToCsvEncoding(s string) (CsvEncoding, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllCsvEncodings, CsvEncoding(s)) {
		return "", fmt.Errorf("invalid csv encoding: %s", s)
	}
	return CsvEncoding(s), nil
}

type JsonCompression string

var (
	JSONCompressionAuto       JsonCompression = "AUTO"
	JSONCompressionGzip       JsonCompression = "GZIP"
	JSONCompressionBz2        JsonCompression = "BZ2"
	JSONCompressionBrotli     JsonCompression = "BROTLI"
	JSONCompressionZstd       JsonCompression = "ZSTD"
	JSONCompressionDeflate    JsonCompression = "DEFLATE"
	JSONCompressionRawDeflate JsonCompression = "RAW_DEFLATE"
	JSONCompressionNone       JsonCompression = "NONE"
)

var AllJsonCompressions = []JsonCompression{
	JSONCompressionAuto,
	JSONCompressionGzip,
	JSONCompressionBz2,
	JSONCompressionBrotli,
	JSONCompressionZstd,
	JSONCompressionDeflate,
	JSONCompressionRawDeflate,
	JSONCompressionNone,
}

func ToJsonCompression(s string) (JsonCompression, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllJsonCompressions, JsonCompression(s)) {
		return "", fmt.Errorf("invalid json compression: %s", s)
	}
	return JsonCompression(s), nil
}

type AvroCompression string

var (
	AvroCompressionAuto       AvroCompression = "AUTO"
	AvroCompressionGzip       AvroCompression = "GZIP"
	AvroCompressionBrotli     AvroCompression = "BROTLI"
	AvroCompressionZstd       AvroCompression = "ZSTD"
	AvroCompressionDeflate    AvroCompression = "DEFLATE"
	AvroCompressionRawDeflate AvroCompression = "RAW_DEFLATE"
	AvroCompressionNone       AvroCompression = "NONE"
)

type ParquetCompression string

var (
	ParquetCompressionAuto   ParquetCompression = "AUTO"
	ParquetCompressionLzo    ParquetCompression = "LZO"
	ParquetCompressionSnappy ParquetCompression = "SNAPPY"
	ParquetCompressionNone   ParquetCompression = "NONE"
)

type XmlCompression string

var (
	XMLCompressionAuto       XmlCompression = "AUTO"
	XMLCompressionGzip       XmlCompression = "GZIP"
	XMLCompressionBz2        XmlCompression = "BZ2"
	XMLCompressionBrotli     XmlCompression = "BROTLI"
	XMLCompressionZstd       XmlCompression = "ZSTD"
	XMLCompressionDeflate    XmlCompression = "DEFLATE"
	XMLCompressionRawDeflate XmlCompression = "RAW_DEFLATE"
	XMLCompressionNone       XmlCompression = "NONE"
)

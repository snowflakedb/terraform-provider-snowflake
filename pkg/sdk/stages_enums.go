package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type InternalStageEncryptionOption string

var (
	InternalStageEncryptionFull InternalStageEncryptionOption = "SNOWFLAKE_FULL"
	InternalStageEncryptionSSE  InternalStageEncryptionOption = "SNOWFLAKE_SSE"
)

type ExternalStageS3EncryptionOption string

var (
	ExternalStageS3EncryptionCSE       ExternalStageS3EncryptionOption = "AWS_CSE"
	ExternalStageS3EncryptionSSES3     ExternalStageS3EncryptionOption = "AWS_SSE_S3"
	ExternalStageS3EncryptionSSEKMS    ExternalStageS3EncryptionOption = "AWS_SSE_KMS"
	ExternalStageS3EncryptionNoneValue ExternalStageS3EncryptionOption = "NONE"
)

type ExternalStageGCSEncryptionOption string

var (
	ExternalStageGCSEncryptionSSEKMS    ExternalStageGCSEncryptionOption = "GCS_SSE_KMS"
	ExternalStageGCSEncryptionNoneValue ExternalStageGCSEncryptionOption = "NONE"
)

type ExternalStageAzureEncryptionOption string

var (
	ExternalStageAzureEncryptionCSE       ExternalStageAzureEncryptionOption = "AZURE_CSE"
	ExternalStageAzureEncryptionNoneValue ExternalStageAzureEncryptionOption = "NONE"
)

// TODO: move to tables
type StageCopyColumnMapOption string

var (
	StageCopyColumnMapCaseSensitive   StageCopyColumnMapOption = "CASE_SENSITIVE"
	StageCopyColumnMapCaseInsensitive StageCopyColumnMapOption = "CASE_INSENSITIVE"
	StageCopyColumnMapCaseNone        StageCopyColumnMapOption = "NONE"
)

type StageFileFormatCsvCompression string

const (
	StageFileFormatCsvCompressionAuto       StageFileFormatCsvCompression = "AUTO"
	StageFileFormatCsvCompressionGzip       StageFileFormatCsvCompression = "GZIP"
	StageFileFormatCsvCompressionBz2        StageFileFormatCsvCompression = "BZ2"
	StageFileFormatCsvCompressionBrotli     StageFileFormatCsvCompression = "BROTLI"
	StageFileFormatCsvCompressionZstd       StageFileFormatCsvCompression = "ZSTD"
	StageFileFormatCsvCompressionDeflate    StageFileFormatCsvCompression = "DEFLATE"
	StageFileFormatCsvCompressionRawDeflate StageFileFormatCsvCompression = "RAW_DEFLATE"
	StageFileFormatCsvCompressionNone       StageFileFormatCsvCompression = "NONE"
)

var AllStageFileFormatCsvCompressions = []StageFileFormatCsvCompression{
	StageFileFormatCsvCompressionAuto,
	StageFileFormatCsvCompressionGzip,
	StageFileFormatCsvCompressionBz2,
	StageFileFormatCsvCompressionBrotli,
	StageFileFormatCsvCompressionZstd,
	StageFileFormatCsvCompressionDeflate,
	StageFileFormatCsvCompressionRawDeflate,
	StageFileFormatCsvCompressionNone,
}

type StageFileFormatJsonCompression string

const (
	StageFileFormatJsonCompressionAuto       StageFileFormatJsonCompression = "AUTO"
	StageFileFormatJsonCompressionGzip       StageFileFormatJsonCompression = "GZIP"
	StageFileFormatJsonCompressionBz2        StageFileFormatJsonCompression = "BZ2"
	StageFileFormatJsonCompressionBrotli     StageFileFormatJsonCompression = "BROTLI"
	StageFileFormatJsonCompressionZstd       StageFileFormatJsonCompression = "ZSTD"
	StageFileFormatJsonCompressionDeflate    StageFileFormatJsonCompression = "DEFLATE"
	StageFileFormatJsonCompressionRawDeflate StageFileFormatJsonCompression = "RAW_DEFLATE"
	StageFileFormatJsonCompressionNone       StageFileFormatJsonCompression = "NONE"
)

var AllStageFileFormatJsonCompressions = []StageFileFormatJsonCompression{
	StageFileFormatJsonCompressionAuto,
	StageFileFormatJsonCompressionGzip,
	StageFileFormatJsonCompressionBz2,
	StageFileFormatJsonCompressionBrotli,
	StageFileFormatJsonCompressionZstd,
	StageFileFormatJsonCompressionDeflate,
	StageFileFormatJsonCompressionRawDeflate,
	StageFileFormatJsonCompressionNone,
}

type StageFileFormatAvroCompression string

const (
	StageFileFormatAvroCompressionAuto       StageFileFormatAvroCompression = "AUTO"
	StageFileFormatAvroCompressionGzip       StageFileFormatAvroCompression = "GZIP"
	StageFileFormatAvroCompressionBrotli     StageFileFormatAvroCompression = "BROTLI"
	StageFileFormatAvroCompressionZstd       StageFileFormatAvroCompression = "ZSTD"
	StageFileFormatAvroCompressionDeflate    StageFileFormatAvroCompression = "DEFLATE"
	StageFileFormatAvroCompressionRawDeflate StageFileFormatAvroCompression = "RAW_DEFLATE"
	StageFileFormatAvroCompressionNone       StageFileFormatAvroCompression = "NONE"
)

var AllStageFileFormatAvroCompressions = []StageFileFormatAvroCompression{
	StageFileFormatAvroCompressionAuto,
	StageFileFormatAvroCompressionGzip,
	StageFileFormatAvroCompressionBrotli,
	StageFileFormatAvroCompressionZstd,
	StageFileFormatAvroCompressionDeflate,
	StageFileFormatAvroCompressionRawDeflate,
	StageFileFormatAvroCompressionNone,
}

type StageFileFormatParquetCompression string

const (
	StageFileFormatParquetCompressionAuto   StageFileFormatParquetCompression = "AUTO"
	StageFileFormatParquetCompressionLzo    StageFileFormatParquetCompression = "LZO"
	StageFileFormatParquetCompressionSnappy StageFileFormatParquetCompression = "SNAPPY"
	StageFileFormatParquetCompressionNone   StageFileFormatParquetCompression = "NONE"
)

var AllStageFileFormatParquetCompressions = []StageFileFormatParquetCompression{
	StageFileFormatParquetCompressionAuto,
	StageFileFormatParquetCompressionLzo,
	StageFileFormatParquetCompressionSnappy,
	StageFileFormatParquetCompressionNone,
}

type StageFileFormatXmlCompression string

const (
	StageFileFormatXmlCompressionAuto       StageFileFormatXmlCompression = "AUTO"
	StageFileFormatXmlCompressionGzip       StageFileFormatXmlCompression = "GZIP"
	StageFileFormatXmlCompressionBz2        StageFileFormatXmlCompression = "BZ2"
	StageFileFormatXmlCompressionBrotli     StageFileFormatXmlCompression = "BROTLI"
	StageFileFormatXmlCompressionZstd       StageFileFormatXmlCompression = "ZSTD"
	StageFileFormatXmlCompressionDeflate    StageFileFormatXmlCompression = "DEFLATE"
	StageFileFormatXmlCompressionRawDeflate StageFileFormatXmlCompression = "RAW_DEFLATE"
	StageFileFormatXmlCompressionNone       StageFileFormatXmlCompression = "NONE"
)

var AllStageFileFormatXmlCompressions = []StageFileFormatXmlCompression{
	StageFileFormatXmlCompressionAuto,
	StageFileFormatXmlCompressionGzip,
	StageFileFormatXmlCompressionBz2,
	StageFileFormatXmlCompressionBrotli,
	StageFileFormatXmlCompressionZstd,
	StageFileFormatXmlCompressionDeflate,
	StageFileFormatXmlCompressionRawDeflate,
	StageFileFormatXmlCompressionNone,
}

type StageFileFormatBinaryFormat string

const (
	StageFileFormatBinaryFormatHex    StageFileFormatBinaryFormat = "HEX"
	StageFileFormatBinaryFormatBase64 StageFileFormatBinaryFormat = "BASE64"
	StageFileFormatBinaryFormatUtf8   StageFileFormatBinaryFormat = "UTF8"
)

var AllStageFileFormatBinaryFormats = []StageFileFormatBinaryFormat{
	StageFileFormatBinaryFormatHex,
	StageFileFormatBinaryFormatBase64,
	StageFileFormatBinaryFormatUtf8,
}

func ToStageFileFormatCsvCompression(s string) (StageFileFormatCsvCompression, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllStageFileFormatCsvCompressions, StageFileFormatCsvCompression(s)) {
		return "", fmt.Errorf("invalid stage file format CSV compression: %s", s)
	}
	return StageFileFormatCsvCompression(s), nil
}

func ToStageFileFormatJsonCompression(s string) (StageFileFormatJsonCompression, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllStageFileFormatJsonCompressions, StageFileFormatJsonCompression(s)) {
		return "", fmt.Errorf("invalid stage file format JSON compression: %s", s)
	}
	return StageFileFormatJsonCompression(s), nil
}

func ToStageFileFormatAvroCompression(s string) (StageFileFormatAvroCompression, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllStageFileFormatAvroCompressions, StageFileFormatAvroCompression(s)) {
		return "", fmt.Errorf("invalid stage file format AVRO compression: %s", s)
	}
	return StageFileFormatAvroCompression(s), nil
}

func ToStageFileFormatParquetCompression(s string) (StageFileFormatParquetCompression, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllStageFileFormatParquetCompressions, StageFileFormatParquetCompression(s)) {
		return "", fmt.Errorf("invalid stage file format Parquet compression: %s", s)
	}
	return StageFileFormatParquetCompression(s), nil
}

func ToStageFileFormatXmlCompression(s string) (StageFileFormatXmlCompression, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllStageFileFormatXmlCompressions, StageFileFormatXmlCompression(s)) {
		return "", fmt.Errorf("invalid stage file format XML compression: %s", s)
	}
	return StageFileFormatXmlCompression(s), nil
}

func ToStageFileFormatBinaryFormat(s string) (StageFileFormatBinaryFormat, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllStageFileFormatBinaryFormats, StageFileFormatBinaryFormat(s)) {
		return "", fmt.Errorf("invalid stage file format binary format: %s", s)
	}
	return StageFileFormatBinaryFormat(s), nil
}

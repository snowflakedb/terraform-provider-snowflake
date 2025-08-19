package main

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"log"
	"os"
	"reflect"
	"slices"
	"strconv"
)

type ConversionConfig[T, R any] struct {
	AdditionalConvertMapping func(row T, convertedValue *R)
}

func NewConversionConfigWithOpts[T, R any](opts ...func(*ConversionConfig[T, R])) *ConversionConfig[T, R] {
	config := new(ConversionConfig[T, R])
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func WithAdditionalConvertMapping[T, R any](mappingFunc func(row T, convertedValue *R)) func(*ConversionConfig[T, R]) {
	return func(config *ConversionConfig[T, R]) {
		config.AdditionalConvertMapping = mappingFunc
	}
}

type ConvertibleRowStructField struct {
	Index int
	Name  string
	Tag   string
	Kind  reflect.Kind
}

func ConvertCsvInput[T sdk.ConvertibleRowDeprecated[R], R any](csvInputFormat [][]string, opts ...func(*ConversionConfig[T, R])) []R {
	if len(csvInputFormat) < 1 {
		log.Println("CSV input is empty")
		os.Exit(1)
	}

	parseConfig := NewConversionConfigWithOpts(opts...)
	convertibleRowFields := make([]ConvertibleRowStructField, 0)
	structType := reflect.TypeFor[T]()

	for i := 0; i < structType.NumField(); i++ {
		convertibleRowFields = append(convertibleRowFields, ConvertibleRowStructField{
			Index: i,
			Name:  structType.Field(i).Name,
			Kind:  structType.Field(i).Type.Kind(),
			Tag:   structType.Field(i).Tag.Get("db"),
		})
	}

	csvHeader := csvInputFormat[0]

	// Remove fields that are not present in the CSV input
	convertibleRowFields = slices.DeleteFunc(convertibleRowFields, func(field ConvertibleRowStructField) bool {
		return !slices.Contains(csvHeader, field.Tag)
	})

	// Sort fields by header order
	slices.SortFunc(convertibleRowFields, func(a, b ConvertibleRowStructField) int {
		return slices.Index(csvHeader, a.Tag) - slices.Index(csvHeader, b.Tag)
	})

	// TODO: Separate and create another function with below logic?

	result := make([]R, 0)
	for _, csvRow := range csvInputFormat[1:] {
		var row T
		for i, csvColumnValue := range csvRow {
			field := reflect.ValueOf(&row).Elem().Field(convertibleRowFields[i].Index)

			switch convertibleRowFields[i].Kind {
			case reflect.String:
				field.SetString(csvColumnValue)
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(csvColumnValue)
				if err != nil {
					log.Printf("Error parsing boolean value for field %s: %v", convertibleRowFields[i].Name, err)
				}
				field.SetBool(boolValue)
			default:
				log.Printf("Unsupported type %s for field %s", convertibleRowFields[i].Kind, convertibleRowFields[i].Name)
			}
		}

		// Calling convert option 1 interface conversion (Requires interface to be exposed by the SDK package and objects to expose exported Convert method)
		//var convertedValue R
		//if v, ok := any(row).(sdk.ConvertibleRowDeprecated[R]); ok {
		//	convertedValue = *v.Convert()
		//}

		// Calling convert option 2 golinkname
		convertedValue := convertByGolinkname[T](row)

		// Calling convert option 3 reflection - doesn't work because it doesn't see unexported methods

		// Calling convert option 4 use something proposed in https://github.com/alangpierce/go-forceexport/blob/master/forceexport.go
		// Ref: https://www.alangpierce.com/blog/2016/03/17/adventures-in-go-accessing-unexported-functions/

		if parseConfig.AdditionalConvertMapping != nil {
			parseConfig.AdditionalConvertMapping(row, convertedValue)
		}

		result = append(result, *convertedValue)
	}
	return result
}

// TODO: The correct mapping function could be passed somewhere earlier, e.g., at the ConvertCsvInput level.
func convertByGolinkname[T sdk.ConvertibleRowDeprecated[R], R any](row T) *R {
	switch typedRow := any(row).(type) {
	case sdk.GrantRow:
		return any(convertGrantRow(&typedRow)).(*R)
	}
	panic(fmt.Sprintf("unsupported type for conversion: %T", row))
}

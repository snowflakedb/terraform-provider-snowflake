package main

import (
	"fmt"
	"log"
	"reflect"
	"slices"
	"strconv"
)

type ConvertibleCsvRow[T any] interface {
	convert() (*T, error)
}

type Converter[T ConvertibleCsvRow[R], R any] struct {
	AdditionalConvertMapping func(row T, convertedValue *R)
}

func NewConversionConfigWithOpts[T ConvertibleCsvRow[R], R any](opts ...func(*Converter[T, R])) *Converter[T, R] {
	config := new(Converter[T, R])
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func WithAdditionalConvertMapping[T ConvertibleCsvRow[R], R any](mappingFunc func(row T, convertedValue *R)) func(*Converter[T, R]) {
	return func(config *Converter[T, R]) {
		config.AdditionalConvertMapping = mappingFunc
	}
}

type ConvertibleRowStructField struct {
	Index []int
	Name  string
	Tag   string
	Kind  reflect.Kind
}

func ConvertCsvInput[T ConvertibleCsvRow[R], R any](csvInputFormat [][]string, opts ...func(*Converter[T, R])) ([]R, error) {
	if len(csvInputFormat) < 1 {
		return nil, fmt.Errorf("CSV input is empty")
	}

	converter := NewConversionConfigWithOpts(opts...)
	convertibleRowFields := make([]ConvertibleRowStructField, 0)
	structType := reflect.TypeFor[T]()

	for i := 0; i < structType.NumField(); i++ {
		convertibleRowFields = append(convertibleRowFields, ConvertibleRowStructField{
			Index: structType.Field(i).Index,
			Name:  structType.Field(i).Name,
			Kind:  structType.Field(i).Type.Kind(),
			Tag:   structType.Field(i).Tag.Get("csv"),
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

	// Based on the headers included for mapping, prepare a list of indices that should be taken into consideration when mapping a CSV row
	columnIndices := make([]int, 0)
	for index, headerValue := range csvHeader {
		if slices.ContainsFunc(convertibleRowFields, func(field ConvertibleRowStructField) bool { return field.Tag == headerValue }) {
			columnIndices = append(columnIndices, index)
		}
	}

	result := make([]R, 0)

csvRowLoop:
	for _, csvRow := range csvInputFormat[1:] {
		var row T
		for i := range columnIndices {
			csvColumnValue := csvRow[columnIndices[i]]
			field := reflect.ValueOf(&row).Elem().FieldByIndex(convertibleRowFields[i].Index)

			switch convertibleRowFields[i].Kind {
			case reflect.String:
				field.SetString(csvColumnValue)
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(csvColumnValue)
				if err != nil {
					log.Printf("Error parsing boolean value for field %s: %v. Skipping this row (will not be included in the final generation).", convertibleRowFields[i].Name, err)
					continue csvRowLoop
				}
				field.SetBool(boolValue)
			default:
				log.Printf("Unsupported type %s for field %s. Skipping this row (will not be included in the final generation).", convertibleRowFields[i].Kind, convertibleRowFields[i].Name)
				continue csvRowLoop
			}
		}

		var convertedValue *R
		if v, err := row.convert(); err != nil {
			return nil, err
		} else {
			convertedValue = v
		}

		if converter.AdditionalConvertMapping != nil {
			converter.AdditionalConvertMapping(row, convertedValue)
		}

		result = append(result, *convertedValue)
	}

	return result, nil
}

package sdk

import (
	"errors"
	"fmt"
)

var (
	_ validatable = (*CreateExternalTableOptions)(nil)
	_ validatable = (*CreateWithManualPartitioningExternalTableOptions)(nil)
	_ validatable = (*CreateDeltaLakeExternalTableOptions)(nil)
	_ validatable = (*CreateExternalTableUsingTemplateOptions)(nil)
	_ validatable = (*AlterExternalTableOptions)(nil)
	_ validatable = (*AlterExternalTablePartitionOptions)(nil)
	_ validatable = (*DropExternalTableOptions)(nil)
	_ validatable = (*ShowExternalTableOptions)(nil)
	_ validatable = (*describeExternalTableColumnsOptions)(nil)
	_ validatable = (*describeExternalTableStageOptions)(nil)
)

func (opts *CreateExternalTableOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateExternalTableOptions", "OrReplace", "IfNotExists"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.Location) {
		errs = append(errs, errNotSet("CreateExternalTableOptions", "Location"))
	}
	if !exactlyOneValueSet(opts.RawFileFormat, opts.FileFormat) {
		errs = append(errs, errExactlyOneOf("CreateExternalTableOptions", "RawFileFormat", "FileFormat"))
	}
	if valueSet(opts.FileFormat) {
		for i, ff := range opts.FileFormat {
			if !valueSet(ff.Name) && !valueSet(ff.Type) {
				errs = append(errs, errNotSet(fmt.Sprintf("CreateExternalTableOptions.FileFormat[%d]", i), "Name or Type"))
			}
			if valueSet(ff.Name) && valueSet(ff.Type) {
				errs = append(errs, errOneOf(fmt.Sprintf("CreateExternalTableOptions.FileFormat[%d]", i), "Name or Type"))
			}
		}
	}
	return errors.Join(errs...)
}

func (opts *CreateWithManualPartitioningExternalTableOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateWithManualPartitioningExternalTableOptions", "OrReplace", "IfNotExists"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.Location) {
		errs = append(errs, errNotSet("CreateWithManualPartitioningExternalTableOptions", "Location"))
	}
	if !exactlyOneValueSet(opts.RawFileFormat, opts.FileFormat) {
		errs = append(errs, errExactlyOneOf("CreateWithManualPartitioningExternalTableOptions", "RawFileFormat", "FileFormat"))
	}
	if valueSet(opts.FileFormat) {
		for i, ff := range opts.FileFormat {
			if !valueSet(ff.Name) && !valueSet(ff.Type) {
				errs = append(errs, errNotSet(fmt.Sprintf("CreateWithManualPartitioningExternalTableOptions.FileFormat[%d]", i), "Name or Type"))
			}
			if valueSet(ff.Name) && valueSet(ff.Type) {
				errs = append(errs, errOneOf(fmt.Sprintf("CreateWithManualPartitioningExternalTableOptions.FileFormat[%d]", i), "Name or Type"))
			}
		}
	}
	return errors.Join(errs...)
}

func (opts *CreateDeltaLakeExternalTableOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateDeltaLakeExternalTableOptions", "OrReplace", "IfNotExists"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.Location) {
		errs = append(errs, errNotSet("CreateDeltaLakeExternalTableOptions", "Location"))
	}
	if !exactlyOneValueSet(opts.RawFileFormat, opts.FileFormat) {
		errs = append(errs, errExactlyOneOf("CreateDeltaLakeExternalTableOptions", "RawFileFormat", "FileFormat"))
	}
	if valueSet(opts.FileFormat) {
		for i, ff := range opts.FileFormat {
			if !valueSet(ff.Name) && !valueSet(ff.Type) {
				errs = append(errs, errNotSet(fmt.Sprintf("CreateDeltaLakeExternalTableOptions.FileFormat[%d]", i), "Name or Type"))
			}
			if valueSet(ff.Name) && valueSet(ff.Type) {
				errs = append(errs, errOneOf(fmt.Sprintf("CreateDeltaLakeExternalTableOptions.FileFormat[%d]", i), "Name or Type"))
			}
		}
	}
	return errors.Join(errs...)
}

func (opts *CreateExternalTableUsingTemplateOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.Query) {
		errs = append(errs, errNotSet("CreateExternalTableUsingTemplateOptions", "Query"))
	}
	if !valueSet(opts.Location) {
		errs = append(errs, errNotSet("CreateExternalTableUsingTemplateOptions", "Location"))
	}
	if !exactlyOneValueSet(opts.RawFileFormat, opts.FileFormat) {
		errs = append(errs, errExactlyOneOf("CreateExternalTableUsingTemplateOptions", "RawFileFormat", "FileFormat"))
	}
	if valueSet(opts.FileFormat) {
		for i, ff := range opts.FileFormat {
			if !valueSet(ff.Name) && !valueSet(ff.Type) {
				errs = append(errs, errNotSet(fmt.Sprintf("CreateExternalTableUsingTemplateOptions.FileFormat[%d]", i), "Name or Type"))
			}
			if valueSet(ff.Name) && valueSet(ff.Type) {
				errs = append(errs, errOneOf(fmt.Sprintf("CreateExternalTableUsingTemplateOptions.FileFormat[%d]", i), "Name or Type"))
			}
		}
	}
	return errors.Join(errs...)
}

func (opts *AlterExternalTableOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Refresh, opts.AddFiles, opts.RemoveFiles, opts.AutoRefresh, opts.SetTag, opts.UnsetTag) {
		errs = append(errs, errExactlyOneOf("AlterExternalTableOptions", "Refresh", "AddFiles", "RemoveFiles", "AutoRefresh", "SetTag", "UnsetTag"))
	}
	return errors.Join(errs...)
}

func (opts *AlterExternalTablePartitionOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.AddPartitions, opts.DropPartition) {
		errs = append(errs, errOneOf("AlterExternalTablePartitionOptions", "AddPartitions", "DropPartition"))
	}
	return errors.Join(errs...)
}

func (opts *DropExternalTableOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.DropOption) {
		if err := opts.DropOption.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (opts *ShowExternalTableOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (opts *describeExternalTableColumnsOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *describeExternalTableStageOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (cpp *CloudProviderParams) validate() error {
	if anyValueSet(cpp.GoogleCloudStorageIntegration, cpp.MicrosoftAzureIntegration) && !exactlyOneValueSet(cpp.GoogleCloudStorageIntegration, cpp.MicrosoftAzureIntegration) {
		return errOneOf("CloudProviderParams", "GoogleCloudStorageIntegration", "MicrosoftAzureIntegration")
	}
	return nil
}

func (opts *ExternalTableFileFormat) validate() error {
	var errs []error
	if everyValueSet(opts.Name, opts.Type) {
		errs = append(errs, errOneOf("ExternalTableFileFormat", "Name", "Type"))
	}
	fields := externalTableFileFormatTypeOptionsFieldsByType(opts.Options)
	for formatType := range fields {
		if *opts.Type == formatType {
			continue
		}
		if anyValueSet(fields[formatType]...) {
			errs = append(errs, fmt.Errorf("cannot set %s fields when TYPE = %s", formatType, *opts.Type))
		}
	}
	return errors.Join(errs...)
}

func (opts *ExternalTableDropOption) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if anyValueSet(opts.Restrict, opts.Cascade) && !exactlyOneValueSet(opts.Restrict, opts.Cascade) {
		return errors.Join(errOneOf("ExternalTableDropOption", "Restrict", "Cascade"))
	}
	return nil
}

func externalTableFileFormatTypeOptionsFieldsByType(opts *ExternalTableFileFormatTypeOptions) map[ExternalTableFileFormatType][]any {
	return map[ExternalTableFileFormatType][]any{
		ExternalTableFileFormatTypeCsv: {
			opts.CSVCompression,
			opts.CSVRecordDelimiter,
			opts.CSVFieldDelimiter,
			opts.CSVSkipHeader,
			opts.CSVSkipBlankLines,
			opts.CSVEscapeUnenclosedField,
			opts.CSVTrimSpace,
			opts.CSVFieldOptionallyEnclosedBy,
			opts.CSVNullIf,
			opts.CSVEmptyFieldAsNull,
			opts.CSVEncoding,
		},
		ExternalTableFileFormatTypeJson: {
			opts.JSONCompression,
			opts.JSONAllowDuplicate,
			opts.JSONStripOuterArray,
			opts.JSONStripNullValues,
			opts.JSONReplaceInvalidCharacters,
		},
		ExternalTableFileFormatTypeAvro: {
			opts.AvroCompression,
			opts.AvroReplaceInvalidCharacters,
		},
		ExternalTableFileFormatTypeOrc: {
			opts.ORCTrimSpace,
			opts.ORCReplaceInvalidCharacters,
			opts.ORCNullIf,
		},
		ExternalTableFileFormatTypeParquet: {
			opts.ParquetCompression,
			opts.ParquetBinaryAsText,
			opts.ParquetReplaceInvalidCharacters,
		},
	}
}

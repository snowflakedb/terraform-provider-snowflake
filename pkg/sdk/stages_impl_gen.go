package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
)

var _ Stages = (*stages)(nil)

type stages struct {
	client *Client
}

func (v *stages) CreateInternal(ctx context.Context, request *CreateInternalStageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) CreateOnS3(ctx context.Context, request *CreateOnS3StageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) CreateOnGCS(ctx context.Context, request *CreateOnGCSStageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) CreateOnAzure(ctx context.Context, request *CreateOnAzureStageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) CreateOnS3Compatible(ctx context.Context, request *CreateOnS3CompatibleStageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) Alter(ctx context.Context, request *AlterStageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) AlterInternalStage(ctx context.Context, request *AlterInternalStageStageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) AlterExternalS3Stage(ctx context.Context, request *AlterExternalS3StageStageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) AlterExternalGCSStage(ctx context.Context, request *AlterExternalGCSStageStageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) AlterExternalAzureStage(ctx context.Context, request *AlterExternalAzureStageStageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) AlterDirectoryTable(ctx context.Context, request *AlterDirectoryTableStageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) Drop(ctx context.Context, request *DropStageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropStageRequest(id).WithIfExists(Bool(true))) }, ctx, id)
}

func (v *stages) Describe(ctx context.Context, id SchemaObjectIdentifier) ([]StageProperty, error) {
	opts := &DescribeStageOptions{
		name: id,
	}
	rows, err := validateAndQuery[stageDescRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[stageDescRow, StageProperty](rows), nil
}

func (v *stages) Show(ctx context.Context, request *ShowStageRequest) ([]Stage, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[stageShowRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[stageShowRow, Stage](dbRows)
	return resultList, nil
}

func (v *stages) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Stage, error) {
	request := NewShowStageRequest().
		WithLike(Like{Pattern: String(id.Name())}).
		WithIn(In{Schema: id.SchemaId()})
	stages, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(stages, func(r Stage) bool { return r.Name == id.Name() })
}

func (v *stages) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Stage, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (s *CreateInternalStageRequest) toOpts() *CreateInternalStageOptions {
	opts := &CreateInternalStageOptions{
		OrReplace:   s.OrReplace,
		Temporary:   s.Temporary,
		IfNotExists: s.IfNotExists,
		name:        s.name,

		Comment: s.Comment,
		Tag:     s.Tag,
	}
	if s.Encryption != nil {
		opts.Encryption = &InternalStageEncryption{
			Type: s.Encryption.Type,
		}
	}
	if s.DirectoryTableOptions != nil {
		opts.DirectoryTableOptions = &InternalDirectoryTableOptions{
			Enable:          s.DirectoryTableOptions.Enable,
			RefreshOnCreate: s.DirectoryTableOptions.RefreshOnCreate,
		}
	}
	if s.FileFormat != nil {
		opts.FileFormat = s.FileFormat.toOpts()
	}
	if s.CopyOptions != nil {
		opts.CopyOptions = &StageCopyOptions{
			SizeLimit:         s.CopyOptions.SizeLimit,
			Purge:             s.CopyOptions.Purge,
			ReturnFailedOnly:  s.CopyOptions.ReturnFailedOnly,
			MatchByColumnName: s.CopyOptions.MatchByColumnName,
			EnforceLength:     s.CopyOptions.EnforceLength,
			Truncatecolumns:   s.CopyOptions.Truncatecolumns,
			Force:             s.CopyOptions.Force,
		}
		if s.CopyOptions.OnError != nil {
			opts.CopyOptions.OnError = &StageCopyOnErrorOptions{
				Continue:       s.CopyOptions.OnError.Continue,
				SkipFile:       s.CopyOptions.OnError.SkipFile,
				AbortStatement: s.CopyOptions.OnError.AbortStatement,
			}
		}
	}
	return opts
}

func (r *CreateOnS3StageRequest) toOpts() *CreateOnS3StageOptions {
	opts := &CreateOnS3StageOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		Comment: r.Comment,
		Tag:     r.Tag,
	}
	if r.ExternalStageParams != nil {
		opts.ExternalStageParams = &ExternalS3StageParams{
			Url:                r.ExternalStageParams.Url,
			StorageIntegration: r.ExternalStageParams.StorageIntegration,
		}
		if r.ExternalStageParams.Credentials != nil {
			opts.ExternalStageParams.Credentials = &ExternalStageS3Credentials{
				AWSKeyId:     r.ExternalStageParams.Credentials.AWSKeyId,
				AWSSecretKey: r.ExternalStageParams.Credentials.AWSSecretKey,
				AWSToken:     r.ExternalStageParams.Credentials.AWSToken,
				AWSRole:      r.ExternalStageParams.Credentials.AWSRole,
			}
		}
		if r.ExternalStageParams.Encryption != nil {
			opts.ExternalStageParams.Encryption = &ExternalStageS3Encryption{
				Type:      r.ExternalStageParams.Encryption.Type,
				MasterKey: r.ExternalStageParams.Encryption.MasterKey,
				KmsKeyId:  r.ExternalStageParams.Encryption.KmsKeyId,
			}
		}
	}
	if r.DirectoryTableOptions != nil {
		opts.DirectoryTableOptions = &ExternalS3DirectoryTableOptions{
			Enable:          r.DirectoryTableOptions.Enable,
			RefreshOnCreate: r.DirectoryTableOptions.RefreshOnCreate,
			AutoRefresh:     r.DirectoryTableOptions.AutoRefresh,
		}
	}
	if r.FileFormat != nil {
		opts.FileFormat = &StageFileFormat{
			FormatName: r.FileFormat.FormatName,
			Type:       r.FileFormat.Type,
			Options:    r.FileFormat.Options.toOpts(),
		}
	}
	if r.CopyOptions != nil {
		opts.CopyOptions = &StageCopyOptions{
			SizeLimit:         r.CopyOptions.SizeLimit,
			Purge:             r.CopyOptions.Purge,
			ReturnFailedOnly:  r.CopyOptions.ReturnFailedOnly,
			MatchByColumnName: r.CopyOptions.MatchByColumnName,
			EnforceLength:     r.CopyOptions.EnforceLength,
			Truncatecolumns:   r.CopyOptions.Truncatecolumns,
			Force:             r.CopyOptions.Force,
		}
		if r.CopyOptions.OnError != nil {
			opts.CopyOptions.OnError = &StageCopyOnErrorOptions{
				Continue:       r.CopyOptions.OnError.Continue,
				SkipFile:       r.CopyOptions.OnError.SkipFile,
				AbortStatement: r.CopyOptions.OnError.AbortStatement,
			}
		}
	}
	return opts
}

func (r *CreateOnGCSStageRequest) toOpts() *CreateOnGCSStageOptions {
	opts := &CreateOnGCSStageOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		Comment: r.Comment,
		Tag:     r.Tag,
	}
	if r.ExternalStageParams != nil {
		opts.ExternalStageParams = &ExternalGCSStageParams{
			Url:                r.ExternalStageParams.Url,
			StorageIntegration: r.ExternalStageParams.StorageIntegration,
		}
		if r.ExternalStageParams.Encryption != nil {
			opts.ExternalStageParams.Encryption = &ExternalStageGCSEncryption{
				Type:     r.ExternalStageParams.Encryption.Type,
				KmsKeyId: r.ExternalStageParams.Encryption.KmsKeyId,
			}
		}
	}
	if r.DirectoryTableOptions != nil {
		opts.DirectoryTableOptions = &ExternalGCSDirectoryTableOptions{
			Enable:                  r.DirectoryTableOptions.Enable,
			RefreshOnCreate:         r.DirectoryTableOptions.RefreshOnCreate,
			AutoRefresh:             r.DirectoryTableOptions.AutoRefresh,
			NotificationIntegration: r.DirectoryTableOptions.NotificationIntegration,
		}
	}
	if r.FileFormat != nil {
		opts.FileFormat = &StageFileFormat{
			FormatName: r.FileFormat.FormatName,
			Type:       r.FileFormat.Type,
			Options:    r.FileFormat.Options.toOpts(),
		}
	}
	if r.CopyOptions != nil {
		opts.CopyOptions = &StageCopyOptions{
			SizeLimit:         r.CopyOptions.SizeLimit,
			Purge:             r.CopyOptions.Purge,
			ReturnFailedOnly:  r.CopyOptions.ReturnFailedOnly,
			MatchByColumnName: r.CopyOptions.MatchByColumnName,
			EnforceLength:     r.CopyOptions.EnforceLength,
			Truncatecolumns:   r.CopyOptions.Truncatecolumns,
			Force:             r.CopyOptions.Force,
		}
		if r.CopyOptions.OnError != nil {
			opts.CopyOptions.OnError = &StageCopyOnErrorOptions{
				Continue:       r.CopyOptions.OnError.Continue,
				SkipFile:       r.CopyOptions.OnError.SkipFile,
				AbortStatement: r.CopyOptions.OnError.AbortStatement,
			}
		}
	}
	return opts
}

func (r *CreateOnAzureStageRequest) toOpts() *CreateOnAzureStageOptions {
	opts := &CreateOnAzureStageOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		Comment: r.Comment,
		Tag:     r.Tag,
	}
	if r.ExternalStageParams != nil {
		opts.ExternalStageParams = &ExternalAzureStageParams{
			Url:                r.ExternalStageParams.Url,
			StorageIntegration: r.ExternalStageParams.StorageIntegration,
		}
		if r.ExternalStageParams.Credentials != nil {
			opts.ExternalStageParams.Credentials = &ExternalStageAzureCredentials{
				AzureSasToken: r.ExternalStageParams.Credentials.AzureSasToken,
			}
		}
		if r.ExternalStageParams.Encryption != nil {
			opts.ExternalStageParams.Encryption = &ExternalStageAzureEncryption{
				Type:      r.ExternalStageParams.Encryption.Type,
				MasterKey: r.ExternalStageParams.Encryption.MasterKey,
			}
		}
	}
	if r.DirectoryTableOptions != nil {
		opts.DirectoryTableOptions = &ExternalAzureDirectoryTableOptions{
			Enable:                  r.DirectoryTableOptions.Enable,
			RefreshOnCreate:         r.DirectoryTableOptions.RefreshOnCreate,
			AutoRefresh:             r.DirectoryTableOptions.AutoRefresh,
			NotificationIntegration: r.DirectoryTableOptions.NotificationIntegration,
		}
	}
	if r.FileFormat != nil {
		opts.FileFormat = &StageFileFormat{
			FormatName: r.FileFormat.FormatName,
			Type:       r.FileFormat.Type,
			Options:    r.FileFormat.Options.toOpts(),
		}
	}
	if r.CopyOptions != nil {
		opts.CopyOptions = &StageCopyOptions{
			SizeLimit:         r.CopyOptions.SizeLimit,
			Purge:             r.CopyOptions.Purge,
			ReturnFailedOnly:  r.CopyOptions.ReturnFailedOnly,
			MatchByColumnName: r.CopyOptions.MatchByColumnName,
			EnforceLength:     r.CopyOptions.EnforceLength,
			Truncatecolumns:   r.CopyOptions.Truncatecolumns,
			Force:             r.CopyOptions.Force,
		}
		if r.CopyOptions.OnError != nil {
			opts.CopyOptions.OnError = &StageCopyOnErrorOptions{
				Continue:       r.CopyOptions.OnError.Continue,
				SkipFile:       r.CopyOptions.OnError.SkipFile,
				AbortStatement: r.CopyOptions.OnError.AbortStatement,
			}
		}
	}
	return opts
}

func (r *CreateOnS3CompatibleStageRequest) toOpts() *CreateOnS3CompatibleStageOptions {
	opts := &CreateOnS3CompatibleStageOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		Url:         r.Url,
		Endpoint:    r.Endpoint,

		Comment: r.Comment,
		Tag:     r.Tag,
	}
	if r.Credentials != nil {
		opts.Credentials = &ExternalStageS3CompatibleCredentials{
			AWSKeyId:     r.Credentials.AWSKeyId,
			AWSSecretKey: r.Credentials.AWSSecretKey,
		}
	}
	if r.DirectoryTableOptions != nil {
		opts.DirectoryTableOptions = &ExternalS3DirectoryTableOptions{
			Enable:          r.DirectoryTableOptions.Enable,
			RefreshOnCreate: r.DirectoryTableOptions.RefreshOnCreate,
			AutoRefresh:     r.DirectoryTableOptions.AutoRefresh,
		}
	}
	if r.FileFormat != nil {
		opts.FileFormat = &StageFileFormat{
			FormatName: r.FileFormat.FormatName,
			Type:       r.FileFormat.Type,
			Options:    r.FileFormat.Options.toOpts(),
		}
	}
	if r.CopyOptions != nil {
		opts.CopyOptions = &StageCopyOptions{
			SizeLimit:         r.CopyOptions.SizeLimit,
			Purge:             r.CopyOptions.Purge,
			ReturnFailedOnly:  r.CopyOptions.ReturnFailedOnly,
			MatchByColumnName: r.CopyOptions.MatchByColumnName,
			EnforceLength:     r.CopyOptions.EnforceLength,
			Truncatecolumns:   r.CopyOptions.Truncatecolumns,
			Force:             r.CopyOptions.Force,
		}
		if r.CopyOptions.OnError != nil {
			opts.CopyOptions.OnError = &StageCopyOnErrorOptions{
				Continue:       r.CopyOptions.OnError.Continue,
				SkipFile:       r.CopyOptions.OnError.SkipFile,
				AbortStatement: r.CopyOptions.OnError.AbortStatement,
			}
		}
	}
	return opts
}

func (r *AlterStageRequest) toOpts() *AlterStageOptions {
	opts := &AlterStageOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		RenameTo:  r.RenameTo,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	return opts
}

func (r *AlterInternalStageStageRequest) toOpts() *AlterInternalStageStageOptions {
	opts := &AlterInternalStageStageOptions{
		IfExists: r.IfExists,
		name:     r.name,

		Comment: r.Comment,
	}
	if r.FileFormat != nil {
		opts.FileFormat = &StageFileFormat{
			FormatName: r.FileFormat.FormatName,
			Type:       r.FileFormat.Type,
			Options:    r.FileFormat.Options.toOpts(),
		}
	}
	if r.CopyOptions != nil {
		opts.CopyOptions = &StageCopyOptions{
			SizeLimit:         r.CopyOptions.SizeLimit,
			Purge:             r.CopyOptions.Purge,
			ReturnFailedOnly:  r.CopyOptions.ReturnFailedOnly,
			MatchByColumnName: r.CopyOptions.MatchByColumnName,
			EnforceLength:     r.CopyOptions.EnforceLength,
			Truncatecolumns:   r.CopyOptions.Truncatecolumns,
			Force:             r.CopyOptions.Force,
		}
		if r.CopyOptions.OnError != nil {
			opts.CopyOptions.OnError = &StageCopyOnErrorOptions{
				Continue:       r.CopyOptions.OnError.Continue,
				SkipFile:       r.CopyOptions.OnError.SkipFile,
				AbortStatement: r.CopyOptions.OnError.AbortStatement,
			}
		}
	}
	return opts
}

func (r *AlterExternalS3StageStageRequest) toOpts() *AlterExternalS3StageStageOptions {
	opts := &AlterExternalS3StageStageOptions{
		IfExists: r.IfExists,
		name:     r.name,

		Comment: r.Comment,
	}
	if r.ExternalStageParams != nil {
		opts.ExternalStageParams = &ExternalS3StageParams{
			Url:                r.ExternalStageParams.Url,
			StorageIntegration: r.ExternalStageParams.StorageIntegration,
		}
		if r.ExternalStageParams.Credentials != nil {
			opts.ExternalStageParams.Credentials = &ExternalStageS3Credentials{
				AWSKeyId:     r.ExternalStageParams.Credentials.AWSKeyId,
				AWSSecretKey: r.ExternalStageParams.Credentials.AWSSecretKey,
				AWSToken:     r.ExternalStageParams.Credentials.AWSToken,
				AWSRole:      r.ExternalStageParams.Credentials.AWSRole,
			}
		}
		if r.ExternalStageParams.Encryption != nil {
			opts.ExternalStageParams.Encryption = &ExternalStageS3Encryption{
				Type:      r.ExternalStageParams.Encryption.Type,
				MasterKey: r.ExternalStageParams.Encryption.MasterKey,
				KmsKeyId:  r.ExternalStageParams.Encryption.KmsKeyId,
			}
		}
	}
	if r.FileFormat != nil {
		opts.FileFormat = &StageFileFormat{
			FormatName: r.FileFormat.FormatName,
			Type:       r.FileFormat.Type,
			Options:    r.FileFormat.Options.toOpts(),
		}
	}
	if r.CopyOptions != nil {
		opts.CopyOptions = &StageCopyOptions{
			SizeLimit:         r.CopyOptions.SizeLimit,
			Purge:             r.CopyOptions.Purge,
			ReturnFailedOnly:  r.CopyOptions.ReturnFailedOnly,
			MatchByColumnName: r.CopyOptions.MatchByColumnName,
			EnforceLength:     r.CopyOptions.EnforceLength,
			Truncatecolumns:   r.CopyOptions.Truncatecolumns,
			Force:             r.CopyOptions.Force,
		}
		if r.CopyOptions.OnError != nil {
			opts.CopyOptions.OnError = &StageCopyOnErrorOptions{
				Continue:       r.CopyOptions.OnError.Continue,
				SkipFile:       r.CopyOptions.OnError.SkipFile,
				AbortStatement: r.CopyOptions.OnError.AbortStatement,
			}
		}
	}
	return opts
}

func (r *AlterExternalGCSStageStageRequest) toOpts() *AlterExternalGCSStageStageOptions {
	opts := &AlterExternalGCSStageStageOptions{
		IfExists: r.IfExists,
		name:     r.name,

		Comment: r.Comment,
	}
	if r.ExternalStageParams != nil {
		opts.ExternalStageParams = &ExternalGCSStageParams{
			Url:                r.ExternalStageParams.Url,
			StorageIntegration: r.ExternalStageParams.StorageIntegration,
		}
		if r.ExternalStageParams.Encryption != nil {
			opts.ExternalStageParams.Encryption = &ExternalStageGCSEncryption{
				Type:     r.ExternalStageParams.Encryption.Type,
				KmsKeyId: r.ExternalStageParams.Encryption.KmsKeyId,
			}
		}
	}
	if r.FileFormat != nil {
		opts.FileFormat = &StageFileFormat{
			FormatName: r.FileFormat.FormatName,
			Type:       r.FileFormat.Type,
			Options:    r.FileFormat.Options.toOpts(),
		}
	}
	if r.CopyOptions != nil {
		opts.CopyOptions = &StageCopyOptions{
			SizeLimit:         r.CopyOptions.SizeLimit,
			Purge:             r.CopyOptions.Purge,
			ReturnFailedOnly:  r.CopyOptions.ReturnFailedOnly,
			MatchByColumnName: r.CopyOptions.MatchByColumnName,
			EnforceLength:     r.CopyOptions.EnforceLength,
			Truncatecolumns:   r.CopyOptions.Truncatecolumns,
			Force:             r.CopyOptions.Force,
		}
		if r.CopyOptions.OnError != nil {
			opts.CopyOptions.OnError = &StageCopyOnErrorOptions{
				Continue:       r.CopyOptions.OnError.Continue,
				SkipFile:       r.CopyOptions.OnError.SkipFile,
				AbortStatement: r.CopyOptions.OnError.AbortStatement,
			}
		}
	}
	return opts
}

func (r *AlterExternalAzureStageStageRequest) toOpts() *AlterExternalAzureStageStageOptions {
	opts := &AlterExternalAzureStageStageOptions{
		IfExists: r.IfExists,
		name:     r.name,

		Comment: r.Comment,
	}
	if r.ExternalStageParams != nil {
		opts.ExternalStageParams = &ExternalAzureStageParams{
			Url:                r.ExternalStageParams.Url,
			StorageIntegration: r.ExternalStageParams.StorageIntegration,
		}
		if r.ExternalStageParams.Credentials != nil {
			opts.ExternalStageParams.Credentials = &ExternalStageAzureCredentials{
				AzureSasToken: r.ExternalStageParams.Credentials.AzureSasToken,
			}
		}
		if r.ExternalStageParams.Encryption != nil {
			opts.ExternalStageParams.Encryption = &ExternalStageAzureEncryption{
				Type:      r.ExternalStageParams.Encryption.Type,
				MasterKey: r.ExternalStageParams.Encryption.MasterKey,
			}
		}
	}
	if r.FileFormat != nil {
		opts.FileFormat = &StageFileFormat{
			FormatName: r.FileFormat.FormatName,
			Type:       r.FileFormat.Type,
			Options:    r.FileFormat.Options.toOpts(),
		}
	}
	if r.CopyOptions != nil {
		opts.CopyOptions = &StageCopyOptions{
			SizeLimit:         r.CopyOptions.SizeLimit,
			Purge:             r.CopyOptions.Purge,
			ReturnFailedOnly:  r.CopyOptions.ReturnFailedOnly,
			MatchByColumnName: r.CopyOptions.MatchByColumnName,
			EnforceLength:     r.CopyOptions.EnforceLength,
			Truncatecolumns:   r.CopyOptions.Truncatecolumns,
			Force:             r.CopyOptions.Force,
		}
		if r.CopyOptions.OnError != nil {
			opts.CopyOptions.OnError = &StageCopyOnErrorOptions{
				Continue:       r.CopyOptions.OnError.Continue,
				SkipFile:       r.CopyOptions.OnError.SkipFile,
				AbortStatement: r.CopyOptions.OnError.AbortStatement,
			}
		}
	}
	return opts
}

func (r *AlterDirectoryTableStageRequest) toOpts() *AlterDirectoryTableStageOptions {
	opts := &AlterDirectoryTableStageOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	if r.SetDirectory != nil {
		opts.SetDirectory = &DirectoryTableSet{
			Enable: r.SetDirectory.Enable,
		}
	}
	if r.Refresh != nil {
		opts.Refresh = &DirectoryTableRefresh{
			Subpath: r.Refresh.Subpath,
		}
	}
	return opts
}

func (r *DropStageRequest) toOpts() *DropStageOptions {
	opts := &DropStageOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *DescribeStageRequest) toOpts() *DescribeStageOptions {
	opts := &DescribeStageOptions{
		name: r.name,
	}
	return opts
}

func (r stageDescRow) convert() *StageProperty {
	stageProp := &StageProperty{
		Parent:  r.ParentProperty,
		Name:    r.Property,
		Type:    r.PropertyType,
		Value:   r.PropertyValue,
		Default: r.PropertyDefault,
	}
	return stageProp
}

func (r *ShowStageRequest) toOpts() *ShowStageOptions {
	opts := &ShowStageOptions{
		Like: r.Like,
		In:   r.In,
	}
	return opts
}

func (r stageShowRow) convert() *Stage {
	stage := &Stage{
		CreatedOn:        r.CreatedOn,
		Name:             r.Name,
		DatabaseName:     r.DatabaseName,
		SchemaName:       r.SchemaName,
		Url:              r.Url,
		HasCredentials:   r.HasCredentials == "Y",
		HasEncryptionKey: r.HasEncryptionKey == "Y",
		Owner:            r.Owner,
		Comment:          r.Comment,
		Type:             r.Type,
		DirectoryEnabled: r.DirectoryEnabled == "Y",
	}
	if r.Region.Valid {
		stage.Region = &r.Region.String
	}
	if r.Cloud.Valid {
		stage.Cloud = &r.Cloud.String
	}
	if r.StorageIntegration.Valid {
		stage.StorageIntegration = &r.StorageIntegration.String
	}
	if r.Endpoint.Valid {
		stage.Endpoint = &r.Endpoint.String
	}
	if r.OwnerRoleType.Valid {
		stage.OwnerRoleType = &r.OwnerRoleType.String
	}
	return stage
}

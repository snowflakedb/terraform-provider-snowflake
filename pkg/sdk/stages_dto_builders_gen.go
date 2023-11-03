// Code generated by dto builder generator; DO NOT EDIT.

package sdk

import ()

func NewCreateInternalStageRequest(
	name SchemaObjectIdentifier,
) *CreateInternalStageRequest {
	s := CreateInternalStageRequest{}
	s.name = name
	return &s
}

func (s *CreateInternalStageRequest) WithOrReplace(OrReplace *bool) *CreateInternalStageRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateInternalStageRequest) WithTemporary(Temporary *bool) *CreateInternalStageRequest {
	s.Temporary = Temporary
	return s
}

func (s *CreateInternalStageRequest) WithIfNotExists(IfNotExists *bool) *CreateInternalStageRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateInternalStageRequest) WithEncryption(Encryption *InternalStageEncryptionRequest) *CreateInternalStageRequest {
	s.Encryption = Encryption
	return s
}

func (s *CreateInternalStageRequest) WithDirectoryTableOptions(DirectoryTableOptions *InternalDirectoryTableOptionsRequest) *CreateInternalStageRequest {
	s.DirectoryTableOptions = DirectoryTableOptions
	return s
}

func (s *CreateInternalStageRequest) WithFileFormat(FileFormat *StageFileFormatRequest) *CreateInternalStageRequest {
	s.FileFormat = FileFormat
	return s
}

func (s *CreateInternalStageRequest) WithCopyOptions(CopyOptions *StageCopyOptionsRequest) *CreateInternalStageRequest {
	s.CopyOptions = CopyOptions
	return s
}

func (s *CreateInternalStageRequest) WithComment(Comment *string) *CreateInternalStageRequest {
	s.Comment = Comment
	return s
}

func (s *CreateInternalStageRequest) WithTag(Tag []TagAssociation) *CreateInternalStageRequest {
	s.Tag = Tag
	return s
}

func NewInternalStageEncryptionRequest(
	Type *InternalStageEncryptionOption,
) *InternalStageEncryptionRequest {
	s := InternalStageEncryptionRequest{}
	s.Type = Type
	return &s
}

func NewInternalDirectoryTableOptionsRequest() *InternalDirectoryTableOptionsRequest {
	return &InternalDirectoryTableOptionsRequest{}
}

func (s *InternalDirectoryTableOptionsRequest) WithEnable(Enable *bool) *InternalDirectoryTableOptionsRequest {
	s.Enable = Enable
	return s
}

func (s *InternalDirectoryTableOptionsRequest) WithRefreshOnCreate(RefreshOnCreate *bool) *InternalDirectoryTableOptionsRequest {
	s.RefreshOnCreate = RefreshOnCreate
	return s
}

func NewStageFileFormatRequest() *StageFileFormatRequest {
	return &StageFileFormatRequest{}
}

func (s *StageFileFormatRequest) WithFormatName(FormatName *string) *StageFileFormatRequest {
	s.FormatName = FormatName
	return s
}

func (s *StageFileFormatRequest) WithType(Type *FileFormatType) *StageFileFormatRequest {
	s.Type = Type
	return s
}

func NewStageCopyOptionsRequest() *StageCopyOptionsRequest {
	return &StageCopyOptionsRequest{}
}

func (s *StageCopyOptionsRequest) WithOnError(OnError *StageCopyOnErrorOptionsRequest) *StageCopyOptionsRequest {
	s.OnError = OnError
	return s
}

func (s *StageCopyOptionsRequest) WithSizeLimit(SizeLimit *int) *StageCopyOptionsRequest {
	s.SizeLimit = SizeLimit
	return s
}

func (s *StageCopyOptionsRequest) WithPurge(Purge *bool) *StageCopyOptionsRequest {
	s.Purge = Purge
	return s
}

func (s *StageCopyOptionsRequest) WithReturnFailedOnly(ReturnFailedOnly *bool) *StageCopyOptionsRequest {
	s.ReturnFailedOnly = ReturnFailedOnly
	return s
}

func (s *StageCopyOptionsRequest) WithMatchByColumnName(MatchByColumnName *StageCopyColumnMapOption) *StageCopyOptionsRequest {
	s.MatchByColumnName = MatchByColumnName
	return s
}

func (s *StageCopyOptionsRequest) WithEnforceLength(EnforceLength *bool) *StageCopyOptionsRequest {
	s.EnforceLength = EnforceLength
	return s
}

func (s *StageCopyOptionsRequest) WithTruncatecolumns(Truncatecolumns *bool) *StageCopyOptionsRequest {
	s.Truncatecolumns = Truncatecolumns
	return s
}

func (s *StageCopyOptionsRequest) WithForce(Force *bool) *StageCopyOptionsRequest {
	s.Force = Force
	return s
}

func NewStageCopyOnErrorOptionsRequest() *StageCopyOnErrorOptionsRequest {
	return &StageCopyOnErrorOptionsRequest{}
}

func (s *StageCopyOnErrorOptionsRequest) WithContinue(Continue *bool) *StageCopyOnErrorOptionsRequest {
	s.Continue = Continue
	return s
}

func (s *StageCopyOnErrorOptionsRequest) WithSkipFile(SkipFile *bool) *StageCopyOnErrorOptionsRequest {
	s.SkipFile = SkipFile
	return s
}

func (s *StageCopyOnErrorOptionsRequest) WithAbortStatement(AbortStatement *bool) *StageCopyOnErrorOptionsRequest {
	s.AbortStatement = AbortStatement
	return s
}

func NewCreateOnS3StageRequest(
	name SchemaObjectIdentifier,
) *CreateOnS3StageRequest {
	s := CreateOnS3StageRequest{}
	s.name = name
	return &s
}

func (s *CreateOnS3StageRequest) WithOrReplace(OrReplace *bool) *CreateOnS3StageRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateOnS3StageRequest) WithTemporary(Temporary *bool) *CreateOnS3StageRequest {
	s.Temporary = Temporary
	return s
}

func (s *CreateOnS3StageRequest) WithIfNotExists(IfNotExists *bool) *CreateOnS3StageRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateOnS3StageRequest) WithExternalStageParams(ExternalStageParams *ExternalS3StageParamsRequest) *CreateOnS3StageRequest {
	s.ExternalStageParams = ExternalStageParams
	return s
}

func (s *CreateOnS3StageRequest) WithDirectoryTableOptions(DirectoryTableOptions *ExternalS3DirectoryTableOptionsRequest) *CreateOnS3StageRequest {
	s.DirectoryTableOptions = DirectoryTableOptions
	return s
}

func (s *CreateOnS3StageRequest) WithFileFormat(FileFormat *StageFileFormatRequest) *CreateOnS3StageRequest {
	s.FileFormat = FileFormat
	return s
}

func (s *CreateOnS3StageRequest) WithCopyOptions(CopyOptions *StageCopyOptionsRequest) *CreateOnS3StageRequest {
	s.CopyOptions = CopyOptions
	return s
}

func (s *CreateOnS3StageRequest) WithComment(Comment *string) *CreateOnS3StageRequest {
	s.Comment = Comment
	return s
}

func (s *CreateOnS3StageRequest) WithTag(Tag []TagAssociation) *CreateOnS3StageRequest {
	s.Tag = Tag
	return s
}

func NewExternalS3StageParamsRequest(
	Url string,
) *ExternalS3StageParamsRequest {
	s := ExternalS3StageParamsRequest{}
	s.Url = Url
	return &s
}

func (s *ExternalS3StageParamsRequest) WithStorageIntegration(StorageIntegration *AccountObjectIdentifier) *ExternalS3StageParamsRequest {
	s.StorageIntegration = StorageIntegration
	return s
}

func (s *ExternalS3StageParamsRequest) WithCredentials(Credentials *ExternalStageS3CredentialsRequest) *ExternalS3StageParamsRequest {
	s.Credentials = Credentials
	return s
}

func (s *ExternalS3StageParamsRequest) WithEncryption(Encryption *ExternalStageS3EncryptionRequest) *ExternalS3StageParamsRequest {
	s.Encryption = Encryption
	return s
}

func NewExternalStageS3CredentialsRequest() *ExternalStageS3CredentialsRequest {
	return &ExternalStageS3CredentialsRequest{}
}

func (s *ExternalStageS3CredentialsRequest) WithAwsKeyId(AwsKeyId *string) *ExternalStageS3CredentialsRequest {
	s.AwsKeyId = AwsKeyId
	return s
}

func (s *ExternalStageS3CredentialsRequest) WithAwsSecretKey(AwsSecretKey *string) *ExternalStageS3CredentialsRequest {
	s.AwsSecretKey = AwsSecretKey
	return s
}

func (s *ExternalStageS3CredentialsRequest) WithAwsToken(AwsToken *string) *ExternalStageS3CredentialsRequest {
	s.AwsToken = AwsToken
	return s
}

func (s *ExternalStageS3CredentialsRequest) WithAwsRole(AwsRole *string) *ExternalStageS3CredentialsRequest {
	s.AwsRole = AwsRole
	return s
}

func NewExternalStageS3EncryptionRequest(
	Type *ExternalStageS3EncryptionOption,
) *ExternalStageS3EncryptionRequest {
	s := ExternalStageS3EncryptionRequest{}
	s.Type = Type
	return &s
}

func (s *ExternalStageS3EncryptionRequest) WithMasterKey(MasterKey *string) *ExternalStageS3EncryptionRequest {
	s.MasterKey = MasterKey
	return s
}

func (s *ExternalStageS3EncryptionRequest) WithKmsKeyId(KmsKeyId *string) *ExternalStageS3EncryptionRequest {
	s.KmsKeyId = KmsKeyId
	return s
}

func NewExternalS3DirectoryTableOptionsRequest() *ExternalS3DirectoryTableOptionsRequest {
	return &ExternalS3DirectoryTableOptionsRequest{}
}

func (s *ExternalS3DirectoryTableOptionsRequest) WithEnable(Enable *bool) *ExternalS3DirectoryTableOptionsRequest {
	s.Enable = Enable
	return s
}

func (s *ExternalS3DirectoryTableOptionsRequest) WithRefreshOnCreate(RefreshOnCreate *bool) *ExternalS3DirectoryTableOptionsRequest {
	s.RefreshOnCreate = RefreshOnCreate
	return s
}

func (s *ExternalS3DirectoryTableOptionsRequest) WithAutoRefresh(AutoRefresh *bool) *ExternalS3DirectoryTableOptionsRequest {
	s.AutoRefresh = AutoRefresh
	return s
}

func NewCreateOnGCSStageRequest(
	name SchemaObjectIdentifier,
) *CreateOnGCSStageRequest {
	s := CreateOnGCSStageRequest{}
	s.name = name
	return &s
}

func (s *CreateOnGCSStageRequest) WithOrReplace(OrReplace *bool) *CreateOnGCSStageRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateOnGCSStageRequest) WithTemporary(Temporary *bool) *CreateOnGCSStageRequest {
	s.Temporary = Temporary
	return s
}

func (s *CreateOnGCSStageRequest) WithIfNotExists(IfNotExists *bool) *CreateOnGCSStageRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateOnGCSStageRequest) WithExternalStageParams(ExternalStageParams *ExternalGCSStageParamsRequest) *CreateOnGCSStageRequest {
	s.ExternalStageParams = ExternalStageParams
	return s
}

func (s *CreateOnGCSStageRequest) WithDirectoryTableOptions(DirectoryTableOptions *ExternalGCSDirectoryTableOptionsRequest) *CreateOnGCSStageRequest {
	s.DirectoryTableOptions = DirectoryTableOptions
	return s
}

func (s *CreateOnGCSStageRequest) WithFileFormat(FileFormat *StageFileFormatRequest) *CreateOnGCSStageRequest {
	s.FileFormat = FileFormat
	return s
}

func (s *CreateOnGCSStageRequest) WithCopyOptions(CopyOptions *StageCopyOptionsRequest) *CreateOnGCSStageRequest {
	s.CopyOptions = CopyOptions
	return s
}

func (s *CreateOnGCSStageRequest) WithComment(Comment *string) *CreateOnGCSStageRequest {
	s.Comment = Comment
	return s
}

func (s *CreateOnGCSStageRequest) WithTag(Tag []TagAssociation) *CreateOnGCSStageRequest {
	s.Tag = Tag
	return s
}

func NewExternalGCSStageParamsRequest(
	Url string,
) *ExternalGCSStageParamsRequest {
	s := ExternalGCSStageParamsRequest{}
	s.Url = Url
	return &s
}

func (s *ExternalGCSStageParamsRequest) WithStorageIntegration(StorageIntegration *AccountObjectIdentifier) *ExternalGCSStageParamsRequest {
	s.StorageIntegration = StorageIntegration
	return s
}

func (s *ExternalGCSStageParamsRequest) WithEncryption(Encryption *ExternalStageGCSEncryptionRequest) *ExternalGCSStageParamsRequest {
	s.Encryption = Encryption
	return s
}

func NewExternalStageGCSEncryptionRequest(
	Type *ExternalStageGCSEncryptionOption,
) *ExternalStageGCSEncryptionRequest {
	s := ExternalStageGCSEncryptionRequest{}
	s.Type = Type
	return &s
}

func (s *ExternalStageGCSEncryptionRequest) WithKmsKeyId(KmsKeyId *string) *ExternalStageGCSEncryptionRequest {
	s.KmsKeyId = KmsKeyId
	return s
}

func NewExternalGCSDirectoryTableOptionsRequest() *ExternalGCSDirectoryTableOptionsRequest {
	return &ExternalGCSDirectoryTableOptionsRequest{}
}

func (s *ExternalGCSDirectoryTableOptionsRequest) WithEnable(Enable *bool) *ExternalGCSDirectoryTableOptionsRequest {
	s.Enable = Enable
	return s
}

func (s *ExternalGCSDirectoryTableOptionsRequest) WithRefreshOnCreate(RefreshOnCreate *bool) *ExternalGCSDirectoryTableOptionsRequest {
	s.RefreshOnCreate = RefreshOnCreate
	return s
}

func (s *ExternalGCSDirectoryTableOptionsRequest) WithAutoRefresh(AutoRefresh *bool) *ExternalGCSDirectoryTableOptionsRequest {
	s.AutoRefresh = AutoRefresh
	return s
}

func (s *ExternalGCSDirectoryTableOptionsRequest) WithNotificationIntegration(NotificationIntegration *string) *ExternalGCSDirectoryTableOptionsRequest {
	s.NotificationIntegration = NotificationIntegration
	return s
}

func NewCreateOnAzureStageRequest(
	name SchemaObjectIdentifier,
) *CreateOnAzureStageRequest {
	s := CreateOnAzureStageRequest{}
	s.name = name
	return &s
}

func (s *CreateOnAzureStageRequest) WithOrReplace(OrReplace *bool) *CreateOnAzureStageRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateOnAzureStageRequest) WithTemporary(Temporary *bool) *CreateOnAzureStageRequest {
	s.Temporary = Temporary
	return s
}

func (s *CreateOnAzureStageRequest) WithIfNotExists(IfNotExists *bool) *CreateOnAzureStageRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateOnAzureStageRequest) WithExternalStageParams(ExternalStageParams *ExternalAzureStageParamsRequest) *CreateOnAzureStageRequest {
	s.ExternalStageParams = ExternalStageParams
	return s
}

func (s *CreateOnAzureStageRequest) WithDirectoryTableOptions(DirectoryTableOptions *ExternalAzureDirectoryTableOptionsRequest) *CreateOnAzureStageRequest {
	s.DirectoryTableOptions = DirectoryTableOptions
	return s
}

func (s *CreateOnAzureStageRequest) WithFileFormat(FileFormat *StageFileFormatRequest) *CreateOnAzureStageRequest {
	s.FileFormat = FileFormat
	return s
}

func (s *CreateOnAzureStageRequest) WithCopyOptions(CopyOptions *StageCopyOptionsRequest) *CreateOnAzureStageRequest {
	s.CopyOptions = CopyOptions
	return s
}

func (s *CreateOnAzureStageRequest) WithComment(Comment *string) *CreateOnAzureStageRequest {
	s.Comment = Comment
	return s
}

func (s *CreateOnAzureStageRequest) WithTag(Tag []TagAssociation) *CreateOnAzureStageRequest {
	s.Tag = Tag
	return s
}

func NewExternalAzureStageParamsRequest(
	Url string,
) *ExternalAzureStageParamsRequest {
	s := ExternalAzureStageParamsRequest{}
	s.Url = Url
	return &s
}

func (s *ExternalAzureStageParamsRequest) WithStorageIntegration(StorageIntegration *AccountObjectIdentifier) *ExternalAzureStageParamsRequest {
	s.StorageIntegration = StorageIntegration
	return s
}

func (s *ExternalAzureStageParamsRequest) WithCredentials(Credentials *ExternalStageAzureCredentialsRequest) *ExternalAzureStageParamsRequest {
	s.Credentials = Credentials
	return s
}

func (s *ExternalAzureStageParamsRequest) WithEncryption(Encryption *ExternalStageAzureEncryptionRequest) *ExternalAzureStageParamsRequest {
	s.Encryption = Encryption
	return s
}

func NewExternalStageAzureCredentialsRequest(
	AzureSasToken string,
) *ExternalStageAzureCredentialsRequest {
	s := ExternalStageAzureCredentialsRequest{}
	s.AzureSasToken = AzureSasToken
	return &s
}

func NewExternalStageAzureEncryptionRequest(
	Type *ExternalStageAzureEncryptionOption,
) *ExternalStageAzureEncryptionRequest {
	s := ExternalStageAzureEncryptionRequest{}
	s.Type = Type
	return &s
}

func (s *ExternalStageAzureEncryptionRequest) WithMasterKey(MasterKey *string) *ExternalStageAzureEncryptionRequest {
	s.MasterKey = MasterKey
	return s
}

func NewExternalAzureDirectoryTableOptionsRequest() *ExternalAzureDirectoryTableOptionsRequest {
	return &ExternalAzureDirectoryTableOptionsRequest{}
}

func (s *ExternalAzureDirectoryTableOptionsRequest) WithEnable(Enable *bool) *ExternalAzureDirectoryTableOptionsRequest {
	s.Enable = Enable
	return s
}

func (s *ExternalAzureDirectoryTableOptionsRequest) WithRefreshOnCreate(RefreshOnCreate *bool) *ExternalAzureDirectoryTableOptionsRequest {
	s.RefreshOnCreate = RefreshOnCreate
	return s
}

func (s *ExternalAzureDirectoryTableOptionsRequest) WithAutoRefresh(AutoRefresh *bool) *ExternalAzureDirectoryTableOptionsRequest {
	s.AutoRefresh = AutoRefresh
	return s
}

func (s *ExternalAzureDirectoryTableOptionsRequest) WithNotificationIntegration(NotificationIntegration *string) *ExternalAzureDirectoryTableOptionsRequest {
	s.NotificationIntegration = NotificationIntegration
	return s
}

func NewCreateOnS3CompatibleStageRequest(
	name SchemaObjectIdentifier,
	Url string,
	Endpoint string,
) *CreateOnS3CompatibleStageRequest {
	s := CreateOnS3CompatibleStageRequest{}
	s.name = name
	s.Url = Url
	s.Endpoint = Endpoint
	return &s
}

func (s *CreateOnS3CompatibleStageRequest) WithOrReplace(OrReplace *bool) *CreateOnS3CompatibleStageRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateOnS3CompatibleStageRequest) WithTemporary(Temporary *bool) *CreateOnS3CompatibleStageRequest {
	s.Temporary = Temporary
	return s
}

func (s *CreateOnS3CompatibleStageRequest) WithIfNotExists(IfNotExists *bool) *CreateOnS3CompatibleStageRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateOnS3CompatibleStageRequest) WithCredentials(Credentials *ExternalStageS3CompatibleCredentialsRequest) *CreateOnS3CompatibleStageRequest {
	s.Credentials = Credentials
	return s
}

func (s *CreateOnS3CompatibleStageRequest) WithDirectoryTableOptions(DirectoryTableOptions *ExternalS3DirectoryTableOptionsRequest) *CreateOnS3CompatibleStageRequest {
	s.DirectoryTableOptions = DirectoryTableOptions
	return s
}

func (s *CreateOnS3CompatibleStageRequest) WithFileFormat(FileFormat *StageFileFormatRequest) *CreateOnS3CompatibleStageRequest {
	s.FileFormat = FileFormat
	return s
}

func (s *CreateOnS3CompatibleStageRequest) WithCopyOptions(CopyOptions *StageCopyOptionsRequest) *CreateOnS3CompatibleStageRequest {
	s.CopyOptions = CopyOptions
	return s
}

func (s *CreateOnS3CompatibleStageRequest) WithComment(Comment *string) *CreateOnS3CompatibleStageRequest {
	s.Comment = Comment
	return s
}

func (s *CreateOnS3CompatibleStageRequest) WithTag(Tag []TagAssociation) *CreateOnS3CompatibleStageRequest {
	s.Tag = Tag
	return s
}

func NewExternalStageS3CompatibleCredentialsRequest(
	AwsKeyId *string,
	AwsSecretKey *string,
) *ExternalStageS3CompatibleCredentialsRequest {
	s := ExternalStageS3CompatibleCredentialsRequest{}
	s.AwsKeyId = AwsKeyId
	s.AwsSecretKey = AwsSecretKey
	return &s
}

func NewAlterStageRequest(
	name SchemaObjectIdentifier,
) *AlterStageRequest {
	s := AlterStageRequest{}
	s.name = name
	return &s
}

func (s *AlterStageRequest) WithIfExists(IfExists *bool) *AlterStageRequest {
	s.IfExists = IfExists
	return s
}

func (s *AlterStageRequest) WithRenameTo(RenameTo *SchemaObjectIdentifier) *AlterStageRequest {
	s.RenameTo = RenameTo
	return s
}

func (s *AlterStageRequest) WithSetTags(SetTags []TagAssociation) *AlterStageRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterStageRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterStageRequest {
	s.UnsetTags = UnsetTags
	return s
}

func NewAlterInternalStageStageRequest(
	name SchemaObjectIdentifier,
) *AlterInternalStageStageRequest {
	s := AlterInternalStageStageRequest{}
	s.name = name
	return &s
}

func (s *AlterInternalStageStageRequest) WithIfExists(IfExists *bool) *AlterInternalStageStageRequest {
	s.IfExists = IfExists
	return s
}

func (s *AlterInternalStageStageRequest) WithFileFormat(FileFormat *StageFileFormatRequest) *AlterInternalStageStageRequest {
	s.FileFormat = FileFormat
	return s
}

func (s *AlterInternalStageStageRequest) WithCopyOptions(CopyOptions *StageCopyOptionsRequest) *AlterInternalStageStageRequest {
	s.CopyOptions = CopyOptions
	return s
}

func (s *AlterInternalStageStageRequest) WithComment(Comment *string) *AlterInternalStageStageRequest {
	s.Comment = Comment
	return s
}

func NewAlterExternalS3StageStageRequest(
	name SchemaObjectIdentifier,
) *AlterExternalS3StageStageRequest {
	s := AlterExternalS3StageStageRequest{}
	s.name = name
	return &s
}

func (s *AlterExternalS3StageStageRequest) WithIfExists(IfExists *bool) *AlterExternalS3StageStageRequest {
	s.IfExists = IfExists
	return s
}

func (s *AlterExternalS3StageStageRequest) WithExternalStageParams(ExternalStageParams *ExternalS3StageParamsRequest) *AlterExternalS3StageStageRequest {
	s.ExternalStageParams = ExternalStageParams
	return s
}

func (s *AlterExternalS3StageStageRequest) WithFileFormat(FileFormat *StageFileFormatRequest) *AlterExternalS3StageStageRequest {
	s.FileFormat = FileFormat
	return s
}

func (s *AlterExternalS3StageStageRequest) WithCopyOptions(CopyOptions *StageCopyOptionsRequest) *AlterExternalS3StageStageRequest {
	s.CopyOptions = CopyOptions
	return s
}

func (s *AlterExternalS3StageStageRequest) WithComment(Comment *string) *AlterExternalS3StageStageRequest {
	s.Comment = Comment
	return s
}

func NewAlterExternalGCSStageStageRequest(
	name SchemaObjectIdentifier,
) *AlterExternalGCSStageStageRequest {
	s := AlterExternalGCSStageStageRequest{}
	s.name = name
	return &s
}

func (s *AlterExternalGCSStageStageRequest) WithIfExists(IfExists *bool) *AlterExternalGCSStageStageRequest {
	s.IfExists = IfExists
	return s
}

func (s *AlterExternalGCSStageStageRequest) WithExternalStageParams(ExternalStageParams *ExternalGCSStageParamsRequest) *AlterExternalGCSStageStageRequest {
	s.ExternalStageParams = ExternalStageParams
	return s
}

func (s *AlterExternalGCSStageStageRequest) WithFileFormat(FileFormat *StageFileFormatRequest) *AlterExternalGCSStageStageRequest {
	s.FileFormat = FileFormat
	return s
}

func (s *AlterExternalGCSStageStageRequest) WithCopyOptions(CopyOptions *StageCopyOptionsRequest) *AlterExternalGCSStageStageRequest {
	s.CopyOptions = CopyOptions
	return s
}

func (s *AlterExternalGCSStageStageRequest) WithComment(Comment *string) *AlterExternalGCSStageStageRequest {
	s.Comment = Comment
	return s
}

func NewAlterExternalAzureStageStageRequest(
	name SchemaObjectIdentifier,
) *AlterExternalAzureStageStageRequest {
	s := AlterExternalAzureStageStageRequest{}
	s.name = name
	return &s
}

func (s *AlterExternalAzureStageStageRequest) WithIfExists(IfExists *bool) *AlterExternalAzureStageStageRequest {
	s.IfExists = IfExists
	return s
}

func (s *AlterExternalAzureStageStageRequest) WithExternalStageParams(ExternalStageParams *ExternalAzureStageParamsRequest) *AlterExternalAzureStageStageRequest {
	s.ExternalStageParams = ExternalStageParams
	return s
}

func (s *AlterExternalAzureStageStageRequest) WithFileFormat(FileFormat *StageFileFormatRequest) *AlterExternalAzureStageStageRequest {
	s.FileFormat = FileFormat
	return s
}

func (s *AlterExternalAzureStageStageRequest) WithCopyOptions(CopyOptions *StageCopyOptionsRequest) *AlterExternalAzureStageStageRequest {
	s.CopyOptions = CopyOptions
	return s
}

func (s *AlterExternalAzureStageStageRequest) WithComment(Comment *string) *AlterExternalAzureStageStageRequest {
	s.Comment = Comment
	return s
}

func NewAlterDirectoryTableStageRequest(
	name SchemaObjectIdentifier,
) *AlterDirectoryTableStageRequest {
	s := AlterDirectoryTableStageRequest{}
	s.name = name
	return &s
}

func (s *AlterDirectoryTableStageRequest) WithIfExists(IfExists *bool) *AlterDirectoryTableStageRequest {
	s.IfExists = IfExists
	return s
}

func (s *AlterDirectoryTableStageRequest) WithSetDirectory(SetDirectory *DirectoryTableSetRequest) *AlterDirectoryTableStageRequest {
	s.SetDirectory = SetDirectory
	return s
}

func (s *AlterDirectoryTableStageRequest) WithRefresh(Refresh *DirectoryTableRefreshRequest) *AlterDirectoryTableStageRequest {
	s.Refresh = Refresh
	return s
}

func NewDirectoryTableSetRequest(
	Enable bool,
) *DirectoryTableSetRequest {
	s := DirectoryTableSetRequest{}
	s.Enable = Enable
	return &s
}

func NewDirectoryTableRefreshRequest() *DirectoryTableRefreshRequest {
	return &DirectoryTableRefreshRequest{}
}

func (s *DirectoryTableRefreshRequest) WithSubpath(Subpath *string) *DirectoryTableRefreshRequest {
	s.Subpath = Subpath
	return s
}

func NewDropStageRequest(
	name SchemaObjectIdentifier,
) *DropStageRequest {
	s := DropStageRequest{}
	s.name = name
	return &s
}

func (s *DropStageRequest) WithIfExists(IfExists *bool) *DropStageRequest {
	s.IfExists = IfExists
	return s
}

func NewDescribeStageRequest(
	name SchemaObjectIdentifier,
) *DescribeStageRequest {
	s := DescribeStageRequest{}
	s.name = name
	return &s
}

func NewShowStageRequest() *ShowStageRequest {
	return &ShowStageRequest{}
}

func (s *ShowStageRequest) WithLike(Like *Like) *ShowStageRequest {
	s.Like = Like
	return s
}

func (s *ShowStageRequest) WithIn(In *In) *ShowStageRequest {
	s.In = In
	return s
}

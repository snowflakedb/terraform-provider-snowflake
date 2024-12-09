// Code generated by dto builder generator; DO NOT EDIT.

package sdk

// imports added manually
import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func NewCreateForJavaProcedureRequest(
	name SchemaObjectIdentifier,
	Returns ProcedureReturnsRequest,
	RuntimeVersion string,
	Packages []ProcedurePackageRequest,
	Handler string,
) *CreateForJavaProcedureRequest {
	s := CreateForJavaProcedureRequest{}
	s.name = name
	s.Returns = Returns
	s.RuntimeVersion = RuntimeVersion
	s.Packages = Packages
	s.Handler = Handler
	return &s
}

func (s *CreateForJavaProcedureRequest) WithOrReplace(OrReplace bool) *CreateForJavaProcedureRequest {
	s.OrReplace = &OrReplace
	return s
}

func (s *CreateForJavaProcedureRequest) WithSecure(Secure bool) *CreateForJavaProcedureRequest {
	s.Secure = &Secure
	return s
}

func (s *CreateForJavaProcedureRequest) WithArguments(Arguments []ProcedureArgumentRequest) *CreateForJavaProcedureRequest {
	s.Arguments = Arguments
	return s
}

func (s *CreateForJavaProcedureRequest) WithCopyGrants(CopyGrants bool) *CreateForJavaProcedureRequest {
	s.CopyGrants = &CopyGrants
	return s
}

func (s *CreateForJavaProcedureRequest) WithImports(Imports []ProcedureImportRequest) *CreateForJavaProcedureRequest {
	s.Imports = Imports
	return s
}

func (s *CreateForJavaProcedureRequest) WithExternalAccessIntegrations(ExternalAccessIntegrations []AccountObjectIdentifier) *CreateForJavaProcedureRequest {
	s.ExternalAccessIntegrations = ExternalAccessIntegrations
	return s
}

func (s *CreateForJavaProcedureRequest) WithSecrets(Secrets []SecretReference) *CreateForJavaProcedureRequest {
	s.Secrets = Secrets
	return s
}

func (s *CreateForJavaProcedureRequest) WithTargetPath(TargetPath string) *CreateForJavaProcedureRequest {
	s.TargetPath = &TargetPath
	return s
}

func (s *CreateForJavaProcedureRequest) WithNullInputBehavior(NullInputBehavior NullInputBehavior) *CreateForJavaProcedureRequest {
	s.NullInputBehavior = &NullInputBehavior
	return s
}

func (s *CreateForJavaProcedureRequest) WithComment(Comment string) *CreateForJavaProcedureRequest {
	s.Comment = &Comment
	return s
}

func (s *CreateForJavaProcedureRequest) WithExecuteAs(ExecuteAs ExecuteAs) *CreateForJavaProcedureRequest {
	s.ExecuteAs = &ExecuteAs
	return s
}

func (s *CreateForJavaProcedureRequest) WithProcedureDefinition(ProcedureDefinition string) *CreateForJavaProcedureRequest {
	s.ProcedureDefinition = &ProcedureDefinition
	return s
}

func NewProcedureArgumentRequest(
	ArgName string,
	ArgDataType datatypes.DataType,
) *ProcedureArgumentRequest {
	s := ProcedureArgumentRequest{}
	s.ArgName = ArgName
	s.ArgDataType = ArgDataType
	return &s
}

func (s *ProcedureArgumentRequest) WithArgDataTypeOld(ArgDataTypeOld DataType) *ProcedureArgumentRequest {
	s.ArgDataTypeOld = ArgDataTypeOld
	return s
}

func (s *ProcedureArgumentRequest) WithDefaultValue(DefaultValue string) *ProcedureArgumentRequest {
	s.DefaultValue = &DefaultValue
	return s
}

func NewProcedureReturnsRequest() *ProcedureReturnsRequest {
	return &ProcedureReturnsRequest{}
}

func (s *ProcedureReturnsRequest) WithResultDataType(ResultDataType ProcedureReturnsResultDataTypeRequest) *ProcedureReturnsRequest {
	s.ResultDataType = &ResultDataType
	return s
}

func (s *ProcedureReturnsRequest) WithTable(Table ProcedureReturnsTableRequest) *ProcedureReturnsRequest {
	s.Table = &Table
	return s
}

func NewProcedureReturnsResultDataTypeRequest(
	ResultDataType datatypes.DataType,
) *ProcedureReturnsResultDataTypeRequest {
	s := ProcedureReturnsResultDataTypeRequest{}
	s.ResultDataType = ResultDataType
	return &s
}

func (s *ProcedureReturnsResultDataTypeRequest) WithResultDataTypeOld(ResultDataTypeOld DataType) *ProcedureReturnsResultDataTypeRequest {
	s.ResultDataTypeOld = ResultDataTypeOld
	return s
}

func (s *ProcedureReturnsResultDataTypeRequest) WithNull(Null bool) *ProcedureReturnsResultDataTypeRequest {
	s.Null = &Null
	return s
}

func (s *ProcedureReturnsResultDataTypeRequest) WithNotNull(NotNull bool) *ProcedureReturnsResultDataTypeRequest {
	s.NotNull = &NotNull
	return s
}

func NewProcedureReturnsTableRequest() *ProcedureReturnsTableRequest {
	return &ProcedureReturnsTableRequest{}
}

func (s *ProcedureReturnsTableRequest) WithColumns(Columns []ProcedureColumnRequest) *ProcedureReturnsTableRequest {
	s.Columns = Columns
	return s
}

func NewProcedureColumnRequest(
	ColumnName string,
	ColumnDataType datatypes.DataType,
) *ProcedureColumnRequest {
	s := ProcedureColumnRequest{}
	s.ColumnName = ColumnName
	s.ColumnDataType = ColumnDataType
	return &s
}

func (s *ProcedureColumnRequest) WithColumnDataTypeOld(ColumnDataTypeOld DataType) *ProcedureColumnRequest {
	s.ColumnDataTypeOld = ColumnDataTypeOld
	return s
}

func NewProcedurePackageRequest(
	Package string,
) *ProcedurePackageRequest {
	s := ProcedurePackageRequest{}
	s.Package = Package
	return &s
}

func NewProcedureImportRequest(
	Import string,
) *ProcedureImportRequest {
	s := ProcedureImportRequest{}
	s.Import = Import
	return &s
}

func NewCreateForJavaScriptProcedureRequest(
	name SchemaObjectIdentifier,
	ResultDataType datatypes.DataType,
	ProcedureDefinition string,
) *CreateForJavaScriptProcedureRequest {
	s := CreateForJavaScriptProcedureRequest{}
	s.name = name
	s.ResultDataType = ResultDataType
	s.ProcedureDefinition = ProcedureDefinition
	return &s
}

func (s *CreateForJavaScriptProcedureRequest) WithOrReplace(OrReplace bool) *CreateForJavaScriptProcedureRequest {
	s.OrReplace = &OrReplace
	return s
}

func (s *CreateForJavaScriptProcedureRequest) WithSecure(Secure bool) *CreateForJavaScriptProcedureRequest {
	s.Secure = &Secure
	return s
}

func (s *CreateForJavaScriptProcedureRequest) WithArguments(Arguments []ProcedureArgumentRequest) *CreateForJavaScriptProcedureRequest {
	s.Arguments = Arguments
	return s
}

func (s *CreateForJavaScriptProcedureRequest) WithCopyGrants(CopyGrants bool) *CreateForJavaScriptProcedureRequest {
	s.CopyGrants = &CopyGrants
	return s
}

func (s *CreateForJavaScriptProcedureRequest) WithResultDataTypeOld(ResultDataTypeOld DataType) *CreateForJavaScriptProcedureRequest {
	s.ResultDataTypeOld = ResultDataTypeOld
	return s
}

func (s *CreateForJavaScriptProcedureRequest) WithNotNull(NotNull bool) *CreateForJavaScriptProcedureRequest {
	s.NotNull = &NotNull
	return s
}

func (s *CreateForJavaScriptProcedureRequest) WithNullInputBehavior(NullInputBehavior NullInputBehavior) *CreateForJavaScriptProcedureRequest {
	s.NullInputBehavior = &NullInputBehavior
	return s
}

func (s *CreateForJavaScriptProcedureRequest) WithComment(Comment string) *CreateForJavaScriptProcedureRequest {
	s.Comment = &Comment
	return s
}

func (s *CreateForJavaScriptProcedureRequest) WithExecuteAs(ExecuteAs ExecuteAs) *CreateForJavaScriptProcedureRequest {
	s.ExecuteAs = &ExecuteAs
	return s
}

func NewCreateForPythonProcedureRequest(
	name SchemaObjectIdentifier,
	Returns ProcedureReturnsRequest,
	RuntimeVersion string,
	Packages []ProcedurePackageRequest,
	Handler string,
) *CreateForPythonProcedureRequest {
	s := CreateForPythonProcedureRequest{}
	s.name = name
	s.Returns = Returns
	s.RuntimeVersion = RuntimeVersion
	s.Packages = Packages
	s.Handler = Handler
	return &s
}

func (s *CreateForPythonProcedureRequest) WithOrReplace(OrReplace bool) *CreateForPythonProcedureRequest {
	s.OrReplace = &OrReplace
	return s
}

func (s *CreateForPythonProcedureRequest) WithSecure(Secure bool) *CreateForPythonProcedureRequest {
	s.Secure = &Secure
	return s
}

func (s *CreateForPythonProcedureRequest) WithArguments(Arguments []ProcedureArgumentRequest) *CreateForPythonProcedureRequest {
	s.Arguments = Arguments
	return s
}

func (s *CreateForPythonProcedureRequest) WithCopyGrants(CopyGrants bool) *CreateForPythonProcedureRequest {
	s.CopyGrants = &CopyGrants
	return s
}

func (s *CreateForPythonProcedureRequest) WithImports(Imports []ProcedureImportRequest) *CreateForPythonProcedureRequest {
	s.Imports = Imports
	return s
}

func (s *CreateForPythonProcedureRequest) WithExternalAccessIntegrations(ExternalAccessIntegrations []AccountObjectIdentifier) *CreateForPythonProcedureRequest {
	s.ExternalAccessIntegrations = ExternalAccessIntegrations
	return s
}

func (s *CreateForPythonProcedureRequest) WithSecrets(Secrets []SecretReference) *CreateForPythonProcedureRequest {
	s.Secrets = Secrets
	return s
}

func (s *CreateForPythonProcedureRequest) WithNullInputBehavior(NullInputBehavior NullInputBehavior) *CreateForPythonProcedureRequest {
	s.NullInputBehavior = &NullInputBehavior
	return s
}

func (s *CreateForPythonProcedureRequest) WithComment(Comment string) *CreateForPythonProcedureRequest {
	s.Comment = &Comment
	return s
}

func (s *CreateForPythonProcedureRequest) WithExecuteAs(ExecuteAs ExecuteAs) *CreateForPythonProcedureRequest {
	s.ExecuteAs = &ExecuteAs
	return s
}

func (s *CreateForPythonProcedureRequest) WithProcedureDefinition(ProcedureDefinition string) *CreateForPythonProcedureRequest {
	s.ProcedureDefinition = &ProcedureDefinition
	return s
}

func NewCreateForScalaProcedureRequest(
	name SchemaObjectIdentifier,
	Returns ProcedureReturnsRequest,
	RuntimeVersion string,
	Packages []ProcedurePackageRequest,
	Handler string,
) *CreateForScalaProcedureRequest {
	s := CreateForScalaProcedureRequest{}
	s.name = name
	s.Returns = Returns
	s.RuntimeVersion = RuntimeVersion
	s.Packages = Packages
	s.Handler = Handler
	return &s
}

func (s *CreateForScalaProcedureRequest) WithOrReplace(OrReplace bool) *CreateForScalaProcedureRequest {
	s.OrReplace = &OrReplace
	return s
}

func (s *CreateForScalaProcedureRequest) WithSecure(Secure bool) *CreateForScalaProcedureRequest {
	s.Secure = &Secure
	return s
}

func (s *CreateForScalaProcedureRequest) WithArguments(Arguments []ProcedureArgumentRequest) *CreateForScalaProcedureRequest {
	s.Arguments = Arguments
	return s
}

func (s *CreateForScalaProcedureRequest) WithCopyGrants(CopyGrants bool) *CreateForScalaProcedureRequest {
	s.CopyGrants = &CopyGrants
	return s
}

func (s *CreateForScalaProcedureRequest) WithImports(Imports []ProcedureImportRequest) *CreateForScalaProcedureRequest {
	s.Imports = Imports
	return s
}

func (s *CreateForScalaProcedureRequest) WithTargetPath(TargetPath string) *CreateForScalaProcedureRequest {
	s.TargetPath = &TargetPath
	return s
}

func (s *CreateForScalaProcedureRequest) WithNullInputBehavior(NullInputBehavior NullInputBehavior) *CreateForScalaProcedureRequest {
	s.NullInputBehavior = &NullInputBehavior
	return s
}

func (s *CreateForScalaProcedureRequest) WithComment(Comment string) *CreateForScalaProcedureRequest {
	s.Comment = &Comment
	return s
}

func (s *CreateForScalaProcedureRequest) WithExecuteAs(ExecuteAs ExecuteAs) *CreateForScalaProcedureRequest {
	s.ExecuteAs = &ExecuteAs
	return s
}

func (s *CreateForScalaProcedureRequest) WithProcedureDefinition(ProcedureDefinition string) *CreateForScalaProcedureRequest {
	s.ProcedureDefinition = &ProcedureDefinition
	return s
}

func NewCreateForSQLProcedureRequest(
	name SchemaObjectIdentifier,
	Returns ProcedureSQLReturnsRequest,
	ProcedureDefinition string,
) *CreateForSQLProcedureRequest {
	s := CreateForSQLProcedureRequest{}
	s.name = name
	s.Returns = Returns
	s.ProcedureDefinition = ProcedureDefinition
	return &s
}

func (s *CreateForSQLProcedureRequest) WithOrReplace(OrReplace bool) *CreateForSQLProcedureRequest {
	s.OrReplace = &OrReplace
	return s
}

func (s *CreateForSQLProcedureRequest) WithSecure(Secure bool) *CreateForSQLProcedureRequest {
	s.Secure = &Secure
	return s
}

func (s *CreateForSQLProcedureRequest) WithArguments(Arguments []ProcedureArgumentRequest) *CreateForSQLProcedureRequest {
	s.Arguments = Arguments
	return s
}

func (s *CreateForSQLProcedureRequest) WithCopyGrants(CopyGrants bool) *CreateForSQLProcedureRequest {
	s.CopyGrants = &CopyGrants
	return s
}

func (s *CreateForSQLProcedureRequest) WithNullInputBehavior(NullInputBehavior NullInputBehavior) *CreateForSQLProcedureRequest {
	s.NullInputBehavior = &NullInputBehavior
	return s
}

func (s *CreateForSQLProcedureRequest) WithComment(Comment string) *CreateForSQLProcedureRequest {
	s.Comment = &Comment
	return s
}

func (s *CreateForSQLProcedureRequest) WithExecuteAs(ExecuteAs ExecuteAs) *CreateForSQLProcedureRequest {
	s.ExecuteAs = &ExecuteAs
	return s
}

func NewProcedureSQLReturnsRequest() *ProcedureSQLReturnsRequest {
	return &ProcedureSQLReturnsRequest{}
}

func (s *ProcedureSQLReturnsRequest) WithResultDataType(ResultDataType ProcedureReturnsResultDataTypeRequest) *ProcedureSQLReturnsRequest {
	s.ResultDataType = &ResultDataType
	return s
}

func (s *ProcedureSQLReturnsRequest) WithTable(Table ProcedureReturnsTableRequest) *ProcedureSQLReturnsRequest {
	s.Table = &Table
	return s
}

func (s *ProcedureSQLReturnsRequest) WithNotNull(NotNull bool) *ProcedureSQLReturnsRequest {
	s.NotNull = &NotNull
	return s
}

func NewAlterProcedureRequest(
	name SchemaObjectIdentifierWithArguments,
) *AlterProcedureRequest {
	s := AlterProcedureRequest{}
	s.name = name
	return &s
}

func (s *AlterProcedureRequest) WithIfExists(IfExists bool) *AlterProcedureRequest {
	s.IfExists = &IfExists
	return s
}

func (s *AlterProcedureRequest) WithRenameTo(RenameTo SchemaObjectIdentifier) *AlterProcedureRequest {
	s.RenameTo = &RenameTo
	return s
}

func (s *AlterProcedureRequest) WithSetComment(SetComment string) *AlterProcedureRequest {
	s.SetComment = &SetComment
	return s
}

func (s *AlterProcedureRequest) WithSetLogLevel(SetLogLevel string) *AlterProcedureRequest {
	s.SetLogLevel = &SetLogLevel
	return s
}

func (s *AlterProcedureRequest) WithSetTraceLevel(SetTraceLevel string) *AlterProcedureRequest {
	s.SetTraceLevel = &SetTraceLevel
	return s
}

func (s *AlterProcedureRequest) WithUnsetComment(UnsetComment bool) *AlterProcedureRequest {
	s.UnsetComment = &UnsetComment
	return s
}

func (s *AlterProcedureRequest) WithSetTags(SetTags []TagAssociation) *AlterProcedureRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterProcedureRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterProcedureRequest {
	s.UnsetTags = UnsetTags
	return s
}

func (s *AlterProcedureRequest) WithExecuteAs(ExecuteAs ExecuteAs) *AlterProcedureRequest {
	s.ExecuteAs = &ExecuteAs
	return s
}

func NewDropProcedureRequest(
	name SchemaObjectIdentifierWithArguments,
) *DropProcedureRequest {
	s := DropProcedureRequest{}
	s.name = name
	return &s
}

func (s *DropProcedureRequest) WithIfExists(IfExists bool) *DropProcedureRequest {
	s.IfExists = &IfExists
	return s
}

func NewShowProcedureRequest() *ShowProcedureRequest {
	return &ShowProcedureRequest{}
}

func (s *ShowProcedureRequest) WithLike(Like Like) *ShowProcedureRequest {
	s.Like = &Like
	return s
}

func (s *ShowProcedureRequest) WithIn(In In) *ShowProcedureRequest {
	s.In = &In
	return s
}

func NewDescribeProcedureRequest(
	name SchemaObjectIdentifierWithArguments,
) *DescribeProcedureRequest {
	s := DescribeProcedureRequest{}
	s.name = name
	return &s
}

func NewCallProcedureRequest(
	name SchemaObjectIdentifier,
) *CallProcedureRequest {
	s := CallProcedureRequest{}
	s.name = name
	return &s
}

func (s *CallProcedureRequest) WithCallArguments(CallArguments []string) *CallProcedureRequest {
	s.CallArguments = CallArguments
	return s
}

func (s *CallProcedureRequest) WithScriptingVariable(ScriptingVariable string) *CallProcedureRequest {
	s.ScriptingVariable = &ScriptingVariable
	return s
}

func NewCreateAndCallForJavaProcedureRequest(
	Name AccountObjectIdentifier,
	Returns ProcedureReturnsRequest,
	RuntimeVersion string,
	Packages []ProcedurePackageRequest,
	Handler string,
	ProcedureName AccountObjectIdentifier,
) *CreateAndCallForJavaProcedureRequest {
	s := CreateAndCallForJavaProcedureRequest{}
	s.Name = Name
	s.Returns = Returns
	s.RuntimeVersion = RuntimeVersion
	s.Packages = Packages
	s.Handler = Handler
	s.ProcedureName = ProcedureName
	return &s
}

func (s *CreateAndCallForJavaProcedureRequest) WithArguments(Arguments []ProcedureArgumentRequest) *CreateAndCallForJavaProcedureRequest {
	s.Arguments = Arguments
	return s
}

func (s *CreateAndCallForJavaProcedureRequest) WithImports(Imports []ProcedureImportRequest) *CreateAndCallForJavaProcedureRequest {
	s.Imports = Imports
	return s
}

func (s *CreateAndCallForJavaProcedureRequest) WithNullInputBehavior(NullInputBehavior NullInputBehavior) *CreateAndCallForJavaProcedureRequest {
	s.NullInputBehavior = &NullInputBehavior
	return s
}

func (s *CreateAndCallForJavaProcedureRequest) WithProcedureDefinition(ProcedureDefinition string) *CreateAndCallForJavaProcedureRequest {
	s.ProcedureDefinition = &ProcedureDefinition
	return s
}

func (s *CreateAndCallForJavaProcedureRequest) WithWithClause(WithClause ProcedureWithClauseRequest) *CreateAndCallForJavaProcedureRequest {
	s.WithClause = &WithClause
	return s
}

func (s *CreateAndCallForJavaProcedureRequest) WithCallArguments(CallArguments []string) *CreateAndCallForJavaProcedureRequest {
	s.CallArguments = CallArguments
	return s
}

func (s *CreateAndCallForJavaProcedureRequest) WithScriptingVariable(ScriptingVariable string) *CreateAndCallForJavaProcedureRequest {
	s.ScriptingVariable = &ScriptingVariable
	return s
}

func NewProcedureWithClauseRequest(
	CteName AccountObjectIdentifier,
	Statement string,
) *ProcedureWithClauseRequest {
	s := ProcedureWithClauseRequest{}
	s.CteName = CteName
	s.Statement = Statement
	return &s
}

func (s *ProcedureWithClauseRequest) WithCteColumns(CteColumns []string) *ProcedureWithClauseRequest {
	s.CteColumns = CteColumns
	return s
}

func NewCreateAndCallForScalaProcedureRequest(
	Name AccountObjectIdentifier,
	Returns ProcedureReturnsRequest,
	RuntimeVersion string,
	Packages []ProcedurePackageRequest,
	Handler string,
	ProcedureName AccountObjectIdentifier,
) *CreateAndCallForScalaProcedureRequest {
	s := CreateAndCallForScalaProcedureRequest{}
	s.Name = Name
	s.Returns = Returns
	s.RuntimeVersion = RuntimeVersion
	s.Packages = Packages
	s.Handler = Handler
	s.ProcedureName = ProcedureName
	return &s
}

func (s *CreateAndCallForScalaProcedureRequest) WithArguments(Arguments []ProcedureArgumentRequest) *CreateAndCallForScalaProcedureRequest {
	s.Arguments = Arguments
	return s
}

func (s *CreateAndCallForScalaProcedureRequest) WithImports(Imports []ProcedureImportRequest) *CreateAndCallForScalaProcedureRequest {
	s.Imports = Imports
	return s
}

func (s *CreateAndCallForScalaProcedureRequest) WithNullInputBehavior(NullInputBehavior NullInputBehavior) *CreateAndCallForScalaProcedureRequest {
	s.NullInputBehavior = &NullInputBehavior
	return s
}

func (s *CreateAndCallForScalaProcedureRequest) WithProcedureDefinition(ProcedureDefinition string) *CreateAndCallForScalaProcedureRequest {
	s.ProcedureDefinition = &ProcedureDefinition
	return s
}

func (s *CreateAndCallForScalaProcedureRequest) WithWithClauses(WithClauses []ProcedureWithClauseRequest) *CreateAndCallForScalaProcedureRequest {
	s.WithClauses = WithClauses
	return s
}

func (s *CreateAndCallForScalaProcedureRequest) WithCallArguments(CallArguments []string) *CreateAndCallForScalaProcedureRequest {
	s.CallArguments = CallArguments
	return s
}

func (s *CreateAndCallForScalaProcedureRequest) WithScriptingVariable(ScriptingVariable string) *CreateAndCallForScalaProcedureRequest {
	s.ScriptingVariable = &ScriptingVariable
	return s
}

func NewCreateAndCallForJavaScriptProcedureRequest(
	Name AccountObjectIdentifier,
	ResultDataType datatypes.DataType,
	ProcedureDefinition string,
	ProcedureName AccountObjectIdentifier,
) *CreateAndCallForJavaScriptProcedureRequest {
	s := CreateAndCallForJavaScriptProcedureRequest{}
	s.Name = Name
	s.ResultDataType = ResultDataType
	s.ProcedureDefinition = ProcedureDefinition
	s.ProcedureName = ProcedureName
	return &s
}

func (s *CreateAndCallForJavaScriptProcedureRequest) WithArguments(Arguments []ProcedureArgumentRequest) *CreateAndCallForJavaScriptProcedureRequest {
	s.Arguments = Arguments
	return s
}

func (s *CreateAndCallForJavaScriptProcedureRequest) WithResultDataTypeOld(ResultDataTypeOld DataType) *CreateAndCallForJavaScriptProcedureRequest {
	s.ResultDataTypeOld = ResultDataTypeOld
	return s
}

func (s *CreateAndCallForJavaScriptProcedureRequest) WithNotNull(NotNull bool) *CreateAndCallForJavaScriptProcedureRequest {
	s.NotNull = &NotNull
	return s
}

func (s *CreateAndCallForJavaScriptProcedureRequest) WithNullInputBehavior(NullInputBehavior NullInputBehavior) *CreateAndCallForJavaScriptProcedureRequest {
	s.NullInputBehavior = &NullInputBehavior
	return s
}

func (s *CreateAndCallForJavaScriptProcedureRequest) WithWithClauses(WithClauses []ProcedureWithClauseRequest) *CreateAndCallForJavaScriptProcedureRequest {
	s.WithClauses = WithClauses
	return s
}

func (s *CreateAndCallForJavaScriptProcedureRequest) WithCallArguments(CallArguments []string) *CreateAndCallForJavaScriptProcedureRequest {
	s.CallArguments = CallArguments
	return s
}

func (s *CreateAndCallForJavaScriptProcedureRequest) WithScriptingVariable(ScriptingVariable string) *CreateAndCallForJavaScriptProcedureRequest {
	s.ScriptingVariable = &ScriptingVariable
	return s
}

func NewCreateAndCallForPythonProcedureRequest(
	Name AccountObjectIdentifier,
	Returns ProcedureReturnsRequest,
	RuntimeVersion string,
	Packages []ProcedurePackageRequest,
	Handler string,
	ProcedureName AccountObjectIdentifier,
) *CreateAndCallForPythonProcedureRequest {
	s := CreateAndCallForPythonProcedureRequest{}
	s.Name = Name
	s.Returns = Returns
	s.RuntimeVersion = RuntimeVersion
	s.Packages = Packages
	s.Handler = Handler
	s.ProcedureName = ProcedureName
	return &s
}

func (s *CreateAndCallForPythonProcedureRequest) WithArguments(Arguments []ProcedureArgumentRequest) *CreateAndCallForPythonProcedureRequest {
	s.Arguments = Arguments
	return s
}

func (s *CreateAndCallForPythonProcedureRequest) WithImports(Imports []ProcedureImportRequest) *CreateAndCallForPythonProcedureRequest {
	s.Imports = Imports
	return s
}

func (s *CreateAndCallForPythonProcedureRequest) WithNullInputBehavior(NullInputBehavior NullInputBehavior) *CreateAndCallForPythonProcedureRequest {
	s.NullInputBehavior = &NullInputBehavior
	return s
}

func (s *CreateAndCallForPythonProcedureRequest) WithProcedureDefinition(ProcedureDefinition string) *CreateAndCallForPythonProcedureRequest {
	s.ProcedureDefinition = &ProcedureDefinition
	return s
}

func (s *CreateAndCallForPythonProcedureRequest) WithWithClauses(WithClauses []ProcedureWithClauseRequest) *CreateAndCallForPythonProcedureRequest {
	s.WithClauses = WithClauses
	return s
}

func (s *CreateAndCallForPythonProcedureRequest) WithCallArguments(CallArguments []string) *CreateAndCallForPythonProcedureRequest {
	s.CallArguments = CallArguments
	return s
}

func (s *CreateAndCallForPythonProcedureRequest) WithScriptingVariable(ScriptingVariable string) *CreateAndCallForPythonProcedureRequest {
	s.ScriptingVariable = &ScriptingVariable
	return s
}

func NewCreateAndCallForSQLProcedureRequest(
	Name AccountObjectIdentifier,
	Returns ProcedureReturnsRequest,
	ProcedureDefinition string,
	ProcedureName AccountObjectIdentifier,
) *CreateAndCallForSQLProcedureRequest {
	s := CreateAndCallForSQLProcedureRequest{}
	s.Name = Name
	s.Returns = Returns
	s.ProcedureDefinition = ProcedureDefinition
	s.ProcedureName = ProcedureName
	return &s
}

func (s *CreateAndCallForSQLProcedureRequest) WithArguments(Arguments []ProcedureArgumentRequest) *CreateAndCallForSQLProcedureRequest {
	s.Arguments = Arguments
	return s
}

func (s *CreateAndCallForSQLProcedureRequest) WithNullInputBehavior(NullInputBehavior NullInputBehavior) *CreateAndCallForSQLProcedureRequest {
	s.NullInputBehavior = &NullInputBehavior
	return s
}

func (s *CreateAndCallForSQLProcedureRequest) WithWithClauses(WithClauses []ProcedureWithClauseRequest) *CreateAndCallForSQLProcedureRequest {
	s.WithClauses = WithClauses
	return s
}

func (s *CreateAndCallForSQLProcedureRequest) WithCallArguments(CallArguments []string) *CreateAndCallForSQLProcedureRequest {
	s.CallArguments = CallArguments
	return s
}

func (s *CreateAndCallForSQLProcedureRequest) WithScriptingVariable(ScriptingVariable string) *CreateAndCallForSQLProcedureRequest {
	s.ScriptingVariable = &ScriptingVariable
	return s
}

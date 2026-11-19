package sdk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
)

func catalogParameterResult(err error) (handled bool, alter bool, ret error) {
	if err == nil {
		return true, true, nil
	}
	if errors.Is(err, ErrParameterNotSupported) {
		return false, false, nil
	}
	return true, false, err
}

func objectParameterCatalogSetResult(key string, err error) (handled bool, alter bool, ret error) {
	if errors.Is(err, ErrEnumValueNotSupported) {
		log.Printf("[WARN] parameter catalog does not recognize an enum value for %s; falling back to the generic object parameter writer: %v", key, err)
		return false, false, nil
	}
	return catalogParameterResult(err)
}

type rawParameterSetter interface {
	SetParameterFromRaw(key, value string) error
}

type rawParameterUnsetter interface {
	UnsetParameterFromRaw(key string) error
}

func catalogSet[ID any, Req rawParameterSetter](ctx context.Context, name ObjectIdentifier, key, value string, newReq func() Req, alter func(context.Context, ID, Req) error) (bool, error) {
	id, ok := name.(ID)
	if !ok {
		return false, nil
	}
	req := newReq()
	handled, doAlter, err := objectParameterCatalogSetResult(key, req.SetParameterFromRaw(key, value))
	if !doAlter {
		return handled, err
	}
	return true, alter(ctx, id, req)
}

type objectParameterCatalogSetter func(*Client, context.Context, ObjectIdentifier, ObjectParameter, string) (bool, error)

func newObjectParameterCatalogSetter[ID any, Req rawParameterSetter](
	newReq func() Req,
	alter func(*Client, context.Context, ID, Req) error,
) objectParameterCatalogSetter {
	return func(client *Client, ctx context.Context, name ObjectIdentifier, parameter ObjectParameter, value string) (bool, error) {
		return catalogSet(ctx, name, string(parameter), value, newReq, func(ctx context.Context, id ID, req Req) error {
			return alter(client, ctx, id, req)
		})
	}
}

var objectParameterCatalogSetters = map[ObjectType]objectParameterCatalogSetter{
	ObjectTypeDatabase: newObjectParameterCatalogSetter(
		NewDatabaseSetRequest,
		func(client *Client, ctx context.Context, id AccountObjectIdentifier, req *DatabaseSetRequest) error {
			return client.Databases.Alter(ctx, NewAlterDatabaseRequest(id).WithSet(*req))
		},
	),
	ObjectTypeSchema: newObjectParameterCatalogSetter(
		NewSchemaSetRequest,
		func(client *Client, ctx context.Context, id DatabaseObjectIdentifier, req *SchemaSetRequest) error {
			return client.Schemas.Alter(ctx, NewAlterSchemaRequest(id).WithSet(*req))
		},
	),
	ObjectTypeWarehouse: newObjectParameterCatalogSetter(
		NewWarehouseSetRequest,
		func(client *Client, ctx context.Context, id AccountObjectIdentifier, req *WarehouseSetRequest) error {
			return client.Warehouses.Alter(ctx, NewAlterWarehouseRequest(id).WithSet(*req))
		},
	),
	ObjectTypeUser: newObjectParameterCatalogSetter(
		NewUserSetRequest,
		func(client *Client, ctx context.Context, id AccountObjectIdentifier, req *UserSetRequest) error {
			return client.Users.Alter(ctx, NewAlterUserRequest(id).WithSet(*req))
		},
	),
	ObjectTypeTask: newObjectParameterCatalogSetter(
		NewTaskSetRequest,
		func(client *Client, ctx context.Context, id SchemaObjectIdentifier, req *TaskSetRequest) error {
			return client.Tasks.Alter(ctx, NewAlterTaskRequest(id).WithSet(*req))
		},
	),
	ObjectTypeService: newObjectParameterCatalogSetter(
		NewServiceSetRequest,
		func(client *Client, ctx context.Context, id SchemaObjectIdentifier, req *ServiceSetRequest) error {
			return client.Services.Alter(ctx, NewAlterServiceRequest(id).WithSet(*req))
		},
	),
	ObjectTypeHybridTable: newObjectParameterCatalogSetter(
		NewHybridTableSetPropertiesRequest,
		func(client *Client, ctx context.Context, id SchemaObjectIdentifier, req *HybridTableSetPropertiesRequest) error {
			return client.HybridTables.Alter(ctx, NewAlterHybridTableRequest(id).WithSet(*req))
		},
	),
	ObjectTypeIcebergTable: newObjectParameterCatalogSetter(
		NewIcebergTableSetPropertiesRequest,
		func(client *Client, ctx context.Context, id SchemaObjectIdentifier, req *IcebergTableSetPropertiesRequest) error {
			return client.IcebergTables.Alter(ctx, NewAlterIcebergTableRequest(id).WithSet(*req))
		},
	),
	ObjectTypeFunction: newObjectParameterCatalogSetter(
		NewFunctionSetRequest,
		func(client *Client, ctx context.Context, id SchemaObjectIdentifierWithArguments, req *FunctionSetRequest) error {
			return client.Functions.Alter(ctx, NewAlterFunctionRequest(id).WithSet(*req))
		},
	),
	ObjectTypeProcedure: newObjectParameterCatalogSetter(
		NewProcedureSetRequest,
		func(client *Client, ctx context.Context, id SchemaObjectIdentifierWithArguments, req *ProcedureSetRequest) error {
			return client.Procedures.Alter(ctx, NewAlterProcedureRequest(id).WithSet(*req))
		},
	),
}

func catalogUnset[ID any, Req rawParameterUnsetter](ctx context.Context, name ObjectIdentifier, key string, newReq func() Req, alter func(context.Context, ID, Req) error) (bool, error) {
	id, ok := name.(ID)
	if !ok {
		return false, nil
	}
	req := newReq()
	handled, doAlter, err := catalogParameterResult(req.UnsetParameterFromRaw(key))
	if !doAlter {
		return handled, err
	}
	return true, alter(ctx, id, req)
}

type objectParameterCatalogUnsetter func(*Client, context.Context, ObjectIdentifier, ObjectParameter) (bool, error)

func newObjectParameterCatalogUnsetter[ID any, Req rawParameterUnsetter](
	newReq func() Req,
	alter func(*Client, context.Context, ID, Req) error,
) objectParameterCatalogUnsetter {
	return func(client *Client, ctx context.Context, name ObjectIdentifier, parameter ObjectParameter) (bool, error) {
		return catalogUnset(ctx, name, string(parameter), newReq, func(ctx context.Context, id ID, req Req) error {
			return alter(client, ctx, id, req)
		})
	}
}

var objectParameterCatalogUnsetters = map[ObjectType]objectParameterCatalogUnsetter{
	ObjectTypeDatabase: newObjectParameterCatalogUnsetter(
		NewDatabaseUnsetRequest,
		func(client *Client, ctx context.Context, id AccountObjectIdentifier, req *DatabaseUnsetRequest) error {
			return client.Databases.Alter(ctx, NewAlterDatabaseRequest(id).WithUnset(*req))
		},
	),
	ObjectTypeSchema: newObjectParameterCatalogUnsetter(
		NewSchemaUnsetRequest,
		func(client *Client, ctx context.Context, id DatabaseObjectIdentifier, req *SchemaUnsetRequest) error {
			return client.Schemas.Alter(ctx, NewAlterSchemaRequest(id).WithUnset(*req))
		},
	),
	ObjectTypeWarehouse: newObjectParameterCatalogUnsetter(
		NewWarehouseUnsetRequest,
		func(client *Client, ctx context.Context, id AccountObjectIdentifier, req *WarehouseUnsetRequest) error {
			return client.Warehouses.Alter(ctx, NewAlterWarehouseRequest(id).WithUnset(*req))
		},
	),
	ObjectTypeUser: newObjectParameterCatalogUnsetter(
		NewUserUnsetRequest,
		func(client *Client, ctx context.Context, id AccountObjectIdentifier, req *UserUnsetRequest) error {
			return client.Users.Alter(ctx, NewAlterUserRequest(id).WithUnset(*req))
		},
	),
	ObjectTypeTask: newObjectParameterCatalogUnsetter(
		NewTaskUnsetRequest,
		func(client *Client, ctx context.Context, id SchemaObjectIdentifier, req *TaskUnsetRequest) error {
			return client.Tasks.Alter(ctx, NewAlterTaskRequest(id).WithUnset(*req))
		},
	),
	ObjectTypeService: newObjectParameterCatalogUnsetter(
		NewServiceUnsetRequest,
		func(client *Client, ctx context.Context, id SchemaObjectIdentifier, req *ServiceUnsetRequest) error {
			return client.Services.Alter(ctx, NewAlterServiceRequest(id).WithUnset(*req))
		},
	),
	ObjectTypeHybridTable: newObjectParameterCatalogUnsetter(
		NewHybridTableUnsetPropertiesRequest,
		func(client *Client, ctx context.Context, id SchemaObjectIdentifier, req *HybridTableUnsetPropertiesRequest) error {
			return client.HybridTables.Alter(ctx, NewAlterHybridTableRequest(id).WithUnset(*req))
		},
	),
	ObjectTypeIcebergTable: newObjectParameterCatalogUnsetter(
		NewIcebergTableUnsetPropertiesRequest,
		func(client *Client, ctx context.Context, id SchemaObjectIdentifier, req *IcebergTableUnsetPropertiesRequest) error {
			return client.IcebergTables.Alter(ctx, NewAlterIcebergTableRequest(id).WithUnset(*req))
		},
	),
	ObjectTypeFunction: newObjectParameterCatalogUnsetter(
		NewFunctionUnsetRequest,
		func(client *Client, ctx context.Context, id SchemaObjectIdentifierWithArguments, req *FunctionUnsetRequest) error {
			return client.Functions.Alter(ctx, NewAlterFunctionRequest(id).WithUnset(*req))
		},
	),
	ObjectTypeProcedure: newObjectParameterCatalogUnsetter(
		NewProcedureUnsetRequest,
		func(client *Client, ctx context.Context, id SchemaObjectIdentifierWithArguments, req *ProcedureUnsetRequest) error {
			return client.Procedures.Alter(ctx, NewAlterProcedureRequest(id).WithUnset(*req))
		},
	),
}

func (v *parameters) SetObjectParameterOnObjectViaCatalog(ctx context.Context, object Object, parameter ObjectParameter, value string) (bool, error) {
	setter, ok := objectParameterCatalogSetters[object.ObjectType]
	if !ok {
		return false, nil
	}
	return setter(v.client, ctx, object.Name, parameter, value)
}

func (v *parameters) UnsetObjectParameterOnObjectViaCatalog(ctx context.Context, object Object, parameter ObjectParameter) (bool, error) {
	unsetter, ok := objectParameterCatalogUnsetters[object.ObjectType]
	if !ok {
		return false, nil
	}
	return unsetter(v.client, ctx, object.Name, parameter)
}

func (v *parameters) SetObjectParameterOnAccountViaCatalog(ctx context.Context, parameter ObjectParameter, value string) (bool, error) {
	req := NewAccountParametersRequest()
	handled, alter, err := objectParameterCatalogSetResult(string(parameter), req.SetParameterFromRaw(string(parameter), value))
	if !alter {
		return handled, err
	}
	return true, v.client.Accounts.Alter(ctx, NewAlterAccountRequest().WithSet(*NewAccountSetRequest().WithParameters(*req)))
}

func (v *parameters) UnsetObjectParameterOnAccountViaCatalog(ctx context.Context, parameter ObjectParameter) (bool, error) {
	req := NewAccountParametersUnsetRequest()
	handled, alter, err := catalogParameterResult(req.UnsetParameterFromRaw(string(parameter)))
	if !alter {
		return handled, err
	}
	return true, v.client.Accounts.Alter(ctx, NewAlterAccountRequest().WithUnset(*NewAccountUnsetRequest().WithParameters(*req)))
}

func (v *SessionParameters) additionalValidations() error {
	var errs []error
	// Do not validate input and output formats, because there are a lot of them, and may be not supported in Go itself.
	// See https://docs.snowflake.com/en/sql-reference/date-time-input-output#supported-formats-for-auto-detection.
	if valueSet(v.ClientPrefetchThreads) {
		if !validateIntGreaterThanOrEqual(*v.ClientPrefetchThreads, 0) {
			errs = append(errs, errIntValue("SessionParameters", "ClientPrefetchThreads", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.ClientResultChunkSize) {
		if !validateIntGreaterThanOrEqual(*v.ClientResultChunkSize, 0) {
			errs = append(errs, errIntValue("SessionParameters", "ClientResultChunkSize", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.ClientSessionKeepAliveHeartbeatFrequency) {
		if !validateIntGreaterThanOrEqual(*v.ClientSessionKeepAliveHeartbeatFrequency, 0) {
			errs = append(errs, errIntValue("SessionParameters", "ClientSessionKeepAliveHeartbeatFrequency", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.HybridTableLockTimeout) {
		if !validateIntGreaterThanOrEqual(*v.HybridTableLockTimeout, 0) {
			errs = append(errs, errIntValue("SessionParameters", "HybridTableLockTimeout", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.JsonIndent) {
		if !validateIntGreaterThanOrEqual(*v.JsonIndent, 0) {
			errs = append(errs, errIntValue("SessionParameters", "JsonIndent", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.LockTimeout) {
		if !validateIntGreaterThanOrEqual(*v.LockTimeout, 0) {
			errs = append(errs, errIntValue("SessionParameters", "LockTimeout", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.MultiStatementCount) {
		if !validateIntGreaterThanOrEqual(*v.MultiStatementCount, 0) {
			errs = append(errs, errIntValue("SessionParameters", "MultiStatementCount", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.RowsPerResultset) {
		if !validateIntGreaterThanOrEqual(*v.RowsPerResultset, 0) {
			errs = append(errs, errIntValue("SessionParameters", "RowsPerResultset", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.StatementQueuedTimeoutInSeconds) {
		if !validateIntGreaterThanOrEqual(*v.StatementQueuedTimeoutInSeconds, 0) {
			errs = append(errs, errIntValue("SessionParameters", "StatementQueuedTimeoutInSeconds", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.StatementTimeoutInSeconds) {
		if !validateIntGreaterThanOrEqual(*v.StatementTimeoutInSeconds, 0) {
			errs = append(errs, errIntValue("SessionParameters", "StatementTimeoutInSeconds", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.TwoDigitCenturyStart) {
		if !validateIntGreaterThanOrEqual(*v.TwoDigitCenturyStart, 1900) {
			errs = append(errs, errIntValue("SessionParameters", "TwoDigitCenturyStart", IntErrGreaterOrEqual, 1900))
		}
	}
	if valueSet(v.WeekOfYearPolicy) {
		if !validateIntGreaterThanOrEqual(*v.WeekOfYearPolicy, 0) {
			errs = append(errs, errIntValue("SessionParameters", "WeekOfYearPolicy", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.WeekStart) {
		if !validateIntGreaterThanOrEqual(*v.WeekStart, 0) {
			errs = append(errs, errIntValue("SessionParameters", "WeekStart", IntErrGreaterOrEqual, 0))
		}
	}
	return errors.Join(errs...)
}

// ShowParametersOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-parameters.
type ShowParametersOptions struct {
	show       bool          `ddl:"static" sql:"SHOW"`
	parameters bool          `ddl:"static" sql:"PARAMETERS"`
	Like       *Like         `ddl:"keyword" sql:"LIKE"`
	In         *ParametersIn `ddl:"keyword" sql:"IN"`
}

func (opts *ShowParametersOptions) validate() error {
	if valueSet(opts.In) {
		if err := opts.In.validate(); err != nil {
			return err
		}
	}
	return nil
}

type ParametersIn struct {
	Session            *bool                               `ddl:"keyword" sql:"SESSION"`
	Account            *bool                               `ddl:"keyword" sql:"ACCOUNT"`
	User               AccountObjectIdentifier             `ddl:"identifier" sql:"USER"`
	Warehouse          AccountObjectIdentifier             `ddl:"identifier" sql:"WAREHOUSE"`
	Database           AccountObjectIdentifier             `ddl:"identifier" sql:"DATABASE"`
	Schema             DatabaseObjectIdentifier            `ddl:"identifier" sql:"SCHEMA"`
	Task               SchemaObjectIdentifier              `ddl:"identifier" sql:"TASK"`
	Table              SchemaObjectIdentifier              `ddl:"identifier" sql:"TABLE"`
	Function           SchemaObjectIdentifierWithArguments `ddl:"identifier" sql:"FUNCTION"`
	Procedure          SchemaObjectIdentifierWithArguments `ddl:"identifier" sql:"PROCEDURE"`
	OpenflowDeployment AccountObjectIdentifier             `ddl:"identifier" sql:"OPENFLOW DEPLOYMENT"`
	Service            SchemaObjectIdentifier              `ddl:"identifier" sql:"SERVICE"`
	Pipe               SchemaObjectIdentifier              `ddl:"identifier" sql:"PIPE"`
	// Object is the generic form (IN <object_type> <name>), used for the object types that have no dedicated field above.
	Object *ParametersInObject `ddl:"keyword"`
}

type ParametersInObject struct {
	ObjectType ObjectType       `ddl:"keyword"`
	Name       ObjectIdentifier `ddl:"identifier"`
}

func (v *ParametersIn) validate() error {
	if !anyValueSet(v.Session, v.Account, v.User, v.Warehouse, v.Database, v.Schema, v.Task, v.Table, v.Function, v.Procedure, v.OpenflowDeployment, v.Service, v.Pipe, v.Object) {
		return errors.Join(errAtLeastOneOf("Session", "Account", "User", "Warehouse", "Database", "Schema", "Task", "Table", "Function", "Procedure", "OpenflowDeployment", "Service", "Pipe", "Object"))
	}
	return nil
}

type parameterRow struct {
	Key         string `db:"key"`
	Value       string `db:"value"`
	Default     string `db:"default"`
	Level       string `db:"level"`
	Description string `db:"description"`
}

func (r parameterRow) convert() (*Parameter, error) {
	level, err := ToParameterType(r.Level)
	if err != nil {
		return nil, fmt.Errorf("parsing parameter type: %w", err)
	}
	return &Parameter{
		Key:         r.Key,
		Value:       r.Value,
		Default:     r.Default,
		Level:       level,
		Description: r.Description,
	}, nil
}

type Parameter struct {
	Key         string
	Value       string
	Default     string
	Level       ParameterType
	Description string
}

func (v *parameters) ShowParameters(ctx context.Context, opts *ShowParametersOptions) ([]*Parameter, error) {
	if opts == nil {
		opts = &ShowParametersOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	rows := []parameterRow{}
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	parameters := make([]*Parameter, len(rows))
	for i, row := range rows {
		parameter, err := row.convert()
		if err != nil {
			return nil, err
		}
		parameters[i] = parameter
	}
	return parameters, nil
}

func (v *parameters) ShowAccountParameter(ctx context.Context, parameter AccountParameter) (*Parameter, error) {
	opts := &ShowParametersOptions{
		Like: &Like{
			Pattern: String(string(parameter)),
		},
		In: &ParametersIn{
			Account: Bool(true),
		},
	}
	parameters, err := v.ShowParameters(ctx, opts)
	if err != nil {
		return nil, err
	}
	return exactParameter(parameters, string(parameter))
}

func (v *parameters) ShowSessionParameter(ctx context.Context, parameter SessionParameter) (*Parameter, error) {
	opts := &ShowParametersOptions{
		Like: &Like{
			Pattern: String(string(parameter)),
		},
		In: &ParametersIn{
			Session: Bool(true),
		},
	}
	parameters, err := v.ShowParameters(ctx, opts)
	if err != nil {
		return nil, err
	}
	return exactParameter(parameters, string(parameter))
}

func (v *parameters) ShowUserParameter(ctx context.Context, parameter UserParameter, userId AccountObjectIdentifier) (*Parameter, error) {
	opts := &ShowParametersOptions{
		Like: &Like{
			Pattern: String(string(parameter)),
		},
		In: &ParametersIn{
			User: userId,
		},
	}
	parameters, err := v.ShowParameters(ctx, opts)
	if err != nil {
		return nil, err
	}
	return exactParameter(parameters, string(parameter))
}

func (v *parameters) ShowObjectParameter(ctx context.Context, parameter ObjectParameter, object Object) (*Parameter, error) {
	opts := &ShowParametersOptions{
		Like: &Like{
			Pattern: String(string(parameter)),
		},
		In: &ParametersIn{},
	}
	switch object.ObjectType {
	case ObjectTypeWarehouse:
		opts.In.Warehouse = object.Name.(AccountObjectIdentifier)
	case ObjectTypeDatabase:
		opts.In.Database = object.Name.(AccountObjectIdentifier)
	case ObjectTypeSchema:
		opts.In.Schema = object.Name.(DatabaseObjectIdentifier)
	case ObjectTypeTask:
		opts.In.Task = object.Name.(SchemaObjectIdentifier)
	// Hybrid and Iceberg tables are shown with IN TABLE. IN HYBRID TABLE and IN ICEBERG TABLE are rejected.
	case ObjectTypeTable, ObjectTypeHybridTable, ObjectTypeIcebergTable:
		opts.In.Table = object.Name.(SchemaObjectIdentifier)
	case ObjectTypeUser:
		opts.In.User = object.Name.(AccountObjectIdentifier)
	case ObjectTypeFunction:
		opts.In.Function = object.Name.(SchemaObjectIdentifierWithArguments)
	case ObjectTypeProcedure:
		opts.In.Procedure = object.Name.(SchemaObjectIdentifierWithArguments)
	case ObjectTypeService:
		opts.In.Service = object.Name.(SchemaObjectIdentifier)
	case ObjectTypePipe:
		opts.In.Pipe = object.Name.(SchemaObjectIdentifier)
	default:
		// Any other object type, so that parameters of the object types that are not known to the provider can be read.
		if object.Name == nil {
			return nil, fmt.Errorf("object identifier is required for object type %s", object.ObjectType)
		}
		opts.In.Object = &ParametersInObject{ObjectType: object.ObjectType, Name: object.Name}
	}
	parameters, err := v.ShowParameters(ctx, opts)
	if err != nil {
		return nil, err
	}
	return exactParameter(parameters, string(parameter))
}

// exactParameter picks the row of the given name. The LIKE pattern used to show a parameter treats "_" as a wildcard,
// so the result may contain other parameters (e.g. when the name is not known to the provider).
func exactParameter(parameters []*Parameter, name string) (*Parameter, error) {
	for _, p := range parameters {
		if strings.EqualFold(p.Key, name) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("parameter %s not found", name)
}

// setWithGenericFallback sends the value with the catalog write path when the catalog handles the parameter (the kind is not used),
// and with the raw write path otherwise, with the value rendered by the kind.
func setWithGenericFallback(parameter string, kind ParameterValueKind, value string, viaCatalog func() (bool, error), raw func(rendered string) error) error {
	handled, err := viaCatalog()
	if err != nil {
		return err
	}
	if handled {
		log.Printf("[WARN] value kind %s is ignored for %s, because the parameter catalog supports the parameter", kind, parameter)
		return nil
	}
	rendered, err := kind.Render(parameter, value)
	if err != nil {
		return err
	}
	return raw(rendered)
}

// unsetWithGenericFallback is the unset counterpart of setWithGenericFallback. Unsetting needs no value, so there is no kind.
func unsetWithGenericFallback(viaCatalog func() (bool, error), raw func() error) error {
	handled, err := viaCatalog()
	if err != nil || handled {
		return err
	}
	return raw()
}

func (v *parameters) SetAccountParameterWithGenericFallback(ctx context.Context, parameter string, kind ParameterValueKind, value string) error {
	return setWithGenericFallback(parameter, kind, value,
		func() (bool, error) {
			return v.SetObjectParameterOnAccountViaCatalog(ctx, ObjectParameter(parameter), value)
		},
		func(rendered string) error {
			return v.SetAccountParameterRaw(ctx, NewSetAccountParameterRawParameterRequest(parameter, rendered))
		},
	)
}

func (v *parameters) UnsetAccountParameterWithGenericFallback(ctx context.Context, parameter string) error {
	return unsetWithGenericFallback(
		func() (bool, error) {
			return v.UnsetObjectParameterOnAccountViaCatalog(ctx, ObjectParameter(parameter))
		},
		func() error {
			return v.UnsetAccountParameterRaw(ctx, NewUnsetAccountParameterRawParameterRequest(parameter))
		},
	)
}

func (v *parameters) SetObjectParameterOnObjectWithGenericFallback(ctx context.Context, object Object, parameter ObjectParameter, kind ParameterValueKind, value string) error {
	return setWithGenericFallback(string(parameter), kind, value,
		func() (bool, error) { return v.SetObjectParameterOnObjectViaCatalog(ctx, object, parameter, value) },
		func(rendered string) error {
			return v.SetObjectParameterOnObjectRaw(ctx, NewSetObjectParameterOnObjectRawParameterRequest(object.ObjectType, object.Name, string(parameter), rendered))
		},
	)
}

func (v *parameters) UnsetObjectParameterOnObjectWithGenericFallback(ctx context.Context, object Object, parameter ObjectParameter) error {
	return unsetWithGenericFallback(
		func() (bool, error) { return v.UnsetObjectParameterOnObjectViaCatalog(ctx, object, parameter) },
		func() error {
			return v.UnsetObjectParameterOnObjectRaw(ctx, NewUnsetObjectParameterOnObjectRawParameterRequest(object.ObjectType, object.Name, string(parameter)))
		},
	)
}

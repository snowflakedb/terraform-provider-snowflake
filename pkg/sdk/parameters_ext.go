package sdk

import (
	"context"
	"errors"
	"fmt"
)

var _ validatable = new(setParameterOnObject)

func (v *parameters) SetAccountParameter(ctx context.Context, parameter AccountParameter, value string) error {
	parameters := NewAccountParametersRequest()
	if err := parameters.SetParameterFromRaw(string(parameter), value); err != nil {
		return err
	}
	return v.client.Accounts.Alter(ctx, NewAlterAccountRequest().WithSet(*NewAccountSetRequest().WithParameters(*parameters)))
}

// TODO(SNOW-1866453): add integration tests
func (v *parameters) UnsetAccountParameter(ctx context.Context, parameter AccountParameter) error {
	parameters := NewAccountParametersUnsetRequest()
	if err := parameters.UnsetParameterFromRaw(string(parameter)); err != nil {
		return err
	}
	return v.client.Accounts.Alter(ctx, NewAlterAccountRequest().WithUnset(*NewAccountUnsetRequest().WithParameters(*parameters)))
}

func (v *parameters) SetObjectParameterOnAccount(ctx context.Context, parameter ObjectParameter, value string) error {
	parameters := NewAccountParametersRequest()
	if err := parameters.SetParameterFromRaw(string(parameter), value); err != nil {
		return err
	}
	return v.client.Accounts.Alter(ctx, NewAlterAccountRequest().WithSet(*NewAccountSetRequest().WithParameters(*parameters)))
}

func (v *parameters) UnsetObjectParameterOnAccount(ctx context.Context, parameter ObjectParameter) error {
	parameters := NewAccountParametersUnsetRequest()
	if err := parameters.UnsetParameterFromRaw(string(parameter)); err != nil {
		return err
	}
	return v.client.Accounts.Alter(ctx, NewAlterAccountRequest().WithUnset(*NewAccountUnsetRequest().WithParameters(*parameters)))
}

type setParameterOnObject struct {
	alter            bool             `ddl:"static" sql:"ALTER"`
	objectType       ObjectType       `ddl:"keyword"`
	objectIdentifier ObjectIdentifier `ddl:"identifier"`
	set              bool             `ddl:"static" sql:"SET"`
	parameterKey     ObjectParameter  `ddl:"keyword"`
	equals           bool             `ddl:"static" sql:"="`
	parameterValue   string           `ddl:"keyword"`
}

func (v *setParameterOnObject) validate() error {
	return nil
}

func (v *parameters) SetObjectParameterOnObject(ctx context.Context, object Object, parameter ObjectParameter, value string) error {
	opts := &setParameterOnObject{
		objectType:       object.ObjectType,
		objectIdentifier: object.Name,
		parameterKey:     parameter,
		parameterValue:   value,
	}
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

type unsetParameterOnObject struct {
	alter            bool             `ddl:"static" sql:"ALTER"`
	objectType       ObjectType       `ddl:"keyword"`
	objectIdentifier ObjectIdentifier `ddl:"identifier"`
	unset            bool             `ddl:"static" sql:"UNSET"`
	parameterKey     ObjectParameter  `ddl:"keyword"`
}

func (v *unsetParameterOnObject) validate() error {
	return nil
}

func (v *parameters) UnsetObjectParameterOnObject(ctx context.Context, object Object, parameter ObjectParameter) error {
	opts := &unsetParameterOnObject{
		objectType:       object.ObjectType,
		objectIdentifier: object.Name,
		parameterKey:     parameter,
	}
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
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
}

func (v *ParametersIn) validate() error {
	if !anyValueSet(v.Session, v.Account, v.User, v.Warehouse, v.Database, v.Schema, v.Task, v.Table, v.Function, v.Procedure, v.OpenflowDeployment, v.Service) {
		return errors.Join(errAtLeastOneOf("Session", "Account", "User", "Warehouse", "Database", "Schema", "Task", "Table", "Function", "Procedure", "OpenflowDeployment", "Service"))
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
	if len(parameters) == 0 {
		return nil, fmt.Errorf("parameter %s not found", parameter)
	}
	return parameters[0], nil
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
	if len(parameters) == 0 {
		return nil, fmt.Errorf("parameter %s not found", parameter)
	}
	return parameters[0], nil
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
	if len(parameters) == 0 {
		return nil, fmt.Errorf("parameter %s not found", parameter)
	}
	return parameters[0], nil
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
	case ObjectTypeTable:
		opts.In.Table = object.Name.(SchemaObjectIdentifier)
	case ObjectTypeUser:
		opts.In.User = object.Name.(AccountObjectIdentifier)
	case ObjectTypeFunction:
		opts.In.Function = object.Name.(SchemaObjectIdentifierWithArguments)
	case ObjectTypeProcedure:
		opts.In.Procedure = object.Name.(SchemaObjectIdentifierWithArguments)
	case ObjectTypeService:
		opts.In.Service = object.Name.(SchemaObjectIdentifier)
	default:
		return nil, fmt.Errorf("unsupported object type %s", object.ObjectType)
	}
	parameters, err := v.ShowParameters(ctx, opts)
	if err != nil {
		return nil, err
	}
	if len(parameters) == 0 {
		return nil, fmt.Errorf("parameter %s not found", parameter)
	}
	return parameters[0], nil
}

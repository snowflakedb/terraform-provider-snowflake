package sdk

import (
	"context"
	"database/sql"
	"time"
)

type OpenflowDeployments interface {
	Create(ctx context.Context, request *CreateOpenflowDeploymentRequest) error
	Alter(ctx context.Context, request *AlterOpenflowDeploymentRequest) error
	Drop(ctx context.Context, request *DropOpenflowDeploymentRequest) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, request *ShowOpenflowDeploymentRequest) ([]OpenflowDeployment, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*OpenflowDeployment, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*OpenflowDeployment, error)
	Describe(ctx context.Context, id AccountObjectIdentifier) (*OpenflowDeploymentDetails, error)
}

// CreateOpenflowDeploymentOptions is based on CREATE OPENFLOW DEPLOYMENT.
type CreateOpenflowDeploymentOptions struct {
	create             bool                    `ddl:"static" sql:"CREATE"`
	openflowDeployment bool                    `ddl:"static" sql:"OPENFLOW DEPLOYMENT"`
	IfNotExists        *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name               AccountObjectIdentifier `ddl:"identifier"`
	DeploymentType     *OpenflowDeploymentType `ddl:"parameter,single_quotes" sql:"DEPLOYMENT_TYPE"`
	VpcType            *OpenflowVpcType        `ddl:"parameter,single_quotes" sql:"VPC_TYPE"`
	CustomIngressHostname *string              `ddl:"parameter,single_quotes" sql:"CUSTOM_INGRESS_HOSTNAME"`
	UsePrivateLink     *bool                   `ddl:"parameter" sql:"USE_PRIVATE_LINK"`
	UseUserAuthOverPrivatelink *bool           `ddl:"parameter" sql:"USE_USER_AUTH_OVER_PRIVATELINK"`
	EventTable         *string                 `ddl:"parameter,single_quotes" sql:"EVENT_TABLE"`
	DisplayName        *string                 `ddl:"parameter,single_quotes" sql:"DISPLAY_NAME"`
	Comment            *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterOpenflowDeploymentOptions is based on ALTER OPENFLOW DEPLOYMENT.
type AlterOpenflowDeploymentOptions struct {
	alter              bool                    `ddl:"static" sql:"ALTER"`
	openflowDeployment bool                    `ddl:"static" sql:"OPENFLOW DEPLOYMENT"`
	name               AccountObjectIdentifier `ddl:"identifier"`
	Set                *OpenflowDeploymentSet  `ddl:"keyword" sql:"SET"`
	Unset              *OpenflowDeploymentUnset `ddl:"list,no_parentheses" sql:"UNSET"`
}

type OpenflowDeploymentSet struct {
	Comment     *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
	DisplayName *string `ddl:"parameter,single_quotes" sql:"DISPLAY_NAME"`
	EventTable  *string `ddl:"parameter,single_quotes" sql:"EVENT_TABLE"`
}

type OpenflowDeploymentUnset struct {
	Comment     *bool `ddl:"keyword" sql:"COMMENT"`
	DisplayName *bool `ddl:"keyword" sql:"DISPLAY_NAME"`
	EventTable  *bool `ddl:"keyword" sql:"EVENT_TABLE"`
}

// DropOpenflowDeploymentOptions is based on DROP OPENFLOW DEPLOYMENT.
type DropOpenflowDeploymentOptions struct {
	drop               bool                    `ddl:"static" sql:"DROP"`
	openflowDeployment bool                    `ddl:"static" sql:"OPENFLOW DEPLOYMENT"`
	IfExists           *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name               AccountObjectIdentifier `ddl:"identifier"`
}

// ShowOpenflowDeploymentOptions is based on SHOW OPENFLOW DEPLOYMENTS.
type ShowOpenflowDeploymentOptions struct {
	show                bool  `ddl:"static" sql:"SHOW"`
	openflowDeployments bool  `ddl:"static" sql:"OPENFLOW DEPLOYMENTS"`
	Like                *Like `ddl:"keyword" sql:"LIKE"`
}

type openflowDeploymentRow struct {
	Name                       string         `db:"name"`
	DeploymentType             string         `db:"type"`
	Status                     string         `db:"status"`
	VpcType                    sql.NullString `db:"vpc_type"`
	DisplayName                sql.NullString `db:"display_name"`
	UsePrivateLink             bool           `db:"use_private_link"`
	UseUserAuthOverPrivatelink bool           `db:"use_user_auth_over_private_link"`
	CustomIngressHostname      sql.NullString `db:"custom_ingress_hostname"`
	Key                        sql.NullString `db:"key"`
	Owner                      string         `db:"owner"`
	Comment                    sql.NullString `db:"comment"`
	CreatedOn                  time.Time      `db:"created_on"`
	UpdatedOn                  time.Time      `db:"updated_on"`
}

type OpenflowDeployment struct {
	Name                       string
	Status                     OpenflowDeploymentStatus
	DeploymentType             OpenflowDeploymentType
	VpcType                    *OpenflowVpcType
	UsePrivateLink             bool
	UseUserAuthOverPrivatelink bool
	CustomIngressHostname      *string
	Key                        *string
	EventTable                 *string
	DisplayName                *string
	Comment                    *string
	Owner                      string
	CreatedOn                  time.Time
	UpdatedOn                  time.Time
}

func (v *OpenflowDeployment) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *OpenflowDeployment) ObjectType() ObjectType {
	return ObjectTypeOpenflowDeployment
}

// DescribeOpenflowDeploymentOptions is based on DESCRIBE OPENFLOW DEPLOYMENT.
type DescribeOpenflowDeploymentOptions struct {
	describe           bool                    `ddl:"static" sql:"DESCRIBE"`
	openflowDeployment bool                    `ddl:"static" sql:"OPENFLOW DEPLOYMENT"`
	name               AccountObjectIdentifier `ddl:"identifier"`
}

type openflowDeploymentDetailsRow struct {
	Name                       string         `db:"name"`
	DeploymentType             string         `db:"type"`
	Status                     string         `db:"status"`
	VpcType                    sql.NullString `db:"vpc_type"`
	DisplayName                sql.NullString `db:"display_name"`
	UsePrivateLink             bool           `db:"use_private_link"`
	UseUserAuthOverPrivatelink bool           `db:"use_user_auth_over_private_link"`
	CustomIngressHostname      sql.NullString `db:"custom_ingress_hostname"`
	Key                        sql.NullString `db:"key"`
	Owner                      string         `db:"owner"`
	Comment                    sql.NullString `db:"comment"`
	CreatedOn                  time.Time      `db:"created_on"`
	UpdatedOn                  time.Time      `db:"updated_on"`
	ErrorCode                  sql.NullString `db:"error_code"`
	StatusMessage              sql.NullString `db:"status_message"`
}

type OpenflowDeploymentDetails struct {
	Name                       string
	Status                     OpenflowDeploymentStatus
	DeploymentType             OpenflowDeploymentType
	VpcType                    *OpenflowVpcType
	UsePrivateLink             bool
	UseUserAuthOverPrivatelink bool
	CustomIngressHostname      *string
	Key                        *string
	EventTable                 *string
	DisplayName                *string
	Comment                    *string
	Owner                      string
	CreatedOn                  time.Time
	UpdatedOn                  time.Time
	ErrorCode                  *string
	StatusMessage              *string
}

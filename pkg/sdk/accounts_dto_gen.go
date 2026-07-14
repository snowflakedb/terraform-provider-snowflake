package sdk

// Stub request types — will be fully generated in Step 3.
// These are minimal definitions to allow compilation during Step 2 verification.

type CreateAccountRequest struct {
	name                     AccountObjectIdentifier
	adminName                string
	adminPassword            *string
	adminRsaPublicKey        *string
	adminUserType            *UserType
	firstName                *string
	lastName                 *string
	email                    string
	mustChangePassword       *bool
	edition                  AccountEdition
	regionGroup              *string
	region                   *string
	comment                  *string
	consumptionBillingEntity *string
	polaris                  *bool
}

func NewCreateAccountRequest(name AccountObjectIdentifier, adminName string, email string, edition AccountEdition) *CreateAccountRequest {
	return &CreateAccountRequest{name: name, adminName: adminName, email: email, edition: edition}
}
func (r *CreateAccountRequest) WithAdminPassword(v string) *CreateAccountRequest {
	r.adminPassword = &v
	return r
}
func (r *CreateAccountRequest) WithAdminRsaPublicKey(v string) *CreateAccountRequest {
	r.adminRsaPublicKey = &v
	return r
}
func (r *CreateAccountRequest) toOpts() *CreateAccountOptions {
	opts := &CreateAccountOptions{
		name:                     r.name,
		AdminName:                r.adminName,
		AdminPassword:            r.adminPassword,
		AdminRSAPublicKey:        r.adminRsaPublicKey,
		AdminUserType:            r.adminUserType,
		FirstName:                r.firstName,
		LastName:                 r.lastName,
		Email:                    r.email,
		MustChangePassword:       r.mustChangePassword,
		Edition:                  r.edition,
		RegionGroup:              r.regionGroup,
		Region:                   r.region,
		Comment:                  r.comment,
		ConsumptionBillingEntity: r.consumptionBillingEntity,
		Polaris:                  r.polaris,
	}
	return opts
}

type AlterAccountRequest struct {
	name     *AccountObjectIdentifier
	set      *AccountSet
	unset    *AccountUnset
	setTag   []TagAssociation
	unsetTag []ObjectIdentifier
	rename   *AccountRename
	drop     *AccountDrop
}

func NewAlterAccountRequest() *AlterAccountRequest                           { return &AlterAccountRequest{} }
func (r *AlterAccountRequest) WithSet(v AccountSet) *AlterAccountRequest     { r.set = &v; return r }
func (r *AlterAccountRequest) WithUnset(v AccountUnset) *AlterAccountRequest { r.unset = &v; return r }
func (r *AlterAccountRequest) toOpts() *AlterAccountOptions {
	opts := &AlterAccountOptions{
		Name:     r.name,
		Set:      r.set,
		Unset:    r.unset,
		SetTag:   r.setTag,
		UnsetTag: r.unsetTag,
		Rename:   r.rename,
		Drop:     r.drop,
	}
	return opts
}

type ShowAccountRequest struct {
	history *bool
	like    *Like
}

func NewShowAccountRequest() *ShowAccountRequest                     { return &ShowAccountRequest{} }
func (r *ShowAccountRequest) WithLike(v Like) *ShowAccountRequest    { r.like = &v; return r }
func (r *ShowAccountRequest) WithHistory(v bool) *ShowAccountRequest { r.history = &v; return r }
func (r *ShowAccountRequest) toOpts() *ShowAccountOptions {
	opts := &ShowAccountOptions{
		History: r.history,
		Like:    r.like,
	}
	return opts
}

type DropAccountRequest struct {
	name              AccountObjectIdentifier
	ifExists          *bool
	gracePeriodInDays *int
}

func NewDropAccountRequest(name AccountObjectIdentifier) *DropAccountRequest {
	return &DropAccountRequest{name: name}
}
func (r *DropAccountRequest) WithIfExists(v bool) *DropAccountRequest { r.ifExists = &v; return r }
func (r *DropAccountRequest) WithGracePeriodInDays(v int) *DropAccountRequest {
	r.gracePeriodInDays = &v
	return r
}
func (r *DropAccountRequest) toOpts() *DropAccountOptions {
	opts := &DropAccountOptions{
		name:              r.name,
		IfExists:          r.ifExists,
		GracePeriodInDays: r.gracePeriodInDays,
	}
	return opts
}

type UndropAccountRequest struct {
	name AccountObjectIdentifier
}

func NewUndropAccountRequest(name AccountObjectIdentifier) *UndropAccountRequest {
	return &UndropAccountRequest{name: name}
}
func (r *UndropAccountRequest) toOpts() *UndropAccountOptions {
	return &UndropAccountOptions{name: r.name}
}

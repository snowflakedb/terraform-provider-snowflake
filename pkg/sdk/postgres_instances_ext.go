package sdk

func (r *CreatePostgresInstanceRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r postgresInstancesRow) convert() (*PostgresInstance, error) {
	pi := &PostgresInstance{
		Name:                    r.Name,
		Owner:                   r.Owner,
		OwnerRoleType:           r.OwnerRoleType,
		CreatedOn:               r.CreatedOn,
		UpdatedOn:               r.UpdatedOn,
		Type:                    r.Type,
		ComputeFamily:           r.ComputeFamily,
		AuthenticationAuthority: r.AuthenticationAuthority,
		StorageSize:             r.StorageSize,
		PostgresVersion:         r.PostgresVersion,
		IsHa:                    r.IsHa == "true",
		RetentionTime:           r.RetentionTime,
	}
	mapNullString(&pi.Origin, r.Origin)
	mapNullString(&pi.Host, r.Host)
	mapNullString(&pi.PrivatelinkServiceIdentifier, r.PrivatelinkServiceIdentifier)
	mapNullString(&pi.PostgresSettings, r.PostgresSettings)
	mapNullString(&pi.Comment, r.Comment)
	mapStringWithMapping(&pi.State, r.State, ToPostgresInstanceState)
	return pi, nil
}

func (r postgresInstanceDetailsRow) convert() (*PostgresInstanceProperty, error) {
	return &PostgresInstanceProperty{
		Property: r.Property,
		Value:    r.Value.String,
	}, nil
}

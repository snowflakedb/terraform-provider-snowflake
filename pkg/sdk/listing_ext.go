package sdk

func (v *ListingDetails) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

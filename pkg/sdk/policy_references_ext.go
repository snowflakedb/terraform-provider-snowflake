package sdk

func (p PolicyReference) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(*p.PolicyDb, *p.PolicySchema, p.PolicyName)
}

func NewGetForEntityPolicyReferenceRequestCustom(
	refEntityName ObjectIdentifier,
	refEntityDomain PolicyEntityDomain,
) *GetForEntityPolicyReferenceRequest {
	return NewGetForEntityPolicyReferenceRequest(
		NewpolicyReferenceParametersRequest(
			NewpolicyReferenceFunctionArgumentsRequest(
				[]ObjectIdentifier{refEntityName},
				refEntityDomain,
			),
		),
	)
}

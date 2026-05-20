package sdk

func NewGetForEntityTagReferenceRequestFull(
	objectName string,
	objectDomain TagReferenceObjectDomain,
) *GetForEntityTagReferenceRequest {
	return NewGetForEntityTagReferenceRequest(
		NewTagReferenceParametersRequest(
			NewTagReferenceFunctionArgumentsRequest(
				objectName,
				objectDomain,
			),
		),
	)
}

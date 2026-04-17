package sdk

func NewGetForEntityTagReferenceRequestFull(
	objectName string,
	objectDomain TagReferenceObjectDomain,
) *GetForEntityTagReferenceRequest {
	return NewGetForEntityTagReferenceRequest(
		NewtagReferenceParametersRequest(
			NewtagReferenceFunctionArgumentsRequest(
				&objectName,
				&objectDomain,
			),
		),
	)
}

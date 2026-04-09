package sdk

var _ optionsProvider[getForEntityTagReferenceOptions] = new(GetForEntityTagReferenceRequest)

//go:generate go run ./dto-builder-generator/main.go

type GetForEntityTagReferenceRequest struct {
	ObjectName   ObjectIdentifier         // required
	ObjectDomain TagReferenceObjectDomain // required
}

func (request *GetForEntityTagReferenceRequest) toOpts() *getForEntityTagReferenceOptions {
	return &getForEntityTagReferenceOptions{
		parameters: &tagReferenceParameters{
			arguments: &tagReferenceFunctionArguments{
				objectName:   Pointer(request.ObjectName.FullyQualifiedName()),
				objectDomain: Pointer(request.ObjectDomain),
			},
		},
	}
}

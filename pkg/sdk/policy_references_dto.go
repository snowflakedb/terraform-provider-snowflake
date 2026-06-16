package sdk

var _ optionsProvider[getForEntityPolicyReferenceOptions] = new(GetForEntityPolicyReferenceRequest)

//go:generate go run ./dto-builder-generator/main.go

type GetForEntityPolicyReferenceRequest struct {
	RefEntityName   ObjectIdentifier   // required
	RefEntityDomain PolicyEntityDomain // required
}

func (request *GetForEntityPolicyReferenceRequest) toOpts() *getForEntityPolicyReferenceOptions {
	return &getForEntityPolicyReferenceOptions{
		parameters: &policyReferenceParameters{
			arguments: &policyReferenceFunctionArguments{
				refEntityName:   []ObjectIdentifier{request.RefEntityName},
				refEntityDomain: new(request.RefEntityDomain),
			},
		},
	}
}

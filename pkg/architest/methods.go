package architest

type (
	MethodReceiver = func(method *Method)
	Methods        []Method
)

func (methods Methods) All(receiver MethodReceiver) {
	for _, method := range methods {
		method := method
		receiver(&method)
	}
}

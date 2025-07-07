package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type HttpServerEmbeddable[T any] struct {
	serverUrl string
	path      string
}

func (r *HttpServerEmbeddable[T]) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	providerContext, ok := request.ProviderData.(*TestProviderContext)
	if !ok {
		response.Diagnostics.AddError("Provider context is broken", "Set up the context correctly in the provider's Configure func.")
		return
	}

	r.serverUrl = providerContext.ServerUrl()
}

func NewHttpServerEmbeddable[T any](path string) *HttpServerEmbeddable[T] {
	return &HttpServerEmbeddable[T]{
		path: path,
	}
}

func (r *HttpServerEmbeddable[T]) SetPath(path string) {
	r.path = path
}

func (r *HttpServerEmbeddable[T]) Get() (*T, error) {
	var target T
	err := Get(r.serverUrl, r.path, &target)
	return &target, err
}

func (r *HttpServerEmbeddable[T]) Post(target T) error {
	return Post(r.serverUrl, r.path, target)
}

package sdk

import "context"

var (
	_ TagReferences                = new(tagReferenceImpl)
	_ convertibleRow[TagReference] = new(tagReferenceDBRow)
)

type tagReferenceImpl struct {
	client *Client
}

func (v *tagReferenceImpl) GetForEntity(ctx context.Context, request *GetForEntityTagReferenceRequest) ([]TagReference, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[tagReferenceDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[tagReferenceDBRow, TagReference](dbRows)
}

package resourceshowoutputassert

func (p *ImageRepositoryShowOutputAssert) HasCreatedOnNotEmpty() *ImageRepositoryShowOutputAssert {
	p.ValuePresent("created_on")
	return p
}

func (p *ImageRepositoryShowOutputAssert) HasRepositoryUrlNotEmpty() *ImageRepositoryShowOutputAssert {
	p.ValuePresent("repository_url")
	return p
}

func (p *ImageRepositoryShowOutputAssert) HasPrivatelinkRepositoryUrlNotEmpty() *ImageRepositoryShowOutputAssert {
	p.ValuePresent("privatelink_repository_url")
	return p
}

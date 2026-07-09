package resourceshowoutputassert

func (p *ImageRepositoryShowOutputAssert) HasCreatedOnNotEmpty() *ImageRepositoryShowOutputAssert {
	p.ValuePresent("created_on")
	return p
}

func (p *ImageRepositoryShowOutputAssert) HasRepositoryUrlNotEmpty() *ImageRepositoryShowOutputAssert {
	p.ValuePresent("repository_url")
	return p
}

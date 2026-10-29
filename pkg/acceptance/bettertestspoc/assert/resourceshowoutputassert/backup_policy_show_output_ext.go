package resourceshowoutputassert

func (b *BackupPolicyShowOutputAssert) HasCreatedOnNotEmpty() *BackupPolicyShowOutputAssert {
	b.ValuePresent("created_on")
	return b
}

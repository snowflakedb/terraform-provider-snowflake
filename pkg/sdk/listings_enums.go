package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type ListingRevision string

const (
	ListingRevisionDraft     ListingRevision = "DRAFT"
	ListingRevisionPublished ListingRevision = "PUBLISHED"
)

type ListingState string

const (
	ListingStateDraft       ListingState = "DRAFT"
	ListingStatePublished   ListingState = "PUBLISHED"
	ListingStateUnpublished ListingState = "UNPUBLISHED"
)

var AllListingStates = []ListingState{
	ListingStateDraft,
	ListingStatePublished,
	ListingStateUnpublished,
}

func ToListingState(s string) (ListingState, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllListingStates, ListingState(s)) {
		return "", fmt.Errorf("invalid listing state: %s", s)
	}
	return ListingState(s), nil
}

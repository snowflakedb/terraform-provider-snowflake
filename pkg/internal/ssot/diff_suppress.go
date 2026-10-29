package ssot

import "slices"

type DiffSuppress interface {
	Suppressors() []DiffSuppressor
}

type DiffSuppressKind string

const (
	SuppressStatement DiffSuppressKind = "statement"
	SuppressShow      DiffSuppressKind = "show"
	SuppressDescribe  DiffSuppressKind = "describe"
)

type DiffSuppressor struct {
	Kind DiffSuppressKind
	Key  string
}

func (kind DiffSuppressKind) Suppressors() []DiffSuppressor {
	return []DiffSuppressor{{Kind: kind}}
}

func (kind DiffSuppressKind) WithKey(key string) DiffSuppressor {
	return DiffSuppressor{Kind: kind, Key: key}
}

func (suppressor DiffSuppressor) Suppressors() []DiffSuppressor {
	return []DiffSuppressor{suppressor}
}

type diffSuppressors []DiffSuppressor

func (suppressors diffSuppressors) Suppressors() []DiffSuppressor {
	return slices.Clone(suppressors)
}

func SuppressIfAny(suppressors ...DiffSuppress) DiffSuppress {
	var combined diffSuppressors
	for _, suppressor := range suppressors {
		if suppressor != nil {
			combined = append(combined, suppressor.Suppressors()...)
		}
	}
	if len(combined) == 0 {
		return nil
	}
	return combined
}

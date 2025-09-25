package rework

// This file's only purpose is to make generated objects compile (or close to compile).
// Later code will be generated inside sdk package, so the objects will be accessible there.

type Client struct{}

type (
	ObjectIdentifier         interface{}
	AccountObjectIdentifier  struct{}
	DatabaseObjectIdentifier struct{}
	ExternalObjectIdentifier struct{}
	SchemaObjectIdentifier   struct{}
	TableColumnIdentifier    struct{}

	In   struct{}
	Like struct{}
)

type ValuesBehavior string

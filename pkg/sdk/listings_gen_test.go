package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func sampleListingManifest() string {
	return `
title: "MyListing"
subtitle: "Subtitle for MyListing"
description: "Description for MyListing"
listing_terms:
   type: "STANDARD"
targets:
    accounts: ["Org1.Account1"]
usage_examples:
    - title: "this is a test sql"
      description: "Simple example"
      query: "select *"
`
}

func TestListings_Create(t *testing.T) {
	id := randomAccountObjectIdentifier()
	stageId := randomSchemaObjectIdentifier()
	manifest := sampleListingManifest()
	var stageLocation Location = &StageLocation{
		stage: stageId,
		path:  "dir/subdir",
	}

	// Minimal valid CreateListingOptions
	defaultOpts := func() *CreateListingOptions {
		return &CreateListingOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateListingOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = invalidAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.As opts.From] should be present - none set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateListingOptions", "As", "From"))
	})

	t.Run("validation: exactly one field from [opts.As opts.From] should be present - two set", func(t *testing.T) {
		opts := defaultOpts()
		opts.As = &manifest
		opts.From = &stageLocation
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateListingOptions", "As", "From"))
	})

	t.Run("validation: exactly one field from [opts.With.Share opts.With.ApplicationPackage] should be present - none set", func(t *testing.T) {
		opts := defaultOpts()
		opts.With = &ListingWith{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateListingOptions.With", "Share", "ApplicationPackage"))
	})

	t.Run("validation: exactly one field from [opts.With.Share opts.With.ApplicationPackage] should be present - two set", func(t *testing.T) {
		shareId := randomAccountObjectIdentifier()
		applicationPackageId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.With = &ListingWith{
			Share:              &shareId,
			ApplicationPackage: &applicationPackageId,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateListingOptions.With", "Share", "ApplicationPackage"))
	})

	t.Run("basic with As", func(t *testing.T) {
		opts := defaultOpts()
		opts.As = &manifest
		assertOptsValidAndSQLEquals(t, opts, "CREATE EXTERNAL LISTING %s AS $$%s$$", opts.name.FullyQualifiedName(), manifest)
	})

	t.Run("basic with From", func(t *testing.T) {
		opts := defaultOpts()
		opts.From = &stageLocation
		assertOptsValidAndSQLEquals(t, opts, "CREATE EXTERNAL LISTING %s FROM @%s/dir/subdir", opts.name.FullyQualifiedName(), stageId.FullyQualifiedName())
	})

	t.Run("all As options with stage", func(t *testing.T) {
		shareId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.With = &ListingWith{
			Share: &shareId,
		}
		opts.As = &manifest
		opts.Publish = Bool(true)
		opts.Review = Bool(true)
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE EXTERNAL LISTING IF NOT EXISTS %s SHARE %s AS $$%s$$ PUBLISH = true REVIEW = true COMMENT = 'comment'", opts.name.FullyQualifiedName(), shareId.FullyQualifiedName(), manifest)
	})

	t.Run("all As options with application package", func(t *testing.T) {
		applicationPackageId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.With = &ListingWith{
			ApplicationPackage: &applicationPackageId,
		}
		opts.As = &manifest
		opts.Publish = Bool(true)
		opts.Review = Bool(true)
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE EXTERNAL LISTING IF NOT EXISTS %s APPLICATION PACKAGE %s AS $$%s$$ PUBLISH = true REVIEW = true COMMENT = 'comment'", opts.name.FullyQualifiedName(), applicationPackageId.FullyQualifiedName(), manifest)
	})

	t.Run("all From options with stage", func(t *testing.T) {
		shareId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.With = &ListingWith{
			Share: &shareId,
		}
		opts.From = &stageLocation
		opts.Publish = Bool(true)
		opts.Review = Bool(true)
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE EXTERNAL LISTING IF NOT EXISTS %s SHARE %s FROM @%s/dir/subdir PUBLISH = true REVIEW = true COMMENT = 'comment'", opts.name.FullyQualifiedName(), shareId.FullyQualifiedName(), stageId.FullyQualifiedName())
	})

	t.Run("all From options with application package", func(t *testing.T) {
		applicationPackageId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.With = &ListingWith{
			ApplicationPackage: &applicationPackageId,
		}
		opts.From = &stageLocation
		opts.Publish = Bool(true)
		opts.Review = Bool(true)
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE EXTERNAL LISTING IF NOT EXISTS %s APPLICATION PACKAGE %s FROM @%s/dir/subdir PUBLISH = true REVIEW = true COMMENT = 'comment'", opts.name.FullyQualifiedName(), applicationPackageId.FullyQualifiedName(), stageId.FullyQualifiedName())
	})
}

func TestListings_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier()
	manifest := sampleListingManifest()

	// Minimal valid AlterListingOptions
	defaultOpts := func() *AlterListingOptions {
		return &AlterListingOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterListingOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = invalidAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfExists opts.AddVersion]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.AddVersion = &AddListingVersion{}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterListingOptions", "IfExists", "AddVersion"))
	})

	t.Run("validation: exactly one field from [opts.Publish opts.Unpublish opts.Review opts.AlterListingAs opts.AddVersion opts.RenameTo opts.Set] should be present - none set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterListingOptions", "Publish", "Unpublish", "Review", "AlterListingAs", "AddVersion", "RenameTo", "Set", "Unset"))
	})

	t.Run("validation: exactly one field from [opts.Publish opts.Unpublish opts.Review opts.AlterListingAs opts.AddVersion opts.RenameTo opts.Set] should be present - two set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Publish = Bool(true)
		opts.Unpublish = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterListingOptions", "Publish", "Unpublish", "Review", "AlterListingAs", "AddVersion", "RenameTo", "Set", "Unset"))
	})

	t.Run("publish", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Publish = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER LISTING IF EXISTS %s PUBLISH", opts.name.FullyQualifiedName())
	})

	t.Run("unpublish", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Unpublish = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER LISTING IF EXISTS %s UNPUBLISH", opts.name.FullyQualifiedName())
	})

	t.Run("review", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Review = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER LISTING IF EXISTS %s REVIEW", opts.name.FullyQualifiedName())
	})

	t.Run("as: basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.AlterListingAs = &AlterListingAs{
			As: manifest,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER LISTING IF EXISTS %s AS $$%s$$", opts.name.FullyQualifiedName(), manifest)
	})

	t.Run("as: complete", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.AlterListingAs = &AlterListingAs{
			As:      manifest,
			Publish: Bool(true),
			Review:  Bool(true),
			Comment: String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER LISTING IF EXISTS %s AS $$%s$$ PUBLISH = true REVIEW = true COMMENT = 'comment'", opts.name.FullyQualifiedName(), manifest)
	})

	t.Run("add version", func(t *testing.T) {
		stageId := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.AddVersion = &AddListingVersion{
			IfNotExists: Bool(true),
			VersionName: "version-name",
			From: StageLocation{
				stage: stageId,
				path:  "dir/subdir",
			},
			Comment: String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER LISTING %s ADD VERSION IF NOT EXISTS version-name FROM @%s/dir/subdir COMMENT = 'comment'", opts.name.FullyQualifiedName(), stageId.FullyQualifiedName())
	})

	t.Run("rename to", func(t *testing.T) {
		newId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.RenameTo = &newId
		assertOptsValidAndSQLEquals(t, opts, "ALTER LISTING %s RENAME TO %s", opts.name.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ListingSet{
			Comment: String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER LISTING %s SET COMMENT = 'comment'", opts.name.FullyQualifiedName())
	})
}

func TestListings_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid DropListingOptions
	defaultOpts := func() *DropListingOptions {
		return &DropListingOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropListingOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = invalidAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP LISTING IF EXISTS %s", opts.name.FullyQualifiedName())
	})
}

func TestListings_Show(t *testing.T) {
	// Minimal valid ShowListingOptions
	defaultOpts := func() *ShowListingOptions {
		return &ShowListingOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowListingOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW LISTINGS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("pattern"),
		}
		opts.StartsWith = String("startsWith")
		opts.Limit = &LimitFrom{
			Rows: Int(10),
			From: String("from"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW LISTINGS LIKE 'pattern' STARTS WITH 'startsWith' LIMIT 10 FROM 'from'")
	})
}

func TestListings_ShowVersions(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid ShowVersionsListingOptions
	defaultOpts := func() *ShowVersionsListingOptions {
		return &ShowVersionsListingOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowVersionsListingOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("invalid identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = invalidAccountObjectIdentifier
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW VERSIONS IN LISTING %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Limit = &LimitFrom{
			Rows: Int(5),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW VERSIONS IN LISTING %s LIMIT 5", id.FullyQualifiedName())
	})
}

func TestListings_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid DescribeListingOptions
	defaultOpts := func() *DescribeListingOptions {
		return &DescribeListingOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeListingOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = invalidAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE LISTING %s", opts.name.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Revision = Pointer(ListingRevisionDraft)
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE LISTING %s REVISION = DRAFT", opts.name.FullyQualifiedName())
	})
}

func Test_Listings_ToListingState(t *testing.T) {
	type test struct {
		input string
		want  ListingState
	}

	valid := []test{
		{input: "draft", want: ListingStateDraft},
		{input: "published", want: ListingStatePublished},
		{input: "unpublished", want: ListingStateUnpublished},
		{input: "DRAft", want: ListingStateDraft},
		{input: "drAFT", want: ListingStateDraft},
		{input: "DRAFT", want: ListingStateDraft},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToListingState(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToListingState(tc.input)
			require.Error(t, err)
		})
	}
}

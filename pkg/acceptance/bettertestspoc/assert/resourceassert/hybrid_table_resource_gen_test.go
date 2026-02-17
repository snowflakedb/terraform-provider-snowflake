package resourceassert

import (
	"testing"
)

// TestHybridTableResourceAssert_FactoryFunctions validates factory functions
func TestHybridTableResourceAssert_FactoryFunctions(t *testing.T) {
	// Test resource factory
	assert := HybridTableResource(t, "snowflake_hybrid_table.test")
	if assert == nil {
		t.Fatal("HybridTableResource returned nil")
	}
	if assert.ResourceAssert == nil {
		t.Fatal("ResourceAssert not initialized")
	}

	// Test imported resource factory
	importedAssert := ImportedHybridTableResource(t, "test_db|test_schema|test_table")
	if importedAssert == nil {
		t.Fatal("ImportedHybridTableResource returned nil")
	}
	if importedAssert.ResourceAssert == nil {
		t.Fatal("ResourceAssert not initialized for imported resource")
	}
}

// TestHybridTableResourceAssert_MethodChaining validates fluent interface
func TestHybridTableResourceAssert_MethodChaining(t *testing.T) {
	assert := HybridTableResource(t, "snowflake_hybrid_table.test")

	// Verify all methods return the assert object for chaining
	result := assert.
		HasDatabaseString("TEST_DB").
		HasSchemaString("TEST_SCHEMA").
		HasNameString("TEST_TABLE").
		HasCommentString("test comment").
		HasDataRetentionTimeInDaysString("7").
		HasColumnCount(3).
		HasIndexCount(1).
		HasPrimaryKeyNotEmpty().
		HasShowOutputNotEmpty()

	if result != assert {
		t.Fatal("Method chaining broken - methods don't return self")
	}
}

// TestHybridTableResourceAssert_AllStringMethods validates all Has*String methods exist
func TestHybridTableResourceAssert_AllStringMethods(t *testing.T) {
	assert := HybridTableResource(t, "snowflake_hybrid_table.test")

	// Test all string setters
	assert.HasDatabaseString("db")
	assert.HasSchemaString("schema")
	assert.HasNameString("name")
	assert.HasOrReplaceString("true")
	assert.HasDataRetentionTimeInDaysString("7")
	assert.HasCommentString("comment")
	assert.HasFullyQualifiedNameString("db.schema.name")
}

// TestHybridTableResourceAssert_AllNoValueMethods validates all HasNo* methods exist
func TestHybridTableResourceAssert_AllNoValueMethods(t *testing.T) {
	assert := HybridTableResource(t, "snowflake_hybrid_table.test")

	// Test all no value checks
	assert.HasNoDatabase()
	assert.HasNoSchema()
	assert.HasNoName()
	assert.HasNoOrReplace()
	assert.HasNoDataRetentionTimeInDays()
	assert.HasNoComment()
}

// TestHybridTableResourceAssert_AllEmptyMethods validates all Has*Empty methods exist
func TestHybridTableResourceAssert_AllEmptyMethods(t *testing.T) {
	assert := HybridTableResource(t, "snowflake_hybrid_table.test")

	// Test all empty checks
	assert.HasCommentEmpty()
	assert.HasFullyQualifiedNameEmpty()
	assert.HasColumnEmpty()
	assert.HasIndexEmpty()
	assert.HasPrimaryKeyEmpty()
	assert.HasUniqueConstraintEmpty()
	assert.HasForeignKeyEmpty()
	assert.HasShowOutputEmpty()
	assert.HasDescribeOutputEmpty()
}

// TestHybridTableResourceAssert_AllNotEmptyMethods validates all Has*NotEmpty methods exist
func TestHybridTableResourceAssert_AllNotEmptyMethods(t *testing.T) {
	assert := HybridTableResource(t, "snowflake_hybrid_table.test")

	// Test all not empty checks
	assert.HasDatabaseNotEmpty()
	assert.HasSchemaNotEmpty()
	assert.HasNameNotEmpty()
	assert.HasCommentNotEmpty()
	assert.HasFullyQualifiedNameNotEmpty()
	assert.HasDataRetentionTimeInDaysNotEmpty()
	assert.HasColumnNotEmpty()
	assert.HasIndexNotEmpty()
	assert.HasPrimaryKeyNotEmpty()
	assert.HasUniqueConstraintNotEmpty()
	assert.HasForeignKeyNotEmpty()
	assert.HasShowOutputNotEmpty()
	assert.HasDescribeOutputNotEmpty()
}

// TestHybridTableResourceAssert_AllCountMethods validates all count methods exist
func TestHybridTableResourceAssert_AllCountMethods(t *testing.T) {
	assert := HybridTableResource(t, "snowflake_hybrid_table.test")

	// Test all count checks
	assert.HasColumnCount(3)
	assert.HasIndexCount(2)
	assert.HasUniqueConstraintCount(1)
	assert.HasForeignKeyCount(1)
}

// TestHybridTableResourceAssert_ColumnAttributeMethods validates column attribute checks
func TestHybridTableResourceAssert_ColumnAttributeMethods(t *testing.T) {
	assert := HybridTableResource(t, "snowflake_hybrid_table.test")

	// Test column attribute checks
	assert.HasColumnName(0, "id")
	assert.HasColumnType(0, "NUMBER(38,0)")
	assert.HasColumnNullable(0, "false")
	assert.HasColumnPrimaryKey(0, "true")
	assert.HasColumnUnique(0, "false")
	assert.HasColumnComment(0, "ID column")
	assert.HasColumnCollate(0, "en_US")
}

// TestHybridTableResourceAssert_IndexAttributeMethods validates index attribute checks
func TestHybridTableResourceAssert_IndexAttributeMethods(t *testing.T) {
	assert := HybridTableResource(t, "snowflake_hybrid_table.test")

	// Test index attribute checks
	assert.HasIndexName(0, "idx_name")
}

// TestHybridTableResourceAssert_PrimaryKeyAttributeMethods validates primary key attribute checks
func TestHybridTableResourceAssert_PrimaryKeyAttributeMethods(t *testing.T) {
	assert := HybridTableResource(t, "snowflake_hybrid_table.test")

	// Test primary key attribute checks
	assert.HasPrimaryKeyName("pk_test")
}

// TestHybridTableResourceAssert_UniqueConstraintAttributeMethods validates unique constraint checks
func TestHybridTableResourceAssert_UniqueConstraintAttributeMethods(t *testing.T) {
	assert := HybridTableResource(t, "snowflake_hybrid_table.test")

	// Test unique constraint attribute checks
	assert.HasUniqueConstraintName(0, "uq_email")
}

// TestHybridTableResourceAssert_ForeignKeyAttributeMethods validates foreign key checks
func TestHybridTableResourceAssert_ForeignKeyAttributeMethods(t *testing.T) {
	assert := HybridTableResource(t, "snowflake_hybrid_table.test")

	// Test foreign key attribute checks
	assert.HasForeignKeyName(0, "fk_customer")
	assert.HasForeignKeyReferencesTable(0, "customers")
}

#!/bin/bash
# =============================================================================
# Migration Script Manual Test Runner
# =============================================================================
# This script orchestrates the full testing workflow for any object type:
#   1. Create test objects on Snowflake
#   2. Fetch objects via data source and generate CSV
#   3. Run migration script and compare output
#
# Usage:
#   ./run_test.sh <object_type>                    # Run test (create, test, keep resources)
#   ./run_test.sh <object_type> --destroy          # Only destroy existing resources
#   ./run_test.sh <object_type> --skip-create      # Skip creation, just run migration test
#   ./run_test.sh --list                           # List available object types
#
# Examples:
#   ./run_test.sh users
#   ./run_test.sh account_roles --skip-create
#   ./run_test.sh users --destroy
#
# Supported object types (must have a folder with test files):
#   - users
# =============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATION_SCRIPT_DIR="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

log_step() {
    echo -e "\n${BLUE}==>${NC} ${GREEN}$1${NC}"
}

log_info() {
    echo -e "${YELLOW}    $1${NC}"
}

log_error() {
    echo -e "${RED}ERROR: $1${NC}"
}

log_header() {
    echo -e "\n${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}  $1${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# List available object types
list_object_types() {
    echo "Available object types:"
    echo ""
    for dir in "$SCRIPT_DIR"/*/; do
        if [ -d "$dir" ] && [ -f "$dir/objects_def.tf" ]; then
            object_type=$(basename "$dir")
            echo "  - $object_type"
        fi
    done
    echo ""
    echo "To add a new object type, create a folder with:"
    echo "  - objects_def.tf      (creates test objects)"
    echo "  - datasource.tf       (fetches objects, generates CSV)"
    echo "  - expected_output.tf  (expected migration script output)"
}

# Parse arguments
OBJECT_TYPE=""
DESTROY_ONLY=false
SKIP_CREATE=false

for arg in "$@"; do
    case $arg in
        --destroy)
            DESTROY_ONLY=true
            ;;
        --skip-create)
            SKIP_CREATE=true
            ;;
        --list)
            list_object_types
            exit 0
            ;;
        --help|-h)
            echo "Usage: $0 <object_type> [OPTIONS]"
            echo ""
            echo "Arguments:"
            echo "  object_type     The type of object to test (e.g., users, account_roles)"
            echo ""
            echo "Options:"
            echo "  --destroy       Only destroy existing resources (no test)"
            echo "  --skip-create   Skip creation, only run migration test"
            echo "  --list          List available object types"
            echo "  --help, -h      Show this help"
            echo ""
            echo "Examples:"
            echo "  $0 users                    # Test users"
            echo "  $0 account_roles --skip-create  # Test with existing objects"
            echo "  $0 users --destroy          # Destroy test resources"
            echo "  $0 --list                   # Show available object types"
            exit 0
            ;;
        -*)
            log_error "Unknown option: $arg"
            exit 1
            ;;
        *)
            if [ -z "$OBJECT_TYPE" ]; then
                OBJECT_TYPE="$arg"
            fi
            ;;
    esac
done

# Validate object type
if [ -z "$OBJECT_TYPE" ]; then
    log_error "Object type is required. Use --list to see available types."
    echo ""
    echo "Usage: $0 <object_type> [OPTIONS]"
    exit 1
fi

OBJECT_DIR="$SCRIPT_DIR/$OBJECT_TYPE"

if [ ! -d "$OBJECT_DIR" ]; then
    log_error "Object type '$OBJECT_TYPE' not found."
    echo ""
    list_object_types
    exit 1
fi

# Check required files exist
if [ ! -f "$OBJECT_DIR/objects_def.tf" ]; then
    log_error "Missing $OBJECT_DIR/objects_def.tf"
    exit 1
fi

if [ ! -f "$OBJECT_DIR/datasource.tf" ]; then
    log_error "Missing $OBJECT_DIR/datasource.tf"
    exit 1
fi

if [ ! -f "$OBJECT_DIR/expected_output.tf" ]; then
    log_error "Missing $OBJECT_DIR/expected_output.tf"
    exit 1
fi

log_header "Testing: $OBJECT_TYPE"

cd "$OBJECT_DIR"

# =============================================================================
# Remove previous actual_output.tf to avoid terraform errors
# =============================================================================
if [ -f "actual_output.tf" ]; then
    rm -f actual_output.tf
    log_info "Removed previous actual_output.tf"
fi

# =============================================================================
# DESTROY ONLY MODE
# =============================================================================
if [ "$DESTROY_ONLY" = true ]; then
    log_step "Destroying all test resources for $OBJECT_TYPE..."
    terraform destroy -auto-approve
    log_step "Cleanup complete!"
    exit 0
fi

# =============================================================================
# STEP 0: Initialize Terraform (if needed)
# =============================================================================
if [ ! -d ".terraform" ]; then
    log_step "Step 0: Initializing Terraform..."
    terraform init
fi

# =============================================================================
# STEP 1: Create test objects on Snowflake
# =============================================================================
if [ "$SKIP_CREATE" = false ]; then
    log_step "Step 1: Creating test objects on Snowflake..."

    # Get all resource names from objects_def.tf and apply them
    RESOURCES=$(grep -E '^resource\s+"' objects_def.tf | sed 's/resource "\([^"]*\)" "\([^"]*\)".*/\1.\2/' | tr '\n' ' ')

    if [ -n "$RESOURCES" ]; then
        log_info "Found resources to create"

        # Build -target arguments
        TARGET_ARGS=""
        for resource in $RESOURCES; do
            TARGET_ARGS="$TARGET_ARGS -target=$resource"
        done

        export SF_TF_ACC_TEST_ENABLE_ALL_PREVIEW_FEATURES=true
        terraform apply -auto-approve $TARGET_ARGS
        log_info "Test objects created successfully!"
    else
        log_error "No resources found in objects_def.tf"
        exit 1
    fi
else
    log_step "Step 1: Skipping object creation (--skip-create)"
fi

# =============================================================================
# STEP 2: Fetch objects via data source and generate CSV
# =============================================================================
log_step "Step 2: Fetching objects via data source and generating CSV..."
log_info "This retrieves object data including DESCRIBE and PARAMETERS output (if applicable)"

# Get data sources and local_file resources from datasource.tf only
# This avoids re-applying objects_def.tf resources
DATASOURCE_TARGETS=""
if [ -f "datasource.tf" ]; then
    # Get data sources
    DATA_SOURCES=$(grep -E '^data\s+"' datasource.tf | sed 's/data "\([^"]*\)" "\([^"]*\)".*/data.\1.\2/' | tr '\n' ' ')
    for ds in $DATA_SOURCES; do
        DATASOURCE_TARGETS="$DATASOURCE_TARGETS -target=$ds"
    done

    # Get local_file resources
    LOCAL_FILES=$(grep -E '^resource\s+"local_file"' datasource.tf | sed 's/resource "local_file" "\([^"]*\)".*/local_file.\1/' | tr '\n' ' ')
    for lf in $LOCAL_FILES; do
        DATASOURCE_TARGETS="$DATASOURCE_TARGETS -target=$lf"
    done
fi

# Apply only the data source and CSV generation
export SF_TF_ACC_TEST_ENABLE_ALL_PREVIEW_FEATURES=true
terraform apply -auto-approve $DATASOURCE_TARGETS

# Check if CSV was generated
if [ -f "objects.csv" ]; then
    LINES=$(wc -l < objects.csv | tr -d ' ')
    log_info "CSV generated with $LINES lines: $OBJECT_DIR/objects.csv"
else
    log_error "CSV file was not generated!"
    exit 1
fi

# =============================================================================
# STEP 3: Run migration script and compare output
# =============================================================================
log_step "Step 3: Running migration script..."

cd "$MIGRATION_SCRIPT_DIR"
go run . -import=block "$OBJECT_TYPE" < "$OBJECT_DIR/objects.csv" > "$OBJECT_DIR/actual_output.tf"

log_info "Migration script output saved to: $OBJECT_DIR/actual_output.tf"

cd "$OBJECT_DIR"

# =============================================================================
# STEP 4: Compare with expected output
# =============================================================================
log_step "Step 4: Comparing actual output with expected output..."

# Check if files exist
if [ ! -f "expected_output.tf" ]; then
    log_error "Expected output file not found: expected_output.tf"
    exit 1
fi

if [ ! -f "actual_output.tf" ]; then
    log_error "Actual output file not found: actual_output.tf"
    exit 1
fi

# Normalize and compare (remove comments and empty lines)
grep -v '^#' expected_output.tf | grep -v '^$' > /tmp/expected_norm.txt
grep -v '^#' actual_output.tf | grep -v '^$' > /tmp/actual_norm.txt

echo ""
echo "Differences (if any):"
echo "====================="

DIFF_RESULT=0
if diff -u /tmp/expected_norm.txt /tmp/actual_norm.txt > /tmp/migration_diff.txt 2>&1; then
    echo -e "${GREEN}✓ No significant differences found!${NC}"
else
    cat /tmp/migration_diff.txt
    echo ""
    echo -e "${YELLOW}Note: Some differences may be expected due to Snowflake defaults.${NC}"
    DIFF_RESULT=1
fi

# Show resource counts
EXPECTED_RESOURCES=$(grep -c '^resource' expected_output.tf 2>/dev/null || echo 0)
ACTUAL_RESOURCES=$(grep -c '^resource' actual_output.tf 2>/dev/null || echo 0)
EXPECTED_IMPORTS=$(grep -c '^import' expected_output.tf 2>/dev/null || echo 0)
ACTUAL_IMPORTS=$(grep -c '^import' actual_output.tf 2>/dev/null || echo 0)

echo ""
echo "Summary:"
echo "  Expected: $EXPECTED_RESOURCES resources, $EXPECTED_IMPORTS imports"
echo "  Actual:   $ACTUAL_RESOURCES resources, $ACTUAL_IMPORTS imports"

# =============================================================================
# DONE
# =============================================================================
echo ""
echo -e "${YELLOW}Resources were NOT destroyed.${NC}"
echo "To cleanup, run: $0 $OBJECT_TYPE --destroy"
echo ""

if [ $DIFF_RESULT -eq 0 ]; then
    log_step "Test completed successfully! ✓"
else
    log_step "Test completed with differences. Review the output above."
    exit 1
fi

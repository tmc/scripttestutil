#!/bin/bash
# Script to run all scripttest self-tests

set -e  # Exit immediately if a command exits with a non-zero status

# Determine script location and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
SCRIPTTEST_BIN="$PROJECT_ROOT/scripttesttool"

echo "Working from directory: $PROJECT_ROOT"

# Build scripttest if needed
if [ ! -f "$SCRIPTTEST_BIN" ] || [ "$1" == "--rebuild" ]; then
    echo "Building scripttest..."
    cd "$PROJECT_ROOT"
    go build -o scripttesttool ./cmd/scripttest
fi

# Function to run a test and report result
run_test() {
    local test_file="$1"
    local test_name=$(basename "$test_file" .txt)
    
    echo "Running $test_name..."
    if "$SCRIPTTEST_BIN" test "$test_file"; then
        echo "✅ $test_name passed"
        return 0
    else
        echo "❌ $test_name failed"
        return 1
    fi
}

# Run each test
cd "$PROJECT_ROOT"

# Basic self-test
run_test "$SCRIPT_DIR/basic_self_test.txt"

# Snapshot test (first update snapshots, then test against them)
echo "Creating snapshots for snapshot tests..."
UPDATE_SNAPSHOTS=1 "$SCRIPTTEST_BIN" test "$SCRIPT_DIR/snapshot_test.txt"
run_test "$SCRIPT_DIR/snapshot_test.txt"

# Auto-toolchain test
run_test "$SCRIPT_DIR/auto_toolchain_test.txt"

# Asciicast test
run_test "$SCRIPT_DIR/asciicast_test.txt"

# Meta test (run_self_tests.txt)
run_test "$SCRIPT_DIR/run_self_tests.txt"

# Docker test (if available)
if command -v docker &> /dev/null; then
    echo "Docker is available. Running Docker tests..."
    if "$SCRIPTTEST_BIN" -docker test "$SCRIPT_DIR/docker_self_test.txt"; then
        echo "✅ Docker test passed"
    else
        echo "❌ Docker test failed"
        exit 1
    fi
else
    echo "⚠️ Docker is not available. Skipping Docker tests."
fi

echo "🎉 All self-tests complete!"
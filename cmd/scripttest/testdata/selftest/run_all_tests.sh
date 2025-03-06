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

echo "Running basic self-test..."
run_test "$SCRIPT_DIR/basic_self_test.txt"

echo "Running auto-toolchain test..."
run_test "$SCRIPT_DIR/auto_toolchain_test.txt"

echo "Running asciicast test..."
run_test "$SCRIPT_DIR/asciicast_test.txt"

echo "Running self-tests meta-test..."
run_test "$SCRIPT_DIR/run_self_tests.txt"

# We skip running the snapshot test since it's difficult to make persistent in the test environment
echo "Skipping snapshot tests as they require persistent storage between test runs"

# We also skip Docker tests since they require Docker to be running
echo "⚠️ Skipping Docker tests for this run."

echo "🎉 All self-tests complete!"
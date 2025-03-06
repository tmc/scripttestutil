package bridge_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tmc/scripttestutil/bridge"
)

// TestBasicUsage demonstrates the basic usage of the bridge package.
func TestBasicUsage(t *testing.T) {
	// Skip this test when running go test
	// Remove this line when you want to actually run the test
	t.Skip("Example test - skipped by default")

	// Setup test options
	opts := bridge.DefaultOptions()
	opts.Verbose = testing.Verbose()

	// Run a single test file
	bridge.RunFile(t, "testdata/basic.txt", opts)
}

// TestWithCustomOptions demonstrates using custom options.
func TestWithCustomOptions(t *testing.T) {
	// Skip this test when running go test
	t.Skip("Example test - skipped by default")

	// Setup custom options
	opts := bridge.DefaultOptions()
	opts.Verbose = true
	opts.UpdateSnapshots = true
	opts.EnvVars = map[string]string{
		"CUSTOM_ENV": "test-value",
		"DEBUG":      "true",
	}

	// Run all tests in a directory
	bridge.RunDir(t, "testdata", opts)
}

// TestPatternMatching demonstrates using pattern matching for test files.
func TestPatternMatching(t *testing.T) {
	// Skip this test when running go test
	t.Skip("Example test - skipped by default")

	// Run all tests matching a pattern
	opts := bridge.DefaultOptions()
	bridge.RunPattern(t, "testdata/*.txt", opts)
}

// TestWithSnapshots demonstrates working with snapshots.
func TestWithSnapshots(t *testing.T) {
	// Skip this test when running go test
	t.Skip("Example test - skipped by default")

	// Setup options for snapshot testing
	opts := bridge.DefaultOptions()
	opts.SnapshotDir = "testdata/snapshots"

	// First run: create/update snapshots
	updateOpts := opts
	updateOpts.UpdateSnapshots = true
	bridge.RunDir(t, "testdata", updateOpts)

	// Second run: verify against snapshots
	bridge.RunDir(t, "testdata", opts)
}

// TestRunnerAPI demonstrates using the Runner API directly.
func TestRunnerAPI(t *testing.T) {
	// Skip this test when running go test
	t.Skip("Example test - skipped by default")

	// Create a runner with custom options
	opts := bridge.DefaultOptions()
	opts.Pattern = "testdata/*.txt"
	opts.Verbose = true

	// Create runner and run tests
	runner := bridge.NewRunner(opts)
	runner.Run(t)
}

// ExampleSetupTestFiles shows how to programmatically create test files.
func ExampleSetupTestFiles() {
	// Create test directory
	os.MkdirAll("testdata", 0755)

	// Create a simple test file
	testContent := `# Simple test
echo "Hello, World!"
stdout 'Hello, World!'
! stderr .
`
	os.WriteFile(filepath.Join("testdata", "simple.txt"), []byte(testContent), 0644)

	// Create a test file with environment variables
	envTestContent := `# Environment variable test
env TEST_VAR=test_value
env | grep TEST_VAR
stdout 'TEST_VAR=test_value'
`
	os.WriteFile(filepath.Join("testdata", "env.txt"), []byte(envTestContent), 0644)
}
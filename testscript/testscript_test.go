package testscript_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tmc/scripttestutil/testscript"
)

// TestBasicUsage demonstrates the basic usage of the testscript package.
func TestBasicUsage(t *testing.T) {
	// Skip this test when running go test
	// Remove this line when you want to actually run the test
	t.Skip("Example test - skipped by default")

	// Setup test options
	opts := testscript.DefaultOptions()
	opts.Verbose = testing.Verbose()

	// Run a single test file
	testscript.RunFile(t, "testdata/basic.txt", opts)
}

// TestWithCustomOptions demonstrates using custom options.
func TestWithCustomOptions(t *testing.T) {
	// Skip this test when running go test
	t.Skip("Example test - skipped by default")

	// Setup custom options
	opts := testscript.DefaultOptions()
	opts.Verbose = true
	opts.UpdateSnapshots = true
	opts.EnvVars = map[string]string{
		"CUSTOM_ENV": "test-value",
		"DEBUG":      "true",
	}

	// Run all tests in a directory
	testscript.RunDir(t, "testdata", opts)
}

// TestPatternMatching demonstrates using pattern matching for test files.
func TestPatternMatching(t *testing.T) {
	// Skip this test when running go test
	t.Skip("Example test - skipped by default")

	// Run all tests matching a pattern
	opts := testscript.DefaultOptions()
	testscript.Run(t, "testdata/*.txt", opts)
}

// TestWithSnapshots demonstrates working with snapshots.
func TestWithSnapshots(t *testing.T) {
	// Skip this test when running go test
	t.Skip("Example test - skipped by default")

	// Setup options for snapshot testing
	opts := testscript.DefaultOptions()
	opts.SnapshotDir = "testdata/snapshots"

	// First run: create/update snapshots
	updateOpts := opts
	updateOpts.UpdateSnapshots = true
	testscript.RunDir(t, "testdata", updateOpts)

	// Second run: verify against snapshots
	testscript.RunDir(t, "testdata", opts)
}

// TestRunnerAPI demonstrates using the Runner API directly.
func TestRunnerAPI(t *testing.T) {
	// Skip this test when running go test
	t.Skip("Example test - skipped by default")

	// Create a runner with custom options
	opts := testscript.DefaultOptions()
	opts.Pattern = "testdata/*.txt"
	opts.Verbose = true

	// Create runner and run tests
	runner := testscript.NewRunner(opts)
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
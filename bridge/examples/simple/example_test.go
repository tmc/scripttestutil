package simple_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tmc/scripttestutil/bridge"
)

func TestMain(m *testing.M) {
	// Set up test environment
	setupTests()
	// Run tests and get exit code
	code := m.Run()
	// Clean up if needed
	os.Exit(code)
}

// setupTests creates test files for demonstration
func setupTests() {
	// Create test directory if it doesn't exist
	os.MkdirAll("testdata", 0755)

	// Create a basic test file
	basicTest := `# Basic echo test
echo "Hello from scripttest bridge!"
stdout 'Hello from scripttest bridge!'
! stderr .

# Test with environment variable
env BRIDGE_TEST=working
env | grep BRIDGE_TEST
stdout 'BRIDGE_TEST=working'
`
	os.WriteFile(filepath.Join("testdata", "basic.txt"), []byte(basicTest), 0644)

	// Create a file test
	fileTest := `# Test file operations
echo "File content" > test_file.txt
cat test_file.txt
stdout 'File content'
exists test_file.txt
`
	os.WriteFile(filepath.Join("testdata", "file_test.txt"), []byte(fileTest), 0644)
}

// TestWithBridge demonstrates how to use the bridge package
func TestWithBridge(t *testing.T) {
	// Create options with verbose output if go test -v is used
	opts := bridge.DefaultOptions()
	opts.Verbose = testing.Verbose()

	// Set some environment variables for the tests
	opts.EnvVars = map[string]string{
		"EXAMPLE_ENV": "example_value",
	}

	// Run all tests in the testdata directory
	bridge.RunDir(t, "testdata", opts)
}

// TestSingleFile demonstrates running a single test file
func TestSingleFile(t *testing.T) {
	opts := bridge.DefaultOptions()
	bridge.RunFile(t, "testdata/basic.txt", opts)
}

// TestPattern demonstrates running tests matching a pattern
func TestPattern(t *testing.T) {
	opts := bridge.DefaultOptions()
	bridge.RunPattern(t, "testdata/*_test.txt", opts)
}
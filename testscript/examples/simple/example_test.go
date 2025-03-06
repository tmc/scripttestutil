package simple_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tmc/scripttestutil/testscript"
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
echo "Hello from testscript!"
stdout 'Hello from testscript!'
! stderr .

# Test with environment variable
env TST_VAR=working
env | grep TST_VAR
stdout 'TST_VAR=working'
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

// TestWithTestscript demonstrates how to use the testscript package
func TestWithTestscript(t *testing.T) {
	// Create options with verbose output if go test -v is used
	opts := testscript.DefaultOptions()
	opts.Verbose = testing.Verbose()

	// Set some environment variables for the tests
	opts.EnvVars = map[string]string{
		"EXAMPLE_ENV": "example_value",
	}

	// Run all tests in the testdata directory
	testscript.RunDir(t, "testdata", opts)
}

// TestSingleFile demonstrates running a single test file
func TestSingleFile(t *testing.T) {
	opts := testscript.DefaultOptions()
	testscript.RunFile(t, "testdata/basic.txt", opts)
}

// TestPattern demonstrates running tests matching a pattern
func TestPattern(t *testing.T) {
	opts := testscript.DefaultOptions()
	testscript.Run(t, "testdata/*_test.txt", opts)
}
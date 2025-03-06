package main_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tmc/scripttestutil/bridge"
)

func TestMain(m *testing.M) {
	// Build the CLI app for testing
	setup()
	// Run tests
	code := m.Run()
	// Clean up
	teardown()
	os.Exit(code)
}

func setup() {
	// Build the CLI tool
	os.MkdirAll("bin", 0755)
	os.MkdirAll("testdata", 0755)
	
	// Build the CLI for testing
	os.Setenv("GOOS", "")
	os.Setenv("GOARCH", "")
	
	// Simple build command - in a real project, use proper build flags
	os.Chdir("..")
	os.Chdir("cli")
	os.Exec("go", "build", "-o", "bin/cli-app", ".")
	os.Chdir("..")
	
	// Create test files
	createTestFiles()
}

func teardown() {
	// Clean up built binaries
	os.RemoveAll("bin")
}

func createTestFiles() {
	// Basic greeting test
	basicTest := `# Test basic greeting
./bin/cli-app
stdout 'Hello, World!'
! stderr .

# Test with custom name
./bin/cli-app -name Alice
stdout 'Hello, Alice!'
! stderr .
`
	os.WriteFile("testdata/basic.txt", []byte(basicTest), 0644)

	// Test flags
	flagsTest := `# Test verbose flag
./bin/cli-app -verbose
stdout 'Hello, World!'
stderr 'About to print greeting 1 times'

# Test count flag
./bin/cli-app -count 3
stdout 'Hello, World!'
stdout 'Hello, World!'
stdout 'Hello, World!'
! stderr .

# Test multiple flags
./bin/cli-app -name Bob -count 2 -verbose
stdout 'Hello, Bob!'
stdout 'Hello, Bob!'
stderr 'About to print greeting 2 times'

# Test extra arguments
./bin/cli-app arg1 arg2
stdout 'Hello, World!'
stdout 'Extra arguments: arg1, arg2'
! stderr .
`
	os.WriteFile("testdata/flags.txt", []byte(flagsTest), 0644)
}

// TestCLIWithBridge demonstrates testing a CLI application with the bridge
func TestCLIWithBridge(t *testing.T) {
	// Configure test options
	opts := bridge.DefaultOptions()
	opts.Verbose = testing.Verbose()
	
	// Run the tests
	bridge.RunPattern(t, "testdata/*.txt", opts)
}

// TestSpecificFeatures demonstrates testing specific features of the CLI
func TestSpecificFeatures(t *testing.T) {
	t.Run("BasicGreeting", func(t *testing.T) {
		opts := bridge.DefaultOptions()
		bridge.RunFile(t, "testdata/basic.txt", opts)
	})
	
	t.Run("FlagHandling", func(t *testing.T) {
		opts := bridge.DefaultOptions()
		bridge.RunFile(t, "testdata/flags.txt", opts)
	})
}
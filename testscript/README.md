# testscript

Package testscript integrates Go's testing framework with `scripttest` tests.

This package makes it easy to run `scripttest` tests as part of your standard Go test suite,
letting you combine traditional Go unit tests with scripttest-based integration tests.

## Installation

```bash
go get github.com/tmc/scripttestutil/testscript
```

## Basic Usage

```go
package myapp_test

import (
	"testing"
	
	"github.com/tmc/scripttestutil/testscript"
)

func TestApp(t *testing.T) {
	// Run all scripttest tests in the testdata directory
	opts := testscript.DefaultOptions()
	testscript.RunDir(t, "testdata", opts)
}
```

## Features

- Run scripttest tests as standard Go tests
- Simple API with smart defaults
- Supports Docker-based testing
- Manages snapshots for verification testing
- Provides clean integration with Go's testing package

## Running Methods

```go
// Run a single test file
testscript.RunFile(t, "testdata/single.txt", opts)

// Run all tests in a directory
testscript.RunDir(t, "testdata", opts)

// Run tests matching a pattern
testscript.Run(t, "testdata/api/*.txt", opts)
```

## Running Tests in This Repository

This repository contains several test files in the testscript package:

```bash
# Run all tests in all categories
go test -v ./testscript -run TestAll

# Run tests by category
go test -v ./testscript -run TestWithTag/gotools   # Run all Go tools tests
go test -v ./testscript -run TestWithTag/spinner   # Run all spinner tests

# Run individual test groups
go test -v ./testscript -run TestGoTools      # Run all Go tools tests
go test -v ./testscript -run TestSpinner      # Run all spinner tests

# Run specific tests
go test -v ./testscript -run TestSpecificGoTool   # Run a specific Go tool test
go test -v ./testscript -run TestIndividualSpinnerTests/Basic  # Run basic spinner test

# Update snapshots
UPDATE_SNAPSHOTS=1 go test -v ./testscript -run TestAll
```

## Options

```go
opts := testscript.DefaultOptions()

// Enable verbose output (respects `go test -v`)
opts.Verbose = testing.Verbose()

// Update test snapshots
opts.UpdateSnapshots = true

// Set environment variables
opts.EnvVars = map[string]string{
	"API_URL": "http://localhost:8080",
}

// Use Docker for testing
opts.UseDocker = true
opts.DockerImage = "golang:latest"
```

## Example Test Files

Create scripttest files in your testdata directory:

```
# testdata/basic.txt
echo "Hello, World!"
stdout 'Hello, World!'
! stderr .

env TEST_VAR=test_value
env | grep TEST_VAR
stdout 'TEST_VAR=test_value'
```

## Advanced Usage

For more control, use the Runner API directly:

```go
runner := testscript.NewRunner(testscript.Options{
	Pattern: "testdata/*.txt",
	Verbose: true,
})
runner.Run(t)
```

## Testing CLI Applications

```go
func TestMain(m *testing.M) {
	// Build your CLI application
	os.MkdirAll("bin", 0755)
	exec.Command("go", "build", "-o", "bin/myapp", ".").Run()
	
	// Run tests
	os.Exit(m.Run())
}

func TestCLI(t *testing.T) {
	testscript.RunDir(t, "testdata", testscript.DefaultOptions())
}
```

Your test files can then use your CLI application:

```
# testdata/cli.txt
bin/myapp --version
stdout 'v1.0.0'
! stderr .

bin/myapp calculate 2 + 3
stdout '5'
```

## License

MIT - See the LICENSE file for details.
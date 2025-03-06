# Scripttest Bridge

This package provides a clean, simple bridge between Go's standard testing framework and scripttest tests. It allows easy integration of scripttest-based tests into regular Go test suites.

## Features

- **Seamless Integration** - Run scripttest tests alongside regular Go tests
- **Multiple Running Options** - Run individual files, patterns, or directories
- **Configuration** - Extensive options for customizing test execution
- **Go Test Features** - Preserves sub-tests, parallelism, and other testing features
- **Docker Support** - Run tests in Docker containers
- **Snapshot Support** - Create and verify test output snapshots

## Installation

```bash
go get github.com/tmc/scripttestutil/bridge
```

## Basic Usage

```go
package main_test

import (
	"testing"
	
	"github.com/tmc/scripttestutil/bridge"
)

func TestMyCode(t *testing.T) {
	// Run all scripttest tests in the testdata directory
	opts := bridge.DefaultOptions()
	opts.Verbose = testing.Verbose()
	bridge.RunDir(t, "testdata", opts)
}
```

## Running Specific Tests

```go
// Run a single test file
bridge.RunFile(t, "testdata/specific_test.txt", opts)

// Run tests matching a pattern
bridge.RunPattern(t, "testdata/api/*.txt", opts)
```

## Custom Options

```go
opts := bridge.DefaultOptions()

// Set verbose output
opts.Verbose = true

// Enable snapshot updating
opts.UpdateSnapshots = true

// Define custom environment variables
opts.EnvVars = map[string]string{
	"API_URL": "http://localhost:8080",
	"DEBUG": "true",
}

// Use Docker for testing
opts.UseDocker = true
opts.DockerImage = "node:16-alpine"

// Run tests with these options
bridge.RunDir(t, "testdata", opts)
```

## Using the Runner API

For more advanced usage, you can use the Runner API directly:

```go
// Create a runner with custom options
runner := bridge.NewRunner(bridge.TestOptions{
	Pattern: "testdata/*.txt",
	Verbose: true,
	EnvVars: map[string]string{"CUSTOM": "value"},
})

// Run all tests with this runner
runner.Run(t)

// Or run a specific test
runner.RunTest(t, "testdata/specific.txt")
```

## Full Example

See the `bridge_test.go` file for complete examples of using this package.

## License

MIT - See the LICENSE file for details.
# Expect Command Example

This example demonstrates how to use the expect command set in scripttest tests.

## Running the Example

To run all tests:

```bash
go run main.go
```

To run a specific test:

```bash
go run main.go -test simple.txt
```

To enable verbose output:

```bash
go run main.go -v
```

## Test Files

- `simple.txt`: A simple test that interacts with the echo command
- `python.txt`: A more complex test that interacts with the Python interpreter

## Understanding the Example

This example shows how to:

1. Register the expect command set with testscript
2. Run scripttest tests with custom commands
3. Interact with interactive programs like Python
4. Use different patterns for testing

The main components are:

- `main.go`: A simple harness that runs the tests
- `testdata/`: Directory containing the test files
- `commands/expect/`: The expect command set implementation

## Using in Your Own Tests

To use the expect commands in your Go tests:

```go
import (
    "testing"
    
    "github.com/tmc/scripttestutil/commands"
    "github.com/tmc/scripttestutil/testscript"
    "rsc.io/script"
)

func TestWithExpect(t *testing.T) {
    opts := testscript.DefaultOptions()
    opts.Verbose = testing.Verbose()
    
    // Register expect commands
    opts.SetupHook = func(cmds map[string]script.Cmd) {
        commands.RegisterExpect(cmds)
    }
    
    // Run your tests
    testscript.RunDir(t, "testdata", opts)
}
```
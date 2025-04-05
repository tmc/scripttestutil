# scripttestutil

scripttestutil is a package containing utilities related to
[scripttest](https://pkg.go.dev/rsc.io/script/scripttest).

Notably, scripttestutil has a command, `scripttest` that is a tool to run scripttest tests on a
codebase.

## scripttest
scripttest is a command-line tool for managing scripttest testing in a codebase. It provides functionality for running, generating, and managing scripttests with AI assistance.

## Installation

1. Install scripttest tool (requires Go):
```shell
go install https://github.com/tmc/scripttestutil/cmd/scripttest@latest
```

## Usage

scripttest usage:

### Examples

1. Infer config:
   ```
   scripttest infer
   ```

2. Run tests:
   ```
   scripttest test
   ```

3. Run tests in Docker:
   ```
   scripttest -docker test
   ```

4. Record a test execution as an asciicast:
   ```
   scripttest record testdata/example.txt recordings/example.cast
   ```

5. Play an asciicast recording:
   ```
   scripttest play-cast recordings/example.cast
   ```

### Self-Tests

The project includes a suite of self-tests that verify scripttest's functionality using scripttest itself. These serve both as tests and as examples of how to use various features.

To run the self-tests:

```
cd cmd/scripttest/testdata/selftest
./run_all_tests.sh
```

See the [self-tests README](cmd/scripttest/testdata/selftest/README.md) for more details.

### Go Test Integration

The project includes a testscript package that provides a clean integration between Go's standard testing framework and scripttest. This allows you to run scripttest tests as part of your regular Go test suite.

Example usage:

```go
func TestMyFeature(t *testing.T) {
    opts := testscript.DefaultOptions()
    opts.Verbose = testing.Verbose()
    
    // Run all scripttest tests in a directory
    testscript.RunDir(t, "testdata", opts)
}
```

See the [testscript README](testscript/README.md) for more details.

### Reusable Command Sets

The project includes a commands package with reusable command sets that can be added to your scripttest tests. These provide specialized functionality for various domains.

#### Available Command Sets:

- **Expect Commands**: Integration with the expect utility for interacting with interactive programs.

Example usage:

```go
import (
    "github.com/tmc/scripttestutil/commands"
    "github.com/tmc/scripttestutil/testscript"
)

func TestWithExpect(t *testing.T) {
    opts := testscript.DefaultOptions()
    
    // Register expect commands
    opts.SetupHook = func(cmds map[string]script.Cmd) {
        commands.RegisterExpect(cmds)
    }
    
    testscript.RunDir(t, "testdata", opts)
}
```

Then in your test file:

```
# Interact with Python using expect
expect:spawn python3
expect:expect ">>>" 5
expect:send "print('Hello, World!')"
expect:expect "Hello, World!" 5
```

See the [expect README](commands/expect/README.md) for more details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


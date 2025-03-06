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

The project includes a bridge package that provides a clean integration between Go's standard testing framework and scripttest tests. This allows you to run scripttest tests as part of your regular Go test suite.

Example usage:

```go
func TestMyFeature(t *testing.T) {
    opts := bridge.DefaultOptions()
    opts.Verbose = testing.Verbose()
    
    // Run all scripttest tests in a directory
    bridge.RunDir(t, "testdata", opts)
}
```

See the [bridge README](bridge/README.md) for more details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


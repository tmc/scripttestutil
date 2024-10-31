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

scriptest usage:

### Examples

1. Infer config:
   ```
   scripttest infer
   ```

1. Run tests:
   ```
   scripttest test
   ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


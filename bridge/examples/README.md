# Scripttest Bridge Examples

This directory contains examples of how to use the scripttest bridge package for testing various types of applications.

## Examples

### Simple

The `simple` directory contains a basic example showing how to:
- Set up test files programmatically
- Run all tests in a directory
- Run individual test files
- Run tests matching a pattern
- Use different bridge options

Run the example with:
```bash
cd simple
go test -v
```

### CLI

The `cli` directory demonstrates testing a command-line application:
- Building the CLI app for tests
- Testing basic functionality
- Testing command-line flags
- Testing output and error messages
- Organizing tests into sub-tests

Run the example with:
```bash
cd cli
go test -v
```

## Best Practices

1. **Use TestMain for Setup/Teardown**: 
   - Build binaries in TestMain
   - Create test files programmatically if needed
   - Clean up after tests

2. **Structure Test Files Logically**:
   - Create multiple test files for different features
   - Organize by functionality, not by implementation

3. **Use Sub-Tests**:
   - Organize tests using Go's sub-test functionality
   - Makes it easier to run specific tests with `go test -run`

4. **Manage Test Environment**:
   - Use bridge.TestOptions.EnvVars to set environment variables
   - Create isolated test environments when needed

5. **Use Verbose Mode Effectively**:
   - Set opts.Verbose = testing.Verbose() to respect the `-v` flag
   - Add detailed comments in test files
   
6. **Create Snapshots for Complex Output**:
   - Use snapshot testing for complex or large outputs
   - Update snapshots with opts.UpdateSnapshots = true
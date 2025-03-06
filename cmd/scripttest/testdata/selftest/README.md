# Scripttest Self-Tests

This directory contains tests that verify scripttest's own functionality using scripttest itself. These tests serve a dual purpose:

1. Ensure scripttest works correctly
2. Demonstrate how to use scripttest features in real-world scenarios

## Available Tests

- **basic_self_test.txt** - Tests core functionality (stdout/stderr matching, file operations, etc.)
- **snapshot_test.txt** - Tests snapshot creation and verification
- **asciicast_test.txt** - Tests asciicast recording capabilities
- **docker_self_test.txt** - Tests Docker execution environment
- **auto_toolchain_test.txt** - Tests Go toolchain auto-installation
- **run_self_tests.txt** - Meta-test to run and verify all other test files

## Running the Tests

You can run all the tests using the provided shell script:

```bash
./run_all_tests.sh
```

Or run individual tests with:

```bash
# From the project root
go build -o scripttesttool ./cmd/scripttest
./scripttesttool test cmd/scripttest/testdata/selftest/basic_self_test.txt

# For Docker tests:
./scripttesttool -docker test cmd/scripttest/testdata/selftest/docker_self_test.txt
```

## Adding New Self-Tests

When adding new features to scripttest, consider adding corresponding self-tests to verify the functionality. This helps ensure that:

1. The feature works as expected
2. The feature doesn't break when other changes are made
3. The feature has a working example for users to reference

## Notes

- Some tests may fail on first run (especially snapshot tests) because snapshots need to be created first
- Run snapshot tests with `UPDATE_SNAPSHOTS=1` first, then run them normally
- Docker tests require Docker to be installed and running
- Asciicast recording tests may be skipped if `asciinema` is not installed
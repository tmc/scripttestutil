# Go Tools Tests

This directory contains test scripts for testing various Go tools and commands.

## Available Tests

- `simple_fmt_test.txt`: Tests gofmt functionality
- `go_vet_test.txt`: Tests go vet functionality
- `build_run_test.txt`: Tests go build and go run
- `modules_test.txt`: Tests go modules functionality

## Running the Tests

To run all tests:

```
scripttest test testdata/gotools/*.txt
```

To run a specific test:

```
scripttest test testdata/gotools/simple_fmt_test.txt
```

## Creating Recordings

To create an asciinema recording of a test:

```
scripttest record testdata/gotools/simple_fmt_test.txt recordings/simple_fmt.cast
```

## Playing Back Recordings

To play back an asciinema recording:

```
scripttest play-cast recordings/simple_fmt.cast
```

## Test Snapshots

Test snapshots are stored in `testdata/__snapshots__/` directory after running tests.

To update snapshots:

```
UPDATE_SNAPSHOTS=1 scripttest test testdata/gotools/*.txt
```

To play back a snapshot:

```
scripttest playback testdata/__snapshots__/gotools_simple_fmt_test.linux
```
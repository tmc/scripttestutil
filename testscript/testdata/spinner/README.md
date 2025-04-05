# Spinner Package Tests

This directory contains test scripts for testing the functionality of the github.com/tmc/spinner package.

## Available Tests

- `basic_test.txt`: Basic usage of the spinner package
- `concurrent_test.txt`: Concurrent usage of spinners with goroutines
- `custom_spinner_test.txt`: Custom spinner characters and colors

## Running the Tests

To run all spinner tests:

```
scripttest test testdata/spinner/*.txt
```

To run a specific test:

```
scripttest test testdata/spinner/basic_test.txt
```

## Creating Recordings

To create an asciinema recording of a spinner test:

```
scripttest record testdata/spinner/basic_test.txt recordings/spinner_basic.cast
```

The spinner recording will be particularly interesting as it captures the animated spinners in action.

## Playing Back Recordings

To play back an asciinema recording:

```
scripttest play-cast recordings/spinner_basic.cast
```

## Test Snapshots

Test snapshots are stored in `testdata/__snapshots__/` directory after running tests.

To update snapshots:

```
UPDATE_SNAPSHOTS=1 scripttest test testdata/spinner/*.txt
```

To play back a snapshot:

```
scripttest playback testdata/__snapshots__/spinner_basic_test.linux
```
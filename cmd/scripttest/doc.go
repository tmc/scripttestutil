/*
The scripttest command assists with running tests against commands.

It is a CLI wrapper over the rsc.io/script and rsc.io/script/scripttest packages.

Usage:

scripttest [-v] [-u] <command> [args...]

Commands:

	test, run    run scripttest files (default pattern: testdata/*.txt)
	             scripttest test                # uses -p or default pattern
	             scripttest test 'custom/*.txt' # overrides pattern

	             Docker Support:
	             - Use -docker flag to run tests in container
	             - Specify custom image with -docker-image
	             - Include Dockerfile in test file with "-- Dockerfile --" marker

	             Snapshot Support:
	             - Use 'snapshot <n>' command in test file to verify output
	             - Run with UPDATE_SNAPSHOTS=1 to update snapshots

	scaffold     create scripttest scaffold in [dir]
	             scripttest scaffold .

	playback     play back a recorded snapshot
	             scripttest playback testdata/__snapshots__/test.linux

	record       record a test execution as an asciicast
	             scripttest record testdata/example.txt recordings/example.cast

	play-cast    play an asciicast recording
	             scripttest play-cast recordings/example.cast

	convert-cast convert snapshot to asciicast format
	             scripttest convert-cast snapshot.json recording.cast


	help         show available commands and conditions
	             scripttest help

SCRIPTTEST FILE FORMAT:

Test files are plain text files with a specific format. Each file contains:

1. Comments (lines starting with #)
2. Command invocations (commands to run)
3. Assertions (expected output checks)
4. Optional file definitions (marked with -- filename --)
5. Optional Dockerfile definitions (marked with -- Dockerfile --)

Example test file:

	# Basic test description
	echo "Hello, world!"
	stdout 'Hello, world!'     # Check stdout contains this string
	
	? ls nonexistent           # ? prefix for commands expected to fail
	stderr 'No such file'      # Check stderr contains this string
	
	env TEST_VAR=value         # Set environment variables
	mkdir -p testdir           # Create test directories
	
	[linux] echo "Linux only"  # Platform-specific commands
	[darwin] echo "Mac only"   # Only runs on macOS

	snapshot 'test-output'     # Record/verify command output

	-- testfile.txt --         # Create a file that will be available during tests
	This is test content.

	-- Dockerfile --           # Define a Dockerfile for Docker-based tests
	FROM golang:latest
	WORKDIR /app
	COPY . .
	RUN go mod download
	CMD ["go", "test", "-v"]

ASSERTIONS:

The following assertions are available:

1. stdout 'text'  - Check if stdout contains the specified text
2. stderr 'text'  - Check if stderr contains the specified text
3. stdout !       - Check if stdout is empty
4. stderr !       - Check if stderr is empty
5. status 1       - Check if command exit status equals specific code
6. exists filepath - Check if file or directory exists
7. contains file 'content' - Check if file contains specific content
8. snapshot 'name' - Record/verify command output against a snapshot

SPECIAL FEATURES:

1. Platform Conditionals:
   [linux] command     # Only runs on Linux
   [darwin] command    # Only runs on macOS
   [windows] command   # Only runs on Windows
   [unix] command      # Runs on any Unix-like OS (Linux, macOS)

2. Docker Support:
   - Use -docker flag when running tests
   - Specify custom image with -docker-image flag
   - Include a Dockerfile section with "-- Dockerfile --" marker
   - For Windows, use "-- Dockerfile.windows --" 

3. Snapshots:
   - Record output with: snapshot 'name'
   - Update snapshots with: UPDATE_SNAPSHOTS=1 scripttest test
   - Playback snapshots with: scripttest playback path/to/snapshot

4. Asciicast Recordings:
   - Record test execution: scripttest record test.txt output.cast
   - Play recordings: scripttest play-cast output.cast
   - Convert snapshots: scripttest convert-cast snapshot.json output.cast

5. Auto Go Toolchain:
   - Automatically downloads and installs Go if not found
   - Enable with -auto-go flag (default: true)
   - Disable with -auto-go=false

6. Environment Variables:
   - Set with: env NAME=value
   - Test with: env NAME

PROJECT TYPES:

scripttest can be used to test any command-line tool or application:

1. Go CLI Applications
   - Tests cmd/main.go execution
   - Validates command output

2. Shell Scripts/Commands
   - Tests bash/sh script execution
   - Validates environment setup

3. Docker Containers
   - Tests application behavior in containers
   - Validates multi-platform compatibility

4. Build Tools
   - Tests compilation processes
   - Validates build artifacts

SETUP AND RUNNING:

1. Install:
   go install github.com/tmc/scripttestutil/cmd/scripttest@latest

2. Create Test Files:
   - Place .txt test files in testdata/ directory
   - Define test commands and assertions

3. Run Tests:
   scripttest test          # Run all tests
   scripttest -v test       # Run with verbose output
   scripttest -docker test  # Run tests in Docker

4. Update Snapshots:
   UPDATE_SNAPSHOTS=1 scripttest test

ADVANCED USAGE:

1. Project Scaffolding:
   scripttest scaffold .    # Create project scaffold

2. Custom Command Inference:
   scripttest infer         # Generate .scripttest_info

3. Snapshot Playback:
   scripttest playback testdata/__snapshots__/test.json

4. Recording Test Sessions:
   scripttest record testdata/test.txt recording.cast

5. Auto-installing Go Toolchain:
   scripttest -auto-go test  # Automatically installs Go if needed

*/
package main
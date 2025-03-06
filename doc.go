/*
Package scripttestutil provides tools for testing command-line programs using the rsc.io/script package.

# Overview

Scripttestutil is designed to simplify testing of command-line applications by allowing
you to write tests as executable scripts with expected outputs. The tests are written as
plain text files with commands and assertions, making them easy to understand and maintain.

The primary interface to scripttestutil is the 'scripttest' command-line tool, which:
  - Runs script tests matching a pattern (scripttest test)
  - Scaffolds new test directories (scripttest scaffold)
  - Creates snapshots of command output (with snapshot command)
  - Supports running tests in Docker containers (scripttest -docker test)
  - Can play back recorded snapshots (scripttest playback)

# Writing Tests

A scripttest test file is a plain text file typically stored in a 'testdata' directory.
The file contains a series of commands to run, along with assertions about their output.

## Basic Test Structure

A basic test file contains commands followed by assertions:

    # Test that echo works
    echo hello
    stdout 'hello'
    ! stderr .
    
    echo -n hello world
    stdout 'hello world'
    
    # Test exit code
    exit 2
    status 2

## Available Commands and Assertions

Commands:
  - Any shell command (e.g., 'ls', 'echo', 'go build')
  - env NAME=VALUE - Set environment variable
  - cd DIR - Change directory
  - ! CMD - Assert command fails
  - exec CMD - Run command without shell interpretation
  - [condition] CMD - Run command only if condition is true
  - snapshot [path] - Record/verify command output

Assertions:
  - stdout PATTERN - Standard output matches pattern
  - stderr PATTERN - Standard error matches pattern
  - status CODE - Exit status equals code
  - exists PATH - File exists
  - ! exists PATH - File doesn't exist
  - grep PATTERN FILE - File contains pattern

## Multiline Output Matching

For multiline output, use stdout/stderr with multiline strings:

    cat file.txt
    stdout '
    line 1
    line 2
    line 3
    '

## Embedding Files

Test files can include embedded files using the -- markers:

    # Create a Go program and run it
    go run main.go
    stdout 'Hello, World!'
    
    -- main.go --
    package main
    
    import "fmt"
    
    func main() {
        fmt.Println("Hello, World!")
    }

## Environment Variables

Set environment variables using the env command:

    env GO111MODULE=on
    env GOOS=linux
    go build main.go
    exists main

## Conditional Tests

Run tests conditionally based on platform:

    [unix] ls -l
    [unix] stdout 'total'
    
    [windows] dir
    [windows] stdout 'Directory'
    
    [darwin] sw_vers
    [darwin] stdout 'macOS'
    
    [linux] uname -a
    [linux] stdout 'Linux'

# Special Features

## Docker Support

Run tests inside Docker containers using one of these approaches:

1. Using the -docker flag with the default Golang image:
   ```
   scripttest -docker test
   ```

2. Specifying a custom Docker image:
   ```
   scripttest -docker -docker-image=node:18 test
   ```

3. Embedding a Dockerfile in your test file:
   ```
   # Test in Docker
   echo "Running in container"
   
   -- Dockerfile --
   FROM alpine:latest
   RUN apk add bash
   WORKDIR /app
   COPY . .
   CMD ["bash"]
   ```

4. Using platform-specific Dockerfiles:
   ```
   -- Dockerfile --
   # Linux-specific Dockerfile
   FROM ubuntu:22.04
   RUN apt-get update && apt-get install -y nodejs
   
   -- Dockerfile.windows --
   # Windows-specific Dockerfile
   FROM mcr.microsoft.com/windows/servercore:ltsc2022
   RUN powershell -Command "Install-PackageProvider -Name NuGet -Force"
   ```

The Docker container will:
- Mount your test directory as /app in the container
- Pass through environment variables like UPDATE_SNAPSHOTS
- Automatically clean up after test completion
- Support snapshot creation and verification

## Snapshots and Recording

### Basic Snapshots
Record and verify command output:

    # Record snapshot
    snapshot mycommand-output
    mycommand --version
    
    # Later runs will compare against snapshot
    snapshot mycommand-output
    mycommand --version

Update snapshots by setting environment variable:
    UPDATE_SNAPSHOTS=1 scripttest test

### Asciicast Recording
Record and play back terminal sessions as asciicast files:

    # Record a test run as an asciicast
    scripttest record testdata/example.txt recordings/example.cast
    
    # Play back an asciicast recording
    scripttest play-cast recordings/example.cast
    
    # Convert a snapshot to asciicast format
    scripttest convert-cast testdata/__snapshots__/test.json recordings/test.cast

Asciicast recordings can be shared and embedded in documentation, providing
interactive terminal playback with timing information.

## Setting Command Info

Create a .scripttest_info file to define available commands:

    [
      {
        "name": "myapp",
        "summary": "My application",
        "args": "[options]"
      }
    ]

# Test Examples

## Example 1: Testing a CLI Tool

    # Test basic functionality
    mycli --version
    stdout 'v1.0.0'
    
    # Test command with arguments
    mycli add 2 3
    stdout '5'
    
    # Test error handling
    ! mycli add two three
    stderr 'error: invalid arguments'
    status 1

## Example 2: Testing with Environment Variables

    # Test environment variable handling
    env CONFIG_PATH=./config.json
    mycli load
    stdout 'Loaded configuration from ./config.json'
    
    -- config.json --
    {
      "setting": "value"
    }

## Example 3: Testing a Web Server

    # Start server in background
    exec myserver --port 8080 &
    
    # Wait for server to start
    sleep 1
    
    # Test HTTP endpoint
    curl -s http://localhost:8080/health
    stdout '{"status":"ok"}'
    
    # Cleanup
    pkill -f "myserver --port 8080"

## Example 4: Testing with Docker

    # Test inside Docker container
    # Run this test with: scripttest -docker test testdata/docker_test.txt
    
    # Check if we're running in Docker
    [ -f /.dockerenv ] && echo "Running in Docker container" || echo "Not in Docker"
    stdout 'Running in Docker container'
    
    # Test container environment
    go version
    stdout 'go version'
    
    # Ensure file operations work inside container
    echo "Hello from Docker" > testfile.txt
    cat testfile.txt
    stdout 'Hello from Docker'
    
    # Test environment variable passing
    env TEST_VAR=docker_value
    env | grep TEST_VAR
    stdout 'TEST_VAR=docker_value'
    
    # Test network access
    ping -c 1 8.8.8.8
    status 0
    
    # Create a snapshot inside Docker
    mkdir -p __snapshots__
    snapshot __snapshots__/docker-test.json
    go version
    
    -- Dockerfile --
    FROM golang:1.22-alpine
    # Install additional tools for testing
    RUN apk add --no-cache bash curl git openssh-client ping
    # Set up working directory
    WORKDIR /app
    # Pre-install Go tools
    RUN go install gotest.tools/gotestsum@latest
    # Copy test files
    COPY . .
    # Default command
    CMD ["go", "test", "-v"]
    
    -- Dockerfile.windows --
    FROM mcr.microsoft.com/windows/servercore:ltsc2022
    # Install Go (for Windows testing)
    SHELL ["powershell", "-Command"]
    RUN Invoke-WebRequest -Uri https://go.dev/dl/go1.22.0.windows-amd64.zip -OutFile go.zip; \
        Expand-Archive -Path go.zip -DestinationPath C:\; \
        $env:Path += ';C:\go\bin'; \
        [Environment]::SetEnvironmentVariable('Path', $env:Path, [EnvironmentVariableTarget]::Machine)
    WORKDIR C:\app
    COPY . .
    CMD ["go", "test", "-v"]

## Example 5: Testing Go Programs

    # Build the program
    go build -o myapp main.go
    exists myapp
    
    # Run the program
    ./myapp input.txt
    stdout 'Processing input.txt'
    ! stderr .
    
    -- main.go --
    package main
    
    import (
        "fmt"
        "os"
    )
    
    func main() {
        if len(os.Args) < 2 {
            fmt.Fprintln(os.Stderr, "missing input file")
            os.Exit(1)
        }
        fmt.Printf("Processing %s\n", os.Args[1])
    }
    
    -- input.txt --
    test data

# Using the scripttest Command

Installation:
    go install github.com/tmc/scripttestutil/cmd/scripttest@latest

Common commands:

    # Run all tests
    scripttest test
    
    # Run specific test pattern
    scripttest test 'testdata/feature_*.txt'
    
    # Run tests in Docker
    scripttest -docker test
    
    # Specify a custom Docker image
    scripttest -docker -docker-image=node:18 test
    
    # Update snapshots
    UPDATE_SNAPSHOTS=1 scripttest test
    
    # Create scaffold in current directory
    scripttest scaffold .
    
    # Playback a recorded snapshot
    scripttest playback testdata/__snapshots__/output.linux
    
    # Record a test as asciicast
    scripttest record testdata/example.txt recordings/example.cast
    
    # Play an asciicast recording
    scripttest play-cast recordings/example.cast
    
    # Convert snapshot to asciicast
    scripttest convert-cast snapshot.json recording.cast
    
    
    # Run with auto Go toolchain installation (default)
    scripttest -auto-go test
    
    # Disable auto Go toolchain installation
    scripttest -auto-go=false test

For more information, run:
    scripttest help
*/
package scripttestutil
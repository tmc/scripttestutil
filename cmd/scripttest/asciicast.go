package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// AsciicastHeader contains metadata for an asciicast recording
type AsciicastHeader struct {
	Version   int               `json:"version"`
	Width     int               `json:"width"`
	Height    int               `json:"height"`
	Timestamp int64             `json:"timestamp"`
	Title     string            `json:"title,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
}

// AsciicastFrame represents a single frame in an asciicast recording
type AsciicastFrame struct {
	Time    float64 `json:"time"`
	Type    string  `json:"type"`
	Content string  `json:"data"`
}

// RecordAsciicast records a test execution as an asciicast file
func recordAsciicast(testFile, outputFile string) error {
	// Check if asciinema is installed
	_, err := exec.LookPath("asciinema")
	if err != nil {
		return fmt.Errorf("asciinema not found. Please install with 'brew install asciinema' or 'pip install asciinema'")
	}

	// Create the output directory if it doesn't exist
	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// First run the test to make sure it passes
	if verbose {
		fmt.Printf("Running test file %s before recording\n", testFile)
	}

	// Run the test with scripttest
	testCmd := exec.Command("go", "run", ".", "test", testFile)
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr
	if err := testCmd.Run(); err != nil {
		return fmt.Errorf("test failed before recording: %v", err)
	}

	// Record the test execution
	if verbose {
		fmt.Printf("Recording test execution to %s\n", outputFile)
	}
	
	// Start asciinema recording
	recordCmd := exec.Command("asciinema", "rec", "--overwrite", "--command", 
		fmt.Sprintf("go run . -v test %s", testFile), outputFile)
	recordCmd.Stdin = os.Stdin
	recordCmd.Stdout = os.Stdout
	recordCmd.Stderr = os.Stderr
	
	if err := recordCmd.Run(); err != nil {
		return fmt.Errorf("recording failed: %v", err)
	}

	if verbose {
		fmt.Printf("Recording saved to %s\n", outputFile)
	}
	return nil
}

// PlayAsciicast plays an asciicast recording
func playAsciicast(asciicastFile string) error {
	// Check if asciinema is installed
	_, err := exec.LookPath("asciinema")
	if err != nil {
		return fmt.Errorf("asciinema not found. Please install with 'brew install asciinema' or 'pip install asciinema'")
	}

	// Play the recording
	playCmd := exec.Command("asciinema", "play", asciicastFile)
	playCmd.Stdout = os.Stdout
	playCmd.Stderr = os.Stderr
	
	return playCmd.Run()
}

// ConvertSnapshotToAsciicast converts a scripttest snapshot to an asciicast format
func convertSnapshotToAsciicast(snapshotFile, outputFile string) error {
	// Read the snapshot file
	data, err := os.ReadFile(snapshotFile)
	if err != nil {
		return fmt.Errorf("failed to read snapshot file: %v", err)
	}

	// Parse the JSON content
	var snapshot map[string]string
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("invalid snapshot format: %v", err)
	}

	// Create the asciicast header
	header := AsciicastHeader{
		Version:   2,
		Width:     80,
		Height:    25,
		Timestamp: time.Now().Unix(),
		Title:     filepath.Base(snapshotFile),
		Env:       map[string]string{"SHELL": "/bin/bash"},
	}

	// Open the output file
	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer out.Close()

	// Write the header
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return fmt.Errorf("failed to marshal header: %v", err)
	}
	if _, err := out.Write(headerBytes); err != nil {
		return fmt.Errorf("failed to write header: %v", err)
	}
	if _, err := out.Write([]byte("\n")); err != nil {
		return fmt.Errorf("failed to write newline: %v", err)
	}

	// Create frames for stdout and stderr
	writer := bufio.NewWriter(out)
	
	// Process stdout
	if stdout, ok := snapshot["stdout"]; ok && stdout != "" {
		frame := AsciicastFrame{
			Time:    0.1,
			Type:    "o",
			Content: stdout,
		}
		frameBytes, err := json.Marshal(frame)
		if err != nil {
			return fmt.Errorf("failed to marshal stdout frame: %v", err)
		}
		if _, err := writer.Write(frameBytes); err != nil {
			return fmt.Errorf("failed to write stdout frame: %v", err)
		}
		if _, err := writer.Write([]byte("\n")); err != nil {
			return fmt.Errorf("failed to write newline: %v", err)
		}
	}

	// Process stderr
	if stderr, ok := snapshot["stderr"]; ok && stderr != "" {
		frame := AsciicastFrame{
			Time:    0.2,
			Type:    "o",
			Content: "\033[31m" + stderr + "\033[0m", // Red color for stderr
		}
		frameBytes, err := json.Marshal(frame)
		if err != nil {
			return fmt.Errorf("failed to marshal stderr frame: %v", err)
		}
		if _, err := writer.Write(frameBytes); err != nil {
			return fmt.Errorf("failed to write stderr frame: %v", err)
		}
		if _, err := writer.Write([]byte("\n")); err != nil {
			return fmt.Errorf("failed to write newline: %v", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %v", err)
	}

	if verbose {
		fmt.Printf("Converted snapshot to asciicast: %s\n", outputFile)
	}
	return nil
}
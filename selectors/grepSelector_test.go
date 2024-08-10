package selectors

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGrepSelector_RunHappyPath(t *testing.T) {
    // Create temporary directory
    tempDir, err := ioutil.TempDir("", "grep_test")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir) // Clean up

    // Create temporary files with content
    file1 := filepath.Join(tempDir, "file1.txt")
    file2 := filepath.Join(tempDir, "file2.txt")
    content1 := "namespace: nmspc\n"
    content2 := "some other content\nnamespace: nmspc\n"

    if err := ioutil.WriteFile(file1, []byte(content1), 0644); err != nil {
        t.Fatalf("Failed to write to file1: %v", err)
    }
    if err := ioutil.WriteFile(file2, []byte(content2), 0644); err != nil {
        t.Fatalf("Failed to write to file2: %v", err)
    }

    // Create GrepSelector instance
    selector := &GrepSelector{
        grepArgs: []string{"-rnw", tempDir, "-e", "namespace: nmspc"},
    }

    // Run the grep command
    output, err := selector.Run()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    // Expected output
    expectedOutput := []string{
        file1,
        file2,
    }

    // Verify the output
    if !reflect.DeepEqual(output, expectedOutput) {
        t.Errorf("Expected %v, got %v", expectedOutput, output)
    }
}

func TestGrepSelector_RunMinusL(t *testing.T) {
    // Create temporary directory
    tempDir, err := ioutil.TempDir("", "grep_test")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir) // Clean up

    // Create temporary files with content
    file1 := filepath.Join(tempDir, "file1.txt")
    file2 := filepath.Join(tempDir, "file2.txt")
    content1 := "namespace: nmspc\n"
    content2 := "some other content\nnamespace: other-nmspc\n"

    if err := ioutil.WriteFile(file1, []byte(content1), 0644); err != nil {
        t.Fatalf("Failed to write to file1: %v", err)
    }
    if err := ioutil.WriteFile(file2, []byte(content2), 0644); err != nil {
        t.Fatalf("Failed to write to file2: %v", err)
    }

    // Create GrepSelector instance
    selector := &GrepSelector{
        grepArgs: []string{"-rnw", tempDir, "-e", "namespace: nmspc", "-l"},
    }

    // Run the grep command
    output, err := selector.Run()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    // Expected output
    expectedOutput := []string{
        file1,
    }

    // Verify the output
    if !reflect.DeepEqual(output, expectedOutput) {
        t.Errorf("Expected %v, got %v", expectedOutput, output)
    }
}

func TestGrepSelector_RunNoHit(t *testing.T) {
    // Create temporary directory
    tempDir, err := ioutil.TempDir("", "grep_test")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir) // Clean up

    // Create temporary files with content
    file1 := filepath.Join(tempDir, "file1.txt")
    file2 := filepath.Join(tempDir, "file2.txt")
    content1 := "namespace: nmspc\n"
    content2 := "some other content\nnamespace: nmspc\n"

    if err := ioutil.WriteFile(file1, []byte(content1), 0644); err != nil {
        t.Fatalf("Failed to write to file1: %v", err)
    }
    if err := ioutil.WriteFile(file2, []byte(content2), 0644); err != nil {
        t.Fatalf("Failed to write to file2: %v", err)
    }

    // Create GrepSelector instance
    selector := &GrepSelector{
        grepArgs: []string{"-rnw", tempDir, "-e", "namespace: other-nmspc"},
    }

    // Run the grep command
    output, err := selector.Run()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    // Expected output
    expectedOutput := []string{}

    // Verify the output
    if !reflect.DeepEqual(output, expectedOutput) {
        t.Errorf("Expected %v, got %v", expectedOutput, output)
    }
}

func TestGrepSelector_RunInvalidArgs(t *testing.T) {
    // Create temporary directory
    tempDir, err := ioutil.TempDir("", "grep_test")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir) // Clean up

    // Create temporary files with content
    file1 := filepath.Join(tempDir, "file1.txt")
    content1 := "namespace: nmspc\n"

    if err := ioutil.WriteFile(file1, []byte(content1), 0644); err != nil {
        t.Fatalf("Failed to write to file1: %v", err)
    }

    // Create GrepSelector instance with invalid arguments
    selector := &GrepSelector{
        grepArgs: []string{"-p"},
    }

    // Run the grep command
    _, err = selector.Run()
    if err == nil {
        t.Fatalf("Expected an error due to invalid arguments, got none")
    }
}
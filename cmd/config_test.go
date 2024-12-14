package cmd

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/SublimeIbanez/todor/common"
)

func TestRootCommand(t *testing.T) {
	// Redirect stdout and stderr
	var out_buffer, err_buffer bytes.Buffer
	og_out := root_command.OutOrStdout()
	og_err := root_command.ErrOrStderr()
	root_command.SetOut(&out_buffer)
	root_command.SetErr(&err_buffer)

	temp_directory, err := os.MkdirTemp("", "todor_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(temp_directory)

	sample_file_path := filepath.Join(temp_directory, "sample.go")
	sample_content := `package main
// TODO: This is a sample TODO comment

func main() {
	// TODO: Implement main function
}
`
	if err := os.WriteFile(sample_file_path, []byte(sample_content), fs.FileMode(common.DEFAULT_FILE_PERMISSIONS)); err != nil {
		t.Fatalf("Failed to write sample file: %v", err)
	}

	// Test Cases
	tests := []struct {
		name          string
		args          []string
		expectError   bool
		expectOutput  string
		expectFile    string
		cleanupOutput bool
	}{
		{
			name:         "Default Execution",
			args:         []string{},
			expectError:  false,
			expectOutput: "",
			expectFile:   "TODOS.md",
		},
		{
			name:         "Specify Input Path",
			args:         []string{"-p", temp_directory},
			expectError:  false,
			expectOutput: "",
			expectFile:   "TODOS.md",
		},
		{
			name:         "Specify Output Path",
			args:         []string{"-o", "output.md"},
			expectError:  false,
			expectOutput: "",
			expectFile:   "output.md",
		},
		{
			name:         "Specify Input and Output Paths",
			args:         []string{"-p", temp_directory, "-o", "output.md"},
			expectError:  false,
			expectOutput: "",
			expectFile:   "output.md",
		},
		{
			name:        "Invalid Input Path",
			args:        []string{"-p", "/invalid/path"},
			expectError: true,
		},
		{
			name:        "Invalid Flag",
			args:        []string{"--invalid"},
			expectError: true,
		},
	}

	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Reset command state
			root_command.SetArgs(test.args)
			out_buffer.Reset()
			err_buffer.Reset()

			err := root_command.Execute()

			if test.expectError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
				return // Skip further checks on error
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(test.expectFile) != 0 {
				output_path := test.expectFile
				if !filepath.IsAbs(output_path) {
					output_path = filepath.Join(temp_directory, output_path)
				}
				if _, err := os.Stat(output_path); os.IsNotExist(err) {
					t.Errorf("Expected output file %s to exist", output_path)
				} else {
					os.Remove(output_path)
				}
			}
		})
	}

	root_command.SetOut(og_out)
	root_command.SetErr(og_err)
}

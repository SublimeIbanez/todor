package configuration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigFilePath(t *testing.T) {
	config_file_path, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected_file_name := CONFIG_FILE_NAME
	if filepath.Base(config_file_path) != expected_file_name {
		t.Errorf("Expected <%s>, got <%s>", expected_file_name, filepath.Base(config_file_path))
	}

	dir := filepath.Dir(config_file_path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("Expected directory <%s> to exist, but it does not", dir)
	}

	if err := os.RemoveAll(dir); err != nil {
		t.Logf("Failed to clean up test directory: %v", err)
	}
}

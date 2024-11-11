package common

import (
	"fmt"
	"os"
)

func Path(dir_path string) (string, error) {
	var err error
	if len(dir_path) == 0 {
		dir_path, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	if _, err := os.Stat(dir_path); err != nil {
		// TODO: perhaps return a specific custom error type
		return "", err
	}

	fmt.Println(dir_path)

	return dir_path, nil
}

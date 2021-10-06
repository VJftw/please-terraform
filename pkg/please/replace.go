package please

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReplaceInDirectory(dir string, search string, replace string) error {
	err := filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		fileContents, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read '%s': %w", path, err)
		}

		newContents := strings.ReplaceAll(string(fileContents), search, replace)
		if err := os.WriteFile(path, []byte(newContents), fi.Mode()); err != nil {
			return fmt.Errorf("could not write file '%s': %w", path, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not walk files: %w", err)
	}

	return nil
}

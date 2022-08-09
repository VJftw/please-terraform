package please

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

// Sync copies files from the src directory to the dest, deleting files which
// are no longer in the srcs and do not match an exception regex.
func Sync(src string, dest string, exceptionRegexes []string) error {
	absSrc, err := filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("could not get absolute path for '%s': %w", src, err)
	}

	absDest, err := filepath.Abs(dest)
	if err != nil {
		return fmt.Errorf("could not get absolute path for '%s': %w", dest, err)
	}

	srcsSet := map[string]struct{}{}

	if err := filepath.WalkDir(absSrc, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == absSrc {
			return nil
		}

		relPath, err := filepath.Rel(absSrc, path)
		if err != nil {
			return fmt.Errorf("could not determine relPath of '%s': %w", path, err)
		}

		destPath := filepath.Join(absDest, relPath)

		if d.IsDir() {
			log.Debug().Str("dest", destPath).Msg("creating directory")
			if err := os.MkdirAll(destPath, 0750); err != nil {
				return err
			}

			return nil
		}

		if d.Type().IsRegular() {
			if err := CopyFile(path, destPath); err != nil {
				return err
			}
			srcsSet[relPath] = struct{}{}
		}

		return nil
	}); err != nil {
		return err
	}

	if err := filepath.WalkDir(absDest, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == absDest {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(absDest, path)
		if err != nil {
			return fmt.Errorf("could not determine relPath of '%s': %w", path, err)
		}

		if !setContains(srcsSet, relPath) && !isFileNameInFilesAndDirsToKeep(
			path,
			exceptionRegexes,
		) {
			log.Debug().Str("path", path).Msg("removing path")
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// BufferSize represents the file copy buffer.
const BufferSize = 32

// CopyFile copies the src file to the destination.
func CopyFile(src string, dest string) error {
	log.Debug().Str("src", src).Str("dest", dest).Msg("copying")
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	if err := os.MkdirAll(filepath.Dir(dest), 0750); err != nil {
		return err
	}

	destination, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create '%s': %w", dest, err)
	}
	defer destination.Close()
	if err := destination.Chmod(sourceFileStat.Mode()); err != nil {
		return err
	}

	buf := make([]byte, BufferSize)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}

func setContains(set map[string]struct{}, elem string) bool {
	_, ok := set[elem]
	return ok
}

func isFileNameInFilesAndDirsToKeep(fileName string, filesAndDirsToKeepRegexes []string) bool {
	for _, keepRegex := range filesAndDirsToKeepRegexes {
		re := regexp.MustCompile(keepRegex)
		if re.MatchString(fileName) {
			return true
		}
	}

	return false
}

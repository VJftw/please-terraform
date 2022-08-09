package please_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/VJftw/please-terraform/pkg/please"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSync(t *testing.T) {
	var tests = []struct {
		description                string
		srcFileNames               []string
		destFileNames              []string
		expectedFileNamesAfterSync []string
	}{
		{
			"src to empty dest",
			[]string{"main.tf", "variables.tf", ".modules/main.tf", ".modules/variables.tf"},
			[]string{},
			[]string{"main.tf", "variables.tf", ".modules/main.tf", ".modules/variables.tf"},
		},
		{
			"src to dirty dest no exceptions",
			[]string{"main.tf", "variables.tf", ".modules/main.tf", ".modules/variables.tf"},
			[]string{"main_2.tf", "main.tf", "variables.tf", ".modules/main_2.tf", ".modules/variables.tf"},
			[]string{"main.tf", "variables.tf", ".modules/main.tf", ".modules/variables.tf"},
		},
		{
			"src to dirty dest with exceptions",
			[]string{"main.tf", "variables.tf", ".modules/main.tf", ".modules/variables.tf"},
			[]string{"main.tf", "terraform.tfstate", "a.tfstate", ".modules/main.tf", ".terraform/foo/bar"},
			[]string{"main.tf", "variables.tf", "terraform.tfstate", "a.tfstate", ".modules/main.tf", ".modules/variables.tf", ".terraform/foo/bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			srcDir := generateSrc(t, tt.srcFileNames)
			destDir := generateDest(t, tt.destFileNames)
			err := please.Sync(srcDir, destDir, []string{`\.terraform.*`, `.*\.tfstate`})
			require.NoError(t, err)

			expectedAbsFileNamesAfterSync := []string{}
			for _, expectedAbsFileNameAfterSync := range tt.expectedFileNamesAfterSync {
				expectedAbsFileNamesAfterSync = append(expectedAbsFileNamesAfterSync, filepath.Join(destDir, expectedAbsFileNameAfterSync))
			}

			syncedDestPaths := []string{}
			err = filepath.WalkDir(destDir, func(path string, d fs.DirEntry, err error) error {
				if d.IsDir() {
					return nil
				}
				if path == destDir {
					return nil
				}
				syncedDestPaths = append(syncedDestPaths, path)
				return nil
			})
			require.NoError(t, err)

			assert.ElementsMatch(t, expectedAbsFileNamesAfterSync, syncedDestPaths)
		})
	}
}

func generateSrc(t *testing.T, srcFileNames []string) string {
	dir, err := os.MkdirTemp("", "test_sync_src_*")
	require.NoError(t, err)

	for _, srcFileName := range srcFileNames {
		absSrcFileName := filepath.Join(dir, srcFileName)
		err := os.MkdirAll(filepath.Dir(absSrcFileName), 0750)
		require.NoError(t, err)

		file, err := os.Create(absSrcFileName)
		require.NoError(t, err)

		file.Close()
	}

	return dir
}

func generateDest(t *testing.T, destFileNames []string) string {
	dir, err := os.MkdirTemp("", "test_sync_dest_*")
	require.NoError(t, err)

	for _, destFileName := range destFileNames {
		absDestFileName := filepath.Join(dir, destFileName)
		err := os.MkdirAll(filepath.Dir(absDestFileName), 0750)
		require.NoError(t, err)

		file, err := os.Create(absDestFileName)
		require.NoError(t, err)

		file.Close()
	}

	return dir
}

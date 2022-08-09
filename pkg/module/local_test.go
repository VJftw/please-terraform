package module_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/VJftw/please-terraform/pkg/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandLocalExecute(t *testing.T) {

	var tests = []struct {
		description      string
		cmdLocal         *module.CommandLocal
		srcFileNames     []string
		expectedMetadata *module.Metadata
	}{
		{
			"build minimum",
			&module.CommandLocal{
				Name: "my_module",
				Pkg:  "my_package",
				Out:  "my_module",
				Opts: &module.Opts{
					MetadataFile: ".please/terraform/module.json",
				},
			},
			[]string{"main.tf", "variables.tf", "outputs.tf"},
			&module.Metadata{
				Target:  "//my_package:my_module",
				Aliases: []string{"//my_package:my_module"},
			},
		},
		// {
		// 	"build with deps",
		// 	&module.CommandLocal{
		// 		Name: "my_module_with_deps",
		// 		Pkg:  "my_package",
		// 		Out:  "my_module",
		// 		Opts: &module.Opts{
		// 			MetadataFile: ".please/terraform/module.json",
		// 		},
		// 		Deps: []string{"my_other_package/my_other_module", "my_other_package/my_other_module_2"},
		// 	},
		// 	[]string{"main.tf", "variables.tf"},
		// 	&module.Metadata{
		// 		Target:  "//my_package:my_module_with_deps",
		// 		Aliases: []string{"//my_package:my_module_with_deps"},
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {

			absSrcs := generateSrcs(t, tt.srcFileNames)
			tt.cmdLocal.Srcs = strings.Join(absSrcs, " ")

			outDir, err := os.MkdirTemp("", "test_command_local_execute_dest_*")
			require.NoError(t, err)
			tt.cmdLocal.Out = outDir

			err = tt.cmdLocal.Execute([]string{})
			assert.NoError(t, err)

			// check that the out has all of the src filenames flattened
			for _, src := range tt.srcFileNames {
				assert.FileExists(t, filepath.Join(outDir, filepath.Base(src)))
			}

			// check that `.please/terraform/module.json` is complete
			outMetadataPath := filepath.Join(outDir, tt.cmdLocal.Opts.MetadataFile)
			assert.FileExists(t, outMetadataPath)
			actualMetadataBytes, err := os.ReadFile(outMetadataPath)
			require.NoError(t, err)
			actualMetadata := &module.Metadata{}
			err = json.Unmarshal(actualMetadataBytes, actualMetadata)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedMetadata, actualMetadata)
		})
	}

}

func generateSrcs(t *testing.T, srcFileNames []string) []string {
	dir, err := os.MkdirTemp("", "test_command_local_execute_src_*")
	require.NoError(t, err)

	absSrcFileNames := []string{}
	for _, srcFileName := range srcFileNames {
		absSrcFileName := filepath.Join(dir, srcFileName)
		err := os.MkdirAll(filepath.Dir(absSrcFileName), 0750)
		require.NoError(t, err)

		file, err := os.Create(absSrcFileName)
		require.NoError(t, err)

		file.Close()

		absSrcFileNames = append(absSrcFileNames, absSrcFileName)
	}

	return absSrcFileNames
}

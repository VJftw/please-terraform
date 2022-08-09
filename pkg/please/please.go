package please

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/VJftw/please-terraform/internal/logging"
)

var log = logging.NewLogger()

// Opts represents the available options to this Please package as a whole.
type Opts struct {
	PlzOutDir string `long:"plz_out_dir" default:"plz-out/" description:"The plz-out directory relative to the repo root."`
}

// MustRepoRoot returns the path to the root of the repository.
func MustRepoRoot(plzOutDir string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("could not get current working directory")
	}

	var plzOutRegex = regexp.MustCompile(fmt.Sprintf(`%c%s.*`, filepath.Separator, plzOutDir))

	repoRoot := plzOutRegex.ReplaceAllString(cwd, "")

	log.Debug().
		Str("path", repoRoot).
		Msg("resolved repo root path")
	return repoRoot
}

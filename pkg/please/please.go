package please

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/VJftw/please-terraform/internal/logging"
)

var log = logging.NewLogger()

func MustRepoRoot(plzPkg string) string {
	repoRoot := filepath.Dir(MustAbsPlzOut(plzPkg))
	log.Debug().
		Str("path", repoRoot).
		Msg("resolved repo root path")
	return repoRoot
}

func MustAbsPlzOut(plzPkg string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("could not get working directory")
	}

	absPlzOutPath := filepath.Dir(filepath.Dir(strings.TrimSuffix(filepath.Dir(cwd), plzPkg)))
	log.Debug().
		Str("path", absPlzOutPath).
		Msg("resolved plz-out path")

	return absPlzOutPath
}

package please

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func MustRepoRoot(plzPkg string) string {
	return filepath.Dir(MustAbsPlzOut(plzPkg))
}

func MustAbsPlzOut(plzPkg string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Dir(filepath.Dir(strings.TrimSuffix(filepath.Dir(cwd), plzPkg)))
}

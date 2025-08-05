package module

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/VJftw/please-terraform/pkg/please"
)

type BuilderOpt func(*Builder) error

func WithBuilderPleaseTargetName(name string) BuilderOpt {
	return func(b *Builder) error {
		b.pleaseTargetName = name

		return nil
	}
}

func WithBuilderPleaseTargetPackage(pkg string) BuilderOpt {
	return func(b *Builder) error {
		b.pleaseTargetPackage = pkg

		return nil
	}
}

func WithBuilderSpaceSeparatedSrcs(spaceSeparatedSrcs string) BuilderOpt {
	return func(b *Builder) error {
		b.srcs = strings.Split(spaceSeparatedSrcs, " ")

		return nil
	}
}

func WithBuilderOutputDir(outputDir string) BuilderOpt {
	return func(b *Builder) error {
		b.outputDir = outputDir

		return nil
	}
}

func WithBuilderAliases(aliases []string) BuilderOpt {
	return func(b *Builder) error {
		b.aliases = aliases

		return nil
	}
}

func WithBuilderStrip(strip []string) BuilderOpt {
	return func(b *Builder) error {
		b.strip = strip

		return nil
	}
}

func WithBuilderDeps(deps []string) BuilderOpt {
	return func(b *Builder) error {
		b.deps = deps

		return nil
	}
}

type Builder struct {
	pleaseTargetName    string
	pleaseTargetPackage string
	srcs                []string
	outputDir           string
	aliases             []string
	strip               []string
	deps                []string
	metadataFile        string
}

func NewBuilder(opts ...BuilderOpt) (*Builder, error) {
	b := &Builder{
		srcs:         []string{},
		aliases:      []string{},
		strip:        []string{},
		deps:         []string{},
		metadataFile: "plz.terraform.metadata.json",
	}

	for _, opt := range opts {
		if err := opt(b); err != nil {
			return nil, fmt.Errorf("apply opt: %w", err)
		}
	}

	return b, nil
}

func (b *Builder) Build(ctx context.Context) error {
	log.Info().
		Str("name", b.pleaseTargetName).
		Str("pkg", b.pleaseTargetPackage).
		Strs("srcs", b.srcs).
		Str("out", b.outputDir).
		Strs("aliases", b.aliases).
		Strs("strip", b.strip).
		Strs("deps", b.deps).
		Msg("building module")

	log.Debug().Str("path", b.outputDir).Msg("ensuring out directory")
	if err := os.MkdirAll(b.outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("ensure out directory '%s': %w", b.outputDir, err)
	}

	log.Debug().Msg("copying files")
	if err := b.CopySrcs(b.outputDir, b.srcs); err != nil {
		return err
	}

	log.Debug().Strs("strip", b.strip).Msg("stripping directories")
	// Strip directories
	if err := b.StripDirs(b.outputDir, b.strip); err != nil {
		return fmt.Errorf("could not strip directories: %w", err)
	}

	log.Debug().Msg("adding aliases")
	m := &Metadata{
		Aliases: b.aliases,
	}
	m.Aliases = append(m.Aliases, []string{
		// Supports referencing by Please target.
		fmt.Sprintf("//%s:%s", b.pleaseTargetPackage, b.pleaseTargetName),
	}...)

	// log.Debug().Strs("deps", b.deps).Msg("colocating modules")
	// colocate modules
	// if err := ColocateModules(b.metadataFile, b.outputDir, b.deps); err != nil {
	// 	return err
	// }

	log.Debug().Str("path", b.metadataFile).Msg("saving metadata")
	if err := m.Save(filepath.Join(b.outputDir, b.metadataFile)); err != nil {
		return err
	}

	return nil
}

// StripDirs strips the configured directories from the module.
func (b *Builder) StripDirs(out string, strip []string) error {
	for _, stripDir := range strip {
		stripPath := filepath.Join(out, stripDir)
		if err := os.RemoveAll(stripPath); err != nil {
			return fmt.Errorf("could not remove directory '%s': %w", stripDir, err)
		}
	}

	return nil
}

func (b *Builder) CopySrcs(outDir string, srcs []string) error {
	for _, src := range srcs {
		fileInfo, err := os.Stat(src)
		if err != nil {
			return fmt.Errorf("stat '%s': %w", src, err)
		}
		if fileInfo.IsDir() {
			// copy contents of dir
			absBasePath, err := filepath.Abs(src)
			if err != nil {
				return fmt.Errorf("abs '%s': %w", src, err)
			}
			if err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() {
					return nil
				}

				absTargetPath, err := filepath.Abs(path)
				if err != nil {
					return fmt.Errorf("abs '%s': %w", path, err)
				}

				relPath, err := filepath.Rel(absBasePath, absTargetPath)
				if err != nil {
					return err
				}

				log.Info().Str("relPath", relPath).Msg("Copying")

				if err := please.CopyFile(
					path,
					filepath.Join(b.outputDir, relPath),
				); err != nil {
					return fmt.Errorf("copy file: %w", err)
				}

				return nil

			}); err != nil {
				return err
			}
			// if err := please.Sync(src, b.outputDir, []string{}); err != nil {
			// 	return fmt.Errorf("sync '%s' to '%s': %w", src, b.outputDir, err)
			// }
		} else {
			relPath, err := filepath.Rel(b.pleaseTargetPackage, src)
			if err != nil {
				return err
			}
			if err := please.CopyFile(src, filepath.Join(b.outputDir, relPath)); err != nil {
				return fmt.Errorf("copy file: %w", err)
			}
		}
	}

	return nil
}

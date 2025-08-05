package module

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/VJftw/please-terraform/internal/logging"
	"github.com/hashicorp/go-getter"
)

type DownloaderOpt func(*Downloader) error

func WithDownloaderRegistry(registry string) DownloaderOpt {
	return func(d *Downloader) error {
		d.registry = registry

		return nil
	}
}

func WithDownloaderNamespace(namespace string) DownloaderOpt {
	return func(d *Downloader) error {
		d.namespace = namespace

		return nil
	}
}

func WithDownloaderModuleName(moduleName string) DownloaderOpt {
	return func(d *Downloader) error {
		d.moduleName = moduleName

		return nil
	}
}

func WithDownloaderVersion(version string) DownloaderOpt {
	return func(d *Downloader) error {
		d.version = version

		return nil
	}
}

func WithDownloaderSystem(system string) DownloaderOpt {
	return func(d *Downloader) error {
		d.system = system

		return nil
	}
}

func WithDownloaderOutDir(outDir string) DownloaderOpt {
	return func(d *Downloader) error {
		d.outDir = outDir

		return nil
	}
}

type Downloader struct {
	registry   string
	namespace  string
	moduleName string
	version    string
	system     string
	outDir     string
}

func NewDownloader(opts ...DownloaderOpt) (*Downloader, error) {
	d := &Downloader{}

	for _, opt := range opts {
		if err := opt(d); err != nil {
			return nil, fmt.Errorf("apply opt: %w", err)
		}
	}

	return d, nil
}

func (d *Downloader) Download(ctx context.Context) error {
	downloadURL, err := d.GetDownloadURL()
	if err != nil {
		return err
	}

	logging.Logger.Info().Str("url", downloadURL).Msg("downloading")
	if err := getter.GetAny(d.outDir, downloadURL); err != nil {
		return fmt.Errorf("could not get '%s': %w", downloadURL, err)
	}

	// purge .git directories for consistent hashes.
	purgeDirs := []string{".git"}
	for _, purgeDir := range purgeDirs {
		if err := os.RemoveAll(filepath.Join(d.outDir, purgeDir)); err != nil {
			return fmt.Errorf("could not prune '%s' folder: %w", purgeDir, err)
		}
	}

	return nil
}

// GetDownloadURL returns the `hashicorp/go-getter` compatible URI.
func (d *Downloader) GetDownloadURL() (string, error) {
	address := fmt.Sprintf("%s/v1/modules/%s/%s/%s/%s/download",
		d.registry,
		d.namespace,
		d.moduleName,
		d.system,
		d.version,
	)

	getterURL, err := url.Parse(address)
	if err != nil {
		return "", fmt.Errorf("could not parse '%s' as URL: %w", address, err)
	}

	log.Info().Str("url", getterURL.String()).Msg("retrieving download url from registry")
	resp, err := http.Get(getterURL.String())
	if err != nil {
		return "", fmt.Errorf("could not get download URL: %w", err)
	}

	downloadURL := resp.Header.Get("X-Terraform-Get")

	return downloadURL, nil
}

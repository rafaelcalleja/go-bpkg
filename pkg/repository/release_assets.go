package repository

import (
	"errors"
	"fmt"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"path/filepath"
)

type ReleaseAssets struct {
	name           string
	version        string
	sourceTarFile  string
	untarFilesPath string
}

func NewReleaseAssetsWith(options ...func(*ReleaseAssets) error) (ReleaseAssets, error) {
	var releaseAssets = new(ReleaseAssets)

	for _, option := range options {
		err := option(releaseAssets)
		if err != nil {
			return ReleaseAssets{}, err
		}
	}

	err := files.UnTargzAll(releaseAssets.sourceTarFile, releaseAssets.untarFilesPath)

	if err != nil {
		return ReleaseAssets{}, errors.New(fmt.Sprintf("Error Downloading Plugin decompressing file %s in %s", releaseAssets.sourceTarFile, releaseAssets.untarFilesPath))
	}

	return *releaseAssets, nil
}

func NewReleaseAssets(name string, version string, sourceTarFile string, untarFilesPath string) ReleaseAssets {
	releaseAsset, _ := NewReleaseAssetsWith(
		ReleaseAssetsWithName(name),
		ReleaseAssetsWithVersion(version),
		ReleaseAssetsWithSourceTarFile(sourceTarFile),
		ReleaseAssetsWithUntarFilePath(untarFilesPath),
	)

	return releaseAsset
}

func ReleaseAssetsWithName(name string) func(*ReleaseAssets) error {
	return func(p *ReleaseAssets) error {
		p.name = name
		return nil
	}
}

func ReleaseAssetsWithVersion(version string) func(*ReleaseAssets) error {
	return func(p *ReleaseAssets) error {
		p.version = version
		return nil
	}
}

func ReleaseAssetsWithSourceTarFile(sourceTarFile string) func(*ReleaseAssets) error {
	return func(p *ReleaseAssets) error {
		p.sourceTarFile = sourceTarFile
		return nil
	}
}

func ReleaseAssetsWithUntarFilePath(untarFilesPath string) func(*ReleaseAssets) error {
	return func(p *ReleaseAssets) error {
		p.untarFilesPath = untarFilesPath
		return nil
	}
}

func (asset *ReleaseAssets) PackageFolder() string {
	return fmt.Sprintf("%s-%s", asset.name, asset.version)
}

func (asset *ReleaseAssets) DecompressPath() string {
	return filepath.Join(asset.untarFilesPath, asset.PackageFolder())
}

func (asset *ReleaseAssets) Install(metadata *PackageInstaller, releaseDir string) error {
	err := metadata.Install(asset.DecompressPath(), filepath.Join(releaseDir, asset.name))

	if err != nil {
		return errors.New(fmt.Sprintf("Error Installing Package %s", err))
	}

	return nil
}

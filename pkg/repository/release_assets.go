package repository

import (
	"errors"
	"fmt"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ReleaseAssets struct {
	name           string
	version        string
	sourceTarFile  string
	untarFilesPath string
	packageFolder  string
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

	err = filepath.Walk(releaseAssets.untarFilesPath, func(path string, info os.FileInfo, err error) error {
		if path == releaseAssets.untarFilesPath {
			return nil
		}

		if info.IsDir() {
			releaseAssets.packageFolder = filepath.Base(path)

			return io.EOF
		}

		return err
	})

	if err != nil && err != io.EOF {
		return ReleaseAssets{}, err
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
	return asset.packageFolder
	//return fmt.Sprintf("%s-%s", asset.name, asset.version)
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

func (asset *ReleaseAssets) Uninstall(metadata *PackageInstaller, releaseDir string) error {
	err := metadata.Uninstall(filepath.Join(releaseDir, asset.name))

	if err != nil {
		return errors.New(fmt.Sprintf("Error Uninstalling Package %s", err))
	}

	return nil
}

func (asset *ReleaseAssets) IsInstalled(metadata *PackageInstaller, releaseDir string) bool {
	var metadataFiles []string

	err := filepath.Walk(releaseDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".json") {
			metadataFiles = append(metadataFiles, path)
		}

		return nil
	})

	if nil != err {
		return false
	}

	for _, file := range metadataFiles {
		other, _ := NewPackageInstallerFromFileName(file)

		if metadata.Equals(other) {
			return true
		}
	}

	return false
}

func (asset *ReleaseAssets) clone() ReleaseAssets {
	var clone = new(ReleaseAssets)

	clone.name = asset.name
	clone.version = asset.version
	clone.sourceTarFile = asset.sourceTarFile
	clone.untarFilesPath = asset.untarFilesPath
	clone.packageFolder = asset.packageFolder

	return *clone
}

func (asset *ReleaseAssets) CopyWithName(name string) ReleaseAssets {
	clone := asset.clone()
	clone.name = name

	return clone
}

func (asset *ReleaseAssets) CopyWithVersion(version string) ReleaseAssets {
	clone := asset.clone()
	clone.version = version

	return clone
}

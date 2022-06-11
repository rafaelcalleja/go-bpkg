package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"io/ioutil"
	"os"
	"path/filepath"
)

type PackageInstaller struct {
	Manifest string
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Scripts  []string `json:"scripts"`
	Files    []string `json:"files"`
	BinDir   string
}

func NewPackageInstallerFromFileName(filePath string) (*PackageInstaller, error) {
	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		return &PackageInstaller{}, errors.New(fmt.Sprintf("Error can't open file %s", filePath))
	}

	data := new(PackageInstaller)
	data.BinDir = "bin"

	data.Manifest = filepath.Base(filePath)

	err = json.Unmarshal(file, &data)

	if err != nil {
		return &PackageInstaller{}, errors.New(fmt.Sprintf("Error unmarsalling %s", filePath))
	}

	return data, nil
}

func (packageMetadata *PackageInstaller) InstallationFiles() []string {
	f := make([]string, 0)
	for _, file := range packageMetadata.Files {
		f = append(f, file)
	}

	for _, file := range packageMetadata.Scripts {
		f = append(f, file)
	}

	f = append(f, packageMetadata.Manifest)

	return f
}

func (packageMetadata *PackageInstaller) LinkFiles() map[string]string {
	f := make(map[string]string)

	for _, file := range packageMetadata.Scripts {
		f[file] = filepath.Join(packageMetadata.BinDir, filepath.Base(file))
	}

	return f
}

func (packageMetadata *PackageInstaller) InstallationFilesCount() int {
	return len(packageMetadata.Files) + len(packageMetadata.Scripts) + 1
}

func (packageMetadata *PackageInstaller) IsInstalled(installPath string) bool {
	for _, file := range packageMetadata.InstallationFiles() {
		if _, err := os.Stat(filepath.Join(installPath, file)); errors.Is(err, os.ErrNotExist) {
			return false
		}
	}

	return true
}

func (packageMetadata *PackageInstaller) Install(sourceDir string, destDir string) error {
	for _, file := range packageMetadata.InstallationFiles() {
		src := filepath.Join(sourceDir, file)
		dst := filepath.Join(destDir, file)

		if _, err := os.Stat(src); errors.Is(err, os.ErrNotExist) {
			return errors.New(fmt.Sprintf("Source File not found %s", src))
		}

		if err := os.MkdirAll(filepath.Dir(dst), 0755); nil != err {
			return errors.New(fmt.Sprintf("Error Creating dir %s", filepath.Dir(dst)))
		}

		if err := files.CopyFile(src, dst); nil != err {
			return errors.New(fmt.Sprintf("Error Coping %s to %s", src, dst))
		}
	}

	binDir := filepath.Join(destDir, packageMetadata.BinDir)
	if err := os.MkdirAll(binDir, 0755); nil != err {
		return errors.New(fmt.Sprintf("Error Creating bin dir %s", binDir))
	}

	for srcFile, linkFile := range packageMetadata.LinkFiles() {
		src := filepath.Join(sourceDir, srcFile)
		dst := filepath.Join(destDir, linkFile)

		symlinkPathTmp := dst + ".tmp"
		if err := os.Remove(symlinkPathTmp); err != nil && !os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("Error Unlinking Symlink from %s", dst))
		}

		if err := os.Symlink(src, symlinkPathTmp); err != nil {
			return errors.New(fmt.Sprintf("Error Creating Symlink from %s to %s", src, dst))
		}

		if err := os.Rename(symlinkPathTmp, dst); err != nil {
			return errors.New(fmt.Sprintf("Error Renaming Symlink from %s to %s", src, dst))
		}
	}

	return nil
}

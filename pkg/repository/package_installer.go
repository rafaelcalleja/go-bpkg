package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type PackageInstaller struct {
	Manifest string
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Scripts  []string `json:"scripts"`
	Files    []string `json:"files"`
	BinDir   string
}

func NewPackageInstallerFromLiteral(metadata string, filePath string) (*PackageInstaller, error) {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); nil != err {
		return &PackageInstaller{}, errors.New(fmt.Sprintf("Error Creating dir %s", filepath.Dir(filePath)))
	}

	err := ioutil.WriteFile(filePath, []byte(metadata), 0644)

	if err != nil {
		return &PackageInstaller{}, errors.New(fmt.Sprintf("Error creating metadata file %s", filePath))
	}

	return NewPackageInstallerFromFileName(filePath)
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

func (packageMetadata *PackageInstaller) LinkFiles() []string {
	f := make([]string, 0)

	for _, file := range packageMetadata.Scripts {
		f = append(f, file)
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

	binDir := filepath.Join(filepath.Dir(destDir), packageMetadata.BinDir)
	if err := os.MkdirAll(binDir, 0755); nil != err {
		return errors.New(fmt.Sprintf("Error Creating bin dir %s", binDir))
	}

	for _, srcFile := range packageMetadata.LinkFiles() {
		src := filepath.Join(sourceDir, srcFile)
		dst := filepath.Join(destDir, srcFile)

		if _, err := os.Stat(src); errors.Is(err, os.ErrNotExist) {
			return errors.New(fmt.Sprintf("Source File not found %s", src))
		}

		if err := os.MkdirAll(filepath.Dir(dst), 0755); nil != err {
			return errors.New(fmt.Sprintf("Error Creating dir %s", filepath.Dir(dst)))
		}

		if err := files.CopyFile(src, dst); nil != err {
			return errors.New(fmt.Sprintf("Error Coping %s to %s", src, dst))
		}

		dstLink := filepath.Join(binDir, filepath.Base(srcFile))

		symlinkPathTmp := dstLink + ".tmp"
		if err := os.Remove(symlinkPathTmp); err != nil && !os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("Error Unlinking Symlink from %s", dst))
		}

		if err := os.Symlink(dst, symlinkPathTmp); err != nil {
			return errors.New(fmt.Sprintf("Error Creating Symlink from %s to %s", src, dst))
		}

		if err := os.Rename(symlinkPathTmp, dstLink); err != nil {
			return errors.New(fmt.Sprintf("Error Renaming Symlink from %s to %s", src, dst))
		}
	}

	return nil
}

func (packageMetadata *PackageInstaller) Equals(other *PackageInstaller) bool {
	return packageMetadata.Name == other.Name &&
		packageMetadata.BinDir == other.BinDir &&
		packageMetadata.equal(packageMetadata.Scripts, other.Scripts) &&
		packageMetadata.equal(packageMetadata.Files, other.Files) &&
		packageMetadata.Manifest == other.Manifest &&
		packageMetadata.Version == other.Version

}

func (packageMetadata *PackageInstaller) equal(a, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)

	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

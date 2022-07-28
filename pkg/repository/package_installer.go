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
	Manifest string   `json:"-"`
	Name     string   `json:"name,omitempty"`
	Version  string   `json:"version,omitempty"`
	Scripts  []string `json:"scripts,omitempty"`
	Files    []string `json:"files,omitempty"`
	BinDir   string   `json:"-"`
}

var (
	DefaultPackageFile                 = "package.json"
	ErrPackageInstallerNameCantBeEmpty = errors.New("package installer name can't be empty")
)

func NewPackageInstallerWith(options ...func(*PackageInstaller) error) (PackageInstaller, error) {
	var packageInstaller = new(PackageInstaller)

	for _, option := range options {
		err := option(packageInstaller)
		if err != nil {
			return PackageInstaller{}, err
		}
	}

	if "" == packageInstaller.BinDir {
		packageInstaller.BinDir = "bin"
	}

	if "" == packageInstaller.Manifest {
		packageInstaller.Manifest = DefaultPackageFile
	}

	if "" == packageInstaller.Name {
		return PackageInstaller{}, ErrPackageInstallerNameCantBeEmpty
	}

	return *packageInstaller, nil
}

func NewPackageInstaller(manifest string, name string, version string, scripts []string, files []string, binDir string) PackageInstaller {
	packageInstaller, _ := NewPackageInstallerWith(
		PackageInstallerWithManifest(manifest),
		PackageInstallerWithName(name),
		PackageInstallerWithVersion(version),
		PackageInstallerWithScripts(scripts),
		PackageInstallerWithFiles(files),
		PackageInstallerWithBinDir(binDir),
	)

	return packageInstaller
}

func PackageInstallerWithScripts(scripts []string) func(*PackageInstaller) error {
	return func(p *PackageInstaller) error {
		p.Scripts = scripts
		return nil
	}
}

func PackageInstallerWithFiles(files []string) func(*PackageInstaller) error {
	return func(p *PackageInstaller) error {
		p.Files = files
		return nil
	}
}

func PackageInstallerWithBinDir(binDir string) func(*PackageInstaller) error {
	return func(p *PackageInstaller) error {
		p.BinDir = binDir
		return nil
	}
}

func PackageInstallerWithManifest(manifest string) func(*PackageInstaller) error {
	return func(p *PackageInstaller) error {
		p.Manifest = manifest
		return nil
	}
}

func PackageInstallerWithVersion(version string) func(*PackageInstaller) error {
	return func(p *PackageInstaller) error {
		p.Version = version
		return nil
	}
}

func PackageInstallerWithName(name string) func(*PackageInstaller) error {
	return func(p *PackageInstaller) error {
		p.Name = name
		return nil
	}
}

func NewPackageInstallerFromLiteral(metadata string) (*PackageInstaller, error) {
	tmpDir, err := os.MkdirTemp("", "temp-pkg-metadata")
	if nil != err {
		return &PackageInstaller{}, errors.New(fmt.Sprintf("Error Creating temporal dir %s", tmpDir))
	}

	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, DefaultPackageFile)
	err = ioutil.WriteFile(filePath, []byte(metadata), 0644)

	if err != nil {
		return &PackageInstaller{}, errors.New(fmt.Sprintf("Error creating temporal metadata file %s", filePath))
	}

	defer os.Remove(filePath)

	return NewPackageInstallerFromFileName(filePath)
}

func NewPackageInstallerFromFileName(filePath string) (*PackageInstaller, error) {
	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		return &PackageInstaller{}, errors.New(fmt.Sprintf("Error can't open file %s", filePath))
	}

	data := new(PackageInstaller)

	data.Manifest = filepath.Base(filePath)

	err = json.Unmarshal(file, &data)

	if err != nil {
		return &PackageInstaller{}, errors.New(fmt.Sprintf("Error unmarsalling %s", filePath))
	}

	newPackageInstaller, err := NewPackageInstallerWith(
		PackageInstallerWithManifest(data.Manifest),
		PackageInstallerWithName(data.Name),
		PackageInstallerWithVersion(data.Version),
		PackageInstallerWithScripts(data.Scripts),
		PackageInstallerWithFiles(data.Files),
	)

	return &newPackageInstaller, err
}

func (packageMetadata *PackageInstaller) InstallationFiles() []string {
	f := make([]string, 0)
	for _, file := range packageMetadata.Files {
		f = append(f, file)
	}

	for _, file := range packageMetadata.Scripts {
		f = append(f, file)
	}

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

	packageMetadata.Name = filepath.Base(destDir)
	metadataFile, err := json.MarshalIndent(packageMetadata, "", " ")
	if nil != err {
		return err
	}

	manifestPath := filepath.Join(destDir, packageMetadata.Manifest)
	if err = ioutil.WriteFile(manifestPath, metadataFile, 0644); nil != err {
		return errors.New(fmt.Sprintf("Error Creating manifest file %s", manifestPath))
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

		//ToDo coverage test
		if false == filepath.IsAbs(dst) {
			if dst, err = filepath.Abs(dst); err != nil {
				return errors.New(fmt.Sprintf("Error getting absolute representation of path %s", dst))
			}
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

func (packageMetadata *PackageInstaller) Uninstall(destDir string) error {
	if false == packageMetadata.IsInstalled(destDir) {
		return errors.New(fmt.Sprintf("Package not installed"))
	}

	for _, file := range packageMetadata.InstallationFiles() {
		src := filepath.Join(destDir, file)

		if _, err := os.Stat(src); errors.Is(err, os.ErrNotExist) {
			return errors.New(fmt.Sprintf("Source File not found %s", src))
		}

		if err := os.Remove(src); nil != err {
			return errors.New(fmt.Sprintf("Error uninstalling file %s: %s", src, err))
		}
	}

	manifestPath := filepath.Join(destDir, packageMetadata.Manifest)
	if err := os.Remove(manifestPath); nil != err {
		return errors.New(fmt.Sprintf("Error uninstalling manifest file %s", manifestPath))
	}

	for _, srcFile := range packageMetadata.LinkFiles() {
		src := filepath.Join(filepath.Dir(destDir), packageMetadata.BinDir, filepath.Base(srcFile))

		if err := os.Remove(src); nil != err {
			return errors.New(fmt.Sprintf("Error uninstalling link file %s", src))
		}
	}

	if err := os.Remove(destDir); nil != err {
		return errors.New(fmt.Sprintf("Error uninstalling directory %s", destDir))
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

func PackagesInstalled(releaseDir string) ([]*PackageInstaller, error) {
	assets := make([]*PackageInstaller, 0)

	metadataFiles, err := metadataFilesFinder(releaseDir)

	if err != nil {
		return []*PackageInstaller{}, err
	}

	for _, file := range metadataFiles {
		other, err := NewPackageInstallerFromFileName(file)

		if err != nil {
			continue
		}

		assets = append(assets, other)
	}

	return assets, nil
}

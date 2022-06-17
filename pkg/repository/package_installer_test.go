package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestUnmarshall(t *testing.T) {
	packageFile, err := NewPackageInstallerFromFileName("testdata/package.json")
	require.Nil(t, err)

	assert.Equal(t, packageFile.Manifest, "package.json")
	assert.Equal(t, packageFile.Name, "test-package")
	assert.Equal(t, packageFile.Version, "0.0.1")

	assert.Equal(t, packageFile.Scripts, []string{"src/scripts/file1", "src/scripts/file2"})
	assert.Equal(t, packageFile.Files, []string{"src/files/file1", "src/files/file2"})

	assert.Equal(t, 5, packageFile.InstallationFilesCount())

	assert.Equal(
		t,
		[]string{"src/files/file1", "src/files/file2", "src/scripts/file1", "src/scripts/file2"},
		packageFile.InstallationFiles(),
	)

	assert.Equal(
		t,
		[]string{"src/scripts/file1", "src/scripts/file2"},
		packageFile.LinkFiles(),
	)

	tempDir, _ := os.MkdirTemp("", "temp-test-plugin-folder")
	installDir, _ := os.MkdirTemp("", "temp-test-plugin-folder")
	defer os.RemoveAll(tempDir)

	assert.False(t, packageFile.IsInstalled(tempDir))

	for _, file := range packageFile.InstallationFiles() {
		err = os.MkdirAll(filepath.Join(tempDir, "test-package", filepath.Dir(file)), 0755)
		require.Nil(t, err)

		_, err = os.Create(filepath.Join(tempDir, "test-package", file))
		require.Nil(t, err)
	}

	assert.True(t, packageFile.IsInstalled(filepath.Join(tempDir, "test-package")))

	err = packageFile.Install(filepath.Join(tempDir, "test-package"), installDir)
	require.Nil(t, err)

	assert.True(t, packageFile.IsInstalled(installDir))
}

func TestNotFound(t *testing.T) {
	_, err := NewPackageInstallerFromFileName("testdata/package_not_found.json")
	assert.NotNil(t, err)
}

func TestFromLiteralUnmarshall(t *testing.T) {
	metadata := "{\"name\":\"test-package\",\"version\":\"0.0.1\",\"description\":\"test-description\",\"files\":[\"src/files/file1\",\"src/files/file2\"],\"scripts\":[\"src/scripts/file1\",\"src/scripts/file2\"]}"

	packageFile, err := NewPackageInstallerFromLiteral(metadata)
	require.Nil(t, err)

	other, err := NewPackageInstallerFromFileName("testdata/package.json")
	require.Nil(t, err)

	assert.True(t, packageFile.Equals(other))

	_, err = NewPackageInstallerFromFileName("testdata/invalid.json")
	assert.NotNil(t, err)
}

func TestRequiredName(t *testing.T) {
	_, err := NewPackageInstallerFromLiteral("{}")
	assert.Equal(t, ErrPackageInstallerNameCantBeEmpty, err)
}

func TestNewPackageInstaller(t *testing.T) {
	packageNew := NewPackageInstaller("a.json", "name", "version", []string{"scripts"}, []string{"files"}, "bin")

	assert.Equal(t, "version", packageNew.Version)
	assert.Equal(t, "a.json", packageNew.Manifest)
	assert.Equal(t, "name", packageNew.Name)
	assert.Equal(t, []string{"scripts"}, packageNew.Scripts)
	assert.Equal(t, []string{"files"}, packageNew.Files)
	assert.Equal(t, "bin", packageNew.BinDir)

	otherPackage := NewPackageInstaller("", "name", "version", []string{}, []string{}, "")
	assert.Equal(t, "version", otherPackage.Version)
	assert.Equal(t, DefaultPackageFile, otherPackage.Manifest)
	assert.Equal(t, "name", otherPackage.Name)
	assert.Equal(t, []string{}, otherPackage.Scripts)
	assert.Equal(t, []string{}, otherPackage.Files)
	assert.Equal(t, "bin", otherPackage.BinDir)

	assert.Equal(
		t,
		PackageInstaller{},
		NewPackageInstaller("a.json", "", "version", []string{"scripts"}, []string{"files"}, "bin"),
	)
}

func TestPackagesInstalled(t *testing.T) {
	packages, err := PackagesInstalled("testdata/packages_installed")
	require.Nil(t, err)

	assert.Equal(t, 2, len(packages))

	_, err = PackagesInstalled("testdata/not_found")
	assert.NotNil(t, err)
}

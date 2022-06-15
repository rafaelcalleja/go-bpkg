package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestInstallFromCustomJson(t *testing.T) {
	tempFolder, _ := os.MkdirTemp("", "temp-test-folder")
	defer os.RemoveAll(tempFolder)

	installFolder, _ := os.MkdirTemp("", "temp-install-folder")
	defer os.RemoveAll(installFolder)

	asset := NewReleaseAssets("testName", "v1.1", "testdata/sourceTarFile.tar.gz", tempFolder)
	assert.Equal(t, asset.PackageFolder(), "assert.sh-1.1")

	packageInstaller, err := NewPackageInstallerWith(
		PackageInstallerWithManifest("demo.json"),
		PackageInstallerWithName(asset.name),
		PackageInstallerWithVersion(asset.version),
		PackageInstallerWithFiles([]string{"assert.sh", "tests.sh"}),
	)

	require.Nil(t, err)

	err = asset.Install(&packageInstaller, installFolder)
	require.Nil(t, err)

	assert.True(t, asset.IsInstalled(&packageInstaller, installFolder))

	notPackageInstaller, _ := NewPackageInstallerWith(
		PackageInstallerWithManifest(packageInstaller.Manifest),
		PackageInstallerWithName(packageInstaller.Name),
		PackageInstallerWithVersion("v2.0"),
		PackageInstallerWithFiles(packageInstaller.Files),
	)

	assert.False(t, asset.IsInstalled(&notPackageInstaller, installFolder))

	err = asset.Uninstall(&packageInstaller, installFolder)
	require.Nil(t, err)

	assert.False(t, asset.IsInstalled(&packageInstaller, installFolder))
}

func TestCloneWith(t *testing.T) {
	tempFolder, _ := os.MkdirTemp("", "temp-test-folder")
	defer os.RemoveAll(tempFolder)

	installFolder, _ := os.MkdirTemp("", "temp-install-folder")
	defer os.RemoveAll(installFolder)

	asset := NewReleaseAssets("testName", "v1.1", "testdata/sourceTarFile.tar.gz", tempFolder)

	expectedName := "newName"
	clone := asset.CopyWithName(expectedName)
	assert.Equal(t, asset, asset.clone())
	assert.Equal(t, expectedName, clone.name)
	assert.NotEqual(t, asset, clone)

	expectedVersion := "v2.0"
	clone = asset.CopyWithVersion(expectedVersion)
	assert.Equal(t, asset, asset.clone())
	assert.Equal(t, expectedVersion, clone.version)
}

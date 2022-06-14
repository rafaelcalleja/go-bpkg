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
		[]string{"src/files/file1", "src/files/file2", "src/scripts/file1", "src/scripts/file2", "package.json"},
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
	releaseDir, _ := os.MkdirTemp("", "temp-test-package-folder")

	packageFile, err := NewPackageInstallerFromLiteral(metadata, filepath.Join(releaseDir, "package.json"))
	require.Nil(t, err)

	other, err := NewPackageInstallerFromFileName("testdata/package.json")
	require.Nil(t, err)

	assert.True(t, packageFile.Equals(other))
}

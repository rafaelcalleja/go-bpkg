package repository

import (
	"errors"
	"fmt"
	ghfactory "github.com/cli/cli/v2/pkg/cmd/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var testReleaseVersion = ReleaseVersion{"ORG", "NAME", "VERSION", "package.json"}

func newAssetProviderMock() *GithubProvider {
	return NewGithubProvider(ghfactory.New("0.0.0"))
}

func TestDownloadPlugin(t *testing.T) {
	releaseDir, _ := os.MkdirTemp("", "temp-test-plugin-folder")
	defer os.RemoveAll(releaseDir)

	testReleaseVersionOther, err := NewReleaseVersionWith(
		ReleaseVersionWithOrganization("rafaelcalleja"),
		ReleaseVersionWithName("assert.sh"),
		ReleaseVersionWithVersion("v1.1"),
	)
	require.Nil(t, err)

	assert.False(t, testReleaseVersionOther.IsInstalled(releaseDir))
	err = testReleaseVersionOther.DownloadAndInstallAsset(newAssetProviderMock(), releaseDir)
	require.Nil(t, err)

	assert.True(t, testReleaseVersionOther.IsInstalled(releaseDir))

	assert.True(t, testReleaseVersionOther.IsVersionInstalled(testReleaseVersionOther.Version(), releaseDir))
}

func TestManifestPath(t *testing.T) {
	releaseDir, _ := os.MkdirTemp("", "temp-test-plugin-folder")
	assert.Contains(t, testReleaseVersion.FilePath(releaseDir), testReleaseVersion.Manifest())
}

func TestSetVersionRemoveV(t *testing.T) {
	versionWithV := "v0.0.1"
	versionWithOutV := "0.0.1"

	mockReleaseVersion := NewReleaseVersion("dummy", "name", versionWithV)

	assert.Equal(t, mockReleaseVersion.Version(), versionWithV)
	assert.Equal(t, mockReleaseVersion.VersionWithV(), versionWithV)
	assert.Equal(t, mockReleaseVersion.VersionWithOutV(), versionWithOutV)
}

func TestReleaseVersionWithFQDN(t *testing.T) {
	t.Run("valid FQDN Format", func(t *testing.T) {
		valid := map[string]ReleaseVersion{
			"organization/name":        {"organization", "name", "", "package.json"},
			"organization/name:v1.0":   {"organization", "name", "v1.0", "package.json"},
			"organization/name:latest": {"organization", "name", "latest", "package.json"},
			"name/organization:2.0":    {"name", "organization", "2.0", "package.json"},
		}

		for fqdn, expected := range valid {
			releaseVersion, err := ReleaseVersionWithFQP(fqdn)
			require.Nil(t, err)

			assert.True(t, releaseVersion.Equals(expected))
		}
	})

	t.Run("Invalid FQDN Format", func(t *testing.T) {
		invalid := []string{
			"organization//name",
			"/organization/name",
			"organization/name::",
			"organization/name:^1.0",
			"organization/name:~1.0",
			"organization/name:1-0",
			"organization /name:1.0",
			"organization/name :1.0",
			"organization\\/name:1.0",
		}

		for _, fqdn := range invalid {
			_, err := ReleaseVersionWithFQP(fqdn)
			assert.Equal(t, err, ErrFullyQualifyPackageInvalidFormat)
		}
	})
}

func TestNewPluginLatestVersion(t *testing.T) {
	finder := newMockReleaseVersionFinder()

	t.Run("Can't find latest version", func(t *testing.T) {
		expected := errors.New(fmt.Sprintf("cant find any version"))

		finder.LatestFn = func(string, string) (string, error) {
			return "", expected
		}

		_, actual := NewReleaseLatestVersion("dummy", "dum", finder)
		assert.Equal(t, fmt.Sprintf("%T", actual), fmt.Sprintf("%T", expected))
	})

	t.Run("Founded latest version", func(t *testing.T) {
		expected := "v3"
		finder.LatestFn = func(string, string) (string, error) {
			return expected, nil
		}

		releaseVersion, err := NewReleaseLatestVersion("dummy", "dum", finder)
		require.NoError(t, err)

		assert.Equal(t, expected, releaseVersion.VersionWithV())
	})

}

func TestDownloadAssetWithMultipleInstallations(t *testing.T) {
	downloadDir, _ := os.MkdirTemp("", "temp-test-download-folder")
	installDirA, _ := os.MkdirTemp("", "temp-test-install-folder")
	installDirB, _ := os.MkdirTemp("", "temp-test-install-b-folder")
	defer os.RemoveAll(downloadDir)
	defer os.RemoveAll(installDirA)
	defer os.RemoveAll(installDirB)

	testReleaseVersionOther, err := NewReleaseVersionWith(
		ReleaseVersionWithOrganization("rafaelcalleja"),
		ReleaseVersionWithName("assert.sh"),
		ReleaseVersionWithVersion("v1.1"),
	)
	require.Nil(t, err)

	assert.False(t, testReleaseVersionOther.IsInstalled(installDirA))
	asset, err := testReleaseVersionOther.DownloadAsset(newAssetProviderMock(), downloadDir)
	require.Nil(t, err)

	err = testReleaseVersionOther.InstallAsset(asset, installDirA)
	require.Nil(t, err)

	assert.True(t, testReleaseVersionOther.IsInstalled(installDirA))

	cloneWithName := testReleaseVersionOther.CopyWithName("copy-of-asset")
	err = cloneWithName.InstallAssetWithName("copy-of-asset", asset, installDirB)

	require.Nil(t, err)

	assert.True(t, cloneWithName.IsInstalled(installDirB))
}

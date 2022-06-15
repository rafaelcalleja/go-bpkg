package repository

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type ReleasesProvider interface {
	Download(releaseVersion *ReleaseVersion, downloadDir string) error
}

type ReleaseVersionFinder interface {
	Latest(organization string, name string) (string, error)
}

type ReleaseVersion struct {
	Organization string
	Name         string
	version      string
	manifest     string
}

func NewReleaseVersionWith(options ...func(*ReleaseVersion) error) (ReleaseVersion, error) {
	var releaseVersion = new(ReleaseVersion)

	for _, option := range options {
		err := option(releaseVersion)
		if err != nil {
			return ReleaseVersion{}, err
		}
	}

	if "" == releaseVersion.manifest {
		releaseVersion.manifest = "package.json"
	}

	return *releaseVersion, nil
}

func NewReleaseVersion(organization string, name string, version string) ReleaseVersion {
	releaseVersion, _ := NewReleaseVersionWith(
		ReleaseVersionWithOrganization(organization),
		ReleaseVersionWithName(name),
		ReleaseVersionWithVersion(version),
	)

	return releaseVersion
}

func ReleaseVersionWithOrganization(organization string) func(*ReleaseVersion) error {
	return func(p *ReleaseVersion) error {
		p.Organization = organization
		return nil
	}
}

func ReleaseVersionWithName(name string) func(*ReleaseVersion) error {
	return func(p *ReleaseVersion) error {
		p.Name = name
		return nil
	}
}

func ReleaseVersionWithVersion(version string) func(*ReleaseVersion) error {
	return func(p *ReleaseVersion) error {
		p.SetVersion(version)
		return nil
	}
}

func ReleaseVersionWithManifest(manifest string) func(*ReleaseVersion) error {
	return func(p *ReleaseVersion) error {
		p.manifest = manifest
		return nil
	}
}

func ReleaseVersionWithFQP(fqp string) (ReleaseVersion, error) {
	fqpVO, err := NewFullyQualifyPackage(fqp)

	if nil != err {
		return ReleaseVersion{}, err
	}

	return NewReleaseVersion(fqpVO.Organization(), fqpVO.Name(), fqpVO.Version()), nil
}

func (releaseVersion *ReleaseVersion) Manifest() string {
	return releaseVersion.manifest
}

func (releaseVersion *ReleaseVersion) FilePath(releaseDir string) string {
	return filepath.Join(releaseDir, releaseVersion.Manifest())
}

func (releaseVersion *ReleaseVersion) IsInstalled(releaseDir string) bool {
	return releaseVersion.HasPackageMetadata(filepath.Join(releaseDir, releaseVersion.NameWithOrganization()))
}

func (releaseVersion *ReleaseVersion) IsVersionInstalled(version string, releaseDir string) bool {
	if false == releaseVersion.HasPackageMetadata(filepath.Join(releaseDir, releaseVersion.NameWithOrganization())) {
		return false
	}

	packageMetadata := releaseVersion.MustPackageMetadata(filepath.Join(releaseDir, releaseVersion.NameWithOrganization()))

	return packageMetadata.Version == version
}

func (releaseVersion *ReleaseVersion) MustPackageMetadata(releaseDir string) *PackageInstaller {
	packageMetadata, _ := releaseVersion.GetPackageMetadata(releaseDir)

	return packageMetadata
}

func (releaseVersion *ReleaseVersion) GetPackageMetadata(releaseDir string) (*PackageInstaller, error) {
	manifestFile := filepath.Join(releaseDir, releaseVersion.Manifest())

	return NewPackageInstallerFromFileName(manifestFile)
}

func (releaseVersion *ReleaseVersion) HasPackageMetadata(releaseDir string) bool {
	_, err := releaseVersion.GetPackageMetadata(releaseDir)

	return err == nil
}

func (releaseVersion *ReleaseVersion) DownloadAsset(provider ReleasesProvider, releaseDir string) (ReleaseAssets, error) {
	var err error

	tempDirectory, err := os.MkdirTemp("", "temp-plugin-folder")
	if err != nil {
		return ReleaseAssets{}, errors.New(fmt.Sprintf("Error Downloading Plugin can't create temp dir %s", tempDirectory))
	}

	err = provider.Download(releaseVersion, tempDirectory)
	if err != nil {
		return ReleaseAssets{}, errors.New(fmt.Sprintf("Error Downloading Plugin executing cmd from provider %s", provider))
	}

	dirFiles, err := ioutil.ReadDir(tempDirectory)
	if err != nil {
		return ReleaseAssets{}, errors.New(fmt.Sprintf("Error Downloading Plugin ioutil.ReadDir failed at %s", tempDirectory))
	}

	var pluginFileTar string
	for _, f := range dirFiles {
		pluginFileTar = filepath.Join(tempDirectory, f.Name())
	}

	return NewReleaseAssets(releaseVersion.NameWithOrganization(), releaseVersion.VersionWithOutV(), pluginFileTar, tempDirectory), nil
}

func (releaseVersion *ReleaseVersion) NameWithOrganization() string {
	return fmt.Sprintf("%s-%s", releaseVersion.Organization, releaseVersion.Name)
}

func (releaseVersion *ReleaseVersion) DownloadAndInstallAsset(provider ReleasesProvider, releaseDir string) error {
	var asset ReleaseAssets
	var err error

	if asset, err = releaseVersion.DownloadAsset(provider, releaseDir); err != nil {
		return err
	}

	err = releaseVersion.InstallAsset(asset, releaseDir)

	return err
}

func (releaseVersion *ReleaseVersion) InstallAsset(asset ReleaseAssets, releaseDir string) error {
	if false == releaseVersion.HasPackageMetadata(asset.DecompressPath()) {
		return errors.New(fmt.Sprintf("Error Package Metadata not found at %s", filepath.Join(asset.DecompressPath(), releaseVersion.Manifest())))
	}

	packageMetadata := releaseVersion.MustPackageMetadata(asset.DecompressPath())

	err := asset.Install(packageMetadata, releaseDir)
	if err != nil {
		return errors.New(fmt.Sprintf("Error Installing Package %s", err))
	}

	return nil
}

func (releaseVersion *ReleaseVersion) InstallAssetWithName(name string, asset ReleaseAssets, releaseDir string) error {
	return releaseVersion.InstallAsset(asset.CopyWithName(fmt.Sprintf("%s-%s", releaseVersion.Organization, name)), releaseDir)
}

func (releaseVersion *ReleaseVersion) SetVersion(version string) {
	releaseVersion.version = version
}

func (releaseVersion *ReleaseVersion) Version() string {
	return releaseVersion.version
}

func (releaseVersion *ReleaseVersion) VersionWithOutV() string {
	return releaseVersion.VersionWithV()[1:]
}

func (releaseVersion *ReleaseVersion) VersionWithV() string {
	if firstCharacter := releaseVersion.version[0:1]; firstCharacter == "v" {
		return releaseVersion.version
	}

	return "v" + releaseVersion.version
}

func (releaseVersion *ReleaseVersion) Equals(otherRelease ReleaseVersion) bool {
	return releaseVersion.version == otherRelease.version &&
		releaseVersion.Name == otherRelease.Name &&
		releaseVersion.Organization == otherRelease.Organization &&
		releaseVersion.manifest == otherRelease.manifest
}

func (releaseVersion ReleaseVersion) String() string {
	if "" != releaseVersion.version {
		return fmt.Sprintf("%s/%s:%s", releaseVersion.Organization, releaseVersion.Name, releaseVersion.version)
	}

	return fmt.Sprintf("%s/%s", releaseVersion.Organization, releaseVersion.Name)
}

func (releaseVersion ReleaseVersion) clone() ReleaseVersion {
	var clone = new(ReleaseVersion)
	clone.Organization = releaseVersion.Organization
	clone.Name = releaseVersion.Name
	clone.version = releaseVersion.version
	clone.manifest = releaseVersion.manifest

	return *clone
}

func (releaseVersion ReleaseVersion) CopyWithName(name string) ReleaseVersion {
	var clone = releaseVersion.clone()
	clone.Name = name

	return clone
}

func (releaseVersion ReleaseVersion) CopyWithManifest(manifest string) ReleaseVersion {
	var clone = releaseVersion.clone()
	clone.manifest = manifest

	return clone
}

func (releaseVersion ReleaseVersion) CopyWithVersion(version string) ReleaseVersion {
	var clone = releaseVersion.clone()
	clone.version = version

	return clone
}

func (releaseVersion *ReleaseVersion) Uninstall(releaseDir string) error {
	join := filepath.Join(releaseDir, releaseVersion.Name+"-*")
	installedPlugins, err := filepath.Glob(join)
	if err != nil {
		return errors.New(fmt.Sprintf("Error Uninstalling Release finding %s", join))
	}

	reg := regexp.MustCompile(`^` + releaseVersion.Name + `[\-]{1}[\d|\.]{4}\d$`)

	for _, file := range installedPlugins {
		if reg.MatchString(filepath.Base(file)) {
			err = os.Remove(file)
		}

		if err != nil {
			return errors.New(fmt.Sprintf("Error Uninstalling Release can't remove file %s", file))
		}
	}

	return nil
}

func NewReleaseLatestVersion(organization string, name string, finder ReleaseVersionFinder) (ReleaseVersion, error) {
	version, err := finder.Latest(organization, name)

	if err != nil {
		return ReleaseVersion{}, errors.New(fmt.Sprintf("Cant find latest release version of %s/%s", organization, name))
	}

	newReleaseVersion, err := NewReleaseVersionWith(
		ReleaseVersionWithOrganization(organization),
		ReleaseVersionWithName(name),
		ReleaseVersionWithVersion(version),
	)

	if err != nil {
		return ReleaseVersion{}, errors.New(fmt.Sprintf("Error creating ReleaseVersion of %s/%s:%s", organization, name, version))
	}

	return newReleaseVersion, nil
}

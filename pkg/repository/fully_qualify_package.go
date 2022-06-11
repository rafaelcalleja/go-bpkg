package repository

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrFullyQualifyPackageInvalidFormat = errors.New("fully qualify package invalid format")
)

type FullyQualifyPackage struct {
	organization string
	name         string
	version      string
}

func NewFullyQualifyPackage(fqp string) (FullyQualifyPackage, error) {
	expression := regexp.MustCompile(`^[\w\-\.]+\/([\w\-\.]+)(\:{1}[\w\.]+)?$`)
	if false == expression.MatchString(fqp) {
		return FullyQualifyPackage{}, ErrFullyQualifyPackageInvalidFormat
	}

	versionSeparatorIndex := strings.Index(fqp, ":")

	components := strings.Split(fqp, "/")
	version := ""
	if versionSeparatorIndex > -1 {
		components = strings.Split(fqp[0:versionSeparatorIndex], "/")
		version = fqp[versionSeparatorIndex+1:]
	}

	organization, name := components[0], components[1]

	return FullyQualifyPackage{
		organization: organization,
		name:         name,
		version:      version,
	}, nil
}

func (fqp FullyQualifyPackage) String() string {
	if "" != fqp.version {
		return fmt.Sprintf("%s/%s:%s", fqp.organization, fqp.name, fqp.version)
	}

	return fmt.Sprintf("%s/%s", fqp.organization, fqp.name)
}

func (fqp FullyQualifyPackage) Organization() string {
	return fqp.organization
}

func (fqp FullyQualifyPackage) Name() string {
	return fqp.name
}

func (fqp FullyQualifyPackage) Version() string {
	return fqp.version
}

func (fqp *FullyQualifyPackage) Equals(other FullyQualifyPackage) bool {
	return fqp.name == other.name &&
		fqp.organization == other.organization &&
		fqp.version == other.version
}

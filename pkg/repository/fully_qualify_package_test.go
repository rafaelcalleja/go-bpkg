package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewFullyQualifyPackage(t *testing.T) {
	t.Run("valid FQP Format", func(t *testing.T) {
		valid := map[string]FullyQualifyPackage{
			"organization/name":        {"organization", "name", ""},
			"organization/name:v1.0":   {"organization", "name", "v1.0"},
			"organization/name:latest": {"organization", "name", "latest"},
			"name/organization:2.0":    {"name", "organization", "2.0"},
		}

		for fqp, expected := range valid {
			fqpVO, err := NewFullyQualifyPackage(fqp)
			require.Nil(t, err)

			assert.True(t, fqpVO.Equals(expected))
		}
	})

	t.Run("Invalid FQP Format", func(t *testing.T) {
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

		for _, fqp := range invalid {
			_, err := NewFullyQualifyPackage(fqp)
			assert.Equal(t, err, ErrFullyQualifyPackageInvalidFormat)
		}
	})
}
